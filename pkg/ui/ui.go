package ui

import "github.com/zwoo-hq/zwooc/pkg/tasks"

func NewRunner(forest tasks.Collection, options ViewOptions) {
	if options.QuiteMode {
		newQuiteRunner(forest, options)
		return
	}

	if options.DisableTUI {
		newStaticTreeRunner(forest, options)
		return
	}

	// try interactive view
	// if err := NewStatusView(task, options); err != nil {
	// 	// fall back to static view
	// 	newStaticRunner(task, options)
	// }
}

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

type StatusUpdate struct {
	NodeID string
	Status TaskStatus
}

type GenericStatusProvider interface {
	Status() StatusUpdate
	Cancel()
}
