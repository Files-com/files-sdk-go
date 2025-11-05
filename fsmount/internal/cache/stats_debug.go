//go:build filescomfs_debug
// +build filescomfs_debug

package cache

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/dustin/go-humanize"
)

// Stats maintains statistics about cache usage.
type Stats struct {
	// CapacityBytes is the maximum size of the cache in bytes. 0 means unbounded.
	//
	// This is set when the cache is created and does not change.
	CapacityBytes int64

	// CapacityBytesRemaining is the remaining size of the cache in bytes.
	//
	// This is calculated as CapacityBytes - SizeBytes.
	CapacityBytesRemaining int64

	// MaxFileCount is the maximum number of files in the cache. 0 means unbounded.
	//
	// This is set when the cache is created and does not change.
	MaxFileCount int64

	// FileCountRemaining is the remaining number of files that can be stored in the cache.
	//
	// This is calculated as MaxFileCount - FileCount.
	FileCountRemaining int64

	// FileCount represents the number of files in the cache.
	FileCount atomic.Int64

	// ReadBytes is the total number of bytes read from the cache.
	// This is incremented on each successful Read call.
	ReadBytes atomic.Int64

	// ReadCount is the total number of Read calls.
	ReadCount atomic.Int64

	// SizeBytes is the total size of all files in the cache.
	SizeBytes atomic.Int64

	// WriteBytes is the total number of bytes written to the cache.
	// This is incremented on each successful Write call.
	WriteBytes atomic.Int64

	// WriteCount is the total number of Write calls.
	WriteCount atomic.Int64

	// LoadDuration is the time taken to refresh in-memory cache statistics from disk during the most recent load.
	LoadDuration time.Duration

	// LruCount is the number of items in the LRU cache.
	LruCount int

	// PinnedCount is the number of files currently pinned in the cache.
	// Pinned files cannot be evicted because they have open file handles.
	PinnedCount int

	// PinnedPaths is a map of pinned file paths to their reference counts.
	// This shows which files are currently pinned and how many handles are open.
	PinnedPaths map[string]int

	// LruKeys is a list of all keys in the LRU cache, ordered from oldest to newest.
	// This shows the eviction order - the first item will be evicted next when space is needed.
	LruKeys []string
}

// ResetCounters sets "counter" type stats to zero.
//
// It does not reset Capacity, CapacityRemaining, FileCount, or SizeBytes.
func (st *Stats) ResetCounters() {
	st.ReadBytes.Store(0)
	st.ReadCount.Store(0)
	st.WriteBytes.Store(0)
	st.WriteCount.Store(0)
}

// MarshalJSON implements the json.Marshaler interface for Stats.
// It ensures CapacityRemaining is calculated as Capacity - SizeBytes.
func (st *Stats) MarshalJSON() ([]byte, error) {
	type alias struct {
		Capacity          int64 `json:"capacity"`
		SizeBytes         int64 `json:"size_bytes"`
		CapacityRemaining int64 `json:"capacity_remaining"`

		CapacityHuman          string `json:"capacity_human"`
		SizeBytesHuman         string `json:"size_bytes_human"`
		CapacityRemainingHuman string `json:"capacity_remaining_human"`

		MaxFileCount       int64 `json:"max_file_count"`
		FileCount          int64 `json:"file_count"`
		FileCountRemaining int64 `json:"file_count_remaining"`

		ReadBytes      int64  `json:"read_bytes"`
		ReadBytesHuman string `json:"read_bytes_human"`
		ReadCount      int64  `json:"read_count"`

		WriteBytes      int64  `json:"write_bytes"`
		WriteBytesHuman string `json:"write_bytes_human"`
		WriteCount      int64  `json:"write_count"`

		LoadDuration string         `json:"load_duration"`
		LruCount     int            `json:"lru_count"`
		PinnedCount  int            `json:"pinned_count"`
		PinnedPaths  map[string]int `json:"pinned_paths"`
		LruKeys      []string       `json:"lru_keys"`
	}
	al := alias{
		Capacity:               st.CapacityBytes,
		CapacityRemaining:      st.CapacityBytes - st.SizeBytes.Load(),
		MaxFileCount:           st.MaxFileCount,
		FileCountRemaining:     st.MaxFileCount - st.FileCount.Load(),
		FileCount:              st.FileCount.Load(),
		ReadBytes:              st.ReadBytes.Load(),
		ReadCount:              st.ReadCount.Load(),
		SizeBytes:              st.SizeBytes.Load(),
		WriteBytes:             st.WriteBytes.Load(),
		WriteCount:             st.WriteCount.Load(),
		LruCount:               st.LruCount,
		PinnedCount:            st.PinnedCount,
		PinnedPaths:            st.PinnedPaths,
		LruKeys:                st.LruKeys,
		LoadDuration:           fmt.Sprintf("%v ms", st.LoadDuration.Milliseconds()),
		CapacityHuman:          humanize.IBytes(uint64(st.CapacityBytes)),
		CapacityRemainingHuman: humanize.IBytes(uint64(st.CapacityBytes - st.SizeBytes.Load())),
		ReadBytesHuman:         humanize.IBytes(uint64(st.ReadBytes.Load())),
		SizeBytesHuman:         humanize.IBytes(uint64(st.SizeBytes.Load())),
		WriteBytesHuman:        humanize.IBytes(uint64(st.WriteBytes.Load())),
	}
	return json.Marshal(&al)
}
