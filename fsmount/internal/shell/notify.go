//go:build !windows

package shell

// NotifyUpdatedDir is a no-op on non-Windows platforms.
func NotifyUpdatedDir(path string) error {
	return nil
}
