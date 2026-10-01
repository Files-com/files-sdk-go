//go:build linux || windows

package disk

func (dc *DiskCache) accountedSize(path string, fallback int64) (int64, bool) {
	dc.accountedBytesMu.Lock()
	defer dc.accountedBytesMu.Unlock()
	size, ok := dc.accountedBytes[path]
	if !ok {
		return fallback, false
	}
	return size, true
}

func (dc *DiskCache) setAccountedSize(path string, size int64) {
	dc.accountedBytesMu.Lock()
	defer dc.accountedBytesMu.Unlock()

	previous, existed := dc.accountedBytes[path]
	dc.accountedBytes[path] = size

	if !existed {
		dc.stats.FileCount.Add(1)
	}
	dc.stats.SizeBytes.Add(size - previous)
}

func (dc *DiskCache) addAccountedBytes(path string, added int64) {
	dc.accountedBytesMu.Lock()
	defer dc.accountedBytesMu.Unlock()

	previous, existed := dc.accountedBytes[path]
	dc.accountedBytes[path] = previous + added

	if !existed {
		dc.stats.FileCount.Add(1)
	}
	dc.stats.SizeBytes.Add(added)
}

func (dc *DiskCache) removeAccountedFile(path string) {
	dc.accountedBytesMu.Lock()
	defer dc.accountedBytesMu.Unlock()

	size, existed := dc.accountedBytes[path]
	if existed {
		delete(dc.accountedBytes, path)
	}

	if !existed {
		return
	}
	dc.stats.FileCount.Add(-1)
	dc.stats.SizeBytes.Add(-size)
}

func (dc *DiskCache) refreshAccountedStats() (size, count int64) {
	dc.accountedBytesMu.Lock()
	defer dc.accountedBytesMu.Unlock()

	for _, accounted := range dc.accountedBytes {
		size += accounted
		count++
	}
	dc.stats.SizeBytes.Store(size)
	dc.stats.FileCount.Store(count)
	return size, count
}
