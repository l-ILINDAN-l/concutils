package concutils

import (
	"sync"
	"sync/atomic"
)

// WorkerPool manages a fixed-size pool of goroutines to execute tasks concurrently.
type WorkerPool struct {
	taskQueue   chan Task
	wg          sync.WaitGroup
	activeTasks atomic.Int64
}

// NewWorkerPool creates and starts a new worker pool with a specified number of workers.
func NewWorkerPool(numWorkers int) *WorkerPool {
	pool := &WorkerPool{
		taskQueue: make(chan Task, numWorkers),
	}

	pool.wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer pool.wg.Done()
			for task := range pool.taskQueue {
				pool.activeTasks.Add(1)
				task.Execute()
				pool.activeTasks.Add(-1)
			}
		}()
	}

	return pool
}

// Submit adds a new task to the pool's queue for execution.
func (p *WorkerPool) Submit(task Task) {
	p.taskQueue <- task
}

// Stop gracefully shuts down the worker pool.
// It waits for all submitted tasks to be completed before stopping the workers.
func (p *WorkerPool) Stop() {
	close(p.taskQueue)
	p.wg.Wait()
}

// ActiveTasks returns the number of tasks currently being processed by workers.
func (p *WorkerPool) ActiveTasks() int64 {
	return p.activeTasks.Load()
}

// IsIdle returns true if no tasks are currently being processed.
func (p *WorkerPool) IsIdle() bool {
	return p.ActiveTasks() == 0
}
