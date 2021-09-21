package file

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/stretchr/testify/assert"
	"github.com/zenthangplus/goccm"
)

func TestFS_Open(t *testing.T) {
	client, r, err := CreateClient("TestFS_Open")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	assert := assert.New(t)

	client.Upload(context.Background(), strings.NewReader("testing 3"), int64(9), files_sdk.FileBeginUploadParams{MkdirParents: lib.Bool(true), Path: filepath.Join("remotefs_test", "1.text")}, func(i int64) {}, goccm.New(1))

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
