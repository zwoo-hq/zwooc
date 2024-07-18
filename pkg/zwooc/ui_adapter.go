package zwooc

import (
	"github.com/zwoo-hq/zwooc/pkg/tasks/runner"
	"github.com/zwoo-hq/zwooc/pkg/ui"
)

func treeToStatus(node *runner.TreeStatusNode) *ui.TreeStatusNode {
	status := &ui.TreeStatusNode{
		Name:     node.Name,
		Status:   node.AggregatedStatus,
		Children: make([]*ui.TreeStatusNode, 0),
	}

	if len(node.PreNodes) > 0 {
		preWrapper := &ui.TreeStatusNode{
			Name:     "Pre",
			ID:       node.ID,
			Error:    node.Error,
			Status:   runner.StatusPending,
			Children: make([]*ui.TreeStatusNode, len(node.PreNodes)),
		}
		for i, child := range node.PreNodes {
			status.Children[i] = treeToStatus(child)
		}
		status.Children = append(status.Children, preWrapper)
	}

	mainNode := &ui.TreeStatusNode{
		Name:   node.MainName,
		Status: node.Status,
	}
	status.Children = append(status.Children, mainNode)

	if len(node.PostNodes) > 0 {
		postWrapper := &ui.TreeStatusNode{
			Name:     "Post",
			Status:   runner.StatusPending,
			Children: make([]*ui.TreeStatusNode, len(node.PostNodes)),
		}
		for i, child := range node.PostNodes {
			status.Children[i] = treeToStatus(child)
		}
		status.Children = append(status.Children, postWrapper)
	}

	return status
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
