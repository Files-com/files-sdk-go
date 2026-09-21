//go:build !darwin

package file

import (
	"path/filepath"
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

	t.Run("it supports a temp path", func(t *testing.T) {
		dir := t.TempDir()
		tempDir := t.TempDir()
		path, err := tmpDownloadPath(explicitDestination(filepath.Join(dir, "find-me")), tempDir)
		require.NoError(t, err)
		assert.Equal(t, filepath.Join(tempDir, ".~files-cli.find-me.download"), path.String())
	})
}
