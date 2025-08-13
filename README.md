# ConcUtils: A Go Concurrency Utilities Package
`concutils` is a small, dependency-free Go library providing generic, type-safe implementations of common concurrency patterns. It helps orchestrate complex goroutine workflows with clear, reusable components.

This library is designed to make concurrent code more readable and less error-prone by abstracting away common boilerplate for channel manipulation. Instead of writing complex `select` statements and `sync.WaitGroup` logic repeatedly, you can use the high-level components provided by this package.

## Features
- Or-Channel (OrCombiner): Monitors multiple "done" channels and signals completion as soon as the first one closes. Ideal for scenarios where you only need one successful result from a group of redundant tasks.

- And-Channel (AndCombiner): Monitors multiple "done" channels and signals completion only after all of them have closed. Perfect for waiting for a whole group of goroutines to finish their work.

- Fan-In / Merge (FanInMerger): Merges multiple input channels of the same type into a single output channel, simplifying consumption of results from multiple producers.

- Fan-Out / Distribute (FanOutDistributor): Distributes items from a single input channel across multiple output channels, typically for parallel processing by a pool of workers.

## Installation
``
go get github.com/l-ILINDAN-l/concutils
``

## Usage
Below are examples for each component.

### Or-Channel
Wait for the fastest goroutine to finish.
```
package main

import (
    "fmt"
    "github.com/l-ILINDAN-l/concutils"
    "time"
)

func main() {
    // sig is a helper function that returns a channel that closes after a duration
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
		sig(1*time.Second), // This one will finish first
		sig(3*time.Second),
	).Out()

	<-done
	fmt.Printf("The fastest goroutine has finished after %v.\n", time.Since(start))
}
```
### Fan-In
Merge results from multiple producers into one stream.
```
package main

import (
    "fmt"
    "github.com/l-ILINDAN-l/concutils"
    "time"
)

func main() {
    // producer creates a channel and sends `count` messages to it
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

	mergedChan := concutils.NewFanInMerger(
		producer("A", 3),
		producer("B", 2),
	).Out()

	for msg := range mergedChan {
		fmt.Println(msg)
	}
}
```
## Running Tests
To run the unit tests for this package, navigate to the package directory and run:
```
go test -v ./...
```
