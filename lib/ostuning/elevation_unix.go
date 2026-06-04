//go:build linux || darwin

package ostuning

import "os"

func currentProcessElevated() bool {
	return os.Geteuid() == 0
}
