package ui

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type quiteView struct {
	tasks         config.TaskList
	currentState  tasks.RunnerStatus
	currentRunner *tasks.TaskRunner
}

// RunStatic runs a config.TaskList with a static ui suited for non TTY environments
func newQuiteRunner(taskList config.TaskList) {
	model := &quiteView{
		tasks: taskList,
	}

	model.setupInterruptHandler()
	execStart := time.Now()

	for _, step := range taskList.Steps {
		model.currentRunner = tasks.NewRunner(step.Name, step.Tasks, step.RunParallel)
		model.currentState = tasks.RunnerStatus{}
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