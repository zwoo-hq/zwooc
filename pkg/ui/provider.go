package ui

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
