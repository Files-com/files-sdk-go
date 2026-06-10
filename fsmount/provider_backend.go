package fsmount

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	path_lib "path"
	"strings"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
)

const (
	ProviderTypeFile      = "file"
	ProviderTypeDirectory = "directory"
)

// errProviderUploadNotSupported is a defensive marker returned from the
// upload methods on providerRemoteBackend. The SDK upload pipeline is bypassed
// for provider mounts (uploads go through RemoteFs.uploadWorkingCopy and
// providerBackend.Write directly), so these methods should be unreachable.
// They exist only to satisfy the remoteBackend interface; if a future code
// path does hit them, returning a sentinel surfaces the regression instead
// of silently failing the upload.
var errProviderUploadNotSupported = errors.New("provider backend does not support SDK upload pipeline")

// ProviderEntry is the minimal file metadata fsmount needs from a direct provider.
type ProviderEntry struct {
	Path        string
	DisplayName string
	Type        string
	Size        int64
	ModTime     time.Time
	CreatedAt   time.Time
	Permissions string
}

// ProviderBackend is the direct provider interface used by fsmount when the
// caller supplies a non-Files.com backend.
type ProviderBackend interface {
	Stat(ctx context.Context, path string) (ProviderEntry, error)
	List(ctx context.Context, path string) ([]ProviderEntry, error)
	Read(ctx context.Context, path string) (io.ReadCloser, int64, error)
	Write(ctx context.Context, path string, reader io.Reader, size int64, modTime time.Time) (ProviderEntry, error)
	Mkdir(ctx context.Context, path string) error
	Delete(ctx context.Context, path string, recursive bool) error
	Rename(ctx context.Context, sourcePath string, destinationPath string) error
}

// ProviderMtimeBackend is optional. Backends that can persist an explicit
// modified time should implement it so FUSE Utimens calls survive remounts.
type ProviderMtimeBackend interface {
	SetMtime(ctx context.Context, path string, modTime time.Time) (ProviderEntry, error)
}

func (e ProviderEntry) sdkFile() files_sdk.File {
	if e.Type == "" {
		e.Type = ProviderTypeFile
	}
	if e.Path == "" {
		e.Path = "/"
	}
	if e.DisplayName == "" {
		e.DisplayName = path_lib.Base(strings.TrimSuffix(e.Path, "/"))
		if e.DisplayName == "." || e.DisplayName == "" {
			e.DisplayName = "/"
		}
	}

	file := files_sdk.File{
		Path:        e.Path,
		DisplayName: e.DisplayName,
		Type:        e.Type,
		Size:        e.Size,
		Permissions: e.Permissions,
	}
	if !e.ModTime.IsZero() {
		modTime := e.ModTime
		file.Mtime = &modTime
	}
	if !e.CreatedAt.IsZero() {
		createdAt := e.CreatedAt
		file.CreatedAt = &createdAt
	}
	return file
}

func providerEntryForPath(path string, entryType string) ProviderEntry {
	now := time.Now()
	return ProviderEntry{
		Path:        cleanProviderPath(path),
		DisplayName: path_lib.Base(strings.TrimSuffix(path, "/")),
		Type:        entryType,
		ModTime:     now,
		CreatedAt:   now,
	}
}

func cleanProviderPath(path string) string {
	clean := path_lib.Clean("/" + strings.TrimPrefix(path, "/"))
	if clean == "." {
		return "/"
	}
	return clean
}

// providerContext intentionally forwards only context-bearing request options.
// Other SDK request/response hooks do not apply to direct provider calls; the
// download adapter still passes opts through BuildResponse for HTTP semantics.
func providerContext(opts ...files_sdk.RequestResponseOption) context.Context {
	return files_sdk.ContextOption(opts)
}

func providerNotFound(path string, err error) error {
	if err == nil {
		return nil
	}
	return files_sdk.ResponseError{
		Type:         string(files_sdk.ErrFileNotFound),
		ErrorMessage: fmt.Sprintf("%s: %v", path, err),
	}
}

type providerRemoteBackend struct {
	provider ProviderBackend
}

// findCurrent is unreachable on the provider path: RemoteFs.Init skips the
// API key lookup when a ProviderBackend is set. The method exists only to
// satisfy the remoteBackend interface and returns a zero value so any future
// accidental call is benign rather than introducing a synthetic UserId.
func (b *providerRemoteBackend) findCurrent(opts ...files_sdk.RequestResponseOption) (files_sdk.ApiKey, error) {
	return files_sdk.ApiKey{}, nil
}

func (b *providerRemoteBackend) find(params files_sdk.FileFindParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	entry, err := b.provider.Stat(providerContext(opts...), cleanProviderPath(params.Path))
	if os.IsNotExist(err) {
		err = providerNotFound(params.Path, err)
	}
	return entry.sdkFile(), err
}

func (b *providerRemoteBackend) listFor(params files_sdk.FolderListForParams, opts ...files_sdk.RequestResponseOption) (remoteFileIter, error) {
	entries, err := b.provider.List(providerContext(opts...), cleanProviderPath(params.Path))
	if os.IsNotExist(err) {
		err = providerNotFound(params.Path, err)
	}
	if err != nil {
		return nil, err
	}
	files := make([]files_sdk.File, 0, len(entries))
	for _, entry := range entries {
		files = append(files, entry.sdkFile())
	}
	return &providerFileIter{files: files}, nil
}

func (b *providerRemoteBackend) createFolder(params files_sdk.FolderCreateParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	ctx := providerContext(opts...)
	path := cleanProviderPath(params.Path)
	if err := b.provider.Mkdir(ctx, path); err != nil {
		if os.IsNotExist(err) {
			err = providerNotFound(params.Path, err)
		}
		return files_sdk.File{}, err
	}
	entry, err := b.provider.Stat(ctx, path)
	if err != nil {
		entry = providerEntryForPath(path, ProviderTypeDirectory)
	}
	return entry.sdkFile(), nil
}

func (b *providerRemoteBackend) move(params files_sdk.FileMoveParams, opts ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error) {
	err := b.provider.Rename(providerContext(opts...), cleanProviderPath(params.Path), cleanProviderPath(params.Destination))
	if err != nil {
		if os.IsNotExist(err) {
			err = providerNotFound(params.Path, err)
		}
		return files_sdk.FileAction{}, err
	}
	return files_sdk.FileAction{Status: "completed"}, nil
}

func (b *providerRemoteBackend) update(params files_sdk.FileUpdateParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	ctx := providerContext(opts...)
	path := cleanProviderPath(params.Path)
	if params.ProvidedMtime != nil {
		if mtimeBackend, ok := b.provider.(ProviderMtimeBackend); ok {
			entry, err := mtimeBackend.SetMtime(ctx, path, *params.ProvidedMtime)
			if os.IsNotExist(err) {
				err = providerNotFound(params.Path, err)
			}
			return entry.sdkFile(), err
		}
	}

	entry, err := b.provider.Stat(ctx, path)
	if os.IsNotExist(err) {
		err = providerNotFound(params.Path, err)
	}
	return entry.sdkFile(), err
}

func (b *providerRemoteBackend) uploadWithResume(opts ...file.UploadOption) (file.UploadResumable, error) {
	return file.UploadResumable{}, errProviderUploadNotSupported
}

func (b *providerRemoteBackend) upload(opts ...file.UploadOption) error {
	return errProviderUploadNotSupported
}

func (b *providerRemoteBackend) downloadToFile(params files_sdk.FileDownloadParams, filePath string, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	ctx := providerContext(opts...)
	path := params.Path
	if path == "" {
		path = params.File.Path
	}
	reader, size, err := b.provider.Read(ctx, cleanProviderPath(path))
	if os.IsNotExist(err) {
		err = providerNotFound(path, err)
	}
	if err != nil {
		return files_sdk.File{}, err
	}
	defer reader.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return files_sdk.File{}, err
	}
	defer out.Close()
	if _, err := io.Copy(out, reader); err != nil {
		return files_sdk.File{}, err
	}
	entry, err := b.provider.Stat(ctx, cleanProviderPath(path))
	if err != nil {
		entry = providerEntryForPath(path, ProviderTypeFile)
		entry.Size = size
	}
	return entry.sdkFile(), nil
}

func (b *providerRemoteBackend) download(params files_sdk.FileDownloadParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	ctx := providerContext(opts...)
	path := params.Path
	if path == "" {
		path = params.File.Path
	}
	reader, size, err := b.provider.Read(ctx, cleanProviderPath(path))
	if os.IsNotExist(err) {
		err = providerNotFound(path, err)
	}
	if err != nil {
		return files_sdk.File{}, err
	}
	resp := &http.Response{StatusCode: http.StatusOK, Body: reader}
	if _, err := files_sdk.BuildResponse(resp, opts...); err != nil {
		return files_sdk.File{}, err
	}
	entry, err := b.provider.Stat(ctx, cleanProviderPath(path))
	if err != nil {
		entry = providerEntryForPath(path, ProviderTypeFile)
		entry.Size = size
	}
	return entry.sdkFile(), nil
}

// The three lock methods below are unreachable on the provider path: Mount
// forces DisableLocking = true when a ProviderBackend is set, and every lock
// call site in RemoteFs is guarded by `if fs.disableLocking`. They panic
// rather than returning empty success so a future regression that bypasses
// the guard surfaces loudly instead of silently dropping locks. Upload methods
// return a sentinel instead because a stray upload call can be surfaced as an
// operation error without killing the mount process.
const providerLockUnreachable = "lock operation invoked on provider backend; DisableLocking should have short-circuited this call"

func (b *providerRemoteBackend) createLock(params files_sdk.LockCreateParams, opts ...files_sdk.RequestResponseOption) (files_sdk.Lock, error) {
	panic(providerLockUnreachable)
}

func (b *providerRemoteBackend) deleteLock(params files_sdk.LockDeleteParams, opts ...files_sdk.RequestResponseOption) error {
	panic(providerLockUnreachable)
}

func (b *providerRemoteBackend) listLocksFor(params files_sdk.LockListForParams, opts ...files_sdk.RequestResponseOption) (remoteLockIter, error) {
	panic(providerLockUnreachable)
}

func (b *providerRemoteBackend) delete(params files_sdk.FileDeleteParams, opts ...files_sdk.RequestResponseOption) error {
	recursive := false
	if params.Recursive != nil {
		recursive = *params.Recursive
	}
	err := b.provider.Delete(providerContext(opts...), cleanProviderPath(params.Path), recursive)
	if os.IsNotExist(err) {
		return providerNotFound(params.Path, err)
	}
	return err
}

func (b *providerRemoteBackend) wait(action files_sdk.FileAction, status func(files_sdk.FileMigration), opts ...files_sdk.RequestResponseOption) (files_sdk.FileMigration, error) {
	migration := files_sdk.FileMigration{Status: "completed"}
	if status != nil {
		status(migration)
	}
	return migration, nil
}

type providerFileIter struct {
	files []files_sdk.File
	index int
	err   error
}

func (i *providerFileIter) Next() bool {
	if i.index >= len(i.files) {
		return false
	}
	i.index++
	return true
}

func (i *providerFileIter) File() files_sdk.File {
	if i.index == 0 || i.index > len(i.files) {
		return files_sdk.File{}
	}
	return i.files[i.index-1]
}

func (i *providerFileIter) Err() error {
	return i.err
}
