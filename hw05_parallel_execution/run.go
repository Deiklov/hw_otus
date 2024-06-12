package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksChan := make(chan Task)

	resultChan := make(chan error)

	var errCounter int
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		defer close(tasksChan)
		for _, task := range tasks {
			select {
			// получили m ошибок, перестаем отправлять задания
			case <-ctx.Done():
				return
			case tasksChan <- task:
			}
		}
	}()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasksChan {
				// получили m ошибок, прерываем не ждем ответа, текущей задачи
				select {
				case <-ctx.Done():
					return
				case resultChan <- task():
				}
			}
		}()
	}

	go func() {
		for err := range resultChan {
			if err != nil {
				errCounter++
				if errCounter >= m {
					cancel()
					break
				}
			}
		}
	}()

	wg.Wait()
	close(resultChan)

	if errCounter >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
