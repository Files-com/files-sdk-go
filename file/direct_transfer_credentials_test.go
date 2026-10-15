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
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/stretchr/testify/require"
)

const testDirectCredentialServerName = "agent-123.agents.files.internal"

func TestDownloadDirectOnlyResponseSucceeds(t *testing.T) {
	info := testDirectTransferInfo(t, "/downloads?jwt=direct-token", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("direct"))
	}))
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(files_sdk.File{
			Path:                 "/download.bin",
			DirectConnectionInfo: info,
		}))
	}))
	t.Cleanup(api.Close)
	client := &Client{Config: files_sdk.Config{
		EndpointOverride: api.URL,
		Environment:      files_sdk.Development,
		Logger:           log.New(io.Discard, "", 0),
	}.Init()}
	var body bytes.Buffer

	_, err := client.Download(
		files_sdk.FileDownloadParams{Path: "/download.bin"},
		files_sdk.ResponseBodyOption(func(responseBody io.ReadCloser) error {
			defer responseBody.Close()
			_, copyErr := io.Copy(&body, responseBody)
			return copyErr
		}),
	)

	require.NoError(t, err)
	require.Equal(t, "direct", body.String())
}

func TestDownloadDirectOnlyResponsePreservesDirectError(t *testing.T) {
	info := testDirectTransferInfo(t, "/downloads?jwt=direct-token", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(files_sdk.File{
			Path:                 "/download.bin",
			DirectConnectionInfo: info,
		}))
	}))
	t.Cleanup(api.Close)
	client := &Client{Config: files_sdk.Config{
		EndpointOverride: api.URL,
		Environment:      files_sdk.Development,
		Logger:           log.New(io.Discard, "", 0),
	}.Init()}

	_, err := client.Download(files_sdk.FileDownloadParams{Path: "/download.bin"})

	var directErr *files_sdk.DirectTransferResponseError
	require.ErrorAs(t, err, &directErr)
	require.Equal(t, http.StatusServiceUnavailable, directErr.StatusCode)
	require.NotContains(t, err.Error(), "unsupported protocol scheme")

	f := (&File{
		File: &files_sdk.File{Path: "/download.bin"},
		FS:   (&FS{Context: context.Background()}).Init(client.Config, false),
	}).Init()
	_, err = f.Read(make([]byte, 1))

	require.ErrorAs(t, err, &directErr)
	require.Equal(t, http.StatusServiceUnavailable, directErr.StatusCode)
}

func TestDownloadRejectsResponseWithoutDownloadTarget(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(files_sdk.File{Path: "/download.bin"}))
	}))
	t.Cleanup(api.Close)
	client := &Client{Config: files_sdk.Config{
		EndpointOverride: api.URL,
		Logger:           log.New(io.Discard, "", 0),
	}.Init()}

	_, err := client.Download(files_sdk.FileDownloadParams{Path: "/download.bin"})

	require.EqualError(t, err, "download response did not include a usable download URL")

	f := (&File{
		File: &files_sdk.File{Path: "/download.bin"},
		FS:   (&FS{Context: context.Background()}).Init(client.Config, false),
	}).Init()
	_, err = f.Read(make([]byte, 1))

	require.EqualError(t, err, "download response did not include a usable download URL")
}

func TestDownloadDirectTransferDoesNotSendSDKCredentials(t *testing.T) {
	var gotHeaders http.Header
	var gotRequestURI string
	var logs bytes.Buffer
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
		Logger:           log.New(&logs, "", 0),
		Debug:            true,
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
	require.NotContains(t, logs.String(), "direct download attempt")
	require.Contains(t, logs.String(), "direct download success")
}

func TestDownloadDirectPartialBodyNeverAppendsProxyBody(t *testing.T) {
	for name, callbackErr := range map[string]func(copyErr error) error{
		"copy error": func(copyErr error) error { return copyErr },
		// A response option's own error may also identify as unavailable;
		// that must not turn consumed bytes into a fallback request.
		"error wrapping unavailable": func(copyErr error) error {
			return fmt.Errorf("%w: %w", files_sdk.ErrDirectTransferUnavailable, copyErr)
		},
	} {
		t.Run(name, func(t *testing.T) {
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
					return callbackErr(copyErr)
				}),
			)
			require.NoError(t, destination.Sync())
			contents, err := os.ReadFile(destination.Name())
			require.NoError(t, err)

			require.ErrorIs(t, downloadErr, files_sdk.ErrDirectTransferResponseStarted)
			require.ErrorIs(t, downloadErr, io.ErrUnexpectedEOF)
			require.Zero(t, proxyRequests.Load())
			require.Equal(t, "prefix", string(contents))
		})
	}
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

	// The proxy request keeps the caller's own client; the direct request never
	// goes through it.
	var callerClientHosts []string
	var callerClientMu sync.Mutex
	callerClient := &http.Client{Transport: uploadV2RoundTripFunc(func(request *http.Request) (*http.Response, error) {
		callerClientMu.Lock()
		callerClientHosts = append(callerClientHosts, request.URL.Host)
		callerClientMu.Unlock()
		return http.DefaultTransport.RoundTrip(request)
	})}
	client := &Client{Config: files_sdk.Config{
		Environment: files_sdk.Development,
		Logger:      log.New(io.Discard, "", 0),
	}.Init().SetCustomClient(callerClient)}
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
	proxyURL, err := url.Parse(proxy.URL)
	require.NoError(t, err)
	require.Equal(t, []string{proxyURL.Host}, callerClientHosts)
}

func TestDownloadDirect429RetriesThenFallsBackToProxy(t *testing.T) {
	// download sends one direct-capable Download whose direct endpoint answers
	// 429 with retryAfter until request number directSucceedsOn, if any.
	download := func(t *testing.T, ctx context.Context, retryAfter string, directSucceedsOn int32, firstDirect chan<- struct{}) (body string, direct, proxy int32, err error) {
		directRequests := atomic.Int32{}
		info := testDirectTransferInfo(t, "/downloads?jwt=direct-token", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			request := directRequests.Add(1)
			if request == directSucceedsOn {
				_, _ = w.Write([]byte("direct"))
				return
			}
			w.Header().Set("Retry-After", retryAfter)
			w.WriteHeader(http.StatusTooManyRequests)
			if request == 1 && firstDirect != nil {
				close(firstDirect)
			}
		}))
		proxyRequests := atomic.Int32{}
		proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			proxyRequests.Add(1)
			_, _ = w.Write([]byte("proxy"))
		}))
		t.Cleanup(proxyServer.Close)

		var output bytes.Buffer
		client := &Client{Config: files_sdk.Config{
			Environment: files_sdk.Development,
			Logger:      log.New(io.Discard, "", 0),
		}.Init()}
		_, err = client.Download(
			files_sdk.FileDownloadParams{File: files_sdk.File{
				Path:                 "/download.bin",
				DownloadUri:          proxyServer.URL + "/download.bin?jwt=proxy-token",
				DirectConnectionInfo: info,
			}},
			files_sdk.WithContext(ctx),
			files_sdk.ResponseBodyOption(func(responseBody io.ReadCloser) error {
				defer responseBody.Close()
				_, copyErr := io.Copy(&output, responseBody)
				return copyErr
			}),
		)
		return output.String(), directRequests.Load(), proxyRequests.Load(), err
	}

	t.Run("falls back after the last attempt", func(t *testing.T) {
		body, direct, proxy, err := download(t, context.Background(), "0", 0, nil)

		require.NoError(t, err)
		require.Equal(t, int32(downloadV2RetryAttempts), direct)
		require.Equal(t, int32(1), proxy)
		require.Equal(t, "proxy", body)
	})

	t.Run("sends directly again after Retry-After", func(t *testing.T) {
		started := time.Now()
		body, direct, proxy, err := download(t, context.Background(), "1", 2, nil)

		require.NoError(t, err)
		require.GreaterOrEqual(t, time.Since(started), time.Second)
		require.Equal(t, int32(2), direct)
		require.Zero(t, proxy)
		require.Equal(t, "direct", body)
	})

	t.Run("cancellation ends the Retry-After wait", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		firstDirect := make(chan struct{})
		go func() {
			<-firstDirect
			cancel()
		}()
		started := time.Now()
		_, direct, proxy, err := download(t, ctx, "60", 0, firstDirect)

		require.ErrorIs(t, err, context.Canceled)
		require.Less(t, time.Since(started), 10*time.Second)
		require.Equal(t, int32(1), direct)
		require.Zero(t, proxy)
	})
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
	headers.Set("Proxy-Authorization", "Basic proxy-secret")
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
	require.Empty(t, headers.Get("Proxy-Authorization"))
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
