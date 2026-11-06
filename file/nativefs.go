package file

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"sort"
	"strings"
	"sync"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/folder"
	"github.com/Files-com/files-sdk-go/v3/lib"
)

// NativeFS provides relative io/fs paths and independently owned read handles.
// Create retains FS's staged upload-on-Close behavior; MkdirTemp is unsupported.
type NativeFS struct {
	source *FS
}

// Native returns an adapter without changing the FS's existing handle or path behavior.
func (f *FS) Native() *NativeFS {
	return &NativeFS{source: f}
}

// WithContext returns an io/fs view using ctx for its API calls and transfers.
func (n *NativeFS) WithContext(ctx context.Context) fs.FS {
	return n.source.WithContext(ctx).(*FS).Native()
}

func (n *NativeFS) remotePath(op, name string) (string, error) {
	if !fs.ValidPath(name) {
		return "", &fs.PathError{Op: op, Path: name, Err: fs.ErrInvalid}
	}
	if name == "." {
		name = ""
	}
	return lib.UrlJoinNoEscape(n.source.Root, name), nil
}

func nativePathError(op, name string, err error) error {
	if err == nil {
		return nil
	}
	if files_sdk.IsNotExist(err) && !errors.Is(err, fs.ErrNotExist) {
		err = errors.Join(fs.ErrNotExist, err)
	}
	return &fs.PathError{Op: op, Path: name, Err: err}
}

// Open returns a fresh handle, including when the source FS has caching enabled.
// Each handle has its own initialized reader and directory cursor.
func (n *NativeFS) Open(name string) (fs.File, error) {
	remote, err := n.remotePath("open", name)
	if err != nil {
		return nil, err
	}
	uncached := (&FS{Context: n.source.Context}).Init(n.source.Config, false)
	file, err := uncached.openRemoteFile(remote)
	if err != nil {
		return nil, nativePathError("open", name, err)
	}
	if remote == "" {
		file.File.DisplayName = "."
	}
	parent := n.source.Context
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)
	uncached.Context = ctx
	handle := &nativeFile{file: file, name: name, parent: parent, cancel: cancel}
	if file.IsDir() {
		handle.directory = &ReadDirFile{File: file}
	}
	return handle, nil
}

// ReadDir returns the directory's entries sorted by name. Open's directory
// handles instead retain the server's order across ReadDir calls.
func (n *NativeFS) ReadDir(name string) ([]fs.DirEntry, error) {
	opened, err := n.Open(name)
	if err != nil {
		return nil, nativePathError("readdir", name, err)
	}
	handle := opened.(*nativeFile)
	entries, err := handle.ReadDir(-1)
	if closeErr := handle.Close(); closeErr != nil {
		err = errors.Join(err, closeErr)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	return entries, err
}

// Create uses the existing disk-staged writer and uploads only when it is closed.
func (n *NativeFS) Create(name string) (io.WriteCloser, error) {
	if _, err := n.remotePath("create", name); err != nil {
		return nil, err
	}
	path := name
	if name == "." {
		path = ""
	}
	writer, err := n.source.Create(path)
	return writer, nativePathError("create", name, err)
}

func (n *NativeFS) MkdirAll(name string, _ fs.FileMode) error {
	_, err := n.remotePath("mkdir", name)
	if err != nil {
		return err
	}
	if name == "." {
		return n.requireDirectory("mkdir", name, nil)
	}
	var relative string
	client := &folder.Client{Config: n.source.Config}
	for _, component := range strings.Split(name, "/") {
		relative = lib.UrlJoinNoEscape(relative, component)
		remote := lib.UrlJoinNoEscape(n.source.Root, relative)
		created, err := client.Create(files_sdk.FolderCreateParams{Path: remote, MkdirParents: lib.Bool(false)}, files_sdk.WithContext(n.source.Context))
		if err != nil && !errors.Is(err, files_sdk.ErrDestinationExists) {
			return nativePathError("mkdir", name, err)
		}
		if err == nil && created.Type == "directory" && lib.NormalizeForComparison(created.Path) == lib.NormalizeForComparison(remote) {
			continue
		}
		if err := n.requireDirectory("mkdir", relative, err); err != nil {
			return nativePathError("mkdir", name, err)
		}
	}
	return nil
}

func (n *NativeFS) requireDirectory(op, name string, cause error) error {
	opened, err := n.Open(name)
	if err != nil {
		return nativePathError(op, name, errors.Join(err, cause))
	}
	info, err := opened.Stat()
	closeErr := opened.Close()
	if err != nil || closeErr != nil {
		return nativePathError(op, name, errors.Join(err, closeErr, cause))
	}
	if !info.IsDir() {
		return nativePathError(op, name, errors.Join(fs.ErrExist, cause))
	}
	return nil
}

func (n *NativeFS) Remove(name string) error {
	remote, err := n.remotePath("remove", name)
	if err != nil {
		return err
	}
	if name == "." {
		return nativePathError("remove", name, fs.ErrInvalid)
	}
	return nativePathError("remove", name, n.source.Remove(remote))
}

func (n *NativeFS) RemoveAll(name string) error {
	remote, err := n.remotePath("removeall", name)
	if err != nil {
		return err
	}
	if name == "." {
		return nativePathError("removeall", name, fs.ErrInvalid)
	}
	err = n.source.RemoveAll(remote)
	if files_sdk.IsNotExist(err) {
		return nil
	}
	return nativePathError("removeall", name, err)
}

// Chtimes updates mtime when atime is zero. Nonzero atime is unsupported;
// when both times are zero it verifies the path without changing metadata.
func (n *NativeFS) Chtimes(name string, atime, mtime time.Time) error {
	remote, err := n.remotePath("chtimes", name)
	if err != nil {
		return err
	}
	if !atime.IsZero() {
		return nativePathError("chtimes", name, fmt.Errorf("setting remote access time is unavailable: %w", errors.ErrUnsupported))
	}
	if mtime.IsZero() {
		opened, err := n.Open(name)
		if err != nil {
			return nativePathError("chtimes", name, err)
		}
		return nativePathError("chtimes", name, opened.Close())
	}
	return nativePathError("chtimes", name, n.source.Chtimes(remote, atime, mtime))
}

// MkdirTemp cannot guarantee exclusive remote allocation and sends no request.
func (n *NativeFS) MkdirTemp(dir, pattern string) (string, error) {
	return n.source.MkdirTemp(dir, pattern)
}

type nativeFile struct {
	operation sync.Mutex
	mutex     sync.Mutex
	file      *File
	directory *ReadDirFile
	name      string
	parent    context.Context
	cancel    context.CancelFunc
	closed    bool
	readEOF   bool
}

func (f *nativeFile) isClosed() bool {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.closed
}

func (f *nativeFile) operationError(op string, err error) error {
	if f.isClosed() {
		return nativePathError(op, f.name, errors.Join(fs.ErrClosed, err))
	}
	if err == io.EOF {
		return err
	}
	return nativePathError(op, f.name, err)
}

func (f *nativeFile) Stat() (fs.FileInfo, error) {
	f.operation.Lock()
	defer f.operation.Unlock()
	if f.isClosed() {
		return nil, nativePathError("stat", f.name, fs.ErrClosed)
	}
	info, err := f.file.Stat()
	return info, f.operationError("stat", err)
}

func (f *nativeFile) Read(p []byte) (int, error) {
	f.operation.Lock()
	defer f.operation.Unlock()
	if f.isClosed() {
		return 0, nativePathError("read", f.name, fs.ErrClosed)
	}
	if len(p) == 0 {
		return 0, nil
	}
	if f.directory != nil {
		return 0, nativePathError("read", f.name, fs.ErrInvalid)
	}
	n, err := f.file.Read(p)
	if err == io.EOF && !f.isClosed() {
		f.readEOF = true
	}
	return n, f.operationError("read", err)
}

func (f *nativeFile) ReadDir(n int) ([]fs.DirEntry, error) {
	f.operation.Lock()
	defer f.operation.Unlock()
	if f.isClosed() {
		return nil, nativePathError("readdir", f.name, fs.ErrClosed)
	}
	if f.directory == nil {
		return nil, nativePathError("readdir", f.name, fs.ErrInvalid)
	}
	entries, err := f.directory.ReadDir(n)
	return entries, f.operationError("readdir", err)
}

func (f *nativeFile) Close() error {
	f.mutex.Lock()
	if f.closed {
		f.mutex.Unlock()
		return nativePathError("close", f.name, fs.ErrClosed)
	}
	f.closed = true
	f.mutex.Unlock()
	f.cancel()
	f.operation.Lock()
	defer f.operation.Unlock()
	// Close's status verification still uses the caller's context. Cancelling
	// this handle's in-flight read must not cancel that separate verification.
	f.file.FS.Context = f.parent
	return nativePathError("close", f.name, f.file.closeRead(f.readEOF))
}
