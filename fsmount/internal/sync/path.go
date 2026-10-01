//go:build linux || windows

package sync

import "sync"

// PathMutex provides per-path locking to allow concurrent operations on different paths
type PathMutex struct {
	mu    sync.Mutex
	locks map[string]*sync.Mutex
}

// PathRWMutex provides shared read locks and exclusive mutation locks per path.
type PathRWMutex struct {
	mu    sync.Mutex
	locks map[string]*sync.RWMutex
}

// NewPathMutex creates a new path mutex
func NewPathMutex() *PathMutex {
	return &PathMutex{
		locks: make(map[string]*sync.Mutex),
	}
}

func NewPathRWMutex() *PathRWMutex {
	return &PathRWMutex{locks: make(map[string]*sync.RWMutex)}
}

func (pm *PathRWMutex) pathLock(path string) *sync.RWMutex {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	lock, ok := pm.locks[path]
	if !ok {
		lock = &sync.RWMutex{}
		pm.locks[path] = lock
	}
	return lock
}

func (pm *PathRWMutex) Lock(path string) {
	pm.pathLock(path).Lock()
}

func (pm *PathRWMutex) Unlock(path string) {
	pm.pathLock(path).Unlock()
}

func (pm *PathRWMutex) RLock(path string) {
	pm.pathLock(path).RLock()
}

func (pm *PathRWMutex) RUnlock(path string) {
	pm.pathLock(path).RUnlock()
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
