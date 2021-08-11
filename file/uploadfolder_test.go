package file

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/zenthangplus/goccm"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/file/status"

	"github.com/stretchr/testify/assert"
)

type MockUploader struct {
	findError files_sdk.ResponseError
}

func (m *MockUploader) Upload(context.Context, io.ReaderAt, int64, files_sdk.FileBeginUploadParams, func(int64), goccm.ConcurrencyManager) (files_sdk.File, error) {
	return files_sdk.File{}, nil
}

func (m *MockUploader) Find(context.Context, string) (files_sdk.File, error) {
	return files_sdk.File{}, m.findError
}

func Test_checkUploadSync(t *testing.T) {
	assert := assert.New(t)
	uploadStatus := UploadStatus{Job: &status.Job{}}
	params := UploadParams{}
	uploader := &MockUploader{}
	var ctx context.Context
	ctx, uploadStatus.CancelFunc = context.WithCancel(context.Background())
	var progressReportStatus status.File
	var progressReportError error
	params.Reporter = func(s status.File, err error) {
		progressReportStatus = s
		progressReportError = err
	}

	// sync not enabled
	params.Sync = false
	assert.Equal(true, checkUpdateSync(ctx, &uploadStatus, &params, uploader))

	// Mtime is the same between server and local
	params.Sync = true
	assert.Equal(false, checkUpdateSync(ctx, &uploadStatus, &params, uploader))
	assert.Equal(uploadStatus.Status, status.Skipped)
	assert.Equal(uploadStatus.Status, progressReportStatus.Status)
	assert.Equal(nil, progressReportError)

	// local version is newer than server
	uploadStatus.Mtime = time.Now()
	assert.Equal(true, checkUpdateSync(ctx, &uploadStatus, &params, uploader))

	// There is no server version
	uploader.findError = files_sdk.ResponseError{Type: "not-found"}
	assert.Equal(true, checkUpdateSync(ctx, &uploadStatus, &params, uploader))
}
