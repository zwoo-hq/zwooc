package tasks

import (
	"io"
)

type functionTask struct {
	name    string
	execute func(cancel <-chan bool, out io.Writer) error
	writer  *multiWriter
}

func NewTask(name string, execute func(cancel <-chan bool, out io.Writer) error) Task {
	return functionTask{
		name:    name,
		execute: execute,
		writer:  newMultiWriter(),
	}
}

func (ft functionTask) Name() string {
	return ft.name
}

func (ft functionTask) Pipe(destination io.Writer) {
	ft.writer.Pipe(destination)
}

func (ft functionTask) Run(cancel <-chan bool) error {
	return ft.execute(cancel, ft.writer)
}
