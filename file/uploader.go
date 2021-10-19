package file

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"

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

func uploader(parentCtx context.Context, c Uploader, params UploaderParams) *status.Job {
	var job *status.Job
	if params.Job == nil {
		job = status.Job{}.Init()
	} else {
		job = params.Job
	}
	SetJobParams(job, direction.UploadType, params)
	jobCtx := job.WithContext(parentCtx)

	fi, statErr := os.Stat(params.LocalPath)

	if statErr == nil && fi.IsDir() {
		job.Type = directory.Dir
	} else {
		job.Type = directory.File
	}
	job.Client = c
	onComplete := make(chan *UploadStatus)
	count := 0
	job.CodeStart = func() {
		job.Scan()

		go enqueueIndexedUploads(job, jobCtx, onComplete)
		metaFile := &UploadStatus{
			job:       job,
			status:    status.Errored,
			localPath: params.LocalPath,
			Sync:      params.Sync,
			Uploader:  c,
			Mutex:     &sync.RWMutex{},
		}
		if statErr != nil {
			job.Add(metaFile)
			job.UpdateStatus(status.Errored, metaFile, statErr)
			job.EndScan()
			job.Finish()
			return
		}
		i, err := ignore.New(params.Ignore...)
		if err != nil {
			job.Add(metaFile)
			job.UpdateStatus(status.Errored, metaFile, err)
			job.EndScan()
			job.Finish()
			return
		}
		job.GitIgnore = i
		if job.Type == directory.File {
			metaFile.file = files_sdk.File{
				DisplayName: filepath.Base(params.LocalPath),
				Type:        job.Direction.Name(),
				Mtime:       fi.ModTime(),
				Size:        fi.Size(),
				Path:        params.RemotePath,
			}
			path, err := remotePath(jobCtx, params.LocalPath, params.RemotePath, c)
			metaFile.remotePath = path
			if err != nil {
				job.Add(metaFile)
				job.UpdateStatus(status.Errored, metaFile, err)
				job.EndScan()
				job.Finish()
				return
			}

			job.Add(metaFile)
			job.UpdateStatus(status.Indexed, metaFile, nil)
			count = 1
		} else {
			var walkErr error
			count, walkErr = walkPaginated(
				jobCtx,
				params.LocalPath,
				params.RemotePath,
				job,
				params,
				c,
			)
			if walkErr != nil {
				job.Add(metaFile)
				job.UpdateStatus(status.Errored, metaFile, walkErr)
			}
		}
		job.EndScan()

		go markUploadOnComplete(count, job, onComplete, jobCtx)
	}

	return job
}

func remotePath(ctx context.Context, localPath string, remotePath string, c Uploader) (string, error) {
	destination := remotePath
	_, localFileName := filepath.Split(localPath)
	if remotePath == "" {
		destination = localFileName
	} else {
		remoteFile, err := c.Find(ctx, files_sdk.FileFindParams{Path: remotePath})
		responseError, ok := err.(files_sdk.ResponseError)
		if remoteFile.Type == "directory" {
			destination = filepath.Join(remotePath, localFileName)
		} else if ok && responseError.Type == "not-found" {
			if destination[len(destination)-1:] == "/" {
				destination = filepath.Join(remotePath, localFileName)
			}
		} else if err != nil {
			return "", err
		}
	}
	return destination, nil
}

func markUploadOnComplete(count int, job *status.Job, onComplete chan *UploadStatus, jobCtx context.Context) {
	for range iter.N(count) {
		<-onComplete
	}
	close(onComplete)
	RetryByPolicy(jobCtx, job, RetryPolicy(job.RetryPolicy), false)
	job.Finish()
}

func enqueueIndexedUploads(job *status.Job, jobCtx context.Context, onComplete chan *UploadStatus) {
	for !job.EndScanning.Called() || job.Count(status.Indexed) > 0 {
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
		finish := func() {
			if localFile != nil {
				localFile.Close()
			}
			job.FilesManager.Done()
			onComplete <- uploadStatus
		}
		if skipOrIgnore(ctx, uploadStatus) {
			finish()
			return
		}
		localFile, err = os.Open(uploadStatus.LocalPath())
		if err != nil {
			uploadStatus.Job().UpdateStatus(status.Errored, uploadStatus, err)
			finish()
			return
		}
		if ctx.Err() != nil {
			job.UpdateStatus(status.Canceled, uploadStatus, nil)
			finish()
			return
		}
		params := UploadIOParams{
			Path:           uploadStatus.RemotePath(),
			Reader:         localFile,
			Size:           uploadStatus.File().Size,
			Progress:       uploadProgress(uploadStatus),
			Manager:        job.FilePartsManager,
			Parts:          uploadStatus.Parts,
			FileUploadPart: uploadStatus.FileUploadPart,
		}
		_, uploadStatus.FileUploadPart, uploadStatus.Parts, err = uploadStatus.UploadIO(ctx, params)
		localFile.Close()
		if err != nil {
			uploadStatus.Job().StatusFromError(uploadStatus, err)
		} else {
			uploadStatus.SetUploadedBytes(uploadStatus.Parts.SuccessfulBytes())
			uploadStatus.Job().UpdateStatus(status.Complete, uploadStatus, nil)
		}

		finish()
	}()
}

func walkPaginated(ctx context.Context, localFolderPath string, destinationRootPath string, job *status.Job, params UploaderParams, c Uploader) (int, error) {
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

			if de.IsDir() || !de.IsRegular() {
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
				job:        job,
				remotePath: destination,
				localPath:  path,
				Sync:       params.Sync,
				status:     status.Indexed,
				file:       files_sdk.File{Type: "file", DisplayName: filepath.Base(destination), Path: destination, Size: info.Size(), Mtime: info.ModTime()},
				Mutex:      &sync.RWMutex{},
			})
			return nil
		},
		Unsorted: true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
	})
	return count, err
}

func skipOrIgnore(downloadCtx context.Context, uploadStatus *UploadStatus) bool {
	if uploadStatus.Job().MatchesPath(uploadStatus.LocalPath()) {
		uploadStatus.Job().UpdateStatus(status.Ignored, uploadStatus, nil)
		return true
	}

	if uploadStatus.Sync {
		file, err := uploadStatus.Find(downloadCtx, files_sdk.FileFindParams{Path: uploadStatus.RemotePath()})
		responseError, ok := err.(files_sdk.ResponseError)
		if ok && responseError.Type == "not-found" {
			return false
		}
		// local is not after server
		if uploadStatus.File().Size == file.Size {
			// Server version is the same or newer
			uploadStatus.Job().UpdateStatus(status.Skipped, uploadStatus, nil)
			return true
		}
	}
	return false
}
