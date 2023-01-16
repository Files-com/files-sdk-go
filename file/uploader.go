package file

import (
	"context"
	"errors"
	"io/fs"
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
)

type Uploader interface {
	UploadIO(context.Context, UploadIOParams) (files_sdk.File, files_sdk.FileUploadPart, Parts, []error, error)
	Find(context.Context, files_sdk.FileFindParams, ...files_sdk.RequestResponseOption) (files_sdk.File, error)
}

func uploader(parentCtx context.Context, c Uploader, params UploaderParams) *status.Job {
	var job *status.Job
	if params.Job == nil {
		job = status.Job{}.Init()
	} else {
		job = params.Job
	}
	SetJobParams(job, direction.UploadType, params, params.Config.Logger(), (&FS{}).Init(params.Config, true))
	job.Config = params.Config
	jobCtx := job.WithContext(parentCtx)

	fi, statErr := os.Stat(params.LocalPath)

	if statErr == nil && fi.IsDir() {
		job.Type = directory.Dir

		//When the local/dest has a trailing slash
		if !(lib.Path{Path: params.LocalPath}).EndingSlash() {
			_, lastDir := filepath.Split(params.LocalPath)
			params.RemotePath = lib.UrlJoinNoEscape(params.RemotePath, lastDir)
		}

	} else {
		job.Type = directory.File
	}
	job.Client = c
	onComplete := make(chan *UploadStatus)
	job.CodeStart = func() {
		job.Scan()

		go enqueueIndexedUploads(job, jobCtx, onComplete)
		status.WaitTellFinished(job, onComplete, func() { RetryByPolicy(jobCtx, job, job.RetryPolicy.(RetryPolicy), false) })

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
			destination = lib.UrlJoinNoEscape(remotePath, localFileName)
		} else if ok && responseError.Type == "not-found" {
			if destination[len(destination)-1:] == "/" {
				destination = lib.UrlJoinNoEscape(remotePath, localFileName)
			}
		} else if err != nil {
			return "", err
		}
	}
	return destination, nil
}

func enqueueIndexedUploads(job *status.Job, jobCtx context.Context, onComplete chan *UploadStatus) {
	for !job.EndScanning.Called() || job.Count(status.Indexed) > 0 {
		f, ok := job.EnqueueNext()
		if ok {
			go enqueueUpload(jobCtx, job, f.(*UploadStatus), onComplete)
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
	func() {
		var localFile *os.File
		var err error
		defer func() {
			if localFile != nil {
				localFile.Close()
			}
			job.FilesManager.Done()
			onComplete <- uploadStatus
		}()
		if skipOrIgnore(uploadStatus) {
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
		var otherErrors []error
		_, uploadStatus.FileUploadPart, uploadStatus.Parts, otherErrors, err = uploadStatus.UploadIO(
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
			if len(otherErrors) > 0 {
				for _, otherErr := range otherErrors {
					uploadStatus.Job().StatusFromError(uploadStatus, otherErr)
				}
			}
			uploadStatus.Job().StatusFromError(uploadStatus, err)
		} else {
			uploadStatus.SetUploadedBytes(uploadStatus.Parts.SuccessfulBytes())
			if localFile != nil {
				localFile.Close()
			}
			uploadStatus.Job().UpdateStatus(status.Complete, uploadStatus, nil)
		}
	}()
}

func walkPaginated(ctx context.Context, localFolderPath string, destinationRootPath string, job *status.Job, params UploaderParams, c Uploader) error {
	err := fs.WalkDir(os.DirFS(localFolderPath), ".", func(name string, d fs.DirEntry, err error) error {
		path := filepath.Join(localFolderPath, name)
		uploadStatus, ok := buildUploadStatus(path, localFolderPath, destinationRootPath, c, job, params)
		if err != nil {
			if errors.Is(err, fs.ErrPermission) {
				job.UpdateStatus(status.Errored, &uploadStatus, err)
				job.Add(&uploadStatus)
				return fs.SkipDir
			}
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if d.IsDir() || !d.Type().IsRegular() {
			return nil
		}

		if !ok {
			return nil
		}

		uploadStatus.file.Type = "file"
		info, err := d.Info()
		if err != nil {
			uploadStatus.missingStat = true
			uploadStatus.error = err
		} else {
			uploadStatus.file.Size = info.Size()
			uploadStatus.file.Mtime = lib.Time(info.ModTime())
		}
		job.Add(&uploadStatus)
		return nil
	})
	return err
}

func buildUploadStatus(path string, localFolderPath string, destinationRootPath string, c Uploader, job *status.Job, params UploaderParams) (UploadStatus, bool) {
	dir, filename := filepath.Split(path)

	if localFolderPath == dir && filename == "" {
		return UploadStatus{}, false
	}

	destination := buildDestination(path, localFolderPath, destinationRootPath, filename)
	uploadStatus := UploadStatus{
		Uploader:   c,
		job:        job,
		remotePath: destination,
		localPath:  path,
		Sync:       params.Sync,
		status:     status.Indexed,
		Mutex:      &sync.RWMutex{},
		file:       files_sdk.File{Path: destination, DisplayName: filepath.Base(destination)},
	}
	return uploadStatus, true
}

func buildDestination(path string, localFolderPath string, destinationRootPath string, filename string) string {
	var destination string
	var baseDestination string
	if localFolderPath != "." {
		baseDestination = strings.TrimPrefix(path, localFolderPath)
	} else if path != "." {
		baseDestination = path
	}
	baseDestination = strings.TrimLeft(baseDestination, string(os.PathSeparator))
	baseDestination = strings.TrimPrefix(baseDestination, string(os.PathSeparator))
	if destinationRootPath == "" {
		destination = baseDestination
	} else {
		destination = filepath.Join(destinationRootPath, baseDestination)
	}

	if destination == "." {
		destination = filename
	}
	return lib.Path{Path: destination}.NormalizePathSystemForAPI().String()
}

func skipOrIgnore(uploadStatus *UploadStatus) bool {
	if uploadStatus.Job().MatchesPath(uploadStatus.LocalPath()) {
		uploadStatus.Job().UpdateStatus(status.Ignored, uploadStatus, nil)
		return true
	}

	if uploadStatus.Sync {
		file, found, err := uploadStatus.Job().FindRemoteFile(uploadStatus)
		responseError, ok := err.(files_sdk.ResponseError)
		if !found || (ok && responseError.Type == "not-found") {
			uploadStatus.Job().Logger.Printf("sync %v not found on destination", uploadStatus.RemotePath())
			return false
		}
		// local is not after server
		if uploadStatus.File().Size == file.Size {
			// Server version is the same or newer
			uploadStatus.Job().UpdateStatus(status.Skipped, uploadStatus, nil)
			return true
		}
		uploadStatus.Job().Logger.Printf("sync %v found on destination with non matching sizes: local: %v, remote: %v", uploadStatus.RemotePath(), uploadStatus.File().Size, file.Size)
	}
	return false
}
