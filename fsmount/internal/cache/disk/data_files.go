//go:build linux || windows

package disk

import "os"

// dataFile is an entry's data file, kept open while the entry is downloaded.
// Opening and closing a file that is being written is expensive on Windows,
// where antivirus filters inspect it on every open and close; reopening it
// for each chunk made those inspections the main cost of a download.
type dataFile struct {
	file    *os.File
	users   int
	retired bool
	closed  chan struct{}
}

// openDataFile returns the entry's open data file, opening it when needed.
// The caller must not wait on any lock before calling releaseDataFile.
func (dc *DiskCache) openDataFile(fqPath string, create bool) (*dataFile, error) {
	dc.dataFilesMu.Lock()
	defer dc.dataFilesMu.Unlock()
	if df := dc.dataFiles[fqPath]; df != nil {
		df.users++
		return df, nil
	}
	flag := os.O_RDWR
	if create {
		flag |= os.O_CREATE
	}
	file, err := os.OpenFile(fqPath, flag, privateFileMode)
	if err != nil {
		return nil, err
	}
	df := &dataFile{file: file, users: 1, closed: make(chan struct{})}
	dc.dataFiles[fqPath] = df
	return df, nil
}

// sharedDataFile returns the data file of an entry being downloaded, if it is
// open, so reads can use it instead of opening their own.
func (dc *DiskCache) sharedDataFile(fqPath string) *dataFile {
	dc.dataFilesMu.Lock()
	defer dc.dataFilesMu.Unlock()
	df := dc.dataFiles[fqPath]
	if df != nil {
		df.users++
	}
	return df
}

func (dc *DiskCache) releaseDataFile(df *dataFile) {
	dc.dataFilesMu.Lock()
	defer dc.dataFilesMu.Unlock()
	df.users--
	if df.retired && df.users == 0 {
		_ = df.file.Close()
		close(df.closed)
	}
}

// closeDataFile closes an entry's kept-open data file. It waits until no
// operation still uses the file, because Windows cannot delete or replace a
// file that is open.
func (dc *DiskCache) closeDataFile(fqPath string) {
	dc.dataFilesMu.Lock()
	df := dc.dataFiles[fqPath]
	if df == nil {
		dc.dataFilesMu.Unlock()
		return
	}
	delete(dc.dataFiles, fqPath)
	df.retired = true
	idle := df.users == 0
	if idle {
		_ = df.file.Close()
		close(df.closed)
	}
	dc.dataFilesMu.Unlock()
	if !idle {
		<-df.closed
	}
}

func (dc *DiskCache) closeAllDataFiles() {
	dc.dataFilesMu.Lock()
	paths := make([]string, 0, len(dc.dataFiles))
	for fqPath := range dc.dataFiles {
		paths = append(paths, fqPath)
	}
	dc.dataFilesMu.Unlock()
	for _, fqPath := range paths {
		dc.closeDataFile(fqPath)
	}
}
