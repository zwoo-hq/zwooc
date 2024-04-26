package runner

// A TaskStatus represents the status of a task.
type TaskStatus int

const (
	// StatusPending indicates that the task is pending.
	StatusPending TaskStatus = iota
	// StatusScheduled indicates that the task is scheduled for execution.
	StatusScheduled
	// StatusRunning indicates that the task is currently running.
	StatusRunning
	// StatusDone indicates that the task has been successfully executed.
	StatusDone
	// StatusError indicates that the task has failed.
	StatusError
	// StatusCanceled indicates that the task has been canceled.
	StatusCanceled
)

// A RunnerStatus represents the status of a runner.
type RunnerStatus int

const (
	// RunnerIdle indicates that the runner is idle i.e. not running any tasks.
	RunnerIdle RunnerStatus = iota
	// RunnerRunning indicates that the runner is currently running pre tasks.
	RunnerPreparing
	// RunnerRunning indicates that the runner is currently running main tasks.
	RunnerRunning
	// RunnerRunning indicates that the runner is currently running post tasks.
	RunnerShuttingDown
	// RunnerCanceled indicates that the runner has been canceled.
	RunnerCanceled
	// RunnerErrored indicates that the runner has encountered an error.
	RunnerErrored
)

// A RunnerUpdate represents an update of the runner status.
type RunnerUpdate struct {
	// Status is the status of the runner.
	Status RunnerStatus
	// UpdatedNode is the status of the current node.
	UpdatedNode *TreeStatusNode
	// StatusRoot is the root node of the task tree.
	StatusRoot *TreeStatusNode
}

// A TreeStatusNode represents the status of a task tree node.
// The status of a task tree is represented in a tree structure od TreeStatusNodes that mirrors the task tree.
type TreeStatusNode struct {
	// ID is a unique identifier of the original node.
	ID string
	// name is the name of the node.
	name string
	// status is the status of the node.
	status TaskStatus
	// PreNodes is a collection of nodes that should be executed before the main task.
	PreNodes []*TreeStatusNode
	// PostNodes is a collection of nodes that should be executed after the main task.
	PostNodes []*TreeStatusNode
	// Parent is the parent node.
	Parent *TreeStatusNode
	// Error is the error that occurred during the execution of the main task.
	Error error
}

func (s *TreeStatusNode) Name() string {
	return s.name
}

func (s *TreeStatusNode) Status() TaskStatus {
	return s.status
}

func (t *TreeStatusNode) Iterate(handler func(node *TreeStatusNode)) {
	for _, pre := range t.PreNodes {
		pre.Iterate(handler)
	}
	handler(t)
	for _, post := range t.PostNodes {
		post.Iterate(handler)
	}
}
