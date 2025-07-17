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
	opts = append(opts, "-o", "nobrowse")
	// TODO: Figure out how to provide the volume icon correctly to fuse
	// opts = append(opts, "-o", "modules=volicon")
	// opts = append(opts, "-o", "iconpath=application.icns")
	// opts = append(opts, "-o", "volicon="+params.IconPath)
	return opts
}
