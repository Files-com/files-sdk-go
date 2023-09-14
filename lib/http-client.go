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

type NullLogger struct{}

func (n NullLogger) Printf(_ string, _ ...interface{}) {}

type Logger interface {
	Printf(string, ...interface{})
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
		if resp != nil && (resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable) {
			if s, ok := resp.Header["Retry-After"]; ok {
				if sleep, err := strconv.ParseInt(s[0], 10, 64); err == nil {
					return time.Second * time.Duration(sleep)
				}
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

func defaultPooledClient() *http.Client {
	return &http.Client{
		Transport: defaultPooledTransport(),
	}
}

func defaultPooledTransport() *Transport {
	transport := &Transport{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			MaxIdleConns:          100,
			IdleConnTimeout:       30 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   75,
			MaxConnsPerHost:       100,
		},
		Connections: make(map[string]int),
		Dialer: &net.Dialer{
			Timeout: 30 * time.Second,
		},
	}

	transport.Transport.DialContext = transport.DialContext

	return transport
}
