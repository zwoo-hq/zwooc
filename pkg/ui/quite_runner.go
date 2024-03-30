package ui

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/zwoo-hq/zwooc/pkg/tasks"
	"github.com/zwoo-hq/zwooc/pkg/tasks/runner"
)

type quiteView struct {
	tasks         tasks.TaskList
	currentState  runner.RunnerStatus
	currentRunner *runner.TaskRunner
}

// RunStatic runs a tasks.TaskList with a static ui suited for non TTY environments
func newQuiteRunner(taskList tasks.TaskList, opts ViewOptions) {
	model := &quiteView{
		tasks: taskList,
	}

	model.setupInterruptHandler()
	execStart := time.Now()

	for _, step := range taskList.Steps {
		model.currentRunner = runner.NewRunner(step.Name, step.Tasks, opts.MaxConcurrency)
		model.currentState = runner.RunnerStatus{}
		if err := model.currentRunner.Run(); err != nil {
			HandleError(err)
		}
	}

	execEnd := time.Now()
	fmt.Printf(" %s %s completed successfully in %s\n", successStyle.Render("✓"), taskList.Name, execEnd.Sub(execStart))
}

func (m *quiteView) setupInterruptHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			if m.currentRunner != nil {
				m.currentRunner.Cancel()
				break
			}
		}
	}()
}
