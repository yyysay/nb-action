package worker

import "sync"

// Task 执行任务
type Task func() error

// Run 并发执行任务
func Run(tasks []Task, concurrency int) error {
	if concurrency <= 0 {
		concurrency = 1
	}

	taskChan := make(chan Task)

	var wg sync.WaitGroup
	var firstErr error
	var mu sync.Mutex

	// 创建 worker
	for i := 0; i < concurrency; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for task := range taskChan {
				if err := task(); err != nil {
					mu.Lock()

					if firstErr == nil {
						firstErr = err
					}

					mu.Unlock()
				}
			}
		}()
	}

	// 投递任务
	for _, task := range tasks {
		taskChan <- task
	}

	close(taskChan)

	wg.Wait()

	return firstErr
}
