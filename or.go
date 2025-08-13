package concutils

import (
	"sync"
)

// OrCombiner waits for any of its input channels to close
type OrCombiner struct {
	out <-chan any
}

// NewOrCombiner creates and starts an OrCombiner.
// The output channel closes as soon as any of the input channels close.
func NewOrCombiner(channels ...<-chan any) *OrCombiner {
	switch len(channels) {
	case 0:
		return &OrCombiner{}
	case 1:
		return &OrCombiner{out: channels[0]}
	}

	out := make(chan any)
	var once sync.Once

	for _, channel := range channels {
		go func(ch <-chan any) {
			select {
			case <-ch:
				once.Do(func() {
					close(out)
				})
			case <-out:
				return
			}
		}(channel)
	}

	return &OrCombiner{out}
}

// Out returns the single output channel that closes when any input channel closes.
func (o *OrCombiner) Out() <-chan any {
	return o.out
}
