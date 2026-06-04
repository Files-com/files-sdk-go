//go:build linux || darwin

package ostuning

import "syscall"

// RaiseCurrentProcessOpenFileLimit raises the current process soft nofile limit
// to the preferred high-throughput limit when the platform supports POSIX
// resource limits.
func RaiseCurrentProcessOpenFileLimit() (OpenFileLimitResult, error) {
	var limit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit); err != nil {
		return OpenFileLimitResult{}, err
	}

	return raiseCurrentProcessOpenFileLimit(limit, func(next syscall.Rlimit) error {
		return syscall.Setrlimit(syscall.RLIMIT_NOFILE, &next)
	})
}

func raiseCurrentProcessOpenFileLimit(limit syscall.Rlimit, setrlimit func(syscall.Rlimit) error) (OpenFileLimitResult, error) {
	result := OpenFileLimitResult{
		Supported:  true,
		BeforeSoft: limit.Cur,
		BeforeHard: limit.Max,
		AfterSoft:  limit.Cur,
	}

	var lastErr error
	for _, target := range openFileLimitRaiseTargets(limit.Cur, limit.Max) {
		next := limit
		next.Cur = target
		if err := setrlimit(next); err != nil {
			lastErr = err
			continue
		}

		result.AfterSoft = target
		result.Changed = target != limit.Cur
		return result, nil
	}

	return result, lastErr
}

func openFileLimitRaiseTargets(current uint64, hard uint64) []uint64 {
	target := minUint64(hard, uint64(PreferredOpenFileLimit))
	if target <= current {
		return nil
	}

	targets := []uint64{target}
	fallback := minUint64(hard, uint64(MinimumOpenFileLimit))
	if fallback > current && fallback < target {
		targets = append(targets, fallback)
	}
	return targets
}

func minUint64(a uint64, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
