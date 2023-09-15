package file

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/directory"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/Files-com/files-sdk-go/v3/ignore"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/lib/direction"
)

type Uploader interface {
	UploadWithResume(...UploadOption) (UploadResumable, error)
	Find(files_sdk.FileFindParams, ...files_sdk.RequestResponseOption) (files_sdk.File, error)
}

func uploader(parentCtx context.Context, c Uploader, params UploaderParams) *Job {
	var job *Job
	if params.Job == nil {
		job = (&Job{}).Init()
	} else {
		job = params.Job
	}
	SetJobParams(job, direction.UploadType, params, params.config.Logger, (&FS{}).Init(params.config, true))
	job.Config = params.config
	jobCtx := job.WithContext(parentCtx)

	fi, statErr := os.Stat(params.LocalPath)

	if len(params.LocalPaths) > 0 {
		job.Type = directory.Files
	} else {
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
	}
	job.Client = c
	onComplete := make(chan *UploadStatus)
	job.CodeStart = func() {
		job.Scan()

		go enqueueIndexedUploads(job, jobCtx, onComplete)
		WaitTellFinished(job, onComplete, func() { RetryByPolicy(jobCtx, job, job.RetryPolicy.(RetryPolicy), false) })

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
		job.Ignore, err = ignore.New(params.Ignore...)
		if err != nil {
			job.Add(metaFile)
			job.UpdateStatus(status.Errored, metaFile, err)
			job.EndScan()
			job.Finish()
			return
		}
		if len(params.Include) > 0 {
			job.Include, err = ignore.New(params.Include...)
			if err != nil {
				job.Add(metaFile)
				job.UpdateStatus(status.Errored, metaFile, err)
				job.EndScan()
				job.Finish()
				return
			}
		}
		if job.Type == directory.File {
			metaFile.file = files_sdk.File{
				DisplayName: filepath.Base(params.LocalPath),
				Type:        "file",
				Mtime:       lib.Time(fi.ModTime()),
				Size:        fi.Size(),
				Path:        params.RemotePath,
			}
			metaFile.remotePath, err = remotePath(jobCtx, params, c, job)
			if err != nil {
				job.Add(metaFile)
				job.UpdateStatus(status.Errored, metaFile, err)
				job.EndScan()
				job.Finish()
				return
			}
			metaFile.file.Path = metaFile.remotePath
			job.Add(metaFile)
			job.UpdateStatus(status.Indexed, metaFile, nil)
		} else if job.Type == directory.Files {
			for _, path := range params.LocalPaths {
				file, err := os.Stat(path)
				uploadStatus, ok := buildUploadStatus(path, params.LocalPath, params.RemotePath, c, job, params)
				if !ok {
					continue
				}

				uploadStatus.file.Type = "file"
				if err != nil {
					uploadStatus.missingStat = true
					uploadStatus.error = err
				} else {
					uploadStatus.file.Size = file.Size()
					uploadStatus.file.Mtime = lib.Time(fi.ModTime())
				}
				job.Add(&uploadStatus)
			}
		} else {
			it := (&lib.Walk[lib.DirEntry]{
				FS:                 os.DirFS(params.LocalPath),
				ConcurrencyManager: job.Manager.DirectoryListingManager,
				WalkFile:           lib.DirEntryWalkFile,
			}).Walk(jobCtx)

			for it.Next() {
				uploadStatus, ok := buildUploadStatus(filepath.Join(params.LocalPath, it.Resource().Path()), params.LocalPath, params.RemotePath, c, job, params)
				if !ok {
					continue
				}

				uploadStatus.file.Type = "file"
				if it.Resource().Err() != nil {
					uploadStatus.missingStat = true
					uploadStatus.error = it.Resource().Err()
				} else {
					uploadStatus.file.Size = it.Resource().FileInfo.Size()
					uploadStatus.file.Mtime = lib.Time(it.Resource().FileInfo.ModTime())
				}
				job.Add(&uploadStatus)
			}

			if it.Err() != nil {
				job.Add(metaFile)
				job.UpdateStatus(status.Errored, metaFile, it.Err())
			}
		}

		job.EndScan()
	}

	return job
}

func remotePath(ctx context.Context, params UploaderParams, c Uploader, job *Job) (string, error) {
	destination := params.RemotePath
	_, localFileName := filepath.Split(params.LocalPath)
	if params.RemotePath == "" {
		destination = localFileName
	} else {
		var err error
		var remoteFile files_sdk.File
		if job.FilePartsManager.WaitWithContext(ctx) {
			remoteFile, err = c.Find(files_sdk.FileFindParams{Path: lib.NewUrlPath(params.RemotePath).PruneEndingSlash().String()}, files_sdk.WithContext(ctx))
			job.FilePartsManager.Done()
		} else {
			return "", ctx.Err()
		}
		var responseError files_sdk.ResponseError
		ok := errors.As(err, &responseError)
		if remoteFile.Type == "directory" {
			destination = lib.UrlJoinNoEscape(params.RemotePath, localFileName)
		} else if ok && responseError.Type == "not-found" {
			if destination[len(destination)-1:] == "/" {
				destination = lib.UrlJoinNoEscape(params.RemotePath, localFileName)
			}
		} else if err != nil {
			return "", err
		}
	}
	return destination, nil
}

func enqueueIndexedUploads(job *Job, jobCtx context.Context, onComplete chan *UploadStatus) {
	for !job.EndScanning.Called() || job.Count(status.Indexed) > 0 {
		if f, ok := job.EnqueueNext(); ok {
			if job.FilesManager.WaitWithContext(jobCtx) {
				go enqueueUpload(jobCtx, job, f.(*UploadStatus), onComplete)
			} else {
				job.UpdateStatus(status.Canceled, f.(*UploadStatus), nil)
				onComplete <- f.(*UploadStatus)
			}
		}
	}
}

func enqueueUpload(ctx context.Context, job *Job, uploadStatus *UploadStatus, onComplete chan *UploadStatus) {
	finish := func() {
		job.FilesManager.Done()
		onComplete <- uploadStatus
	}
	if uploadStatus.error != nil || uploadStatus.missingStat {
		job.UpdateStatus(status.Errored, uploadStatus, uploadStatus.RecentError())
		finish()
		return
	}
	func() {
		var localFile *os.File
		var err error
		defer func() {
			if localFile != nil {
				localFile.Close()
			}
			finish()
		}()
		if skipOrIgnore(uploadStatus, job.Config.FeatureFlag("incremental-updates")) {
			return
		}
		if uploadStatus.dryRun {
			uploadStatus.Job().UpdateStatus(status.Complete, uploadStatus, nil)
			return
		}
		localFile, err = os.Open(uploadStatus.LocalPath())
		if err != nil {
			uploadStatus.Job().UpdateStatus(status.Errored, uploadStatus, err)
			return
		}
		uploadStatus.UploadResumable, err = uploadStatus.UploadWithResume(
			UploadWithContext(ctx),
			UploadWithManager(job.FilePartsManager),
			UploadWithReaderAt(localFile),
			UploadWithSize(uploadStatus.File().Size),
			UploadWithResume(uploadStatus.UploadResumable),
			UploadWithProgress(uploadProgress(uploadStatus)),
			UploadWithProvidedMtime(*uploadStatus.File().Mtime),
			UploadWithDestinationPath(uploadStatus.RemotePath()),
		)
		if err != nil {
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

func buildUploadStatus(path string, localFolderPath string, destinationRootPath string, c Uploader, job *Job, params UploaderParams) (UploadStatus, bool) {
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
		dryRun:     params.DryRun,
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

func skipOrIgnore(uploadStatus *UploadStatus, incrementalUpdates bool) bool {
	if uploadStatus.Job().Ignore.MatchesPath(uploadStatus.LocalPath()) {
		uploadStatus.Job().UpdateStatus(status.Ignored, uploadStatus, nil)
		return true
	}

	if uploadStatus.Job().Include != nil && !uploadStatus.Job().Include.MatchesPath(uploadStatus.LocalPath()) {
		uploadStatus.Job().UpdateStatus(status.Ignored, uploadStatus, nil)
		return true
	}

	if uploadStatus.Sync {
		file, found, err := uploadStatus.Job().FindRemoteFile(uploadStatus)
		var responseError files_sdk.ResponseError
		ok := errors.As(err, &responseError)
		if !found || (ok && responseError.Type == "not-found") {
			uploadStatus.Job().Logger.Printf("sync %v not found on destination", uploadStatus.RemotePath())
			return false
		}
		// local is not after server
		if uploadStatus.File().Size == file.Size {
			// Server version is the same or newer
			uploadStatus.Job().UpdateStatus(status.Skipped, uploadStatus, nil)
			uploadStatus.Job().Logger.Printf("sync %v size match", uploadStatus.RemotePath())
			return true
		}
		if incrementalUpdates && file.Mtime != nil && !file.Mtime.IsZero() {
			if file.Mtime.After(*uploadStatus.File().Mtime) {
				uploadStatus.Job().Logger.Printf("sync incremental-updates %v server has a newer version", uploadStatus.RemotePath())
				return true
			}
			if file.Mtime.Truncate(time.Minute).Equal(uploadStatus.File().Mtime.Truncate(time.Minute)) {
				uploadStatus.Job().Logger.Printf("sync incremental-updates %v both times are within the same minute", uploadStatus.RemotePath())
				return true
			}
		}
		uploadStatus.Job().Logger.Printf("sync %v found on destination with non matching sizes: local: %v, remote: %v", uploadStatus.RemotePath(), uploadStatus.File().Size, file.Size)
	}
	return false
}
