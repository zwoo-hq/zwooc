package ui

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/zwoo-hq/zwooc/pkg/tasks"
	"github.com/zwoo-hq/zwooc/pkg/tasks/runner"
)

type staticTreeView struct {
	forest        tasks.Collection
	currentRunner *runner.TaskTreeRunner
	wasCanceled   bool
	wg            sync.WaitGroup
	mu            sync.RWMutex
}

func newStaticTreeRunner(forest tasks.Collection, opts ViewOptions) {
	model := &staticTreeView{
		forest: forest,
	}

	fmt.Printf("%s - %s\n", zwoocBranding, forest.GetName())
	model.setupInterruptHandler()
	execStart := time.Now()
	provider := runner.NewSharedProvider(opts.MaxConcurrency)
	hasError := false

	for _, tree := range forest {
		start := time.Now()
		outputs := map[string]*tasks.CommandCapturer{}

		tree.Iterate(func(t *tasks.TaskTreeNode) {
			cap := tasks.NewCapturer()
			outputs[t.NodeID()] = cap
			t.Main.Pipe(cap)
			if opts.InlineOutput {
				if opts.DisablePrefix {
					t.Main.Pipe(tasks.NewPrefixer("│  ", os.Stdout))
				} else {
					t.Main.Pipe(tasks.NewPrefixer("│  "+t.Name+" ", os.Stdout))
				}
			}
		})

		// setup new runner
		model.currentRunner = runner.NewTaskTreeRunner(tree, provider)
		model.wg = sync.WaitGroup{}
		model.wg.Add(1)
		go model.ReceiveUpdates(model.currentRunner.Updates(), "│ ")
		fmt.Printf("╭─── running %s\n", stepStyle.Render(tree.Name))
		err := model.currentRunner.Start()
		if err != nil {
			hasError = true
		}

		// wait until everything is completed
		model.wg.Wait()
		end := time.Now()

		if err != nil {
			// handle runner error
			fmt.Printf("╰─── %s failed\n", errorIcon)
			model.currentRunner.Status().Iterate(func(node *runner.TreeStatusNode) {
				if node.AggregatedStatus == runner.StatusError {
					fmt.Printf(" %s %s failed\n", errorIcon, node.Name)
					fmt.Printf(" %s error: %s\n", errorIcon, err)
					fmt.Printf(" %s stdout:\n", errorIcon)
					// ligloss does some messy things to the string and cant handle \r\n on windows...
					wrapper := canceledStyle.Render("===")
					parts := strings.Split(wrapper, "===")
					fmt.Printf(parts[0])
					fmt.Println(strings.TrimSpace(outputs[node.ID].String()))
					fmt.Printf(parts[1])
				}
			})
			fmt.Printf("%s %s %s failed after %s\n", zwoocBranding, errorIcon, tree.Name, end.Sub(execStart))
		} else if model.wasCanceled {
			fmt.Printf("╰─── %s was canceled - stopping execution\n", cancelIcon)
			fmt.Printf("%s %s %s canceled after %s\n", zwoocBranding, cancelIcon, tree.Name, end.Sub(execStart))
		} else {
			fmt.Printf("╰─── %s successfully ran %s\n", successIcon, end.Sub(start))
		}
	}

	execEnd := time.Now()
	if hasError {
		fmt.Printf("%s %s %s errored after %s\n", zwoocBranding, errorIcon, forest.GetName(), execEnd.Sub(execStart))
		os.Exit(1)
	} else if model.wasCanceled {
		fmt.Printf("%s %s %s was canceled after %s\n", zwoocBranding, cancelIcon, forest.GetName(), execEnd.Sub(execStart))
	} else {
		fmt.Printf("%s %s %s completed in %s\n", zwoocBranding, successIcon, forest.GetName(), execEnd.Sub(execStart))
	}
}

func (m *staticTreeView) ReceiveUpdates(c <-chan *runner.TreeStatusNode, prefix string) {
	for node := range c {
		switch node.AggregatedStatus {
		case runner.StatusPending:
			fmt.Printf("%s %s %s\n", prefix, node.Name, pendingStyle.Render("was scheduled"))
		case runner.StatusRunning:
			fmt.Printf("%s %s %s\n", prefix, node.Name, runningStyle.Render("started running"))
		case runner.StatusDone:
			fmt.Printf("%s %s %s\n", prefix, node.Name, successStyle.Render("finished"))
		case runner.StatusError:
			fmt.Printf("%s %s %s\n", prefix, node.Name, errorStyle.Render("failed"))
		case runner.StatusCanceled:
			fmt.Printf("%s %s %s\n", prefix, node.Name, canceledStyle.Render("was canceled"))
		}
	}
	m.wg.Done()
}

func (m *staticTreeView) setupInterruptHandler() {
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
