package runner

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/zwoo-hq/zwooc/pkg/tasks"
	"golang.org/x/exp/maps"
)

type Scheduler struct {
	wasCanceled    atomic.Bool
	cancelForwards map[string]chan bool
	errs           []error
	status         SchedulerStatus
	updates        chan SchedulerStatus
	wg             sync.WaitGroup
	mu             sync.RWMutex
}

type SchedulerStatus = map[string]TaskStatus

func NewScheduler() *Scheduler {
	return &Scheduler{
		wasCanceled:    atomic.Bool{},
		cancelForwards: make(map[string]chan bool),
		status:         make(SchedulerStatus),
		errs:           make([]error, 0),
		updates:        make(chan SchedulerStatus, 1),
		wg:             sync.WaitGroup{},
		mu:             sync.RWMutex{},
	}
}

func (s *Scheduler) Status() SchedulerStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return maps.Clone(s.status)
}

func (s *Scheduler) Updates() <-chan SchedulerStatus {
	return s.updates
}

func (s *Scheduler) updateTaskStatus(task tasks.Task, status TaskStatus) {
	s.mu.Lock()
	s.status[task.Name()] = status
	// s.updates <- maps.Clone(s.status)
	s.mu.Unlock()
}

func (s *Scheduler) CancelTask(name string) {
	s.mu.Lock()
	if cancel, ok := s.cancelForwards[name]; ok {
		cancel <- true
	}
	s.mu.Unlock()
}

func (s *Scheduler) Schedule(task tasks.Task) {
	cancel := make(chan bool, 1)
	s.mu.Lock()
	s.wg.Add(1)
	s.mu.Unlock()
	s.updateTaskStatus(task, StatusPending)
	go func() {
		s.cancelForwards[task.Name()] = cancel
		defer s.wg.Done()
		s.updateTaskStatus(task, StatusRunning)
		err := task.Run(cancel)
		delete(s.cancelForwards, task.Name())
		if err != nil {
			s.updateTaskStatus(task, StatusError)
			s.mu.Lock()
			s.errs = append(s.errs, err)
			s.mu.Unlock()
		} else if !s.wasCanceled.Load() {
			s.updateTaskStatus(task, StatusCanceled)
		} else {
			s.updateTaskStatus(task, StatusDone)
		}
	}()
}

func (s *Scheduler) Cancel() error {
	s.wasCanceled.Store(true)
	s.mu.RLock()
	for _, cancel := range s.cancelForwards {
		cancel <- true
	}
	s.mu.RUnlock()
	s.wg.Wait()
	close(s.updates)
	if len(s.errs) == 0 {
		return nil
	}
	return errors.Join(s.errs...)
}
