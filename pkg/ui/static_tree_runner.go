package ui

import (
	"fmt"
	"os"
	"os/signal"
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
	fmt.Printf("running %s (tree)\n", tree.Name)
	execStart := time.Now()

	// setup new runner
	model.currentRunner = runner.NewTaskTreeRunner(tree, opts.MaxConcurrency)
	model.wg = sync.WaitGroup{}
	model.wg.Add(1)
	go model.ReceiveUpdates(model.currentRunner.Updates(), "│ ")

	start := time.Now()
	fmt.Printf("╭─── running main tree\n")
	err := model.currentRunner.Start()
	end := time.Now()
	// wait until everything is completed
	model.wg.Wait()

	if err != nil {
		// handle runner error
		fmt.Printf("╰─── %s failed\n", errorStyle.Render("✗"))
	}
	if model.wasCanceled {
		fmt.Printf("╰─── %s was canceled - stopping execution\n", canceledStyle.Render("-"))
		return
	}
	fmt.Printf("╰─── %s successfully ran %s\n", successStyle.Render("✓"), end.Sub(start))

	execEnd := time.Now()
	fmt.Printf(" %s %s completed successfully in %s\n", successStyle.Render("✓"), tree.Name, execEnd.Sub(execStart))
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
