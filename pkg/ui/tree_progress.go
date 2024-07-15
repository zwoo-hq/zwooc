package ui

import (
	"github.com/zwoo-hq/zwooc/pkg/tasks"
	"github.com/zwoo-hq/zwooc/pkg/tasks/runner"
)

type TreeProgressView struct {
	runner []*runner.TaskTreeRunner
	status []*treeStatusNode
}

type treeStatusNode struct {
	name     string
	status   runner.TaskStatus
	children []*treeStatusNode
}

func NewTreeProgressView(forest tasks.Collection, opts ViewOptions) *TreeProgressView {
	runners := make([]*runner.TaskTreeRunner, len(forest))
	statuses := make([]*treeStatusNode, len(forest))
	concurrencyProvider := runner.NewSharedProvider(opts.MaxConcurrency)
	for i, node := range forest {
		statuses[i] = treeToStatus(node)
		runners[i] = runner.NewTaskTreeRunner(node, concurrencyProvider)
	}
	return &TreeProgressView{
		runner: runners,
		status: statuses,
	}
}

func treeToStatus(node *tasks.TaskTreeNode) *treeStatusNode {
	status := &treeStatusNode{
		name:     node.Name,
		status:   runner.StatusPending,
		children: make([]*treeStatusNode, 0),
	}

	if len(node.Pre) > 0 {
		preWrapper := &treeStatusNode{
			name:     "Pre",
			status:   runner.StatusPending,
			children: make([]*treeStatusNode, len(node.Pre)),
		}
		for i, child := range node.Pre {
			status.children[i] = treeToStatus(child)
		}
		status.children = append(status.children, preWrapper)
	}

	mainNode := &treeStatusNode{
		name:   node.Main.Name(),
		status: runner.StatusPending,
	}
	status.children = append(status.children, mainNode)

	if len(node.Post) > 0 {
		postWrapper := &treeStatusNode{
			name:     "Post",
			status:   runner.StatusPending,
			children: make([]*treeStatusNode, len(node.Post)),
		}
		for i, child := range node.Post {
			status.children[i] = treeToStatus(child)
		}
		status.children = append(status.children, postWrapper)
	}

	return status
}
