//go:build linux

package fsmount

import (
	"fmt"
	"os"
	path_lib "path"
	"regexp"
)

// libfuse hides a deleted file that is still open by renaming it, in the same
// directory, to ".fuse_hidden" followed by two 8-digit hex counters.
var libfuseHiddenNameRe = regexp.MustCompile(`^\.fuse_hidden[0-9a-f]{16}$`)

func mountPoint(mountPoint string, _ bool) (string, error) {
	// TODO: build a path to the mount point that is OS specific
	// for now, require that the mount point is provided and
	// exists.
	if mountPoint == "" {
		return "", fmt.Errorf("mount point cannot be empty")
	}
	if _, err := os.Stat(mountPoint); os.IsNotExist(err) {
		return "", fmt.Errorf("mount point does not exist: %w", err)
	}
	return mountPoint, nil
}

// libfuse3 rejects the hard_remove and volname options from defaultMountOpts,
// and cgofuse already passes -f. With auto_unmount, fusermount3 3.15 and later
// removes the mount when the process exits without unmounting, such as after a
// crash or SIGKILL, instead of leaving a mount point that fails with "Transport
// endpoint is not connected". Older versions check the mount point as root,
// which FUSE denies, and leave the mount in place.
func mountOpts(params MountParams) []string {
	opts := []string{"-o", "attr_timeout=1", "-o", "fsname=Files.com", "-o", "auto_unmount"}
	if params.DebugFuse {
		opts = append(opts, "-o", "debug")
	}
	return opts
}

func isLibfuseHiddenRename(oldpath, newpath string) bool {
	return path_lib.Dir(oldpath) == path_lib.Dir(newpath) && libfuseHiddenNameRe.MatchString(path_lib.Base(newpath))
}

// LibreOffice lock files
// .~lock.*#
func additionalIgnorePatterns() []string {
	return []string{
		".~lock.*#",
	}
}
