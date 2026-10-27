package files_sdk_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"testing/iotest"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/behavior"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/site"
	"github.com/Files-com/files-sdk-go/v3/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type recordedRequest struct {
	method           string
	path             string
	contentType      string
	contentLength    int64
	transferEncoding []string
	body             []byte
	stagedFiles      int
}

// attachmentServer records every request and answers with the given
// statuses in turn, then 200. A 307 redirects to the same path.
type attachmentServer struct {
	*httptest.Server
	tempDir  string
	statuses []int
	// errorBody, if set, is the body of each answer that is not 2xx.
	errorBody string
	mu        sync.Mutex
	requests  []recordedRequest
}

func newAttachmentServer(t *testing.T, tempDir string, statuses ...int) *attachmentServer {
	t.Helper()
	server := &attachmentServer{tempDir: tempDir, statuses: statuses}
	server.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		staged, err := os.ReadDir(server.tempDir)
		assert.NoError(t, err)
		server.mu.Lock()
		server.requests = append(server.requests, recordedRequest{
			method:           r.Method,
			path:             r.URL.Path,
			contentType:      r.Header.Get("Content-Type"),
			contentLength:    r.ContentLength,
			transferEncoding: r.TransferEncoding,
			body:             body,
			stagedFiles:      len(staged),
		})
		status := http.StatusOK
		if len(server.requests) <= len(server.statuses) {
			status = server.statuses[len(server.requests)-1]
		}
		server.mu.Unlock()
		if status == http.StatusTemporaryRedirect {
			w.Header().Set("Location", r.URL.Path+"?redirected=true")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if status >= 300 && server.errorBody != "" {
			_, _ = w.Write([]byte(server.errorBody))
			return
		}
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)
	return server
}

func (s *attachmentServer) config() files_sdk.Config {
	return files_sdk.Config{APIKey: "attachment-test", EndpointOverride: s.URL}.Init()
}

func (s *attachmentServer) recorded() []recordedRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]recordedRequest(nil), s.requests...)
}

// isolateTempDir points the SDK's temporary files at a new, empty directory.
func isolateTempDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("TMPDIR", dir)
	t.Setenv("TMP", dir)
	t.Setenv("TEMP", dir)
	return dir
}

func assertNoFiles(t *testing.T, dir string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	assert.Empty(t, entries, "the SDK left temporary files behind")
}

// callerReader is a reader the caller owns. Like a pipe or a response body,
// it can only be read forward.
type callerReader struct {
	reader io.Reader
	read   int
	closed bool
	// afterRead, if set, runs after each read.
	afterRead func()
}

func (r *callerReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	r.read += n
	if r.afterRead != nil {
		r.afterRead()
	}
	return n, err
}

func (r *callerReader) Close() error {
	r.closed = true
	return nil
}

func newCallerReader(content []byte) *callerReader {
	return &callerReader{reader: bytes.NewReader(content)}
}

func jsonMembers(t *testing.T, body []byte) map[string]json.RawMessage {
	t.Helper()
	var members map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(body, &members), "the body is one JSON object")
	return members
}

type fileObject struct {
	content     []byte
	filename    string
	contentType string
}

func decodeFileObject(t *testing.T, raw json.RawMessage) fileObject {
	t.Helper()
	members := jsonMembers(t, raw)
	var keys []string
	for key := range members {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	require.Equal(t, []string{"encoded_content", "filename", "type"}, keys)
	var encoded string
	var object fileObject
	require.NoError(t, json.Unmarshal(members["encoded_content"], &encoded))
	require.NoError(t, json.Unmarshal(members["filename"], &object.filename))
	require.NoError(t, json.Unmarshal(members["type"], &object.contentType))
	content, err := base64.StdEncoding.DecodeString(encoded)
	require.NoError(t, err)
	object.content = content
	return object
}

func TestUserUpdateSendsAvatarAttachmentWithTheOtherParameters(t *testing.T) {
	avatar := []byte("\x89PNG\r\n\x1a\n avatar \x00\xff")
	for name, debug := range map[string]bool{"without debug logging": false, "with debug logging": true} {
		t.Run(name, func(t *testing.T) {
			tempDir := isolateTempDir(t)
			server := newAttachmentServer(t, tempDir)
			config := server.config()
			var logged bytes.Buffer
			config.Debug = debug
			config.Logger = log.New(&logged, "", 0)
			client := user.Client{Config: config}

			params := files_sdk.UserUpdateParams{Id: 7, Name: "Avatar Owner", Company: "Files.com", ReceiveAdminAlerts: lib.Bool(false)}
			_, err := client.Update(params)
			require.NoError(t, err)
			source := newCallerReader(avatar)
			// The filename is sent exactly as given, percent sign and backslash included.
			params.AvatarFileAttachment = lib.NewJSONAttachment(source, int64(len(avatar)), `50% off\avatar.png`, "image/png")
			_, err = client.Update(params)
			require.NoError(t, err)

			assert.False(t, source.closed, "the SDK does not close the caller's reader")
			requests := server.recorded()
			require.Len(t, requests, 2)
			sent := requests[1]
			assert.Equal(t, http.MethodPatch, sent.method)
			assert.Equal(t, "/api/rest/v1/users/7", sent.path)
			assert.Equal(t, "application/json", sent.contentType)
			assert.Equal(t, int64(len(sent.body)), sent.contentLength, "the request declares its real length")
			assert.Empty(t, sent.transferEncoding, "the body is not chunked")
			members := jsonMembers(t, sent.body)
			assert.Equal(t, fileObject{content: avatar, filename: `50% off\avatar.png`, contentType: "image/png"}, decodeFileObject(t, members["avatar_file"]))
			delete(members, "avatar_file")
			assert.Equal(t, jsonMembers(t, requests[0].body), members, "the other parameters are sent as they are without an attachment")
			assertNoFiles(t, tempDir)

			if debug {
				assert.Contains(t, logged.String(), "with attachments, not logged")
				assert.NotContains(t, logged.String(), base64.StdEncoding.EncodeToString(avatar))
				// Each call's response is logged once, the attachment call's after its request.
				assert.Equal(t, 2, strings.Count(logged.String(), "HTTP/1.1 200 OK"))
				assert.Greater(t, strings.LastIndex(logged.String(), "HTTP/1.1 200 OK"), strings.Index(logged.String(), "with attachments, not logged"),
					"the attachment call's response is logged")
			}
		})
	}
}

func TestBehaviorCreateSendsNestedValueWithAForwardOnlyAttachment(t *testing.T) {
	tempDir := isolateTempDir(t)
	server := newAttachmentServer(t, tempDir)
	watermark := bytes.Repeat([]byte("watermark "), 1000)
	source := newCallerReader(watermark)

	_, err := (&behavior.Client{Config: server.config()}).Create(files_sdk.BehaviorCreateParams{
		Path:     "reports",
		Behavior: "watermark",
		Value: map[string]interface{}{
			"none":         nil,
			"empty_object": map[string]interface{}{},
			"empty_list":   []interface{}{},
			"off":          false,
			"zero":         0,
			"empty_text":   "",
			"nested":       map[string]interface{}{"list": []interface{}{1, "two", nil, false}},
		},
		Recursive:                lib.Bool(false),
		AttachmentFileAttachment: lib.NewJSONAttachmentUpTo(source, 1<<20, "watermark.png", "image/png"),
	})
	require.NoError(t, err)

	requests := server.recorded()
	require.Len(t, requests, 1)
	members := jsonMembers(t, requests[0].body)
	assert.Equal(t, fileObject{content: watermark, filename: "watermark.png", contentType: "image/png"}, decodeFileObject(t, members["attachment_file"]))
	assert.JSONEq(t, `{"none":null,"empty_object":{},"empty_list":[],"off":false,"zero":0,"empty_text":"","nested":{"list":[1,"two",null,false]}}`, string(members["value"]))
	assert.JSONEq(t, `false`, string(members["recursive"]))
	assert.JSONEq(t, `"reports"`, string(members["path"]))
	assert.JSONEq(t, `"watermark"`, string(members["behavior"]))
	assert.Len(t, members, 5)
	assert.False(t, source.closed)
	assertNoFiles(t, tempDir)
}

func TestUpdateWithMapSendsAnAttachmentBesideAPreEncodedFile(t *testing.T) {
	tempDir := isolateTempDir(t)
	server := newAttachmentServer(t, tempDir)
	logo := []byte("logo bytes")
	preEncoded := map[string]interface{}{
		"encoded_content": base64.StdEncoding.EncodeToString([]byte("icon bytes")),
		"filename":        "icon.png",
		"type":            "image/png",
	}

	_, err := (&site.Client{Config: server.config()}).UpdateWithMap(map[string]interface{}{
		"logo_file":             lib.NewJSONAttachment(bytes.NewReader(logo), int64(len(logo)), "logo.png", "image/png"),
		"icon16_file":           preEncoded,
		"name":                  "",
		"welcome_custom_css":    nil,
		"session_expiry":        0,
		"allowed_ips":           []interface{}{},
		"disable_notifications": false,
		"custom_settings":       map[string]interface{}{},
	})
	require.NoError(t, err)

	requests := server.recorded()
	require.Len(t, requests, 1)
	members := jsonMembers(t, requests[0].body)
	assert.Equal(t, fileObject{content: logo, filename: "logo.png", contentType: "image/png"}, decodeFileObject(t, members["logo_file"]))
	delete(members, "logo_file")
	body, err := json.Marshal(members)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"icon16_file": {"encoded_content": "aWNvbiBieXRlcw==", "filename": "icon.png", "type": "image/png"},
		"name": "",
		"welcome_custom_css": null,
		"session_expiry": 0,
		"allowed_ips": [],
		"disable_notifications": false,
		"custom_settings": {}
	}`, string(body))
}

func TestAttachmentRetriesResendTheSameBody(t *testing.T) {
	tempDir := isolateTempDir(t)
	server := newAttachmentServer(t, tempDir, http.StatusServiceUnavailable)
	// Large enough to be kept in a temporary file rather than in memory.
	background := make([]byte, 3<<20)
	_, err := rand.Read(background)
	require.NoError(t, err)
	icon := []byte("icon bytes")

	_, err = (&site.Client{Config: server.config()}).Update(files_sdk.SiteUpdateParams{
		Name:                                   "Retried Site",
		LoginPageBackgroundImageFileAttachment: lib.NewJSONAttachmentUpTo(newCallerReader(background), 4<<20, "background.jpg", "image/jpeg"),
		Icon16FileAttachment:                   lib.NewJSONAttachment(bytes.NewReader(icon), int64(len(icon)), "icon16.png", "image/png"),
	})
	require.NoError(t, err)

	requests := server.recorded()
	require.Len(t, requests, 2, "the 503 was retried")
	assert.True(t, bytes.Equal(requests[0].body, requests[1].body), "the retry sent the same bytes")
	for _, request := range requests {
		assert.Equal(t, int64(len(request.body)), request.contentLength)
		assert.Equal(t, 1, request.stagedFiles, "the large attachment was kept in one temporary file")
	}
	members := jsonMembers(t, requests[1].body)
	assert.True(t, bytes.Equal(background, decodeFileObject(t, members["login_page_background_image_file"]).content))
	assert.Equal(t, fileObject{content: icon, filename: "icon16.png", contentType: "image/png"}, decodeFileObject(t, members["icon16_file"]))
	assert.JSONEq(t, `"Retried Site"`, string(members["name"]))
	assertNoFiles(t, tempDir)
}

func TestRequestOptionCanChangeTheBodyOfARequestWithAnAttachment(t *testing.T) {
	avatar := bytes.Repeat([]byte("avatar "), 200_000) // kept in a temporary file
	readBody := func(t *testing.T, req *http.Request) []byte {
		body, err := io.ReadAll(req.Body)
		require.NoError(t, err)
		return body
	}
	// renameTo replaces the body as a caller that rewrites requests can,
	// setting Body, GetBody and ContentLength together.
	renameTo := func(name string) func(*testing.T, *http.Request) {
		return func(t *testing.T, req *http.Request) {
			body := bytes.Replace(readBody(t, req), []byte(`"name":"Before"`), []byte(`"name":"`+name+`"`), 1)
			req.Body = io.NopCloser(bytes.NewReader(body))
			req.GetBody = func() (io.ReadCloser, error) {
				return io.NopCloser(bytes.NewReader(body)), nil
			}
			req.ContentLength = int64(len(body))
		}
	}
	tests := map[string]struct {
		option   func(*testing.T, *http.Request)
		wantName string
	}{
		"an option that only sets a header": {
			option:   func(_ *testing.T, req *http.Request) { req.Header.Set("X-Request-Tag", "tagged") },
			wantName: "Before",
		},
		"an option that reads the body and restores it": {
			option: func(t *testing.T, req *http.Request) {
				readBody(t, req)
				body, err := req.GetBody()
				require.NoError(t, err)
				req.Body = body
			},
			wantName: "Before",
		},
		"an option that replaces the body with one of the same length": {option: renameTo("Edited"), wantName: "Edited"},
		"an option that replaces the body with a longer one":           {option: renameTo("Edited again"), wantName: "Edited again"},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			for variant, withAttachment := range map[string]bool{"with an attachment": true, "without an attachment": false} {
				t.Run(variant, func(t *testing.T) {
					tempDir := isolateTempDir(t)
					server := newAttachmentServer(t, tempDir, http.StatusServiceUnavailable, http.StatusTemporaryRedirect)
					params := files_sdk.UserUpdateParams{Id: 7, Name: "Before"}
					if withAttachment {
						params.AvatarFileAttachment = lib.NewJSONAttachment(bytes.NewReader(avatar), int64(len(avatar)), "avatar.png", "image/png")
					}

					_, err := (&user.Client{Config: server.config()}).Update(params, files_sdk.RequestOption(func(req *http.Request) error {
						test.option(t, req)
						return nil
					}))
					require.NoError(t, err)

					requests := server.recorded()
					require.Len(t, requests, 3, "the 503 was retried and the 307 followed")
					for _, request := range requests {
						assert.True(t, bytes.Equal(requests[0].body, request.body), "every attempt sent the same bytes")
						assert.Equal(t, int64(len(request.body)), request.contentLength, "the request declares its real length")
					}
					members := jsonMembers(t, requests[0].body)
					var sentName string
					require.NoError(t, json.Unmarshal(members["name"], &sentName))
					assert.Equal(t, test.wantName, sentName)
					if withAttachment {
						assert.True(t, bytes.Equal(avatar, decodeFileObject(t, members["avatar_file"]).content))
						assert.NotZero(t, requests[2].stagedFiles, "the SDK kept its prepared copy until the call returned")
					}
					assertNoFiles(t, tempDir)
				})
			}
		})
	}
}

// attachmentRegionMismatch is the API's typed lockout region mismatch naming eu.files.com.
const attachmentRegionMismatch = `{"type":"not-authenticated/lockout-region-mismatch","http-code":403,"error":"Lockout region mismatch","data":{"host":"eu.files.com"}}`

// A lockout region mismatch moves a request with an attachment to the host the API
// names. The relocated request sends the same bytes, with their real length, from
// the copy the SDK keeps, including a body a RequestOption replaced, and the
// copies are removed once the call returns.
func TestAttachmentRegionMismatchResendsTheSameBodyToTheNamedHost(t *testing.T) {
	avatar := bytes.Repeat([]byte("avatar "), 200_000) // kept in a temporary file
	for name, rename := range map[string]bool{"the prepared body": false, "a body a RequestOption replaced": true} {
		t.Run(name, func(t *testing.T) {
			tempDir := isolateTempDir(t)
			type sent struct {
				host, method, path, key string
				body                    []byte
				declared                int64
				stagedFiles             int
			}
			var requests []sent
			config := files_sdk.Config{APIKey: "attachment-test", EndpointOverride: "https://us.files.com"}.Init()
			config = config.SetCustomClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				body, err := io.ReadAll(req.Body)
				if closeErr := req.Body.Close(); err == nil {
					err = closeErr
				}
				if err != nil {
					return nil, err
				}
				staged, err := os.ReadDir(tempDir)
				if err != nil {
					return nil, err
				}
				requests = append(requests, sent{host: req.URL.Host, method: req.Method, path: req.URL.Path, key: req.Header.Get("X-FilesAPI-Key"), body: body, declared: req.ContentLength, stagedFiles: len(staged)})
				if req.URL.Host == "us.files.com" {
					return &http.Response{StatusCode: http.StatusForbidden, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(attachmentRegionMismatch)), Request: req}, nil
				}
				return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(`{}`)), Request: req}, nil
			})})
			var opts []files_sdk.RequestResponseOption
			if rename {
				opts = append(opts, files_sdk.RequestOption(func(req *http.Request) error {
					body, err := io.ReadAll(req.Body)
					if err != nil {
						return err
					}
					body = bytes.Replace(body, []byte(`"name":"Before"`), []byte(`"name":"Edited"`), 1)
					req.Body = io.NopCloser(bytes.NewReader(body))
					req.GetBody = func() (io.ReadCloser, error) {
						return io.NopCloser(bytes.NewReader(body)), nil
					}
					req.ContentLength = int64(len(body))
					return nil
				}))
			}
			source := newCallerReader(avatar)

			_, err := (&user.Client{Config: config}).Update(files_sdk.UserUpdateParams{
				Id:                   7,
				Name:                 "Before",
				AvatarFileAttachment: lib.NewJSONAttachment(source, int64(len(avatar)), "avatar.png", "image/png"),
			}, opts...)
			require.NoError(t, err)

			require.Len(t, requests, 2, "the mismatch was followed once")
			assert.Equal(t, "us.files.com", requests[0].host)
			assert.Equal(t, "eu.files.com", requests[1].host)
			for _, request := range requests {
				assert.Equal(t, http.MethodPatch, request.method)
				assert.Equal(t, "/api/rest/v1/users/7", request.path)
				assert.Equal(t, "attachment-test", request.key, "the relocated request keeps its credentials")
				assert.Equal(t, int64(len(request.body)), request.declared, "the request declares its real length")
				assert.NotZero(t, request.stagedFiles, "the SDK kept its copy until the call returned")
			}
			assert.True(t, bytes.Equal(requests[0].body, requests[1].body), "the relocated request sent the same bytes")
			members := jsonMembers(t, requests[1].body)
			var sentName string
			require.NoError(t, json.Unmarshal(members["name"], &sentName))
			wantName := "Before"
			if rename {
				wantName = "Edited"
			}
			assert.Equal(t, wantName, sentName)
			assert.True(t, bytes.Equal(avatar, decodeFileObject(t, members["avatar_file"]).content))
			assert.False(t, source.closed, "the SDK does not close the caller's reader")
			assertNoFiles(t, tempDir)
		})
	}
}

// A lockout region mismatch that is not followed, here for want of a retry,
// fails a request with an attachment with the mismatch described and the API's
// error still reachable, and removes the prepared copy.
func TestAttachmentRegionMismatchThatIsNotFollowedKeepsItsErrorAndRemovesTheCopy(t *testing.T) {
	tempDir := isolateTempDir(t)
	server := newAttachmentServer(t, tempDir, http.StatusForbidden)
	server.errorBody = attachmentRegionMismatch
	config := server.config()
	config.RetryMax = 0
	avatar := bytes.Repeat([]byte("avatar "), 200_000) // kept in a temporary file

	_, err := (&user.Client{Config: config}).Update(files_sdk.UserUpdateParams{
		Id:                   7,
		AvatarFileAttachment: lib.NewJSONAttachment(bytes.NewReader(avatar), int64(len(avatar)), "avatar.png", "image/png"),
	})

	var mismatch *files_sdk.RegionMismatchError
	require.ErrorAs(t, err, &mismatch)
	assert.Equal(t, files_sdk.RegionMismatchRetryBudget, mismatch.Reason)
	assert.Equal(t, "eu.files.com", mismatch.Host)
	assert.Equal(t, "https://eu.files.com", mismatch.Origin)
	assert.Equal(t, server.URL, mismatch.From)
	assert.True(t, mismatch.Approved)
	var responseError files_sdk.ResponseError
	require.ErrorAs(t, err, &responseError)
	assert.Equal(t, "not-authenticated/lockout-region-mismatch", responseError.Type)
	assert.Equal(t, responseError.Error(), err.Error(), "the call's own error message is kept")
	requests := server.recorded()
	require.Len(t, requests, 1, "nothing is sent to the named host")
	assert.Equal(t, 1, requests[0].stagedFiles)
	assertNoFiles(t, tempDir)
}

// zeroReader is a finite source of zeros that allocates nothing to read.
type zeroReader struct{ left int64 }

func (z *zeroReader) Read(p []byte) (int, error) {
	if z.left == 0 {
		return 0, io.EOF
	}
	n := min(int64(len(p)), z.left)
	clear(p[:n])
	z.left -= n
	return int(n), nil
}

// tapBody observes the reads of the request body it wraps, as a
// RequestOption that counts or hashes the body can.
type tapBody struct {
	io.ReadCloser
	read        int64
	largestRead int
}

func (b *tapBody) Read(p []byte) (int, error) {
	b.largestRead = max(b.largestRead, len(p))
	n, err := b.ReadCloser.Read(p)
	b.read += int64(n)
	return n, err
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// liveHeap returns the bytes of heap in use after a garbage collection.
func liveHeap() int64 {
	runtime.GC()
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	return int64(stats.HeapAlloc)
}

func TestAttachmentBodyStaysOutOfMemoryWhenAnOptionWrapsIt(t *testing.T) {
	const size = 8 << 20
	encoded := int64(base64.StdEncoding.EncodedLen(size))
	for name, wrap := range map[string]bool{"without an option": false, "with an option that wraps the body": true} {
		t.Run(name, func(t *testing.T) {
			tempDir := isolateTempDir(t)
			type attempt struct {
				sum        [sha256.Size]byte
				length     int64
				declared   int64
				heapGrowth int64
			}
			var attempts []attempt
			var baseline int64
			config := files_sdk.Config{APIKey: "attachment-test", EndpointOverride: "https://files.invalid"}.Init()
			// The transport hashes each attempt's body without keeping it, and
			// measures the heap in use before reading it.
			config = config.SetCustomClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				growth := liveHeap() - baseline
				hash := sha256.New()
				length, err := io.Copy(hash, req.Body)
				// As any RoundTripper must, the transport closes the body.
				if closeErr := req.Body.Close(); err == nil {
					err = closeErr
				}
				if err != nil {
					return nil, err
				}
				attempts = append(attempts, attempt{sum: [sha256.Size]byte(hash.Sum(nil)), length: length, declared: req.ContentLength, heapGrowth: growth})
				status := http.StatusOK
				if len(attempts) == 1 {
					status = http.StatusServiceUnavailable
				}
				return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(`{}`)), Request: req}, nil
			})})
			var tap *tapBody
			var opts []files_sdk.RequestResponseOption
			if wrap {
				opts = append(opts, files_sdk.RequestOption(func(req *http.Request) error {
					tap = &tapBody{ReadCloser: req.Body}
					req.Body = tap
					return nil
				}))
			}

			baseline = liveHeap()
			_, err := (&user.Client{Config: config}).Update(files_sdk.UserUpdateParams{
				Id:                   7,
				AvatarFileAttachment: lib.NewJSONAttachment(&zeroReader{left: size}, size, "zeros.bin", "application/octet-stream"),
			}, opts...)
			require.NoError(t, err)

			require.Len(t, attempts, 2, "the 503 was retried")
			for i, a := range attempts {
				t.Logf("attempt %d: %d-byte body, heap in use %d bytes above the start", i+1, a.length, a.heapGrowth)
				assert.Equal(t, attempts[0].sum, a.sum, "the retry sent the same bytes")
				assert.Equal(t, a.declared, a.length, "the request declares its real length")
				assert.Greater(t, a.length, encoded)
				assert.Less(t, a.heapGrowth, encoded/4, "the %d-byte body is not held in memory", a.length)
			}
			if wrap {
				t.Logf("the option's wrapper was read for %d bytes, at most %d at a time", tap.read, tap.largestRead)
				assert.Equal(t, attempts[0].length, tap.read, "the wrapped body was read once, to its end")
				assert.LessOrEqual(t, tap.largestRead, 64<<10, "the wrapped body is read a little at a time")
			}
			assertNoFiles(t, tempDir)
		})
	}
}

// net/http lets a transport read and close a request body after RoundTrip,
// and so after the call, has returned.
func TestAttachmentBodyStaysReadableUntilTheTransportClosesIt(t *testing.T) {
	avatar := bytes.Repeat([]byte("avatar "), 200_000) // kept in a temporary file
	for name, test := range map[string]struct {
		afterReturn bool
		wrap        bool
		redirect    bool
	}{
		"read and closed during the call":                                               {},
		"read and closed after the call returns":                                        {afterReturn: true},
		"read and closed after the call returns, with a wrapped body":                   {afterReturn: true, wrap: true},
		"read and closed after the call returns, after a redirect":                      {afterReturn: true, redirect: true},
		"read and closed after the call returns, after a redirect, with a wrapped body": {afterReturn: true, wrap: true, redirect: true},
	} {
		t.Run(name, func(t *testing.T) {
			tempDir := isolateTempDir(t)
			type sent struct {
				body     []byte
				declared int64
				err      error
			}
			returned := make(chan struct{})
			redirected := make(chan sent, 1)
			done := make(chan sent, 1)
			var attempts atomic.Int32
			config := files_sdk.Config{APIKey: "attachment-test", EndpointOverride: "https://files.invalid"}.Init()
			config = config.SetCustomClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				consume := func() sent {
					body, err := io.ReadAll(req.Body)
					if closeErr := req.Body.Close(); err == nil {
						err = closeErr
					}
					return sent{body: body, declared: req.ContentLength, err: err}
				}
				if attempts.Add(1) == 1 && test.redirect {
					redirected <- consume()
					return &http.Response{StatusCode: http.StatusTemporaryRedirect, Header: http.Header{"Location": {req.URL.Path + "?redirected=true"}}, Body: http.NoBody, Request: req}, nil
				}
				if test.afterReturn {
					// Answer before sending the body, then send it once the call has returned.
					go func() {
						<-returned
						done <- consume()
					}()
				} else {
					done <- consume()
				}
				return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(`{}`)), Request: req}, nil
			})})
			var opts []files_sdk.RequestResponseOption
			if test.wrap {
				opts = append(opts, files_sdk.RequestOption(func(req *http.Request) error {
					req.Body = &tapBody{ReadCloser: req.Body}
					return nil
				}))
			}
			source := newCallerReader(avatar)

			_, err := (&user.Client{Config: config}).Update(files_sdk.UserUpdateParams{
				Id:                   7,
				Name:                 "Late Reader",
				ReceiveAdminAlerts:   lib.Bool(false),
				AvatarFileAttachment: lib.NewJSONAttachment(source, int64(len(avatar)), "avatar.png", "image/png"),
			}, opts...)
			close(returned)
			require.NoError(t, err)

			var got sent
			select {
			case got = <-done:
			case <-time.After(30 * time.Second):
				t.Fatal("the transport did not finish sending the request body")
			}
			require.NoError(t, got.err, "the transport read the whole body")
			assert.Equal(t, got.declared, int64(len(got.body)), "the request declares its real length")
			if test.redirect {
				assert.Equal(t, int32(2), attempts.Load())
				first := <-redirected
				require.NoError(t, first.err)
				assert.True(t, bytes.Equal(first.body, got.body), "the redirect resent the same bytes")
			} else {
				assert.Equal(t, int32(1), attempts.Load())
			}
			members := jsonMembers(t, got.body)
			assert.Equal(t, fileObject{content: avatar, filename: "avatar.png", contentType: "image/png"}, decodeFileObject(t, members["avatar_file"]))
			assert.JSONEq(t, `"Late Reader"`, string(members["name"]))
			assert.JSONEq(t, `false`, string(members["receive_admin_alerts"]))
			assert.False(t, source.closed, "the SDK does not close the caller's reader")
			// Once the transport has closed the body, nothing the SDK kept remains.
			assertNoFiles(t, tempDir)
		})
	}
}

// net/http opens the body for a redirect with GetBody before it calls
// CheckRedirect. When CheckRedirect stops the redirect, net/http closes the
// request's body except after ErrUseLastResponse, and never closes a body
// CheckRedirect replaced.
func TestAttachmentBodyOfAStoppedRedirectDoesNotKeepTheFiles(t *testing.T) {
	stopped := errors.New("the caller stopped the redirect")
	for name, test := range map[string]struct {
		// change changes the redirect's request as CheckRedirect can, and
		// returns the body it set there, if any.
		change func(req *http.Request) *redirectBody
		err    error
	}{
		"CheckRedirect leaves the body": {
			change: func(*http.Request) *redirectBody { return nil },
			err:    http.ErrUseLastResponse,
		},
		"CheckRedirect removes the body and returns its error": {
			change: func(req *http.Request) *redirectBody {
				req.Body, req.GetBody, req.ContentLength = http.NoBody, nil, 0
				return nil
			},
			err: stopped,
		},
		"CheckRedirect replaces the body and returns its error": {
			change: func(req *http.Request) *redirectBody {
				body := &redirectBody{Reader: strings.NewReader(`{}`)}
				req.Body, req.GetBody, req.ContentLength = body, nil, 2
				return body
			},
			err: stopped,
		},
		"CheckRedirect wraps the body and returns its error": {
			change: func(req *http.Request) *redirectBody {
				body := &redirectBody{Reader: req.Body, closer: req.Body}
				req.Body = body
				return body
			},
			err: stopped,
		},
	} {
		t.Run(name, func(t *testing.T) {
			tempDir := isolateTempDir(t)
			server := newAttachmentServer(t, tempDir, http.StatusTemporaryRedirect)
			var redirectBodies []io.ReadCloser
			var setBodies []*redirectBody
			config := server.config().SetCustomClient(&http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
				redirectBodies = append(redirectBodies, req.Body)
				if body := test.change(req); body != nil {
					setBodies = append(setBodies, body)
				}
				return test.err
			}})
			config.RetryMax = 0
			avatar := bytes.Repeat([]byte("avatar "), 200_000) // kept in a temporary file

			_, err := (&user.Client{Config: config}).Update(files_sdk.UserUpdateParams{
				Id:                   7,
				AvatarFileAttachment: lib.NewJSONAttachment(bytes.NewReader(avatar), int64(len(avatar)), "avatar.png", "image/png"),
			})
			// What the call does with the response ErrUseLastResponse keeps is
			// not what this checks.
			if test.err != http.ErrUseLastResponse {
				assert.ErrorIs(t, err, test.err, "the call returns CheckRedirect's error")
			}

			requests := server.recorded()
			require.NotEmpty(t, requests)
			assert.True(t, bytes.Equal(avatar, decodeFileObject(t, jsonMembers(t, requests[0].body)["avatar_file"]).content))
			require.NotEmpty(t, redirectBodies, "CheckRedirect stopped the redirect")
			assertNoFiles(t, tempDir)
			for _, body := range redirectBodies {
				require.NotNil(t, body)
				// On Unix a removed file can still be open; reading shows it is closed.
				_, err := io.ReadAll(body)
				assert.ErrorIs(t, err, os.ErrClosed, "the files behind the unsent body are closed")
				_ = body.Close()
			}
			for _, body := range setBodies {
				assert.Equal(t, 1, body.closes, "the body CheckRedirect set is closed once")
			}
		})
	}
}

// When CheckRedirect lets a redirect go ahead with a body of its own, the
// transport sends and closes that body, not the one net/http opened for the
// redirect with GetBody.
func TestAttachmentBodyOfAFollowedRedirectDoesNotKeepTheFiles(t *testing.T) {
	for name, test := range map[string]struct {
		// change changes the redirect's request as CheckRedirect can, and
		// returns the body it set there, if any.
		change func(req *http.Request) *redirectBody
		// resends is whether the redirect sends the attachment request's
		// body again, through the body CheckRedirect set; if not, it sends
		// sends.
		resends bool
		sends   string
		// late has the transport send the redirect's body after the call
		// returns.
		late bool
	}{
		"CheckRedirect removes the body": {
			change: func(req *http.Request) *redirectBody {
				req.Body, req.GetBody, req.ContentLength = http.NoBody, nil, 0
				return nil
			},
		},
		"CheckRedirect replaces the body": {
			change: func(req *http.Request) *redirectBody {
				body := &redirectBody{Reader: strings.NewReader(`{}`)}
				req.Body, req.GetBody, req.ContentLength = body, nil, 2
				return body
			},
			sends: `{}`,
		},
		"CheckRedirect wraps the body, sent after the call returns": {
			change: func(req *http.Request) *redirectBody {
				body := &redirectBody{Reader: req.Body, closer: req.Body}
				req.Body = body
				return body
			},
			resends: true,
			late:    true,
		},
		"CheckRedirect wraps the body without closing it, sent after the call returns": {
			change: func(req *http.Request) *redirectBody {
				body := &redirectBody{Reader: req.Body}
				req.Body = body
				return body
			},
			resends: true,
			late:    true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			tempDir := isolateTempDir(t)
			type sent struct {
				body     []byte
				declared int64
				err      error
			}
			returned := make(chan struct{})
			done := make(chan sent, 2)
			var attempts atomic.Int32
			var original io.ReadCloser
			var set *redirectBody
			config := files_sdk.Config{APIKey: "attachment-test", EndpointOverride: "https://files.invalid"}.Init()
			config = config.SetCustomClient(&http.Client{
				Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					consume := func() {
						body, err := io.ReadAll(req.Body)
						if closeErr := req.Body.Close(); err == nil {
							err = closeErr
						}
						done <- sent{body: body, declared: req.ContentLength, err: err}
					}
					if attempts.Add(1) == 1 {
						consume()
						return &http.Response{StatusCode: http.StatusTemporaryRedirect, Header: http.Header{"Location": {req.URL.Path + "?redirected=true"}}, Body: http.NoBody, Request: req}, nil
					}
					if test.late {
						// Answer before sending the body, then send it once the call has returned.
						go func() {
							<-returned
							consume()
						}()
					} else {
						consume()
					}
					return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(`{}`)), Request: req}, nil
				}),
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					original = req.Body
					set = test.change(req)
					return nil
				},
			})
			config.RetryMax = 0
			avatar := bytes.Repeat([]byte("avatar "), 200_000) // kept in a temporary file
			source := newCallerReader(avatar)

			_, err := (&user.Client{Config: config}).Update(files_sdk.UserUpdateParams{
				Id:                   7,
				AvatarFileAttachment: lib.NewJSONAttachment(source, int64(len(avatar)), "avatar.png", "image/png"),
			})
			close(returned)
			require.NoError(t, err)

			var first, redirected sent
			for _, got := range []*sent{&first, &redirected} {
				select {
				case *got = <-done:
				case <-time.After(30 * time.Second):
					t.Fatal("the transport did not finish sending the request bodies")
				}
			}
			assert.Equal(t, int32(2), attempts.Load())
			require.NoError(t, first.err)
			require.NoError(t, redirected.err, "the transport read the redirect's whole body")
			assert.True(t, bytes.Equal(avatar, decodeFileObject(t, jsonMembers(t, first.body)["avatar_file"]).content))
			assert.Equal(t, redirected.declared, int64(len(redirected.body)), "the redirect declares its real length")
			if test.resends {
				assert.True(t, bytes.Equal(first.body, redirected.body), "the redirect resent the same bytes")
			} else {
				assert.Equal(t, test.sends, string(redirected.body))
			}
			if set != nil {
				assert.Equal(t, 1, set.closes, "the body CheckRedirect set is closed once")
			}
			assert.False(t, source.closed, "the SDK does not close the caller's reader")
			assertNoFiles(t, tempDir)
			require.NotNil(t, original)
			if !test.resends {
				// On Unix a removed file can still be open; reading shows it is closed.
				_, err := io.ReadAll(original)
				assert.ErrorIs(t, err, os.ErrClosed, "the files behind the body the redirect did not send are closed")
			}
			_ = original.Close()
		})
	}
}

// redirectBody is a body a CheckRedirect sets on a redirect's request. It
// counts its Close calls and closes what it wraps, if anything.
type redirectBody struct {
	io.Reader
	closer io.Closer
	closes int
}

func (b *redirectBody) Close() error {
	b.closes++
	if b.closer != nil {
		return b.closer.Close()
	}
	return nil
}

func TestAttachmentErrorResponseKeepsItsErrorAndRemovesTheCopy(t *testing.T) {
	tempDir := isolateTempDir(t)
	server := newAttachmentServer(t, tempDir, http.StatusUnprocessableEntity)
	server.errorBody = `{"type":"processing-failure/invalid-image","http-code":422,"error":"Avatar must be an image"}`
	avatar := bytes.Repeat([]byte("not an image "), 200_000) // kept in a temporary file

	_, err := (&user.Client{Config: server.config()}).Update(files_sdk.UserUpdateParams{
		Id:                   7,
		AvatarFileAttachment: lib.NewJSONAttachment(bytes.NewReader(avatar), int64(len(avatar)), "avatar.png", "image/png"),
	})

	responseError, ok := err.(files_sdk.ResponseError)
	require.True(t, ok, "the error is the API's ResponseError, as without an attachment: %#v", err)
	assert.Equal(t, "processing-failure/invalid-image", responseError.Type)
	assert.True(t, files_sdk.IsProcessingFailure(err))
	requests := server.recorded()
	require.Len(t, requests, 1, "a 422 is not retried")
	assert.Equal(t, 1, requests[0].stagedFiles)
	assertNoFiles(t, tempDir)
}

func TestAttachmentFailuresSendNothing(t *testing.T) {
	// A reader's error can include the signed URL of its source.
	readFailure := &url.Error{Op: "Get", URL: "https://storage.example.invalid/avatar.png?label='avatar'&signature=secret", Err: io.ErrUnexpectedEOF}
	pathFailure := &fs.PathError{Op: "read", Path: "/Volumes/Photos/avatar.png", Err: errors.New("input/output error")}
	large := bytes.Repeat([]byte("large "), 600_000) // 3.6 MB, kept in a temporary file

	tests := map[string]struct {
		attachment func(t *testing.T, cancel context.CancelFunc) (*lib.JSONAttachment, *callerReader)
		params     func(*files_sdk.UserUpdateParams)
		expect     func(t *testing.T, err error, source *callerReader)
	}{
		"a reader longer than NewJSONAttachmentUpTo allows": {
			attachment: func(t *testing.T, _ context.CancelFunc) (*lib.JSONAttachment, *callerReader) {
				source := newCallerReader(large)
				return lib.NewJSONAttachmentUpTo(source, 2<<20, "large.png", "image/png"), source
			},
			expect: func(t *testing.T, err error, source *callerReader) {
				assert.ErrorIs(t, err, lib.ErrAttachmentTooLarge)
				assert.Equal(t, 2<<20+1, source.read, "the SDK reads one byte past the limit")
			},
		},
		"a reader shorter than its size": {
			attachment: func(t *testing.T, _ context.CancelFunc) (*lib.JSONAttachment, *callerReader) {
				source := newCallerReader([]byte("short"))
				return lib.NewJSONAttachment(source, 10, "short.png", "image/png"), source
			},
			expect: func(t *testing.T, err error, _ *callerReader) {
				assert.ErrorIs(t, err, lib.ErrAttachmentSizeMismatch)
				assert.ErrorContains(t, err, "ended after 5 of its 10 bytes")
			},
		},
		"a reader longer than its size": {
			attachment: func(t *testing.T, _ context.CancelFunc) (*lib.JSONAttachment, *callerReader) {
				source := newCallerReader(large)
				return lib.NewJSONAttachment(source, 3<<20, "changed.png", "image/png"), source
			},
			expect: func(t *testing.T, err error, _ *callerReader) {
				assert.ErrorIs(t, err, lib.ErrAttachmentSizeMismatch)
				assert.ErrorContains(t, err, "more than its size")
			},
		},
		"a reader that fails": {
			attachment: func(t *testing.T, _ context.CancelFunc) (*lib.JSONAttachment, *callerReader) {
				source := &callerReader{reader: io.MultiReader(bytes.NewReader(large), iotest.ErrReader(readFailure))}
				return lib.NewJSONAttachmentUpTo(source, 4<<20, "failing.png", "image/png"), source
			},
			expect: func(t *testing.T, err error, _ *callerReader) {
				assert.NotContains(t, err.Error(), "signature", "the message leaves out the signed URL")
				assert.EqualError(t, err, "attachment for avatar_file: reading: the reader returned an error")
				assert.ErrorIs(t, err, readFailure)
				assert.ErrorIs(t, err, io.ErrUnexpectedEOF)
				var urlError *url.Error
				require.ErrorAs(t, err, &urlError)
				assert.Equal(t, readFailure.URL, urlError.URL, "the reader's error itself is unchanged")
			},
		},
		"a reader that fails with a file path": {
			attachment: func(t *testing.T, _ context.CancelFunc) (*lib.JSONAttachment, *callerReader) {
				source := &callerReader{reader: iotest.ErrReader(pathFailure)}
				return lib.NewJSONAttachmentUpTo(source, 4<<20, "failing.png", "image/png"), source
			},
			expect: func(t *testing.T, err error, _ *callerReader) {
				assert.EqualError(t, err, "attachment for avatar_file: reading: the reader returned an error")
				var failure *fs.PathError
				require.ErrorAs(t, err, &failure)
				assert.Equal(t, "/Volumes/Photos/avatar.png", failure.Path)
			},
		},
		"no directory for the temporary copy": {
			attachment: func(t *testing.T, _ context.CancelFunc) (*lib.JSONAttachment, *callerReader) {
				missing := filepath.Join(t.TempDir(), "missing")
				t.Setenv("TMPDIR", missing)
				t.Setenv("TMP", missing)
				t.Setenv("TEMP", missing)
				source := newCallerReader(large)
				return lib.NewJSONAttachmentUpTo(source, 4<<20, "large.png", "image/png"), source
			},
			expect: func(t *testing.T, err error, source *callerReader) {
				// The SDK's own errors keep their messages, local paths included.
				assert.ErrorIs(t, err, fs.ErrNotExist)
				assert.ErrorContains(t, err, "attachment for avatar_file: creating a temporary file: ")
				assert.ErrorContains(t, err, filepath.Join("missing", "files-sdk-attachment-"))
				assert.Zero(t, source.read)
			},
		},
		"a context canceled while reading": {
			attachment: func(t *testing.T, cancel context.CancelFunc) (*lib.JSONAttachment, *callerReader) {
				source := newCallerReader(large)
				source.afterRead = cancel
				return lib.NewJSONAttachmentUpTo(source, 4<<20, "canceled.png", "image/png"), source
			},
			expect: func(t *testing.T, err error, source *callerReader) {
				assert.ErrorIs(t, err, context.Canceled)
				assert.Equal(t, 32<<10, source.read, "reading stopped at the next read after the cancellation")
			},
		},
		"the io.Writer field also set": {
			attachment: func(t *testing.T, _ context.CancelFunc) (*lib.JSONAttachment, *callerReader) {
				source := newCallerReader([]byte("avatar"))
				return lib.NewJSONAttachment(source, 6, "avatar.png", "image/png"), source
			},
			params: func(params *files_sdk.UserUpdateParams) {
				params.AvatarFile = &bytes.Buffer{}
			},
			expect: func(t *testing.T, err error, source *callerReader) {
				assert.EqualError(t, err, "AvatarFile and AvatarFileAttachment are both set; leave AvatarFile nil to send the attachment")
				assert.Zero(t, source.read, "nothing is read from a request that cannot be sent")
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tempDir := isolateTempDir(t)
			server := newAttachmentServer(t, tempDir)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			attachment, source := test.attachment(t, cancel)
			params := files_sdk.UserUpdateParams{Id: 7, Name: "Unchanged", AvatarFileAttachment: attachment}
			if test.params != nil {
				test.params(&params)
			}

			_, err := (&user.Client{Config: server.config()}).Update(params, files_sdk.WithContext(ctx))

			require.Error(t, err)
			test.expect(t, err, source)
			assert.Empty(t, server.recorded(), "no request was sent")
			assert.False(t, source.closed, "the SDK does not close the caller's reader")
			assertNoFiles(t, tempDir)
		})
	}

	t.Run("an attachment sent a second time", func(t *testing.T) {
		tempDir := isolateTempDir(t)
		server := newAttachmentServer(t, tempDir)
		client := user.Client{Config: server.config()}
		params := files_sdk.UserUpdateParams{Id: 7, AvatarFileAttachment: lib.NewJSONAttachment(bytes.NewReader([]byte("avatar")), 6, "avatar.png", "image/png")}
		_, err := client.Update(params)
		require.NoError(t, err)

		_, err = client.Update(params)

		assert.EqualError(t, err, "attachment for avatar_file: already sent; make a new JSONAttachment for each request")
		assert.Len(t, server.recorded(), 1)
	})
}
