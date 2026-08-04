package fsmount

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/disk"
)

// DiskCacheSize returns the total persistent cache size across this process's
// active mounts. Cache directories not owned by this process are not included.
func DiskCacheSize() (int64, error) {
	activeCaches, err := snapshotActiveDiskCaches()
	if err != nil {
		return 0, err
	}
	var sizeBytes int64
	for _, cache := range activeCaches {
		sizeBytes += cache.SizeBytes()
	}
	return sizeBytes, nil
}

// ClearDiskCache clears this process's active caches. Cache directories not
// owned by this process are left untouched. It returns the bytes still pinned
// by open file handles.
func ClearDiskCache() (int64, error) {
	// Keep mount creation paused until the clear finishes. A new mount can reuse
	// a recently unmounted cache path, so allowing one to start while an old
	// cache object is being cleared could delete the new mount's data.
	mountMu.Lock()
	defer mountMu.Unlock()

	activeCaches, err := activeDiskCacheSnapshot()
	if err != nil {
		return 0, err
	}
	var clearErrors []error
	for _, cache := range activeCaches {
		if err := cache.Clear(); err != nil {
			clearErrors = append(clearErrors, err)
		}
	}

	var remainingBytes int64
	for _, cache := range activeCaches {
		remainingBytes += cache.SizeBytes()
	}
	return remainingBytes, errors.Join(clearErrors...)
}

func snapshotActiveDiskCaches() ([]*disk.DiskCache, error) {
	mountMu.Lock()
	defer mountMu.Unlock()
	return activeDiskCacheSnapshot()
}

// activeDiskCacheSnapshot returns a stable set of cache pointers without
// retaining the registry lock. The caller must hold mountMu.
func activeDiskCacheSnapshot() ([]*disk.DiskCache, error) {
	var caches []*disk.DiskCache
	for _, cache := range registeredDiskCaches() {
		alreadyAdded := false
		for _, existing := range caches {
			sameRoot, err := sameDiskCacheRoot(existing.CacheRoot, cache.CacheRoot)
			if err != nil {
				return nil, err
			}
			if !sameRoot {
				continue
			}
			if existing != cache {
				return nil, fmt.Errorf("%w: %s", errDiskCachePathInUse, cache.CacheRoot)
			}
			alreadyAdded = true
			break
		}
		if !alreadyAdded {
			caches = append(caches, cache)
		}
	}
	return caches, nil
}

func ensureDiskCacheRootAvailable(cacheRoot string) error {
	for _, cache := range registeredDiskCaches() {
		sameRoot, err := sameDiskCacheRoot(cache.CacheRoot, cacheRoot)
		if err != nil {
			return err
		}
		if sameRoot {
			return fmt.Errorf("%w: %s", errDiskCachePathInUse, cacheRoot)
		}
	}
	return nil
}

func registeredDiskCaches() []*disk.DiskCache {
	if mntRegistry == nil {
		return nil
	}

	mntRegistry.hostsMu.Lock()
	defer mntRegistry.hostsMu.Unlock()
	var caches []*disk.DiskCache
	for _, host := range mntRegistry.hosts {
		if host == nil || host.fs == nil || host.fs.remote == nil {
			continue
		}
		cache, ok := host.fs.remote.cacheStore.(*disk.DiskCache)
		if !ok || cache.Disabled {
			continue
		}
		caches = append(caches, cache)
	}
	return caches
}

func canonicalDiskCacheRoot(path string) (string, error) {
	root, err := filepath.EvalSymlinks(filepath.Clean(path))
	if err != nil {
		return "", fmt.Errorf("failed to resolve disk cache path: %w", err)
	}
	return filepath.Clean(root), nil
}

func sameDiskCacheRoot(left, right string) (bool, error) {
	leftInfo, err := os.Stat(left)
	if err != nil {
		return false, fmt.Errorf("failed to stat disk cache path %s: %w", left, err)
	}
	rightInfo, err := os.Stat(right)
	if err != nil {
		return false, fmt.Errorf("failed to stat disk cache path %s: %w", right, err)
	}
	return os.SameFile(leftInfo, rightInfo), nil
}
