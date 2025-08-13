package concutils

import "sync"

// FanInMerger merges multiple input channels into a single output channel
type FanInMerger[T any] struct {
	out <-chan T
}

// NewFanInMerger creates and starts a FanInMerger
func NewFanInMerger[T any](channels ...<-chan T) *FanInMerger[T] {
	out := make(chan T)
	var wg sync.WaitGroup

	for _, channel := range channels {
		wg.Add(1)
		go func(ch <-chan T) {
			defer wg.Done()
			for val := range ch {
				out <- val
			}
		}(channel)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return &FanInMerger[T]{out: out}
}

// Out returns the single merged output channel
func (f *FanInMerger[T]) Out() <-chan T {
	return f.out
}
