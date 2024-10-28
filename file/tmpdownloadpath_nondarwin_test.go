//go:build !darwin

package file

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_tmpDownloadPath(t *testing.T) {
	t.Run("base case", func(t *testing.T) {
		dir := t.TempDir()
		path, err := tmpDownloadPath(filepath.Join(dir, "you-wont-find-me"), "")
		require.NoError(t, err)
		assert.Equal(t, filepath.Join(dir, "you-wont-find-me.download"), path)
	})

	t.Run("it increments a number", func(t *testing.T) {
		dir := t.TempDir()
		file, err := os.Create(filepath.Join(dir, "find-me.download"))
		_, err = file.Write([]byte("hello"))
		require.NoError(t, err)
		err = file.Close()
		if err != nil {
			panic(err)
		}
		path, err := tmpDownloadPath(filepath.Join(dir, "find-me"), "")
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf(filepath.Join(dir, "find-me (1).download")), path, "it increments a number")
	})

	t.Run("it increments a number lots of times", func(t *testing.T) {
		dir := t.TempDir()
		for i := 0; i < 11; i++ {
			path, err := tmpDownloadPath(filepath.Join(dir, "find-me"), "")
			require.NoError(t, err)
			file, err := os.Create(path)
			require.NoError(t, err)
			file.Close()
		}

		path, err := tmpDownloadPath(filepath.Join(dir, "find-me"), "")
		require.NoError(t, err)
		assert.NotEqual(t, fmt.Sprintf(filepath.Join(dir, "find-me (11).download")), path)
	})

	t.Run("it supports a temp path", func(t *testing.T) {
		dir := t.TempDir()
		tempDir := t.TempDir()
		path, err := tmpDownloadPath(filepath.Join(dir, "find-me"), tempDir)
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf(filepath.Join(tempDir, "find-me.download")), path)
	})
}
