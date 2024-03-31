package runner

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type TreeStatusNode struct {
	name      string
	status    TaskStatus
	PreNodes  []*TreeStatusNode
	PostNodes []*TreeStatusNode
	Parent    *TreeStatusNode
}

func (s *TreeStatusNode) Name() string {
	return s.name
}

func (s *TreeStatusNode) Status() TaskStatus {
	return s.status
}

type TaskTreeRunner struct {
	root           *tasks.TaskTreeNode
	scheduledNodes chan *tasks.TaskTreeNode
	status         *TreeStatusNode
	updates        chan *TreeStatusNode
	cancel         chan bool
	cancelComplete chan error
	mutex          sync.RWMutex
	maxConcurrency int
}

func NewTaskTreeRunner(root *tasks.TaskTreeNode, maxConcurrency int) *TaskTreeRunner {
	ticketAmount := maxConcurrency
	if ticketAmount < 1 {
		ticketAmount = runtime.NumCPU()
	}
	status := buildStatus(root)

	return &TaskTreeRunner{
		root:           root,
		scheduledNodes: make(chan *tasks.TaskTreeNode, 16),
		status:         status,
		updates:        make(chan *TreeStatusNode, 1),
		cancel:         make(chan bool),
		cancelComplete: make(chan error),
		maxConcurrency: ticketAmount,
		mutex:          sync.RWMutex{},
	}
}

func (r *TaskTreeRunner) Updates() <-chan *TreeStatusNode {
	return r.updates
}

func (r *TaskTreeRunner) Status() *TreeStatusNode {
	return r.status
}

func (r *TaskTreeRunner) Cancel() error {
	r.cancel <- true
	close(r.cancel)
	return <-r.cancelComplete
}

func (r *TaskTreeRunner) updateTaskStatus(node *tasks.TaskTreeNode, status TaskStatus) {
	r.mutex.Lock()
	statusNode := findStatus(r.status, node)
	statusNode.status = status
	r.updates <- statusNode
	r.mutex.Unlock()
}

func (r *TaskTreeRunner) Start() error {
	wasCanceled := atomic.Bool{}
	forwardCancel := []chan bool{}
	tickets := make(chan int, r.maxConcurrency)
	done := make(chan bool, 1)
	errs := []error{}
	errMu := sync.Mutex{}
	wg := sync.WaitGroup{}

	for i := 0; i < r.maxConcurrency; i++ {
		tickets <- i
	}

	// cleanup
	defer func() {
		if wasCanceled.Load() {
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
		close(tickets)
		for _, cancel := range forwardCancel {
			close(cancel)
		}
	}()

	// scheduler
	go func() {
		for scheduledNode := range r.scheduledNodes {
			if wasCanceled.Load() {
				break
			}

			wg.Add(1)
			taskCancel := make(chan bool, 1)
			forwardCancel = append(forwardCancel, taskCancel)

			go func(task *tasks.TaskTreeNode, cancel <-chan bool) {
				// acquire a ticket to run the task
				ticket := <-tickets
				r.updateTaskStatus(task, StatusRunning)
				if err := task.Main.Run(cancel); err != nil {
					errMu.Lock()
					errs = append(errs, err)
					errMu.Unlock()
					r.updateTaskStatus(task, StatusError)
					close(r.scheduledNodes)
				} else if wasCanceled.Load() {
					r.updateTaskStatus(task, StatusCanceled)
				} else {
					r.updateTaskStatus(task, StatusDone)
					// continue execution
					if len(errs) == 0 {
						r.scheduleNext(task)
					}
				}
				// release the ticket to be used by another channel
				tickets <- ticket
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
			wasCanceled.Store(true)
			close(r.scheduledNodes)
			for _, cancel := range forwardCancel {
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
	statusNode := findStatus(r.status, node)
	if isPre(statusNode) && helper.All(statusNode.Parent.PreNodes, func(n *TreeStatusNode) bool {
		return n.status == StatusDone
	}) {
		r.scheduledNodes <- node.Parent
	} else if isMain(statusNode) {
		for _, post := range node.Post {
			for _, scheduled := range getStartingNodes(post) {
				r.scheduledNodes <- scheduled
			}
		}
	} else if allDone(r.status) {
		close(r.scheduledNodes)
	}
}

func isPre(node *TreeStatusNode) bool {
	return node.Parent != nil && helper.IncludesBy(node.Parent.PreNodes, func(n *TreeStatusNode) bool {
		return n.name == node.name
	})
}

func isPost(node *TreeStatusNode) bool {
	return node.Parent != nil && helper.IncludesBy(node.Parent.PostNodes, func(n *TreeStatusNode) bool {
		return n.name == node.name
	})
}

func isMain(node *TreeStatusNode) bool {
	return !isPre(node) && !isPost(node)
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
		name:      root.Name,
		status:    StatusPending,
		PreNodes:  []*TreeStatusNode{},
		PostNodes: []*TreeStatusNode{},
	}

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
	if status.status != StatusPending && status.status != StatusRunning {
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
		if current.name == name {
			return current
		}

		for _, pre := range current.PreNodes {
			if pre.name == name {
				if i == len(path)-2 {
					return pre
				}
				current = pre
				continue outer
			}
		}

		for _, post := range current.PostNodes {
			if post.name == name {
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
