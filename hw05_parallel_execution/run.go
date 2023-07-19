package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(tasks <-chan Task, errors chan<- error) {
	for task := range tasks {
		if len(errors) == cap(errors) {
			return
		}
		if err := task(); err != nil {
			if len(errors) == cap(errors) {
				return
			}
			errors <- err
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	var result error
	if m < 0 {
		m = 0
	}
	wg.Add(n)
	tasksChan := make(chan Task, len(tasks))
	errorChan := make(chan error, m)
	for _, task := range tasks {
		tasksChan <- task
	}
	close(tasksChan)
	for i := 1; i <= n; i++ {
		go func() {
			worker(tasksChan, errorChan)
			wg.Done()
		}()
	}
	wg.Wait()
	close(errorChan)
	if len(errorChan) == cap(errorChan) {
		result = ErrErrorsLimitExceeded
	}
	return result
}
