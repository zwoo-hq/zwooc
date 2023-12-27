package tasks

import "io"

type emptyTask struct{}

var _ Task = (*emptyTask)(nil)

func (e *emptyTask) Name() string {
	return "#empty#"
}

func (e *emptyTask) Pipe(destination io.Writer) {
	// ignore
}

func (e *emptyTask) Run(cancel <-chan bool) error {
	return nil
}

func Empty() Task {
	return &emptyTask{}
}

func IsEmptyTask(task Task) bool {
	_, ok := task.(*emptyTask)
	return ok
}
