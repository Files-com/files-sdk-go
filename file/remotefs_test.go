package file

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestFS_Open(t *testing.T) {
	client, r, err := CreateClient("TestFS_Open")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	assert := assert.New(t)
	client.Upload(
		UploadWithReader(strings.NewReader("testing 3")),
		UploadWithSize(9),
		UploadWithDestinationPath(filepath.Join("remotefs_test", "1.text")),
	)

	fs := (&FS{}).Init(client.Config, true)
	fs = fs.WithContext(context.TODO()).(*FS)
	f, err := fs.Open("remotefs_test")
	assert.NoError(err)
	rf, ok := f.(*ReadDirFile)
	assert.True(ok)
	entry, err := rf.ReadDir(0)
	assert.NoError(err)
	assert.Equal(1, len(entry))
	assert.False(entry[0].IsDir())
	info, err := entry[0].Info()
	assert.NoError(err)
	assert.Equal("1.text", info.Name())
	fsFile, ok := entry[0].(*File)
	assert.True(ok)
	buf := make([]byte, 8)
	_, err = fsFile.Read(buf)
	assert.NoError(err)
	err = fsFile.Close()
	assert.NoError(err)
	assert.Equal("testing ", string(buf))

	buf = make([]byte, 9)
	_, err = fsFile.Read(buf)
	assert.NoError(err)
	err = fsFile.Close()
	assert.NoError(err)

	assert.Equal("testing 3", string(buf))

	_, err = fs.ReadDir(".")
	require.NoError(t, err)

	_, err = fs.Open(".")
	require.NoError(t, err)
}
