//go:build linux || windows

package disk

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
)

type rangeWriteState struct {
	metadata   cache.EntryMetadata
	dirtyBytes int64
	sparse     bool
}

// RangeMetadata returns the validator and coverage belonging to a listing identity.
func (dc *DiskCache) RangeMetadata(path string, meta cache.EntryMetadata) (cache.EntryMetadata, bool, error) {
	dc.writeMu.RLock()
	defer dc.writeMu.RUnlock()
	return dc.rangeMetadata(path, meta)
}

// InvalidateRanges waits for readers before removing metadata. Windows can
// reject deletion while another reader has the metadata file open.
func (dc *DiskCache) InvalidateRanges(path string) error {
	dc.writeMu.Lock()
	defer dc.writeMu.Unlock()
	dc.deleteRangeWritePath(path)
	return dc.deleteMetadata(path)
}

// ReadRange reads only when every requested byte is valid for the remote file version.
func (dc *DiskCache) ReadRange(path string, meta cache.EntryMetadata, buff []byte, ofst int64) (int, error) {
	if dc.Disabled || len(buff) == 0 {
		return 0, nil
	}
	requestedMeta, err := remoteVersionMetadata(path, meta)
	if err != nil {
		return 0, err
	}
	requested, err := boundedBufferRange(ofst, len(buff), requestedMeta.Size, true)
	if err != nil {
		return 0, err
	}
	if requested.Start == requested.End {
		return 0, nil
	}

	dc.writeMu.RLock()
	defer dc.writeMu.RUnlock()

	stored, ok, err := dc.rangeMetadata(path, requestedMeta)
	if err != nil || !ok {
		return 0, err
	}
	covered, err := cache.ByteRangeCovered(stored.Ranges, requested)
	if err != nil {
		return 0, err
	}
	if !covered {
		return 0, nil
	}

	want := requested.End - requested.Start
	n, err := dc.read(path, buff[:want], requested.Start)
	if err != nil {
		return 0, err
	}
	if int64(n) != want {
		_ = dc.Delete(path)
		return 0, nil
	}
	return n, nil
}

// MissingRanges returns the requested ranges not present for the remote file version.
func (dc *DiskCache) MissingRanges(path string, meta cache.EntryMetadata, requested cache.ByteRange) ([]cache.ByteRange, error) {
	requestedMeta, err := remoteVersionMetadata(path, meta)
	if err != nil {
		return nil, err
	}
	if err := requested.Validate(); err != nil {
		return nil, err
	}
	if requested.End > requestedMeta.Size {
		return nil, fmt.Errorf("%w: requested range [%d, %d) exceeds file size %d", cache.ErrInvalidByteRange, requested.Start, requested.End, requestedMeta.Size)
	}
	if dc.Disabled {
		if requested.Start == requested.End {
			return nil, nil
		}
		return []cache.ByteRange{requested}, nil
	}

	dc.writeMu.RLock()
	defer dc.writeMu.RUnlock()
	stored, ok, err := dc.rangeMetadata(path, requestedMeta)
	if err != nil {
		return nil, err
	}
	if !ok {
		if requested.Start == requested.End {
			return nil, nil
		}
		return []cache.ByteRange{requested}, nil
	}
	return cache.MissingByteRanges(stored.Ranges, requested)
}

// RangeEntryComplete reports whether the cached ranges cover the remote file.
func (dc *DiskCache) RangeEntryComplete(path string, meta cache.EntryMetadata) (bool, error) {
	if dc.Disabled {
		return false, nil
	}
	requestedMeta, err := remoteVersionMetadata(path, meta)
	if err != nil {
		return false, err
	}

	dc.writeMu.RLock()
	defer dc.writeMu.RUnlock()
	stored, ok, err := dc.rangeMetadata(path, requestedMeta)
	if err != nil || !ok {
		return false, err
	}
	return stored.IsComplete(), nil
}

// WriteRange makes received bytes readable immediately. Durable coverage is
// published every RangeFlushBytes of new data; FlushRanges publishes the rest
// when a download stops.
func (dc *DiskCache) WriteRange(path string, meta cache.EntryMetadata, buff []byte, ofst int64) (int, error) {
	return dc.writeRange(path, meta, buff, ofst, true)
}

// WriteCompleteRange shares the publication policy without requiring sparse
// file support, since a complete response fills every byte through EOF.
func (dc *DiskCache) WriteCompleteRange(path string, meta cache.EntryMetadata, buff []byte, ofst int64) (int, error) {
	return dc.writeRange(path, meta, buff, ofst, false)
}

func (dc *DiskCache) writeRange(path string, meta cache.EntryMetadata, buff []byte, ofst int64, sparse bool) (int, error) {
	if dc.Disabled || len(buff) == 0 {
		return 0, nil
	}
	requested, err := remoteVersionMetadata(path, meta)
	if err != nil {
		return 0, err
	}
	written, err := boundedBufferRange(ofst, len(buff), requested.Size, false)
	if err != nil {
		return 0, err
	}

	// Retain the file while syncing outside the cache-wide lock. Other files
	// can publish concurrently; writes to this entry remain serialized.
	dc.Pin(path)
	defer dc.Unpin(path)
	dc.rangeWriteLocks.Lock(dc.entryPath(path))
	defer dc.rangeWriteLocks.Unlock(dc.entryPath(path))
	dc.writeMu.Lock()
	state, err := dc.writeRangeChunk(path, requested, buff, written, sparse)
	dc.writeMu.Unlock()
	if err != nil {
		return 0, err
	}
	if state.dirtyBytes >= dc.RangeFlushBytes || state.metadata.IsComplete() {
		if err := dc.flushRangeWrite(path, state); err != nil {
			return 0, err
		}
	}
	return len(buff), nil
}

// writeRangeChunk runs under writeMu. Published state values are immutable so
// readers and a sync in progress can retain their own snapshot safely.
func (dc *DiskCache) writeRangeChunk(path string, requested cache.EntryMetadata, buff []byte, written cache.ByteRange, sparse bool) (*rangeWriteState, error) {
	var state rangeWriteState
	if current, ok := dc.rangeWrite(path, requested); ok && current.metadata.SameContent(requested) {
		state = *current
	} else {
		dc.deleteRangeWritePath(path)
		stored, reused, err := dc.prepareRangeEntry(path, requested)
		if err != nil {
			return nil, err
		}
		if !reused {
			// A clear that found this entry pinned deletes it on the last
			// Unpin. Bytes written after the clear replace that entry.
			dc.pinnedFilesMu.Lock()
			delete(dc.clearPending, dc.entryPath(path))
			dc.pinnedFilesMu.Unlock()
		}
		state.metadata = stored
	}
	missing, err := cache.MissingByteRanges(state.metadata.Ranges, written)
	if err != nil {
		return nil, err
	}
	newBytes := byteRangeCount(missing)
	if newBytes == 0 {
		if current, ok := dc.rangeWrite(path, requested); ok {
			return current, nil
		}
		return &state, nil
	}
	fqPath := dc.entryPath(path)
	_, tracked := dc.accountedSize(fqPath, 0)
	dc.makeCapacity(newBytes, !tracked, fqPath)
	if err := os.MkdirAll(filepath.Dir(fqPath), privateDirMode); err != nil {
		return nil, err
	}
	df, err := dc.openDataFile(fqPath, true)
	if err != nil {
		return nil, err
	}
	if sparse && !state.sparse {
		if err := markSparseFile(df.file); err != nil {
			dc.releaseDataFile(df)
			if len(state.metadata.Ranges) == 0 {
				dc.closeDataFile(fqPath)
				_ = os.Remove(fqPath)
			}
			return nil, err
		}
		state.sparse = true
	}
	for _, gap := range missing {
		start, end := gap.Start-written.Start, gap.End-written.Start
		n, err := df.file.WriteAt(buff[start:end], gap.Start)
		if err != nil || int64(n) != gap.End-gap.Start {
			dc.releaseDataFile(df)
			if err == nil {
				err = io.ErrShortWrite
			}
			return nil, err
		}
	}
	dc.releaseDataFile(df)
	previous := byteRangeCount(state.metadata.Ranges)
	updated, err := state.metadata.WithRanges(append(state.metadata.Ranges, written))
	if err != nil {
		return nil, err
	}
	state.metadata = updated
	state.dirtyBytes += newBytes
	dc.setRangeWrite(path, &state)
	if tracked {
		dc.addAccountedBytes(fqPath, newBytes)
	} else {
		// The startup scan has not reached this entry, so account for the
		// bytes it already had as well.
		dc.setAccountedSize(fqPath, previous+newBytes)
	}
	dc.stats.WriteCount.Add(1)
	dc.stats.WriteBytes.Add(newBytes)
	dc.lru.Add(fqPath, struct{}{})
	dc.lruDirty.Store(true)
	return &state, nil
}

// FlushRanges makes all currently readable bytes durable at request completion.
func (dc *DiskCache) FlushRanges(path string, meta cache.EntryMetadata) error {
	dc.Pin(path)
	defer dc.Unpin(path)
	dc.rangeWriteLocks.Lock(dc.entryPath(path))
	defer dc.rangeWriteLocks.Unlock(dc.entryPath(path))
	dc.writeMu.RLock()
	state, ok := dc.rangeWrite(path, meta)
	dc.writeMu.RUnlock()
	// The download that called this is done with the entry's data file.
	defer dc.closeDataFile(dc.entryPath(path))
	if !ok || !state.metadata.SameContent(meta) {
		return nil
	}
	return dc.flushRangeWrite(path, state)
}

// The entry lock stays held across sync, but the cache-wide lock only protects
// the final metadata rename. Clear, eviction, or replacement invalidates the
// snapshot and prevents it from publishing after retirement.
func (dc *DiskCache) flushRangeWrite(path string, state *rangeWriteState) error {
	if state.dirtyBytes == 0 {
		return nil
	}
	df, err := dc.openDataFile(dc.entryPath(path), false)
	if err != nil {
		return err
	}
	// Release before taking writeMu: closing the file waits for its users.
	syncErr := df.file.Sync()
	dc.releaseDataFile(df)
	if syncErr != nil {
		return syncErr
	}
	tmpPath, err := dc.prepareEntryMetadata(state.metadata)
	if err != nil {
		return err
	}
	defer os.Remove(tmpPath)
	dc.writeMu.Lock()
	defer dc.writeMu.Unlock()
	current, ok := dc.rangeWrite(path, state.metadata)
	if !ok || current != state {
		return context.Canceled
	}
	if err := os.Rename(tmpPath, dc.metadataPath(path)); err != nil {
		return err
	}
	if state.metadata.IsComplete() {
		dc.deleteRangeWrite(path, state.metadata)
	} else {
		durable := *state
		durable.dirtyBytes = 0
		dc.setRangeWrite(path, &durable)
	}
	return nil
}

func (dc *DiskCache) rangeWrite(path string, meta cache.EntryMetadata) (*rangeWriteState, bool) {
	dc.rangeWritesMu.Lock()
	defer dc.rangeWritesMu.Unlock()
	state, ok := dc.rangeWrites[dc.entryPath(path)]
	if !ok || !state.metadata.Matches(meta) {
		return nil, false
	}
	return state, true
}

func (dc *DiskCache) setRangeWrite(path string, state *rangeWriteState) {
	dc.rangeWritesMu.Lock()
	dc.rangeWrites[dc.entryPath(path)] = state
	dc.rangeWritesMu.Unlock()
}

func (dc *DiskCache) deleteRangeWrite(path string, meta cache.EntryMetadata) {
	path = dc.entryPath(path)
	dc.rangeWritesMu.Lock()
	deleted := false
	if state, ok := dc.rangeWrites[path]; ok && state.metadata.Matches(meta) {
		delete(dc.rangeWrites, path)
		deleted = true
	}
	dc.rangeWritesMu.Unlock()
	if deleted {
		dc.closeDataFile(path)
	}
}

// deleteRangeWritePath forgets an entry's download progress and closes its
// data file, which a replacement or removal may need to delete.
func (dc *DiskCache) deleteRangeWritePath(path string) {
	fqPath := dc.entryPath(path)
	dc.rangeWritesMu.Lock()
	delete(dc.rangeWrites, fqPath)
	dc.rangeWritesMu.Unlock()
	dc.closeDataFile(fqPath)
}

func (dc *DiskCache) prepareRangeEntry(path string, requested cache.EntryMetadata) (cache.EntryMetadata, bool, error) {
	stored, err := dc.readEntryMetadata(path)
	if err == nil && stored.SameContent(requested) {
		info, statErr := os.Stat(dc.entryPath(path))
		if statErr == nil {
			validData := true
			for _, r := range stored.Ranges {
				if r.End > info.Size() {
					validData = false
					break
				}
			}
			if validData {
				return stored, true, nil
			}
		} else if !errors.Is(statErr, os.ErrNotExist) {
			return cache.EntryMetadata{}, false, statErr
		}
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) && !errors.Is(err, cache.ErrInvalidEntryMetadata) {
		return cache.EntryMetadata{}, false, err
	}

	fqPath := dc.entryPath(path)
	_ = dc.deleteMetadata(path)
	// A pin prevents capacity eviction, but it cannot preserve bytes from a
	// different remote version. Replace stale data while retaining the handle's
	// pin count for the new entry.
	if err := dc.deleteFile(fqPath); err != nil {
		return cache.EntryMetadata{}, false, err
	}
	empty, err := requested.WithRanges(nil)
	return empty, false, err
}

func (dc *DiskCache) rangeMetadata(path string, requested cache.EntryMetadata) (cache.EntryMetadata, bool, error) {
	if state, ok := dc.rangeWrite(path, requested); ok {
		return state.metadata, true, nil
	}
	stored, err := dc.readEntryMetadata(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cache.EntryMetadata{}, false, nil
		}
		if errors.Is(err, cache.ErrInvalidEntryMetadata) {
			_ = dc.Delete(path)
			return cache.EntryMetadata{}, false, nil
		}
		return cache.EntryMetadata{}, false, err
	}
	if !stored.Matches(requested) {
		// An older open handle can still ask for its version after a newer
		// version has replaced the entry. That stale request must miss without
		// invalidating the newer cached bytes.
		if !stored.ModTime.After(requested.ModTime) {
			_ = dc.Delete(path)
		}
		return cache.EntryMetadata{}, false, nil
	}

	fqPath := dc.entryPath(path)
	info, err := os.Stat(fqPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			_ = dc.deleteMetadata(path)
			return cache.EntryMetadata{}, false, nil
		}
		return cache.EntryMetadata{}, false, err
	}
	if deleted, err := dc.deleteIfExpired(fqPath, info); err != nil || deleted {
		return cache.EntryMetadata{}, false, err
	}
	for _, r := range stored.Ranges {
		if r.End > info.Size() {
			_ = dc.Delete(path)
			return cache.EntryMetadata{}, false, nil
		}
	}
	return stored, true, nil
}

func remoteVersionMetadata(path string, meta cache.EntryMetadata) (cache.EntryMetadata, error) {
	meta.Path = path
	meta.Ranges = nil
	if err := meta.Validate(); err != nil {
		return cache.EntryMetadata{}, err
	}
	return meta, nil
}

func boundedBufferRange(ofst int64, length int, fileSize int64, clampEOF bool) (cache.ByteRange, error) {
	if ofst < 0 || int64(length) > fileSize-ofst {
		if !clampEOF || ofst < 0 || ofst > fileSize {
			return cache.ByteRange{}, fmt.Errorf("%w: offset %d and length %d exceed file size %d", cache.ErrInvalidByteRange, ofst, length, fileSize)
		}
		length = int(fileSize - ofst)
	}
	return cache.ByteRange{Start: ofst, End: ofst + int64(length)}, nil
}

func byteRangeCount(ranges []cache.ByteRange) int64 {
	var count int64
	for _, r := range ranges {
		count += r.End - r.Start
	}
	return count
}

func (dc *DiskCache) makeCapacity(delta int64, newFile bool, currentPath string) {
	evictionAttempts := 0
	for !dc.hasCapacityDelta(delta, newFile) {
		oldestKey, _, ok := dc.lru.GetOldest()
		if !ok {
			dc.deleteOneByMtime()
			break
		}
		if oldestKey == currentPath {
			break
		}
		if dc.isPinned(oldestKey) {
			evictionAttempts++
			if evictionAttempts >= maxEvictionAttempts {
				break
			}
			dc.lru.Remove(oldestKey)
			continue
		}
		evictionAttempts = 0
		dc.lru.Remove(oldestKey)
	}
}
