package file

import (
	"context"
	"os/user"
	"path/filepath"

	"github.com/Files-com/files-sdk-go/v2/lib/direction"

	"github.com/Files-com/files-sdk-go/v2/file/manager"
	"github.com/Files-com/files-sdk-go/v2/file/status"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
)

func (c *Client) UploadRetry(ctx context.Context, job status.Job) *status.Job {
	newJob := job.ClearStatuses()
	return c.Uploader(ctx,
		UploaderParams{
			Sync:           newJob.Sync,
			LocalPath:      newJob.LocalPath,
			RemotePath:     newJob.RemotePath,
			EventsReporter: newJob.EventsReporter,
			Manager:        newJob.Manager,
			RetryPolicy:    RetryPolicy(newJob.RetryPolicy),
			Ignore:         newJob.Params.(UploaderParams).Ignore,
		},
	)
}

type UploaderParams struct {
	Ignore []string
	*status.Job
	Sync       bool
	LocalPath  string
	RemotePath string
	RetryPolicy
	status.EventsReporter
	*manager.Manager
}

func expand(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, path[1:]), nil
}

func (c *Client) Uploader(ctx context.Context, params UploaderParams) *status.Job {
	job := status.Job{}.Init()
	SetJobParams(job, direction.UploadType, params)
	job.CodeStart = func() {
		params.Job = job
		job.Params = params
		file := &UploadStatus{file: files_sdk.File{}, remotePath: params.RemotePath, localPath: params.LocalPath, status: status.Queued, job: job}
		expandedPath, err := expand(params.LocalPath)
		if err != nil {
			job.Add(file)
			job.UpdateStatus(status.Errored, file, err)
			return
		}
		absolutePath, err := filepath.Abs(expandedPath)
		if err != nil {
			job.Add(file)
			job.UpdateStatus(status.Errored, file, err)
			return
		}
		params.LocalPath = absolutePath

		uploader(ctx, c, params).CodeStart()
	}

	return job
}
