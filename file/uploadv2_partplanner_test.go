package file

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/directory"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/lib/uploadchecksum"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploadV2PartPlannerKnownSizeS3(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")

	t.Run("tiny file uses a single final part", func(t *testing.T) {
		size := int64(1)
		offsets := collectUploadV2Offsets(t, part, &size)
		require.Len(t, offsets, 1)
		assert.Equal(t, int64(1), offsets[0].len)
	})

	t.Run("zero byte file uses a single empty final part", func(t *testing.T) {
		size := int64(0)
		offsets := collectUploadV2Offsets(t, part, &size)
		require.Len(t, offsets, 1)
		assert.Equal(t, int64(0), offsets[0].len)
	})

	t.Run("resume from end of known size stops", func(t *testing.T) {
		size := int64(1)
		iterator, ok := newUploadV2PartIterator(part, &size, 1, 1)
		require.True(t, ok)
		offset, next, index := iterator()
		assert.Equal(t, OffSet{}, offset)
		assert.Nil(t, next)
		assert.Equal(t, 1, index)
	})

	t.Run("exact five mib uses one part", func(t *testing.T) {
		size := 5 * uploadV2MiB
		offsets := collectUploadV2Offsets(t, part, &size)
		require.Len(t, offsets, 1)
		assert.Equal(t, 5*uploadV2MiB, offsets[0].len)
	})

	t.Run("two hundred fifty six mib uses balanced s3 parts for concurrency", func(t *testing.T) {
		size := int64(256) * uploadV2MiB
		offsets := collectUploadV2Offsets(t, part, &size)
		require.Len(t, offsets, 16)
		assert.Equal(t, 16*uploadV2MiB, offsets[0].len)
		assert.Equal(t, size, sumOffsets(offsets))
	})

	t.Run("one hundred gib uses the preferred sixty four mib part size", func(t *testing.T) {
		size := int64(100) * uploadV2GiB
		offsets := collectUploadV2Offsets(t, part, &size)
		require.Len(t, offsets, 1600)
		assert.Equal(t, 64*uploadV2MiB, offsets[0].len)
		assert.Equal(t, size, sumOffsets(offsets))
	})

	t.Run("five tib stays within s3 part count", func(t *testing.T) {
		size := 5 * uploadV2TiB
		offsets := collectUploadV2Offsets(t, part, &size)
		assert.LessOrEqual(t, len(offsets), int(uploadV2S3MaxPartCount))
		assert.Equal(t, 525*uploadV2MiB, offsets[0].len)
		assert.Equal(t, size, sumOffsets(offsets))
	})

	t.Run("larger than five tib falls back out of v2 planning", func(t *testing.T) {
		size := 5*uploadV2TiB + 1
		_, ok := newUploadV2PartIterator(part, &size, 0, 0)
		assert.False(t, ok)
	})
}

func TestUploadV2S3KnownSizePreferredPartSizeTiers(t *testing.T) {
	tests := []struct {
		name      string
		totalSize int64
		want      int64
	}{
		{name: "small stays sixteen mib", totalSize: uploadV2GiB - 1, want: 16 * uploadV2MiB},
		{name: "one gib moves to thirty two mib", totalSize: uploadV2GiB, want: 32 * uploadV2MiB},
		{name: "sixteen gib moves to sixty four mib", totalSize: 16 * uploadV2GiB, want: 64 * uploadV2MiB},
		{name: "five hundred twelve gib moves to one hundred twenty eight mib", totalSize: 512 * uploadV2GiB, want: 128 * uploadV2MiB},
		{name: "one tib moves to two hundred fifty six mib", totalSize: uploadV2TiB, want: 256 * uploadV2MiB},
		{name: "two tib moves to five hundred twelve mib", totalSize: 2 * uploadV2TiB, want: 512 * uploadV2MiB},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, s3KnownSizePreferredPartSize(test.totalSize))
		})
	}
}

func TestUploadV2PartPlannerKnownSizeTargetDefaults(t *testing.T) {
	size := int64(100) * uploadV2MiB

	tests := []struct {
		name    string
		uri     string
		wantLen int64
	}{
		{name: "default", uri: "https://uploads.example.com/upload/file", wantLen: 16 * uploadV2MiB},
		{name: "custom named host still uses default", uri: "https://gateway.example.com/upload/file", wantLen: 16 * uploadV2MiB},
		{name: "parallel non s3 defaults to default", uri: "https://files.example.com/upload/file", wantLen: 16 * uploadV2MiB},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			part := uploadV2TestPart(test.uri)
			offsets := collectUploadV2Offsets(t, part, &size)
			require.NotEmpty(t, offsets)
			assert.Equal(t, test.wantLen, offsets[0].len)
			assert.Equal(t, size, sumOffsets(offsets))
		})
	}
}

func TestUploadV2KnownSizeNonS3PreferredPartSizeTiers(t *testing.T) {
	tests := []struct {
		name      string
		target    TransferV2TargetClass
		totalSize int64
		want      int64
	}{
		{name: "default medium uses sixteen mib", target: uploadV2TargetDefault, totalSize: 7 * uploadV2GiB, want: 16 * uploadV2MiB},
		{name: "default large uses thirty two mib", target: uploadV2TargetDefault, totalSize: 8 * uploadV2GiB, want: 32 * uploadV2MiB},
		{name: "default huge uses sixty four mib", target: uploadV2TargetDefault, totalSize: 128 * uploadV2GiB, want: 64 * uploadV2MiB},
		{name: "custom follows default tiers", target: TransferV2TargetClass("custom"), totalSize: 128 * uploadV2GiB, want: 64 * uploadV2MiB},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			plan, ok, reason := newUploadV2PartPlan(test.target, &test.totalSize)
			require.True(t, ok, reason)
			assert.Equal(t, test.want, plan.partSize)
		})
	}
}

func TestUploadV2TargetClassifierCanReturnCustomTarget(t *testing.T) {
	part := uploadV2TestPart("https://uploads.example.com/upload/file")
	size := int64(100) * uploadV2MiB

	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size, func(files_sdk.FileUploadPart) TransferV2TargetClass {
		return "custom"
	})

	require.True(t, ok, reason)
	assert.Equal(t, TransferV2TargetClass("custom"), plan.target)
	assert.Equal(t, int64(16)*uploadV2MiB, plan.partSize)
}

func TestUploadV2PartPlannerUnknownSizeGrowth(t *testing.T) {
	t.Run("default caps at thirty two mib", func(t *testing.T) {
		plan, ok, _ := newUploadV2PartPlan(uploadV2TargetDefault, nil)
		require.True(t, ok)
		assert.Equal(t, 8*uploadV2MiB, plan.partSizeForIndex(0))
		assert.Equal(t, 8*uploadV2MiB, plan.partSizeForIndex(127))
		assert.Equal(t, 16*uploadV2MiB, plan.partSizeForIndex(128))
		assert.Equal(t, 32*uploadV2MiB, plan.partSizeForIndex(256))
		assert.Equal(t, 32*uploadV2MiB, plan.partSizeForIndex(384))
	})

	t.Run("custom caps at thirty two mib", func(t *testing.T) {
		plan, ok, _ := newUploadV2PartPlan(TransferV2TargetClass("custom"), nil)
		require.True(t, ok)
		assert.Equal(t, 8*uploadV2MiB, plan.partSizeForIndex(0))
		assert.Equal(t, 8*uploadV2MiB, plan.partSizeForIndex(127))
		assert.Equal(t, 16*uploadV2MiB, plan.partSizeForIndex(128))
		assert.Equal(t, 32*uploadV2MiB, plan.partSizeForIndex(256))
		assert.Equal(t, 32*uploadV2MiB, plan.partSizeForIndex(384))
	})

	t.Run("s3 unknown size uses v1 fallback", func(t *testing.T) {
		part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
		_, ok := newUploadV2PartIterator(part, nil, 0, 0)
		assert.False(t, ok)
	})

	t.Run("non parallel upload uses v1 fallback", func(t *testing.T) {
		size := int64(100) * uploadV2MiB
		part := files_sdk.FileUploadPart{
			UploadUri:     "https://uploads.example.com/upload/file",
			ParallelParts: lib.Bool(false),
		}
		_, ok := newUploadV2PartIterator(part, &size, 0, 0)
		assert.False(t, ok)
	})
}

func TestUploadV2ProgressBatcherFlushesExactBytes(t *testing.T) {
	var progress int64
	batcher := newUploadV2ProgressBatcher(func(delta int64) {
		progress += delta
	})

	batcher.Add(512 * 1024)
	assert.Equal(t, int64(0), progress)

	batcher.Add(600 * 1024)
	assert.Equal(t, int64(1112*1024), progress)

	batcher.Add(7)
	assert.Equal(t, int64(1112*1024), progress)

	batcher.Flush()
	assert.Equal(t, int64(1112*1024+7), progress)
}

func TestUploadV2ProgressBatcherRewindCancelsPendingAndFlushedBytes(t *testing.T) {
	var progress int64
	batcher := newUploadV2ProgressBatcher(func(delta int64) {
		progress += delta
	})

	batcher.Add(512 * 1024)
	batcher.Add(-512 * 1024)
	assert.Equal(t, int64(0), progress)

	batcher.Add(2 * uploadV2MiB)
	assert.Equal(t, 2*uploadV2MiB, progress)

	batcher.Add(512 * 1024)
	batcher.Add(-(2*uploadV2MiB + 512*1024))
	assert.Equal(t, int64(0), progress)
}

func TestUploadV2PartPlannerIsGated(t *testing.T) {
	size := int64(100) * uploadV2MiB
	v1 := uploadIO{
		ByteOffset: ByteOffset{PartSizes: lib.PartSizes},
		Size:       &size,
		FileUploadPart: files_sdk.FileUploadPart{
			UploadUri:     "https://uploads.example.com/upload/file",
			ParallelParts: lib.Bool(true),
		},
	}
	v1Offset, _, _ := v1.ByteOffset.Resume(v1.Size, 0, 0)()
	assert.Equal(t, 10*uploadV2MiB, v1Offset.len)

	v2Iterator, ok := newUploadV2PartIterator(v1.FileUploadPart, v1.Size, 0, 0)
	require.True(t, ok)
	v2Offset, _, _ := v2Iterator()
	assert.Equal(t, 16*uploadV2MiB, v2Offset.len)
}

func TestUploadWithV2EnablesAdaptiveConcurrencyAndTelemetry(t *testing.T) {
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	server.MockFiles["v2-observable.txt"] = mockFile{File: files_sdk.File{Size: 1}}

	var logs bytes.Buffer
	client := server.Client()
	client.Config.Logger = log.New(&logs, "", 0)

	_, err := client.UploadWithResume(
		UploadWithV2(),
		UploadWithReaderAt(bytes.NewReader([]byte("x"))),
		UploadWithDestinationPath("v2-observable.txt"),
		UploadWithSize(1),
		UploadWithManager(lib.NewConstrainedWorkGroup(75)),
	)

	require.NoError(t, err)
	logText := logs.String()
	assert.Contains(t, logText, "event: upload v2 enabled")
	assert.Contains(t, logText, "target_class: default")
	assert.Contains(t, logText, "part_size_mode: known_size")
	assert.Contains(t, logText, "adaptive_initial_target: 8")
	assert.Contains(t, logText, "adaptive_max_target: 75")
	assert.Contains(t, logText, "ready_runway_parts: 4")
	assert.Contains(t, logText, "ready_runway_bytes: 268435456")
	assert.Contains(t, logText, "upload_http_client_adjusted: false")
	assert.Contains(t, logText, "event: upload v2 complete")
	assert.Contains(t, logText, "adaptive_success_total: 1")
	assert.Contains(t, logText, "read_duration_ms:")
	assert.Contains(t, logText, "read_duration_ns:")
	assert.Contains(t, logText, "read_throughput_bytes_per_second:")
	assert.Contains(t, logText, "success: true")

	uploadRequests := server.TrackRequest["/upload/*path"]
	require.Len(t, uploadRequests, 1)
	assert.Contains(t, uploadRequests[0], "part_number=1")
	assert.Contains(t, uploadRequests[0], "part_offset=0")
}

func TestUploadWithV2ReadyRunwayOverridesDefaults(t *testing.T) {
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	server.MockFiles["v2-runway.txt"] = mockFile{File: files_sdk.File{Size: 1}}

	var logs bytes.Buffer
	client := server.Client()
	client.Config.Logger = log.New(&logs, "", 0)

	_, err := client.UploadWithResume(
		UploadWithV2(),
		UploadWithV2ReadyRunway(2, 64*uploadV2MiB),
		UploadWithReaderAt(bytes.NewReader([]byte("x"))),
		UploadWithDestinationPath("v2-runway.txt"),
		UploadWithSize(1),
		UploadWithManager(lib.NewConstrainedWorkGroup(4)),
	)

	require.NoError(t, err)
	logText := logs.String()
	assert.Contains(t, logText, "ready_runway_parts: 2")
	assert.Contains(t, logText, "ready_runway_bytes: 67108864")
}

func TestUploadWithV2KeepsStatusQueuedUntilAdaptivePartSlot(t *testing.T) {
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	server.MockFiles["v2-queued.txt"] = mockFile{File: files_sdk.File{Size: 1}}

	client := server.Client()
	job := (&Job{
		Config: client.Config,
		Logger: client.Config.Logger,
	}).Init()
	uploadStatus := &UploadStatus{
		Mutex:      &sync.RWMutex{},
		status:     status.Queued,
		job:        job,
		localPath:  "v2-queued.txt",
		remotePath: "v2-queued.txt",
		file: files_sdk.File{
			Path: "v2-queued.txt",
			Size: 1,
		},
	}
	blockedManager := lib.NewAdaptiveConcurrencyManagerWithConfig(lib.AdaptiveConcurrencyConfig{
		MaxConcurrency: 1,
		InitialTarget:  1,
		MinTarget:      1,
	})
	blockedManager.Wait()
	released := false
	done := make(chan error, 1)
	defer func() {
		if !released {
			released = true
			blockedManager.DoneNeutral()
			select {
			case <-done:
			case <-time.After(time.Second):
			}
		}
	}()

	go func() {
		_, err := client.UploadWithResume(
			UploadWithV2(),
			UploadWithReaderAt(bytes.NewReader([]byte("x"))),
			UploadWithDestinationPath("v2-queued.txt"),
			UploadWithSize(1),
			UploadWithManager(lib.NewConstrainedWorkGroup(2)),
			UploadWithProgress(uploadProgress(uploadStatus)),
			uploadWithV2AdaptiveManagerProvider(func(uploadV2PartPlan, int, UploadV2Tuning) *lib.AdaptiveConcurrencyManager {
				return blockedManager
			}),
		)
		done <- err
	}()

	trackedRequestCount := func(route string) int {
		server.traceMutex.Lock()
		defer server.traceMutex.Unlock()
		return len(server.TrackRequest[route])
	}
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) && trackedRequestCount("/api/rest/v1/file_actions/begin_upload/*path") == 0 {
		time.Sleep(10 * time.Millisecond)
	}
	require.NotZero(t, trackedRequestCount("/api/rest/v1/file_actions/begin_upload/*path"))

	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, status.Queued, uploadStatus.Status())
	assert.Zero(t, trackedRequestCount("/upload/*path"))

	released = true
	blockedManager.DoneNeutral()
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for upload v2 to finish")
	}
	assert.Equal(t, status.Uploading, uploadStatus.Status())
	assert.Equal(t, 1, trackedRequestCount("/upload/*path"))
}

func TestUploadWithV2ReadyRunwayRejectsNegativeValues(t *testing.T) {
	_, err := UploadWithV2ReadyRunway(-1, 0)(uploadIO{})
	require.Error(t, err)

	_, err = UploadWithV2ReadyRunway(1, -1)(uploadIO{})
	require.Error(t, err)
}

func TestUploadV2FeatureFlagAloneDoesNotUpgradeExistingSDKUploads(t *testing.T) {
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	server.MockFiles["v2-feature-flag.txt"] = mockFile{File: files_sdk.File{Size: 1}}

	var logs bytes.Buffer
	client := server.Client()
	client.Config.Logger = log.New(&logs, "", 0)
	client.Config.FeatureFlags[files_sdk.FeatureFlagAdaptiveUploadV2] = true

	_, err := client.UploadWithResume(
		UploadWithReaderAt(bytes.NewReader([]byte("x"))),
		UploadWithDestinationPath("v2-feature-flag.txt"),
		UploadWithSize(1),
		UploadWithManager(lib.NewConstrainedWorkGroup(9)),
	)

	require.NoError(t, err)
	logText := logs.String()
	assert.NotContains(t, logText, "event: upload v2 enabled")
	assert.NotContains(t, logText, "adaptive_max_target:")

	uploadRequests := server.TrackRequest["/upload/*path"]
	require.Len(t, uploadRequests, 1)
	assert.NotContains(t, uploadRequests[0], "part_offset=0")
}

func TestUploadV2ChecksumTrailerFeatureFlagUsesAWSChunkedForSignedS3URL(t *testing.T) {
	var decoded bytes.Buffer
	expectedAlgorithm := uploadV2ChecksumTrailerAlgorithm.OrBestForPlatform()
	expectedTrailerHeader, err := expectedAlgorithm.TrailerHeader()
	require.NoError(t, err)
	transport := uploadV2RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, uploadchecksum.ContentEncodingAWSChunked, req.Header.Get("Content-Encoding"))
		assert.Equal(t, uploadchecksum.ContentSHA256StreamingUnsignedTrailer, req.Header.Get("x-amz-content-sha256"))
		assert.Equal(t, string(expectedAlgorithm), req.Header.Get("x-amz-sdk-checksum-algorithm"))
		assert.Equal(t, expectedTrailerHeader, req.Header.Get("x-amz-trailer"))
		assert.Equal(t, "5", req.Header.Get("x-amz-decoded-content-length"))
		assert.NotEqual(t, "5", req.Header.Get("Content-Length"))
		assert.Equal(t, req.ContentLength, mustParseInt64(t, req.Header.Get("Content-Length")))

		decoder, err := uploadchecksum.NewAWSChunkedDecoderForHeaders(req.Body, req.Header)
		require.NoError(t, err)
		_, err = io.Copy(&decoded, decoder)
		require.NoError(t, err)
		assert.NotEmpty(t, decoder.TrailerValue())

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Etag": []string{"checksum-etag"}},
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    req,
		}, nil
	})
	part := uploadV2TestPart(uploadV2ChecksumTrailerSignedS3URL())
	part.HttpMethod = "PUT"
	part.PartNumber = 1
	config := files_sdk.Config{}.Init().SetCustomClient(&http.Client{Transport: transport})
	config.FeatureFlags[files_sdk.FeatureFlagUploadV2ChecksumTrailer] = true
	engine := newUploadV2TestEngineWithClient(t, part, &Client{Config: config})
	engine.u.Progress = func(int64) {}

	uploadPart := &uploadV2Part{
		uploadV2PartDescriptor: uploadV2PartDescriptor{
			number: 1,
			offset: OffSet{off: 0, len: 5},
			upload: part,
			legacy: &Part{},
		},
		reader: &ProxyReaderAt{
			ReaderAt: bytes.NewReader([]byte("hello")),
			off:      0,
			len:      5,
			onRead:   func(int64) {},
		},
	}

	etag, _, _, _, err := engine.uploadPart(context.Background(), uploadPart)

	require.NoError(t, err)
	assert.Equal(t, files_sdk.EtagsParam{Part: "1", Etag: "checksum-etag"}, etag)
	assert.Equal(t, "hello", decoded.String())
}

func TestUploadV2ChecksumTrailerFeatureFlagSkipsUnsignedS3URL(t *testing.T) {
	var raw bytes.Buffer
	transport := uploadV2RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		assert.Empty(t, req.Header.Get("Content-Encoding"))
		assert.Empty(t, req.Header.Get("x-amz-content-sha256"))
		assert.Empty(t, req.Header.Get("x-amz-decoded-content-length"))
		assert.Empty(t, req.Header.Get("x-amz-trailer"))
		assert.Equal(t, "5", req.Header.Get("Content-Length"))
		_, err := io.Copy(&raw, req.Body)
		require.NoError(t, err)

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Etag": []string{"raw-etag"}},
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    req,
		}, nil
	})
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1&X-Amz-SignedHeaders=host")
	part.HttpMethod = "PUT"
	part.PartNumber = 1
	config := files_sdk.Config{}.Init().SetCustomClient(&http.Client{Transport: transport})
	config.FeatureFlags[files_sdk.FeatureFlagUploadV2ChecksumTrailer] = true
	engine := newUploadV2TestEngineWithClient(t, part, &Client{Config: config})
	engine.u.Progress = func(int64) {}

	uploadPart := &uploadV2Part{
		uploadV2PartDescriptor: uploadV2PartDescriptor{
			number: 1,
			offset: OffSet{off: 0, len: 5},
			upload: part,
			legacy: &Part{},
		},
		reader: &ProxyReaderAt{
			ReaderAt: bytes.NewReader([]byte("hello")),
			off:      0,
			len:      5,
			onRead:   func(int64) {},
		},
	}

	etag, _, _, _, err := engine.uploadPart(context.Background(), uploadPart)

	require.NoError(t, err)
	assert.Equal(t, files_sdk.EtagsParam{Part: "1", Etag: "raw-etag"}, etag)
	assert.Equal(t, "hello", raw.String())
}

func TestUploadV2ChecksumTrailerFeatureFlagSkipsUnsupportedDestination(t *testing.T) {
	var raw bytes.Buffer
	transport := uploadV2RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		assert.Empty(t, req.Header.Get("Content-Encoding"))
		assert.Empty(t, req.Header.Get("x-amz-content-sha256"))
		assert.Empty(t, req.Header.Get("x-amz-decoded-content-length"))
		assert.Empty(t, req.Header.Get("x-amz-trailer"))
		assert.Equal(t, "5", req.Header.Get("Content-Length"))
		_, err := io.Copy(&raw, req.Body)
		require.NoError(t, err)

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Etag": []string{"raw-etag"}},
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    req,
		}, nil
	})
	part := uploadV2TestPart("https://uploads.example.com/upload/key?part_number=1")
	part.HttpMethod = "POST"
	part.PartNumber = 1
	config := files_sdk.Config{}.Init().SetCustomClient(&http.Client{Transport: transport})
	config.FeatureFlags[files_sdk.FeatureFlagUploadV2ChecksumTrailer] = true
	engine := newUploadV2TestEngineWithClient(t, part, &Client{Config: config})
	engine.u.Progress = func(int64) {}

	uploadPart := &uploadV2Part{
		uploadV2PartDescriptor: uploadV2PartDescriptor{
			number: 1,
			offset: OffSet{off: 0, len: 5},
			upload: part,
			legacy: &Part{},
		},
		reader: &ProxyReaderAt{
			ReaderAt: bytes.NewReader([]byte("hello")),
			off:      0,
			len:      5,
			onRead:   func(int64) {},
		},
	}

	etag, _, _, _, err := engine.uploadPart(context.Background(), uploadPart)

	require.NoError(t, err)
	assert.Equal(t, files_sdk.EtagsParam{Part: "1", Etag: "raw-etag"}, etag)
	assert.Equal(t, "hello", raw.String())
	assert.Equal(t, "unsupported_destination", engine.checksumTrailerSkipReason)
}

func TestUploadV2CancelDuringFinalizeReturnsResumableQuickly(t *testing.T) {
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	server.MockFiles["v2-cancel-finalize.txt"] = mockFile{File: files_sdk.File{Size: 1}}
	server.MockRoute("/api/rest/v1/file_actions/begin_upload/v2-cancel-finalize.txt", func(ctx *gin.Context, model interface{}) bool {
		beginUpload := model.(files_sdk.FileBeginUploadParams)
		if beginUpload.Part == 0 {
			beginUpload.Part = 1
		}
		ctx.JSON(http.StatusOK, files_sdk.FileUploadPartCollection{
			files_sdk.FileUploadPart{
				HttpMethod:    "POST",
				Path:          beginUpload.Path,
				Ref:           "v2-cancel-finalize-ref",
				UploadUri:     server.Server.URL + "/upload/v2-cancel-finalize.txt?part_number=1",
				ParallelParts: lib.Bool(true),
				Expires:       time.Now().Add(time.Hour).Format(time.RFC3339),
				PartNumber:    beginUpload.Part,
			},
		})
		return true
	})

	finalizeStarted := make(chan struct{})
	var blockFinalize atomic.Bool
	var finalizeStartedOnce atomic.Bool
	blockFinalize.Store(true)
	server.MockRoute("/api/rest/v1/files/v2-cancel-finalize.txt", func(ctx *gin.Context, model interface{}) bool {
		if blockFinalize.Load() {
			if finalizeStartedOnce.CompareAndSwap(false, true) {
				close(finalizeStarted)
			}
			select {
			case <-ctx.Request.Context().Done():
				return true
			case <-time.After(2 * time.Second):
				ctx.JSON(http.StatusOK, files_sdk.File{Path: "v2-cancel-finalize.txt", Size: 1})
				return true
			}
		}
		ctx.JSON(http.StatusOK, files_sdk.File{Path: "v2-cancel-finalize.txt", Size: 1})
		return true
	})

	client := server.Client()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan struct{})
	var resumable UploadResumable
	var err error
	go func() {
		defer close(done)
		resumable, err = client.UploadWithResume(
			UploadWithV2(),
			UploadWithContext(ctx),
			UploadWithReaderAt(bytes.NewReader([]byte("x"))),
			UploadWithDestinationPath("v2-cancel-finalize.txt"),
			UploadWithSize(1),
			UploadWithManager(lib.NewConstrainedWorkGroup(2)),
		)
	}()

	select {
	case <-finalizeStarted:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for finalize to start")
	}
	start := time.Now()
	cancel()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for canceled finalize to return")
	}

	require.ErrorIs(t, err, context.Canceled)
	assert.Less(t, time.Since(start), 500*time.Millisecond)
	require.Len(t, resumable.Parts, 1)
	assert.True(t, resumable.Parts[0].Successful())
	assert.Equal(t, "v2-cancel-finalize-ref", resumable.FileUploadPart.Ref)

	blockFinalize.Store(false)
	resumable, err = client.UploadWithResume(
		UploadWithV2(),
		UploadWithReaderAt(bytes.NewReader([]byte("x"))),
		UploadWithDestinationPath("v2-cancel-finalize.txt"),
		UploadWithSize(1),
		UploadWithResume(resumable),
		UploadWithManager(lib.NewConstrainedWorkGroup(2)),
	)

	require.NoError(t, err)
	assert.Equal(t, int64(1), resumable.File.Size)
	uploadRequests := server.TrackRequest["/upload/*path"]
	require.Len(t, uploadRequests, 1)
	assert.Contains(t, uploadRequests[0], "part_offset=0")
}

func TestUploadV2FinalizeNotFoundInvalidatesResumeState(t *testing.T) {
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	server.MockFiles["v2-finalize-not-found.txt"] = mockFile{File: files_sdk.File{Size: 1}}
	server.MockRoute("/api/rest/v1/files/v2-finalize-not-found.txt", func(ctx *gin.Context, model interface{}) bool {
		ctx.JSON(http.StatusNotFound, files_sdk.ResponseError{
			Type:         string(files_sdk.ErrFileUploadNotFound),
			Title:        "File Upload Not Found",
			ErrorMessage: "File Upload for id stale-ref not found.",
		})
		return true
	})

	client := server.Client()
	resumable, err := client.UploadWithResume(
		UploadWithV2(),
		UploadWithReaderAt(bytes.NewReader([]byte("x"))),
		UploadWithDestinationPath("v2-finalize-not-found.txt"),
		UploadWithSize(1),
		UploadWithManager(lib.NewConstrainedWorkGroup(2)),
	)

	require.Error(t, err)
	assert.True(t, files_sdk.IsNotExist(err), err.Error())
	assert.Empty(t, resumable.Parts)
	assert.Empty(t, resumable.FileUploadPart.Ref)
}

func TestUploadWithV2FallsBackForNonParallelUpload(t *testing.T) {
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	server.MockFiles["v2-serial.txt"] = mockFile{File: files_sdk.File{Size: 1}}

	server.MockRoute("/api/rest/v1/file_actions/begin_upload/v2-serial.txt", func(ctx *gin.Context, model interface{}) bool {
		beginUpload := model.(files_sdk.FileBeginUploadParams)
		if beginUpload.Part == 0 {
			beginUpload.Part = 1
		}
		ctx.JSON(http.StatusOK, files_sdk.FileUploadPartCollection{
			files_sdk.FileUploadPart{
				HttpMethod:    "POST",
				Path:          beginUpload.Path,
				UploadUri:     server.Server.URL + "/upload/v2-serial.txt?part_number=1",
				ParallelParts: lib.Bool(false),
				Expires:       time.Now().Add(time.Hour).Format(time.RFC3339),
				PartNumber:    beginUpload.Part,
			},
		})
		return true
	})

	var logs bytes.Buffer
	client := server.Client()
	client.Config.Logger = log.New(&logs, "", 0)

	_, err := client.UploadWithResume(
		UploadWithV2(),
		UploadWithReader(strings.NewReader("x")),
		UploadWithDestinationPath("v2-serial.txt"),
		UploadWithSize(1),
	)

	require.NoError(t, err)
	assert.Contains(t, logs.String(), "event: upload v2 fallback")
	assert.Contains(t, logs.String(), "reason: parallel_parts_disabled")
	assert.NotContains(t, logs.String(), "event: upload v2 complete")
}

func TestUploadV2PartOffsetQuery(t *testing.T) {
	part := files_sdk.FileUploadPart{
		UploadUri:     "https://uploads.example.com/upload/file?partNumber=7&offset=old",
		ParallelParts: lib.Bool(true),
		PartNumber:    7,
	}
	engine := newUploadV2TestEngine(t, part)

	engine.decorateUploadURL(&part, 7, 33554432)

	query := uploadV2TestQuery(t, part.UploadUri)
	assert.Equal(t, "7", query.Get("part_number"))
	assert.Equal(t, "33554432", query.Get("part_offset"))
	assert.Empty(t, query.Get("partNumber"))
	assert.Empty(t, query.Get("offset"))
}

func TestUploadV2PartOffsetQueryPreservesRawQueryAndFragment(t *testing.T) {
	got := uploadURLWithPartOffset(
		"https://uploads.example.com/upload/file?keep=a%2Fb&offset=old&partNumber=7&empty=&sig=abc%2Bdef#fragment",
		8,
		123456,
	)

	assert.Equal(t, "https://uploads.example.com/upload/file?keep=a%2Fb&empty=&sig=abc%2Bdef&part_number=8&part_offset=123456#fragment", got)
}

func TestUploadV2PartOffsetQueryAddsQueryToBareURL(t *testing.T) {
	got := uploadURLWithPartOffset("https://uploads.example.com/upload/file", 8, 123456)

	assert.Equal(t, "https://uploads.example.com/upload/file?part_number=8&part_offset=123456", got)
}

func TestUploadV2PartOffsetQuerySkipsS3(t *testing.T) {
	part := files_sdk.FileUploadPart{
		UploadUri:     "https://s3.amazonaws.com/bucket/key?partNumber=7&X-Amz-Signature=signed",
		ParallelParts: lib.Bool(true),
		PartNumber:    7,
	}
	engine := newUploadV2TestEngine(t, part)

	engine.decorateUploadURL(&part, 7, 33554432)

	query := uploadV2TestQuery(t, part.UploadUri)
	assert.Equal(t, "7", query.Get("partNumber"))
	assert.Empty(t, query.Get("part_offset"))
	assert.Equal(t, "signed", query.Get("X-Amz-Signature"))
}

func TestUploadV2EngineCreationDoesNotReplaceManagerBeforeRun(t *testing.T) {
	part := uploadV2TestPart("https://uploads.example.com/upload/file")
	size := int64(100)
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)
	originalManager := lib.NewConstrainedWorkGroup(7)
	u := &uploadIO{FileUploadPart: part, Manager: originalManager, managerSet: true}

	engine := newUploadV2Engine(u, plan)

	gotManager, ok := u.Manager.(*lib.ConstrainedWorkGroup)
	require.True(t, ok)
	assert.Same(t, originalManager, gotManager)
	assert.Same(t, originalManager, engine.globalManager)
}

func TestUploadV2S3UsesTargetDefaultConcurrencyWithoutExplicitManagerCap(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(256) * uploadV2MiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)

	engine := newUploadV2Engine(&uploadIO{FileUploadPart: part}, plan)

	assert.Equal(t, uploadV2S3MaxConcurrency, engine.manager.Max())
	assert.Equal(t, 16, engine.manager.Target())
}

func TestUploadV2RaisesDefaultHTTPTransportCapForSDKClient(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(256) * uploadV2MiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)
	client := &Client{Config: files_sdk.Config{}.Init()}
	originalTransport, ok := client.Config.Client.HTTPClient.Transport.(*lib.Transport)
	require.True(t, ok)
	require.Equal(t, 75, originalTransport.MaxConnsPerHost)

	engine := newUploadV2Engine(&uploadIO{
		Client:         client,
		FileUploadPart: part,
		Size:           &size,
	}, plan)

	adjustedTransport, ok := engine.u.Config.Client.HTTPClient.Transport.(*lib.Transport)
	require.True(t, ok)
	assert.NotSame(t, originalTransport, adjustedTransport)
	assert.NotSame(t, client, engine.u.Client)
	assert.Equal(t, 75, originalTransport.MaxConnsPerHost)
	assert.Equal(t, uploadV2S3GrowthCeiling, adjustedTransport.MaxConnsPerHost)
	assert.Equal(t, uploadV2DefaultHTTPIdleConnectionCap, adjustedTransport.MaxIdleConns)
	assert.Equal(t, uploadV2DefaultHTTPIdleConnectionCap, adjustedTransport.MaxIdleConnsPerHost)
	assert.True(t, engine.httpClientLimits.adjusted)
	assert.True(t, engine.httpClientLimits.available)
	assert.Equal(t, uploadV2S3GrowthCeiling, engine.httpClientLimits.maxConnsPerHost)
	assert.Equal(t, uploadV2DefaultHTTPIdleConnectionCap, engine.httpClientLimits.maxIdleConns)
	assert.Equal(t, uploadV2DefaultHTTPIdleConnectionCap, engine.httpClientLimits.maxIdleConnsPerHost)
	attrs := engine.uploadV2EnabledLogAttrs(uploadV2ReadyRunwayConfig{})
	assert.Equal(t, true, attrs["upload_http_client_adjusted"])
	assert.Equal(t, uploadV2S3GrowthCeiling, attrs["upload_max_conns_per_host"])
	assert.Equal(t, uploadV2DefaultHTTPIdleConnectionCap, attrs["upload_max_idle_conns"])
	assert.Equal(t, uploadV2DefaultHTTPIdleConnectionCap, attrs["upload_max_idle_conns_per_host"])
}

func TestUploadV2LowersPreconfiguredHTTPTransportCapToS3GrowthCeiling(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(200) * uploadV2MiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)
	client := &Client{Config: files_sdk.Config{}.Init()}
	originalTransport, ok := client.Config.Client.HTTPClient.Transport.(*lib.Transport)
	require.True(t, ok)
	originalTransport.MaxConnsPerHost = uploadV2S3MaxConcurrency

	engine := newUploadV2Engine(&uploadIO{
		Client:         client,
		FileUploadPart: part,
		Size:           &size,
		uploadV2Tuning: UploadV2Tuning{S3WorkloadBytes: int64(200*200) * uploadV2MiB},
	}, plan)

	adjustedTransport, ok := engine.u.Config.Client.HTTPClient.Transport.(*lib.Transport)
	require.True(t, ok)
	assert.NotSame(t, originalTransport, adjustedTransport)
	assert.Equal(t, uploadV2S3MaxConcurrency, originalTransport.MaxConnsPerHost)
	assert.Equal(t, uploadV2S3GrowthCeiling, adjustedTransport.MaxConnsPerHost)
	assert.Equal(t, uploadV2DefaultHTTPIdleConnectionCap, adjustedTransport.MaxIdleConns)
	assert.Equal(t, uploadV2DefaultHTTPIdleConnectionCap, adjustedTransport.MaxIdleConnsPerHost)
	assert.Equal(t, uploadV2S3GrowthCeiling, engine.httpClientLimits.maxConnsPerHost)
	attrs := engine.uploadV2EnabledLogAttrs(uploadV2ReadyRunwayConfig{})
	assert.Equal(t, uploadV2S3GrowthCeiling, attrs["upload_max_conns_per_host"])
}

func TestUploadV2RaisesHTTPTransportCapForLargeS3Workload(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(256) * uploadV2MiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)
	client := &Client{Config: files_sdk.Config{}.Init()}

	engine := newUploadV2Engine(&uploadIO{
		Client:         client,
		FileUploadPart: part,
		Size:           &size,
		uploadV2Tuning: UploadV2Tuning{S3WorkloadBytes: uploadV2S3GrowthCeilingProbeBytes},
	}, plan)

	assert.Equal(t, uploadV2S3MaxConcurrency, engine.httpClientLimits.maxConnsPerHost)
}

func TestUploadV2HTTPTransportCapTracksS3GrowthCeiling(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")

	cases := []struct {
		name          string
		fileSize      int64
		workloadBytes int64
		tunePlan      bool
		want          int
	}{
		{name: "twenty by two hundred mib", fileSize: int64(200) * uploadV2MiB, workloadBytes: int64(20*200) * uploadV2MiB, want: uploadV2S3GrowthCeiling},
		{name: "two hundred by two hundred mib", fileSize: int64(200) * uploadV2MiB, workloadBytes: int64(200*200) * uploadV2MiB, tunePlan: true, want: uploadV2S3GrowthCeiling},
		{name: "single twenty gib file", fileSize: int64(20) * uploadV2GiB, workloadBytes: int64(20) * uploadV2GiB, want: uploadV2S3GrowthCeiling},
		{name: "large enough to probe above ceiling", fileSize: int64(200) * uploadV2MiB, workloadBytes: uploadV2S3GrowthCeilingProbeBytes, want: uploadV2S3MaxConcurrency},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			plan, ok, reason := newUploadV2PartPlanForUpload(part, &tt.fileSize)
			require.True(t, ok, reason)
			tuning := UploadV2Tuning{S3WorkloadBytes: tt.workloadBytes}
			testPlan := plan
			if tt.tunePlan {
				var ok bool
				testPlan, ok, reason = testPlan.withTuning(tuning)
				require.True(t, ok, reason)
			}
			got := uploadV2HTTPMaxConnsPerHost(testPlan, tuning, uploadV2S3MaxConcurrency)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUploadV2HTTPTransportCapOpensForManySmallS3Parts(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(8) * uploadV2MiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)
	tuning := UploadV2Tuning{
		S3WorkloadBytes:               int64(1000*8) * uploadV2MiB,
		S3GrowthCeilingProbeSuccesses: 384,
	}
	plan, ok, reason = plan.withTuning(tuning)
	require.True(t, ok, reason)

	got := uploadV2HTTPMaxConnsPerHost(plan, tuning, uploadV2S3MaxConcurrency)

	assert.Equal(t, int64(8)*uploadV2MiB, plan.partSize)
	assert.Equal(t, uploadV2S3MaxConcurrency, got)
}

func TestUploadV2LargeS3WorkloadOpensTransportBeforeAdaptiveGrowthUnlocks(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := uploadV2S3GrowthCeilingProbeBytes
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)
	tuning := UploadV2Tuning{S3WorkloadBytes: uploadV2S3GrowthCeilingProbeBytes}

	maxConnsPerHost := uploadV2HTTPMaxConnsPerHost(plan, tuning, uploadV2S3MaxConcurrency)
	manager := lib.NewAdaptiveConcurrencyManagerWithConfig(uploadV2AdaptiveConcurrencyConfigWithInitial(
		plan,
		uploadV2S3MaxConcurrency,
		uploadV2InitialConcurrencyForPlan(plan, uploadV2S3MaxConcurrency, tuning),
		tuning,
	))
	snapshot := manager.Snapshot()

	assert.Equal(t, uploadV2S3MaxConcurrency, maxConnsPerHost)
	assert.Equal(t, uploadV2S3InitialConcurrency, snapshot.Target)
	assert.Equal(t, uploadV2S3GrowthCeiling, snapshot.GrowthCeiling)
	assert.False(t, snapshot.GrowthCeilingUnlocked)
}

func TestUploadV2JobHTTPClientUsesS3GrowthCeilingForBenchmarkWorkload(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(200) * uploadV2MiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)
	client := &Client{Config: files_sdk.Config{}.Init()}
	job := (&Job{}).Init()
	tuning := UploadV2Tuning{S3WorkloadBytes: int64(200*200) * uploadV2MiB}
	maxConnsPerHost := uploadV2HTTPMaxConnsPerHost(plan, tuning, uploadV2S3MaxConcurrency)
	maxIdleConnsPerHost := uploadV2HTTPIdleConnectionCap(plan, maxConnsPerHost)

	_, limits, ok := job.uploadV2HTTPClient(client, plan, maxConnsPerHost, maxIdleConnsPerHost)

	require.True(t, ok)
	assert.Equal(t, uploadV2S3GrowthCeiling, maxConnsPerHost)
	assert.Equal(t, uploadV2S3GrowthCeiling, limits.maxConnsPerHost)
	assert.Equal(t, uploadV2DefaultHTTPIdleConnectionCap, limits.maxIdleConnsPerHost)
}

func TestUploadV2JobSharesAdjustedHTTPClientAcrossFiles(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(256) * uploadV2MiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)
	client := &Client{Config: files_sdk.Config{}.Init()}
	job := (&Job{}).Init()

	first := newUploadV2Engine(&uploadIO{
		Client:                     client,
		FileUploadPart:             part,
		Size:                       &size,
		uploadV2HTTPClientProvider: job.uploadV2HTTPClient,
	}, plan)
	second := newUploadV2Engine(&uploadIO{
		Client:                     client,
		FileUploadPart:             part,
		Size:                       &size,
		uploadV2HTTPClientProvider: job.uploadV2HTTPClient,
	}, plan)

	assert.Same(t, first.u.Client, second.u.Client)
	assert.NotSame(t, client, first.u.Client)
	assert.True(t, first.httpClientLimits.adjusted)
	assert.True(t, second.httpClientLimits.adjusted)
	assert.Equal(t, uploadV2DefaultHTTPIdleConnectionCap, first.httpClientLimits.maxIdleConnsPerHost)
	assert.Equal(t, 1, len(job.adaptiveUploadV2Clients))
}

func TestUploadV2S3RespectsExplicitManagerCap(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(256) * uploadV2MiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)

	engine := newUploadV2Engine(&uploadIO{
		FileUploadPart: part,
		Manager:        lib.NewConstrainedWorkGroup(50),
		managerSet:     true,
	}, plan)

	assert.Equal(t, 50, engine.manager.Max())
	assert.Equal(t, 16, engine.manager.Target())
}

func TestUploadV2S3IgnoresSchedulingOnlyManagerCap(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(256) * uploadV2MiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)

	engine := newUploadV2Engine(&uploadIO{
		FileUploadPart:            part,
		Manager:                   lib.NewConstrainedWorkGroup(50),
		managerSet:                true,
		uploadV2UseSDKDefaultCaps: true,
	}, plan)

	assert.Equal(t, uploadV2S3MaxConcurrency, engine.manager.Max())
	assert.Nil(t, engine.globalManager)
}

func TestUploadV2ResultsBufferUsesAdaptiveMax(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(10) * uploadV2GiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)

	engine := newUploadV2Engine(&uploadIO{FileUploadPart: part}, plan)

	assert.Equal(t, engine.manager.Max(), engine.resultsBufferSize())
	assert.GreaterOrEqual(t, engine.resultsBufferSize(), engine.manager.Target())
}

func TestUploadV2SmallS3KnownFileCapsPerUploadConcurrency(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(162) * uploadV2MiB
	tuning := UploadV2Tuning{S3WorkloadBytes: size}
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)
	plan, ok, reason = plan.withTuning(tuning)
	require.True(t, ok, reason)
	require.Equal(t, int64(8)*uploadV2MiB, plan.partSize)
	require.Equal(t, 21, plan.estimatedPartCount())

	engine := newUploadV2Engine(&uploadIO{
		FileUploadPart: part,
		Size:           &size,
		uploadV2Tuning: tuning,
	}, plan)

	assert.Equal(t, 11, engine.partConcurrencyLimit())
	assert.Equal(t, 11, engine.resultsBufferSize())
	attrs := engine.uploadV2EnabledLogAttrs(uploadV2ReadyRunwayConfig{})
	assert.Equal(t, 11, attrs["adaptive_part_target"])
	assert.Equal(t, 21, attrs["adaptive_planned_parts"])
}

func TestUploadV2SmallS3KnownFileAllowsSmallPlannedPartSet(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(162) * uploadV2MiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)
	require.Equal(t, 11, plan.estimatedPartCount())

	engine := newUploadV2Engine(&uploadIO{FileUploadPart: part, Size: &size}, plan)

	assert.Equal(t, 11, engine.partConcurrencyLimit())
	assert.Equal(t, 11, engine.resultsBufferSize())
}

func TestUploadV2LargeS3KnownFileKeepsAdaptiveMaxPerUploadConcurrency(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(10) * uploadV2GiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)

	engine := newUploadV2Engine(&uploadIO{FileUploadPart: part, Size: &size}, plan)

	assert.Equal(t, engine.manager.Max(), engine.partConcurrencyLimit())
	assert.Equal(t, engine.manager.Max(), engine.resultsBufferSize())
}

func TestUploadV2PartConcurrencyGateReportsAdaptiveSamples(t *testing.T) {
	manager := lib.NewAdaptiveConcurrencyManagerWithConfig(lib.AdaptiveConcurrencyConfig{
		MaxConcurrency: 10,
		InitialTarget:  10,
	})
	gate := newUploadV2PartConcurrencyGate(manager.NewSubWorker(), 2)

	require.True(t, gate.WaitWithContext(context.Background()))
	gate.DoneWithSample(lib.AdaptiveConcurrencySample{
		Success:  true,
		Bytes:    5,
		Duration: time.Millisecond,
	})
	gate.WaitAllDone()

	snapshot := manager.Snapshot()
	assert.Equal(t, 1, snapshot.SuccessTotal)
	assert.Equal(t, int64(5), snapshot.BytesTotal)
}

func TestUploadV2PartConcurrencyGateReleasesLocalSlotWhenParentWaitCancels(t *testing.T) {
	parent := &uploadV2GateTestParent{}
	gate := newUploadV2PartConcurrencyGate(parent, 1)

	require.False(t, gate.WaitWithContext(context.Background()))
	assert.Equal(t, 0, gate.RunningCount())
	gate.WaitAllDone()

	parent.allowWait.Store(true)
	require.True(t, gate.WaitWithContext(context.Background()))
	assert.Equal(t, 1, gate.RunningCount())
	assert.Equal(t, 1, parent.RunningCount())

	gate.Done()
	gate.WaitAllDone()
	assert.Equal(t, 0, gate.RunningCount())
	assert.Equal(t, 0, parent.RunningCount())
}

func TestUploadV2PartConcurrencyGateReportsFailureSamples(t *testing.T) {
	manager := lib.NewAdaptiveConcurrencyManagerWithConfig(lib.AdaptiveConcurrencyConfig{
		MaxConcurrency: 10,
		InitialTarget:  10,
	})
	gate := newUploadV2PartConcurrencyGate(manager.NewSubWorker(), 2)

	require.True(t, gate.WaitWithContext(context.Background()))
	gate.DoneWithSample(lib.AdaptiveConcurrencySample{
		Success:      false,
		Bytes:        7,
		Duration:     time.Millisecond,
		BackPressure: true,
		StatusCode:   http.StatusTooManyRequests,
	})
	gate.WaitAllDone()

	snapshot := manager.Snapshot()
	assert.Equal(t, 0, snapshot.SuccessTotal)
	assert.Equal(t, 1, snapshot.FailureTotal)
	assert.Equal(t, 1, snapshot.BackPressureTotal)
	assert.Equal(t, int64(7), snapshot.BytesTotal)
}

type uploadV2GateTestParent struct {
	allowWait atomic.Bool
	running   atomic.Int32
}

func (p *uploadV2GateTestParent) Wait() {
	p.running.Add(1)
}

func (p *uploadV2GateTestParent) Done() {
	p.running.Add(-1)
}

func (p *uploadV2GateTestParent) WaitAllDone() {}

func (p *uploadV2GateTestParent) RunningCount() int {
	return int(p.running.Load())
}

func (p *uploadV2GateTestParent) WaitWithContext(context.Context) bool {
	if !p.allowWait.Load() {
		return false
	}
	p.running.Add(1)
	return true
}

func (p *uploadV2GateTestParent) WaitForADone() bool {
	return false
}

func (p *uploadV2GateTestParent) WaitForADoneWithContext(context.Context) bool {
	return false
}

func TestUploadV2PreallocationIsBoundedByConcurrency(t *testing.T) {
	part := uploadV2TestPart("https://uploads.example.com/upload/file")
	size := int64(1) << 50
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)
	require.Greater(t, plan.estimatedPartCount(), 1_000_000)

	engine := newUploadV2Engine(&uploadIO{
		FileUploadPart: part,
		Size:           &size,
		readerAt:       bytes.NewReader([]byte("fixture")),
	}, plan)

	assert.Equal(t, engine.manager.Max()*uploadV2PreallocateConcurrencyMultiplier, cap(engine.u.Parts))
	assert.Less(t, cap(engine.u.Parts), plan.estimatedPartCount())
}

func TestUploadV2S3LargeUploadStartsHigherAndGrowsQuickly(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(10) * uploadV2GiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)

	manager := lib.NewAdaptiveConcurrencyManagerWithConfig(uploadV2AdaptiveConcurrencyConfig(plan, 50))

	assert.Equal(t, 50, manager.Target())
	for i := 0; i < 16; i++ {
		manager.Wait()
		manager.DoneWithSample(lib.AdaptiveConcurrencySample{Success: true})
	}
	assert.Equal(t, 50, manager.Target())
}

func TestUploadV2S3DefaultUsesMeasuredEnterpriseSoftCeiling(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(10) * uploadV2GiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)

	manager := lib.NewAdaptiveConcurrencyManagerWithConfig(uploadV2AdaptiveConcurrencyConfig(plan, uploadV2S3MaxConcurrency))

	assert.Equal(t, uploadV2S3InitialConcurrency, manager.Target())
	for i := 0; i < 16; i++ {
		manager.Wait()
		manager.DoneWithSample(lib.AdaptiveConcurrencySample{Success: true})
	}
	assert.Equal(t, uploadV2S3InitialConcurrency, manager.Target())
	assert.Equal(t, uploadV2S3MaxConcurrency, manager.Max())
	snapshot := manager.Snapshot()
	assert.Equal(t, uploadV2S3GrowthCeiling, snapshot.GrowthCeiling)
	assert.False(t, snapshot.GrowthCeilingUnlocked)
}

func TestUploadV2S3AdaptiveConfigUsesEnterprisePlateauEconomics(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(10) * uploadV2GiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)

	config := uploadV2AdaptiveConcurrencyConfig(plan, 1024)

	assert.Equal(t, uploadV2S3ThroughputProbeFloor, config.ThroughputProbeFloor)
	assert.Equal(t, uploadV2S3ThroughputProbePlateau, config.ThroughputProbePlateauTarget)
	assert.Equal(t, uploadV2S3ThroughputProbeMinGainPerTargetPercent, config.ThroughputProbeMinGainPerTargetPercent)
	assert.Equal(t, uploadV2S3ThroughputProbeLossTolerancePercent, config.ThroughputProbeLossTolerancePercent)
	assert.Equal(t, uploadV2S3GrowthCeiling, config.GrowthCeiling)
	assert.Equal(t, int64(uploadV2S3GrowthCeilingProbeBytes), config.GrowthCeilingProbeBytes)
	assert.Equal(t, uploadV2S3GrowthCeilingProbeSuccesses, config.GrowthCeilingProbeSuccesses)
	assert.Equal(t, float64(uploadV2S3GrowthCeilingProbeRate), config.GrowthCeilingProbeRate)
	assert.Equal(t, 8, config.ThroughputShrinkPercent)
	assert.Equal(t, 8, config.LatencyShrinkPercent)
	assert.Equal(t, float64(uploadV2S3LatencyQueueHigh), config.LatencyQueueHigh)
	assert.Equal(t, float64(uploadV2S3LatencyGrowthQueueHigh), config.LatencyGrowthQueueHigh)
}

func TestUploadV2S3AdaptiveConfigAppliesDiagnosticTuning(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(10) * uploadV2GiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)
	tuning := UploadV2Tuning{
		S3InitialTarget:                          180,
		S3AdaptiveFloor:                          120,
		S3GrowEvery:                              12,
		S3GrowStep:                               8,
		S3ThroughputWindow:                       64,
		S3ThroughputMinGainPercent:               2,
		S3ThroughputProbeMinWindows:              4,
		S3ThroughputProbeFloor:                   170,
		S3ThroughputProbeFloorRateBytesPerSecond: 123456,
		S3ThroughputProbePlateau:                 240,
		S3ThroughputShrinkPercent:                5,
		S3ThroughputHoldWindows:                  3,
		S3ThroughputProbeMinGainPerTargetPercent: 0.05,
		S3ThroughputProbeLossTolerancePercent:    3,
		S3GrowthCeiling:                          190,
		S3GrowthCeilingProbeBytes:                123456789,
		S3GrowthCeilingProbeSuccesses:            77,
		S3GrowthCeilingProbeRateBytesPerSecond:   987654321,
		S3LatencyQueueHigh:                       130,
		S3LatencyGrowthQueueHigh:                 140,
	}

	config := uploadV2SharedAdaptiveConcurrencyConfig(plan, 1024, tuning)

	assert.Equal(t, 180, config.InitialTarget)
	assert.Equal(t, 120, config.ThroughputFloor)
	assert.Equal(t, 120, config.LatencyFloor)
	assert.Equal(t, 12, config.GrowEvery)
	assert.Equal(t, 8, config.GrowStep)
	assert.Equal(t, 64, config.ThroughputWindow)
	assert.Equal(t, 2, config.ThroughputMinGainPercent)
	assert.Equal(t, 4, config.ThroughputProbeMinWindows)
	assert.Equal(t, 170, config.ThroughputProbeFloor)
	assert.Equal(t, float64(123456), config.ThroughputProbeFloorRate)
	assert.Equal(t, 240, config.ThroughputProbePlateauTarget)
	assert.Equal(t, 5, config.ThroughputShrinkPercent)
	assert.Equal(t, 3, config.ThroughputHoldWindows)
	assert.Equal(t, 0.05, config.ThroughputProbeMinGainPerTargetPercent)
	assert.Equal(t, 3, config.ThroughputProbeLossTolerancePercent)
	assert.Equal(t, 190, config.GrowthCeiling)
	assert.Equal(t, int64(123456789), config.GrowthCeilingProbeBytes)
	assert.Equal(t, 77, config.GrowthCeilingProbeSuccesses)
	assert.Equal(t, float64(987654321), config.GrowthCeilingProbeRate)
	assert.Equal(t, float64(130), config.LatencyQueueHigh)
	assert.Equal(t, float64(140), config.LatencyGrowthQueueHigh)
}

func TestUploadV2S3PartSizeDiagnosticTuningOverridesKnownSizePlan(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(100) * uploadV2GiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)
	assert.Equal(t, int64(64)*uploadV2MiB, plan.partSize)

	tuned, ok, reason := plan.withTuning(UploadV2Tuning{S3PartSizeMiB: 32})

	require.True(t, ok, reason)
	assert.Equal(t, int64(32)*uploadV2MiB, tuned.partSize)
	assert.Equal(t, "known_size_tuned", tuned.mode)
}

func TestUploadV2S3WorkloadTuningShrinksShortAggregateJob(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(200) * uploadV2MiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)
	assert.Equal(t, int64(16)*uploadV2MiB, plan.partSize)

	tuned, ok, reason := plan.withTuning(UploadV2Tuning{S3WorkloadBytes: int64(20*200) * uploadV2MiB})

	require.True(t, ok, reason)
	assert.Equal(t, int64(8)*uploadV2MiB, tuned.partSize)
	assert.Equal(t, "known_size_workload_tuned", tuned.mode)
}

func TestUploadV2S3WorkloadTuningUsesDiagnosticMultiplierAndMinPartSize(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(20) * uploadV2GiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)

	tuned, ok, reason := plan.withTuning(UploadV2Tuning{
		S3WorkloadBytes:                size,
		S3WorkloadTargetPartMultiplier: 16,
		S3WorkloadMinPartSizeMiB:       8,
	})

	require.True(t, ok, reason)
	assert.Equal(t, int64(8)*uploadV2MiB, tuned.partSize)
	assert.Equal(t, "known_size_workload_tuned", tuned.mode)
}

func TestUploadV2S3WorkloadTuningKeepsLargeAggregateJob(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(200) * uploadV2MiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)

	tuned, ok, reason := plan.withTuning(UploadV2Tuning{S3WorkloadBytes: int64(200*200) * uploadV2MiB})

	require.True(t, ok, reason)
	assert.Equal(t, plan.partSize, tuned.partSize)
	assert.Equal(t, plan.mode, tuned.mode)
}

func TestUploadV2S3WorkloadTuningShrinksSingleLargeFile(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(20) * uploadV2GiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)
	assert.Equal(t, int64(64)*uploadV2MiB, plan.partSize)

	tuned, ok, reason := plan.withTuning(UploadV2Tuning{S3WorkloadBytes: size})

	require.True(t, ok, reason)
	assert.Equal(t, int64(16)*uploadV2MiB, tuned.partSize)
	assert.Equal(t, "known_size_workload_tuned", tuned.mode)
}

func TestUploadV2S3ExplicitPartSizeTuningOverridesWorkload(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(20) * uploadV2GiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)

	tuned, ok, reason := plan.withTuning(UploadV2Tuning{
		S3PartSizeMiB:   32,
		S3WorkloadBytes: size,
	})

	require.True(t, ok, reason)
	assert.Equal(t, int64(32)*uploadV2MiB, tuned.partSize)
	assert.Equal(t, "known_size_tuned", tuned.mode)
}

func TestUploadV2JobWorkloadBytesUsesIndexedStatuses(t *testing.T) {
	job := (&Job{}).Init()
	job.Type = directory.Files
	job.Add(&UploadStatus{
		Mutex:  &sync.RWMutex{},
		status: status.Indexed,
		file:   files_sdk.File{Size: 200 * uploadV2MiB},
	})
	job.Add(&UploadStatus{
		Mutex:  &sync.RWMutex{},
		status: status.Uploading,
		file:   files_sdk.File{Size: 300 * uploadV2MiB},
	})
	job.Add(&UploadStatus{
		Mutex:  &sync.RWMutex{},
		status: status.Canceled,
		file:   files_sdk.File{Size: 400 * uploadV2MiB},
	})
	job.EndScan()

	assert.Equal(t, int64(500)*uploadV2MiB, job.uploadV2WorkloadBytes(0, UploadV2Tuning{}))
	assert.Equal(t, int64(700)*uploadV2MiB, job.uploadV2WorkloadBytes(700*uploadV2MiB, UploadV2Tuning{}))
}

func TestUploadV2JobWorkloadBytesWaitsForDirectoryScanCompletion(t *testing.T) {
	job := (&Job{}).Init()
	job.Type = directory.Dir
	job.Add(&UploadStatus{
		Mutex:  &sync.RWMutex{},
		status: status.Indexed,
		file:   files_sdk.File{Size: 200 * uploadV2MiB},
	})

	assert.Equal(t, int64(0), job.uploadV2WorkloadBytes(200*uploadV2MiB, UploadV2Tuning{S3WorkloadScanWaitMillis: 1}))

	job.EndScan()
	assert.Equal(t, int64(200)*uploadV2MiB, job.uploadV2WorkloadBytes(200*uploadV2MiB, UploadV2Tuning{}))
}

func TestUploadV2JobWorkloadBytesWaitsBrieflyForScanToEnd(t *testing.T) {
	job := (&Job{}).Init()
	job.Type = directory.Dir
	job.Add(&UploadStatus{
		Mutex:  &sync.RWMutex{},
		status: status.Indexed,
		file:   files_sdk.File{Size: 200 * uploadV2MiB},
	})
	go func() {
		time.Sleep(5 * time.Millisecond)
		job.EndScan()
	}()

	assert.Equal(t, int64(200)*uploadV2MiB, job.uploadV2WorkloadBytes(200*uploadV2MiB, UploadV2Tuning{S3WorkloadScanWaitMillis: 100}))
}

func TestUploadV2JobWorkloadBytesAllowsSingleFileBeforeScanCompletion(t *testing.T) {
	job := (&Job{}).Init()
	job.Type = directory.File
	job.Add(&UploadStatus{
		Mutex:  &sync.RWMutex{},
		status: status.Indexed,
		file:   files_sdk.File{Size: 20 * uploadV2GiB},
	})

	assert.Equal(t, int64(20)*uploadV2GiB, job.uploadV2WorkloadBytes(20*uploadV2GiB, UploadV2Tuning{}))
}

func TestUploadV2S3FastThroughputCanReachEnterpriseConcurrency(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(10) * uploadV2GiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)

	manager := lib.NewAdaptiveConcurrencyManagerWithConfig(uploadV2AdaptiveConcurrencyConfig(plan, 1024))

	for i := 0; i < 80; i++ {
		manager.Wait()
		manager.DoneWithSample(lib.AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    16 * uploadV2MiB,
			Duration: 30 * time.Millisecond,
		})
	}

	snapshot := manager.Snapshot()
	assert.GreaterOrEqual(t, snapshot.Target, uploadV2S3ThroughputProbeFloor)
	assert.GreaterOrEqual(t, snapshot.PeakTarget, uploadV2S3ThroughputProbeFloor)
	assert.Equal(t, 0, snapshot.ThroughputBackoffTotal)
	assert.Greater(t, snapshot.BestThroughputBytesPerSecond, float64(uploadV2S3ThroughputProbeFloorRate))
}

func TestUploadV2S3LargeUnlockedProbeClimbsToPlateauOnNeutralThroughput(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(80) * uploadV2GiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)
	tuning := UploadV2Tuning{S3GrowthCeilingProbeBytes: 1}

	manager := lib.NewAdaptiveConcurrencyManagerWithConfig(uploadV2SharedAdaptiveConcurrencyConfig(plan, uploadV2S3MaxConcurrency, tuning))

	for i := 0; i < 360; i++ {
		manager.Wait()
		manager.DoneWithSample(lib.AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    plan.partSize,
			Duration: 40 * time.Millisecond,
		})
	}

	snapshot := manager.Snapshot()
	assert.True(t, snapshot.GrowthCeilingUnlocked)
	assert.GreaterOrEqual(t, snapshot.PeakTarget, uploadV2S3ThroughputProbePlateau)
	assert.Greater(t, snapshot.Target, uploadV2S3GrowthCeiling)
}

func TestUploadV2S3SeekableRunwayScalesWithInitialTarget(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(20) * uploadV2GiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)

	engine := newUploadV2Engine(&uploadIO{
		FileUploadPart: part,
		Size:           &size,
		readerAt:       bytes.NewReader(make([]byte, 1)),
	}, plan)

	runway := engine.readyRunwayConfig()
	assert.Equal(t, uploadV2S3InitialConcurrency/uploadV2DefaultSeekableS3ReadyRunwayTargetDivisor, runway.parts)
	assert.Equal(t, uploadV2DefaultReadyRunwayBytes, runway.bytes)
}

func TestUploadV2ExplicitReadyRunwayIsRespectedForS3(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(20) * uploadV2GiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)

	engine := newUploadV2Engine(&uploadIO{
		FileUploadPart: part,
		Size:           &size,
		readerAt:       bytes.NewReader(make([]byte, 1)),
		uploadV2ReadyRunway: uploadV2ReadyRunwayConfig{
			parts:      6,
			bytes:      7,
			configured: true,
		},
	}, plan)

	runway := engine.readyRunwayConfig()
	assert.Equal(t, 6, runway.parts)
	assert.Equal(t, int64(7), runway.bytes)
}

func TestUploadV2JobSharesAdaptiveManagerAcrossFiles(t *testing.T) {
	resetUploadV2SharedAdaptiveManagersForTest()
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(10) * uploadV2GiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)
	job := (&Job{}).Init()

	first := job.uploadV2AdaptiveManager(plan, 50, UploadV2Tuning{})
	first.Wait()
	first.DoneWithSample(lib.AdaptiveConcurrencySample{Success: true})
	first.Wait()
	first.DoneWithSample(lib.AdaptiveConcurrencySample{Success: true})
	second := job.uploadV2AdaptiveManager(plan, 50, UploadV2Tuning{})

	assert.Same(t, first, second)
	assert.Equal(t, 50, second.Target())
}

func TestUploadV2JobsShareAdaptiveManagerGlobally(t *testing.T) {
	resetUploadV2SharedAdaptiveManagersForTest()
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(10) * uploadV2GiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)
	firstJob := (&Job{}).Init()
	secondJob := (&Job{}).Init()

	first := firstJob.uploadV2AdaptiveManager(plan, 50, UploadV2Tuning{})
	second := secondJob.uploadV2AdaptiveManager(plan, 50, UploadV2Tuning{})

	assert.Same(t, first, second)
}

func TestAdaptiveTransferStatsReportsSharedUploadManagers(t *testing.T) {
	resetUploadV2SharedAdaptiveManagersForTest()
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(10) * uploadV2GiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)
	firstJob := (&Job{}).Init()
	secondJob := (&Job{}).Init()

	first := firstJob.uploadV2AdaptiveManager(plan, 50, UploadV2Tuning{})
	second := secondJob.uploadV2AdaptiveManager(plan, 50, UploadV2Tuning{})
	idle := secondJob.uploadV2AdaptiveManager(plan, 25, UploadV2Tuning{})
	require.Same(t, first, second)
	require.NotSame(t, first, idle)
	first.Wait()
	t.Cleanup(first.Done)

	stats := AdaptiveTransferStats()

	assert.Equal(t, 1, stats.Upload.Active)
	assert.Equal(t, 50, stats.Upload.Max)
	assert.Equal(t, 0, stats.Download.Active)
	assert.Equal(t, 0, stats.Download.Max)
}

func TestUploadV2SharedAdaptiveManagerKeepsExplicitCapsSeparate(t *testing.T) {
	resetUploadV2SharedAdaptiveManagersForTest()
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(10) * uploadV2GiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)
	firstJob := (&Job{}).Init()
	secondJob := (&Job{}).Init()

	limited := firstJob.uploadV2AdaptiveManager(plan, 50, UploadV2Tuning{})
	defaulted := secondJob.uploadV2AdaptiveManager(plan, uploadV2S3MaxConcurrency, UploadV2Tuning{})

	assert.NotSame(t, limited, defaulted)
	assert.Equal(t, 50, limited.Max())
	assert.Equal(t, uploadV2S3MaxConcurrency, defaulted.Max())
}

func TestUploadV2JobSharedManagerIgnoresPlannerOnlyTuning(t *testing.T) {
	resetUploadV2SharedAdaptiveManagersForTest()
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(20) * uploadV2GiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)
	job := (&Job{}).Init()

	first := job.uploadV2AdaptiveManager(plan, 50, UploadV2Tuning{
		S3PartSizeMiB:   16,
		S3WorkloadBytes: size,
	})
	second := job.uploadV2AdaptiveManager(plan, 50, UploadV2Tuning{
		S3PartSizeMiB:   32,
		S3WorkloadBytes: int64(200*200) * uploadV2MiB,
	})

	assert.Same(t, first, second)
}

func TestUploadV2JobSharedManagerDoesNotStartAtTinyFilePartCount(t *testing.T) {
	resetUploadV2SharedAdaptiveManagersForTest()
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	tinySize := int64(1)
	tinyPlan, ok, reason := newUploadV2PartPlanForUpload(part, &tinySize)
	require.True(t, ok, reason)
	largeSize := int64(10) * uploadV2GiB
	largePlan, ok, reason := newUploadV2PartPlanForUpload(part, &largeSize)
	require.True(t, ok, reason)
	job := (&Job{}).Init()

	first := job.uploadV2AdaptiveManager(tinyPlan, 50, UploadV2Tuning{})
	second := job.uploadV2AdaptiveManager(largePlan, 50, UploadV2Tuning{})

	assert.Same(t, first, second)
	assert.Equal(t, 50, second.Target())
}

func TestUploadV2ClearStatusesKeepsSharedAdaptiveManagers(t *testing.T) {
	resetUploadV2SharedAdaptiveManagersForTest()
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(10) * uploadV2GiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)
	job := (&Job{}).Init()
	first := job.uploadV2AdaptiveManager(plan, 50, UploadV2Tuning{})
	first.Wait()
	first.DoneWithSample(lib.AdaptiveConcurrencySample{Success: true})
	first.Wait()
	first.DoneWithSample(lib.AdaptiveConcurrencySample{Success: true})

	retryJob := job.ClearStatuses()
	second := retryJob.uploadV2AdaptiveManager(plan, 50, UploadV2Tuning{})

	assert.Same(t, first, second)
	assert.Equal(t, 50, second.Target())
}

func TestUploadV2NetworkErrorsCountAsBackPressure(t *testing.T) {
	result := uploadV2PartResult{err: errors.New("write tcp: write: broken pipe")}

	sample := result.sample()

	assert.False(t, sample.Success)
	assert.True(t, sample.BackPressure)
}

func TestUploadV2BufferedReaderPartsAreReplayable(t *testing.T) {
	part := uploadV2TestPart("https://uploads.example.com/upload/file")
	engine := newUploadV2TestEngine(t, part)
	size := int64(5)
	engine.u.Size = &size
	engine.u.reader = bytes.NewBufferString("hello")

	reader, progress, err := engine.buildReader(OffSet{off: 0, len: size})
	require.NoError(t, err)
	defer reader.Close()

	first, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, "hello", string(first))
	progress.Flush()

	require.True(t, reader.Rewind())
	second, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, "hello", string(second))
}

func TestUploadV2ReadyRunwayMemoryCostOnlyCountsBufferedParts(t *testing.T) {
	size := int64(5)
	readerAtEngine := &uploadV2Engine{u: &uploadIO{
		Size:     &size,
		readerAt: bytes.NewReader([]byte("hello")),
	}}
	assert.Equal(t, int64(0), readerAtEngine.readyRunwayMemoryCost(size))

	streamEngine := &uploadV2Engine{u: &uploadIO{
		Size:   &size,
		reader: bytes.NewBufferString("hello"),
	}}
	assert.Equal(t, size, streamEngine.readyRunwayMemoryCost(size))
}

func TestUploadV2ReadyRunwayByteCapBoundsBufferedQueue(t *testing.T) {
	size := int64(10)
	engine := &uploadV2Engine{u: &uploadIO{
		Size:   &size,
		reader: bytes.NewBufferString("0123456789"),
	}}
	runway := uploadV2ReadyRunwayConfig{parts: 4, bytes: 2}

	assert.False(t, engine.readyRunwayCanPrepare(nil, 0, 3, runway))
	assert.True(t, engine.readyRunwayCanPrepare(nil, 0, 2, runway))
	assert.True(t, engine.readyRunwayCanPrepare([]uploadV2PreparedPart{{memoryBytes: 1}}, 1, 1, runway))
	assert.False(t, engine.readyRunwayCanPrepare([]uploadV2PreparedPart{{memoryBytes: 1}}, 1, 2, runway))
	assert.False(t, engine.readyRunwayCanPrepare([]uploadV2PreparedPart{{memoryBytes: 1}, {memoryBytes: 1}}, 2, 0, uploadV2ReadyRunwayConfig{parts: 2, bytes: 0}))
}

func TestUploadV2SchedulerStatsExposeTuningSignals(t *testing.T) {
	var stats uploadV2SchedulerStats
	stats.recordDirectScheduled()
	stats.recordPrepared(32*uploadV2MiB, 0)
	stats.recordDispatchedPrepared()
	stats.recordRunwayPartCapBlock()
	stats.recordReadyDepth(4, 0)
	stats.recordReaderBuild(32*uploadV2MiB, time.Millisecond, nil)
	stats.recordAdaptiveWait(2*time.Millisecond, true)
	stats.recordUploadURLRefresh(3*time.Millisecond, nil)
	stats.recordHTTPCall(4*time.Millisecond, 32*uploadV2MiB, nil)
	stats.recordPartComplete(32*uploadV2MiB, 11*time.Second, nil)

	attrs := map[string]any{}
	stats.addLogAttrs(attrs)

	assert.Equal(t, 1, attrs["scheduler_direct_scheduled_parts"])
	assert.Equal(t, 1, attrs["scheduler_prepared_parts"])
	assert.Equal(t, 1, attrs["scheduler_dispatched_prepared_parts"])
	assert.Equal(t, 1, attrs["scheduler_runway_part_cap_blocks"])
	assert.Equal(t, 4, attrs["scheduler_max_ready_parts"])
	assert.Equal(t, int64(time.Millisecond), attrs["scheduler_reader_build_duration_ns"])
	assert.Equal(t, int64(2*time.Millisecond), attrs["scheduler_adaptive_wait_duration_ns"])
	assert.Equal(t, int64(3*time.Millisecond), attrs["scheduler_upload_url_refresh_duration_ns"])
	assert.Equal(t, int64(4*time.Millisecond), attrs["scheduler_http_call_duration_ns"])
	assert.Equal(t, int64(32*uploadV2MiB), attrs["scheduler_http_call_bytes"])
	assert.Equal(t, 1, attrs["scheduler_part_complete_count"])
	assert.Equal(t, int64(11*time.Second), attrs["scheduler_part_complete_duration_ns"])
	assert.Equal(t, int64(11*time.Second), attrs["scheduler_part_complete_max_duration_ns"])
	assert.Equal(t, int64(uploadV2SlowPartThreshold), attrs["scheduler_part_slow_threshold_ns"])
	assert.Equal(t, 1, attrs["scheduler_part_slow_count"])
	assert.Equal(t, int64(11*time.Second), attrs["scheduler_part_slow_duration_ns"])
	assert.Equal(t, int64(11*time.Second), attrs["scheduler_part_slow_max_duration_ns"])
}

func TestUploadV2EtagsForCompleteSortsByPartNumber(t *testing.T) {
	engine := newUploadV2TestEngine(t, uploadV2TestPart("https://uploads.example.com/upload/file"))
	engine.recordSuccess(3, files_sdk.EtagsParam{Part: "3", Etag: "three"}, 1)
	engine.recordSuccess(1, files_sdk.EtagsParam{Part: "1", Etag: "one"}, 1)
	engine.recordSuccess(2, files_sdk.EtagsParam{Part: "2", Etag: "two"}, 1)

	got := engine.etagsForComplete()

	assert.Equal(t, []files_sdk.EtagsParam{
		{Part: "1", Etag: "one"},
		{Part: "2", Etag: "two"},
		{Part: "3", Etag: "three"},
	}, got)
}

func TestUploadV2ResumeValidation(t *testing.T) {
	engine := newUploadV2TestEngine(t, uploadV2TestPart("https://uploads.example.com/upload/file"))

	t.Run("accepts contiguous offset parts", func(t *testing.T) {
		engine.u.Parts = Parts{
			NewPart(1, 0, 16*uploadV2MiB, 16*uploadV2MiB, "one", "1", ""),
			NewPart(2, 16*uploadV2MiB, 16*uploadV2MiB, 0, "", "", "interrupted"),
		}
		engine.u.Parts[0].FileUploadPart.UploadUri = "https://uploads.example.com/upload/file?part_number=1&part_offset=0"
		assert.Empty(t, engine.resumeResetReason())
	})

	t.Run("rejects offset gaps", func(t *testing.T) {
		engine.u.Parts = Parts{
			NewPart(1, 0, 16*uploadV2MiB, 16*uploadV2MiB, "one", "1", ""),
			NewPart(2, 32*uploadV2MiB, 16*uploadV2MiB, 0, "", "", "interrupted"),
		}
		engine.u.Parts[0].FileUploadPart.UploadUri = "https://uploads.example.com/upload/file?part_number=1&part_offset=0"
		assert.Equal(t, "non_contiguous_part_offset", engine.resumeResetReason())
	})

	t.Run("rejects successful offset upload parts without part offset", func(t *testing.T) {
		engine.u.Parts = Parts{
			NewPart(1, 0, 16*uploadV2MiB, 16*uploadV2MiB, "one", "1", ""),
		}
		engine.u.Parts[0].FileUploadPart.UploadUri = "https://uploads.example.com/upload/file?part_number=1"
		assert.Equal(t, "successful_offset_part_missing_part_offset", engine.resumeResetReason())
	})

	t.Run("rejects planner mismatch", func(t *testing.T) {
		engine.u.Parts = Parts{
			NewPart(1, 0, uploadV2MiB, uploadV2MiB, "one", "1", ""),
		}
		engine.u.Parts[0].FileUploadPart.UploadUri = "https://uploads.example.com/upload/file?part_number=1&part_offset=0"
		assert.Equal(t, "part_size_plan_mismatch", engine.resumeResetReason())
	})
}

func TestUploadV2RefreshesPartURLWithStableOffset(t *testing.T) {
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	server.MockFiles["v2-refresh.txt"] = mockFile{File: files_sdk.File{Size: 3}}

	client := server.Client()
	part := files_sdk.FileUploadPart{
		HttpMethod:    "POST",
		Path:          "v2-refresh.txt",
		Ref:           "put-refresh",
		ParallelParts: lib.Bool(true),
		PartNumber:    1,
		UploadUri:     server.Server.URL + "/upload/v2-refresh.txt?part_number=1",
		Expires:       time.Now().Add(time.Hour).Format(time.RFC3339),
	}
	engine := newUploadV2TestEngine(t, part)
	engine.u.Client = client
	engine.u.Path = "v2-refresh.txt"
	engine.u.readerAt = bytes.NewReader([]byte("abc"))
	engine.u.Progress = func(int64) {}

	uploadPart := &uploadV2Part{
		uploadV2PartDescriptor: uploadV2PartDescriptor{
			number: 2,
			offset: OffSet{off: 1, len: 1},
			upload: files_sdk.FileUploadPart{
				HttpMethod:    "POST",
				Path:          "v2-refresh.txt",
				Ref:           "put-refresh",
				ParallelParts: lib.Bool(true),
				PartNumber:    2,
			},
			legacy: &Part{},
		},
		reader: &ProxyReaderAt{
			ReaderAt: bytes.NewReader([]byte("abc")),
			off:      1,
			len:      1,
			onRead:   func(int64) {},
		},
	}

	_, _, _, _, err := engine.uploadPart(context.Background(), uploadPart)

	require.NoError(t, err)
	uploadRequests := server.TrackRequest["/upload/*path"]
	require.Len(t, uploadRequests, 1)
	assert.Contains(t, uploadRequests[0], "part_number=2")
	assert.Contains(t, uploadRequests[0], "part_offset=1")
}

func TestUploadV2PartRetryUsesSameOffsetAndWinningETag(t *testing.T) {
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	server.MockFiles["v2-retry.txt"] = mockFile{File: files_sdk.File{Size: 1}}
	attempts := 0
	server.MockRoute("/upload/v2-retry.txt", func(ctx *gin.Context, model interface{}) bool {
		attempts++
		if attempts == 1 {
			ctx.Header("Retry-After", "1")
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "busy"})
			return true
		}
		ctx.Header("Etag", "winning-etag")
		ctx.Status(http.StatusOK)
		return true
	})

	client := server.Client()
	part := files_sdk.FileUploadPart{
		HttpMethod:    "POST",
		Path:          "v2-retry.txt",
		Ref:           "put-retry",
		ParallelParts: lib.Bool(true),
		PartNumber:    1,
		UploadUri:     server.Server.URL + "/upload/v2-retry.txt?part_number=1",
		Expires:       time.Now().Add(time.Hour).Format(time.RFC3339),
	}
	engine := newUploadV2TestEngine(t, part)
	engine.u.Client = client
	engine.u.Path = "v2-retry.txt"
	engine.u.Progress = func(int64) {}
	descriptor := uploadV2PartDescriptor{
		number: 1,
		offset: OffSet{off: 0, len: 1},
		upload: part,
		legacy: &Part{},
	}
	engine.decorateUploadURL(&descriptor.upload, 1, 0)
	uploadPart := &uploadV2Part{
		uploadV2PartDescriptor: descriptor,
		reader: &ProxyReaderAt{
			ReaderAt: bytes.NewReader([]byte("x")),
			off:      0,
			len:      1,
			onRead:   func(int64) {},
		},
	}

	result := engine.runPart(context.Background(), uploadPart)

	require.NoError(t, result.err)
	assert.Equal(t, "winning-etag", result.etag.Etag)
	assert.True(t, result.backPressure)
	assert.Equal(t, time.Second, result.retryAfter)
	assert.Equal(t, 2, attempts)
	uploadRequests := server.TrackRequest["/upload/*path"]
	require.Len(t, uploadRequests, 2)
	assert.Contains(t, uploadRequests[0], "part_number=1")
	assert.Contains(t, uploadRequests[0], "part_offset=0")
	assert.Contains(t, uploadRequests[1], "part_number=1")
	assert.Contains(t, uploadRequests[1], "part_offset=0")
}

func TestUploadV2UploadPartClosesErrorResponseBody(t *testing.T) {
	closed := &atomic.Bool{}
	transport := uploadV2RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Header: http.Header{
				"Content-Type": []string{"text/html"},
				"Server":       []string{"nginx"},
				"X-Request-Id": []string{"request-id"},
			},
			Body:    uploadV2CloseTrackingBody{Reader: strings.NewReader("bad request"), closed: closed},
			Request: req,
		}, nil
	})
	part := uploadV2TestPart("https://uploads.example.com/upload/file?part_number=1")
	engine := newUploadV2TestEngine(t, part)
	engine.u.Client = &Client{Config: files_sdk.Config{}.Init().SetCustomClient(&http.Client{Transport: transport})}
	engine.u.Progress = func(int64) {}
	descriptor := uploadV2PartDescriptor{
		number: 1,
		offset: OffSet{off: 0, len: 1},
		upload: files_sdk.FileUploadPart{
			HttpMethod:    "POST",
			UploadUri:     "https://uploads.example.com/upload/file?part_number=1",
			ParallelParts: lib.Bool(true),
			PartNumber:    1,
			Expires:       time.Now().Add(time.Hour).Format(time.RFC3339),
		},
		legacy: &Part{},
	}
	engine.decorateUploadURL(&descriptor.upload, 1, 0)
	uploadPart := &uploadV2Part{
		uploadV2PartDescriptor: descriptor,
		reader: &ProxyReaderAt{
			ReaderAt: bytes.NewReader([]byte("x")),
			off:      0,
			len:      1,
			onRead:   func(int64) {},
		},
	}

	_, _, _, _, err := engine.uploadPart(context.Background(), uploadPart)

	require.Error(t, err)
	assert.True(t, closed.Load())
}

func collectUploadV2Offsets(t *testing.T, part files_sdk.FileUploadPart, size *int64) []OffSet {
	t.Helper()
	iterator, ok := newUploadV2PartIterator(part, size, 0, 0)
	require.True(t, ok)

	var offsets []OffSet
	for iterator != nil {
		var offset OffSet
		offset, iterator, _ = iterator()
		offsets = append(offsets, offset)
	}
	return offsets
}

func sumOffsets(offsets []OffSet) int64 {
	var total int64
	for _, offset := range offsets {
		total += offset.len
	}
	return total
}

func uploadV2TestPart(uploadURI string) files_sdk.FileUploadPart {
	return files_sdk.FileUploadPart{
		UploadUri:     uploadURI,
		ParallelParts: lib.Bool(true),
	}
}

func resetUploadV2SharedAdaptiveManagersForTest() {
	uploadV2SharedAdaptiveManagers.resetForTest()
}

func uploadV2TestQuery(t *testing.T, uploadURI string) url.Values {
	t.Helper()
	parsed, err := url.Parse(uploadURI)
	require.NoError(t, err)
	return parsed.Query()
}

func uploadV2ChecksumTrailerSignedS3URL() string {
	return "https://s3.amazonaws.com/bucket/key?partNumber=1&X-Amz-SignedHeaders=content-encoding%3Bhost%3Bx-amz-content-sha256%3Bx-amz-decoded-content-length%3Bx-amz-sdk-checksum-algorithm%3Bx-amz-trailer"
}

func mustParseInt64(t *testing.T, value string) int64 {
	t.Helper()
	parsed, err := strconv.ParseInt(value, 10, 64)
	require.NoError(t, err)
	return parsed
}

func newUploadV2TestEngine(t *testing.T, part files_sdk.FileUploadPart) *uploadV2Engine {
	t.Helper()
	return newUploadV2TestEngineWithClient(t, part, nil)
}

func newUploadV2TestEngineWithClient(t *testing.T, part files_sdk.FileUploadPart, client *Client) *uploadV2Engine {
	t.Helper()
	size := int64(100)
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)
	return newUploadV2Engine(&uploadIO{FileUploadPart: part, Manager: lib.NewConstrainedWorkGroup(1), Client: client}, plan)
}

type uploadV2RoundTripFunc func(*http.Request) (*http.Response, error)

func (f uploadV2RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type uploadV2CloseTrackingBody struct {
	*strings.Reader
	closed *atomic.Bool
}

func (b uploadV2CloseTrackingBody) Close() error {
	b.closed.Store(true)
	return nil
}
