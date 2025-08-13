// This file provides runnable examples for the concutils library.
package main

import (
	"fmt"
	"github.com/l-ILINDAN-l/concutils"
	"sync"
	"time"
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
}

// runOrExample demonstrates how to use OrCombiner.
func runOrExample() {
	// sig is a helper function that returns a channel that closes after a duration.
	sig := func(after time.Duration) <-chan any {
		c := make(chan any)
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()

	// Create a combiner that waits for the fastest channel to close.
	done := concutils.NewOrCombiner(
		sig(2*time.Second),
		sig(1*time.Second), // This one will close first.
		sig(3*time.Second),
	).Out()

	<-done // Block until the combined channel is closed.
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

	// Create a combiner that waits for all channels to close.
	done := concutils.NewAndCombiner(
		sig(1*time.Second),
		sig(2*time.Second),
		sig(3*time.Second), // It will wait for this slowest one.
	).Out()

	<-done
	fmt.Printf("AndCombiner finished after %v\n", time.Since(start))
}

// runFanInExample demonstrates how to use FanInMerger.
func runFanInExample() {
	// producer creates a channel and sends `count` messages to it.
	producer := func(id string, count int) <-chan string {
		out := make(chan string)
		go func() {
			defer close(out)
			for i := 0; i < count; i++ {
				out <- fmt.Sprintf("Producer %s: message %d", id, i)
				time.Sleep(100 * time.Millisecond)
			}
		}()
		return out
	}

	// Merge the output of three producers into a single channel.
	mergedChan := concutils.NewFanInMerger(
		producer("A", 3),
		producer("B", 2),
		producer("C", 4),
	).Out()

	// Read all messages from the merged channel.
	for msg := range mergedChan {
		fmt.Println(msg)
	}
}

// runFanOutExample demonstrates how to use FanOutDistributor.
func runFanOutExample() {
	in := make(chan int)
	numWorkers := 3

	// Start a producer to send 10 numbers into the input channel.
	go func() {
		defer close(in)
		for i := 0; i < 10; i++ {
			in <- i
		}
	}()

	// Distribute the numbers to the worker channels.
	distributor := concutils.NewFanOutDistributor(in, numWorkers)
	workerChannels := distributor.Outs()

	// Start consumers to process the numbers from their respective channels.
	var wg sync.WaitGroup
	for i, ch := range workerChannels {
		wg.Add(1)
		go func(workerID int, workerChan <-chan int) {
			defer wg.Done()
			for num := range workerChan {
				fmt.Printf("Worker %d received: %d\n", workerID, num)
			}
		}(i, ch)
	}

	wg.Wait()
}
