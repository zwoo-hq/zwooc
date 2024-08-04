package ui

type Scheduler interface {
	Schedule(command string, id string)
	OnSchedule(handler func(command string, id string))
	Shutdown()
	OnShutdown(handler func())
	// TODO: implement stop/restart etc of individual tasks
	// TODO: needs notifier for:
	// - new task started (when scheduled)
}

type SchedulerStatusProvider struct {
	*SimpleStatusProvider
	schedule func(command string, id string)
	shutdown func()
}

var _ StatusProvider = &SchedulerStatusProvider{}
var _ Scheduler = &SchedulerStatusProvider{}

func NewSchedulerStatusProvider() *SchedulerStatusProvider {
	statusProvider := NewSimpleStatusProvider()
	return &SchedulerStatusProvider{
		SimpleStatusProvider: statusProvider,
	}
}

func (g SchedulerStatusProvider) Schedule(command string, id string) {
	g.schedule(command, id)
}

func (g *SchedulerStatusProvider) OnSchedule(handler func(command string, id string)) {
	g.schedule = handler
}

func (g SchedulerStatusProvider) Shutdown() {
	g.shutdown()
}

func (g *SchedulerStatusProvider) OnShutdown(handler func()) {
	g.shutdown = handler
}
