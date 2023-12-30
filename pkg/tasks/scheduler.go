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

func (s *Scheduler) Schedule(task Task) {
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
	if len(s.errs) == 0 {
		return nil
	}
	return errors.Join(s.errs...)
}
