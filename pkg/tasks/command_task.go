package tasks

import (
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/zwoo-hq/zwooc/pkg/helper"
)

type commandTask struct {
	name   string
	cmd    *exec.Cmd
	writer *multiWriter
	stdIn  io.WriteCloser
}

func NewCommandTask(name string, cmd *exec.Cmd) Task {
	writer := newMultiWriter()
	cmd.Stdout = writer
	cmd.Stderr = writer
	in, _ := cmd.StdinPipe()
	return commandTask{
		name:   name,
		cmd:    cmd,
		writer: writer,
		stdIn:  in,
	}
}

func NewBasicCommandTask(name string, command string, dir string, args []string) Task {
	writer := newMultiWriter()
	fullCommand := strings.Join(append([]string{command}, args...), " ")
	cmd := exec.Command("sh", "-c", fullCommand)
	if dir != "" {
		cmd.Dir = dir
	}
	in, _ := cmd.StdinPipe()
	cmd.Stdout = writer
	cmd.Stderr = writer
	return commandTask{
		name:   name,
		cmd:    cmd,
		writer: writer,
		stdIn:  in,
	}
}

func (ct commandTask) Name() string {
	return ct.name
}

func (ct commandTask) Pipe(destination io.Writer) {
	ct.writer.Pipe(destination)
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
		ct.writer.Write([]byte(fmt.Sprintf("pid: %d\n", ct.cmd.Process.Pid)))
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
		var err error

		if runtime.GOOS == "windows" {
			err = exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(ct.cmd.Process.Pid)).Run()
		} else {
			err = exec.Command("pkill", "-P", strconv.Itoa(ct.cmd.Process.Pid)).Run()
		}

		if err != nil {
			// fall back to builtin kill
			if err := ct.cmd.Process.Kill(); err != nil {
				return err
			}
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
