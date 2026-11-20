//go:build linux || windows

package fsmount

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	iofs "io/fs"
	"net/http"
	path_lib "path"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	"github.com/winfsp/cgofuse/fuse"
)

// offlineStorage is the storage of the offline tests: files by path, and whether
// it can be reached.
type offlineStorage struct {
	mu    sync.Mutex
	files map[string][]byte
	down  bool
	mtime time.Time

	// Per folder: an error instead of its listing, an error after its items, a
	// pause before one of its items, and a gate its next listing waits at.
	listErr   map[string]error
	listBreak map[string]error
	listPause map[string]offlinePause
	listGate  map[string]*offlineGate
	// Per file: a gate its next download waits at, before any bytes.
	downloadGate map[string]*offlineGate
}

type offlinePause struct {
	before int
	pause  time.Duration
}

// offlineGate holds a listing: reached closes when it gets there, and it goes
// on when open is closed.
type offlineGate struct {
	reached chan struct{}
	open    chan struct{}
}

// offlineListIter hands out a listing the way paged storage does: it can wait
// before an item, stopping when its context ends, and fail after the items.
type offlineListIter struct {
	ctx    context.Context
	files  []files_sdk.File
	idx    int
	pause  offlinePause
	broken error
	err    error
}

func (it *offlineListIter) Next() bool {
	if it.pause.pause > 0 && it.idx == it.pause.before {
		select {
		case <-time.After(it.pause.pause):
		case <-it.ctx.Done():
			it.err = it.ctx.Err()
			return false
		}
	}
	if it.idx < len(it.files) {
		it.idx++
		return true
	}
	it.err = it.broken
	return false
}

func (it *offlineListIter) File() files_sdk.File { return it.files[it.idx-1] }

func (it *offlineListIter) Err() error { return it.err }

var errOfflineTestUnreachable = errors.New("dial tcp 127.0.0.1:9: connect: connection refused")

func newOfflineStorage(files map[string]string) *offlineStorage {
	s := &offlineStorage{
		files:     make(map[string][]byte),
		mtime:     time.Now().Add(-time.Hour).Truncate(time.Second),
		listErr:   make(map[string]error),
		listBreak: make(map[string]error),
		listPause: make(map[string]offlinePause),
		listGate:  make(map[string]*offlineGate),

		downloadGate: make(map[string]*offlineGate),
	}
	for path, content := range files {
		s.files[path] = []byte(content)
	}
	return s
}

func (s *offlineStorage) set(path, content string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.files[path] = []byte(content)
}

func (s *offlineStorage) setDown(down bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.down = down
}

// with changes the storage's settings while it isn't answering.
func (s *offlineStorage) with(change func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	change()
}

func (s *offlineStorage) file(path string, content []byte) files_sdk.File {
	mtime := s.mtime
	return files_sdk.File{Path: path, DisplayName: path_lib.Base(path), Type: "file", Size: int64(len(content)), Mtime: &mtime}
}

func (s *offlineStorage) folder(path string) files_sdk.File {
	mtime := s.mtime
	return files_sdk.File{Path: path, DisplayName: path_lib.Base(path), Type: "directory", Mtime: &mtime}
}

func (s *offlineStorage) find(path string) (files_sdk.File, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.down {
		return files_sdk.File{}, errOfflineTestUnreachable
	}
	if content, ok := s.files[path]; ok {
		return s.file(path, content), nil
	}
	for name := range s.files {
		if path == "/" || strings.HasPrefix(name, path+"/") {
			return s.folder(path), nil
		}
	}
	return files_sdk.File{}, files_sdk.ResponseError{Type: string(files_sdk.ErrFileNotFound)}
}

func (s *offlineStorage) list(ctx context.Context, path string) (remoteFileIter, error) {
	s.mu.Lock()
	if gate := s.listGate[path]; gate != nil {
		delete(s.listGate, path)
		s.mu.Unlock()
		close(gate.reached)
		<-gate.open
		s.mu.Lock()
	}
	defer s.mu.Unlock()
	if s.down {
		return nil, errOfflineTestUnreachable
	}
	if err := s.listErr[path]; err != nil {
		return nil, err
	}
	prefix := strings.TrimSuffix(path, "/") + "/"
	children := make(map[string]files_sdk.File)
	for name, content := range s.files {
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		rest := strings.TrimPrefix(name, prefix)
		if child, _, nested := strings.Cut(rest, "/"); nested {
			children[child] = s.folder(prefix + child)
		} else {
			children[rest] = s.file(name, content)
		}
	}
	items := make([]files_sdk.File, 0, len(children))
	for _, item := range children {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].DisplayName < items[j].DisplayName })
	return &offlineListIter{ctx: ctx, files: items, pause: s.listPause[path], broken: s.listBreak[path]}, nil
}

func (s *offlineStorage) move(from, to string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.down {
		return errOfflineTestUnreachable
	}
	moved := make(map[string][]byte)
	for name, content := range s.files {
		if name == from || strings.HasPrefix(name, from+"/") {
			moved[to+strings.TrimPrefix(name, from)] = content
			delete(s.files, name)
		}
	}
	for name, content := range moved {
		s.files[name] = content
	}
	return nil
}

func (s *offlineStorage) content(params files_sdk.FileDownloadParams) ([]byte, error) {
	path := params.Path
	if path == "" {
		path = params.File.Path
	}
	s.mu.Lock()
	if gate := s.downloadGate[path]; gate != nil {
		delete(s.downloadGate, path)
		s.mu.Unlock()
		close(gate.reached)
		<-gate.open
		s.mu.Lock()
	}
	defer s.mu.Unlock()
	if s.down {
		return nil, errOfflineTestUnreachable
	}
	content, ok := s.files[path]
	if !ok {
		return nil, files_sdk.ResponseError{Type: string(files_sdk.ErrFileNotFound)}
	}
	return content, nil
}

func (s *offlineStorage) backend() *fakeRemoteBackend {
	return &fakeRemoteBackend{
		findFunc: func(params files_sdk.FileFindParams, _ ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
			return s.find(params.Path)
		},
		listForFunc: func(params files_sdk.FolderListForParams, opts ...files_sdk.RequestResponseOption) (remoteFileIter, error) {
			return s.list(files_sdk.ContextOption(opts), params.Path)
		},
		moveFunc: func(params files_sdk.FileMoveParams, _ ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error) {
			return files_sdk.FileAction{}, s.move(params.Path, params.Destination)
		},
		downloadRangeFunc: func(params files_sdk.FileDownloadParams, requested cache.ByteRange, _ ...files_sdk.RequestResponseOption) (remoteRangeResponse, error) {
			content, err := s.content(params)
			if err != nil {
				return remoteRangeResponse{}, err
			}
			sum := sha256.Sum256(content)
			start, end := min(requested.Start, int64(len(content))), min(requested.End, int64(len(content)))
			return remoteRangeResponse{
				ETag:      `"` + hex.EncodeToString(sum[:8]) + `"`,
				Returned:  cache.ByteRange{Start: start, End: end},
				TotalSize: int64(len(content)),
				Partial:   true,
				Body:      io.NopCloser(bytes.NewReader(content[start:end])),
			}, nil
		},
		downloadFunc: func(params files_sdk.FileDownloadParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
			content, err := s.content(params)
			if err != nil {
				return params.File, err
			}
			_, err = files_sdk.BuildResponse(&http.Response{StatusCode: http.StatusOK, ContentLength: int64(len(content)), Body: io.NopCloser(bytes.NewReader(content))}, opts...)
			return params.File, err
		},
	}
}

func newOfflineTestFs(t *testing.T, storage *offlineStorage) *RemoteFs {
	t.Helper()
	fs, vfs, _ := newTestRemoteFs(t)
	fs.cacheStore = newTestDiskCache(t)
	fs.backend = storage.backend()
	keeper, err := newOfflineKeeper(fs)
	if err != nil {
		t.Fatal(err)
	}
	fs.offline = keeper
	keeper.start()
	t.Cleanup(func() {
		fs.offline.close()
		fs.closeRangeDownloads()
		vfs.destroy()
	})
	return fs
}

// waitForStates waits until each path is in its state.
func waitForStates(t *testing.T, fs *RemoteFs, want map[string]OfflineState) {
	t.Helper()
	paths := make([]string, 0, len(want))
	for path := range want {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	deadline := time.Now().Add(10 * time.Second)
	for {
		states := fs.OfflineStates(paths)
		mismatch := ""
		for i, path := range paths {
			if states[i] != want[path] {
				mismatch += " " + path + "=" + states[i].String() + " (want " + want[path].String() + ")"
			}
		}
		if mismatch == "" {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("states:%s", mismatch)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func childNames(childPaths map[string]struct{}) []string {
	names := make([]string, 0, len(childPaths))
	for childPath := range childPaths {
		names = append(names, path_lib.Base(childPath))
	}
	sort.Strings(names)
	return names
}

func readMountFile(t *testing.T, fs *RemoteFs, path string) string {
	t.Helper()
	errno, fh := fs.Open(path, fuse.O_RDONLY)
	if errno != 0 {
		t.Fatalf("Open %s: errno %d", path, errno)
	}
	defer fs.Release(path, fh)
	buf := make([]byte, 64)
	n := fs.Read(path, buf, 0, fh)
	if n < 0 {
		t.Fatalf("Read %s: errno %d", path, n)
	}
	return string(buf[:n])
}

func keptDocs(t *testing.T) (*offlineStorage, *RemoteFs) {
	t.Helper()
	storage := newOfflineStorage(map[string]string{
		"/docs/a.txt":     "alpha",
		"/docs/sub/b.txt": "bravo",
		"/other/c.txt":    "charlie",
	})
	fs := newOfflineTestFs(t, storage)
	if err := fs.KeepOffline("/docs", true); err != nil {
		t.Fatal(err)
	}
	waitForStates(t, fs, map[string]OfflineState{
		"/docs":           OfflineKept,
		"/docs/a.txt":     OfflineKept,
		"/docs/sub":       OfflineKept,
		"/docs/sub/b.txt": OfflineKept,
	})
	return storage, fs
}

func TestOfflineKeepDownloadsAFolder(t *testing.T) {
	_, fs := keptDocs(t)

	waitForStates(t, fs, map[string]OfflineState{
		"/":            OfflinePartlyKept,
		"/other":       OfflineNone,
		"/other/c.txt": OfflineNone,
	})
	for _, path := range []string{"/docs/a.txt", "/docs/sub/b.txt"} {
		if !fs.cacheStore.IsKept(path) {
			t.Errorf("%s is downloaded but not kept in the cache", path)
		}
	}
	if fs.cacheStore.IsKept("/other/c.txt") {
		t.Error("/other/c.txt is kept in the cache")
	}
}

func TestOfflineShowsKeptFoldersWithoutTheStorage(t *testing.T) {
	storage, fs := keptDocs(t)
	storage.setDown(true)

	childPaths, err := fs.listDirOrSaved("/docs")
	if err != nil {
		t.Fatalf("listing a kept folder without the storage: %v", err)
	}
	if got := strings.Join(childNames(childPaths), ","); got != "a.txt,sub" {
		t.Fatalf("kept folder lists %q without the storage, want a.txt,sub", got)
	}
	if !fs.offline.unreachable() {
		t.Fatal("the storage isn't known to be unreachable after a failed listing")
	}

	started := time.Now()
	if _, err := fs.listDirOrSaved("/other"); !errors.Is(err, errStorageUnreachable) {
		t.Fatalf("listing a folder that isn't kept without the storage: %v, want errStorageUnreachable", err)
	}
	if elapsed := time.Since(started); elapsed > time.Second {
		t.Fatalf("a folder that isn't kept took %v to fail without the storage", elapsed)
	}

	if got := readMountFile(t, fs, "/docs/sub/b.txt"); got != "bravo" {
		t.Fatalf("kept file reads %q without the storage, want bravo", got)
	}
}

func TestOfflineAsksTheStorageAgainInTheBackground(t *testing.T) {
	storage, fs := keptDocs(t)
	k := fs.offline
	storage.setDown(true)
	if _, err := fs.listDirOrSaved("/docs"); err != nil {
		t.Fatal(err)
	}

	storage.set("/docs/new.txt", "new")
	storage.setDown(false)
	k.mu.Lock()
	k.unreachableUntil = time.Now().Add(-time.Second)
	k.mu.Unlock()

	// The saved listing answers at once while the storage is asked in the background.
	childPaths, err := fs.listDirOrSaved("/docs")
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(childNames(childPaths), ","); got != "a.txt,sub" {
		t.Fatalf("listing while the storage is checked = %q, want the saved a.txt,sub", got)
	}
	deadline := time.Now().Add(5 * time.Second)
	for {
		k.mu.Lock()
		reachable := !k.wasUnreachable
		k.mu.Unlock()
		if reachable {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("the storage wasn't found reachable again in the background")
		}
		time.Sleep(20 * time.Millisecond)
	}

	childPaths, err = fs.listDirOrSaved("/docs")
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(childNames(childPaths), ","); got != "a.txt,new.txt,sub" {
		t.Fatalf("listing after the storage came back = %q, want a.txt,new.txt,sub", got)
	}
	waitForStates(t, fs, map[string]OfflineState{"/docs/new.txt": OfflineKept})
}

func TestOfflineReleasingASubfolderKeepsTheRest(t *testing.T) {
	_, fs := keptDocs(t)

	if err := fs.KeepOffline("/docs/sub", false); err != nil {
		t.Fatal(err)
	}
	waitForStates(t, fs, map[string]OfflineState{
		"/docs":           OfflinePartlyKept,
		"/docs/a.txt":     OfflineKept,
		"/docs/sub":       OfflineNone,
		"/docs/sub/b.txt": OfflineNone,
	})
	if fs.cacheStore.IsKept("/docs/sub/b.txt") {
		t.Error("released file is still kept in the cache")
	}
	if !fs.cacheStore.IsKept("/docs/a.txt") {
		t.Error("file next to the released folder is no longer kept in the cache")
	}

	if err := fs.KeepOffline("/docs", false); err != nil {
		t.Fatal(err)
	}
	waitForStates(t, fs, map[string]OfflineState{"/": OfflineNone, "/docs": OfflineNone, "/docs/a.txt": OfflineNone})
}

func TestOfflineKeptPathsAreRememberedAfterARestart(t *testing.T) {
	storage, fs := keptDocs(t)
	fs.offline.close()

	restarted, err := newOfflineKeeper(fs)
	if err != nil {
		t.Fatal(err)
	}
	fs.offline = restarted
	restarted.start()
	storage.setDown(true)

	waitForStates(t, fs, map[string]OfflineState{"/docs": OfflineKept, "/docs/a.txt": OfflineKept})
	childPaths, err := fs.listDirOrSaved("/docs")
	if err != nil {
		t.Fatalf("listing a kept folder without the storage after a restart: %v", err)
	}
	if got := strings.Join(childNames(childPaths), ","); got != "a.txt,sub" {
		t.Fatalf("kept folder lists %q after a restart, want a.txt,sub", got)
	}
}

func cachedOnDisk(fs *RemoteFs, path string) bool {
	buf := make([]byte, 1)
	n, err := fs.cacheStore.Read(path, buf, 0)
	return err == nil && n > 0
}

// waitUntilDeleted waits until path is no longer on this computer. A download
// stream that is still finishing counts as an open file, and an idle one stops
// after the streams' 30 second pause timeout.
func waitUntilDeleted(t *testing.T, fs *RemoteFs, path, why string) {
	t.Helper()
	deadline := time.Now().Add(45 * time.Second)
	for cachedOnDisk(fs, path) {
		if time.Now().After(deadline) {
			t.Fatal(why)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func TestOfflineOnlineOnlyFreesTheSpace(t *testing.T) {
	_, fs := keptDocs(t)
	errno, fh := fs.Open("/docs/a.txt", fuse.O_RDONLY)
	if errno != 0 {
		t.Fatalf("Open: errno %d", errno)
	}

	if err := fs.KeepOffline("/docs", false); err != nil {
		t.Fatal(err)
	}
	waitForStates(t, fs, map[string]OfflineState{"/docs": OfflineNone, "/docs/a.txt": OfflineNone})

	waitUntilDeleted(t, fs, "/docs/sub/b.txt", "a released file is still on this computer")
	if !cachedOnDisk(fs, "/docs/a.txt") {
		t.Fatal("an open file lost its data when it was released")
	}
	fs.Release("/docs/a.txt", fh)
	waitUntilDeleted(t, fs, "/docs/a.txt", "a released file is still on this computer after it was closed")
}

func TestOfflineKeepingAgainBeforeCloseKeepsTheData(t *testing.T) {
	_, fs := keptDocs(t)
	errno, fh := fs.Open("/docs/a.txt", fuse.O_RDONLY)
	if errno != 0 {
		t.Fatalf("Open: errno %d", errno)
	}
	if err := fs.KeepOffline("/docs", false); err != nil {
		t.Fatal(err)
	}
	if err := fs.KeepOffline("/docs", true); err != nil {
		t.Fatal(err)
	}
	waitForStates(t, fs, map[string]OfflineState{"/docs": OfflineKept, "/docs/a.txt": OfflineKept})
	fs.Release("/docs/a.txt", fh)

	if !cachedOnDisk(fs, "/docs/a.txt") {
		t.Fatal("a file kept again before it was closed was deleted when it closed")
	}
	waitForStates(t, fs, map[string]OfflineState{"/docs/a.txt": OfflineKept})
}

// waitForWalks waits until the keeper isn't walking any folder.
func waitForWalks(t *testing.T, k *offlineKeeper) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for {
		k.mu.Lock()
		walking := len(k.walks)
		k.mu.Unlock()
		if walking == 0 {
			return
		}
		if time.Now().After(deadline) {
			t.Fatal("the keeper is still walking")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func TestOfflineAFailedRefreshKeepsTheSavedVersion(t *testing.T) {
	storage, fs := keptDocs(t)
	storage.with(func() {
		storage.files["/docs/a.txt"] = []byte("alpha, changed")
		storage.listBreak["/docs"] = errOfflineTestUnreachable
	})

	// The listing shows the changed file, then fails.
	if _, err := fs.listDirOrSaved("/docs"); err != nil {
		t.Fatal(err)
	}
	storage.setDown(true)
	if got := readMountFile(t, fs, "/docs/a.txt"); got != "alpha" {
		t.Fatalf("kept file reads %q without the storage after a failed refresh, want alpha", got)
	}

	storage.with(func() {
		storage.down = false
		delete(storage.listBreak, "/docs")
	})
	fs.offline.markReachable()
	if _, err := fs.listDirOrSaved("/docs"); err != nil {
		t.Fatal(err)
	}
	waitForStates(t, fs, map[string]OfflineState{"/docs/a.txt": OfflineKept})
	if got := readMountFile(t, fs, "/docs/a.txt"); got != "alpha, changed" {
		t.Fatalf("kept file reads %q after a refresh, want the new version", got)
	}
}

func newOfflineGate() *offlineGate {
	return &offlineGate{reached: make(chan struct{}), open: make(chan struct{})}
}

func TestOfflineChangedFileKeepsItsLastVersionUntilTheNewOneIsIn(t *testing.T) {
	storage, fs := keptDocs(t)
	gate := newOfflineGate()
	storage.with(func() {
		storage.files["/docs/a.txt"] = []byte("alpha, changed")
		storage.downloadGate["/docs/a.txt"] = gate
	})

	// The refresh finds the change, and the new version's download waits.
	if _, err := fs.listDirOrSaved("/docs"); err != nil {
		t.Fatal(err)
	}
	<-gate.reached
	waitForStates(t, fs, map[string]OfflineState{"/docs/a.txt": OfflineDownloading})

	storage.setDown(true)
	if _, err := fs.listDirOrSaved("/docs"); err != nil {
		t.Fatal(err)
	}
	if got := readMountFile(t, fs, "/docs/a.txt"); got != "alpha" {
		t.Fatalf("changed kept file reads %q without the storage before its new version is in, want alpha", got)
	}

	storage.setDown(false)
	fs.offline.markReachable()
	close(gate.open)
	waitForStates(t, fs, map[string]OfflineState{"/docs/a.txt": OfflineKept})
	storage.setDown(true)
	if _, err := fs.listDirOrSaved("/docs"); err != nil {
		t.Fatal(err)
	}
	if got := readMountFile(t, fs, "/docs/a.txt"); got != "alpha, changed" {
		t.Fatalf("changed kept file reads %q without the storage after its new version is in, want the new version", got)
	}
}

func TestOfflineNewVersionDoesNotChangeASavedListingInUse(t *testing.T) {
	storage, fs := keptDocs(t)
	gate := newOfflineGate()
	storage.with(func() {
		storage.files["/docs/a.txt"] = []byte("alpha, changed")
		storage.downloadGate["/docs/a.txt"] = gate
	})
	if _, err := fs.listDirOrSaved("/docs"); err != nil {
		t.Fatal(err)
	}
	<-gate.reached

	// The saved listing is shown while the new version comes in. The race
	// detector catches a listing changed under its reader.
	done := make(chan struct{})
	var showing sync.WaitGroup
	showing.Go(func() {
		for {
			select {
			case <-done:
				return
			default:
				fs.offline.savedChildren("/docs")
			}
		}
	})
	close(gate.open)
	waitForStates(t, fs, map[string]OfflineState{"/docs/a.txt": OfflineKept})
	close(done)
	showing.Wait()
}

func TestOfflineRefreshKeepsANewSubfolder(t *testing.T) {
	storage, fs := keptDocs(t)
	storage.set("/docs/new/d.txt", "delta")

	if _, err := fs.listDirOrSaved("/docs"); err != nil {
		t.Fatal(err)
	}
	waitForStates(t, fs, map[string]OfflineState{"/docs/new": OfflineKept, "/docs/new/d.txt": OfflineKept})

	storage.setDown(true)
	childPaths, err := fs.listDirOrSaved("/docs/new")
	if err != nil {
		t.Fatalf("listing a new kept subfolder without the storage: %v", err)
	}
	if got := strings.Join(childNames(childPaths), ","); got != "d.txt" {
		t.Fatalf("new kept subfolder lists %q without the storage, want d.txt", got)
	}
	if got := readMountFile(t, fs, "/docs/new/d.txt"); got != "delta" {
		t.Fatalf("file in a new kept subfolder reads %q without the storage, want delta", got)
	}
}

func TestOfflineSubfolderIsNotKeptBeforeItsListingIsSaved(t *testing.T) {
	storage, fs := keptDocs(t)
	storage.with(func() {
		storage.files["/docs/new/d.txt"] = []byte("delta")
		storage.listErr["/docs/new"] = fmt.Errorf("listing /docs/new: %w", iofs.ErrPermission)
	})

	if _, err := fs.listDirOrSaved("/docs"); err != nil {
		t.Fatal(err)
	}
	waitForWalks(t, fs.offline)
	if state := fs.OfflineStates([]string{"/docs/new"})[0]; state != OfflineDownloading {
		t.Fatalf("a kept subfolder that couldn't be listed is %v, want downloading", state)
	}
}

func TestOfflineFolderIsNotKeptWhileASubfolderIsNotSaved(t *testing.T) {
	storage, fs := keptDocs(t)
	storage.with(func() {
		storage.files["/docs/new/d.txt"] = []byte("delta")
		storage.listErr["/docs/new"] = fmt.Errorf("listing /docs/new: %w", iofs.ErrPermission)
	})

	if _, err := fs.listDirOrSaved("/docs"); err != nil {
		t.Fatal(err)
	}
	waitForWalks(t, fs.offline)
	waitForStates(t, fs, map[string]OfflineState{"/docs": OfflineDownloading, "/docs/new": OfflineDownloading})

	// The next refresh saves the subfolder.
	storage.with(func() { delete(storage.listErr, "/docs/new") })
	if _, err := fs.listDirOrSaved("/docs"); err != nil {
		t.Fatal(err)
	}
	waitForStates(t, fs, map[string]OfflineState{"/docs": OfflineKept, "/docs/new": OfflineKept, "/docs/new/d.txt": OfflineKept})
}

func TestOfflineRemovedFolderStopsItsDownloads(t *testing.T) {
	storage, fs := keptDocs(t)
	gate := newOfflineGate()
	t.Cleanup(func() { close(gate.open) })
	storage.with(func() {
		storage.files["/docs/sub/x.txt"] = []byte("x-ray")
		storage.downloadGate["/docs/sub/x.txt"] = gate
	})
	if _, err := fs.listDirOrSaved("/docs/sub"); err != nil {
		t.Fatal(err)
	}
	<-gate.reached

	// The folder is removed on the storage while its file downloads.
	storage.with(func() {
		delete(storage.files, "/docs/sub/b.txt")
		delete(storage.files, "/docs/sub/x.txt")
	})
	if _, err := fs.listDirOrSaved("/docs"); err != nil {
		t.Fatal(err)
	}
	waitForStates(t, fs, map[string]OfflineState{"/docs": OfflineKept})
}

func TestOfflineRenamedFolderLeavesNothingUnderItsOldName(t *testing.T) {
	storage, fs := keptDocs(t)
	gate := newOfflineGate()
	t.Cleanup(func() { close(gate.open) })
	storage.with(func() {
		storage.files["/docs/sub/x.txt"] = []byte("x-ray")
		storage.downloadGate["/docs/sub/x.txt"] = gate
	})
	if _, err := fs.listDirOrSaved("/docs/sub"); err != nil {
		t.Fatal(err)
	}
	<-gate.reached

	// The folder is renamed while one of its files downloads.
	if errc := fs.Rename("/docs/sub", "/docs/renamed"); errc != 0 {
		t.Fatalf("Rename: errc %d", errc)
	}
	waitForStates(t, fs, map[string]OfflineState{
		"/docs":               OfflineKept,
		"/docs/renamed":       OfflineKept,
		"/docs/renamed/b.txt": OfflineKept,
		"/docs/renamed/x.txt": OfflineKept,
	})
	if fs.cacheStore.IsKept("/docs/sub/b.txt") {
		t.Fatal("a file is still kept under its folder's old name")
	}
	deadline := time.Now().Add(10 * time.Second)
	for cachedOnDisk(fs, "/docs/sub/b.txt") {
		if time.Now().After(deadline) {
			t.Fatal("a file is still on this computer under its folder's old name")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func TestOfflineRenamedFolderShowsItsNewNameOffline(t *testing.T) {
	storage, fs := keptDocs(t)

	if errc := fs.Rename("/docs/sub", "/docs/renamed"); errc != 0 {
		t.Fatalf("Rename: errc %d", errc)
	}
	waitForStates(t, fs, map[string]OfflineState{"/docs/renamed": OfflineKept, "/docs/renamed/b.txt": OfflineKept})

	storage.setDown(true)
	childPaths, err := fs.listDirOrSaved("/docs")
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(childNames(childPaths), ","); got != "a.txt,renamed" {
		t.Fatalf("kept folder lists %q without the storage after a rename, want a.txt,renamed", got)
	}
	if got := readMountFile(t, fs, "/docs/renamed/b.txt"); got != "bravo" {
		t.Fatalf("renamed kept file reads %q without the storage, want bravo", got)
	}
}

func TestOfflineRenamingAFolderKeepsTheFileKeptInIt(t *testing.T) {
	storage := newOfflineStorage(map[string]string{
		"/docs/a.txt":     "alpha",
		"/docs/sub/b.txt": "bravo",
	})
	fs := newOfflineTestFs(t, storage)
	if err := fs.KeepOffline("/docs/sub/b.txt", true); err != nil {
		t.Fatal(err)
	}
	waitForStates(t, fs, map[string]OfflineState{"/docs/sub/b.txt": OfflineKept})

	if errc := fs.Rename("/docs/sub", "/docs/renamed"); errc != 0 {
		t.Fatalf("Rename: errc %d", errc)
	}
	waitForStates(t, fs, map[string]OfflineState{"/docs/renamed": OfflinePartlyKept, "/docs/renamed/b.txt": OfflineKept})

	storage.setDown(true)
	if _, err := fs.listDirOrSaved("/docs/renamed"); err != nil {
		t.Fatalf("listing the renamed folder without the storage: %v", err)
	}
	if got := readMountFile(t, fs, "/docs/renamed/b.txt"); got != "bravo" {
		t.Fatalf("kept file in a renamed folder reads %q without the storage, want bravo", got)
	}
}

func TestOfflineOnlineOnlyWinsOverAKeepInProgress(t *testing.T) {
	storage := newOfflineStorage(map[string]string{"/docs/a.txt": "alpha"})
	fs := newOfflineTestFs(t, storage)
	gate := &offlineGate{reached: make(chan struct{}), open: make(chan struct{})}
	storage.with(func() { storage.listGate["/docs"] = gate })

	if err := fs.KeepOffline("/docs/a.txt", true); err != nil {
		t.Fatal(err)
	}
	<-gate.reached
	if err := fs.KeepOffline("/docs/a.txt", false); err != nil {
		t.Fatal(err)
	}
	close(gate.open)
	waitForWalks(t, fs.offline)

	if state := fs.OfflineStates([]string{"/docs/a.txt"})[0]; state != OfflineNone {
		t.Fatalf("a file made Online Only during its keep is %v, want none", state)
	}
	if fs.cacheStore.IsKept("/docs/a.txt") {
		t.Fatal("a file made Online Only during its keep is kept in the cache")
	}
	fs.offline.mu.Lock()
	_, pending := fs.offline.pending["/docs/a.txt"]
	fs.offline.mu.Unlock()
	if pending {
		t.Fatal("a file made Online Only during its keep is waiting to download")
	}
}

func TestOfflineOnlineOnlyDeletesAnOpenFileItsDownloadWroteTo(t *testing.T) {
	storage := newOfflineStorage(map[string]string{"/docs/a.txt": "alpha"})
	fs := newOfflineTestFs(t, storage)
	if _, err := fs.listDirOrSaved("/docs"); err != nil {
		t.Fatal(err)
	}
	errno, fh := fs.Open("/docs/a.txt", fuse.O_RDONLY)
	if errno != 0 {
		t.Fatalf("Open: errno %d", errno)
	}
	gate := newOfflineGate()
	storage.with(func() { storage.downloadGate["/docs/a.txt"] = gate })

	// Online Only comes while the keep's download waits for its first bytes.
	if err := fs.KeepOffline("/docs", true); err != nil {
		t.Fatal(err)
	}
	<-gate.reached
	if err := fs.KeepOffline("/docs", false); err != nil {
		t.Fatal(err)
	}
	close(gate.open)
	downloaded := savedEntry{Size: int64(len("alpha")), ModTime: storage.mtime}
	deadline := time.Now().Add(10 * time.Second)
	for {
		fs.offline.mu.Lock()
		complete := fs.offline.complete("/docs/a.txt", downloaded)
		fs.offline.mu.Unlock()
		if complete {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("the download didn't finish")
		}
		time.Sleep(20 * time.Millisecond)
	}

	fs.Release("/docs/a.txt", fh)
	waitUntilDeleted(t, fs, "/docs/a.txt", "a file made Online Only is still on this computer after it was closed")
}

func TestOfflineAccessDeniedIsNotTakenForNoStorage(t *testing.T) {
	storage, fs := keptDocs(t)
	storage.set("/secret/s.txt", "secret")
	if err := fs.KeepOffline("/secret", true); err != nil {
		t.Fatal(err)
	}
	waitForStates(t, fs, map[string]OfflineState{"/secret/s.txt": OfflineKept})
	storage.with(func() {
		storage.listErr["/secret"] = fmt.Errorf("listing /secret: %w", iofs.ErrPermission)
	})

	if _, err := fs.listDirOrSaved("/secret"); err == nil {
		t.Fatal("a kept folder the storage refuses lists from its saved listing")
	}
	childPaths, err := fs.listDirOrSaved("/other")
	if err != nil {
		t.Fatalf("listing another folder after one was refused: %v", err)
	}
	if got := strings.Join(childNames(childPaths), ","); got != "c.txt" {
		t.Fatalf("other folder lists %q, want c.txt", got)
	}
}

func TestOfflineAPIRefusalIsNotTakenForNoStorage(t *testing.T) {
	storage, fs := keptDocs(t)
	storage.set("/secret/s.txt", "secret")
	if err := fs.KeepOffline("/secret", true); err != nil {
		t.Fatal(err)
	}
	waitForStates(t, fs, map[string]OfflineState{"/secret/s.txt": OfflineKept})
	storage.with(func() {
		storage.listErr["/secret"] = files_sdk.ResponseError{
			Type:         string(files_sdk.ErrReadPermissionRequired),
			ErrorMessage: "Read permission required",
			HttpCode:     http.StatusForbidden,
		}
	})

	if _, err := fs.listDirOrSaved("/secret"); err == nil {
		t.Fatal("a kept folder the API refuses lists from its saved listing")
	}
	childPaths, err := fs.listDirOrSaved("/other")
	if err != nil {
		t.Fatalf("listing another folder after the API refused one: %v", err)
	}
	if got := strings.Join(childNames(childPaths), ","); got != "c.txt" {
		t.Fatalf("other folder lists %q, want c.txt", got)
	}

	// An unavailable service is shown from what is saved.
	storage.with(func() {
		storage.listErr["/secret"] = files_sdk.ResponseError{
			Type:         string(files_sdk.ErrMigrationInProgress),
			ErrorMessage: "Migration in progress",
			HttpCode:     http.StatusServiceUnavailable,
		}
	})
	childPaths, err = fs.listDirOrSaved("/secret")
	if err != nil {
		t.Fatalf("listing a kept folder while the service is unavailable: %v", err)
	}
	if got := strings.Join(childNames(childPaths), ","); got != "s.txt" {
		t.Fatalf("kept folder lists %q while the service is unavailable, want the saved s.txt", got)
	}
}

func TestOfflineSlowListingIsNotTakenForNoStorage(t *testing.T) {
	timeout := offlineRemoteTimeout
	offlineRemoteTimeout = 100 * time.Millisecond
	t.Cleanup(func() { offlineRemoteTimeout = timeout })
	storage, fs := keptDocs(t)

	// The storage answers at once, then takes longer than the timeout for
	// the rest of the listing.
	storage.set("/docs/z.txt", "zulu")
	storage.with(func() { storage.listPause["/docs"] = offlinePause{before: 1, pause: 300 * time.Millisecond} })
	childPaths, err := fs.listDirOrSaved("/docs")
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(childNames(childPaths), ","); got != "a.txt,sub,z.txt" {
		t.Fatalf("slow listing = %q, want a.txt,sub,z.txt", got)
	}
	if fs.offline.unreachable() {
		t.Fatal("a slow listing marked the storage unreachable")
	}

	// Storage that doesn't answer at all still gives up after the timeout.
	storage.with(func() { storage.listPause["/docs"] = offlinePause{before: 0, pause: 10 * time.Second} })
	started := time.Now()
	childPaths, err = fs.listDirOrSaved("/docs")
	if err != nil {
		t.Fatal(err)
	}
	if elapsed := time.Since(started); elapsed > 2*time.Second {
		t.Fatalf("storage that didn't answer took %v to give up", elapsed)
	}
	if got := strings.Join(childNames(childPaths), ","); got != "a.txt,sub,z.txt" {
		t.Fatalf("listing without an answer = %q, want the saved a.txt,sub,z.txt", got)
	}
	if !fs.offline.unreachable() {
		t.Fatal("storage that didn't answer isn't known to be unreachable")
	}
}

func TestOfflineAnswerDeadlineReachesDerivedContexts(t *testing.T) {
	deadline, cancel := newAnswerDeadline(20 * time.Millisecond)
	defer cancel()
	// Storage clients make their own contexts from the one they are given.
	derived, cancelDerived := context.WithTimeout(deadline, time.Minute)
	defer cancelDerived()

	<-derived.Done()
	if !errors.Is(derived.Err(), context.DeadlineExceeded) {
		t.Fatalf("a context made from the answer deadline ended with %v, want the deadline", derived.Err())
	}
}
