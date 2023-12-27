package config

import "github.com/zwoo-hq/zwooc/pkg/tasks"

type TaskList struct {
	Name  string
	Steps []ExecutionStep
}

func NewTaskList(name string, steps []ExecutionStep) TaskList {
	return TaskList{
		Name:  name,
		Steps: steps,
	}
}

func (t *TaskList) InsertBefore(other TaskList) {
	t.Steps = append(other.Steps, t.Steps...)
}

func (t *TaskList) InsertAfter(other TaskList) {
	t.Steps = append(t.Steps, other.Steps...)
}

func (t *TaskList) RemoveEmptyStagesAndTasks() {
	steps := []ExecutionStep{}
	// remove empty steps
	for _, step := range t.Steps {
		if len(step.Tasks) > 0 {
			// remove empty tasks
			for i := 0; i < len(step.Tasks); i++ {
				if tasks.IsEmptyTask(step.Tasks[i]) {
					step.Tasks = append(step.Tasks[:i], step.Tasks[i+1:]...)
					i--
				}
			}
			steps = append(steps, step)
		}
	}
	t.Steps = steps
}

type ExecutionStep struct {
	Name        string
	Tasks       []tasks.Task
	RunParallel bool
}

func NewExecutionStep(name string, tasks []tasks.Task, runParallel bool) ExecutionStep {
	return ExecutionStep{
		Name:        name,
		Tasks:       tasks,
		RunParallel: runParallel,
	}
}
