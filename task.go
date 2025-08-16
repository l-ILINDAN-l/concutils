package concutils

// Task defines the interface for a unit of work that can be executed by a WorkerPool.
type Task interface {
	// Execute performs the work of the task.
	Execute()
}
