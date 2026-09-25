package file

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/directory"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/Files-com/files-sdk-go/v3/ignore"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/lib/direction"
	"github.com/Files-com/files-sdk-go/v3/lib/keyvalue"
	gitignore "github.com/sabhiram/go-gitignore"
)

func downloader(ctx context.Context, fileSys fs.FS, params DownloaderParams, opts ...files_sdk.RequestResponseOption) *Job {
	job := (&Job{}).Init()
	SetJobParams(job, direction.DownloadType, params, params.config.Logger, fileSys)
	job.Config = params.config
	if params.PriorJobCheckpoint != nil {
		job.CompletedPaths = make(map[string]struct{}, len(params.PriorJobCheckpoint.CompletedPaths))
		for _, p := range params.PriorJobCheckpoint.CompletedPaths {
			job.CompletedPaths[p] = struct{}{}
		}
	}
	jobCtx := job.withDirectTransferClientCache(ctx)
	remoteFs, ok := fileSys.(lib.FSWithContext)
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
		job.RemotePath = lib.NewUrlPath(params.RemotePath).PruneStartingSlash().String()
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
		} else if err == nil {
			if stats.IsDir() {
				localType = directory.Dir
			} else {
				localType = directory.File
			}
		} else {
			// Propagating this error is difficult, but this error will happen again in CodeStart.
		}
		if (!lib.NewUrlPath(params.RemotePath).EndingSlash() && localType == directory.Dir) || remoteType == directory.File && localType == directory.Dir {
			job.LocalPath = filepath.Join(job.LocalPath, lib.NewUrlPath(job.RemotePath).SwitchPathSeparator(string(os.PathSeparator)).Pop())
			if remoteType == directory.File {
				localType = directory.File
			}
		}

		// Use relative path
		if job.LocalPath == "" {
			job.LocalPath = lib.NewUrlPath(job.RemotePath).SwitchPathSeparator(string(os.PathSeparator)).Pop()
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
		WaitTellFinished(job, onComplete, func() {
			RetryByPolicy(jobCtx, job, job.RetryPolicy.(RetryPolicy), false)
			runSyncAfterEmptyFolders(job, params.SyncAfterActions, params.DryRun, params.config, opts...)
		})

		// ignore.New only returns an error if run on an unsupported OS
		job.Ignore, _ = ignore.New(params.Ignore...)
		if len(params.Include) > 0 {
			job.Include, _ = ignore.New(params.Include...)
		}

		it := (&lib.Walk[lib.DirEntry]{
			FS:                 fileSys,
			Root:               lib.UrlJoinNoEscape(job.RemotePath),
			ConcurrencyManager: job.Manager.FilePartsManager,
			WalkFile:           lib.DirEntryWalkFile,
			ListDirectories:    true,
		}).Walk(jobCtx)

		for it.Next() {
			if it.Resource().Err() != nil {
				createIndexedStatus(Entity{error: it.Resource().Err()}, params, job)
			} else {
				f, err := fileSys.Open(it.Resource().Path())
				createIndexedStatus(Entity{error: err, File: f, FS: fileSys}, params, job)
			}
		}

		if it.Err() != nil {
			metaFile := &DownloadStatus{
				job:         job,
				status:      status.Errored,
				localPath:   params.LocalPath,
				remotePath:  params.RemotePath,
				tempPath:    params.TempPath,
				Sync:        params.Sync,
				NoOverwrite: params.NoOverwrite,
				Mutex:       &sync.RWMutex{},
			}
			metaFile.file = files_sdk.File{
				DisplayName: filepath.Base(params.LocalPath),
				Type:        job.Direction.Name(),
				Path:        params.RemotePath,
			}
			job.Add(metaFile)
			job.UpdateStatus(status.Errored, metaFile, it.Err())
			onComplete <- metaFile
		}

		job.EndScan()
	}

	return job
}

func enqueueIndexedDownloads(job *Job, jobCtx context.Context, onComplete chan *DownloadStatus) {
	params, _ := job.Params.(DownloaderParams)
	if zipBatchDisabled(params, job) {
		enqueueIndexedDownloadsDirect(job, jobCtx, onComplete)
		return
	}

	batcher := newZipBatchDownloader(job, jobCtx, params, onComplete)
	for !job.EndScanning.Called() || job.Count(status.Indexed) > 0 {
		if jobCtx.Err() != nil {
			// The job is canceled or paused: drain the remaining queue in one
			// batch instead of one EnqueueNext scan per file. Files already
			// offered to the batcher are finalized by flushEnd below.
			drainCanceledIndexed(job, onComplete)
			break
		}
		if f, ok := job.EnqueueNext(); ok {
			downloadStatus := f.(*DownloadStatus)
			if batcher.offer(downloadStatus) {
				continue
			}
			enqueueIndexedDownloadDirect(job, jobCtx, downloadStatus, onComplete)
		}
	}
	batcher.flushEnd()
}

func enqueueIndexedDownloadsDirect(job *Job, jobCtx context.Context, onComplete chan *DownloadStatus) {
	for !job.EndScanning.Called() || job.Count(status.Indexed) > 0 {
		if jobCtx.Err() != nil {
			// The job is canceled or paused: drain the remaining queue in one
			// batch instead of one EnqueueNext scan per file. In-flight files
			// keep their normal cancellation path via enqueueDownload.
			drainCanceledIndexed(job, onComplete)
			return
		}
		if f, ok := job.EnqueueNext(); ok {
			enqueueIndexedDownloadDirect(job, jobCtx, f.(*DownloadStatus), onComplete)
		}
	}
}

func enqueueIndexedDownloadDirect(job *Job, jobCtx context.Context, downloadStatus *DownloadStatus, onComplete chan *DownloadStatus) {
	params, _ := job.Params.(DownloaderParams)
	adaptiveAdmission := params.AdaptiveConcurrency && (params.AdaptiveConcurrencyUseSDKDefaultCaps || params.Manager == nil)
	adaptiveTargets := downloadV2JobAdmissionTargets{job: job}
	fileAdmissionManager := job.fileAdmissionManager()
	gateDownload := shouldGateAdaptiveDownloadAdmission(downloadStatus, adaptiveAdmission)
	admitted := !gateDownload ||
		waitForAdaptiveFileAdmission(jobCtx, fileAdmissionManager, adaptiveTargets, adaptiveFileAdmissionInitialTarget())
	if admitted && fileAdmissionManager.WaitWithContext(jobCtx) {
		go enqueueDownload(jobCtx, job, downloadStatus, onComplete)
	} else {
		job.UpdateStatus(status.Canceled, downloadStatus, nil)
		onComplete <- downloadStatus
	}
}

func normalizePath(rootDestination string) string {
	if rootDestination != "" && rootDestination[len(rootDestination)-1:] == string(os.PathSeparator) {
	} else {
		rootDestination, _ = filepath.Abs(rootDestination)
	}
	return rootDestination
}

func createIndexedStatus(f Entity, params DownloaderParams, job *Job) {
	s := &DownloadStatus{
		error:         f.error,
		fsFile:        f.File,
		FS:            f.FS,
		job:           job,
		localPath:     params.LocalPath,
		remotePath:    params.RemotePath,
		tempPath:      params.TempPath,
		Sync:          params.Sync,
		NoOverwrite:   params.NoOverwrite,
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
			s.remotePath = s.file.Path
			s.TmpPath = params.ResumeTmpPath
			s.destination, err = downloadDestination(job, s.file)
			if err == nil {
				s.localPath = s.destination.String()
			} else {
				s.error, s.indexErr = err, err
			}
		} else {
			s.SetStatus(status.Errored, err)
		}
	} else {
		s.indexErr = f.error
	}

	job.Add(s)
}

func enqueueDownload(ctx context.Context, job *Job, downloadStatus *DownloadStatus, signal chan *DownloadStatus) {
	if err := downloadStatus.indexErr; err != nil || downloadStatus.error != nil || downloadStatus.fsFile == nil {
		if err == nil {
			err = downloadStatus.RecentError()
		}
		job.UpdateStatus(status.Errored, downloadStatus, err)
		job.fileAdmissionManager().Done()
		signal <- downloadStatus
		return
	}
	if _, ok := job.CompletedPaths[downloadStatus.LocalPath()]; ok {
		job.UpdateStatus(status.Skipped, downloadStatus, nil)
		job.fileAdmissionManager().Done()
		signal <- downloadStatus
		return
	}
	if ignoreDownloadJob(job, downloadStatus) {
		job.UpdateStatus(status.Ignored, downloadStatus, nil)
		job.fileAdmissionManager().Done()
		signal <- downloadStatus
		return
	}

	downloadFolderItem(ctx, signal, downloadStatus)
}

func ignoreDownloadJob(job *Job, downloadStatus *DownloadStatus) bool {
	// A temporary download is never transferred, whatever the caller's rules
	// say. Otherwise a remote file could take over the unfinished download of
	// the file it is named after and be delivered in its place. The remote path
	// is checked as well as the local one, because a temporary download folder
	// selected as the folder to download loses its own name on the way to the
	// local destination.
	if isReservedTempDownloadPath(downloadStatus.RemotePath()) {
		return true
	}
	// The local path is also checked through the name the filesystem stores, so
	// a second name for a temporary download does not get past this. A path
	// that exists but cannot be resolved is refused, not allowed.
	reaches, err := localPathReachesTempDownload(downloadStatus.LocalPath())
	if err != nil {
		job.Logger.Printf("not transferring %v: its local path could not be resolved: %v", downloadStatus.LocalPath(), err)
		return true
	}
	if reaches {
		return true
	}
	return ignorePath(downloadStatus.RemotePath(), job.Ignore, job.Include)
}

func ignorePath(path string, ignored, included *gitignore.GitIgnore) bool {
	// if the ignore matches, or the include doesn't match, we skip the file
	if (ignored != nil && ignored.MatchesPath(path)) || (included != nil && !included.MatchesPath(path)) {
		return true
	}
	return false
}

func downloadFolderItem(ctx context.Context, signal chan *DownloadStatus, s *DownloadStatus) {
	defer func() {
		s.job.fileAdmissionManager().Done()
		signal <- s
	}()
	runDownloadFolderItem(ctx, s)
}

func runDownloadFolderItem(ctx context.Context, reportStatus *DownloadStatus) {
	remoteStat, ok := prepareDownloadFolderItem(reportStatus)
	if !ok {
		return
	}

	var startOffset int64
	// Each attempt establishes the temporary file again, so the identity
	// recorded by an earlier attempt must not be carried over.
	reportStatus.tmpIdentity = nil
	// Only a paused temporary download is continued from: one a run ended at
	// the bytes it received and then renamed. Whatever else a run left, under
	// the checkpointed path or the canonical in-flight name, is discarded.
	tmp, resuming := reportStatus.tmpDownloadToResume()
	if !resuming {
		if checkpointed := reportStatus.tmpDownloadPath(); checkpointed != "" {
			reportStatus.Job().Logger.Printf("tmp download file not found or not paused, starting over: %v", checkpointed)
			reportStatus.discardUnpausedCheckpoint()
		}
		reportStatus.setTmpDownloadPath("")
		var err error
		tmp, err = tmpDownloadPath(reportStatus.destination, reportStatus.tempPath)
		if err != nil {
			reportStatus.Job().UpdateStatus(status.Errored, reportStatus, err)
			return
		}
	} else {
		fi, err := tmp.stat()
		if err == nil {
			startOffset = fi.Size()
			reportStatus.tmpIdentity = fi
		}
		if startOffset > remoteStat.Size() {
			startOffset = 0
		}
		if startOffset == remoteStat.Size() {
			// The paused file holds the whole transfer: publish it without a
			// network request. A publication that fails because of a pause
			// leaves it paused for the next run, and the status names it.
			reportStatus.setTmpDownloadPath(tmp.String())
			if err := finalizeTmpDownload(ctx, tmp, reportStatus.destination, reportStatus.tmpIdentity); err != nil {
				if !errors.Is(context.Cause(ctx), ErrJobPaused) {
					removeTmpDownload(tmp)
				}
				reportStatus.Job().UpdateStatus(status.Errored, reportStatus, err)
			} else {
				reportStatus.SetFinalSize(startOffset)
				reportStatus.Job().UpdateStatus(status.Complete, reportStatus, nil)
			}
			return
		}
		// The file goes back under an in-flight name before anything is
		// written to it again: a process that stops before the next pause must
		// not leave it looking continued from.
		tmp, err = activateTmpDownload(tmp, reportStatus.destination, reportStatus.tempPath)
		if err != nil {
			reportStatus.Job().UpdateStatus(status.Errored, reportStatus, err)
			return
		}
	}
	reportStatus.setTmpDownloadPath(tmp.String())
	// Whichever engine opens the temporary file records which file it opened
	// (recordTmpDownloadIdentity), so publishing can tell this transfer's file
	// from any other file that reaches the same path while the transfer runs.
	if startOffset > 0 {
		reportStatus.IncrementTransferBytes(startOffset)
	}
	var finalSize int64
	// engineErr is what stopped the engine early, if anything. It is reported
	// after the temporary file has been dealt with, so the status the caller
	// sees names the file as it is then.
	var engineErr error
	downloadV2Used, downloadV2FinalSize, planStat, downloadV2Err := runDownloadV2IfSupported(ctx, reportStatus, remoteStat, tmp, startOffset)
	if downloadV2Used {
		finalSize = downloadV2FinalSize
		engineErr = downloadV2Err
	} else {
		writer, err := openFile(tmp, reportStatus, startOffset)
		if err != nil {
			reportStatus.Job().UpdateStatus(status.Errored, reportStatus, err)
			return
		}
		downloadParts := (&DownloadParts{}).Init(
			reportStatus.fsFile,
			planStat,
			reportStatus.Job().Manager.FilePartsManager,
			writer,
			reportStatus.Job().Config,
			startOffset,
		)
		engineErr = downloadParts.Run(ctx)
		if engineErr == nil {
			engineErr = downloadParts.CloseError
		} else if downloadParts.CloseError != nil {
			engineErr = errors.Join(engineErr, downloadParts.CloseError)
		}
		finalSize = downloadParts.FinalSize()
	}

	cause := context.Cause(ctx)

	// parked is set once the temporary file has been kept for a resume under
	// its paused name. That happens only when the engine stopped for the
	// pause itself and ended the file at its received prefix; a file the
	// engine could not trim, or one of a transfer the source rejected, is not
	// kept, pause or not, because a resume would continue from its length.
	parked := false
	var parkErr error
	notTrimmed, stageUntrimmed := tmpDownloadNotTrimmed(engineErr)
	if engineErr != nil && errors.Is(cause, ErrJobPaused) && !stageUntrimmed && stoppedForPauseOnly(engineErr) {
		if paused, err := parkTmpDownload(tmp, reportStatus.tmpIdentity); err != nil {
			parkErr = fmt.Errorf("download: the paused download %v could not be kept for a resume: %w", tmp, err)
		} else {
			parked = true
			// The status reported next names the paused file, which is what a
			// caller keeps for the resume.
			reportStatus.setTmpDownloadPath(paused.String())
		}
	}
	// A file that could not be kept is reported instead of the pause, once:
	// that the resume will start over is what the caller must know.
	switch {
	case stageUntrimmed:
		reportStatus.Job().UpdateStatus(status.Errored, reportStatus, notTrimmed)
	case parkErr != nil:
		reportStatus.Job().UpdateStatus(status.Errored, reportStatus, parkErr)
	case engineErr != nil:
		reportStatus.Job().UpdateStatus(status.Errored, reportStatus, engineErr)
	}

	if reportStatus.Status().Is(status.Valid...) {
		err := completeTmpDownload(ctx, reportStatus, tmp, finalSize)
		if err != nil && !errors.Is(context.Cause(ctx), ErrJobPaused) {
			removeTmpDownload(tmp)
		}

		if err != nil {
			reportStatus.Job().UpdateStatus(status.Errored, reportStatus, err)
		} else if reportStatus.Status().Is(status.Downloading) {
			reportStatus.Job().UpdateStatus(status.Complete, reportStatus, nil)
		}
	} else {
		if !parked {
			err := removeTmpDownload(tmp) // Clean up on invalid download
			if err != nil && !errors.Is(err, fs.ErrNotExist) {
				reportStatus.Job().UpdateStatus(status.Errored, reportStatus, err)
			}
			if downloadSourceChanged(reportStatus.Err()) {
				// The rejected download request cannot serve any range. A retry
				// within the file retry budget must start a new lifecycle.
				if remoteFile, ok := reportStatus.fsFile.(*File); ok {
					remoteFile.restartDownload()
				}
			}
		}
	}
}

// stoppedForPauseOnly reports whether err is the engine stopping for the
// cancellation and nothing else: every error in it, however wrapped or joined,
// ends in the cancellation. A transfer the source rejected, a range response
// that contradicted the plan, or an unrelated failure such as a write or close
// error joined with the cancellation says the file's state is not simply
// "stopped here", and such a file is not continued from.
func stoppedForPauseOnly(err error) bool {
	if err == nil {
		return false
	}
	switch unwrapped := err.(type) {
	case interface{ Unwrap() []error }:
		for _, joined := range unwrapped.Unwrap() {
			if !stoppedForPauseOnly(joined) {
				return false
			}
		}
		return true
	case interface{ Unwrap() error }:
		return stoppedForPauseOnly(unwrapped.Unwrap())
	}
	return err == ErrJobPaused || err == context.Canceled
}

func prepareDownloadFolderItem(reportStatus *DownloadStatus) (fs.FileInfo, bool) {
	destination := reportStatus.destination
	if !reportStatus.dryRun {
		if err := destination.parent().mkdirAll(); err != nil {
			reportStatus.Job().UpdateStatus(status.Errored, reportStatus, err)
			return nil, false
		}
	}

	remoteStat, remoteStatErr := reportStatus.fsFile.Stat()
	if remoteStatErr != nil {
		reportStatus.Job().UpdateStatus(status.Errored, reportStatus, remoteStatErr)
		return nil, false
	}

	if reportStatus.NoOverwrite {
		_, localStatErr := destination.stat()
		if localStatErr == nil {
			reportStatus.Job().UpdateStatus(status.FileExists, reportStatus, localStatErr)
			return nil, false
		}
		if !os.IsNotExist(localStatErr) {
			reportStatus.Job().UpdateStatus(status.Errored, reportStatus, localStatErr)
			return nil, false
		}
	}

	if reportStatus.Job().Sync {
		localStat, localStatErr := destination.stat()
		if localStatErr != nil && !os.IsNotExist(localStatErr) {
			reportStatus.Job().UpdateStatus(status.Errored, reportStatus, localStatErr)
			return nil, false
		}
		// server is not after local
		if !os.IsNotExist(localStatErr) && remoteStat.Size() == localStat.Size() {
			// Local version is the same or newer
			reportStatus.Job().UpdateStatus(status.Skipped, reportStatus, nil)
			return nil, false
		}
		reportStatus.Job().UpdateStatus(status.Compared, reportStatus, nil)
	}

	if reportStatus.dryRun {
		reportStatus.Job().UpdateStatus(status.Complete, reportStatus, nil)
		return nil, false
	}

	if reportStatus.File().IsDir() {
		err := destination.mkdirAll()
		if err != nil {
			reportStatus.Job().UpdateStatus(status.Errored, reportStatus, err)
		} else {
			reportStatus.Job().UpdateStatus(status.FolderCreated, reportStatus, nil)
		}
		return nil, false
	}

	return remoteStat, true
}

func completeTmpDownload(ctx context.Context, reportStatus *DownloadStatus, tmp destinationPath, finalSize int64) error {
	reportStatus.SetFinalSize(finalSize)
	if err := finalizeTmpDownload(ctx, tmp, reportStatus.destination, reportStatus.tmpIdentity); err != nil {
		if errors.Is(context.Cause(ctx), ErrJobPaused) {
			// finalizeTmpDownload keeps a file it could not publish because
			// of a pause under its paused name; report that name.
			if paused, ok := pausedTmpDownloadOf(tmp); ok {
				reportStatus.setTmpDownloadPath(paused.String())
			}
		}
		return err
	}

	if reportStatus.PreserveTimes {
		var t time.Time
		if reportStatus.file.ProvidedMtime != nil {
			t = *reportStatus.file.ProvidedMtime
		} else if reportStatus.file.Mtime != nil {
			t = *reportStatus.file.Mtime
		}
		if !t.IsZero() {
			return reportStatus.destination.chtimes(t.Local())
		}
	}

	return nil
}

func openFile(tmp destinationPath, reportStatus *DownloadStatus, startOffset int64) (lib.ProgressWriter, error) {
	var out *os.File
	var err error
	if startOffset > 0 {
		out, err = tmp.openFile(os.O_WRONLY|os.O_CREATE, 0644)
	} else {
		out, err = tmp.create()
	}
	if err != nil {
		return lib.ProgressWriter{}, err
	}
	if err := reportStatus.recordTmpDownloadIdentity(out); err != nil {
		out.Close()
		return lib.ProgressWriter{}, err
	}
	writer := lib.ProgressWriter{WriterAndAt: out}
	writer.ProgressWatcher = func(incDownloadedBytes int64) {
		reportStatus.Job().UpdateStatusWithBytes(status.Downloading, reportStatus, incDownloadedBytes)
	}
	return writer, nil
}
