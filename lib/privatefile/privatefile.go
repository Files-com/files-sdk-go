// Package privatefile creates and maintains local files and directories that
// only the current user can read.
//
// The SDK keeps data on the local disk on its own behalf: cached file
// contents, scratch files a mounted file system stores outside Files.com,
// working copies of files being uploaded, and downloads that are still being
// written. None of that is meant for other users of the same computer, so it
// is created private and kept private until it is deleted or, for a download,
// deliberately published.
//
// "Private" means the process owner on POSIX systems (mode 0700 for
// directories and 0600 for files), and on Windows an access control list that
// grants access only to the owning account, SYSTEM and Administrators. On
// Windows the mode bits Go accepts do not restrict anything, so the Windows
// implementation sets discretionary access control lists instead.
package privatefile

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// MkdirTemp creates a new private directory with a unique name inside dir and
// returns its path. The name is pattern with the final "*" replaced by a
// random string, like os.MkdirTemp, but the directory is private on every
// platform, including Windows where os.MkdirTemp inherits the parent's access.
func MkdirTemp(dir string, pattern string) (string, error) {
	prefix, suffix, err := tempPattern(pattern)
	if err != nil {
		return "", err
	}
	if dir == "" {
		dir = os.TempDir()
	}
	for attempt := 0; ; attempt++ {
		path := filepath.Join(dir, prefix+rand.Text()+suffix)
		err := Mkdir(path)
		if err == nil {
			return path, nil
		}
		if !errors.Is(err, fs.ErrExist) || attempt >= 100 {
			return "", err
		}
	}
}

// tempPattern splits a MkdirTemp pattern around its final "*".
func tempPattern(pattern string) (string, string, error) {
	prefix, suffix, _ := strings.Cut(pattern, "*")
	if strings.ContainsAny(prefix, `/\`) || strings.ContainsAny(suffix, `/\`) {
		return "", "", &fs.PathError{Op: "mkdirtemp", Path: pattern, Err: errors.New("pattern contains path separator")}
	}
	return prefix, suffix, nil
}

// Publication is the access a privately staged file receives when it is
// published at its final path: what a freshly created file would get there,
// never more than an existing file at that path already had, and never more
// than the caller prepared for the staged file.
//
// The staged file keeps its private access through every rename or copy that
// moves it into place. Only once it is at its final path does it receive the
// published access, through a handle whose identity is checked, so bytes are
// never readable by others while publication can still fail, and access is
// never applied to a file that took the staged file's place.
type Publication struct {
	platform publicationPlatform
}

// PublicationFor determines the access a file published as name inside root
// receives. probeName is a name inside the directory of name that nothing
// else uses; on POSIX systems an empty file is created and removed under that
// name to observe the process umask without the process-global umask call.
// limit is the most access the published file may have on account of the
// staged file itself: fs.ModePerm when the stage was created by the caller's
// own code and carries no wishes of its own, or the mode a caller-prepared
// stage had before it was made private.
func PublicationFor(root *os.Root, name string, probeName string, limit fs.FileMode) (Publication, error) {
	platform, err := publicationFor(root, name, probeName, limit)
	if err != nil {
		return Publication{}, err
	}
	return Publication{platform: platform}, nil
}

// ApplyToPublished gives the file now at name inside root the published
// access. published identifies the file that was moved there; a different
// file at that name is left alone and reported.
func (p Publication) ApplyToPublished(root *os.Root, name string, published fs.FileInfo) error {
	return p.platform.applyToPublished(root, name, published)
}

// errPublishedReplaced reports that the file at the published name is not the
// file that was published there.
var errPublishedReplaced = errors.New("privatefile: the published file was replaced by another file")

// CheckContainer reports whether private storage can safely be created inside
// dir, a directory the caller selected. Private storage is only as private as
// the directories above it: a user who can rename or delete entries in dir or
// in any directory above it can replace the private storage with a tree they
// control. So dir and every directory above it must be owned by the current
// user or the system, and must not let other users modify their entries.
// World-writable directories with the sticky bit, such as the system
// temporary directory, are accepted: there only the owner of an entry can
// remove it.
//
// Every name on the way is examined as it is given and as it resolves: a
// directory component must itself be safe, and a link component must be owned
// by the current user or the system, must sit in a safe directory, and must
// lead, through any further links, only through safe directories. Otherwise a
// user who can replace the link, or a directory the link passes through, could
// redirect the path.
//
// A configuration that fails this check is refused with a descriptive error
// rather than used, because there is no way to make storage inside it private.
func CheckContainer(dir string) error {
	return checkContainer(dir, unheld)
}

// containerHold says what a caller holds open while it uses a container, which
// on some file systems keeps directories in place regardless of their access.
type containerHold int

const (
	// unheld: nothing is held; every directory is checked in full.
	unheld containerHold = iota
	// heldBelow: a directory below the one being checked is held open.
	heldBelow
	// heldHere: the directory being checked is itself held open.
	heldHere
)

// MkdirTempHeld creates a private directory inside dir like MkdirTemp, after
// checking dir like CheckContainer, and returns it open. An empty dir means
// the system temporary directory, as for MkdirTemp. The caller keeps the
// returned file open for as long as it uses the directory, and closes it
// before removing the directory.
//
// On Windows the directory is created and opened in one step, inside the
// directory dir was opened as, with its private access already applied and
// without sharing delete access, so there is no moment at which it exists
// unprotected or could be swapped for another directory. On a local NTFS
// volume, Windows also does not rename or move a directory while anything
// inside it is open. So while it is held, no account can take away the
// directory or any directory above it, whatever access it has to them, and
// those directories are not checked: only links on the way still are, and the
// directories above a link, which the held directory is not inside. Elsewhere,
// and on POSIX systems, dir is checked in full, like CheckContainer.
func MkdirTempHeld(dir string, pattern string) (string, *os.File, error) {
	prefix, suffix, err := tempPattern(pattern)
	if err != nil {
		return "", nil, err
	}
	if dir == "" {
		dir = os.TempDir()
	}
	container, err := openContainer(dir)
	if err != nil {
		return "", nil, err
	}
	defer container.Close()
	hold := unheld
	if heldDirKeepsParents(container) {
		hold = heldBelow
	}
	if err := checkContainer(dir, hold); err != nil {
		return "", nil, err
	}
	opened, err := container.Stat()
	if err != nil {
		return "", nil, err
	}
	info, err := os.Stat(dir)
	if err != nil {
		return "", nil, err
	}
	if !os.SameFile(opened, info) {
		return "", nil, &fs.PathError{Op: "mkdirtemp", Path: dir, Err: errors.New("directory was replaced while it was being checked")}
	}
	for attempt := 0; ; attempt++ {
		name := prefix + rand.Text() + suffix
		held, err := mkdirHeld(container, name)
		if err == nil {
			return filepath.Join(dir, name), held, nil
		}
		if !errors.Is(err, fs.ErrExist) || attempt >= 100 {
			return "", nil, err
		}
	}
}

// CheckContainerHeld checks dir like CheckContainer for a caller that keeps
// dir itself open for as long as it uses it. If dir passes CheckContainer,
// nothing needs to be held: the returned file is nil. Otherwise, on a local
// NTFS volume, dir is opened and accepted when the directories above it are
// left out of the check and deleting or renaming dir itself is allowed: while
// dir is open, no account can take it or any directory above it away (see
// MkdirTempHeld). The rights to remove entries from dir or to change its
// security are still refused. The caller keeps the returned file open while it
// uses dir and closes it afterwards.
func CheckContainerHeld(dir string) (*os.File, error) {
	unheldErr := CheckContainer(dir)
	if unheldErr == nil {
		return nil, nil
	}
	held, err := os.Open(dir)
	if err != nil {
		return nil, unheldErr
	}
	if !heldDirKeepsParents(held) {
		held.Close()
		return nil, unheldErr
	}
	if err := checkContainer(dir, heldHere); err != nil {
		held.Close()
		return nil, err
	}
	opened, err := held.Stat()
	if err != nil {
		held.Close()
		return nil, err
	}
	info, err := os.Stat(dir)
	if err != nil {
		held.Close()
		return nil, err
	}
	if !os.SameFile(opened, info) {
		held.Close()
		return nil, &fs.PathError{Op: "checkcontainer", Path: dir, Err: errors.New("directory was replaced while it was being checked")}
	}
	return held, nil
}

// EnsureDirIn makes name, inside the directory root that CheckContainerHeld
// returned, a private directory, going through the open root rather than its
// path. A missing directory is created private. An existing one must already
// be private, owned by this user or the system, and not a link: one that is
// not could have been put there by another account, which might still hold it
// open, so it is refused rather than made private.
func EnsureDirIn(root *os.File, name string) error {
	return ensureDirIn(root, name)
}

// checkContainer is CheckContainer with what the caller holds open.
func checkContainer(dir string, hold containerHold) error {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return &fs.PathError{Op: "checkcontainer", Path: dir, Err: errors.New("not a directory")}
	}
	return checkContainerPath(abs, 0, hold)
}

// maxLinkDepth bounds the chains of links CheckContainer follows.
const maxLinkDepth = 40

// checkContainerPath checks every component of path, following each link
// component into the path it names. A held directory keeps in place only the
// directories it is actually inside: the path a link names, but not the
// directories above the link itself, so those are checked as unheld.
func checkContainerPath(path string, depth int, hold containerHold) error {
	if depth > maxLinkDepth {
		return &fs.PathError{Op: "checkcontainer", Path: path, Err: errors.New("too many levels of symbolic links")}
	}
	for _, component := range pathComponents(path) {
		info, err := os.Lstat(component)
		if err != nil {
			return err
		}
		if !isLink(info) {
			if err := checkContainerComponent(component, hold); err != nil {
				return err
			}
			if hold == heldHere {
				hold = heldBelow
			}
			continue
		}
		if err := checkLinkComponent(component); err != nil {
			return err
		}
		target, err := os.Readlink(component)
		if err != nil {
			return &fs.PathError{Op: "checkcontainer", Path: component, Err: err}
		}
		if !filepath.IsAbs(target) {
			target = filepath.Join(filepath.Dir(component), target)
		}
		if err := checkContainerPath(target, depth+1, hold); err != nil {
			return err
		}
		hold = unheld
	}
	return nil
}

// pathComponents returns path and every directory above it, nearest first,
// ending at the root of the volume.
func pathComponents(path string) []string {
	var components []string
	path = filepath.Clean(path)
	for {
		components = append(components, path)
		parent := filepath.Dir(path)
		if parent == path {
			return components
		}
		path = parent
	}
}

func errNotDirectory(path string, info fs.FileInfo) error {
	kind := "not a directory"
	if info.Mode()&fs.ModeSymlink != 0 {
		kind = "a symbolic link"
	} else if info.Mode()&fs.ModeIrregular != 0 {
		kind = "a reparse point"
	}
	return &fs.PathError{Op: "ensuredir", Path: path, Err: fmt.Errorf("must be a directory that this process owns, but it is %s", kind)}
}

// lstatDirectory returns the Lstat of path when it is a real directory, and
// an error otherwise. Symbolic links and other reparse points are refused
// because an attacker who can write the parent could point one at a tree they
// control or at unrelated data.
func lstatDirectory(path string) (fs.FileInfo, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() || info.Mode()&(fs.ModeSymlink|fs.ModeIrregular) != 0 {
		return nil, errNotDirectory(path, info)
	}
	return info, nil
}

// isLink reports whether info describes a symbolic link or, on Windows, any
// other reparse point.
func isLink(info fs.FileInfo) bool {
	return info.Mode()&(fs.ModeSymlink|fs.ModeIrregular) != 0
}
