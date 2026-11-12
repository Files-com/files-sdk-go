package file

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const finalizationTestRef = "finalization-test-ref"

var finalizationTestData = []byte("0123456789")

// finalizationServer is a MockAPIServer whose uploads to path get finalizationTestRef and whose requests finishing the
// upload are answered by answers in turn, the last one again for every later request.
type finalizationServer struct {
	*MockAPIServer
	path     string
	mutex    sync.Mutex
	requests []finalizationRequest
	answers  []func(*gin.Context)
}

type finalizationRequest struct {
	params files_sdk.FileCreateParams
	header http.Header
	host   string
	at     time.Time
}

func newFinalizationServer(t *testing.T, path string, answers ...func(*gin.Context)) *finalizationServer {
	server := (&MockAPIServer{T: t}).Do()
	t.Cleanup(server.Shutdown)
	f := &finalizationServer{MockAPIServer: server, path: path, answers: answers}
	server.MockFiles[path] = mockFile{File: files_sdk.File{Size: int64(len(finalizationTestData))}}
	server.MockRoute("/api/rest/v1/file_actions/begin_upload/"+path, func(ctx *gin.Context, model interface{}) bool {
		beginUpload, ok := model.(files_sdk.FileBeginUploadParams)
		if !ok {
			t.Errorf("begin_upload bound %T", model)
			return false
		}
		part := beginUpload.Part
		if part == 0 {
			part = 1
		}
		ctx.JSON(http.StatusOK, files_sdk.FileUploadPartCollection{{
			HttpMethod:    "POST",
			Path:          path,
			Ref:           finalizationTestRef,
			UploadUri:     server.Server.URL + "/upload/" + path + "?part_number=" + strconv.FormatInt(part, 10),
			ParallelParts: lib.Bool(true),
			Expires:       time.Now().Add(time.Hour).Format(time.RFC3339),
			PartNumber:    part,
		}})
		return true
	})
	server.MockRoute("/api/rest/v1/files/"+path, func(ctx *gin.Context, model interface{}) bool {
		params, ok := model.(files_sdk.FileCreateParams)
		if !ok {
			return false // GET, the file's metadata
		}
		f.mutex.Lock()
		f.requests = append(f.requests, finalizationRequest{params: params, header: ctx.Request.Header.Clone(), host: ctx.Request.Host, at: time.Now()})
		answer := f.answers[min(len(f.requests), len(f.answers))-1]
		f.mutex.Unlock()
		answer(ctx)
		return true
	})
	return f
}

func (f *finalizationServer) finalizes() []finalizationRequest {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return append([]finalizationRequest(nil), f.requests...)
}

func (f *finalizationServer) setAnswers(answers ...func(*gin.Context)) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.requests = nil
	f.answers = answers
}

func (f *finalizationServer) partUploads() int {
	f.traceMutex.Lock()
	defer f.traceMutex.Unlock()
	return len(f.TrackRequest["/upload/*path"])
}

func (f *finalizationServer) upload(opts ...UploadOption) (UploadResumable, error) {
	return f.uploadWith(f.Client(), opts...)
}

func (f *finalizationServer) uploadWith(client *Client, opts ...UploadOption) (UploadResumable, error) {
	return client.UploadWithResume(append([]UploadOption{
		UploadWithReaderAt(bytes.NewReader(finalizationTestData)),
		UploadWithDestinationPath(f.path),
		UploadWithSize(int64(len(finalizationTestData))),
		UploadWithManager(lib.NewConstrainedWorkGroup(2)),
	}, opts...)...)
}

// pendingAnswer answers that Files.com is still finishing the upload, with progress when it is not nil.
func pendingAnswer(ref any, progress any) func(*gin.Context) {
	return func(ctx *gin.Context) {
		body := map[string]any{"status": "pending", "ref": ref}
		if progress != nil {
			body["progress"] = progress
		}
		ctx.JSON(http.StatusAccepted, body)
	}
}

func fileAnswer(path string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusCreated, files_sdk.File{Path: path, DisplayName: filepath.Base(path), Size: int64(len(finalizationTestData)), Type: "file"})
	}
}

func errorAnswer(code int, errorType files_sdk.ResponseErrorType, message string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		ctx.JSON(code, files_sdk.ResponseError{Type: string(errorType), Title: message, ErrorMessage: message, HttpCode: code})
	}
}

func TestUploadFinalizationSendsAsyncFinalizeAndReturnsAnImmediateFile(t *testing.T) {
	for _, engine := range []struct {
		name string
		opts []UploadOption
	}{{"v1", nil}, {"v2", []UploadOption{UploadWithV2()}}} {
		t.Run(engine.name, func(t *testing.T) {
			path := "immediate-" + engine.name + ".txt"
			server := newFinalizationServer(t, path, fileAnswer(path))

			resumable, err := server.upload(engine.opts...)

			require.NoError(t, err)
			assert.Equal(t, path, resumable.File.Path)
			assert.Equal(t, int64(len(finalizationTestData)), resumable.File.Size)
			finalizes := server.finalizes()
			require.Len(t, finalizes, 1)
			assert.Equal(t, "end", finalizes[0].params.Action)
			assert.Equal(t, finalizationTestRef, finalizes[0].params.Ref)
			assert.Equal(t, lib.Bool(true), finalizes[0].params.AsyncFinalize)
		})
	}
}

func TestUploadFinalizationReportsPendingAnswersUntilTheFile(t *testing.T) {
	path := "pending.txt"
	server := newFinalizationServer(t, path, pendingAnswer(finalizationTestRef, nil), pendingAnswer(finalizationTestRef, 42), pendingAnswer(finalizationTestRef, 100),
		pendingAnswer(finalizationTestRef, 100), fileAnswer(path))
	local := filepath.Join(t.TempDir(), "pending.txt")
	require.NoError(t, os.WriteFile(local, finalizationTestData, 0o600))
	client := server.Client()
	client.Config.APIKey = "finalization-test-key"
	client.Config.WorkspaceId = 7

	var reportsMutex sync.Mutex
	var reports []JobFile
	reporter := func(file JobFile) {
		reportsMutex.Lock()
		defer reportsMutex.Unlock()
		reports = append(reports, file)
	}
	job := client.Uploader(UploaderParams{
		LocalPath:      local,
		RemotePath:     path,
		EventsReporter: CreateFileEvents(reporter, OnBytesChange(status.Uploading), status.Complete),
	})
	job.Start()
	job.Wait()

	require.Len(t, job.Statuses, 1)
	uploaded, ok := job.Statuses[0].(*UploadStatus)
	require.True(t, ok)
	require.True(t, uploaded.Status().Is(status.Complete), "%v %v", uploaded.Status(), uploaded.Err())
	assert.Equal(t, path, uploaded.UploadResumable.File.Path)
	assert.Equal(t, int64(len(finalizationTestData)), uploaded.UploadResumable.File.Size)
	final := ToStatusFile(uploaded)
	assert.False(t, final.Finalizing)
	assert.Nil(t, final.FinalizationProgress)

	finalizes := server.finalizes()
	require.Len(t, finalizes, 5)
	assert.Equal(t, lib.Bool(true), finalizes[0].params.AsyncFinalize)
	// One, two, then three seconds after each pending answer, a pending 100 included.
	waits := []time.Duration{time.Second, 2 * time.Second, 3 * time.Second, 3 * time.Second}
	for i, request := range finalizes[1:] {
		assert.Equal(t, finalizes[0].params, request.params, "request %d is not the first one sent again", i+2)
		for _, header := range []string{"X-Filesapi-Key", "X-Files-Workspace-Id", "User-Agent"} {
			assert.NotEmpty(t, finalizes[0].header.Get(header), header)
			assert.Equal(t, finalizes[0].header.Values(header), request.header.Values(header), "request %d %s", i+2, header)
		}
		gap := request.at.Sub(finalizes[i].at)
		assert.True(t, gap >= waits[i] && gap < waits[i]+900*time.Millisecond, "the wait before request %d was %v, not %v", i+2, gap, waits[i])
	}

	// Another upload's finalization starts again from one second.
	again := newFinalizationServer(t, path, pendingAnswer(finalizationTestRef, nil), fileAnswer(path))
	_, err := again.upload()
	require.NoError(t, err)
	next := again.finalizes()
	require.Len(t, next, 2)
	gap := next[1].at.Sub(next[0].at)
	assert.True(t, gap >= time.Second && gap < time.Second+900*time.Millisecond, "the next upload's first wait was %v, not 1s", gap)

	reportsMutex.Lock()
	defer reportsMutex.Unlock()
	var progress []*int64
	var since []time.Time
	for _, report := range reports {
		if !report.Finalizing {
			continue
		}
		progress = append(progress, report.FinalizationProgress)
		since = append(since, report.FinalizingSince)
		assert.Equal(t, "uploading", report.StatusName)
		assert.Equal(t, int64(len(finalizationTestData)), report.TransferBytes, "a finalizing report changed the transferred bytes")
	}
	assert.Equal(t, []*int64{nil, lib.Int64(42), lib.Int64(100), lib.Int64(100)}, progress)
	require.Len(t, since, 4)
	assert.False(t, since[0].IsZero())
	assert.Equal(t, []time.Time{since[0], since[0], since[0], since[0]}, since, "FinalizingSince is the first pending answer's time")
	last := reports[len(reports)-1]
	assert.Equal(t, "complete", last.StatusName)
	assert.False(t, last.Finalizing)
	assert.True(t, last.FinalizingSince.IsZero())
	assert.Equal(t, int64(len(finalizationTestData)), last.TransferBytes)
}

func TestUploadFinalizationKeepsTheFirstRequestsConfig(t *testing.T) {
	path := "pending-config.txt"
	server := newFinalizationServer(t, path, pendingAnswer(finalizationTestRef, 42), fileAnswer(path))
	local := filepath.Join(t.TempDir(), "pending-config.txt")
	require.NoError(t, os.WriteFile(local, finalizationTestData, 0o600))
	client := server.Client()
	client.Config.APIKey = "first-key"
	client.Config.WorkspaceId = 7
	client.Config.UserAgent = "first-agent"
	client.Config.EndpointOverride = "https://first.example.com"
	client.Config.AdditionalHeaders = map[string]string{"X-Finalization-Test": "first"}
	unreachable := client.Config.SetCustomClient(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("the changed Config's client was used")
	})})

	changed := false
	reporter := func(f JobFile) {
		if !f.Finalizing || changed {
			return
		}
		changed = true
		// An upload job's reporter runs between the accepted request and the next one; it changes the Client's Config.
		client.Config.AdditionalHeaders["X-Finalization-Test"] = "changed"
		other := unreachable
		other.APIKey = "other-key"
		other.WorkspaceId = 99
		other.UserAgent = "other-agent"
		other.EndpointOverride = "https://other.example.com"
		client.Config = other
	}
	job := client.Uploader(UploaderParams{
		LocalPath:      local,
		RemotePath:     path,
		EventsReporter: CreateFileEvents(reporter, OnBytesChange(status.Uploading)),
	})
	job.Start()
	job.Wait()

	require.True(t, changed)
	require.Len(t, job.Statuses, 1)
	assert.True(t, job.Statuses[0].Status().Is(status.Complete), "%v %v", job.Statuses[0].Status().Name, job.Statuses[0].Err())
	finalizes := server.finalizes()
	require.Len(t, finalizes, 2)
	for i, request := range finalizes {
		assert.Equal(t, "first-key", request.header.Get("X-Filesapi-Key"), "request %d", i+1)
		assert.Equal(t, "7", request.header.Get("X-Files-Workspace-Id"), "request %d", i+1)
		assert.Equal(t, "first-agent", request.header.Get("User-Agent"), "request %d", i+1)
		assert.Equal(t, "first", request.header.Get("X-Finalization-Test"), "request %d", i+1)
		assert.Equal(t, "first.example.com", request.host, "request %d", i+1)
	}
}

// TestUploadFinalizationKeepsWhatTheFirstRequestResolved covers what Config resolves when a request is built rather than
// holding it: FILES_API_KEY and the default HTTP client. A change to either after Files.com accepted the first request
// must not reach the requests sent again.
func TestUploadFinalizationKeepsWhatTheFirstRequestResolved(t *testing.T) {
	path := "pending-resolved.txt"
	t.Run("FILES_API_KEY", func(t *testing.T) {
		t.Setenv("FILES_API_KEY", "first-env-key")
		server := newFinalizationServer(t, path, pendingAnswer(finalizationTestRef, nil), fileAnswer(path))
		client := server.Client()
		client.Config.APIKey, client.Config.SessionId = "", ""

		_, err := server.uploadWith(client, UploadWithFinalizationProgress(func(*int64) {
			require.NoError(t, os.Setenv("FILES_API_KEY", "other-env-key"))
		}))

		require.NoError(t, err)
		finalizes := server.finalizes()
		require.Len(t, finalizes, 2)
		for i, request := range finalizes {
			assert.Equal(t, "first-env-key", request.header.Get("X-Filesapi-Key"), "request %d", i+1)
		}
	})

	t.Run("a session ID keeps its precedence over FILES_API_KEY", func(t *testing.T) {
		t.Setenv("FILES_API_KEY", "first-env-key")
		server := newFinalizationServer(t, path, pendingAnswer(finalizationTestRef, nil), fileAnswer(path))
		client := server.Client()
		client.Config.APIKey, client.Config.SessionId = "", "first-session"

		_, err := server.uploadWith(client, UploadWithFinalizationProgress(func(*int64) {
			require.NoError(t, os.Setenv("FILES_API_KEY", "other-env-key"))
		}))

		require.NoError(t, err)
		finalizes := server.finalizes()
		require.Len(t, finalizes, 2)
		for i, request := range finalizes {
			assert.Equal(t, "first-session", request.header.Get("X-Filesapi-Auth"), "request %d", i+1)
			assert.Empty(t, request.header.Values("X-Filesapi-Key"), "request %d", i+1)
		}
	})

	t.Run("the default HTTP client", func(t *testing.T) {
		server := newFinalizationServer(t, path, pendingAnswer(finalizationTestRef, nil), fileAnswer(path))
		client := server.Client()
		original := lib.DefaultClient
		t.Cleanup(func() { lib.DefaultClient = original })
		// A Config without a client sends through lib.DefaultClient.
		lib.DefaultClient = client.Config.Client.HTTPClient
		client.Config.Client = nil

		_, err := server.uploadWith(client, UploadWithFinalizationProgress(func(*int64) {
			lib.DefaultClient = &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return nil, errors.New("the changed default client was used")
			})}
		}))

		require.NoError(t, err)
		assert.Len(t, server.finalizes(), 2)
	})
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// TestUploadFinalizationReportsProgressToADirectUpload covers UploadWithFinalizationProgress on Client.UploadWithResume:
// one call for each pending answer, with an invalid percentage reported as none, and the File as the result.
func TestUploadFinalizationReportsProgressToADirectUpload(t *testing.T) {
	path := "direct-progress.txt"
	server := newFinalizationServer(t, path, pendingAnswer(finalizationTestRef, 42.5), pendingAnswer(finalizationTestRef, 42), fileAnswer(path))
	var reported []*int64

	resumable, err := server.upload(UploadWithFinalizationProgress(func(percent *int64) { reported = append(reported, percent) }))

	require.NoError(t, err)
	assert.Equal(t, path, resumable.File.Path)
	assert.Equal(t, []*int64{nil, lib.Int64(42)}, reported)
	assert.Len(t, server.finalizes(), 3)
}

func TestPendingUploadFinalizationProgress(t *testing.T) {
	for _, c := range []struct {
		progress string
		want     *int64
	}{
		{`0`, lib.Int64(0)},
		{`42`, lib.Int64(42)},
		{`100`, lib.Int64(100)},
		{`101`, nil},
		{`-1`, nil},
		{`42.5`, nil},
		{`4.2e1`, nil},
		{`"42"`, nil},
		{`true`, nil},
		{`null`, nil},
		{`{}`, nil},
	} {
		t.Run(c.progress, func(t *testing.T) {
			percent, err := pendingUploadFinalization([]byte(`{"status":"pending","ref":"r","progress":`+c.progress+`}`), "r")

			require.NoError(t, err)
			assert.Equal(t, c.want, percent)
		})
	}
	percent, err := pendingUploadFinalization([]byte(`{"status":"pending","ref":"r"}`), "r")
	require.NoError(t, err)
	assert.Nil(t, percent)
}

func TestUploadFinalizationInvalidPendingAnswerIsCompletionUnknown(t *testing.T) {
	for _, c := range []struct {
		name   string
		answer func(*gin.Context)
	}{
		{"another status", func(ctx *gin.Context) {
			ctx.JSON(http.StatusAccepted, map[string]any{"status": "completed", "ref": finalizationTestRef})
		}},
		{"no status", func(ctx *gin.Context) { ctx.JSON(http.StatusAccepted, map[string]any{"ref": finalizationTestRef}) }},
		{"no ref", func(ctx *gin.Context) { ctx.JSON(http.StatusAccepted, map[string]any{"status": "pending"}) }},
		{"another ref", pendingAnswer("another-ref", 42)},
		{"a ref that is not a string", pendingAnswer(7, nil)},
		{"not an object", func(ctx *gin.Context) { ctx.Data(http.StatusAccepted, "application/json", []byte(`["pending"]`)) }},
		{"not JSON", func(ctx *gin.Context) { ctx.Data(http.StatusAccepted, "text/plain", []byte("Accepted")) }},
	} {
		t.Run(c.name, func(t *testing.T) {
			path := "invalid-pending.txt"
			server := newFinalizationServer(t, path, c.answer, fileAnswer(path))

			_, err := server.upload()

			var unknown *UploadCompletionUnknownError
			require.ErrorAs(t, err, &unknown)
			assert.Equal(t, finalizationTestRef, unknown.Ref)
			assert.Equal(t, path, unknown.Path)
			assert.Contains(t, err.Error(), "may still complete")
			assert.Len(t, server.finalizes(), 1, "an unreadable answer was followed by another request")
		})
	}
}

func TestUploadFinalizationErrorAfterPendingIsTheOutcome(t *testing.T) {
	for _, c := range []struct {
		name   string
		answer func(*gin.Context)
		check  func(*testing.T, error)
	}{
		{"failure", errorAnswer(http.StatusUnprocessableEntity, "processing-failure/remote-server-error", "The remote server rejected the file."), func(t *testing.T, err error) {
			assert.ErrorIs(t, err, files_sdk.ResponseErrorType("processing-failure/remote-server-error"))
		}},
		{"not found", errorAnswer(http.StatusNotFound, files_sdk.ErrFileUploadNotFound, "File Upload Not Found"), func(t *testing.T, err error) {
			assert.True(t, files_sdk.IsNotExist(err), err.Error())
		}},
	} {
		t.Run(c.name, func(t *testing.T) {
			path := "pending-error.txt"
			server := newFinalizationServer(t, path, pendingAnswer(finalizationTestRef, nil), c.answer)

			_, err := server.upload()

			require.Error(t, err)
			var responseError files_sdk.ResponseError
			assert.ErrorAs(t, err, &responseError)
			var unknown *UploadCompletionUnknownError
			assert.False(t, errors.As(err, &unknown), "Files.com's answer was reported as an unknown outcome: %v", err)
			c.check(t, err)
			assert.Len(t, server.finalizes(), 2)
		})
	}
}

func TestUploadFinalizationRetriesAFailedRepeatedRequest(t *testing.T) {
	t.Run("a dropped connection is sent again", func(t *testing.T) {
		path := "pending-retried.txt"
		dropConnection := func(ctx *gin.Context) { panic(http.ErrAbortHandler) }
		server := newFinalizationServer(t, path, pendingAnswer(finalizationTestRef, nil), dropConnection, fileAnswer(path))

		resumable, err := server.upload()

		require.NoError(t, err)
		assert.Equal(t, path, resumable.File.Path)
		finalizes := server.finalizes()
		require.Len(t, finalizes, 3)
		assert.Equal(t, finalizes[0].params, finalizes[2].params)
	})

	t.Run("server errors past the retries are an unknown outcome", func(t *testing.T) {
		path := "pending-unavailable.txt"
		server := newFinalizationServer(t, path, pendingAnswer(finalizationTestRef, nil), errorAnswer(http.StatusServiceUnavailable, "service-unavailable", "Service Unavailable"))

		_, err := server.upload()

		var unknown *UploadCompletionUnknownError
		require.ErrorAs(t, err, &unknown)
		assert.Equal(t, finalizationTestRef, unknown.Ref)
		// The first request, then the repeated one and the client's three retries of it.
		assert.Len(t, server.finalizes(), 5)
	})
}

func TestUploadFinalizationCancelAfterPendingStopsOnlyTheWait(t *testing.T) {
	t.Run("v1 context ends during the first request", func(t *testing.T) {
		path := "pending-cancel-v1.txt"
		ctx, cancel := context.WithCancelCause(context.Background())
		defer cancel(nil)
		server := newFinalizationServer(t, path, func(c *gin.Context) {
			cancel(ErrJobPaused) // the job is paused while Files.com accepts the request
			pendingAnswer(finalizationTestRef, nil)(c)
		}, pendingAnswer(finalizationTestRef, nil), fileAnswer(path))

		start := time.Now()
		resumable, err := server.upload(UploadWithContext(ctx))

		var unknown *UploadCompletionUnknownError
		require.ErrorAs(t, err, &unknown)
		assert.ErrorIs(t, err, context.Canceled)
		assert.ErrorIs(t, err, ErrJobPaused)
		assert.Contains(t, err.Error(), "may still complete")
		assert.Less(t, time.Since(start), 900*time.Millisecond, "the client waited after the context ended")
		require.Len(t, server.finalizes(), 1, "a request was sent after the context ended")
		first := server.finalizes()[0]

		// The upload's checkpoint asks for the outcome again rather than sending the parts again.
		assert.Equal(t, finalizationTestRef, resumable.FileUploadPart.Ref)
		require.Len(t, resumable.Parts, 1)
		assert.True(t, resumable.Parts[0].Successful())
		parts := server.partUploads()
		server.setAnswers(fileAnswer(path))

		resumed, err := server.upload(UploadWithResume(resumable))

		require.NoError(t, err)
		assert.Equal(t, path, resumed.File.Path)
		assert.Equal(t, parts, server.partUploads(), "the resumed upload sent its parts again")
		finalizes := server.finalizes()
		require.Len(t, finalizes, 1)
		assert.Equal(t, first.params, finalizes[0].params)
	})

	t.Run("a canceled upload job keeps the file pending in its checkpoint", func(t *testing.T) {
		path := "pending-cancel-job.txt"
		var job atomic.Pointer[Job]
		server := newFinalizationServer(t, path, func(c *gin.Context) {
			job.Load().Cancel()
			pendingAnswer(finalizationTestRef, nil)(c)
		}, pendingAnswer(finalizationTestRef, nil), fileAnswer(path))
		// A folder's contents, whose files a checkpoint resumes.
		folder := t.TempDir()
		local := filepath.Join(folder, path)
		require.NoError(t, os.WriteFile(local, finalizationTestData, 0o600))
		client := server.Client()

		canceled := client.Uploader(UploaderParams{LocalPath: folder + string(os.PathSeparator)})
		job.Store(canceled)
		canceled.Start()
		canceled.Wait()

		require.Len(t, canceled.Statuses, 1)
		assert.True(t, canceled.Statuses[0].Status().Is(status.Canceled), canceled.Statuses[0].Status().Name)
		ended := ToStatusFile(canceled.Statuses[0])
		assert.True(t, ended.Finalizing, "a file whose outcome is unknown stays finalizing")
		assert.False(t, ended.FinalizingSince.IsZero())
		require.Len(t, server.finalizes(), 1, "a request was sent after the job was canceled")
		first := server.finalizes()[0]
		checkpoint := canceled.UploadCheckpoint()
		assert.NotContains(t, checkpoint.CompletedPaths, local)
		require.Contains(t, checkpoint.PendingParts, local)
		assert.Equal(t, finalizationTestRef, checkpoint.PendingParts[local].FileUploadPart.Ref)
		parts := server.partUploads()
		server.setAnswers(fileAnswer(path))

		resumed := client.Uploader(UploaderParams{LocalPath: folder + string(os.PathSeparator), PriorJobCheckpoint: checkpoint})
		resumed.Start()
		resumed.Wait()

		require.Len(t, resumed.Statuses, 1)
		assert.True(t, resumed.Statuses[0].Status().Is(status.Complete), "%v %v", resumed.Statuses[0].Status().Name, resumed.Statuses[0].Err())
		assert.Equal(t, parts, server.partUploads(), "the resumed job sent the file's parts again")
		finalizes := server.finalizes()
		require.Len(t, finalizes, 1)
		assert.Equal(t, first.params, finalizes[0].params)
	})

	t.Run("v2 context ends during the wait", func(t *testing.T) {
		path := "pending-cancel-v2.txt"
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		server := newFinalizationServer(t, path, pendingAnswer(finalizationTestRef, 42), pendingAnswer(finalizationTestRef, 42), fileAnswer(path))

		resumable, err := server.upload(UploadWithV2(), UploadWithContext(ctx), UploadWithFinalizationProgress(func(*int64) { cancel() }))

		var unknown *UploadCompletionUnknownError
		require.ErrorAs(t, err, &unknown)
		assert.ErrorIs(t, err, context.Canceled)
		assert.Len(t, server.finalizes(), 1)
		assert.Equal(t, finalizationTestRef, resumable.FileUploadPart.Ref)
		require.Len(t, resumable.Parts, 1)
		assert.True(t, resumable.Parts[0].Successful())
	})

	t.Run("v2 context ends during a repeated request", func(t *testing.T) {
		path := "pending-cancel-request.txt"
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		server := newFinalizationServer(t, path, pendingAnswer(finalizationTestRef, nil), func(c *gin.Context) {
			cancel()
			select {
			case <-c.Request.Context().Done():
			case <-time.After(5 * time.Second):
				t.Errorf("the repeated request was not stopped")
			}
		}, fileAnswer(path))

		start := time.Now()
		resumable, err := server.upload(UploadWithV2(), UploadWithContext(ctx))

		var unknown *UploadCompletionUnknownError
		require.ErrorAs(t, err, &unknown)
		assert.ErrorIs(t, err, context.Canceled)
		assert.Less(t, time.Since(start), 3*time.Second)
		assert.Len(t, server.finalizes(), 2)
		assert.Equal(t, finalizationTestRef, resumable.FileUploadPart.Ref)
	})
}

func TestUploadFinalizationProgressPanicStopsTheWait(t *testing.T) {
	path := "pending-panic.txt"
	server := newFinalizationServer(t, path, pendingAnswer(finalizationTestRef, nil), fileAnswer(path))

	assert.PanicsWithValue(t, "progress callback", func() {
		_, _ = server.upload(UploadWithFinalizationProgress(func(*int64) { panic("progress callback") }))
	})
	assert.Len(t, server.finalizes(), 1)
}
