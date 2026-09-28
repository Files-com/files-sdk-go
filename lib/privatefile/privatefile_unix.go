//go:build !windows

package privatefile

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
)

// privateDirMode is the widest mode a private entry may have: the owner's
// bits only. The process umask can only remove bits from it.
const privateDirMode fs.FileMode = 0o700

// Mkdir creates a directory only the current user can access. It fails if
// path already exists.
func Mkdir(path string) error {
	return os.Mkdir(path, privateDirMode)
}

// EnsureDir makes path a private directory this process owns: it creates a
// missing one, and removes group and other access from an existing one.
//
// An existing path must be a real directory (not a symbolic link) owned by the
// current user. Anything else is refused rather than used, because chmod on a
// link or a directory another user prepared could change unrelated data or
// hand this process's data to that user.
func EnsureDir(path string) error {
	err := Mkdir(path)
	if err == nil || !errors.Is(err, fs.ErrExist) {
		return err
	}
	info, err := lstatDirectory(path)
	if err != nil {
		return err
	}
	if err := ownedByCurrentUser("ensuredir", path, info); err != nil {
		return err
	}
	if info.Mode().Perm()&^privateDirMode == 0 {
		return nil
	}
	return os.Chmod(path, info.Mode().Perm()&privateDirMode)
}

// CheckEntry reports whether one entry inside a directory this package is
// about to make private may be made private along with it. info must be the
// Lstat of path. Only entries this package could have created qualify: a
// regular file or directory owned by the current user with a single name. A
// symbolic link is refused because chmod would follow it to whatever it points
// at; a file with several names is refused because its mode is shared with a
// name somewhere this process does not control; an entry owned by someone
// else is refused because it is not this process's data and could not be made
// private anyway. Nothing is modified.
func CheckEntry(path string, info fs.FileInfo) error {
	if isLink(info) {
		return &fs.PathError{Op: "chmod", Path: path, Err: errors.New("unexpected symbolic link inside a private directory")}
	}
	if err := ownedByCurrentUser("chmod", path, info); err != nil {
		return err
	}
	if stat, ok := info.Sys().(*syscall.Stat_t); ok && !info.IsDir() && stat.Nlink > 1 {
		return &fs.PathError{Op: "chmod", Path: path, Err: fmt.Errorf("file has %d names; a private file must have only this one", stat.Nlink)}
	}
	return nil
}

// EnsureEntryPrivate removes group and other access from one entry inside a
// directory that EnsureDir made private, after CheckEntry accepts it. info
// must be the Lstat of path.
func EnsureEntryPrivate(path string, info fs.FileInfo) error {
	if err := CheckEntry(path, info); err != nil {
		return err
	}
	if info.Mode().Perm()&^privateDirMode == 0 {
		return nil
	}
	return os.Chmod(path, info.Mode().Perm()&privateDirMode)
}

// checkLinkComponent refuses a symbolic link on the way to a container unless
// the current user or root owns it: anyone else who owns a link can point it
// at a tree they control.
func checkLinkComponent(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if err := ownedByCurrentUserOrRoot(path, info); err != nil {
		return &fs.PathError{Op: "checkcontainer", Path: path, Err: fmt.Errorf("symbolic link on the way to the storage location: %w", err)}
	}
	return nil
}

// checkContainerComponent refuses a directory that another user could modify
// entries in. Only the current user and root may own it, and it may not be
// writable by its group or by everyone unless the sticky bit limits removal to
// each entry's owner. A group-writable directory is refused even when the
// current user owns it: nothing here can tell whether the group holds anyone
// else, and private storage must not depend on that.
func checkContainerComponent(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return &fs.PathError{Op: "checkcontainer", Path: path, Err: syscall.ENOTDIR}
	}
	if err := ownedByCurrentUserOrRoot(path, info); err != nil {
		return &fs.PathError{Op: "checkcontainer", Path: path, Err: err}
	}
	if info.Mode()&fs.ModeSticky != 0 {
		return nil
	}
	if info.Mode().Perm()&0o002 != 0 {
		return &fs.PathError{Op: "checkcontainer", Path: path, Err: errors.New("directory is writable by every user without the sticky bit; private storage below it could be replaced by another user")}
	}
	if info.Mode().Perm()&0o020 != 0 {
		return &fs.PathError{Op: "checkcontainer", Path: path, Err: errors.New("directory is writable by a group without the sticky bit; private storage below it could be replaced by another member of that group")}
	}
	return nil
}

// ownedByCurrentUser refuses an entry another user owns: this process could
// not make it private, and using it would hand that user the data.
func ownedByCurrentUser(op string, path string, info fs.FileInfo) error {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return &fs.PathError{Op: op, Path: path, Err: errors.New("cannot determine the owner")}
	}
	if int(stat.Uid) != os.Geteuid() {
		return &fs.PathError{Op: op, Path: path, Err: fmt.Errorf("owned by uid %d, not by this process (uid %d)", stat.Uid, os.Geteuid())}
	}
	return nil
}

func ownedByCurrentUserOrRoot(path string, info fs.FileInfo) error {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return errors.New("cannot determine the owner")
	}
	if stat.Uid != 0 && int(stat.Uid) != os.Geteuid() {
		return fmt.Errorf("owned by uid %d, which is neither this process (uid %d) nor root", stat.Uid, os.Geteuid())
	}
	return nil
}

// OpenFile opens name inside root as a file only the current user can read.
// perm is the mode for a file this call creates; it must not grant group or
// other access. A file that already existed is made private before OpenFile
// returns, by removing its group and other bits, so no data written through
// the returned file is ever exposed.
func OpenFile(root *os.Root, name string, flag int, perm fs.FileMode) (*os.File, error) {
	if perm&^privateDirMode != 0 {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fmt.Errorf("private file mode %04o grants access beyond the owner", perm)}
	}
	f, err := root.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	if flag&os.O_EXCL != 0 {
		// Created by this call, with the private mode.
		return f, nil
	}
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	if mode := info.Mode().Perm(); mode&^privateDirMode != 0 {
		// Chmod the open file, not a path that could have been replaced.
		if err := f.Chmod(mode & privateDirMode); err != nil {
			f.Close()
			return nil, err
		}
	}
	return f, nil
}

// publicationPlatform on POSIX is the mode the published file receives.
type publicationPlatform struct {
	mode fs.FileMode
}

// publicationFor computes the published mode: the mode a file created the
// usual way (0666 before the process umask) would get, intersected with the
// mode of the file being replaced when there is one, and with limit. The
// usual mode is observed by creating an empty probe file: it is the only way
// to apply the umask without the process-global umask call, and the probe
// never holds data.
func publicationFor(root *os.Root, name string, probeName string, limit fs.FileMode) (publicationPlatform, error) {
	fresh, err := freshFileMode(root, probeName)
	if err != nil {
		return publicationPlatform{}, err
	}
	mode := fresh & limit.Perm()
	info, err := root.Lstat(name)
	if err == nil && info.Mode().IsRegular() {
		mode &= info.Mode().Perm()
	} else if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return publicationPlatform{}, err
	}
	return publicationPlatform{mode: mode}, nil
}

func freshFileMode(root *os.Root, probeName string) (fs.FileMode, error) {
	probe, err := root.OpenFile(probeName, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o666)
	if err != nil {
		return 0, err
	}
	info, statErr := probe.Stat()
	closeErr := probe.Close()
	removeErr := root.Remove(probeName)
	if err := errors.Join(statErr, closeErr, removeErr); err != nil {
		return 0, err
	}
	return info.Mode().Perm(), nil
}

// applyToPublished sets the published mode on the file at name once it is in
// place, through a handle checked to be the published file. Until this point
// the file kept its private mode through every rename or copy.
func (p publicationPlatform) applyToPublished(root *os.Root, name string, published fs.FileInfo) error {
	f, err := root.OpenFile(name, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	got, err := f.Stat()
	if err != nil {
		return err
	}
	if published == nil || !os.SameFile(published, got) {
		return &fs.PathError{Op: "chmod", Path: filepath.Join(root.Name(), name), Err: errPublishedReplaced}
	}
	return f.Chmod(p.mode)
}
