//go:build linux

package fsmount

func platformFileMode(mode uint32) uint32 {
	return mode
}
