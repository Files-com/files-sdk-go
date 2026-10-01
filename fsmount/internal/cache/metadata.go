package cache

import (
	"errors"
	"fmt"
	"time"
)

const EntryMetadataVersion = 2

var ErrInvalidEntryMetadata = errors.New("invalid cache entry metadata")

// EntryMetadata identifies cached ranges for a specific remote file version.
type EntryMetadata struct {
	Version int         `json:"version"`
	Path    string      `json:"path"`
	Size    int64       `json:"size"`
	ModTime time.Time   `json:"mtime"`
	ETag    string      `json:"etag,omitempty"`
	Ranges  []ByteRange `json:"ranges"`
}

func NewEntryMetadata(path string, size int64, modTime time.Time) EntryMetadata {
	return EntryMetadata{
		Version: EntryMetadataVersion,
		Path:    path,
		Size:    size,
		ModTime: modTime,
	}
}

// Validate checks the schema version, file size, and normalized range list.
func (m EntryMetadata) Validate() error {
	if m.Version != EntryMetadataVersion {
		return fmt.Errorf("%w: schema version %d", ErrInvalidEntryMetadata, m.Version)
	}
	if m.Path == "" {
		return fmt.Errorf("%w: empty remote path", ErrInvalidEntryMetadata)
	}
	if m.Size < 0 {
		return fmt.Errorf("%w: negative file size %d", ErrInvalidEntryMetadata, m.Size)
	}

	normalized, err := NormalizeByteRanges(m.Ranges)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidEntryMetadata, err)
	}
	if !equalByteRanges(m.Ranges, normalized) {
		return fmt.Errorf("%w: ranges are not normalized", ErrInvalidEntryMetadata)
	}
	for _, r := range normalized {
		if r.End > m.Size {
			return fmt.Errorf("%w: range [%d, %d) exceeds file size %d", ErrInvalidEntryMetadata, r.Start, r.End, m.Size)
		}
	}
	return nil
}

// WithRanges returns metadata with a validated, normalized range list.
func (m EntryMetadata) WithRanges(ranges []ByteRange) (EntryMetadata, error) {
	normalized, err := NormalizeByteRanges(ranges)
	if err != nil {
		return EntryMetadata{}, fmt.Errorf("%w: %v", ErrInvalidEntryMetadata, err)
	}
	m.Ranges = normalized
	if err := m.Validate(); err != nil {
		return EntryMetadata{}, err
	}
	return m, nil
}

// WithCompleteRange returns metadata covering the complete file.
func (m EntryMetadata) WithCompleteRange() (EntryMetadata, error) {
	if m.Size == 0 {
		return m.WithRanges(nil)
	}
	return m.WithRanges([]ByteRange{{Start: 0, End: m.Size}})
}

// Matches compares the listing identity. HTTP validators are checked separately
// because callers looking up cached bytes have only listing metadata.
func (m EntryMetadata) Matches(other EntryMetadata) bool {
	return m.Validate() == nil &&
		other.Validate() == nil &&
		m.Path == other.Path &&
		m.Size == other.Size &&
		m.ModTime.Equal(other.ModTime)
}

// SameContent also requires the HTTP validator to match before merging bytes.
func (m EntryMetadata) SameContent(other EntryMetadata) bool {
	return m.Matches(other) && m.ETag == other.ETag
}

// IsComplete reports whether the valid ranges cover the complete file.
func (m EntryMetadata) IsComplete() bool {
	if err := m.Validate(); err != nil {
		return false
	}
	complete, err := ByteRangesCoverFile(m.Ranges, m.Size)
	return err == nil && complete
}

// CachedByteCount returns the number of valid bytes represented by the ranges.
func (m EntryMetadata) CachedByteCount() (int64, error) {
	if err := m.Validate(); err != nil {
		return 0, err
	}
	var count int64
	for _, r := range m.Ranges {
		count += r.End - r.Start
	}
	return count, nil
}

func equalByteRanges(left, right []ByteRange) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}
