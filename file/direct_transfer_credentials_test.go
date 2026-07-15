package file

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/stretchr/testify/require"
)

const testDirectCredentialServerName = "agent-123.agents.files.internal"

func TestDownloadDirectTransferDoesNotSendSDKCredentials(t *testing.T) {
	var gotHeaders http.Header
	var gotRequestURI string
	info := testDirectTransferInfo(t, "/downloads?jwt=direct-token", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeaders = r.Header.Clone()
		gotRequestURI = r.URL.RequestURI()
		_, _ = w.Write([]byte("ok"))
	}))
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "proxy fallback should not be used", http.StatusBadGateway)
	}))
	t.Cleanup(proxy.Close)
	headers := testDirectTransferCredentialHeaders()

	client := &Client{Config: files_sdk.Config{
		APIKey:           "config-api-key",
		EndpointOverride: proxy.URL,
		WorkspaceId:      9,
		Environment:      files_sdk.Development,
		Logger:           log.New(io.Discard, "", 0),
		AdditionalHeaders: map[string]string{
			"Authorization": "Bearer config-token",
			"Cookie":        "config_session=config-cookie",
		},
	}.Init()}
	_, err := client.Download(
		files_sdk.FileDownloadParams{File: files_sdk.File{
			Path:                 "/download.bin",
			DownloadUri:          proxy.URL + "/download.bin?jwt=proxy-token",
			DirectConnectionInfo: info,
		}},
		files_sdk.RequestHeadersOption(headers),
		files_sdk.ResponseBodyOption(func(body io.ReadCloser) error {
			return body.Close()
		}),
	)

	require.NoError(t, err)
	require.Equal(t, "/downloads?jwt=direct-token", gotRequestURI)
	requireDirectTransferCredentialsStripped(t, gotHeaders)
}

func TestDownloadDirectPartialBodyNeverAppendsProxyBody(t *testing.T) {
	info := testDirectTransferInfo(t, "/downloads?jwt=direct-token", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Length", "12")
		_, _ = w.Write([]byte("prefix"))
	}))
	proxyRequests := atomic.Int32{}
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		proxyRequests.Add(1)
		_, _ = w.Write([]byte("complete"))
	}))
	t.Cleanup(proxy.Close)
	destination, err := os.CreateTemp(t.TempDir(), "direct-download-*")
	require.NoError(t, err)
	t.Cleanup(func() { _ = destination.Close() })

	client := &Client{Config: files_sdk.Config{
		Environment: files_sdk.Development,
		Logger:      log.New(io.Discard, "", 0),
	}.Init()}
	_, downloadErr := client.Download(
		files_sdk.FileDownloadParams{File: files_sdk.File{
			Path:                 "/download.bin",
			DownloadUri:          proxy.URL + "/download.bin?jwt=proxy-token",
			DirectConnectionInfo: info,
		}},
		files_sdk.ResponseBodyOption(func(body io.ReadCloser) error {
			defer body.Close()
			_, copyErr := io.Copy(destination, body)
			return copyErr
		}),
	)
	require.NoError(t, destination.Sync())
	contents, err := os.ReadFile(destination.Name())
	require.NoError(t, err)

	if downloadErr == nil {
		require.Equal(t, "complete", string(contents))
	} else {
		require.Zero(t, proxyRequests.Load())
	}
	require.NotEqual(t, "prefixcomplete", string(contents))
}

func TestDownloadDirectRedirectFallsBackWithoutContactingTarget(t *testing.T) {
	targetRequests := atomic.Int32{}
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		targetRequests.Add(1)
		_, _ = w.Write([]byte("redirected"))
	}))
	t.Cleanup(target.Close)
	info := testDirectTransferInfo(t, "/downloads?jwt=direct-token", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Location", target.URL)
		w.WriteHeader(http.StatusFound)
	}))
	proxyRequests := atomic.Int32{}
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		proxyRequests.Add(1)
		_, _ = w.Write([]byte("proxy"))
	}))
	t.Cleanup(proxy.Close)

	client := &Client{Config: files_sdk.Config{
		Environment: files_sdk.Development,
		Logger:      log.New(io.Discard, "", 0),
	}.Init()}
	var body bytes.Buffer
	_, err := client.Download(
		files_sdk.FileDownloadParams{File: files_sdk.File{
			Path:                 "/download.bin",
			DownloadUri:          proxy.URL + "/download.bin?jwt=proxy-token",
			DirectConnectionInfo: info,
		}},
		files_sdk.ResponseBodyOption(func(responseBody io.ReadCloser) error {
			defer responseBody.Close()
			_, copyErr := io.Copy(&body, responseBody)
			return copyErr
		}),
	)

	require.NoError(t, err)
	require.Equal(t, "proxy", body.String())
	require.Equal(t, int32(1), proxyRequests.Load())
	require.Zero(t, targetRequests.Load())
}

func TestUploadDirectTransferDoesNotSendSDKCredentials(t *testing.T) {
	var gotHeaders http.Header
	var gotRequestURI string
	info := testDirectTransferInfo(t, "/uploads?jwt=direct-token", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeaders = r.Header.Clone()
		gotRequestURI = r.URL.RequestURI()
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, "abc", string(body))
		w.WriteHeader(http.StatusOK)
	}))
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "proxy fallback should not be used", http.StatusBadGateway)
	}))
	t.Cleanup(proxy.Close)
	headers := testDirectTransferCredentialHeaders()
	headers.Set("Content-Length", "3")

	u := &uploadIO{
		Client: &Client{Config: files_sdk.Config{
			Environment: files_sdk.Development,
			Logger:      log.New(io.Discard, "", 0),
		}.Init()},
		Path: "/upload.bin",
	}
	part := &Part{FileUploadPart: files_sdk.FileUploadPart{
		PartNumber:           1,
		UploadUri:            proxy.URL + "/uploads?jwt=proxy-token",
		DirectConnectionInfo: info,
	}}
	params := &files_sdk.CallParams{
		Method:  http.MethodPut,
		Config:  u.Config,
		Uri:     part.UploadUri,
		BodyIo:  io.NopCloser(bytes.NewReader([]byte("abc"))),
		Headers: headers,
		Context: context.Background(),
	}

	res, err := u.callUploadPart(context.Background(), part, params, func() bool { return true })
	require.NoError(t, err)
	t.Cleanup(func() { _ = res.Body.Close() })
	require.Equal(t, "/uploads?jwt=direct-token", gotRequestURI)
	requireDirectTransferCredentialsStripped(t, gotHeaders)
	require.Equal(t, "option-api-key", params.Headers.Get("X-FilesAPI-Key"))
	require.Equal(t, "files_session=option-cookie", params.Headers.Get("Cookie"))
}

func TestUploadV2DirectTransferDoesNotSendSDKCredentials(t *testing.T) {
	var gotHeaders http.Header
	var gotRequestURI string
	info := testDirectTransferInfo(t, "/uploads?jwt=direct-token", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeaders = r.Header.Clone()
		gotRequestURI = r.URL.RequestURI()
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, "abc", string(body))
		w.WriteHeader(http.StatusOK)
	}))
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "proxy fallback should not be used", http.StatusBadGateway)
	}))
	t.Cleanup(proxy.Close)
	headers := testDirectTransferCredentialHeaders()
	headers.Set("Content-Length", "3")

	u := &uploadIO{
		Client: &Client{Config: files_sdk.Config{
			Environment: files_sdk.Development,
			Logger:      log.New(io.Discard, "", 0),
		}.Init()},
		Path: "/upload.bin",
	}
	part := &uploadV2Part{uploadV2PartDescriptor: uploadV2PartDescriptor{
		number: 1,
		upload: files_sdk.FileUploadPart{
			PartNumber:           1,
			UploadUri:            proxy.URL + "/uploads?jwt=proxy-token",
			DirectConnectionInfo: info,
		},
	}}
	params := &files_sdk.CallParams{
		Method:  http.MethodPut,
		Config:  u.Config,
		Uri:     part.upload.UploadUri,
		BodyIo:  io.NopCloser(bytes.NewReader([]byte("abc"))),
		Headers: headers,
		Context: context.Background(),
	}

	res, err := (&uploadV2Engine{u: u}).callUploadV2Part(context.Background(), part, params, func() bool { return true }, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = res.Body.Close() })
	require.Equal(t, "/uploads?jwt=direct-token", gotRequestURI)
	requireDirectTransferCredentialsStripped(t, gotHeaders)
	require.Equal(t, "option-api-key", params.Headers.Get("X-FilesAPI-Key"))
	require.Equal(t, "files_session=option-cookie", params.Headers.Get("Cookie"))
}

func testDirectTransferCredentialHeaders() *http.Header {
	headers := &http.Header{}
	headers.Set("X-FilesAPI-Key", "option-api-key")
	headers.Set("X-FilesAPI-Auth", "option-session-id")
	headers.Set("X-Files-Reauthentication", "password:secret")
	headers.Set("X-Files-Workspace-Id", "123")
	headers.Set("Authorization", "Bearer option-token")
	headers.Set("Cookie", "files_session=option-cookie")
	return headers
}

func requireDirectTransferCredentialsStripped(t *testing.T, headers http.Header) {
	t.Helper()
	require.Empty(t, headers.Get("X-FilesAPI-Key"))
	require.Empty(t, headers.Get("X-FilesAPI-Auth"))
	require.Empty(t, headers.Get("X-Files-Reauthentication"))
	require.Empty(t, headers.Get("X-Files-Workspace-Id"))
	require.Empty(t, headers.Get("Authorization"))
	require.Empty(t, headers.Get("Cookie"))
	require.NotEmpty(t, headers.Get("User-Agent"))
}

func testDirectTransferInfo(t *testing.T, path string, handler http.HandlerFunc) files_sdk.DirectConnectionInfo {
	t.Helper()
	caPEM, cert := testDirectTransferCertificate(t, testDirectCredentialServerName)
	server := httptest.NewUnstartedServer(handler)
	server.TLS = &tls.Config{Certificates: []tls.Certificate{cert}}
	server.StartTLS()
	t.Cleanup(server.Close)

	_, port, err := net.SplitHostPort(server.Listener.Addr().String())
	require.NoError(t, err)
	return files_sdk.DirectConnectionInfo{
		Version:    1,
		ServerName: testDirectCredentialServerName,
		Addresses:  []string{"tcp://" + server.Listener.Addr().String()},
		DirectUri:  "https://" + net.JoinHostPort(testDirectCredentialServerName, port) + path,
		CaPem:      string(caPEM),
	}
}

func testDirectTransferCertificate(t *testing.T, serverName string) ([]byte, tls.Certificate) {
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
