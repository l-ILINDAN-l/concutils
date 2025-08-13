package concutils

import (
	"testing"
	"time"
)

// TestAndCombiner checks th logic of the AndCombiner
func TestAndCombiner(t *testing.T) {
	// sig is a helper function that returns a channel that closes after a specified duration
	sig := func(after time.Duration) <-chan any {
		c := make(chan any)
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	// Arrange: Record the start time
	start := time.Now()

	// Act: Create a new combiner with channels that have different closing times
	// We expect the combiner to close when the slowest channel (10s) closes
	combiner := NewAndCombiner(
		sig(10*time.Second), // This channel
		sig(100*time.Millisecond),
		sig(5*time.Second),
	)

	// Wait closing combined channel
	<-combiner.Out()

	// Assert: Check if the elapsed time is reasonable
	elapsed := time.Since(start)

	if elapsed < 10*time.Second {
		t.Errorf("AndCombiner closed too early. Got %v, expected >10s", elapsed)
	}
	if elapsed > 11*time.Second {
		t.Errorf("AndCombiner closed too late. Got %v, expected ~10s", elapsed)
	}

	t.Logf("Combiner correctly closed after %v", elapsed)
}
