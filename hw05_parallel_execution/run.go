package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Process(task Task, results chan<- error) {
	results <- task()
}

func Work(tasks <-chan Task, results chan<- error, stop <-chan struct{}) {
	for task := range tasks {
		res := make(chan error)
		go Process(task, res)
	Loop:
		for {
			select {
			case r := <-res:
				results <- r
				break Loop
			case <-stop:
				for {
					results <- <-res
					return
				}
			}
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	jobs := make(chan Task, len(tasks))
	results := make(chan error, len(tasks))
	out := make(chan struct{}, n)
	wg := sync.WaitGroup{}

	// run workers
	for w := 1; w <= n; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Work(jobs, results, out)
		}()
	}

	// send jobs
	for _, task := range tasks {
		jobs <- task
	}

	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	// Receive results
	for a := 1; a <= len(tasks); a++ {
		if res, ok := <-results; ok && res != nil {
			m--
			if m == 1 {
				for w := 1; w <= n; w++ {
					out <- struct{}{}
				}
				for {
					if _, ok := <-results; ok {
						continue
					}
					break
				}
				return ErrErrorsLimitExceeded
			}
		}
	}

	return nil
}
