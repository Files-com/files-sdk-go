//go:build windows

package file

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var getShortPathNameW = syscall.NewLazyDLL("kernel32.dll").NewProc("GetShortPathNameW")

// shortName is the 8.3 name NTFS keeps for path, or "" when the volume does not
// keep one.
func shortName(t *testing.T, path string) string {
	t.Helper()
	wide, err := syscall.UTF16PtrFromString(path)
	require.NoError(t, err)
	buf := make([]uint16, 1024)
	n, _, errno := getShortPathNameW.Call(uintptr(unsafe.Pointer(wide)), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if n == 0 {
		t.Skipf("this volume does not report an 8.3 name: %v", errno)
	}
	short := filepath.Base(syscall.UTF16ToString(buf[:n]))
	if short == filepath.Base(path) {
		t.Skip("8.3 name creation is disabled on this volume")
	}
	return short
}

// NTFS keeps a second, 8.3 name for every name too long to be one, and that
// name reaches the same file. A temporary download's second name is outside the
// reserved namespace, so without resolving it an ordinary-looking name could be
// used to reach an unfinished download.
func Test_localPathReachesTempDownload_throughAnNTFSShortName(t *testing.T) {
	dir := t.TempDir()
	stage := filepath.Join(dir, ".~files-cli.victim.bin.download")
	require.NoError(t, os.WriteFile(stage, []byte("partial"), 0600))
	alias := filepath.Join(dir, shortName(t, stage))

	aliasInfo, err := os.Stat(alias)
	require.NoError(t, err)
	stageInfo, err := os.Stat(stage)
	require.NoError(t, err)
	require.True(t, os.SameFile(aliasInfo, stageInfo), "the 8.3 name is expected to reach the same file")

	assert.False(t, isReservedTempDownloadPath(alias), "the 8.3 name is not reserved by its spelling")
	reaches, err := localPathReachesTempDownload(alias)
	require.NoError(t, err)
	assert.True(t, reaches, "the 8.3 name reaches the temporary download and must not be transferred")
}

// An ordinary file whose name merely looks like an 8.3 name stays transferable:
// only a name that actually reaches a temporary download is refused.
func Test_localPathReachesTempDownload_keepsOrdinaryShortNames(t *testing.T) {
	dir := t.TempDir()
	ordinary := filepath.Join(dir, "~FILES~1.DOW")
	require.NoError(t, os.WriteFile(ordinary, []byte("an ordinary file"), 0600))

	reaches, err := localPathReachesTempDownload(ordinary)
	require.NoError(t, err)
	assert.False(t, reaches)

	// A destination that does not exist yet cannot be another name for
	// anything, so it is kept as it is rather than refused.
	reaches, err = localPathReachesTempDownload(filepath.Join(dir, "not-created-yet.bin"))
	require.NoError(t, err)
	assert.False(t, reaches)
}
