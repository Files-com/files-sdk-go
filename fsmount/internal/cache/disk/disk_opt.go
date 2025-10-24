package disk

import (
	"time"

	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/log"
)

// Option configures a DiskCache.
type Option func(*DiskCache)

// WithCapacityBytes sets the maximum total size (in bytes) for the cache.
func WithCapacityBytes(n int64) Option {
	return func(dc *DiskCache) {
		dc.Capacity = n
	}
}

// WithDisabled disables the cache (no I/O, no directory creation).
func WithDisabled(disabled bool) Option {
	return func(dc *DiskCache) {
		dc.Disabled = disabled
	}
}

// WithMaxAge sets the maximum age for cache files.
func WithMaxAge(d time.Duration) Option {
	return func(dc *DiskCache) {
		dc.MaxAge = d
	}
}

// WithMaxFileCount sets the maximum number of files for the cache.
func WithMaxFileCount(n int64) Option {
	return func(dc *DiskCache) {
		dc.MaxFileCount = n
	}
}

// WithMaintenanceInterval sets the interval for cache maintenance operations.
func WithMaintenanceInterval(d time.Duration) Option {
	return func(dc *DiskCache) {
		dc.MaintenanceInterval = d
	}
}

// WithLruFlushInterval sets the interval for cache maintenance operations.
func WithLruFlushInterval(d time.Duration) Option {
	return func(dc *DiskCache) {
		dc.LruFlushInterval = d
	}
}

// WithLogger sets the logger for the cache.
func WithLogger(logger log.Logger) Option {
	return func(dc *DiskCache) {
		dc.log = logger
	}
}
