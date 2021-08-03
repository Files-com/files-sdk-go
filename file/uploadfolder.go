package file

import (
	"context"
	"fmt"
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

func uploadFolder(ctx context.Context, c Uploader, config files_sdk.Config, params *UploadParams) status.Job {
	if params.Reporter == nil {
		params.Reporter = func(uploadStatus status.File, err error) {}
	}
	if params.Manager == nil {
		params.Manager = manager.Default()
	}
	var uploadFiles []*UploadStatus
	localFolderPath := params.Source
	destinationRootPath := params.Destination
	directoriesToCreate := make(map[string]UploadStatus)
	job := status.Job{}.Init(params.JobId)

	metaFile := UploadStatus{
		File:      files_sdk.File{DisplayName: filepath.Base(localFolderPath)},
		Job:       job,
		Status:    status.Errored,
		LocalPath: localFolderPath,
		Sync:      false,
	}
	metaStata, statErr := os.Stat(localFolderPath)
	if statErr != nil {
		job.Add(metaFile)
		params.Reporter(metaFile.ToStatusFile(), statErr)
		return *job
	}
	if metaStata.IsDir() {
		metaFile.File.Type = "directory"
	} else {
		metaFile.File.Type = "file"
	}

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
			Job:        job,
			RemotePath: destination,
			LocalPath:  path,
			Sync:       params.Sync,
			File:       files_sdk.File{DisplayName: filepath.Base(destination), Path: destination, Size: info.Size(), Mtime: info.ModTime()},
		}
		if info.IsDir() {
			file.File.Type = "directory"
			directoriesToCreate[destination] = *file
		} else {
			file.File.Type = "file"
			job.Add(file)
			uploadFiles = append(uploadFiles, file)
			file.SetStatus(status.Queued)
			params.Reporter(file.ToStatusFile(), nil)
		}
		return nil
	}
	walkErr := filepath.Walk(localFolderPath, addUploads)

	if walkErr != nil {
		job.Add(metaFile)
		params.Reporter(metaFile.ToStatusFile(), walkErr)
		return *job
	}

	if len(uploadFiles) == 0 {
		if metaFile.File.Type == "File" {
			if addUploads(localFolderPath, metaStata, nil) != nil {
				job.Add(metaFile)
				params.Reporter(metaFile.ToStatusFile(), walkErr)
				return *job
			}
		}
	}

	if destinationRootPath != "" {
		folderClient := folder.Client{Config: config}
		_, err := folderClient.Create(ctx, files_sdk.FolderCreateParams{Path: filepath.Clean(destinationRootPath)})
		responseError, ok := (err).(files_sdk.ResponseError)
		if err != nil && ok && responseError.ErrorMessage != "The destination exists." {
			job.Add(metaFile)
			params.Reporter(metaFile.ToStatusFile(), walkErr)
			return *job
		}
	}

	someMapMutex := sync.RWMutex{}
	onComplete := make(chan *UploadStatus)
	for i := range uploadFiles {
		uploadStatus := uploadFiles[i]
		downloadCtx, cancel := context.WithCancel(ctx)
		uploadStatus.CancelFunc = cancel
		uploadStatus.SetStatus(status.Queued)
		params.Manager.FilesManager.Wait()
		go func(ctx context.Context, uploadStatus *UploadStatus) {
			defer func() {
				params.Manager.FilesManager.Done()
				onComplete <- uploadStatus
			}()
			if !checkUpdateSync(ctx, uploadStatus, params, c) {
				return
			}

			dir, _ := filepath.Split(uploadStatus.File.Path)
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
			localFile, err := os.Open(uploadStatus.LocalPath)
			defer localFile.Close()
			if dealWithDBasicError(uploadStatus, err, params) {
				return
			}

			file, err := c.Upload(ctx, localFile, uploadStatus.Size, files_sdk.FileActionBeginUploadParams{Path: uploadStatus.RemotePath, MkdirParents: lib.Bool(true)}, uploadProgress(params, uploadStatus), params.FilePartsManager)
			dealWithCanceledError(ctx, uploadStatus, err, file, params)
		}(downloadCtx, uploadStatus)
	}

	for range uploadFiles {
		s := <-onComplete
		if s.Running() {
			panic(fmt.Sprintf("<- Signal id: %v, status: %v\n", s.Id(), s.String()))
		}
	}

	return *job
}

func checkUpdateSync(downloadCtx context.Context, uploadStatus *UploadStatus, params *UploadParams, c Uploader) bool {
	if params.Sync {
		file, err := c.Find(downloadCtx, uploadStatus.RemotePath)
		responseError, ok := err.(files_sdk.ResponseError)
		if ok && responseError.Type == "not-found" {
			return true
		}
		// local is not after server
		if !uploadStatus.Mtime.After(file.Mtime) {
			// Server version is the same or newer
			uploadStatus.SetStatus(status.Skipped)
			params.Reporter(uploadStatus.ToStatusFile(), nil)
			return false
		}
	}
	return true
}
