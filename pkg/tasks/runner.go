package tasks

import (
	"errors"
	"sync"
	"sync/atomic"

	"golang.org/x/exp/maps"
)

const (
	StatusPending  = 1
	StatusRunning  = 2
	StatusDone     = 3
	StatusError    = 4
	StatusCanceled = 5
)

type RunnerStatus = map[string]int

type TaskRunner struct {
	name           string
	tasks          []Task
	status         RunnerStatus
	updates        chan RunnerStatus
	cancel         chan bool
	cancelComplete chan error
	RunParallel    bool
	mutex          sync.RWMutex
}

func NewRunner(name string, tasks []Task, parallel bool) *TaskRunner {
	status := make(RunnerStatus)
	for _, task := range tasks {
		status[task.Name()] = StatusPending
	}

	return &TaskRunner{
		name:           name,
		tasks:          tasks,
		status:         status,
		RunParallel:    parallel,
		updates:        make(chan RunnerStatus, len(tasks)*5),
		cancel:         make(chan bool),
		cancelComplete: make(chan error),
	}
}

func NewParallelRunner(name string, tasks []Task) *TaskRunner {
	return NewRunner(name, tasks, true)
}

func NewSequentialRunner(name string, tasks []Task) *TaskRunner {
	return NewRunner(name, tasks, false)
}

func (tr *TaskRunner) Name() string {
	return tr.name
}

func (tr *TaskRunner) Cancel() error {
	tr.cancel <- true
	close(tr.cancel)
	return <-tr.cancelComplete
}

func (tr *TaskRunner) Run() error {
	if tr.RunParallel {
		return tr.runParallel()
	}
	return tr.runSequential()
}

func (tr *TaskRunner) Status() RunnerStatus {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()
	return maps.Clone(tr.status)
}

func (tr *TaskRunner) Updates() <-chan RunnerStatus {
	return tr.updates
}

func (tr *TaskRunner) updateTaskStatus(task Task, status int) {
	tr.mutex.Lock()
	tr.status[task.Name()] = status
	tr.updates <- maps.Clone(tr.status)
	tr.mutex.Unlock()
}

func (tr *TaskRunner) runSequential() error {
	wasCanceled := atomic.Bool{}
	forwardCancel := make(chan bool, 1)
	done := make(chan bool, 1)
	var err error

	defer func() {
		if wasCanceled.Load() {
			tr.cancelComplete <- err
		}
		close(tr.updates)
		close(tr.cancelComplete)
	}()

	go func() {
		select {
		case <-tr.cancel:
			// run was canceled - forward cancel to the task
			wasCanceled.Store(true)
			forwardCancel <- true
			close(forwardCancel)
			return
		case <-done:
			// stop the goroutine
			return
		}
	}()

	for _, task := range tr.tasks {
		tr.updateTaskStatus(task, StatusRunning)

		if err = task.Run(forwardCancel); err != nil {
			tr.updateTaskStatus(task, StatusError)
			return err
		}

		if wasCanceled.Load() {
			tr.updateTaskStatus(task, StatusCanceled)
			break
		} else {
			tr.updateTaskStatus(task, StatusDone)
		}
	}

	done <- true
	close(done)
	return nil
}

func (tr *TaskRunner) runParallel() error {
	wasCanceled := atomic.Bool{}
	forwardCancel := []chan bool{}
	done := make(chan bool, 1)
	errs := []error{}
	wg := sync.WaitGroup{}

	defer func() {
		if wasCanceled.Load() {
			if len(errs) > 0 {
				tr.cancelComplete <- errors.Join(errs...)
			} else {
				tr.cancelComplete <- nil
			}
		}
		close(tr.updates)
		close(tr.cancelComplete)
		for _, cancel := range forwardCancel {
			close(cancel)
		}
	}()

	go func() {
		select {
		case <-tr.cancel:
			// run was canceled - forward cancel to all tasks
			wasCanceled.Store(true)
			for _, cancel := range forwardCancel {
				cancel <- true
			}
			return
		case <-done:
			// stop the goroutine
			return
		}
	}()

	for _, task := range tr.tasks {
		wg.Add(1)
		taskCancel := make(chan bool, 1)
		forwardCancel = append(forwardCancel, taskCancel)

		go func(task Task, cancel <-chan bool) {
			tr.updateTaskStatus(task, StatusRunning)
			if err := task.Run(cancel); err != nil {
				errs = append(errs, err)
				tr.updateTaskStatus(task, StatusError)
			} else if wasCanceled.Load() {
				tr.updateTaskStatus(task, StatusCanceled)
			} else {
				tr.updateTaskStatus(task, StatusDone)

			}
			wg.Done()
		}(task, taskCancel)
	}

	wg.Wait()
	done <- true
	close(done)

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
