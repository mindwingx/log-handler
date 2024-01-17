package utils

import "sync"

type (
	Task func()

	Scheduler struct {
		wg      sync.WaitGroup
		tasks   []Task
		workers int
	}
)

func InitScheduler(workers int) *Scheduler {
	return &Scheduler{
		wg:      sync.WaitGroup{},
		workers: workers,
	}
}

func (s *Scheduler) AddTask(task Task) {
	s.tasks = append(s.tasks, task)
}

func (s *Scheduler) Start() {
	taskCh := make(chan Task, len(s.tasks))
	resultCh := make(chan struct{})

	// Start workers
	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
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
		s.wg.Wait()
		close(resultCh)
	}()

	// Wait for the result
	<-resultCh
}
