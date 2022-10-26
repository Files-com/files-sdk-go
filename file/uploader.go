package file

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Files-com/files-sdk-go/v2/lib"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/directory"
	"github.com/Files-com/files-sdk-go/v2/file/manager"
	"github.com/Files-com/files-sdk-go/v2/file/status"
	"github.com/Files-com/files-sdk-go/v2/ignore"
	"github.com/Files-com/files-sdk-go/v2/lib/direction"
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

		//When the local/dest has a trailing slash
		if !(lib.Path{Path: params.LocalPath}).EndingSlash() {
			_, lastDir := filepath.Split(params.LocalPath)
			params.RemotePath = filepath.Join(params.RemotePath, lastDir)
		}

	} else {
		job.Type = directory.File
	}
	job.Client = c
	onComplete := make(chan *UploadStatus)
	job.CodeStart = func() {
		job.Scan()

		go enqueueIndexedUploads(job, jobCtx, onComplete)
		status.WaitTellFinished(job, onComplete, func() { RetryByPolicy(jobCtx, job, RetryPolicy(job.RetryPolicy), false) })

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
		var err error
		job.GitIgnore, err = ignore.New(params.Ignore...)
		if err != nil {
			job.Add(metaFile)
			job.UpdateStatus(status.Errored, metaFile, err)
			job.EndScan()
			job.Finish()
			return
		}
		if job.Type == directory.File {
			metaFile.file = files_sdk.File{
				DisplayName: filepath.Base(params.LocalPath),
				Type:        job.Direction.Name(),
				Mtime:       lib.Time(fi.ModTime()),
				Size:        fi.Size(),
				Path:        params.RemotePath,
			}
			metaFile.remotePath, err = remotePath(jobCtx, params.LocalPath, params.RemotePath, c)
			if err != nil {
				job.Add(metaFile)
				job.UpdateStatus(status.Errored, metaFile, err)
				job.EndScan()
				job.Finish()
				return
			}

			job.Add(metaFile)
			job.UpdateStatus(status.Indexed, metaFile, nil)
		} else {
			var walkErr error
			walkErr = walkPaginated(
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

func enqueueIndexedUploads(job *status.Job, jobCtx context.Context, onComplete chan *UploadStatus) {
	for !job.EndScanning.Called() || job.Count(status.Indexed) > 0 {
		f, ok := job.Find(status.Indexed)
		if ok {
			enqueueUpload(jobCtx, job, f.(*UploadStatus), onComplete)
		}
	}
}

func enqueueUpload(ctx context.Context, job *status.Job, uploadStatus *UploadStatus, onComplete chan *UploadStatus) {
	if uploadStatus.error != nil || uploadStatus.missingStat {
		job.UpdateStatus(status.Errored, uploadStatus, uploadStatus.RecentError())
		onComplete <- uploadStatus
		return
	}
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
		localFile, err = os.Open(uploadStatus.LocalPath())
		if err != nil {
			uploadStatus.Job().UpdateStatus(status.Errored, uploadStatus, err)
			return
		}
		if ctx.Err() != nil {
			job.UpdateStatus(status.Canceled, uploadStatus, nil)
			return
		}

		stats, err := os.Stat(uploadStatus.LocalPath())
		if err != nil {
			uploadStatus.Job().UpdateStatus(status.Errored, uploadStatus, err)
			return
		}
		_, uploadStatus.FileUploadPart, uploadStatus.Parts, err = uploadStatus.UploadIO(
			ctx, UploadIOParams{
				Path:           uploadStatus.RemotePath(),
				Reader:         localFile,
				Size:           uploadStatus.File().Size,
				Progress:       uploadProgress(uploadStatus),
				Manager:        job.FilePartsManager,
				Parts:          uploadStatus.Parts,
				FileUploadPart: uploadStatus.FileUploadPart,
				ProvidedMtime:  stats.ModTime(),
			})
		if err != nil {
			uploadStatus.Job().StatusFromError(uploadStatus, err)
		} else {
			uploadStatus.SetUploadedBytes(uploadStatus.Parts.SuccessfulBytes())
			uploadStatus.Job().UpdateStatus(status.Complete, uploadStatus, nil)
		}
	}()
}

func walkPaginated(ctx context.Context, localFolderPath string, destinationRootPath string, job *status.Job, params UploaderParams, c Uploader) error {
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

			uploadStatus := UploadStatus{
				Uploader:   c,
				job:        job,
				remotePath: destination,
				localPath:  path,
				Sync:       params.Sync,
				status:     status.Indexed,
				Mutex:      &sync.RWMutex{},
				error:      err,
			}

			if err != nil {
				uploadStatus.file = files_sdk.File{Type: "file", DisplayName: filepath.Base(destination), Path: destination}
				uploadStatus.missingStat = true
				uploadStatus.error = err
			} else {
				uploadStatus.file = files_sdk.File{Type: "file", DisplayName: filepath.Base(destination), Path: destination, Size: info.Size(), Mtime: lib.Time(info.ModTime())}
			}
			job.Add(&uploadStatus)
			return nil
		},
		Unsorted: true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
	})
	return err
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
