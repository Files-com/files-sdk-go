package files_sdk

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDirectTransferSessionOwnsItsConnections(t *testing.T) {
	connections := atomic.Int32{}
	closed := atomic.Int32{}
	info := directTransferSessionTestInfo(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "ok")
	}), func(_ net.Conn, state http.ConnState) {
		switch state {
		case http.StateNew:
			connections.Add(1)
		case http.StateClosed:
			closed.Add(1)
		}
	})
	session := NewDirectTransferSession(DirectTransferSessionOptions{})

	for range 3 {
		request, err := http.NewRequest(http.MethodGet, "https://proxy.example/downloads", nil)
		require.NoError(t, err)
		response, err := session.Send(Config{Environment: Development}, info, request)
		require.NoError(t, err)
		_, err = io.Copy(io.Discard, response.Body)
		require.NoError(t, err)
		require.NoError(t, response.Body.Close())
	}

	require.Equal(t, int32(1), connections.Load())
	session.Close()
	require.Eventually(t, func() bool { return closed.Load() == connections.Load() }, time.Second, 10*time.Millisecond)
	require.Equal(t, DirectTransferSessionStats{Attempts: 3, Successes: 3}, session.Stats())
}

func TestDirectTransferSessionFallsBackAfterBodyFailureButNotCancellation(t *testing.T) {
	requests := atomic.Int32{}
	firstRequestStarted := make(chan struct{})
	info := directTransferSessionTestInfo(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch requests.Add(1) {
		case 1:
			// Hold the response until the caller cancels before its headers.
			close(firstRequestStarted)
			<-r.Context().Done()
		case 2:
			// Send the headers and a first chunk, then hold the body until the
			// caller cancels.
			_, _ = io.WriteString(w, "ok")
			w.(http.Flusher).Flush()
			<-r.Context().Done()
		default:
			w.Header().Set("Content-Length", "10")
			_, _ = io.WriteString(w, "short")
		}
	}), nil)
	session := NewDirectTransferSession(DirectTransferSessionOptions{})
	t.Cleanup(session.Close)
	config := Config{Environment: Development}
	newRequest := func(ctx context.Context) *http.Request {
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://proxy.example/downloads", nil)
		require.NoError(t, err)
		return request
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-firstRequestStarted
		cancel()
	}()
	_, err := session.Send(config, info, newRequest(ctx))
	require.ErrorIs(t, err, context.Canceled)
	require.Empty(t, session.Stats().DisabledReason, "canceling before the response does not disable direct transfer")

	ctx, cancel = context.WithCancel(context.Background())
	response, err := session.Send(config, info, newRequest(ctx))
	require.NoError(t, err)
	_, err = io.ReadFull(response.Body, make([]byte, 2))
	require.NoError(t, err)
	cancel()
	_, err = io.ReadAll(response.Body)
	require.ErrorIs(t, err, context.Canceled)
	_ = response.Body.Close()
	require.Empty(t, session.Stats().DisabledReason, "canceling while reading the body does not disable direct transfer")

	response, err = session.Send(config, info, newRequest(context.Background()))
	require.NoError(t, err)
	_, err = io.ReadAll(response.Body)
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
	require.NoError(t, response.Body.Close())
	require.Equal(t, "direct_body_read_failed", session.Stats().DisabledReason)

	_, err = session.Send(config, info, newRequest(context.Background()))
	require.ErrorIs(t, err, ErrDirectTransferUnavailable)
	require.Equal(t, int32(3), requests.Load(), "a disabled session does not contact the direct endpoint")
	require.Equal(t, DirectTransferSessionStats{Attempts: 3, Successes: 2, DisabledReason: "direct_body_read_failed"}, session.Stats())
}

func TestDirectTransferSessionCancelCauseReportsBodyFailure(t *testing.T) {
	info := directTransferSessionTestInfo(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "ok")
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	}), nil)
	session := NewDirectTransferSession(DirectTransferSessionOptions{})
	t.Cleanup(session.Close)
	ctx, cancel := context.WithCancelCause(context.Background())
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://proxy.example/downloads", nil)
	require.NoError(t, err)

	response, err := session.Send(Config{Environment: Development}, info, request)
	require.NoError(t, err)
	_, err = io.ReadFull(response.Body, make([]byte, 2))
	require.NoError(t, err)
	// A caller's own stall detection cancels with its reason.
	cancel(errors.New("no data received within the inactivity period"))
	_, err = io.ReadAll(response.Body)
	require.Error(t, err)
	_ = response.Body.Close()

	require.Equal(t, "direct_body_read_failed", session.Stats().DisabledReason)
}

func TestDirectTransferSessionUsesContextFromRequestOptions(t *testing.T) {
	info := directTransferSessionTestInfo(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "ok")
	}), nil)
	session := NewDirectTransferSession(DirectTransferSessionOptions{})
	t.Cleanup(session.Close)
	config := Config{Environment: Development}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()

	request, err := http.NewRequestWithContext(canceled, http.MethodGet, "https://proxy.example/downloads", nil)
	require.NoError(t, err)
	response, err := session.Send(config, info, request, WithContext(context.Background()))
	require.NoError(t, err, "WithContext replaces the request's canceled context")
	require.NoError(t, response.Body.Close())

	request, err = http.NewRequest(http.MethodGet, "https://proxy.example/downloads", nil)
	require.NoError(t, err)
	_, err = session.Send(config, info, request, WithContext(canceled))
	require.ErrorIs(t, err, context.Canceled)
	require.Equal(t, DirectTransferSessionStats{Attempts: 1, Successes: 1}, session.Stats())
}

func TestDirectTransferSessionWrapTransportSeesSanitizedDirectRequests(t *testing.T) {
	targetRequests := atomic.Int32{}
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		targetRequests.Add(1)
	}))
	t.Cleanup(target.Close)
	connections := atomic.Int32{}
	closed := atomic.Int32{}
	directRequests := atomic.Int32{}
	info := directTransferSessionTestInfo(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if directRequests.Add(1) == 2 {
			// No body, so the connection stays open for Close to release.
			w.Header().Set("Location", target.URL)
			w.WriteHeader(http.StatusFound)
			return
		}
		_, _ = io.WriteString(w, "ok")
	}), func(_ net.Conn, state http.ConnState) {
		switch state {
		case http.StateNew:
			connections.Add(1)
		case http.StateClosed:
			closed.Add(1)
		}
	})
	info.DirectUri += "?jwt=direct-token"
	var wraps atomic.Int32
	var seen []*http.Request
	var seenMu sync.Mutex
	session := NewDirectTransferSession(DirectTransferSessionOptions{
		WrapTransport: func(direct http.RoundTripper) http.RoundTripper {
			wraps.Add(1)
			return directTransferRoundTripFunc(func(request *http.Request) (*http.Response, error) {
				seenMu.Lock()
				seen = append(seen, request.Clone(context.Background()))
				seenMu.Unlock()
				return direct.RoundTrip(request)
			})
		},
	})
	config := Config{Environment: Development}
	newRequest := func() *http.Request {
		request, err := http.NewRequest(http.MethodGet, "https://proxy.example/downloads?proxy=secret", nil)
		require.NoError(t, err)
		request.Header.Set(apiKeyHeader, "api-key")
		request.Header.Set(sessionIdHeader, "session-id")
		request.Header.Set("Authorization", "Bearer token")
		request.Header.Set("Cookie", "session=secret")
		request.Header.Set("Range", "bytes=0-1")
		return request
	}

	response, err := session.Send(config, info, newRequest())
	require.NoError(t, err)
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Equal(t, "ok", string(body))

	_, err = session.Send(config, info, newRequest())
	var responseErr *DirectTransferResponseError
	require.ErrorAs(t, err, &responseErr)
	require.Equal(t, http.StatusFound, responseErr.StatusCode)
	require.Zero(t, targetRequests.Load(), "the wrapped client does not follow redirects")

	require.Equal(t, int32(1), wraps.Load())
	require.Len(t, seen, 2)
	for _, request := range seen {
		require.Equal(t, info.DirectUri, request.URL.String())
		require.Empty(t, request.Header.Get(apiKeyHeader))
		require.Empty(t, request.Header.Get(sessionIdHeader))
		require.Empty(t, request.Header.Get("Authorization"))
		require.Empty(t, request.Header.Get("Cookie"))
		require.Equal(t, "bytes=0-1", request.Header.Get("Range"))
	}

	// The wrapper does not forward CloseIdleConnections; Close still closes
	// the SDK transport's idle connection.
	require.Equal(t, int32(1), connections.Load())
	require.Zero(t, closed.Load())
	session.Close()
	require.Eventually(t, func() bool { return closed.Load() == 1 }, time.Second, 10*time.Millisecond)
}

func TestDirectTransferSessionFailureDiagnosticsOmitTransferToken(t *testing.T) {
	for name, wrapTransport := range map[string]func(http.RoundTripper) http.RoundTripper{
		"SDK transport": nil,
		// A caller's own http.Client reports its own *url.Error, which the
		// SDK's client then wraps in another.
		"caller-owned http.Client": func(direct http.RoundTripper) http.RoundTripper {
			callerClient := &http.Client{Transport: direct, CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			}}
			return directTransferRoundTripFunc(callerClient.Do)
		},
	} {
		t.Run(name, func(t *testing.T) {
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err)
			unreachable := listener.Addr().String()
			require.NoError(t, listener.Close())
			caPEM, _ := directTransferTestCertificate(t, directTransferTestServerName)
			info := DirectConnectionInfo{
				Version:    1,
				ServerName: directTransferTestServerName,
				Addresses:  []string{directTransferTCPLocator(unreachable)},
				DirectUri:  "https://" + directTransferTestServerName + "/downloads?jwt=direct-token",
				CaPem:      string(caPEM),
			}
			var logs bytes.Buffer
			logger := log.New(&logs, "", 0)
			session := NewDirectTransferSession(DirectTransferSessionOptions{
				OnDisable: func(reason string, err error) {
					logger.Printf("direct transfer disabled reason=%s error=%v", reason, err)
				},
				WrapTransport: wrapTransport,
			})
			t.Cleanup(session.Close)
			request, err := http.NewRequest(http.MethodGet, "https://proxy.example/downloads", nil)
			require.NoError(t, err)

			_, err = session.Send(Config{Environment: Development, Logger: logger}, info, request)

			require.ErrorIs(t, err, ErrDirectTransferUnavailable)
			var dialErr *net.OpError
			require.ErrorAs(t, err, &dialErr, "the underlying failure stays inspectable")
			require.Contains(t, logs.String(), "[DEBUG] GET https://"+directTransferTestServerName+"/downloads")
			require.Contains(t, logs.String(), "direct transfer disabled reason=direct_request_failed")
			require.Contains(t, logs.String(), unreachable)
			require.NotContains(t, logs.String(), "direct-token")
			require.NotContains(t, err.Error(), "direct-token")
		})
	}
}

func directTransferSessionTestInfo(t *testing.T, handler http.Handler, connState func(net.Conn, http.ConnState)) DirectConnectionInfo {
	t.Helper()
	caPEM, cert := directTransferTestCertificate(t, directTransferTestServerName)
	server := httptest.NewUnstartedServer(handler)
	server.Config.ConnState = connState
	server.TLS = &tls.Config{Certificates: []tls.Certificate{cert}}
	server.StartTLS()
	t.Cleanup(server.Close)

	_, port, err := net.SplitHostPort(server.Listener.Addr().String())
	require.NoError(t, err)
	return DirectConnectionInfo{
		Version:    1,
		ServerName: directTransferTestServerName,
		Addresses:  []string{directTransferTCPLocator(server.Listener.Addr().String())},
		DirectUri:  "https://" + net.JoinHostPort(directTransferTestServerName, port) + "/downloads",
		CaPem:      string(caPEM),
	}
}
