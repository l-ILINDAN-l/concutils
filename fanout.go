package concutils

// FanOutDistributor distributes items from one input channel to multiple output channels
type FanOutDistributor[T any] struct {
	outs []<-chan T
}

// NewFanOutDistributor creates and starts a FanOutDistributor
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
			for _, c := range outs {
				close(c)
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

// Outs returns a slice of read-only output channels
func (f *FanOutDistributor[T]) Outs() []<-chan T {
	return f.outs
}
