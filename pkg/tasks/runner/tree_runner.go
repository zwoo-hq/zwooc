package runner

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

// A TaskTreeRunner represents a runner for a task tree.
type TaskTreeRunner struct {
	// root is the root node of the task tree that should be executed.
	root *tasks.TaskTreeNode
	// status indicates the status of the runner.
	status RunnerStatus
	// statusTree is the mirrored statusTree tree of the task tree.
	statusTree *TreeStatusNode

	// scheduledNodes is a channel that is used to schedule nodes for execution.
	scheduledNodes chan *tasks.TaskTreeNode
	// forwardCancel is a collection of channels that are used to forward cancel signals to running tasks.
	forwardCancel map[string]chan bool
	// tickets is a concurrency provider that is used to limit the amount of concurrently running tasks.
	tickets ConcurrencyProvider

	// updates is a channel that is used to send updates of the status tree.
	updates chan *TreeStatusNode
	// wasCanceled is a flag that indicates whether the execution of the task tree was canceled.
	wasCanceled atomic.Bool
	// cancel is a channel that is used to cancel the execution of the task tree.
	cancel chan bool
	// cancelComplete is a channel that is used to signal that the cancel operation has completed.
	cancelComplete chan error
	// hasError is a flag that indicates whether an error occurred during the execution of the task tree.
	hasError atomic.Bool

	// mutex is used to synchronize access to the status tree.
	mutex sync.RWMutex
}

func NewTaskTreeRunner(root *tasks.TaskTreeNode, p ConcurrencyProvider) *TaskTreeRunner {
	status := buildStatus(root)

	return &TaskTreeRunner{
		root:       root,
		status:     RunnerIdle,
		statusTree: status,

		scheduledNodes: make(chan *tasks.TaskTreeNode, 16),
		forwardCancel:  map[string]chan bool{},
		tickets:        p,

		updates:        make(chan *TreeStatusNode, 1000),
		wasCanceled:    atomic.Bool{},
		cancel:         make(chan bool),
		cancelComplete: make(chan error),

		mutex: sync.RWMutex{},
	}
}

func (r *TaskTreeRunner) Updates() <-chan *TreeStatusNode {
	return r.updates
}

func (r *TaskTreeRunner) Status() *TreeStatusNode {
	return r.statusTree
}

// ShutdownGracefully cancels the execution of the task tree.
func (r *TaskTreeRunner) ShutdownGracefully() {
	// if r.status == RunnerIdle {
	// 	return
	// }

	// if r.status == RunnerPreparing {
	// 	r.wasCanceled.Store(true)
	// }

	// if cancel, ok := r.forwardCancel[r.root.NodeID()]; ok {
	// 	cancel <- true
	// }
}

func (r *TaskTreeRunner) Cancel() error {
	r.cancel <- true
	close(r.cancel)
	return <-r.cancelComplete
}

func (r *TaskTreeRunner) updateTaskStatus(node *tasks.TaskTreeNode, status TaskStatus) {
	r.mutex.Lock()
	statusNode := findStatus(r.statusTree, node)
	statusNode.Status = status
	r.updates <- statusNode
	r.mutex.Unlock()
}

func (r *TaskTreeRunner) setError(node *tasks.TaskTreeNode, err error) {
	r.mutex.Lock()
	statusNode := findStatus(r.statusTree, node)
	statusNode.Error = err
	r.hasError.Store(true)
	r.mutex.Unlock()
}

func (r *TaskTreeRunner) Start() error {
	done := make(chan bool, 1)
	errs := []error{}
	errMu := sync.Mutex{}
	wg := sync.WaitGroup{}

	// cleanup
	defer func() {
		if r.wasCanceled.Load() {
			errMu.Lock()
			if len(errs) > 0 {
				r.cancelComplete <- errors.Join(errs...)
			} else {
				r.cancelComplete <- nil
			}
			errMu.Unlock()
		}
		close(r.updates)
		close(r.cancelComplete)
		for _, cancel := range r.forwardCancel {
			close(cancel)
		}
	}()

	// scheduler
	go func() {
		for scheduledNode := range r.scheduledNodes {
			if r.wasCanceled.Load() {
				break
			}

			wg.Add(1)
			taskCancel := make(chan bool, 1)
			r.mutex.Lock()
			r.forwardCancel[scheduledNode.NodeID()] = taskCancel
			r.mutex.Unlock()
			go func(task *tasks.TaskTreeNode, cancel <-chan bool) {
				// acquire a ticket to run the task
				fmt.Println("wait", task.NodeID())
				ticket := r.tickets.Acquire()
				fmt.Println("start", task.NodeID())
				r.updateTaskStatus(task, StatusRunning)
				if err := task.Main.Run(cancel); err != nil {
					fmt.Println("error", task.NodeID(), err)
					errMu.Lock()
					errs = append(errs, err)
					if !r.hasError.Load() {
						// first node erroring -> close the channel
						close(r.scheduledNodes)
					}
					r.setError(task, err)
					r.updateTaskStatus(task, StatusError)
					errMu.Unlock()
				} else if r.wasCanceled.Load() {
					r.updateTaskStatus(task, StatusCanceled)
				} else {
					r.updateTaskStatus(task, StatusDone)
					// continue execution
					if len(errs) == 0 {
						r.scheduleNext(task)
					}
				}
				r.mutex.Lock()
				delete(r.forwardCancel, task.NodeID())
				r.mutex.Unlock()
				fmt.Println("release ticket", task.NodeID())
				// release the ticket to be used by another channel
				r.tickets.Release(ticket)
				wg.Done()
			}(scheduledNode, taskCancel)

		}
		wg.Done()
	}()

	// cancel
	go func() {
		select {
		case <-r.cancel:
			// run was canceled - forward cancel to all tasks
			r.wasCanceled.Store(true)
			for id, cancel := range r.forwardCancel {
				fmt.Println("cancel", id)
				cancel <- true
			}
			return
		case <-done:
			// stop the goroutine
			return
		}
	}()

	// start scheduling
	wg.Add(1)
	startingNodes := getStartingNodes(r.root)
	for _, node := range startingNodes {
		r.scheduledNodes <- node
	}

	wg.Wait()
	done <- true
	close(done)

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (r *TaskTreeRunner) scheduleNext(node *tasks.TaskTreeNode) {
	if r.hasError.Load() || r.wasCanceled.Load() {
		return // fail fast
	}

	statusNode := findStatus(r.statusTree, node)
	if isPre(statusNode) && helper.All(statusNode.Parent.PreNodes, func(n *TreeStatusNode) bool {
		return n.Status == StatusDone
	}) {
		r.scheduledNodes <- node.Parent
	} else if isMain(statusNode) {
		for _, post := range node.Post {
			for _, scheduled := range getStartingNodes(post) {
				r.scheduledNodes <- scheduled
			}
		}
	} else if allDone(r.statusTree) {
		close(r.scheduledNodes)
	}
}

func isMain(node *TreeStatusNode) bool {
	return node.Parent != nil && node.Parent.Main.Name == node.Name
}

func isPre(node *TreeStatusNode) bool {
	return node.Parent != nil && helper.IncludesBy(node.Parent.PreNodes, func(n *TreeStatusNode) bool {
		return n.Name == node.Name
	})
}

func isPost(node *TreeStatusNode) bool {
	return node.Parent != nil && helper.IncludesBy(node.Parent.PostNodes, func(n *TreeStatusNode) bool {
		return n.Name == node.Name
	})
}

func isWrapper(node *TreeStatusNode) bool {
	return !isPre(node) && !isPost(node) && !isMain(node)
}

func getStartingNodes(root *tasks.TaskTreeNode) []*tasks.TaskTreeNode {
	if len(root.Pre) == 0 {
		return []*tasks.TaskTreeNode{root}
	}

	allNodes := []*tasks.TaskTreeNode{}
	for _, pre := range root.Pre {
		allNodes = append(allNodes, getStartingNodes(pre)...)
	}

	return allNodes
}

func buildStatus(root *tasks.TaskTreeNode) *TreeStatusNode {
	status := &TreeStatusNode{
		Name:      root.Name,
		Status:    StatusPending,
		PreNodes:  []*TreeStatusNode{},
		PostNodes: []*TreeStatusNode{},
		ID:        root.NodeID(),
	}

	main := &TreeStatusNode{
		ID:     helper.BuildName(root.NodeID(), "main"),
		Name:   root.Main.Name(),
		Status: StatusPending,
	}
	status.Main = main
	main.Parent = status

	for _, pre := range root.Pre {
		preStatus := buildStatus(pre)
		preStatus.Parent = status
		status.PreNodes = append(status.PreNodes, preStatus)
	}

	for _, post := range root.Post {
		postStatus := buildStatus(post)
		postStatus.Parent = status
		status.PostNodes = append(status.PostNodes, postStatus)
	}

	return status
}

func allDone(status *TreeStatusNode) bool {
	if status.Status != StatusPending && status.Status != StatusRunning {
		return false
	}

	for _, pre := range status.PreNodes {
		if !allDone(pre) {
			return false
		}
	}

	for _, post := range status.PostNodes {
		if !allDone(post) {
			return false
		}
	}

	return true
}

func findStatus(status *TreeStatusNode, target *tasks.TaskTreeNode) *TreeStatusNode {
	path := []string{target.Name}
	for target.Parent != nil {
		path = append([]string{target.Parent.Name}, path...)
		target = target.Parent
	}
	if len(path) == 1 {
		return status
	}

	current := status
outer:
	for i, name := range path[1:] {
		if current.Name == name {
			return current
		}

		for _, pre := range current.PreNodes {
			if pre.Name == name {
				if i == len(path)-2 {
					return pre
				}
				current = pre
				continue outer
			}
		}

		for _, post := range current.PostNodes {
			if post.Name == name {
				if i == len(path)-2 {
					return post
				}
				current = post
				continue outer
			}
		}
		break
	}
	return nil
}
