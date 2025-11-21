package sync_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	lim "github.com/Files-com/files-sdk-go/v3/fsmount/internal/limit"
)

func TestFuseOpLimiter_Basic(t *testing.T) {
	limiter := lim.NewFuseOpLimiter(map[lim.FuseOpType]int64{
		lim.FuseOpUpload:   2,
		lim.FuseOpDownload: 3,
	}, 0)

	ctx := context.Background()
	called := false

	err := limiter.WithLimit(ctx, lim.FuseOpUpload, func(ctx context.Context) error {
		called = true
		return nil
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !called {
		t.Fatal("expected function to be called")
	}
}

func TestFuseOpLimiter_FunctionError(t *testing.T) {
	limiter := lim.NewFuseOpLimiter(map[lim.FuseOpType]int64{
		lim.FuseOpUpload: 1,
	}, 0)

	ctx := context.Background()
	expectedErr := errors.New("test error")

	err := limiter.WithLimit(ctx, lim.FuseOpUpload, func(ctx context.Context) error {
		return expectedErr
	})

	if err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}

func TestFuseOpLimiter_PerClassLimit(t *testing.T) {
	limiter := lim.NewFuseOpLimiter(map[lim.FuseOpType]int64{
		lim.FuseOpUpload: 2,
	}, 0)

	ctx := context.Background()
	var running int32
	var maxRunning int32

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = limiter.WithLimit(ctx, lim.FuseOpUpload, func(ctx context.Context) error {
				current := atomic.AddInt32(&running, 1)
				defer atomic.AddInt32(&running, -1)

				// Track max concurrent operations
				for {
					max := atomic.LoadInt32(&maxRunning)
					if current <= max || atomic.CompareAndSwapInt32(&maxRunning, max, current) {
						break
					}
				}

				time.Sleep(10 * time.Millisecond)
				return nil
			})
		}()
	}

	wg.Wait()

	max := atomic.LoadInt32(&maxRunning)
	if max > 2 {
		t.Fatalf("expected max 2 concurrent operations, got %d", max)
	}
	if max < 2 {
		t.Logf("warning: only saw %d concurrent operations, expected to see 2", max)
	}
}

func TestFuseOpLimiter_GlobalLimit(t *testing.T) {
	limiter := lim.NewFuseOpLimiter(map[lim.FuseOpType]int64{
		lim.FuseOpUpload:   5,
		lim.FuseOpDownload: 5,
	}, 3) // global limit of 3

	ctx := context.Background()
	var running int32
	var maxRunning int32

	var wg sync.WaitGroup
	// Mix of upload and download operations
	for i := 0; i < 4; i++ {
		wg.Add(1)
		opType := lim.FuseOpUpload
		if i%2 == 0 {
			opType = lim.FuseOpDownload
		}
		go func(op lim.FuseOpType) {
			defer wg.Done()
			_ = limiter.WithLimit(ctx, op, func(ctx context.Context) error {
				current := atomic.AddInt32(&running, 1)
				defer atomic.AddInt32(&running, -1)

				// Track max concurrent operations
				for {
					max := atomic.LoadInt32(&maxRunning)
					if current <= max || atomic.CompareAndSwapInt32(&maxRunning, max, current) {
						break
					}
				}

				time.Sleep(10 * time.Millisecond)
				return nil
			})
		}(opType)
	}

	wg.Wait()

	max := atomic.LoadInt32(&maxRunning)
	if max > 3 {
		t.Fatalf("expected max 3 concurrent operations (global limit), got %d", max)
	}
	if max < 3 {
		t.Logf("warning: only saw %d concurrent operations, expected to see 3", max)
	}
}

func TestFuseOpLimiter_Reentrancy(t *testing.T) {
	limiter := lim.NewFuseOpLimiter(map[lim.FuseOpType]int64{
		lim.FuseOpUpload: 1,
	}, 0)

	ctx := context.Background()
	var depth int

	err := limiter.WithLimit(ctx, lim.FuseOpUpload, func(ctx1 context.Context) error {
		depth++
		// Nested call with same operation type should not block
		return limiter.WithLimit(ctx1, lim.FuseOpUpload, func(ctx2 context.Context) error {
			depth++
			// Third level
			return limiter.WithLimit(ctx2, lim.FuseOpUpload, func(ctx3 context.Context) error {
				depth++
				return nil
			})
		})
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if depth != 3 {
		t.Fatalf("expected depth 3, got %d", depth)
	}
}

func TestFuseOpLimiter_ReentrancyDoesNotBypassLimit(t *testing.T) {
	limiter := lim.NewFuseOpLimiter(map[lim.FuseOpType]int64{
		lim.FuseOpUpload:   1,
		lim.FuseOpDownload: 1,
	}, 0)

	ctx := context.Background()
	blocked := make(chan struct{})
	released := make(chan struct{})

	go func() {
		_ = limiter.WithLimit(ctx, lim.FuseOpUpload, func(ctx1 context.Context) error {
			close(blocked)
			<-released
			return nil
		})
	}()

	// Wait for first goroutine to acquire the limit
	<-blocked

	// Try to acquire from another goroutine (should block)
	acquired := make(chan struct{})
	go func() {
		_ = limiter.WithLimit(context.Background(), lim.FuseOpUpload, func(ctx context.Context) error {
			close(acquired)
			return nil
		})
	}()

	// Give it a moment to try to acquire
	select {
	case <-acquired:
		t.Fatal("second goroutine should not have acquired the limit")
	case <-time.After(50 * time.Millisecond):
		// Expected: still blocked
	}

	// Release the first goroutine
	close(released)

	// Now the second should be able to acquire
	select {
	case <-acquired:
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("second goroutine should have acquired the limit after release")
	}
}

func TestFuseOpLimiter_ContextCancellation(t *testing.T) {
	limiter := lim.NewFuseOpLimiter(map[lim.FuseOpType]int64{
		lim.FuseOpUpload: 1,
	}, 0)

	// First, occupy the limit
	ctx1 := context.Background()
	blocked := make(chan struct{})
	released := make(chan struct{})

	go func() {
		_ = limiter.WithLimit(ctx1, lim.FuseOpUpload, func(ctx context.Context) error {
			close(blocked)
			<-released
			return nil
		})
	}()

	<-blocked

	// Try to acquire with a cancelled context
	ctx2, cancel := context.WithCancel(context.Background())
	cancel()

	err := limiter.WithLimit(ctx2, lim.FuseOpUpload, func(ctx context.Context) error {
		t.Fatal("function should not be called with cancelled context")
		return nil
	})

	if err != context.Canceled {
		t.Fatalf("expected context.Canceled error, got %v", err)
	}

	close(released)
}

func TestFuseOpLimiter_NoLimit(t *testing.T) {
	// Create limiter with no limits set
	limiter := lim.NewFuseOpLimiter(map[lim.FuseOpType]int64{}, 0)

	ctx := context.Background()
	var running int32
	var maxRunning int32

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = limiter.WithLimit(ctx, lim.FuseOpUpload, func(ctx context.Context) error {
				current := atomic.AddInt32(&running, 1)
				defer atomic.AddInt32(&running, -1)

				// Track max concurrent operations
				for {
					max := atomic.LoadInt32(&maxRunning)
					if current <= max || atomic.CompareAndSwapInt32(&maxRunning, max, current) {
						break
					}
				}

				time.Sleep(5 * time.Millisecond)
				return nil
			})
		}()
	}

	wg.Wait()

	max := atomic.LoadInt32(&maxRunning)
	// With no limit, we should see high concurrency
	if max < 5 {
		t.Logf("warning: only saw %d concurrent operations, expected more without limits", max)
	}
}

func TestFuseOpLimiter_DifferentClassesSeparate(t *testing.T) {
	limiter := lim.NewFuseOpLimiter(map[lim.FuseOpType]int64{
		lim.FuseOpUpload:   1,
		lim.FuseOpDownload: 1,
	}, 0)

	ctx := context.Background()
	var uploadRunning int32
	var downloadRunning int32

	uploadBlocked := make(chan struct{})
	downloadBlocked := make(chan struct{})
	released := make(chan struct{})

	// Start upload operation
	go func() {
		_ = limiter.WithLimit(ctx, lim.FuseOpUpload, func(ctx context.Context) error {
			atomic.StoreInt32(&uploadRunning, 1)
			close(uploadBlocked)
			<-released
			atomic.StoreInt32(&uploadRunning, 0)
			return nil
		})
	}()

	<-uploadBlocked

	// Start download operation - should not be blocked by upload
	go func() {
		_ = limiter.WithLimit(ctx, lim.FuseOpDownload, func(ctx context.Context) error {
			atomic.StoreInt32(&downloadRunning, 1)
			close(downloadBlocked)
			<-released
			atomic.StoreInt32(&downloadRunning, 0)
			return nil
		})
	}()

	<-downloadBlocked

	// Both should be running simultaneously
	if atomic.LoadInt32(&uploadRunning) != 1 {
		t.Fatal("upload should be running")
	}
	if atomic.LoadInt32(&downloadRunning) != 1 {
		t.Fatal("download should be running")
	}

	close(released)
}

func TestFuseOpLimiter_GlobalAcquiredFirst(t *testing.T) {
	// This test verifies that global limit is acquired before class limit
	// to avoid potential deadlocks
	limiter := lim.NewFuseOpLimiter(map[lim.FuseOpType]int64{
		lim.FuseOpUpload: 2,
	}, 1) // global limit of 1

	ctx := context.Background()
	var running int32

	blocked := make(chan struct{})
	released := make(chan struct{})

	go func() {
		_ = limiter.WithLimit(ctx, lim.FuseOpUpload, func(ctx context.Context) error {
			atomic.StoreInt32(&running, 1)
			close(blocked)
			<-released
			atomic.StoreInt32(&running, 0)
			return nil
		})
	}()

	<-blocked

	// Try to acquire another upload operation
	// Global limit (1) should block this even though class limit (2) has room
	acquired := make(chan struct{})
	go func() {
		_ = limiter.WithLimit(context.Background(), lim.FuseOpUpload, func(ctx context.Context) error {
			close(acquired)
			return nil
		})
	}()

	select {
	case <-acquired:
		t.Fatal("second operation should be blocked by global limit")
	case <-time.After(50 * time.Millisecond):
		// Expected: blocked by global limit
	}

	close(released)

	// Now should be able to proceed
	select {
	case <-acquired:
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("second operation should have proceeded after release")
	}
}
