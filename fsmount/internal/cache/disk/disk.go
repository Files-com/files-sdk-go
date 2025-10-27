// Package disk implements a disk-based cache for file data.
package disk

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/log"

	lru "github.com/hashicorp/golang-lru/v2"
)

const (
	// DefaultCacheCapacity is the maximum size of the local cache in bytes if not provided
	// at cache initialization. 0 means unbounded. Equivalent to 2GiB.
	DefaultCapacity = 2 * (1 << 30) // (1 << 30) == 1GiB

	// DefaultMaxAge is the maximum age of cache files, if not provided
	// at cache initialization. 0 means no expiration. Equivalent to 7 days.
	DefaultMaxAge = 7 * 24 * time.Hour

	// DefaultMaxFileCount is the maximum number of files in the cache if not provided
	// at cache initialization. 0 means unbounded.
	DefaultMaxFileCount = 10_000

	// DefaultMaintenanceInterval is the interval at which cache maintenance tasks are performed,
	// if not provided at cache initialization.
	DefaultMaintenanceInterval = 5 * time.Minute

	// DefaultLruFlushInterval is the interval at which the LRU state is persisted to disk,
	// if not provided at cache initialization.
	DefaultLruFlushInterval = 1 * time.Minute

	// unboundedFileCount is a large number used to represent an unbounded file count in the LRU cache.
	//
	// Techically the LRU cache is bounded, so this is a sufficiently large number to approximate
	// unbounded behavior. Equivalent to 4,294,967,295
	unboundedFileCount = 1 << 32

	// stateDir is the directory within the cache root where the lru.json file will be stored
	stateDir = "state"

	// dataDir is the directory within the cache root where the downloaded file data will be stored
	dataDir = "data"

	lruStateFile = "lru.json"
)

// DiskCache implements a simple disk-based cache for file data with an interface that roughly
// matches the FUSE Read/Write methods.
//
// TODO:
//   - allow concurrent reads/writes to different files in the cache by moving lock coordination
//     out of the caller and lock by path within the cache with a RWMutex map. This will require
//     maintaining a reservation of bytes so that concurrent goroutines can't exhaust memory that
//     a different goroutine evicted for its Write operation.
//   - consider limiting the max size of individual cache files so that a single large file
//     cannot evict the entire cache
//   - allow callers to distiguish between "fatal" vs "non fatal" errors with custom error types
//     e.g. if a read fails, the caller can just read from the source instead of treating it as a
//     fatal error but a write failure because the cache is at capacity and could not make space is
//     more serious
//   - figure out platform specific way to disable indexing of cache files e.g. on macOS
//     set the com.apple.metadata:com_apple_backup_excludeItem attribute
//     or use .noindex directory for cache root
//   - move Stats#MarshalJSON to a separate file to and apply the filescomfs_debug tag
type DiskCache struct {
	// CacheRoot is the root directory for the cache.
	CacheRoot string

	// Capacity is the maximum size of the cache in bytes. 0 means unbounded.
	Capacity int64

	// Disabled indicates whether the cache is disabled.
	//
	// A disabled cache:
	//  - performs no I/O and returns no data.
	//  - allows the caller to operate without caching, while not requiring conditional code around
	//  all cache operations.
	Disabled bool

	// MaxAge is the maximum age of cache files. 0 means no expiration.
	MaxAge time.Duration

	// MaxFileCount is the maximum number of files in the cache. 0 means unbounded.
	MaxFileCount int64

	// MaintenanceInterval is the interval at which cache maintenance tasks are performed.
	//
	// Maintenance tasks include:
	//   - Stats are reloaded from disk
	//   - Old files are deleted based on MaxAge
	MaintenanceInterval time.Duration

	// LruFlushInterval is the interval at which the LRU state is persisted to disk.
	LruFlushInterval time.Duration

	// current cache stats
	stats cache.Stats

	// protects from concurrent write operations
	writeMu sync.Mutex

	// protects from concurrent delete operations
	delMu sync.Mutex

	// used to log cache operations
	log log.Logger

	// LRU cache to track file access for eviction
	lru *lru.Cache[string, struct{}]

	maintMu     sync.Mutex
	maintActive bool
	maintCancel context.CancelFunc
	wg          sync.WaitGroup

	stateDir string
	dataDir  string

	lruDirty atomic.Bool
}

// NewDiskCache creates a DiskCache rooted at path and applies any options. The provided path must
// be an absolute path to a directory that already exists and is writable by the current process.
//
// If not disabled, it ensures the directory exists and initializes stats by scanning it.
//
// Defaults:
//   - Disabled: false
//   - Capacity: DefaultCapacity
//   - MaxAge: DefaultMaxAge
//   - MaxFileCount: DefaultMaxFileCount
func NewDiskCache(path string, opts ...Option) (*DiskCache, error) {
	// make sure the path is writable
	if err := validateCachePath(path); err != nil {
		return nil, err
	}

	// ensure data and state directories exist
	dataDir := filepath.Join(path, dataDir)
	stateDir := filepath.Join(path, stateDir)
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("diskCache: error creating data directory in cache root %s: %w", path, err)
	}
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		return nil, fmt.Errorf("diskCache: error creating state directory in cache root %s: %w", path, err)
	}

	dc := &DiskCache{
		CacheRoot:           path,
		dataDir:             dataDir,
		stateDir:            stateDir,
		Capacity:            DefaultCapacity,
		Disabled:            false,
		LruFlushInterval:    DefaultLruFlushInterval,
		MaintenanceInterval: DefaultMaintenanceInterval,
		MaxAge:              DefaultMaxAge,
		MaxFileCount:        DefaultMaxFileCount,
		log:                 nil,
		lru:                 nil,
	}

	// apply options
	for _, opt := range opts {
		opt(dc)
	}

	// If disabled, do not create directory or scan
	if dc.Disabled {
		return dc, nil
	}

	// validate the options before initializing LRU/scan
	if err := dc.validateOpts(); err != nil {
		return nil, err
	}

	// do this after applying options in case the caller set a different max file count
	// or provided their own LRU cache
	if dc.lru == nil {
		// protect against zero or negative max file count
		var trackCap int64
		if dc.MaxFileCount > 0 {
			trackCap = dc.MaxFileCount
		} else {
			// Unbounded file-count; use a large number since the LRU cache cannot be truly unbounded.
			trackCap = unboundedFileCount
		}
		lru, err := lru.NewWithEvict(int(trackCap), dc.onEvict)
		if err != nil {
			return nil, fmt.Errorf("diskCache: error creating LRU cache: %w", err)
		}
		dc.lru = lru
	}

	// restore LRU state from disk if present
	dc.restoreLRUState()

	if err := dc.loadStats(); err != nil {
		return nil, fmt.Errorf("diskCache: error initializing cache stats: %w", err)
	}

	return dc, nil
}

// Read reads data from the cached file at the given path into buff starting at the provided offset.
//
// It returns the number of bytes read, or 0 if the file is not in the cache.
func (dc *DiskCache) Read(path string, buff []byte, ofst int64) (n int, err error) {
	if dc.Disabled {
		return 0, nil
	}
	dc.stats.ReadCount.Add(1)
	fqPath := dc.entryPath(path)

	_, ok := dc.lru.Get(fqPath)
	if !ok {
		// file is not in the cache
		dc.log.Trace("diskCache: LRU does not contain path %s", path)
	}

	file, err := os.Open(fqPath)
	if err != nil {
		// this is not really an error - the file isn't cached yet
		if errors.Is(err, os.ErrNotExist) {
			return 0, nil
		}
		return 0, fmt.Errorf("diskCache: error opening cached file %s: %v", fqPath, err)
	}
	defer file.Close()

	// get the file info and delete if expired
	info, err := file.Stat()
	if err != nil {
		return 0, fmt.Errorf("diskCache: error stating cached file %s: %v", fqPath, err)
	}
	deleted, err := dc.deleteIfExpired(fqPath, info)
	if err != nil {
		return 0, fmt.Errorf("diskCache: error checking expiration for cached file %s: %v", fqPath, err)
	}
	if deleted {
		// file was expired and deleted
		return 0, nil
	}

	n, err = file.ReadAt(buff, ofst)
	if err != nil && !errors.Is(err, io.EOF) {
		return 0, fmt.Errorf("diskCache: error reading cached file %s at offset %d: %v", fqPath, ofst, err)
	}
	dc.stats.ReadBytes.Add(int64(n))

	// the file exists, but may not have been in the LRU, calling Add will ensure it is in the LRU
	// and is safe to call if the key was already present
	dc.lru.Add(fqPath, struct{}{})
	dc.lruDirty.Store(true)

	return n, nil
}

// Write writes data from buff to the cached file at the given path starting at offset ofst.
// Writing at an offset past the end of the file will grow the file and fill the gap with zeros.
//
// It returns the number of bytes written, or 0 if the cache is not enabled.
func (dc *DiskCache) Write(path string, buff []byte, ofst int64) (n int, err error) {
	if dc.Disabled {
		return 0, nil
	}
	dc.writeMu.Lock()
	defer dc.writeMu.Unlock()
	dc.stats.WriteCount.Add(1)

	fqPath := dc.entryPath(path)
	if err := os.MkdirAll(filepath.Dir(fqPath), 0o755); err != nil {
		return 0, fmt.Errorf("diskCache: error creating directories for cached file %s: %v", fqPath, err)
	}

	st, err := os.Stat(fqPath)
	newFile := errors.Is(err, os.ErrNotExist)
	var currSize int64
	if err == nil && !st.IsDir() {
		currSize = st.Size()
	}
	projected := max(ofst+int64(len(buff)), currSize)
	delta := projected - currSize

	if dc.Capacity > 0 || dc.MaxFileCount > 0 {
		// evict the least recently used files until there is enough space, or the file count is under the limit
		for !dc.hasCapacityDelta(delta, newFile) {
			if _, _, ok := dc.lru.RemoveOldest(); ok {
				continue
			}
			// LRU is empty → fall back to disk scan deleting ONE file
			if !dc.deleteOneByMtime() {
				return 0, fmt.Errorf("diskCache: at capacity, no entries to evict, and nothing deletable")
			}
		}
	}

	file, err := os.OpenFile(fqPath, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return 0, fmt.Errorf("diskCache: error opening cached file %s: %v", fqPath, err)
	}
	defer file.Close()

	// WriteAt will write at the specified offset, even if it's past the end of the file,
	// and grow the file to ofst + n bytes, filling the gap with zeros.
	n, err = file.WriteAt(buff, ofst)
	if err != nil {
		return n, fmt.Errorf("diskCache: error writing cached file %s at offset %d: %v", fqPath, ofst, err)
	}
	dc.stats.WriteBytes.Add(int64(n))
	if newFile && projected > 0 {
		dc.stats.FileCount.Add(1)
	}
	dc.stats.SizeBytes.Add(delta)

	// the file exists, but may not have been in the LRU, calling Add will ensure it is in the LRU
	// and is safe to call if the key was already present
	dc.lru.Add(fqPath, struct{}{})
	dc.lruDirty.Store(true)

	return n, nil
}

// Delete removes the cached file from the cache. It returns true if the file was deleted.
func (dc *DiskCache) Delete(path string) bool {
	if dc.Disabled {
		return false
	}

	fqPath := dc.entryPath(path)
	deleted := dc.lru.Remove(fqPath)
	return deleted
}

// StartMaintenance starts the maintenance goroutine if it is not already running.
func (dc *DiskCache) StartMaintenance() {
	dc.maintMu.Lock()
	if dc.maintActive {
		dc.maintMu.Unlock()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	dc.maintCancel = cancel
	dc.maintActive = true
	dc.maintMu.Unlock()
	dc.wg.Add(1)
	go func() {
		defer dc.wg.Done()
		dc.maintenanceLoop(ctx)
	}()
}

// StopMaintenance stops the maintenance goroutine if it is running.
func (dc *DiskCache) StopMaintenance() {
	dc.maintMu.Lock()
	if !dc.maintActive {
		dc.maintMu.Unlock()
		return
	}
	cancel := dc.maintCancel
	dc.maintMu.Unlock()
	// Trigger shutdown and wait for the goroutine to exit cleanly.
	if cancel != nil {
		cancel()
	}
	dc.wg.Wait()

	dc.maintMu.Lock()
	// Clear struct state under the lock so a concurrent Start can proceed.
	dc.maintCancel = nil
	dc.maintActive = false
	dc.maintMu.Unlock()

	dc.persistLRUState()
}

// Stats returns the current cache statistics.
func (dc *DiskCache) Stats() *cache.Stats {
	s := &cache.Stats{
		CapacityBytes: dc.stats.CapacityBytes,
		MaxFileCount:  dc.stats.MaxFileCount,
		LoadDuration:  dc.stats.LoadDuration,
		LruCount:      dc.stats.LruCount,
	}

	s.SizeBytes.Store(dc.stats.SizeBytes.Load())
	s.FileCount.Store(dc.stats.FileCount.Load())
	s.ReadBytes.Store(dc.stats.ReadBytes.Load())
	s.ReadCount.Store(dc.stats.ReadCount.Load())
	s.WriteBytes.Store(dc.stats.WriteBytes.Load())
	s.WriteCount.Store(dc.stats.WriteCount.Load())

	s.CapacityBytesRemaining = s.CapacityBytes - s.SizeBytes.Load()
	s.FileCountRemaining = s.MaxFileCount - s.FileCount.Load()
	return s
}

// onEvict is called when a key is evicted from the LRU. This includes when the key is removed
// explicitly via Remove or RemoveOldest.
func (dc *DiskCache) onEvict(path string, value struct{}) {
	if err := dc.deleteFile(path); err != nil {
		dc.log.Info("diskCache: error deleting evicted cached file %s: %v", path, err)
	}
	dc.lruDirty.Store(true)
}

func (dc *DiskCache) hasCapacityDelta(delta int64, newFile bool) bool {
	bytesOK := dc.Capacity == 0 || dc.stats.SizeBytes.Load()+delta <= dc.Capacity
	filesOK := dc.MaxFileCount == 0 || !newFile || dc.stats.FileCount.Load() < dc.MaxFileCount
	return bytesOK && filesOK
}

func (dc *DiskCache) deleteFile(path string) error {
	dc.delMu.Lock()
	defer dc.delMu.Unlock()
	fqPath := dc.entryPath(path)
	// os.RemoveAll does not return an error if the path does not exist
	// so stat the file first to avoid updating the stats incorrectly
	st, err := os.Stat(fqPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// no-op - file already deleted
			return nil
		}
		return err
	}
	if st.IsDir() {
		return nil
	}

	if err := os.RemoveAll(fqPath); err != nil {
		dc.log.Info("diskCache: error deleting evicted cached file %s: %v", fqPath, err)
		return err
	}

	dc.stats.FileCount.Add(-1)
	dc.stats.SizeBytes.Add(-st.Size())
	return nil
}

func (dc *DiskCache) loadStats() error {
	start := time.Now()

	dc.writeMu.Lock()
	dc.delMu.Lock()
	defer func() {
		dc.writeMu.Unlock()
		dc.delMu.Unlock()
	}()

	// initialize stats by scanning files in the cache directory
	var totalSize int64
	var fileCount int64
	err := filepath.Walk(dc.dataDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			fileCount++
			totalSize += info.Size()
			if !dc.lru.Contains(p) {
				// TODO: decide what to do in this case - for now just log it
				dc.log.Debug("diskCache: loadStats: LRU does not contain %s", p)
			}
		}
		return nil
	})
	if err != nil {
		dc.log.Debug("diskCache: loadStats: failed to load stats: %v", err)
		return err
	}
	dc.stats.CapacityBytes = dc.Capacity
	dc.stats.CapacityBytesRemaining = dc.Capacity - totalSize
	dc.stats.MaxFileCount = dc.MaxFileCount
	dc.stats.FileCountRemaining = dc.MaxFileCount - fileCount
	dc.stats.FileCount.Store(fileCount)
	dc.stats.SizeBytes.Store(totalSize)
	dc.stats.LoadDuration = time.Since(start)
	dc.stats.LruCount = len(dc.lru.Keys())
	return nil
}

// validateCachePath checks that the provided cache path meets the requirements:
//   - is not empty
//   - must be an absolute path
//   - must be an existing directory
//   - must be writable by the current process
func validateCachePath(path string) error {
	if path == "" {
		return errors.New("diskCache: cache root path cannot be empty")
	}
	if !filepath.IsAbs(path) {
		return fmt.Errorf("diskCache: cache root path must be absolute: %s", path)
	}
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("diskCache: error stating cache root path %s: %w", path, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("diskCache: cache root path is not a directory: %s", path)
	}
	testFile := filepath.Join(path, ".cache_write_test")
	f, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("diskCache: cache root path is not writable: %s", path)
	}
	f.Close()
	return os.Remove(testFile)
}

func (dc *DiskCache) validateOpts() error {
	// Hard errors: invalid limits
	if dc.Capacity < 0 {
		return fmt.Errorf("diskCache: capacity cannot be negative: %d", dc.Capacity)
	}
	if dc.MaxFileCount < 0 {
		return fmt.Errorf("diskCache: max file count cannot be negative: %d", dc.MaxFileCount)
	}
	if dc.MaxAge < 0 {
		return fmt.Errorf("diskCache: max age cannot be negative: %v", dc.MaxAge)
	}

	if dc.MaintenanceInterval <= 0 {
		return fmt.Errorf("diskCache: maintenance interval cannot be negative: %v", dc.MaxAge)
	}
	if dc.LruFlushInterval <= 0 {
		return fmt.Errorf("diskCache: LRU flush interval cannot be negative: %v", dc.MaxAge)
	}

	// Ensure logger
	if dc.log == nil {
		dc.log = &log.NoOpLogger{}
	}
	return nil
}

// adjust the length of the shard prefix to change the number of subdirectories used to store cached files.
// the math is something like:
//
// files per directory == (total number of cache files / 16 ^ shardPrefixLen)
//
// e.g.
//   - 100,000 files / 16^2 == 390 files per directory
//   - 100,000 files / 16^3 == 24 files per directory
//   - 1,000,000 files / 16^3 == 244 files per directory
const shardPrefixLen = 2

// entryPath returns the full file path for the cached file based on a hash of the original path.
func (dc *DiskCache) entryPath(path string) string {
	if path != dc.dataDir && strings.HasPrefix(path, dc.dataDir) {
		// path is already a full path in the cache
		return path
	}
	sum := sha256.Sum256([]byte(path))
	h := hex.EncodeToString(sum[:])
	n := min(shardPrefixLen, len(h))
	dir := h[:n]
	name := h[n:] + "-" + filepath.Base(path)
	return filepath.Join(dc.dataDir, dir, name)
}

type lruState struct {
	Keys []string `json:"keys"`
}

// restoreLRUState loads any existing LRU state stored on disk. it is expected to be
// called after initializing dc.lru in NewDiskCache
func (dc *DiskCache) restoreLRUState() {
	_ = dc.loadLRUState()
}

func (dc *DiskCache) loadLRUState() error {
	if dc.lru == nil {
		return nil
	}

	statePath := filepath.Join(dc.stateDir, lruStateFile)
	f, err := os.Open(statePath)
	if err != nil {
		// a non-existent state file is not an error
		return nil
	}
	defer f.Close()
	var state lruState
	dec := json.NewDecoder(f)
	if err := dec.Decode(&state); err != nil {
		return err
	}
	// add keys in order to LRU
	for _, k := range state.Keys {
		dc.lru.Add(k, struct{}{})
		dc.lruDirty.Store(true)
	}

	dc.writeMu.Lock()
	defer dc.writeMu.Unlock()
	dc.stats.LruCount = len(state.Keys)
	return nil
}

func (dc *DiskCache) saveLRUState() error {
	if dc.lru == nil {
		return nil
	}
	// Keys returns a slice of the keys in the LRU, from oldest to newest.
	// when restoring, we add them in order to preserve LRU order.
	keys := dc.lru.Keys()
	state := lruState{Keys: keys}
	statePath := filepath.Join(dc.stateDir, lruStateFile)
	f, err := os.Create(statePath)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	return enc.Encode(state)
}

// persistLRUState saves the current LRU state to disk. it is called periodically by the maintenance loop.
func (dc *DiskCache) persistLRUState() {
	if !dc.lruDirty.Load() {
		return
	}
	if err := dc.saveLRUState(); err != nil {
		dc.log.Info("diskCache: error saving LRU state: %v", err)
	}
	dc.lruDirty.Store(false)
}

// maintenanceLoop runs periodic maintenance tasks until the context is cancelled.
func (dc *DiskCache) maintenanceLoop(ctx context.Context) {
	ticker := time.NewTicker(dc.MaintenanceInterval)
	persistTicker := time.NewTicker(dc.LruFlushInterval)
	defer ticker.Stop()
	defer persistTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			dc.runMaintenanceOnce(ctx)
		case <-persistTicker.C:
			dc.persistLRUState()
		}
	}
}

type fileMeta struct {
	path string
	size int64
	mod  time.Time
}

func (dc *DiskCache) runMaintenanceOnce(ctx context.Context) {
	if dc.Disabled {
		return
	}

	var agedOut []fileMeta
	var retainedSize, retainedCount int64
	filesOnDisk := make(map[string]struct{}, 1024)
	var notInLru []string

	// create a snapshot of files to delete along with the predicted size and count of files
	// that will be left after deletions
	dc.writeMu.Lock()
	dc.log.Debug("diskCache: performing maintenance")
	start := time.Now()
	_ = filepath.WalkDir(dc.dataDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		info, ierr := d.Info()
		if ierr != nil {
			// the maintenance is best-effort, so just log and continue
			dc.log.Info("diskCache: maintenance: error stating file %s: %v", p, ierr)
			return nil
		}

		filesOnDisk[p] = struct{}{}

		expired := dc.MaxAge > 0 && time.Since(info.ModTime()) > dc.MaxAge
		if expired {
			agedOut = append(agedOut, fileMeta{path: p, size: info.Size(), mod: info.ModTime()})
			return nil
		}
		retainedCount++
		retainedSize += info.Size()

		if !dc.lru.Contains(p) {
			// these will be added later outside the lock
			notInLru = append(notInLru, p)
		}
		return nil
	})
	dc.writeMu.Unlock()

	// apply deletions, no need to hold any locks as deleteFile holds delMu
	for _, old := range agedOut {
		// prefer LRU so onEvict -> deleteFile runs.
		if removed := dc.lru.Remove(old.path); !removed {
			// not in LRU: delete directly, and ignore errors because
			// maintenance is best-effort
			_ = dc.deleteFile(old.path)
		}
		// check context for responsiveness
		select {
		case <-ctx.Done():
			return
		default:
		}
	}

	// drop LRU entries for files that are no longer on disk
	for _, k := range dc.lru.Keys() {
		if _, ok := filesOnDisk[k]; ok {
			continue
		}
		// double-check that the file is really gone for robustness
		if _, err := os.Stat(k); err != nil {
			// most likely os.ErrNotExist, remove the entry from the LRU
			// the onEvict callback is a no-op if the file is already gone
			dc.lru.Remove(k)
		}
		select {
		case <-ctx.Done():
			return
		default:
		}
	}

	// add anything that was on disk but not in the LRU
	// this technically could change the LRU order, but since maintenance is
	// best-effort, it's acceptable
	for _, p := range notInLru {
		dc.lru.Add(p, struct{}{})
		select {
		case <-ctx.Done():
			return
		default:
		}
	}

	duration := time.Since(start)
	// commit authoritative stats
	dc.writeMu.Lock()
	dc.stats.CapacityBytes = dc.Capacity
	dc.stats.MaxFileCount = dc.MaxFileCount
	dc.stats.SizeBytes.Store(retainedSize)
	dc.stats.FileCount.Store(retainedCount)
	dc.stats.CapacityBytesRemaining = dc.Capacity - retainedSize
	dc.stats.FileCountRemaining = dc.MaxFileCount - retainedCount
	dc.stats.LoadDuration = duration
	dc.stats.LruCount = len(dc.lru.Keys())
	dc.writeMu.Unlock()

	dc.log.Debug("diskCache: maintenance complete, duration=%vms", duration.Milliseconds())
}

// delete if expired based on MaxAge
func (dc *DiskCache) deleteIfExpired(path string, info os.FileInfo) (deleted bool, err error) {
	if info.IsDir() {
		return false, nil
	}
	fqPath := dc.entryPath(path)
	if dc.MaxAge > 0 {
		age := time.Since(info.ModTime())
		if age > dc.MaxAge {
			if removed := dc.lru.Remove(fqPath); removed {
				// deleted via LRU evict callback
				return removed, nil
			}
			// not in LRU for some reason, just delete the file directly
			if err := dc.deleteFile(fqPath); err != nil {
				return false, err
			}
			return true, nil
		}
	}
	return false, nil
}

// deleteOneByMtime finds the oldest file and deletes just that one.
// Returns true if it deleted something.
func (dc *DiskCache) deleteOneByMtime() bool {
	type fi struct {
		path string
		t    time.Time
	}
	var oldest *fi
	_ = filepath.Walk(dc.dataDir, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if oldest == nil || info.ModTime().Before(oldest.t) {
			oldest = &fi{path: p, t: info.ModTime()}
		}
		return nil
	})
	if oldest == nil {
		return false
	}
	// Prefer going through LRU if present to trigger onEvict; otherwise delete directly.
	if removed := dc.lru.Remove(oldest.path); removed {
		return true
	}
	if err := dc.deleteFile(oldest.path); err != nil {
		dc.log.Info("diskCache: error deleting oldest by mtime %s: %v", oldest.path, err)
		return false
	}
	return true
}
