package tasks

import (
	"bytes"
	"os/exec"
	"sync"

	"github.com/zwoo-hq/zwooc/pkg/helper"
)

type commandTask struct {
	name       string
	cmd        *exec.Cmd
	outWrapper *bytes.Buffer
}

func NewCommandTask(name string, cmd *exec.Cmd) Task {
	buf := &bytes.Buffer{}
	cmd.Stdout = buf
	cmd.Stderr = buf
	return commandTask{
		name:       name,
		cmd:        cmd,
		outWrapper: buf,
	}
}

func NewBasicCommandTask(name string, command string, dir string) Task {
	buf := &bytes.Buffer{}
	cmd := exec.Command("sh", "-c", command)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdout = buf
	cmd.Stderr = buf
	return commandTask{
		name:       name,
		cmd:        cmd,
		outWrapper: buf,
	}
}

func (ct commandTask) Name() string {
	return ct.name
}

func (ct commandTask) Out() bytes.Buffer {
	return *ct.outWrapper
}

func (ct commandTask) Run(cancel <-chan bool) error {
	// start the command
	if err := ct.cmd.Start(); err != nil {
		return err
	}

	errChan := make(chan error, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		// wait until the command is finished
		err := ct.cmd.Wait()
		if err != nil {
			errChan <- err
		}
		close(errChan)
		wg.Done()
	}()

	select {
	case <-cancel:
		// task go cancelled
		if err := ct.cmd.Process.Kill(); err != nil {
			return err
		}
	case <-helper.WaitFor(&wg):
		// task finished
		for err := range errChan {
			// there is at most one error
			return err
		}
		return nil
	}
	return nil
}
