package cache

import "time"

// EntryMetadata identifies a completed cache entry for a specific remote file version.
type EntryMetadata struct {
	Path     string    `json:"path"`
	Size     int64     `json:"size"`
	ModTime  time.Time `json:"mtime"`
	Complete bool      `json:"complete"`
}

func NewEntryMetadata(path string, size int64, modTime time.Time) EntryMetadata {
	return EntryMetadata{
		Path:     path,
		Size:     size,
		ModTime:  modTime,
		Complete: true,
	}
}

func (m EntryMetadata) Matches(other EntryMetadata) bool {
	return m.Complete &&
		other.Complete &&
		m.Path == other.Path &&
		m.Size == other.Size &&
		m.ModTime.Equal(other.ModTime)
}
