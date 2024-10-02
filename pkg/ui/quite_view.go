package ui

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type quiteTreeView struct {
	tasks    tasks.Collection
	provider *SimpleStatusProvider
}

func newQuiteTreeView(forest tasks.Collection, provider *SimpleStatusProvider) {
	model := &quiteTreeView{
		tasks:    forest,
		provider: provider,
	}
	model.setupInterruptHandler()

	execStart := time.Now()
	provider.Start()

	err := <-provider.done
	execEnd := time.Now()

	var failedError *tasks.MultiTaskError
	if errors.As(err, &failedError) {
		// handle runner error
		for nodeId, err := range failedError.Errors {
			fmt.Printf(" %s %s failed: %s\n", errorIcon, nodeId, err)
		}
		fmt.Printf("%s %s %s failed after %s\n", zwoocBranding, errorIcon, forest.GetName(), execEnd.Sub(execStart))
		os.Exit(1)
	} else if errors.Is(err, tasks.ErrCancelled) {
		fmt.Printf("%s %s %s canceled after %s\n", zwoocBranding, cancelIcon, forest.GetName(), execEnd.Sub(execStart))
	} else {
		fmt.Printf("%s %s %s completed after  %s\n", zwoocBranding, successIcon, forest.GetName(), execEnd.Sub(execStart))
	}
}

func (m *quiteTreeView) setupInterruptHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		m.provider.Cancel()
	}()
}
