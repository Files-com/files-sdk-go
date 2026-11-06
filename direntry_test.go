package files_sdk

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDirEntryReportsDirectoryType(t *testing.T) {
	for _, kind := range []string{"directory", "file"} {
		t.Run(kind, func(t *testing.T) {
			entry := DirEntry{File: File{Type: kind}}
			info, err := entry.Info()
			require.NoError(t, err)
			require.Equal(t, entry.IsDir(), entry.Type().IsDir())
			require.Equal(t, info.IsDir(), info.Mode().IsDir())
			require.Equal(t, info.Mode().Type(), entry.Type())
			require.Zero(t, entry.Type()&fs.ModePerm, "DirEntry.Type contains type bits only")
		})
	}
}
