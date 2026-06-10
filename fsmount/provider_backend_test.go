package fsmount

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/log"
)

type testProviderBackend struct {
	statFunc   func(context.Context, string) (ProviderEntry, error)
	listFunc   func(context.Context, string) ([]ProviderEntry, error)
	readFunc   func(context.Context, string) (io.ReadCloser, int64, error)
	writeFunc  func(context.Context, string, io.Reader, int64, time.Time) (ProviderEntry, error)
	mkdirFunc  func(context.Context, string) error
	deleteFunc func(context.Context, string, bool) error
	renameFunc func(context.Context, string, string) error
}

func (b *testProviderBackend) Stat(ctx context.Context, path string) (ProviderEntry, error) {
	if b.statFunc != nil {
		return b.statFunc(ctx, path)
	}
	return ProviderEntry{}, nil
}

func (b *testProviderBackend) List(ctx context.Context, path string) ([]ProviderEntry, error) {
	if b.listFunc != nil {
		return b.listFunc(ctx, path)
	}
	return nil, nil
}

func (b *testProviderBackend) Read(ctx context.Context, path string) (io.ReadCloser, int64, error) {
	if b.readFunc != nil {
		return b.readFunc(ctx, path)
	}
	return io.NopCloser(strings.NewReader("")), 0, nil
}

func (b *testProviderBackend) Write(ctx context.Context, path string, reader io.Reader, size int64, modTime time.Time) (ProviderEntry, error) {
	if b.writeFunc != nil {
		return b.writeFunc(ctx, path, reader, size, modTime)
	}
	return ProviderEntry{}, nil
}

func (b *testProviderBackend) Mkdir(ctx context.Context, path string) error {
	if b.mkdirFunc != nil {
		return b.mkdirFunc(ctx, path)
	}
	return nil
}

func (b *testProviderBackend) Delete(ctx context.Context, path string, recursive bool) error {
	if b.deleteFunc != nil {
		return b.deleteFunc(ctx, path, recursive)
	}
	return nil
}

func (b *testProviderBackend) Rename(ctx context.Context, sourcePath string, destinationPath string) error {
	if b.renameFunc != nil {
		return b.renameFunc(ctx, sourcePath, destinationPath)
	}
	return nil
}

type statErrorUploadReader struct {
	*strings.Reader
	err error
}

func (r statErrorUploadReader) Stat() (os.FileInfo, error) {
	return nil, r.err
}

type testMtimeProviderBackend struct {
	*testProviderBackend
	mtimeFunc func(context.Context, string, time.Time) (ProviderEntry, error)
}

func (b *testMtimeProviderBackend) SetMtime(ctx context.Context, path string, modTime time.Time) (ProviderEntry, error) {
	if b.mtimeFunc != nil {
		return b.mtimeFunc(ctx, path, modTime)
	}
	return ProviderEntry{}, nil
}

func testUploadedMetadata(size int64, modTime time.Time) uploadedFileMetadata {
	return uploadedFileMetadata{
		size:    size,
		modTime: modTime,
	}
}

func TestProviderRemoteBackendAdaptsMetadataOperations(t *testing.T) {
	modTime := time.Date(2026, 5, 29, 10, 30, 0, 0, time.UTC)
	var mkdirPath string
	var deletePath string
	var deleteRecursive bool
	var renameSource string
	var renameDestination string

	provider := &testProviderBackend{
		statFunc: func(_ context.Context, path string) (ProviderEntry, error) {
			switch path {
			case "/docs/report.txt":
				return ProviderEntry{
					Path:    path,
					Type:    ProviderTypeFile,
					Size:    123,
					ModTime: modTime,
				}, nil
			case "/new-folder":
				return ProviderEntry{
					Path:    path,
					Type:    ProviderTypeDirectory,
					ModTime: modTime,
				}, nil
			default:
				return ProviderEntry{}, os.ErrNotExist
			}
		},
		listFunc: func(_ context.Context, path string) ([]ProviderEntry, error) {
			if path != "/docs" {
				t.Fatalf("List path = %q, want /docs", path)
			}
			return []ProviderEntry{
				{Path: "/docs/report.txt", Type: ProviderTypeFile, Size: 123},
				{Path: "/docs/archive", Type: ProviderTypeDirectory},
			}, nil
		},
		mkdirFunc: func(_ context.Context, path string) error {
			mkdirPath = path
			return nil
		},
		deleteFunc: func(_ context.Context, path string, recursive bool) error {
			deletePath = path
			deleteRecursive = recursive
			return nil
		},
		renameFunc: func(_ context.Context, sourcePath string, destinationPath string) error {
			renameSource = sourcePath
			renameDestination = destinationPath
			return nil
		},
	}
	backend := &providerRemoteBackend{provider: provider}

	file, err := backend.find(files_sdk.FileFindParams{Path: "docs/report.txt"})
	if err != nil {
		t.Fatalf("find failed: %v", err)
	}
	if file.Path != "/docs/report.txt" {
		t.Fatalf("find path = %q, want /docs/report.txt", file.Path)
	}
	if file.DisplayName != "report.txt" {
		t.Fatalf("find display name = %q, want report.txt", file.DisplayName)
	}
	if file.Size != 123 {
		t.Fatalf("find size = %d, want 123", file.Size)
	}
	if file.Mtime == nil || !file.Mtime.Equal(modTime) {
		t.Fatalf("find mtime = %v, want %v", file.Mtime, modTime)
	}

	iter, err := backend.listFor(files_sdk.FolderListForParams{Path: "docs"})
	if err != nil {
		t.Fatalf("listFor failed: %v", err)
	}
	var paths []string
	for iter.Next() {
		paths = append(paths, iter.File().Path)
	}
	if err := iter.Err(); err != nil {
		t.Fatalf("list iterator failed: %v", err)
	}
	if strings.Join(paths, ",") != "/docs/report.txt,/docs/archive" {
		t.Fatalf("listed paths = %v", paths)
	}

	folder, err := backend.createFolder(files_sdk.FolderCreateParams{Path: "new-folder"})
	if err != nil {
		t.Fatalf("createFolder failed: %v", err)
	}
	if mkdirPath != "/new-folder" {
		t.Fatalf("mkdir path = %q, want /new-folder", mkdirPath)
	}
	if folder.Path != "/new-folder" || folder.Type != ProviderTypeDirectory {
		t.Fatalf("created folder = %#v", folder)
	}

	action, err := backend.move(files_sdk.FileMoveParams{Path: "old-name", Destination: "new-name"})
	if err != nil {
		t.Fatalf("move failed: %v", err)
	}
	if action.Status != "completed" {
		t.Fatalf("move status = %q, want completed", action.Status)
	}
	if renameSource != "/old-name" || renameDestination != "/new-name" {
		t.Fatalf("rename source/destination = %q/%q", renameSource, renameDestination)
	}

	recursive := true
	if err := backend.delete(files_sdk.FileDeleteParams{Path: "old-folder", Recursive: &recursive}); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if deletePath != "/old-folder" || !deleteRecursive {
		t.Fatalf("delete path/recursive = %q/%v", deletePath, deleteRecursive)
	}
}

func TestProviderRemoteBackendDownloadToFile(t *testing.T) {
	content := "downloaded provider payload"
	modTime := time.Date(2026, 5, 29, 11, 0, 0, 0, time.UTC)
	var readPath string

	provider := &testProviderBackend{
		readFunc: func(_ context.Context, path string) (io.ReadCloser, int64, error) {
			readPath = path
			return io.NopCloser(strings.NewReader(content)), int64(len(content)), nil
		},
		statFunc: func(_ context.Context, path string) (ProviderEntry, error) {
			return ProviderEntry{
				Path:    path,
				Type:    ProviderTypeFile,
				Size:    int64(len(content)),
				ModTime: modTime,
			}, nil
		},
	}
	backend := &providerRemoteBackend{provider: provider}
	outputPath := filepath.Join(t.TempDir(), "payload.txt")

	file, err := backend.downloadToFile(files_sdk.FileDownloadParams{Path: "docs/payload.txt"}, outputPath)
	if err != nil {
		t.Fatalf("downloadToFile failed: %v", err)
	}
	if readPath != "/docs/payload.txt" {
		t.Fatalf("read path = %q, want /docs/payload.txt", readPath)
	}
	if file.Path != "/docs/payload.txt" || file.Size != int64(len(content)) {
		t.Fatalf("downloaded file metadata = %#v", file)
	}
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if string(data) != content {
		t.Fatalf("downloaded content = %q, want %q", string(data), content)
	}
}

func TestProviderRemoteBackendMapsNotFoundErrors(t *testing.T) {
	backend := &providerRemoteBackend{
		provider: &testProviderBackend{
			statFunc: func(context.Context, string) (ProviderEntry, error) {
				return ProviderEntry{}, os.ErrNotExist
			},
		},
	}

	_, err := backend.find(files_sdk.FileFindParams{Path: "missing.txt"})
	if err == nil {
		t.Fatal("expected find to fail")
	}
	if !files_sdk.IsNotExist(err) {
		t.Fatalf("find error = %v, want not-found response error", err)
	}
}

func TestProviderRemoteBackendCreateFolderMapsNotFoundError(t *testing.T) {
	backend := &providerRemoteBackend{
		provider: &testProviderBackend{
			mkdirFunc: func(context.Context, string) error {
				return os.ErrNotExist
			},
		},
	}

	_, err := backend.createFolder(files_sdk.FolderCreateParams{Path: "missing-parent/new-folder"})
	if err == nil {
		t.Fatal("expected createFolder to fail")
	}
	if !files_sdk.IsNotExist(err) {
		t.Fatalf("createFolder error = %v, want not-found response error", err)
	}
}

func TestProviderRemoteBackendMoveMapsNotFoundError(t *testing.T) {
	backend := &providerRemoteBackend{
		provider: &testProviderBackend{
			renameFunc: func(context.Context, string, string) error {
				return os.ErrNotExist
			},
		},
	}

	_, err := backend.move(files_sdk.FileMoveParams{Path: "missing.txt", Destination: "renamed.txt"})
	if err == nil {
		t.Fatal("expected move to fail")
	}
	if !files_sdk.IsNotExist(err) {
		t.Fatalf("move error = %v, want not-found response error", err)
	}
}

func TestProviderRemoteBackendCreateFolderFallsBackWhenStatFails(t *testing.T) {
	var mkdirPath string
	provider := &testProviderBackend{
		mkdirFunc: func(_ context.Context, path string) error {
			mkdirPath = path
			return nil
		},
		statFunc: func(context.Context, string) (ProviderEntry, error) {
			return ProviderEntry{}, errors.New("stat unsupported")
		},
	}
	backend := &providerRemoteBackend{provider: provider}

	folder, err := backend.createFolder(files_sdk.FolderCreateParams{Path: "new-folder"})
	if err != nil {
		t.Fatalf("createFolder failed: %v", err)
	}
	if mkdirPath != "/new-folder" {
		t.Fatalf("mkdir path = %q, want /new-folder", mkdirPath)
	}
	if folder.Path != "/new-folder" {
		t.Fatalf("folder path = %q, want /new-folder", folder.Path)
	}
	if folder.Type != ProviderTypeDirectory {
		t.Fatalf("folder type = %q, want %q", folder.Type, ProviderTypeDirectory)
	}
	if folder.DisplayName != "new-folder" {
		t.Fatalf("folder display name = %q, want new-folder", folder.DisplayName)
	}
	if folder.Mtime == nil || folder.Mtime.IsZero() {
		t.Fatalf("folder mtime = %v, want a non-zero fallback", folder.Mtime)
	}
}

func TestNewRemoteFsProviderModeWiresProviderBackend(t *testing.T) {
	_, vfs, cacheStore := newTestRemoteFs(t)
	defer vfs.destroy()

	payload := "provider upload payload"
	inputPath := filepath.Join(t.TempDir(), "payload.txt")
	if err := os.WriteFile(inputPath, []byte(payload), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	input, err := os.Open(inputPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer input.Close()

	modTime := time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC)
	providerModTime := modTime.Add(30 * time.Second)
	var writtenPath string
	var writtenSize int64
	var writtenModTime time.Time
	var writtenContent string
	provider := &testProviderBackend{
		writeFunc: func(_ context.Context, path string, reader io.Reader, size int64, requestedModTime time.Time) (ProviderEntry, error) {
			data, err := io.ReadAll(reader)
			if err != nil {
				return ProviderEntry{}, err
			}
			writtenPath = path
			writtenSize = size
			writtenModTime = requestedModTime
			writtenContent = string(data)
			return ProviderEntry{
				Path:    path,
				Type:    ProviderTypeFile,
				Size:    int64(len(data)),
				ModTime: providerModTime,
			}, nil
		},
	}

	fs, err := newRemoteFs(MountParams{
		Config:          &files_sdk.Config{},
		ProviderBackend: provider,
		TmpFsPath:       t.TempDir(),
	}, vfs, &log.NoOpLogger{}, cacheStore)
	if err != nil {
		t.Fatalf("newRemoteFs failed: %v", err)
	}
	if fs.providerBackend != provider {
		t.Fatal("provider backend was not stored on RemoteFs")
	}
	if !fs.disableLocking {
		t.Fatal("provider mode should disable Files.com locks")
	}
	providerBackend, ok := fs.backend.(*providerRemoteBackend)
	if !ok {
		t.Fatalf("backend type = %T, want *providerRemoteBackend", fs.backend)
	}
	if providerBackend.provider != provider {
		t.Fatal("provider backend was not passed to remote backend adapter")
	}
	if fs.uploadWorkingCopy == nil {
		t.Fatal("provider mode should install a direct upload working-copy hook")
	}

	uploaded, err := fs.uploadWorkingCopy(context.Background(), nil, "/payload.txt", input, modTime, 0)
	if err != nil {
		t.Fatalf("uploadWorkingCopy failed: %v", err)
	}
	if uploaded.size != int64(len(payload)) {
		t.Fatalf("uploaded size = %d, want %d", uploaded.size, len(payload))
	}
	if !uploaded.modTime.Equal(providerModTime) {
		t.Fatalf("uploaded mod time = %v, want %v", uploaded.modTime, providerModTime)
	}
	if writtenPath != "/payload.txt" {
		t.Fatalf("provider write path = %q, want /payload.txt", writtenPath)
	}
	if writtenSize != int64(len(payload)) {
		t.Fatalf("provider write size = %d, want %d", writtenSize, len(payload))
	}
	if !writtenModTime.Equal(modTime) {
		t.Fatalf("provider write mod time = %v, want %v", writtenModTime, modTime)
	}
	if writtenContent != payload {
		t.Fatalf("provider write content = %q, want %q", writtenContent, payload)
	}
}

func TestProviderRemoteBackendUpdateUsesOptionalMtimeBackend(t *testing.T) {
	modTime := time.Date(2026, 6, 5, 9, 30, 0, 0, time.UTC)
	var touchedPath string
	var touchedTime time.Time
	statCalled := false
	provider := &testMtimeProviderBackend{
		testProviderBackend: &testProviderBackend{
			statFunc: func(context.Context, string) (ProviderEntry, error) {
				statCalled = true
				return ProviderEntry{}, nil
			},
		},
		mtimeFunc: func(_ context.Context, path string, providerModTime time.Time) (ProviderEntry, error) {
			touchedPath = path
			touchedTime = providerModTime
			return ProviderEntry{
				Path:    path,
				Type:    ProviderTypeFile,
				Size:    4,
				ModTime: providerModTime,
			}, nil
		},
	}
	backend := &providerRemoteBackend{provider: provider}

	file, err := backend.update(files_sdk.FileUpdateParams{
		Path:          "mtime.txt",
		ProvidedMtime: &modTime,
	})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if statCalled {
		t.Fatal("Stat should not be called when ProviderMtimeBackend handles the update")
	}
	if touchedPath != "/mtime.txt" {
		t.Fatalf("touched path = %q, want /mtime.txt", touchedPath)
	}
	if !touchedTime.Equal(modTime) {
		t.Fatalf("touched mtime = %v, want %v", touchedTime, modTime)
	}
	if file.Mtime == nil || !file.Mtime.Equal(modTime) {
		t.Fatalf("file mtime = %v, want %v", file.Mtime, modTime)
	}
}

func TestProviderRemoteBackendUpdateFallsBackToStatWithoutMtimeBackend(t *testing.T) {
	requestedModTime := time.Date(2026, 6, 5, 9, 30, 0, 0, time.UTC)
	existingModTime := requestedModTime.Add(-time.Hour)
	var statPath string
	provider := &testProviderBackend{
		statFunc: func(_ context.Context, path string) (ProviderEntry, error) {
			statPath = path
			return ProviderEntry{
				Path:    path,
				Type:    ProviderTypeFile,
				Size:    4,
				ModTime: existingModTime,
			}, nil
		},
	}
	backend := &providerRemoteBackend{provider: provider}

	file, err := backend.update(files_sdk.FileUpdateParams{
		Path:          "mtime.txt",
		ProvidedMtime: &requestedModTime,
	})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if statPath != "/mtime.txt" {
		t.Fatalf("stat path = %q, want /mtime.txt", statPath)
	}
	if file.Mtime == nil || !file.Mtime.Equal(existingModTime) {
		t.Fatalf("file mtime = %v, want stat mtime %v", file.Mtime, existingModTime)
	}
}

func TestProviderModeUploadWorkingCopyReturnsStatError(t *testing.T) {
	_, vfs, cacheStore := newTestRemoteFs(t)
	defer vfs.destroy()

	writeCalled := false
	provider := &testProviderBackend{
		writeFunc: func(context.Context, string, io.Reader, int64, time.Time) (ProviderEntry, error) {
			writeCalled = true
			return ProviderEntry{}, nil
		},
	}
	fs, err := newRemoteFs(MountParams{
		Config:          &files_sdk.Config{},
		ProviderBackend: provider,
		TmpFsPath:       t.TempDir(),
	}, vfs, &log.NoOpLogger{}, cacheStore)
	if err != nil {
		t.Fatalf("newRemoteFs failed: %v", err)
	}

	statErr := errors.New("stat failed")
	_, err = fs.uploadWorkingCopy(context.Background(), nil, "/payload.txt", statErrorUploadReader{
		Reader: strings.NewReader("payload"),
		err:    statErr,
	}, time.Now(), 0)
	if !errors.Is(err, statErr) {
		t.Fatalf("uploadWorkingCopy error = %v, want %v", err, statErr)
	}
	if writeCalled {
		t.Fatal("provider Write should not be called when reader.Stat fails")
	}
}

func TestProviderModeUploadWorkingCopyReturnsProviderWriteError(t *testing.T) {
	_, vfs, cacheStore := newTestRemoteFs(t)
	defer vfs.destroy()

	inputPath := filepath.Join(t.TempDir(), "payload.txt")
	if err := os.WriteFile(inputPath, []byte("payload"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	input, err := os.Open(inputPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer input.Close()

	writeErr := errors.New("provider write failed")
	provider := &testProviderBackend{
		writeFunc: func(context.Context, string, io.Reader, int64, time.Time) (ProviderEntry, error) {
			return ProviderEntry{}, writeErr
		},
	}
	fs, err := newRemoteFs(MountParams{
		Config:          &files_sdk.Config{},
		ProviderBackend: provider,
		TmpFsPath:       t.TempDir(),
	}, vfs, &log.NoOpLogger{}, cacheStore)
	if err != nil {
		t.Fatalf("newRemoteFs failed: %v", err)
	}

	_, err = fs.uploadWorkingCopy(context.Background(), nil, "/payload.txt", input, time.Now(), 0)
	if !errors.Is(err, writeErr) {
		t.Fatalf("uploadWorkingCopy error = %v, want %v", err, writeErr)
	}
}
