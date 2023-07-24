package file

import (
	"io"
	"io/fs"
	"os"

	"github.com/Files-com/files-sdk-go/v2/file/manager"
	"github.com/Files-com/files-sdk-go/v2/file/status"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
)

func (c *Client) DownloadRetry(job status.Job, opts ...files_sdk.RequestResponseOption) *status.Job {
	newJob := job.ClearStatuses()
	return c.Downloader(
		DownloaderParams{
			RemotePath:     newJob.RemotePath,
			Sync:           newJob.Sync,
			Manager:        newJob.Manager,
			LocalPath:      newJob.LocalPath,
			RetryPolicy:    newJob.RetryPolicy.(RetryPolicy),
			EventsReporter: newJob.EventsReporter,
		},
		opts...)
}

func (c *Client) DownloadToFile(params files_sdk.FileDownloadParams, filePath string, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	out, err := os.Create(filePath)
	if err != nil {
		return files_sdk.File{}, err
	}
	return c.Download(params, append(opts, files_sdk.ResponseBodyOption(func(closer io.ReadCloser) error {
		_, err := io.Copy(out, closer)
		return err
	}))...)
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
	DryRun bool
}

func (c *Client) Downloader(params DownloaderParams, opts ...files_sdk.RequestResponseOption) *status.Job {
	params.Config = c.Config
	return downloader(files_sdk.ContextOption(opts), (&FS{}).Init(c.Config, true), params)
}

type Entity struct {
	fs.File
	fs.FS
	error
}
