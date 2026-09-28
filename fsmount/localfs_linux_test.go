package fsmount

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/stretchr/testify/require"
	"github.com/winfsp/cgofuse/fuse"
)

func TestLinuxLocalGetattrUsesCachePath(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, ".~lock.document#"), []byte("owner"), 0600))
	fs := &LocalFs{localFsRoot: root, log: lib.NewLeveledLogger(lib.NullLogger{})}
	var stat fuse.Stat_t
	require.Zero(t, fs.Getattr("/.~lock.document#", &stat, 0))
	require.EqualValues(t, 5, stat.Size)
	require.EqualValues(t, fuse.S_IFREG, stat.Mode&fuse.S_IFMT)
}
