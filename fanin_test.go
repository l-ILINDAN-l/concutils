package concutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestFanInMerger checks that it correctly merges multiple channels
func TestFanInMerger(t *testing.T) {
	// producer is a helper function that creates a channel and sends `count` messages to it
	producer := func(count int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for i := 0; i < count; i++ {
				out <- i
			}
		}()
		return out
	}

	// Arrange: Create a merger for three producers
	merger := NewFanInMerger(
		producer(3),
		producer(5),
		producer(2),
	)

	// Act: Read all messages from the merged channel
	count := 0
	for range merger.Out() {
		count++
	}

	// Assert: The total count should be the sum of messages from all producers
	expectedCount := 3 + 5 + 2
	assert.Equal(t, expectedCount, count, "should receive all messages from all producers")
}
