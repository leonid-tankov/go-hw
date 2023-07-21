package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(tasks <-chan Task, errorCounter *int32, maxError int) {
	for task := range tasks {
		if int(atomic.LoadInt32(errorCounter)) >= maxError {
			return
		}
		if err := task(); err != nil {
			atomic.AddInt32(errorCounter, 1)
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	if m < 0 {
		m = 0
	}
	wg.Add(n)
	tasksChan := make(chan Task, len(tasks))
	for _, task := range tasks {
		tasksChan <- task
	}
	close(tasksChan)
	var errorCounter int32
	for i := 1; i <= n; i++ {
		go func() {
			worker(tasksChan, &errorCounter, m)
			wg.Done()
		}()
	}
	wg.Wait()
	if int(atomic.LoadInt32(&errorCounter)) >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
