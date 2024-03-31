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
	tasks         *tasks.TaskTreeNode
	currentRunner *runner.TaskTreeRunner
	wasCanceled   bool
	wg            sync.WaitGroup
	mu            sync.RWMutex
}

func NewStaticTreeRunner(tree *tasks.TaskTreeNode, opts ViewOptions) {
	model := &staticTreeView{
		tasks: tree,
	}

	model.setupInterruptHandler()
	fmt.Printf("%s - %s\n", zwoocBranding, tree.Name)
	execStart := time.Now()
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
	model.currentRunner = runner.NewTaskTreeRunner(tree, opts.MaxConcurrency)
	model.wg = sync.WaitGroup{}
	model.wg.Add(1)
	go model.ReceiveUpdates(model.currentRunner.Updates(), "│ ")

	start := time.Now()
	fmt.Printf("╭─── running %s\n", stepStyle.Render(tree.Name))
	err := model.currentRunner.Start()
	end := time.Now()
	// wait until everything is completed
	model.wg.Wait()
	execEnd := time.Now()

	if err != nil {
		// handle runner error
		fmt.Printf("╰─── %s failed\n", errorStyle.Render("✗"))
		model.currentRunner.Status().Iterate(func(node *runner.TreeStatusNode) {
			if node.Status() == runner.StatusError {
				fmt.Printf(" %s %s failed\n", errorStyle.Render("✗"), node.Name())
				fmt.Printf(" %s error: %s\n", errorStyle.Render("✗"), err)
				fmt.Printf(" %s stdout:\n", errorStyle.Render("✗"))
				// ligloss does some messy things to the string and cant handle \r\n on windows...
				wrapper := canceledStyle.Render("===")
				parts := strings.Split(wrapper, "===")
				fmt.Printf(parts[0])
				fmt.Println(strings.TrimSpace(outputs[node.ID].String()))
				fmt.Printf(parts[1])
			}
		})
		fmt.Printf("%s %s %s failed after %s\n", zwoocBranding, errorStyle.Render("✗"), tree.Name, execEnd.Sub(execStart))
		os.Exit(1)
	} else if model.wasCanceled {
		fmt.Printf("╰─── %s was canceled - stopping execution\n", canceledStyle.Render("-"))
		fmt.Printf("%s %s %s canceled after %s\n", zwoocBranding, canceledStyle.Render("-"), tree.Name, execEnd.Sub(execStart))
	} else {
		fmt.Printf("╰─── %s successfully ran %s\n", successStyle.Render("✓"), end.Sub(start))
		fmt.Printf("%s %s %s completed in %s\n", zwoocBranding, successStyle.Render("✓"), tree.Name, execEnd.Sub(execStart))
	}
}

func (m *staticTreeView) ReceiveUpdates(c <-chan *runner.TreeStatusNode, prefix string) {
	for node := range c {
		switch node.Status() {
		case runner.StatusPending:
			fmt.Printf("%s %s %s\n", prefix, node.Name(), pendingStyle.Render("was scheduled"))
		case runner.StatusRunning:
			fmt.Printf("%s %s %s\n", prefix, node.Name(), runningStyle.Render("started running"))
		case runner.StatusDone:
			fmt.Printf("%s %s %s\n", prefix, node.Name(), successStyle.Render("finished"))
		case runner.StatusError:
			fmt.Printf("%s %s %s\n", prefix, node.Name(), errorStyle.Render("failed"))
		case runner.StatusCanceled:
			fmt.Printf("%s %s %s\n", prefix, node.Name(), canceledStyle.Render("was canceled"))
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
