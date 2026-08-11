//go:build (linux || windows) && !filescomfs_debug
// +build linux windows
// +build !filescomfs_debug

package fsmount

func (fs *mountRegistry) startPprof() {
	// no-op if not built with the filescomfs_debug build tag
}

func (reg *mountRegistry) stopPprof() {
	// no-op if not built with the filescomfs_debug build tag
}
