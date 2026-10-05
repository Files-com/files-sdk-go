//go:build linux || windows

package fsmount

import (
	"testing"

	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/stretchr/testify/require"
	"github.com/winfsp/cgofuse/fuse"
)

func TestStatfsReportsMaxNameLength(t *testing.T) {
	fs := &Filescomfs{log: lib.NewLeveledLogger(lib.NullLogger{})}

	var stat fuse.Statfs_t
	require.Equal(t, 0, fs.Statfs("/", &stat))
	require.Equal(t, uint64(maxNameLength), stat.Namemax)
}
