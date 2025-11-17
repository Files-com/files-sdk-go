//go:build darwin

package fsmount

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

func mountPoint(mntRoot string, _ ...bool) (string, error) {
	mntRoot = filepath.Clean(mntRoot)
	if mntRoot == "" {
		return "", fmt.Errorf("mount point cannot be empty")
	}

	// ensure the parent directory exists
	if _, err := os.Stat(mntRoot); os.IsNotExist(err) {
		return "", fmt.Errorf("parent directory for creating mount point does not exist: %w", err)
	}

	dirs, err := os.ReadDir(mntRoot)
	if err != nil {
		return "", fmt.Errorf("failed to read mount root directory: %w", err)
	}
	// if there are directories in the mount root, check if any of them are
	// available for use as a mount point
	for _, dir := range dirs {
		if dir.IsDir() {
			candidate := filepath.Join(mntRoot, dir.Name())

			// check if the candidate is already reserved
			if _, ok := mntRegistry.get(candidate); ok {
				continue
			}

			// check if the mount point is already in use
			inUse, err := mountInUse(candidate)
			if err != nil {
				continue
			}

			if !inUse {
				// return the candidate
				return candidate, nil
			}
		}
	}

	// no existing directories are available,

	// find a suitable name for the mount point that is not already in use
	// by checking for directories named A-Z in the mount root
	// and checking if they are already mounted
	// and checking if they are already reserved in this process
	for l := 'A'; l <= 'Z'; l++ {
		subdir := string(l)

		// generate a candidate mount point
		candidate := filepath.Join(mntRoot, subdir)

		// check if the candidate is already reserved
		if _, ok := mntRegistry.get(candidate); ok {
			continue
		}

		// check if the mount point is already in use
		inUse, err := mountInUse(candidate)
		if err != nil {
			continue
		}

		if !inUse {
			// make sure it exists
			if _, err := os.Stat(candidate); os.IsNotExist(err) {
				// create the directory if it doesn't exist
				if err := os.Mkdir(candidate, 0o700); err != nil {
					return "", fmt.Errorf("failed to create mount point directory: %w", err)
				}
			}

			// return the candidate
			return candidate, nil
		}
	}

	return "", fmt.Errorf("no available mount points")
}

func mountOpts(params MountParams) []string {
	opts := defaultMountOpts(params)
	// uses the smb implementation from fuse-t/go-nfsv4
	opts = append(opts, "-o", "backend=smb")

	// sets the name that is displayed in the Finder sidebar
	opts = append(opts, "-o", "location=Files")
	return opts
}

// macOS safe-save temp DIRECTORIES (Word/Pages/TextEdit)
// Use trailing slash so only directories match.
// *.sb-*/
// Office scratch temp files (occasionally appear on macOS)
// ~WR*.tmp
// Per-volume Trash directories (external/network volumes)
// .Trash-*/
func additionalIgnorePatterns() []string {
	return []string{
		"*.doc.sb-*",
		"*.docx.sb-*",
		"*.dotx.sb-*",
		"*.ppt.sb-*",
		"*.pptm.sb-*",
		"*.pptx.sb-*",
		"*.pdf.sb-*",
		"*.rtf.sb-*",
		"*.csv.sb-*",
		"*.xls.sb-*",
		"*.xlsb.sb-*",
		"*.xlsm.sb-*",
		"*.xlsx.sb-*",
		".Trash-*/",
		"~WR*.tmp",
		"*.smbdelete*",
	}
}

// mountInUse reports whether base/rel is currently the mount point of some file system.
func mountInUse(mntPnt string) (bool, error) {
	p := filepath.Clean(mntPnt)
	// Fast path: statfs on the path and compare mntonname.
	var s unix.Statfs_t
	if err := unix.Statfs(p, &s); err == nil {
		mntonname := cArrayToString(s.Mntonname[:])
		if filepath.Clean(mntonname) == p {
			return true, nil
		}
	}
	// Fallback: enumerate mounts (covers cases where the dir doesn't exist yet).
	n, err := unix.Getfsstat(nil, unix.MNT_NOWAIT)
	if err != nil {
		return false, err
	}
	buf := make([]unix.Statfs_t, n)
	_, err = unix.Getfsstat(buf, unix.MNT_NOWAIT)
	if err != nil {
		return false, err
	}
	for i := range buf {
		if filepath.Clean(cArrayToString(buf[i].Mntonname[:])) == p {
			return true, nil
		}
	}
	return false, nil
}

func cArrayToString(arr []byte) string {
	n := 0
	for n < len(arr) && arr[n] != 0 {
		n++
	}
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte(arr[i])
	}
	return string(b)
}
