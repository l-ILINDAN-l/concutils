package concutils

import "sync"

// AndCombiner waits for all of its input "done" channels to close.
// It is useful for waiting for a group of goroutines to complete their work.
type AndCombiner struct {
	out <-chan any
}

// NewAndCombiner creates and starts an AndCombiner.
// The output channel closes only after all input channels have closed.
func NewAndCombiner(channels ...<-chan any) *AndCombiner {
	// Handle edge cases
	switch len(channels) {
	case 0:
		return &AndCombiner{}
	case 1:
		return &AndCombiner{out: channels[0]}
	}

	out := make(chan any)
	var wg sync.WaitGroup
	wg.Add(len(channels))

	// Start a listener goroutine for each input channel.
	for _, channel := range channels {
		go func(ch <-chan any) {
			<-ch // Wait for the channel to close
			wg.Done()
		}(channel)
	}

	// Start a final goroutine that waits for the counter to reach zero.
	go func() {
		wg.Wait()
		close(out)
	}()

	return &AndCombiner{out: out}
}

// Out returns the single output channel that closes when all input channels close.
func (a *AndCombiner) Out() <-chan any {
	return a.out
}
