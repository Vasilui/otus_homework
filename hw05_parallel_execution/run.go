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

func Work(tasks <-chan Task, results chan<- error, stop <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		res := make(chan error)
		br := true
		go Process(task, res)

		for br {
			select {
			case r := <-res:
				results <- r
				br = false
			case <-stop:
				for {
					select {
					case r := <-res:
						results <- r
						return
					}
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
		go Work(jobs, results, out, &wg)
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
		if res, ok := <-results; ok {
			if res != nil {
				m--
				if m == 1 {
					for w := 1; w <= n; w++ {
						out <- struct{}{}
					}
					return ErrErrorsLimitExceeded
				}
			}
		}
	}

	return nil
}
