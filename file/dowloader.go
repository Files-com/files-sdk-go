package file

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Files-com/files-sdk-go/v2/directory"
	"github.com/Files-com/files-sdk-go/v2/file/manager"
	"github.com/Files-com/files-sdk-go/v2/file/status"
	"github.com/Files-com/files-sdk-go/v2/lib/direction"
	"github.com/bradfitz/iter"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/lib"
)

func downloader(ctx context.Context, fileSys fs.FS, params DownloadFolderParams) *status.Job {
	job := status.Job{}.Init()
	SetJobParams(job, direction.DownloadType, params)
	jobCtx := job.WithContext(ctx)
	remoteFs, ok := fileSys.(FS)
	if ok {
		fileSys = remoteFs.WithContext(jobCtx)
	}
	if params.RemoteFile.Path != "" {
		job.LocalPath = params.LocalPath
		job.RemotePath = params.RemoteFile.Path

		if params.RemoteFile.Type == "directory" {
			job.Type = directory.Dir
		} else {
			job.Type = directory.File
		}
	} else {
		job.LocalPath = params.LocalPath
		job.RemotePath = params.RemotePath
		if job.RemotePath == "" {
			job.RemotePath = "."
		}
		stats, err := os.Stat(params.LocalPath)
		if os.IsNotExist(err) {
			if params.LocalPath == "" || params.LocalPath[len(params.LocalPath)-1:] == string(os.PathSeparator) {
				job.Type = directory.Dir
			} else {
				job.Type = directory.File
			}
		} else {
			if stats.IsDir() {
				job.Type = directory.Dir
			} else {
				job.Type = directory.File
			}
		}
	}
	onComplete := make(chan *DownloadStatus)
	count := 0
	job.Start = func() {
		job.Scanning = true
		go enqueueIndexedDownloads(job, jobCtx, onComplete)
		job.StartTime = time.Now()
		fs.WalkDir(fileSys, job.RemotePath, func(path string, d fs.DirEntry, err error) error {
			if job.Stopped {
				return jobCtx.Err()
			}
			if err != nil {
				count += 1
				createIndexedStatus(Entity{error: err}, params, job)
				return err
			}
			if !d.IsDir() {
				count += 1
				f, err := fileSys.Open(path)
				createIndexedStatus(Entity{error: err, File: f}, params, job)
			}

			return nil
		})
		job.Scanning = false
		go markDownloadOnComplete(count, onComplete, jobCtx, job)
	}

	job.Wait = func() {
		for job.EndTime.IsZero() {
		}
	}

	return job
}

func markDownloadOnComplete(count int, onComplete chan *DownloadStatus, jobCtx context.Context, job *status.Job) {
	for range iter.N(count) {
		select {
		case <-jobCtx.Done():
			break
		case <-onComplete:
		}
	}
	close(onComplete)
	RetryTransfers(jobCtx, job)
	job.EndTime = time.Now()
}

func enqueueIndexedDownloads(job *status.Job, jobCtx context.Context, onComplete chan *DownloadStatus) {
	for job.Scanning || job.Count(status.Indexed) > 0 {
		f, ok := job.Find(status.Indexed)
		if ok {
			enqueueDownload(jobCtx, job, f.(*DownloadStatus), onComplete)
		}
	}
}

func normalizePath(rootDestination string) string {
	if rootDestination != "" && rootDestination[len(rootDestination)-1:] == string(os.PathSeparator) {
	} else {
		rootDestination, _ = filepath.Abs(rootDestination)
	}
	return rootDestination
}

func createIndexedStatus(f Entity, params DownloadFolderParams, job *status.Job) {
	var fi files_sdk.File
	if f.error == nil {
		info, err := f.File.Stat()
		if err != nil {
			panic(err)
		}
		fi = info.Sys().(files_sdk.File)
	}

	s := &DownloadStatus{
		error:      f.error,
		fsFile:     f.File,
		File:       fi,
		LocalPath:  localPath(fi, *job),
		RemotePath: fi.Path,
		Job:        job,
		Sync:       params.Sync,
		Status:     status.Indexed,
	}
	job.Add(s)
}

func enqueueDownload(ctx context.Context, job *status.Job, downloadStatus *DownloadStatus, signal chan *DownloadStatus) {
	if downloadStatus.error != nil {
		job.UpdateStatus(status.Errored, downloadStatus, downloadStatus.error)
		signal <- downloadStatus
		return
	}
	job.UpdateStatus(status.Queued, downloadStatus, nil)
	if manager.Wait(ctx, job.FilesManager) {
		go downloadFolderItem(ctx, signal, downloadStatus)
	} else {
		job.UpdateStatus(status.Canceled, downloadStatus, nil)
		signal <- downloadStatus
	}
}

func downloadFolderItem(ctx context.Context, signal chan *DownloadStatus, s *DownloadStatus) {
	func(ctx context.Context, reportStatus *DownloadStatus) {
		defer func() {
			s.Job.FilesManager.Done()
			signal <- reportStatus
		}()
		dir, _ := filepath.Split(reportStatus.LocalPath)
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				reportStatus.Job.UpdateStatus(status.Errored, reportStatus, err)
				return
			}
		}
		localStat, err := os.Stat(reportStatus.LocalPath)
		if err != nil && !os.IsNotExist(err) {
			reportStatus.Job.UpdateStatus(status.Errored, reportStatus, err)
			return
		}
		remoteStat, err := reportStatus.fsFile.Stat()
		if err != nil {
			reportStatus.Job.UpdateStatus(status.Errored, reportStatus, err)
			return
		}
		// server is not after local
		if !os.IsNotExist(err) && reportStatus.Job.Sync && !remoteStat.ModTime().After(localStat.ModTime()) {
			// Local version is the same or newer
			reportStatus.Job.UpdateStatus(status.Skipped, reportStatus, nil)
			return
		}
		downloadParams := files_sdk.FileDownloadParams{Path: reportStatus.RemotePath}

		tmpName := tmpDownloadPath(reportStatus.LocalPath)
		var out *os.File
		out, downloadParams.Writer = openFile(tmpName, reportStatus)
		written, err := io.Copy(downloadParams.Writer, lib.NewReader(ctx, s.fsFile))
		if err != nil {
			reportStatus.Job.StatusFromError(reportStatus, err)
		} else {
			reportStatus.Size = written
		}
		closeErr := out.Close()

		if closeErr != nil {
			reportStatus.Job.UpdateStatus(status.Errored, reportStatus, closeErr)
		}

		closeErr = s.fsFile.Close()

		if closeErr != nil {
			reportStatus.Job.UpdateStatus(status.Errored, reportStatus, closeErr)
		}

		if !reportStatus.Is(status.Valid...) {
			os.Remove(tmpName) // Clean up on invalid download
		} else {
			err = os.Rename(tmpName, reportStatus.LocalPath)
			if err != nil {
				reportStatus.Job.UpdateStatus(status.Errored, reportStatus, err)
			} else if reportStatus.Is(status.Downloading) {
				reportStatus.Job.UpdateStatus(status.Complete, reportStatus, nil)
			}
		}
	}(ctx, s)
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

func openFile(partName string, reportStatus *DownloadStatus) (*os.File, lib.ProgressWriter) {
	out, createErr := os.Create(partName)
	if createErr != nil {
		reportStatus.Job.UpdateStatus(status.Errored, reportStatus, createErr)
	}
	writer := lib.ProgressWriter{Writer: out}
	writer.ProgressWatcher = func(incDownloadedBytes int64) {
		reportStatus.Job.UpdateStatus(status.Downloading, reportStatus, nil)
		reportStatus.incrementDownloadedBytes(incDownloadedBytes)
	}
	return out, writer
}

func localPath(file files_sdk.File, job status.Job) string {
	var path string
	if job.Type == directory.File {
		path = job.LocalPath
	} else {
		path = filepath.Join(normalizePath(job.LocalPath), compactPath(job, file))
	}

	return path
}

func compactPath(job status.Job, file files_sdk.File) string {
	sourceRootLen := len(strings.Split(job.RemotePath, "/"))
	sep := strings.Split(file.Path, "/")
	r := int(math.Min(float64(len(sep)-1), float64(sourceRootLen)))
	filePathCompacted := strings.Join(sep[r:], string(os.PathSeparator))
	return filePathCompacted
}
