//go:build !darwin

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
		assert.Equal(t, filepath.Join(dir, ".~files-cli.you-wont-find-me.download"), path.String())
	})

	// Inside a temp directory the file is identified by its whole path, and
	// still carries its own name.
	t.Run("it supports a temp path", func(t *testing.T) {
		dir := t.TempDir()
		tempDir := t.TempDir()
		path, err := tmpDownloadPath(explicitDestination(filepath.Join(dir, "find-me")), tempDir)
		require.NoError(t, err)
		assert.Equal(t, tempDir, filepath.Dir(path.String()))
		element := filepath.Base(path.String())
		assert.True(t, strings.HasPrefix(element, ".~files-cli~.path-"), element)
		assert.True(t, strings.HasSuffix(element, ".find-me.download"), element)
	})
}

// writeLegacyExternalStage leaves what a client from before path identities
// left in an external temp directory for final: a file under the element name.
func writeLegacyExternalStage(t *testing.T, element string, _ destinationPath, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(element, []byte(content), 0600))
}

func readLegacyExternalStage(t *testing.T, element string, _ destinationPath) string {
	t.Helper()
	content, err := os.ReadFile(element)
	require.NoError(t, err)
	return string(content)
}
