package lib

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/Files-com/files-sdk-go/v3/downloadurl"
	"github.com/hashicorp/go-retryablehttp"
)

var DefaultClient *http.Client

func init() {
	DefaultClient = defaultPooledClient()
}

func DefaultRetryableHttp(logger Logger, client ...*http.Client) *retryablehttp.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = logger
	if len(client) == 1 {
		retryClient.HTTPClient = client[0]
	} else {
		retryClient.HTTPClient = DefaultClient
	}
	retryClient.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		if resp != nil {
			if sleep, ok := retryAfterDuration(resp.Header.Get("Retry-After"), time.Now()); ok {
				return sleep
			}
		}

		// Current uses do not improve with adding backoff it will only make things worse, things will become expired.
		return time.Duration(0)
	}
	retryClient.ErrorHandler = retryablehttp.PassthroughErrorHandler
	retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		ok, err := retryablehttp.DefaultRetryPolicy(ctx, resp, err)

		if ok && resp != nil {
			// Don't waste time retrying an expired URL
			if url, err := downloadurl.New(resp.Request.URL.String()); err == nil && !url.IsZero() {
				if _, valid := url.Valid(time.Millisecond * 100); !valid {
					return false, err
				}
			}
		}
		return ok, err
	}
	retryClient.RetryMax = 3
	return retryClient
}

func retryAfterDuration(value string, now time.Time) (time.Duration, bool) {
	if value == "" {
		return 0, false
	}

	if seconds, err := strconv.ParseInt(value, 10, 64); err == nil {
		if seconds < 0 {
			return 0, false
		}
		return time.Second * time.Duration(seconds), true
	}

	retryAt, err := http.ParseTime(value)
	if err != nil {
		return 0, false
	}
	if retryAt.Before(now) {
		return 0, true
	}
	return retryAt.Sub(now), true
}

func defaultPooledClient() *http.Client {
	return &http.Client{
		Transport: DefaultPooledTransport(),
		// Don't use 'Timeout' since it applies to the entire request/response.
	}
}

func DefaultPooledTransport() *Transport {
	transport := &Transport{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			MaxIdleConns:          100,
			ResponseHeaderTimeout: 60 * time.Second,
			IdleConnTimeout:       30 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   75,
			MaxConnsPerHost:       75,
		},
		connections: make(map[string]*int32),
		Dialer: &net.Dialer{
			Timeout: 30 * time.Second,
		},
	}

	transport.Transport.DialContext = transport.DialContext

	return transport
}
