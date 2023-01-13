package file

import (
	"context"
	"io/fs"
	"sync"
	"testing"
	"testing/fstest"

	"github.com/Files-com/files-sdk-go/v2/directory"

	"github.com/Files-com/files-sdk-go/v2/ignore"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/file/status"

	"github.com/stretchr/testify/assert"
)

type MockUploader struct {
	files_sdk.File
	findError   files_sdk.ResponseError
	uploadError error
}

func (m *MockUploader) UploadIO(context.Context, UploadIOParams) (files_sdk.File, files_sdk.FileUploadPart, Parts, []error, error) {
	return files_sdk.File{}, files_sdk.FileUploadPart{}, Parts{}, []error{}, m.uploadError
}

func (m *MockUploader) Find(context.Context, files_sdk.FileFindParams, ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	return m.File, m.findError
}

func Test_skipOrIgnore(t *testing.T) {
	assert := assert.New(t)
	job := status.Job{Logger: (&files_sdk.Config{}).Logger()}.Init()
	job.GitIgnore, _ = ignore.New()
	job.Params = UploaderParams{}
	uploadStatus := &UploadStatus{job: job, Mutex: &sync.RWMutex{}, file: files_sdk.File{Path: "test"}, remotePath: "test"}
	uploadStatus.Job().Add(uploadStatus)
	uploader := &MockUploader{}
	uploadStatus.Uploader = uploader
	var progressReportError error
	uploadStatus.job.EventsReporter = Reporter(func(s status.File) {
		progressReportError = s.Err
	})

	mockFs := make(fstest.MapFS)
	job.RemoteFs = mockFs

	// sync not enabled and sizes do match
	mockFs["test"] = &fstest.MapFile{
		Sys: files_sdk.File{Size: 10},
	}
	uploadStatus.file.Size = 10
	uploadStatus.Sync = false
	assert.Equal(false, skipOrIgnore(uploadStatus))

	// when sizes don't match
	uploadStatus.Sync = true
	mockFs["test"] = &fstest.MapFile{
		Sys: files_sdk.File{Size: 9},
	}
	uploadStatus.file.Size = 10
	assert.Equal(false, skipOrIgnore(uploadStatus))
	assert.Equal(nil, progressReportError)

	// when sizes do match
	mockFs["test"] = &fstest.MapFile{
		Sys: files_sdk.File{Size: 10},
	}
	uploadStatus.file.Size = 10
	assert.Equal(true, skipOrIgnore(uploadStatus))

	// There is no server version
	delete(mockFs, "test")
	assert.Equal(false, skipOrIgnore(uploadStatus))

	// when sizes do match on a deeply nested path
	oldUploadStatus := *uploadStatus
	uploadStatus = &UploadStatus{job: job, Mutex: &sync.RWMutex{}, file: files_sdk.File{Path: "test/path/test"}, remotePath: "test/path/test"}
	uploadStatus.Sync = true
	mockFs["test/path/test"] = &fstest.MapFile{
		Sys: files_sdk.File{Size: 10},
	}
	mockFs["test/path"] = &fstest.MapFile{
		Sys:  files_sdk.File{Size: 10},
		Mode: fs.ModeDir,
	}
	uploadStatus.file.Size = 10
	assert.Equal(true, skipOrIgnore(uploadStatus))
	uploadStatus = &oldUploadStatus

	// when transfer is a single file
	uploadStatus.job.Type = directory.File
	uploadStatus.Sync = true
	mockFs["test"] = &fstest.MapFile{
		Sys: files_sdk.File{Size: 10},
	}
	uploadStatus.file.Size = 10
	assert.Equal(true, skipOrIgnore(uploadStatus))
	uploadStatus.job.Type = directory.Dir
	uploadStatus.Sync = false

	// Ignore files
	job.GitIgnore, _ = ignore.New([]string{"*.css"}...)
	uploadStatus.localPath = "main.css"
	assert.Equal(true, skipOrIgnore(uploadStatus))

	uploadStatus.localPath = "main.php"
	assert.Equal(false, skipOrIgnore(uploadStatus))

	job.GitIgnore, _ = ignore.New([]string{"*.css", "*.php"}...)
	uploadStatus.localPath = "main.css"
	assert.Equal(true, skipOrIgnore(uploadStatus))
}
