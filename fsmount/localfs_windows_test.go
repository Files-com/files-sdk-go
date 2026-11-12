package fsmount

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/stretchr/testify/require"
	"github.com/winfsp/cgofuse/fuse"
)

func TestWindowsLocalGetattrReportsLocalDirectoriesAsOwnerOnlyDirectories(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(root, "$RECYCLE.BIN"), 0o700))
	vfs := &virtualfs{nodes: make(map[string]*fsNode)}
	fs := &LocalFs{localFsRoot: root, vfs: vfs, log: lib.NewLeveledLogger(lib.NullLogger{})}
	var stat fuse.Stat_t
	require.Zero(t, fs.Getattr("/$RECYCLE.BIN", &stat, 0))
	require.EqualValues(t, fuse.S_IFDIR, stat.Mode&fuse.S_IFMT)
	require.Zero(t, stat.Mode&0o077, "group and world permissions must stay cleared")
}
