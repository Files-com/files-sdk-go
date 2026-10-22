package files_sdk

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/Files-com/files-sdk-go/v3/lib"
	libLog "github.com/Files-com/files-sdk-go/v3/lib/logpath"
	"github.com/hashicorp/go-retryablehttp"
)

var VERSION = "3.3.321"
var defaultUserAgent = fmt.Sprintf("%v %v", UserAgent, strings.TrimSpace(VERSION))

const (
	UserAgent   = "Files.com Go SDK"
	DefaultSite = "app"
	APIPath     = "/api/rest/v1"
)

const (
	apiKeyHeader           = "X-FilesAPI-Key"
	sessionIdHeader        = "X-FilesAPI-Auth"
	reauthenticationHeader = "X-Files-Reauthentication"
	workspaceIdHeader      = "X-Files-Workspace-Id"
)

var GlobalConfig Config

func init() {
	GlobalConfig = Config{}.Init()
}

// Config controls Files.com API authentication, workspace scoping, and transport behavior.
type Config struct {
	APIKey    string `header:"X-FilesAPI-Key" json:"api_key"`
	SessionId string `header:"X-FilesAPI-Auth" json:"session_id"`
	// WorkspaceId scopes Files.com API requests when non-zero.
	WorkspaceId      int64  `header:"X-Files-Workspace-Id" json:"workspace_id,omitempty"`
	Language         string `header:"Accept-Language"`
	Subdomain        string `json:"subdomain"`
	EndpointOverride string `json:"endpoint_override"`
	*retryablehttp.Client
	AdditionalHeaders map[string]string `json:"additional_headers"`
	// Logger receives SDK diagnostics. Init uses a logger that discards output when this is nil.
	lib.Logger
	// Debug enables detailed request and response diagnostics, which may contain credentials
	// and signed URLs. FILES_SDK_DEBUG also enables these diagnostics. SDK retry log entries
	// at INFO, WARN, and ERROR still redact request URLs and credentials.
	Debug        bool   `json:"debug"`
	UserAgent    string `json:"user_agents"`
	Environment  `json:"environment"`
	FeatureFlags map[string]bool `json:"feature_flags"`
	// DisableDirectTransfers forces upload and download requests to use proxied transfer paths.
	DisableDirectTransfers bool `json:"disable_direct_transfers"`
}

// Init returns a copy of c with defaults for its logger, retry client, feature flags,
// and user agent. It preserves values already supplied by the caller.
func (c Config) Init() Config {
	if c.Logger == nil {
		c.Logger = lib.NullLogger{}
	}
	if c.Client == nil {
		c.Client = lib.DefaultRetryableHttp(c)
	}

	if c.FeatureFlags == nil {
		c.FeatureFlags = FeatureFlags()
	}

	if c.UserAgent == "" {
		c.UserAgent = defaultUserAgent
	}

	return c
}

func (c Config) Endpoint() string {
	if c.EndpointOverride != "" && !strings.HasPrefix(c.EndpointOverride, "https://") && !strings.HasPrefix(c.EndpointOverride, "http://") {
		c.EndpointOverride = "https://" + c.EndpointOverride
	}

	return lib.DefaultString(
		c.EndpointOverride,
		strings.Replace(c.Environment.Endpoint(), "{{SUBDOMAIN}}", lib.DefaultString(c.Subdomain, DefaultSite), 1),
	)
}

// Do sends req through the configured retry client with redirect handling.
// The caller must close the returned response body. Outside debug mode, returned
// URL errors redact the request URL while preserving the operation and cause.
// Retry log entries above DEBUG redact request URLs and credentials in either mode.
func (c Config) Do(req *http.Request) (*http.Response, error) {
	client := c.redirectSafeClient()
	client.Logger = c.retryLogFor(client.Logger, req)
	guardRetryableResponses(client, true)
	response, err := (&http.Client{Transport: &retryableRoundTripper{client: client}}).Do(req)
	if c.InDebug() {
		return response, err
	}
	// Outside DEBUG, keep request URLs, which can be signed download URLs, out of errors, as CallRaw does. The error
	// keeps its type, operation and cause, so cancellation and timeout checks still work.
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		err = &url.Error{Op: urlErr.Op, URL: "[redacted]", Err: urlErr.Err}
	}
	return response, err
}

type retryableRoundTripper struct{ client *retryablehttp.Client }

func (r *retryableRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	retryable, err := retryablehttp.FromRequest(req)
	if err != nil {
		return nil, err
	}
	response, err := r.client.Do(retryable)
	return finishRetryableRoundTrip(response, err)
}

func finishRetryableRoundTrip(response *http.Response, err error) (*http.Response, error) {
	if err != nil {
		// The guarded error handler has released the supplied final body.
		// net/http cannot consume a response returned with a RoundTrip error.
		response = nil
	}
	// Match retryablehttp.RoundTripper: net/http adds its own URL error.
	if _, ok := err.(*url.Error); ok {
		return response, errors.Unwrap(err)
	}
	return response, err
}

func guardRetryableResponses(client *retryablehttp.Client, closeOnError bool) {
	checkRetry := client.CheckRetry
	if checkRetry != nil {
		client.CheckRetry = func(ctx context.Context, response *http.Response, doErr error) (bool, error) {
			retry, policyErr := checkRetry(ctx, response, doErr)
			if retry || policyErr != nil || doErr != nil {
				// Install before the dependency drains a retry response. It may
				// return that response if PrepareRetry or the next body fails.
				_, redirectError := doErr.(*url.Error)
				guardRetryableBody(response, redirectError)
			}
			return retry, policyErr
		}
	}
	if !closeOnError {
		return
	}
	handler := client.ErrorHandler
	if handler == nil {
		return // retryablehttp's default handler already closes and discards.
	}
	client.ErrorHandler = func(response *http.Response, err error, attempts int) (*http.Response, error) {
		body := guardRetryableBody(response, false)
		processed, err := handler(response, err, attempts)
		if err != nil && body != nil {
			_ = body.Close()
		}
		// A handler replacing the response or its body owns that replacement's
		// cleanup contract. Only closes of the supplied body can be observed here.
		return processed, err
	}
}

func guardRetryableBody(response *http.Response, closed bool) *retryableErrorBody {
	if response == nil || response.Body == nil {
		return nil
	}
	if body, ok := response.Body.(*retryableErrorBody); ok {
		return body
	}
	body := &retryableErrorBody{ReadCloser: response.Body, closed: closed}
	response.Body = body
	return body
}

// This guard exists only on retry and error responses. The dependency or its
// ErrorHandler can close the body or hand it back to its caller.
type retryableErrorBody struct {
	io.ReadCloser
	closed bool
}

func (b *retryableErrorBody) Close() error {
	if b.closed {
		return nil
	}
	b.closed = true
	return b.ReadCloser.Close()
}

func (b *retryableErrorBody) WriteTo(w io.Writer) (int64, error) {
	return io.Copy(w, b.ReadCloser)
}

// retryLogFor returns the logger for one request's retry client in place of logger, the client's own. Entries at
// DEBUG level, such as the request line with its full URL, pass through unchanged. Entries above DEBUG, such as
// "request failed" or a log hook's own, keep their event and context with the request's URL and credentials
// removed. A logger without levels gets DEBUG entries only in debug mode, its only level control. The client's own
// logger is not changed.
func (c Config) retryLogFor(logger interface{}, req *http.Request) interface{} {
	redact := c.secretsOf(req)
	if sdkConfig, ok := logger.(Config); ok {
		// SDK retry clients log through the Config they were built from, whose Logger may have levels.
		return printfRetryLog{logger: sdkConfig.Logger, debug: c.InDebug(), redact: redact}
	}
	switch v := logger.(type) {
	case retryablehttp.LeveledLogger:
		return leveledRetryLog{logger: v, redact: redact}
	case retryablehttp.Logger:
		return printfRetryLog{logger: v, debug: c.InDebug(), redact: redact}
	}
	return logger
}

// printfRetryLog takes go-retryablehttp's Printf entries, whose level is the "[DEBUG]" or "[ERR]" prefix of their
// format. Other entries come from a log hook, which go-retryablehttp logs at Info for a logger with levels.
type printfRetryLog struct {
	logger lib.Logger
	debug  bool
	redact requestSecrets
}

func (l printfRetryLog) Printf(format string, args ...interface{}) {
	leveled, hasLevels := l.logger.(lib.LeveledLogger)
	if strings.HasPrefix(format, "[DEBUG] ") {
		if hasLevels {
			leveled.Debug(strings.TrimPrefix(format, "[DEBUG] "), args...)
		} else if l.debug {
			l.logger.Printf(format, args...)
		}
		return
	}
	entry := l.redact.text(fmt.Sprintf(format, l.redact.values(args)...))
	switch {
	case !hasLevels:
		l.logger.Printf("%s", entry)
	case strings.HasPrefix(entry, "[ERR] "):
		leveled.Error("%s", strings.TrimPrefix(entry, "[ERR] "))
	default:
		leveled.Info("%s", entry)
	}
}

// leveledRetryLog takes go-retryablehttp's structured entries, a message with key/value pairs, for a client logger
// with levels.
type leveledRetryLog struct {
	logger retryablehttp.LeveledLogger
	redact requestSecrets
}

func (l leveledRetryLog) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.Debug(msg, keysAndValues...)
}

func (l leveledRetryLog) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Info(l.redact.text(msg), l.redact.values(keysAndValues)...)
}

func (l leveledRetryLog) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.Warn(l.redact.text(msg), l.redact.values(keysAndValues)...)
}

func (l leveledRetryLog) Error(msg string, keysAndValues ...interface{}) {
	l.logger.Error(l.redact.text(msg), l.redact.values(keysAndValues)...)
}

// requestSecrets are one request's values that must not appear in log entries above DEBUG: its full URL and query,
// which can be signed, and the Files.com API key, session ID and reauthentication it carries or is configured with.
type requestSecrets []string

func (c Config) secretsOf(req *http.Request) requestSecrets {
	values := []string{c.GetAPIKey(), c.SessionId}
	if req != nil {
		values = append([]string{req.Header.Get(apiKeyHeader), req.Header.Get(sessionIdHeader), req.Header.Get(reauthenticationHeader)}, values...)
		// A request without a URL still reaches net/http, which returns its usual error for it.
		if req.URL != nil {
			// The full URL comes first, so it is replaced whole rather than around its query.
			values = append([]string{req.URL.String(), req.URL.RawQuery}, values...)
		}
	}
	var secrets requestSecrets
	for _, value := range values {
		if value != "" {
			secrets = append(secrets, value)
		}
	}
	return secrets
}

// text returns text with each secret replaced.
func (s requestSecrets) text(text string) string {
	for _, secret := range s {
		text = strings.ReplaceAll(text, secret, "[redacted]")
	}
	return text
}

// values returns a copy of args for logging: an absolute URL is replaced, an error becomes one with the same text
// with secrets replaced and, for a *url.Error, its URL redacted, and other strings have secrets replaced. args and
// the errors in it are not changed.
func (s requestSecrets) values(args []interface{}) []interface{} {
	clean := make([]interface{}, len(args))
	for i, arg := range args {
		clean[i] = arg
		switch v := arg.(type) {
		case error:
			var urlErr *url.Error
			if errors.As(v, &urlErr) {
				v = &url.Error{Op: urlErr.Op, URL: "[redacted]", Err: urlErr.Err}
			}
			clean[i] = errors.New(s.text(v.Error()))
		case string:
			if parsed, err := url.Parse(v); err == nil && parsed.Scheme != "" && parsed.Host != "" {
				clean[i] = "[redacted]"
			} else {
				clean[i] = s.text(v)
			}
		}
	}
	return clean
}

// replayableSender sends requests as Do does, but without holding their
// bodies in memory, where Do copies a body into memory to resend it. An
// unread replayableBody is read again from the start by each attempt,
// including each retry and redirect. Any other body, such as one a
// RequestOption set or wrapped, is read to its end first, as Do does, and
// copied as attachment content is kept; every retry resends that copy, and a
// redirect uses the request's GetBody, as with Do. Each attempt's body, and a
// redirect's replayableBody from GetBody, is opened with
// lib.PreparedJSON.Open, because the transport can read and close it after the
// call returns. Close removes the copies; any files such a body still reads
// stay open until it is closed.
type replayableSender struct {
	config Config
	client *retryablehttp.Client
	// logger is the client's own logger, which Do wraps for each request as
	// Config.Do does.
	logger interface{}
	mu     sync.Mutex
	copies []*lib.PreparedJSON
}

func (c Config) newReplayableSender() *replayableSender {
	client := c.redirectSafeClient()
	guardRetryableResponses(client, true)
	// http.Client opens a redirect's body with GetBody before it calls
	// CheckRedirect, which can stop the redirect or change its body. When
	// CheckRedirect stops the redirect, http.Client closes the request's
	// body, whatever CheckRedirect set it to, except after
	// ErrUseLastResponse; it never closes a body CheckRedirect replaced. An
	// open reader keeps the prepared files open, so the sender takes the
	// reader it opened for the redirect before CheckRedirect runs and sees
	// that it is closed once nothing can read it. Closing it again is safe,
	// and the sender closes no other body itself. redirectSafeClient returns
	// a new http.Client for each sender.
	checkRedirect := client.HTTPClient.CheckRedirect
	client.HTTPClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		opened, _ := req.Body.(*openedBody)
		err := checkRedirect(req, via)
		switch {
		case opened == nil || (err == nil && req.Body == io.ReadCloser(opened)):
			// The transport sends the sender's reader and closes it.
		case err != nil || req.Body == nil || req.Body == http.NoBody:
			// Nothing sends the sender's reader.
			_ = opened.Close()
		default:
			// The redirect sends CheckRedirect's body, which can read the
			// sender's reader, so that reader is closed with it.
			req.Body = &substituteBody{ReadCloser: req.Body, opened: opened}
		}
		return err
	}
	return &replayableSender{config: c, client: client, logger: client.Logger}
}

// substituteBody is a body CheckRedirect set on a followed redirect in place
// of the reader the sender opened for it, which it can still read. Closing it
// closes that reader too.
type substituteBody struct {
	io.ReadCloser
	opened *openedBody
}

func (b *substituteBody) Close() error {
	err := b.ReadCloser.Close()
	return errors.Join(err, b.opened.Close())
}

// openedBody is a request body the sender opened with lib.PreparedJSON.Open,
// which lets it tell its own readers from a caller's. Closing it more than
// once is safe.
type openedBody struct{ io.ReadCloser }

func openBody(content *lib.PreparedJSON) io.ReadCloser {
	return &openedBody{ReadCloser: content.Open()}
}

// Do sends req. Like Config.Do, it keeps the request's URL and credentials out
// of retry log entries above DEBUG and, outside debug mode, out of a returned
// *url.Error, which keeps its operation and cause.
func (s *replayableSender) Do(req *http.Request) (*http.Response, error) {
	s.client.Logger = s.config.retryLogFor(s.logger, req)
	response, err := (&http.Client{Transport: s}).Do(req)
	if s.config.InDebug() {
		return response, err
	}
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		err = &url.Error{Op: urlErr.Op, URL: "[redacted]", Err: urlErr.Err}
	}
	return response, err
}

// Close removes the copies of request bodies; call it once the request is
// done.
func (s *replayableSender) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	var errs []error
	for _, copied := range s.copies {
		errs = append(errs, copied.Close())
	}
	s.copies = nil
	return errors.Join(errs...)
}

// RoundTrip is retryablehttp.RoundTripper, except for how it keeps the body.
func (s *replayableSender) RoundTrip(req *http.Request) (*http.Response, error) {
	retryable := &retryablehttp.Request{Request: req.Clone(req.Context())}
	if body, ok := req.Body.(*replayableBody); ok && !body.read {
		err := retryable.SetBody(retryablehttp.ReaderFunc(func() (io.Reader, error) {
			return openBody(body.content), nil
		}))
		if err != nil {
			return nil, err
		}
	} else if req.Body != nil && req.Body != http.NoBody {
		copied, err := lib.CopyBody(req.Context(), req.Body)
		_ = req.Body.Close()
		if err != nil {
			return nil, err
		}
		s.mu.Lock()
		s.copies = append(s.copies, copied)
		s.mu.Unlock()
		newReader := func() io.Reader { return openBody(copied) }
		if copied.Len() == 0 {
			newReader = func() io.Reader { return http.NoBody }
		}
		err = retryable.SetBody(retryablehttp.ReaderFunc(func() (io.Reader, error) {
			return newReader(), nil
		}))
		if err != nil {
			return nil, err
		}
		retryable.GetBody = req.GetBody
		if getBody := req.GetBody; getBody != nil {
			// A redirect sends what GetBody returns through the transport, so
			// the SDK's own replayableBody is opened as an attempt's body is.
			retryable.GetBody = func() (io.ReadCloser, error) {
				body, err := getBody()
				if replayable, ok := body.(*replayableBody); ok && err == nil && !replayable.read {
					return openBody(replayable.content), nil
				}
				return body, err
			}
		}
	}
	retryable.ContentLength = req.ContentLength
	response, err := s.client.Do(retryable)
	return finishRetryableRoundTrip(response, err)
}

// replayableBody is a request body over prepared content, which each attempt
// reads again from the start.
type replayableBody struct {
	reader  io.Reader
	content *lib.PreparedJSON
	read    bool
}

func newReplayableBody(content *lib.PreparedJSON) *replayableBody {
	return &replayableBody{reader: content.NewReader(), content: content}
}

func (b *replayableBody) Read(p []byte) (int, error) {
	b.read = true
	return b.reader.Read(p)
}

func (b *replayableBody) Close() error {
	return nil
}

func (c Config) SetCustomClient(client *http.Client) Config {
	c.Client = lib.DefaultRetryableHttp(c, client)
	return c
}

// InDebug reports whether Debug is set or FILES_SDK_DEBUG is non-empty.
// A logger with levels controls which retry DEBUG entries it displays independently
// of this setting. A Printf-only logger receives retry DEBUG entries only in debug mode.
func (c Config) InDebug() bool {
	return c.Debug || (os.Getenv("FILES_SDK_DEBUG") != "")
}

// debugf logs a DEBUG diagnostic, which may hold full URLs and headers: at Debug level when the Logger has levels,
// otherwise through Printf as before.
func (c Config) debugf(format string, args ...any) {
	if leveled, ok := c.Logger.(lib.LeveledLogger); ok {
		leveled.Debug(format, args...)
		return
	}
	c.Printf(format, args...)
}

func (c Config) LogPath(path string, args map[string]interface{}) {
	c.Logger.Printf(libLog.New(path, args))
}

func (c Config) RootPath() string {
	return c.Endpoint() + APIPath
}

func (c Config) GetAPIKey() string {
	if c.SessionId != "" && c.APIKey == "" {
		return ""
	}
	return lib.DefaultString(c.APIKey, os.Getenv("FILES_API_KEY"))
}

func (c Config) SetUserAgentHeader(headers *http.Header) {
	if headers.Get("User-Agent") == "" {
		headers.Set("User-Agent", cmp.Or(c.UserAgent, defaultUserAgent))
	}
}

func (c Config) SetHeaders(headers *http.Header) {
	c.setHeadersForURL(headers, c.RootPath())
}

func (c Config) SetHeadersForRequest(req *http.Request) {
	c.setHeadersForURL(&req.Header, req.URL.String())
}

func (c Config) setHeadersForURL(headers *http.Header, rawURL string) {
	headers.Set("User-Agent", cmp.Or(c.UserAgent, defaultUserAgent))
	apiKey := c.GetAPIKey()
	if apiKey != "" {
		headers.Set(apiKeyHeader, apiKey)
	} else if c.SessionId != "" {
		headers.Set(sessionIdHeader, c.SessionId)
	}
	if c.Language != "" {
		headers.Set("Accept-Language", c.Language)
	}
	if c.WorkspaceId != 0 {
		headers.Set(workspaceIdHeader, strconv.FormatInt(c.WorkspaceId, 10))
	}
	for key, value := range c.AdditionalHeaders {
		headers.Set(key, value)
	}
	if !c.shouldSendAuthHeaders(rawURL) {
		clearAuthHeaders(headers)
	}
}

func (c Config) redirectSafeClient() *retryablehttp.Client {
	retrySource := c.Client
	if retrySource == nil {
		initialized := c.Init()
		retrySource = initialized.Client
	}
	return c.redirectSafeCopy(retrySource)
}

// redirectSafeCopy returns a copy of retrySource, with the same settings and hooks, whose HTTP client re-applies the
// URL-aware headers to each redirected request.
func (c Config) redirectSafeCopy(retrySource *retryablehttp.Client) *retryablehttp.Client {
	httpClient := http.Client{}
	if retrySource.HTTPClient != nil {
		httpClient = *retrySource.HTTPClient
	}
	originalCheckRedirect := httpClient.CheckRedirect
	// Go copies custom headers during redirects, including cross-origin redirects.
	// Re-apply URL-aware headers to each redirected request so Files auth is
	// stripped if a same-origin download URL redirects to a storage provider.
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		c.SetHeadersForRequest(req)
		if originalCheckRedirect != nil {
			return originalCheckRedirect(req, via)
		}
		// Keep net/http's default policy for a nil CheckRedirect: stop after 10 consecutive requests, with its error
		// text, which retryablehttp's default retry policy does not retry.
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		return nil
	}

	return &retryablehttp.Client{
		HTTPClient:      &httpClient,
		Logger:          retrySource.Logger,
		RetryWaitMin:    retrySource.RetryWaitMin,
		RetryWaitMax:    retrySource.RetryWaitMax,
		RetryMax:        retrySource.RetryMax,
		RequestLogHook:  retrySource.RequestLogHook,
		ResponseLogHook: retrySource.ResponseLogHook,
		CheckRetry:      retrySource.CheckRetry,
		Backoff:         retrySource.Backoff,
		ErrorHandler:    retrySource.ErrorHandler,
		PrepareRetry:    retrySource.PrepareRetry,
	}
}

func (c Config) shouldSendAuthHeaders(rawURL string) bool {
	if rawURL == "" {
		return false
	}

	destination, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	if destination.Host == "" {
		return true
	}

	endpoint, err := url.Parse(c.Endpoint())
	if err != nil {
		return false
	}

	return strings.EqualFold(destination.Scheme, endpoint.Scheme) && normalizedURLHost(destination) == normalizedURLHost(endpoint)
}

func clearAuthHeaders(headers *http.Header) {
	headers.Del(apiKeyHeader)
	headers.Del(sessionIdHeader)
	headers.Del(reauthenticationHeader)
	headers.Del(workspaceIdHeader)
}

// hasAuthHeaders reports whether headers carry any of the Files.com headers that clearAuthHeaders removes.
func hasAuthHeaders(headers http.Header) bool {
	return headers.Get(apiKeyHeader) != "" || headers.Get(sessionIdHeader) != "" || headers.Get(reauthenticationHeader) != "" || headers.Get(workspaceIdHeader) != ""
}

func normalizedURLHost(u *url.URL) string {
	host := strings.ToLower(u.Hostname())
	port := u.Port()
	// Treat default ports as equivalent to an omitted port for origin comparison.
	if port == "" || (u.Scheme == "https" && port == "443") || (u.Scheme == "http" && port == "80") {
		return host
	}
	return net.JoinHostPort(host, port)
}

func (c Config) FeatureFlag(flag string) bool {
	value, ok := c.FeatureFlags[flag]
	if !ok {
		panic(fmt.Sprintf("feature flag `%v` is not a valid flag", flag))
	}
	return value
}

const (
	FeatureFlagAdaptiveUploadV2        = "adaptive-upload-v2"
	FeatureFlagSparseRangeReads        = "sparse-range-reads"
	FeatureFlagUploadV2ChecksumTrailer = "upload-v2-checksum-trailer"
)

func FeatureFlags() map[string]bool {
	return map[string]bool{
		FeatureFlagAdaptiveUploadV2:        false,
		FeatureFlagSparseRangeReads:        false,
		FeatureFlagUploadV2ChecksumTrailer: false,
		"incremental-updates":              false,
	}
}
