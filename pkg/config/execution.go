package config

import (
	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

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
	for _, step := range t.Steps {
		// remove empty tasks
		for i := 0; i < len(step.Tasks); i++ {
			if tasks.IsEmptyTask(step.Tasks[i]) {
				step.Tasks = append(step.Tasks[:i], step.Tasks[i+1:]...)
				i--
			}
		}
		// remove empty steps
		if len(step.Tasks) > 0 {
			steps = append(steps, step)
		}
	}
	t.Steps = steps
}

func (t *TaskList) IsEmpty() bool {
	return len(t.Steps) == 0
}

func (t *TaskList) Split() (pre TaskList, main ExecutionStep, post TaskList) {
	pre = NewTaskList(helper.BuildName(t.Name, KeyPre), []ExecutionStep{})
	post = NewTaskList(helper.BuildName(t.Name, KeyPost), []ExecutionStep{})
	wasMain := false
	for _, step := range t.Steps {
		if step.IsLongRunning {
			main = step
			wasMain = true
		} else if wasMain {
			post.Steps = append(post.Steps, step)
		} else {
			pre.Steps = append(pre.Steps, step)
		}
	}
	return pre, main, post
}

type ExecutionStep struct {
	Name          string
	Tasks         []tasks.Task
	IsLongRunning bool
}

func NewExecutionStep(name string, tasks []tasks.Task, isLongRunning bool) ExecutionStep {
	return ExecutionStep{
		Name:          name,
		Tasks:         tasks,
		IsLongRunning: isLongRunning,
	}
}
