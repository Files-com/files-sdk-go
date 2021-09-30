package file

import (
	"context"
	"testing"
	"time"

	"github.com/Files-com/files-sdk-go/v2/ignore"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/file/status"

	"github.com/stretchr/testify/assert"
)

type MockUploader struct {
	findError files_sdk.ResponseError
}

func (m *MockUploader) UploadIO(context.Context, UploadIOParams) (files_sdk.File, files_sdk.FileUploadPart, Parts, error) {
	return files_sdk.File{}, files_sdk.FileUploadPart{}, Parts{}, nil
}

func (m *MockUploader) Find(context.Context, files_sdk.FileFindParams) (files_sdk.File, error) {
	return files_sdk.File{}, m.findError
}

func Test_skipOrIgnore(t *testing.T) {
	assert := assert.New(t)
	job := status.Job{}.Init()
	job.GitIgnore, _ = ignore.New()
	job.Params = UploadParams{}
	uploadStatus := &UploadStatus{Job: job}
	uploadStatus.Job.Add(uploadStatus)
	uploader := &MockUploader{}
	uploadStatus.Uploader = uploader
	ctx := context.Background()
	var progressReportStatus status.File
	var progressReportError error
	uploadStatus.Job.EventsReporter = Reporter(func(s status.File) {
		progressReportStatus = s
		progressReportError = s.Err
	})

	// sync not enabled
	uploadStatus.Sync = false
	assert.Equal(false, skipOrIgnore(ctx, uploadStatus))

	// Mtime is the same between server and local
	uploadStatus.Sync = true
	assert.Equal(true, skipOrIgnore(ctx, uploadStatus))
	assert.Equal(status.Skipped, uploadStatus.Status)
	assert.Equal(uploadStatus.Status, progressReportStatus.Status)
	assert.Equal(nil, progressReportError)

	// local version is newer than server
	uploadStatus.Mtime = time.Now()
	assert.Equal(false, skipOrIgnore(ctx, uploadStatus))

	// There is no server version
	uploader.findError = files_sdk.ResponseError{Type: "not-found"}
	assert.Equal(false, skipOrIgnore(ctx, uploadStatus))

	// Ignore files
	job.GitIgnore, _ = ignore.New([]string{"*.css"}...)
	uploadStatus.LocalPath = "main.css"
	assert.Equal(true, skipOrIgnore(ctx, uploadStatus))

	uploadStatus.LocalPath = "main.php"
	assert.Equal(false, skipOrIgnore(ctx, uploadStatus))

	job.GitIgnore, _ = ignore.New([]string{"*.css", "*.php"}...)
	uploadStatus.LocalPath = "main.css"
	assert.Equal(true, skipOrIgnore(ctx, uploadStatus))
}
