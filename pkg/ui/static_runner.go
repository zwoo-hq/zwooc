package ui

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type staticView struct {
	tasks         config.TaskList
	currentState  tasks.RunnerStatus
	currentRunner *tasks.TaskRunner
	wg            *sync.WaitGroup
}

// RunStatic runs a config.TaskList with a static ui suited for non TTY environments
func newStaticRunner(taskList config.TaskList) {
	model := &staticView{
		tasks: taskList,
	}

	model.setupInterruptHandler()
	fmt.Printf("running %s (%d steps)\n", taskList.Name, len(taskList.Steps))
	execStart := time.Now()

	for i, step := range taskList.Steps {
		start := time.Now()
		fmt.Printf("╭─── running step %s (%d/%d)\n", lipgloss.NewStyle().Foreground(lipgloss.Color("93")).Bold(true).Render(step.Name), i+1, len(taskList.Steps))

		for _, t := range step.Tasks {
			t.Pipe(os.Stdout)
		}

		model.currentRunner = tasks.NewRunner(step.Name, step.Tasks, step.RunParallel)
		model.currentState = tasks.RunnerStatus{}
		model.wg = &sync.WaitGroup{}
		model.wg.Add(1)
		go model.ReceiveUpdates(model.currentRunner.Updates(), "│ ")

		if err := model.currentRunner.Run(); err != nil {
			fmt.Printf("╰─── %s %s failed\n", errorStyle.Render("✗"), step.Name)
			HandleError(err)
		}

		end := time.Now()
		model.wg.Wait()
		fmt.Printf("╰─── %s %s successfully ran %s\n", successStyle.Render("✓"), step.Name, end.Sub(start))
	}

	execEnd := time.Now()
	fmt.Printf(" %s %s completed successfully in %s\n", successStyle.Render("✓"), taskList.Name, execEnd.Sub(execStart))
}

func (m *staticView) ReceiveUpdates(c <-chan tasks.RunnerStatus, prefix string) {
	for update := range c {
		for name, status := range update {
			if m.currentState[name] != status {
				m.currentState[name] = status
				switch status {
				case tasks.StatusPending:
					fmt.Printf("%s %s %s\n", prefix, name, pendingStyle.Render("was scheduled"))
				case tasks.StatusRunning:
					fmt.Printf("%s %s %s\n", prefix, name, runningStyle.Render("started running"))
				case tasks.StatusDone:
					fmt.Printf("%s %s %s\n", prefix, name, successStyle.Render("finished"))
				case tasks.StatusError:
					fmt.Printf("%s %s %s\n", prefix, name, errorStyle.Render("failed"))
				case tasks.StatusCanceled:
					fmt.Printf("%s %s %s\n", prefix, name, canceledStyle.Render("was canceled"))
				}
			}
		}
	}
	m.wg.Done()
}

func (m *staticView) setupInterruptHandler() {
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
