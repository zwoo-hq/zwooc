package zwooc

import (
	"sync"

	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/runner"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
	"github.com/zwoo-hq/zwooc/pkg/ui"
	"golang.org/x/sync/errgroup"
)

type statusAdapter struct {
	options config.RunnerOptions

	scheduler *ui.SchedulerStatusProvider

	tasks               tasks.Collection
	runners             []*runner.TaskTreeRunner
	concurrencyProvider runner.ConcurrencyProvider

	isStarted bool
	updates   sync.WaitGroup
	errs      errgroup.Group
}

func newStatusAdapter(forest tasks.Collection, options config.RunnerOptions) *statusAdapter {
	concurrencyProvider := runner.NewSharedProvider(options.MaxConcurrency)
	scheduler := ui.NewSchedulerStatusProvider()
	adapter := &statusAdapter{
		options:             options,
		scheduler:           scheduler,
		concurrencyProvider: concurrencyProvider,
		tasks:               tasks.NewCollection(),
		runners:             []*runner.TaskTreeRunner{},
	}

	// map scheduler events to adapter
	scheduler.OnStart(adapter.start)
	scheduler.OnCancel(adapter.cancel)
	scheduler.OnShutdown(adapter.shutdownGracefully)
	scheduler.OnSchedule(adapter.schedule)

	// schedule initial tasks
	for _, node := range forest {
		adapter.addTask(node)
	}

	go func() {
		// collect ends of updates
		adapter.updates.Wait()
		scheduler.CloseUpdates()
	}()

	return adapter
}

func (a *statusAdapter) addTask(node *tasks.TaskTreeNode) {
	// create a new runner
	runner := runner.NewTreeRunner(node, a.concurrencyProvider)
	a.runners = append(a.runners, runner)
	a.tasks = append(a.tasks, node)

	if a.isStarted {
		// manually start the runner if the scheduler already started
		a.errs.Go(func() error {
			return runner.Start()
			// TODO: remove from runners
		})
	}

	// collect runner updates
	a.updates.Add(1)
	go func() {
		for update := range runner.Updates() {
			a.scheduler.UpdateStatus(runnerToStatusProvider(update))
		}
		a.updates.Done()
	}()

}

func (a *statusAdapter) start() {
	a.isStarted = true
	// start all known runners
	for _, r := range a.runners {
		currentRunner := r
		a.errs.Go(func() error {
			return currentRunner.Start()
		})
	}

	// collect done
	go func() {
		err := a.errs.Wait()
		a.scheduler.Done(err)
	}()
}

func (a *statusAdapter) cancel() {
	for _, r := range a.runners {
		r.Cancel()
	}
}

func (a *statusAdapter) shutdownGracefully() {
	for _, r := range a.runners {
		r.ShutdownGracefully()
	}
}

func (a *statusAdapter) schedule(command string, id string) {
	//TODO: implement scheduling
	// - check whether task is already scheduled (dont schedule twice)
}

func runnerToStatusProvider(updatedNode *runner.TreeStatusNode) ui.StatusUpdate {
	node := ui.StatusUpdate{
		NodeID:           updatedNode.ID,
		Status:           runnerStatusToUi(updatedNode.Status),
		AggregatedStatus: runnerStatusToUi(updatedNode.AggregatedStatus),
		Error:            updatedNode.Error,
	}

	if updatedNode.Parent != nil {
		parentUpdate := runnerToStatusProvider(updatedNode.Parent)
		node.Parent = &parentUpdate
	}
	return node
}

func runnerStatusToUi(status runner.TaskStatus) ui.TaskStatus {
	switch status {
	case runner.StatusPending:
		return ui.StatusPending
	case runner.StatusScheduled:
		return ui.StatusScheduled
	case runner.StatusRunning:
		return ui.StatusRunning
	case runner.StatusDone:
		return ui.StatusDone
	case runner.StatusError:
		return ui.StatusError
	case runner.StatusCanceled:
		return ui.StatusCanceled
	default:
		return ui.StatusPending
	}
}
