/*
Package concutils provides generic, type-safe implementations of common
Go concurrency patterns. It helps orchestrate complex goroutine workflows
with clear, reusable components.

This library is designed to make concurrent code more readable and less
error-prone by abstracting away common boilerplate for channel manipulation.
Instead of writing complex select statements and sync.WaitGroup logic
repeatedly, you can use the high-level components provided by this package.

Features include:
  - Or-Channel (OrCombiner): Waits for the first of many "done" channels to close.
  - And-Channel (AndCombiner): Waits for all "done" channels to close.
  - Fan-In / Merge (FanInMerger): Merges multiple channels into a single channel.
  - Fan-Out / Distribute (FanOutDistributor): Distributes items from one channel
    to multiple worker channels.
*/
package concutils
