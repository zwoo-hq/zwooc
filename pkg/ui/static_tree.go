package ui

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type staticTreeView struct {
	forest      tasks.Collection
	provider    SimpleStatusProvider
	wasCanceled bool
	err         error
	wg          sync.WaitGroup
	mu          sync.RWMutex
}

func newStaticTreeRunner(forest tasks.Collection, provider SimpleStatusProvider, opts ViewOptions) {
	model := &staticTreeView{
		forest:   forest,
		provider: provider,
	}

	model.setupInterruptHandler()

	execStart := time.Now()
	outputs := map[string]*tasks.CommandCapturer{}

	// setup task pipes
	for _, tree := range forest {
		tree.Iterate(func(t *tasks.TaskTreeNode) {
			cap := tasks.NewCapturer()
			outputs[t.NodeID()] = cap
			t.Main.Pipe(cap)
			if opts.InlineOutput {
				if opts.DisablePrefix {
					t.Main.Pipe(tasks.NewPrefixer("  ", os.Stdout))
				} else {
					t.Main.Pipe(tasks.NewPrefixer("  "+t.Name+" ", os.Stdout))
				}
			}
		})
	}

	// start the runner
	fmt.Printf("%s running %s\n", zwoocBranding, forest.GetName())
	model.wg.Add(2)
	go model.ReceiveUpdates(provider.status, "")
	go model.WaitForDone()
	provider.Start()

	// wait until everything is completed
	model.wg.Wait()
	execEnd := time.Now()

	var failedError *tasks.MultiTaskError
	if errors.As(model.err, &failedError) {
		// handle runner error
		for nodeId, err := range failedError.Errors {
			fmt.Printf(" %s %s failed: %s\n", errorIcon, nodeId, err)
			fmt.Printf(" %s stdout:\n", errorIcon)
			wrapper := canceledStyle.Render("===")
			parts := strings.Split(wrapper, "===")
			fmt.Printf(parts[0])
			fmt.Println(strings.TrimSpace(outputs[nodeId].String()))
			fmt.Printf(parts[1])
		}
		fmt.Printf("%s %s %s failed after %s\n", zwoocBranding, errorIcon, forest.GetName(), execEnd.Sub(execStart))
	} else if model.wasCanceled || errors.Is(model.err, tasks.ErrCancelled) {
		fmt.Printf("%s %s %s canceled after %s\n", zwoocBranding, cancelIcon, forest.GetName(), execEnd.Sub(execStart))
	} else {
		fmt.Printf("%s %s %s completed after  %s\n", zwoocBranding, successIcon, forest.GetName(), execEnd.Sub(execStart))
	}

	// if hasError {
	// 	fmt.Printf("%s %s %s errored after %s\n", zwoocBranding, errorIcon, forest.GetName(), execEnd.Sub(execStart))
	// 	os.Exit(1)
	// } else if model.wasCanceled {
	// 	fmt.Printf("%s %s %s was canceled after %s\n", zwoocBranding, cancelIcon, forest.GetName(), execEnd.Sub(execStart))
	// } else {
	// 	fmt.Printf("%s %s %s completed in %s\n", zwoocBranding, successIcon, forest.GetName(), execEnd.Sub(execStart))
	// }
}

func (m *staticTreeView) ReceiveUpdates(c <-chan StatusUpdate, prefix string) {
	for node := range c {
		switch node.Status {
		case StatusPending:
			fmt.Printf("%s %s %s\n", prefix, node.NodeID, pendingStyle.Render("was scheduled"))
		case StatusRunning:
			fmt.Printf("%s %s %s\n", prefix, node.NodeID, runningStyle.Render("started running"))
		case StatusDone:
			fmt.Printf("%s %s %s\n", prefix, node.NodeID, successStyle.Render("finished"))
		case StatusError:
			fmt.Printf("%s %s %s\n", prefix, node.NodeID, errorStyle.Render("failed"))
		case StatusCanceled:
			fmt.Printf("%s %s %s\n", prefix, node.NodeID, canceledStyle.Render("was canceled"))
		}
	}
	m.wg.Done()
}

func (m *staticTreeView) WaitForDone() {
	err := <-m.provider.done
	if err != nil {
		m.mu.Lock()
		m.err = err
		m.mu.Unlock()
	}
	m.wg.Done()
}

func (m *staticTreeView) printFinalStatus() {

}

func (m *staticTreeView) setupInterruptHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		m.provider.Cancel()
	}()
}
