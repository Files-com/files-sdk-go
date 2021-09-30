package file

import (
	"context"
	"os"
	"os/user"
	"path/filepath"

	"github.com/Files-com/files-sdk-go/v2/ignore"

	"github.com/Files-com/files-sdk-go/v2/directory"
	"github.com/Files-com/files-sdk-go/v2/lib/direction"

	"github.com/Files-com/files-sdk-go/v2/file/manager"
	"github.com/Files-com/files-sdk-go/v2/file/status"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/lib"
)

func (c *Client) UploadRetry(ctx context.Context, job status.Job) *status.Job {
	newJob := job.ClearStatuses()
	return c.Uploader(ctx,
		UploadParams{
			Sync:           newJob.Sync,
			LocalPath:      newJob.LocalPath,
			RemotePath:     newJob.RemotePath,
			EventsReporter: newJob.EventsReporter,
			Manager:        newJob.Manager,
			RetryPolicy:    RetryPolicy(newJob.RetryPolicy),
			Ignore:         newJob.Params.(UploadParams).Ignore,
		},
	)
}

type UploadParams struct {
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

func (c *Client) Uploader(ctx context.Context, params UploadParams) *status.Job {
	job := status.Job{}.Init()
	SetJobParams(job, direction.UploadType, params)
	job.CodeStart = func() {
		params.Job = job
		job.Params = params
		file := &status.File{File: files_sdk.File{}, RemotePath: params.RemotePath, LocalPath: params.LocalPath, Status: status.Queued, Job: job}
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
		fi, err := os.Stat(params.LocalPath)
		if err != nil {
			job.Add(file)
			job.UpdateStatus(status.Errored, file, err)
			return
		}
		if fi.IsDir() {
			c.UploadFolder(ctx, params).CodeStart()
		} else {
			c.UploadFile(ctx, params).CodeStart()
		}
	}

	return job
}

func (c *Client) UploadFolder(ctx context.Context, params UploadParams) *status.Job {
	return uploadFolder(ctx, c, params)
}

func (c *Client) UploadFile(parentCtx context.Context, params UploadParams) *status.Job {
	var job *status.Job
	if params.Job == nil {
		job = status.Job{}.Init()
	} else {
		job = params.Job
	}
	SetJobParams(job, direction.UploadType, params)
	job.Client = c
	jobCtx := job.WithContext(parentCtx)
	job.Type = directory.File

	job.CodeStart = func() {
		var localFile *os.File
		var err error
		defer func() {
			if localFile != nil {
				localFile.Close()
			}
			job.Finish()
			job.FilesManager.Done()
		}()
		uploadStatus := &UploadStatus{
			Job:        job,
			LocalPath:  params.LocalPath,
			RemotePath: params.RemotePath,
			Sync:       params.Sync,
		}

		beginUpload := files_sdk.FileBeginUploadParams{MkdirParents: lib.Bool(true)}
		destination := params.RemotePath
		_, localFileName := filepath.Split(params.LocalPath)
		if params.RemotePath == "" {
			destination = localFileName
		} else {
			remoteFile, err := c.Find(jobCtx, files_sdk.FileFindParams{Path: params.RemotePath})
			responseError, ok := err.(files_sdk.ResponseError)
			if remoteFile.Type == "directory" {
				destination = filepath.Join(params.RemotePath, localFileName)
			} else if ok && responseError.Type == "not-found" {
				if destination[len(destination)-1:] == "/" {
					destination = filepath.Join(params.RemotePath, localFileName)
				}
			} else if err != nil {
				job.UpdateStatus(status.Errored, uploadStatus, err)
			}
		}
		job.FilesManager.Wait()
		fi, err := os.Stat(params.LocalPath)
		if err != nil {
			job.Add(uploadStatus)
			job.UpdateStatus(status.Errored, uploadStatus, err)
			return
		}
		uploadStatus.RemotePath = destination
		uploadStatus.File = files_sdk.File{
			DisplayName: filepath.Base(destination),
			Path:        destination,
			Type:        "file",
			Mtime:       fi.ModTime(),
			Size:        fi.Size(),
		}
		uploadStatus.Uploader = c
		job.Add(uploadStatus)
		job.UpdateStatus(status.Queued, uploadStatus, nil)

		job.GitIgnore, err = ignore.New(params.Ignore...)
		if err != nil {
			job.Add(uploadStatus)
			job.UpdateStatus(status.Errored, uploadStatus, err)
			return
		}
		if skipOrIgnore(jobCtx, uploadStatus) {
			return
		}
		localFile, err = os.Open(params.LocalPath)
		if err != nil {
			uploadStatus.Job.UpdateStatus(status.Errored, uploadStatus, err)
			return
		}
		beginUpload.Path = destination
		uParams := UploadIOParams{
			Path:           uploadStatus.RemotePath,
			Reader:         localFile,
			Size:           uploadStatus.ToStatusFile().Size,
			Progress:       uploadProgress(uploadStatus),
			Manager:        job.FilePartsManager,
			Parts:          uploadStatus.Parts,
			FileUploadPart: uploadStatus.FileUploadPart,
		}
		var file files_sdk.File
		file, uploadStatus.FileUploadPart, uploadStatus.Parts, err = uploadStatus.UploadIO(jobCtx, uParams)
		if dealWithCanceledError(uploadStatus, err) {
			uploadStatus.File = file
		}
		RetryByPolicy(jobCtx, job, RetryPolicy(job.RetryPolicy), false)
	}
	job.Wait = func() {
		for !job.Finished.Called {
		}
	}

	return job
}
