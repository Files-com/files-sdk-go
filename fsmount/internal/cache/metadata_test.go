package cache

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestEntryMetadataValidate(t *testing.T) {
	modTime := time.Date(2026, time.August, 13, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name    string
		meta    EntryMetadata
		wantErr bool
	}{
		{name: "remote version without cached ranges", meta: NewEntryMetadata("/file", 10, modTime)},
		{name: "normalized ranges", meta: EntryMetadata{Version: EntryMetadataVersion, Path: "/file", Size: 10, ModTime: modTime, Ranges: []ByteRange{{Start: 0, End: 4}, {Start: 6, End: 10}}}},
		{name: "old schema", meta: EntryMetadata{Path: "/file", Size: 10, ModTime: modTime}, wantErr: true},
		{name: "empty remote path", meta: EntryMetadata{Version: EntryMetadataVersion, Size: 10, ModTime: modTime}, wantErr: true},
		{name: "negative size", meta: EntryMetadata{Version: EntryMetadataVersion, Path: "/file", Size: -1, ModTime: modTime}, wantErr: true},
		{name: "invalid range", meta: EntryMetadata{Version: EntryMetadataVersion, Path: "/file", Size: 10, ModTime: modTime, Ranges: []ByteRange{{Start: 5, End: 4}}}, wantErr: true},
		{name: "range past file size", meta: EntryMetadata{Version: EntryMetadataVersion, Path: "/file", Size: 10, ModTime: modTime, Ranges: []ByteRange{{Start: 0, End: 11}}}, wantErr: true},
		{name: "out of order ranges", meta: EntryMetadata{Version: EntryMetadataVersion, Path: "/file", Size: 10, ModTime: modTime, Ranges: []ByteRange{{Start: 6, End: 10}, {Start: 0, End: 4}}}, wantErr: true},
		{name: "adjacent unmerged ranges", meta: EntryMetadata{Version: EntryMetadataVersion, Path: "/file", Size: 10, ModTime: modTime, Ranges: []ByteRange{{Start: 0, End: 5}, {Start: 5, End: 10}}}, wantErr: true},
		{name: "empty range", meta: EntryMetadata{Version: EntryMetadataVersion, Path: "/file", Size: 10, ModTime: modTime, Ranges: []ByteRange{{Start: 4, End: 4}}}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.meta.Validate()
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidEntryMetadata) {
					t.Fatalf("Validate() error = %v, want ErrInvalidEntryMetadata", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Validate() error = %v", err)
			}
		})
	}
}

func TestEntryMetadataWithRangesNormalizes(t *testing.T) {
	meta := NewEntryMetadata("/file", 20, time.Now())
	ranges := []ByteRange{{Start: 10, End: 15}, {Start: 0, End: 5}, {Start: 4, End: 10}}
	original := append([]ByteRange(nil), ranges...)

	got, err := meta.WithRanges(ranges)
	if err != nil {
		t.Fatalf("WithRanges() error = %v", err)
	}
	want := []ByteRange{{Start: 0, End: 15}}
	if !reflect.DeepEqual(got.Ranges, want) {
		t.Errorf("WithRanges() ranges = %#v, want %#v", got.Ranges, want)
	}
	if !reflect.DeepEqual(ranges, original) {
		t.Errorf("WithRanges() modified input: got %#v, want %#v", ranges, original)
	}
}

func TestEntryMetadataCompleteness(t *testing.T) {
	modTime := time.Now()
	tests := []struct {
		name   string
		size   int64
		ranges []ByteRange
		want   bool
	}{
		{name: "empty file", want: true},
		{name: "complete", size: 10, ranges: []ByteRange{{Start: 0, End: 10}}, want: true},
		{name: "gap", size: 10, ranges: []ByteRange{{Start: 0, End: 4}, {Start: 6, End: 10}}},
		{name: "no ranges", size: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta, err := NewEntryMetadata("/file", tt.size, modTime).WithRanges(tt.ranges)
			if err != nil {
				t.Fatalf("WithRanges() error = %v", err)
			}
			if got := meta.IsComplete(); got != tt.want {
				t.Errorf("IsComplete() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestEntryMetadataMatchesRemoteVersion(t *testing.T) {
	modTime := time.Now().Round(0)
	stored, err := NewEntryMetadata("/file", 10, modTime).WithCompleteRange()
	if err != nil {
		t.Fatalf("WithCompleteRange() error = %v", err)
	}

	tests := []struct {
		name  string
		other EntryMetadata
		want  bool
	}{
		{name: "same remote version", other: NewEntryMetadata("/file", 10, modTime), want: true},
		{name: "different path", other: NewEntryMetadata("/other", 10, modTime)},
		{name: "different size", other: NewEntryMetadata("/file", 11, modTime)},
		{name: "different modification time", other: NewEntryMetadata("/file", 10, modTime.Add(time.Second))},
		{name: "old schema", other: EntryMetadata{Path: "/file", Size: 10, ModTime: modTime}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stored.Matches(tt.other); got != tt.want {
				t.Errorf("Matches() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestEntryMetadataJSONRoundTrip(t *testing.T) {
	meta, err := NewEntryMetadata("/file", 10, time.Now().UTC().Round(0)).WithCompleteRange()
	if err != nil {
		t.Fatalf("WithCompleteRange() error = %v", err)
	}

	data, err := json.Marshal(meta)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if strings.Contains(string(data), "complete") {
		t.Fatalf("metadata JSON contains obsolete completeness field: %s", data)
	}

	var decoded EntryMetadata
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if !reflect.DeepEqual(decoded, meta) {
		t.Errorf("round trip = %#v, want %#v", decoded, meta)
	}
	if err := decoded.Validate(); err != nil {
		t.Fatalf("decoded metadata is invalid: %v", err)
	}
}

func TestEntryMetadataCachedByteCount(t *testing.T) {
	meta, err := NewEntryMetadata("/file", 20, time.Now()).WithRanges([]ByteRange{
		{Start: 0, End: 5},
		{Start: 10, End: 17},
	})
	if err != nil {
		t.Fatalf("WithRanges() error = %v", err)
	}
	got, err := meta.CachedByteCount()
	if err != nil {
		t.Fatalf("CachedByteCount() error = %v", err)
	}
	if got != 12 {
		t.Errorf("CachedByteCount() = %d, want 12", got)
	}
}
