package zwooc

import (
	"github.com/zwoo-hq/zwooc/pkg/runner"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
	"github.com/zwoo-hq/zwooc/pkg/ui"
	"golang.org/x/sync/errgroup"
)

func createSimpleForestRunner(forest tasks.Collection, maxConcurrency int) ui.SimpleStatusProvider {
	concurrencyProvider := runner.NewSharedProvider(maxConcurrency)
	runners := []*runner.TaskTreeRunner{}
	// create a new error group
	errs := errgroup.Group{}

	for _, tree := range forest {
		runners = append(runners, runner.NewTaskTreeRunner(tree, concurrencyProvider))
	}

	statusProvider := ui.NewSimpleStatusProvider()

	// forward start
	statusProvider.OnStart(func() {
		for _, r := range runners {
			currentRunner := r
			errs.Go(func() error {
				return currentRunner.Start()
			})
		}
	})

	// forward cancel
	statusProvider.OnCancel(func() {
		for _, r := range runners {
			r.Cancel()
		}
	})

	// forward updates
	for _, r := range runners {
		currentRunner := r
		go func() {
			for update := range currentRunner.Updates() {
				statusProvider.UpdateStatus(runnerToStatusProvider(update))
			}
		}()
	}

	// collect done
	go func() {
		err := errs.Wait()
		statusProvider.Done(err)
	}()

	return statusProvider
}

func runnerToStatusProvider(updatedNode *runner.TreeStatusNode) ui.StatusUpdate {
	return ui.StatusUpdate{
		NodeID: updatedNode.ID,
		Status: runnerStatusToUi(updatedNode.Status),
		Error:  updatedNode.Error,
	}
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