//go:build linux || windows

package fsmount

import (
	"errors"
	"slices"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/winfsp/cgofuse/fuse"
)

func refreshDirectory(t *testing.T, fs *RemoteFs, path string) []string {
	t.Helper()
	node, ok := fs.vfs.fetch(path)
	if !ok {
		t.Fatalf("directory %q is not cached", path)
	}
	node.expireInfo()
	node.childPathsExpires = timeZero
	var names []string
	if errc := fs.Readdir(path, func(name string, _ *fuse.Stat_t, _ int64) bool {
		if name != "." && name != ".." {
			names = append(names, name)
		}
		return true
	}, 0, ^uint64(0)); errc != 0 {
		t.Fatalf("Readdir(%q) returned %d", path, errc)
	}
	return names
}

func TestRemoteRefreshAllowsReusingMovedNames(t *testing.T) {
	for _, kind := range []string{"file", "directory"} {
		for _, operation := range []string{"create", "rename"} {
			t.Run(kind+"/"+operation, func(t *testing.T) {
				fs, vfs, _ := newTestRemoteFs(t)
				defer vfs.destroy()
				items := []files_sdk.File{
					{DisplayName: "vacated", Type: kind},
					{DisplayName: "remaining", Type: kind},
				}
				backend := fs.backend.(*fakeRemoteBackend)
				backend.listForFunc = func(_ files_sdk.FolderListForParams, _ ...files_sdk.RequestResponseOption) (remoteFileIter, error) {
					return &fakeFileIter{files: items}, nil
				}
				if names := refreshDirectory(t, fs, "/"); !slices.Equal(names, []string{"remaining", "vacated"}) {
					t.Fatalf("initial listing = %v", names)
				}

				// Another client moves the first item away. Its cached metadata may
				// still be fresh when Explorer refreshes the containing directory.
				items = items[1:]
				if names := refreshDirectory(t, fs, "/"); !slices.Equal(names, []string{"remaining"}) {
					t.Fatalf("refreshed listing = %v", names)
				}
				var stat fuse.Stat_t
				if errc := fs.Getattr("/vacated", &stat, ^uint64(0)); errc != -fuse.ENOENT {
					t.Fatalf("Getattr removed name = %d, want ENOENT", errc)
				}
				if errc, _ := fs.Open("/vacated", fuse.O_RDONLY); errc != -fuse.ENOENT {
					t.Fatalf("Open removed name = %d, want ENOENT", errc)
				}

				if operation == "rename" {
					if errc := fs.Rename("/remaining", "/vacated"); errc != 0 {
						t.Fatalf("Rename to removed name = %d", errc)
					}
				} else if kind == "directory" {
					if errc := fs.Mkdir("/vacated", 0o755); errc != 0 {
						t.Fatalf("Mkdir at removed name = %d", errc)
					}
				} else {
					errc, fh := fs.Create("/vacated", fuse.O_RDWR|fuse.O_CREAT|fuse.O_EXCL, 0o644)
					if errc != 0 {
						t.Fatalf("Create at removed name = %d", errc)
					}
					if errc := fs.Release("/vacated", fh); errc != 0 {
						t.Fatalf("Release recreated file = %d", errc)
					}
					// Reusing the name must create a local placeholder, so an immediate
					// rename still works without a nonexistent remote source.
					backend.moveFunc = func(_ files_sdk.FileMoveParams, _ ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error) {
						t.Fatal("empty recreated file must be renamed locally")
						return files_sdk.FileAction{}, nil
					}
					if errc := fs.Rename("/vacated", "/new-name"); errc != 0 {
						t.Fatalf("Rename recreated empty file = %d", errc)
					}
				}
			})
		}
	}
}

func TestRemoteRefreshPreservesOpenFileUntilReleased(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()
	fs.localFsRoot = "/.local"
	items := []files_sdk.File{{DisplayName: "open.txt", Type: "file"}}
	fs.backend.(*fakeRemoteBackend).listForFunc = func(_ files_sdk.FolderListForParams, _ ...files_sdk.RequestResponseOption) (remoteFileIter, error) {
		return &fakeFileIter{files: items}, nil
	}
	refreshDirectory(t, fs, "/")
	errc, fh := fs.Open("/open.txt", fuse.O_RDONLY)
	if errc != 0 {
		t.Fatalf("Open = %d", errc)
	}
	items = nil
	if names := refreshDirectory(t, fs, "/"); !slices.Equal(names, []string{"open.txt"}) {
		t.Fatalf("open file disappeared: %v", names)
	}
	if errc := fs.Release("/open.txt", fh); errc != 0 {
		t.Fatalf("Release = %d", errc)
	}
	if names := refreshDirectory(t, fs, "/"); len(names) != 0 {
		t.Fatalf("released file still listed: %v", names)
	}
	var stat fuse.Stat_t
	if errc := fs.Getattr("/open.txt", &stat, ^uint64(0)); errc != -fuse.ENOENT {
		t.Fatalf("released remote-removed file Getattr = %d, want ENOENT", errc)
	}
}

type failedFileIter struct {
	fakeFileIter
}

func (*failedFileIter) Err() error { return errors.New("listing interrupted") }

func TestRemoteRefreshDoesNotForgetNamesAfterPartialListing(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()
	backend := fs.backend.(*fakeRemoteBackend)
	backend.listForFunc = func(_ files_sdk.FolderListForParams, _ ...files_sdk.RequestResponseOption) (remoteFileIter, error) {
		return &fakeFileIter{files: []files_sdk.File{{DisplayName: "retained.txt", Type: "file"}}}, nil
	}
	refreshDirectory(t, fs, "/")
	backend.listForFunc = func(_ files_sdk.FolderListForParams, _ ...files_sdk.RequestResponseOption) (remoteFileIter, error) {
		return &failedFileIter{fakeFileIter{files: []files_sdk.File{{DisplayName: "other.txt", Type: "file"}}}}, nil
	}
	root, _ := vfs.fetch("/")
	root.expireInfo()
	root.childPathsExpires = timeZero
	if errc := fs.Readdir("/", func(string, *fuse.Stat_t, int64) bool { return true }, 0, ^uint64(0)); errc == 0 {
		t.Fatal("partial listing unexpectedly succeeded")
	}
	var stat fuse.Stat_t
	if errc := fs.Getattr("/retained.txt", &stat, ^uint64(0)); errc != 0 {
		t.Fatalf("previously cached file lost after failed listing: %d", errc)
	}
}
