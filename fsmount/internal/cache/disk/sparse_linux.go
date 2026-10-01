//go:build linux

package disk

import "os"

func markSparseFile(_ *os.File) error {
	return nil
}
