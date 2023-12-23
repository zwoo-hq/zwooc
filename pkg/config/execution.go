package config

import "github.com/zwoo-hq/zwooc/pkg/tasks"

type TaskList struct {
	Name  string
	Steps []ExecutionStep
}

type ExecutionStep struct {
	Name        string
	Tasks       []tasks.Task
	RunParallel bool
}
