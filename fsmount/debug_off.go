//go:build !filescomfs_debug
// +build !filescomfs_debug

package fsmount

func (fs *Filescomfs) startPprof() {
	// no-op if not built with the filescomfs_debug build tag
}
