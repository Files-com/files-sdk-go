//go:build !linux && !darwin && !windows

package ostuning

func currentProcessElevated() bool {
	return false
}
