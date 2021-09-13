package file

import (
	"context"
	"io"
	"math"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Files-com/files-sdk-go/v2/ignore"

	"github.com/Files-com/files-sdk-go/v2/directory"
	"github.com/Files-com/files-sdk-go/v2/lib/direction"

	"github.com/Files-com/files-sdk-go/v2/file/manager"
	"github.com/zenthangplus/goccm"

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
	job.Start = func() {
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
			c.UploadFolder(ctx, params).Start()
		} else {
			c.UploadFile(ctx, params).Start()
		}
	}

	return job
}

func (c *Client) UploadFolder(ctx context.Context, params UploadParams) *status.Job {
	return uploadFolder(ctx, c, params)
}

func (c *Client) UploadFile(ctx context.Context, params UploadParams) *status.Job {
	var job *status.Job
	if params.Job == nil {
		job = status.Job{}.Init()
	} else {
		job = params.Job
	}
	SetJobParams(job, direction.UploadType, params)
	job.Client = c
	jobCtx := job.WithContext(ctx)
	job.Type = directory.File

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
		_, err := c.Find(jobCtx, params.RemotePath)
		responseError, ok := err.(files_sdk.ResponseError)
		if ok && responseError.Type == "bad-request/cannot-download-directory" {
			destination = filepath.Join(params.RemotePath, localFileName)
		} else if ok && responseError.Type == "not-found" {
			if destination[len(destination)-1:] == "/" {
				destination = filepath.Join(params.RemotePath, localFileName)
			}
		} else if err != nil {
			job.UpdateStatus(status.Errored, uploadStatus, err)
			return job
		}
	}
	job.FilesManager.Wait()
	defer job.FilesManager.Done()
	fi, _ := os.Stat(params.LocalPath)
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

	job.Start = func() {
		job.StartTime = time.Now()
		defer func() {
			job.EndTime = time.Now()
		}()
		var err error
		job.GitIgnore, err = ignore.New(params.Ignore...)
		if err != nil {
			job.Add(uploadStatus)
			job.UpdateStatus(status.Errored, uploadStatus, err)
			return
		}
		if skipOrIgnore(jobCtx, uploadStatus) {
			return
		}
		localFile, err := os.Open(params.LocalPath)
		defer localFile.Close()
		if dealWithDBasicError(uploadStatus, err) {
			return
		}
		uploadStatus.Size = fi.Size()
		beginUpload.Path = destination
		file, err := uploadStatus.Upload(jobCtx, localFile, fi.Size(), beginUpload, uploadProgress(uploadStatus), job.FilePartsManager)
		dealWithCanceledError(uploadStatus, err, file)
		RetryTransfers(jobCtx, job)
	}

	job.Wait = func() {
		for job.EndTime.IsZero() {
		}
	}

	return job
}

func (c *Client) Upload(parentCtx context.Context, reader io.ReaderAt, size int64, params files_sdk.FileBeginUploadParams, progress func(int64), cm goccm.ConcurrencyManager) (files_sdk.File, error) {
	onComplete := make(chan files_sdk.EtagsParam)
	onError := make(chan error)
	bytesWritten := int64(0)
	etags := make([]files_sdk.EtagsParam, 0)
	fileUploadPart, err := c.startUpload(parentCtx, params)
	if err != nil {
		return files_sdk.File{}, err
	}
	fileUploadPart.Path = params.Path
	count := int64(0)
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	run := func(off int64, len int64, fileUploadPart files_sdk.FileUploadPart) {
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
	}
	byteOffset(
		size,
		fileUploadPart.Partsize,
		func(off int64, len int64) {
			count += len
			if *fileUploadPart.ParallelParts {
				cm.Wait()
				go run(off, len, fileUploadPart)
			} else {
				run(off, len, fileUploadPart)
			}
			fileUploadPart.PartNumber += 1
		},
	)

	n := int64(0)
	for n < fileUploadPart.PartNumber-1 {
		n++
		select {
		case err := <-onError:
			cancel()
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
			ctx, files_sdk.FileBeginUploadParams{Path: fileUploadPart.Path, Ref: fileUploadPart.Ref, Part: fileUploadPart.PartNumber, MkdirParents: lib.Bool(true)},
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

func dealWithCanceledError(uploadStatus *UploadStatus, err error, file files_sdk.File) {
	if err != nil {
		uploadStatus.Job.StatusFromError(uploadStatus, err)
	} else {
		uploadStatus.File = file
		uploadStatus.Job.UpdateStatus(status.Complete, uploadStatus, nil)
	}
}

func dealWithDBasicError(uploadStatus *UploadStatus, err error) bool {
	if err != nil {
		uploadStatus.Job.UpdateStatus(status.Errored, uploadStatus, err)
		return true
	}
	return false
}

func uploadProgress(uploadStatus *UploadStatus) func(bytesCount int64) {
	return func(bytesCount int64) {
		uploadStatus.incrementDownloadedBytes(bytesCount)
		uploadStatus.Job.UpdateStatus(status.Uploading, uploadStatus, nil)
	}
}
