package concutils

import "sync"

// OrCombiner waits for any of its input "done" channels to close.
// It is useful for scenarios where you need to proceed as soon as the first
// of several concurrent operations completes.
type OrCombiner struct {
	out <-chan any
}

// NewOrCombiner creates and starts an OrCombiner.
// It takes one or more read-only "done" channels as input.
// The output channel closes as soon as any of the input channels close.
func NewOrCombiner(channels ...<-chan any) *OrCombiner {
	// Handle edge cases for 0 or 1 input channels for efficiency.
	switch len(channels) {
	case 0:
		// Return a combiner with a nil channel, which blocks forever.
		return &OrCombiner{}
	case 1:
		// If there's only one channel, no complex logic is needed.
		return &OrCombiner{out: channels[0]}
	}

	out := make(chan any)
	var once sync.Once

	// Start a listener goroutine for each input channel.
	for _, channel := range channels {
		go func(ch <-chan any) {
			select {
			case <-ch:
				// This input channel closed, try to trigger the main close.
				once.Do(func() {
					close(out)
				})
			case <-out:
				// The main channel already closed, so just exit.
				return
			}
		}(channel)
	}

	return &OrCombiner{out: out}
}

// Out returns the single output channel that closes when any input channel closes.
func (o *OrCombiner) Out() <-chan any {
	return o.out
}
