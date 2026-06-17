package file

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestUploadIOMarkTransferStartedIsConcurrentOneShot(t *testing.T) {
	u := &uploadIO{
		Progress:        func(int64) {},
		transferStarted: &atomic.Bool{},
	}

	var progressCalls atomic.Int64
	u.Progress = func(int64) {
		progressCalls.Add(1)
	}

	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			u.markTransferStarted()
		}()
	}

	close(start)
	wg.Wait()

	if got := progressCalls.Load(); got != 1 {
		t.Fatalf("expected one progress start signal, got %d", got)
	}
}
