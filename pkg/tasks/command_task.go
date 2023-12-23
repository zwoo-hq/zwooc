package tasks

import (
	"os/exec"
	"sync"

	"github.com/zwoo-hq/zwooc/pkg/helper"
)

type commandTask struct {
	name string
	cmd  *exec.Cmd
}

func NewCommandTask(name string, cmd *exec.Cmd) Task {
	return commandTask{
		name: name,
		cmd:  cmd,
	}
}

func NewBasicCommandTask(name string, cmd string, dir string) Task {
	return commandTask{
		name: name,
		cmd:  helper.NewCommand(cmd, dir),
	}
}

func (ct commandTask) Name() string {
	return ct.name
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
			close(errChan)
		}
	}()

	select {
	case <-cancel:
		// task go cancelled
		if err := ct.cmd.Process.Kill(); err != nil {
			return err
		}
	case <-helper.WaitFor(&wg):
		// task finished
		close(errChan)
		for err := range errChan {
			// there is at most one error
			return err
		}
		return nil
	}
	return nil
}
