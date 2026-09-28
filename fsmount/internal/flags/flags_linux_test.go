package flags_test

import (
	"os"
	"testing"

	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/flags"
	"github.com/winfsp/cgofuse/fuse"
)

func TestLinuxReadOnlyWithStatusFlags(t *testing.T) {
	for _, test := range []struct {
		name     string
		mode     int
		readOnly bool
	}{
		{"read only", fuse.O_RDONLY, true},
		{"write only", fuse.O_WRONLY, false},
		{"read write", fuse.O_RDWR, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := flags.NewFuseFlags(test.mode | os.O_SYNC).IsReadOnly()
			if got != test.readOnly {
				t.Fatalf("IsReadOnly() = %v, want %v", got, test.readOnly)
			}
		})
	}
}
