package file

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/directory"
	"github.com/Files-com/files-sdk-go/v2/file/manager"
	"github.com/Files-com/files-sdk-go/v2/file/status"
	"github.com/Files-com/files-sdk-go/v2/ignore"
	"github.com/Files-com/files-sdk-go/v2/lib/direction"
	"github.com/bradfitz/iter"
	"github.com/karrick/godirwalk"
)

type Uploader interface {
	UploadIO(context.Context, UploadIOParams) (files_sdk.File, files_sdk.FileUploadPart, Parts, error)
	Find(context.Context, files_sdk.FileFindParams) (files_sdk.File, error)
}

func uploadFolder(parentCtx context.Context, c Uploader, params UploadParams) *status.Job {
	var job *status.Job
	if params.Job == nil {
		job = status.Job{}.Init()
	} else {
		job = params.Job
	}
	SetJobParams(job, direction.UploadType, params)
	jobCtx := job.WithContext(parentCtx)
	job.Client = c
	job.Type = directory.Dir
	metaFile := &UploadStatus{
		Job:       job,
		Status:    status.Errored,
		LocalPath: params.LocalPath,
		Sync:      params.Sync,
	}
	metaFile.File = files_sdk.File{DisplayName: filepath.Base(params.LocalPath), Type: "directory"}

	onComplete := make(chan *UploadStatus)
	count := 0
	job.CodeStart = func() {
		job.Scan()
		go enqueueIndexedUploads(job, jobCtx, onComplete)
		i, err := ignore.New(params.Ignore...)
		if err != nil {
			job.Add(metaFile)
			job.UpdateStatus(status.Errored, metaFile, err)
			job.Timer.Stop()
			return
		}
		job.GitIgnore = i
		var walkErr error
		count, walkErr = walkPaginated(
			jobCtx,
			params.LocalPath,
			params.RemotePath,
			job,
			params,
			c,
		)
		job.EndScan()
		if walkErr != nil {
			job.Add(metaFile)
			job.UpdateStatus(status.Errored, metaFile, walkErr)
		}

		go markUploadOnComplete(count, job, metaFile, onComplete, jobCtx)
	}

	job.Wait = func() {
		for !job.Finished.Called {
		}
	}

	return job
}

func markUploadOnComplete(count int, job *status.Job, metaFile *UploadStatus, onComplete chan *UploadStatus, jobCtx context.Context) {
	if count == 0 {
		job.Add(metaFile)
		job.UpdateStatus(status.Complete, metaFile, nil)
	}
	for range iter.N(count) {
		<-onComplete
	}
	close(onComplete)
	RetryByPolicy(jobCtx, job, RetryPolicy(job.RetryPolicy), false)
	job.Finish()
}

func enqueueIndexedUploads(job *status.Job, jobCtx context.Context, onComplete chan *UploadStatus) {
	for !job.EndScanning.Called || job.Count(status.Indexed) > 0 {
		f, ok := job.Find(status.Indexed)
		if ok {
			enqueueUpload(jobCtx, job, f.(*UploadStatus), onComplete)
		}
	}
}

func enqueueUpload(ctx context.Context, job *status.Job, uploadStatus *UploadStatus, onComplete chan *UploadStatus) {
	job.UpdateStatus(status.Queued, uploadStatus, nil)
	if !manager.Wait(ctx, job.FilesManager) {
		job.UpdateStatus(status.Canceled, uploadStatus, nil)
		onComplete <- uploadStatus
		return
	}
	go func() {
		var localFile *os.File
		var err error
		defer func() {
			if localFile != nil {
				localFile.Close()
			}
			job.FilesManager.Done()
			onComplete <- uploadStatus
		}()
		if skipOrIgnore(ctx, uploadStatus) {
			return
		}
		localFile, err = os.Open(uploadStatus.LocalPath)
		if err != nil {
			uploadStatus.Job.UpdateStatus(status.Errored, uploadStatus, err)
			return
		}
		if ctx.Err() != nil {
			job.UpdateStatus(status.Canceled, uploadStatus, nil)
			return
		}
		params := UploadIOParams{
			Path:           uploadStatus.RemotePath,
			Reader:         localFile,
			Size:           uploadStatus.ToStatusFile().Size,
			Progress:       uploadProgress(uploadStatus),
			Manager:        job.FilePartsManager,
			Parts:          uploadStatus.Parts,
			FileUploadPart: uploadStatus.FileUploadPart,
		}
		var file files_sdk.File
		file, uploadStatus.FileUploadPart, uploadStatus.Parts, err = uploadStatus.UploadIO(ctx, params)
		if dealWithCanceledError(uploadStatus, err) {
			uploadStatus.File = file
		}
	}()
}

func walkPaginated(ctx context.Context, localFolderPath string, destinationRootPath string, job *status.Job, params UploadParams, c Uploader) (int, error) {
	count := 0
	err := godirwalk.Walk(localFolderPath, &godirwalk.Options{
		Callback: func(path string, de *godirwalk.Dirent) error {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			dir, filename := filepath.Split(path)

			if localFolderPath == dir && filename == "" {
				return nil
			}

			if de.IsDir() {
				return nil
			}

			var destination string
			var baseDestination string
			if localFolderPath != "." {
				baseDestination = strings.TrimPrefix(path, localFolderPath)
			} else if path != "." {
				baseDestination = path
			}
			baseDestination = strings.TrimLeft(baseDestination, "/")
			baseDestination = strings.TrimPrefix(baseDestination, "/")
			if destinationRootPath == "" {
				destination = baseDestination
			} else {
				destination = filepath.Join(destinationRootPath, baseDestination)
			}

			if destination == "." {
				destination = filename
			}

			info, err := os.Stat(filepath.Join(path))
			if err != nil {
				return err
			}

			count += 1
			job.Add(&UploadStatus{
				Uploader:   c,
				Job:        job,
				RemotePath: destination,
				LocalPath:  path,
				Sync:       params.Sync,
				Status:     status.Indexed,
				File:       files_sdk.File{Type: "file", DisplayName: filepath.Base(destination), Path: destination, Size: info.Size(), Mtime: info.ModTime()},
			})
			return nil
		},
		Unsorted: true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
	})
	return count, err
}

func skipOrIgnore(downloadCtx context.Context, uploadStatus *UploadStatus) bool {
	if uploadStatus.Job.MatchesPath(uploadStatus.LocalPath) {
		uploadStatus.Job.UpdateStatus(status.Ignored, uploadStatus, nil)
		return true
	}

	if uploadStatus.Sync {
		file, err := uploadStatus.Find(downloadCtx, files_sdk.FileFindParams{Path: uploadStatus.RemotePath})
		responseError, ok := err.(files_sdk.ResponseError)
		if ok && responseError.Type == "not-found" {
			return false
		}
		// local is not after server
		if !uploadStatus.ToStatusFile().Mtime.After(file.Mtime) {
			// Server version is the same or newer
			uploadStatus.Job.UpdateStatus(status.Skipped, uploadStatus, nil)
			return true
		}
	}
	return false
}
