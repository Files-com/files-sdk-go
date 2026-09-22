//go:build windows

package file

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows"
)

// canonicalLocalPath returns path spelled the way the filesystem stores it.
//
// NTFS keeps a second, 8.3 name for every name that is not already an 8.3 name,
// and opening either one reaches the same file. A temporary download's second
// name is outside the reserved namespace, so without this an ordinary-looking
// name could reach a temporary download. Only the part of path that exists can
// be resolved; what does not exist yet cannot be another name for anything, so
// it is kept as it is.
func canonicalLocalPath(path string) (string, error) {
	if path == "" {
		return path, nil
	}
	resolved, rest, err := longestResolvablePrefix(path)
	if err != nil {
		return "", err
	}
	if rest == "" {
		return resolved, nil
	}
	return filepath.Join(resolved, rest), nil
}

// longestResolvablePrefix resolves the longest leading part of path that the
// filesystem knows, and returns it with the remainder that does not exist.
func longestResolvablePrefix(path string) (resolved string, rest string, err error) {
	remaining := path
	for {
		long, err := longPathName(remaining)
		if err == nil {
			return long, rest, nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			// Refuse rather than accept a name we could not resolve.
			return "", "", err
		}
		parent, element := filepath.Split(remaining)
		element = strings.TrimRight(element, string(filepath.Separator))
		if element == "" {
			// A volume root or an empty element: nothing further to resolve.
			return path, "", nil
		}
		rest = filepath.Join(element, rest)
		remaining = strings.TrimRight(parent, string(filepath.Separator))
		if remaining == "" || remaining == filepath.VolumeName(path) {
			return path, "", nil
		}
	}
}

func longPathName(path string) (string, error) {
	wide, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return "", err
	}
	buf := make([]uint16, windows.MAX_LONG_PATH)
	for {
		n, err := windows.GetLongPathName(wide, &buf[0], uint32(len(buf)))
		if err != nil {
			return "", &fs.PathError{Op: "GetLongPathName", Path: path, Err: err}
		}
		if n <= uint32(len(buf)) {
			return windows.UTF16ToString(buf[:n]), nil
		}
		buf = make([]uint16, n)
	}
}
