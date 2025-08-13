package concutils

// FanOutDistributor distributes items from one input channel to multiple output channels
// in a round-robin fashion. This is useful for distributing work among a pool of workers.
type FanOutDistributor[T any] struct {
	outs []<-chan T
}

// NewFanOutDistributor creates and starts a FanOutDistributor.
func NewFanOutDistributor[T any](in <-chan T, numOut int) *FanOutDistributor[T] {
	if numOut <= 0 {
		return &FanOutDistributor[T]{outs: nil}
	}

	outs := make([]chan T, numOut)
	for i := 0; i < numOut; i++ {
		outs[i] = make(chan T)
	}

	go func() {
		defer func() {
			for _, out := range outs {
				close(out)
			}
		}()

		for i := 0; ; i++ {
			val, ok := <-in
			if !ok {
				return
			}
			outs[i%numOut] <- val
		}
	}()

	readOnlyOuts := make([]<-chan T, numOut)
	for i := 0; i < numOut; i++ {
		readOnlyOuts[i] = outs[i]
	}
	return &FanOutDistributor[T]{outs: readOnlyOuts}
}

// Outs returns a slice of read-only output channels for the consumers.
func (f *FanOutDistributor[T]) Outs() []<-chan T {
	return f.outs
}
