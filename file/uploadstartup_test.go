package file

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// startupS3 is a Files.com API and an Amazon S3 bucket for uploads: it negotiates each part's URL, records each part
// request's size, and holds part requests until release while holding.
type startupS3 struct {
	mu       sync.Mutex
	holding  bool
	release  chan struct{}
	sending  int
	peak     int
	partSize map[string][]int64
}

func newStartupS3(holding bool) *startupS3 {
	return &startupS3{holding: holding, release: make(chan struct{}), partSize: map[string][]int64{}}
}

func (s *startupS3) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const beginUpload, files, bucket = "/api/rest/v1/file_actions/begin_upload/", "/api/rest/v1/files/", "/bucket/"
	path := r.URL.Path
	switch {
	case r.Method == http.MethodPost && strings.HasPrefix(path, beginUpload):
		remote := strings.TrimPrefix(path, beginUpload)
		var params files_sdk.FileBeginUploadParams
		_ = json.NewDecoder(r.Body).Decode(&params)
		part := max(params.Part, 1)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(files_sdk.FileUploadPartCollection{{
			HttpMethod:    http.MethodPut,
			Path:          remote,
			Ref:           "ref-" + remote,
			UploadUri:     fmt.Sprintf("https://s3.amazonaws.com%s%s?partNumber=%d", bucket, remote, part),
			ParallelParts: lib.Bool(true),
			Expires:       time.Now().Add(time.Hour).Format(time.RFC3339),
			PartNumber:    part,
		}})
	case r.Method == http.MethodPut && strings.HasPrefix(path, bucket):
		size, _ := io.Copy(io.Discard, r.Body)
		remote := strings.TrimPrefix(path, bucket)
		s.mu.Lock()
		s.partSize[remote] = append(s.partSize[remote], size)
		s.sending++
		s.peak = max(s.peak, s.sending)
		holding := s.holding
		s.mu.Unlock()
		if holding {
			select {
			case <-s.release:
			case <-r.Context().Done():
			}
		}
		s.mu.Lock()
		s.sending--
		s.mu.Unlock()
		w.Header().Set("Etag", `"etag"`)
		w.WriteHeader(http.StatusOK)
	case r.Method == http.MethodPost && strings.HasPrefix(path, files):
		remote := strings.TrimPrefix(path, files)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(files_sdk.File{Path: remote, DisplayName: filepath.Base(remote), Type: "file"})
	default:
		http.NotFound(w, r)
	}
}

// Sending is how many part requests the bucket holds or is receiving now.
func (s *startupS3) Sending() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sending
}

func (s *startupS3) releaseAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.holding {
		s.holding = false
		close(s.release)
	}
}

// PartSizes are the sizes of the part requests the bucket received, sorted.
func (s *startupS3) PartSizes() []int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	var sizes []int64
	for _, file := range s.partSize {
		sizes = append(sizes, file...)
	}
	sort.Slice(sizes, func(i, j int) bool { return sizes[i] < sizes[j] })
	return sizes
}

func startupS3Client(t *testing.T, storage *startupS3) *Client {
	t.Helper()
	server := httptest.NewServer(storage)
	t.Cleanup(server.Close)
	t.Cleanup(storage.releaseAll)
	serverURL, err := url.Parse(server.URL)
	require.NoError(t, err)
	client := &Client{Config: files_sdk.Config{}.Init()}
	client.Config.Logger = lib.NullLogger{}
	client.Config = client.Config.SetCustomClient(&http.Client{Transport: &CustomTransport{URL: serverURL}})
	client.Config.Client.RetryMax = 0
	return client
}

func onePartFiles(t *testing.T, count int, size int) string {
	t.Helper()
	local := t.TempDir()
	for i := range count {
		require.NoError(t, os.WriteFile(filepath.Join(local, fmt.Sprintf("%02d.bin", i)), bytes.Repeat([]byte{byte(i)}, size), 0600))
	}
	return local + string(os.PathSeparator)
}

// A directory upload to S3 with SDK default caps, of files of one part each. By default, while storage has not
// answered, it sends one part at a time, where the aggressive startup sends a part for every file it admits (after
// its admission's first recheck, a second in); once answers come fast enough to show a fast link, the files' shared
// target goes straight to the policy's start. Every file's parts taught that one target, and the job completes.
func TestUploadJobStartupSendsOnePartUntilAnswersShowAFastLink(t *testing.T) {
	for _, tt := range []struct {
		startup UploadStartup
		// wantHeld is how many part requests storage holds at once before answering any.
		wantHeld func(t *testing.T, held int)
	}{
		{startup: "", wantHeld: func(t *testing.T, held int) { assert.Equal(t, 1, held) }},
		{startup: UploadStartupAggressive, wantHeld: func(t *testing.T, held int) { assert.GreaterOrEqual(t, held, 20) }},
	} {
		t.Run(string(tt.startup), func(t *testing.T) {
			storage := newStartupS3(true)
			client := startupS3Client(t, storage)
			job := client.Uploader(UploaderParams{
				LocalPath:             onePartFiles(t, 40, 1<<20),
				RemotePath:            "startup",
				AdaptiveConcurrency:   true,
				AdaptiveUploadStartup: tt.startup,
			})
			job.Start()

			require.Eventually(t, func() bool { return storage.Sending() >= 1 }, 10*time.Second, 5*time.Millisecond)
			time.Sleep(1500 * time.Millisecond)
			tt.wantHeld(t, storage.Sending())
			storage.releaseAll()
			job.Wait()

			assert.Equal(t, 40, job.Count(status.Complete))
			require.Len(t, job.adaptiveUploadV2Managers, 1, "the files share one target")
			for _, manager := range job.adaptiveUploadV2Managers {
				snapshot := manager.Snapshot()
				assert.Equal(t, 40, snapshot.SuccessTotal, "every file's part taught it")
				if tt.startup == UploadStartupAggressive {
					assert.False(t, snapshot.WarmStarted)
					continue
				}
				assert.Equal(t, "fast_link", snapshot.WarmStartEnd)
				assert.Equal(t, AdaptiveTransferS3InitialTarget, snapshot.PeakTarget)
			}
		})
	}
}

// Starting with one part does not plan large files in fewer, larger parts: a 256 MiB upload in a job of that size is
// planned for the policy's start, as eight parts for each of its 150 requests at most, so its parts are 8 MiB under
// either startup (planned for a start of one, they would be 16 MiB).
func TestUploadJobStartupPlansTheSamePartsForEitherStartup(t *testing.T) {
	local := t.TempDir()
	sparse, err := os.Create(filepath.Join(local, "large.bin"))
	require.NoError(t, err)
	require.NoError(t, sparse.Truncate(256<<20))
	require.NoError(t, sparse.Close())

	for _, startup := range []UploadStartup{UploadStartupDefault, UploadStartupAggressive} {
		storage := newStartupS3(false)
		client := startupS3Client(t, storage)
		job := client.Uploader(UploaderParams{
			LocalPath:             local + string(os.PathSeparator),
			RemotePath:            "large",
			AdaptiveConcurrency:   true,
			AdaptiveUploadStartup: startup,
		})
		job.Start()
		job.Wait()

		require.Equal(t, 1, job.Count(status.Complete), startup)
		sizes := storage.PartSizes()
		assert.Len(t, sizes, 32, startup)
		for _, size := range sizes {
			assert.Equal(t, int64(8<<20), size, startup)
		}
	}
}

// An explicit start in the tuning is where the upload starts, whichever startup is selected; a startup the SDK does
// not know fails the job before anything is sent.
func TestUploadStartupPrecedenceAndValidation(t *testing.T) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := int64(10) * uploadV2GiB
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	require.True(t, ok, reason)
	for _, tuning := range []UploadV2Tuning{{S3InitialTarget: 20}, {InitialTarget: 20}} {
		engine := newUploadV2Engine(&uploadIO{FileUploadPart: part, uploadV2Tuning: tuning}, plan)
		assert.Equal(t, lib.AdaptiveConcurrencyStartupState{Target: 20, StartTarget: 20}, engine.manager.StartupState())
	}

	storage := newStartupS3(false)
	client := startupS3Client(t, storage)
	job := client.Uploader(UploaderParams{
		LocalPath:             onePartFiles(t, 2, 16),
		RemotePath:            "invalid",
		AdaptiveConcurrency:   true,
		AdaptiveUploadStartup: "eager",
	})
	job.Start()
	job.Wait()
	assert.Equal(t, 1, job.Count(status.Errored))
	assert.Empty(t, storage.PartSizes())
}
