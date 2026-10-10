package files_sdk

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildRequestJSONBodyIsReplayable(t *testing.T) {
	expectedBody := `{"path":"reports/daily.csv","size":42}`
	params := lib.Params{Params: struct {
		Path string `json:"path"`
		Size int    `json:"size"`
	}{
		Path: "reports/daily.csv",
		Size: 42,
	}}

	tests := map[string]context.Context{
		"without context": nil,
		"with context":    context.Background(),
	}
	for name, requestContext := range tests {
		t.Run(name, func(t *testing.T) {
			request, err := buildRequest(&CallParams{
				Method:  http.MethodPost,
				Config:  Config{}.Init(),
				Uri:     "https://app.files.com/api/rest/v1/files/begin_upload",
				Params:  params,
				Headers: &http.Header{},
				Context: requestContext,
			})
			require.NoError(t, err)
			require.NotNil(t, request.GetBody)

			body, err := io.ReadAll(request.Body)
			require.NoError(t, err)
			assert.Equal(t, expectedBody, string(body))

			replayedBody, err := request.GetBody()
			require.NoError(t, err)
			defer replayedBody.Close()
			replayedPayload, err := io.ReadAll(replayedBody)
			require.NoError(t, err)
			assert.Equal(t, expectedBody, string(replayedPayload))
		})
	}
}

func TestBuildRequestSendsFileParamsAsReplayableMultipartForm(t *testing.T) {
	content := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x00, 0xff}
	params := struct {
		Paths                   []string  `url:"paths,omitempty" json:"paths,omitempty" path:"paths"`
		Description             string    `url:"description,omitempty" json:"description,omitempty" path:"description"`
		WatermarkAttachmentFile io.Writer `url:"watermark_attachment_file,omitempty" json:"watermark_attachment_file,omitempty" path:"watermark_attachment_file"`
	}{
		Paths:                   []string{"reports/a.csv", "reports/b.csv"},
		Description:             "Quarterly reports",
		WatermarkAttachmentFile: lib.NewAttachment("logo.png", "image/png", content),
	}
	callParams := func() *CallParams {
		return &CallParams{
			Method:  http.MethodPost,
			Config:  Config{}.Init(),
			Uri:     "https://app.files.com/api/rest/v1/bundles",
			Params:  lib.Params{Params: params},
			Headers: &http.Header{},
		}
	}

	request, err := buildRequest(callParams())
	require.NoError(t, err)
	require.NotNil(t, request.GetBody)
	mediaType, mediaParams, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
	require.NoError(t, err)
	assert.Equal(t, "multipart/form-data", mediaType)

	body, err := io.ReadAll(request.Body)
	require.NoError(t, err)
	replayedBody, err := request.GetBody()
	require.NoError(t, err)
	defer replayedBody.Close()
	replayedPayload, err := io.ReadAll(replayedBody)
	require.NoError(t, err)
	assert.Equal(t, body, replayedPayload)

	form, err := multipart.NewReader(bytes.NewReader(body), mediaParams["boundary"]).ReadForm(1 << 20)
	require.NoError(t, err)
	assert.Equal(t, []string{"Quarterly reports"}, form.Value["description"])
	assert.Equal(t, []string{"reports/a.csv"}, form.Value["paths[0]"])
	assert.Equal(t, []string{"reports/b.csv"}, form.Value["paths[1]"])
	require.Len(t, form.File["watermark_attachment_file"], 1)
	file := form.File["watermark_attachment_file"][0]
	assert.Equal(t, "logo.png", file.Filename)
	assert.Equal(t, "image/png", file.Header.Get("Content-Type"))
	opened, err := file.Open()
	require.NoError(t, err)
	defer opened.Close()
	sent, err := io.ReadAll(opened)
	require.NoError(t, err)
	assert.Equal(t, content, sent)

	params.WatermarkAttachmentFile = io.Discard
	_, err = buildRequest(callParams())
	assert.ErrorContains(t, err, "watermark_attachment_file")
}

func TestCallLogsDebugResponseAndKeepsItsBody(t *testing.T) {
	body := `[{"id":231,"description":"ordinary-debug-response-marker"}]`
	retriedBody := `{"error":"ordinary-debug-retried-marker","http-code":503}`
	// In the second case the server declares more bytes than it sends, so reading the body fails after the
	// valid JSON. In the third the server first answers 503, which the client retries.
	tests := map[string]struct {
		contentLength string
		readErr       error
		retried       bool
	}{
		"complete body":                        {},
		"body shorter than its Content-Length": {contentLength: strconv.Itoa(len(body) + 10), readErr: io.ErrUnexpectedEOF},
		"503 retried, then 200":                {retried: true},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if test.retried && requests.Add(1) == 1 {
					w.Header().Set("Content-Type", "application/json")
					w.Header().Set("X-Request-Id", "debug-retried-request")
					w.WriteHeader(http.StatusServiceUnavailable)
					w.Write([]byte(retriedBody))
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-Request-Id", "debug-response-request")
				if test.contentLength != "" {
					w.Header().Set("Content-Length", test.contentLength)
				}
				w.Write([]byte(body))
			}))
			defer server.Close()

			var logged bytes.Buffer
			config := Config{
				APIKey:           "debug-log-fixture",
				Debug:            true,
				Logger:           log.New(&logged, "", 0),
				EndpointOverride: server.URL,
			}.Init()
			// A response hook the caller supplied still runs, once for each response.
			var statuses []int
			if test.retried {
				config.Client.ResponseLogHook = func(_ retryablehttp.Logger, res *http.Response) {
					statuses = append(statuses, res.StatusCode)
				}
			}

			data, res, err := Call(http.MethodGet, config, "/bundles", nil)
			require.ErrorIs(t, err, test.readErr, "the caller gets the body's read result")
			defer res.Body.Close()

			assert.Equal(t, body, string(*data), "the caller still gets the body bytes")
			assert.Contains(t, logged.String(), "200 OK")
			assert.Contains(t, logged.String(), "X-Request-Id: debug-response-request")
			assert.Contains(t, logged.String(), body)
			if test.readErr != nil {
				assert.Contains(t, logged.String(), test.readErr.Error())
			}
			if test.retried {
				retried := strings.Index(logged.String(), "503 Service Unavailable")
				require.NotEqual(t, -1, retried, "the retried response is logged before it is discarded")
				assert.Contains(t, logged.String(), "X-Request-Id: debug-retried-request")
				assert.Contains(t, logged.String(), retriedBody)
				assert.Less(t, retried, strings.Index(logged.String(), "200 OK"), "responses are logged in the order they arrive")
				assert.Equal(t, 1, strings.Count(logged.String(), body), "the final response is logged once")
				assert.Equal(t, []int{http.StatusServiceUnavailable, http.StatusOK}, statuses)
			}
		})
	}
}

func TestCallRawSetsConfiguredUserAgent(t *testing.T) {
	var gotUserAgent, gotAPIKey string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUserAgent = r.UserAgent()
		gotAPIKey = r.Header.Get("X-FilesAPI-Key")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := Config{
		APIKey:    "secret",
		UserAgent: "Files.com Desktop Helper 1.2.3",
	}.Init()

	res, err := CallRaw(&CallParams{
		Method:  http.MethodGet,
		Config:  config,
		Uri:     server.URL,
		Headers: &http.Header{},
	})
	require.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, "Files.com Desktop Helper 1.2.3", gotUserAgent)
	assert.Empty(t, gotAPIKey)
}

func TestCallRawPreservesExplicitUserAgentHeader(t *testing.T) {
	var gotUserAgent string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUserAgent = r.UserAgent()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := Config{UserAgent: "Files.com Go SDK test"}.Init()
	headers := http.Header{}
	headers.Set("User-Agent", "Custom Transfer Client")

	res, err := CallRaw(&CallParams{
		Method:  http.MethodGet,
		Config:  config,
		Uri:     server.URL,
		Headers: &headers,
	})
	require.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, "Custom Transfer Client", gotUserAgent)
}

func TestCallRawDoesNotSendAuthHeadersAfterOffOriginRedirect(t *testing.T) {
	var firstRequestHeaders, redirectedRequestHeaders http.Header
	storageServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirectedRequestHeaders = r.Header.Clone()
		w.WriteHeader(http.StatusOK)
	}))
	defer storageServer.Close()

	filesServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		firstRequestHeaders = r.Header.Clone()
		http.Redirect(w, r, storageServer.URL+"/object", http.StatusFound)
	}))
	defer filesServer.Close()

	config := Config{
		APIKey:            "api-key",
		EndpointOverride:  filesServer.URL,
		AdditionalHeaders: map[string]string{"X-Files-Reauthentication": "password:additional-password"},
	}.Init()
	// The caller's own redirect policy still runs.
	redirects := 0
	config.Client.HTTPClient = &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error {
		redirects++
		return nil
	}}
	headers := &http.Header{}
	config.SetHeaders(headers)

	res, err := CallRaw(&CallParams{
		Method: http.MethodPost,
		Config: config,
		Uri:    config.RootPath() + "/zip_downloads",
		Params: lib.Params{Params: struct {
			Paths []string `json:"paths"`
		}{Paths: []string{"report.csv"}}},
		Headers: headers,
	})
	require.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, "api-key", firstRequestHeaders.Get("X-FilesAPI-Key"))
	assert.Equal(t, "password:additional-password", firstRequestHeaders.Get("X-Files-Reauthentication"))
	assert.Empty(t, redirectedRequestHeaders.Get("X-FilesAPI-Key"))
	assert.Empty(t, redirectedRequestHeaders.Get("X-Files-Reauthentication"))
	assert.Equal(t, 1, redirects)
}

func TestCallRawStopsAuthenticatedRedirectCycleAtDefaultLimit(t *testing.T) {
	var requests atomic.Int32
	filesServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		http.Redirect(w, r, r.URL.Path, http.StatusFound)
	}))
	defer filesServer.Close()

	config := Config{APIKey: "api-key", EndpointOverride: filesServer.URL}.Init()
	headers := &http.Header{}
	config.SetHeaders(headers)

	_, err := CallRaw(&CallParams{
		Method:  http.MethodGet,
		Config:  config,
		Uri:     config.RootPath() + "/zip_downloads",
		Headers: headers,
	})

	// With no caller redirect policy, net/http's default stops after 10 consecutive requests, and the client does
	// not retry that error.
	require.ErrorContains(t, err, "stopped after 10 redirects")
	assert.Equal(t, int32(10), requests.Load())
}

type retryLogTestTransport struct {
	attempts *atomic.Int32
	err      error
}

func (t retryLogTestTransport) RoundTrip(*http.Request) (*http.Response, error) {
	t.attempts.Add(1)
	return nil, t.err
}

// levelRecorder is a Logger with levels that keeps each entry with the level it was logged at.
type levelRecorder struct {
	entries []string
}

func (r *levelRecorder) Printf(format string, args ...any) { r.add("PRINTF", format, args...) }
func (r *levelRecorder) Error(format string, args ...any)  { r.add("ERROR", format, args...) }
func (r *levelRecorder) Warn(format string, args ...any)   { r.add("WARN", format, args...) }
func (r *levelRecorder) Info(format string, args ...any)   { r.add("INFO", format, args...) }
func (r *levelRecorder) Debug(format string, args ...any)  { r.add("DEBUG", format, args...) }
func (r *levelRecorder) Trace(format string, args ...any)  { r.add("TRACE", format, args...) }

func (r *levelRecorder) add(level string, format string, args ...any) {
	r.entries = append(r.entries, level+" "+fmt.Sprintf(format, args...))
}

// structuredRecorder is a go-retryablehttp LeveledLogger, whose entries are a message with key/value pairs.
type structuredRecorder struct {
	entries []string
}

func (r *structuredRecorder) Error(msg string, keysAndValues ...interface{}) {
	r.add("ERROR", msg, keysAndValues)
}
func (r *structuredRecorder) Info(msg string, keysAndValues ...interface{}) {
	r.add("INFO", msg, keysAndValues)
}
func (r *structuredRecorder) Debug(msg string, keysAndValues ...interface{}) {
	r.add("DEBUG", msg, keysAndValues)
}
func (r *structuredRecorder) Warn(msg string, keysAndValues ...interface{}) {
	r.add("WARN", msg, keysAndValues)
}

func (r *structuredRecorder) add(level string, msg string, keysAndValues []interface{}) {
	r.entries = append(r.entries, level+" "+msg+" "+fmt.Sprintf("%v", keysAndValues))
}

func TestCallRawKeepsSecretsToDebugLogEntries(t *testing.T) {
	signedURL := "https://transfer.example.test/object?X-Amz-Signature=signature-sentinel"
	secrets := []string{"signature-sentinel", "api-key-sentinel", "session-sentinel"}
	// The transport's error text carries the signed URL and the request's credentials, as a proxy's error might.
	cause := errors.New("connection refused for " + signedURL + " with key api-key-sentinel and session session-sentinel")
	// A logger without levels marks retryablehttp's error entries "[ERR]"; with Debug on, its other entries are
	// DEBUG diagnostics, except the request hook's own.
	tests := map[string]struct {
		leveled    bool
		structured bool
		debug      bool
	}{
		"logger without levels, debug off":   {},
		"logger without levels, debug on":    {debug: true},
		"logger with levels, debug on":       {leveled: true, debug: true},
		"structured client logger, debug on": {structured: true, debug: true},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var printed bytes.Buffer
			recorder := &levelRecorder{}
			structured := &structuredRecorder{}
			var logger lib.Logger = log.New(&printed, "", 0)
			if test.leveled {
				logger = recorder
			}
			config := Config{Logger: logger, Debug: test.debug}.Init()
			var attempts atomic.Int32
			config.Client.RetryMax = 1
			config.Client.HTTPClient = &http.Client{Transport: retryLogTestTransport{attempts: &attempts, err: cause}}
			config.Client.RequestLogHook = func(hookLog retryablehttp.Logger, req *http.Request, attempt int) {
				hookLog.Printf("sending attempt %d to %s", attempt, req.URL)
			}
			if test.structured {
				config.Client.Logger = structured
			}
			headers := &http.Header{}
			headers.Set("X-FilesAPI-Key", "api-key-sentinel")
			headers.Set("X-FilesAPI-Auth", "session-sentinel")

			_, err := CallRaw(&CallParams{Method: http.MethodGet, Config: config, Uri: signedURL, Headers: headers})

			require.ErrorIs(t, err, cause)
			assert.Equal(t, int32(2), attempts.Load(), "the configured retry still runs")
			if test.structured {
				assert.Same(t, structured, config.Client.Logger, "the shared client's logger is not replaced")
			} else {
				_, sharedLogger := config.Client.Logger.(Config)
				assert.True(t, sharedLogger, "the shared client's logger is not replaced")
			}

			entries := recorder.entries
			if test.structured {
				entries = structured.entries
			} else if !test.leveled {
				entries = strings.Split(strings.TrimSpace(printed.String()), "\n")
			}
			loggedFailure, loggedHook := false, false
			for _, entry := range entries {
				debugEntry := strings.HasPrefix(entry, "DEBUG ")
				if !test.leveled && !test.structured && test.debug {
					debugEntry = !strings.HasPrefix(entry, "[ERR] ") && !strings.HasPrefix(entry, "sending attempt")
				}
				if debugEntry {
					continue
				}
				for _, secret := range secrets {
					assert.NotContains(t, entry, secret)
				}
				loggedFailure = loggedFailure || (strings.Contains(entry, "request failed") && strings.Contains(entry, "connection refused"))
				loggedHook = loggedHook || strings.Contains(entry, "sending attempt")
			}
			assert.True(t, loggedFailure, "the failed request and its cause are still logged above DEBUG")
			assert.True(t, loggedHook, "the request hook's entries are still logged")
			if test.debug {
				assert.Contains(t, strings.Join(entries, "\n"), "signature-sentinel", "DEBUG entries keep the full URL")
			}
			if test.leveled {
				assert.NotContains(t, strings.Join(entries, "\n"), "PRINTF ", "entries go to their own level")
			}
		})
	}
}

func TestConfigDoReturnsNetHTTPErrorForRequestWithoutURL(t *testing.T) {
	config := Config{APIKey: "api-key"}.Init()

	_, err := config.Do(&http.Request{})

	require.ErrorContains(t, err, "http: nil Request.URL")
}

func TestSetHeadersUsesSessionIdWhenAPIKeyOnlyComesFromEnvironment(t *testing.T) {
	t.Setenv("FILES_API_KEY", "env-secret")

	headers := http.Header{}
	config := Config{SessionId: "session-secret"}.Init()
	config.SetHeaders(&headers)

	assert.Empty(t, headers.Get(apiKeyHeader))
	assert.Equal(t, "session-secret", headers.Get(sessionIdHeader))
}

func TestGetAPIKeyIgnoresEnvironmentWhenSessionIdIsConfigured(t *testing.T) {
	t.Setenv("FILES_API_KEY", "env-secret")

	config := Config{SessionId: "session-secret"}.Init()

	assert.Empty(t, config.GetAPIKey())
}

func TestSetHeadersUsesExplicitAPIKeyBeforeSessionId(t *testing.T) {
	t.Setenv("FILES_API_KEY", "env-secret")

	headers := http.Header{}
	config := Config{APIKey: "explicit-secret", SessionId: "session-secret"}.Init()
	config.SetHeaders(&headers)

	assert.Equal(t, "explicit-secret", headers.Get(apiKeyHeader))
	assert.Empty(t, headers.Get(sessionIdHeader))
}

func TestSetHeadersUsesConfiguredWorkspaceId(t *testing.T) {
	headers := http.Header{}
	config := Config{WorkspaceId: 123}.Init()
	config.SetHeaders(&headers)

	assert.Equal(t, "123", headers.Get(workspaceIdHeader))
}

func TestSetHeadersOmitsWorkspaceIdWhenUnset(t *testing.T) {
	headers := http.Header{}
	config := Config{}.Init()
	config.SetHeaders(&headers)

	assert.Empty(t, headers.Get(workspaceIdHeader))
}
