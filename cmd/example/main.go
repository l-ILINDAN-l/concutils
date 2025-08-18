// This file provides runnable examples for the concutils library.
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	// Use the correct import path for your package.
	"github.com/ilindan-dev/concutils"
)

func main() {
	fmt.Println("--- OrCombiner Example ---")
	runOrExample()

	fmt.Println("\n--- AndCombiner Example ---")
	runAndExample()

	fmt.Println("\n--- FanInMerger Example ---")
	runFanInExample()

	fmt.Println("\n--- FanOutDistributor Example ---")
	runFanOutExample()

	fmt.Println("\n--- WorkerPool Example ---")
	runWorkerPoolExample()
}

// runOrExample demonstrates how to use OrCombiner.
func runOrExample() {
	sig := func(after time.Duration) <-chan any {
		c := make(chan any)
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	start := time.Now()
	done := concutils.NewOrCombiner(
		sig(2*time.Second),
		sig(1*time.Second), // This one will close first.
		sig(3*time.Second),
	).Out()
	<-done
	fmt.Printf("OrCombiner finished after %v\n", time.Since(start))
}

// runAndExample demonstrates how to use AndCombiner.
func runAndExample() {
	sig := func(after time.Duration) <-chan any {
		c := make(chan any)
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	start := time.Now()
	done := concutils.NewAndCombiner(
		sig(1*time.Second),
		sig(2*time.Second), // It will wait for this one.
	).Out()
	<-done
	fmt.Printf("AndCombiner finished after %v\n", time.Since(start))
}

// runFanInExample demonstrates how to use FanInMerger.
func runFanInExample() {
	producer := func(id string, count int) <-chan string {
		out := make(chan string)
		go func() {
			defer close(out)
			for i := 0; i < count; i++ {
				out <- fmt.Sprintf("Producer %s: message %d", id, i)
				time.Sleep(50 * time.Millisecond)
			}
		}()
		return out
	}
	mergedChan := concutils.NewFanInMerger(
		producer("A", 3),
		producer("B", 2),
	).Out()
	for msg := range mergedChan {
		fmt.Println(msg)
	}
}

// runFanOutExample demonstrates how to use FanOutDistributor.
func runFanOutExample() {
	in := make(chan int)
	numWorkers := 3
	go func() {
		defer close(in)
		for i := 0; i < 10; i++ {
			in <- i
		}
	}()
	distributor := concutils.NewFanOutDistributor(in, numWorkers)
	workerChannels := distributor.Outs()
	var wg sync.WaitGroup
	for i, ch := range workerChannels {
		wg.Add(1)
		go func(workerID int, workerChan <-chan int) {
			defer wg.Done()
			for num := range workerChan {
				fmt.Printf("Worker %d (FanOut) received: %d\n", workerID, num)
			}
		}(i, ch)
	}
	wg.Wait()
}

// --- WorkerPool Example ---

// myTask is a simple implementation of the Task interface for the example.
type myTask struct {
	id      int
	counter *int64
}

// Execute performs the work of the task.
func (t *myTask) Execute() {
	fmt.Printf("WorkerPool: Executing task %d\n", t.id)
	time.Sleep(100 * time.Millisecond)
	atomic.AddInt64(t.counter, 1)
}

// runWorkerPoolExample demonstrates how to use WorkerPool.
func runWorkerPoolExample() {
	pool := concutils.NewWorkerPool(4) // Create a pool with 4 workers
	var completedTasks int64

	fmt.Println("Submitting 20 tasks to the pool...")
	for i := 0; i < 20; i++ {
		pool.Submit(&myTask{id: i, counter: &completedTasks})
	}

	fmt.Println("All tasks submitted. Waiting for workers to finish...")
	pool.Stop() // Wait for all tasks to complete
	fmt.Printf("All %d tasks have been completed.\n", completedTasks)
}
