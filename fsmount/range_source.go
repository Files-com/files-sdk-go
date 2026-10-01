//go:build linux || windows

package fsmount

import (
	"context"
	"errors"
	"io"
	"net/http"
	"slices"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/fsmount/events"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/disk"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/mem"
	lim "github.com/Files-com/files-sdk-go/v3/fsmount/internal/limit"
	"github.com/winfsp/cgofuse/fuse"
)

type remoteFsRangeSource struct {
	fs *RemoteFs
}

func (s *remoteFsRangeSource) DownloadRange(ctx context.Context, path string, meta cache.EntryMetadata, requested cache.ByteRange) (remoteRangeResponse, error) {
	response, err := s.downloadWithLimit(ctx, func(ctx context.Context) (remoteRangeResponse, error) {
		params := files_sdk.FileDownloadParams{File: files_sdk.File{
			Path:  s.fs.remotePath(path),
			Size:  meta.Size,
			Mtime: &meta.ModTime,
		}}
		if node, ok := s.fs.vfs.fetch(path); ok {
			params.File.DownloadUri = node.getDownloadURI()
		}

		response, err := s.fs.backend.downloadRange(params, requested, files_sdk.WithContext(ctx), rangeVersionOption(meta))
		if errors.Is(err, errRangeDownloadUnsupported) {
			response, err = s.downloadComplete(ctx, params, meta, requested)
		}
		return response, err
	})
	if err != nil {
		return remoteRangeResponse{}, err
	}
	return response, nil
}

func (s *remoteFsRangeSource) DownloadComplete(ctx context.Context, path string, meta cache.EntryMetadata, requested cache.ByteRange) (remoteRangeResponse, error) {
	response, err := s.downloadWithLimit(ctx, func(ctx context.Context) (remoteRangeResponse, error) {
		params := files_sdk.FileDownloadParams{File: files_sdk.File{
			Path:  s.fs.remotePath(path),
			Size:  meta.Size,
			Mtime: &meta.ModTime,
		}}
		if node, ok := s.fs.vfs.fetch(path); ok {
			params.File.DownloadUri = node.getDownloadURI()
		}

		return s.downloadComplete(ctx, params, meta, requested)
	})
	if err != nil {
		return remoteRangeResponse{}, err
	}
	return response, nil
}

func (s *remoteFsRangeSource) downloadWithLimit(ctx context.Context, download func(context.Context) (remoteRangeResponse, error)) (remoteRangeResponse, error) {
	ctx, release, err := s.fs.ops.Acquire(ctx, lim.FuseOpDownload)
	if err != nil {
		return remoteRangeResponse{}, err
	}
	response, err := download(ctx)
	if err != nil || response.Body == nil {
		if response.Body != nil {
			_ = response.Body.Close()
		}
		release()
		return response, err
	}
	response.Body = &limitedDownloadBody{ReadCloser: response.Body, release: release}
	return response, nil
}

type limitedDownloadBody struct {
	io.ReadCloser
	release func()
}

func (b *limitedDownloadBody) Close() error {
	defer b.release()
	return b.ReadCloser.Close()
}

func (s *remoteFsRangeSource) downloadComplete(ctx context.Context, params files_sdk.FileDownloadParams, meta cache.EntryMetadata, requested cache.ByteRange) (remoteRangeResponse, error) {
	result := remoteRangeResponse{TotalSize: -1}
	file, err := s.fs.backend.download(
		params,
		files_sdk.WithContext(ctx),
		rangeVersionOption(meta),
		files_sdk.ResponseOption(func(response *http.Response) error {
			return captureRangeResponse(&result, params.File.Size, requested, response)
		}),
	)
	result.File = file
	if err != nil {
		if result.Body != nil {
			_ = result.Body.Close()
		}
		return remoteRangeResponse{}, err
	}
	return result, nil
}

func rangeVersionOption(meta cache.EntryMetadata) files_sdk.RequestResponseOption {
	headers := &http.Header{}
	if meta.ETag != "" {
		headers.Set("If-Match", meta.ETag)
	}
	return files_sdk.RequestHeadersOption(headers)
}

func (fs *RemoteFs) rangeDownloads() *rangeStreams {
	fs.rangesMu.Lock()
	defer fs.rangesMu.Unlock()
	if fs.rangesClosed || fs.rangesPaused {
		return nil
	}
	if fs.ranges == nil {
		policy := defaultStreamPolicy(fs.sparseRangeReads, cacheCapacity(fs.cacheStore))
		fs.ranges = newRangeStreams(fs.cacheStore, &remoteFsRangeSource{fs: fs}, policy)
		fs.ranges.responseAccepted = func(path, uri string) {
			if uri != "" {
				if node, ok := fs.vfs.fetch(path); ok {
					node.setDownloadURI(uri)
				}
			}
		}
		fs.ranges.versionChanged = func(path string) {
			if node, ok := fs.vfs.fetch(path); ok {
				node.setDownloadURI("")
				node.expireInfo()
			}
		}
		if fs.log != nil {
			fs.ranges.logger = fs.log
		}
		fs.ranges.newReporter = func(path string, meta cache.EntryMetadata) *transferReporter {
			return fs.newTransferReporter(events.TransferDirectionDownload, path, meta.Size)
		}
	}
	return fs.ranges
}

// activeRangeDownloads returns the downloads without starting them.
func (fs *RemoteFs) activeRangeDownloads() *rangeStreams {
	fs.rangesMu.Lock()
	defer fs.rangesMu.Unlock()
	return fs.ranges
}

func cacheCapacity(store cacheStore) int64 {
	switch store := store.(type) {
	case *disk.DiskCache:
		return store.Capacity
	case *mem.MemoryCache:
		return store.Capacity
	}
	return 0
}

func (fs *RemoteFs) closeRangeDownloads() {
	fs.rangesMu.Lock()
	fs.rangesClosed = true
	ranges := fs.ranges
	fs.ranges = nil
	fs.rangesMu.Unlock()
	if ranges != nil {
		ranges.Close()
	}
}

// Pause rejects new downloads and retires the current coordinator before a
// cache clear waits for readers. Existing readers are woken by cancellation.
func (fs *RemoteFs) pauseRangeDownloads() {
	fs.rangesMu.Lock()
	fs.rangesPaused = true
	if fs.rangesResumed == nil {
		fs.rangesResumed = make(chan struct{})
	}
	ranges := fs.ranges
	fs.ranges = nil
	fs.rangesMu.Unlock()
	if ranges != nil {
		ranges.Close()
	}
}

func (fs *RemoteFs) resumeRangeDownloads() {
	fs.rangesMu.Lock()
	fs.rangesPaused = false
	if fs.rangesResumed != nil {
		close(fs.rangesResumed)
		fs.rangesResumed = nil
	}
	fs.rangesMu.Unlock()
}

// errnoRetryAfterCacheClear is returned by an operation whose download a
// cache clear stopped. It never reaches the file system; see
// retryAfterCacheClear.
const errnoRetryAfterCacheClear = -(1 << 30)

// retryAfterCacheClear runs op again after each cache clear that stopped a
// download op needed, so the application sees a slower operation instead of an
// I/O error. op must return with its locks released, because the clear waits
// for them before it removes cached files.
func (fs *RemoteFs) retryAfterCacheClear(op func() int) int {
	for {
		errc := op()
		if errc != errnoRetryAfterCacheClear {
			return errc
		}
		if !fs.waitForCacheClear() {
			return -fuse.EIO
		}
	}
}

// waitForCacheClear waits for a cache clear to resume downloads. It reports
// false when downloads were closed for unmount instead.
func (fs *RemoteFs) waitForCacheClear() bool {
	fs.rangesMu.Lock()
	resumed := fs.rangesResumed
	fs.rangesMu.Unlock()
	if resumed != nil {
		<-resumed
	}
	fs.rangesMu.Lock()
	defer fs.rangesMu.Unlock()
	return !fs.rangesClosed
}

func (fs *RemoteFs) lockRangeRead(path string) func() {
	fs.rangeOperationsMu.RLock()
	fs.rangePathMutexes.RLock(path)
	return func() {
		fs.rangePathMutexes.RUnlock(path)
		fs.rangeOperationsMu.RUnlock()
	}
}

func (fs *RemoteFs) lockRangeMutation(paths ...string) func() {
	paths = append([]string(nil), paths...)
	slices.Sort(paths)
	paths = slices.Compact(paths)
	fs.rangeOperationsMu.RLock()
	for _, path := range paths {
		fs.rangePathMutexes.Lock(path)
	}
	return func() {
		for index := len(paths) - 1; index >= 0; index-- {
			fs.rangePathMutexes.Unlock(paths[index])
		}
		fs.rangeOperationsMu.RUnlock()
	}
}

func (fs *RemoteFs) cancelRangeDownloads(paths ...string) {
	fs.rangesMu.Lock()
	ranges := fs.ranges
	fs.rangesMu.Unlock()
	if ranges == nil {
		return
	}
	for _, path := range paths {
		ranges.CancelPath(path)
	}
}

func (fs *RemoteFs) cancelAllRangeDownloads() {
	fs.rangesMu.Lock()
	ranges := fs.ranges
	fs.rangesMu.Unlock()
	if ranges != nil {
		ranges.CancelAll()
	}
}

var _ rangeSource = (*remoteFsRangeSource)(nil)
var _ completeRangeSource = (*remoteFsRangeSource)(nil)
