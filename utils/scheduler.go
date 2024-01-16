package utils

import "sync"

type (
	Task func()

	Scheduler struct {
		tasks   []Task
		workers int
	}
)

func InitScheduler(workers int) *Scheduler {
	return &Scheduler{
		workers: workers,
	}
}

func (s *Scheduler) AddTask(task Task) {
	s.tasks = append(s.tasks, task)
}

func (s *Scheduler) Start() {
	var wg sync.WaitGroup

	taskCh := make(chan Task, len(s.tasks))
	resultCh := make(chan struct{})

	// Start workers
	for i := 0; i < s.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskCh {
				task()
			}
		}()
	}

	// Add tasks to the channel
	go func() {
		for _, task := range s.tasks {
			taskCh <- task
		}
		close(taskCh)
	}()

	// Wait for all tasks to complete
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Wait for the result
	<-resultCh
}
