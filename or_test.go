package concutils

import (
	"testing"
	"time"
)

// TestOrCombiner checks th logic of the OrCombiner
func TestOrCombiner(t *testing.T) {
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
	// We expect the combiner to close when the fastest channel (100ms) closes
	combiner := NewOrCombiner(
		sig(1*time.Hour),
		sig(1*time.Minute),
		sig(100*time.Millisecond), // This channel
		sig(10*time.Second),
	)

	// Wait closing combined channel
	<-combiner.Out()

	// Assert: Check if the elapsed time is reasonable
	elapsed := time.Since(start)

	if elapsed < 100*time.Millisecond {
		t.Errorf("OrCombiner closed too early. Got %v, expected >100ms", elapsed)
	}
	if elapsed > 500*time.Millisecond { // Using a generous upper bound
		t.Errorf("OrCombiner closed too late. Got %v, expected ~100ms", elapsed)
	}

	t.Logf("Combiner correctly closed after %v", elapsed)
}
