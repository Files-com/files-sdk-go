//go:build linux || windows

package mem

import (
	"fmt"
	"time"

	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
)

// RangeMetadata returns the validator and coverage belonging to a listing identity.
func (mc *MemoryCache) RangeMetadata(path string, meta cache.EntryMetadata) (cache.EntryMetadata, bool, error) {
	mc.filesMu.RLock()
	defer mc.filesMu.RUnlock()
	ent, ok := mc.files[path]
	if !ok || !ent.metadata.Matches(meta) || !entryBacksRanges(ent, ent.metadata.Ranges) {
		return cache.EntryMetadata{}, false, nil
	}
	result := ent.metadata
	result.Ranges = append([]cache.ByteRange(nil), result.Ranges...)
	return result, true, nil
}

// InvalidateRanges removes the completed-version identity, including when pinned.
func (mc *MemoryCache) InvalidateRanges(path string) error {
	mc.Delete(path)
	return nil
}

// ReadRange reads only when every requested byte is valid for the remote file version.
func (mc *MemoryCache) ReadRange(path string, meta cache.EntryMetadata, buff []byte, ofst int64) (int, error) {
	if len(buff) == 0 {
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

	mc.filesMu.RLock()
	ent, ok := mc.files[path]
	if !ok || ent.metadata.Version == 0 {
		mc.filesMu.RUnlock()
		return 0, nil
	}
	if !ent.metadata.Matches(requestedMeta) {
		stored := ent.metadata
		mc.filesMu.RUnlock()
		if !stored.ModTime.After(requestedMeta.ModTime) {
			mc.deleteRangeEntryIfVersion(path, stored)
		}
		return 0, nil
	}
	covered, err := cache.ByteRangeCovered(ent.metadata.Ranges, requested)
	if err != nil {
		mc.filesMu.RUnlock()
		return 0, err
	}
	if !covered {
		mc.filesMu.RUnlock()
		return 0, nil
	}
	if !entryBacksRange(ent, requested) {
		stored := ent.metadata
		mc.filesMu.RUnlock()
		mc.deleteRangeEntryIfVersion(path, stored)
		return 0, nil
	}

	readEntryRange(ent, buff, requested)
	mc.filesMu.RUnlock()

	mc.lru.Add(path, struct{}{})
	mc.stats.ReadCount.Add(1)
	mc.stats.ReadBytes.Add(requested.End - requested.Start)
	return int(requested.End - requested.Start), nil
}

// MissingRanges returns the requested ranges not present for the remote file version.
func (mc *MemoryCache) MissingRanges(path string, meta cache.EntryMetadata, requested cache.ByteRange) ([]cache.ByteRange, error) {
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

	mc.filesMu.RLock()
	ent, ok := mc.files[path]
	if !ok || ent.metadata.Version == 0 {
		mc.filesMu.RUnlock()
		return uncachedRange(requested), nil
	}
	if !ent.metadata.Matches(requestedMeta) {
		stored := ent.metadata
		mc.filesMu.RUnlock()
		if !stored.ModTime.After(requestedMeta.ModTime) {
			mc.deleteRangeEntryIfVersion(path, stored)
		}
		return uncachedRange(requested), nil
	}
	if !entryBacksRanges(ent, ent.metadata.Ranges) {
		stored := ent.metadata
		mc.filesMu.RUnlock()
		mc.deleteRangeEntryIfVersion(path, stored)
		return uncachedRange(requested), nil
	}
	missing, err := cache.MissingByteRanges(ent.metadata.Ranges, requested)
	mc.filesMu.RUnlock()
	return missing, err
}

// RangeEntryComplete reports whether the cached ranges cover the remote file.
func (mc *MemoryCache) RangeEntryComplete(path string, meta cache.EntryMetadata) (bool, error) {
	requestedMeta, err := remoteVersionMetadata(path, meta)
	if err != nil {
		return false, err
	}

	mc.filesMu.RLock()
	ent, ok := mc.files[path]
	if !ok || ent.metadata.Version == 0 {
		mc.filesMu.RUnlock()
		return false, nil
	}
	if !ent.metadata.Matches(requestedMeta) {
		stored := ent.metadata
		mc.filesMu.RUnlock()
		if !stored.ModTime.After(requestedMeta.ModTime) {
			mc.deleteRangeEntryIfVersion(path, stored)
		}
		return false, nil
	}
	metadataComplete := ent.metadata.IsComplete()
	backed := entryBacksRanges(ent, ent.metadata.Ranges)
	complete := metadataComplete && backed
	stored := ent.metadata
	mc.filesMu.RUnlock()
	if !backed {
		mc.deleteRangeEntryIfVersion(path, stored)
	}
	return complete, nil
}

// deleteRangeEntryIfVersion invalidates only the entry that was inspected.
// A stale reader must not remove a newer version published after its lookup.
func (mc *MemoryCache) deleteRangeEntryIfVersion(path string, observed cache.EntryMetadata) {
	mc.writeMu.Lock()
	defer mc.writeMu.Unlock()

	mc.filesMu.RLock()
	current, ok := mc.files[path]
	unchanged := ok && sameRangeEntryMetadata(current.metadata, observed)
	mc.filesMu.RUnlock()
	if unchanged {
		mc.deleteLocked(path)
	}
}

func sameRangeEntryMetadata(left, right cache.EntryMetadata) bool {
	if !left.SameContent(right) || len(left.Ranges) != len(right.Ranges) {
		return false
	}
	for i := range left.Ranges {
		if left.Ranges[i] != right.Ranges[i] {
			return false
		}
	}
	return true
}

// WriteRange stores only missing bytes and then publishes their valid ranges.
func (mc *MemoryCache) WriteRange(path string, meta cache.EntryMetadata, buff []byte, ofst int64) (int, error) {
	if len(buff) == 0 {
		return 0, nil
	}
	requestedMeta, err := remoteVersionMetadata(path, meta)
	if err != nil {
		return 0, err
	}
	writtenRange, err := boundedBufferRange(ofst, len(buff), requestedMeta.Size, false)
	if err != nil {
		return 0, err
	}

	mc.writeMu.Lock()
	defer mc.writeMu.Unlock()

	stored, err := mc.prepareRangeEntry(path, requestedMeta)
	if err != nil {
		return 0, err
	}
	missing, err := cache.MissingByteRanges(stored.Ranges, writtenRange)
	if err != nil {
		return 0, err
	}

	var newBytes int64
	for _, gap := range missing {
		start := gap.Start - writtenRange.Start
		end := gap.End - writtenRange.Start
		n, err := mc.writeLocked(path, buff[start:end], gap.Start)
		if err != nil {
			return 0, err
		}
		if int64(n) != gap.End-gap.Start {
			return 0, fmt.Errorf("memoryCache: ranged write at [%d, %d) wrote %d bytes", gap.Start, gap.End, n)
		}
		newBytes += int64(n)
	}

	updated, err := stored.WithRanges(append(stored.Ranges, writtenRange))
	if err != nil {
		return 0, err
	}
	mc.filesMu.Lock()
	ent, ok := mc.files[path]
	if !ok {
		mc.filesMu.Unlock()
		return 0, fmt.Errorf("memoryCache: ranged cache entry disappeared: %s", path)
	}
	ent.metadata = updated
	ent.mod = time.Now()
	mc.filesMu.Unlock()

	mc.lru.Add(path, struct{}{})
	mc.stats.WriteCount.Add(1)
	mc.stats.WriteBytes.Add(newBytes)
	return len(buff), nil
}

// WriteCompleteRange uses the same page store as a partial range. It exists so
// the coordinator can use one contract for dense disk fallback and memory.
func (mc *MemoryCache) WriteCompleteRange(path string, meta cache.EntryMetadata, buff []byte, ofst int64) (int, error) {
	return mc.WriteRange(path, meta, buff, ofst)
}

func (mc *MemoryCache) prepareRangeEntry(path string, requested cache.EntryMetadata) (cache.EntryMetadata, error) {
	mc.filesMu.RLock()
	ent, found := mc.files[path]
	if found && ent.metadata.SameContent(requested) && entryBacksRanges(ent, ent.metadata.Ranges) {
		stored := ent.metadata
		mc.filesMu.RUnlock()
		return stored, nil
	}
	mc.filesMu.RUnlock()

	if found {
		// A pin prevents capacity eviction, but it cannot preserve bytes from a
		// different remote version. Replace stale data while retaining the
		// handle's pin count for the new entry.
		mc.filesMu.Lock()
		if current, ok := mc.files[path]; ok {
			var freed int64
			for _, page := range current.pages {
				freed += int64(len(page))
			}
			delete(mc.files, path)
			mc.stats.SizeBytes.Add(-freed)
			mc.stats.FileCount.Add(-1)
		}
		mc.filesMu.Unlock()
	}
	return requested.WithRanges(nil)
}

func entryBacksRanges(ent *entry, ranges []cache.ByteRange) bool {
	for _, r := range ranges {
		if !entryBacksRange(ent, r) {
			return false
		}
	}
	return true
}

func entryBacksRange(ent *entry, requested cache.ByteRange) bool {
	if requested.End > ent.size {
		return false
	}
	for pos := requested.Start; pos < requested.End; {
		pageIdx := pos / pageSize
		if _, ok := ent.pages[pageIdx]; !ok {
			return false
		}
		pos = min(requested.End, (pageIdx+1)*pageSize)
	}
	return true
}

func readEntryRange(ent *entry, buff []byte, requested cache.ByteRange) {
	read := 0
	for pos := requested.Start; pos < requested.End; {
		pageIdx := pos / pageSize
		pageOff := int(pos % pageSize)
		chunkLen := int(min(requested.End-pos, int64(pageSize-pageOff)))
		copy(buff[read:read+chunkLen], ent.pages[pageIdx][pageOff:pageOff+chunkLen])
		read += chunkLen
		pos += int64(chunkLen)
	}
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

func uncachedRange(requested cache.ByteRange) []cache.ByteRange {
	if requested.Start == requested.End {
		return nil
	}
	return []cache.ByteRange{requested}
}
