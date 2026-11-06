package file

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/Files-com/files-sdk-go/v3/file/status"
)

// ErrJobPaused is the cancel cause set when a job is paused rather than fully canceled.
// Use errors.Is(context.Cause(ctx), ErrJobPaused) to distinguish pause from cancel.
var ErrJobPaused = errors.New("job paused")

type JobDownloadCheckpoint struct {
	CompletedPaths []string
	// Bindings optionally identify paused HTTP prefixes. Older checkpoints omit
	// them and remain readable; their HTTP downloads restart from byte zero.
	Bindings []DownloadBinding `json:"bindings,omitempty"`
}

// DownloadBinding ties a paused HTTP prefix to its strong response validator
// and exact source/destination. It contains no signed URL or credentials. Keep
// it together with the paused file and return it in PriorJobCheckpoint.
type DownloadBinding struct {
	Source     string `json:"source"`
	RemotePath string `json:"remote_path"`
	LocalPath  string `json:"local_path"`
	TmpPath    string `json:"tmp_path"`
	ETag       string `json:"etag"`
	Total      int64  `json:"total"`
	Prefix     int64  `json:"prefix"`
}

func downloadSource(f *File) string {
	u, err := url.Parse(f.Config.Endpoint())
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s://%s%s#workspace=%d", u.Scheme, u.Host, u.EscapedPath(), f.Config.WorkspaceId)
}

func checkpointLocalPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return ""
	}
	return filepath.Clean(abs)
}

func (d *DownloadStatus) downloadBinding() *DownloadBinding {
	if !d.Status().Is(status.Canceled) || errors.Is(d.Err(), errDownloadPauseProof) {
		return nil
	}
	tmp := explicitTmpDownload(d.tmpDownloadPath())
	paused, ok := pausedTmpDownloadOf(tmp)
	if !ok || paused.String() != tmp.String() {
		return nil
	}
	return d.bindingForStage(tmp, d.tmpIdentity)
}

// bindingForStage is called only after the engine and output close have
// settled. Status publication happens later, once the artifact is durable.
func (d *DownloadStatus) bindingForStage(tmp destinationPath, expected os.FileInfo) *DownloadBinding {
	f, ok := d.fsFile.(*File)
	if !ok {
		return nil
	}
	info, err := tmp.lstat()
	if err != nil || !info.Mode().IsRegular() || expected == nil || !os.SameFile(info, expected) {
		return nil
	}
	f.fileMutex.Lock()
	representation, path := f.representation, f.File.Path
	f.fileMutex.Unlock()
	etag, total, whole := representation.snapshot()
	if etag == "" || whole || info.Size() > total || downloadSource(f) == "" {
		return nil
	}
	return &DownloadBinding{Source: downloadSource(f), RemotePath: path, LocalPath: checkpointLocalPath(d.LocalPath()), TmpPath: checkpointLocalPath(tmp.String()), ETag: etag, Total: total, Prefix: info.Size()}
}

func (d *DownloadStatus) restoreDownloadBinding(tmp destinationPath, prefix int64) bool {
	f, ok := d.fsFile.(*File)
	if !ok {
		return true
	}
	f.fileMutex.Lock()
	representation := f.representation
	f.fileMutex.Unlock()
	etag, total, _ := representation.snapshot()
	if etag != "" && prefix <= total {
		return true
	}
	params := d.job.Params.(DownloaderParams)
	if params.PriorJobCheckpoint != nil {
		for _, b := range params.PriorJobCheckpoint.Bindings {
			if d.matchesDownloadBinding(b, tmp, prefix) {
				f.fileMutex.Lock()
				f.representation = &downloadRepresentation{etag: b.ETag, total: b.Total}
				f.fileMutex.Unlock()
				return true
			}
		}
	}
	f.fileMutex.Lock()
	f.representation = &downloadRepresentation{whole: true}
	f.fileMutex.Unlock()
	return false
}

func (d *DownloadStatus) matchesDownloadBinding(b DownloadBinding, tmp destinationPath, prefix int64) bool {
	f, ok := d.fsFile.(*File)
	if !ok {
		return false
	}
	f.fileMutex.Lock()
	path := f.File.Path
	f.fileMutex.Unlock()
	return b.Source != "" && b.Source == downloadSource(f) && b.RemotePath == path && b.LocalPath == checkpointLocalPath(d.LocalPath()) && b.TmpPath == checkpointLocalPath(tmp.String()) && b.Prefix == prefix && b.Total >= prefix && strongDownloadETag(b.ETag) == b.ETag && b.ETag != ""
}

// Bound artifacts have a disjoint name. Legacy discovery/admission never
// interprets a directory as a prefix; only this matching owner restores it.
func (d *DownloadStatus) boundTmpDownloadToResume() (*boundPausedDownload, error) {
	_, ok := d.fsFile.(*File)
	if !ok {
		return nil, nil
	}
	var candidates []destinationPath
	if explicit, ok := d.checkpointedTmpDownloadPath(); ok {
		if _, bound := boundPausedDownloadContainer(explicit); bound {
			candidates = append(candidates, explicit)
		}
	}
	if d.job != nil {
		naming, err := tmpDownloadNamingFor(d.destination, d.tempPath)
		if err != nil {
			return nil, err
		}
		canonical := naming.candidate(tempDownloadPausedToken)
		_, digest, _, _ := encodedTmpDownloadParts(canonical.base())
		dir := canonical.parent()
		for _, name := range d.job.pausedListing(dir)[digest] {
			if isBoundPausedDownloadElement(name) {
				candidates = append(candidates, boundPausedDataIn(dir.withName(filepath.Join(dir.name, name))))
			}
		}
	}
	seen := map[string]bool{}
	var selected *boundPausedDownload
	var selectedAt time.Time
	for _, tmp := range candidates {
		key := checkpointLocalPath(tmp.String())
		if seen[key] {
			continue
		}
		seen[key] = true
		a, err := readBoundPausedDownload(tmp)
		if err != nil {
			return nil, err
		}
		if a != nil && d.matchesDownloadBinding(a.binding, tmp, a.dataIdentity.Size()) && (selected == nil || a.dataIdentity.ModTime().After(selectedAt)) {
			selected, selectedAt = a, a.dataIdentity.ModTime()
		}
	}
	return selected, nil
}

func (d *DownloadStatus) retainPausedCustodyError(err error) {
	var custody pausedDownloadCustodyError
	if errors.As(err, &custody) && custody.terminal {
		d.indexErr = err // a retry must not bless an observed replacement
	}
}

func (d *DownloadStatus) parkSettledDownload(tmp destinationPath) (destinationPath, error) {
	if binding := d.bindingForStage(tmp, d.tmpIdentity); binding != nil {
		return parkBoundTmpDownload(tmp, d.tmpIdentity, *binding)
	}
	return parkTmpDownload(tmp, d.tmpIdentity)
}

func (d *DownloadStatus) recordPublicationPause(ctx context.Context, tmp destinationPath, err error) error {
	if paused, ok := pausedPathFromError(err); ok {
		d.setTmpDownloadPath(paused.String())
		return err
	}
	if !errors.Is(context.Cause(ctx), ErrJobPaused) || !stoppedForPauseOnly(err) {
		return err
	}
	paused, ok := pausedTmpDownloadOf(tmp)
	if !ok {
		return err
	}
	if verifyErr := verifyPausedIdentity(paused, d.tmpIdentity); verifyErr != nil {
		return errors.Join(err, pauseProofError(verifyErr))
	}
	kept, proofErr := d.parkSettledDownload(paused)
	if kept.name != "" {
		d.setTmpDownloadPath(kept.String())
	}
	return errors.Join(err, proofErr)
}

// DownloadCheckpoint builds a JobDownloadCheckpoint from the job's settled file statuses.
// Call after Wait returns or Finished is signaled; Canceled only indicates that cancellation was requested.
func (j *Job) DownloadCheckpoint() *JobDownloadCheckpoint {
	completed := make(map[string]struct{})
	var bindings []DownloadBinding
	for p := range j.CompletedPaths {
		completed[p] = struct{}{}
	}

	j.statusesMutex.RLock()
	for _, f := range j.Statuses {
		if d, ok := f.(*DownloadStatus); ok {
			if b := d.downloadBinding(); b != nil {
				bindings = append(bindings, *b)
			}
		}
		if f.Status().Is(status.Complete) {
			completed[f.LocalPath()] = struct{}{}
		}
	}
	j.statusesMutex.RUnlock()

	paths := make([]string, 0, len(completed))
	for p := range completed {
		paths = append(paths, p)
	}
	return &JobDownloadCheckpoint{CompletedPaths: paths, Bindings: bindings}
}
