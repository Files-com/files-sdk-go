//go:build linux || windows

package disk

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// keptFileName holds the entries kept for offline use, by their path in the
// cache, so they stay kept after a restart.
const keptFileName = "kept.json"

// offlineDirName is where the mount saves what it needs to show kept folders
// when the storage can't be reached.
const offlineDirName = "offline"

// SetKept keeps the entries of the given paths for offline use, or stops
// keeping them. A kept entry is never removed to stay under the cache's size,
// age or file count limits: it stays until it is deleted, such as when the file
// changes or is removed on the storage, or until it is no longer kept.
func (dc *DiskCache) SetKept(paths []string, kept bool) error {
	if dc.Disabled || len(paths) == 0 {
		return nil
	}
	dc.keptMu.Lock()
	changed := false
	var released []string
	for _, path := range paths {
		fqPath := dc.entryPath(path)
		_, already := dc.kept[fqPath]
		switch {
		case kept && !already:
			dc.kept[fqPath] = struct{}{}
			changed = true
		case !kept && already:
			delete(dc.kept, fqPath)
			released = append(released, fqPath)
			changed = true
		}
	}
	var err error
	if changed {
		err = dc.saveKeptLocked()
	}
	dc.keptMu.Unlock()

	if kept {
		// Kept entries live outside the LRU, so its limits never reach them.
		// Removing them calls onEvict, which leaves kept data in place.
		for _, path := range paths {
			dc.lru.Remove(dc.entryPath(path))
		}
		// An entry kept again before its last handle closed stays.
		dc.pinnedFilesMu.Lock()
		for _, path := range paths {
			delete(dc.clearPending, dc.entryPath(path))
			delete(dc.deleteOnClose, dc.entryPath(path))
		}
		dc.pinnedFilesMu.Unlock()
	}
	// An entry that is no longer kept is managed by the limits again.
	for _, fqPath := range released {
		if _, statErr := os.Stat(fqPath); statErr == nil {
			dc.lru.Add(fqPath, struct{}{})
		}
	}
	return err
}

// DeleteWhenClosed marks the entry of path to be deleted when its last handle
// closes, if it is open, and reports whether it was. What is written to the
// entry before then is deleted too. An entry that isn't open is left for the
// caller to delete.
func (dc *DiskCache) DeleteWhenClosed(path string) bool {
	if dc.Disabled {
		return false
	}
	fqPath := dc.entryPath(path)
	dc.pinnedFilesMu.Lock()
	defer dc.pinnedFilesMu.Unlock()
	if dc.pinnedFiles[fqPath] <= 0 {
		return false
	}
	dc.deleteOnClose[fqPath] = struct{}{}
	return true
}

// IsKept reports whether the entry of path is kept for offline use.
func (dc *DiskCache) IsKept(path string) bool {
	return dc.isKept(dc.entryPath(path))
}

// OfflineDir returns a private directory in the cache's state where the mount
// saves what it needs to work without the storage.
func (dc *DiskCache) OfflineDir() (string, error) {
	if dc.Disabled {
		return "", errors.New("the disk cache is disabled")
	}
	dir := filepath.Join(dc.stateDir, offlineDirName)
	if err := os.MkdirAll(dir, privateDirMode); err != nil {
		return "", err
	}
	return dir, nil
}

func (dc *DiskCache) isKept(fqPath string) bool {
	dc.keptMu.Lock()
	defer dc.keptMu.Unlock()
	_, ok := dc.kept[fqPath]
	return ok
}

// isPinnedOrKept reports whether the size, age and file count limits must leave
// an entry alone.
func (dc *DiskCache) isPinnedOrKept(fqPath string) bool {
	return dc.isPinned(fqPath) || dc.isKept(fqPath)
}

// lruAdd tracks an entry for eviction, unless it is kept.
func (dc *DiskCache) lruAdd(fqPath string) {
	if dc.isKept(fqPath) {
		return
	}
	dc.lru.Add(fqPath, struct{}{})
}

func (dc *DiskCache) keptStatePath() string {
	return filepath.Join(dc.stateDir, keptFileName)
}

func (dc *DiskCache) loadKept() error {
	dc.keptMu.Lock()
	defer dc.keptMu.Unlock()
	dc.kept = make(map[string]struct{})
	data, err := os.ReadFile(dc.keptStatePath())
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	var entries []string
	if err := json.Unmarshal(data, &entries); err != nil {
		return fmt.Errorf("diskCache: unreadable kept entries: %w", err)
	}
	for _, name := range entries {
		// Saved names are relative to the data directory, so a moved cache root
		// keeps them.
		dc.kept[filepath.Join(dc.dataDir, name)] = struct{}{}
	}
	return nil
}

func (dc *DiskCache) saveKeptLocked() error {
	entries := make([]string, 0, len(dc.kept))
	for fqPath := range dc.kept {
		name, err := filepath.Rel(dc.dataDir, fqPath)
		if err != nil {
			continue
		}
		entries = append(entries, name)
	}
	sort.Strings(entries)
	data, err := json.Marshal(entries)
	if err != nil {
		return err
	}
	return os.WriteFile(dc.keptStatePath(), data, privateFileMode)
}

// restoreKeptMetadataLocked moves the metadata of kept entries from a metadata
// directory a clear rotated away back into the current one. It runs under
// writeMu.
func (dc *DiskCache) restoreKeptMetadataLocked(oldMetaDir string) {
	dc.keptMu.Lock()
	kept := make([]string, 0, len(dc.kept))
	for fqPath := range dc.kept {
		kept = append(kept, fqPath)
	}
	dc.keptMu.Unlock()
	for _, fqPath := range kept {
		current, ok := dc.metadataPathForEntryPath(fqPath)
		if !ok {
			continue
		}
		rel, err := filepath.Rel(dc.metaDir, current)
		if err != nil {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(current), privateDirMode); err != nil {
			continue
		}
		_ = os.Rename(filepath.Join(oldMetaDir, rel), current)
	}
}

// OfflineDirFor returns the directory OfflineDir uses for the cache rooted at
// root, without opening the cache.
func OfflineDirFor(root string) string {
	return filepath.Join(root, stateDir, offlineDirName)
}
