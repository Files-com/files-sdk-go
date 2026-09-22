//go:build darwin

package file

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_tmpDownloadPath(t *testing.T) {
	t.Run("base case", func(t *testing.T) {
		dir := t.TempDir()
		path, err := tmpDownloadPath(explicitDestination(filepath.Join(dir, "you-wont-find-me")), "")
		require.NoError(t, err)
		assert.Equal(t, filepath.Join(dir, ".~files-cli.you-wont-find-me.download/you-wont-find-me"), path.String())
	})

	// Inside a temp directory the folder is identified by the file's whole
	// path, and still holds the file under its own name.
	t.Run("it supports a temp path", func(t *testing.T) {
		dir := t.TempDir()
		tempDir := t.TempDir()
		path, err := tmpDownloadPath(explicitDestination(filepath.Join(dir, "find-me")), tempDir)
		require.NoError(t, err)
		assert.Equal(t, "find-me", filepath.Base(path.String()))
		folder := filepath.Base(filepath.Dir(path.String()))
		assert.Equal(t, tempDir, filepath.Dir(filepath.Dir(path.String())))
		assert.True(t, strings.HasPrefix(folder, ".~files-cli~.path-"), folder)
		assert.True(t, strings.HasSuffix(folder, ".find-me.download"), folder)
	})
}

// writeLegacyExternalStage leaves what a client from before path identities
// left in an external temp directory for final: a folder under the element
// name, holding the file under its own name.
func writeLegacyExternalStage(t *testing.T, element string, final destinationPath, content string) {
	t.Helper()
	require.NoError(t, os.Mkdir(element, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(element, final.base()), []byte(content), 0600))
}

func readLegacyExternalStage(t *testing.T, element string, final destinationPath) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(element, final.base()))
	require.NoError(t, err)
	return string(content)
}
