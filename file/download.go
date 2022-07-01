package file

import (
	"context"
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
			RetryPolicy:    RetryPolicy(newJob.RetryPolicy),
			EventsReporter: newJob.EventsReporter,
		})
}

func (c *Client) DownloadToFile(ctx context.Context, params files_sdk.FileDownloadParams, filePath string) (files_sdk.File, error) {
	out, err := os.Create(filePath)
	if err != nil {
		return files_sdk.File{}, err
	}
	params.Writer = out
	return c.Download(ctx, params)
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
}

func (c *Client) Downloader(ctx context.Context, params DownloaderParams) *status.Job {
	return downloader(ctx, FS{}.Init(c.Config), params)
}

type Entity struct {
	fs.File
	error
}
