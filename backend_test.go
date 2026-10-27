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
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"testing/iotest"

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

func TestCallKeepsSuccessfulResponsesAndReleasesDiscardedErrorResponses(t *testing.T) {
	for _, attachments := range []bool{false, true} {
		for _, outcome := range []string{"success", "read failure", "policy cancellation", "redirect rejected"} {
			t.Run(fmt.Sprintf("attachments=%v/%s", attachments, outcome), func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				readErr := errors.New("response read interrupted")
				var source io.Reader = strings.NewReader(`{}`)
				if outcome == "read failure" {
					source = io.MultiReader(source, iotest.ErrReader(readErr))
				}
				body := &optionResponseBody{Reader: source}
				var sentResponse *http.Response
				config := Config{EndpointOverride: "https://files.test"}.Init().SetCustomClient(&http.Client{
					Transport: responseOptionTransport(func(req *http.Request) (*http.Response, error) {
						if req.Body != nil {
							_, _ = io.Copy(io.Discard, req.Body)
							_ = req.Body.Close()
						}
						status := http.StatusOK
						header := http.Header{"Content-Type": []string{"application/json"}}
						if outcome == "policy cancellation" {
							status = http.StatusInternalServerError
							cancel()
						}
						if outcome == "redirect rejected" {
							status = http.StatusTemporaryRedirect
							header.Set("Location", "https://files.test/redirected")
						}
						sentResponse = &http.Response{StatusCode: status, Body: body, ContentLength: 2, Header: header, Request: req}
						return sentResponse, nil
					}),
					CheckRedirect: func(*http.Request, []*http.Request) error { return errors.New("stopped after 10 redirects") },
				})
				config.RetryMax = 0
				params := lib.Params{}
				if attachments {
					params.Params = map[string]interface{}{"avatar_file": lib.NewJSONAttachment(strings.NewReader("a"), 1, "avatar.png", "image/png")}
				}
				data, response, err := Call(http.MethodPost, config, "/ownership", params, WithContext(ctx))
				switch outcome {
				case "success", "read failure":
					if outcome == "success" {
						require.NoError(t, err)
					} else {
						require.ErrorIs(t, err, readErr)
					}
					require.Same(t, sentResponse, response)
					require.NotNil(t, data)
					assert.Equal(t, `{}`, string(*data))
					assert.Zero(t, body.closes)
					require.NoError(t, response.Body.Close())
				case "policy cancellation":
					require.ErrorIs(t, err, context.Canceled)
					require.Nil(t, response)
					require.Nil(t, data)
				default:
					require.ErrorContains(t, err, "stopped after 10 redirects")
					require.Nil(t, response)
					require.Nil(t, data)
				}
				assert.Equal(t, 1, body.closes)
			})
		}
	}
}

func TestCallRawLeavesFinalResponseForItsCaller(t *testing.T) {
	for _, canceled := range []bool{false, true} {
		t.Run(map[bool]string{false: "retry exhaustion without transport error", true: "policy cancellation"}[canceled], func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			body := &optionResponseBody{Reader: strings.NewReader("response")}
			config := Config{}.Init().SetCustomClient(&http.Client{Transport: responseOptionTransport(func(req *http.Request) (*http.Response, error) {
				if canceled {
					cancel()
				}
				return &http.Response{StatusCode: http.StatusInternalServerError, Header: http.Header{}, Body: body, Request: req}, nil
			})})
			config.RetryMax = 0
			response, err := CallRaw(&CallParams{Method: http.MethodGet, Config: config, Uri: "https://upload.test/file", Context: ctx})
			if canceled {
				require.ErrorIs(t, err, context.Canceled)
			} else {
				require.NoError(t, err)
			}
			require.NotNil(t, response)
			assert.Zero(t, body.closes)
			require.NoError(t, response.Body.Close())
			assert.Equal(t, 1, body.closes)
		})
	}
}

func TestCustomErrorHandlersKeepBodyOwnershipAndReplacementDiagnostics(t *testing.T) {
	for _, sender := range []string{"standard", "attachments", "raw"} {
		for _, owner := range []string{"caller", "handler", "replacement"} {
			t.Run(sender+"/"+owner, func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				handlerErr := errors.New("custom handler rejected response")
				closeErr := errors.New("original close error")
				original := &optionResponseBody{Reader: strings.NewReader("original"), closeErr: closeErr}
				replacementBody := &optionResponseBody{Reader: strings.NewReader("replacement diagnostics")}
				replacement := &http.Response{StatusCode: http.StatusBadGateway, Header: http.Header{"X-Diagnostic": []string{"retained"}}, Body: replacementBody}
				config := Config{EndpointOverride: "https://files.test"}.Init().SetCustomClient(&http.Client{Transport: responseOptionTransport(func(req *http.Request) (*http.Response, error) {
					if req.Body != nil {
						_, _ = io.Copy(io.Discard, req.Body)
						_ = req.Body.Close()
					}
					cancel()
					return &http.Response{StatusCode: http.StatusInternalServerError, Body: original, Header: http.Header{}, Request: req}, nil
				})})
				config.RetryMax = 0
				config.Client.ErrorHandler = func(response *http.Response, err error, _ int) (*http.Response, error) {
					require.ErrorIs(t, err, context.Canceled)
					if owner != "caller" {
						require.ErrorIs(t, response.Body.Close(), closeErr)
					}
					if owner == "replacement" {
						require.NoError(t, replacementBody.Close())
						return replacement, handlerErr
					}
					return response, handlerErr
				}
				var response *http.Response
				var err error
				if sender == "raw" {
					response, err = CallRaw(&CallParams{Method: http.MethodGet, Config: config, Uri: "https://files.test/file", Context: ctx})
					require.Same(t, handlerErr, err)
					require.NotNil(t, response)
					if owner == "replacement" {
						require.Same(t, replacement, response)
						require.Same(t, replacementBody, response.Body)
						assert.Equal(t, "retained", response.Header.Get("X-Diagnostic"))
						data, readErr := io.ReadAll(response.Body)
						require.NoError(t, readErr)
						assert.Equal(t, "replacement diagnostics", string(data))
					} else if owner == "caller" {
						assert.Zero(t, original.closes)
						require.ErrorIs(t, response.Body.Close(), closeErr)
					} else {
						require.NoError(t, response.Body.Close())
					}
				} else {
					params := lib.Params{}
					if sender == "attachments" {
						params.Params = map[string]interface{}{"avatar_file": lib.NewJSONAttachment(strings.NewReader("a"), 1, "avatar.png", "image/png")}
					}
					_, response, err = Call(http.MethodPost, config, "/ownership", params, WithContext(ctx))
					require.ErrorIs(t, err, handlerErr)
					require.Nil(t, response)
				}
				assert.Equal(t, 1, original.closes)
				if owner == "replacement" {
					assert.Equal(t, 1, replacementBody.closes)
				}
			})
		}
	}
}

func TestCallRawKeepsAlreadyDrainedResponseWhenRetryPreparationFails(t *testing.T) {
	for _, failure := range []string{"next body", "prepare retry"} {
		t.Run(failure, func(t *testing.T) {
			retryErr := errors.New("cannot prepare another upload attempt")
			body := &optionResponseBody{Reader: strings.NewReader("failed response")}
			requests, handlerCalls := 0, 0
			var sentResponse *http.Response
			config := Config{}.Init().SetCustomClient(&http.Client{Transport: responseOptionTransport(func(req *http.Request) (*http.Response, error) {
				if req.Body != nil {
					_, _ = io.Copy(io.Discard, req.Body)
					_ = req.Body.Close()
				}
				requests++
				sentResponse = &http.Response{StatusCode: http.StatusInternalServerError, Header: http.Header{}, Body: body, Request: req}
				return sentResponse, nil
			})})
			config.RetryMax = 1
			config.Client.ErrorHandler = func(response *http.Response, err error, attempts int) (*http.Response, error) {
				handlerCalls++
				return retryablehttp.PassthroughErrorHandler(response, err, attempts)
			}
			if failure == "prepare retry" {
				config.Client.PrepareRetry = func(*http.Request) error { return retryErr }
			}
			response, err := CallRaw(&CallParams{Method: http.MethodPut, Config: config, Uri: "https://upload.test/file", Params: lib.Params{}, RetryableBody: retryablehttp.ReaderFunc(func() (io.Reader, error) {
				if failure == "next body" && requests > 0 {
					return nil, retryErr
				}
				return strings.NewReader("a"), nil
			})})
			require.ErrorIs(t, err, retryErr)
			require.Same(t, sentResponse, response)
			require.Equal(t, 1, requests)
			assert.Equal(t, 1, body.closes, "the retry dependency already closed the previous response")
			if failure == "next body" {
				assert.Zero(t, handlerCalls)
			} else {
				assert.Equal(t, 1, handlerCalls)
			}
			require.NoError(t, response.Body.Close())
			assert.Equal(t, 1, body.closes)
		})
	}
}

func TestPolicyURLErrorDoesNotMarkAnOpenResponseAsClosed(t *testing.T) {
	for _, sender := range []string{"standard", "attachments", "raw"} {
		t.Run(sender, func(t *testing.T) {
			policyCause := errors.New("policy rejected response")
			body := &optionResponseBody{Reader: strings.NewReader("response")}
			config := Config{EndpointOverride: "https://files.test"}.Init().SetCustomClient(&http.Client{Transport: responseOptionTransport(func(req *http.Request) (*http.Response, error) {
				if req.Body != nil {
					_, _ = io.Copy(io.Discard, req.Body)
					_ = req.Body.Close()
				}
				return &http.Response{StatusCode: http.StatusOK, Header: http.Header{}, Body: body, Request: req}, nil
			})})
			config.Client.CheckRetry = func(_ context.Context, _ *http.Response, doErr error) (bool, error) {
				require.NoError(t, doErr)
				return false, &url.Error{Op: "policy", URL: "https://files.test/file", Err: policyCause}
			}
			var response *http.Response
			var err error
			if sender == "raw" {
				response, err = CallRaw(&CallParams{Method: http.MethodGet, Config: config, Uri: "https://files.test/file"})
				require.NotNil(t, response)
				assert.Zero(t, body.closes)
				require.NoError(t, response.Body.Close())
			} else {
				params := lib.Params{}
				if sender == "attachments" {
					params.Params = map[string]interface{}{"avatar_file": lib.NewJSONAttachment(strings.NewReader("a"), 1, "avatar.png", "image/png")}
				}
				_, response, err = Call(http.MethodPost, config, "/ownership", params)
				require.Nil(t, response)
			}
			require.ErrorIs(t, err, policyCause)
			assert.Equal(t, 1, body.closes)
		})
	}
}

// regionTestContextKey marks a request's context, so that a test can see each attempt carry it.
type regionTestContextKey struct{}

// regionTestSent is one request a regionTestSites transport received, as the site it was sent to saw it.
type regionTestSent struct {
	host, method, path, query       string
	key, session, reauth, workspace string
	tag                             string
	ctxValue                        interface{}
	body                            string
}

// regionTestSites answers each request with the handler for its host, as separate sites would, without changing the
// request the SDK sent, and records each request.
type regionTestSites struct {
	handlers map[string]func(*http.Request) (*http.Response, error)
	sent     []regionTestSent
}

func (s *regionTestSites) RoundTrip(req *http.Request) (*http.Response, error) {
	var body []byte
	if req.Body != nil {
		body, _ = io.ReadAll(req.Body)
		_ = req.Body.Close()
	}
	s.sent = append(s.sent, regionTestSent{
		host:      req.URL.Host,
		method:    req.Method,
		path:      req.URL.EscapedPath(),
		query:     req.URL.RawQuery,
		key:       req.Header.Get(apiKeyHeader),
		session:   req.Header.Get(sessionIdHeader),
		reauth:    req.Header.Get(reauthenticationHeader),
		workspace: req.Header.Get(workspaceIdHeader),
		tag:       req.Header.Get("X-Request-Tag"),
		ctxValue:  req.Context().Value(regionTestContextKey{}),
		body:      string(body),
	})
	handler, ok := s.handlers[req.URL.Host]
	if !ok {
		return nil, fmt.Errorf("no site at %s", req.URL.Host)
	}
	return handler(req)
}

func (s *regionTestSites) hosts() []string {
	hosts := make([]string, 0, len(s.sent))
	for _, sent := range s.sent {
		hosts = append(hosts, sent.host)
	}
	return hosts
}

// regionTestConfig is a configuration for https://us.files.com, with the SDK's default retry client, whose requests go
// to sites.
func regionTestConfig(sites *regionTestSites) Config {
	return Config{APIKey: "region-key", EndpointOverride: "https://us.files.com"}.Init().SetCustomClient(&http.Client{Transport: sites})
}

// regionTestPayload is the API's typed lockout region mismatch naming host.
func regionTestPayload(host string) string {
	return fmt.Sprintf(`{"type":"not-authenticated/lockout-region-mismatch","http-code":403,"error":"Lockout region mismatch","data":{"host":%q}}`, host)
}

// regionTestAnswer answers with status, contentType (none when empty), payload and header.
func regionTestAnswer(status int, contentType, payload string, header http.Header) func(*http.Request) (*http.Response, error) {
	return func(req *http.Request) (*http.Response, error) {
		answer := http.Header{}
		for name, values := range header {
			answer[name] = append([]string(nil), values...)
		}
		if contentType != "" {
			answer.Set("Content-Type", contentType)
		}
		return &http.Response{StatusCode: status, Status: fmt.Sprintf("%d %s", status, http.StatusText(status)), Header: answer, Body: io.NopCloser(strings.NewReader(payload)), Request: req}, nil
	}
}

// regionTestMismatch answers with the typed lockout region mismatch naming host.
func regionTestMismatch(host string) func(*http.Request) (*http.Response, error) {
	return regionTestAnswer(http.StatusForbidden, "application/json", regionTestPayload(host), nil)
}

var regionTestOK = regionTestAnswer(http.StatusOK, "application/json", `{}`, nil)

// regionTestBackoff returns backoff with observe called first for each wait it is asked for; when observe returns
// true, the wait is the longest the client allows instead.
func regionTestBackoff[D any](backoff func(D, D, int, *http.Response) D, observe func(attempt int, resp *http.Response) bool) func(D, D, int, *http.Response) D {
	return func(min, max D, attempt int, resp *http.Response) D {
		if observe(attempt, resp) {
			return max
		}
		return backoff(min, max, attempt, resp)
	}
}

// A lockout region mismatch moves each kind of API request to the host it names, sending the same request again: its
// method, escaped path, query, body, context, request-option headers and credentials, with each request and response
// option applied once.
func TestRegionMismatchMovesEachSenderWithTheSameRequest(t *testing.T) {
	for _, sender := range []string{"standard GET", "standard POST", "attachments", "raw API without a body", "raw API with a body factory"} {
		t.Run(sender, func(t *testing.T) {
			sites := &regionTestSites{handlers: map[string]func(*http.Request) (*http.Response, error){
				"us.files.com": regionTestMismatch("eu.files.com"),
				"eu.files.com": regionTestOK,
			}}
			config := regionTestConfig(sites)
			ctx := context.WithValue(context.Background(), regionTestContextKey{}, "kept")
			requestOptions, responseOptions := 0, 0
			opts := []RequestResponseOption{
				WithContext(ctx),
				RequestHeadersOption(&http.Header{"X-Request-Tag": []string{"tagged"}}),
				RequestOption(func(*http.Request) error { requestOptions++; return nil }),
				ResponseOption(func(*http.Response) error { responseOptions++; return nil }),
			}
			raw := strings.HasPrefix(sender, "raw")
			var response *http.Response
			var err error
			switch sender {
			case "standard GET":
				_, response, err = Call(http.MethodGet, config, "/files/report%20a.csv?cursor=c%201", nil, opts...)
			case "standard POST":
				_, response, err = Call(http.MethodPost, config, "/folders", lib.Params{Params: map[string]interface{}{"path": "Reports"}}, opts...)
			case "attachments":
				_, response, err = Call(http.MethodPost, config, "/users/7", lib.Params{Params: map[string]interface{}{
					"name":        "Reports",
					"avatar_file": lib.NewJSONAttachment(strings.NewReader("avatar"), 6, "avatar.png", "image/png"),
				}}, opts...)
			default:
				headers := &http.Header{"X-Request-Tag": []string{"tagged"}}
				config.SetHeaders(headers)
				params := &CallParams{Method: http.MethodGet, Config: config, Uri: config.RootPath() + "/files/report%20a.csv?cursor=c%201", Headers: headers, Context: ctx}
				if sender == "raw API with a body factory" {
					body := lib.Params{Params: map[string]interface{}{"paths": []string{}, "encoded_paths": []string{"a%20b"}}}
					params.Method = http.MethodPost
					params.Params = body
					params.RetryableBody = func() (io.Reader, error) { return body.ToJSON() }
				}
				response, err = CallRawAPI(params)
				assert.True(t, params.RelocatedAPIContext().Relocated())
				assert.Equal(t, "https://eu.files.com", params.RelocatedAPIContext().Origin())
			}
			require.NoError(t, err)
			require.NotNil(t, response)
			require.NoError(t, response.Body.Close())

			require.Equal(t, []string{"us.files.com", "eu.files.com"}, sites.hosts())
			first, moved := sites.sent[0], sites.sent[1]
			assert.Equal(t, first.method, moved.method)
			assert.Equal(t, first.path, moved.path)
			assert.Equal(t, first.query, moved.query)
			assert.Equal(t, first.body, moved.body, "the same body is sent again")
			for _, sent := range sites.sent {
				assert.Equal(t, "region-key", sent.key)
				assert.Equal(t, "tagged", sent.tag)
				assert.Equal(t, "kept", sent.ctxValue)
			}
			if sender == "standard GET" || sender == "raw API without a body" {
				assert.Equal(t, "/api/rest/v1/files/report%20a.csv", moved.path)
				assert.Equal(t, "cursor=c%201", moved.query)
				assert.Empty(t, moved.body)
			} else {
				assert.NotEmpty(t, moved.body)
			}
			if !raw {
				assert.Equal(t, 1, requestOptions)
				assert.Equal(t, 1, responseOptions)
			}
		})
	}
}

// The configured retry policy sees each attempt once, a followable mismatch as the real 403 with a *RegionMismatch, and
// its refusal or error keeps the request from moving. Backoff, PrepareRetry, which sees the moved URL, and the request
// hook run as for any retry, and a PrepareRetry failure keeps its own error.
func TestRegionMismatchReachesTheConfiguredRetryHooks(t *testing.T) {
	policyErr := errors.New("policy refused the move")
	prepareErr := errors.New("cannot sign the moved request")
	for _, decision := range []string{"follow", "veto", "policy error", "prepare error"} {
		t.Run(decision, func(t *testing.T) {
			sites := &regionTestSites{handlers: map[string]func(*http.Request) (*http.Response, error){
				"us.files.com": regionTestMismatch("eu.files.com"),
				"eu.files.com": regionTestOK,
			}}
			config := regionTestConfig(sites)
			type check struct {
				status int
				cause  *RegionMismatch
			}
			var checks []check
			defaultPolicy := config.Client.CheckRetry
			config.Client.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
				var cause *RegionMismatch
				errors.As(err, &cause)
				checks = append(checks, check{status: resp.StatusCode, cause: cause})
				if cause != nil {
					switch decision {
					case "veto":
						return false, nil
					case "policy error":
						return false, policyErr
					}
				}
				return defaultPolicy(ctx, resp, err)
			}
			var waits []string
			config.Client.Backoff = regionTestBackoff(config.Client.Backoff, func(attempt int, resp *http.Response) bool {
				waits = append(waits, fmt.Sprintf("%d %d", attempt, resp.StatusCode))
				return false
			})
			var sends []string
			config.Client.RequestLogHook = func(_ retryablehttp.Logger, req *http.Request, attempt int) {
				sends = append(sends, fmt.Sprintf("%d %s", attempt, req.URL.Host))
			}
			var prepared []string
			config.Client.PrepareRetry = func(req *http.Request) error {
				prepared = append(prepared, req.URL.Host)
				if decision == "prepare error" {
					return prepareErr
				}
				return nil
			}

			_, response, err := Call(http.MethodGet, config, "/files/report.csv", nil)

			require.NotEmpty(t, checks)
			assert.Equal(t, http.StatusForbidden, checks[0].status, "the policy sees the real 403")
			require.NotNil(t, checks[0].cause, "the policy sees the mismatch it may follow")
			assert.Equal(t, "https://eu.files.com", checks[0].cause.Origin)
			assert.Equal(t, "https://us.files.com", checks[0].cause.From)
			assert.Empty(t, checks[0].cause.Reason)
			var mismatch *RegionMismatchError
			switch decision {
			case "follow":
				require.NoError(t, err)
				require.NoError(t, response.Body.Close())
				assert.Equal(t, []string{"us.files.com", "eu.files.com"}, sites.hosts())
				require.Len(t, checks, 2)
				assert.Equal(t, http.StatusOK, checks[1].status)
				assert.Nil(t, checks[1].cause)
				assert.Equal(t, []string{"0 403"}, waits, "the move waits as the configured Backoff says, at its attempt")
				assert.Equal(t, []string{"0 us.files.com", "1 eu.files.com"}, sends)
				assert.Equal(t, []string{"eu.files.com"}, prepared, "PrepareRetry sees the moved URL")
			case "veto":
				require.ErrorAs(t, err, &mismatch)
				assert.Equal(t, RegionMismatchRetryPolicy, mismatch.Reason)
				require.NotNil(t, response)
				require.NoError(t, response.Body.Close())
				assert.Equal(t, []string{"us.files.com"}, sites.hosts())
				assert.Len(t, checks, 1)
				assert.Empty(t, waits)
				assert.Empty(t, prepared)
			case "policy error":
				require.ErrorIs(t, err, policyErr)
				require.ErrorAs(t, err, &mismatch)
				assert.Equal(t, RegionMismatchRetryPolicy, mismatch.Reason)
				assert.Nil(t, response)
				assert.Equal(t, []string{"us.files.com"}, sites.hosts())
				assert.Len(t, checks, 1)
				assert.Empty(t, waits)
				assert.Empty(t, prepared)
			case "prepare error":
				require.ErrorIs(t, err, prepareErr)
				assert.False(t, errors.As(err, &mismatch), "a later failure carries no earlier mismatch")
				assert.Nil(t, response)
				assert.Equal(t, []string{"us.files.com"}, sites.hosts())
				assert.Equal(t, []string{"0 403"}, waits)
				assert.Equal(t, []string{"eu.files.com"}, prepared)
			}
		})
	}
}

// A recognized mismatch that ends a call, here for want of a retry, keeps the error and the bodies each sender and
// error handler leave: the default and a handler keeping the 403 return it with the mismatch, the retry client's own
// handler drains it, and a handler's error stays the call's. A response the handler selects is described as it now is:
// a success, another 403, or the mismatch response once the handler has changed it so it is no longer a mismatch, is
// the outcome without the old mismatch, whether or not the handler returns its own error beside it, while one that is
// still a mismatch, such as the mismatch response whose X-Files-Host the handler changed, is the mismatch of the host
// it now names.
func TestRegionMismatchTerminalOutcomesKeepErrorsAndBodies(t *testing.T) {
	handlerErr := errors.New("custom handler rejected the mismatch")
	for _, sender := range []string{"standard", "attachments", "raw"} {
		for _, handler := range []string{"default", "nil", "same", "replacement", "unrelated", "error", "unrelated with an error", "same with another Content-Type", "same with its X-Files-Host changed in place"} {
			t.Run(sender+"/"+handler, func(t *testing.T) {
				mismatchPayload, mismatchHeader := regionTestPayload("eu.files.com"), http.Header{"Content-Type": []string{"application/json"}}
				if handler == "same with its X-Files-Host changed in place" {
					// The mismatch names its host only in its X-Files-Host header.
					mismatchPayload = regionTestPayload("")
					mismatchHeader.Set("X-Files-Host", "eu.files.com")
				}
				original := &optionResponseBody{Reader: strings.NewReader(mismatchPayload)}
				var replacementBody *optionResponseBody
				sites := &regionTestSites{handlers: map[string]func(*http.Request) (*http.Response, error){
					"us.files.com": func(req *http.Request) (*http.Response, error) {
						return &http.Response{StatusCode: http.StatusForbidden, Status: "403 Forbidden", Header: mismatchHeader, Body: original, Request: req}, nil
					},
				}}
				config := regionTestConfig(sites)
				config.RetryMax = 0
				if handler != "default" {
					// A policy that retries any 403 reaches the error handler once no retry is left.
					config.Client.CheckRetry = func(context.Context, *http.Response, error) (bool, error) { return true, nil }
				}
				switch handler {
				case "nil":
					config.Client.ErrorHandler = nil
				case "same":
					config.Client.ErrorHandler = func(resp *http.Response, _ error, _ int) (*http.Response, error) { return resp, nil }
				case "replacement", "unrelated", "unrelated with an error":
					status, payload := http.StatusOK, `{}`
					if handler != "replacement" {
						status, payload = http.StatusForbidden, `{"type":"not-authenticated/invalid-credentials","http-code":403,"error":"Invalid credentials"}`
					}
					config.Client.ErrorHandler = func(resp *http.Response, _ error, _ int) (*http.Response, error) {
						if err := resp.Body.Close(); err != nil {
							return nil, err
						}
						replacementBody = &optionResponseBody{Reader: strings.NewReader(payload)}
						replacement := &http.Response{StatusCode: status, Status: fmt.Sprintf("%d %s", status, http.StatusText(status)), Header: http.Header{"Content-Type": []string{"application/json"}}, Body: replacementBody, Request: resp.Request}
						if handler == "unrelated with an error" {
							return replacement, handlerErr
						}
						return replacement, nil
					}
				case "error":
					config.Client.ErrorHandler = func(resp *http.Response, _ error, _ int) (*http.Response, error) { return resp, handlerErr }
				case "same with another Content-Type":
					config.Client.ErrorHandler = func(resp *http.Response, _ error, _ int) (*http.Response, error) {
						resp.Header.Set("Content-Type", "text/html")
						return resp, nil
					}
				case "same with its X-Files-Host changed in place":
					config.Client.ErrorHandler = func(resp *http.Response, _ error, _ int) (*http.Response, error) {
						// The handler changes the host the response names in place, in the slice its header holds.
						resp.Header.Values(regionHostHeader)[0] = "eu.example.com"
						return resp, nil
					}
				}
				var response *http.Response
				var err error
				if sender == "raw" {
					headers := &http.Header{}
					config.SetHeaders(headers)
					response, err = CallRawAPI(&CallParams{Method: http.MethodGet, Config: config, Uri: config.RootPath() + "/files/report.csv", Headers: headers})
				} else {
					params := lib.Params{}
					if sender == "attachments" {
						params.Params = map[string]interface{}{"avatar_file": lib.NewJSONAttachment(strings.NewReader("a"), 1, "avatar.png", "image/png")}
					}
					_, response, err = Call(http.MethodPost, config, "/files", params)
				}

				assert.Equal(t, []string{"us.files.com"}, sites.hosts(), "nothing is sent to the named host")
				var mismatch *RegionMismatchError
				var responseError ResponseError
				switch handler {
				case "default", "same":
					require.ErrorAs(t, err, &mismatch)
					assert.Equal(t, RegionMismatchRetryBudget, mismatch.Reason)
					assert.Equal(t, "https://eu.files.com", mismatch.Origin)
					require.ErrorAs(t, err, &responseError)
					assert.Equal(t, regionMismatchType, responseError.Type)
					require.NotNil(t, response, "the mismatch response is the caller's")
					if sender == "raw" {
						data, readErr := io.ReadAll(response.Body)
						require.NoError(t, readErr)
						assert.JSONEq(t, regionTestPayload("eu.files.com"), string(data), "its body is still readable")
					}
					assert.Zero(t, original.closes)
					require.NoError(t, response.Body.Close())
				case "nil":
					require.ErrorAs(t, err, &mismatch)
					assert.Equal(t, RegionMismatchRetryBudget, mismatch.Reason)
					assert.ErrorContains(t, err, "giving up after 1 attempt(s)")
					require.ErrorAs(t, err, &responseError, "the API's error for the mismatch stays reachable")
					assert.Nil(t, response, "the retry client drained the mismatch response")
				case "replacement":
					require.NoError(t, err)
					require.NotNil(t, response)
					assert.Same(t, replacementBody, response.Body)
					require.NoError(t, response.Body.Close())
					assert.Equal(t, 1, replacementBody.closes)
				case "unrelated":
					assert.False(t, errors.As(err, &mismatch), "another 403 the handler selected is not the mismatch")
					if sender == "raw" {
						require.NoError(t, err)
					} else {
						require.ErrorAs(t, err, &responseError)
						assert.Equal(t, "not-authenticated/invalid-credentials", responseError.Type)
					}
					require.NotNil(t, response)
					require.NoError(t, response.Body.Close())
					assert.Equal(t, 1, replacementBody.closes)
				case "error":
					require.ErrorIs(t, err, handlerErr)
					require.ErrorAs(t, err, &mismatch)
					assert.Equal(t, RegionMismatchRetryBudget, mismatch.Reason)
					if sender == "raw" {
						require.NotNil(t, response, "a raw caller owns the response the handler returned")
						assert.Zero(t, original.closes)
						require.NoError(t, response.Body.Close())
					} else {
						assert.Nil(t, response)
					}
				case "unrelated with an error":
					require.ErrorIs(t, err, handlerErr, "the handler's error stays the call's")
					assert.False(t, errors.As(err, &mismatch), "the mismatch does not describe another response the handler selected beside its error")
					assert.False(t, errors.As(err, &responseError), "nor does the API's error for the mismatch")
					if sender == "raw" {
						require.Same(t, handlerErr, err)
						require.NotNil(t, response, "a raw caller owns the response the handler returned")
						data, readErr := io.ReadAll(response.Body)
						require.NoError(t, readErr)
						assert.JSONEq(t, `{"type":"not-authenticated/invalid-credentials","http-code":403,"error":"Invalid credentials"}`, string(data), "its diagnostics are still readable")
						require.NoError(t, response.Body.Close())
						assert.Equal(t, 1, replacementBody.closes)
					} else {
						assert.Nil(t, response)
						assert.Zero(t, replacementBody.closes, "a replacement returned beside the handler's error stays the handler's to close")
					}
				case "same with another Content-Type":
					assert.False(t, errors.As(err, &mismatch), "a mismatch response the handler changed is no longer the mismatch")
					if sender == "raw" {
						require.NoError(t, err)
					} else {
						require.Error(t, err)
						assert.False(t, errors.As(err, &responseError), "an HTML 403 is not the API's error")
					}
					require.NotNil(t, response)
					require.NoError(t, response.Body.Close())
				case "same with its X-Files-Host changed in place":
					require.ErrorAs(t, err, &mismatch, "a mismatch response the handler changed is the mismatch of the host it now names")
					assert.Equal(t, "eu.example.com", mismatch.Host)
					assert.Equal(t, "https://eu.example.com", mismatch.Origin)
					assert.Equal(t, "https://us.files.com", mismatch.From)
					assert.False(t, mismatch.Approved)
					assert.Equal(t, RegionMismatchUntrustedHost, mismatch.Reason)
					require.ErrorAs(t, err, &responseError)
					assert.Equal(t, regionMismatchType, responseError.Type)
					require.NotNil(t, response, "the mismatch response is the caller's")
					if sender == "raw" {
						data, readErr := io.ReadAll(response.Body)
						require.NoError(t, readErr)
						assert.JSONEq(t, mismatchPayload, string(data), "its body is still readable")
					}
					assert.Zero(t, original.closes)
					require.NoError(t, response.Body.Close())
				}
				assert.Equal(t, 1, original.closes, "the mismatch body is closed once")
			})
		}
	}
}

// An earlier mismatch never describes a later failure: a wait the caller cancels, a body that cannot be opened again,
// and a transport failure or another status at the moved host keep their own errors.
func TestRegionMismatchIsNotAttachedToLaterFailures(t *testing.T) {
	transportErr := errors.New("connection reset at the moved host")
	bodyErr := errors.New("cannot open the body again")
	for _, failure := range []string{"canceled wait", "next body", "transport at the moved host", "status at the moved host"} {
		t.Run(failure, func(t *testing.T) {
			sites := &regionTestSites{handlers: map[string]func(*http.Request) (*http.Response, error){
				"us.files.com": regionTestMismatch("eu.files.com"),
				"eu.files.com": func(req *http.Request) (*http.Response, error) {
					if failure == "transport at the moved host" {
						return nil, transportErr
					}
					return regionTestAnswer(http.StatusServiceUnavailable, "application/json", `{"error":"unavailable","http-code":503,"type":"service-unavailable"}`, nil)(req)
				},
			}}
			config := regionTestConfig(sites)
			config.RetryMax = 1
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if failure == "canceled wait" {
				config.Client.Backoff = regionTestBackoff(config.Client.Backoff, func(int, *http.Response) bool {
					cancel()
					return true
				})
			}
			var response *http.Response
			var err error
			if failure == "next body" {
				opens := 0
				body := lib.Params{Params: map[string]interface{}{"paths": []string{}}}
				headers := &http.Header{}
				config.SetHeaders(headers)
				response, err = CallRawAPI(&CallParams{Method: http.MethodPost, Config: config, Uri: config.RootPath() + "/zip_downloads", Params: body, Headers: headers, RetryableBody: func() (io.Reader, error) {
					opens++
					if opens > 2 {
						return nil, bodyErr
					}
					return body.ToJSON()
				}})
			} else {
				_, response, err = Call(http.MethodGet, config, "/files/report.csv", nil, WithContext(ctx))
			}

			var mismatch *RegionMismatchError
			assert.False(t, errors.As(err, &mismatch), "the earlier mismatch does not describe %v", err)
			switch failure {
			case "canceled wait":
				require.ErrorIs(t, err, context.Canceled)
				assert.Nil(t, response)
				assert.Equal(t, []string{"us.files.com"}, sites.hosts())
			case "next body":
				require.ErrorIs(t, err, bodyErr)
				require.NotNil(t, response, "the dependency returns the previous, drained response")
				require.NoError(t, response.Body.Close())
				assert.Equal(t, []string{"us.files.com"}, sites.hosts())
			case "transport at the moved host":
				require.ErrorIs(t, err, transportErr)
				assert.Nil(t, response)
				assert.Equal(t, []string{"us.files.com", "eu.files.com"}, sites.hosts())
			case "status at the moved host":
				require.Error(t, err)
				require.NotNil(t, response)
				assert.Equal(t, http.StatusServiceUnavailable, response.StatusCode)
				require.NoError(t, response.Body.Close())
				assert.Equal(t, []string{"us.files.com", "eu.files.com"}, sites.hosts())
			}
		})
	}
}

// A moved request sends the Files.com authentication its answered request actually sent: what the caller's request
// hook set for the first attempt, never the environment's again after a redirect at the origin, and none when it sent
// none.
func TestRegionMismatchSendsTheAuthenticationTheIssuerReceived(t *testing.T) {
	t.Run("credentials a request hook set", func(t *testing.T) {
		sites := &regionTestSites{handlers: map[string]func(*http.Request) (*http.Response, error){
			"us.files.com": regionTestMismatch("eu.files.com"),
			"eu.files.com": regionTestOK,
		}}
		config := regionTestConfig(sites)
		config.Client.RequestLogHook = func(_ retryablehttp.Logger, req *http.Request, attempt int) {
			if attempt == 0 {
				req.Header.Set(apiKeyHeader, "hook-key")
				req.Header.Set(workspaceIdHeader, "77")
			}
		}

		_, response, err := Call(http.MethodGet, config, "/files/report.csv", nil)

		require.NoError(t, err)
		require.NoError(t, response.Body.Close())
		require.Equal(t, []string{"us.files.com", "eu.files.com"}, sites.hosts())
		for _, sent := range sites.sent {
			assert.Equal(t, "hook-key", sent.key)
			assert.Equal(t, "77", sent.workspace)
		}
	})
	t.Run("a redirect at the origin after the environment changed", func(t *testing.T) {
		t.Setenv("FILES_API_KEY", "first-key")
		sites := &regionTestSites{handlers: map[string]func(*http.Request) (*http.Response, error){
			"us.files.com": func(req *http.Request) (*http.Response, error) {
				if req.URL.Path == "/api/rest/v1/files/report.csv" {
					t.Setenv("FILES_API_KEY", "changed-key")
					return regionTestAnswer(http.StatusFound, "", "", http.Header{"Location": []string{"/api/rest/v1/files/moved.csv"}})(req)
				}
				return regionTestMismatch("eu.files.com")(req)
			},
			"eu.files.com": regionTestOK,
		}}
		config := Config{EndpointOverride: "https://us.files.com"}.Init().SetCustomClient(&http.Client{Transport: sites})

		_, response, err := Call(http.MethodGet, config, "/files/report.csv", nil)

		require.NoError(t, err)
		require.NoError(t, response.Body.Close())
		require.Equal(t, []string{"us.files.com", "us.files.com", "eu.files.com"}, sites.hosts())
		for _, sent := range sites.sent {
			assert.Equal(t, "first-key", sent.key, "the operation never takes the environment's new key")
		}
		assert.Equal(t, "/api/rest/v1/files/report.csv", sites.sent[2].path, "the original request is sent again")
	})
	t.Run("no authentication", func(t *testing.T) {
		t.Setenv("FILES_API_KEY", "")
		sites := &regionTestSites{handlers: map[string]func(*http.Request) (*http.Response, error){
			"us.files.com": func(req *http.Request) (*http.Response, error) {
				t.Setenv("FILES_API_KEY", "late-key")
				return regionTestMismatch("eu.files.com")(req)
			},
			"eu.files.com": regionTestOK,
		}}
		config := Config{EndpointOverride: "https://us.files.com"}.Init().SetCustomClient(&http.Client{Transport: sites})

		_, response, err := Call(http.MethodGet, config, "/files/report.csv", nil)

		require.NoError(t, err)
		require.NoError(t, response.Body.Close())
		require.Equal(t, []string{"us.files.com", "eu.files.com"}, sites.hosts())
		for _, sent := range sites.sent {
			assert.Empty(t, sent.key)
			assert.Empty(t, sent.session)
		}
	})
}

// After a move, a redirect or a retry that the caller's hooks, CheckRedirect, PrepareRetry or RequestLogHook, send to
// another origin carries none of the four Files.com authentication headers, and its answer cannot move the request
// again.
func TestRegionMismatchGuardsTheFinalDestinationAfterCallerHooks(t *testing.T) {
	for _, hook := range []string{"CheckRedirect", "PrepareRetry", "RequestLogHook"} {
		t.Run(hook, func(t *testing.T) {
			sites := &regionTestSites{handlers: map[string]func(*http.Request) (*http.Response, error){
				"us.files.com":  regionTestMismatch("eu.files.com"),
				"eu.files.com":  regionTestAnswer(http.StatusFound, "", "", http.Header{"Location": []string{"/api/rest/v1/files/next.csv"}}),
				"d.example.com": regionTestMismatch("ap.files.com"),
				"ap.files.com":  regionTestOK,
			}}
			httpClient := &http.Client{Transport: sites}
			if hook == "CheckRedirect" {
				httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					// The caller sends an otherwise in-origin redirect elsewhere.
					req.URL.Host = "d.example.com"
					return nil
				}
			}
			config := Config{APIKey: "region-key", WorkspaceId: 9, EndpointOverride: "https://us.files.com", AdditionalHeaders: map[string]string{
				"X-FilesAPI-Auth":          "region-session",
				"X-Files-Reauthentication": "password:region",
			}}.Init().SetCustomClient(httpClient)
			if hook == "PrepareRetry" {
				config.Client.PrepareRetry = func(req *http.Request) error {
					relocated := *req.URL
					relocated.Host = "d.example.com"
					req.URL = &relocated
					return nil
				}
			}
			if hook == "RequestLogHook" {
				config.Client.RequestLogHook = func(_ retryablehttp.Logger, req *http.Request, attempt int) {
					// The caller sends the moved retry elsewhere, changing its URL in place.
					if attempt > 0 {
						req.URL.Host = "d.example.com"
					}
				}
			}

			_, response, err := Call(http.MethodGet, config, "/files/report.csv", nil)

			var mismatch *RegionMismatchError
			require.ErrorAs(t, err, &mismatch)
			assert.Equal(t, RegionMismatchCrossedOrigin, mismatch.Reason, "an answer from another origin moves nothing")
			require.NotNil(t, response)
			require.NoError(t, response.Body.Close())
			want := []string{"us.files.com", "d.example.com"}
			if hook == "CheckRedirect" {
				want = []string{"us.files.com", "eu.files.com", "d.example.com"}
			}
			require.Equal(t, want, sites.hosts())
			for _, sent := range sites.sent[:len(sites.sent)-1] {
				assert.Equal(t, "region-key", sent.key)
				assert.Equal(t, "region-session", sent.session)
				assert.Equal(t, "password:region", sent.reauth)
				assert.Equal(t, "9", sent.workspace)
			}
			foreign := sites.sent[len(sites.sent)-1]
			assert.Empty(t, foreign.key)
			assert.Empty(t, foreign.session)
			assert.Empty(t, foreign.reauth)
			assert.Empty(t, foreign.workspace)
		})
	}
}

// regionRewritingSites is regionTestSites behind a transport that, as MockAPIServer's does, rewrites the request it is
// given in place: once the site the SDK sent it to has answered, its URL leads to physical.test.
type regionRewritingSites struct{ *regionTestSites }

func (s regionRewritingSites) RoundTrip(req *http.Request) (*http.Response, error) {
	response, err := s.regionTestSites.RoundTrip(req)
	req.URL.Scheme, req.URL.Host = "http", "physical.test"
	return response, err
}

// An ordinary retry through a transport that rewrites the request it is given in place goes where the SDK sent the
// attempt, with the same method, path, query, body and authentication, and a mismatch answering it moves the request
// within the same retry budget.
func TestRegionMismatchOrdinaryRetriesKeepTheirDestinationBelowTheTransport(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodPost} {
		t.Run(method, func(t *testing.T) {
			usAnswers := 0
			sites := &regionTestSites{handlers: map[string]func(*http.Request) (*http.Response, error){
				"us.files.com": func(req *http.Request) (*http.Response, error) {
					usAnswers++
					if usAnswers == 1 {
						return regionTestAnswer(http.StatusServiceUnavailable, "application/json", `{"error":"unavailable","http-code":503,"type":"service-unavailable"}`, nil)(req)
					}
					return regionTestMismatch("eu.files.com")(req)
				},
				"eu.files.com": regionTestOK,
			}}
			config := Config{APIKey: "region-key", WorkspaceId: 9, EndpointOverride: "https://us.files.com"}.Init().SetCustomClient(&http.Client{Transport: regionRewritingSites{sites}})
			var params lib.Values
			if method == http.MethodPost {
				params = lib.Params{Params: map[string]interface{}{"path": "Reports"}}
			}

			_, response, err := Call(method, config, "/files/report%20a.csv?cursor=c%201", params)

			require.NoError(t, err)
			require.NoError(t, response.Body.Close())
			require.Equal(t, []string{"us.files.com", "us.files.com", "eu.files.com"}, sites.hosts(), "the ordinary retry and the move share the default retries")
			for _, sent := range sites.sent {
				assert.Equal(t, method, sent.method)
				assert.Equal(t, "/api/rest/v1/files/report%20a.csv", sent.path)
				assert.Equal(t, "cursor=c%201", sent.query)
				assert.Equal(t, sites.sent[0].body, sent.body)
				assert.Equal(t, "region-key", sent.key)
				assert.Equal(t, "9", sent.workspace)
			}
		})
	}
}

// A redirect of an API request gets the configured headers other than Files.com authentication again before the
// caller's CheckRedirect, where its sender set them before it was enrolled: Call's always, a raw request's when it
// carried Files.com authentication. The caller's own changes in CheckRedirect are sent, and the four Files.com
// authentication headers stay the operation's own and do not leave its origin.
func TestRegionMismatchRedirectsKeepTheConfiguredHeadersOtherThanFilesAuthentication(t *testing.T) {
	for _, sender := range []string{"standard", "raw with Files.com authentication", "raw without Files.com authentication"} {
		t.Run(sender, func(t *testing.T) {
			t.Setenv("FILES_API_KEY", "")
			authenticated := sender != "raw without Files.com authentication"
			var atCallback, atStorage http.Header
			storage := func(req *http.Request) (*http.Response, error) {
				atStorage = req.Header.Clone()
				return regionTestOK(req)
			}
			sites := &regionTestSites{handlers: map[string]func(*http.Request) (*http.Response, error){
				"us.files.com":        regionTestAnswer(http.StatusFound, "", "", http.Header{"Location": []string{"https://storage.example.com/object"}}),
				"storage.example.com": storage,
			}}
			config := Config{UserAgent: "configured-agent", Language: "fr", EndpointOverride: "https://us.files.com", AdditionalHeaders: map[string]string{
				"X-Custom":      "configured",
				"Authorization": "Bearer configured",
				"Cookie":        "session=configured",
			}}
			if authenticated {
				config.APIKey = "region-key"
				config.AdditionalHeaders["X-FilesAPI-Auth"] = "configured-session"
			}
			config = config.Init().SetCustomClient(&http.Client{Transport: sites, CheckRedirect: func(req *http.Request, via []*http.Request) error {
				atCallback = req.Header.Clone()
				// The caller signs the redirect and sets one of the configured headers itself.
				req.Header.Set("X-Signature", "signed")
				req.Header.Set("X-Custom", "callback")
				return nil
			}})
			overrides := http.Header{"User-Agent": []string{"request-agent"}, "X-Custom": []string{"request"}}
			var response *http.Response
			var err error
			if sender == "standard" {
				_, response, err = Call(http.MethodGet, config, "/files/report.csv", nil, RequestHeadersOption(&overrides))
			} else {
				headers := &http.Header{}
				config.SetHeaders(headers)
				for name, values := range overrides {
					headers.Set(name, values[0])
				}
				response, err = CallRawAPI(&CallParams{Method: http.MethodGet, Config: config, Uri: config.RootPath() + "/files/report.csv", Headers: headers})
			}

			require.NoError(t, err)
			require.NoError(t, response.Body.Close())
			require.Equal(t, []string{"us.files.com", "storage.example.com"}, sites.hosts())
			wantAgent, wantCustom := "request-agent", "request"
			if authenticated {
				assert.Equal(t, "region-key", sites.sent[0].key)
				assert.Equal(t, "configured-session", sites.sent[0].session)
				// The configured headers are set again before the caller's CheckRedirect, as they were before.
				wantAgent, wantCustom = "configured-agent", "configured"
				assert.Equal(t, "Bearer configured", atCallback.Get("Authorization"))
				assert.Equal(t, "session=configured", atCallback.Get("Cookie"))
				assert.Equal(t, "Bearer configured", atStorage.Get("Authorization"))
				assert.Equal(t, "session=configured", atStorage.Get("Cookie"))
			} else {
				// A raw request without Files.com authentication is redirected as net/http redirects it, as before.
				assert.Empty(t, atCallback.Get("Authorization"))
				assert.Empty(t, atStorage.Get("Authorization"))
				assert.Empty(t, atStorage.Get("Cookie"))
			}
			assert.Equal(t, wantAgent, atCallback.Get("User-Agent"))
			assert.Equal(t, wantCustom, atCallback.Get("X-Custom"))
			assert.Equal(t, "fr", atCallback.Get("Accept-Language"))
			assert.Equal(t, wantAgent, atStorage.Get("User-Agent"))
			assert.Equal(t, "callback", atStorage.Get("X-Custom"), "the caller's own change is sent")
			assert.Equal(t, "signed", atStorage.Get("X-Signature"))
			for _, name := range []string{apiKeyHeader, sessionIdHeader, reauthenticationHeader, workspaceIdHeader} {
				assert.Empty(t, atCallback.Get(name), "the redirect leaves the origin without %s", name)
				assert.Empty(t, atStorage.Get(name))
			}
		})
	}
}

// regionTailBody is a response body whose tail never arrives, as an unfinished body's does not: once its head has been
// read, every further read is counted and fails.
type regionTailBody struct {
	head          io.Reader
	firstErr      error
	firstErrReads int
	tailReads     int
	closes        int
}

var errRegionTailRead = errors.New("the unfinished tail of the body was read")

// Read waits, as a reader may, until it can fill p: a read the head can fill is answered from it, and one asking for
// more than is left of the head waits for the tail, which is counted and fails after the rest of the head. A body
// given firstErr instead answers its first read asking for any byte with one byte of the head and that error, once, as
// a reader may return an error beside the data it read.
func (b *regionTailBody) Read(p []byte) (int, error) {
	if b.firstErr != nil && len(p) > 0 {
		err := b.firstErr
		b.firstErr = nil
		b.firstErrReads++
		n, _ := b.head.Read(p[:1])
		return n, err
	}
	n, _ := io.ReadFull(b.head, p)
	if n < len(p) {
		b.tailReads++
		return n, errRegionTailRead
	}
	return n, nil
}

func (b *regionTailBody) Close() error {
	b.closes++
	return nil
}

// Recognizing a mismatch never waits for more of an unrelated 403 than the SDK's parsers do: an nginx page with its
// request ID keeps its error without its body being waited for, even when its length is unknown and it names a host,
// which recognition reads only as far as the first byte that differs from the legacy prefix, and a plain-text 403
// naming a host is read only as far as its own error. What recognition read is put back, so an nginx page still reads
// as it arrived: its head and then its tail's failure, or, when its first byte came with an error, that byte and then
// that same error.
func TestRegionMismatchRecognitionLeavesUnrelatedBodiesUnread(t *testing.T) {
	errFirstRead := errors.New("the body failed with its first byte")
	for name, test := range map[string]struct {
		header        http.Header
		head          string
		unknownLength bool
		unread        bool
		firstErr      error
		wantErr       string
	}{
		"an nginx page": {
			header:  http.Header{"Content-Type": {"text/html"}, "Server": {"nginx"}, "X-Request-Id": {"probe"}},
			head:    "<html>denied</html>",
			unread:  true,
			wantErr: "files.com Server error - request id: probe",
		},
		"a plain-text 403 naming a host": {
			header:  http.Header{"Content-Type": {"text/plain"}, "X-Files-Host": {"eu.files.com"}},
			head:    "Forbidden: this site does not allow the request",
			wantErr: "Forbidden: this site does not allow the request",
		},
		"an nginx page of unknown length naming a host": {
			header:        http.Header{"Content-Type": {"text/html"}, "Server": {"nginx"}, "X-Request-Id": {"probe"}, "X-Files-Host": {"eu.files.com"}},
			head:          "<html>denied</html>",
			unknownLength: true,
			unread:        true,
			wantErr:       "files.com Server error - request id: probe",
		},
		"an nginx page naming a host whose first byte comes with an error": {
			header:        http.Header{"Content-Type": {"text/html"}, "Server": {"nginx"}, "X-Request-Id": {"probe"}, "X-Files-Host": {"eu.files.com"}},
			head:          "<",
			unknownLength: true,
			unread:        true,
			firstErr:      errFirstRead,
			wantErr:       "files.com Server error - request id: probe",
		},
	} {
		t.Run(name, func(t *testing.T) {
			body := &regionTailBody{head: strings.NewReader(test.head), firstErr: test.firstErr}
			length := int64(len(test.head))
			if test.unknownLength {
				length = -1
			}
			sites := &regionTestSites{handlers: map[string]func(*http.Request) (*http.Response, error){
				"us.files.com": func(req *http.Request) (*http.Response, error) {
					return &http.Response{StatusCode: http.StatusForbidden, Status: "403 Forbidden", Header: test.header, Body: body, ContentLength: length, Request: req}, nil
				},
			}}

			_, response, err := Call(http.MethodGet, regionTestConfig(sites), "/files/report.csv", nil)

			require.Error(t, err)
			assert.Contains(t, err.Error(), test.wantErr)
			var mismatch *RegionMismatchError
			assert.False(t, errors.As(err, &mismatch))
			assert.Zero(t, body.tailReads, "the unfinished tail is never waited for")
			assert.Equal(t, []string{"us.files.com"}, sites.hosts())
			require.NotNil(t, response)
			if test.firstErr != nil {
				assert.Equal(t, 1, body.firstErrReads, "recognition read the byte that came with the error, before the caller")
			}
			if test.unread {
				// The body reads as it arrived: what recognition read, then the rest as the original gave it.
				data, readErr := io.ReadAll(response.Body)
				assert.Equal(t, test.head, string(data))
				if test.firstErr != nil {
					// Recognition read the head's only byte with the error that came with it: the body gives back both,
					// the same error, and reads nothing more of the original.
					assert.Same(t, test.firstErr, readErr)
					assert.Zero(t, body.tailReads)
					assert.Equal(t, 1, body.firstErrReads)
				} else {
					// The rest of the head with the tail's failure, as the original's last read returned them together.
					assert.ErrorIs(t, readErr, errRegionTailRead)
					assert.Equal(t, 1, body.tailReads)
				}
			}
			require.NoError(t, response.Body.Close())
			assert.Equal(t, 1, body.closes)
		})
	}
}

// CallParams sent again describe only their own call: when the next request cannot even be built, it fails with its
// own error, not an earlier call's mismatch, sends nothing, and has no relocated answer, whichever sender sends it.
func TestCallRawAPIParamsSentAgainDescribeOnlyTheirOwnCall(t *testing.T) {
	for _, earlier := range []string{"a refused mismatch", "a move"} {
		for _, next := range []string{"CallRawAPI", "CallRaw"} {
			t.Run(earlier+"/"+next, func(t *testing.T) {
				sites := &regionTestSites{handlers: map[string]func(*http.Request) (*http.Response, error){
					"us.files.com": regionTestMismatch("eu.files.com"),
					"eu.files.com": regionTestOK,
				}}
				config := regionTestConfig(sites)
				if earlier == "a refused mismatch" {
					config.RetryMax = 0
				}
				headers := &http.Header{}
				config.SetHeaders(headers)
				params := &CallParams{Method: http.MethodGet, Config: config, Uri: config.RootPath() + "/files/report.csv", Headers: headers}
				response, err := CallRawAPI(params)
				var mismatch *RegionMismatchError
				if earlier == "a refused mismatch" {
					require.ErrorAs(t, err, &mismatch)
				} else {
					require.NoError(t, err)
					require.True(t, params.RelocatedAPIContext().Relocated())
				}
				require.NoError(t, response.Body.Close())
				sent := len(sites.sent)

				params.Uri = config.RootPath() + "/files/%zz"
				if next == "CallRawAPI" {
					_, err = CallRawAPI(params)
				} else {
					_, err = CallRaw(params)
				}

				var urlErr *url.Error
				require.ErrorAs(t, err, &urlErr, "the request's own build error")
				assert.Equal(t, "parse", urlErr.Op)
				assert.Equal(t, "[redacted]", urlErr.URL)
				assert.False(t, errors.As(err, &mismatch), "an earlier call's mismatch does not describe it")
				var responseError ResponseError
				assert.False(t, errors.As(err, &responseError))
				assert.Len(t, sites.sent, sent, "nothing is sent")
				assert.False(t, params.RelocatedAPIContext().Relocated(), "nor is an earlier answer its own")
			})
		}
	}
}

// A raw API request's relocated answer is the answer the call returned: none when the moved request failed without
// one or the ErrorHandler selected an unrelated response, and the moved request's own when the handler replaced it.
func TestCallRawAPIRelocatedContextIsTheAnswerTheCallReturned(t *testing.T) {
	unavailable := regionTestAnswer(http.StatusServiceUnavailable, "application/json", `{"error":"unavailable","http-code":503,"type":"service-unavailable"}`, nil)
	for name, test := range map[string]struct {
		moved     func(*http.Request) (*http.Response, error)
		selected  func(answer *http.Response) *http.Response
		relocated bool
	}{
		"a transport failure at the moved host": {
			moved: func(*http.Request) (*http.Response, error) {
				return nil, errors.New("connection reset at the moved host")
			},
		},
		"an unrelated response the handler selected": {
			moved: unavailable,
			selected: func(*http.Response) *http.Response {
				return &http.Response{StatusCode: http.StatusOK, Status: "200 OK", Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(`{}`))}
			},
		},
		"the handler's replacement for the moved answer": {
			moved: unavailable,
			selected: func(answer *http.Response) *http.Response {
				replacement, _ := regionTestOK(answer.Request)
				return replacement
			},
			relocated: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			sites := &regionTestSites{handlers: map[string]func(*http.Request) (*http.Response, error){
				"us.files.com": regionTestMismatch("eu.files.com"),
				"eu.files.com": test.moved,
			}}
			config := regionTestConfig(sites)
			// The move takes the one retry, so the moved request's outcome is the call's.
			config.RetryMax = 1
			if test.selected != nil {
				config.Client.ErrorHandler = func(answer *http.Response, _ error, _ int) (*http.Response, error) {
					if err := answer.Body.Close(); err != nil {
						return nil, err
					}
					return test.selected(answer), nil
				}
			}
			headers := &http.Header{}
			config.SetHeaders(headers)
			params := &CallParams{Method: http.MethodGet, Config: config, Uri: config.RootPath() + "/files/report.csv", Headers: headers}

			response, err := CallRawAPI(params)

			assert.Equal(t, []string{"us.files.com", "eu.files.com"}, sites.hosts())
			if test.selected == nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NoError(t, response.Body.Close())
			}
			assert.Equal(t, test.relocated, params.RelocatedAPIContext().Relocated())
			if test.relocated {
				assert.Equal(t, "https://eu.files.com", params.RelocatedAPIContext().Origin())
			}
		})
	}
}

// Only a recognized mismatch from the operation's origin can move a request, and only to a trusted HTTPS origin; any
// other answer, the opt-out or a refusal leaves the request where it was, a recognized mismatch described.
func TestRegionMismatchRecognitionAndTrust(t *testing.T) {
	large := strings.Repeat("x", 100_000)
	trustCustom := func(origin, from string) bool {
		return from == "https://us.files.com" && (origin == "https://eu.example.com" || origin == "https://10.0.0.5:8443")
	}
	for name, test := range map[string]struct {
		answer       func(*http.Request) (*http.Response, error)
		modify       func(*Config)
		hosts        []string
		unrecognized bool
		reason       string
		origin       string
		approved     bool
	}{
		"typed":                                 {answer: regionTestMismatch("eu.files.com"), hosts: []string{"us.files.com", "eu.files.com"}},
		"typed with a large body":               {answer: regionTestAnswer(http.StatusForbidden, "application/json", fmt.Sprintf(`{"type":%q,"http-code":403,"error":%q,"data":{"host":"eu.files.com"}}`, regionMismatchType, large), nil), hosts: []string{"us.files.com", "eu.files.com"}},
		"typed naming its host in X-Files-Host": {answer: regionTestAnswer(http.StatusForbidden, "application/json", `{"type":"not-authenticated/lockout-region-mismatch","http-code":403,"error":"Lockout region mismatch"}`, http.Header{"X-Files-Host": []string{"eu.files.com"}}), hosts: []string{"us.files.com", "eu.files.com"}},
		"legacy without a Content-Type":         {answer: regionTestAnswer(http.StatusForbidden, "", "You have connected to a URL that does not serve this site.", http.Header{"X-Files-Host": []string{"eu.files.com"}}), hosts: []string{"us.files.com", "eu.files.com"}},
		"legacy in an HTML page":                {answer: regionTestAnswer(http.StatusForbidden, "text/html", "You have connected to a URL that does not serve this site.", http.Header{"X-Files-Host": []string{"eu.files.com"}}), hosts: []string{"us.files.com", "eu.files.com"}},
		"legacy without X-Files-Host":           {answer: regionTestAnswer(http.StatusForbidden, "text/plain", "You have connected to a URL that does not serve this site.", nil), hosts: []string{"us.files.com"}, unrecognized: true},
		"another 403 naming a host":             {answer: regionTestAnswer(http.StatusForbidden, "application/json", `{"type":"not-authenticated/invalid-credentials","http-code":403,"error":"Invalid credentials"}`, http.Header{"X-Files-Host": []string{"eu.files.com"}}), hosts: []string{"us.files.com"}, unrecognized: true},
		"an untrusted host":                     {answer: regionTestMismatch("eu.example.com"), hosts: []string{"us.files.com"}, reason: RegionMismatchUntrustedHost, origin: "https://eu.example.com"},
		"a custom host the caller trusts":       {answer: regionTestMismatch("eu.example.com"), modify: func(c *Config) { c.TrustRegionOrigin = trustCustom }, hosts: []string{"us.files.com", "eu.example.com"}},
		"an address and port the caller trusts": {answer: regionTestMismatch("https://10.0.0.5:8443"), modify: func(c *Config) { c.TrustRegionOrigin = trustCustom }, hosts: []string{"us.files.com", "10.0.0.5:8443"}},
		"an HTTP origin": {answer: regionTestMismatch("http://eu.files.com"), modify: func(c *Config) {
			c.TrustRegionOrigin = func(string, string) bool { panic("an HTTP origin is never offered for trust") }
		}, hosts: []string{"us.files.com"}, reason: RegionMismatchUntrustedHost, origin: "http://eu.files.com"},
		"no host":         {answer: regionTestAnswer(http.StatusForbidden, "application/json", `{"type":"not-authenticated/lockout-region-mismatch","http-code":403,"error":"Lockout region mismatch"}`, nil), hosts: []string{"us.files.com"}, reason: RegionMismatchNoHost},
		"an invalid host": {answer: regionTestMismatch("eu.files.com/elsewhere"), hosts: []string{"us.files.com"}, reason: RegionMismatchInvalidHost},
		"an invalid data.host beside a valid header": {answer: regionTestAnswer(http.StatusForbidden, "application/json", regionTestPayload("eu.files.com/elsewhere"), http.Header{"X-Files-Host": []string{"eu.files.com"}}), hosts: []string{"us.files.com"}, reason: RegionMismatchInvalidHost},
		"the same origin": {answer: regionTestMismatch("US.files.com"), hosts: []string{"us.files.com"}, reason: RegionMismatchSameOrigin, origin: "https://us.files.com", approved: true},
		"the opt-out":     {answer: regionTestMismatch("eu.files.com"), modify: func(c *Config) { c.DisableRegionRelocation = true }, hosts: []string{"us.files.com"}, reason: RegionMismatchDisabled, origin: "https://eu.files.com", approved: true},
	} {
		t.Run(name, func(t *testing.T) {
			sites := &regionTestSites{handlers: map[string]func(*http.Request) (*http.Response, error){
				"us.files.com":   test.answer,
				"eu.files.com":   regionTestOK,
				"eu.example.com": regionTestOK,
				"10.0.0.5:8443":  regionTestOK,
			}}
			config := regionTestConfig(sites)
			if test.modify != nil {
				test.modify(&config)
			}

			_, response, err := Call(http.MethodGet, config, "/files/report.csv", nil)

			if response != nil && response.Body != nil {
				require.NoError(t, response.Body.Close())
			}
			assert.Equal(t, test.hosts, sites.hosts())
			var mismatch *RegionMismatchError
			switch {
			case len(test.hosts) == 2:
				require.NoError(t, err)
			case test.unrecognized:
				require.Error(t, err)
				assert.False(t, errors.As(err, &mismatch), "an answer that is not a mismatch is left as it was")
			default:
				require.ErrorAs(t, err, &mismatch)
				assert.Equal(t, test.reason, mismatch.Reason)
				assert.Equal(t, "https://us.files.com", mismatch.From)
				if test.origin != "" {
					assert.Equal(t, test.origin, mismatch.Origin)
				}
				assert.Equal(t, test.approved, mismatch.Approved)
				var responseError ResponseError
				require.ErrorAs(t, err, &responseError, "the API's error stays reachable")
				assert.Equal(t, regionMismatchType, responseError.Type)
			}
		})
	}
}

// Successive valid moves and cycles share the client's one retry budget, and a raw API request whose body cannot be
// sent again, or a CallRaw request, is never moved.
func TestRegionMismatchMovesShareTheRetryBudget(t *testing.T) {
	t.Run("successive moves", func(t *testing.T) {
		sites := &regionTestSites{handlers: map[string]func(*http.Request) (*http.Response, error){
			"us.files.com": regionTestMismatch("eu.files.com"),
			"eu.files.com": regionTestMismatch("ap.files.com"),
			"ap.files.com": regionTestOK,
		}}

		_, response, err := Call(http.MethodGet, regionTestConfig(sites), "/files/report.csv", nil)

		require.NoError(t, err)
		require.NoError(t, response.Body.Close())
		assert.Equal(t, []string{"us.files.com", "eu.files.com", "ap.files.com"}, sites.hosts())
		for _, sent := range sites.sent {
			assert.Equal(t, "region-key", sent.key)
		}
	})
	t.Run("a cycle", func(t *testing.T) {
		sites := &regionTestSites{handlers: map[string]func(*http.Request) (*http.Response, error){
			"us.files.com": regionTestMismatch("eu.files.com"),
			"eu.files.com": regionTestMismatch("us.files.com"),
		}}

		_, response, err := Call(http.MethodGet, regionTestConfig(sites), "/files/report.csv", nil)

		var mismatch *RegionMismatchError
		require.ErrorAs(t, err, &mismatch)
		assert.Equal(t, RegionMismatchRetryBudget, mismatch.Reason)
		assert.Equal(t, "https://eu.files.com", mismatch.From)
		assert.Equal(t, "https://us.files.com", mismatch.Origin)
		require.NotNil(t, response)
		require.NoError(t, response.Body.Close())
		assert.Equal(t, []string{"us.files.com", "eu.files.com", "us.files.com", "eu.files.com"}, sites.hosts(), "the default three retries")
	})
	t.Run("a raw API body without a factory", func(t *testing.T) {
		sites := &regionTestSites{handlers: map[string]func(*http.Request) (*http.Response, error){
			"us.files.com": regionTestMismatch("eu.files.com"),
			"eu.files.com": regionTestOK,
		}}
		config := regionTestConfig(sites)
		headers := &http.Header{}
		config.SetHeaders(headers)
		params := &CallParams{Method: http.MethodPost, Config: config, Uri: config.RootPath() + "/zip_downloads", Params: lib.Params{Params: map[string]interface{}{"paths": []string{}}}, Headers: headers}

		response, err := CallRawAPI(params)

		var mismatch *RegionMismatchError
		require.ErrorAs(t, err, &mismatch)
		assert.Equal(t, RegionMismatchBodyNotReplayable, mismatch.Reason)
		require.NotNil(t, response)
		require.NoError(t, response.Body.Close())
		assert.Equal(t, []string{"us.files.com"}, sites.hosts())
		assert.False(t, params.RelocatedAPIContext().Relocated())
	})
	t.Run("CallRaw", func(t *testing.T) {
		sites := &regionTestSites{handlers: map[string]func(*http.Request) (*http.Response, error){
			"us.files.com": regionTestMismatch("eu.files.com"),
			"eu.files.com": regionTestOK,
		}}
		config := regionTestConfig(sites)
		headers := &http.Header{}
		config.SetHeaders(headers)

		response, err := CallRaw(&CallParams{Method: http.MethodGet, Config: config, Uri: config.RootPath() + "/files/report.csv", Headers: headers})

		require.NoError(t, err, "CallRaw leaves the response to its caller as before")
		require.NotNil(t, response)
		assert.Equal(t, http.StatusForbidden, response.StatusCode)
		require.NoError(t, response.Body.Close())
		assert.Equal(t, []string{"us.files.com"}, sites.hosts())
	})
}
