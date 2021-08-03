package file

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Files-com/files-sdk-go/file/manager"

	"github.com/Files-com/files-sdk-go/file/status"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/lib"
)

type Downloader interface {
	Download(context.Context, files_sdk.FileDownloadParams) (files_sdk.File, error)
}

func downloadFolder(ctx context.Context, files []Entity, c Downloader, params DownloadFolderParams) status.Job {
	if params.Reporter == nil {
		params.Reporter = func(_ status.File, _ error) {}
	}
	if params.Manager == nil {
		params.Manager = manager.Default()
	}
	job := status.Job{}.Init(params.JobId)
	jobCtx, cancel := context.WithCancel(ctx)
	job.CancelFunc = cancel
	rootDestination := params.RootDestination
	rootDestinationIsDir := false
	if rootDestination != "" && rootDestination[len(rootDestination)-1:] == string(os.PathSeparator) {
		rootDestinationIsDir = true
	} else {
		rootDestination, _ = filepath.Abs(rootDestination)
		fi, err := os.Stat(rootDestination)
		if err == nil && fi.Mode().IsDir() {
			rootDestinationIsDir = true
		}
	}

	metaDownload := DownloadStatus{
		File:       files_sdk.File{},
		LocalPath:  params.RootDestination,
		RemotePath: params.Path,
		Sync:       params.Sync,
		Job:        job,
	}
	for _, entity := range files {
		if entity.error != nil {
			metaDownload.SetStatus(status.Errored)
			job.Add(metaDownload)
			params.Reporter(metaDownload.ToStatusFile(), entity.error)
			return *job
		}
	}

	if len(files) > 1 {
		rootDestinationIsDir = true
	}

	sourceRootLen := len(strings.Split(params.Path, "/"))
	signal := make(chan *DownloadStatus)
	go runEach(jobCtx, files, c, params, sourceRootLen, rootDestinationIsDir, rootDestination, job, signal)

	for range files {
		s := <-signal
		if !s.Ended() {
			panic(fmt.Sprintf("<- Signal id: %v, status: %v\n", s.Id(), s.String()))
		}
	}

	return *job
}

func runEach(ctx context.Context, files []Entity, c Downloader, params DownloadFolderParams, sourceRootLen int, rootDestinationIsDir bool, rootDestination string, job *status.Job, signal chan *DownloadStatus) {
	for _, entity := range files {
		downloadCtx, cancel := context.WithCancel(ctx)
		s := &DownloadStatus{
			File:       entity.file,
			LocalPath:  destinationPath(entity.file, sourceRootLen, rootDestinationIsDir, rootDestination),
			RemotePath: entity.file.Path,
			Job:        job,
			CancelFunc: cancel,
			Sync:       params.Sync,
		}
		job.Add(s)
		s.SetStatus(status.Queued)
		params.Reporter(s.ToStatusFile(), nil)
		params.Manager.FilesManager.Wait()
		go downloadFolderItem(params, signal, c, downloadCtx, s)
	}
}

func downloadFolderItem(params DownloadFolderParams, signal chan *DownloadStatus, c Downloader, downloadCtx context.Context, s *DownloadStatus) {
	func(ctx context.Context, reportStatus *DownloadStatus) {
		defer func() {
			params.Manager.FilesManager.Done()
			signal <- reportStatus
		}()
		dir, _ := filepath.Split(reportStatus.LocalPath)
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				params.Reporter(reportStatus.ToStatusFile(), err)
				return
			}
		}
		fileInfo, err := os.Stat(reportStatus.LocalPath)
		if err != nil && !os.IsNotExist(err) {
			params.Reporter(reportStatus.ToStatusFile(), err)
			return
		}
		// server is not after local
		if !os.IsNotExist(err) && params.Sync && !reportStatus.Mtime.After(fileInfo.ModTime()) {
			// Local version is the same or newer
			reportStatus.SetStatus(status.Skipped)
			params.Reporter(reportStatus.ToStatusFile(), nil)
			return
		}
		downloadParams := files_sdk.FileDownloadParams{Path: reportStatus.RemotePath}

		tmpName := tmpDownloadPath(reportStatus.LocalPath)
		var out *os.File
		out, downloadParams.Writer = openFile(tmpName, reportStatus, params)

		downloadParams.OnDownload = func(response *http.Response) {
			reportStatus.SetStatus(status.Downloading)
			reportStatus.Size = response.ContentLength
			params.Reporter(reportStatus.ToStatusFile(), nil)
		}

		_, err = c.Download(ctx, downloadParams)
		if IsStreamError(err) {
			out.Close()
			// retry
			out, downloadParams.Writer = openFile(tmpName, reportStatus, params)
			_, err = c.Download(ctx, downloadParams)
		}
		if err != nil {
			if ctx.Err() == nil {
				reportStatus.SetStatus(status.Errored)
			} else {
				reportStatus.SetStatus(status.Canceled)
			}
			params.Reporter(reportStatus.ToStatusFile(), err)
		}

		closeErr := out.Close()

		if closeErr != nil {
			params.Reporter(reportStatus.ToStatusFile(), closeErr)
		}

		if reportStatus.Invalid() {
			os.Remove(tmpName) // Clean up on invalid download
		} else {
			err = os.Rename(tmpName, reportStatus.LocalPath)
			if err != nil {
				reportStatus.SetStatus(status.Errored)
			} else if reportStatus.Downloading() {
				reportStatus.SetStatus(status.Complete)
			}
			params.Reporter(reportStatus.ToStatusFile(), err)
		}
	}(downloadCtx, s)
}

func tmpDownloadPath(path string) string {
	return _tmpDownloadPath(path, 0)
}

func _tmpDownloadPath(path string, index int) string {
	var name string

	if index == 0 {
		name = fmt.Sprintf("%v.download", path)
	} else {
		name = fmt.Sprintf("%v.download (%v)", path, index)
	}
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return name
	}
	return _tmpDownloadPath(path, index+1)
}

func openFile(partName string, reportStatus *DownloadStatus, params DownloadFolderParams) (*os.File, lib.ProgressWriter) {
	out, createErr := os.Create(partName)
	if createErr != nil {
		reportStatus.SetStatus(status.Errored)
		params.Reporter(reportStatus.ToStatusFile(), createErr)
	}
	writer := lib.ProgressWriter{Writer: out}
	writer.ProgressWatcher = func(incDownloadedBytes int64) {
		reportStatus.SetStatus(status.Downloading)
		reportStatus.incrementDownloadedBytes(incDownloadedBytes)
		params.Reporter(reportStatus.ToStatusFile(), nil)
	}
	return out, writer
}

func destinationPath(file files_sdk.File, sourceRootLen int, rootDestinationIsDir bool, rootDestination string) string {
	sep := strings.Split(file.Path, "/")
	r := int(math.Min(float64(len(sep)-1), float64(sourceRootLen)))
	filePathCompacted := strings.Join(sep[r:], string(os.PathSeparator))
	filePath, fileName := filepath.Split(filePathCompacted)
	var path string
	if rootDestinationIsDir {
		path = filepath.Join(rootDestination, filePath, fileName)
	} else {
		path = filepath.Join(rootDestination, filePath)
	}
	return path
}

func IsStreamError(err error) bool {
	if err != nil && strings.Contains(err.Error(), "stream error: stream ID 1") {
		return true
	}
	return false
}
