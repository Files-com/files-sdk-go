//go:build windows

package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Reserved device names without an extension are rejected on every Windows
// version. A reserved name with an extension, such as NUL.txt, is a valid file
// name since Windows 11 and is decided by the running system, so it is not
// asserted here.
func TestDownloadFolder_rejectsWindowsPathForms(t *testing.T) {
	cases := map[string]string{
		"backslash traversal":   `reports/..\..\outside.txt`,
		"drive prefix":          `reports/C:outside.txt`,
		"alternate data stream": `reports/file.txt:stream`,
		"UNC path":              `reports/\\server\share\outside.txt`,
		"device name":           `reports/CON`,
		"lowercase device name": `reports/nul`,
		"serial port name":      `reports/COM1`,
		"console handle name":   `reports/CONIN$`,
	}
	for name, serverPath := range cases {
		t.Run(name, func(t *testing.T) {
			setup := NewTestSetup()
			defer setup.TearDown()
			setup.MapFS["reports"] = mapDir("reports")
			setup.MapFS["reports/entry"] = mapFileAt("reports/entry", serverPath, 10)
			destination := filepath.Join(setup.tempDir, "downloads")
			setup.DownloaderParams = DownloaderParams{RemotePath: "reports/", LocalPath: destination + string(os.PathSeparator)}

			job := setup.Call()

			rejected := findStatus(t, job, serverPath)
			assert.Equal(t, status.Errored, rejected.Status())
			require.ErrorContains(t, rejected.Err(), "not a valid local path")
			entries, err := os.ReadDir(destination)
			require.NoError(t, err)
			assert.Empty(t, entries, "nothing may be written for the rejected path")
			assert.True(t, job.All(status.Ended...))
		})
	}
}
