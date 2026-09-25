package file

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"testing/fstest"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// rangedFS is a remote file system for tests whose one file is served from a
// deterministic byte pattern, by range the way a parallel download reads it or
// as one stream. One range, or one position of the stream, can be held back
// until the test releases it, so a range at a higher offset reaches the disk
// before a lower one, which is what a network does whenever parts finish out of
// order, and so a transfer can be paused while bytes are still outstanding.
type rangedFS struct {
	fstest.MapFS
	name string
	size int64

	holdRange  int64 // the range starting here is held back; -1 for none
	holdStream int64 // a stream read reaching this offset is held back; 0 for none
	holdErr    error // what a held read reports once its context ends, instead of the context's error
	released   chan struct{}
	releaseFn  sync.Once
	holding    chan struct{} // closed once a reader is waiting
	holdingFn  sync.Once

	mu        sync.Mutex
	requested []int64 // range start offsets, in request order
}

func newRangedFS(name string, size int64) *rangedFS {
	r := &rangedFS{MapFS: fstest.MapFS{}, name: name, size: size, holdRange: -1, released: make(chan struct{}), holding: make(chan struct{})}
	r.MapFS[name] = &fstest.MapFile{Mode: fs.ModePerm, Sys: r.sdkFile()}
	return r
}

func (r *rangedFS) sdkFile() files_sdk.File {
	return files_sdk.File{DisplayName: filepath.Base(r.name), Path: r.name, Type: "file", Size: r.size}
}

// content is the byte the file holds at offset i: distinct enough that a range
// written at the wrong place, or left unwritten, changes the file.
func (r *rangedFS) content(i int64) byte { return byte(i*7 + i>>12) }

func (r *rangedFS) expected() []byte {
	b := make([]byte, r.size)
	for i := range b {
		b[i] = r.content(int64(i))
	}
	return b
}

func (r *rangedFS) release() { r.releaseFn.Do(func() { close(r.released) }) }

func (r *rangedFS) nowHolding() { r.holdingFn.Do(func() { close(r.holding) }) }

func (r *rangedFS) Open(name string) (fs.File, error) {
	if name != r.name {
		return r.MapFS.Open(name)
	}
	return &rangedFile{fs: r, ctx: context.Background()}, nil
}

func (r *rangedFS) requestedOffsets() []int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]int64(nil), r.requested...)
}

// waitHeld blocks until the held reader may go on or ctx ends.
func (r *rangedFS) waitHeld(ctx context.Context) error {
	r.nowHolding()
	select {
	case <-r.released:
		return nil
	case <-ctx.Done():
		if r.holdErr != nil {
			return r.holdErr
		}
		return ctx.Err()
	}
}

type rangedFile struct {
	fs  *rangedFS
	ctx context.Context
	pos int64
}

var (
	_ ReaderRange           = (*rangedFile)(nil)
	_ lib.FileWithContext   = (*rangedFile)(nil)
	_ downloadV2URIProvider = (*rangedFile)(nil)
)

func (f *rangedFile) Stat() (fs.FileInfo, error) {
	return Info{File: f.fs.sdkFile(), sizeTrust: TrustedSizeValue}, nil
}

// Read serves the file as one stream, holding at holdStream once.
func (f *rangedFile) Read(p []byte) (int, error) {
	if f.pos >= f.fs.size {
		return 0, io.EOF
	}
	if f.fs.holdStream > 0 && f.pos == f.fs.holdStream {
		if err := f.fs.waitHeld(f.ctx); err != nil {
			return 0, err
		}
	}
	n := 0
	for n < len(p) && f.pos < f.fs.size && !(f.fs.holdStream > 0 && f.pos == f.fs.holdStream && n > 0) {
		p[n] = f.fs.content(f.pos)
		n++
		f.pos++
	}
	return n, nil
}

func (f *rangedFile) Close() error { return nil }

func (f *rangedFile) WithContext(ctx context.Context) fs.File {
	return &rangedFile{fs: f.fs, ctx: ctx}
}

func (f *rangedFile) downloadV2URI(context.Context) (string, error) {
	return "https://ranged.test/" + f.fs.name, nil
}

func (f *rangedFile) ReaderRange(off int64, end int64) (io.ReadCloser, error) {
	f.fs.mu.Lock()
	f.fs.requested = append(f.fs.requested, off)
	f.fs.mu.Unlock()
	reader := &rangeReader{fs: f.fs, ctx: f.ctx, pos: off, end: min(end, f.fs.size-1)}
	if off == f.fs.holdRange {
		reader.held = true
	}
	return reader, nil
}

// rangeReader serves one range. A held range waits for the release or for its
// context to end, like a connection that is still waiting for its first byte
// when the transfer is paused.
type rangeReader struct {
	fs   *rangedFS
	ctx  context.Context
	pos  int64
	end  int64
	held bool
}

func (r *rangeReader) Read(p []byte) (int, error) {
	if r.held {
		if err := r.fs.waitHeld(r.ctx); err != nil {
			return 0, err
		}
		r.held = false
	}
	if r.pos > r.end {
		return 0, io.EOF
	}
	n := 0
	for n < len(p) && r.pos <= r.end {
		p[n] = r.fs.content(r.pos)
		n++
		r.pos++
	}
	return n, nil
}

func (r *rangeReader) Close() error { return nil }

// pauseResumeCase drives one download to a pause while remote bytes are still
// outstanding, then resumes it the way a caller does, with the temporary file
// named in the ended file event the paused job emitted, and checks the
// delivered bytes.
type pauseResumeCase struct {
	remote *rangedFS
	params func(*rangedFS) DownloaderParams
	// ready reports when the job has reached the state to pause in.
	ready func(job *Job, parts *manager.Manager) bool
	// leavesInFlight is set when the scenario leaves a file under the
	// canonical in-flight name that is not this transfer's to remove.
	leavesInFlight bool
}

// pause runs the download until ready, pauses it and returns the job and the
// TmpPath carried by the file's ended event, which is all a caller such as the
// desktop helper keeps for the resume.
func (c pauseResumeCase) pause(t *testing.T, local string) (paused *Job, eventTmpPath string) {
	t.Helper()
	parts := manager.New(1, 3, 1)
	ctx, cancel := context.WithCancelCause(context.Background())
	params := c.params(c.remote)
	params.LocalPath = local
	params.Manager = parts
	paused = downloader(ctx, c.remote, params)
	var mu sync.Mutex
	var ended []JobFile
	paused.RegisterFileEvent(func(f JobFile) {
		mu.Lock()
		defer mu.Unlock()
		ended = append(ended, f)
	}, status.Ended...)
	t.Cleanup(func() {
		// A failed assertion must not leave the job's workers blocked.
		cancel(context.Canceled)
		c.remote.release()
		paused.Wait()
	})
	paused.Start()
	require.Eventually(t, func() bool { return c.ready(paused, parts) }, 30*time.Second, 5*time.Millisecond, "waiting for the download to reach the pause point")
	cancel(ErrJobPaused)
	paused.Wait()
	c.remote.release()
	mu.Lock()
	defer mu.Unlock()
	require.Len(t, ended, 1, "the file ends exactly once")
	return paused, ended[0].TmpPath
}

func (c pauseResumeCase) resume(t *testing.T, local string, tmpPath string) *rangedFS {
	t.Helper()
	resumed := newRangedFS(c.remote.name, c.remote.size)
	params := c.params(resumed)
	params.LocalPath = local
	params.ResumeTmpPath = tmpPath
	params.Manager = manager.New(1, 3, 1)
	job := downloader(context.Background(), resumed, params)
	job.Start()
	job.Wait()
	require.NoError(t, job.Statuses[0].Err())
	require.Equal(t, 1, job.Count(status.Complete))
	requireExactBytes(t, local, c.remote.expected())
	_, pausedLeft := pausedTmpDownloadPath(explicitDestination(local), params.TempPath)
	assert.False(t, pausedLeft, "nothing is left paused after a completed download")
	if !c.leavesInFlight {
		_, inFlight := existingTmpDownloadPath(explicitDestination(local), params.TempPath)
		assert.False(t, inFlight, "nothing is left in flight after a completed download")
	}
	return resumed
}

// parkStage pauses a temporary download a test wrote, as an orderly pause of
// the transfer that wrote it would.
func parkStage(t *testing.T, tmp destinationPath) destinationPath {
	t.Helper()
	identity, err := tmp.lstat()
	require.NoError(t, err)
	paused, err := parkTmpDownload(tmp, identity)
	require.NoError(t, err)
	return paused
}

func requireExactBytes(t *testing.T, local string, expected []byte) {
	t.Helper()
	delivered, err := os.ReadFile(local)
	require.NoError(t, err)
	if !bytes.Equal(delivered, expected) {
		first := -1
		for i := range expected {
			if i >= len(delivered) || delivered[i] != expected[i] {
				first = i
				break
			}
		}
		t.Fatalf("delivered file differs from the source: %d bytes delivered, first mismatch at byte %d", len(delivered), first)
	}
}

func requireNothingStaged(t *testing.T, local string, tempPath string) {
	t.Helper()
	_, inFlight := existingTmpDownloadPath(explicitDestination(local), tempPath)
	_, paused := pausedTmpDownloadPath(explicitDestination(local), tempPath)
	assert.False(t, inFlight, "nothing is left in flight after a completed download")
	assert.False(t, paused, "nothing is left paused after a completed download")
}

// heldRangeCase pauses a download of three ranges while the second is still
// outstanding and the third has already been written.
func heldRangeCase(part int64, params func(*rangedFS) DownloaderParams) pauseResumeCase {
	remote := newRangedFS("big.bin", 3*part-1234)
	remote.holdRange = part
	return pauseResumeCase{
		remote: remote,
		params: params,
		ready: func(job *Job, parts *manager.Manager) bool {
			// Every range has been requested and only the held one is running.
			return len(remote.requestedOffsets()) == 3 && parts.FilePartsManager.RunningCount() <= 1 && job.TransferBytes() >= 2*part-1234
		},
	}
}

func basicParams(remote *rangedFS) DownloaderParams {
	return DownloaderParams{RemotePath: remote.name, config: files_sdk.Config{}.Init()}
}

// adaptiveParams selects the adaptive (V2) engine for any ranged file.
func adaptiveParams(remote *rangedFS) DownloaderParams {
	params := basicParams(remote)
	params.AdaptiveConcurrency = true
	params.AdaptiveDownloadV2TargetClassifier = func(string) TransferV2TargetClass { return downloadV2TargetDefault }
	return params
}

// requirePausedAt checks the state a pause leaves behind: the file named in the
// ended event exists under its paused name, and nothing is left in flight.
func requirePausedAt(t *testing.T, local string, tempPath string, eventTmpPath string) {
	t.Helper()
	require.NotEmpty(t, eventTmpPath, "the ended event names the paused file")
	_, err := os.Stat(eventTmpPath)
	require.NoError(t, err, "the file the ended event names exists")
	paused, ok := pausedTmpDownloadPath(explicitDestination(local), tempPath)
	require.True(t, ok, "the pause leaves a paused temporary download")
	assert.Equal(t, paused.String(), eventTmpPath, "the ended event names the paused file, not the in-flight one")
	_, inFlight := existingTmpDownloadPath(explicitDestination(local), tempPath)
	assert.False(t, inFlight, "nothing is left under the in-flight name")
}

// A parallel download pauses while a range in the middle of the file is still
// outstanding and a later range has already been written. Resuming must
// deliver the file's exact bytes: the temporary file may only be continued
// from bytes that were actually received, not from wherever the furthest
// range happened to end. The bytes received before the held range are kept.
func TestDownloadPauseResume_parallelRangesOutOfOrder(t *testing.T) {
	// Files over 10 MiB are fetched by the parts engine in 10 MiB ranges.
	const part int64 = 10 * 1024 * 1024
	c := heldRangeCase(part, basicParams)
	local := filepath.Join(t.TempDir(), "big.bin")

	paused, tmpPath := c.pause(t, local)
	require.Equal(t, 1, paused.Count(status.Canceled), "the pause leaves the file canceled, not errored")
	requirePausedAt(t, local, "", tmpPath)

	resumed := c.resume(t, local, tmpPath)
	offsets := resumed.requestedOffsets()
	require.NotEmpty(t, offsets, "the missing middle range must be fetched")
	assert.Equal(t, part, offsets[0], "the resume continues from where the received prefix ends")
}

// The adaptive engine fetches larger ranges the same way.
func TestDownloadPauseResume_adaptiveRangesOutOfOrder(t *testing.T) {
	// Files under 1 GiB are fetched by the adaptive engine in 16 MiB ranges.
	const part int64 = 16 * 1024 * 1024
	c := heldRangeCase(part, adaptiveParams)
	local := filepath.Join(t.TempDir(), "big.bin")

	paused, tmpPath := c.pause(t, local)
	require.Equal(t, 1, paused.Count(status.Canceled))
	requirePausedAt(t, local, "", tmpPath)

	resumed := c.resume(t, local, tmpPath)
	offsets := resumed.requestedOffsets()
	require.NotEmpty(t, offsets)
	assert.Equal(t, part, offsets[0], "the resume continues from where the received prefix ends")
}

// A file small enough to be fetched as one stream is written in order, so the
// bytes received before a pause are all valid and the resume continues after
// them rather than starting over.
func TestDownloadPauseResume_singleStreamKeepsReceivedBytes(t *testing.T) {
	const held int64 = 1024*1024 + 12345
	remote := newRangedFS("small.bin", 4*1024*1024-77)
	remote.holdStream = held
	c := pauseResumeCase{
		remote: remote,
		params: basicParams,
		ready: func(job *Job, _ *manager.Manager) bool {
			select {
			case <-remote.holding:
				return job.TransferBytes() == held
			default:
				return false
			}
		},
	}
	local := filepath.Join(t.TempDir(), "small.bin")

	paused, tmpPath := c.pause(t, local)
	require.Equal(t, 1, paused.Count(status.Canceled))
	requirePausedAt(t, local, "", tmpPath)

	resumed := c.resume(t, local, tmpPath)
	offsets := resumed.requestedOffsets()
	require.NotEmpty(t, offsets, "the resume asks for the rest of the file by range")
	assert.Equal(t, held, offsets[0], "the bytes received before the pause are kept")
}

// A second pause, after a resume, keeps working the same way: the resumed
// transfer's own pause is what the next resume continues from.
func TestDownloadPauseResume_pausedTwice(t *testing.T) {
	const part int64 = 10 * 1024 * 1024
	local := filepath.Join(t.TempDir(), "big.bin")
	first := heldRangeCase(part, basicParams)
	_, tmpPath := first.pause(t, local)
	requirePausedAt(t, local, "", tmpPath)

	// Resume, and hold the range after the next one, so the resumed transfer
	// pauses with one more range received and one written past a gap.
	second := heldRangeCase(part, basicParams)
	second.remote.holdRange = 2 * part
	second.ready = func(job *Job, parts *manager.Manager) bool {
		return len(second.remote.requestedOffsets()) == 2 && parts.FilePartsManager.RunningCount() <= 1 && job.TransferBytes() >= 2*part
	}
	ctxParams := second.params
	second.params = func(r *rangedFS) DownloaderParams {
		p := ctxParams(r)
		p.ResumeTmpPath = tmpPath
		return p
	}
	paused, tmpPath2 := second.pause(t, local)
	require.Equal(t, 1, paused.Count(status.Canceled))
	requirePausedAt(t, local, "", tmpPath2)
	assert.Equal(t, []int64{part, 2 * part}, second.remote.requestedOffsets(), "the resumed transfer fetched from the first pause point")

	resumed := second.resume(t, local, tmpPath2)
	require.NotEmpty(t, resumed.requestedOffsets())
	assert.Equal(t, 2*part, resumed.requestedOffsets()[0], "the final resume continues from the second pause point")
}

// With an external temp directory, the paused file lives there under a name
// tied to the destination's whole path, and a resume without a checkpoint
// still finds it.
func TestDownloadPauseResume_externalTempPath(t *testing.T) {
	const part int64 = 10 * 1024 * 1024
	temp := t.TempDir()
	withTemp := func(r *rangedFS) DownloaderParams {
		p := basicParams(r)
		p.TempPath = temp
		return p
	}
	c := heldRangeCase(part, withTemp)
	local := filepath.Join(t.TempDir(), "big.bin")

	_, tmpPath := c.pause(t, local)
	requirePausedAt(t, local, temp, tmpPath)
	assert.Equal(t, temp, filepath.Dir(tmpDownloadEntry(explicitTmpDownload(tmpPath)).String()), "the paused file is inside the temp directory")

	resumed := c.resume(t, local, "") // no checkpoint: found by its canonical paused name
	offsets := resumed.requestedOffsets()
	require.NotEmpty(t, offsets)
	assert.Equal(t, part, offsets[0])
	entries, err := os.ReadDir(temp)
	require.NoError(t, err)
	assert.Empty(t, entries, "the temp directory is left clean")
}

// A temporary download still under its in-flight name was left by a process
// that stopped while it was writing, or by an earlier version of this client.
// Its length says nothing about which bytes were received, so it is not
// continued from and not published: the download starts over and delivers the
// exact bytes, whether the file is short or already at the full size. The
// stale file is removed only when the caller's checkpoint names it; a file
// under the canonical in-flight name that no checkpoint names may belong to
// another attempt and is left alone.
func TestDownloadResume_inFlightStageIsNotTrusted(t *testing.T) {
	const part int64 = 10 * 1024 * 1024
	for name, tc := range map[string]struct {
		leave        func(t *testing.T, c pauseResumeCase, local string) (checkpoint string)
		staleRemains bool
	}{
		// A paused download was resumed, and the process died while the
		// resumed transfer was writing; the caller still holds the paused
		// name from the pause.
		"activated after a pause": {
			leave: func(t *testing.T, c pauseResumeCase, local string) string {
				_, tmpPath := c.pause(t, local)
				paused, ok := pausedTmpDownloadPath(explicitDestination(local), "")
				require.True(t, ok)
				_, err := activateTmpDownload(paused, explicitDestination(local), "")
				require.NoError(t, err)
				return tmpPath
			},
			staleRemains: true,
		},
		// The whole length is there, as after the adaptive engine preallocated
		// the file, but with a gap that was never received; the caller
		// checkpointed the in-flight name from a progress event before the
		// process died.
		"full length with a gap, checkpointed": {
			leave: func(t *testing.T, c pauseResumeCase, local string) string {
				tmp, err := tmpDownloadPath(explicitDestination(local), "")
				require.NoError(t, err)
				wrong := c.remote.expected()
				for i := part; i < 2*part; i++ {
					wrong[i] = 0
				}
				require.NoError(t, os.WriteFile(tmp.String(), wrong, 0o644))
				return tmp.String()
			},
			staleRemains: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			c := heldRangeCase(part, basicParams)
			local := filepath.Join(t.TempDir(), "big.bin")
			checkpoint := tc.leave(t, c, local)
			stale, inFlight := existingTmpDownloadPath(explicitDestination(local), "")
			require.True(t, inFlight, "the scenario leaves a file under the in-flight name")
			c.leavesInFlight = tc.staleRemains

			resumed := c.resume(t, local, checkpoint)
			offsets := resumed.requestedOffsets()
			require.NotEmpty(t, offsets, "the content must be fetched")
			assert.Equal(t, int64(0), offsets[0], "the download starts over")
			_, err := stale.lstat()
			if tc.staleRemains {
				assert.NoError(t, err, "a stale file no checkpoint names is left where it was")
			} else {
				assert.ErrorIs(t, err, fs.ErrNotExist, "the checkpointed stale file is discarded")
			}
		})
	}
}

// When the temporary file cannot be trimmed to its received bytes, keeping it
// would let the next attempt continue from bytes that never arrived. The pause
// then reports the failure and does not keep the file; the next attempt starts
// over and still delivers the exact bytes.
func TestDownloadPause_stageThatCannotBeTrimmedIsNotKept(t *testing.T) {
	for name, tc := range map[string]struct {
		part   int64
		params func(*rangedFS) DownloaderParams
	}{
		"parts engine":    {part: 10 * 1024 * 1024, params: basicParams},
		"adaptive engine": {part: 16 * 1024 * 1024, params: adaptiveParams},
	} {
		t.Run(name, func(t *testing.T) {
			trimFailure := errors.New("disk refused the truncate")
			previous := truncateTmpDownload
			truncateTmpDownload = func(interface{ Truncate(int64) error }, int64) error { return trimFailure }
			t.Cleanup(func() { truncateTmpDownload = previous })

			c := heldRangeCase(tc.part, tc.params)
			local := filepath.Join(t.TempDir(), "big.bin")
			paused, tmpPath := c.pause(t, local)
			require.Equal(t, 1, paused.Count(status.Errored), "the pause reports that the file could not be kept")
			require.ErrorIs(t, paused.Statuses[0].Err(), errTmpDownloadNotTrimmed)
			require.ErrorIs(t, paused.Statuses[0].Err(), trimFailure)
			_, err := os.Stat(tmpPath)
			require.ErrorIs(t, err, fs.ErrNotExist, "the untrimmed temporary file is not kept for a resume")
			requireNothingStaged(t, local, "")

			truncateTmpDownload = previous
			resumed := c.resume(t, local, tmpPath)
			offsets := resumed.requestedOffsets()
			require.NotEmpty(t, offsets)
			assert.Equal(t, int64(0), offsets[0], "the next attempt starts over")
		})
	}
}

// Another transfer can put its own file at the temporary path while this one
// runs, as a rename through a second name for the file would do. At a pause
// that file must not become the checkpoint a later run continues from: it is
// removed, as it is when a completed download is published, and the next
// attempt starts over with the exact bytes.
func TestDownloadPause_replacedStageIsNotKeptAsCheckpoint(t *testing.T) {
	const part int64 = 10 * 1024 * 1024
	c := heldRangeCase(part, basicParams)
	local := filepath.Join(t.TempDir(), "big.bin")
	ready := c.ready
	c.ready = func(job *Job, parts *manager.Manager) bool {
		if !ready(job, parts) {
			return false
		}
		require.NoError(t, replaceStage(explicitDestination(local), []byte("theirs")))
		return true
	}

	paused, tmpPath := c.pause(t, local)
	require.Equal(t, 1, paused.Count(status.Errored), "the pause reports that its file could not be kept")
	require.ErrorIs(t, paused.Statuses[0].Err(), errTempDownloadReplaced)
	_, kept := pausedTmpDownloadPath(explicitDestination(local), "")
	assert.False(t, kept, "a replaced temporary file is not kept as a checkpoint")
	_, inFlight := existingTmpDownloadPath(explicitDestination(local), "")
	assert.False(t, inFlight, "the replacement is removed, as at publication")

	resumed := c.resume(t, local, tmpPath)
	require.NotEmpty(t, resumed.requestedOffsets())
	assert.Equal(t, int64(0), resumed.requestedOffsets()[0], "the next attempt starts over")
}

// A checkpoint that lexically lies below the caller-selected directory but
// leads outside it through a link is refused by the lookup; it must not be
// deleted through an unconfined path either.
func TestDownloadResume_checkpointLeavingTheRootIsNeitherUsedNorDeleted(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	require.NoError(t, os.Symlink(outside, filepath.Join(root, "link")))
	remote := newRangedFS("small.bin", 4*1024*1024-77)
	local := filepath.Join(root, "small.bin")
	// A file spelled like this transfer's in-flight temporary download, but
	// reachable only through the link.
	stage, err := tmpDownloadPath(explicitDestination(filepath.Join(outside, "small.bin")), "")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(stage.String(), []byte("sentinel that must survive"), 0o644))
	checkpoint := filepath.Join(root, "link", strings.TrimPrefix(stage.String(), outside+string(filepath.Separator)))

	params := basicParams(remote)
	params.LocalPath = local
	params.ResumeTmpPath = checkpoint
	job := downloader(context.Background(), remote, params)
	job.Start()
	job.Wait()

	require.NoError(t, job.Statuses[0].Err())
	requireExactBytes(t, local, remote.expected())
	content, err := os.ReadFile(stage.String())
	require.NoError(t, err)
	assert.Equal(t, "sentinel that must survive", string(content), "nothing outside the root is touched")
}

// While another file occupies the canonical in-flight name (another attempt,
// or what a stopped process left), a transfer stages under a uniqueness token
// and pauses under a paused name of its own. The next run continues from it
// whether it has the checkpoint or, like the CLI, only the destination, and
// whether or not the occupant is still there; the occupant is left alone.
func TestDownloadPauseResume_besideAnotherAttemptsStage(t *testing.T) {
	const part int64 = 10 * 1024 * 1024
	for name, tc := range map[string]struct {
		withCheckpoint  bool
		occupantRemoved bool
	}{
		"resumed through the checkpoint":                {withCheckpoint: true},
		"resumed by discovery":                          {},
		"resumed by discovery, occupant gone":           {occupantRemoved: true},
		"resumed through the checkpoint, occupant gone": {withCheckpoint: true, occupantRemoved: true},
	} {
		t.Run(name, func(t *testing.T) {
			c := heldRangeCase(part, basicParams)
			c.leavesInFlight = !tc.occupantRemoved
			local := filepath.Join(t.TempDir(), "big.bin")
			other, err := tmpDownloadPath(explicitDestination(local), "")
			require.NoError(t, err)
			require.NoError(t, os.WriteFile(other.String(), []byte("another attempt is writing this"), 0o644))

			paused, tmpPath := c.pause(t, local)
			require.Equal(t, 1, paused.Count(status.Canceled))
			require.NotEmpty(t, tmpPath)
			_, err = os.Stat(tmpPath)
			require.NoError(t, err, "the ended event names the paused file")
			assert.True(t, isPausedTmpDownloadElement(tmpDownloadEntry(explicitTmpDownload(tmpPath)).base()))
			_, canonical := pausedTmpDownloadPath(explicitDestination(local), "")
			assert.False(t, canonical, "the pause keeps its own tokenized name, not the canonical one")
			assert.NotEqual(t, other.String(), tmpPath)
			content, err := os.ReadFile(other.String())
			require.NoError(t, err)
			assert.Equal(t, "another attempt is writing this", string(content), "the other attempt's file is untouched")
			if tc.occupantRemoved {
				require.NoError(t, os.Remove(other.String()))
			}

			checkpoint := tmpPath
			if !tc.withCheckpoint {
				checkpoint = ""
			}
			resumed := c.resume(t, local, checkpoint)
			require.NotEmpty(t, resumed.requestedOffsets())
			assert.Equal(t, part, resumed.requestedOffsets()[0], "the resume continues from the received prefix")
			if !tc.occupantRemoved {
				require.NoError(t, os.Remove(other.String()))
			}
		})
	}
}

// A pause that arrives while the source has rejected the transfer's download
// request must not turn the bytes of that request into a checkpoint.
func TestDownloadPause_rejectedTransferIsNotKept(t *testing.T) {
	const part int64 = 10 * 1024 * 1024
	c := heldRangeCase(part, basicParams)
	c.remote.holdErr = files_sdk.ErrDownloadSourceChanged
	local := filepath.Join(t.TempDir(), "big.bin")

	paused, tmpPath := c.pause(t, local)
	require.Equal(t, 0, paused.Count(status.Complete))
	requireNothingStaged(t, local, "")
	if tmpPath != "" {
		_, err := os.Stat(tmpPath)
		assert.ErrorIs(t, err, fs.ErrNotExist, "nothing named by the ended event is kept")
	}

	resumed := c.resume(t, local, tmpPath)
	require.NotEmpty(t, resumed.requestedOffsets())
	assert.Equal(t, int64(0), resumed.requestedOffsets()[0], "the next attempt starts over")
}

// Only a stop caused by the cancellation alone lets a paused file be kept.
// The cancellation may be wrapped the way transports report it, but anything
// else joined with it, such as a failed write or close, means the file's
// state is not simply "stopped here".
func Test_stoppedForPauseOnly(t *testing.T) {
	ioFailure := errors.New("write /tmp/x: no space left on device")
	for name, tc := range map[string]struct {
		err  error
		want bool
	}{
		"the pause cause":                   {ErrJobPaused, true},
		"a plain cancellation":              {context.Canceled, true},
		"a wrapped cancellation":            {fmt.Errorf("range 2: %w", context.Canceled), true},
		"a transport-wrapped one":           {&url.Error{Op: "Get", URL: "https://x", Err: context.Canceled}, true},
		"two cancellations joined":          {errors.Join(context.Canceled, fmt.Errorf("part: %w", ErrJobPaused)), true},
		"nothing":                           {nil, false},
		"a cancellation plus an i/o error":  {errors.Join(context.Canceled, ioFailure), false},
		"a cancellation plus a close error": {errors.Join(ErrJobPaused, &fs.PathError{Op: "close", Path: "x", Err: ioFailure}), false},
		"a rejected source":                 {errors.Join(context.Canceled, files_sdk.ErrDownloadSourceChanged), false},
		"a size conflict":                   {fmt.Errorf("%w: part 2", errDownloadSizeConflict), false},
		"an unrelated error":                {ioFailure, false},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, stoppedForPauseOnly(tc.err))
		})
	}
}
