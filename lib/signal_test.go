package lib

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignalConcurrentCallsCloseItOnce(t *testing.T) {
	for attempt := 0; attempt < 2000; attempt++ {
		signal := (&Signal{}).Init()
		// A copy made after Init shares the same channel.
		copied := *signal
		ready := make(chan struct{})
		panics := make(chan any, 32)
		var wg sync.WaitGroup
		for n := 0; n < 32; n++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer func() {
					if v := recover(); v != nil {
						panics <- v
					}
				}()
				<-ready
				if n%2 == 0 {
					signal.Call()
				} else {
					copied.Call()
				}
			}()
		}
		close(ready)
		wg.Wait()
		close(panics)
		for v := range panics {
			t.Fatalf("concurrent Call panicked on attempt %d: %v", attempt, v)
		}
		assert.True(t, signal.Called())
	}

	// Once the calls have settled, Clear starts a new generation.
	signal := (&Signal{}).Init()
	signal.Call()
	signal.Clear()
	assert.False(t, signal.Called())
	signal.Call()
	assert.True(t, signal.Called())
}
