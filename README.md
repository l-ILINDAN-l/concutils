# ConcUtils: A Go Concurrency Utilities Package

[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go CI](https://github.com/l-ILINDAN-l/concutils/actions/workflows/ci.yml/badge.svg)](https://github.com/l-ILINDAN-l/concutils/actions/workflows/ci.yml)

`concutils` is a small, dependency-free Go library providing generic, type-safe implementations of common concurrency patterns. It helps orchestrate complex goroutine workflows with clear, reusable components.

This library is designed to make concurrent code more readable and less error-prone by abstracting away common boilerplate for channel manipulation. Instead of writing complex `select` statements and `sync.WaitGroup` logic repeatedly, you can use the high-level components provided by this package.

## Features

- **Or-Channel (`OrCombiner`):** Monitors multiple "done" channels and signals completion as soon as the *first* one closes.
- **And-Channel (`AndCombiner`):** Monitors multiple "done" channels and signals completion only after *all* of them have closed.
- **Fan-In / Merge (`FanInMerger`):** Merges multiple input channels of the same type into a single output channel.
- **Fan-Out / Distribute (`FanOutDistributor`):** Distributes items from one input channel across multiple output channels.
- **Worker Pool (`WorkerPool`):** Manages a fixed-size pool of goroutines to process tasks from a queue efficiently.

## Installation

```sh
go get github.com/ilindan-dev/concutils
```

## Usage
A runnable example demonstrating all components can be found in the `cmd/example` directory.

### Or-Channel
Wait for the fastest goroutine to finish.
``` Go
package main

import (
    "fmt"
    "github.com/ilindan-dev/concutils"
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

``` Go
package main

import (
    "fmt"
    "github.com/ilindan-dev/concutils"
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

### Worker Pool
Manage a pool of workers to process a large number of tasks concurrently.
``` Go
package main

import (
    "fmt"
    "github.com/ilindan-dev/concutils"
    "sync/atomic"
    "time"
)

// myTask is a simple implementation of the Task interface.
type myTask struct {
    id      int
    counter *int64
}

func (t *myTask) Execute() {
    fmt.Printf("Executing task %d\n", t.id)
    time.Sleep(100 * time.Millisecond)
    atomic.AddInt64(t.counter, 1)
}

func main() {
    pool := concutils.NewWorkerPool(4) // Create a pool with 4 workers
    var completedTasks int64

    for i := 0; i < 20; i++ {
        pool.Submit(&myTask{id: i, counter: &completedTasks})
    }

    pool.Stop() // Wait for all tasks to complete
    fmt.Printf("All %d tasks have been completed.\n", completedTasks)
}
```

## Running Tests
To run the unit tests for this package, navigate to the package directory and run:
``` sh
go test -v ./...
```
