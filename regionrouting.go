package files_sdk

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/hashicorp/go-retryablehttp"
)

// Lockout region mismatch routing for Files.com API requests.
//
// Call, its requests with attachments and CallRawAPI enroll each request as one operation once the caller's request
// options have run, just before it is sent. The operation lives in the request's context, and only the retry client
// copy made for that one request composes its hooks, so storage, upload part, download and other requests that are not
// enrolled are never moved. Following a mismatch is a retry of the same retryablehttp loop, within its RetryMax, and
// the configured CheckRetry, Backoff, PrepareRetry, log hooks and ErrorHandler keep their places and their invocation
// counts.

const (
	regionMismatchType         = "not-authenticated/lockout-region-mismatch"
	legacyRegionMismatchPrefix = "You have connected to a URL"
	regionHostHeader           = "X-Files-Host"
)

// Reasons a lockout region mismatch was not followed, as RegionMismatch.Reason gives them.
const (
	// RegionMismatchCrossedOrigin is an answer that came after the request, or one of its redirects, was sent to
	// another origin than the operation's, so it does not speak for the operation's origin.
	RegionMismatchCrossedOrigin = "crossed-origin"
	// RegionMismatchNoHost is a mismatch that named no host.
	RegionMismatchNoHost = "no-host"
	// RegionMismatchInvalidHost is a mismatch whose host is not an origin.
	RegionMismatchInvalidHost = "invalid-host"
	// RegionMismatchUntrustedHost is a mismatch whose origin may not receive the request's credentials.
	RegionMismatchUntrustedHost = "untrusted-host"
	// RegionMismatchSameOrigin is a mismatch that names the origin the request was sent to.
	RegionMismatchSameOrigin = "same-origin"
	// RegionMismatchDisabled is a mismatch not followed because Config.DisableRegionRelocation is set.
	RegionMismatchDisabled = "disabled"
	// RegionMismatchRetryBudget is a mismatch the client's RetryMax left no retry to follow.
	RegionMismatchRetryBudget = "retry-budget"
	// RegionMismatchBodyNotReplayable is a mismatch answering a request whose body cannot be sent again.
	RegionMismatchBodyNotReplayable = "body-not-replayable"
	// RegionMismatchRetryPolicy is a mismatch the client's CheckRetry refused to follow.
	RegionMismatchRetryPolicy = "retry-policy"
)

// RegionMismatch is a lockout region mismatch the API answered a request with: the site is served from another host.
// A configured CheckRetry receives it as its error, beside the real 403 response, when the SDK may follow it, and can
// refuse the move by returning false or an error. It holds no credentials.
type RegionMismatch struct {
	// Host is the host the API named, in data.host or X-Files-Host, or "" when it named none.
	Host string
	// Origin is the origin made of Host, such as https://site.files.com, or "" when Host is not one.
	Origin string
	// From is the origin the request was sent to.
	From string
	// Approved reports whether Origin may receive the request: an HTTPS origin that is files.com or one of its
	// subdomains on the default port, or one Config.TrustRegionOrigin approves.
	Approved bool
	// Reason is why the mismatch was not followed, one of the RegionMismatch reason constants, or "" when it may be.
	Reason string
}

func (m *RegionMismatch) Error() string {
	if m.Reason == "" {
		return fmt.Sprintf("lockout region mismatch: the site is served from %s", m.Origin)
	}
	return fmt.Sprintf("lockout region mismatch from %s naming %q was not followed: %s", m.From, m.Host, m.Reason)
}

// RegionMismatchError is the error of a call that a recognized lockout region mismatch ended: one that was not
// followed, including one no retry was left to follow. Err is the call's own error, such as the API's ResponseError, the
// error of the client's ErrorHandler or CheckRetry, or the retry client's own. Error returns Err's message, and
// errors.Is and errors.As reach both Err and the API's error for the mismatch: a ResponseError for the typed JSON
// mismatch, and the SDK's usual non-OK error, a lib.ResponseError, for the legacy plain-text one.
type RegionMismatchError struct {
	RegionMismatch
	Err   error
	cause error
}

func (e *RegionMismatchError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.RegionMismatch.Error()
}

// Unwrap returns the call's error and the API's error for the mismatch.
func (e *RegionMismatchError) Unwrap() []error {
	errs := make([]error, 0, 2)
	if e.Err != nil {
		errs = append(errs, e.Err)
	}
	if e.cause != nil {
		errs = append(errs, e.cause)
	}
	return errs
}

// regionAuthHeaders are the Files.com authentication headers an operation scopes to its origin.
var regionAuthHeaders = [4]string{apiKeyHeader, sessionIdHeader, reauthenticationHeader, workspaceIdHeader}

// regionAuthSnapshot is the Files.com authentication a request sent: the values of each of regionAuthHeaders, nil for
// a header it did not send. It never prints the values.
type regionAuthSnapshot struct {
	values [4][]string
}

// regionAuthIndex returns the index of name in regionAuthHeaders, compared as HTTP compares header names, or -1.
func regionAuthIndex(name string) int {
	for i, header := range regionAuthHeaders {
		if strings.EqualFold(name, header) {
			return i
		}
	}
	return -1
}

func snapshotRegionAuth(header http.Header) regionAuthSnapshot {
	var snapshot regionAuthSnapshot
	for name, values := range header {
		if i := regionAuthIndex(name); i >= 0 && len(values) > 0 {
			snapshot.values[i] = append(snapshot.values[i], values...)
		}
	}
	return snapshot
}

// apply makes header send exactly the snapshot's authentication, a header the snapshot did not send removed.
func (s regionAuthSnapshot) apply(header http.Header) {
	stripRegionAuth(header)
	for i, name := range regionAuthHeaders {
		if s.values[i] != nil {
			header[http.CanonicalHeaderKey(name)] = append([]string(nil), s.values[i]...)
		}
	}
}

func (s regionAuthSnapshot) secrets() []string {
	var secrets []string
	for _, values := range s.values {
		for _, value := range values {
			if value != "" {
				secrets = append(secrets, value)
			}
		}
	}
	return secrets
}

func (s regionAuthSnapshot) String() string { return "[redacted]" }

func (s regionAuthSnapshot) GoString() string { return "[redacted]" }

// stripRegionAuth removes every Files.com authentication header from header, whatever the case of its name.
func stripRegionAuth(header http.Header) {
	for name := range header {
		if regionAuthIndex(name) >= 0 {
			delete(header, name)
		}
	}
}

func hasRegionAuth(header http.Header) bool {
	for name := range header {
		if regionAuthIndex(name) >= 0 {
			return true
		}
	}
	return false
}

// regionOrigin returns the origin of u in one form: its scheme and host in lower case, and its port only when it is not
// the scheme's default.
func regionOrigin(u *url.URL) string {
	if u == nil {
		return ""
	}
	scheme := strings.ToLower(u.Scheme)
	host := strings.ToLower(u.Hostname())
	port := u.Port()
	if port == "" || (scheme == "https" && port == "443") || (scheme == "http" && port == "80") {
		if strings.Contains(host, ":") {
			host = "[" + host + "]"
		}
		return scheme + "://" + host
	}
	return scheme + "://" + net.JoinHostPort(host, port)
}

// regionOriginOfHost returns the origin a mismatch's host names: a host, HTTPS on its default port unless it gives
// another, or an HTTP or HTTPS origin with a path of at most "/". Anything else is not an origin.
func regionOriginOfHost(host string) (string, bool) {
	if host == "" {
		return "", false
	}
	raw := host
	if !strings.Contains(host, "://") {
		raw = "https://" + host
	}
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Opaque != "" || parsed.User != nil || parsed.Host == "" || parsed.Hostname() == "" ||
		parsed.RawQuery != "" || parsed.ForceQuery || parsed.Fragment != "" || (parsed.Path != "" && parsed.Path != "/") {
		return "", false
	}
	if scheme := strings.ToLower(parsed.Scheme); scheme != "https" && scheme != "http" {
		return "", false
	}
	return regionOrigin(parsed), true
}

// regionTrusted reports whether origin may receive a request moved from from: an HTTPS origin, files.com or one of
// its subdomains on the default port, or one trust approves. An HTTP origin is never trusted, and trust is not asked
// about one.
func regionTrusted(origin, from string, trust func(origin, from string) bool) bool {
	parsed, err := url.Parse(origin)
	if err != nil || parsed.Scheme != "https" {
		return false
	}
	host := parsed.Hostname()
	if parsed.Port() == "" && (host == "files.com" || strings.HasSuffix(host, ".files.com")) {
		return true
	}
	return trust != nil && trust(origin, from)
}

// regionAnswer is a recognized lockout region mismatch: the host it names, or "", and the API's error for it, as the
// SDK's own parsers make it.
type regionAnswer struct {
	host  string
	cause error
}

// regionLegacyDiagnosticBytes is the most of a body the SDK's non-OK error reports, which a recognized legacy mismatch's
// error is made from.
const regionLegacyDiagnosticBytes = 512

// classifyRegionAnswer recognizes a lockout region mismatch in resp: a 403 that is the API's typed
// not-authenticated/lockout-region-mismatch error, as APIError decodes it, whose host is its data.host or else its
// single X-Files-Host header; or the legacy 403 whose body begins "You have connected to a URL" and that has an
// X-Files-Host header. A JSON body is read in full, as APIError reads it, whatever its size. Of any other body, only one
// with an X-Files-Host header is read, a byte at a time and no further than the first byte that differs from the legacy
// prefix or, when it begins with that prefix, as far as the SDK's non-OK error reports. What is read is put back in
// front of the rest of the body, which reads as it would have, ending with the same error if reading failed, and closes
// the original body.
func classifyRegionAnswer(resp *http.Response) (*regionAnswer, bool) {
	if resp == nil || resp.StatusCode != http.StatusForbidden || resp.Body == nil || resp.Body == http.NoBody {
		return nil, false
	}
	hosts := resp.Header.Values(regionHostHeader)
	if lib.IsJSON(resp) {
		data, err := peekRegionBody(resp, -1)
		if err != nil {
			return nil, false
		}
		probe := *resp
		probe.Body = io.NopCloser(bytes.NewReader(data))
		var typed ResponseError
		if errors.As(APIError()(&probe), &typed) && typed.Type == regionMismatchType {
			host := typed.Data.Host
			if host == "" && len(hosts) == 1 {
				host = hosts[0]
			}
			return &regionAnswer{host: host, cause: typed}, true
		}
		return legacyRegionAnswer(resp, hosts, data)
	}
	if len(hosts) == 0 {
		return nil, false
	}
	matched, err := peekRegionPrefix(resp, []byte(legacyRegionMismatchPrefix))
	if err != nil || !matched {
		return nil, false
	}
	data, err := peekRegionBody(resp, regionLegacyDiagnosticBytes)
	if err != nil {
		return nil, false
	}
	return legacyRegionAnswer(resp, hosts, data)
}

// legacyRegionAnswer recognizes the legacy mismatch in resp, whose X-Files-Host headers are hosts, from data, the start
// of its body, and makes its error as the SDK's non-OK error does.
func legacyRegionAnswer(resp *http.Response, hosts []string, data []byte) (*regionAnswer, bool) {
	if len(hosts) == 0 || !bytes.HasPrefix(data, []byte(legacyRegionMismatchPrefix)) {
		return nil, false
	}
	host := ""
	if len(hosts) == 1 {
		host = hosts[0]
	}
	probe := *resp
	probe.Body = io.NopCloser(bytes.NewReader(data))
	return &regionAnswer{host: host, cause: lib.NonOkError(&probe)}, true
}

// regionPeekedBody is a response body the classifier has read from: the bytes it read, then the rest of the original
// body, which it closes.
type regionPeekedBody struct {
	io.Reader
	io.Closer
}

// peekRegionBody reads up to limit bytes of resp's body, or all of it when limit is negative, and puts them back in front
// of the rest, so that the body reads as it would have, ending with the same error if reading failed, and closing it
// closes the original body.
func peekRegionBody(resp *http.Response, limit int64) ([]byte, error) {
	original := resp.Body
	var data []byte
	var err error
	if limit < 0 {
		data, err = io.ReadAll(original)
	} else {
		data, err = io.ReadAll(io.LimitReader(original, limit))
	}
	var rest io.Reader = original
	if limit < 0 || err != nil {
		rest = readResult{err}
	}
	resp.Body = &regionPeekedBody{Reader: io.MultiReader(bytes.NewReader(data), rest), Closer: original}
	return data, err
}

// peekRegionPrefix reports whether resp's body begins with prefix, reading it a byte at a time and no further than the
// first byte that differs, so that it never waits for more of a body that cannot begin with prefix. It puts what it
// read back in front of the rest as peekRegionBody does, and returns the body's error, other than its end.
func peekRegionPrefix(resp *http.Response, prefix []byte) (bool, error) {
	original := resp.Body
	var data []byte
	var err error
	for len(data) < len(prefix) && err == nil && bytes.HasPrefix(prefix, data) {
		var next [1]byte
		var n int
		n, err = original.Read(next[:])
		data = append(data, next[:n]...)
	}
	var rest io.Reader = original
	if err != nil {
		rest = readResult{err}
	}
	resp.Body = &regionPeekedBody{Reader: io.MultiReader(bytes.NewReader(data), rest), Closer: original}
	if err == io.EOF {
		err = nil
	}
	return bytes.Equal(data, prefix), err
}

type regionOperationKey struct{}

func regionOperationOf(ctx context.Context) *regionOperation {
	if ctx == nil {
		return nil
	}
	op, _ := ctx.Value(regionOperationKey{}).(*regionOperation)
	return op
}

// regionPhase is where an operation's current attempt is in the retry loop.
type regionPhase int

const (
	regionPhaseEnrolled regionPhase = iota
	// regionPhaseOpening: an attempt's body is being opened.
	regionPhaseOpening
	// regionPhaseSending: an attempt has been handed to the HTTP client.
	regionPhaseSending
	// regionPhaseDecided: the retry policy has decided on the attempt's response or failure.
	regionPhaseDecided
	// regionPhaseWaiting: an accepted retry is waiting its backoff.
	regionPhaseWaiting
	// regionPhasePreparing: an accepted retry is being prepared.
	regionPhasePreparing
	// regionPhaseHandled: the ErrorHandler has selected the call's outcome.
	regionPhaseHandled
)

// regionVerdict is the current attempt's recognized mismatch that was not followed, which ends the call unless a later
// phase begins: its description, the API's error for it, and the response it was recognized in, or nil once the call's
// outcome is an error, with what that response's classification saw: its status, the body it put back and the headers
// that decide a mismatch.
type regionVerdict struct {
	mismatch    RegionMismatch
	cause       error
	response    *http.Response
	status      int
	body        *regionPeekedBody
	contentType string
	hosts       []string
}

func newRegionVerdict(mismatch RegionMismatch, cause error, resp *http.Response) *regionVerdict {
	verdict := &regionVerdict{mismatch: mismatch, cause: cause, response: resp}
	if resp != nil {
		verdict.status = resp.StatusCode
		verdict.body, _ = resp.Body.(*regionPeekedBody)
		verdict.contentType = resp.Header.Get("Content-Type")
		verdict.hosts = slices.Clone(resp.Header.Values(regionHostHeader))
	}
	return verdict
}

// describes reports whether resp is still the response the verdict was recognized in, as its classification saw it:
// the same response, with the same status, Content-Type and X-Files-Host headers, and the body classification put back,
// which the retry client may have guarded since. Its body need not be read again, and may already be closed.
func (v *regionVerdict) describes(resp *http.Response) bool {
	if resp == nil || resp != v.response || v.body == nil || resp.StatusCode != v.status ||
		resp.Header.Get("Content-Type") != v.contentType || !slices.Equal(resp.Header.Values(regionHostHeader), v.hosts) {
		return false
	}
	body := resp.Body
	if guarded, ok := body.(*retryableErrorBody); ok {
		body = guarded.ReadCloser
	}
	peeked, ok := body.(*regionPeekedBody)
	return ok && peeked == v.body
}

// regionMove is an approved move: the origin the next attempt goes to, with the authentication the answered request
// sent.
type regionMove struct {
	origin string
	auth   regionAuthSnapshot
}

// regionOperation is one enrolled API request, from its first attempt to its outcome. The retry loop runs its hooks one
// at a time; the mutex keeps the logger's reads safe too.
type regionOperation struct {
	mu         sync.Mutex
	disabled   bool
	trust      func(origin, from string) bool
	replayable bool
	retryMax   int
	// configHeaders reports whether the operation's redirects get the configuration's headers other than Files.com
	// authentication again, as its sender's redirects did before the operation was enrolled.
	configHeaders bool

	// origin is the operation's current API origin, the configured one until a mismatch moves it.
	origin  string
	moved   bool
	attempt int
	phase   regionPhase
	// crossed records that the current attempt, or one of its redirects, was sent to another origin.
	crossed bool
	// sentURL and sentHost are where the current attempt was sent, as the caller's hooks and the SDK's guard left its
	// request before the HTTP client had it.
	sentURL  *url.URL
	sentHost string
	// lastAuth is the authentication the current attempt's latest request sent, after the caller's hooks and the
	// SDK's guard; originAuth the latest it sent at the operation's origin.
	lastAuth   regionAuthSnapshot
	originAuth regionAuthSnapshot
	pending    *regionMove
	verdict    *regionVerdict
	// answered is the current attempt's answer as the retry policy received it; handled records that the ErrorHandler
	// selected the call's outcome; final is the context of the answer the call returned, once CallRawAPI settles it.
	answered *regionAnswered
	handled  bool
	final    RelocatedAPIContext
	// secrets are values retry log entries above DEBUG leave out: relocated URLs and credentials sent.
	secrets []string
}

// regionAnswered is an attempt's answer: the response and request the retry policy received, the request's URL as it
// was then, the origin the attempt was sent to, whether it or one of its redirects left that origin, and the
// authentication it sent.
type regionAnswered struct {
	response   *http.Response
	request    *http.Request
	requestURL string
	origin     string
	crossed    bool
	auth       regionAuthSnapshot
}

// carries reports whether response is this answer: a response for the answer's request, whose URL is still the one it
// was answered at, or, when the answer had no request, the answer itself, returned without an ErrorHandler selecting it.
func (a *regionAnswered) carries(response *http.Response, handled bool) bool {
	if response == nil {
		return false
	}
	if a.request != nil {
		return response.Request == a.request && a.request.URL != nil && a.request.URL.String() == a.requestURL
	}
	return response == a.response && !handled
}

// enrollRegionOperation returns req as a request of a new operation, which replayable says can send its body again.
// Its redirects get the configuration's headers other than Files.com authentication again, as Config.Do's and the
// attachment sender's always did.
func (c Config) enrollRegionOperation(req *http.Request, replayable bool) (*http.Request, *regionOperation) {
	if req == nil || req.URL == nil {
		return req, nil
	}
	op := &regionOperation{
		disabled:      c.DisableRegionRelocation,
		trust:         c.TrustRegionOrigin,
		replayable:    replayable,
		configHeaders: true,
		origin:        regionOrigin(req.URL),
		lastAuth:      snapshotRegionAuth(req.Header),
	}
	op.originAuth = op.lastAuth
	return req.WithContext(context.WithValue(req.Context(), regionOperationKey{}, op)), op
}

// composeRegionClient installs op's hooks on client, the retry client copy made for op's request alone, around the
// hooks it already has.
func composeRegionClient(client *retryablehttp.Client, op *regionOperation) {
	if client == nil || op == nil {
		return
	}
	op.mu.Lock()
	op.retryMax = client.RetryMax
	op.mu.Unlock()
	if checkRetry := client.CheckRetry; checkRetry != nil {
		client.CheckRetry = op.checkRetry(checkRetry)
	}
	if backoff := client.Backoff; backoff != nil {
		client.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
			op.enter(regionPhaseWaiting)
			return backoff(min, max, attemptNum, resp)
		}
	}
	requestLogHook := client.RequestLogHook
	client.RequestLogHook = func(logger retryablehttp.Logger, req *http.Request, attempt int) {
		if requestLogHook != nil {
			requestLogHook(logger, req, attempt)
		}
		op.beginSend(req, attempt)
	}
	prepareRetry := client.PrepareRetry
	client.PrepareRetry = func(req *http.Request) error {
		op.enter(regionPhasePreparing)
		op.restoreSend(req)
		op.applyMove(req)
		if prepareRetry != nil {
			return prepareRetry(req)
		}
		return nil
	}
	if handler := client.ErrorHandler; handler != nil {
		client.ErrorHandler = func(resp *http.Response, err error, attempts int) (*http.Response, error) {
			selected, selectedErr := handler(resp, err, attempts)
			op.reconcile(selected, selectedErr)
			return selected, selectedErr
		}
	}
}

// enter records that the operation has reached phase: an attempt's body, or an accepted retry's wait or preparation,
// after which the previous decision no longer ends the call and the previous answer is not the one it returns.
func (op *regionOperation) enter(phase regionPhase) {
	if op == nil {
		return
	}
	op.mu.Lock()
	defer op.mu.Unlock()
	op.phase = phase
	op.verdict = nil
	op.answered = nil
}

// beginSend runs once the caller's RequestLogHook has had the attempt's request, the last change before it is sent:
// the request keeps no Files.com authentication when it leaves the operation's origin, and where it goes and what it
// sends are recorded.
func (op *regionOperation) beginSend(req *http.Request, attempt int) {
	op.mu.Lock()
	defer op.mu.Unlock()
	op.attempt = attempt
	op.phase = regionPhaseSending
	op.verdict = nil
	op.answered = nil
	if req == nil {
		return
	}
	op.crossed = regionOrigin(req.URL) != op.origin
	if op.crossed && hasRegionAuth(req.Header) {
		header := req.Header.Clone()
		stripRegionAuth(header)
		req.Header = header
	}
	op.sentURL = nil
	if req.URL != nil {
		sent := *req.URL
		op.sentURL = &sent
	}
	op.sentHost = req.Host
	op.lastAuth = snapshotRegionAuth(req.Header)
	if !op.crossed {
		op.originAuth = op.lastAuth
	}
	op.secrets = append(op.secrets, op.lastAuth.secrets()...)
}

// restoreSend gives an accepted retry the destination its attempt was sent to, as the caller's hooks and the SDK's
// guard left it, before the retry's own hooks run: a transport that rewrote the request it was given in place, below
// the SDK, does not choose where the retry goes or make it leave the operation's origin. The caller's PrepareRetry and
// RequestLogHook can still send it elsewhere, guarded as any attempt is.
func (op *regionOperation) restoreSend(req *http.Request) {
	op.mu.Lock()
	defer op.mu.Unlock()
	if req == nil || op.sentURL == nil {
		return
	}
	sent := *op.sentURL
	req.URL = &sent
	req.Host = op.sentHost
}

// scopeRedirect gives a redirect the operation's own authentication, in place of the configuration's and the
// environment's: what its attempt sent at the operation's origin there, and none elsewhere.
func (op *regionOperation) scopeRedirect(req *http.Request) {
	op.mu.Lock()
	defer op.mu.Unlock()
	if req.Header == nil {
		req.Header = http.Header{}
	}
	if regionOrigin(req.URL) == op.origin {
		op.originAuth.apply(req.Header)
	} else {
		stripRegionAuth(req.Header)
	}
}

// guardRedirect runs once the caller's CheckRedirect has let a redirect go ahead: a redirect that now leads away from
// the operation's origin keeps no Files.com authentication, and its answer cannot move the operation.
func (op *regionOperation) guardRedirect(req *http.Request) {
	if op == nil || req == nil {
		return
	}
	op.mu.Lock()
	defer op.mu.Unlock()
	if regionOrigin(req.URL) != op.origin {
		stripRegionAuth(req.Header)
		op.crossed = true
	}
	op.lastAuth = snapshotRegionAuth(req.Header)
	if !op.crossed {
		op.originAuth = op.lastAuth
	}
	op.secrets = append(op.secrets, op.lastAuth.secrets()...)
}

// checkRetry returns inner with each recognized mismatch decided first. A mismatch that may be followed reaches inner as
// its error, beside the real 403, and is followed when inner retries; any other answer reaches inner unchanged.
func (op *regionOperation) checkRetry(inner retryablehttp.CheckRetry) retryablehttp.CheckRetry {
	return func(ctx context.Context, resp *http.Response, doErr error) (bool, error) {
		op.mu.Lock()
		op.phase = regionPhaseDecided
		op.verdict = nil
		op.pending = nil
		op.answered = nil
		if doErr == nil && resp != nil {
			op.answered = &regionAnswered{response: resp, request: resp.Request, origin: op.origin, crossed: op.crossed, auth: op.lastAuth}
			if resp.Request != nil && resp.Request.URL != nil {
				op.answered.requestURL = resp.Request.URL.String()
			}
		}
		op.mu.Unlock()
		if doErr != nil || resp == nil || resp.StatusCode != http.StatusForbidden || ctx.Err() != nil {
			return inner(ctx, resp, doErr)
		}
		answer, ok := classifyRegionAnswer(resp)
		if !ok {
			return inner(ctx, resp, doErr)
		}
		mismatch := op.decide(answer)
		if mismatch.Reason != "" {
			op.record(mismatch, answer.cause, resp)
			return inner(ctx, resp, doErr)
		}
		cause := mismatch
		retry, err := inner(ctx, resp, &cause)
		if retry {
			op.mu.Lock()
			op.pending = &regionMove{origin: mismatch.Origin, auth: op.lastAuth}
			op.mu.Unlock()
			return retry, err
		}
		mismatch.Reason = RegionMismatchRetryPolicy
		op.record(mismatch, answer.cause, resp)
		return retry, err
	}
}

// decide describes answer, a mismatch answering the current attempt, with why it may not be followed, if it may not.
func (op *regionOperation) decide(answer *regionAnswer) RegionMismatch {
	op.mu.Lock()
	from, crossed, budget := op.origin, op.crossed, op.retryMax-op.attempt
	op.mu.Unlock()
	mismatch := RegionMismatch{Host: answer.host, From: from}
	if origin, ok := regionOriginOfHost(answer.host); ok {
		mismatch.Origin = origin
		mismatch.Approved = regionTrusted(origin, from, op.trust)
	}
	switch {
	case crossed:
		mismatch.Reason = RegionMismatchCrossedOrigin
	case mismatch.Host == "":
		mismatch.Reason = RegionMismatchNoHost
	case mismatch.Origin == "":
		mismatch.Reason = RegionMismatchInvalidHost
	case !mismatch.Approved:
		mismatch.Reason = RegionMismatchUntrustedHost
	case mismatch.Origin == from:
		mismatch.Reason = RegionMismatchSameOrigin
	case op.disabled:
		mismatch.Reason = RegionMismatchDisabled
	case budget <= 0:
		mismatch.Reason = RegionMismatchRetryBudget
	case !op.replayable:
		mismatch.Reason = RegionMismatchBodyNotReplayable
	}
	return mismatch
}

func (op *regionOperation) record(mismatch RegionMismatch, cause error, resp *http.Response) {
	op.mu.Lock()
	defer op.mu.Unlock()
	op.phase = regionPhaseDecided
	op.verdict = newRegionVerdict(mismatch, cause, resp)
}

// applyMove sends the next attempt to the approved origin, if a move is pending: its URL with only the scheme and host
// changed, no explicit Host, and the authentication the answered request sent.
func (op *regionOperation) applyMove(req *http.Request) {
	op.mu.Lock()
	defer op.mu.Unlock()
	move := op.pending
	op.pending = nil
	if move == nil || req == nil || req.URL == nil {
		return
	}
	target, err := url.Parse(move.origin)
	if err != nil {
		return
	}
	relocated := *req.URL
	relocated.Scheme = target.Scheme
	relocated.Host = target.Host
	req.URL = &relocated
	req.Host = ""
	header := req.Header.Clone()
	if header == nil {
		header = http.Header{}
	}
	move.auth.apply(header)
	req.Header = header
	op.origin = move.origin
	op.moved = true
	op.originAuth = move.auth
	op.secrets = append(op.secrets, relocated.String())
}

// verdictLocked returns the decision that ends the call, if the call has gone no further than deciding on it.
func (op *regionOperation) verdictLocked() *regionVerdict {
	if op.phase != regionPhaseDecided && op.phase != regionPhaseHandled {
		return nil
	}
	return op.verdict
}

// reconcile keeps the current mismatch with the outcome the ErrorHandler selected, by what that outcome is, its
// response and its error together: the mismatch response as it was classified, or no response, keeps it; any other
// response, or the mismatch response once the handler has changed what classified it, keeps a mismatch only while it
// is itself a recognized one, whatever error came with it, so a converted success or another response drops it. The
// mismatch stays beside a response only when the handler returned no error.
func (op *regionOperation) reconcile(selected *http.Response, selectedErr error) {
	op.mu.Lock()
	op.handled = true
	verdict := op.verdictLocked()
	if verdict == nil {
		op.mu.Unlock()
		return
	}
	op.phase = regionPhaseHandled
	if selected == nil || verdict.describes(selected) {
		kept := *verdict
		if selected == nil || selectedErr != nil {
			kept.response = nil
		}
		op.verdict = &kept
		op.mu.Unlock()
		return
	}
	op.mu.Unlock()
	answer, ok := classifyRegionAnswer(selected)
	if !ok {
		op.mu.Lock()
		op.verdict = nil
		op.mu.Unlock()
		return
	}
	mismatch := op.decide(answer)
	if mismatch.Reason == "" {
		mismatch.Reason = verdict.mismatch.Reason
	}
	next := newRegionVerdict(mismatch, answer.cause, selected)
	if selectedErr != nil {
		next.response = nil
	}
	op.mu.Lock()
	op.verdict = next
	op.mu.Unlock()
}

// terminalError returns err, the call's error, with the mismatch that ended the call, if one did. response is the
// response err was parsed from, or nil for an error the request itself ended with.
func (op *regionOperation) terminalError(response *http.Response, err error) error {
	if op == nil || err == nil {
		return err
	}
	op.mu.Lock()
	verdict := op.verdictLocked()
	op.mu.Unlock()
	if verdict == nil || (response != nil && verdict.response != response) {
		return err
	}
	return &RegionMismatchError{RegionMismatch: verdict.mismatch, Err: err, cause: verdict.cause}
}

// rawTerminal is terminalError for CallRawAPI, whose recognized mismatch response is returned with its error: the
// API's error for the mismatch.
func (op *regionOperation) rawTerminal(response *http.Response, err error) error {
	if op == nil {
		return err
	}
	if err != nil {
		return op.terminalError(nil, err)
	}
	op.mu.Lock()
	verdict := op.verdictLocked()
	op.mu.Unlock()
	if verdict == nil || response == nil || verdict.response != response {
		return nil
	}
	return &RegionMismatchError{RegionMismatch: verdict.mismatch, Err: verdict.cause}
}

// settle fixes, once CallRawAPI has returned response, the context of the answer it is: the origin a lockout region
// mismatch moved the request to and the authentication the answered attempt sent there, when response carries that
// attempt's answer, as the attempt received it or in a replacement the ErrorHandler made for it. A call that failed
// without an answer, returned another response, or was answered from another origin has none.
func (op *regionOperation) settle(response *http.Response) {
	if op == nil {
		return
	}
	op.mu.Lock()
	defer op.mu.Unlock()
	op.final = RelocatedAPIContext{}
	answered := op.answered
	if !op.moved || answered == nil || answered.crossed || answered.origin != op.origin || !answered.carries(response, op.handled) {
		return
	}
	op.final = RelocatedAPIContext{origin: answered.origin, auth: answered.auth}
}

func (op *regionOperation) relocatedContext() RelocatedAPIContext {
	op.mu.Lock()
	defer op.mu.Unlock()
	return op.final
}

// setConfigHeaders sets whether the operation's redirects get the configuration's headers other than Files.com
// authentication again.
func (op *regionOperation) setConfigHeaders(reapply bool) {
	if op == nil {
		return
	}
	op.mu.Lock()
	defer op.mu.Unlock()
	op.configHeaders = reapply
}

// reappliesConfigHeaders reports whether the operation's redirects get the configuration's headers other than
// Files.com authentication again.
func (op *regionOperation) reappliesConfigHeaders() bool {
	op.mu.Lock()
	defer op.mu.Unlock()
	return op.configHeaders
}

func (op *regionOperation) relocatedSecrets() []string {
	op.mu.Lock()
	defer op.mu.Unlock()
	return append([]string(nil), op.secrets...)
}

// regionSecrets leaves an operation's relocated URLs and the credentials it sent out of its retry log entries above
// DEBUG, besides the request's own secrets.
type regionSecrets struct {
	base requestSecrets
	op   *regionOperation
}

func (s regionSecrets) all() requestSecrets {
	return append(requestSecrets(s.op.relocatedSecrets()), s.base...)
}

func (s regionSecrets) text(text string) string {
	return s.all().text(text)
}

func (s regionSecrets) values(args []interface{}) []interface{} {
	return s.all().values(args)
}

// RelocatedAPIContext is the context of an API answer from an origin a lockout region mismatch moved its request to:
// that origin, and the Files.com authentication headers the request sent there, a header it did not send included.
// It holds no response and never prints the credentials. The zero value is an answer from the configured endpoint.
type RelocatedAPIContext struct {
	origin string
	auth   regionAuthSnapshot
}

// Relocated reports whether a lockout region mismatch moved the request to the origin that answered.
func (c RelocatedAPIContext) Relocated() bool { return c.origin != "" }

// Origin returns the origin that answered a moved request, or "".
func (c RelocatedAPIContext) Origin() string { return c.origin }

// String returns the origin, without the credentials.
func (c RelocatedAPIContext) String() string { return c.origin }

// GoString returns the origin, without the credentials.
func (c RelocatedAPIContext) GoString() string {
	return fmt.Sprintf("files_sdk.RelocatedAPIContext(%q)", c.origin)
}

// ResolveReference resolves reference, a URI the answer gave, for a request that follows it, such as a ZIP download's
// download_uri. An absolute reference is returned as it parses, with config unchanged, and so is a relative one
// resolved against config.Endpoint() when the request was not moved. A relative reference after a move is resolved
// against config.Endpoint() with only its scheme and host replaced by the origin that answered; the Config returned is
// then a copy of config, for that request alone, that sends the Files.com authentication the moved request sent, a
// header it did not send included, and never config's own or the environment's.
func (c RelocatedAPIContext) ResolveReference(config Config, reference string) (string, Config, error) {
	parsed, err := url.Parse(reference)
	if err != nil {
		return "", config, err
	}
	if parsed.IsAbs() {
		return parsed.String(), config, nil
	}
	base, err := url.Parse(config.Endpoint())
	if err != nil {
		return "", config, err
	}
	if !c.Relocated() {
		return base.ResolveReference(parsed).String(), config, nil
	}
	origin, err := url.Parse(c.origin)
	if err != nil {
		return "", config, err
	}
	base.Scheme, base.Host = origin.Scheme, origin.Host
	scoped := config
	scoped.EndpointOverride = base.String()
	auth := c.auth
	scoped.regionAuth = &auth
	return base.ResolveReference(parsed).String(), scoped, nil
}

// RelocatedAPIContext returns, once CallRawAPI has sent p, the context of the answer it returned when a lockout region
// mismatch moved it there; otherwise, including when building or sending the request failed, and after CallRaw, the
// zero RelocatedAPIContext.
func (p *CallParams) RelocatedAPIContext() RelocatedAPIContext {
	if p == nil || p.region == nil {
		return RelocatedAPIContext{}
	}
	return p.region.relocatedContext()
}
