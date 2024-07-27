package runner

import "github.com/zwoo-hq/zwooc/pkg/helper"

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
	// Name is the name of the node.
	Name string
	// AggregatedStatus is the status of the node.
	AggregatedStatus TaskStatus
	// MainName is the name of the main task
	MainName string
	// Status is the status of the main task
	Status TaskStatus
	// PreNodes is a collection of nodes that should be executed before the main task.
	PreNodes []*TreeStatusNode
	// PostNodes is a collection of nodes that should be executed after the main task.
	PostNodes []*TreeStatusNode
	// Parent is the parent node.
	Parent *TreeStatusNode
	// Error is the error that occurred during the execution of the main task.
	Error error
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

func (t *TreeStatusNode) GetDirectChildren() []*TreeStatusNode {
	children := []*TreeStatusNode{}
	children = append(children, t.PreNodes...)
	children = append(children, t.PostNodes...)
	return children
}

func (t *TreeStatusNode) IsDone() bool {
	return t.AggregatedStatus == StatusDone || t.AggregatedStatus == StatusError || t.AggregatedStatus == StatusCanceled
}

func (t *TreeStatusNode) Update() {
	defer func() {
		if t.Parent != nil {
			t.Parent.Update()
		}
	}()

	children := t.GetDirectChildren()
	if len(children) == 0 {
		t.AggregatedStatus = t.Status
		return
	}

	// applies the status based on children by precedence -> the order of the if statements matters
	if someChildWithStatus(children, StatusError) {
		t.AggregatedStatus = StatusError
	} else if someChildWithStatus(children, StatusCanceled) {
		t.AggregatedStatus = StatusCanceled
	} else if allChildrenWithStatus(children, StatusDone) {
		t.AggregatedStatus = StatusRunning
	} else if someChildWithStatus(children, StatusRunning) {
		t.AggregatedStatus = StatusRunning
	} else if someChildWithStatus(children, StatusScheduled) {
		t.AggregatedStatus = StatusScheduled
	}
}

func (t *TreeStatusNode) IsPre() bool {
	return t.Parent != nil && helper.IncludesBy(t.Parent.PreNodes, func(n *TreeStatusNode) bool {
		return n.ID == t.ID
	})
}

func (t *TreeStatusNode) IsPost() bool {
	return t.Parent != nil && helper.IncludesBy(t.Parent.PostNodes, func(n *TreeStatusNode) bool {
		return n.ID == t.ID
	})
}

func someChildWithStatus(children []*TreeStatusNode, status TaskStatus) bool {
	return helper.Some(children, func(n *TreeStatusNode) bool {
		return n.AggregatedStatus == status
	})
}

func allChildrenWithStatus(children []*TreeStatusNode, status TaskStatus) bool {
	return helper.All(children, func(n *TreeStatusNode) bool {
		return n.AggregatedStatus == status
	})
}
