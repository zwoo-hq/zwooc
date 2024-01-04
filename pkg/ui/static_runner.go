package ui

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
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
	wasCanceled   bool
	wg            sync.WaitGroup
	mu            sync.RWMutex
}

// RunStatic runs a config.TaskList with a static ui suited for non TTY environments
func newStaticRunner(taskList config.TaskList, opts ViewOptions) {
	model := &staticView{
		tasks: taskList,
	}

	model.setupInterruptHandler()
	fmt.Printf("running %s (%d steps)\n", taskList.Name, len(taskList.Steps))
	execStart := time.Now()

	for i, step := range taskList.Steps {
		// capture the output of each step
		outputs := map[string]*tasks.CommandCapturer{}
		for _, t := range step.Tasks {
			cap := tasks.NewCapturer()
			outputs[t.Name()] = cap
			t.Pipe(cap)
			if opts.InlineOutput {
				if opts.DisablePrefix {
					t.Pipe(tasks.NewPrefixer("│  ", os.Stdout))
				} else {
					t.Pipe(tasks.NewPrefixer("│  "+t.Name()+" ", os.Stdout))
				}
			}
		}

		// setup new runner
		model.currentRunner = tasks.NewRunner(step.Name, step.Tasks, opts.MaxConcurrency)
		model.currentState = tasks.RunnerStatus{}
		model.wg = sync.WaitGroup{}
		model.wg.Add(1)
		go model.ReceiveUpdates(model.currentRunner.Updates(), "│ ")

		start := time.Now()
		fmt.Printf("╭─── running step %s (%d/%d)\n", lipgloss.NewStyle().Foreground(lipgloss.Color("93")).Bold(true).Render(step.Name), i+1, len(taskList.Steps))
		err := model.currentRunner.Run()
		end := time.Now()
		// wait until everything is completed
		model.wg.Wait()

		if err != nil {
			// handle runner error
			fmt.Printf("╰─── %s %s failed\n", errorStyle.Render("✗"), step.Name)
			for key, status := range model.currentRunner.Status() {
				if status == tasks.StatusError {
					fmt.Printf(" %s %s failed\n", errorStyle.Render("✗"), key)
					fmt.Printf(" %s error: %s\n", errorStyle.Render("✗"), err)
					fmt.Printf(" %s stdout:\n", errorStyle.Render("✗"))
					// ligloss does some messy things to the string and cant handle \r\n on windows...
					wrapper := canceledStyle.Render("===")
					parts := strings.Split(wrapper, "===")
					fmt.Printf(parts[0])
					fmt.Println(strings.TrimSpace(outputs[key].String()))
					fmt.Printf(parts[1])
				}
			}
			return
		}

		if model.wasCanceled {
			fmt.Printf("╰─── %s %s was canceled - stopping execution\n", canceledStyle.Render("-"), step.Name)
			return
		}
		fmt.Printf("╰─── %s %s successfully ran %s\n", successStyle.Render("✓"), step.Name, end.Sub(start))

	}

	execEnd := time.Now()
	fmt.Printf(" %s %s completed successfully in %s\n", successStyle.Render("✓"), taskList.Name, execEnd.Sub(execStart))
}

func (m *staticView) ReceiveUpdates(c <-chan tasks.RunnerStatus, prefix string) {
	for update := range c {
		m.mu.Lock()
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
		m.mu.Unlock()
	}
	m.wg.Done()
}

func (m *staticView) setupInterruptHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			m.mu.Lock()
			if m.currentRunner != nil {
				m.currentRunner.Cancel()
				m.wasCanceled = true
				m.mu.Unlock()
				break
			}
			m.mu.Unlock()
		}
	}()
}
