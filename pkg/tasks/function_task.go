package tasks

import (
	"bytes"
	"io"
)

type functionTask struct {
	name       string
	execute    func(cancel <-chan bool, out io.Writer) error
	outWrapper *bytes.Buffer
}

func NewTask(name string, execute func(cancel <-chan bool, out io.Writer) error) Task {
	return functionTask{
		name:       name,
		execute:    execute,
		outWrapper: &bytes.Buffer{},
	}
}

func (ft functionTask) Name() string {
	return ft.name
}

func (ft functionTask) Out() bytes.Buffer {
	return *ft.outWrapper
}

func (ft functionTask) Run(cancel <-chan bool) error {
	return ft.execute(cancel, ft.outWrapper)
}
