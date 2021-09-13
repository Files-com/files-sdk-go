package file

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Files-com/files-sdk-go/v2/directory"
	"github.com/Files-com/files-sdk-go/v2/file/manager"
	"github.com/Files-com/files-sdk-go/v2/ignore"
	"github.com/Files-com/files-sdk-go/v2/lib/direction"
	"github.com/bradfitz/iter"
	"github.com/zenthangplus/goccm"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/file/status"
	"github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/karrick/godirwalk"
)

type Uploader interface {
	Upload(context.Context, io.ReaderAt, int64, files_sdk.FileBeginUploadParams, func(int64), goccm.ConcurrencyManager) (files_sdk.File, error)
	Find(context.Context, string) (files_sdk.File, error)
}

func uploadFolder(ctx context.Context, c Uploader, params UploadParams) *status.Job {
	var job *status.Job
	if params.Job == nil {
		job = status.Job{}.Init()
	} else {
		job = params.Job
	}
	SetJobParams(job, direction.UploadType, params)
	jobCtx := job.WithContext(ctx)
	job.Client = c
	job.Type = directory.Dir
	metaFile := &UploadStatus{
		File:      files_sdk.File{DisplayName: filepath.Base(params.LocalPath), Type: "directory"},
		Job:       job,
		Status:    status.Errored,
		LocalPath: params.LocalPath,
		Sync:      params.Sync,
	}

	onComplete := make(chan *UploadStatus)
	count := 0
	job.Start = func() {
		job.Scanning = true
		go enqueueIndexedUploads(job, jobCtx, onComplete)
		job.StartTime = time.Now()
		i, err := ignore.New(params.Ignore...)
		if err != nil {
			job.Add(metaFile)
			job.UpdateStatus(status.Errored, metaFile, err)
			job.EndTime = time.Now()
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
		job.Scanning = false
		if walkErr != nil {
			job.Add(metaFile)
			job.UpdateStatus(status.Errored, metaFile, walkErr)
		}

		go markUploadOnComplete(count, job, metaFile, onComplete, jobCtx)
	}

	job.Wait = func() {
		for job.EndTime.IsZero() {
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
		select {
		case <-jobCtx.Done():
			break
		case <-onComplete:
		}
	}
	close(onComplete)
	RetryTransfers(jobCtx, job)
	job.EndTime = time.Now()
}

func enqueueIndexedUploads(job *status.Job, jobCtx context.Context, onComplete chan *UploadStatus) {
	for job.Scanning || job.Count(status.Indexed) > 0 {
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
		defer func() {
			job.FilesManager.Done()
			onComplete <- uploadStatus
		}()
		if skipOrIgnore(ctx, uploadStatus) {
			return
		}
		localFile, err := os.Open(uploadStatus.LocalPath)
		defer localFile.Close()
		if dealWithDBasicError(uploadStatus, err) {
			return
		}
		if ctx.Err() != nil {
			job.UpdateStatus(status.Canceled, uploadStatus, nil)
			return
		}
		file, err := uploadStatus.Upload(ctx, localFile, uploadStatus.Size, files_sdk.FileBeginUploadParams{Path: uploadStatus.RemotePath, MkdirParents: lib.Bool(true)}, uploadProgress(uploadStatus), job.FilePartsManager)
		dealWithCanceledError(uploadStatus, err, file)
	}()
}

func walkPaginated(ctx context.Context, localFolderPath string, destinationRootPath string, job *status.Job, params UploadParams, c Uploader) (int, error) {
	count := 0
	job.Scanning = true
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
	job.Scanning = false
	return count, err
}

func skipOrIgnore(downloadCtx context.Context, uploadStatus *UploadStatus) bool {
	if uploadStatus.Job.MatchesPath(uploadStatus.LocalPath) {
		uploadStatus.Job.UpdateStatus(status.Ignored, uploadStatus, nil)
		return true
	}

	if uploadStatus.Sync {
		file, err := uploadStatus.Find(downloadCtx, uploadStatus.RemotePath)
		responseError, ok := err.(files_sdk.ResponseError)
		if ok && responseError.Type == "not-found" {
			return false
		}
		// local is not after server
		if !uploadStatus.Mtime.After(file.Mtime) {
			// Server version is the same or newer
			uploadStatus.Job.UpdateStatus(status.Skipped, uploadStatus, nil)
			return true
		}
	}
	return false
}
