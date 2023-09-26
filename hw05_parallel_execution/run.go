package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded  = errors.New("errors limit exceeded")
	ErrErrorsCountOfWorkers = errors.New("errors limit count of workers")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n < 1 {
		return ErrErrorsCountOfWorkers
	}

	if len(tasks) < n {
		n = len(tasks)
	}

	countErrors := int32(0)
	ch := make(chan Task)
	wg := sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range ch {
				if err := t(); err != nil {
					atomic.AddInt32(&countErrors, 1)
				}
			}
		}()
	}

	for _, t := range tasks {
		if m > 0 && atomic.LoadInt32(&countErrors) >= int32(m) {
			break
		}
		ch <- t
	}

	close(ch)
	wg.Wait()

	if m > 0 && atomic.LoadInt32(&countErrors) >= int32(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
