package ui

type StatusUpdate struct {
	NodeID           string
	Status           TaskStatus
	AggregatedStatus TaskStatus
	Error            error
	Parent           *StatusUpdate
}

type StatusProvider interface {
	Start()
	OnStart(handler func())
	Cancel()
	OnCancel(handler func())
	UpdateStatus(update StatusUpdate)
	CloseUpdates()
	Done(err error)
}

type SimpleStatusProvider struct {
	start       func()
	cancel      func()
	status      chan StatusUpdate
	wasCanceled bool
	done        chan error
}

var _ StatusProvider = &SimpleStatusProvider{}

func NewSimpleStatusProvider() *SimpleStatusProvider {
	status := make(chan StatusUpdate)
	done := make(chan error)
	return &SimpleStatusProvider{
		status: status,
		done:   done,
	}
}

func (g SimpleStatusProvider) Start() {
	g.start()
}

func (g *SimpleStatusProvider) Cancel() {
	if !g.wasCanceled {
		g.wasCanceled = true
		g.cancel()
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

func (g *SimpleStatusProvider) OnStart(handler func()) {
	g.start = handler
}

func (g *SimpleStatusProvider) OnCancel(handler func()) {
	g.cancel = handler
}
