package file

import (
	"context"
	"errors"
	"fmt"
	"io"
	fs2 "io/fs"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/Files-com/files-sdk-go/v3/folder"
	"github.com/Files-com/files-sdk-go/v3/ignore"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/lib/test"
	recorder "github.com/dnaeon/go-vcr/recorder"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type FileStatus struct {
	name   string
	status string
	size   int
	dir    bool
}

func CreateClient(fixture string) (client *Client, r *recorder.Recorder, err error) {
	client = &Client{}
	client.Config, r, err = test.CreateConfig(fixture)
	// Recorded fixtures predate direct connection info; focused tests cover the default-on behavior.
	client.Config.DisableDirectTransfers = true

	return client, r, err
}

func TestBoundedResponsesCloseOnceAndPreserveErrors(t *testing.T) {
	closeErr := errors.New("closing response failed")
	readErr := errors.New("reading response failed")
	for _, tc := range []struct {
		name, operation            string
		status                     int
		readErr, closeErr, wantErr error
	}{
		{name: "stats unsupported HEAD", operation: "stats", status: 422, closeErr: closeErr},
		{name: "stats close error", operation: "stats", status: 200, closeErr: closeErr, wantErr: closeErr},
		{name: "stats status error", operation: "stats", status: 503},
		{name: "read-at ignores close error", operation: "read", status: 206, closeErr: closeErr},
		{name: "read-at read error", operation: "read", status: 206, readErr: readErr, wantErr: readErr},
		{name: "read-at status error", operation: "read", status: 503},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := &ownedHTTPResponseBody{Reader: strings.NewReader("abc"), readErr: tc.readErr, closeErr: tc.closeErr}
			config := ownedResponseConfig(body, tc.status, http.Header{"Content-Type": []string{"text/plain"}})
			remote := files_sdk.File{Path: "file", DownloadUri: "https://download.test/file"}
			var err error
			if tc.operation == "stats" {
				_, err = (&Client{Config: config}).FileStats(remote)
			} else {
				fs := (&FS{Context: context.Background()}).Init(config, false)
				reader := (&File{File: &remote, FS: fs}).Init()
				_, err = reader.ReadAt(make([]byte, 3), 0)
			}
			if tc.status == 503 {
				var responseErr lib.ResponseError
				require.ErrorAs(t, err, &responseErr)
				assert.Equal(t, 503, responseErr.StatusCode)
				assert.Equal(t, "503 Service Unavailable", responseErr.Response.Status)
				assert.Equal(t, "text/plain", responseErr.Response.Header.Get("Content-Type"))
				assert.Contains(t, responseErr.Error(), "abc")
			} else if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, 1, body.closes)
		})
	}
}

func TestUploadV1ClosesErrorResponsesOnce(t *testing.T) {
	for _, tc := range []struct{ name, payload, contentType string }{
		{name: "nginx", payload: "denied", contentType: "text/html"},
		{name: "structured API", payload: `{"error":"denied","type":"permission-denied","http-code":"403"}`, contentType: "application/json"},
		{name: "API fallback", payload: "{broken", contentType: "application/json"},
		{name: "S3 XML", payload: `<Error><Code>AccessDenied</Code><Message>denied</Message></Error>`, contentType: "application/xml"},
		{name: "other status error", payload: "denied", contentType: "text/plain"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := &ownedHTTPResponseBody{Reader: strings.NewReader(tc.payload)}
			header := http.Header{"Content-Type": []string{tc.contentType}}
			if tc.name == "nginx" {
				header.Set("Server", "nginx")
				header.Set("X-Request-Id", "probe")
			}
			config := ownedResponseConfig(body, 403, header)
			err := (&Client{Config: config}).Upload(UploadWithReader(strings.NewReader("a")), UploadWithSize(1), UploadWithDestinationPath("file"), UploadWithResume(UploadResumable{FileUploadPart: files_sdk.FileUploadPart{Ref: "ref", Path: "file", PartNumber: 1, UploadUri: "https://upload.test/file", HttpMethod: "PUT", ParallelParts: lib.Bool(false)}}))
			require.Error(t, err)
			assert.Equal(t, 1, body.closes)
			if tc.name == "API fallback" || tc.name == "other status error" {
				var responseErr lib.ResponseError
				require.ErrorAs(t, err, &responseErr)
				assert.Equal(t, 403, responseErr.StatusCode)
				assert.Equal(t, "403 Forbidden", responseErr.Response.Status)
				assert.Contains(t, responseErr.Error(), tc.payload)
			}
		})
	}
}

func TestGetClosesOwnedResponsesOnceAndKeepsErrorDiagnostics(t *testing.T) {
	readErr := errors.New("get response interrupted")
	for _, tc := range []struct {
		name, payload, contentType string
		server, requestID          string
		status                     int
		readErr                    error
	}{
		{name: "success", status: 200, payload: `{"path":"file"}`, contentType: "application/json"},
		{name: "read failure", status: 200, payload: `{"path":"file"}`, contentType: "application/json", readErr: readErr},
		{name: "structured API", status: 403, payload: `{"error":"denied","type":"permission-denied","http-code":"403"}`, contentType: "application/json"},
		{name: "nginx", status: 403, payload: "denied", contentType: "text/html", server: "nginx", requestID: "probe"},
		{name: "versioned nginx", status: 403, payload: "denied", contentType: "text/html", server: "nginx/1.27", requestID: "probe"},
		{name: "Files gateway", status: 403, payload: "denied", contentType: "text/html", server: "files.com", requestID: "probe"},
		{name: "headerless gateway", status: 403, payload: "denied", contentType: "text/html"},
		{name: "API fallback", status: 403, payload: "{broken", contentType: "application/json"},
		{name: "generic status", status: 403, payload: "denied", contentType: "text/plain"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := &ownedHTTPResponseBody{Reader: strings.NewReader(tc.payload), readErr: tc.readErr}
			header := http.Header{"Content-Type": []string{tc.contentType}}
			if tc.server != "" {
				header.Set("Server", tc.server)
			}
			if tc.requestID != "" {
				header.Set("X-Request-Id", tc.requestID)
			}
			value, err := (&Client{Config: ownedResponseConfig(body, tc.status, header)}).Get("file")
			if tc.name == "success" {
				require.NoError(t, err)
				assert.Equal(t, "file", value.Path)
			} else if tc.readErr != nil {
				require.ErrorIs(t, err, tc.readErr)
			} else {
				require.Error(t, err)
				if tc.name == "nginx" {
					assert.EqualError(t, err, "files.com Server error - request id: probe")
				} else if tc.contentType == "text/html" {
					assert.NotContains(t, err.Error(), "files.com Server error")
				}
				if tc.name == "API fallback" {
					var responseErr lib.ResponseError
					require.ErrorAs(t, err, &responseErr)
					assert.Equal(t, "403 Forbidden", responseErr.Response.Status)
					assert.Contains(t, responseErr.Error(), tc.payload)
				} else if tc.name == "generic status" {
					assert.Contains(t, err.Error(), tc.payload)
				}
			}
			assert.Equal(t, 1, body.closes)
		})
	}
}

type ownedHTTPResponseBody struct {
	*strings.Reader
	readErr, closeErr error
	closes            int
}

func (b *ownedHTTPResponseBody) Read(p []byte) (int, error) {
	if b.readErr != nil {
		return 0, b.readErr
	}
	return b.Reader.Read(p)
}
func (b *ownedHTTPResponseBody) Close() error { b.closes++; return b.closeErr }

type ownedHTTPResponseTransport func(*http.Request) (*http.Response, error)

func (f ownedHTTPResponseTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
func ownedResponseConfig(body *ownedHTTPResponseBody, status int, header http.Header) files_sdk.Config {
	config := files_sdk.Config{DisableDirectTransfers: true}.Init().SetCustomClient(&http.Client{Transport: ownedHTTPResponseTransport(func(req *http.Request) (*http.Response, error) {
		if req.Body != nil {
			_, _ = io.Copy(io.Discard, req.Body)
			_ = req.Body.Close()
		}
		return &http.Response{StatusCode: status, Status: fmt.Sprintf("%d %s", status, http.StatusText(status)), Header: header, Body: body, ContentLength: int64(body.Len()), Request: req}, nil
	})})
	config.RetryMax = 0
	return config
}

func TestClient_DownloadUriRequestsDirectConnectionInfoByDefault(t *testing.T) {
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	server.MockFiles["remote.txt"] = mockFile{}

	_, err := server.Client().DownloadUri(files_sdk.FileDownloadParams{Path: "remote.txt"})

	require.NoError(t, err)
	require.Len(t, server.TrackRequest["/api/rest/v1/files/*path"], 1)
	assert.Contains(t, server.TrackRequest["/api/rest/v1/files/*path"][0], "with_direct_connection_info=true")
}

func TestClient_DownloadUriOmitsDirectConnectionInfoWhenDisabled(t *testing.T) {
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	server.MockFiles["remote.txt"] = mockFile{}
	client := server.Client()
	client.Config.DisableDirectTransfers = true

	_, err := client.DownloadUri(files_sdk.FileDownloadParams{Path: "remote.txt"})

	require.NoError(t, err)
	require.Len(t, server.TrackRequest["/api/rest/v1/files/*path"], 1)
	assert.NotContains(t, server.TrackRequest["/api/rest/v1/files/*path"][0], "with_direct_connection_info")
}

func TestClient_DownloadDoesNotSendAuthHeadersToOffOriginDownloadUri(t *testing.T) {
	var gotHeaders http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeaders = r.Header.Clone()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &Client{Config: files_sdk.Config{
		APIKey: "api-key",
		AdditionalHeaders: map[string]string{
			"X-FilesAPI-Auth":          "additional-session-id",
			"X-Files-Reauthentication": "password:additional-password",
		},
	}.Init()}
	_, err := client.Download(
		files_sdk.FileDownloadParams{File: files_sdk.File{Path: "remote.txt", DownloadUri: server.URL + "/download"}},
		files_sdk.ResponseBodyOption(func(body io.ReadCloser) error {
			return body.Close()
		}),
	)

	require.NoError(t, err)
	assert.Empty(t, gotHeaders.Get("X-FilesAPI-Key"))
	assert.Empty(t, gotHeaders.Get("X-FilesAPI-Auth"))
	assert.Empty(t, gotHeaders.Get("X-Files-Reauthentication"))
	assert.NotEmpty(t, gotHeaders.Get("User-Agent"))
}

func TestClient_DownloadSendsAuthHeadersToConfiguredEndpoint(t *testing.T) {
	var gotHeaders http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeaders = r.Header.Clone()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &Client{Config: files_sdk.Config{APIKey: "api-key", EndpointOverride: server.URL}.Init()}
	_, err := client.Download(
		files_sdk.FileDownloadParams{File: files_sdk.File{Path: "remote.txt", DownloadUri: server.URL + "/download"}},
		files_sdk.ResponseBodyOption(func(body io.ReadCloser) error {
			return body.Close()
		}),
	)

	require.NoError(t, err)
	assert.Equal(t, "api-key", gotHeaders.Get("X-FilesAPI-Key"))
	assert.Empty(t, gotHeaders.Get("X-FilesAPI-Auth"))
}

func TestClient_DownloadDoesNotSendAuthHeadersAfterOffOriginRedirect(t *testing.T) {
	var firstRequestHeaders http.Header
	var redirectedRequestHeaders http.Header

	storageServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirectedRequestHeaders = r.Header.Clone()
		w.WriteHeader(http.StatusOK)
	}))
	defer storageServer.Close()

	filesServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		firstRequestHeaders = r.Header.Clone()
		http.Redirect(w, r, storageServer.URL+"/object", http.StatusFound)
	}))
	defer filesServer.Close()

	client := &Client{Config: files_sdk.Config{
		APIKey:           "api-key",
		EndpointOverride: filesServer.URL,
		AdditionalHeaders: map[string]string{
			"X-FilesAPI-Auth":          "additional-session-id",
			"X-Files-Reauthentication": "password:additional-password",
		},
	}.Init()}
	_, err := client.Download(
		files_sdk.FileDownloadParams{File: files_sdk.File{Path: "remote.txt", DownloadUri: filesServer.URL + "/download"}},
		files_sdk.ResponseBodyOption(func(body io.ReadCloser) error {
			return body.Close()
		}),
	)

	require.NoError(t, err)
	assert.Equal(t, "api-key", firstRequestHeaders.Get("X-FilesAPI-Key"))
	assert.Equal(t, "additional-session-id", firstRequestHeaders.Get("X-FilesAPI-Auth"))
	assert.Equal(t, "password:additional-password", firstRequestHeaders.Get("X-Files-Reauthentication"))
	assert.Empty(t, redirectedRequestHeaders.Get("X-FilesAPI-Key"))
	assert.Empty(t, redirectedRequestHeaders.Get("X-FilesAPI-Auth"))
	assert.Empty(t, redirectedRequestHeaders.Get("X-Files-Reauthentication"))
}

func TestClient_DownloadReturnsParseErrorForMalformedDownloadUri(t *testing.T) {
	client := &Client{Config: files_sdk.Config{APIKey: "api-key"}.Init()}
	_, err := client.Download(files_sdk.FileDownloadParams{File: files_sdk.File{Path: "remote.txt", DownloadUri: "https://storage.example.test/%zz"}})

	var urlErr *url.Error
	require.ErrorAs(t, err, &urlErr)
	assert.Equal(t, "parse", urlErr.Op)
}

func TestClient_DownloadRequestStatusDoesNotSendAuthHeadersToOffOriginStatusURL(t *testing.T) {
	var gotHeaders http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeaders = r.Header.Clone()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"status":"completed","bytes_transferred":12}}`))
	}))
	defer server.Close()

	client := &Client{Config: files_sdk.Config{
		SessionId: "session-id",
		AdditionalHeaders: map[string]string{
			"x-filesapi-key":           "additional-api-key",
			"X-Files-Reauthentication": "password:additional-password",
		},
	}.Init()}
	status, err := client.DownloadRequestStatus(server.URL+"/download", "request-id")

	require.NoError(t, err)
	assert.Equal(t, "completed", status.Data.Status)
	assert.Empty(t, gotHeaders.Get("X-FilesAPI-Key"))
	assert.Empty(t, gotHeaders.Get("X-FilesAPI-Auth"))
	assert.Empty(t, gotHeaders.Get("X-Files-Reauthentication"))
	assert.NotEmpty(t, gotHeaders.Get("User-Agent"))
}

func createFileStructure(t *testing.T, tmpDir string, filesAndStatus []FileStatus) (filePaths []string) {
	for index, file := range filesAndStatus {
		file.name = filepath.Join(tmpDir, file.name)
		filesAndStatus[index] = file
		if file.dir {
			require.NoError(t, os.MkdirAll(file.name, 0750))
		} else {
			require.NoError(t, os.MkdirAll(filepath.Dir(file.name), 0750))
			f, err := os.Create(file.name)
			assert.NoError(t, err)
			f.Write([]byte("hello how are you doing?"))
			f.Close()
			if file.status == "complete" {
				filePaths = append(filePaths, f.Name())
			}
		}
	}
	return
}

func deletePath(client *Client, path string) {
	err := client.Delete(files_sdk.FileDeleteParams{Path: path})
	var responseError files_sdk.ResponseError
	ok := errors.As(err, &responseError)
	if ok && responseError.Type == "not-found" {
	} else if ok && responseError.Type == "processing-failure/folder-not-empty" {
		err = client.Delete(files_sdk.FileDeleteParams{Path: path, Recursive: lib.Bool(true)})
		ok = errors.As(err, &responseError)
		if ok && responseError.Type == "not-found" {
			//noop
		} else if ok {
			panic(err)
		}
	} else if ok {
		panic(err)
	}
}

func ignoreSomeErrors(err error) {
	if err != nil && !files_sdk.IsExist(err) {
		panic(err)
	}
}

func buildScenario(base string, client *Client) {
	folderClient := folder.Client{Config: client.Config}

	_, err := folderClient.Create(files_sdk.FolderCreateParams{Path: base})
	ignoreSomeErrors(err)
	_, err = folderClient.Create(files_sdk.FolderCreateParams{Path: lib.UrlJoinNoEscape(base, "nested_1")})
	ignoreSomeErrors(err)
	_, err = folderClient.Create(files_sdk.FolderCreateParams{Path: lib.UrlJoinNoEscape(base, "nested_1", "nested_2")})
	ignoreSomeErrors(err)
	_, err = folderClient.Create(files_sdk.FolderCreateParams{Path: lib.UrlJoinNoEscape(base, "nested_1", "nested_2", "nested_3")})
	ignoreSomeErrors(err)

	client.Upload()

	err = client.Upload(
		UploadWithSize(9),
		UploadWithDestinationPath(lib.UrlJoinNoEscape(base, "nested_1", "nested_2", "3.text")),
		UploadWithProvidedMtime(time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC)),
		UploadWithReader(strings.NewReader("testing 3")),
	)

	ignoreSomeErrors(err)

	err = client.Upload(
		UploadWithSize(9),
		UploadWithDestinationPath(lib.UrlJoinNoEscape(base, "nested_1", "nested_2", "nested_3", "4.text")),
		UploadWithProvidedMtime(time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC)),
		UploadWithReader(strings.NewReader("testing 3")),
	)
	ignoreSomeErrors(err)
}

func runDownloadScenario(path string, destination string, client *Client) map[string][]JobFile {
	m := &sync.Mutex{}
	results := make(map[string][]JobFile)

	reporter := func(r JobFile) {
		m.Lock()
		results[r.File.Path] = append(results[r.File.Path], r)
		m.Unlock()
	}

	job := client.Downloader(
		DownloaderParams{RemotePath: path, LocalPath: destination, EventsReporter: CreateReporter(reporter)},
	)

	job.Start()
	job.Wait()

	return results
}

func CreateReporter(callback Reporter) EventsReporter {
	return CreateFileEvents(callback, append(status.Excluded, append(status.Included, OnBytesChange(status.Uploading))...)...)
}

func TestClient_UploadFolder(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFolder")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	assert := assert.New(t)
	resultsMapMutex := sync.RWMutex{}
	results := make(map[string][]ReporterCall)

	job := client.Uploader(
		UploaderParams{
			LocalPath:  "../lib",
			RemotePath: "golib",
			EventsReporter: CreateReporter(func(status JobFile) {
				resultsMapMutex.Lock()
				results[status.RemotePath] = append(results[status.RemotePath], ReporterCall{JobFile: status, err: status.Err})
				resultsMapMutex.Unlock()
			}),
			Manager: manager.Default(),
		},
	)

	job.Start()
	job.Wait()
	files, err := os.ReadDir("../lib")
	assert.NoError(err)
	gitIgnore, err := ignore.New()
	assert.NoError(err)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if gitIgnore.MatchesPath(f.Name()) {
			continue
		}
		remotePath := fmt.Sprintf("golib/lib/%v", f.Name())
		assert.Contains(results, remotePath)
		lastStatuses, ok := results[remotePath]
		if !ok {
			continue
		}
		lastStatus := lastStatuses[len(lastStatuses)-1]
		if lastStatus.Err != nil && strings.Contains(lastStatus.Err.Error(), "Requested interaction not found") {
			assert.Equal(status.Errored, lastStatus.Status)
		} else {
			assert.Equal(status.Complete, lastStatus.Status)
			assert.NoError(lastStatus.Err)
		}
	}
	deletePath(client, "golib")
}

func TestClient_UploadFolder_Dot(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFolder_Dot")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	assert := assert.New(t)
	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	resultsMapMutex := sync.RWMutex{}
	results := make(map[string][]int64)
	err = os.MkdirAll(filepath.Join(tmpDir, "dot"), 0755)
	assert.NoError(err)

	f1, _ := os.Create(filepath.Join(tmpDir, "dot/1.text"))
	f1.Write([]byte("hello 1"))
	f1.Close()

	f2, _ := os.Create(filepath.Join(tmpDir, "dot/2.text"))
	f2.Write([]byte("hello 2"))
	f2.Close()

	f3, _ := os.Create(filepath.Join(tmpDir, "dot/3.text"))
	f3.Write([]byte("hello 3"))
	f3.Close()

	currentPwd, _ := os.Getwd()
	err = os.Chdir(filepath.Join(tmpDir, "dot"))
	defer func() {
		os.Chdir(currentPwd)
		os.RemoveAll(filepath.Join(tmpDir, "/dot"))
	}()
	assert.NoError(err)

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	job := client.Uploader(
		UploaderParams{
			LocalPath:  "." + string(os.PathSeparator),
			RemotePath: "go-from-dot",
			EventsReporter: CreateReporter(func(s JobFile) {
				resultsMapMutex.Lock()
				assert.NoError(s.Err)

				results[s.File.Path] = append(results[s.File.Path], s.TransferBytes)
				resultsMapMutex.Unlock()
			}),
		}, files_sdk.WithContext(ctx))
	job.Start()
	job.Wait()
	assert.Contains(results, "go-from-dot/1.text")
	assert.Contains(results, "go-from-dot/2.text")
	assert.Contains(results, "go-from-dot/3.text")
	assert.Equal(int64(7), job.Statuses[0].TransferBytes())
	assert.Equal(int64(7), job.Statuses[1].TransferBytes())
	assert.Equal(int64(7), job.Statuses[2].TransferBytes())

	deletePath(client, "go-from-dot")
}

func TestClient_UploadFolder_Relative(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFolder_Relative")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	assert := assert.New(t)
	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	resultsMapMutex := sync.RWMutex{}
	results := make(map[string][]int64)

	err = os.MkdirAll(filepath.Join(tmpDir, "relative"), 0755)
	assert.NoError(err)

	f1, _ := os.Create(filepath.Join(tmpDir, "relative", "1.text"))
	f1.Write([]byte("hello 1"))
	f1.Close()

	f2, _ := os.Create(filepath.Join(tmpDir, "relative", "2.text"))
	f2.Write([]byte("hello 2"))
	f2.Close()

	f3, _ := os.Create(filepath.Join(tmpDir, "relative", "3.text"))
	f3.Write([]byte("hello 3"))
	f3.Close()

	currentPwd, _ := os.Getwd()
	err = os.Chdir(tmpDir)
	defer os.Chdir(currentPwd)
	assert.NoError(err)

	job := client.Uploader(
		UploaderParams{
			LocalPath:  "relative" + string(os.PathSeparator),
			RemotePath: "relative",
			EventsReporter: CreateReporter(func(status JobFile) {
				resultsMapMutex.Lock()
				results[status.File.Path] = append(results[status.File.Path], status.TransferBytes)
				resultsMapMutex.Unlock()
			}),
		})
	job.Start()
	job.Wait()
	assert.Contains(results, "relative/1.text")
	assert.Contains(results, "relative/2.text")
	assert.Contains(results, "relative/3.text")
	assert.Equal(int64(7), job.Statuses[0].TransferBytes())
	assert.Equal(int64(7), job.Statuses[1].TransferBytes())
	assert.Equal(int64(7), job.Statuses[2].TransferBytes())
	assert.Equal(int64(21), job.TotalBytes(status.Valid...))
	assert.True(job.All(status.Ended...), "All jobs should have ended")

	deletePath(client, "relative")
}

func TestClient_Uploader(t *testing.T) {
	client, r, err := CreateClient("TestClient_Uploader")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)

	uploadPath := ".." + string(os.PathSeparator) + "LICENSE"
	job := client.Uploader(UploaderParams{LocalPath: uploadPath})
	job.Start()
	job.Wait()
	assert.True(job.Started.Called(), "Job did not start")
	assert.True(job.Scanning.Called(), "Job did not start scanning")
	assert.True(job.EndScanning.Called(), "Job did not end scanning")
	assert.True(job.Finished.Called(), "Job was not finished")
	assert.False(job.Canceled.Called(), "Job should not have been cancelled")
	assert.Equal("LICENSE", job.Files()[0].DisplayName)
	assert.Equal(1, job.Count())
	assert.Equal(int64(1102), job.TotalBytes())
	assert.True(job.All(status.Ended...), "All jobs should have ended")

	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tempFile, err := os.CreateTemp(tmpDir, "LICENSE")
	if err != nil {
		panic(err)
	}
	file, err := client.DownloadToFile(files_sdk.FileDownloadParams{Path: "LICENSE"}, tempFile.Name())
	assert.NoError(err)

	assert.Equal(file.DisplayName, "LICENSE")

	downloadData, err := os.ReadFile(tempFile.Name())
	if err != nil {
		panic(err)
	}
	localData, err := os.ReadFile(uploadPath)
	if err != nil {
		panic(err)
	}
	assert.Equal(string(downloadData), string(localData))
	defer os.Remove(tempFile.Name())
}

func TestClient_UploadFile_To_Existing_Dir(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFile_To_Existing_Dir")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)

	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	folderClient := folder.Client{Config: client.Config}
	_, err = folderClient.Create(files_sdk.FolderCreateParams{Path: "docs"})
	defer deletePath(client, "docs")

	assert.NoError(err)
	uploadPath := ".." + string(os.PathSeparator) + "LICENSE"
	job := client.Uploader(UploaderParams{LocalPath: uploadPath, RemotePath: "docs"})
	job.Start()
	job.Wait()
	tempFile, err := os.CreateTemp(tmpDir, "LICENSE")
	if err != nil {
		panic(err)
	}
	file, err := client.DownloadToFile(files_sdk.FileDownloadParams{Path: "docs/LICENSE"}, tempFile.Name())
	assert.NoError(err)

	assert.Equal(file.DisplayName, "LICENSE")

	downloadData, err := os.ReadFile(tempFile.Name())
	if err != nil {
		panic(err)
	}
	localData, err := os.ReadFile(uploadPath)
	if err != nil {
		panic(err)
	}
	assert.Equal(string(downloadData), string(localData))
	defer os.Remove(tempFile.Name())
}

func TestClient_UploadFile_To_NonExistingPath(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFile_To_NonExistingPath")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)

	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	deletePath(client, "taco")
	uploadPath := ".." + string(os.PathSeparator) + "LICENSE"
	job := client.Uploader(UploaderParams{LocalPath: uploadPath, RemotePath: "taco"})
	defer deletePath(client, "taco")
	job.Start()
	job.Wait()
	tempFile, _ := os.CreateTemp(tmpDir, "LICENSE")
	_, err = filepath.Abs(filepath.Dir(tempFile.Name()))
	if err != nil {
		panic(err)
	}
	file, err := client.DownloadToFile(files_sdk.FileDownloadParams{Path: "taco"}, tempFile.Name())
	assert.NoError(err)

	assert.Equal("taco", file.DisplayName, "because the docs did not exist as a folder it becomes the file")

	downloadData, err := os.ReadFile(tempFile.Name())
	assert.NoError(err)
	localData, err := os.ReadFile(uploadPath)
	if err != nil {
		panic(err)
	}
	assert.Equal(string(downloadData), string(localData))
	defer os.Remove(tempFile.Name())
}

func TestClient_UploadFile_To_NonExistingPath_WithSlash(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFile_To_NonExistingPath_WithSlash")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)

	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	assert.NoError(err)
	uploadPath := ".." + string(os.PathSeparator) + "LICENSE"
	deletePath(client, "docs")
	job := client.Uploader(UploaderParams{LocalPath: uploadPath, RemotePath: "docs/"})
	defer deletePath(client, "docs")
	job.Start()
	job.Wait()
	tempFile, err := os.CreateTemp(tmpDir, "LICENSE")
	if err != nil {
		panic(err)
	}
	file, err := client.DownloadToFile(files_sdk.FileDownloadParams{Path: "docs/LICENSE"}, tempFile.Name())
	assert.NoError(err)

	assert.Equal("file", file.Type)
	assert.Equal("LICENSE", file.DisplayName)

	downloadData, err := os.ReadFile(tempFile.Name())
	if err != nil {
		panic(err)
	}
	localData, err := os.ReadFile(uploadPath)
	if err != nil {
		panic(err)
	}
	assert.Equal(string(downloadData), string(localData))
	defer os.Remove(tempFile.Name())
}

func TestClient_UploadFolder_as_file2(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFolder_as_file2")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)
	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	uploadPath := ".." + string(os.PathSeparator) + "LICENSE"
	job := client.Uploader(UploaderParams{LocalPath: uploadPath})
	job.Start()
	job.Wait()

	assert.Equal(int64(1102), job.TransferBytes())
	assert.Equal(int64(1102), job.TotalBytes())

	tempFile, err := os.CreateTemp(tmpDir, "LICENSE")
	if err != nil {
		panic(err)
	}
	file, err := client.DownloadToFile(files_sdk.FileDownloadParams{Path: "LICENSE"}, tempFile.Name())
	assert.NoError(err)

	assert.Equal(file.DisplayName, "LICENSE")

	downloadData, err := os.ReadFile(tempFile.Name())
	if err != nil {
		panic(err)
	}
	localData, err := os.ReadFile(uploadPath)
	if err != nil {
		panic(err)
	}
	assert.Equal(string(downloadData), string(localData))
	defer os.Remove(tempFile.Name())
}

func TestClient_UploadFolder_RemotePathWithStartingSlash(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFolder_RemotePathWithStartingSlash")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	err = os.MkdirAll(filepath.Join(tmpDir, "test"), 0755)
	assert.NoError(t, err)

	f1, _ := os.Create(filepath.Join(tmpDir, "test", "1.text"))
	f1.Write([]byte("hello 1"))
	f1.Close()

	f2, _ := os.Create(filepath.Join(tmpDir, "test", "2.text"))
	f2.Write([]byte("hello 2"))
	f2.Close()

	f3, _ := os.Create(filepath.Join(tmpDir, "test", "3.text"))
	f3.Write([]byte("hello 3"))
	f3.Close()
	job := client.Uploader(UploaderParams{LocalPath: filepath.Join(tmpDir, "test"), RemotePath: "/test", Manager: manager.Sync()})
	job.Start()
	job.Wait()
	assert.Len(t, job.Statuses, 4)
	assert.NoError(t, job.Statuses[1].Err())
	dir, _ := filepath.Split(job.Statuses[1].RemotePath())
	assert.Equal(t, "test/test/", dir)
}

func TestClient_UploadFolder_ZeroByteFile(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFolder_ZeroByteFile")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	err = os.MkdirAll(filepath.Join(tmpDir, "zero_byte_folder"), 0755)
	assert.NoError(t, err)

	f1, _ := os.Create(filepath.Join(tmpDir, "zero_byte_folder", "zero-byte-file.text"))
	f1.Close()

	job := client.Uploader(UploaderParams{LocalPath: filepath.Join(tmpDir, "zero_byte_folder"), RemotePath: "", Manager: manager.Sync()})
	job.Start()
	job.Wait()
	require.Len(t, job.Statuses, 2)
	assert.NoError(t, job.Statuses[0].Err())
	assert.NoError(t, job.Statuses[1].Err())
	assert.Equal(t, "zero_byte_folder", job.Statuses[0].RemotePath())
	assert.Equal(t, "zero_byte_folder/zero-byte-file.text", job.Statuses[1].RemotePath())

	job = client.Downloader(DownloaderParams{RemotePath: "zero_byte_folder", LocalPath: tmpDir + string(os.PathSeparator)})
	job.Start()
	job.Wait()
	require.Len(t, job.Statuses, 2)
	assert.NoError(t, job.Statuses[0].Err())
	assert.NoError(t, job.Statuses[1].Err())
}

func TestClient_DownloadFolder(t *testing.T) {
	client, r, err := CreateClient("TestClient_DownloadFolder")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	buildScenario("TestClient_DownloadFolder", client)

	assert := assert.New(t)

	it, err := client.ListFor(files_sdk.FolderListForParams{
		ListParams: files_sdk.ListParams{PerPage: 1},
		Path:       "TestClient_DownloadFolder/nested_1/nested_2",
	})

	assert.NoError(err)
	files := files_sdk.FileCollection{}
	for it.Next() {
		files = append(files, it.File())
	}

	assert.Len(files, 2, "something is wrong with cursor")

	results := runDownloadScenario("TestClient_DownloadFolder", "download/", client)
	assert.NoError(err)
	assert.Equal(int64(9), results["TestClient_DownloadFolder/nested_1/nested_2/3.text"][2].TransferBytes)
	assert.Equal(int64(9), results["TestClient_DownloadFolder/nested_1/nested_2/nested_3/4.text"][2].TransferBytes)
	os.RemoveAll("download")
}

func TestClient_DownloadFolder_Smart(t *testing.T) {
	client, r, err := CreateClient("TestClient_DownloadFolder_Smart")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	defer os.RemoveAll("download")

	buildScenario("TestClient_DownloadFolder_Smart", client)

	results := runDownloadScenario(lib.UrlJoinNoEscape("TestClient_DownloadFolder_Smart", "nested_1", "nested_2", "3.text"), "download"+string(os.PathSeparator), client)
	for _, result := range results {
		for _, f := range result {
			require.NoError(t, f.Err)
		}
	}

	assert.Len(t, results["TestClient_DownloadFolder_Smart/nested_1/nested_2/3.text"], 3)
	assert.Equal(t, int64(9), results["TestClient_DownloadFolder_Smart/nested_1/nested_2/3.text"][2].TransferBytes)

	results2 := runDownloadScenario(lib.UrlJoinNoEscape("TestClient_DownloadFolder_Smart", "nested_1", "nested_2")+"/", "download", client)

	assert.NoError(t, err)

	path, err := os.Getwd()
	assert.NoError(t, err)

	require.Len(t, results2["TestClient_DownloadFolder_Smart/nested_1/nested_2/3.text"], 3)
	assert.Equal(t, int64(9), results2["TestClient_DownloadFolder_Smart/nested_1/nested_2/3.text"][2].TransferBytes)
	assert.Equal(t, filepath.Join(path, "download", "3.text"), results2["TestClient_DownloadFolder_Smart/nested_1/nested_2/3.text"][2].LocalPath)
	fileInfo, err := os.Stat(filepath.Join(path, "download", "3.text"))
	require.NoError(t, err)
	assert.Equal(t, int64(9), fileInfo.Size())
	require.Len(t, results2["TestClient_DownloadFolder_Smart/nested_1/nested_2/nested_3/4.text"], 3)
	assert.Equal(t, int64(9), results2["TestClient_DownloadFolder_Smart/nested_1/nested_2/nested_3/4.text"][2].TransferBytes)
	assert.Equal(t, filepath.Join(path, "download", "nested_3", "4.text"), results2["TestClient_DownloadFolder_Smart/nested_1/nested_2/nested_3/4.text"][2].LocalPath)
}

func TestClient_DownloadFolder_file_to_file(t *testing.T) {
	client, r, err := CreateClient("TestClient_DownloadFolder_file_to_file")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	buildScenario("TestClient_DownloadFolder_file_to_file", client)
	assert := assert.New(t)

	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	results := runDownloadScenario(lib.UrlJoinNoEscape("TestClient_DownloadFolder_file_to_file", "nested_1", "nested_2", "3.text"), filepath.Join(tmpDir, "3.text"), client)
	assert.NoError(err)
	for _, result := range results {
		for _, f := range result {
			require.NoError(t, f.Err)
		}
	}
	assert.Equal(int64(9), results["TestClient_DownloadFolder_file_to_file/nested_1/nested_2/3.text"][2].TransferBytes)
}

func TestClient_DownloadFolder_file_to_implicit(t *testing.T) {
	client, r, err := CreateClient("TestClient_DownloadFolder_file_to_implicit")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	buildScenario("file_to_implicit", client)
	assert := assert.New(t)
	results := runDownloadScenario(lib.UrlJoinNoEscape("file_to_implicit", "nested_1", "nested_2", "3.text"), "", client)
	assert.NoError(err)
	for _, result := range results {
		for _, f := range result {
			require.NoError(t, f.Err)
		}
	}
	assert.Equal(int64(9), results["file_to_implicit/nested_1/nested_2/3.text"][2].TransferBytes)
	os.RemoveAll("3.text")
}

func TestClient_DownloadFolder_file_only(t *testing.T) {
	client, r, err := CreateClient("TestClient_DownloadFolder_file_only")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	err = client.Upload(
		UploadWithSize(5),
		UploadWithDestinationPath("i am at the root.text"),
		UploadWithProvidedMtime(time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC)),
		UploadWithReader(strings.NewReader("hello")),
	)

	require.NoError(t, err)

	assert := assert.New(t)
	results := runDownloadScenario("i am at the root.text", "", client)
	assert.NoError(err)
	for _, result := range results {
		for _, f := range result {
			require.NoError(t, f.Err)
		}
	}
	require.Len(t, results["i am at the root.text"], 3)
	assert.Equal(int64(5), results["i am at the root.text"][2].TransferBytes)
	os.RemoveAll("i am at the root.text")
}

func TestClient_Downloader_Delete_Source(t *testing.T) {
	client, r, err := CreateClient("TestClient_Downloader_Delete_Source")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)

	folderClient := folder.Client{Config: client.Config}

	folderClient.Create(files_sdk.FolderCreateParams{Path: "test-delete-source"})

	err = client.Upload(
		UploadWithSize(9),
		UploadWithDestinationPath(lib.UrlJoinNoEscape("test-delete-source", "test.text")),
		UploadWithProvidedMtime(time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC)),
		UploadWithReader(strings.NewReader("testing 3")),
	)

	require.NoError(t, err)
	localPath, err := os.MkdirTemp("", "TestClient_Downloader_Delete_Source")
	require.NoError(t, err)

	job := client.Downloader(
		DownloaderParams{RemotePath: "test-delete-source", LocalPath: localPath},
	)
	var fi JobFile
	var log status.Log
	job.RegisterFileEvent(func(f JobFile) {
		fi = f
		log, err = DeleteSource{Config: client.Config, Direction: job.Direction}.Call(f)
	}, status.Complete)
	job.Start()
	<-job.Finished.C
	assert.NoError(err)
	assert.Equal("delete source", log.Action)
	assert.Equal(fi.RemotePath, log.Path)

	_, err = client.Find(files_sdk.FileFindParams{Path: lib.UrlJoinNoEscape("test-delete-source", "test.text")})
	require.NotNil(t, err)
	assert.Equal("Not Found - `Not Found`", err.Error())
	os.RemoveAll("test.text")
}

func TestClient_Downloader_Delete_Source_Folder(t *testing.T) {
	client, r, err := CreateClient("TestClient_Downloader_Delete_Source_Folder")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)

	folderClient := folder.Client{Config: client.Config}

	folderClient.Create(files_sdk.FolderCreateParams{Path: "test-delete-source-folder"})

	err = client.Upload(
		UploadWithSize(9),
		UploadWithDestinationPath(lib.UrlJoinNoEscape("test-delete-source-folder", "test.text")),
		UploadWithProvidedMtime(time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC)),
		UploadWithReader(strings.NewReader("testing 3")),
	)

	require.NoError(t, err)
	localPath, err := os.MkdirTemp("", "TestClient_Downloader_Delete_Source_Folder")
	require.NoError(t, err)

	job := client.Downloader(
		DownloaderParams{RemotePath: "test-delete-source-folder", LocalPath: localPath},
	)
	var fi JobFile
	var log status.Log

	job.RegisterFileEvent(func(f JobFile) {
		fi = f
		DeleteSource{Config: client.Config, Direction: job.Direction}.Call(f)
	}, status.Complete)
	job.Start()
	<-job.Finished.C

	log, err = DeleteEmptySourceFolders{Config: client.Config, Direction: job.Direction}.Call(*job)
	assert.NoError(err)
	assert.Equal(filepath.Dir(fi.RemotePath), log.Path)

	_, err = client.Find(files_sdk.FileFindParams{Path: lib.UrlJoinNoEscape("test-delete-source", "test.text")})
	require.NotNil(t, err)
	assert.Equal("Not Found - `Not Found.  This may be related to your permissions.`", err.Error())
	os.RemoveAll("test.text")
}

func TestClient_Downloader_Move_Source(t *testing.T) {
	client, r, err := CreateClient("TestClient_Downloader_Move_Source")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)

	folderClient := folder.Client{Config: client.Config}

	folderClient.Create(files_sdk.FolderCreateParams{Path: "test-move-source"})

	err = client.Upload(
		UploadWithSize(9),
		UploadWithDestinationPath(lib.UrlJoinNoEscape("test-move-source", "test.text")),
		UploadWithProvidedMtime(time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC)),
		UploadWithReader(strings.NewReader("testing 3")),
	)
	require.NoError(t, err)
	localPath, err := os.MkdirTemp("", "TestClient_Downloader_Move_Source")
	require.NoError(t, err)
	job := client.Downloader(
		DownloaderParams{RemotePath: "test-move-source", LocalPath: localPath},
	)
	var log status.Log
	job.RegisterFileEvent(func(f JobFile) {
		log, err = MoveSource{Config: client.Config, Direction: job.Direction, Path: "test-moved-source"}.Call(f)
	}, status.Complete)
	job.Start()
	job.Wait()

	require.NoError(t, err)
	assert.Equal("move source", log.Action)
	assert.Equal(filepath.Join("test-moved-source", "test.text"), log.Path)

	_, err = client.Find(files_sdk.FileFindParams{Path: lib.UrlJoinNoEscape("test-move-source", "test.text")})
	assert.Equal("Not Found - `Not Found`", err.Error())
	_, err = client.Find(files_sdk.FileFindParams{Path: lib.UrlJoinNoEscape("test-moved-source", "test.text")})
	assert.NoError(err)
	os.RemoveAll("test.text")
}

func TestClient_UploadFolder_Move_Source(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFolder_Move_Source")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)
	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	var log status.Log

	defer os.RemoveAll(tmpDir)
	err = os.MkdirAll(filepath.Join(tmpDir, "move-source"), 0755)
	assert.NoError(err)

	tempFile, err := os.Create(filepath.Join(tmpDir, "move-source", "upload-move-source.text"))
	assert.NoError(err)
	tempFile.Write([]byte("testing"))
	require.NoError(t, tempFile.Close())
	job := client.Uploader(UploaderParams{LocalPath: tempFile.Name()})
	job.RegisterFileEvent(func(f JobFile) {
		fmt.Println("RegisterFileEvent")
		log, err = MoveSource{Config: client.Config, Direction: job.Direction, Path: filepath.Join(tmpDir, "move-source", "test-moved-source.text")}.Call(f)
	}, status.Complete)
	job.Start()
	job.Wait()
	assert.False(job.Any(status.Errored), "A job errored")
	assert.NoError(err)
	assert.Equal("move source", log.Action)
	assert.Equal(filepath.Join(tmpDir, "move-source", "test-moved-source.text"), log.Path)
	stat, err := os.Stat(log.Path)
	require.NoError(t, err)
	assert.Equal("test-moved-source.text", stat.Name())
	tempFile.Close()
}

func TestClient_UploadFolder_Move_Source_Missing_Dir(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFolder_Move_Source_Missing_Dir")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)

	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	var log status.Log
	err = os.MkdirAll(filepath.Join(tmpDir, "move-source-dir"), 0755)
	assert.NoError(err)
	tempFile, err := os.Create(filepath.Join(tmpDir, "move-source-dir", "upload-move-source.text"))
	assert.NoError(err)
	tempFile.Write([]byte("testing"))
	tempFile.Close()
	job := client.Uploader(UploaderParams{LocalPath: filepath.Join(tmpDir, "move-source-dir") + string(os.PathSeparator)})
	job.RegisterFileEvent(func(f JobFile) {
		log, err = MoveSource{Config: client.Config, Direction: job.Direction, Path: filepath.Join(tmpDir, "moved-source-dir")}.Call(f)
	}, status.Complete)
	job.Start()
	job.Wait()
	erroredFile, ok := job.Find(status.Errored)
	if ok {
		assert.NoError(erroredFile.Err(), erroredFile.LocalPath())
	}
	require.Equal(t, false, job.Any(status.Errored), "")
	assert.NoError(err)
	assert.Equal("move source", log.Action)
	assert.Equal(filepath.Join(tmpDir, "moved-source-dir", "upload-move-source.text"), log.Path)
	stat, err := os.Stat(log.Path)
	assert.NoError(err)
	assert.Equal("upload-move-source.text", stat.Name())
	assert.Falsef(stat.IsDir(), "%s is a directory.", log.Path)

	_, err = os.Stat(filepath.Join(tmpDir, "move-source-dir", "upload-move-source.text"))
	assert.Truef(os.IsNotExist(err), "%s/%s was not deleted", "move-source-dir", "upload-move-source.text")
}

func TestClient_Downloader_Move_Source_Missing_Dir(t *testing.T) {
	client, r, err := CreateClient("TestClient_Downloader_Move_Source_Missing_Dir")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)

	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	buildScenario("TestClient_Downloader_Move_Source_Missing_Dir", client)
	os.MkdirAll(filepath.Join(tmpDir, "TestClient_Downloader_Move_Source_Missing_Dir"), 0755)

	logChan := make(chan status.Log)
	errChan := make(chan error)
	job := client.Downloader(
		DownloaderParams{
			RemotePath: "TestClient_Downloader_Move_Source_Missing_Dir",
			LocalPath:  filepath.Join(tmpDir, "TestClient_Downloader_Move_Source_Missing_Dir") + string(os.PathSeparator),
		},
	)
	job.RegisterFileEvent(func(f JobFile) {
		moveLog, err := MoveSource{Config: client.Config, Direction: job.Direction, Path: "TestClient_Downloader_Move_Source_Missing_Dir-moved"}.Call(f)
		logChan <- moveLog
		errChan <- err
	}, status.Complete)
	job.Start()

	fileValues := []string{
		"TestClient_Downloader_Move_Source_Missing_Dir-moved/nested_1/nested_2/3.text",
		"TestClient_Downloader_Move_Source_Missing_Dir-moved/nested_1/nested_2/nested_3/4.text",
	}
	files := make(map[string]string)
	files[fileValues[0]] = fileValues[0]
	files[fileValues[1]] = fileValues[1]

	log1 := <-logChan
	assert.Equal("move source", log1.Action)
	assert.Equal(files[lib.Path{Path: log1.Path}.NormalizePathSystemForAPI().String()], lib.Path{Path: log1.Path}.NormalizePathSystemForAPI().String())
	assert.NoError(<-errChan)

	log2 := <-logChan
	assert.Equal("move source", log2.Action)
	assert.Equal(files[lib.Path{Path: log2.Path}.NormalizePathSystemForAPI().String()], lib.Path{Path: log2.Path}.NormalizePathSystemForAPI().String())
	assert.NoError(<-errChan)

	job.Wait()

	assert.False(job.Any(status.Errored), "A job errored")

	movedDir, err := client.Find(files_sdk.FileFindParams{Path: "TestClient_Downloader_Move_Source_Missing_Dir-moved"})
	assert.NoError(err)
	assert.Equal("TestClient_Downloader_Move_Source_Missing_Dir-moved", movedDir.Path)
	assert.Equal("directory", movedDir.Type)
}

func TestClient_UploadFile_Delete_Source(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFile_Delete_Source")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)
	tmpDir, err := os.MkdirTemp(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	var log status.Log
	tempFile, err := os.Create(filepath.Join(tmpDir, "upload-delete-source.text"))
	assert.NoError(err)
	tempFile.Write([]byte("testing"))
	require.NoError(t, tempFile.Close())
	var fi JobFile
	job := client.Uploader(UploaderParams{LocalPath: tempFile.Name()})
	job.RegisterFileEvent(func(f JobFile) {
		fi = f
		log, err = DeleteSource{Config: client.Config, Direction: job.Direction}.Call(f)
	}, status.Complete)
	job.Start()
	job.Wait()
	assert.False(job.Any(status.Errored), "A job errored")
	assert.NoError(err)
	assert.Equal("delete source", log.Action)
	assert.Equal(fi.LocalPath, log.Path)
	tempFile.Close()
	os.Remove(tempFile.Name())
}

func TestClient_Uploader_Files(t *testing.T) {
	client, r, err := CreateClient("TestClient_Uploader_Files")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)
	tmpDir := t.TempDir()

	filesAndStatus := []FileStatus{
		{name: "1 (1).text", status: "complete", size: 24},
		{name: "2.text", status: "complete", size: 24},
		{name: "3.pdf", status: "ignored"},
	}
	filePaths := createFileStructure(t, tmpDir, filesAndStatus)

	job := client.Uploader(UploaderParams{LocalPath: tmpDir + string(os.PathSeparator), LocalPaths: filePaths})
	job.Start()
	job.Wait()

	assert.Len(job.Statuses, 2)
	assert.Equal(filePaths[0], job.Statuses[0].LocalPath())
	assert.Equal(status.Complete, job.Statuses[0].Status())
	assert.NoError(job.Statuses[0].Err())

	assert.Equal(filePaths[1], job.Statuses[1].LocalPath())
	assert.Equal(status.Complete, job.Statuses[1].Status())
	assert.NoError(job.Statuses[1].Err())
}

func TestClient_Uploader_File_Remote_Folder(t *testing.T) {
	client, r, err := CreateClient("TestClient_Uploader_File_Remote_Folder")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)
	tmpDir := t.TempDir()

	fileName := "file.text"
	dirName := "remote-folder"

	folderClient := folder.Client{Config: client.Config}
	_, err = folderClient.Create(files_sdk.FolderCreateParams{Path: dirName})
	defer deletePath(client, dirName)

	filesAndStatus := []FileStatus{{name: fileName, status: "complete", size: 24}}
	createFileStructure(t, tmpDir, filesAndStatus)

	job := client.Uploader(UploaderParams{LocalPath: filepath.Join(tmpDir, fileName), RemotePath: dirName})
	job.Start()
	job.Wait()

	assert.Len(job.Statuses, 1)
	assert.Equal(filepath.Join(dirName, fileName), job.Statuses[0].RemotePath())
	assert.Equal(status.Complete, job.Statuses[0].Status())
	assert.NoError(job.Statuses[0].Err())
}

func TestClient_Uploader_File_Remote_Path_Error(t *testing.T) {
	client, r, err := CreateClient("TestClient_Uploader_File_Remote_Path_Error")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)
	tmpDir := t.TempDir()

	fileName := "file.text"
	dirName := "remote-folder"

	filesAndStatus := []FileStatus{{name: fileName, status: "complete", size: 24}}
	createFileStructure(t, tmpDir, filesAndStatus)

	job := client.Uploader(UploaderParams{
		LocalPath:   filepath.Join(tmpDir, fileName),
		RemotePath:  dirName,
		RetryPolicy: RetryPolicy{Type: RetryUnfinished, RetryCount: 2},
	})
	job.Start()
	job.Wait()

	assert.Len(job.Statuses, 1)
	assert.Equal(status.Errored, job.Statuses[0].Status())
	assert.ErrorContains(job.Statuses[0].Err(), "token provided could not")
}

func TestClient_Uploader_Directories(t *testing.T) {
	client, r, err := CreateClient("TestClient_Uploader_Directories")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)
	tmpDir := t.TempDir()

	filesAndStatus := []FileStatus{
		{name: "A/1.text", status: "complete", size: 24},
		{name: "B", status: "complete", dir: true},
		{name: "B/2.text", status: "complete", size: 24},
		{name: "B/Y", status: "complete", dir: true},
		{name: "B/Z", status: "complete", dir: true},
		{name: "B/Z/4.text", status: "complete", size: 24},
		{name: "C", status: "complete", dir: true},
		{name: "3.text", status: "complete", size: 24},
	}
	createFileStructure(t, tmpDir, filesAndStatus)

	job := client.Uploader(
		UploaderParams{
			LocalPath: tmpDir + string(os.PathSeparator),
			LocalPaths: []string{
				filepath.Join(tmpDir, "A") + string(os.PathSeparator),
				filepath.Join(tmpDir, "B"),
				filepath.Join(tmpDir, "C"),
				filepath.Join(tmpDir, "3.text"),
			},
		},
	)
	job.Start()
	job.Wait()

	require.Equal(t, len(filesAndStatus), len(job.Statuses), "the right number of files did not upload")

	assert.Equal("1.text", job.Statuses[0].RemotePath())
	assert.Equal("B", job.Statuses[1].RemotePath())
	assert.Equal("B/2.text", job.Statuses[2].RemotePath())
	assert.Equal("B/Y", job.Statuses[3].RemotePath())
	assert.Equal("B/Z", job.Statuses[4].RemotePath())
	assert.Equal("B/Z/4.text", job.Statuses[5].RemotePath())
	assert.Equal("C", job.Statuses[6].RemotePath())
	assert.Equal("3.text", job.Statuses[7].RemotePath())

	for index, file := range filesAndStatus {
		assert.Equal(file.name, job.Statuses[index].LocalPath())
		expectedStatus := status.Complete
		if file.dir {
			expectedStatus = status.FolderCreated
		}
		assert.Equal(expectedStatus, job.Statuses[index].Status(), "file %s", file.name)
		assert.NoError(job.Statuses[index].Err(), "file %s", file.name)
	}
}

func TestClient_Uploader_Empty_Directory(t *testing.T) {
	client, r, err := CreateClient("TestClient_Uploader_Empty_Directory")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)
	tmpDir := t.TempDir()

	emptyDir := filepath.Join(tmpDir, "empty_dir")
	require.NoError(t, os.Mkdir(emptyDir, 0750))

	job := client.Uploader(
		UploaderParams{
			LocalPath: emptyDir,
		},
	)
	job.Start()
	job.Wait()

	require.Equal(t, 1, len(job.Statuses), "the right number of files did not upload")

	assert.Equal("empty_dir", job.Statuses[0].RemotePath())
	assert.Equal(emptyDir, job.Statuses[0].LocalPath())
	assert.Equal(status.FolderCreated, job.Statuses[0].Status())
	assert.NoError(job.Statuses[0].Err())
}

// This starts with the following directory structure in tmp
// - empty_dir/
// - A/
// - - 1.text
// - - empty_dir/
// - B/
// - - 2.text
// - - Z/
// - - - 4.text
// - - - empty_dir/
//
// A and B should be completely removed because there was a file uploaded from those hierarchies
// empty_dir/ does get removed because there were no files in that directory and it's part of the source path.
func TestClient_Uploader_Directories_With_DeleteSource(t *testing.T) {
	client, r, err := CreateClient("TestClient_Uploader_Directories_With_DeleteSource")
	if err != nil {
		t.Fatal(err)
	}

	defer r.Stop()
	assert := assert.New(t)
	tmpDir := t.TempDir()

	require.NoError(t, os.Mkdir(filepath.Join(tmpDir, "empty_dir"), 0750))

	filesAndStatus := []FileStatus{
		{name: "A/1.text", status: "complete", size: 24},
		{name: "A/empty_dir", status: "complete", dir: true},
		{name: "B/2.text", status: "complete", size: 24},
		{name: "B/Z", status: "complete", dir: true},
		{name: "B/Z/4.text", status: "complete", size: 24},
		{name: "B/Z/empty_dir", status: "complete", dir: true},
	}
	createFileStructure(t, tmpDir, filesAndStatus)

	job := client.Uploader(
		UploaderParams{
			LocalPath: tmpDir + string(os.PathSeparator),
			LocalPaths: []string{
				filepath.Join(tmpDir, "A") + string(os.PathSeparator),
				filepath.Join(tmpDir, "B") + string(os.PathSeparator),
			},
		},
	)

	job.RegisterFileEvent(func(f JobFile) {
		DeleteSource{Config: client.Config, Direction: job.Direction}.Call(f)
	}, status.Complete)

	job.Start()
	job.Wait()

	require.Equal(t, 6, len(job.Statuses), "the right number of files did not upload")

	assert.Equal("1.text", job.Statuses[0].RemotePath())
	assert.Equal("empty_dir", job.Statuses[1].RemotePath())
	assert.Equal("2.text", job.Statuses[2].RemotePath())
	assert.Equal("Z", job.Statuses[3].RemotePath())
	assert.Equal("Z/4.text", job.Statuses[4].RemotePath())
	assert.Equal("Z/empty_dir", job.Statuses[5].RemotePath())

	for index, file := range filesAndStatus {
		assert.Equal(file.name, job.Statuses[index].LocalPath())
		expectedStatus := status.Complete
		if file.dir {
			expectedStatus = status.FolderCreated
		}
		assert.Equal(expectedStatus, job.Statuses[index].Status(), "file %s", file.name)
		assert.NoError(job.Statuses[index].Err(), "file %s", file.name)
	}

	DeleteEmptySourceFolders{Config: client.Config, Direction: job.Direction}.Call(*job)

	for _, path := range []string{filepath.Join(tmpDir, "empty_dir"), filepath.Join("A", "empty_dir"), filepath.Join("A", "1.text"), filepath.Join("B", "1.text"), filepath.Join("B", "Z", "1.text"), filepath.Join("B", "Z", "empty_dir"), filepath.Join("B", "Z")} {
		_, dirErr := os.Stat(filepath.Join(tmpDir, path))
		assert.Truef(os.IsNotExist(dirErr), "%s should not exist", path)
	}
}

func TestClient_ListForRecursive(t *testing.T) {
	type recursiveRequest struct {
		path      string
		recursive string
		cursor    string
		fileType  string
	}
	requests := make(chan recursiveRequest, 10)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cursor := r.URL.Query().Get("cursor")
		requests <- recursiveRequest{
			path:      r.URL.Path,
			recursive: r.URL.Query().Get("recursive"),
			cursor:    cursor,
			fileType:  r.URL.Query().Get("type"),
		}
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/rest/v1/file_actions/metadata/reports" {
			_, _ = w.Write([]byte(`{"path":"Reports","type":"directory"}`))
			return
		}
		if r.URL.Path != "/api/rest/v1/folders/Reports" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.URL.Query().Get("type") == "file" {
			_, _ = w.Write([]byte(`[{"path":"Reports/first.txt","type":"file"}]`))
			return
		}
		switch cursor {
		case "":
			w.Header().Set("X-Files-Cursor", "cursor-1")
			_, _ = w.Write([]byte(`[{"path":"Reports/first.txt","type":"file"}]`))
		case "cursor-1":
			w.Header().Set("X-Files-Cursor", "cursor-2")
			_, _ = w.Write([]byte(`[]`))
		case "cursor-2":
			_, _ = w.Write([]byte(`[{"path":"Reports/nested","type":"directory"},{"path":"Reports/nested/last.txt","type":"file"}]`))
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer server.Close()

	config := files_sdk.Config{}.Init()
	config.EndpointOverride = server.URL
	client := &Client{Config: config}
	it, err := client.ListForRecursive(files_sdk.FolderListForParams{Path: "reports"})
	require.NoError(t, err)

	var paths []string
	for it.Next() {
		item, ok := it.Current().(RecursiveItem)
		require.True(t, ok)
		assert.Equal(t, item, it.Resource())
		assert.NoError(t, item.Err())
		paths = append(paths, item.Path)
	}

	require.NoError(t, it.Err())
	assert.Equal(t, []string{"Reports", "Reports/first.txt", "Reports/nested", "Reports/nested/last.txt"}, paths)

	resumed, err := client.ListForRecursive(files_sdk.FolderListForParams{
		Path: "reports",
		ListParams: files_sdk.ListParams{
			Cursor: "cursor-2",
		},
	})
	require.NoError(t, err)
	var resumedPaths []string
	for resumed.Next() {
		resumedPaths = append(resumedPaths, resumed.Resource().Path)
	}
	require.NoError(t, resumed.Err())
	assert.Equal(t, []string{"Reports/nested", "Reports/nested/last.txt"}, resumedPaths)

	filesOnly, err := client.ListForRecursive(files_sdk.FolderListForParams{Path: "reports", Type: "file"})
	require.NoError(t, err)
	var filePaths []string
	for filesOnly.Next() {
		filePaths = append(filePaths, filesOnly.Resource().Path)
	}
	require.NoError(t, filesOnly.Err())
	assert.Equal(t, []string{"Reports/first.txt"}, filePaths)

	require.Equal(t, 8, len(requests))
	for i, expected := range []recursiveRequest{
		{path: "/api/rest/v1/file_actions/metadata/reports"},
		{path: "/api/rest/v1/folders/Reports"},
		{path: "/api/rest/v1/folders/Reports", cursor: "cursor-1"},
		{path: "/api/rest/v1/folders/Reports", cursor: "cursor-2"},
		{path: "/api/rest/v1/file_actions/metadata/reports"},
		{path: "/api/rest/v1/folders/Reports", cursor: "cursor-2"},
		{path: "/api/rest/v1/file_actions/metadata/reports"},
		{path: "/api/rest/v1/folders/Reports", fileType: "file"},
	} {
		request := <-requests
		assert.Equal(t, expected.path, request.path, "request %d", i+1)
		assert.Equal(t, expected.cursor, request.cursor, "request %d", i+1)
		assert.Equal(t, expected.fileType, request.fileType, "request %d", i+1)
		if strings.Contains(request.path, "/file_actions/metadata/") {
			assert.Empty(t, request.recursive, "request %d", i+1)
		} else {
			assert.Equal(t, "true", request.recursive, "request %d", i+1)
		}
	}
}

func TestClient_ListForRecursiveRoot(t *testing.T) {
	requests := make(chan string, 3)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests <- r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path != "/api/rest/v1/folders/" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = w.Write([]byte(`[{"path":"first.txt","type":"file"}]`))
	}))
	defer server.Close()

	config := files_sdk.Config{}.Init()
	config.EndpointOverride = server.URL
	client := &Client{Config: config}
	for _, path := range []string{"", "/", "."} {
		it, err := client.ListForRecursive(files_sdk.FolderListForParams{Path: path})
		require.NoError(t, err)
		require.True(t, it.Next())
		assert.Equal(t, "first.txt", it.Resource().Path)
		assert.False(t, it.Next())
		require.NoError(t, it.Err())
	}

	require.Equal(t, 3, len(requests))
	assert.Equal(t, "/api/rest/v1/folders/", <-requests)
	assert.Equal(t, "/api/rest/v1/folders/", <-requests)
	assert.Equal(t, "/api/rest/v1/folders/", <-requests)
}

func TestClient_ListForRecursiveRejectsUnsupportedParams(t *testing.T) {
	modifiedAt := time.Now()
	for name, params := range map[string]files_sdk.FolderListForParams{
		"sort_by":                    {SortBy: "size"},
		"search":                     {Search: "report"},
		"search_all":                 {SearchAll: lib.Bool(true)},
		"search_custom_metadata_key": {SearchCustomMetadataKey: "category"},
		"modified_at_datetime":       {ModifiedAtDatetime: &modifiedAt},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := (&Client{}).ListForRecursive(params)
			assert.EqualError(t, err, "recursive listings do not support sort_by, search, search_all, search_custom_metadata_key, or modified_at_datetime")
		})
	}
}

func TestClient_UploadFile(t *testing.T) {
	client := remoteFSSimulatorClient(t)
	fileName := filepath.Join(t.TempDir(), "anything")
	file, err := os.Create(fileName)
	require.NoError(t, err)
	_, err = file.Write([]byte("anything"))
	require.NoError(t, err)
	require.NoError(t, file.Close())

	err = client.UploadFile(fileName, "test/anything")
	require.NoError(t, err)

	sdkFile, err := client.Find(files_sdk.FileFindParams{Path: "test/anything"})
	require.NoError(t, err)
	require.Equal(t, sdkFile.DisplayName, "anything")
	fs := (&FS{}).Init(client.Config, false).WithContext(context.Background()).(fs2.FS)
	fsFile, err := fs.Open("test/anything")
	require.NoError(t, err)
	defer fsFile.Close()
	fileBytes, err := io.ReadAll(fsFile)
	require.NoError(t, err)
	require.Equal(t, "anything", string(fileBytes))
}

func TestClient_ListFor(t *testing.T) {
	client, r, err := CreateClient("TestClient_ListFor")
	if err != nil {
		t.Fatal(err)
	}

	buildScenario("TestClient_DownloadFolder", client)

	assert := assert.New(t)

	it, err := client.ListFor(files_sdk.FolderListForParams{
		ListParams: files_sdk.ListParams{PerPage: 1},
		Path:       "TestClient_DownloadFolder/nested_1/nested_2",
	})

	assert.NoError(err)
	var files []files_sdk.File
	for it.Next() {
		if it.File().Type == "directory" {
			subIt, _ := it.Iterate(it.File().Identifier())
			for subIt.Next() {
				subFile := subIt.Current().(files_sdk.File)
				files = append(files, subFile)
				loadedFile, err := it.LoadResource(subFile.Identifier())
				assert.NoError(err)
				assert.Equal(subFile, loadedFile, "LoadResource with Identifier matches file from list")
			}
		}
		files = append(files, it.File())

	}
	paths := lo.Map[files_sdk.File, string](files, func(item files_sdk.File, index int) string {
		return item.Path
	})
	assert.Equal([]string{
		"TestClient_DownloadFolder/nested_1/nested_2/nested_3/4.text",
		"TestClient_DownloadFolder/nested_1/nested_2/nested_3",
		"TestClient_DownloadFolder/nested_1/nested_2/3.text",
	}, paths)

	r.Stop()
}

func TestUploadsCloseErroredResponsesOnce(t *testing.T) {
	for _, version := range []string{"v1", "v2"} {
		for _, mode := range []string{"redirect-rejected", "cancel-after-response", "plain-status-error"} {
			t.Run(version+"/"+mode, func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				var bodies []*ownedHTTPResponseBody
				config := files_sdk.Config{DisableDirectTransfers: true}.Init().SetCustomClient(&http.Client{
					Transport: ownedHTTPResponseTransport(func(req *http.Request) (*http.Response, error) {
						if req.Body != nil {
							_, _ = io.Copy(io.Discard, req.Body)
							_ = req.Body.Close()
						}
						body := &ownedHTTPResponseBody{Reader: strings.NewReader("response")}
						bodies = append(bodies, body)
						status := http.StatusOK
						header := http.Header{}
						if mode == "redirect-rejected" {
							status = http.StatusTemporaryRedirect
							header.Set("Location", "https://upload.test/again")
						} else if mode == "cancel-after-response" {
							cancel()
						} else {
							status = http.StatusForbidden
							header.Set("Content-Type", "text/plain")
						}
						return &http.Response{StatusCode: status, Status: fmt.Sprintf("%d %s", status, http.StatusText(status)), Header: header, Body: body, ContentLength: int64(body.Len()), Request: req}, nil
					}),
					CheckRedirect: func(*http.Request, []*http.Request) error { return errors.New("stopped after 10 redirects") },
				})
				config.RetryMax = 0
				config.FeatureFlags[files_sdk.FeatureFlagUploadV2ChecksumTrailer] = false
				upload := files_sdk.FileUploadPart{PartNumber: 1, Path: "file", UploadUri: "https://upload.test/file", HttpMethod: "PUT"}
				reader := &ProxyReaderAt{ReaderAt: strings.NewReader("a"), len: 1}
				u := &uploadIO{Client: &Client{Config: config}}
				var err error
				if version == "v1" {
					_, err = u.uploadPart(ctx, &Part{FileUploadPart: upload, ProxyReader: reader})
				} else {
					e := &uploadV2Engine{u: u}
					_, _, _, _, err = e.uploadPart(ctx, &uploadV2Part{uploadV2PartDescriptor: uploadV2PartDescriptor{number: 1, upload: upload}, reader: reader})
				}
				require.Error(t, err)
				require.Len(t, bodies, 1)
				if mode == "cancel-after-response" {
					require.ErrorIs(t, err, context.Canceled)
				}
				assert.Equal(t, 1, bodies[0].closes)
			})
		}
	}
}

func TestDownloadErrorsCloseResponseOnce(t *testing.T) {
	for _, operation := range []string{"stream", "range", "v2-range"} {
		for _, status := range []int{http.StatusNotFound, http.StatusServiceUnavailable} {
			t.Run(fmt.Sprintf("%s/%d", operation, status), func(t *testing.T) {
				body := &ownedHTTPResponseBody{Reader: strings.NewReader("response")}
				config := ownedResponseConfig(body, status, http.Header{"Content-Type": []string{"text/plain"}})
				remote := files_sdk.File{Path: "file", DownloadUri: "https://download.test/file"}
				reader := (&File{File: &remote, FS: (&FS{Context: context.Background()}).Init(config, false)}).Init()
				var err error
				switch operation {
				case "stream":
					_, err = reader.readCloserInit()
				case "range":
					_, err = reader.ReaderRange(0, 2)
				case "v2-range":
					_, _, err = reader.downloadV2ReaderRange(context.Background(), 0, 2, false)
				}
				require.Error(t, err)
				var responseErr lib.ResponseError
				require.ErrorAs(t, err, &responseErr)
				require.Equal(t, status, responseErr.StatusCode)
				assert.Equal(t, 1, body.closes)
			})
		}
	}
}

func TestGetReleasesFinalResponseWhenPolicyObservesCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	body := &ownedHTTPResponseBody{Reader: strings.NewReader(`{"error":"unavailable","type":"internal_server_error","http-code":"500"}`)}
	config := files_sdk.Config{DisableDirectTransfers: true}.Init().SetCustomClient(&http.Client{Transport: ownedHTTPResponseTransport(func(req *http.Request) (*http.Response, error) {
		if req.Body != nil {
			_, _ = io.Copy(io.Discard, req.Body)
			_ = req.Body.Close()
		}
		cancel()
		return &http.Response{StatusCode: 500, Status: "500 Internal Server Error", Header: http.Header{"Content-Type": []string{"application/json"}}, Body: body, ContentLength: int64(body.Len()), Request: req}, nil
	})})
	config.RetryMax = 0
	_, err := (&Client{Config: config}).Get("file", files_sdk.WithContext(ctx))
	require.ErrorIs(t, err, context.Canceled)
	assert.Equal(t, 1, body.closes)
}

func TestZipMetadataReleasesFinalResponseOnPolicyCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	body := &ownedHTTPResponseBody{Reader: strings.NewReader("response")}
	config := files_sdk.Config{EndpointOverride: "https://files.test"}.Init().SetCustomClient(&http.Client{Transport: ownedHTTPResponseTransport(func(req *http.Request) (*http.Response, error) {
		if req.Body != nil {
			_, _ = io.Copy(io.Discard, req.Body)
			_ = req.Body.Close()
		}
		cancel()
		return &http.Response{StatusCode: http.StatusInternalServerError, Header: http.Header{}, Body: body, Request: req}, nil
	})})
	config.RetryMax = 0
	batch := &zipBatchDownloader{ctx: ctx, job: &Job{Config: config}}
	download, err := batch.createZipDownloadOnce(nil)
	require.ErrorIs(t, err, context.Canceled)
	assert.Empty(t, download.uri)
	assert.Equal(t, 1, body.closes)
}

// zipRegionMismatch is the API's typed lockout region mismatch naming eu.files.com.
const zipRegionMismatch = `{"type":"not-authenticated/lockout-region-mismatch","http-code":403,"error":"Lockout region mismatch","data":{"host":"eu.files.com"}}`

// zipRegionSent is one request a zipRegionSites transport received, with the headers it carried.
type zipRegionSent struct {
	host, method, path              string
	key, session, reauth, workspace string
	custom                          string
	body                            string
}

func zipRegionSentOf(req *http.Request, body []byte) zipRegionSent {
	return zipRegionSent{
		host:      req.URL.Host,
		method:    req.Method,
		path:      req.URL.Path,
		key:       req.Header.Get("X-FilesAPI-Key"),
		session:   req.Header.Get("X-FilesAPI-Auth"),
		reauth:    req.Header.Get("X-Files-Reauthentication"),
		workspace: req.Header.Get("X-Files-Workspace-Id"),
		custom:    req.Header.Get("X-Custom"),
		body:      string(body),
	}
}

// zipRegionSites answers each request with the handler for its host, as separate sites would, without changing the
// request the SDK sent, and records each request.
type zipRegionSites struct {
	mu       sync.Mutex
	handlers map[string]func(*http.Request) *http.Response
	sent     []zipRegionSent
}

func (s *zipRegionSites) RoundTrip(req *http.Request) (*http.Response, error) {
	var body []byte
	if req.Body != nil {
		body, _ = io.ReadAll(req.Body)
		_ = req.Body.Close()
	}
	s.mu.Lock()
	s.sent = append(s.sent, zipRegionSentOf(req, body))
	handler, ok := s.handlers[req.URL.Host]
	s.mu.Unlock()
	if !ok {
		return nil, fmt.Errorf("no site at %s", req.URL.Host)
	}
	return handler(req), nil
}

func (s *zipRegionSites) sentTo(host, method string) []zipRegionSent {
	s.mu.Lock()
	defer s.mu.Unlock()
	var sent []zipRegionSent
	for _, request := range s.sent {
		if request.host == host && request.method == method {
			sent = append(sent, request)
		}
	}
	return sent
}

// zipRegionAnswer answers with status, contentType (none when empty), payload and header.
func zipRegionAnswer(status int, contentType, payload string, header http.Header) func(*http.Request) *http.Response {
	return func(req *http.Request) *http.Response {
		answer := http.Header{}
		for name, values := range header {
			answer[name] = append([]string(nil), values...)
		}
		if contentType != "" {
			answer.Set("Content-Type", contentType)
		}
		return &http.Response{StatusCode: status, Status: fmt.Sprintf("%d %s", status, http.StatusText(status)), Header: answer, Body: io.NopCloser(strings.NewReader(payload)), ContentLength: int64(len(payload)), Request: req}
	}
}

// An upload part is a request to storage: a typed lockout region mismatch, even
// with a retry budget left, never sends it to the host it names.
func TestUploadPartsAreNeverMovedByARegionMismatch(t *testing.T) {
	for _, version := range []string{"v1", "v2"} {
		t.Run(version, func(t *testing.T) {
			var bodies []*ownedHTTPResponseBody
			var hosts []string
			config := files_sdk.Config{DisableDirectTransfers: true}.Init().SetCustomClient(&http.Client{Transport: ownedHTTPResponseTransport(func(req *http.Request) (*http.Response, error) {
				if req.Body != nil {
					_, _ = io.Copy(io.Discard, req.Body)
					_ = req.Body.Close()
				}
				hosts = append(hosts, req.URL.Host)
				body := &ownedHTTPResponseBody{Reader: strings.NewReader(zipRegionMismatch)}
				bodies = append(bodies, body)
				header := http.Header{"Content-Type": []string{"application/json"}, "X-Files-Host": []string{"eu.files.com"}}
				return &http.Response{StatusCode: http.StatusForbidden, Status: "403 Forbidden", Header: header, Body: body, ContentLength: int64(body.Len()), Request: req}, nil
			})})
			// A retry is left, so that only the kind of request keeps it from moving.
			config.RetryMax = 2
			config.FeatureFlags[files_sdk.FeatureFlagUploadV2ChecksumTrailer] = false
			upload := files_sdk.FileUploadPart{PartNumber: 1, Path: "file", UploadUri: "https://us.files.com/upload/file", HttpMethod: "PUT"}
			reader := &ProxyReaderAt{ReaderAt: strings.NewReader("a"), len: 1}
			u := &uploadIO{Client: &Client{Config: config}}
			var err error
			if version == "v1" {
				_, err = u.uploadPart(context.Background(), &Part{FileUploadPart: upload, ProxyReader: reader})
			} else {
				e := &uploadV2Engine{u: u}
				_, _, _, _, err = e.uploadPart(context.Background(), &uploadV2Part{uploadV2PartDescriptor: uploadV2PartDescriptor{number: 1, upload: upload}, reader: reader})
			}
			require.Error(t, err)
			require.NotEmpty(t, hosts)
			assert.NotContains(t, hosts, "eu.files.com", "an upload part is never sent to the host a mismatch names")
			var mismatch *files_sdk.RegionMismatchError
			assert.False(t, errors.As(err, &mismatch), "an upload part is not routed as an API request")
			var responseErr files_sdk.ResponseError
			require.ErrorAs(t, err, &responseErr)
			assert.Equal(t, "not-authenticated/lockout-region-mismatch", responseErr.Type)
			for _, body := range bodies {
				assert.Equal(t, 1, body.closes)
			}
		})
	}
}

// The ZIP create request is an API request: a lockout region mismatch moves it,
// with the same body and credentials, and its result keeps the answer's origin.
func TestZipCreateFollowsARegionMismatchAndKeepsItsAnswer(t *testing.T) {
	sites := &zipRegionSites{handlers: map[string]func(*http.Request) *http.Response{
		"us.files.com": zipRegionAnswer(http.StatusForbidden, "application/json", zipRegionMismatch, nil),
		"eu.files.com": zipRegionAnswer(http.StatusOK, "application/json", `{"download_uri":"/zip/download?z=1"}`, nil),
	}}
	config := files_sdk.Config{APIKey: "zip-key", EndpointOverride: "https://us.files.com"}.Init().SetCustomClient(&http.Client{Transport: sites})
	batch := &zipBatchDownloader{ctx: context.Background(), job: &Job{Config: config}}

	download, err := batch.createZipDownloadOnce([]*DownloadStatus{{remotePath: "batch/a.txt"}})

	require.NoError(t, err)
	assert.Equal(t, "/zip/download?z=1", download.uri)
	assert.True(t, download.api.Relocated())
	assert.Equal(t, "https://eu.files.com", download.api.Origin())
	require.Len(t, sites.sent, 2)
	assert.Equal(t, "us.files.com", sites.sent[0].host)
	assert.Equal(t, "eu.files.com", sites.sent[1].host)
	for _, sent := range sites.sent {
		assert.Equal(t, http.MethodPost, sent.method)
		assert.Equal(t, "/api/rest/v1/zip_downloads", sent.path)
		assert.Equal(t, "zip-key", sent.key)
		assert.JSONEq(t, `{"paths":[],"encoded_paths":["batch%2Fa.txt"]}`, sent.body)
	}
	assert.Equal(t, "https://us.files.com", batch.job.Config.Endpoint(), "the job's configuration is not changed")
}

// A ZIP create mismatch that is not followed, here for want of a retry, is
// returned with the API's error reachable, and the response is closed once.
func TestZipCreateRegionMismatchThatIsNotFollowedReturnsItWithTheAPIError(t *testing.T) {
	body := &ownedHTTPResponseBody{Reader: strings.NewReader(zipRegionMismatch)}
	var hosts []string
	config := files_sdk.Config{APIKey: "zip-key", EndpointOverride: "https://us.files.com"}.Init().SetCustomClient(&http.Client{Transport: ownedHTTPResponseTransport(func(req *http.Request) (*http.Response, error) {
		if req.Body != nil {
			_, _ = io.Copy(io.Discard, req.Body)
			_ = req.Body.Close()
		}
		hosts = append(hosts, req.URL.Host)
		return &http.Response{StatusCode: http.StatusForbidden, Status: "403 Forbidden", Header: http.Header{"Content-Type": []string{"application/json"}}, Body: body, Request: req}, nil
	})})
	config.RetryMax = 0
	batch := &zipBatchDownloader{ctx: context.Background(), job: &Job{Config: config}}

	download, err := batch.createZipDownloadOnce(nil)

	var mismatch *files_sdk.RegionMismatchError
	require.ErrorAs(t, err, &mismatch)
	assert.Equal(t, files_sdk.RegionMismatchRetryBudget, mismatch.Reason)
	assert.Equal(t, "https://eu.files.com", mismatch.Origin)
	assert.Equal(t, "https://us.files.com", mismatch.From)
	var responseErr files_sdk.ResponseError
	require.ErrorAs(t, err, &responseErr)
	assert.Equal(t, "eu.files.com", responseErr.Data.Host)
	assert.Empty(t, download.uri)
	assert.False(t, download.api.Relocated())
	assert.Equal(t, []string{"us.files.com"}, hosts)
	assert.Equal(t, 1, body.closes)
}

// A relative ZIP download URI from a moved create request is fetched with the
// Files.com authentication the moved request sent, including none, whatever the
// environment or the configuration says by then.
func TestZipRelativeStreamSendsTheAuthenticationTheMovedCreateSent(t *testing.T) {
	for name, test := range map[string]struct {
		config   files_sdk.Config
		env      string
		laterEnv string
		later    func(files_sdk.Config) files_sdk.Config
		want     zipRegionSent
	}{
		"an API key from the environment": {env: "create-key", laterEnv: "stream-key", want: zipRegionSent{key: "create-key"}},
		"a session, a workspace and additional headers": {
			config:   files_sdk.Config{SessionId: "create-session", WorkspaceId: 12, AdditionalHeaders: map[string]string{"X-Files-Reauthentication": "password:create", "X-Custom": "kept"}},
			laterEnv: "stream-key",
			later: func(config files_sdk.Config) files_sdk.Config {
				config.SessionId = "later-session"
				config.WorkspaceId = 34
				config.AdditionalHeaders = map[string]string{"X-Files-Reauthentication": "password:later", "X-Custom": "kept"}
				return config
			},
			want: zipRegionSent{session: "create-session", workspace: "12", reauth: "password:create", custom: "kept"},
		},
		"no authentication": {laterEnv: "late-key", want: zipRegionSent{}},
	} {
		t.Run(name, func(t *testing.T) {
			t.Setenv("FILES_API_KEY", test.env)
			sites := &zipRegionSites{handlers: map[string]func(*http.Request) *http.Response{
				"us.files.com": zipRegionAnswer(http.StatusForbidden, "application/json", zipRegionMismatch, nil),
				"eu.files.com": func(req *http.Request) *http.Response {
					if req.Method == http.MethodPost {
						return zipRegionAnswer(http.StatusOK, "application/json", `{"download_uri":"/zip/download?z=1"}`, nil)(req)
					}
					return zipRegionAnswer(http.StatusOK, "application/zip", "zip-bytes", nil)(req)
				},
			}}
			config := test.config
			config.EndpointOverride = "https://us.files.com"
			config = config.Init().SetCustomClient(&http.Client{Transport: sites})
			batch := &zipBatchDownloader{ctx: context.Background(), job: &Job{Config: config}}
			download, err := batch.createZipDownloadOnce(nil)
			require.NoError(t, err)

			// The environment and the configuration change before the archive is fetched.
			t.Setenv("FILES_API_KEY", test.laterEnv)
			later := config
			if test.later != nil {
				later = test.later(config)
			}
			streamURI, streamConfig, err := zipBatchStreamTarget(later, download)
			require.NoError(t, err)
			assert.Equal(t, "https://eu.files.com/zip/download?z=1", streamURI)
			zipFile := newZipBatchFile(streamConfig, context.Background(), streamURI)
			data, err := io.ReadAll(zipFile)
			require.NoError(t, err)
			require.NoError(t, zipFile.Close())
			assert.Equal(t, "zip-bytes", string(data))

			creates := sites.sentTo("eu.files.com", http.MethodPost)
			streams := sites.sentTo("eu.files.com", http.MethodGet)
			require.Len(t, creates, 1)
			require.Len(t, streams, 1)
			for _, sent := range []zipRegionSent{creates[0], streams[0]} {
				assert.Equal(t, test.want.key, sent.key)
				assert.Equal(t, test.want.session, sent.session)
				assert.Equal(t, test.want.workspace, sent.workspace)
				assert.Equal(t, test.want.reauth, sent.reauth)
				assert.Equal(t, test.want.custom, sent.custom)
			}
			assert.Equal(t, "https://us.files.com", config.Endpoint(), "the job's configuration is not changed")
		})
	}
}

// A moved ZIP create's relocated context belongs to the answer the call returned: the ErrorHandler's replacement for the
// moved request's own answer keeps it, and its relative download URI is fetched from that origin with the
// authentication the moved create sent; an unrelated response the handler selected, or the moved answer once its
// request leads elsewhere, has none, and its relative URI is fetched as an unmoved answer's is, with the configuration.
func TestZipCreateRelocatedContextBelongsToTheSelectedAnswer(t *testing.T) {
	for name, test := range map[string]struct {
		selected  func(answer *http.Response) *http.Response
		relocated bool
		stream    string
		host      string
		key       string
	}{
		"a replacement for the moved answer": {
			selected: func(answer *http.Response) *http.Response {
				return zipRegionAnswer(http.StatusOK, "application/json", `{"download_uri":"/zip/download?z=b"}`, nil)(answer.Request)
			},
			relocated: true,
			stream:    "https://eu.files.com/zip/download?z=b",
			host:      "eu.files.com",
			key:       "create-key",
		},
		"an unrelated response": {
			selected: func(*http.Response) *http.Response {
				return zipRegionAnswer(http.StatusOK, "application/json", `{"download_uri":"/zip/download?z=c"}`, nil)(nil)
			},
			stream: "https://us.files.com/zip/download?z=c",
			host:   "us.files.com",
			key:    "later-key",
		},
		"the moved answer once its request leads elsewhere": {
			selected: func(answer *http.Response) *http.Response {
				answer.Request.URL.Host = "c.example.com"
				return zipRegionAnswer(http.StatusOK, "application/json", `{"download_uri":"/zip/download?z=c"}`, nil)(answer.Request)
			},
			stream: "https://us.files.com/zip/download?z=c",
			host:   "us.files.com",
			key:    "later-key",
		},
	} {
		t.Run(name, func(t *testing.T) {
			archive := zipRegionAnswer(http.StatusOK, "application/zip", "zip-bytes", nil)
			sites := &zipRegionSites{handlers: map[string]func(*http.Request) *http.Response{
				"us.files.com": func(req *http.Request) *http.Response {
					if req.Method == http.MethodPost {
						return zipRegionAnswer(http.StatusForbidden, "application/json", zipRegionMismatch, nil)(req)
					}
					return archive(req)
				},
				"eu.files.com": func(req *http.Request) *http.Response {
					if req.Method == http.MethodPost {
						return zipRegionAnswer(http.StatusServiceUnavailable, "application/json", `{"error":"unavailable","http-code":503,"type":"service-unavailable"}`, nil)(req)
					}
					return archive(req)
				},
			}}
			config := files_sdk.Config{APIKey: "create-key", EndpointOverride: "https://us.files.com"}.Init().SetCustomClient(&http.Client{Transport: sites})
			// The move takes the one retry, so the moved create's 503 reaches the ErrorHandler.
			config.RetryMax = 1
			config.Client.ErrorHandler = func(answer *http.Response, _ error, _ int) (*http.Response, error) {
				if err := answer.Body.Close(); err != nil {
					return nil, err
				}
				return test.selected(answer), nil
			}
			batch := &zipBatchDownloader{ctx: context.Background(), job: &Job{Config: config}}

			download, err := batch.createZipDownloadOnce(nil)

			require.NoError(t, err)
			assert.Equal(t, test.relocated, download.api.Relocated())
			require.Len(t, sites.sentTo("us.files.com", http.MethodPost), 1)
			require.Len(t, sites.sentTo("eu.files.com", http.MethodPost), 1)

			// The configuration's key changes before the archive is fetched.
			later := config
			later.APIKey = "later-key"
			streamURI, streamConfig, err := zipBatchStreamTarget(later, download)
			require.NoError(t, err)
			assert.Equal(t, test.stream, streamURI)
			zipFile := newZipBatchFile(streamConfig, context.Background(), streamURI)
			data, err := io.ReadAll(zipFile)
			require.NoError(t, err)
			require.NoError(t, zipFile.Close())
			assert.Equal(t, "zip-bytes", string(data))
			streams := sites.sentTo(test.host, http.MethodGet)
			require.Len(t, streams, 1)
			assert.Equal(t, test.key, streams[0].key)
		})
	}
}
