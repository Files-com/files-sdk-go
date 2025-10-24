package mem

import (
	"time"

	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/log"
)

// Option configures a MemoryCache.
type Option func(*MemoryCache)

// WithCapacityBytes sets the maximum total size (in bytes) for the cache.
func WithCapacityBytes(n int64) Option {
	return func(dc *MemoryCache) {
		dc.Capacity = n
	}
}

// WithMaxAge sets the maximum age for cache files.
func WithMaxAge(d time.Duration) Option {
	return func(dc *MemoryCache) {
		dc.MaxAge = d
	}
}

// WithMaxFileCount sets the maximum number of files for the cache.
func WithMaxFileCount(n int64) Option {
	return func(dc *MemoryCache) {
		dc.MaxFileCount = n
	}
}

// WithMaintenanceInterval sets the interval for cache maintenance operations.
func WithMaintenanceInterval(d time.Duration) Option {
	return func(dc *MemoryCache) {
		dc.MaintenanceInterval = d
	}
}

// WithLogger sets the logger for the cache.
func WithLogger(logger log.Logger) Option {
	return func(dc *MemoryCache) {
		dc.log = logger
	}
}
