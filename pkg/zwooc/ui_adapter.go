package zwooc

import (
	"github.com/zwoo-hq/zwooc/pkg/tasks/runner"
	"github.com/zwoo-hq/zwooc/pkg/ui"
)

func treeToStatus(node *runner.TreeStatusNode) *ui.TreeStatusNode {
	status := &ui.TreeStatusNode{
		Name:     node.Name(),
		Status:   node.Status(),
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
		name:   node.ID,
		status: runner.StatusPending,
	}
	status.children = append(status.children, mainNode)

	if len(node.Post) > 0 {
		postWrapper := &ui.TreeStatusNode{
			name:     "Post",
			status:   runner.StatusPending,
			children: make([]*ui.TreeStatusNode, len(node.Post)),
		}
		for i, child := range node.Post {
			status.children[i] = treeToStatus(child)
		}
		status.children = append(status.children, postWrapper)
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
