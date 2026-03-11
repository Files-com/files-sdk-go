//go:build !windows

package fsmount

import "os"

func openLocalFile(path string, flags int, mode os.FileMode) (*os.File, error) {
	return os.OpenFile(path, flags, mode)
}
