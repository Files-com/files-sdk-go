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

func downloadFolder(ctx context.Context, files []Entity, c Downloader, params DownloadFolderParams) (status.Job, error) {
	if params.Reporter == nil {
		params.Reporter = func(_ status.Report, _ error) {}
	}
	if params.Manager == nil {
		params.Manager = manager.Default()
	}
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

	if len(files) > 1 {
		rootDestinationIsDir = true
	}
	job := status.Job{Id: params.JobId}.Init()

	for _, entity := range files {
		if entity.error != nil {
			return *job, entity.error
		}
	}
	sourceRootLen := len(strings.Split(params.Path, "/"))
	signal := make(chan *DownloadStatus)
	go runEach(ctx, files, c, params, sourceRootLen, rootDestinationIsDir, rootDestination, job, signal)

	for range files {
		s := <-signal
		if !s.Ended() {
			panic(fmt.Sprintf("<- Signal id: %v, status: %v\n", s.Id(), s.String()))
		}
	}

	return *job, nil
}

func runEach(ctx context.Context, files []Entity, c Downloader, params DownloadFolderParams, sourceRootLen int, rootDestinationIsDir bool, rootDestination string, job *status.Job, signal chan *DownloadStatus) {
	for _, entity := range files {
		downloadCtx, cancel := context.WithCancel(ctx)
		s := &DownloadStatus{
			file:        entity.file,
			destination: destinationPath(entity.file, sourceRootLen, rootDestinationIsDir, rootDestination),
			runStats:    job,
			cancel:      cancel,
		}
		job.Add(s)
		s.SetStatus(status.Queued)
		params.Reporter(*s, nil)
		params.Manager.FilesManager.Wait()
		go func(ctx context.Context, reportStatus *DownloadStatus) {
			defer func() {
				params.Manager.FilesManager.Done()
				signal <- reportStatus
			}()
			dir, _ := filepath.Split(reportStatus.Destination())
			_, err := os.Stat(dir)
			if os.IsNotExist(err) {
				os.MkdirAll(dir, 0755)
			}
			fileInfo, err := os.Stat(reportStatus.Destination())
			if err != nil && !os.IsNotExist(err) {
				go params.Reporter(*reportStatus, err)
			}
			// server is not after local
			if !os.IsNotExist(err) && params.Sync && !reportStatus.File().Mtime.After(fileInfo.ModTime()) {
				// Local version is the same or newer
				reportStatus.SetStatus(status.Skipped)
				params.Reporter(*reportStatus, nil)
				reportStatus.Cancel()
				return
			}
			partName := reportStatus.Destination() + ".part"
			out, createErr := os.Create(partName)

			if createErr != nil {
				reportStatus.SetStatus(status.Errored)
				params.Reporter(*reportStatus, createErr)

				return
			}
			downloadParams := files_sdk.FileDownloadParams{Path: reportStatus.File().Path}
			writer := lib.ProgressWriter{Writer: out}
			writer.ProgressWatcher = func(incDownloadedBytes int64) {
				reportStatus.SetStatus(status.Downloading)
				reportStatus.incrementDownloadedBytes(incDownloadedBytes)
				params.Reporter(*reportStatus, nil)
			}
			downloadParams.Writer = writer
			downloadParams.OnDownload = func(response *http.Response) {
				reportStatus.SetStatus(status.Downloading)
				reportStatus.file.Size = response.ContentLength
				params.Reporter(*reportStatus, nil)
			}

			_, err = c.Download(ctx, downloadParams)
			if err != nil {
				if ctx.Err() == nil {
					reportStatus.SetStatus(status.Errored)
				} else {
					reportStatus.SetStatus(status.Canceled)
				}
				params.Reporter(*reportStatus, err)
			}

			closeErr := out.Close()

			if closeErr != nil {
				go params.Reporter(*reportStatus, closeErr)
			}

			if reportStatus.Invalid() {
				os.Remove(partName) // Clean up on invalid download
			} else {
				err = os.Rename(partName, reportStatus.Destination())
				if err != nil {
					reportStatus.SetStatus(status.Errored)
				} else if reportStatus.Downloading() {
					reportStatus.SetStatus(status.Complete)
				}
				params.Reporter(*reportStatus, err)
			}
		}(downloadCtx, s)
	}
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
