package file

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFS_Open(t *testing.T) {
	client, r, err := CreateClient("TestFS_Open")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	assert := assert.New(t)
	client.UploadIO(context.Background(), UploadIOParams{Path: filepath.Join("remotefs_test", "1.text"), Reader: strings.NewReader("testing 3"), Size: int64(9)})

	fs := FS{}.Init(client.Config)
	fs = fs.WithContext(context.TODO())
	f, err := fs.Open("remotefs_test")
	assert.NoError(err)
	rf, ok := f.(*ReadDirFile)
	assert.Equal(true, ok)
	entry, err := rf.ReadDir(0)
	assert.NoError(err)
	assert.Equal(1, len(entry))
	assert.Equal(false, entry[0].IsDir())
	info, err := entry[0].Info()
	assert.NoError(err)
	assert.Equal("1.text", info.Name())
	fsFile, ok := entry[0].(File)
	assert.Equal(true, ok)
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
}
