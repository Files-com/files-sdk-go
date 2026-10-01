package cache

import (
	"errors"
	"reflect"
	"testing"
)

func TestByteRangeValidate(t *testing.T) {
	tests := []struct {
		name    string
		r       ByteRange
		wantErr bool
	}{
		{name: "non-empty", r: ByteRange{Start: 3, End: 8}},
		{name: "empty", r: ByteRange{Start: 3, End: 3}},
		{name: "negative start", r: ByteRange{Start: -1, End: 3}, wantErr: true},
		{name: "negative end", r: ByteRange{Start: -2, End: -1}, wantErr: true},
		{name: "reversed", r: ByteRange{Start: 8, End: 3}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.r.Validate()
			if tt.wantErr && !errors.Is(err, ErrInvalidByteRange) {
				t.Fatalf("Validate() error = %v, want ErrInvalidByteRange", err)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("Validate() error = %v", err)
			}
		})
	}
}

func TestNormalizeByteRanges(t *testing.T) {
	tests := []struct {
		name    string
		ranges  []ByteRange
		want    []ByteRange
		wantErr bool
	}{
		{name: "nil"},
		{
			name:   "sort merge and omit empty",
			ranges: []ByteRange{{Start: 20, End: 30}, {Start: 0, End: 5}, {Start: 4, End: 10}, {Start: 10, End: 12}, {Start: 7, End: 8}, {Start: 15, End: 15}},
			want:   []ByteRange{{Start: 0, End: 12}, {Start: 20, End: 30}},
		},
		{
			name:   "same start and duplicate",
			ranges: []ByteRange{{Start: 2, End: 4}, {Start: 2, End: 9}, {Start: 2, End: 4}},
			want:   []ByteRange{{Start: 2, End: 9}},
		},
		{
			name:    "invalid",
			ranges:  []ByteRange{{Start: 4, End: 3}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := append([]ByteRange(nil), tt.ranges...)
			got, err := NormalizeByteRanges(tt.ranges)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidByteRange) {
					t.Fatalf("NormalizeByteRanges() error = %v, want ErrInvalidByteRange", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("NormalizeByteRanges() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NormalizeByteRanges() = %#v, want %#v", got, tt.want)
			}
			if !reflect.DeepEqual(tt.ranges, original) {
				t.Errorf("NormalizeByteRanges() modified input: got %#v, want %#v", tt.ranges, original)
			}
		})
	}
}

func TestMissingByteRanges(t *testing.T) {
	tests := []struct {
		name      string
		ranges    []ByteRange
		requested ByteRange
		want      []ByteRange
	}{
		{
			name:      "nothing cached",
			requested: ByteRange{Start: 10, End: 20},
			want:      []ByteRange{{Start: 10, End: 20}},
		},
		{
			name:      "covered by out of order ranges",
			ranges:    []ByteRange{{Start: 15, End: 25}, {Start: 5, End: 15}},
			requested: ByteRange{Start: 10, End: 20},
		},
		{
			name:      "clip gaps to request",
			ranges:    []ByteRange{{Start: 0, End: 4}, {Start: 8, End: 12}, {Start: 16, End: 30}},
			requested: ByteRange{Start: 2, End: 20},
			want:      []ByteRange{{Start: 4, End: 8}, {Start: 12, End: 16}},
		},
		{
			name:      "leading and trailing gaps",
			ranges:    []ByteRange{{Start: 12, End: 18}},
			requested: ByteRange{Start: 10, End: 20},
			want:      []ByteRange{{Start: 10, End: 12}, {Start: 18, End: 20}},
		},
		{
			name:      "empty request",
			ranges:    []ByteRange{{Start: 0, End: 10}},
			requested: ByteRange{Start: 5, End: 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MissingByteRanges(tt.ranges, tt.requested)
			if err != nil {
				t.Fatalf("MissingByteRanges() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MissingByteRanges() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestMissingByteRangesRejectsInvalidRange(t *testing.T) {
	tests := []struct {
		name      string
		ranges    []ByteRange
		requested ByteRange
	}{
		{name: "invalid cached range", ranges: []ByteRange{{Start: 5, End: 4}}, requested: ByteRange{Start: 0, End: 10}},
		{name: "invalid cached range with empty request", ranges: []ByteRange{{Start: 5, End: 4}}, requested: ByteRange{Start: 0, End: 0}},
		{name: "invalid requested range", requested: ByteRange{Start: -1, End: 10}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := MissingByteRanges(tt.ranges, tt.requested)
			if !errors.Is(err, ErrInvalidByteRange) {
				t.Fatalf("MissingByteRanges() error = %v, want ErrInvalidByteRange", err)
			}
		})
	}
}

func TestByteRangeCovered(t *testing.T) {
	tests := []struct {
		name      string
		ranges    []ByteRange
		requested ByteRange
		want      bool
	}{
		{name: "covered", ranges: []ByteRange{{Start: 0, End: 5}, {Start: 5, End: 10}}, requested: ByteRange{Start: 2, End: 8}, want: true},
		{name: "gap", ranges: []ByteRange{{Start: 0, End: 4}, {Start: 5, End: 10}}, requested: ByteRange{Start: 2, End: 8}},
		{name: "empty", requested: ByteRange{Start: 2, End: 2}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ByteRangeCovered(tt.ranges, tt.requested)
			if err != nil {
				t.Fatalf("ByteRangeCovered() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("ByteRangeCovered() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestByteRangesCoverFile(t *testing.T) {
	tests := []struct {
		name     string
		ranges   []ByteRange
		fileSize int64
		want     bool
		wantErr  bool
	}{
		{name: "empty file", fileSize: 0, want: true},
		{name: "complete", ranges: []ByteRange{{Start: 5, End: 12}, {Start: 0, End: 5}}, fileSize: 10, want: true},
		{name: "gap", ranges: []ByteRange{{Start: 0, End: 5}, {Start: 6, End: 10}}, fileSize: 10},
		{name: "negative size", fileSize: -1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ByteRangesCoverFile(tt.ranges, tt.fileSize)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidByteRange) {
					t.Fatalf("ByteRangesCoverFile() error = %v, want ErrInvalidByteRange", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("ByteRangesCoverFile() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("ByteRangesCoverFile() = %t, want %t", got, tt.want)
			}
		})
	}
}
