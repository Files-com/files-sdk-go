//go:build linux || windows

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
	fssync "github.com/Files-com/files-sdk-go/v3/fsmount/internal/sync"
	"github.com/Files-com/files-sdk-go/v3/lib/privatefile"

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

	// DefaultRangeFlushBytes is how many newly written bytes an entry may hold
	// before they are synced and recorded as durable. Downloads also flush when
	// they stop, so this only bounds the progress a crash can lose.
	DefaultRangeFlushBytes = 256 * (1 << 20)

	// unboundedFileCount is a large number used to represent an unbounded file count in the LRU cache.
	//
	// Techically the LRU cache is bounded, so this is a sufficiently large number to approximate
	// unbounded behavior. Equivalent to 4,294,967,295
	unboundedFileCount = 1 << 32

	// stateDir is the directory within the cache root where the lru.json file will be stored
	stateDir = "state"

	// dataDir is the directory within the cache root where the downloaded file data will be stored
	dataDir = "data"

	// partialDir is the directory within the cache root where in-progress file data will be stored
	partialDir = "partial"

	// metadataDir is the directory within the cache root where completed cache entry metadata is stored
	metadataDir = "metadata"

	lruStateFile = "lru.json"

	// maxEvictionAttempts is the maximum number of consecutive pinned files
	// checked during eviction before concluding that all files are pinned.
	maxEvictionAttempts = 10
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
	//   - Unpinned data files without completed metadata are deleted
	MaintenanceInterval time.Duration

	// LruFlushInterval is the interval at which the LRU state is persisted to disk.
	LruFlushInterval time.Duration

	// RangeFlushBytes is how many newly written bytes an entry may hold before
	// they are synced and recorded as durable.
	RangeFlushBytes int64

	// current cache stats
	stats cache.Stats

	// accountedBytes stores the valid byte count used by cache capacity and
	// eviction accounting. Sparse file length cannot be used for this purpose.
	accountedBytes   map[string]int64
	accountedBytesMu sync.Mutex

	// protects from concurrent read and write operations.
	// Write acquires the exclusive lock; Read acquires the shared (read) lock.
	// This prevents concurrent reads from observing partially-written file data.
	writeMu sync.RWMutex

	// protects from concurrent delete operations
	delMu sync.Mutex

	// scanDone closes when the startup scan has accounted existing entries.
	scanDone chan struct{}

	// rangeWrites makes download progress readable immediately while
	// persistent range metadata advances only after each durable batch.
	rangeWriteLocks *fssync.PathMutex
	rangeWritesMu   sync.Mutex
	rangeWrites     map[string]*rangeWriteState

	// dataFiles holds the data files of entries being downloaded open.
	dataFilesMu sync.Mutex
	dataFiles   map[string]*dataFile

	// used to log cache operations
	log log.Logger

	// LRU cache to track file access for eviction
	lru *lru.Cache[string, struct{}]

	// Track pinned files (reference counted) to prevent eviction of files with open handles
	pinnedFiles   map[string]int
	pinnedFilesMu sync.Mutex
	clearPending  map[string]struct{}

	maintenanceRunMu sync.Mutex
	maintMu          sync.Mutex
	maintActive      bool
	maintCancel      context.CancelFunc
	maintDone        chan struct{}

	stateDir string
	dataDir  string
	partDir  string
	metaDir  string

	lruDirty atomic.Bool
}

// NewDiskCache creates a DiskCache rooted at path and applies any options. The provided path must
// be an absolute path to a directory that already exists and is writable by the current process.
//
// If not disabled, it ensures the directory exists and initializes stats by scanning it.
//
// The cache root itself belongs to the caller and is left as it is. Everything the cache stores
// below it (file data, partial data, metadata and LRU state) is readable only by the current user:
// the data, partial and state directories the cache owns are created private, and an existing
// cache is made private before any of its state is read, so a cache written by an earlier version
// that created world-readable entries is protected as well. See lib/privatefile for what private
// means on each platform. A root that another user could replace entries in is refused.
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
	if err := privatefile.CheckContainer(path); err != nil {
		return nil, fmt.Errorf("diskCache: cache root path cannot hold private data: %w", err)
	}

	// ensure the private data and state directories exist
	dataDir := filepath.Join(path, dataDir)
	partDir := filepath.Join(path, partialDir)
	stateDir := filepath.Join(path, stateDir)
	metaDir := filepath.Join(stateDir, metadataDir)
	for _, dir := range []string{dataDir, partDir, stateDir} {
		// Every existing entry is checked before the directory's access is
		// changed: on Windows, changing a directory's access control list
		// rewrites what its descendants inherit, which would reach a file that
		// also has a name outside the cache.
		if err := checkTree(dir); err != nil {
			return nil, fmt.Errorf("diskCache: existing cache state in %s cannot be made private: %w", path, err)
		}
		if err := privatefile.EnsureDir(dir); err != nil {
			return nil, fmt.Errorf("diskCache: error creating private directory in cache root %s: %w", path, err)
		}
		if err := makeTreePrivate(dir); err != nil {
			return nil, fmt.Errorf("diskCache: error protecting existing cache state in %s: %w", path, err)
		}
	}
	if err := os.MkdirAll(metaDir, privateDirMode); err != nil {
		return nil, fmt.Errorf("diskCache: error creating metadata directory in cache root %s: %w", path, err)
	}

	dc := &DiskCache{
		CacheRoot:           path,
		dataDir:             dataDir,
		partDir:             partDir,
		stateDir:            stateDir,
		metaDir:             metaDir,
		Capacity:            DefaultCapacity,
		Disabled:            false,
		LruFlushInterval:    DefaultLruFlushInterval,
		RangeFlushBytes:     DefaultRangeFlushBytes,
		MaintenanceInterval: DefaultMaintenanceInterval,
		MaxAge:              DefaultMaxAge,
		MaxFileCount:        DefaultMaxFileCount,
		log:                 nil,
		lru:                 nil,
		pinnedFiles:         make(map[string]int),
		clearPending:        make(map[string]struct{}),
		accountedBytes:      make(map[string]int64),
		rangeWrites:         make(map[string]*rangeWriteState),
		rangeWriteLocks:     fssync.NewPathMutex(),
		dataFiles:           make(map[string]*dataFile),
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

	if err := dc.removeAbandonedPartials(); err != nil {
		return nil, fmt.Errorf("diskCache: error removing abandoned partial data: %w", err)
	}
	dc.stats.CapacityBytes = dc.Capacity
	dc.stats.MaxFileCount = dc.MaxFileCount
	// List entries now, so the background scan only judges data a previous
	// process left, never files this process creates while it runs.
	paths, err := dc.listEntries()
	if err != nil {
		return nil, fmt.Errorf("diskCache: error listing cached entries: %w", err)
	}
	dc.scanDone = make(chan struct{})
	go dc.loadStats(paths)

	return dc, nil
}

// Read reads data from the cached file at the given path into buff starting at the provided offset.
//
// It returns the number of bytes read, or 0 if the file is not in the cache.
func (dc *DiskCache) Read(path string, buff []byte, ofst int64) (n int, err error) {
	if dc.Disabled {
		return 0, nil
	}
	// Acquire the shared read lock to prevent reading partially-written data
	// from a concurrent Write. Without this, a Read at an offset within a chunk
	// being written by WriteAt can observe zeros (OS-extended but not yet written
	// pages), which the early cache check in remotefs.go treats as valid data.
	dc.writeMu.RLock()
	defer dc.writeMu.RUnlock()
	return dc.read(path, buff, ofst)
}

// read reads a cache entry while the caller holds writeMu for reading.
func (dc *DiskCache) read(path string, buff []byte, ofst int64) (n int, err error) {
	dc.stats.ReadCount.Add(1)
	fqPath := dc.entryPath(path)

	_, ok := dc.lru.Get(fqPath)
	if !ok {
		// file is not in the cache
		dc.log.Trace("DiskCache: LRU does not contain path %s", path)
	}

	var file *os.File
	if df := dc.sharedDataFile(fqPath); df != nil {
		// The entry is being downloaded, so it cannot be expired.
		defer dc.releaseDataFile(df)
		file = df.file
	} else {
		file, err = os.Open(fqPath)
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

func (dc *DiskCache) ReadComplete(path string, meta cache.EntryMetadata, buff []byte, ofst int64) (n int, err error) {
	if dc.Disabled {
		return 0, nil
	}
	dc.writeMu.RLock()
	defer dc.writeMu.RUnlock()

	fqPath := dc.entryPath(path)
	stored, found, err := dc.rangeMetadata(path, meta)
	if err != nil || !found {
		return 0, err
	}
	if !stored.IsComplete() {
		return 0, nil
	}
	info, err := os.Stat(fqPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			_ = dc.deleteMetadata(path)
			return 0, nil
		}
		return 0, err
	}
	if info.Size() != meta.Size {
		_ = dc.Delete(path)
		return 0, nil
	}
	return dc.read(path, buff, ofst)
}

func (dc *DiskCache) ReadPartial(path string, buff []byte, ofst int64) (n int, err error) {
	return dc.Read(dc.partialEntryPath(path), buff, ofst)
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
	dc.deleteRangeWritePath(path)
	dc.stats.WriteCount.Add(1)

	fqPath := dc.entryPath(path)
	if err := os.MkdirAll(filepath.Dir(fqPath), privateDirMode); err != nil {
		return 0, fmt.Errorf("diskCache: error creating directories for cached file %s: %v", fqPath, err)
	}

	st, err := os.Stat(fqPath)
	var currSize int64
	if err == nil && !st.IsDir() {
		currSize = st.Size()
	}
	projected := max(ofst+int64(len(buff)), currSize)
	accounted, tracked := dc.accountedSize(fqPath, currSize)
	delta := projected - accounted
	newFile := !tracked

	// Enforce capacity limits when configured by evicting old unpinned files.
	// The limits are "soft", so writes should always succeed, but a best effort is made to evict
	// unpinned files to stay near the configured limits.
	if dc.Capacity > 0 || dc.MaxFileCount > 0 {
		evictionAttempts := 0

		// Try to evict unpinned files to make room, but don't fail if eviction is not possible.
		for !dc.hasCapacityDelta(delta, newFile) {
			oldestKey, _, ok := dc.lru.GetOldest()
			if !ok {
				// LRU is empty - try disk scan as fallback
				dc.deleteOneByMtime()
				// Whether deletion succeeded or not, allow the write to proceed
				break
			}

			// Don't evict the file currently being written to
			if oldestKey == fqPath {
				break
			}

			// Check if the oldest file is pinned
			if dc.isPinned(oldestKey) {
				evictionAttempts++
				if evictionAttempts >= maxEvictionAttempts {
					// Likely all files are pinned. Allow the write to proceed.
					break
				}
				// Try to remove it anyway - onEvict will re-add it and provide the next oldest
				dc.lru.Remove(oldestKey)
				continue
			}

			// Remove the oldest entry (which is not the current file and is not pinned)
			evictionAttempts = 0 // Reset counter since an evictable file was found
			dc.lru.Remove(oldestKey)
		}
	}

	file, err := os.OpenFile(fqPath, os.O_CREATE|os.O_WRONLY, privateFileMode)
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
	dc.setAccountedSize(fqPath, projected)

	// the file exists, but may not have been in the LRU, calling Add will ensure it is in the LRU
	// and is safe to call if the key was already present
	dc.lru.Add(fqPath, struct{}{})
	dc.lruDirty.Store(true)

	return n, nil
}

func (dc *DiskCache) WritePartial(path string, buff []byte, ofst int64) (n int, err error) {
	return dc.Write(dc.partialEntryPath(path), buff, ofst)
}

func (dc *DiskCache) Commit(path string, meta cache.EntryMetadata) error {
	if dc.Disabled {
		return nil
	}
	meta.Path = path
	var err error
	meta, err = meta.WithCompleteRange()
	if err != nil {
		return err
	}

	// The data must be durable before metadata that calls it complete, or a
	// crash could leave a complete entry whose bytes never reached the disk.
	// Both syncs run before writeMu, which every cache read also waits on.
	fqPath := dc.entryPath(path)
	if err := syncFile(fqPath); err != nil {
		return err
	}
	tmpPath, err := dc.prepareEntryMetadata(meta)
	if err != nil {
		return err
	}
	defer os.Remove(tmpPath)

	dc.writeMu.Lock()
	defer dc.writeMu.Unlock()
	dc.deleteRangeWritePath(path)

	info, err := os.Stat(fqPath)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("diskCache: cannot commit directory cache entry %s", fqPath)
	}
	if info.Size() < meta.Size {
		return fmt.Errorf("diskCache: cache entry size mismatch for %s: got %d, want %d", path, info.Size(), meta.Size)
	}
	if info.Size() > meta.Size {
		if err := os.Truncate(fqPath, meta.Size); err != nil {
			return fmt.Errorf("diskCache: truncating cache entry %s to %d failed: %w", path, meta.Size, err)
		}
	}
	dc.setAccountedSize(fqPath, meta.Size)
	if err := os.Rename(tmpPath, dc.metadataPath(path)); err != nil {
		return err
	}
	// A clear that found this entry pinned deletes it on the last Unpin. An
	// entry committed after the clear replaces that one.
	dc.pinnedFilesMu.Lock()
	delete(dc.clearPending, fqPath)
	dc.pinnedFilesMu.Unlock()
	return nil
}

// syncFile flushes a file's written data to disk.
func syncFile(path string) error {
	// Windows only flushes a file opened for writing.
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	syncErr := file.Sync()
	closeErr := file.Close()
	if syncErr != nil {
		return syncErr
	}
	return closeErr
}

// Prepare a unique metadata file so a retired publication cannot overwrite a
// replacement's temporary file while its final rename waits for writeMu.
func (dc *DiskCache) prepareEntryMetadata(meta cache.EntryMetadata) (string, error) {
	if err := meta.Validate(); err != nil {
		return "", err
	}
	metaPath := dc.metadataPath(meta.Path)
	if err := os.MkdirAll(filepath.Dir(metaPath), privateDirMode); err != nil {
		return "", err
	}
	file, err := os.CreateTemp(filepath.Dir(metaPath), ".range-metadata-*")
	if err != nil {
		return "", err
	}
	tmpPath := file.Name()
	if err = json.NewEncoder(file).Encode(meta); err == nil {
		err = file.Sync()
	}
	closeErr := file.Close()
	if err == nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(tmpPath)
		return "", err
	}
	return tmpPath, nil
}

// Delete removes the cached file from the cache. It returns true if the file was deleted.
func (dc *DiskCache) Delete(path string) bool {
	if dc.Disabled {
		return false
	}
	dc.deleteRangeWritePath(path)

	fqPath := dc.entryPath(path)
	// Invalidate completed-version metadata even when an open handle pins the
	// data file. The pinned bytes remain readable through Read, but cannot seed
	// a later ReadComplete for a different remote version.
	_ = dc.deleteMetadata(path)
	if dc.isPinned(fqPath) {
		return false
	}
	_, statErr := os.Stat(fqPath)
	exists := statErr == nil
	deleted := dc.lru.Remove(fqPath)
	if !deleted && exists {
		if err := dc.deleteFile(fqPath); err == nil {
			deleted = true
		}
	}
	return deleted
}

func (dc *DiskCache) DeletePartial(path string) bool {
	return dc.Delete(dc.partialEntryPath(path))
}

// Clear removes all unpinned file data from the cache. Pinned entries are
// invalidated immediately and removed when their final open handle closes.
func (dc *DiskCache) Clear() error {
	invalidated, err := dc.Invalidate()
	if err != nil {
		return err
	}
	return invalidated.Remove()
}

// InvalidatedEntries is the data Invalidate made unreadable, which Remove
// deletes.
type InvalidatedEntries struct {
	dc         *DiskCache
	oldMetaDir string
}

// Invalidate makes every cached entry unreadable in one short exclusive step:
// download progress is forgotten and the metadata directory is replaced with
// an empty one. Reads and writes can continue right away, and Remove then
// deletes the invalidated data without blocking them.
func (dc *DiskCache) Invalidate() (*InvalidatedEntries, error) {
	if dc.Disabled {
		return &InvalidatedEntries{}, nil
	}
	dc.writeMu.Lock()
	dc.rangeWritesMu.Lock()
	clear(dc.rangeWrites)
	dc.rangeWritesMu.Unlock()
	oldMetaDir, err := dc.rotateMetadataDirectoryLocked()
	dc.writeMu.Unlock()
	// Closing waits for the files' current users, so it runs without writeMu.
	dc.closeAllDataFiles()
	if err != nil {
		return nil, err
	}
	return &InvalidatedEntries{dc: dc, oldMetaDir: oldMetaDir}, nil
}

// Remove deletes the data files Invalidate made unreadable. A file with an
// open handle is deleted when its last handle closes, and a file written
// since the invalidation is kept.
func (invalidated *InvalidatedEntries) Remove() error {
	dc := invalidated.dc
	if dc == nil {
		return nil
	}
	// Do not let maintenance operate on a stale snapshot while clear removes
	// entries. Normal reads and writes remain available throughout the walk.
	dc.maintenanceRunMu.Lock()
	defer dc.maintenanceRunMu.Unlock()

	var clearErrors []error
	var paths []string
	for _, root := range []string{dc.dataDir, dc.partDir} {
		err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
			if errors.Is(walkErr, os.ErrNotExist) {
				return nil
			}
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() {
				return nil
			}
			paths = append(paths, path)
			return nil
		})
		if err != nil {
			clearErrors = append(clearErrors, err)
		}
	}

	for _, path := range paths {
		dc.writeMu.Lock()
		if dc.writtenSinceInvalidationLocked(path) {
			dc.writeMu.Unlock()
			continue
		}
		dc.pinnedFilesMu.Lock()
		pinned := dc.pinnedFiles[path] > 0
		if pinned {
			dc.clearPending[path] = struct{}{}
		}
		dc.pinnedFilesMu.Unlock()
		if pinned {
			dc.writeMu.Unlock()
			continue
		}
		if err := dc.deleteFile(path); err != nil {
			clearErrors = append(clearErrors, err)
		} else {
			dc.lru.Remove(path)
		}
		dc.writeMu.Unlock()
	}

	if invalidated.oldMetaDir != "" {
		if err := os.RemoveAll(invalidated.oldMetaDir); err != nil {
			clearErrors = append(clearErrors, err)
		}
	}
	dc.lruDirty.Store(true)
	return errors.Join(clearErrors...)
}

// writtenSinceInvalidationLocked reports whether the entry at fqPath has
// download progress or metadata, which only data written after the
// invalidation can have. It runs under writeMu.
func (dc *DiskCache) writtenSinceInvalidationLocked(fqPath string) bool {
	dc.rangeWritesMu.Lock()
	_, downloading := dc.rangeWrites[fqPath]
	dc.rangeWritesMu.Unlock()
	if downloading {
		return true
	}
	metaPath, ok := dc.metadataPathForEntryPath(fqPath)
	if !ok {
		return false
	}
	_, err := os.Stat(metaPath)
	return err == nil
}

// rotateMetadataDirectoryLocked invalidates every completed cache entry with
// one short operation under writeMu. The old metadata tree can then be removed
// without blocking cache reads or writes.
func (dc *DiskCache) rotateMetadataDirectoryLocked() (string, error) {
	oldMetaDir := fmt.Sprintf("%s.clear-%d", dc.metaDir, time.Now().UnixNano())
	if err := os.Rename(dc.metaDir, oldMetaDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", os.MkdirAll(dc.metaDir, privateDirMode)
		}
		return "", err
	}
	if err := os.MkdirAll(dc.metaDir, privateDirMode); err != nil {
		_ = os.Rename(oldMetaDir, dc.metaDir)
		return "", err
	}
	return oldMetaDir, nil
}

// SizeBytes returns the cached byte count, once the startup scan has
// accounted the entries already on disk.
func (dc *DiskCache) SizeBytes() int64 {
	dc.waitForScan()
	return dc.stats.SizeBytes.Load()
}

// StartMaintenance starts the maintenance goroutine if it is not already running.
func (dc *DiskCache) StartMaintenance() {
	dc.maintMu.Lock()
	defer dc.maintMu.Unlock()
	if dc.maintActive {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	dc.maintCancel = cancel
	dc.maintDone = done
	dc.maintActive = true
	go func() {
		defer close(done)
		dc.maintenanceLoop(ctx)
	}()
}

// StopMaintenance stops the maintenance goroutine if it is running, and closes
// data files kept open for downloads. Files reopen if the cache is used again.
func (dc *DiskCache) StopMaintenance() {
	dc.maintMu.Lock()
	if !dc.maintActive {
		dc.maintMu.Unlock()
		dc.closeAllDataFiles()
		return
	}
	cancel := dc.maintCancel
	done := dc.maintDone
	if cancel != nil {
		cancel()
	}
	if done != nil {
		<-done
	}
	dc.maintCancel = nil
	dc.maintDone = nil
	dc.maintActive = false
	dc.maintMu.Unlock()

	dc.persistLRUState()
	dc.closeAllDataFiles()
}

// Stats returns the current cache statistics, once the startup scan has
// accounted the entries already on disk.
func (dc *DiskCache) Stats() *cache.Stats {
	dc.waitForScan()
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

	// Snapshot pinned files information
	dc.pinnedFilesMu.Lock()
	s.PinnedCount = len(dc.pinnedFiles)
	s.PinnedPaths = make(map[string]int, len(dc.pinnedFiles))
	for path, count := range dc.pinnedFiles {
		s.PinnedPaths[path] = count
	}
	dc.pinnedFilesMu.Unlock()

	// Get LRU keys ordered from oldest to newest
	// lru.Keys() already returns a copy, so no need to copy again
	s.LruKeys = dc.lru.Keys()

	s.CapacityBytesRemaining = s.CapacityBytes - s.SizeBytes.Load()
	s.FileCountRemaining = s.MaxFileCount - s.FileCount.Load()
	return s
}

// Pin increments the reference count for a file, preventing it from being evicted.
// This should be called when a file handle is opened.
func (dc *DiskCache) Pin(path string) {
	dc.writeMu.RLock()
	defer dc.writeMu.RUnlock()
	dc.pinnedFilesMu.Lock()
	defer dc.pinnedFilesMu.Unlock()

	// Use fqPath as the key to match LRU keys
	fqPath := dc.entryPath(path)
	dc.pinnedFiles[fqPath]++
	dc.log.Trace("DiskCache: pinned %s (fqPath: %s, count: %d)", path, fqPath, dc.pinnedFiles[fqPath])
}

func (dc *DiskCache) PinPartial(path string) {
	dc.Pin(dc.partialEntryPath(path))
}

// Unpin decrements the reference count for a file.
// When the count reaches zero, the file becomes eligible for eviction.
// This should be called when a file handle is closed.
func (dc *DiskCache) Unpin(path string) {
	dc.pinnedFilesMu.Lock()

	// Use fqPath as the key to match LRU keys
	fqPath := dc.entryPath(path)
	dc.pinnedFiles[fqPath]--
	fullyUnpinned := dc.pinnedFiles[fqPath] <= 0
	if fullyUnpinned {
		delete(dc.pinnedFiles, fqPath)
		dc.log.Trace("DiskCache: unpinned %s (fqPath: %s, fully released)", path, fqPath)
	} else {
		dc.log.Trace("DiskCache: unpinned %s (fqPath: %s, count: %d)", path, fqPath, dc.pinnedFiles[fqPath])
	}
	_, clearPending := dc.clearPending[fqPath]

	dc.pinnedFilesMu.Unlock()
	if !fullyUnpinned {
		return
	}
	if !clearPending && dc.Capacity == 0 && dc.MaxFileCount == 0 {
		return
	}

	dc.writeMu.Lock()
	defer dc.writeMu.Unlock()

	// A new handle may have pinned the entry while this close waited for the
	// exclusive lock. Leave a pending clear in place for its eventual close.
	dc.pinnedFilesMu.Lock()
	if dc.pinnedFiles[fqPath] > 0 {
		dc.pinnedFilesMu.Unlock()
		return
	}
	_, clearPending = dc.clearPending[fqPath]
	if clearPending {
		delete(dc.clearPending, fqPath)
	}
	dc.pinnedFilesMu.Unlock()

	if clearPending {
		_ = dc.deleteMetadataForEntryPath(fqPath)
		if err := dc.deleteFile(fqPath); err != nil {
			dc.log.Trace("DiskCache: error deleting cleared cache file %s: %v", fqPath, err)
		}
		dc.lru.Remove(fqPath)
		return
	}

	// If this file was fully unpinned and the cache is over capacity, try to trim the cache
	// back to the configured limit to respect the user's settings. This ensures the
	// cache doesn't stay bloated after files are closed.
	if dc.Capacity > 0 || dc.MaxFileCount > 0 {
		// First, try to evict the file that was just unpinned since it just became eligible
		if !dc.hasCapacityDelta(0, false) {
			dc.lru.Remove(fqPath)
		}

		// Then evict other oldest files until the cache is back under capacity
		evictionAttempts := 0
		for !dc.hasCapacityDelta(0, false) {
			oldestKey, _, ok := dc.lru.GetOldest()
			if !ok {
				// LRU is empty, nothing more to evict
				break
			}

			// Don't attempt to evict pinned files
			if dc.isPinned(oldestKey) {
				evictionAttempts++
				if evictionAttempts >= maxEvictionAttempts {
					// Likely all remaining files are pinned, cache cannot be reduced further
					break
				}
				// Try to remove it anyway - onEvict will re-add it and provide the next oldest
				dc.lru.Remove(oldestKey)
				continue
			}

			// Remove the oldest entry (which is not pinned)
			evictionAttempts = 0 // Reset counter since an evictable file was found
			dc.lru.Remove(oldestKey)
		}
	}
}

func (dc *DiskCache) UnpinPartial(path string) {
	dc.Unpin(dc.partialEntryPath(path))
}

// isPinned checks if a file is currently pinned (has open file handles).
// The key parameter should be the fqPath (same as LRU keys).
func (dc *DiskCache) isPinned(key string) bool {
	dc.pinnedFilesMu.Lock()
	defer dc.pinnedFilesMu.Unlock()

	count, exists := dc.pinnedFiles[key]
	return exists && count > 0
}

// onEvict is called when a key is evicted from the LRU. This includes when the key is removed
// explicitly via Remove or RemoveOldest.
func (dc *DiskCache) onEvict(path string, value struct{}) {
	// NEVER evict pinned files (files with open handles)
	if dc.isPinned(path) {
		dc.log.Trace("DiskCache: skipping eviction of pinned file: %s", path)
		// Re-add to LRU to allow retry later
		_ = dc.lru.Add(path, struct{}{})
		return
	}

	if err := dc.deleteFile(path); err != nil {
		dc.log.Trace("DiskCache: error deleting evicted cached file %s: %v", path, err)
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
	// Eviction must discard process-local coverage as well as durable metadata,
	// including when the data file has already disappeared.
	dc.deleteRangeWritePath(fqPath)
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

	_ = dc.deleteMetadataForEntryPath(fqPath)
	dc.closeDataFile(fqPath)
	if err := os.RemoveAll(fqPath); err != nil {
		dc.log.Trace("DiskCache: error deleting evicted cached file %s: %v", fqPath, err)
		return err
	}

	dc.removeAccountedFile(fqPath)
	return nil
}

// removeAbandonedPartials deletes partial data left by a previous process.
// It runs before the cache is returned, while nothing can be writing there.
func (dc *DiskCache) removeAbandonedPartials() error {
	return filepath.Walk(dc.partDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return os.Remove(p)
		}
		return nil
	})
}

// listEntries returns the data files currently in the cache.
func (dc *DiskCache) listEntries() ([]string, error) {
	var paths []string
	err := filepath.WalkDir(dc.dataDir, func(p string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			paths = append(paths, p)
		}
		return nil
	})
	return paths, err
}

// loadStats accounts the entries a previous process left in the cache and
// removes data without valid metadata. It runs in the background after the
// cache opens, because reading every entry's metadata can take a noticeable
// time with a full cache. Each entry is checked under a short lock, and an
// entry written before the scan reaches it has already accounted for itself.
func (dc *DiskCache) loadStats(paths []string) {
	defer close(dc.scanDone)
	dc.maintenanceRunMu.Lock()
	defer dc.maintenanceRunMu.Unlock()
	start := time.Now()
	for _, p := range paths {
		dc.writeMu.Lock()
		dc.scanEntryLocked(p)
		dc.writeMu.Unlock()
	}

	dc.writeMu.Lock()
	dc.stats.LoadDuration = time.Since(start)
	dc.stats.LruCount = len(dc.lru.Keys())
	dc.writeMu.Unlock()
}

// scanEntryLocked runs under writeMu.
func (dc *DiskCache) scanEntryLocked(p string) {
	if _, tracked := dc.accountedSize(p, 0); tracked {
		return
	}
	info, err := os.Stat(p)
	if err != nil {
		return
	}
	meta, ok := dc.validMetadataForEntryPath(p, info.Size())
	if !ok {
		if dc.isPinned(p) {
			// An open handle may be about to write this entry.
			return
		}
		_ = dc.deleteMetadataForEntryPath(p)
		if err := dc.deleteFile(p); err != nil {
			dc.log.Debug("DiskCache: loadStats: failed to remove abandoned data %s: %v", p, err)
		}
		return
	}
	size, err := meta.CachedByteCount()
	if err != nil {
		return
	}
	dc.setAccountedSize(p, size)
	if !dc.lru.Contains(p) {
		dc.log.Trace("DiskCache: loadStats: LRU does not contain %s", p)
	}
}

// waitForScan blocks until the startup scan has accounted existing entries.
func (dc *DiskCache) waitForScan() {
	if dc.scanDone != nil {
		<-dc.scanDone
	}
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
	f, err := os.CreateTemp(path, ".cache_write_test-*")
	if err != nil {
		return fmt.Errorf("diskCache: cache root path is not writable: %s", path)
	}
	f.Close()
	return os.Remove(f.Name())
}

// privateDirMode and privateFileMode are the modes for everything the cache creates below its
// root. They grant access to the owner only; on Windows access comes from the private directories
// created by NewDiskCache, which hand their access control list down to what is created inside.
const (
	privateDirMode  = 0o700
	privateFileMode = 0o600
)

// checkTree reports whether every existing entry below dir, which may not exist yet, may be made
// private: an ordinary file or directory this user owns, with a single name. Links, files with
// other names, and entries owned by someone else are reported so the cache is not used and
// nothing they share is changed. Nothing is modified.
func checkTree(dir string) error {
	return walkTree(dir, privatefile.CheckEntry)
}

// makeTreePrivate makes every existing entry below dir, a private directory the cache owns,
// readable only by the current user. A cache written before entries were created private has
// world-readable data, partial data, metadata and LRU state; on Windows an entry can also carry
// its own permissive access control entries that the directory's do not override. Each entry is
// examined as it is (not through a link) and repaired only when checkTree would accept it.
func makeTreePrivate(dir string) error {
	return walkTree(dir, privatefile.EnsureEntryPrivate)
}

func walkTree(dir string, visit func(path string, info fs.FileInfo) error) error {
	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == dir {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		return visit(path, info)
	})
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
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
	if dc.RangeFlushBytes <= 0 {
		return fmt.Errorf("diskCache: range flush size must be positive: %d", dc.RangeFlushBytes)
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
	if path != dc.partDir && strings.HasPrefix(path, dc.partDir) {
		// path is already a full path in the cache
		return path
	}
	return dc.cachePath(dc.dataDir, path)
}

func (dc *DiskCache) partialEntryPath(path string) string {
	if path != dc.partDir && strings.HasPrefix(path, dc.partDir) {
		return path
	}
	return dc.cachePath(dc.partDir, path)
}

func (dc *DiskCache) cachePath(root string, path string) string {
	sum := sha256.Sum256([]byte(path))
	h := hex.EncodeToString(sum[:])
	n := min(shardPrefixLen, len(h))
	dir := h[:n]
	name := h[n:] + "-" + filepath.Base(path)
	return filepath.Join(root, dir, name)
}

func (dc *DiskCache) metadataPath(path string) string {
	sum := sha256.Sum256([]byte(path))
	h := hex.EncodeToString(sum[:])
	n := min(shardPrefixLen, len(h))
	dir := h[:n]
	name := h[n:] + ".json"
	return filepath.Join(dc.metaDir, dir, name)
}

func (dc *DiskCache) metadataPathForEntryPath(path string) (string, bool) {
	fqPath := dc.entryPath(path)
	rel, err := filepath.Rel(dc.dataDir, fqPath)
	if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", false
	}
	shard := filepath.Dir(rel)
	if shard == "." || shard == "" || strings.Contains(shard, string(os.PathSeparator)) {
		return "", false
	}
	hashRest, _, ok := strings.Cut(filepath.Base(rel), "-")
	if !ok || hashRest == "" {
		return "", false
	}
	return filepath.Join(dc.metaDir, shard, hashRest+".json"), true
}

func (dc *DiskCache) validMetadataForEntryPath(path string, size int64) (cache.EntryMetadata, bool) {
	dc.rangeWritesMu.Lock()
	state, buffered := dc.rangeWrites[dc.entryPath(path)]
	dc.rangeWritesMu.Unlock()
	metaPath, ok := dc.metadataPathForEntryPath(path)
	if !ok {
		return cache.EntryMetadata{}, false
	}
	var meta cache.EntryMetadata
	if buffered {
		meta = state.metadata
	} else {
		var err error
		meta, err = readMetadataFile(metaPath)
		if err != nil {
			return cache.EntryMetadata{}, false
		}
	}
	if dc.entryPath(meta.Path) != dc.entryPath(path) {
		return cache.EntryMetadata{}, false
	}
	if buffered {
		// This process wrote those bytes. While a download holds the file
		// open, Windows can list it with a stale size, so size only checks
		// stored metadata.
		return meta, true
	}
	for _, r := range meta.Ranges {
		if r.End > size {
			return cache.EntryMetadata{}, false
		}
	}
	return meta, true
}

func (dc *DiskCache) readEntryMetadata(path string) (cache.EntryMetadata, error) {
	return readMetadataFile(dc.metadataPath(path))
}

func readMetadataFile(path string) (cache.EntryMetadata, error) {
	var meta cache.EntryMetadata
	data, err := os.ReadFile(path)
	if err != nil {
		return meta, err
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return meta, fmt.Errorf("%w: %v", cache.ErrInvalidEntryMetadata, err)
	}
	if err := meta.Validate(); err != nil {
		return meta, err
	}
	return meta, nil
}

func (dc *DiskCache) deleteMetadata(path string) error {
	err := os.Remove(dc.metadataPath(path))
	if err == nil || errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func (dc *DiskCache) deleteMetadataForEntryPath(path string) error {
	metaPath, ok := dc.metadataPathForEntryPath(path)
	if !ok {
		return nil
	}
	err := os.Remove(metaPath)
	if err == nil || errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
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
	// When restoring, keys are added in order to preserve LRU order.
	keys := dc.lru.Keys()
	state := lruState{Keys: keys}
	statePath := filepath.Join(dc.stateDir, lruStateFile)
	f, err := os.OpenFile(statePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, privateFileMode)
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
		dc.log.Debug("DiskCache: error saving LRU state: %v", err)
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
	// Skip this cycle instead of queueing behind a running Clear. Blocking
	// here would force StopMaintenance, and therefore unmount, to wait out
	// the full clear.
	if !dc.maintenanceRunMu.TryLock() {
		return
	}
	defer dc.maintenanceRunMu.Unlock()
	select {
	case <-ctx.Done():
		return
	default:
	}

	// Remove metadata trees that Clear rotated aside but a crash prevented
	// from being deleted. Clear holds maintenanceRunMu while it removes the
	// rotated tree, so anything matching here is an orphan.
	staleMetaDirs, err := filepath.Glob(dc.metaDir + ".clear-*")
	if err == nil {
		for _, staleMetaDir := range staleMetaDirs {
			if err := os.RemoveAll(staleMetaDir); err != nil {
				dc.log.Debug("DiskCache: maintenance: error removing stale metadata directory %s: %v", staleMetaDir, err)
			}
		}
	}

	var deleteCandidates []fileMeta
	var retainedSize, retainedCount int64
	filesOnDisk := make(map[string]struct{}, 1024)
	var notInLru []string

	// create a snapshot of files to delete along with the predicted size and count of files
	// that will be left after deletions
	dc.writeMu.Lock()
	dc.log.Debug("DiskCache: performing maintenance")
	start := time.Now()
	for _, root := range []string{dc.dataDir, dc.partDir} {
		_ = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
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
				dc.log.Debug("DiskCache: maintenance: error stating file %s: %v", p, ierr)
				return nil
			}

			filesOnDisk[p] = struct{}{}

			meta, validMetadata := dc.validMetadataForEntryPath(p, info.Size())
			if !validMetadata {
				if dc.isPinned(p) {
					accounted, _ := dc.accountedSize(p, info.Size())
					retainedCount++
					retainedSize += accounted
					if !dc.lru.Contains(p) {
						notInLru = append(notInLru, p)
					}
					return nil
				}
				deleteCandidates = append(deleteCandidates, fileMeta{path: p, size: info.Size(), mod: info.ModTime()})
				return nil
			}
			accounted, metaErr := meta.CachedByteCount()
			if metaErr != nil {
				deleteCandidates = append(deleteCandidates, fileMeta{path: p, size: info.Size(), mod: info.ModTime()})
				return nil
			}

			expired := dc.MaxAge > 0 && time.Since(info.ModTime()) > dc.MaxAge
			if expired {
				if dc.isPinned(p) {
					retainedCount++
					retainedSize += accounted
					if !dc.lru.Contains(p) {
						notInLru = append(notInLru, p)
					}
					return nil
				}
				deleteCandidates = append(deleteCandidates, fileMeta{path: p, size: info.Size(), mod: info.ModTime()})
				return nil
			}
			retainedCount++
			retainedSize += accounted

			if !dc.lru.Contains(p) {
				// these will be added later outside the lock
				notInLru = append(notInLru, p)
			}
			return nil
		})
	}
	dc.writeMu.Unlock()

	// apply deletions, no need to hold any locks as deleteFile holds delMu
	for _, old := range deleteCandidates {
		_ = dc.deleteMetadataForEntryPath(old.path)
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
	// Rebuild from the live accounting map. Range writes can finish after the
	// filesystem snapshot above, so committing the snapshot would lose those
	// concurrent counter updates.
	dc.writeMu.Lock()
	retainedSize, retainedCount = dc.refreshAccountedStats()
	dc.stats.CapacityBytes = dc.Capacity
	dc.stats.MaxFileCount = dc.MaxFileCount
	dc.stats.CapacityBytesRemaining = dc.Capacity - retainedSize
	dc.stats.FileCountRemaining = dc.MaxFileCount - retainedCount
	dc.stats.LoadDuration = duration
	dc.stats.LruCount = len(dc.lru.Keys())
	dc.writeMu.Unlock()

	dc.log.Debug("DiskCache: maintenance complete, duration=%vms", duration.Milliseconds())
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
			if dc.isPinned(fqPath) {
				return false, nil
			}
			_ = dc.deleteMetadataForEntryPath(fqPath)
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
	for _, root := range []string{dc.dataDir, dc.partDir} {
		_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if dc.isPinned(p) {
				return nil
			}
			if oldest == nil || info.ModTime().Before(oldest.t) {
				oldest = &fi{path: p, t: info.ModTime()}
			}
			return nil
		})
	}
	if oldest == nil {
		return false
	}
	// Prefer going through LRU if present to trigger onEvict; otherwise delete directly.
	if removed := dc.lru.Remove(oldest.path); removed {
		return true
	}
	if err := dc.deleteFile(oldest.path); err != nil {
		dc.log.Trace("DiskCache: error deleting oldest by mtime %s: %v", oldest.path, err)
		return false
	}
	return true
}
