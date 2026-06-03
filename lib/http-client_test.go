package lib

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDefaultRetryableHttpBackoffHonorsRetryAfterForRetryableStatuses(t *testing.T) {
	client := DefaultRetryableHttp(nil)
	for _, status := range []int{
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
	} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			resp := &http.Response{
				StatusCode: status,
				Header:     http.Header{"Retry-After": []string{"7"}},
			}

			if got := client.Backoff(0, 0, 1, resp); got != 7*time.Second {
				t.Fatalf("Backoff() = %v, want 7s", got)
			}
		})
	}
}

func TestRetryAfterDurationParsesHTTPDate(t *testing.T) {
	now := time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC)
	retryAt := now.Add(9 * time.Second).Format(http.TimeFormat)

	got, ok := retryAfterDuration(retryAt, now)
	if !ok {
		t.Fatal("retryAfterDuration() ok = false, want true")
	}
	if got != 9*time.Second {
		t.Fatalf("retryAfterDuration() = %v, want 9s", got)
	}
}

func TestUploadRetryableHttpRetriesS3XMLRetryableCodes(t *testing.T) {
	tests := []string{
		"RequestTimeout",
		"InternalError",
		"ServiceUnavailable",
		"SlowDown",
		"RequestLimitExceeded",
		"DatabaseTimeout",
	}

	for _, code := range tests {
		t.Run(code, func(t *testing.T) {
			client := UploadRetryableHttp(DefaultRetryableHttp(nil), func() bool { return true })
			ok, err := client.CheckRetry(context.Background(), s3XMLResponse(http.StatusBadRequest, code), nil)

			require.NoError(t, err)
			require.True(t, ok)
		})
	}
}

func TestUploadRetryableHttpDoesNotRetryS3XMLTerminalCodes(t *testing.T) {
	tests := []struct {
		code   string
		status int
	}{
		{"AccessDenied", http.StatusForbidden},
		{"NoSuchKey", http.StatusNotFound},
		{"InvalidRequest", http.StatusBadRequest},
		{"SignatureDoesNotMatch", http.StatusForbidden},
	}

	for _, test := range tests {
		t.Run(test.code, func(t *testing.T) {
			client := UploadRetryableHttp(DefaultRetryableHttp(nil), func() bool { return true })
			ok, err := client.CheckRetry(context.Background(), s3XMLResponse(test.status, test.code), nil)

			require.NoError(t, err)
			require.False(t, ok)
		})
	}
}

func TestUploadRetryableHttpDoesNotRetryNonXMLBody(t *testing.T) {
	client := UploadRetryableHttp(DefaultRetryableHttp(nil), func() bool { return true })
	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"error":"RequestTimeout"}`)),
	}

	ok, err := client.CheckRetry(context.Background(), resp, nil)

	require.NoError(t, err)
	require.False(t, ok)
}

func TestUploadRetryableHttpLeavesBodyReadableAfterPeek(t *testing.T) {
	body := `<?xml version="1.0" encoding="UTF-8"?><Error><Code>RequestTimeout</Code><Message>provider message</Message></Error>`
	client := UploadRetryableHttp(DefaultRetryableHttp(nil), func() bool { return true })
	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Header:     http.Header{"Content-Type": []string{"application/xml"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}

	ok, err := client.CheckRetry(context.Background(), resp, nil)
	readBody, readErr := io.ReadAll(resp.Body)

	require.NoError(t, err)
	require.True(t, ok)
	require.NoError(t, readErr)
	require.Equal(t, body, string(readBody))
}

func TestUploadRetryableHttpClosesOriginalBodyAfterPeek(t *testing.T) {
	body := `<?xml version="1.0" encoding="UTF-8"?><Error><Code>RequestTimeout</Code><Message>provider message</Message></Error>`
	bodyCloser := &closeTrackingReadCloser{Reader: strings.NewReader(body)}
	client := UploadRetryableHttp(DefaultRetryableHttp(nil), func() bool { return true })
	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Header:     http.Header{"Content-Type": []string{"application/xml"}},
		Body:       bodyCloser,
	}

	ok, err := client.CheckRetry(context.Background(), resp, nil)
	closeErr := resp.Body.Close()

	require.NoError(t, err)
	require.True(t, ok)
	require.NoError(t, closeErr)
	require.True(t, bodyCloser.closed)
}

func TestUploadRetryableHttpDoesNotRetryWhenBodyCannotReplay(t *testing.T) {
	body := `<?xml version="1.0" encoding="UTF-8"?><Error><Code>RequestTimeout</Code><Message>provider message</Message></Error>`
	client := UploadRetryableHttp(DefaultRetryableHttp(nil), func() bool { return false })
	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Header:     http.Header{"Content-Type": []string{"application/xml"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}

	ok, err := client.CheckRetry(context.Background(), resp, nil)
	readBody, readErr := io.ReadAll(resp.Body)

	require.NoError(t, err)
	require.False(t, ok)
	require.NoError(t, readErr)
	require.Equal(t, body, string(readBody))
}

func TestUploadRetryableHttpDoesNotRetryExpiredSignedURLForS3XMLRetryableCode(t *testing.T) {
	client := UploadRetryableHttp(DefaultRetryableHttp(nil), func() bool { return true })
	resp := s3XMLResponse(http.StatusBadRequest, "RequestTimeout")
	req, err := http.NewRequest(http.MethodPut, expiredS3URL(), nil)
	require.NoError(t, err)
	resp.Request = req

	ok, err := client.CheckRetry(context.Background(), resp, nil)

	require.NoError(t, err)
	require.False(t, ok)
}

func TestUploadRetryableHttpBackoffHonorsRetryAfter(t *testing.T) {
	client := UploadRetryableHttp(DefaultRetryableHttp(nil), func() bool { return true })
	resp := &http.Response{
		StatusCode: http.StatusServiceUnavailable,
		Header:     http.Header{"Retry-After": []string{"12"}},
	}

	require.Equal(t, 12*time.Second, client.Backoff(0, 0, 0, resp))
}

func TestUploadRetryableHttpBackoffJittersRateLimitS3XML(t *testing.T) {
	client := UploadRetryableHttp(DefaultRetryableHttp(nil), func() bool { return true })
	resp := s3XMLResponse(http.StatusBadRequest, "SlowDown")

	ok, err := client.CheckRetry(context.Background(), resp, nil)

	require.NoError(t, err)
	require.True(t, ok)
	assertUploadRetryBackoffBetween(t, client.Backoff(0, 0, 0, resp), 500*time.Millisecond, time.Second)
	assertUploadRetryBackoffBetween(t, client.Backoff(0, 0, 2, resp), time.Second, 2*time.Second)
	assertUploadRetryBackoffBetween(t, client.Backoff(0, 0, 3, resp), 2*time.Second, 4*time.Second)
}

func TestUploadRetryableHttpBackoffJittersRateLimitS3XMLTooManyRequests(t *testing.T) {
	client := UploadRetryableHttp(DefaultRetryableHttp(nil), func() bool { return true })
	resp := s3XMLResponse(http.StatusTooManyRequests, "SlowDown")

	ok, err := client.CheckRetry(context.Background(), resp, nil)

	require.NoError(t, err)
	require.True(t, ok)
	assertUploadRetryBackoffBetween(t, client.Backoff(0, 0, 0, resp), 500*time.Millisecond, time.Second)
}

func TestUploadRetryableHttpBackoffDoesNotJitterRequestTimeoutS3XML(t *testing.T) {
	client := UploadRetryableHttp(DefaultRetryableHttp(nil), func() bool { return true })
	resp := s3XMLResponse(http.StatusBadRequest, "RequestTimeout")

	ok, err := client.CheckRetry(context.Background(), resp, nil)

	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, time.Duration(0), client.Backoff(0, 0, 0, resp))
}

func TestCloneHTTPClientWithMaxConnsPerHostRaisesTransportWithoutMutatingOriginal(t *testing.T) {
	transport := DefaultPooledTransport()
	client := &http.Client{Transport: transport}

	cloned, ok := CloneHTTPClientWithMaxConnsPerHost(client, 1024)

	require.True(t, ok)
	require.NotSame(t, client, cloned)
	require.NotSame(t, transport, cloned.Transport)
	require.Equal(t, 75, transport.MaxConnsPerHost)
	clonedTransport, ok := cloned.Transport.(*Transport)
	require.True(t, ok)
	require.Equal(t, 1024, clonedTransport.MaxConnsPerHost)
	require.Equal(t, 1024, clonedTransport.MaxIdleConns)
	require.Equal(t, 1024, clonedTransport.MaxIdleConnsPerHost)
}

func TestCloneHTTPClientWithMaxConnsPerHostSharesConnectionStats(t *testing.T) {
	transport := DefaultPooledTransport()
	client := &http.Client{Transport: transport}

	cloned, ok := CloneHTTPClientWithMaxConnsPerHost(client, 1024)

	require.True(t, ok)
	clonedTransport, ok := cloned.Transport.(*Transport)
	require.True(t, ok)
	require.Same(t, transport.connectionStats(), clonedTransport.connectionStats())

	counter := int32(3)
	stats := clonedTransport.connectionStats()
	stats.mu.Lock()
	stats.connections["uploads.files.com:443"] = &counter
	stats.mu.Unlock()

	connectionStats, ok := GetConnectionStatsFromClient(client)
	require.True(t, ok)
	require.Equal(t, map[string]int{"uploads.files.com:443": 3}, connectionStats)
}

func TestCloneHTTPClientWithExactMaxConnsPerHostLowersTransportWithoutMutatingOriginal(t *testing.T) {
	transport := DefaultPooledTransport()
	transport.MaxConnsPerHost = 1024
	transport.MaxIdleConns = 1024
	transport.MaxIdleConnsPerHost = 1024
	client := &http.Client{Transport: transport}

	cloned, ok := CloneHTTPClientWithExactMaxConnsPerHost(client, 150)

	require.True(t, ok)
	require.NotSame(t, client, cloned)
	require.NotSame(t, transport, cloned.Transport)
	require.Equal(t, 1024, transport.MaxConnsPerHost)
	clonedTransport, ok := cloned.Transport.(*Transport)
	require.True(t, ok)
	require.Equal(t, 150, clonedTransport.MaxConnsPerHost)
	require.Equal(t, 1024, clonedTransport.MaxIdleConns)
	require.Equal(t, 1024, clonedTransport.MaxIdleConnsPerHost)
}

func TestCloneHTTPClientWithExactMaxConnsPerHostSharesConnectionStats(t *testing.T) {
	transport := DefaultPooledTransport()
	transport.MaxConnsPerHost = 1024
	client := &http.Client{Transport: transport}

	cloned, ok := CloneHTTPClientWithExactMaxConnsPerHost(client, 150)

	require.True(t, ok)
	clonedTransport, ok := cloned.Transport.(*Transport)
	require.True(t, ok)
	require.Same(t, transport.connectionStats(), clonedTransport.connectionStats())

	counter := int32(2)
	stats := clonedTransport.connectionStats()
	stats.mu.Lock()
	stats.connections["s3.amazonaws.com:443"] = &counter
	stats.mu.Unlock()

	connectionStats, ok := GetConnectionStatsFromClient(client)
	require.True(t, ok)
	require.Equal(t, map[string]int{"s3.amazonaws.com:443": 2}, connectionStats)
}

func TestTransportConnectionStatsLazyInitIsRaceSafe(t *testing.T) {
	transport := &Transport{}
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			require.NotNil(t, transport.connectionStats())
		}()
	}

	close(start)
	wg.Wait()
	require.NotNil(t, transport.connectionStats())
}

func assertUploadRetryBackoffBetween(t *testing.T, delay time.Duration, minDelay time.Duration, maxDelay time.Duration) {
	t.Helper()

	require.GreaterOrEqual(t, delay, minDelay)
	require.Less(t, delay, maxDelay)
}

func s3XMLResponse(status int, code string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/xml"}},
		Body:       io.NopCloser(strings.NewReader(`<?xml version="1.0" encoding="UTF-8"?><Error><Code>` + code + `</Code><Message>provider message</Message></Error>`)),
	}
}

func expiredS3URL() string {
	return "https://example.com/upload?X-Amz-Date=" + time.Now().Add(-time.Hour).UTC().Format("20060102T150405Z") + "&X-Amz-Expires=60"
}

type closeTrackingReadCloser struct {
	*strings.Reader
	closed bool
}

func (c *closeTrackingReadCloser) Close() error {
	c.closed = true
	return nil
}
