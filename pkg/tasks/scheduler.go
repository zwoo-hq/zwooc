package tasks

import (
	"errors"
	"sync"
	"sync/atomic"

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

type SchedulerStatus = map[string]int

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

func (s *Scheduler) updateTaskStatus(task Task, status int) {
	s.mu.Lock()
	s.status[task.Name()] = status
	s.updates <- maps.Clone(s.status)
	s.mu.Unlock()
}

func (s *Scheduler) CancelTask(name string) {
	s.mu.Lock()
	if cancel, ok := s.cancelForwards[name]; ok {
		cancel <- true
	}
	s.mu.Unlock()
}

func (s *Scheduler) Schedule(tasks Task) {
	cancel := make(chan bool, 1)
	s.mu.Lock()
	s.cancelForwards[tasks.Name()] = cancel
	s.wg.Add(1)
	s.mu.Unlock()
	s.updateTaskStatus(tasks, StatusPending)
	go func() {
		defer s.wg.Done()
		s.updateTaskStatus(tasks, StatusRunning)
		err := tasks.Run(cancel)
		if err != nil {
			s.updateTaskStatus(tasks, StatusError)
			s.mu.Lock()
			s.errs = append(s.errs, err)
			s.mu.Unlock()
		} else if !s.wasCanceled.Load() {
			s.updateTaskStatus(tasks, StatusCanceled)
		} else {
			s.updateTaskStatus(tasks, StatusDone)
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
	if len(s.errs) == 0 {
		return nil
	}
	return errors.Join(s.errs...)
}
