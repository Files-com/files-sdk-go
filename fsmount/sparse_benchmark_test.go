//go:build linux || windows

package fsmount

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/winfsp/cgofuse/fuse"
)

const (
	syntheticBenchmarkFileSize = int64(128 * 1024 * 1024)
	syntheticBenchmarkReadSize = int64(4 * 1024 * 1024)
	syntheticBenchmarkChunk    = int64(256 * 1024)
)

type syntheticNetworkProfile struct {
	name               string
	requestDelay       time.Duration
	bytesPerSecond     int64
	aggregateBandwidth bool
}

type syntheticRangeServer struct {
	profile syntheticNetworkProfile
	size    int64
	server  *httptest.Server

	requestCount     atomic.Int64
	requestBytes     atomic.Int64
	transferredBytes atomic.Int64
	active           atomic.Int64
	maxActive        atomic.Int64

	rateMu sync.Mutex
	next   time.Time
}

type syntheticRangeServerStats struct {
	requests    int64
	bytes       int64
	transferred int64
	maxActive   int64
}

func newSyntheticRangeServer(size int64, profile syntheticNetworkProfile) *syntheticRangeServer {
	rangeServer := &syntheticRangeServer{size: size, profile: profile}
	rangeServer.server = httptest.NewServer(rangeServer)
	return rangeServer
}

func (s *syntheticRangeServer) Close() {
	s.server.Close()
}

func (s *syntheticRangeServer) URL() string {
	return s.server.URL
}

func (s *syntheticRangeServer) reset() {
	s.requestCount.Store(0)
	s.requestBytes.Store(0)
	s.transferredBytes.Store(0)
	s.active.Store(0)
	s.maxActive.Store(0)
	s.rateMu.Lock()
	s.next = time.Time{}
	s.rateMu.Unlock()
}

func (s *syntheticRangeServer) stats() syntheticRangeServerStats {
	return syntheticRangeServerStats{
		requests:    s.requestCount.Load(),
		bytes:       s.requestBytes.Load(),
		transferred: s.transferredBytes.Load(),
		maxActive:   s.maxActive.Load(),
	}
}

func (s *syntheticRangeServer) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("ETag", `"synthetic"`)
	start, end, partial, err := syntheticRequestRange(request.Header.Get("Range"), s.size)
	if err != nil {
		response.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	if s.profile.requestDelay > 0 {
		timer := time.NewTimer(s.profile.requestDelay)
		select {
		case <-timer.C:
		case <-request.Context().Done():
			if !timer.Stop() {
				<-timer.C
			}
			return
		}
	}

	length := end - start + 1
	s.requestCount.Add(1)
	s.requestBytes.Add(length)
	active := s.active.Add(1)
	defer s.active.Add(-1)
	for {
		maximum := s.maxActive.Load()
		if active <= maximum || s.maxActive.CompareAndSwap(maximum, active) {
			break
		}
	}

	response.Header().Set("Accept-Ranges", "bytes")
	response.Header().Set("Content-Length", fmt.Sprintf("%d", length))
	if partial {
		response.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, s.size))
		response.WriteHeader(http.StatusPartialContent)
	}

	buffer := make([]byte, syntheticBenchmarkChunk)
	perConnectionNext := time.Time{}
	for remaining := length; remaining > 0; {
		chunk := min(remaining, int64(len(buffer)))
		if s.profile.bytesPerSecond > 0 {
			if s.profile.aggregateBandwidth {
				if !s.waitForAggregateBandwidth(request.Context(), chunk) {
					return
				}
			} else {
				var completed bool
				perConnectionNext, completed = waitForSyntheticBandwidth(request.Context(), perConnectionNext, chunk, s.profile.bytesPerSecond)
				if !completed {
					return
				}
			}
		}
		written, err := response.Write(buffer[:chunk])
		s.transferredBytes.Add(int64(written))
		if err != nil {
			return
		}
		remaining -= chunk
	}
}

func syntheticRequestRange(header string, size int64) (start, end int64, partial bool, err error) {
	if header == "" {
		return 0, size - 1, false, nil
	}
	if _, err = fmt.Sscanf(header, "bytes=%d-%d", &start, &end); err != nil {
		return 0, 0, false, err
	}
	if start < 0 || end < start || end >= size {
		return 0, 0, false, fmt.Errorf("invalid range %q for size %d", header, size)
	}
	return start, end, true, nil
}

func (s *syntheticRangeServer) waitForAggregateBandwidth(ctx context.Context, bytes int64) bool {
	s.rateMu.Lock()
	now := time.Now()
	if s.next.Before(now) {
		s.next = now
	}
	start := s.next
	end := start.Add(time.Duration(float64(time.Second) * float64(bytes) / float64(s.profile.bytesPerSecond)))
	s.next = end
	s.rateMu.Unlock()

	timer := time.NewTimer(time.Until(start))
	select {
	case <-timer.C:
		return true
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
		s.rateMu.Lock()
		if s.next.Equal(end) {
			s.next = start
		}
		s.rateMu.Unlock()
		return false
	}
}

func waitForSyntheticBandwidth(ctx context.Context, next time.Time, bytes, bytesPerSecond int64) (time.Time, bool) {
	now := time.Now()
	if next.Before(now) {
		next = now
	}
	timer := time.NewTimer(time.Until(next))
	select {
	case <-timer.C:
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
		return next, false
	}
	return next.Add(time.Duration(float64(time.Second) * float64(bytes) / float64(bytesPerSecond))), true
}

func BenchmarkRemoteFsSyntheticSequentialRead(b *testing.B) {
	profiles := []syntheticNetworkProfile{
		{name: "loopback"},
		{
			name:               "fixed-aggregate-64MiBps-10ms",
			requestDelay:       10 * time.Millisecond,
			bytesPerSecond:     64 * 1024 * 1024,
			aggregateBandwidth: true,
		},
		{
			name:           "per-connection-32MiBps-10ms",
			requestDelay:   10 * time.Millisecond,
			bytesPerSecond: 32 * 1024 * 1024,
		},
	}
	policies := []struct {
		name   string
		sparse bool
	}{
		{name: "complete-plan", sparse: false},
		{name: "sparse-plan", sparse: true},
	}

	for _, profile := range profiles {
		b.Run(profile.name, func(b *testing.B) {
			for _, policy := range policies {
				b.Run(policy.name, func(b *testing.B) {
					rangeServer := newSyntheticRangeServer(syntheticBenchmarkFileSize, profile)
					defer rangeServer.Close()
					b.SetBytes(syntheticBenchmarkFileSize)

					var totalRequests int64
					var totalRemoteBytes int64
					var totalMaxActive int64
					for iteration := 0; iteration < b.N; iteration++ {
						b.StopTimer()
						rangeServer.reset()
						fs, vfs, _ := newTestRemoteFs(b)
						fs.sparseRangeReads = policy.sparse
						fs.backend = newRangeTestSDKBackend(files_sdk.Config{DisableDirectTransfers: true})
						vfs.cacheTTL = time.Hour

						path := fmt.Sprintf("/synthetic-%d.bin", iteration)
						modTime := time.Unix(1, 0)
						node := vfs.getOrCreate(path, nodeTypeFile)
						node.updateInfo(fsNodeInfo{
							nodeType:     nodeTypeFile,
							size:         syntheticBenchmarkFileSize,
							modTime:      modTime,
							creationTime: modTime,
						})
						node.setDownloadURI(rangeServer.URL())

						errno, handle := fs.Open(path, fuse.O_RDONLY)
						if errno != 0 {
							b.Fatalf("Open returned %d", errno)
						}

						buffer := make([]byte, syntheticBenchmarkReadSize)
						b.StartTimer()
						for offset := int64(0); offset < syntheticBenchmarkFileSize; offset += syntheticBenchmarkReadSize {
							if read := fs.Read(path, buffer, offset, handle); read != len(buffer) {
								b.Fatalf("Read at %d returned %d, want %d", offset, read, len(buffer))
							}
						}
						b.StopTimer()

						fs.Release(path, handle)
						fs.closeRangeDownloads()
						vfs.destroy()

						stats := rangeServer.stats()
						wantRequests := int64(1)
						if policy.sparse {
							// 128 MiB grows through 16, 32, 64, then the final 16 MiB.
							wantRequests = 4
						}
						if stats.requests != wantRequests {
							b.Fatalf("remote requests = %d, want %d", stats.requests, wantRequests)
						}
						if stats.bytes != syntheticBenchmarkFileSize {
							b.Fatalf("remote bytes = %d, want %d", stats.bytes, syntheticBenchmarkFileSize)
						}
						totalRequests += stats.requests
						totalRemoteBytes += stats.bytes
						totalMaxActive += stats.maxActive
						fs = nil
						vfs = nil
						buffer = nil
						runtime.GC()
					}

					b.ReportMetric(float64(totalRequests)/float64(b.N), "requests/op")
					b.ReportMetric(float64(totalRemoteBytes)/float64(b.N)/(1024*1024), "remote-MiB/op")
					b.ReportMetric(float64(totalMaxActive)/float64(b.N), "max-inflight/op")
				})
			}
		})
	}
}
