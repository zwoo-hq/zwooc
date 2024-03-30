package runner

import "github.com/zwoo-hq/zwooc/pkg/tasks"

type TaskTreeRunner struct {
	root *tasks.TaskTreeNode
}

func NewTaskTreeRunner(root *tasks.TaskTreeNode) *TaskTreeRunner {
	return &TaskTreeRunner{root: root}
}
