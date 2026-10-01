//go:build linux || windows

package fsmount

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/fsmount/events"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/disk"
	"github.com/winfsp/cgofuse/fuse"
)

func BenchmarkRemoteFsDiskSequentialRead(b *testing.B) {
	const size = 64 * 1024 * 1024
	payload := bytes.Repeat([]byte("x"), size)
	for _, sparse := range []bool{false, true} {
		b.Run(fmt.Sprintf("sparse=%t", sparse), func(b *testing.B) {
			b.SetBytes(size)
			var totalRequests int64
			for iteration := 0; iteration < b.N; iteration++ {
				b.StopTimer()
				var requests atomic.Int64
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("ETag", `"fixture"`)
					start, end, partial, err := syntheticRequestRange(r.Header.Get("Range"), size)
					if err != nil {
						b.Error(err)
						w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
						return
					}
					requests.Add(1)
					w.Header().Set("Content-Length", fmt.Sprint(end-start+1))
					if partial {
						w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, size))
						w.WriteHeader(http.StatusPartialContent)
					}
					_, _ = w.Write(payload[start : end+1])
				}))
				fs, vfs, _ := newTestRemoteFs(b)
				dc, err := disk.NewDiskCache(b.TempDir())
				if err != nil {
					b.Fatal(err)
				}
				b.Cleanup(dc.StopMaintenance)
				fs.cacheStore = dc
				fs.sparseRangeReads = sparse
				fs.backend = newRangeTestSDKBackend(files_sdk.Config{DisableDirectTransfers: true})
				vfs.cacheTTL = time.Hour
				path := "/benchmark"
				now := time.Now().Truncate(time.Second)
				node := vfs.getOrCreate(path, nodeTypeFile)
				node.updateInfo(fsNodeInfo{nodeType: nodeTypeFile, size: size, modTime: now})
				node.setDownloadURI(server.URL)
				errno, fh := fs.Open(path, fuse.O_RDONLY)
				if errno != 0 {
					b.Fatal(errno)
				}
				buffer := make([]byte, 4*1024*1024)
				b.StartTimer()
				for offset := int64(0); offset < size; offset += int64(len(buffer)) {
					if n := fs.Read(path, buffer, offset, fh); n != len(buffer) || !bytes.Equal(buffer, payload[offset:offset+int64(len(buffer))]) {
						b.Fatalf("invalid bytes at %d (n=%d)", offset, n)
					}
				}
				deadline := time.Now().Add(10 * time.Second)
				for fs.rangeDownloads().diagnosticsSnapshot().ActiveRequests != 0 {
					if time.Now().After(deadline) {
						b.Fatal("download did not finish")
					}
					time.Sleep(time.Millisecond)
				}
				complete, err := dc.RangeEntryComplete(path, cacheEntryMetadata(path, size, now))
				if err != nil || !complete {
					b.Fatalf("incomplete cache: %t %v", complete, err)
				}
				b.StopTimer()
				totalRequests += requests.Load()
				fs.Release(path, fh)
				fs.closeRangeDownloads()
				vfs.destroy()
				server.Close()
			}
			b.ReportMetric(float64(totalRequests)/float64(b.N), "requests/op")
		})
	}
}

func TestRemoteFsEvictionDuringBackgroundDownloadKeepsCacheIntact(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()
	defer fs.closeRangeDownloads()
	fs.sparseRangeReads = false
	const size = 2 * cacheWriteSize
	dc, err := disk.NewDiskCache(t.TempDir(), disk.WithCapacityBytes(size+1))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(dc.StopMaintenance)
	fs.cacheStore = dc
	path := "/download"
	now := time.Now().Truncate(time.Second)
	node := vfs.getOrCreate(path, nodeTypeFile)
	node.updateInfo(fsNodeInfo{nodeType: nodeTypeFile, size: size, modTime: now})
	reader, writer := io.Pipe()
	release := make(chan struct{})
	var once sync.Once
	resume := func() { once.Do(func() { close(release) }) }
	defer resume()
	go func() {
		_, _ = writer.Write(bytes.Repeat([]byte("A"), cacheWriteSize))
		<-release
		_, _ = writer.Write(bytes.Repeat([]byte("B"), cacheWriteSize))
		_ = writer.Close()
	}()
	fs.backend = &fakeRemoteBackend{downloadFunc: func(params files_sdk.FileDownloadParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
		_, err := files_sdk.BuildResponse(&http.Response{StatusCode: http.StatusOK, ContentLength: size, Body: reader}, opts...)
		return params.File, err
	}}
	errno, fh := fs.Open(path, fuse.O_RDONLY)
	if errno != 0 {
		t.Fatal(errno)
	}
	buffer := make([]byte, 1)
	if n := fs.Read(path, buffer, 0, fh); n != 1 || buffer[0] != 'A' {
		t.Fatalf("initial read: %d %q", n, buffer)
	}
	fs.Release(path, fh)
	other := cacheEntryMetadata("/other", size, now)
	if _, err := dc.WriteCompleteRange(other.Path, other, bytes.Repeat([]byte("x"), size), 0); err != nil {
		t.Fatal(err)
	}
	if n, err := dc.ReadRange(path, cacheEntryMetadata(path, size, now), buffer, 0); n != 1 || err != nil || buffer[0] != 'A' {
		t.Fatalf("active download lost its prefix under cache pressure: n=%d err=%v", n, err)
	}
	// Relieve pressure before completion so the finished entry can be retained.
	dc.Delete(other.Path)
	resume()
	deadline := time.Now().Add(3 * time.Second)
	for fs.rangeDownloads().diagnosticsSnapshot().ActiveRequests != 0 {
		if time.Now().After(deadline) {
			t.Fatal("download did not finish")
		}
		time.Sleep(time.Millisecond)
	}
	errno, reopened := fs.Open(path, fuse.O_RDONLY)
	if errno != 0 {
		t.Fatal(errno)
	}
	n := fs.Read(path, buffer, 0, reopened)
	fs.Release(path, reopened)
	if n != 1 || buffer[0] != 'A' {
		t.Fatalf("read-close-evict-reopen returned %d bytes: %q, want A", n, buffer)
	}
	if !dc.Delete(path) {
		t.Fatal("completed download retained its cache pin")
	}
}

func TestRemoteFsURLRefreshRejectsChangedFileVersion(t *testing.T) {
	fs, vfs, store := newTestRemoteFs(t)
	defer vfs.destroy()
	defer fs.closeRangeDownloads()
	path := "/versions"
	oldTime := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	newTime := oldTime.Add(time.Minute)
	meta := cacheEntryMetadata(path, rangeVersionTestSize, oldTime)
	cached := meta
	cached.ETag = `"fixture"`
	if _, err := store.WriteRange(path, cached, []byte("AAAA"), 0); err != nil {
		t.Fatal(err)
	}
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/stale":
			w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", rangeVersionTestSize))
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		case "/fresh":
			serveRangeVersionTest(t, w, r, `"fixture"`, http.StatusPartialContent, 'B')
		default:
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprintf(w, `{"path":"versions","size":%d,"mtime":%q,"download_uri":%q}`, rangeVersionTestSize, newTime.Format(time.RFC3339), server.URL+"/fresh")
		}
	}))
	defer server.Close()
	node := vfs.getOrCreate(path, nodeTypeFile)
	node.updateInfo(fsNodeInfo{nodeType: nodeTypeFile, size: rangeVersionTestSize, modTime: oldTime})
	node.setDownloadURI(server.URL + "/stale")
	fs.backend = newRangeTestSDKBackend(files_sdk.Config{DisableDirectTransfers: true, EndpointOverride: server.URL})
	requested := cache.ByteRange{Start: 2 << 20, End: 2<<20 + 4}
	err := fs.rangeDownloads().Wait(context.Background(), 1, path, meta, requested)
	if !errors.Is(err, errRangeVersionChanged) {
		t.Fatalf("expected version error after refreshed mtime, got %v", err)
	}

}

func TestRemoteFsURLRefreshRejectsChangedVersionForCompleteDownload(t *testing.T) {
	fs, vfs, store := newTestRemoteFs(t)
	defer vfs.destroy()
	defer fs.closeRangeDownloads()
	fs.sparseRangeReads = false
	path := "/versions"
	oldTime := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	meta := cacheEntryMetadata(path, 8, oldTime)
	meta.ETag = `"fixture"`
	if _, err := store.WriteRange(path, meta, []byte("AAAA"), 0); err != nil {
		t.Fatal(err)
	}
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", `"fixture"`)
		if r.URL.Path == "/fresh" {
			if r.Header.Get("Range") != "" {
				t.Error("disabled mode sent a Range header")
			}
			w.Header().Set("Content-Length", "8")
			_, _ = w.Write([]byte("BBBBBBBB"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"path":"versions","size":8,"mtime":%q,"download_uri":%q}`, oldTime.Add(time.Minute).Format(time.RFC3339), server.URL+"/fresh")
	}))
	defer server.Close()
	node := vfs.getOrCreate(path, nodeTypeFile)
	node.updateInfo(fsNodeInfo{nodeType: nodeTypeFile, size: 8, modTime: oldTime})
	fs.backend = newRangeTestSDKBackend(files_sdk.Config{DisableDirectTransfers: true, EndpointOverride: server.URL})
	if err := fs.ensureFullyCached(path, 8); !errors.Is(err, errRangeVersionChanged) {
		t.Fatalf("expected version error for complete download, got %v", err)
	}
}

func TestRemoteFsDownloadLimitCoversOpenResponseBodies(t *testing.T) {
	for _, sparse := range []bool{false, true} {
		t.Run(fmt.Sprintf("sparse=%t", sparse), func(t *testing.T) {
			fs, vfs, _ := newTestRemoteFs(t)
			defer vfs.destroy()
			const total = downloadOpLimit + 2
			started := make(chan struct{}, total)
			release := make(chan struct{})
			unblock := sync.OnceFunc(func() { close(release) })
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("ETag", `"fixture"`)
				w.Header().Set("Content-Length", "1")
				if r.Header.Get("Range") != "" {
					w.Header().Set("Content-Range", "bytes 0-0/1")
					w.WriteHeader(http.StatusPartialContent)
				} else {
					w.WriteHeader(http.StatusOK)
				}
				w.(http.Flusher).Flush()
				started <- struct{}{}
				select {
				case <-release:
					_, _ = w.Write([]byte("x"))
				case <-r.Context().Done():
				}
			}))
			defer server.Close()
			fs.backend = newRangeTestSDKBackend(files_sdk.Config{DisableDirectTransfers: true})
			ctx, cancel := context.WithCancel(context.Background())
			var wg sync.WaitGroup
			defer func() { unblock(); cancel(); wg.Wait() }()
			results := make(chan error, total)
			for i := 0; i < total; i++ {
				path := fmt.Sprintf("/limit-%d", i)
				node := vfs.getOrCreate(path, nodeTypeFile)
				node.setDownloadURI(server.URL)
				wg.Add(1)
				go func() {
					defer wg.Done()
					meta := cacheEntryMetadata(path, 1, time.Now())
					meta.ETag = `"fixture"`
					source := &remoteFsRangeSource{fs: fs}
					download := source.DownloadComplete
					if sparse {
						download = source.DownloadRange
					}
					response, err := download(ctx, path, meta, cache.ByteRange{End: 1})
					if err == nil {
						_, err = io.Copy(io.Discard, response.Body)
						_ = response.Body.Close()
						_ = response.Body.Close() // Repeated close must not release another slot.
					}
					results <- err
				}()
			}
			for i := 0; i < downloadOpLimit; i++ {
				select {
				case <-started:
				case <-time.After(3 * time.Second):
					t.Fatal("download slots were not filled")
				}
			}
			select {
			case <-started:
				t.Fatalf("more than %d bodies active", downloadOpLimit)
			case <-time.After(100 * time.Millisecond):
			}
			unblock()
			for i := 0; i < total; i++ {
				select {
				case err := <-results:
					if err != nil {
						t.Fatal(err)
					}
				case <-time.After(3 * time.Second):
					t.Fatal("queued download did not finish after body close")
				}
			}
		})
	}
}

// newCacheClearTestRemoteFs returns a mount with one uncached file of size
// bytes, registered so ClearDiskCache clears its disk cache.
func newCacheClearTestRemoteFs(t *testing.T, sparse bool, size int64) (*RemoteFs, string, time.Time) {
	t.Helper()
	dc, err := disk.NewDiskCache(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(dc.StopMaintenance)
	fs, vfs, _ := newTestRemoteFs(t)
	t.Cleanup(vfs.destroy)
	t.Cleanup(fs.closeRangeDownloads)
	fs.sparseRangeReads = sparse
	fs.cacheStore = dc
	installCacheControlTestRegistry(t, map[string]*Host{"mount": {fs: &Filescomfs{remote: fs}}})
	path := "/blocked"
	now := time.Now().Truncate(time.Second)
	node := vfs.getOrCreate(path, nodeTypeFile)
	node.updateInfo(fsNodeInfo{nodeType: nodeTypeFile, size: size, modTime: now})
	return fs, path, now
}

// blockFirstDownload returns a backend whose first download waits until its
// request is canceled. Later downloads return zeros.
func blockFirstDownload(size int64) (*fakeRemoteBackend, <-chan struct{}, *atomic.Int64) {
	started := make(chan struct{})
	calls := &atomic.Int64{}
	wait := func(opts []files_sdk.RequestResponseOption) error {
		if calls.Add(1) > 1 {
			return nil
		}
		close(started)
		ctx := files_sdk.ContextOption(opts)
		<-ctx.Done()
		return ctx.Err()
	}
	backend := &fakeRemoteBackend{
		downloadRangeFunc: func(_ files_sdk.FileDownloadParams, requested cache.ByteRange, opts ...files_sdk.RequestResponseOption) (remoteRangeResponse, error) {
			if err := wait(opts); err != nil {
				return remoteRangeResponse{}, err
			}
			return zeroRangeResponse(size, requested), nil
		},
		downloadFunc: func(params files_sdk.FileDownloadParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
			if err := wait(opts); err != nil {
				return params.File, err
			}
			_, err := files_sdk.BuildResponse(&http.Response{StatusCode: http.StatusOK, ContentLength: size, Body: io.NopCloser(io.LimitReader(zeroReader{}, size))}, opts...)
			return params.File, err
		},
	}
	return backend, started, calls
}

// clearDiskCacheWithin clears the cache, failing the test when the clear waits
// for an operation blocked on a download instead of stopping that download.
func clearDiskCacheWithin(t *testing.T, fs *RemoteFs, blocked <-chan int) {
	t.Helper()
	clearDone := make(chan error, 1)
	go func() { _, err := ClearDiskCache(); clearDone <- err }()
	select {
	case err := <-clearDone:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		fs.cancelAllRangeDownloads()
		<-blocked
		<-clearDone
		t.Fatal("ClearDiskCache waited for a blocked operation instead of stopping its download")
	}
}

func TestClearDiskCacheRestartsBlockedRead(t *testing.T) {
	for _, sparse := range []bool{false, true} {
		t.Run(fmt.Sprintf("sparse=%t", sparse), func(t *testing.T) {
			const size = 4 << 20
			fs, path, now := newCacheClearTestRemoteFs(t, sparse, size)
			backend, started, calls := blockFirstDownload(size)
			fs.backend = backend
			publisher := &captureEventPublisher{}
			fs.events = publisher
			errno, fh := fs.Open(path, fuse.O_RDONLY)
			if errno != 0 {
				t.Fatal(errno)
			}
			defer fs.Release(path, fh)
			readDone := make(chan int, 1)
			go func() { readDone <- fs.Read(path, make([]byte, 1), 0, fh) }()
			<-started
			oldDownloads := fs.rangeDownloads()
			clearDiskCacheWithin(t, fs, readDone)
			select {
			case n := <-readDone:
				if n != 1 {
					t.Fatalf("read interrupted by a cache clear returned %d, want 1", n)
				}
			case <-time.After(3 * time.Second):
				t.Fatal("read interrupted by a cache clear did not finish")
			}
			if got := calls.Load(); got < 2 {
				t.Fatalf("downloads = %d, want the read to download again after the clear", got)
			}
			// The stopped download is not a failed transfer; the retry completes.
			deadline := time.Now().Add(3 * time.Second)
			for {
				statuses := map[events.TransferStatus]int{}
				for _, transfer := range publisher.transferEvents() {
					statuses[transfer.Status]++
				}
				if statuses[events.TransferStatusErrored] > 0 {
					t.Fatalf("a download stopped by the clear was reported as failed: %#v", publisher.transferEvents())
				}
				if statuses[events.TransferStatusCanceled] > 0 && statuses[events.TransferStatusComplete] > 0 {
					break
				}
				if time.Now().After(deadline) {
					t.Fatalf("transfer events = %#v, want the stopped download canceled and the retry complete", publisher.transferEvents())
				}
				time.Sleep(10 * time.Millisecond)
			}
			meta := cacheEntryMetadata(path, size, now)
			if err := oldDownloads.Wait(context.Background(), fh, path, meta, cache.ByteRange{End: 1}); !errors.Is(err, errRangeStreamsClosed) {
				t.Fatalf("retired downloads accepted work: %v", err)
			}
		})
	}
}

// Deleting a large cleared cache takes a while, and a read the clear
// interrupted does not wait for it: the drive is usable once the cleared
// entries are unreadable.
func TestClearDiskCacheResumesReadsBeforeDeletingFiles(t *testing.T) {
	const size = 1 << 20
	fs, path, _ := newCacheClearTestRemoteFs(t, false, size)
	backend, started, _ := blockFirstDownload(size)
	fs.backend = backend
	dc := fs.cacheStore.(*disk.DiskCache)
	for i := range 20000 {
		if _, err := dc.Write(fmt.Sprintf("/filler/%d", i), []byte{1}, 0); err != nil {
			t.Fatal(err)
		}
	}
	for _, nodePath := range []string{"/", path} {
		if node, ok := fs.vfs.fetch(nodePath); ok {
			node.extendTtl()
		}
	}
	errno, fh := fs.Open(path, fuse.O_RDONLY)
	if errno != 0 {
		t.Fatal(errno)
	}
	defer fs.Release(path, fh)
	readDone := make(chan int, 1)
	go func() { readDone <- fs.Read(path, make([]byte, 1), 0, fh) }()
	<-started

	clearDone := make(chan error, 1)
	go func() { _, err := ClearDiskCache(); clearDone <- err }()
	select {
	case n := <-readDone:
		if n != 1 {
			t.Fatalf("read interrupted by a cache clear returned %d, want 1", n)
		}
	case err := <-clearDone:
		t.Fatalf("the clear deleted every file (err %v) before the read it interrupted resumed", err)
	case <-time.After(10 * time.Second):
		t.Fatal("read interrupted by a cache clear did not finish")
	}
	if err := <-clearDone; err != nil {
		t.Fatal(err)
	}
}

func TestClearDiskCacheRestartsBlockedEditBaseline(t *testing.T) {
	for _, operation := range []string{"write", "truncate"} {
		for _, sparse := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/sparse=%t", operation, sparse), func(t *testing.T) {
				const size = 4 << 20
				fs, path, _ := newCacheClearTestRemoteFs(t, sparse, size)
				backend, started, calls := blockFirstDownload(size)
				fs.backend = backend
				fs.uploadWorkingCopy = func(_ context.Context, _ *fsNode, _ string, reader uploadWorkingCopyReader, mtime time.Time, _ uint64) (uploadedFileMetadata, error) {
					n, err := io.Copy(io.Discard, reader)
					return testUploadedMetadata(n, mtime), err
				}
				errno, fh := fs.Open(path, fuse.O_RDWR)
				if errno != 0 {
					t.Fatal(errno)
				}
				defer fs.Release(path, fh)
				want := 1
				if operation == "truncate" {
					want = 0
				}
				done := make(chan int, 1)
				go func() {
					if operation == "write" {
						done <- fs.Write(path, []byte("x"), size-1, fh)
					} else {
						done <- fs.Truncate(path, size/2, fh)
					}
				}()
				<-started
				clearDiskCacheWithin(t, fs, done)
				select {
				case n := <-done:
					if n != want {
						t.Fatalf("%s interrupted by a cache clear returned %d, want %d", operation, n, want)
					}
				case <-time.After(3 * time.Second):
					t.Fatalf("%s interrupted by a cache clear did not finish", operation)
				}
				if got := calls.Load(); got < 2 {
					t.Fatalf("downloads = %d, want the baseline to download again after the clear", got)
				}
				buffer := []byte{'?'}
				if n := fs.Read(path, buffer, 1, fh); n != 1 || buffer[0] != 0 {
					t.Fatalf("working copy read = %d %q, want the downloaded baseline", n, buffer)
				}
			})
		}
	}
}

func TestReadDuringCacheClearWaitsForIt(t *testing.T) {
	const size = 4 << 20
	fs, path, _ := newCacheClearTestRemoteFs(t, false, size)
	backend, _, calls := blockFirstDownload(size)
	calls.Store(1) // Do not block: this read starts while the clear already runs.
	fs.backend = backend
	errno, fh := fs.Open(path, fuse.O_RDONLY)
	if errno != 0 {
		t.Fatal(errno)
	}
	defer fs.Release(path, fh)
	fs.pauseRangeDownloads()
	readDone := make(chan int, 1)
	go func() { readDone <- fs.Read(path, make([]byte, 1), 0, fh) }()
	select {
	case n := <-readDone:
		t.Fatalf("read returned %d while a cache clear was running", n)
	case <-time.After(100 * time.Millisecond):
	}
	fs.resumeRangeDownloads()
	select {
	case n := <-readDone:
		if n != 1 {
			t.Fatalf("read after the cache clear returned %d, want 1", n)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("read did not continue after the cache clear")
	}
}

func TestUnmountFailsBlockedRead(t *testing.T) {
	const size = 4 << 20
	fs, path, _ := newCacheClearTestRemoteFs(t, false, size)
	backend, started, _ := blockFirstDownload(size)
	fs.backend = backend
	errno, fh := fs.Open(path, fuse.O_RDONLY)
	if errno != 0 {
		t.Fatal(errno)
	}
	defer fs.Release(path, fh)
	readDone := make(chan int, 1)
	go func() { readDone <- fs.Read(path, make([]byte, 1), 0, fh) }()
	<-started
	fs.closeRangeDownloads()
	select {
	case n := <-readDone:
		if n >= 0 {
			t.Fatalf("read during unmount returned %d, want an error", n)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("read kept waiting after unmount")
	}
}

// A read waiting on a download holds no lock that a rename or a directory
// listing needs, so neither waits for the network. Each stops the download
// instead: the read then fails for the renamed file, as it did before reads
// took that lock, and reads the version the listing reported for a changed one.
func TestChangesDoNotWaitForReadsBlockedOnDownloads(t *testing.T) {
	for _, change := range []string{"rename", "listing"} {
		t.Run(change, func(t *testing.T) {
			const size = 4 << 20
			fs, path, now := newCacheClearTestRemoteFs(t, false, size)
			backend, started, _ := blockFirstDownload(size)
			changedModTime := now.Add(time.Minute)
			backend.listForFunc = func(_ files_sdk.FolderListForParams, _ ...files_sdk.RequestResponseOption) (remoteFileIter, error) {
				return &fakeFileIter{files: []files_sdk.File{{DisplayName: "blocked", Type: "file", Size: size, Mtime: &changedModTime}}}, nil
			}
			fs.backend = backend
			errno, fh := fs.Open(path, fuse.O_RDONLY)
			if errno != 0 {
				t.Fatal(errno)
			}
			defer fs.Release(path, fh)
			readDone := make(chan int, 1)
			go func() { readDone <- fs.Read(path, make([]byte, 4096), size-4096, fh) }()
			<-started

			changeDone := make(chan int, 1)
			go func() {
				if change == "rename" {
					changeDone <- fs.Rename(path, "/renamed")
					return
				}
				root, _ := fs.vfs.fetch("/")
				root.expireInfo()
				root.childPathsExpires = timeZero
				changeDone <- fs.Readdir("/", func(string, *fuse.Stat_t, int64) bool { return true }, 0, ^uint64(0))
			}()
			select {
			case errc := <-changeDone:
				if errc != 0 {
					t.Fatalf("%s returned %d", change, errc)
				}
			case <-time.After(3 * time.Second):
				fs.cancelAllRangeDownloads()
				<-changeDone
				t.Fatalf("%s waited for a read blocked on a download", change)
			}
			select {
			case n := <-readDone:
				if change == "rename" && n >= 0 {
					t.Fatalf("read of a renamed file returned %d, want an error", n)
				}
				if change == "listing" && n != 4096 {
					t.Fatalf("read of a changed file returned %d, want 4096", n)
				}
			case <-time.After(3 * time.Second):
				t.Fatal("read kept waiting after the change stopped its download")
			}
		})
	}
}

// completeClaimingCache reports every entry complete, as if bytes were lost
// after hydration checked completeness.
type completeClaimingCache struct {
	cacheStore
}

func (c *completeClaimingCache) RangeEntryComplete(string, cache.EntryMetadata) (bool, error) {
	return true, nil
}

func TestRemoteFsEditBaselineRefusesCacheEntryWithHoles(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()
	defer fs.closeRangeDownloads()
	dc, err := disk.NewDiskCache(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(dc.StopMaintenance)
	fs.cacheStore = &completeClaimingCache{cacheStore: dc}
	path := "/holes"
	now := time.Now().Truncate(time.Second)
	meta := cacheEntryMetadata(path, 8, now)
	meta.ETag = `"fixture"`
	if _, err := dc.WriteRange(path, meta, []byte("AB"), 0); err != nil {
		t.Fatal(err)
	}
	node := vfs.getOrCreate(path, nodeTypeFile)
	node.updateInfo(fsNodeInfo{nodeType: nodeTypeFile, size: meta.Size, modTime: now})
	var uploads atomic.Int32
	fs.uploadWorkingCopy = func(_ context.Context, _ *fsNode, _ string, reader uploadWorkingCopyReader, mtime time.Time, _ uint64) (uploadedFileMetadata, error) {
		uploads.Add(1)
		data, err := io.ReadAll(reader)
		return testUploadedMetadata(int64(len(data)), mtime), err
	}
	errno, fh := fs.Open(path, fuse.O_RDWR)
	if errno != 0 {
		t.Fatal(errno)
	}
	defer fs.Release(path, fh)

	if n := fs.Write(path, []byte("Z"), 3, fh); n >= 0 {
		t.Fatalf("Write over a cache entry with holes returned %d, want an error", n)
	}
	_ = fs.Flush(path, fh)
	if got := uploads.Load(); got != 0 {
		t.Fatalf("uploaded %d times from an incomplete baseline", got)
	}
}

func TestRemoteFsWriteHydrationUsesCompleteDownloadWhenSparseReadsDisabled(t *testing.T) {
	for _, operation := range []string{"write", "truncate"} {
		for _, cachedPrefix := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/cached=%t", operation, cachedPrefix), func(t *testing.T) {
				fs, vfs, _ := newTestRemoteFs(t)
				defer vfs.destroy()
				defer fs.closeRangeDownloads()
				fs.sparseRangeReads = false
				logger := &captureMountLogger{}
				fs.log = logger
				dc, err := disk.NewDiskCache(t.TempDir())
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(dc.StopMaintenance)
				fs.cacheStore = dc
				path := "/editing"
				payload := []byte("ABCDEFGH")
				requests := make(chan string, 10)
				now := time.Now().Truncate(time.Second)
				var server *httptest.Server
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("ETag", `"fixture"`)
					if r.URL.Path != "/data" {
						w.Header().Set("Content-Type", "application/json")
						_, _ = fmt.Fprintf(w, `{"path":"editing","size":8,"mtime":%q,"download_uri":%q}`, now.Format(time.RFC3339), server.URL+"/data")
						return
					}
					requests <- r.Header.Get("Range")
					w.Header().Set("Content-Length", fmt.Sprint(len(payload)))
					_, _ = w.Write(payload)
				}))
				defer server.Close()
				meta := cacheEntryMetadata(path, int64(len(payload)), now)
				meta.ETag = `"fixture"`
				node := vfs.getOrCreate(path, nodeTypeFile)
				node.updateInfo(fsNodeInfo{nodeType: nodeTypeFile, size: meta.Size, modTime: now})
				node.setDownloadURI(server.URL + "/data")
				fs.backend = newRangeTestSDKBackend(files_sdk.Config{DisableDirectTransfers: true, EndpointOverride: server.URL})
				if cachedPrefix {
					if _, err := dc.WriteRange(path, meta, payload[:2], 0); err != nil {
						t.Fatal(err)
					}
				}
				var uploaded []byte
				fs.uploadWorkingCopy = func(_ context.Context, _ *fsNode, _ string, reader uploadWorkingCopyReader, mtime time.Time, _ uint64) (uploadedFileMetadata, error) {
					var err error
					uploaded, err = io.ReadAll(reader)
					return testUploadedMetadata(int64(len(uploaded)), mtime), err
				}
				errno, fh := fs.Open(path, fuse.O_RDWR)
				if errno != 0 {
					t.Fatal(errno)
				}
				defer fs.Release(path, fh)
				want := "ABCZEFGH"
				if operation == "write" {
					if n := fs.Write(path, []byte("Z"), 3, fh); n != 1 {
						t.Fatalf("Write returned %d: %s", n, logger.allJoined())
					}
				} else {
					want = "ABC"
					if errno := fs.Truncate(path, 3, fh); errno != 0 {
						t.Fatalf("Truncate returned %d: %s", errno, logger.allJoined())
					}
				}
				if errno := fs.Flush(path, fh); errno != 0 {
					t.Fatalf("Flush returned %d", errno)
				}
				if string(uploaded) != want {
					t.Fatalf("uploaded %q, want %q", uploaded, want)
				}
				if len(requests) != 1 {
					t.Fatalf("got %d requests, want one full download", len(requests))
				}
				if header := <-requests; header != "" {
					t.Fatalf("sparse disabled but sent Range: %s", header)
				}
			})
		}
	}
}
