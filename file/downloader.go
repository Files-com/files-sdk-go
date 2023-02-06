package file

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Files-com/files-sdk-go/v2/lib/keyvalue"

	"github.com/Files-com/files-sdk-go/v2/lib/direction"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/directory"
	"github.com/Files-com/files-sdk-go/v2/file/manager"
	"github.com/Files-com/files-sdk-go/v2/file/status"
	"github.com/Files-com/files-sdk-go/v2/lib"
)

func downloader(ctx context.Context, fileSys fs.FS, params DownloaderParams) *status.Job {
	job := status.Job{}.Init()
	SetJobParams(job, direction.DownloadType, params, params.Config.Logger(), fileSys)
	job.Config = params.Config
	jobCtx := job.WithContext(ctx)
	remoteFs, ok := fileSys.(WithContext)
	if ok {
		fileSys = remoteFs.WithContext(jobCtx).(*FS)
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
		job.RemotePath = lib.Path{Path: params.RemotePath}.PruneStartingSlash().String()
		if job.RemotePath == "" {
			job.RemotePath = "."
		}
		var remoteType directory.Type
		remoteFile, err := fileSys.Open(params.RemotePath)
		remoteType = directory.Dir // default to Dir not found error will have to be dealt with downstream
		if err == nil {
			remoteStat, err := remoteFile.Stat()
			if err == nil {
				if remoteStat.IsDir() {
					remoteType = directory.Dir
				} else {
					remoteType = directory.File
				}
			}
		}
		job.LocalPath = lib.ExpandTilde(job.LocalPath)
		var localType directory.Type
		stats, err := os.Stat(job.LocalPath)
		if os.IsNotExist(err) {
			if (lib.Path{Path: job.LocalPath}).EndingSlash() { // explicit directory
				localType = directory.Dir
			} else if remoteType == directory.File {
				localType = directory.File
			} else {
				localType = directory.Dir // implicit directory
			}
		} else {
			if stats.IsDir() {
				localType = directory.Dir
			} else {
				localType = directory.File
			}
		}
		if (!(lib.Path{Path: params.RemotePath}).EndingSlash() && localType == directory.Dir) || remoteType == directory.File && localType == directory.Dir {
			job.LocalPath = filepath.Join(job.LocalPath, (lib.Path{Path: job.RemotePath}).Pop())
			if remoteType == directory.File {
				localType = directory.File
			}
		}

		// Use relative path
		if job.LocalPath == "" {
			job.LocalPath = lib.Path{Path: job.RemotePath}.Pop()
		}

		job.Type = localType
		job.Logger.Printf(keyvalue.New(map[string]interface{}{
			"LocalPath":  job.LocalPath,
			"RemotePath": job.RemotePath,
		}))
	}
	onComplete := make(chan *DownloadStatus)
	job.CodeStart = func() {
		job.Scan()
		go enqueueIndexedDownloads(job, jobCtx, onComplete)
		status.WaitTellFinished(job, onComplete, func() { RetryByPolicy(jobCtx, job, job.RetryPolicy.(RetryPolicy), false) })

		it := (&lib.Walk[lib.DirEntry]{
			FS:                 fileSys,
			Root:               lib.UrlJoinNoEscape(job.RemotePath),
			ConcurrencyManager: job.Manager.FilePartsManager,
			WalkFile:           lib.DirEntryWalkFile,
		}).Walk(jobCtx)

		for it.Next() {
			dirEntry := it.Current()
			if it.Err() != nil {
				createIndexedStatus(Entity{error: it.Err()}, params, job)
			}

			f, err := fileSys.Open(dirEntry.Path())
			createIndexedStatus(Entity{error: err, File: f, FS: fileSys}, params, job)
		}

		if it.Err() != nil {
			metaFile := &DownloadStatus{
				job:        job,
				status:     status.Errored,
				localPath:  params.LocalPath,
				remotePath: params.RemotePath,
				Sync:       params.Sync,
				Mutex:      &sync.RWMutex{},
			}
			metaFile.file = files_sdk.File{
				DisplayName: filepath.Base(params.LocalPath),
				Type:        job.Direction.Name(),
				Path:        params.RemotePath,
			}
			job.Add(metaFile)
			job.UpdateStatus(status.Errored, metaFile, it.Err())
		}

		job.EndScan()
	}

	return job
}

func enqueueIndexedDownloads(job *status.Job, jobCtx context.Context, onComplete chan *DownloadStatus) {
	for !job.EndScanning.Called() || job.Count(status.Indexed) > 0 {
		f, ok := job.EnqueueNext()
		if ok {
			go enqueueDownload(jobCtx, job, f.(*DownloadStatus), onComplete)
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

func createIndexedStatus(f Entity, params DownloaderParams, job *status.Job) {
	s := &DownloadStatus{
		error:         f.error,
		fsFile:        f.File,
		FS:            f.FS,
		job:           job,
		Sync:          params.Sync,
		status:        status.Indexed,
		Mutex:         &sync.RWMutex{},
		PreserveTimes: params.PreserveTimes,
		dryRun:        params.DryRun,
	}
	var err error
	if f.error == nil {
		s.FileInfo, err = f.File.Stat()
		if err == nil {
			s.file = s.FileInfo.Sys().(files_sdk.File)
			s.localPath = localPath(s.file, *job)
			s.remotePath = s.file.Path
		} else {
			s.SetStatus(status.Errored, err)
		}
	}

	job.Add(s)
}

func enqueueDownload(ctx context.Context, job *status.Job, downloadStatus *DownloadStatus, signal chan *DownloadStatus) {
	if downloadStatus.error != nil || downloadStatus.fsFile == nil {
		job.UpdateStatus(status.Errored, downloadStatus, downloadStatus.RecentError())
		signal <- downloadStatus
		return
	}
	if manager.Wait(ctx, job.FilesManager) {
		downloadFolderItem(ctx, signal, downloadStatus)
	} else {
		job.UpdateStatus(status.Canceled, downloadStatus, nil)
		signal <- downloadStatus
	}
}

func downloadFolderItem(ctx context.Context, signal chan *DownloadStatus, s *DownloadStatus) {
	func(ctx context.Context, reportStatus *DownloadStatus) {
		defer func() {
			s.job.FilesManager.Done()
			signal <- reportStatus
		}()
		dir, _ := filepath.Split(reportStatus.LocalPath())
		if dir != "" {
			_, err := os.Stat(dir)
			if os.IsNotExist(err) {
				err = os.MkdirAll(dir, 0755)
				if err != nil {
					reportStatus.Job().UpdateStatus(status.Errored, reportStatus, err)
					return
				}
			}
		}

		remoteStat, remoteStatErr := reportStatus.fsFile.Stat()
		if remoteStatErr != nil {
			reportStatus.Job().UpdateStatus(status.Errored, reportStatus, remoteStatErr)
			return
		}

		if reportStatus.Job().Sync {
			localStat, localStatErr := os.Stat(reportStatus.LocalPath())
			if localStatErr != nil && !os.IsNotExist(localStatErr) {
				reportStatus.Job().UpdateStatus(status.Errored, reportStatus, localStatErr)
				return
			}
			// server is not after local
			if !os.IsNotExist(localStatErr) && reportStatus.Job().Sync && remoteStat.Size() == localStat.Size() {
				// Local version is the same or newer
				reportStatus.Job().UpdateStatus(status.Skipped, reportStatus, nil)
				return
			}
		}

		if reportStatus.dryRun {
			reportStatus.Job().UpdateStatus(status.Complete, reportStatus, nil)
			return
		}

		tmpName := tmpDownloadPath(reportStatus.LocalPath())
		config := reportStatus.Job().Config.(files_sdk.Config)
		config.LogPath(
			reportStatus.RemotePath(),
			map[string]interface{}{
				"LocalTempPath": tmpName,
			},
		)
		writer := openFile(tmpName, reportStatus)
		downloadParts := (&DownloadParts{}).Init(
			reportStatus.fsFile,
			remoteStat,
			reportStatus.Job().Manager.FilePartsManager,
			writer,
			config,
		)

		lib.AnyError(func(err error) {
			reportStatus.Job().UpdateStatus(status.Errored, reportStatus, err)
		},
			func() error { return downloadParts.Run(ctx) },
			func() error { return downloadParts.CloseError },
		)

		if reportStatus.Status().Is(status.Valid...) {
			reportStatus.SetFinalSize(downloadParts.FinalSize())
			config := reportStatus.Job().Config.(files_sdk.Config)
			config.LogPath(
				reportStatus.RemotePath(),
				map[string]interface{}{
					"LocalTempPath": tmpName,
					"FinalSize":     downloadParts.FinalSize(),
				},
			)
			err := finalizeTmpDownload(tmpName, reportStatus.LocalPath())

			if err == nil && reportStatus.PreserveTimes {
				var t time.Time
				if s.file.ProvidedMtime != nil {
					t = *s.file.ProvidedMtime
				} else if s.file.Mtime != nil {
					t = *s.file.Mtime
				}
				if !t.IsZero() {
					err = os.Chtimes(reportStatus.LocalPath(), t.Local(), t.Local())
				}
			}

			if err != nil {
				reportStatus.Job().UpdateStatus(status.Errored, reportStatus, err)
			} else if reportStatus.Status().Is(status.Downloading) {
				reportStatus.Job().UpdateStatus(status.Complete, reportStatus, nil)
			}
		} else {
			err := os.Remove(tmpName) // Clean up on invalid download
			if err != nil {
				reportStatus.Job().UpdateStatus(status.Errored, reportStatus, err)
			}
		}
	}(ctx, s)
}

func openFile(partName string, reportStatus *DownloadStatus) lib.ProgressWriter {
	out, createErr := os.Create(partName)
	if createErr != nil {
		reportStatus.Job().UpdateStatus(status.Errored, reportStatus, createErr)
	}
	writer := lib.ProgressWriter{WriterAndAt: out}
	writer.ProgressWatcher = func(incDownloadedBytes int64) {
		reportStatus.Job().UpdateStatus(status.Downloading, reportStatus, nil)
		reportStatus.incrementDownloadedBytes(incDownloadedBytes)
	}
	return writer
}

func localPath(file files_sdk.File, job status.Job) string {
	var path string
	if job.Type == directory.File {
		path = job.LocalPath
	} else {
		path = filepath.Join(normalizePath(job.LocalPath), relativePath(job, file))
	}

	return path
}

func relativePath(job status.Job, file files_sdk.File) string {
	relativePath, err := filepath.Rel(job.RemotePath, file.Path)
	if err != nil {
		panic(err)
	}
	return relativePath
}
