package concutils

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

// TestFanOutDistributor checks that it correctly distributes items to multiple channels
func TestFanOutDistributor(t *testing.T) {
	// Arrange
	inChan := make(chan int)
	numItems := 100
	numWorkers := 5

	// Act
	// 1. Create the distributor
	distributor := NewFanOutDistributor(inChan, numWorkers)
	outputChans := distributor.Outs()

	// 2. Start a producer to send items to the input channel
	go func() {
		defer close(inChan)
		for i := 0; i < numItems; i++ {
			inChan <- i
		}
	}()

	// 3. Start consumers to read from the output channels
	var wg sync.WaitGroup
	var receivedCount int
	var mu sync.Mutex

	for _, out := range outputChans {
		wg.Add(1)
		go func(ch <-chan int) {
			defer wg.Done()
			for range ch {
				mu.Lock()
				receivedCount++
				mu.Unlock()
			}
		}(out)
	}

	// 4. Wait for all consumers to finish
	wg.Wait()

	// Assert
	assert.Equal(t, numItems, receivedCount, "the total number of received items should match the number of sent items")
}
