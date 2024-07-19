package runner

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/zwoo-hq/zwooc/pkg/tasks"
	"golang.org/x/exp/maps"
)

type TaskRunnerStatus = map[string]TaskStatus

type TaskRunner struct {
	name           string
	tasks          []tasks.Task
	status         TaskRunnerStatus
	updates        chan TaskRunnerStatus
	cancel         chan bool
	cancelComplete chan error
	mutex          sync.RWMutex
	maxConcurrency int
}

func NewRunner(name string, tasks []tasks.Task, maxConcurrency int) *TaskRunner {
	status := make(TaskRunnerStatus)
	for _, task := range tasks {
		status[task.Name()] = StatusPending
	}

	ticketAmount := maxConcurrency
	if ticketAmount < 1 {
		ticketAmount = len(tasks)
	}

	return &TaskRunner{
		name:           name,
		tasks:          tasks,
		status:         status,
		updates:        make(chan TaskRunnerStatus, len(tasks)*5),
		cancel:         make(chan bool),
		cancelComplete: make(chan error),
		maxConcurrency: ticketAmount,
	}
}

func NewParallelRunner(name string, tasks []tasks.Task, maxConcurrency int) *TaskRunner {
	return NewRunner(name, tasks, maxConcurrency)
}

func NewSequentialRunner(name string, tasks []tasks.Task) *TaskRunner {
	return NewRunner(name, tasks, 1)
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
	if tr.maxConcurrency == 1 {
		return tr.runSequential()
	}
	return tr.runParallel()
}

func (tr *TaskRunner) Status() TaskRunnerStatus {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()
	return maps.Clone(tr.status)
}

func (tr *TaskRunner) Updates() <-chan TaskRunnerStatus {
	return tr.updates
}

func (tr *TaskRunner) updateTaskStatus(task tasks.Task, status TaskStatus) {
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
	tickets := make(chan int, tr.maxConcurrency)
	done := make(chan bool, 1)
	errs := []error{}
	errMu := sync.Mutex{}
	wg := sync.WaitGroup{}

	// allocate tickets for limiting concurrency
	for i := 0; i < tr.maxConcurrency; i++ {
		tickets <- i
	}

	defer func() {
		if wasCanceled.Load() {
			errMu.Lock()
			if len(errs) > 0 {
				tr.cancelComplete <- errors.Join(errs...)
			} else {
				tr.cancelComplete <- nil
			}
			errMu.Unlock()
		}
		close(tr.updates)
		close(tr.cancelComplete)
		close(tickets)
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

		go func(task tasks.Task, cancel <-chan bool) {
			// acquire a ticket to run the task
			ticket := <-tickets
			tr.updateTaskStatus(task, StatusRunning)
			if err := task.Run(cancel); err != nil {
				errMu.Lock()
				errs = append(errs, err)
				errMu.Unlock()
				tr.updateTaskStatus(task, StatusError)
			} else if wasCanceled.Load() {
				tr.updateTaskStatus(task, StatusCanceled)
			} else {
				tr.updateTaskStatus(task, StatusDone)
			}
			// release the ticket to be used by another channel
			tickets <- ticket
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
