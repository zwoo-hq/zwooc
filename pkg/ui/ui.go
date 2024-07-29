package ui

import "github.com/zwoo-hq/zwooc/pkg/tasks"

func NewRunner(forest tasks.Collection, provider SimpleStatusProvider, options ViewOptions) {
	if options.QuiteMode {
		// TODO: use provided runner
		newQuiteRunner(forest, options)
		return
	}

	if options.DisableTUI {
		newStaticTreeRunner(forest, provider, options)
		return
	}

	// try interactive view
	if err := NewTreeProgressView(forest, provider, options); err != nil {
		// fall back to static view
		newStaticTreeRunner(forest, provider, options)
	}
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
	NodeID           string
	Status           TaskStatus
	AggregatedStatus TaskStatus
	Error            error
	Parent           *StatusUpdate
}

type SimpleStatusProvider struct {
	status      chan StatusUpdate
	cancel      chan struct{}
	wasCanceled bool
	start       chan struct{}
	done        chan error
}

func (g SimpleStatusProvider) Start() {
	g.start <- struct{}{}
	close(g.start)
}

func (g *SimpleStatusProvider) Cancel() {
	if !g.wasCanceled {
		g.wasCanceled = true
		g.cancel <- struct{}{}
		close(g.cancel)
	}
}

func (g SimpleStatusProvider) UpdateStatus(update StatusUpdate) {
	g.status <- update
}

func (g SimpleStatusProvider) CloseUpdates() {
	close(g.status)
}

func (g SimpleStatusProvider) Done(err error) {
	g.done <- err
	close(g.done)
}

func (g SimpleStatusProvider) OnStart(handler func()) {
	go func() {
		<-g.start
		handler()
	}()
}

func (g SimpleStatusProvider) OnCancel(handler func()) {
	go func() {
		<-g.cancel
		handler()
	}()
}

func NewSimpleStatusProvider() SimpleStatusProvider {
	status := make(chan StatusUpdate)
	cancel := make(chan struct{})
	done := make(chan error)
	start := make(chan struct{})
	return SimpleStatusProvider{
		status: status,
		cancel: cancel,
		done:   done,
		start:  start,
	}
}
