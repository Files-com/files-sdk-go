package cache

import (
	"errors"
	"fmt"
	"sort"
)

var (
	ErrInvalidByteRange       = errors.New("invalid byte range")
	ErrSparseFilesUnsupported = errors.New("sparse files are not supported by the cache volume")
)

// ByteRange identifies a half-open range of bytes: [Start, End).
type ByteRange struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

// Validate checks that the range has non-negative, ordered offsets.
func (r ByteRange) Validate() error {
	if r.Start < 0 || r.End < r.Start {
		return fmt.Errorf("%w: [%d, %d)", ErrInvalidByteRange, r.Start, r.End)
	}
	return nil
}

// NormalizeByteRanges sorts ranges and merges ranges that overlap or touch.
// Empty ranges are omitted. The input slice is not modified.
func NormalizeByteRanges(ranges []ByteRange) ([]ByteRange, error) {
	if len(ranges) == 0 {
		return nil, nil
	}

	normalized := make([]ByteRange, 0, len(ranges))
	for _, r := range ranges {
		if err := r.Validate(); err != nil {
			return nil, err
		}
		if r.Start != r.End {
			normalized = append(normalized, r)
		}
	}
	if len(normalized) == 0 {
		return nil, nil
	}

	sort.Slice(normalized, func(i, j int) bool {
		if normalized[i].Start == normalized[j].Start {
			return normalized[i].End < normalized[j].End
		}
		return normalized[i].Start < normalized[j].Start
	})

	merged := normalized[:1]
	for _, next := range normalized[1:] {
		last := &merged[len(merged)-1]
		if next.Start <= last.End {
			if next.End > last.End {
				last.End = next.End
			}
			continue
		}
		merged = append(merged, next)
	}
	return merged, nil
}

// MissingByteRanges returns the parts of requested that ranges do not cover.
func MissingByteRanges(ranges []ByteRange, requested ByteRange) ([]ByteRange, error) {
	if err := requested.Validate(); err != nil {
		return nil, err
	}

	normalized, err := NormalizeByteRanges(ranges)
	if err != nil {
		return nil, err
	}
	if requested.Start == requested.End {
		return nil, nil
	}

	var missing []ByteRange
	cursor := requested.Start
	for _, cached := range normalized {
		if cached.End <= cursor {
			continue
		}
		if cached.Start >= requested.End {
			break
		}
		if cached.Start > cursor {
			end := min(cached.Start, requested.End)
			missing = append(missing, ByteRange{Start: cursor, End: end})
		}
		if cached.End > cursor {
			cursor = min(cached.End, requested.End)
		}
		if cursor == requested.End {
			return missing, nil
		}
	}
	if cursor < requested.End {
		missing = append(missing, ByteRange{Start: cursor, End: requested.End})
	}
	return missing, nil
}

// ByteRangeCovered reports whether ranges completely cover requested.
func ByteRangeCovered(ranges []ByteRange, requested ByteRange) (bool, error) {
	missing, err := MissingByteRanges(ranges, requested)
	if err != nil {
		return false, err
	}
	return len(missing) == 0, nil
}

// ByteRangesCoverFile reports whether ranges cover every byte in a file.
func ByteRangesCoverFile(ranges []ByteRange, fileSize int64) (bool, error) {
	if fileSize < 0 {
		return false, fmt.Errorf("%w: negative file size %d", ErrInvalidByteRange, fileSize)
	}
	return ByteRangeCovered(ranges, ByteRange{Start: 0, End: fileSize})
}
