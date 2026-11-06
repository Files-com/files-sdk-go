package file

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestNativeFSOpenOwnsIndependentReadersAndTerminalClose(t *testing.T) {
	var bodies []*ownedHTTPResponseBody
	config := files_sdk.Config{EndpointOverride: "https://api.test", DisableDirectTransfers: true}.Init().SetCustomClient(&http.Client{Transport: ownedHTTPResponseTransport(func(request *http.Request) (*http.Response, error) {
		if request.URL.Path == "/content" {
			body := &ownedHTTPResponseBody{Reader: strings.NewReader("abc")}
			bodies = append(bodies, body)
			return &http.Response{StatusCode: http.StatusOK, Header: http.Header{}, Body: body, ContentLength: 3, Request: request}, nil
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"path":"data.bin","display_name":"data.bin","type":"file","size":3,"download_uri":"https://storage.test/content"}`)),
			Request:    request,
		}, nil
	})})
	config.RetryMax = 0
	native := (&FS{}).Init(config, true).Native()
	first, err := native.Open("data.bin")
	require.NoError(t, err)
	second, err := native.Open("data.bin")
	require.NoError(t, err)

	for _, handle := range []fs.File{first, second} {
		buffer := make([]byte, 1)
		_, err := io.ReadFull(handle, buffer)
		require.NoError(t, err)
		require.Equal(t, "a", string(buffer))
	}
	for _, handle := range []fs.File{first, second} {
		remaining, err := io.ReadAll(handle)
		require.NoError(t, err)
		require.Equal(t, "bc", string(remaining))
		_, err = handle.Read(make([]byte, 1))
		require.ErrorIs(t, err, io.EOF)
	}
	for _, handle := range []fs.File{first, second} {
		require.NoError(t, handle.Close())
		_, err = handle.Read(make([]byte, 1))
		require.ErrorIs(t, err, fs.ErrClosed)
		_, err = handle.Stat()
		require.ErrorIs(t, err, fs.ErrClosed)
		require.ErrorIs(t, handle.Close(), fs.ErrClosed)
	}
	require.Len(t, bodies, 2)
	for _, body := range bodies {
		require.Equal(t, 1, body.closes)
	}
}

func TestNativeFSUsesOneRootForListingsAndMutations(t *testing.T) {
	server := (&MockAPIServer{T: t}).Do()
	t.Cleanup(server.Shutdown)
	server.MockFiles["scope"] = mockFile{File: files_sdk.File{Path: "scope", DisplayName: "scope", Type: "directory"}}
	for _, name := range []string{"a", "z"} {
		server.MockFiles["scope/"+name] = mockFile{File: files_sdk.File{Path: "scope/" + name, DisplayName: name, Type: "file"}}
	}
	server.MockRoute("/api/rest/v1/folders/scope", func(c *gin.Context, _ interface{}) bool {
		if c.Request.Method != http.MethodGet {
			return false
		}
		c.JSON(http.StatusOK, []files_sdk.File{
			{Path: "scope/z", DisplayName: "z", Type: "file"},
			{Path: "scope/a", DisplayName: "a", Type: "file"},
		})
		return true
	})
	server.MockRoute("/api/rest/v1/folders/scope/newdir", func(_ *gin.Context, _ interface{}) bool {
		server.MockFiles["scope/newdir"] = mockFile{File: files_sdk.File{Path: "scope/newdir", DisplayName: "newdir", Type: "directory"}}
		return false
	})
	var deleted []string
	server.GetRouter().DELETE("/api/rest/v1/files/*path", func(c *gin.Context) {
		deleted = append(deleted, c.Param("path"))
		c.Status(http.StatusNoContent)
	})
	requests := recordCreateRequests(server, "scope/created.txt")
	legacy := (&FS{Root: "scope"}).Init(server.Client().Config, true)
	native := legacy.WithContext(context.Background()).(*FS).Native()

	listed, err := native.ReadDir(".")
	require.NoError(t, err)
	require.Len(t, listed, 2)
	require.Equal(t, "a", listed[0].Name())
	require.Equal(t, "z", listed[1].Name())
	opened, err := native.Open(".")
	require.NoError(t, err)
	directory := opened.(fs.ReadDirFile)
	copied := directory
	for i, handle := range []fs.ReadDirFile{directory, copied} {
		entries, err := handle.ReadDir(1)
		require.NoError(t, err)
		require.Len(t, entries, 1)
		require.Equal(t, []string{"z", "a"}[i], entries[0].Name())
	}
	_, err = directory.ReadDir(1)
	require.ErrorIs(t, err, io.EOF)
	other, err := native.Open(".")
	require.NoError(t, err)
	entries, err := other.(fs.ReadDirFile).ReadDir(1)
	require.NoError(t, err)
	require.Equal(t, "z", entries[0].Name())
	require.NoError(t, other.Close())
	require.NoError(t, directory.Close())
	_, err = copied.ReadDir(1)
	require.ErrorIs(t, err, fs.ErrClosed)
	_, err = directory.Read(make([]byte, 1))
	require.ErrorIs(t, err, fs.ErrClosed)

	file, err := native.Open("z")
	require.NoError(t, err)
	info, err := file.Stat()
	require.NoError(t, err)
	require.Equal(t, "scope/z", info.Sys().(files_sdk.File).Path)
	require.NoError(t, file.Close())
	writer, err := native.Create("created.txt")
	require.NoError(t, err)
	_, err = writer.Write([]byte("payload"))
	require.NoError(t, err)
	require.Zero(t, requests.seen().begins)
	require.NoError(t, writer.Close())
	require.Equal(t, "payload", string(requests.seen().content))
	require.NoError(t, native.MkdirAll("newdir", 0750))
	require.NoError(t, native.Remove("a"))
	require.NoError(t, native.RemoveAll("sub"))
	require.Equal(t, []string{"/scope/a", "/scope/sub"}, deleted)
}

func TestNativeFSRejectsInvalidNamesBeforeRequestsOrStaging(t *testing.T) {
	var requests atomic.Int32
	config := files_sdk.Config{}.Init().SetCustomClient(&http.Client{Transport: ownedHTTPResponseTransport(func(*http.Request) (*http.Response, error) {
		requests.Add(1)
		return nil, errors.New("unexpected request")
	})})
	native := (&FS{Root: "scope"}).Init(config, true).Native()
	staging := t.TempDir()
	stageCreateFilesIn(t, staging)
	for _, name := range []string{"", "/absolute", "..", "../outside", "a/../b", "a//b", "a/"} {
		t.Run(name, func(t *testing.T) {
			_, err := native.Open(name)
			require.ErrorIs(t, err, fs.ErrInvalid)
			_, err = native.ReadDir(name)
			require.ErrorIs(t, err, fs.ErrInvalid)
			_, err = native.Create(name)
			require.ErrorIs(t, err, fs.ErrInvalid)
			require.ErrorIs(t, native.MkdirAll(name, 0750), fs.ErrInvalid)
			require.ErrorIs(t, native.Remove(name), fs.ErrInvalid)
			require.ErrorIs(t, native.RemoveAll(name), fs.ErrInvalid)
			require.ErrorIs(t, native.Chtimes(name, time.Time{}, time.Time{}), fs.ErrInvalid)
		})
	}
	_, err := native.MkdirTemp("parent", "job-*")
	require.ErrorIs(t, err, errors.ErrUnsupported)
	require.Zero(t, requests.Load())
	require.Empty(t, stagingEntries(t, staging))
}

func TestNativeFSChtimesPreservesZeroTimesAndRejectsAccessTimeChanges(t *testing.T) {
	var mutations []string
	var providedMtime time.Time
	var lookups int
	config := files_sdk.Config{EndpointOverride: "https://api.test"}.Init().SetCustomClient(&http.Client{Transport: ownedHTTPResponseTransport(func(request *http.Request) (*http.Response, error) {
		status := http.StatusOK
		body := `{"path":"scope/file","display_name":"file","type":"file"}`
		if request.Method == http.MethodGet {
			lookups++
			if strings.HasSuffix(request.URL.Path, "/missing") {
				status = http.StatusNotFound
				encoded, err := json.Marshal(files_sdk.ResponseError{Type: string(files_sdk.ErrFileNotFound), ErrorMessage: "file not found", HttpCode: status})
				require.NoError(t, err)
				body = string(encoded)
			}
		} else {
			mutations = append(mutations, request.URL.Path)
			var params struct {
				ProvidedMtime time.Time `json:"provided_mtime"`
			}
			require.NoError(t, json.NewDecoder(request.Body).Decode(&params))
			providedMtime = params.ProvidedMtime
		}
		return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), ContentLength: int64(len(body)), Request: request}, nil
	})})
	config.RetryMax = 0
	native := (&FS{Root: "scope"}).Init(config, true).Native()
	mtime := time.Date(2026, 11, 4, 12, 0, 0, 0, time.UTC)

	require.NoError(t, native.Chtimes("file", time.Time{}, mtime))
	require.Equal(t, mtime, providedMtime)
	require.Equal(t, []string{"/api/rest/v1/files/scope/file"}, mutations)
	require.ErrorIs(t, native.Chtimes("file", mtime, mtime), errors.ErrUnsupported)
	require.Len(t, mutations, 1)
	require.Zero(t, lookups)
	require.NoError(t, native.Chtimes("file", time.Time{}, time.Time{}))
	require.Equal(t, 1, lookups)
	err := native.Chtimes("missing", time.Time{}, time.Time{})
	require.ErrorIs(t, err, fs.ErrNotExist)
	require.True(t, files_sdk.IsNotExist(err), "the SDK's structured not-found cause is retained")
	require.Len(t, mutations, 1, "zero times never mutate metadata")
}

func TestNativeFSMkdirAllRejectsFilesAndRemoveAllOnlyIgnoresNotFound(t *testing.T) {
	config := files_sdk.Config{EndpointOverride: "https://api.test"}.Init().SetCustomClient(&http.Client{Transport: ownedHTTPResponseTransport(func(request *http.Request) (*http.Response, error) {
		status := http.StatusOK
		body := `{"path":"scope/file","display_name":"file","type":"file"}`
		var apiErr *files_sdk.ResponseError
		if request.Method == http.MethodPost {
			status = http.StatusUnprocessableEntity
			apiErr = &files_sdk.ResponseError{Type: string(files_sdk.ErrDestinationExists), ErrorMessage: "destination exists", HttpCode: status}
		} else if request.Method == http.MethodDelete {
			status = http.StatusForbidden
			apiErr = &files_sdk.ResponseError{Type: "permission-denied", ErrorMessage: "permission denied", HttpCode: status}
			if strings.HasSuffix(request.URL.Path, "/missing") {
				status = http.StatusNotFound
				apiErr = &files_sdk.ResponseError{Type: string(files_sdk.ErrFileNotFound), ErrorMessage: "file not found", HttpCode: status}
			}
		}
		if apiErr != nil {
			encoded, err := json.Marshal(apiErr)
			require.NoError(t, err)
			body = string(encoded)
		}
		return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), ContentLength: int64(len(body)), Request: request}, nil
	})})
	config.RetryMax = 0
	native := (&FS{Root: "scope"}).Init(config, true).Native()

	err := native.MkdirAll("file", 0750)
	require.ErrorIs(t, err, fs.ErrExist)
	require.ErrorIs(t, err, files_sdk.ErrDestinationExists)
	require.NoError(t, native.RemoveAll("missing"))
	err = native.RemoveAll("denied")
	require.True(t, files_sdk.IsErrorType(err, files_sdk.ResponseErrorType("permission-denied")))
}

type nativeHeldBody struct {
	ctx     context.Context
	entered chan struct{}
	once    sync.Once
	closes  atomic.Int32
}

func (b *nativeHeldBody) Read([]byte) (int, error) {
	b.once.Do(func() { close(b.entered) })
	<-b.ctx.Done()
	return 0, b.ctx.Err()
}

func (b *nativeHeldBody) Close() error {
	b.closes.Add(1)
	return nil
}

func TestNativeFSCloseInterruptsAnOwnedBlockedRead(t *testing.T) {
	entered := make(chan struct{})
	var body *nativeHeldBody
	config := files_sdk.Config{EndpointOverride: "https://api.test", DisableDirectTransfers: true}.Init().SetCustomClient(&http.Client{Transport: ownedHTTPResponseTransport(func(request *http.Request) (*http.Response, error) {
		if request.URL.Path == "/content" {
			body = &nativeHeldBody{ctx: request.Context(), entered: entered}
			return &http.Response{StatusCode: http.StatusOK, Header: http.Header{}, Body: body, ContentLength: 1, Request: request}, nil
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"path":"held","display_name":"held","type":"file","size":1,"download_uri":"https://storage.test/content"}`)),
			Request:    request,
		}, nil
	})})
	config.RetryMax = 0
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	handle, err := (&FS{}).Init(config, true).Native().WithContext(ctx).Open("held")
	require.NoError(t, err)
	readResult := make(chan error, 1)
	go func() {
		_, err := handle.Read(make([]byte, 1))
		readResult <- err
	}()
	select {
	case <-entered:
	case <-time.After(10 * time.Second):
		t.Fatal("the owned body was never read")
	}
	closed := make(chan error, 1)
	go func() { closed <- handle.Close() }()
	select {
	case err := <-closed:
		require.NoError(t, err)
	case <-time.After(10 * time.Second):
		t.Fatal("Close did not interrupt the owned read")
	}
	select {
	case err := <-readResult:
		require.ErrorIs(t, err, fs.ErrClosed)
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(10 * time.Second):
		t.Fatal("the cancelled read did not return")
	}
	_, err = handle.Read(make([]byte, 1))
	require.ErrorIs(t, err, fs.ErrClosed)
	require.EqualValues(t, 1, body.closes.Load())
}

func TestNativeFSNestedRootCreationAndDotRemovalStayInTheirNamespace(t *testing.T) {
	var created []string
	var lookups, deletes int
	rootExists, childExists := true, true
	config := files_sdk.Config{EndpointOverride: "https://api.test"}.Init().SetCustomClient(&http.Client{Transport: ownedHTTPResponseTransport(func(request *http.Request) (*http.Response, error) {
		status := http.StatusCreated
		body := ""
		switch request.Method {
		case http.MethodPost:
			created = append(created, request.URL.Path)
			path := strings.TrimPrefix(request.URL.Path, "/api/rest/v1/folders/")
			if !strings.HasPrefix(path, "outer/scope/") {
				status = http.StatusForbidden
				body = `{"type":"permission-denied","error":"ancestor creation denied","http-code":403}`
			} else {
				encoded, err := json.Marshal(files_sdk.File{Path: path, Type: "directory"})
				require.NoError(t, err)
				body = string(encoded)
			}
		case http.MethodGet:
			lookups++
			status = http.StatusOK
			body = `{"path":"outer/scope","display_name":"scope","type":"directory"}`
			if !strings.HasSuffix(request.URL.Path, "/outer/scope") {
				status = http.StatusNotFound
				encoded, err := json.Marshal(files_sdk.ResponseError{Type: string(files_sdk.ErrFileNotFound), ErrorMessage: "file not found", HttpCode: status})
				require.NoError(t, err)
				body = string(encoded)
			}
		case http.MethodDelete:
			deletes++
			rootExists, childExists = false, false
			status = http.StatusNoContent
		}
		return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), ContentLength: int64(len(body)), Request: request}, nil
	})})
	config.RetryMax = 0
	native := (&FS{Root: "outer/scope"}).Init(config, true).Native()

	require.NoError(t, native.MkdirAll("child/grandchild", 0750))
	require.Equal(t, []string{"/api/rest/v1/folders/outer/scope/child", "/api/rest/v1/folders/outer/scope/child/grandchild"}, created)
	require.Zero(t, lookups, "a successful directory result does not add a metadata-read requirement")
	require.NoError(t, native.MkdirAll(".", 0750))
	require.Len(t, created, 2, "root existence checks never create ancestors")
	require.Equal(t, 1, lookups)
	require.ErrorIs(t, native.Remove("."), fs.ErrInvalid)
	require.ErrorIs(t, native.RemoveAll("."), fs.ErrInvalid)
	require.Zero(t, deletes)
	require.True(t, rootExists)
	require.True(t, childExists)
	missing := (&FS{Root: "outer/missing"}).Init(config, false).Native()
	err := missing.MkdirAll(".", 0750)
	require.ErrorIs(t, err, fs.ErrNotExist)
	require.True(t, files_sdk.IsNotExist(err))
	require.Len(t, created, 2)
}

func TestNativeFSCloseDistinguishesEarlyClosureFromCompletedReadVerification(t *testing.T) {
	closeErr := errors.New("body close failed")
	for _, test := range []struct {
		name       string
		readAll    bool
		sent       int64
		closeErr   error
		wantErr    error
		wantStatus int
	}{
		{name: "partial read", sent: 3},
		{name: "partial read close failure", sent: 3, closeErr: closeErr, wantErr: closeErr},
		{name: "complete read", readAll: true, sent: 3, wantStatus: 1},
		{name: "complete read mismatch", readAll: true, sent: 4, wantErr: UntrustedSizeRangeRequestSizeSentReceived, wantStatus: 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			body := &ownedHTTPResponseBody{Reader: strings.NewReader("abc"), closeErr: test.closeErr}
			var statusRequests int
			config := files_sdk.Config{EndpointOverride: "https://api.test", DisableDirectTransfers: true}.Init().SetCustomClient(&http.Client{Transport: ownedHTTPResponseTransport(func(request *http.Request) (*http.Response, error) {
				if request.URL.Path == "/content" {
					return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"X-Files-Download-Request-Id": []string{"request"}}, Body: body, ContentLength: -1, Request: request}, nil
				}
				payload := `{"path":"file","display_name":"file","type":"file","size":3,"download_uri":"https://storage.test/content"}`
				if request.URL.Path == "/content/request" {
					statusRequests++
					encoded, err := json.Marshal(map[string]any{"data": map[string]any{"status": "completed", "bytes_transferred": test.sent}})
					require.NoError(t, err)
					payload = string(encoded)
				}
				return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(payload)), Request: request}, nil
			})})
			config.RetryMax = 0
			handle, err := (&FS{}).Init(config, true).Native().Open("file")
			require.NoError(t, err)
			if test.readAll {
				content, err := io.ReadAll(handle)
				require.NoError(t, err)
				require.Equal(t, "abc", string(content))
			} else {
				first := make([]byte, 1)
				_, err := io.ReadFull(handle, first)
				require.NoError(t, err)
				require.Equal(t, "a", string(first))
			}
			err = handle.Close()
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, test.wantStatus, statusRequests)
			require.Equal(t, 1, body.closes)
		})
	}
}

func TestNativeFSDownloaderCancellationReachesPendingDiscovery(t *testing.T) {
	entered := make(chan struct{})
	settled := make(chan error, 1)
	var requests atomic.Int32
	var first sync.Once
	config := files_sdk.Config{EndpointOverride: "https://api.test"}.Init().SetCustomClient(&http.Client{Transport: ownedHTTPResponseTransport(func(request *http.Request) (*http.Response, error) {
		requests.Add(1)
		first.Do(func() { close(entered) })
		<-request.Context().Done()
		err := request.Context().Err()
		settled <- err
		return nil, err
	})})
	config.RetryMax = 0
	sourceCtx, releaseSource := context.WithCancel(context.Background())
	defer releaseSource()
	callerCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	native := (&FS{Context: sourceCtx}).Init(config, true).Native()
	job := downloader(callerCtx, native, DownloaderParams{
		RemoteFile: files_sdk.File{Path: "pending", Type: "file"},
		LocalPath:  filepath.Join(t.TempDir(), "pending"),
		config:     config,
	})
	stopped := make(chan struct{})
	t.Cleanup(func() {
		cancel()
		releaseSource()
		if !job.Finished.Called() {
			job.Cancel()
		}
		select {
		case <-stopped:
		case <-time.After(10 * time.Second):
			t.Error("the downloader did not settle during cleanup")
		}
	})
	go func() {
		job.Start()
		job.Wait()
		close(stopped)
	}()
	select {
	case <-entered:
	case <-time.After(10 * time.Second):
		t.Fatal("the downloader never requested discovery metadata")
	}
	cancel()
	select {
	case err := <-settled:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(10 * time.Second):
		t.Fatal("caller cancellation did not reach the downloader's metadata request")
	}
	select {
	case <-stopped:
	case <-time.After(10 * time.Second):
		t.Fatal("the cancelled downloader did not settle")
	}
	require.True(t, job.Finished.Called())
	require.EqualValues(t, 1, requests.Load())
}
