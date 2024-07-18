package ui

import (
	"github.com/zwoo-hq/zwooc/pkg/tasks"
	"github.com/zwoo-hq/zwooc/pkg/tasks/runner"
)

type TreeStatusProvider interface {
	Status() []*TreeStatusNode
}

type TreeProgressView struct {
	runner []*runner.TaskTreeRunner
	status []*TreeStatusNode
}

type TreeStatusNode struct {
	ID       string
	Name     string
	Error    error
	Status   runner.TaskStatus
	Children []*TreeStatusNode
}

func NewTreeProgressView(forest tasks.Collection, opts ViewOptions) *TreeProgressView {
	runners := make([]*runner.TaskTreeRunner, len(forest))
	statuses := make([]*TreeStatusNode, len(forest))
	concurrencyProvider := runner.NewSharedProvider(opts.MaxConcurrency)
	for i, node := range forest {
		// statuses[i] = treeToStatus(node)
		runners[i] = runner.NewTaskTreeRunner(node, concurrencyProvider)
	}
	return &TreeProgressView{
		runner: runners,
		status: statuses,
	}
}
