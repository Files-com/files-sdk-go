//go:build windows

package fsmount

// platformFileMode removes group and world permissions before WinFSP maps the
// FUSE mode to a Windows ACL. Without this, the UNC mount grants other local
// users access through the group and Everyone ACL entries.
func platformFileMode(mode uint32) uint32 {
	return mode &^ 0o077
}
