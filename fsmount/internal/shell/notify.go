//go:build linux

package shell

// NotifyUpdatedDir is a no-op on Linux.
func NotifyUpdatedDir(path string) error {
	return nil
}
