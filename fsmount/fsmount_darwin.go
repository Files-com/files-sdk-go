//go:build darwin

package fsmount

import (
	"fmt"
	"os"
)

func mountPoint(mountPoint string) (string, error) {
	// TODO: build a path to the mount point that is OS specific. For now,
	// require that the mount point is provided and exists.
	if mountPoint == "" {
		return "", fmt.Errorf("mount point cannot be empty")
	}
	if _, err := os.Stat(mountPoint); os.IsNotExist(err) {
		return "", fmt.Errorf("mount point does not exist: %w", err)
	}
	return mountPoint, nil
}

func mountOpts(params MountParams) []string {
	opts := defaultMountOpts(params)
	// uses the smb implementation from fuse-t/go-nfsv4
	opts = append(opts, "-o", "backend=smb")

	// sets the name that is displayed in the Finder sidebar
	opts = append(opts, "-o", "location=Files")
	return opts
}
