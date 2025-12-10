package sync

import "sync"

// PathMutex provides per-path locking to allow concurrent operations on different paths
type PathMutex struct {
	mu    sync.Mutex
	locks map[string]*sync.Mutex
}

// NewPathMutex creates a new path mutex
func NewPathMutex() *PathMutex {
	return &PathMutex{
		locks: make(map[string]*sync.Mutex),
	}
}

// Lock acquires a lock for the given path
func (pm *PathMutex) Lock(path string) {
	pm.mu.Lock()
	mu, ok := pm.locks[path]
	if !ok {
		mu = &sync.Mutex{}
		pm.locks[path] = mu
	}
	pm.mu.Unlock()
	mu.Lock()
}

// Unlock releases the lock for the given path
func (pm *PathMutex) Unlock(path string) {
	pm.mu.Lock()
	mu, ok := pm.locks[path]
	pm.mu.Unlock()
	if ok {
		mu.Unlock()
	}
}
