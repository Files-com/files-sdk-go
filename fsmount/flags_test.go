package fsmount

import (
	"testing"

	"github.com/winfsp/cgofuse/fuse"
)

func TestFuseFlags(t *testing.T) {
	tests := []struct {
		name                    string
		flags                   int
		expectedReadOnly        bool
		expectedWriteOnly       bool
		expectedReadWrite       bool
		expectedCreate          bool
		expectedExclusive       bool
		expectedTruncate        bool
		expectedAppend          bool
		expectedCreateExclusive bool
		expectedString          string
	}{
		{
			name:                    "ReadOnly",
			flags:                   fuse.O_RDONLY,
			expectedReadOnly:        true,
			expectedWriteOnly:       false,
			expectedReadWrite:       false,
			expectedCreate:          false,
			expectedExclusive:       false,
			expectedTruncate:        false,
			expectedAppend:          false,
			expectedCreateExclusive: false,
			expectedString:          "FuseFlags{O_RDONLY}",
		},
		{
			name:                    "WriteOnly",
			flags:                   fuse.O_WRONLY,
			expectedReadOnly:        false,
			expectedWriteOnly:       true,
			expectedReadWrite:       false,
			expectedCreate:          false,
			expectedExclusive:       false,
			expectedTruncate:        false,
			expectedAppend:          false,
			expectedCreateExclusive: false,
			expectedString:          "FuseFlags{O_WRONLY}",
		},
		{
			name:                    "ReadWrite",
			flags:                   fuse.O_RDWR,
			expectedReadOnly:        false,
			expectedWriteOnly:       false,
			expectedReadWrite:       true,
			expectedCreate:          false,
			expectedExclusive:       false,
			expectedTruncate:        false,
			expectedAppend:          false,
			expectedCreateExclusive: false,
			expectedString:          "FuseFlags{O_RDWR}",
		},
		{
			name:                    "CreateExclusive",
			flags:                   fuse.O_CREAT | fuse.O_EXCL,
			expectedReadOnly:        false,
			expectedWriteOnly:       false,
			expectedReadWrite:       false,
			expectedCreate:          true,
			expectedExclusive:       true,
			expectedTruncate:        false,
			expectedAppend:          false,
			expectedCreateExclusive: true,
			expectedString:          "FuseFlags{O_CREAT|O_EXCL}",
		},
		{
			name:                    "TruncateAppend",
			flags:                   fuse.O_TRUNC | fuse.O_APPEND,
			expectedReadOnly:        false,
			expectedWriteOnly:       false,
			expectedReadWrite:       false,
			expectedCreate:          false,
			expectedExclusive:       false,
			expectedTruncate:        true,
			expectedAppend:          true,
			expectedCreateExclusive: false,
			expectedString:          "FuseFlags{O_TRUNC|O_APPEND}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ff := NewFuseFlags(tt.flags)

			if got := ff.IsReadOnly(); got != tt.expectedReadOnly {
				t.Errorf("IsReadOnly() = %v, want %v", got, tt.expectedReadOnly)
			}
			if got := ff.IsWriteOnly(); got != tt.expectedWriteOnly {
				t.Errorf("IsWriteOnly() = %v, want %v", got, tt.expectedWriteOnly)
			}
			if got := ff.IsReadWrite(); got != tt.expectedReadWrite {
				t.Errorf("IsReadWrite() = %v, want %v", got, tt.expectedReadWrite)
			}
			if got := ff.IsCreate(); got != tt.expectedCreate {
				t.Errorf("IsCreate() = %v, want %v", got, tt.expectedCreate)
			}
			if got := ff.IsExclusive(); got != tt.expectedExclusive {
				t.Errorf("IsExclusive() = %v, want %v", got, tt.expectedExclusive)
			}
			if got := ff.IsTruncate(); got != tt.expectedTruncate {
				t.Errorf("IsTruncate() = %v, want %v", got, tt.expectedTruncate)
			}
			if got := ff.IsAppend(); got != tt.expectedAppend {
				t.Errorf("IsAppend() = %v, want %v", got, tt.expectedAppend)
			}
			if got := ff.IsCreateExclusive(); got != tt.expectedCreateExclusive {
				t.Errorf("IsCreateExclusive() = %v, want %v", got, tt.expectedCreateExclusive)
			}
			if got := ff.String(); got != tt.expectedString {
				t.Errorf("String() = %v, want %v", got, tt.expectedString)
			}
		})
	}
}
