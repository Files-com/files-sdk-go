//go:build linux || windows

package fsmount

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/disk"
)

// rangeVersionTestSize is above the small-file cutoff, so reads use ranges.
const rangeVersionTestSize = 3 << 20

// serveRangeVersionTest answers a Range request with bytes of fill.
func serveRangeVersionTest(t *testing.T, w http.ResponseWriter, r *http.Request, etag string, status int, fill byte) {
	start, end, partial, err := syntheticRequestRange(r.Header.Get("Range"), rangeVersionTestSize)
	if err != nil {
		t.Error(err)
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}
	w.Header().Set("ETag", etag)
	w.Header().Set("Content-Length", fmt.Sprint(end-start+1))
	if partial {
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, rangeVersionTestSize))
		if status == http.StatusOK {
			status = http.StatusPartialContent
		}
	}
	w.WriteHeader(status)
	_, _ = w.Write(bytes.Repeat([]byte{fill}, int(end-start+1)))
}

func TestRangeDownloadResumesOnlyTheCachedRepresentation(t *testing.T) {
	for _, test := range []struct {
		name    string
		etag    string
		status  int
		changed bool
	}{
		{"same version", `"original"`, http.StatusPartialContent, false},
		{"changed despite If-Match", `"replacement"`, http.StatusPartialContent, true},
		{"precondition rejected", `"replacement"`, http.StatusPreconditionFailed, true},
		{"missing validator", "", http.StatusPartialContent, true},
		{"weak validator", `W/"original"`, http.StatusPartialContent, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			fs, vfs, _ := newTestRemoteFs(t)
			defer vfs.destroy()
			defer fs.closeRangeDownloads()
			root := t.TempDir()
			store, err := disk.NewDiskCache(root)
			if err != nil {
				t.Fatal(err)
			}
			meta := cacheEntryMetadata("/versioned", rangeVersionTestSize, time.Now().Truncate(time.Second))
			cached := meta
			cached.ETag = `"original"`
			if _, err := store.WriteRange(meta.Path, cached, bytes.Repeat([]byte("A"), 1<<20), 0); err != nil {
				t.Fatal(err)
			}
			if err := store.FlushRanges(meta.Path, cached); err != nil {
				t.Fatal(err)
			}
			// Reopening the cache must retain the validator, not only its byte ranges.
			store, err = disk.NewDiskCache(root)
			if err != nil {
				t.Fatal(err)
			}
			fs.cacheStore = store
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if got := r.Header.Get("If-Match"); got != cached.ETag {
					t.Errorf("If-Match = %q", got)
				}
				serveRangeVersionTest(t, w, r, test.etag, test.status, 'A')
			}))
			defer server.Close()
			node := vfs.getOrCreate(meta.Path, nodeTypeFile)
			node.updateInfo(fsNodeInfo{nodeType: nodeTypeFile, size: meta.Size, modTime: meta.ModTime})
			node.setDownloadURI(server.URL)
			fs.backend = newRangeTestSDKBackend(files_sdk.Config{DisableDirectTransfers: true})
			requested := cache.ByteRange{Start: 2 << 20, End: 2<<20 + 4}
			err = fs.rangeDownloads().Wait(context.Background(), 1, meta.Path, meta, requested)
			if test.changed {
				if !errors.Is(err, errRangeVersionChanged) {
					t.Fatalf("expected version error, got %v", err)
				}
				if node.getDownloadURI() != "" {
					t.Fatal("changed URL was retained")
				}
				if n, err := store.ReadRange(meta.Path, meta, make([]byte, 4), 0); n != 0 || err != nil {
					t.Fatalf("old prefix remained readable: %d, %v", n, err)
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				data := make([]byte, 4)
				if n, err := store.ReadRange(meta.Path, meta, data, requested.Start); n != 4 || err != nil || string(data) != "AAAA" {
					t.Fatalf("resumed bytes = %q, n=%d, err=%v", data, n, err)
				}
			}
		})
	}
}

func TestRangeDownloadWithoutStrongValidatorUsesOneFullResponse(t *testing.T) {
	for _, etag := range []string{"", `W/"weak"`} {
		t.Run(fmt.Sprintf("etag=%s", etag), func(t *testing.T) {
			fs, vfs, store := newTestRemoteFs(t)
			defer vfs.destroy()
			defer fs.closeRangeDownloads()
			var ranged, full atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fill := byte('B')
				if r.Header.Get("Range") != "" {
					ranged.Add(1)
					fill = 'A'
				} else {
					full.Add(1)
				}
				serveRangeVersionTest(t, w, r, etag, http.StatusOK, fill)
			}))
			defer server.Close()
			meta := cacheEntryMetadata("/unvalidated", rangeVersionTestSize, time.Now().Truncate(time.Second))
			node := vfs.getOrCreate(meta.Path, nodeTypeFile)
			node.updateInfo(fsNodeInfo{nodeType: nodeTypeFile, size: meta.Size, modTime: meta.ModTime})
			node.setDownloadURI(server.URL)
			fs.backend = newRangeTestSDKBackend(files_sdk.Config{DisableDirectTransfers: true})
			wanted := cache.ByteRange{Start: 2 << 20, End: 2<<20 + 4}
			if err := fs.rangeDownloads().Wait(context.Background(), 1, meta.Path, meta, wanted); err != nil {
				t.Fatal(err)
			}
			waitFor(t, func() bool { return !fs.rangeDownloads().hasStreams() })
			data := make([]byte, meta.Size)
			if n, err := store.ReadRange(meta.Path, meta, data, 0); int64(n) != meta.Size || err != nil || !bytes.Equal(data, bytes.Repeat([]byte("B"), int(meta.Size))) {
				t.Fatalf("full fallback n=%d, err=%v, all from the full response=%t", n, err, bytes.Equal(data, bytes.Repeat([]byte("B"), int(meta.Size))))
			}
			if ranged.Load() != 1 || full.Load() != 1 {
				t.Fatalf("range=%d full=%d", ranged.Load(), full.Load())
			}
		})
	}
}

// Two first requests sent before any validator is known must not both publish:
// whichever response is bound first defines the version for the whole file.
func TestConcurrentFirstRangesCannotPublishDifferentVersions(t *testing.T) {
	source := newFakeStreamSource(8 << 20)
	release := make(chan struct{})
	h := newStreamHarness(t, source, testStreamPolicy())
	h.streams.source = rangeSourceFunc(func(ctx context.Context, path string, meta cache.EntryMetadata, wanted cache.ByteRange) (remoteRangeResponse, error) {
		if wanted.Start > 0 {
			select {
			case <-release:
			case <-ctx.Done():
				return remoteRangeResponse{}, ctx.Err()
			}
			source.setETag(`"second"`)
		}
		return source.DownloadRange(ctx, path, meta, wanted)
	})
	tail := make(chan error, 1)
	go func() { tail <- h.read(2, 6<<20, 4096) }()
	waitFor(t, func() bool { return h.streams.diagnosticsSnapshot().ActiveRequests == 1 })
	h.mustRead(1, 0, 4096)
	close(release)

	select {
	case err := <-tail:
		if !errors.Is(err, errRangeVersionChanged) {
			t.Fatalf("tail read = %v, want %v", err, errRangeVersionChanged)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("tail read did not finish")
	}
	if cached := h.cachedBytes(); cached != 0 {
		t.Fatalf("%d bytes stayed cached after two versions were seen", cached)
	}
}

type rangeSourceFunc func(ctx context.Context, path string, meta cache.EntryMetadata, requested cache.ByteRange) (remoteRangeResponse, error)

func (f rangeSourceFunc) DownloadRange(ctx context.Context, path string, meta cache.EntryMetadata, requested cache.ByteRange) (remoteRangeResponse, error) {
	return f(ctx, path, meta, requested)
}
