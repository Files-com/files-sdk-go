package file

import (
	"context"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Files-com/files-sdk-go/file/manager"
	"github.com/zenthangplus/goccm"

	"github.com/Files-com/files-sdk-go/file/status"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/folder"
	"github.com/Files-com/files-sdk-go/lib"
)

type UploadRetryParams struct {
	status.File
	*manager.Manager
	status.Reporter
}

func (c *Client) UploadRetry(ctx context.Context, params UploadRetryParams) status.File {
	return c.UploadFile(ctx,
		&UploadParams{
			JobId:       params.Job.Id,
			Sync:        params.Sync,
			Source:      params.LocalPath,
			Destination: params.RemotePath,
			Reporter:    params.Reporter,
			Manager:     params.Manager,
		},
	)
}

type UploadParams struct {
	JobId       string
	Sync        bool
	Source      string
	Destination string
	status.Reporter
	*manager.Manager
}

func (c *Client) UploadFolderOrFile(ctx context.Context, params *UploadParams) status.Job {
	file := status.File{File: files_sdk.File{}, LocalPath: params.Source, Status: status.Errored, Job: status.Job{}.Init(params.JobId)}
	absoluteSource, err := filepath.Abs(params.Source)
	if err != nil {
		params.Reporter(file, err)
		return *file.Job
	}
	fi, err := os.Stat(absoluteSource)
	if err != nil {
		params.Reporter(file, err)
		return *file.Job
	}

	if fi.IsDir() {
		return c.UploadFolder(ctx, params)
	} else {
		return *c.UploadFile(ctx, params).Job
	}
}

func (c *Client) UploadFolder(ctx context.Context, params *UploadParams) status.Job {
	return uploadFolder(ctx, c, c.Config, params)
}

func (c *Client) UploadFile(ctx context.Context, params *UploadParams) status.File {
	if params.Reporter == nil {
		params.Reporter = func(uploadStatus status.File, err error) {}
	}
	if params.Manager == nil {
		params.Manager = manager.Default()
	}
	job := status.Job{}.Init(params.JobId)
	jobCtx, cancel := context.WithCancel(ctx)
	job.CancelFunc = cancel
	defer job.Cancel()
	uploadStatus := &UploadStatus{
		CancelFunc: cancel,
		Job:        job,
		LocalPath:  params.Source,
		Sync:       params.Sync,
	}

	beginUpload := files_sdk.FileBeginUploadParams{}
	destination := params.Destination
	_, localFileName := filepath.Split(params.Source)
	if params.Destination == "" {
		destination = localFileName
	} else {
		_, err := c.Find(ctx, params.Destination)
		responseError, ok := err.(files_sdk.ResponseError)
		if ok && responseError.Type == "bad-request/cannot-download-directory" {
			destination = filepath.Join(params.Destination, localFileName)
		} else if ok && responseError.Type == "not-found" {
			if destination[len(destination)-1:] == "/" {
				destination = filepath.Join(params.Destination, localFileName)
				beginUpload.MkdirParents = lib.Bool(true)
			}
		} else if err != nil {
			params.Reporter(uploadStatus.ToStatusFile(), err)
			return uploadStatus.ToStatusFile()
		}
	}
	params.FilesManager.Wait()
	defer params.FilesManager.Done()
	fi, err := os.Stat(params.Source)
	uploadStatus.RemotePath = destination
	uploadStatus.File = files_sdk.File{
		DisplayName: filepath.Base(destination),
		Path:        destination,
		Type:        "file",
		Mtime:       fi.ModTime(),
		Size:        fi.Size(),
	}
	uploadStatus.Status = status.Queued
	job.Add(uploadStatus)
	if !checkUpdateSync(ctx, uploadStatus, params, c) {
		return uploadStatus.ToStatusFile()
	}
	params.Reporter(uploadStatus.ToStatusFile(), nil) // Only block on queued so user can wait on locks
	localFile, err := os.Open(params.Source)
	defer localFile.Close()
	if dealWithDBasicError(uploadStatus, err, params) {
		return uploadStatus.ToStatusFile()
	}
	uploadStatus.Size = fi.Size()
	beginUpload.Path = destination
	file, err := c.Upload(jobCtx, localFile, fi.Size(), beginUpload, uploadProgress(params, uploadStatus), params.FilePartsManager)
	dealWithCanceledError(ctx, uploadStatus, err, file, params)

	return uploadStatus.ToStatusFile()
}

func UploadFile(ctx context.Context, params *UploadParams) status.File {
	return (&Client{}).UploadFile(ctx, params)
}

func Upload(ctx context.Context, reader io.ReaderAt, size int64, beginUpload files_sdk.FileBeginUploadParams, progress func(int64), cm goccm.ConcurrencyManager) (files_sdk.File, error) {
	return (&Client{}).Upload(ctx, reader, size, beginUpload, progress, cm)
}

func (c *Client) Upload(ctx context.Context, reader io.ReaderAt, size int64, params files_sdk.FileBeginUploadParams, progress func(int64), cm goccm.ConcurrencyManager) (files_sdk.File, error) {
	onComplete := make(chan files_sdk.EtagsParam)
	onError := make(chan error)
	bytesWritten := int64(0)
	etags := make([]files_sdk.EtagsParam, 0)
	fileUploadPart, err := c.startUpload(ctx, params)
	if err != nil {
		return files_sdk.File{}, err
	}
	partReturnedError := false
	fileUploadPart.Path = params.Path
	count := int64(0)
	byteOffset(
		size,
		fileUploadPart.Partsize,
		func(off int64, len int64) {
			count += len
			if *fileUploadPart.ParallelParts {
				cm.Wait()
			}

			if partReturnedError {
				return
			}
			go func(off int64, len int64, fileUploadPart files_sdk.FileUploadPart) {
				proxyReader := &ProxyReader{
					ReaderAt: reader,
					off:      off,
					len:      len,
					onRead:   progress,
				}

				etag, bytesRead, err := c.createPart(ctx, proxyReader, len, fileUploadPart)
				if err != nil {
					if *fileUploadPart.ParallelParts {
						cm.Done()
					}
					onError <- err
					return
				}
				bytesWritten += bytesRead
				if *fileUploadPart.ParallelParts {
					cm.Done()
				}
				onComplete <- etag
			}(off, len, fileUploadPart)

			fileUploadPart.PartNumber += 1
		},
	)

	n := int64(0)
	for n < fileUploadPart.PartNumber-1 {
		n++
		select {
		case err := <-onError:
			partReturnedError = true
			return files_sdk.File{}, err
		case etag := <-onComplete:
			etags = append(etags, etag)
		}
	}

	return c.completeUpload(ctx, etags, bytesWritten, fileUploadPart.Path, fileUploadPart.Ref)
}

func (c *Client) startUpload(ctx context.Context, beginUpload files_sdk.FileBeginUploadParams) (files_sdk.FileUploadPart, error) {
	uploads, err := c.BeginUpload(ctx, beginUpload)
	if err != nil {
		return files_sdk.FileUploadPart{}, err
	}
	return uploads[0], err
}

func (c *Client) completeUpload(ctx context.Context, etags []files_sdk.EtagsParam, bytesWritten int64, path string, ref string) (files_sdk.File, error) {
	return c.Create(ctx, files_sdk.FileCreateParams{
		ProvidedMtime: time.Now(),
		EtagsParam:    etags,
		Action:        "end",
		Path:          path,
		Ref:           ref,
		Size:          bytesWritten,
	})
}

func byteOffset(size int64, blockSize int64, callback func(off int64, len int64)) {
	off := int64(0)
	endRange := blockSize
	for {
		if off < size {
			endRange = int64(math.Min(float64(endRange), float64(size)))
			callback(off, endRange-off)
			off = endRange
			endRange = off + blockSize
		} else {
			break
		}
	}
}

func (c *Client) createPart(ctx context.Context, reader io.ReadCloser, len int64, fileUploadPart files_sdk.FileUploadPart) (files_sdk.EtagsParam, int64, error) {
	var err error
	if fileUploadPart.PartNumber != 1 {
		fileUploadPart, err = c.startUpload(
			ctx, files_sdk.FileBeginUploadParams{Path: fileUploadPart.Path, Ref: fileUploadPart.Ref, Part: fileUploadPart.PartNumber},
		)
		if err != nil {
			return files_sdk.EtagsParam{}, int64(0), err
		}
	}

	headers := http.Header{}
	headers.Add("Content-Length", strconv.FormatInt(len, 10))
	res, err := files_sdk.CallRaw(
		&files_sdk.CallParams{
			Method:  fileUploadPart.HttpMethod,
			Config:  c.Config,
			Uri:     fileUploadPart.UploadUri,
			BodyIo:  reader,
			Headers: &headers,
			Context: ctx,
		},
	)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return files_sdk.EtagsParam{}, len, err
	}

	return files_sdk.EtagsParam{
		Etag: res.Header.Get("Etag"),
		Part: strconv.FormatInt(fileUploadPart.PartNumber, 10),
	}, len, nil
}

func dealWithCanceledError(ctx context.Context, uploadStatus *UploadStatus, err error, file files_sdk.File, params *UploadParams) {
	if err != nil {
		if ctx.Err() == nil {
			uploadStatus.SetStatus(status.Errored)
		} else {
			uploadStatus.SetStatus(status.Canceled)
		}
	} else {
		uploadStatus.File = file
		uploadStatus.SetStatus(status.Complete)
	}
	// Block on finishing report
	params.Reporter(uploadStatus.ToStatusFile(), err)
}

func dealWithDBasicError(uploadStatus *UploadStatus, err error, params *UploadParams) bool {
	if err != nil {
		uploadStatus.SetStatus(status.Errored)
		params.Reporter(uploadStatus.ToStatusFile(), err)
		return true
	}
	return false
}

func maybeCreateFolder(ctx context.Context, file UploadStatus, config files_sdk.Config) error {
	client := folder.Client{Config: config}
	createdFolder, err := client.Create(ctx, files_sdk.FolderCreateParams{Path: file.RemotePath + "/"})
	responseError, ok := (err).(files_sdk.ResponseError)
	if err != nil && ok && responseError.ErrorMessage != "The destination exists." {
		return err
	} else {
		file.File = createdFolder
	}

	return nil
}

func uploadProgress(params *UploadParams, uploadStatus *UploadStatus) func(bytesCount int64) {
	return func(bytesCount int64) {
		uploadStatus.UploadedBytes += bytesCount
		uploadStatus.SetStatus(status.Uploading)
		params.Reporter(uploadStatus.ToStatusFile(), nil)
	}
}
