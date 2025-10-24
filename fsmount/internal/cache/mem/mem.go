package mem

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/log"

	lru "github.com/hashicorp/golang-lru/v2"
)

const (
	// DefaultCacheCapacity is the maximum size of the cache in bytes if not provided
	// at cache initialization. Equivalent to 1GiB
	DefaultCapacity = 1 << 30

	// MinCapacity is the minimum size of the cache. Equivalent to 256MiB
	MinCapacity = 256 << 20

	// MaxCapacity is the maximum size of the cache. Equivalent to 5GiB
	MaxCapacity = 5 << 30

	// DefaultMaxAge is the maximum age of cache files, if not provided
	// at cache initialization. 0 means no expiration. Equivalent to 12 hours.
	DefaultMaxAge = 12 * time.Hour

	// DefaultMaxFileCount is the maximum number of files in the cache if not provided
	// at cache initialization. 0 means unbounded.
	DefaultMaxFileCount = 1000

	// DefaultMaintenanceInterval is the interval at which cache maintenance tasks are performed,
	// if not provided at cache initialization.
	DefaultMaintenanceInterval = 5 * time.Minute

	// unboundedFileCount is a large number used to represent an unbounded file count in the LRU cache.
	//
	// Techically the LRU cache is bounded, so this is a sufficiently large number to approximate
	// unbounded behavior. Equivalent to 4,294,967,295
	unboundedFileCount = 1 << 32

	// pageSize is the size of each byte slice that is allocated to hold the contents of a file.
	// Equivalent to 128KiB.
	pageSize = 128 << 10
)

// page is a convenient way to think about a chunk of an individual file
type page = []byte

// entry represents a single item in the cache
type entry struct {
	// pages is the map of offsets -> pages
	pages map[int64]page

	// logical file size (max written ofst)
	size int64

	// the last time the file was written to the cache
	mod time.Time
}

// MemoryCache implements a simple in memory LRU cache for file data with an interface that roughly
// matches the FUSE Read/Write methods.
//
// TODO:
//   - allow concurrent reads/writes to different files in the cache by moving lock coordination
//     out of the caller and lock by path within the cache with a RWMutex map. This will require
//     maintaining a reservation of bytes so that concurrent goroutines can't exhaust memory that
//     a different goroutine evicted for its Write operation.
type MemoryCache struct {
	Capacity            int64
	MaxFileCount        int64
	MaxAge              time.Duration
	MaintenanceInterval time.Duration

	writeMu sync.Mutex

	// filesMu protects access to the files map
	filesMu sync.RWMutex

	// files represents the mapping from a path to its entry
	files map[string]*entry

	// Global LRU across *pages* (or across files—see below)
	lru *lru.Cache[string, struct{}]

	stats cache.Stats

	// used to log cache operations
	log log.Logger

	maintMu     sync.Mutex
	maintActive bool
	maintCancel context.CancelFunc
	maintDone   chan struct{}
}

func NewMemoryCache(opts ...Option) (*MemoryCache, error) {
	mc := &MemoryCache{
		Capacity:            DefaultCapacity,
		MaintenanceInterval: DefaultMaintenanceInterval,
		MaxAge:              DefaultMaxAge,
		MaxFileCount:        DefaultMaxFileCount,
		files:               make(map[string]*entry),
	}

	// apply options
	for _, opt := range opts {
		opt(mc)
	}

	// validate the options before initializing LRU/scan
	if err := mc.validateOpts(); err != nil {
		return nil, err
	}

	// do this after applying options in case the caller set a different max file count
	// or provided their own LRU cache
	if mc.lru == nil {
		// protect against zero or negative max file count
		var trackCap int64
		if mc.MaxFileCount > 0 {
			trackCap = mc.MaxFileCount
		} else {
			// Unbounded file-count; use a large number since the LRU cache cannot be truly unbounded.
			trackCap = unboundedFileCount
		}
		lru, err := lru.NewWithEvict(int(trackCap), mc.onEvict)
		if err != nil {
			return nil, fmt.Errorf("memoryCache: error creating LRU cache: %w", err)
		}
		mc.lru = lru
	}
	mc.stats = cache.Stats{
		CapacityBytes: mc.Capacity,
		MaxFileCount:  mc.MaxFileCount,
	}

	return mc, nil
}

func (mc *MemoryCache) Write(path string, buff []byte, ofst int64) (int, error) {
	if len(buff) == 0 {
		return 0, nil
	}

	// prevent concurrent access to the map of files
	mc.filesMu.Lock()
	ent, ok := mc.files[path]
	newFile := false
	if !ok {
		newFile = true
		ent = &entry{pages: make(map[int64]page), mod: time.Now()}
	}
	mc.filesMu.Unlock()

	// make sure the entry is not toward the end of the LRU before trimming for capacity
	_, _ = mc.lru.Get(path)

	// lock writes to avoid another goroutine modifying the capacity during this Write
	mc.writeMu.Lock()
	defer mc.writeMu.Unlock()

	// enure there is enough capacity in the cache for the incoming write in terms of bytes and file count
	need := bytesNeededForWrite(ent, ofst, int64(len(buff)))
	for !mc.hasCapacityDelta(need, newFile) {
		if _, _, ok := mc.lru.RemoveOldest(); !ok {
			return 0, fmt.Errorf("cache: at capacity, and nothing deletable")
		}
	}

	mc.filesMu.Lock()
	// Get was used above to avoid trimming the current Write's data, but the file may not have
	// been in the cache to being with
	if newFile {
		mc.files[path] = ent
		mc.stats.FileCount.Add(1)
	}

	// Write into pages (allocating as needed, updating SizeBytes on *new* pages)
	n := 0
	for n < len(buff) {
		abs := ofst + int64(n)
		pageIdx := abs / pageSize
		pageOff := int(abs % pageSize)

		// allocate page if missing
		page, ok := ent.pages[pageIdx]
		if !ok {
			page = make([]byte, pageSize)
			ent.pages[pageIdx] = page
			mc.stats.SizeBytes.Add(int64(len(page))) // count newly allocated memory
		}

		remainingInSrc := len(buff) - n
		remainingInPage := pageSize - pageOff

		// the number of bytes that should get copied from buff to the page
		// this iteration
		chunkLen := min(remainingInSrc, remainingInPage)
		copy(page[pageOff:pageOff+chunkLen], buff[n:n+chunkLen])
		n += chunkLen
	}

	if ofst+int64(n) > ent.size {
		ent.size = ofst + int64(n)
	}
	ent.mod = time.Now()
	mc.filesMu.Unlock()
	_ = mc.lru.Add(path, struct{}{})
	mc.stats.WriteCount.Add(1)
	mc.stats.WriteBytes.Add(int64(n))
	return n, nil
}

// Read reads data from the cached file at the given path into buff starting at the provided offset.
func (mc *MemoryCache) Read(path string, buff []byte, ofst int64) (n int, err error) {
	if len(buff) == 0 {
		return 0, nil
	}

	mc.filesMu.RLock()
	ent, ok := mc.files[path]
	if !ok {
		mc.filesMu.RUnlock()
		return 0, nil // cache miss
	}
	if ofst >= ent.size {
		mc.filesMu.RUnlock()
		return 0, nil // EOF
	}

	// clamp requested length to EOF
	maxReadable := ent.size - ofst
	want := min(int64(len(buff)), maxReadable)

	read := 0
	for read < int(want) {
		abs := ofst + int64(read)
		pageIdx := abs / pageSize
		pageOff := int(abs % pageSize)

		page, ok := ent.pages[pageIdx]
		if !ok {
			// hit a hole (not cached yet) – stop so caller can fetch from remote
			break
		}

		can := min(int(want)-read, pageSize-pageOff, len(buff)-read)
		copy(buff[read:read+can], page[pageOff:pageOff+can])
		read += can
	}
	mc.filesMu.RUnlock()

	if read > 0 {
		mc.lru.Add(path, struct{}{})
	}
	mc.stats.ReadCount.Add(1)
	mc.stats.ReadBytes.Add(int64(read))
	return read, nil
}

// Delete removes the cached file at the given path from the cache.
func (mc *MemoryCache) Delete(path string) bool {
	// Try to remove via LRU so onEvict frees memory & updates stats
	if ok := mc.lru.Remove(path); ok {
		return true
	}

	mc.filesMu.Lock()
	defer mc.filesMu.Unlock()

	// If it wasn't in the LRU (edge case), free manually.
	if ent, ok := mc.files[path]; ok {
		var freed int64
		for _, pg := range ent.pages {
			freed += int64(len(pg))
		}
		delete(mc.files, path)
		mc.stats.SizeBytes.Add(-freed)
		mc.stats.FileCount.Add(-1)
		return true
	}
	return false
}

// Stats returns the current cache statistics.
func (mc *MemoryCache) Stats() *cache.Stats {
	s := &cache.Stats{
		CapacityBytes: mc.stats.CapacityBytes,
		MaxFileCount:  mc.stats.MaxFileCount,
		LoadDuration:  mc.stats.LoadDuration,
		LruCount:      mc.stats.LruCount,
	}

	s.SizeBytes.Store(mc.stats.SizeBytes.Load())
	s.FileCount.Store(mc.stats.FileCount.Load())
	s.ReadBytes.Store(mc.stats.ReadBytes.Load())
	s.ReadCount.Store(mc.stats.ReadCount.Load())
	s.WriteBytes.Store(mc.stats.WriteBytes.Load())
	s.WriteCount.Store(mc.stats.WriteCount.Load())

	s.CapacityBytesRemaining = s.CapacityBytes - s.SizeBytes.Load()
	s.FileCountRemaining = s.MaxFileCount - s.FileCount.Load()
	return s
}

// StartMaintenance starts the maintenance goroutine if it is not already running.
func (mc *MemoryCache) StartMaintenance() {
	mc.maintMu.Lock()
	if mc.maintActive {
		mc.maintMu.Unlock()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	mc.maintCancel = cancel
	mc.maintDone = make(chan struct{})
	mc.maintActive = true
	mc.maintMu.Unlock()

	go mc.maintenanceLoop(ctx)
}

// StopMaintenance stops the maintenance goroutine if it is running.
func (mc *MemoryCache) StopMaintenance() {
	mc.maintMu.Lock()
	if !mc.maintActive {
		mc.maintMu.Unlock()
		return
	}
	cancel := mc.maintCancel
	done := mc.maintDone

	// Clear struct state under the lock so a concurrent Start can proceed.
	mc.maintCancel = nil
	mc.maintDone = nil
	mc.maintActive = false
	mc.maintMu.Unlock()

	// Trigger shutdown and wait for the goroutine to exit cleanly.
	if cancel != nil {
		cancel()
	}
	if done != nil {
		<-done
	}
}

func (dc *MemoryCache) hasCapacityDelta(delta int64, newFile bool) bool {
	bytesOK := dc.Capacity == 0 || dc.stats.SizeBytes.Load()+delta <= dc.Capacity
	filesOK := dc.MaxFileCount == 0 || !newFile || dc.stats.FileCount.Load() < dc.MaxFileCount
	return bytesOK && filesOK
}

func bytesNeededForWrite(ent *entry, ofst, n int64) int64 {
	var need, pos int64
	for pos < n {
		abs := ofst + pos
		pageIdx := abs / pageSize
		pageOff := int(abs % pageSize)
		toCopy := int64(min(int(n-pos), pageSize-pageOff))
		if _, ok := ent.pages[pageIdx]; !ok {
			need += int64(pageSize)
		}
		pos += toCopy
	}
	return need
}

// maintenanceLoop runs periodic maintenance tasks until the context is cancelled.
func (mc *MemoryCache) maintenanceLoop(ctx context.Context) {
	ticker := time.NewTicker(mc.MaintenanceInterval)
	defer func() {
		ticker.Stop()
		mc.maintMu.Lock()
		// Nothing else to do here; Start/Stop already manage flags,
		// but this ensures done is always closed exactly once.
		done := mc.maintDone
		mc.maintMu.Unlock()

		if done != nil {
			close(done)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			mc.runMaintenanceOnce(ctx)
		}
	}
}

func (mc *MemoryCache) runMaintenanceOnce(_ context.Context) {
	mc.log.Debug("memoryCache: performing maintenance")
	start := time.Now()

	if mc.MaxAge == 0 {
		mc.log.Debug("memoryCache: maintenance complete, duration=%v", time.Since(start))
		return
	}

	// Snapshot keys under read lock, then remove outside the lock
	mc.filesMu.RLock()
	agedOut := make([]string, 0, len(mc.files))
	for path, ent := range mc.files {
		if time.Since(ent.mod) > mc.MaxAge {
			agedOut = append(agedOut, path)
		}
	}
	mc.filesMu.RUnlock()

	for _, path := range agedOut {
		mc.lru.Remove(path)
	}
	mc.log.Debug("memoryCache: maintenance complete, duration=%v", time.Since(start))
}

func (mc *MemoryCache) onEvict(path string, value struct{}) {
	mc.log.Debug("memoryCache: evicting %s", path)
	// look up entry in files map
	mc.filesMu.Lock()
	defer mc.filesMu.Unlock()
	ent, ok := mc.files[path]
	if !ok {
		mc.log.Debug("memoryCache: evicting %s but not found in files map", path)
		return
	}
	// calculate size to be freed
	var sizeFreed int64
	for _, page := range ent.pages {
		sizeFreed += int64(len(page))
	}
	// delete entry from files map
	delete(mc.files, path)
	// update stats
	mc.stats.SizeBytes.Add(-sizeFreed)
	mc.stats.FileCount.Add(-1)
}

func (mc *MemoryCache) validateOpts() error {
	// Hard errors: invalid limits
	if mc.Capacity < MinCapacity {
		return fmt.Errorf("memoryCache: capacity cannot be less than the minimum: %d, provided: %d", MinCapacity, mc.Capacity)
	}
	if mc.Capacity > MaxCapacity {
		return fmt.Errorf("memoryCache: capacity cannot be greater than the maximum: %d, provided: %d", MaxCapacity, mc.Capacity)
	}
	// Ensure logger
	if mc.log == nil {
		mc.log = &log.NoOpLogger{}
	}
	return nil
}
