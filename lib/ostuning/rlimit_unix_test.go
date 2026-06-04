//go:build linux || darwin

package ostuning

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenFileLimitRaiseTargetsBoundHugeHardLimit(t *testing.T) {
	targets := openFileLimitRaiseTargets(256, ^uint64(0))

	require.Equal(t, []uint64{PreferredOpenFileLimit, MinimumOpenFileLimit}, targets)
}

func TestOpenFileLimitRaiseTargetsDoNotLowerExistingSoftLimit(t *testing.T) {
	targets := openFileLimitRaiseTargets(PreferredOpenFileLimit+1, ^uint64(0))

	require.Empty(t, targets)
}

func TestOpenFileLimitRaiseTargetsRespectHardLimitBelowPreferred(t *testing.T) {
	targets := openFileLimitRaiseTargets(1024, 32768)

	require.Equal(t, []uint64{32768, MinimumOpenFileLimit}, targets)
}

func TestRaiseCurrentProcessOpenFileLimitFallsBackToMinimum(t *testing.T) {
	var calls []uint64
	result, err := raiseCurrentProcessOpenFileLimit(syscall.Rlimit{
		Cur: 1024,
		Max: ^uint64(0),
	}, func(next syscall.Rlimit) error {
		calls = append(calls, next.Cur)
		if next.Cur == PreferredOpenFileLimit {
			return syscall.EINVAL
		}
		return nil
	})

	require.NoError(t, err)
	require.True(t, result.Changed)
	require.Equal(t, uint64(MinimumOpenFileLimit), result.AfterSoft)
	require.Equal(t, []uint64{PreferredOpenFileLimit, MinimumOpenFileLimit}, calls)
}

func TestRaiseCurrentProcessOpenFileLimitDoesNotCallSetrlimitWhenAlreadyHighEnough(t *testing.T) {
	result, err := raiseCurrentProcessOpenFileLimit(syscall.Rlimit{
		Cur: PreferredOpenFileLimit + 1,
		Max: ^uint64(0),
	}, func(syscall.Rlimit) error {
		t.Fatal("setrlimit should not be called when the soft limit already exceeds the preferred limit")
		return nil
	})

	require.NoError(t, err)
	require.False(t, result.Changed)
	require.Equal(t, uint64(PreferredOpenFileLimit+1), result.AfterSoft)
}
