package fsmount

import (
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/stretchr/testify/require"
	"github.com/winfsp/cgofuse/fuse"
)

func TestLinuxDeletingOpenFileDeletesRemoteFileWithoutDownloading(t *testing.T) {
	remote, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()
	ignore, err := ignoreFromPatterns(nil)
	require.NoError(t, err)
	logger := lib.NewLeveledLogger(lib.NullLogger{})
	fs := &Filescomfs{remote: remote, local: newLocalFs(MountParams{TmpFsPath: t.TempDir()}, vfs, logger), vfs: vfs, log: logger, ignore: ignore}

	remote.disableLocking = false
	var locks, unlocks, deletes int
	remote.backend = &fakeRemoteBackend{
		createLockFunc: func(params files_sdk.LockCreateParams, opts ...files_sdk.RequestResponseOption) (files_sdk.Lock, error) {
			locks++
			return files_sdk.Lock{Path: params.Path, Token: "token"}, nil
		},
		deleteLockFunc: func(params files_sdk.LockDeleteParams, opts ...files_sdk.RequestResponseOption) error {
			unlocks++
			return nil
		},
		deleteFunc: func(params files_sdk.FileDeleteParams, opts ...files_sdk.RequestResponseOption) error {
			deletes++
			return nil
		},
		downloadToFileFunc: func(params files_sdk.FileDownloadParams, filePath string, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
			t.Fatalf("deleting an open file downloaded %v", params.Path)
			return files_sdk.File{}, nil
		},
	}
	remote.createNode("/report.xlsx", files_sdk.File{Path: "/report.xlsx", Type: "file", Size: 1024})

	errc, fh := fs.Open("/report.xlsx", fuse.O_RDWR)
	require.Zero(t, errc)

	// libfuse's sequence for unlinking a file while a handle is open.
	hidden := "/.fuse_hidden0000000200000001"
	require.Zero(t, fs.Rename("/report.xlsx", hidden))
	require.Equal(t, 1, deletes)
	require.Zero(t, fs.Release(hidden, fh))
	require.Zero(t, fs.Unlink(hidden))

	require.Equal(t, 1, locks)
	require.Equal(t, locks, unlocks, "the open handle's lock outlived it")
}
