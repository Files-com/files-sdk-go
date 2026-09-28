package fsmount

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/stretchr/testify/require"
)

type directoryMountBackend struct{ root string }

func (b directoryMountBackend) path(path string) string {
	return filepath.Join(b.root, filepath.FromSlash(path))
}

func (b directoryMountBackend) Stat(_ context.Context, path string) (ProviderEntry, error) {
	info, err := os.Stat(b.path(path))
	if err != nil {
		return ProviderEntry{}, err
	}
	kind := ProviderTypeFile
	if info.IsDir() {
		kind = ProviderTypeDirectory
	}
	return ProviderEntry{Path: path, DisplayName: info.Name(), Type: kind, Size: info.Size(), ModTime: info.ModTime(), Permissions: "read,write,delete"}, nil
}

func (b directoryMountBackend) List(ctx context.Context, path string) ([]ProviderEntry, error) {
	entries, err := os.ReadDir(b.path(path))
	if err != nil {
		return nil, err
	}
	result := make([]ProviderEntry, 0, len(entries))
	for _, entry := range entries {
		item, err := b.Stat(ctx, filepath.ToSlash(filepath.Join(path, entry.Name())))
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}

func (b directoryMountBackend) Read(ctx context.Context, path string) (io.ReadCloser, int64, error) {
	entry, err := b.Stat(ctx, path)
	if err != nil {
		return nil, 0, err
	}
	file, err := os.Open(b.path(path))
	return file, entry.Size, err
}

func (b directoryMountBackend) Write(ctx context.Context, path string, reader io.Reader, _ int64, _ time.Time) (ProviderEntry, error) {
	file, err := os.Create(b.path(path))
	if err != nil {
		return ProviderEntry{}, err
	}
	_, err = io.Copy(file, reader)
	closeErr := file.Close()
	if err != nil {
		return ProviderEntry{}, err
	}
	if closeErr != nil {
		return ProviderEntry{}, closeErr
	}
	return b.Stat(ctx, path)
}

func (b directoryMountBackend) Mkdir(_ context.Context, path string) error {
	return os.Mkdir(b.path(path), 0700)
}

func (b directoryMountBackend) Delete(_ context.Context, path string, recursive bool) error {
	if recursive {
		return os.RemoveAll(b.path(path))
	}
	return os.Remove(b.path(path))
}

func (b directoryMountBackend) Rename(_ context.Context, source, destination string) error {
	return os.Rename(b.path(source), b.path(destination))
}

func TestLinuxMountedFileOperations(t *testing.T) {
	if os.Getenv("FILESCOMFS_TEST_MOUNT") != "1" {
		t.Skip("set FILESCOMFS_TEST_MOUNT=1 with /dev/fuse available")
	}
	for _, disk := range []bool{false, true} {
		name := "memory"
		if disk {
			name = "disk"
		}
		t.Run(name, func(t *testing.T) {
			backend := directoryMountBackend{root: t.TempDir()}
			mount := t.TempDir()
			host, err := Mount(MountParams{Config: &files_sdk.Config{Logger: lib.NullLogger{}}, ProviderBackend: backend, MountPoint: mount, TmpFsPath: t.TempDir(), DiskCacheEnabled: disk, DiskCachePath: t.TempDir()})
			require.NoError(t, err)
			t.Cleanup(func() { require.True(t, host.Unmount()) })
			require.NoError(t, os.Mkdir(filepath.Join(mount, "documents"), 0700))
			file := filepath.Join(mount, "documents", "test.txt")
			for _, content := range []string{"first version", "edited version"} {
				require.NoError(t, os.WriteFile(file, []byte(content), 0600))
				data, err := os.ReadFile(file)
				require.NoError(t, err)
				require.Equal(t, content, string(data))
				require.Eventually(t, func() bool {
					data, err := os.ReadFile(filepath.Join(backend.root, "documents", "test.txt"))
					return err == nil && string(data) == content
				}, 5*time.Second, 10*time.Millisecond)
			}
			renamed := filepath.Join(mount, "documents", "renamed.txt")
			require.NoError(t, os.Rename(file, renamed))
			require.FileExists(t, filepath.Join(backend.root, "documents", "renamed.txt"))
			open, err := os.Open(renamed)
			require.NoError(t, err)
			require.NoError(t, os.Remove(renamed))
			require.NoFileExists(t, filepath.Join(backend.root, "documents", "renamed.txt"))
			require.NoError(t, open.Close())
			require.NoError(t, os.Remove(filepath.Join(mount, "documents")))
			require.NoDirExists(t, filepath.Join(backend.root, "documents"))
			lock := filepath.Join(mount, ".~lock.document#")
			require.NoError(t, os.WriteFile(lock, []byte("owner"), 0600))
			data, err := os.ReadFile(lock)
			require.NoError(t, err)
			require.Equal(t, "owner", string(data))
			require.NoFileExists(t, filepath.Join(backend.root, ".~lock.document#"))
			require.NoError(t, os.Remove(lock))
		})
	}
}
