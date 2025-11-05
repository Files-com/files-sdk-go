//go:build !filescomfs_debug
// +build !filescomfs_debug

package cache

import (
	"sync/atomic"
	"time"
)

// Stats maintains statistics about cache usage.
// In non-debug builds, this type exists but has no methods for retrieving statistics.
type Stats struct {
	CapacityBytes          int64
	CapacityBytesRemaining int64
	MaxFileCount           int64
	FileCountRemaining     int64
	FileCount              atomic.Int64
	ReadBytes              atomic.Int64
	ReadCount              atomic.Int64
	SizeBytes              atomic.Int64
	WriteBytes             atomic.Int64
	WriteCount             atomic.Int64
	LoadDuration           time.Duration
	LruCount               int
	PinnedCount            int
	PinnedPaths            map[string]int
	LruKeys                []string
}

// ResetCounters sets "counter" type stats to zero.
func (st *Stats) ResetCounters() {
	st.ReadBytes.Store(0)
	st.ReadCount.Store(0)
	st.WriteBytes.Store(0)
	st.WriteCount.Store(0)
}
