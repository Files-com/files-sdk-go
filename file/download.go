package file

import (
	"context"
	"io"
	"io/fs"
	"os"

	"github.com/Files-com/files-sdk-go/v2/file/manager"
	"github.com/Files-com/files-sdk-go/v2/file/status"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
)

func (c *Client) DownloadRetry(ctx context.Context, job status.Job) *status.Job {
	newJob := job.ClearStatuses()
	return c.Downloader(ctx,
		DownloaderParams{
			RemotePath:     newJob.RemotePath,
			Sync:           newJob.Sync,
			Manager:        newJob.Manager,
			LocalPath:      newJob.LocalPath,
			RetryPolicy:    newJob.RetryPolicy.(RetryPolicy),
			EventsReporter: newJob.EventsReporter,
		})
}

func (c *Client) DownloadToFile(ctx context.Context, params files_sdk.FileDownloadParams, filePath string) (files_sdk.File, error) {
	out, err := os.Create(filePath)
	if err != nil {
		return files_sdk.File{}, err
	}
	return c.Download(ctx, params, files_sdk.ResponseBodyOption(func(closer io.ReadCloser) error {
		_, err := io.Copy(out, closer)
		return err
	}))
}

type DownloaderParams struct {
	RemotePath    string
	RemoteFile    files_sdk.File
	LocalPath     string
	Sync          bool
	PreserveTimes bool
	RetryPolicy
	*manager.Manager
	status.EventsReporter
	files_sdk.Config
}

func (c *Client) Downloader(ctx context.Context, params DownloaderParams) *status.Job {
	return downloader(ctx, (&FS{}).Init(c.Config, true), params)
}

type Entity struct {
	fs.File
	fs.FS
	error
}
