package concutils

import "sync"

// AndCombiner waits for all of its input channels to close
type AndCombiner struct {
	out <-chan any
}

// NewAndCombiner creates and starts an AndCombiner.
// The output channel closes only after all input channels have closed
func NewAndCombiner(channels ...<-chan any) *AndCombiner {
	switch len(channels) {
	case 0:
		return &AndCombiner{}
	case 1:
		return &AndCombiner{out: channels[0]}
	}

	out := make(chan any)
	var wg sync.WaitGroup
	wg.Add(len(channels))

	for _, channel := range channels {
		go func(ch <-chan any) {
			<-ch
			wg.Done()
		}(channel)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return &AndCombiner{out: out}
}

// Out returns the single output channel that closes when all input channels close
func (o *AndCombiner) Out() <-chan any {
	return o.out
}
