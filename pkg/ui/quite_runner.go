package ui

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/zwoo-hq/zwooc/pkg/runner"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type quiteView struct {
	tasks   tasks.Collection
	runners []*runner.TaskTreeRunner
	errs    []error
	mu      sync.Mutex
}

// RunStatic runs a tasks.TaskList with a static ui suited for non TTY environments
func newQuiteRunner(forest tasks.Collection, opts ViewOptions) {
	model := &quiteView{
		tasks: forest,
		errs:  []error{},
		mu:    sync.Mutex{},
	}

	model.setupInterruptHandler()
	execStart := time.Now()

	concurrencyProvider := runner.NewSharedProvider(opts.MaxConcurrency)
	wg := sync.WaitGroup{}

	for _, tree := range forest {
		runner := runner.NewTaskTreeRunner(tree, concurrencyProvider)
		model.runners = append(model.runners, runner)
		wg.Add(1)
		go func() {
			if err := runner.Start(); err != nil {
				model.mu.Lock()
				model.errs = append(model.errs, err)
				model.mu.Unlock()
			}
			wg.Done()
		}()
	}

	wg.Wait()

	for _, err := range model.errs {
		HandleError(err)
	}

	execEnd := time.Now()
	fmt.Printf(" %s %s completed successfully in %s\n", successIcon, forest.GetName(), execEnd.Sub(execStart))
}

func (m *quiteView) setupInterruptHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			for _, runner := range m.runners {
				runner.Cancel()
			}
			break
		}
	}()
}
