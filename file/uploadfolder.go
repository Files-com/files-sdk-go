package file

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Files-com/files-sdk-go/file/manager"
	"github.com/zenthangplus/goccm"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/file/status"
	"github.com/Files-com/files-sdk-go/folder"
	"github.com/Files-com/files-sdk-go/ignore"
	"github.com/Files-com/files-sdk-go/lib"
)

type Uploader interface {
	Upload(context.Context, io.ReaderAt, int64, files_sdk.FileActionBeginUploadParams, func(int64), goccm.ConcurrencyManager) (files_sdk.File, error)
	Find(context.Context, string) (files_sdk.File, error)
}

func uploadFolder(ctx context.Context, c Uploader, config files_sdk.Config, params *UploadParams) (status.Job, error) {
	if params.Reporter == nil {
		params.Reporter = func(uploadStatus status.Report, err error) {}
	}
	if params.Manager == nil {
		params.Manager = manager.Default()
	}
	var uploadFiles []*UploadStatus
	localFolderPath := params.Source
	destinationRootPath := params.Destination
	directoriesToCreate := make(map[string]UploadStatus)
	job := status.Job{Id: params.JobId}.Init()
	addUploads := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		dir, filename := filepath.Split(path)

		if localFolderPath == dir && filename == "" {
			return nil
		}

		i, err := ignore.New()
		if err == nil && i.MatchesPath(filename) {
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

		file := &UploadStatus{
			job:         job,
			Source:      path,
			destination: destination,
			file:        files_sdk.File{DisplayName: filepath.Base(destination), Path: destination, Size: info.Size(), Mtime: info.ModTime()},
		}
		if info.IsDir() {
			file.file.Type = "directory"
			directoriesToCreate[destination] = *file
		} else {
			file.file.Type = "file"
			job.Add(file)
			uploadFiles = append(uploadFiles, file)
			file.SetStatus(status.Queued)
			params.Reporter(*file, nil) // Only block on queued so user can wait on locks
		}
		return nil
	}
	err := filepath.Walk(localFolderPath, addUploads)

	if err != nil {
		return *job, err
	}

	if len(uploadFiles) == 0 {
		fi, err := os.Stat(localFolderPath)
		if err != nil {
			return *job, err
		}
		if !fi.IsDir() {
			if addUploads(localFolderPath, fi, nil) != nil {
				return *job, err
			}
		}
	}

	if destinationRootPath != "" {
		folderClient := folder.Client{Config: config}
		_, err := folderClient.Create(ctx, files_sdk.FolderCreateParams{Path: filepath.Clean(destinationRootPath)})
		responseError, ok := (err).(files_sdk.ResponseError)
		if err != nil && ok && responseError.ErrorMessage != "The destination exists." {
			return *job, err
		}
	}

	someMapMutex := sync.RWMutex{}
	onComplete := make(chan bool)
	for i := range uploadFiles {
		uploadStatus := uploadFiles[i]
		downloadCtx, cancel := context.WithCancel(ctx)
		uploadStatus.cancel = cancel
		uploadStatus.SetStatus(status.Queued)
		params.Manager.FilesManager.Wait()
		go func(ctx context.Context, uploadStatus *UploadStatus) {
			defer func() {
				params.Manager.FilesManager.Done()
				onComplete <- true
			}()
			if !checkUpdateSync(ctx, uploadStatus, params, c) {
				return
			}

			dir, _ := filepath.Split(uploadStatus.File().Path)
			someMapMutex.RLock()
			dirFile, ok := directoriesToCreate[filepath.Clean(dir)]
			someMapMutex.RUnlock()
			if ok {
				err := maybeCreateFolder(ctx, dirFile, config)
				dealWithDBasicError(uploadStatus, err, params)
				someMapMutex.Lock()
				delete(directoriesToCreate, filepath.Clean(dir))
				someMapMutex.Unlock()
			}
			localFile, err := os.Open(uploadStatus.Source)
			defer localFile.Close()
			if dealWithDBasicError(uploadStatus, err, params) {
				return
			}

			file, err := c.Upload(ctx, localFile, uploadStatus.File().Size, files_sdk.FileActionBeginUploadParams{Path: uploadStatus.File().Path, MkdirParents: lib.Bool(true)}, uploadProgress(params, uploadStatus), params.FilePartsManager)
			dealWithCanceledError(ctx, uploadStatus, err, file, params)
		}(downloadCtx, uploadStatus)
	}

	for !job.AllEnded() {
	}

	return *job, err
}

func checkUpdateSync(downloadCtx context.Context, uploadStatus *UploadStatus, params *UploadParams, c Uploader) bool {
	if params.Sync {
		file, err := c.Find(downloadCtx, uploadStatus.File().Path)
		responseError, ok := err.(files_sdk.ResponseError)
		if ok && responseError.Type == "not-found" {
			return true
		}
		// local is not after server
		if !uploadStatus.File().Mtime.After(file.Mtime) {
			// Server version is the same or newer
			uploadStatus.SetStatus(status.Skipped)
			uploadStatus.Cancel()
			go params.Reporter(*uploadStatus, nil)
			return false
		}
	}
	return true
}
