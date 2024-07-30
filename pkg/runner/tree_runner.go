package runner

import (
	"sync"
	"sync/atomic"

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
	cancelComplete chan bool
	// hasError is a flag that indicates whether an error occurred during the execution of the task tree.
	hasError atomic.Bool

	// mutex is used to synchronize access to the status tree.
	mutex sync.RWMutex
}

func NewTreeRunner(root *tasks.TaskTreeNode, p ConcurrencyProvider) *TaskTreeRunner {
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
		cancelComplete: make(chan bool),

		mutex: sync.RWMutex{},
	}
}

func (r *TaskTreeRunner) Updates() <-chan *TreeStatusNode {
	return r.updates
}

func (r *TaskTreeRunner) Status() *TreeStatusNode {
	return r.statusTree
}

func (r *TaskTreeRunner) Cancel() {
	r.cancel <- true
	close(r.cancel)
	<-r.cancelComplete
}

func (r *TaskTreeRunner) updateTaskStatus(node *tasks.TaskTreeNode, status TaskStatus) {
	r.mutex.Lock()
	statusNode := findStatus(r.statusTree, node)
	statusNode.Status = status
	statusNode.Update()
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
	errs := map[string]error{}
	errMu := sync.Mutex{}
	wg := sync.WaitGroup{}

	// cleanup
	defer func() {
		if r.wasCanceled.Load() {
			r.cancelComplete <- true
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
				ticket := r.tickets.Acquire()
				r.updateTaskStatus(task, StatusRunning)
				if err := task.Main.Run(cancel); err != nil {
					errMu.Lock()
					errs[task.NodeID()] = err
					if !r.hasError.Load() && !r.wasCanceled.Load() {
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
			close(r.scheduledNodes)
			for _, cancel := range r.forwardCancel {
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

	if r.wasCanceled.Load() {
		return tasks.ErrCancelled
	}

	if len(errs) > 0 {
		return tasks.NewMultiTaskError(errs)
	}
	return nil
}

func (r *TaskTreeRunner) scheduleNext(node *tasks.TaskTreeNode) {
	if r.hasError.Load() || r.wasCanceled.Load() {
		return // fail fast
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()
	statusNode := findStatus(r.statusTree, node)
	if allDone(r.statusTree) {
		close(r.scheduledNodes)
	} else if len(statusNode.PostNodes) > 0 {
		for _, post := range node.Post {
			for _, scheduled := range getStartingNodes(post) {
				r.scheduledNodes <- scheduled
			}
		}
	} else if statusNode.IsPre() && allChildrenWithStatus(statusNode.Parent.PreNodes, StatusDone) {
		r.scheduledNodes <- node.Parent
	} else if statusNode.IsPost() && allChildrenWithStatus(statusNode.Parent.PostNodes, StatusDone) {
		r.scheduledNodes <- node.Parent.Parent
	}
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
		Name:             root.Name,
		AggregatedStatus: StatusPending,
		MainName:         root.Main.Name(),
		Status:           StatusPending,
		PreNodes:         []*TreeStatusNode{},
		PostNodes:        []*TreeStatusNode{},
		ID:               root.NodeID(),
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
	if status.Status == StatusPending || status.Status == StatusRunning || status.Status == StatusScheduled {
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
