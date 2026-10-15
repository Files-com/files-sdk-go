package files_sdk

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/netip"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/hashicorp/go-retryablehttp"
	quic "github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

// The connect budgets are package variables instead of constants so tests can
// shorten them; production code treats them as constants.
var (
	// directTransferConnectTimeout caps the total time spent connecting
	// to all direct candidates before the caller falls back to the proxy URL.
	directTransferConnectTimeout = 5 * time.Second
	// directTransferPerCandidateTimeout caps the connect/handshake time
	// for each candidate address so a single stalled candidate cannot consume
	// the whole connect budget and starve the remaining candidates.
	directTransferPerCandidateTimeout   = 2 * time.Second
	directTransferResponseHeaderTimeout = 60 * time.Second
)

var (
	ErrDirectTransferUnavailable     = errors.New("direct transfer unavailable")
	ErrDirectTransferResponseStarted = errors.New("direct transfer response processing started")
)

// DirectTransferResponseError describes an unsuccessful response from an Agent's
// direct endpoint. StatusCode contains the HTTP status; RetryAfter contains the
// parsed Retry-After delay, or zero when no valid delay is available.
type DirectTransferResponseError struct {
	StatusCode int
	RetryAfter time.Duration
}

// Error returns a message containing the direct endpoint's HTTP status.
func (e *DirectTransferResponseError) Error() string {
	return fmt.Sprintf("direct transfer returned status %d", e.StatusCode)
}

// Unwrap returns ErrDirectTransferUnavailable for errors.Is matching.
func (e *DirectTransferResponseError) Unwrap() error {
	return ErrDirectTransferUnavailable
}

type directTransferClientCacheKey struct {
	serverName string
	caPEM      string
	protocol   string
	candidates string
}

type directTransferClientCache struct {
	mu         sync.Mutex
	key        directTransferClientCacheKey
	httpClient *http.Client
}

type directTransferClientCacheContextKey struct{}

// WithDirectTransferClientCache returns a context that reuses a direct HTTP
// client within one transfer, and a cleanup function that closes idle connections.
// Call the cleanup function when the transfer finishes. If ctx already contains
// a cache, the returned cleanup function does nothing.
func WithDirectTransferClientCache(ctx context.Context) (context.Context, func()) {
	if _, ok := ctx.Value(directTransferClientCacheContextKey{}).(*directTransferClientCache); ok {
		return ctx, func() {}
	}
	cache := &directTransferClientCache{}
	return context.WithValue(ctx, directTransferClientCacheContextKey{}, cache), cache.close
}

func (c *directTransferClientCache) client(key directTransferClientCacheKey, build func() (*http.Client, error)) (*http.Client, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.httpClient != nil {
		if c.key != key {
			return nil, ErrDirectTransferUnavailable
		}
		return c.httpClient, nil
	}
	client, err := build()
	if err == nil {
		c.key = key
		c.httpClient = client
	}
	return client, err
}

func (c *directTransferClientCache) close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.httpClient != nil {
		c.httpClient.CloseIdleConnections()
		c.httpClient = nil
	}
}

var directTransferBlockedIPPrefixes = []netip.Prefix{
	netip.MustParsePrefix("0.0.0.0/8"),
	netip.MustParsePrefix("100.64.0.0/10"),
	netip.MustParsePrefix("127.0.0.0/8"),
	netip.MustParsePrefix("169.254.0.0/16"),
	netip.MustParsePrefix("192.0.0.0/24"),
	netip.MustParsePrefix("192.0.2.0/24"),
	netip.MustParsePrefix("198.18.0.0/15"),
	netip.MustParsePrefix("198.51.100.0/24"),
	netip.MustParsePrefix("203.0.113.0/24"),
	netip.MustParsePrefix("224.0.0.0/4"),
	netip.MustParsePrefix("240.0.0.0/4"),
	netip.MustParsePrefix("::/128"),
	netip.MustParsePrefix("::1/128"),
	netip.MustParsePrefix("64:ff9b:1::/48"),
	netip.MustParsePrefix("100::/64"),
	netip.MustParsePrefix("2001::/23"),
	netip.MustParsePrefix("2001:2::/48"),
	netip.MustParsePrefix("2001:db8::/32"),
	netip.MustParsePrefix("fc00::/7"),
	netip.MustParsePrefix("fe80::/10"),
	netip.MustParsePrefix("ff00::/8"),
}

// DirectConnectionInfoPresent reports whether info has the supported version
// and the fields required for a direct connection. It does not validate the
// addresses, certificate, or reachability of the Agent.
func DirectConnectionInfoPresent(info DirectConnectionInfo) bool {
	return info.Version == 1 &&
		info.ServerName != "" &&
		len(info.Addresses) > 0 &&
		info.DirectUri != "" &&
		info.CaPem != ""
}

// DirectTransferRetryableClient returns the direct URL and a retryable HTTP client
// for info. It returns an error when direct transfers are disabled or info cannot
// be used to configure a direct connection.
func DirectTransferRetryableClient(ctx context.Context, config Config, info DirectConnectionInfo) (string, *retryablehttp.Client, error) {
	directURL, httpClient, err := directTransferHTTPClient(ctx, config, info)
	if err != nil {
		return "", nil, err
	}

	return directURL, lib.DefaultRetryableHttp(config.Logger, httpClient), nil
}

// WrapDirectTransferOptions sends a copy of request to the Agent's direct endpoint
// and applies the supplied request and response options. The caller must close the
// body of a successful response. An unsuccessful HTTP response is closed here and
// returned with a DirectTransferResponseError.
func WrapDirectTransferOptions(config Config, info DirectConnectionInfo, request *http.Request, opts ...RequestResponseOption) (*http.Response, error) {
	if config.DisableDirectTransfers || !DirectConnectionInfoPresent(info) {
		return nil, ErrDirectTransferUnavailable
	}

	directAttemptRequest := request.Clone(request.Context())
	modifiedRequest, err := BuildRequest(directAttemptRequest, opts...)
	if err != nil {
		return nil, err
	}

	response, err := sendDirectTransferRequest(config, info, nil, nil, modifiedRequest)
	if err != nil {
		return response, err
	}

	processedResponse, err := BuildResponse(response, opts...)
	if err != nil {
		lib.CloseBody(response)
		return response, errors.Join(ErrDirectTransferResponseStarted, err)
	}
	return processedResponse, nil
}

// sendDirectTransferRequest sends a credential-free copy of request to info's
// direct endpoint, using cache when the request context carries none. An
// unsuccessful HTTP response is closed and returned with a
// DirectTransferResponseError.
func sendDirectTransferRequest(config Config, info DirectConnectionInfo, cache *directTransferClientCache, wrapTransport func(http.RoundTripper) http.RoundTripper, request *http.Request) (*http.Response, error) {
	directURL, httpClient, err := directTransferHTTPClientWithCache(request.Context(), config, info, cache, wrapTransport)
	if err != nil {
		return nil, err
	}

	directRequest, err := cloneRequestForDirectTransfer(request, directURL)
	if err != nil {
		return nil, err
	}

	logDirectTransferRequest(config.Logger, directRequest)
	response, err := httpClient.Do(directRequest)
	if err != nil {
		redactDirectTransferURLErrors(err)
		return response, err
	}

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		directErr := &DirectTransferResponseError{
			StatusCode: response.StatusCode,
			RetryAfter: directTransferRetryAfter(response.Header.Get("Retry-After")),
		}
		lib.CloseBody(response)
		return response, directErr
	}
	return response, nil
}

func logDirectTransferRequest(logger lib.Logger, request *http.Request) {
	if logger == nil || request == nil || request.URL == nil {
		return
	}

	logURL := redactedDirectTransferURL(*request.URL)
	if leveledLogger, ok := logger.(retryablehttp.LeveledLogger); ok {
		leveledLogger.Debug("performing request", "method", request.Method, "url", logURL)
		return
	}
	logger.Printf("[DEBUG] %s %s", request.Method, logURL)
}

// redactedDirectTransferURL returns directURL without the query, which
// carries the transfer's token, and without user information.
func redactedDirectTransferURL(directURL url.URL) string {
	directURL.User = nil
	directURL.RawQuery = ""
	directURL.ForceQuery = false
	directURL.Fragment = ""
	return directURL.String()
}

// redactDirectTransferURLErrors removes the token from the URL of every
// *url.Error in err's chain, so the error can be logged. The SDK's client
// reports one, and a WrapTransport round tripper that sends through its own
// http.Client reports another inside it.
func redactDirectTransferURLErrors(err error) {
	for ; err != nil; err = errors.Unwrap(err) {
		urlErr, ok := err.(*url.Error)
		if !ok {
			continue
		}
		if parsed, parseErr := url.Parse(urlErr.URL); parseErr == nil {
			urlErr.URL = redactedDirectTransferURL(*parsed)
		} else {
			urlErr.URL = ""
		}
	}
}

func directTransferRetryAfter(value string) time.Duration {
	if seconds, err := strconv.ParseInt(value, 10, 64); err == nil && seconds >= 0 {
		return time.Duration(seconds) * time.Second
	}
	if retryAt, err := http.ParseTime(value); err == nil {
		return max(time.Until(retryAt), 0)
	}
	return 0
}

func directTransferHTTPClient(ctx context.Context, config Config, info DirectConnectionInfo) (string, *http.Client, error) {
	return directTransferHTTPClientWithCache(ctx, config, info, nil, nil)
}

// directTransferHTTPClientWithCache builds the direct client for info, reusing
// the cache carried by ctx, or else cache when it is not nil. A client whose
// transport wrapTransport wraps is built only for cache, never shared through
// ctx, because other users of a context cache may wrap it differently.
func directTransferHTTPClientWithCache(ctx context.Context, config Config, info DirectConnectionInfo, cache *directTransferClientCache, wrapTransport func(http.RoundTripper) http.RoundTripper) (string, *http.Client, error) {
	if config.DisableDirectTransfers || !DirectConnectionInfoPresent(info) {
		return "", nil, ErrDirectTransferUnavailable
	}

	directURL, err := directTransferDirectURL(info)
	if err != nil {
		return "", nil, err
	}

	protocol, candidates, err := directTransferCandidateAddresses(config, info)
	if err != nil {
		return "", nil, err
	}
	build := func() (*http.Client, error) {
		roots, err := directTransferCertPool(info.CaPem)
		if err != nil {
			return nil, err
		}

		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    roots,
			ServerName: info.ServerName,
		}

		var transport http.RoundTripper
		switch protocol {
		case "tcp":
			transport = &http.Transport{
				DialTLSContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
					return directTransferDialTLSCandidates(ctx, network, candidates, tlsConfig)
				},
				ResponseHeaderTimeout: directTransferResponseHeaderTimeout,
				IdleConnTimeout:       30 * time.Second,
				// The Agent direct listener is an HTTP/1.1 endpoint.
				TLSNextProto: map[string]func(string, *tls.Conn) http.RoundTripper{},
			}
		case "quic":
			tlsConfig.MinVersion = tls.VersionTLS13
			transport = &directTransferHTTP3Transport{Transport: &http3.Transport{
				TLSClientConfig: tlsConfig,
				Dial: func(ctx context.Context, _ string, tlsConfig *tls.Config, quicConfig *quic.Config) (*quic.Conn, error) {
					return directTransferDialQUICCandidates(ctx, candidates, tlsConfig, quicConfig)
				},
			}}
		default:
			return nil, ErrDirectTransferUnavailable
		}
		if wrapTransport != nil {
			wrapped := wrapTransport(transport)
			if wrapped == nil {
				return nil, fmt.Errorf("%w: WrapTransport returned no round tripper", ErrDirectTransferUnavailable)
			}
			transport = directTransferWrappedTransport{RoundTripper: wrapped, sdk: transport}
		}

		return &http.Client{
			Transport: transport,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}, nil
	}

	if contextCache, ok := ctx.Value(directTransferClientCacheContextKey{}).(*directTransferClientCache); ok && wrapTransport == nil {
		cache = contextCache
	}
	if cache != nil {
		client, err := cache.client(directTransferClientCacheKey{
			serverName: info.ServerName,
			caPEM:      info.CaPem,
			protocol:   protocol,
			candidates: strings.Join(candidates, ","),
		}, build)
		return directURL, client, err
	}

	client, err := build()
	return directURL, client, err
}

func directTransferDirectURL(info DirectConnectionInfo) (string, error) {
	if !DirectConnectionInfoPresent(info) {
		return "", ErrDirectTransferUnavailable
	}

	transferURL, err := url.Parse(info.DirectUri)
	if err != nil {
		// The parse error repeats direct_uri, token included.
		return "", fmt.Errorf("%w: invalid direct_uri", ErrDirectTransferUnavailable)
	}
	if transferURL.Scheme != "https" {
		return "", fmt.Errorf("%w: direct_uri must use https scheme", ErrDirectTransferUnavailable)
	}
	if !strings.EqualFold(transferURL.Hostname(), info.ServerName) {
		return "", fmt.Errorf("%w: direct_uri host does not match server_name", ErrDirectTransferUnavailable)
	}
	transferURL.User = nil
	return transferURL.String(), nil
}

func directTransferCertPool(caPEM string) (*x509.CertPool, error) {
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM([]byte(caPEM)) {
		return nil, fmt.Errorf("%w: invalid ca_pem", ErrDirectTransferUnavailable)
	}
	return roots, nil
}

func directTransferCandidateAddresses(config Config, info DirectConnectionInfo) (string, []string, error) {
	allowPrivateCandidates := config.Environment != Production
	protocol := ""
	candidates := make([]string, 0, len(info.Addresses))
	for _, rawAddress := range info.Addresses {
		candidateProtocol, candidate, ip, err := directTransferCandidateAddress(rawAddress)
		if err != nil || protocol != "" && candidateProtocol != protocol {
			return "", nil, ErrDirectTransferUnavailable
		}
		protocol = candidateProtocol
		if directTransferCandidateIP(ip, allowPrivateCandidates) {
			candidates = append(candidates, candidate)
		}
	}
	if protocol == "" || len(candidates) == 0 {
		return "", nil, ErrDirectTransferUnavailable
	}
	return protocol, candidates, nil
}

func directTransferCandidateAddress(rawAddress string) (string, string, netip.Addr, error) {
	locator, err := url.Parse(rawAddress)
	if err != nil || locator.Scheme != "tcp" && locator.Scheme != "quic" || rawAddress != locator.Scheme+"://"+locator.Host {
		return "", "", netip.Addr{}, ErrDirectTransferUnavailable
	}

	addrPort, err := netip.ParseAddrPort(locator.Host)
	if err != nil || !addrPort.Addr().Is4() || addrPort.Port() == 0 {
		return "", "", netip.Addr{}, ErrDirectTransferUnavailable
	}

	return locator.Scheme, addrPort.String(), addrPort.Addr(), nil
}

func directTransferCandidateIP(ip netip.Addr, allowPrivate bool) bool {
	if !ip.IsValid() {
		return false
	}
	if allowPrivate && (ip.IsLoopback() || ip.IsPrivate()) {
		return true
	}
	if !ip.IsGlobalUnicast() || ip.IsPrivate() {
		return false
	}
	for _, prefix := range directTransferBlockedIPPrefixes {
		if prefix.Contains(ip) {
			return false
		}
	}
	return true
}

func directTransferDialTLSCandidates(ctx context.Context, network string, candidates []string, tlsConfig *tls.Config) (net.Conn, error) {
	if len(candidates) == 0 {
		return nil, ErrDirectTransferUnavailable
	}

	connectCtx, cancel := context.WithTimeout(ctx, directTransferConnectTimeout)
	defer cancel()

	var errs []error
	dialer := &net.Dialer{}
	for _, candidate := range candidates {
		if ctxErr := connectCtx.Err(); ctxErr != nil {
			errs = append(errs, ctxErr)
			break
		}

		candidateCtx, cancelCandidate := context.WithTimeout(connectCtx, directTransferPerCandidateTimeout)
		conn, err := dialer.DialContext(candidateCtx, network, candidate)
		if err == nil {
			tlsConn := tls.Client(conn, tlsConfig.Clone())
			err = tlsConn.HandshakeContext(candidateCtx)
			if err == nil {
				cancelCandidate()
				return tlsConn, nil
			}
			_ = conn.Close()
		}
		cancelCandidate()
		errs = append(errs, fmt.Errorf("%s: %w", candidate, err))
	}

	if len(errs) == 0 {
		return nil, ErrDirectTransferUnavailable
	}
	return nil, errors.Join(errs...)
}

func directTransferDialQUICCandidates(ctx context.Context, candidates []string, tlsConfig *tls.Config, quicConfig *quic.Config) (*quic.Conn, error) {
	if len(candidates) == 0 {
		return nil, ErrDirectTransferUnavailable
	}

	connectCtx, cancel := context.WithTimeout(ctx, directTransferConnectTimeout)
	defer cancel()

	var errs []error
	for _, candidate := range candidates {
		if ctxErr := connectCtx.Err(); ctxErr != nil {
			errs = append(errs, ctxErr)
			break
		}

		candidateCtx, cancelCandidate := context.WithTimeout(connectCtx, directTransferPerCandidateTimeout)
		conn, err := quic.DialAddr(candidateCtx, candidate, tlsConfig.Clone(), quicConfig.Clone())
		cancelCandidate()
		if err == nil {
			return conn, nil
		}
		errs = append(errs, fmt.Errorf("%s: %w", candidate, err))
	}

	return nil, errors.Join(errs...)
}

type directTransferHTTP3Transport struct {
	*http3.Transport
}

func (t *directTransferHTTP3Transport) RoundTrip(request *http.Request) (*http.Response, error) {
	requestContext, cancel := context.WithCancel(request.Context())
	var timerMu sync.Mutex
	var timer *time.Timer
	roundTripDone := false
	trace := &httptrace.ClientTrace{WroteRequest: func(httptrace.WroteRequestInfo) {
		timerMu.Lock()
		if timer == nil && !roundTripDone {
			timer = time.AfterFunc(directTransferResponseHeaderTimeout, cancel)
		}
		timerMu.Unlock()
	}}
	response, err := t.Transport.RoundTrip(request.Clone(httptrace.WithClientTrace(requestContext, trace)))
	timerMu.Lock()
	roundTripDone = true
	if timer != nil {
		timer.Stop()
	}
	timerMu.Unlock()
	if err != nil {
		cancel()
		return response, err
	}
	response.Body = &directTransferCancelBody{ReadCloser: response.Body, cancel: cancel}
	return response, nil
}

type directTransferCancelBody struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (b *directTransferCancelBody) Close() error {
	err := b.ReadCloser.Close()
	b.cancel()
	return err
}

// directTransferWrappedTransport sends through a caller's wrapper while keeping
// the SDK transport's idle connections closable with the client.
type directTransferWrappedTransport struct {
	http.RoundTripper
	sdk http.RoundTripper
}

func (t directTransferWrappedTransport) CloseIdleConnections() {
	for _, transport := range []http.RoundTripper{t.RoundTripper, t.sdk} {
		if closer, ok := transport.(interface{ CloseIdleConnections() }); ok {
			closer.CloseIdleConnections()
		}
	}
}

func cloneRequestForDirectTransfer(request *http.Request, rawURL string) (*http.Request, error) {
	directURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	cloned := request.Clone(request.Context())
	cloned.URL = directURL
	cloned.Host = directURL.Host
	cloned.Header = *DirectTransferRequestHeaders(&request.Header)
	return cloned, nil
}

// DirectTransferRequestHeaders copies headers and removes Files.com credential
// and workspace headers, Authorization, Proxy-Authorization, and Cookie.
// A nil input returns an empty header map.
func DirectTransferRequestHeaders(headers *http.Header) *http.Header {
	cloned := http.Header{}
	if headers != nil {
		cloned = headers.Clone()
	}
	clearAuthHeaders(&cloned)
	cloned.Del("Authorization")
	cloned.Del("Proxy-Authorization")
	cloned.Del("Cookie")
	return &cloned
}
