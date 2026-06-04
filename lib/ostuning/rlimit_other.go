//go:build !linux && !darwin

package ostuning

// RaiseCurrentProcessOpenFileLimit is unsupported on platforms where the SDK
// does not manage POSIX nofile resource limits.
func RaiseCurrentProcessOpenFileLimit() (OpenFileLimitResult, error) {
	return OpenFileLimitResult{}, nil
}
