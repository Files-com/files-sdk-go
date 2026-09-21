package file

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/directory"
	"github.com/Files-com/files-sdk-go/v3/lib"
)

// destinationPath is a local path split into the directory the caller selected
// and the path below it. The path below the directory can come from a server
// response, so every filesystem operation opens the directory with os.Root and
// resolves the rest inside it. A server path or a local link that leads outside
// the directory makes the operation fail instead of being followed.
type destinationPath struct {
	dir  string // caller-selected directory, as the caller gave it
	name string // local path below dir; "." is the directory itself
}

type downloadPathError struct {
	error
}

func (e downloadPathError) Unwrap() error {
	return e.error
}

func (e downloadPathError) ErrorType() string {
	return "invalid_path"
}

func (e downloadPathError) PublicError() string {
	return "Cannot download this file because its path is not valid on this computer."
}

// explicitDestination is a path the caller selected in full, such as the
// output file of a single-file download. Only its last element is resolved
// inside its directory.
func explicitDestination(path string) destinationPath {
	return destinationPath{dir: filepath.Dir(path), name: filepath.Base(path)}
}

// downloadDestination returns where a listed remote file is written. A folder
// download joins the destination directory with the file's path below the
// requested remote folder; a single-file download uses the caller's output path.
func downloadDestination(job *Job, file files_sdk.File) (destinationPath, error) {
	if job.Type == directory.File {
		return explicitDestination(job.LocalPath), nil
	}
	name, err := localNameBelow(job.RemotePath, file.Path)
	if err != nil {
		return destinationPath{}, downloadPathError{err}
	}
	return destinationPath{dir: normalizePath(job.LocalPath), name: name}, nil
}

// localNameBelow converts the part of a server response path below the
// requested remote folder into a local relative path. Remote paths are
// slash-separated and compared the way Files.com compares them. The result
// must be valid for the host platform: a path that is not a descendant of the
// remote folder, contains "." or ".." elements, or uses forms the platform
// reserves (on Windows: backslashes, drive or stream colons and device names)
// is an error.
func localNameBelow(remoteFolder, remotePath string) (string, error) {
	folder := remoteSegments(remoteFolder)
	segments := remoteSegments(remotePath)
	if len(segments) < len(folder) {
		return "", fmt.Errorf("download: server path %q is outside the requested folder %q", remotePath, remoteFolder)
	}
	for i, segment := range folder {
		if lib.NormalizeForComparison(segment) != lib.NormalizeForComparison(segments[i]) {
			return "", fmt.Errorf("download: server path %q is outside the requested folder %q", remotePath, remoteFolder)
		}
	}
	below := strings.Join(segments[len(folder):], "/")
	if below == "" {
		return ".", nil
	}
	name, err := filepath.Localize(below)
	if err != nil {
		return "", fmt.Errorf("download: server path %q is not a valid local path below %q: %w", remotePath, remoteFolder, err)
	}
	return name, nil
}

// remoteSegments splits a slash-separated remote path into its elements. The
// root folder ("", "." or "/") has no elements.
func remoteSegments(remotePath string) []string {
	remotePath = strings.Trim(remotePath, "/")
	if remotePath == "" || remotePath == "." {
		return nil
	}
	return strings.Split(remotePath, "/")
}

func (p destinationPath) String() string {
	return filepath.Join(p.dir, p.name)
}

// withName returns the path of name inside the same caller-selected directory.
func (p destinationPath) withName(name string) destinationPath {
	return destinationPath{dir: p.dir, name: name}
}

// base is the last element of the path below the directory.
func (p destinationPath) base() string {
	return filepath.Base(p.name)
}

// parent is the directory that holds p, or the caller-selected directory itself.
func (p destinationPath) parent() destinationPath {
	return p.withName(filepath.Dir(p.name))
}

// inRoot runs op with the caller-selected directory opened as an os.Root.
// Errors from op are rewritten to name the full local path, which is the path
// the download reports to callers.
func inRoot[T any](p destinationPath, op func(*os.Root) (T, error)) (T, error) {
	root, err := os.OpenRoot(p.dir)
	if err != nil {
		var zero T
		return zero, err
	}
	defer root.Close()
	result, err := op(root)
	return result, p.fullPathError(err)
}

func (p destinationPath) do(op func(*os.Root) error) error {
	_, err := inRoot(p, func(root *os.Root) (struct{}, error) {
		return struct{}{}, op(root)
	})
	return err
}

func (p destinationPath) fullPathError(err error) error {
	var pathErr *fs.PathError
	if errors.As(err, &pathErr) {
		pathErr.Path = filepath.Join(p.dir, pathErr.Path)
	}
	var linkErr *os.LinkError
	if errors.As(err, &linkErr) {
		linkErr.Old = filepath.Join(p.dir, linkErr.Old)
		linkErr.New = filepath.Join(p.dir, linkErr.New)
	}
	return err
}

func (p destinationPath) stat() (fs.FileInfo, error) {
	return inRoot(p, func(root *os.Root) (fs.FileInfo, error) { return root.Stat(p.name) })
}

func (p destinationPath) lstat() (fs.FileInfo, error) {
	return inRoot(p, func(root *os.Root) (fs.FileInfo, error) { return root.Lstat(p.name) })
}

func (p destinationPath) openFile(flag int, perm fs.FileMode) (*os.File, error) {
	return inRoot(p, func(root *os.Root) (*os.File, error) { return root.OpenFile(p.name, flag, perm) })
}

func (p destinationPath) create() (*os.File, error) {
	return p.openFile(os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
}

func (p destinationPath) open() (*os.File, error) {
	return p.openFile(os.O_RDONLY, 0)
}

// mkdirAll creates the path below the directory. The directory itself is the
// caller's choice, so a missing one is created directly, as it always was.
func (p destinationPath) mkdirAll() error {
	if err := os.MkdirAll(p.dir, 0755); err != nil {
		return err
	}
	return p.do(func(root *os.Root) error { return root.MkdirAll(p.name, 0755) })
}

func (p destinationPath) remove() error {
	return p.do(func(root *os.Root) error { return root.Remove(p.name) })
}

func (p destinationPath) chtimes(t time.Time) error {
	return p.do(func(root *os.Root) error { return root.Chtimes(p.name, t, t) })
}

// rename renames p to final, which must be inside the same caller-selected
// directory.
func (p destinationPath) rename(final destinationPath) error {
	return p.do(func(root *os.Root) error { return root.Rename(p.name, final.name) })
}

// moveTo renames the completed temporary file p to final. Inside one
// caller-selected directory this is a single confined rename. Across
// directories (an external TempPath) the file is staged through the top level
// of both directories: the .download folder on macOS and the server-derived
// part of final are only ever resolved by os.Root, and the one plain rename
// between the directories names caller-selected directories plus names this
// code created. If the destination root is read-only, a confined copy stages
// the file in its writable parent instead, keeping the original until the
// destination replacement succeeds.
func (p destinationPath) moveTo(ctx context.Context, final destinationPath) (err error) {
	if filepath.Clean(p.dir) == filepath.Clean(final.dir) {
		return p.rename(final)
	}
	target, err := reserveMoveName(final.dir)
	if errors.Is(err, fs.ErrPermission) {
		return p.copyTo(ctx, final)
	}
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = target.remove()
		}
	}()

	source := p
	if filepath.Dir(p.name) != "." {
		top, reserveErr := reserveMoveName(p.dir)
		if reserveErr != nil {
			return reserveErr
		}
		if renameErr := p.rename(top); renameErr != nil {
			_ = top.remove()
			return renameErr
		}
		source = top
		defer func() {
			if err != nil {
				_ = source.rename(p)
			}
		}()
	}
	if err = os.Rename(source.String(), target.String()); err != nil {
		return err
	}
	if err = target.rename(final); err != nil {
		_ = os.Rename(target.String(), source.String())
		return err
	}
	return nil
}

// copyTo is the fallback for an external temp file when the selected root is
// read-only but the existing destination parent is writable. Pin that parent
// through os.Root so creation, replacement and cleanup never resolve it again
// through a plain path. The extra copy needs space for both complete files.
func (p destinationPath) copyTo(ctx context.Context, final destinationPath) (err error) {
	in, err := p.open()
	if err != nil {
		return err
	}
	defer in.Close()
	info, err := in.Stat()
	if err != nil {
		return err
	}
	parent, err := inRoot(final, func(root *os.Root) (*os.Root, error) {
		return root.OpenRoot(filepath.Dir(final.name))
	})
	if err != nil {
		return err
	}
	defer parent.Close()

	// Keep the staging name independent of the final basename, which may
	// already occupy the filesystem's entire component length allowance.
	stage := "." + TempDownloadExtension + "-" + rand.Text()
	out, err := parent.OpenFile(stage, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cleanupErr := parent.Remove(stage); cleanupErr != nil && !errors.Is(cleanupErr, fs.ErrNotExist) {
				err = errors.Join(err, cleanupErr)
			}
		}
	}()
	_, copyErr := copyWithContext(ctx, out, in)
	// Chmod the open file, not a path that could have become a symlink.
	if copyErr == nil {
		copyErr = out.Chmod(info.Mode().Perm())
	}
	closeErr := out.Close()
	if err = errors.Join(copyErr, closeErr); err != nil {
		return errors.Join(err, context.Cause(ctx))
	}
	if err = context.Cause(ctx); err != nil {
		return err
	}
	if err = parent.Rename(stage, final.base()); err != nil {
		return err
	}
	// Windows cannot remove an open source file. Report removal failures;
	// the completed destination and original source both remain available.
	if err = in.Close(); err != nil {
		return err
	}
	return p.remove()
}

// reserveMoveName creates an empty file with a fresh name directly inside dir
// for a file moving through the directory. The name does not include the final
// file name, so it fits wherever the ".download" name fit.
func reserveMoveName(dir string) (destinationPath, error) {
	file, err := os.CreateTemp(dir, "."+TempDownloadExtension+"-*")
	if err != nil {
		return destinationPath{}, err
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(file.Name())
		return destinationPath{}, err
	}
	return destinationPath{dir: dir, name: filepath.Base(file.Name())}, nil
}
