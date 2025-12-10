package sync

import (
	"sync"
	"testing"
	"time"
)

func TestPathMutex_BasicLockUnlock(t *testing.T) {
	pm := NewPathMutex()
	path := "/test/path"

	// Should be able to lock and unlock without blocking
	pm.Lock(path)
	pm.Unlock(path)
}

func TestPathMutex_ConcurrentAccessDifferentPaths(t *testing.T) {
	pm := NewPathMutex()
	path1 := "/test/path1"
	path2 := "/test/path2"

	var wg sync.WaitGroup
	wg.Add(2)

	// Both goroutines should run concurrently without blocking each other
	start := time.Now()

	go func() {
		defer wg.Done()
		pm.Lock(path1)
		time.Sleep(50 * time.Millisecond)
		pm.Unlock(path1)
	}()

	go func() {
		defer wg.Done()
		pm.Lock(path2)
		time.Sleep(50 * time.Millisecond)
		pm.Unlock(path2)
	}()

	wg.Wait()
	elapsed := time.Since(start)

	// If locks were serialized, it would take ~100ms. With concurrent access, it should be ~50ms
	if elapsed > 80*time.Millisecond {
		t.Errorf("Expected concurrent execution (~50ms), but took %v", elapsed)
	}
}

func TestPathMutex_ConcurrentAccessSamePath(t *testing.T) {
	pm := NewPathMutex()
	path := "/test/path"

	var wg sync.WaitGroup
	wg.Add(2)

	// Both goroutines access the same path, so they should serialize
	start := time.Now()

	go func() {
		defer wg.Done()
		pm.Lock(path)
		time.Sleep(50 * time.Millisecond)
		pm.Unlock(path)
	}()

	go func() {
		defer wg.Done()
		time.Sleep(10 * time.Millisecond) // Ensure first goroutine gets the lock first
		pm.Lock(path)
		time.Sleep(50 * time.Millisecond)
		pm.Unlock(path)
	}()

	wg.Wait()
	elapsed := time.Since(start)

	// With serialized access, it should take at least 100ms
	if elapsed < 100*time.Millisecond {
		t.Errorf("Expected serialized execution (>100ms), but took %v", elapsed)
	}
}

func TestPathMutex_MultiplePathsNoDeadlock(t *testing.T) {
	pm := NewPathMutex()
	paths := []string{"/path1", "/path2", "/path3", "/path4", "/path5"}

	var wg sync.WaitGroup
	for _, path := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				pm.Lock(p)
				time.Sleep(time.Millisecond)
				pm.Unlock(p)
			}
		}(path)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Test passed - no deadlock
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out - possible deadlock")
	}
}
