package zwooc

import (
	"github.com/zwoo-hq/zwooc/pkg/runner"
	"github.com/zwoo-hq/zwooc/pkg/ui"
)

func runnerToStatusProvider(updatedNode *runner.TreeStatusNode, provider ui.GenericStatusProvider) ui.StatusUpdate {
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
