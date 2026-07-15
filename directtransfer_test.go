package files_sdk

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const directTransferTestServerName = "agent-123.agents.files.internal"

func setDirectTransferTestTimeouts(t *testing.T, total, perCandidate time.Duration) {
	t.Helper()
	previousTotal := directTransferConnectTimeout
	previousPerCandidate := directTransferPerCandidateTimeout
	directTransferConnectTimeout = total
	directTransferPerCandidateTimeout = perCandidate
	t.Cleanup(func() {
		directTransferConnectTimeout = previousTotal
		directTransferPerCandidateTimeout = previousPerCandidate
	})
}

func directTransferTestInfo(path string) DirectConnectionInfo {
	return DirectConnectionInfo{
		Version:    1,
		ServerName: directTransferTestServerName,
		Addresses:  []string{directTransferTCPLocator("8.8.8.8:4001")},
		DirectUri:  "https://" + net.JoinHostPort(directTransferTestServerName, "4001") + path,
		CaPem:      "test-ca",
	}
}

func directTransferTCPLocator(address string) string {
	return "tcp://" + address
}

func TestDirectConnectionInfoPresent(t *testing.T) {
	require.False(t, DirectConnectionInfoPresent(DirectConnectionInfo{}))
	require.True(t, DirectConnectionInfoPresent(directTransferTestInfo("/downloads/file.txt?jwt=download-token")))

	info := directTransferTestInfo("/downloads/file.txt?jwt=download-token")
	info.Version = 2
	require.False(t, DirectConnectionInfoPresent(info))
}

func TestDirectTransferUsesDirectURI(t *testing.T) {
	info := directTransferTestInfo("/downloads/file.txt?jwt=download-token")

	directURL, err := directTransferDirectURL(info)
	require.NoError(t, err)
	require.Equal(t, info.DirectUri, directURL)

	info.DirectUri = "http://" + net.JoinHostPort(info.ServerName, "4001") + "/downloads/file.txt?jwt=download-token"
	_, err = directTransferDirectURL(info)
	require.ErrorIs(t, err, ErrDirectTransferUnavailable)

	info = directTransferTestInfo("/downloads/file.txt?jwt=download-token")
	info.DirectUri = "https://other-agent.agents.files.internal:4001/downloads/file.txt?jwt=download-token"
	_, err = directTransferDirectURL(info)
	require.ErrorIs(t, err, ErrDirectTransferUnavailable)
}

func TestDirectTransferOptionsDoNotMutateFallbackRequest(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "https://proxy.example/downloads/file.txt?jwt=proxy-token", nil)
	require.NoError(t, err)

	info := directTransferTestInfo("/downloads/file.txt?jwt=direct-token")
	addTransferHeader := RequestOption(func(req *http.Request) error {
		req.Header.Add("X-Transfer-Attempt", "1")
		return nil
	})

	_, err = WrapDirectTransferOptions(Config{}, info, request, addTransferHeader)
	require.Error(t, err)
	require.Empty(t, request.Header.Values("X-Transfer-Attempt"))

	_, err = BuildRequest(request, addTransferHeader)
	require.NoError(t, err)
	require.Equal(t, []string{"1"}, request.Header.Values("X-Transfer-Attempt"))
}

func TestDirectTransferHTTPClientDialsAddressWithServerName(t *testing.T) {
	caPEM, cert := directTransferTestCertificate(t, directTransferTestServerName)

	var seenHost, seenSNI string
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenHost = r.Host
		if r.TLS != nil {
			seenSNI = r.TLS.ServerName
		}
		require.Equal(t, "/downloads/file.txt", r.URL.Path)
		_, _ = io.WriteString(w, "direct")
	}))
	server.TLS = &tls.Config{Certificates: []tls.Certificate{cert}}
	server.StartTLS()
	t.Cleanup(server.Close)

	_, port, err := net.SplitHostPort(server.Listener.Addr().String())
	require.NoError(t, err)
	directHost := net.JoinHostPort(directTransferTestServerName, port)
	info := DirectConnectionInfo{
		Version:    1,
		ServerName: directTransferTestServerName,
		Addresses:  []string{directTransferTCPLocator(server.Listener.Addr().String())},
		DirectUri:  "https://" + directHost + "/downloads/file.txt?jwt=direct-token",
		CaPem:      string(caPEM),
	}

	directURL, httpClient, err := directTransferHTTPClient(
		context.Background(),
		Config{Environment: Development},
		info,
	)
	require.NoError(t, err)
	require.Equal(t, info.DirectUri, directURL)

	req, err := http.NewRequest(http.MethodGet, directURL, nil)
	require.NoError(t, err)
	res, err := httpClient.Do(req)
	require.NoError(t, err)
	t.Cleanup(func() { _ = res.Body.Close() })

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "direct", string(body))
	require.Equal(t, directHost, seenHost)
	require.Equal(t, directTransferTestServerName, seenSNI)
}

func TestDirectTransferHTTPClientIgnoresEnvironmentProxy(t *testing.T) {
	caPEM, cert := directTransferTestCertificate(t, directTransferTestServerName)

	proxyHit := make(chan struct{}, 1)
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		select {
		case proxyHit <- struct{}{}:
		default:
		}
		http.Error(w, "proxy should not be used", http.StatusBadGateway)
	}))
	t.Cleanup(proxy.Close)
	t.Setenv("HTTP_PROXY", proxy.URL)
	t.Setenv("HTTPS_PROXY", proxy.URL)

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/downloads/file.txt", r.URL.Path)
		_, _ = io.WriteString(w, "direct")
	}))
	server.TLS = &tls.Config{Certificates: []tls.Certificate{cert}}
	server.StartTLS()
	t.Cleanup(server.Close)

	_, port, err := net.SplitHostPort(server.Listener.Addr().String())
	require.NoError(t, err)
	directHost := net.JoinHostPort(directTransferTestServerName, port)
	info := DirectConnectionInfo{
		Version:    1,
		ServerName: directTransferTestServerName,
		Addresses:  []string{directTransferTCPLocator(server.Listener.Addr().String())},
		DirectUri:  "https://" + directHost + "/downloads/file.txt?jwt=direct-token",
		CaPem:      string(caPEM),
	}

	directURL, httpClient, err := directTransferHTTPClient(
		context.Background(),
		Config{Environment: Development},
		info,
	)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, directURL, nil)
	require.NoError(t, err)
	res, err := httpClient.Do(req)
	require.NoError(t, err)
	t.Cleanup(func() { _ = res.Body.Close() })

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "direct", string(body))
	select {
	case <-proxyHit:
		t.Fatal("direct transfer used the environment proxy")
	default:
	}
}

func TestDirectTransferClientPolicyAndCache(t *testing.T) {
	caPEM, _ := directTransferTestCertificate(t, directTransferTestServerName)
	info := directTransferTestInfo("/downloads")
	info.Addresses = []string{directTransferTCPLocator("127.0.0.1:4001")}
	info.CaPem = string(caPEM)
	ctx, closeClients := WithDirectTransferClientCache(context.Background())
	defer closeClients()

	_, first, err := directTransferHTTPClient(ctx, Config{Environment: Development}, info)
	require.NoError(t, err)
	_, second, err := directTransferHTTPClient(ctx, Config{Environment: Development}, info)
	require.NoError(t, err)
	require.Same(t, first, second)

	transport, ok := first.Transport.(*http.Transport)
	require.True(t, ok)
	require.Equal(t, directTransferResponseHeaderTimeout, transport.ResponseHeaderTimeout)
	require.Equal(t, 30*time.Second, transport.IdleConnTimeout)
	require.ErrorIs(t, first.CheckRedirect(&http.Request{}, nil), http.ErrUseLastResponse)
}

func TestWrapDirectTransferOptionsClosesBodyAfterResponseProcessingStarts(t *testing.T) {
	info := directTransferTestInfo("/downloads")
	ctx, closeClients := WithDirectTransferClientCache(context.Background())
	defer closeClients()
	cache := ctx.Value(directTransferClientCacheContextKey{}).(*directTransferClientCache)
	body := &directTransferTrackingBody{}
	cache.key = directTransferClientCacheKey{
		serverName: info.ServerName,
		caPEM:      info.CaPem,
		candidates: "8.8.8.8:4001",
	}
	cache.httpClient = &http.Client{Transport: directTransferRoundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: body, Header: http.Header{}}, nil
	})}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://proxy.example/downloads", nil)
	require.NoError(t, err)

	_, err = WrapDirectTransferOptions(
		Config{Environment: Development},
		info,
		request,
		ResponseBodyOption(func(responseBody io.ReadCloser) error {
			_, readErr := responseBody.Read(make([]byte, 16))
			return readErr
		}),
	)

	require.ErrorIs(t, err, ErrDirectTransferResponseStarted)
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
	require.True(t, body.closed.Load())
}

type directTransferRoundTripFunc func(*http.Request) (*http.Response, error)

func (f directTransferRoundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

type directTransferTrackingBody struct {
	read   atomic.Bool
	closed atomic.Bool
}

func (b *directTransferTrackingBody) Read(p []byte) (int, error) {
	if !b.read.CompareAndSwap(false, true) {
		return 0, io.EOF
	}
	return copy(p, "prefix"), io.ErrUnexpectedEOF
}

func (b *directTransferTrackingBody) Close() error {
	b.closed.Store(true)
	return nil
}

func TestDirectTransferClientReusesAndClosesConnections(t *testing.T) {
	caPEM, cert := directTransferTestCertificate(t, directTransferTestServerName)
	connections := atomic.Int32{}
	closed := atomic.Int32{}
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "ok")
	}))
	server.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		switch state {
		case http.StateNew:
			connections.Add(1)
		case http.StateClosed:
			closed.Add(1)
		}
	}
	server.TLS = &tls.Config{Certificates: []tls.Certificate{cert}}
	server.StartTLS()
	t.Cleanup(server.Close)

	_, port, err := net.SplitHostPort(server.Listener.Addr().String())
	require.NoError(t, err)
	info := DirectConnectionInfo{
		Version:    1,
		ServerName: directTransferTestServerName,
		Addresses:  []string{directTransferTCPLocator(server.Listener.Addr().String())},
		DirectUri:  "https://" + net.JoinHostPort(directTransferTestServerName, port) + "/downloads",
		CaPem:      string(caPEM),
	}
	ctx, closeClients := WithDirectTransferClientCache(context.Background())

	for range 3 {
		directURL, client, clientErr := directTransferHTTPClient(ctx, Config{Environment: Development}, info)
		require.NoError(t, clientErr)
		response, requestErr := client.Get(directURL)
		require.NoError(t, requestErr)
		_, requestErr = io.Copy(io.Discard, response.Body)
		require.NoError(t, requestErr)
		require.NoError(t, response.Body.Close())
	}

	require.Equal(t, int32(1), connections.Load())
	closeClients()
	require.Eventually(t, func() bool { return closed.Load() == connections.Load() }, time.Second, 10*time.Millisecond)
}

func TestDirectTransferResponseHeaderTimeout(t *testing.T) {
	previousTimeout := directTransferResponseHeaderTimeout
	directTransferResponseHeaderTimeout = 100 * time.Millisecond
	t.Cleanup(func() { directTransferResponseHeaderTimeout = previousTimeout })

	caPEM, cert := directTransferTestCertificate(t, directTransferTestServerName)
	requestStarted := make(chan struct{})
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		close(requestStarted)
		<-r.Context().Done()
	}))
	server.TLS = &tls.Config{Certificates: []tls.Certificate{cert}}
	server.StartTLS()
	t.Cleanup(server.Close)

	_, port, err := net.SplitHostPort(server.Listener.Addr().String())
	require.NoError(t, err)
	info := DirectConnectionInfo{
		Version:    1,
		ServerName: directTransferTestServerName,
		Addresses:  []string{directTransferTCPLocator(server.Listener.Addr().String())},
		DirectUri:  "https://" + net.JoinHostPort(directTransferTestServerName, port) + "/downloads",
		CaPem:      string(caPEM),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://proxy.example/downloads", nil)
	require.NoError(t, err)

	startedAt := time.Now()
	_, err = WrapDirectTransferOptions(Config{Environment: Development}, info, request)
	require.Error(t, err)
	require.Less(t, time.Since(startedAt), time.Second)
	select {
	case <-requestStarted:
	default:
		t.Fatal("direct request did not reach the TLS server")
	}
}

func TestDirectTransferCandidateAddressesFiltersProductionPrivateAddresses(t *testing.T) {
	info := directTransferTestInfo("/downloads/file.txt?jwt=download-token")
	info.Addresses = []string{
		"tcp://127.0.0.1:4001",
		"tcp://10.1.2.3:4001",
		"tcp://8.8.8.8:4001",
		"tcp://100.64.0.9:4001",
		"8.8.8.8:4002",
		"quic://8.8.8.8:4001",
		"tcp://8.8.8.8:0",
		"tcp://8.8.8.8:65536",
		"tcp://8.8.8.8:4001/path",
		"not-an-address",
	}

	production, err := directTransferCandidateAddresses(Config{Environment: Production}, info)
	require.NoError(t, err)
	require.Equal(t, []string{"8.8.8.8:4001"}, production)

	development, err := directTransferCandidateAddresses(Config{Environment: Development}, info)
	require.NoError(t, err)
	require.Equal(t, []string{"127.0.0.1:4001", "10.1.2.3:4001", "8.8.8.8:4001"}, development)
}

func TestDirectTransferPerCandidateTimeoutSkipsStalledCandidate(t *testing.T) {
	setDirectTransferTestTimeouts(t, 3*time.Second, 300*time.Millisecond)

	caPEM, cert := directTransferTestCertificate(t, directTransferTestServerName)
	stalledAddress, accepted := newStalledDirectTransferTLSCandidate(t)
	healthy := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "ok")
	}))
	healthy.TLS = &tls.Config{Certificates: []tls.Certificate{cert}}
	healthy.StartTLS()
	t.Cleanup(healthy.Close)

	candidates := []string{stalledAddress, healthy.Listener.Addr().String()}
	info := DirectConnectionInfo{
		Version:    1,
		ServerName: directTransferTestServerName,
		Addresses:  []string{directTransferTCPLocator(stalledAddress), directTransferTCPLocator(healthy.Listener.Addr().String())},
		DirectUri:  "https://" + net.JoinHostPort(directTransferTestServerName, "4001") + "/downloads/file.txt?jwt=direct-token",
		CaPem:      string(caPEM),
	}
	roots, err := directTransferCertPool(info.CaPem)
	require.NoError(t, err)

	startedAt := time.Now()
	conn, err := directTransferDialTLSCandidates(context.Background(), "tcp", candidates, &tls.Config{
		RootCAs:    roots,
		ServerName: info.ServerName,
		MinVersion: tls.VersionTLS12,
	})
	elapsed := time.Since(startedAt)
	require.NoError(t, err)
	require.NoError(t, conn.Close())

	select {
	case <-accepted:
	case <-time.After(time.Second):
		t.Fatal("stalled first candidate was never dialed")
	}

	require.GreaterOrEqual(t, elapsed, 250*time.Millisecond)
	require.Less(t, elapsed, 3*time.Second)
}

func TestDirectTransferTimeoutDefaults(t *testing.T) {
	require.Equal(t, 5*time.Second, directTransferConnectTimeout)
	require.Equal(t, 2*time.Second, directTransferPerCandidateTimeout)
}

func newStalledDirectTransferTLSCandidate(t *testing.T) (string, <-chan struct{}) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	accepted := make(chan struct{}, 16)
	t.Cleanup(func() {
		_ = listener.Close()
	})

	go func() {
		for {
			conn, acceptErr := listener.Accept()
			if acceptErr != nil {
				return
			}
			select {
			case accepted <- struct{}{}:
			default:
			}
			go func() {
				<-time.After(5 * time.Second)
				_ = conn.Close()
			}()
		}
	}()

	return listener.Addr().String(), accepted
}

func directTransferTestCertificate(t *testing.T, serverName string) ([]byte, tls.Certificate) {
	t.Helper()

	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	caTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "Files Direct Transfer Test CA"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	require.NoError(t, err)

	leafKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	leafTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name{CommonName: serverName},
		DNSNames:     []string{serverName},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	leafDER, err := x509.CreateCertificate(rand.Reader, leafTemplate, caTemplate, &leafKey.PublicKey, caKey)
	require.NoError(t, err)

	leafKeyDER, err := x509.MarshalECPrivateKey(leafKey)
	require.NoError(t, err)
	cert, err := tls.X509KeyPair(
		pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: leafDER}),
		pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: leafKeyDER}),
	)
	require.NoError(t, err)

	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDER}), cert
}

func TestDirectTransferDialTLSCandidatesReturnsUnavailableWithoutCandidates(t *testing.T) {
	_, err := directTransferDialTLSCandidates(context.Background(), "tcp", nil, &tls.Config{})
	require.True(t, errors.Is(err, ErrDirectTransferUnavailable), "unexpected error: %v", err)
}
