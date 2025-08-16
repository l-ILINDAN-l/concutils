package concutils

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// testTask is a simple implementation of the Task interface for testing.
type testTask struct {
	counter *int64
}

// Execute implements the Task interface.
func (t *testTask) Execute() {
	time.Sleep(1 * time.Millisecond)
	atomic.AddInt64(t.counter, 1)
}

// TestWorkerPool checks the core functionality of the WorkerPool.
func TestWorkerPool(t *testing.T) {
	// Arrange
	numWorkers := 4
	numTasks := 100
	pool := NewWorkerPool(numWorkers)

	var completedTasks int64

	// Act
	// Submit tasks to the pool.
	for i := 0; i < numTasks; i++ {
		pool.Submit(&testTask{counter: &completedTasks})
	}

	// Check active tasks while running
	assert.LessOrEqual(t, pool.ActiveTasks(), int64(numWorkers))

	// Stop the pool and wait for all tasks to complete.
	pool.Stop()

	// Assert
	assert.Equal(t, int64(numTasks), completedTasks, "all submitted tasks should be completed")
	assert.True(t, pool.IsIdle(), "pool should be idle after stopping")
}
