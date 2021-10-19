package file

import (
	"context"
	"sync"
	"testing"

	"github.com/Files-com/files-sdk-go/v2/ignore"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/file/status"

	"github.com/stretchr/testify/assert"
)

type MockUploader struct {
	files_sdk.File
	findError files_sdk.ResponseError
}

func (m *MockUploader) UploadIO(context.Context, UploadIOParams) (files_sdk.File, files_sdk.FileUploadPart, Parts, error) {
	return files_sdk.File{}, files_sdk.FileUploadPart{}, Parts{}, nil
}

func (m *MockUploader) Find(context.Context, files_sdk.FileFindParams) (files_sdk.File, error) {
	return m.File, m.findError
}

func Test_skipOrIgnore(t *testing.T) {
	assert := assert.New(t)
	job := status.Job{}.Init()
	job.GitIgnore, _ = ignore.New()
	job.Params = UploaderParams{}
	uploadStatus := &UploadStatus{job: job, Mutex: &sync.RWMutex{}}
	uploadStatus.Job().Add(uploadStatus)
	uploader := &MockUploader{}
	uploadStatus.Uploader = uploader
	ctx := context.Background()
	var progressReportStatus status.File
	var progressReportError error
	uploadStatus.job.EventsReporter = Reporter(func(s status.File) {
		progressReportStatus = s
		progressReportError = s.Err
	})

	// sync not enabled and sizes do match
	uploader.File.Size = 10
	uploadStatus.file.Size = 10
	uploadStatus.Sync = false
	assert.Equal(false, skipOrIgnore(ctx, uploadStatus))

	// Mtime is the same between server and local
	uploadStatus.Sync = true
	assert.Equal(true, skipOrIgnore(ctx, uploadStatus))
	assert.Equal(status.Skipped, uploadStatus.Status())
	assert.Equal(uploadStatus.Status(), progressReportStatus.Status)
	assert.Equal(nil, progressReportError)

	// when sizes don't match
	uploader.File.Size = 9
	uploadStatus.file.Size = 10
	assert.Equal(false, skipOrIgnore(ctx, uploadStatus))
	assert.Equal(status.Skipped, uploadStatus.Status())
	assert.Equal(uploadStatus.Status(), progressReportStatus.Status)
	assert.Equal(nil, progressReportError)

	// when sizes do match
	uploader.File.Size = 10
	uploadStatus.file.Size = 10
	assert.Equal(true, skipOrIgnore(ctx, uploadStatus))

	// There is no server version
	uploader.findError = files_sdk.ResponseError{Type: "not-found"}
	assert.Equal(false, skipOrIgnore(ctx, uploadStatus))

	// Ignore files
	job.GitIgnore, _ = ignore.New([]string{"*.css"}...)
	uploadStatus.localPath = "main.css"
	assert.Equal(true, skipOrIgnore(ctx, uploadStatus))

	uploadStatus.localPath = "main.php"
	assert.Equal(false, skipOrIgnore(ctx, uploadStatus))

	job.GitIgnore, _ = ignore.New([]string{"*.css", "*.php"}...)
	uploadStatus.localPath = "main.css"
	assert.Equal(true, skipOrIgnore(ctx, uploadStatus))
}
