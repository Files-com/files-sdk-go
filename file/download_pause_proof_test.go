package file

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/stretchr/testify/require"
)

func TestPausedDownloadProofPreservesBytesAndEnforcesEntireRecordLimit(t *testing.T) {
	binding := DownloadBinding{Source: "https://source.test#workspace=1", RemotePath: "remote.bin", LocalPath: "local\xff.bin", TmpPath: "stage/DATA.BIN", ETag: "\"opaque\xff\"", Total: 1000, Prefix: 123}
	encoded, err := encodePausedDownloadProof(binding)
	require.NoError(t, err)
	restored, err := decodePausedDownloadProof(encoded)
	require.NoError(t, err)
	require.Equal(t, binding, restored)

	// The producer and reader share the selected ceiling for the entire record.
	remaining := maxPausedDownloadProofBytes - len(encoded)
	binding.ETag = binding.ETag[:len(binding.ETag)-1] + strings.Repeat("x", remaining) + `"`
	encoded, err = encodePausedDownloadProof(binding)
	require.NoError(t, err)
	require.Len(t, encoded, maxPausedDownloadProofBytes)
	restored, err = decodePausedDownloadProof(encoded)
	require.NoError(t, err)
	require.True(t, bytes.Equal([]byte(binding.ETag), []byte(restored.ETag)))
	binding.ETag = binding.ETag[:len(binding.ETag)-1] + `x"`
	_, err = encodePausedDownloadProof(binding)
	require.Error(t, err)
	_, err = decodePausedDownloadProof(append(encoded, 0))
	require.Error(t, err)
}

type proofTestStage struct {
	source     *versionDownloadFixture
	client     *Client
	final, tmp destinationPath
	identity   fs.FileInfo
	binding    DownloadBinding
}

func newProofTestStage(t *testing.T, complete bool) proofTestStage {
	t.Helper()
	source := &versionDownloadFixture{data: bytes.Repeat([]byte("same version bytes;"), 4096)}
	server := httptest.NewServer(http.HandlerFunc(source.serve))
	t.Cleanup(server.Close)
	source.url = server.URL
	client := source.client()
	final := explicitDestination(filepath.Join(t.TempDir(), "output.bin"))
	tmp, err := tmpDownloadPath(final, "")
	require.NoError(t, err)
	prefix := 2048
	if complete {
		prefix = len(source.data)
	}
	identity := writeStage(t, tmp, source.data[:prefix])
	binding := DownloadBinding{Source: downloadSource(&File{FS: &FS{Config: client.Config}}), RemotePath: "remote.bin", LocalPath: checkpointLocalPath(final.String()), TmpPath: checkpointLocalPath(tmp.String()), ETag: `"original"`, Total: int64(len(source.data)), Prefix: int64(prefix)}
	return proofTestStage{source: source, client: client, final: final, tmp: tmp, identity: identity, binding: binding}
}

func (s proofTestStage) download(t *testing.T) *Job {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	job := s.client.Downloader(DownloaderParams{RemotePath: "remote.bin", LocalPath: s.final.String(), Manager: manager.Build(1, 1, true)}, files_sdk.WithContext(ctx))
	job.Start()
	waitForDownloadTestJob(t, ctx, job)
	require.Len(t, job.Statuses, 1)
	return job
}

func TestPausedDownloadProofAutomaticallyRestoresFreshClientsAndCompletePrefixes(t *testing.T) {
	for _, complete := range []bool{false, true} {
		t.Run(fmt.Sprint(complete), func(t *testing.T) {
			s := newProofTestStage(t, complete)
			paused, err := parkBoundTmpDownload(s.tmp, s.identity, s.binding)
			require.NoError(t, err)
			raw, err := os.ReadFile(paused.String())
			require.NoError(t, err)
			require.True(t, bytes.Equal(s.source.data[:s.binding.Prefix], raw), "proof bytes never enter the payload")
			a, err := readBoundPausedDownload(paused)
			require.NoError(t, err)
			require.NotNil(t, a)
			token, _ := tmpDownloadToken(a.container.base())
			require.False(t, strings.HasPrefix(token, tempDownloadPausedToken), "old p- discovery ignores the new reserved container")
			require.True(t, isReservedTempDownloadPath(paused.String()))
			job := s.download(t) // no explicit checkpoint or ResumeTmpPath
			require.NoError(t, job.Statuses[0].Err())
			requireExactBytes(t, s.final.String(), s.source.data)
			s.source.mu.Lock()
			requests := append([]versionDownloadRequest(nil), s.source.requests...)
			s.source.mu.Unlock()
			if complete {
				require.Empty(t, requests, "a complete bound prefix publishes without fetching")
			} else {
				require.Len(t, requests, 1)
				require.Equal(t, fmt.Sprintf("bytes=%d-%d", s.binding.Prefix, s.binding.Total-1), requests[0].ranges)
				require.Equal(t, `"original"`, requests[0].match)
			}
			entries, err := os.ReadDir(s.final.dir)
			require.NoError(t, err)
			require.Len(t, entries, 1, "publication removes only the owned payload/proof container")
		})
	}
}

type faultingPausedProofFile struct {
	pausedProofFile
	writeErr, syncErr, closeErr error
	closed                      *atomic.Int32
}

func (f faultingPausedProofFile) Write(p []byte) (int, error) {
	if f.writeErr != nil {
		return 0, f.writeErr
	}
	return f.pausedProofFile.Write(p)
}
func (f faultingPausedProofFile) Sync() error {
	if f.syncErr != nil {
		return f.syncErr
	}
	return f.pausedProofFile.Sync()
}
func (f faultingPausedProofFile) Close() error {
	if f.closed != nil {
		f.closed.Add(1)
	}
	return errors.Join(f.pausedProofFile.Close(), f.closeErr)
}

func TestPausedDownloadProofFailuresKeepVerifiedLegacyBytesAndRestartWhole(t *testing.T) {
	cause := errors.New("synthetic proof IO failure")
	for _, kind := range []string{"write", "sync", "close", "oversize"} {
		t.Run(kind, func(t *testing.T) {
			s := newProofTestStage(t, false)
			original := openPausedProofFile
			defer func() { openPausedProofFile = original }()
			var opened atomic.Int32
			openPausedProofFile = func(root *os.Root, name string, flags int, mode fs.FileMode) (pausedProofFile, error) {
				file, err := original(root, name, flags, mode)
				if err != nil || flags&os.O_WRONLY == 0 {
					return file, err
				}
				opened.Add(1)
				fault := faultingPausedProofFile{pausedProofFile: file}
				switch kind {
				case "write":
					fault.writeErr = cause
				case "sync":
					fault.syncErr = cause
				case "close":
					fault.closeErr = cause
				}
				return fault, nil
			}
			if kind == "oversize" {
				s.binding.ETag = `"` + strings.Repeat("v", maxPausedDownloadProofBytes) + `"`
			}
			paused, err := parkBoundTmpDownload(s.tmp, s.identity, s.binding)
			require.ErrorIs(t, err, errDownloadPauseProof)
			if kind == "oversize" {
				require.Zero(t, opened.Load(), "producer checks the entire bound before metadata allocation/publication")
			} else {
				require.ErrorIs(t, err, cause)
			}
			require.NotEmpty(t, paused.name, "settled owned bytes are kept only after verified legacy parking")
			require.False(t, isBoundPausedDownloadElement(tmpDownloadEntry(paused).base()))
			raw, err := os.ReadFile(paused.String())
			require.NoError(t, err)
			require.True(t, bytes.Equal(s.source.data[:s.binding.Prefix], raw))
			job := s.download(t)
			require.NoError(t, job.Statuses[0].Err())
			requireExactBytes(t, s.final.String(), s.source.data)
			s.source.mu.Lock()
			defer s.source.mu.Unlock()
			require.Len(t, s.source.requests, 1)
			require.Empty(t, s.source.requests[0].ranges, "a no-proof legacy prefix restarts from byte zero")
		})
	}
}

func TestPublicPauseProofCloseFailureSurvivesCancellationAndReporter(t *testing.T) {
	data := bytes.Repeat([]byte("pause-safe bytes;"), 1<<15)
	var endpoint string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"path": "remote.bin", "display_name": "remote.bin", "type": "file", "size": len(data), "download_uri": endpoint + "/content"})
			return
		}
		w.Header().Set("ETag", `"version"`)
		w.Header().Set("Content-Length", fmt.Sprint(len(data)))
		_, _ = w.Write(data[:64<<10])
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	}))
	defer server.Close()
	endpoint = server.URL
	testCtx, stop := context.WithTimeout(context.Background(), 30*time.Second)
	defer stop()
	ctx, pause := context.WithCancelCause(testCtx)
	defer pause(nil)
	cause := fmt.Errorf("closing proof: %w", context.Canceled)
	var closes atomic.Int32
	original := openPausedProofFile
	defer func() { openPausedProofFile = original }()
	openPausedProofFile = func(root *os.Root, name string, flags int, mode fs.FileMode) (pausedProofFile, error) {
		file, err := original(root, name, flags, mode)
		if err != nil || flags&os.O_WRONLY == 0 {
			return file, err
		}
		return faultingPausedProofFile{pausedProofFile: file, closeErr: cause, closed: &closes}, nil
	}
	destination := filepath.Join(t.TempDir(), "output.bin")
	client := &Client{Config: files_sdk.Config{EndpointOverride: endpoint, DisableDirectTransfers: true}.Init().SetCustomClient(server.Client())}
	job := client.Downloader(DownloaderParams{RemotePath: "remote.bin", LocalPath: destination, Manager: manager.Build(1, 1, true)}, files_sdk.WithContext(ctx))
	var once sync.Once
	job.RegisterFileEvent(func(file JobFile) {
		if file.TransferBytes > 0 {
			once.Do(func() { pause(ErrJobPaused) })
		}
	}, OnBytesChange(status.Downloading))
	var reported error
	job.RegisterFileEvent(func(file JobFile) { reported = file.Err }, status.Ended...)
	job.Start()
	waitForDownloadTestJob(t, testCtx, job)
	require.Len(t, job.Statuses, 1)
	require.Equal(t, status.Errored, job.Statuses[0].Status())
	require.ErrorIs(t, job.Statuses[0].Err(), cause)
	require.ErrorIs(t, reported, cause)
	require.ErrorIs(t, reported, errDownloadPauseProof)
	require.EqualValues(t, 1, closes.Load(), "the original metadata close is attempted once")
	require.Empty(t, job.DownloadCheckpoint().Bindings, "failed proof publication is never exported as a trusted pause")
	ended := ToStatusFile(job.Statuses[0])
	raw, err := os.ReadFile(ended.TmpPath)
	require.NoError(t, err)
	require.NotEmpty(t, raw)
	require.LessOrEqual(t, len(raw), 64<<10)
	require.True(t, bytes.Equal(data[:len(raw)], raw))
	require.NoFileExists(t, destination)
}

func TestPausedDownloadProofReadCloseFailureDoesNotRestartOrRenew(t *testing.T) {
	s := newProofTestStage(t, false)
	paused, err := parkBoundTmpDownload(s.tmp, s.identity, s.binding)
	require.NoError(t, err)
	original := openPausedProofFile
	defer func() { openPausedProofFile = original }()
	cause := fmt.Errorf("closing read proof: %w", context.Canceled)
	openPausedProofFile = func(root *os.Root, name string, flags int, mode fs.FileMode) (pausedProofFile, error) {
		file, err := original(root, name, flags, mode)
		if err != nil || flags != os.O_RDONLY {
			return file, err
		}
		return faultingPausedProofFile{pausedProofFile: file, closeErr: cause}, nil
	}
	job := s.download(t)
	require.ErrorIs(t, job.Statuses[0].Err(), cause)
	require.ErrorIs(t, job.Statuses[0].Err(), errDownloadPauseProof)
	s.source.mu.Lock()
	defer s.source.mu.Unlock()
	require.Empty(t, s.source.requests)
	raw, err := os.ReadFile(paused.String())
	require.NoError(t, err)
	require.True(t, bytes.Equal(s.source.data[:s.binding.Prefix], raw))
}

func TestPausedDownloadProofContextMismatchLeavesForeignStageUntouched(t *testing.T) {
	s := newProofTestStage(t, false)
	s.binding.Source = "https://another-source.test#workspace=1"
	paused, err := parkBoundTmpDownload(s.tmp, s.identity, s.binding)
	require.NoError(t, err)
	proof := filepath.Join(filepath.Dir(paused.String()), boundPausedDownloadRecord)
	before, err := os.ReadFile(proof)
	require.NoError(t, err)
	job := s.download(t)
	require.NoError(t, job.Statuses[0].Err())
	requireExactBytes(t, s.final.String(), s.source.data)
	after, err := os.ReadFile(proof)
	require.NoError(t, err)
	require.True(t, bytes.Equal(before, after), "another source's record is not consumed")
	raw, err := os.ReadFile(paused.String())
	require.NoError(t, err)
	require.True(t, bytes.Equal(s.source.data[:s.binding.Prefix], raw))
	s.source.mu.Lock()
	defer s.source.mu.Unlock()
	require.Len(t, s.source.requests, 1)
	require.Empty(t, s.source.requests[0].ranges)
}

func TestMalformedAndOversizedPausedProofsCannotValidateAStoredPrefix(t *testing.T) {
	for _, oversized := range []bool{false, true} {
		t.Run(fmt.Sprint(oversized), func(t *testing.T) {
			s := newProofTestStage(t, false)
			paused, err := parkBoundTmpDownload(s.tmp, s.identity, s.binding)
			require.NoError(t, err)
			proof := filepath.Join(filepath.Dir(paused.String()), boundPausedDownloadRecord)
			invalid := []byte("unknown or interrupted proof")
			if oversized {
				invalid = bytes.Repeat([]byte("x"), maxPausedDownloadProofBytes+1)
			}
			require.NoError(t, os.WriteFile(proof, invalid, 0600))
			job := s.download(t)
			require.NoError(t, job.Statuses[0].Err())
			requireExactBytes(t, s.final.String(), s.source.data)
			kept, err := os.ReadFile(proof)
			require.NoError(t, err)
			require.True(t, bytes.Equal(invalid, kept), "unverified artifacts are not removed as owned state")
			s.source.mu.Lock()
			defer s.source.mu.Unlock()
			require.Len(t, s.source.requests, 1)
			require.Empty(t, s.source.requests[0].ranges)
		})
	}
}

func TestPausedDownloadProofRefusesForeignChildrenBeforeClaim(t *testing.T) {
	s := newProofTestStage(t, true)
	paused, err := parkBoundTmpDownload(s.tmp, s.identity, s.binding)
	require.NoError(t, err)
	foreign := filepath.Join(filepath.Dir(paused.String()), "FOREIGN")
	content := []byte("another owner's file")
	require.NoError(t, os.WriteFile(foreign, content, 0600))
	err = finalizeTmpDownload(context.Background(), paused, s.final, s.identity)
	require.ErrorIs(t, err, errDownloadPauseProof)
	kept, err := os.ReadFile(foreign)
	require.NoError(t, err)
	require.True(t, bytes.Equal(content, kept))
	requireExactBytes(t, paused.String(), s.source.data)
	require.NoFileExists(t, s.final.String())
}

func TestPausedDownloadProofRollbackRefusesTheOccupiedOldName(t *testing.T) {
	s := newProofTestStage(t, false)
	original := openPausedProofFile
	defer func() { openPausedProofFile = original }()
	cause := errors.New("metadata write failure after old-name replacement")
	foreign := []byte("another owner's bytes")
	openPausedProofFile = func(root *os.Root, name string, flags int, mode fs.FileMode) (pausedProofFile, error) {
		file, err := original(root, name, flags, mode)
		if err != nil || flags&os.O_WRONLY == 0 {
			return file, err
		}
		require.NoError(t, os.WriteFile(s.tmp.String(), foreign, 0600))
		return faultingPausedProofFile{pausedProofFile: file, writeErr: cause}, nil
	}
	kept, err := parkBoundTmpDownload(s.tmp, s.identity, s.binding)
	require.ErrorIs(t, err, cause)
	require.ErrorIs(t, err, errDownloadPauseProof)
	require.NotEmpty(t, kept.name)
	require.NotEqual(t, s.tmp.String(), kept.String())
	actual, err := os.ReadFile(s.tmp.String())
	require.NoError(t, err)
	require.True(t, bytes.Equal(foreign, actual), "rollback never overwrites the newly occupied old name")
	raw, err := os.ReadFile(kept.String())
	require.NoError(t, err)
	require.True(t, bytes.Equal(s.source.data[:s.binding.Prefix], raw))
}

func TestCompleteDownloadPausePublishesBoundProofBeforeReporting(t *testing.T) {
	for _, alreadyBound := range []bool{false, true} {
		t.Run(fmt.Sprint(alreadyBound), func(t *testing.T) {
			s := newProofTestStage(t, true)
			tmp := s.tmp
			var err error
			if alreadyBound {
				tmp, err = parkBoundTmpDownload(s.tmp, s.identity, s.binding)
				require.NoError(t, err)
			}
			ctx, cancel := context.WithCancelCause(context.Background())
			cancel(ErrJobPaused)
			err = finalizeTmpDownload(ctx, tmp, s.final, s.identity, &s.binding)
			require.ErrorIs(t, err, ErrJobPaused)
			require.NotErrorIs(t, err, errDownloadPauseProof)
			paused, ok := pausedPathFromError(err)
			require.True(t, ok)
			a, err := readBoundPausedDownload(paused)
			require.NoError(t, err)
			require.NotNil(t, a)
			require.Equal(t, a.binding.Total, a.binding.Prefix)
			require.NoFileExists(t, s.final.String())
			job := s.download(t)
			require.NoError(t, job.Statuses[0].Err())
			requireExactBytes(t, s.final.String(), s.source.data)
			s.source.mu.Lock()
			defer s.source.mu.Unlock()
			require.Empty(t, s.source.requests, "complete-prefix pause needs no replay on the next run")
		})
	}
}

func TestCanceledCompleteDownloadDoesNotPublish(t *testing.T) {
	for _, alreadyBound := range []bool{false, true} {
		t.Run(fmt.Sprint(alreadyBound), func(t *testing.T) {
			s := newProofTestStage(t, true)
			tmp := s.tmp
			var err error
			if alreadyBound {
				tmp, err = parkBoundTmpDownload(tmp, s.identity, s.binding)
				require.NoError(t, err)
			}
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			err = finalizeTmpDownload(ctx, tmp, s.final, s.identity, &s.binding)
			require.ErrorIs(t, err, context.Canceled)
			require.NotErrorIs(t, err, errDownloadPauseProof, "pure cancellation retains its original status semantics")
			require.NoFileExists(t, s.final.String())
			_, bound := pausedPathFromError(err)
			require.False(t, bound)
			require.NoError(t, removeTmpDownload(tmp))
			entries, err := os.ReadDir(s.final.dir)
			require.NoError(t, err)
			require.Empty(t, entries)
		})
	}
}

type replacingRecoveredProofFile struct {
	pausedProofFile
	afterClose func() error
}

func (f replacingRecoveredProofFile) Close() error {
	if err := f.pausedProofFile.Close(); err != nil {
		return err
	}
	return f.afterClose()
}

func TestRecoveredPausedStageReplacementIsNeverAdopted(t *testing.T) {
	for _, scenario := range []struct {
		name                   string
		complete, rewriteProof bool
	}{
		{name: "partial DATA replacement"},
		{name: "complete DATA replacement", complete: true},
		{name: "in-place BINDING rewrite", rewriteProof: true},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			s := newProofTestStage(t, scenario.complete)
			paused, err := parkBoundTmpDownload(s.tmp, s.identity, s.binding)
			require.NoError(t, err)
			before, err := readBoundPausedDownload(paused)
			require.NoError(t, err)
			require.NotNil(t, before)
			foreign := bytes.Repeat([]byte("X"), int(s.binding.Prefix))
			expectedKept := foreign
			if scenario.rewriteProof {
				expectedKept = s.source.data[:s.binding.Prefix]
			}
			original := openPausedProofFile
			defer func() { openPausedProofFile = original }()
			var once sync.Once
			var replaced atomic.Bool
			var unchangedIDs atomic.Bool
			var replacementErr error
			openPausedProofFile = func(root *os.Root, name string, flags int, mode fs.FileMode) (pausedProofFile, error) {
				file, err := original(root, name, flags, mode)
				if err != nil || flags != os.O_RDONLY {
					return file, err
				}
				return replacingRecoveredProofFile{pausedProofFile: file, afterClose: func() error {
					once.Do(func() {
						if scenario.rewriteProof {
							encoded, err := os.ReadFile(filepath.Join(filepath.Dir(paused.String()), boundPausedDownloadRecord))
							if err != nil {
								replacementErr = err
								return
							}
							binding, err := decodePausedDownloadProof(encoded)
							if err != nil {
								replacementErr = err
								return
							}
							binding.ETag = `"changed!"` // same length; all inode identities remain unchanged
							encoded, err = encodePausedDownloadProof(binding)
							if err != nil {
								replacementErr = err
								return
							}
							if err := root.WriteFile(boundPausedDownloadRecord, encoded, 0600); err != nil {
								replacementErr = err
								return
							}
							dataInfo, dataErr := root.Lstat(boundPausedDownloadData)
							proofInfo, proofErr := root.Lstat(boundPausedDownloadRecord)
							containerInfo, containerErr := root.Stat(".")
							replacementErr = errors.Join(dataErr, proofErr, containerErr)
							unchangedIDs.Store(replacementErr == nil && os.SameFile(before.dataIdentity, dataInfo) && os.SameFile(before.proofIdentity, proofInfo) && os.SameFile(before.containerIdentity, containerInfo))
							replaced.Store(replacementErr == nil)
							return
						}
						// Replace the same-sized DATA after its recovery identity was
						// captured, before the controller can observe/claim it.
						if err := root.WriteFile("NEW.BIN", foreign, 0600); err != nil {
							replacementErr = err
							return
						}
						replacementErr = root.Rename("NEW.BIN", boundPausedDownloadData)
						replaced.Store(replacementErr == nil)
					})
					return replacementErr
				}}, nil
			}
			job := s.download(t)
			require.True(t, replaced.Load(), "the real recovery/claim boundary must have been exercised")
			if scenario.rewriteProof {
				require.True(t, unchangedIDs.Load(), "the rewrite must bypass inode replacement checks to exercise decoded-binding custody")
			}
			require.ErrorIs(t, job.Statuses[0].Err(), errTempDownloadReplaced)
			require.ErrorIs(t, job.Statuses[0].Err(), errDownloadPauseProof)
			require.NoFileExists(t, s.final.String(), "neither complete foreign bytes nor a foreign prefix may publish")
			s.source.mu.Lock()
			requests := append([]versionDownloadRequest(nil), s.source.requests...)
			s.source.mu.Unlock()
			require.Empty(t, requests, "custody refusal precedes any continuation")
			_, err = os.Lstat(filepath.Dir(paused.String()))
			require.ErrorIs(t, err, fs.ErrNotExist, "the observed-invalid b1 cannot be blessed by a later retry")
			entries, err := os.ReadDir(s.final.dir)
			require.NoError(t, err)
			require.Len(t, entries, 1)
			require.False(t, isBoundPausedDownloadElement(entries[0].Name()))
			claim := filepath.Join(s.final.dir, entries[0].Name(), boundPausedDownloadData)
			require.True(t, isReservedTempDownloadPath(claim))
			kept, err := os.ReadFile(claim)
			require.NoError(t, err)
			require.True(t, bytes.Equal(expectedKept, kept), "quarantine preserves all payload bytes")
		})
	}
}

func TestPausedDownloadProofRejectsUnknownTruncatedAndOversizedFields(t *testing.T) {
	binding := DownloadBinding{Source: "source", RemotePath: "remote", LocalPath: "local", TmpPath: "tmp", ETag: `"version"`, Total: 100, Prefix: 20}
	encoded, err := encodePausedDownloadProof(binding)
	require.NoError(t, err)
	for _, invalid := range [][]byte{encoded[:len(encoded)-1], append(append([]byte(nil), encoded...), 0), []byte("unknown version")} {
		_, err := decodePausedDownloadProof(invalid)
		require.Error(t, err)
	}
	invalid := append([]byte(nil), encoded...)
	binary.BigEndian.PutUint32(invalid[len(pausedDownloadProofMagic):], ^uint32(0))
	_, err = decodePausedDownloadProof(invalid)
	require.Error(t, err, "declared field lengths cannot allocate past the bounded record")
}
