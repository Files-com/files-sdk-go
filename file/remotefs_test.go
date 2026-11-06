package file

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	goFs "io/fs"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

// createRequests records what reaches one path of a MockAPIServer while a file is written there: begin_upload requests,
// lookups, the parts' bodies, which it answers itself so that it keeps them, and the size each finalizing request gave.
type createRequests struct {
	mutex sync.Mutex
	log   createRequestLog
}

type createRequestLog struct {
	begins    int
	finds     int
	content   []byte
	finalized []int64
}

func recordCreateRequests(server *MockAPIServer, path string) *createRequests {
	requests := &createRequests{}
	server.MockRoute("/api/rest/v1/file_actions/begin_upload/"+path, func(*gin.Context, interface{}) bool {
		requests.mutex.Lock()
		defer requests.mutex.Unlock()
		requests.log.begins++
		return false
	})
	server.MockRoute("/upload/"+path, func(c *gin.Context, _ interface{}) bool {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return true
		}
		requests.mutex.Lock()
		requests.log.content = append(requests.log.content, body...)
		requests.mutex.Unlock()
		c.Header("Etag", "part")
		c.Status(http.StatusOK)
		return true
	})
	server.MockRoute("/api/rest/v1/files/"+path, func(c *gin.Context, model interface{}) bool {
		requests.mutex.Lock()
		defer requests.mutex.Unlock()
		if c.Request.Method == http.MethodGet {
			requests.log.finds++
		} else if create, ok := model.(files_sdk.FileCreateParams); ok {
			requests.log.finalized = append(requests.log.finalized, create.Size)
		}
		return false
	})
	server.MockRoute("/api/rest/v1/file_actions/metadata/"+path, func(*gin.Context, interface{}) bool {
		requests.mutex.Lock()
		defer requests.mutex.Unlock()
		requests.log.finds++
		return false
	})
	return requests
}

// seen returns what has reached the path so far.
func (r *createRequests) seen() createRequestLog {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	copied := r.log
	copied.content = bytes.Clone(r.log.content)
	copied.finalized = slices.Clone(r.log.finalized)
	return copied
}

// stageCreateFilesIn makes dir the directory os.TempDir names, where Create keeps what is written, for the rest of the
// test.
func stageCreateFilesIn(t *testing.T, dir string) {
	t.Helper()
	for _, name := range []string{"TMPDIR", "TMP", "TEMP"} {
		t.Setenv(name, dir)
	}
}

func stagingEntries(t *testing.T, dir string) []os.DirEntry {
	t.Helper()
	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	return entries
}

// Create's writer keeps what it is given in a temporary file and sends nothing until Close, which uploads exactly those
// bytes, none for an empty file, finalizes once and removes the temporary file. A second Close sends nothing more, and a
// Write after Close fails.
func TestFS_CreateUploadsWhatWasWrittenOnlyOnClose(t *testing.T) {
	for _, test := range []struct{ name, content string }{
		{name: "content", content: "written in two pieces, the second copied"},
		{name: "empty", content: ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			server := (&MockAPIServer{T: t}).Do()
			defer server.Shutdown()
			staging := t.TempDir()
			stageCreateFilesIn(t, staging)
			requests := recordCreateRequests(server, "created.txt")

			w, err := (&FS{}).Init(server.Client().Config, false).Create("created.txt")
			require.NoError(t, err)
			half := len(test.content) / 2
			n, err := w.Write([]byte(test.content[:half]))
			require.NoError(t, err)
			require.Equal(t, half, n)
			_, err = io.Copy(w, strings.NewReader(test.content[half:]))
			require.NoError(t, err)
			require.Len(t, stagingEntries(t, staging), 1)
			require.Zero(t, requests.seen().begins, "nothing is sent before Close")

			require.NoError(t, w.Close())
			uploaded := requests.seen()
			require.Equal(t, 1, uploaded.begins)
			require.Equal(t, test.content, string(uploaded.content))
			require.Equal(t, []int64{int64(len(test.content))}, uploaded.finalized)
			require.Empty(t, stagingEntries(t, staging))

			require.NoError(t, w.Close())
			require.Equal(t, uploaded, requests.seen(), "a second Close sends nothing")
			_, err = w.Write([]byte("late"))
			require.ErrorIs(t, err, os.ErrClosed)
		})
	}
}

// A temporary file that cannot be created fails Create. After a write fails, or once the writer's context has ended,
// Close sends nothing, returns that error and removes the temporary file.
func TestFS_CreateUploadsNothingAfterAFailure(t *testing.T) {
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	config := server.Client().Config

	t.Run("temporary file cannot be created", func(t *testing.T) {
		stageCreateFilesIn(t, filepath.Join(t.TempDir(), "missing"))

		w, err := (&FS{}).Init(config, false).Create("uncreated.txt")
		require.ErrorIs(t, err, os.ErrNotExist)
		require.Nil(t, w)
	})

	t.Run("a write failed", func(t *testing.T) {
		staging := t.TempDir()
		stageCreateFilesIn(t, staging)
		requests := recordCreateRequests(server, "partial.txt")
		w, err := (&FS{}).Init(config, false).Create("partial.txt")
		require.NoError(t, err)
		_, err = w.Write([]byte("accepted"))
		require.NoError(t, err)
		// Closing the temporary file under the writer makes its next write fail, as a full disk would.
		require.NoError(t, w.(*stagedRemoteFile).staging.Close())
		_, writeErr := w.Write([]byte("refused"))
		require.Error(t, writeErr)

		require.ErrorIs(t, w.Close(), writeErr)
		require.Zero(t, requests.seen().begins)
		require.Empty(t, stagingEntries(t, staging))
	})

	t.Run("context ended before Close", func(t *testing.T) {
		staging := t.TempDir()
		stageCreateFilesIn(t, staging)
		requests := recordCreateRequests(server, "cancelled.txt")
		ctx, cancel := context.WithCancel(context.Background())
		w, err := (&FS{}).Init(config, false).WithContext(ctx).(*FS).Create("cancelled.txt")
		require.NoError(t, err)
		_, err = w.Write([]byte("never sent"))
		require.NoError(t, err)
		cancel()

		require.ErrorIs(t, w.Close(), context.Canceled)
		require.Zero(t, requests.seen().begins)
		require.Empty(t, stagingEntries(t, staging))
	})
}

// When the writer's context ends while Close's upload is sending a part, Close returns the upload's error once the part's
// request has been given up, without finalizing the upload, and removes the temporary file.
func TestFS_CreateCloseSettlesWhenItsContextEndsDuringTheUpload(t *testing.T) {
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	staging := t.TempDir()
	stageCreateFilesIn(t, staging)
	requests := recordCreateRequests(server, "held.txt")
	arrived := make(chan struct{}, 1)
	server.MockRoute("/upload/held.txt", func(c *gin.Context, _ interface{}) bool {
		// Reading the whole body lets the server notice the client giving the request up.
		_, _ = io.Copy(io.Discard, c.Request.Body)
		select {
		case arrived <- struct{}{}:
		default:
		}
		select {
		case <-c.Request.Context().Done():
		case <-time.After(10 * time.Second):
		}
		c.Status(http.StatusServiceUnavailable)
		return true
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w, err := (&FS{}).Init(server.Client().Config, false).WithContext(ctx).(*FS).Create("held.txt")
	require.NoError(t, err)
	_, err = w.Write([]byte("held in flight"))
	require.NoError(t, err)

	closed := make(chan error, 1)
	go func() { closed <- w.Close() }()
	select {
	case <-arrived:
	case err := <-closed:
		t.Fatalf("Close returned before its part reached storage: %v", err)
	case <-time.After(10 * time.Second):
		t.Fatal("the part never reached storage")
	}
	cancel()
	select {
	case err = <-closed:
	case <-time.After(10 * time.Second):
		t.Fatal("Close did not return after its context ended")
	}

	require.Error(t, err)
	require.Empty(t, requests.seen().finalized)
	require.Empty(t, stagingEntries(t, staging))
}

// After Close attempts its upload, an FS with its cache on and an empty Root looks the file up again instead of
// returning the one it cached before, even when the writer came from an FS WithContext made from it.
func TestFS_CreateCloseDropsCachedFilesOfEveryFSSharingTheCache(t *testing.T) {
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	server.MockFiles["cached.txt"] = mockFile{File: files_sdk.File{Path: "cached.txt", DisplayName: "cached.txt", Type: "file", Size: 3}}
	requests := recordCreateRequests(server, "cached.txt")
	fs := (&FS{}).Init(server.Client().Config, true)

	before, err := fs.Open("cached.txt")
	require.NoError(t, err)
	cached, err := fs.Open("cached.txt")
	require.NoError(t, err)
	require.Same(t, before, cached)

	w, err := fs.WithContext(context.Background()).(*FS).Create("cached.txt")
	require.NoError(t, err)
	_, err = w.Write([]byte("replaced"))
	require.NoError(t, err)
	require.NoError(t, w.Close())

	after, err := fs.Open("cached.txt")
	require.NoError(t, err)
	require.NotSame(t, before, after)
	require.Equal(t, 2, requests.seen().finds)
}

// An FS's Root stays with the FS WithContext makes from it, and Create writes under it, where Open finds the file.
func TestFS_CreateWritesUnderTheRootOpenUses(t *testing.T) {
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	requests := recordCreateRequests(server, "root/rooted.txt")
	fs := (&FS{}).Init(server.Client().Config, false)
	fs.Root = "root"
	rooted := fs.WithContext(context.Background()).(*FS)

	w, err := rooted.Create("rooted.txt")
	require.NoError(t, err)
	_, err = w.Write([]byte("under root"))
	require.NoError(t, err)
	require.NoError(t, w.Close())
	require.Equal(t, "under root", string(requests.seen().content))

	opened, err := rooted.Open("rooted.txt")
	require.NoError(t, err)
	info, err := opened.Stat()
	require.NoError(t, err)
	require.Equal(t, "root/rooted.txt", info.Sys().(files_sdk.File).Path)
}

func TestFileReadUsesStatusOnlyWithRequestID(t *testing.T) {
	for _, requestID := range []string{"", "request-id"} {
		t.Run("request_id="+requestID, func(t *testing.T) {
			var downloadRequests, statusRequests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.URL.Path {
				case "/download":
					downloadRequests.Add(1)
					w.Header().Set("X-Files-Download-Request-Id", requestID)
					w.WriteHeader(http.StatusBadRequest)
					_, _ = io.WriteString(w, `{"error":"original download failure"}`)
				case "/download/request-id":
					statusRequests.Add(1)
					_, _ = io.WriteString(w, `{"error":"detailed status failure","data":{"status":"failed"}}`)
				default:
					t.Errorf("unexpected request: %s", r.URL.Path)
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			t.Cleanup(server.Close)
			config := files_sdk.Config{Logger: log.New(io.Discard, "", 0)}.Init()
			f := (&File{
				File: &files_sdk.File{Path: "/download.bin", DownloadUri: server.URL + "/download"},
				FS:   (&FS{Context: context.Background()}).Init(config, false),
			}).Init()

			_, err := f.Read(make([]byte, 1))

			require.Error(t, err)
			require.EqualValues(t, 1, downloadRequests.Load())
			if requestID == "" {
				require.ErrorContains(t, err, "original download failure")
				require.Zero(t, statusRequests.Load())
			} else {
				require.EqualError(t, err, "detailed status failure")
				require.EqualValues(t, 1, statusRequests.Load())
			}
		})
	}
}

// The variables targets/mock-server/with-scenario.rb gives a test it runs against a simulator it started: the
// simulator's origin, schema and instance, and the synthetic API key to send it.
var remoteFSSimulatorVariables = []string{
	"FILES_MOCK_SERVER_URL",
	"FILES_MOCK_SERVER_HOST",
	"FILES_MOCK_SERVER_PORT",
	"FILES_MOCK_SCHEMA_SHA256",
	"FILES_MOCK_INSTANCE",
	"FILES_MOCK_API_KEY",
}

// remoteFSSimulatorClient returns a client whose requests go only to the Files.com mock server's simulator that
// targets/mock-server/with-scenario.rb started for this test, with the synthetic API key it gave, and whose transfers
// go through the API as they always have for TestFS_Open. It checks that the launcher's variables name that simulator
// on a loopback address and that its readiness report shows it ready, in simulation, at contract 3, with the
// launcher's schema and instance, and still as the launcher's reset left it: at epoch 1, under the api_scenarios
// profile of the api-scenarios-canonical preset, whose root listing TestFS_Open reads. So it never resets the
// simulator, and every file test that uses it can share one launcher run: each works under remote paths of its own.
// Without the launcher the test is skipped, unless FILES_API_SCENARIOS is "required"; a missing or invalid variable
// fails it. It never falls back to recorded responses, Files.com or another API key.
func remoteFSSimulatorClient(t *testing.T) *Client {
	t.Helper()
	launcher := map[string]string{}
	var missing []string
	for _, name := range remoteFSSimulatorVariables {
		launcher[name] = strings.TrimSpace(os.Getenv(name))
		if launcher[name] == "" {
			missing = append(missing, name)
		}
	}
	// FILES_MOCK_SERVER_HOST and FILES_MOCK_SERVER_PORT alone may name another mock server, so only the variables
	// with-scenario.rb alone sets show that it launched this test.
	launched := launcher["FILES_MOCK_SERVER_URL"] != "" || launcher["FILES_MOCK_SCHEMA_SHA256"] != "" ||
		launcher["FILES_MOCK_INSTANCE"] != "" || launcher["FILES_MOCK_API_KEY"] != ""
	if !launched {
		why := "Run this test through targets/mock-server/with-scenario.rb with the api-scenarios-canonical preset, which starts the Files.com mock server's simulator for it; it never uses Files.com."
		if os.Getenv("FILES_API_SCENARIOS") == "required" {
			t.Fatalf("FILES_API_SCENARIOS is required, so this test may not be skipped: %s", why)
		}
		t.Skip(why)
	}
	if len(missing) > 0 {
		t.Fatalf("with-scenario.rb sets %s, but %s is missing.", strings.Join(remoteFSSimulatorVariables, ", "), strings.Join(missing, ", "))
	}
	host, port := launcher["FILES_MOCK_SERVER_HOST"], launcher["FILES_MOCK_SERVER_PORT"]
	if !slices.Contains([]string{"127.0.0.1", "::1", "localhost"}, host) {
		t.Fatalf("FILES_MOCK_SERVER_HOST must be a loopback address (127.0.0.1, ::1, localhost), not %q.", host)
	}
	if number, err := strconv.Atoi(port); err != nil || len(port) > 5 || number < 1 || number > 65535 {
		t.Fatalf("FILES_MOCK_SERVER_PORT must be a port number, not %q.", port)
	}
	origin := "http://" + net.JoinHostPort(host, port)
	if launcher["FILES_MOCK_SERVER_URL"] != origin {
		t.Fatalf("FILES_MOCK_SERVER_URL is %q, not %s, the simulator FILES_MOCK_SERVER_HOST and FILES_MOCK_SERVER_PORT name.", launcher["FILES_MOCK_SERVER_URL"], origin)
	}

	// with-scenario.rb runs the test only once the simulator has answered its readiness, so one request is enough.
	transport := http.DefaultTransport.(*http.Transport).Clone()
	t.Cleanup(transport.CloseIdleConnections)
	control := &http.Client{Transport: transport, Timeout: 10 * time.Second}
	response, err := control.Get(origin + "/__files_mock/v1/ready")
	if err != nil {
		t.Fatalf("Nothing answered at %s for its readiness report: %v", origin, err)
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	response.Body.Close()
	if err != nil || response.StatusCode != http.StatusOK {
		t.Fatalf("%s answered %d for its readiness report, so it is not the simulator (%v).", origin, response.StatusCode, err)
	}
	var reported struct {
		Status          string `json:"status"`
		Mode            string `json:"mode"`
		ContractVersion int    `json:"contract_version"`
		SchemaSHA256    string `json:"schema_sha256"`
		Instance        string `json:"instance"`
		Epoch           int    `json:"epoch"`
		Profile         struct {
			APIScenarios map[string]any `json:"api_scenarios"`
		} `json:"profile"`
	}
	if err := json.Unmarshal(body, &reported); err != nil {
		t.Fatalf("%s answered its readiness report with something other than the simulator's JSON: %v", origin, err)
	}
	if reported.Status != "ready" || reported.Mode != "simulation" || reported.ContractVersion != 3 ||
		reported.SchemaSHA256 != launcher["FILES_MOCK_SCHEMA_SHA256"] || reported.Instance != launcher["FILES_MOCK_INSTANCE"] {
		t.Fatalf("%s is not the simulator with-scenario.rb started: it reports %+v, not ready, simulation, contract 3, schema %s and instance %s.", origin, reported, launcher["FILES_MOCK_SCHEMA_SHA256"], launcher["FILES_MOCK_INSTANCE"])
	}
	if reported.Epoch != 1 || reported.Profile.APIScenarios == nil {
		t.Fatalf("The simulator at %s is at epoch %d with the api_scenarios profile %v, not as the launcher's api-scenarios-canonical reset left it: at epoch 1, under an api_scenarios profile, whose root listing this test reads.", origin, reported.Epoch, reported.Profile.APIScenarios)
	}
	return &Client{Config: files_sdk.Config{APIKey: launcher["FILES_MOCK_API_KEY"], EndpointOverride: origin, DisableDirectTransfers: true}.Init()}
}

// failSimulatorRequests makes the simulator that client's requests go to answer each of the next count requests that
// match rule with the Files.com error rule describes, by adding a fault rule at the simulator's control endpoint, out
// of band from the SDK under test. The rule selects every matching request until count have failed, so it changes
// nothing for any request after them. It returns a check, to call once the test's requests have been made, that
// requires exactly count matching requests to have arrived and each to have failed.
func failSimulatorRequests(t *testing.T, client *Client, count int, rule map[string]any) func() {
	t.Helper()
	transport := http.DefaultTransport.(*http.Transport).Clone()
	t.Cleanup(transport.CloseIdleConnections)
	control := &http.Client{Transport: transport, Timeout: 10 * time.Second}
	origin := client.Config.EndpointOverride
	var keys []string
	for key := range rule["match"].(map[string]any) {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	spec := map[string]any{"random": map[string]any{"seed": t.Name(), "rate": 1, "key": keys, "max_faults": count}}
	for field, value := range rule {
		spec[field] = value
	}
	encoded, err := json.Marshal(spec)
	require.NoError(t, err)
	response, err := control.Post(origin+"/__files_mock/v1/faults", "application/json", bytes.NewReader(encoded))
	require.NoError(t, err)
	answer, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	response.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, response.StatusCode, "the simulator refused the fault rule: %s", answer)
	var added struct {
		ID int `json:"id"`
	}
	require.NoError(t, json.Unmarshal(answer, &added))
	require.Positive(t, added.ID, "the simulator answered the fault rule without its ID: %s", answer)

	return func() {
		t.Helper()
		response, err := control.Get(origin + "/__files_mock/v1/faults")
		require.NoError(t, err)
		answer, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
		response.Body.Close()
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode, "the simulator did not list its fault rules: %s", answer)
		var listed struct {
			Faults []struct {
				ID     int `json:"id"`
				Random struct {
					Outcomes map[string]int `json:"outcomes"`
				} `json:"random"`
			} `json:"faults"`
		}
		require.NoError(t, json.Unmarshal(answer, &listed))
		for _, fault := range listed.Faults {
			if fault.ID == added.ID {
				want := map[string]int{"fault": count, "pass": 0, "overridden": 0, "unavailable": 0, "budget-exhausted": 0}
				require.Equal(t, want, fault.Random.Outcomes, "the requests fault rule %d matched", added.ID)
				return
			}
		}
		t.Fatalf("The simulator no longer lists fault rule %d: %s", added.ID, answer)
	}
}

func TestFS_Open(t *testing.T) {
	client := remoteFSSimulatorClient(t)

	assert := assert.New(t)
	require.NoError(t, client.Upload(
		UploadWithReader(strings.NewReader("testing 3")),
		UploadWithSize(9),
		UploadWithDestinationPath(filepath.Join("remotefs_test", "1.text")),
	))

	fs := (&FS{}).Init(client.Config, true)
	fs = fs.WithContext(context.TODO()).(*FS)
	f, err := fs.Open("remotefs_test")
	assert.NoError(err)
	rf, ok := f.(*ReadDirFile)
	assert.True(ok)
	entry, err := rf.ReadDir(0)
	assert.NoError(err)
	assert.Equal(1, len(entry))
	assert.False(entry[0].IsDir())
	info, err := entry[0].Info()
	assert.NoError(err)
	assert.Equal("1.text", info.Name())
	fsFile, ok := entry[0].(*File)
	assert.True(ok)
	buf := make([]byte, 8)
	_, err = io.ReadFull(fsFile, buf)
	assert.NoError(err)
	err = fsFile.Close()
	assert.NoError(err)
	assert.Equal("testing ", string(buf))

	buf = make([]byte, 9)
	_, err = io.ReadFull(fsFile, buf)
	assert.NoError(err)
	err = fsFile.Close()
	assert.NoError(err)

	assert.Equal("testing 3", string(buf))

	_, err = fs.ReadDir(".")
	require.NoError(t, err)

	_, err = fs.Open(".")
	require.NoError(t, err)
}

func TestFS_OpenUsesServerUnicodeIdentityInCache(t *testing.T) {
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	client := server.Client()
	for _, name := range []string{"q\u0301/カ.txt", "q/か.txt "} {
		server.MockFiles[name] = mockFile{File: files_sdk.File{Path: name, DisplayName: filepath.Base(name), Type: "file"}}
	}
	fs := (&FS{}).Init(client.Config, true)
	original, err := fs.Open("q\u0301/カ.txt")
	require.NoError(t, err)
	equivalent, err := fs.Open("q/か.txt")
	require.NoError(t, err)
	require.Same(t, original, equivalent)
	distinct, err := fs.Open("q/か.txt ")
	require.NoError(t, err)
	require.NotSame(t, original, distinct)
	info, err := distinct.Stat()
	require.NoError(t, err)
	require.Equal(t, "q/か.txt ", info.Sys().(files_sdk.File).Path)
}

func TestFS_ReadDirContinuesAcrossCalls(t *testing.T) {
	for _, test := range []struct {
		name    string
		counts  []int
		entries [][]string
	}{
		{name: "bounded batches", counts: []int{1, 1, 1}, entries: [][]string{{"a"}, {"b"}, {"c"}}},
		{name: "all remaining", counts: []int{1, -1}, entries: [][]string{{"a"}, {"b", "c"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				w.Header().Set("Content-Type", "application/json")
				if r.URL.Query().Get("cursor") == "" {
					w.Header().Set("X-Files-Cursor-Next", "second")
					_, _ = io.WriteString(w, `[{"path":"a","display_name":"a","type":"file"},{"path":"nested/skipped","display_name":"skipped","type":"file"}]`)
					return
				}
				_, _ = io.WriteString(w, `[{"path":"b","display_name":"b","type":"file"},{"path":"c","display_name":"c","type":"directory"}]`)
			}))
			t.Cleanup(server.Close)
			fs := (&FS{}).Init(files_sdk.Config{EndpointOverride: server.URL}.Init(), true)
			opened, err := fs.Open(".")
			require.NoError(t, err)
			directory := opened.(goFs.ReadDirFile)

			for i, count := range test.counts {
				entries, err := directory.ReadDir(count)
				require.NoError(t, err)
				var names []string
				for _, entry := range entries {
					names = append(names, entry.Name())
				}
				require.Equal(t, test.entries[i], names)
				if i == 0 {
					require.EqualValues(t, 1, requests.Load(), "the first batch must not request the next page")
					other, err := fs.Open(".")
					require.NoError(t, err)
					first, err := other.(goFs.ReadDirFile).ReadDir(1)
					require.NoError(t, err)
					require.Len(t, first, 1)
					require.Equal(t, "a", first[0].Name(), "a separate open has its own directory cursor")
				}
			}
			for _, count := range []int{1, 2, 0, -1} {
				entries, err := directory.ReadDir(count)
				if count > 0 {
					require.ErrorIs(t, err, io.EOF)
				} else {
					require.NoError(t, err)
				}
				require.NotNil(t, entries)
				require.Empty(t, entries)
			}
			require.EqualValues(t, 3, requests.Load(), "exhausted reads do not restart a listing")
		})
	}
}

func TestFS_ReadDirReturnsEntriesBeforeALaterPageFails(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Get("cursor") == "" {
			w.Header().Set("X-Files-Cursor-Next", "second")
			_, _ = io.WriteString(w, `[{"path":"a","display_name":"a","type":"file"}]`)
			return
		}
		w.WriteHeader(http.StatusForbidden)
		_, _ = io.WriteString(w, `{"type":"permission-denied","error":"listing denied","http-code":403}`)
	}))
	t.Cleanup(server.Close)
	config := files_sdk.Config{EndpointOverride: server.URL}.Init()
	config.RetryMax = 0
	opened, err := (&FS{}).Init(config, false).Open(".")
	require.NoError(t, err)
	directory := opened.(goFs.ReadDirFile)

	entries, err := directory.ReadDir(2)
	require.Len(t, entries, 1)
	require.Equal(t, "a", entries[0].Name())
	var pathErr *goFs.PathError
	require.ErrorAs(t, err, &pathErr)
	require.Equal(t, "readdir", pathErr.Op)
	require.True(t, files_sdk.IsErrorType(err, files_sdk.ResponseErrorType("permission-denied")))
	entries, laterErr := directory.ReadDir(1)
	require.Empty(t, entries)
	require.Equal(t, err, laterErr)
	require.EqualValues(t, 2, requests.Load())
}

func TestFS_MkdirAllPreservesCreationErrors(t *testing.T) {
	transportErr := errors.New("transport failed")
	permissionErr := files_sdk.ResponseError{Type: "permission-denied", ErrorMessage: "creation denied"}
	for _, test := range []struct {
		name string
		err  error
	}{
		{name: "transport", err: transportErr},
		{name: "cancellation", err: context.Canceled},
		{name: "wrapped API error", err: fmt.Errorf("wrapped: %w", permissionErr)},
	} {
		t.Run(test.name, func(t *testing.T) {
			var paths []string
			config := files_sdk.Config{}.Init().SetCustomClient(&http.Client{Transport: ownedHTTPResponseTransport(func(request *http.Request) (*http.Response, error) {
				paths = append(paths, request.URL.Path)
				if len(paths) > 1 {
					return nil, test.err
				}
				return &http.Response{
					StatusCode: http.StatusUnprocessableEntity,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(`{"type":"processing-failure/destination-exists","error":"exists","http-code":422}`)),
					Request:    request,
				}, nil
			})})
			config.RetryMax = 0

			err := (&FS{}).Init(config, false).MkdirAll("existing/failed/not-submitted", 0750)

			require.ErrorIs(t, err, test.err)
			require.Equal(t, []string{"/api/rest/v1/folders/existing", "/api/rest/v1/folders/existing/failed"}, paths)
		})
	}
}

func TestFS_MkdirAllAcceptsWrappedDestinationExists(t *testing.T) {
	var paths []string
	config := files_sdk.Config{}.Init().SetCustomClient(&http.Client{Transport: ownedHTTPResponseTransport(func(request *http.Request) (*http.Response, error) {
		paths = append(paths, request.URL.Path)
		if len(paths) == 1 {
			return nil, fmt.Errorf("wrapped: %w", files_sdk.ResponseError{Type: string(files_sdk.ErrDestinationExists)})
		}
		return &http.Response{
			StatusCode: http.StatusCreated,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"path":"existing/child","type":"directory"}`)),
			Request:    request,
		}, nil
	})})
	config.RetryMax = 0

	require.NoError(t, (&FS{}).Init(config, false).MkdirAll("existing/child", 0750))
	require.Equal(t, []string{"/api/rest/v1/folders/existing", "/api/rest/v1/folders/existing/child"}, paths)
}

func TestFileMetadataReportsDirectoryType(t *testing.T) {
	for _, kind := range []string{"directory", "file"} {
		t.Run(kind, func(t *testing.T) {
			file := (&File{File: &files_sdk.File{Type: kind}}).Init()
			info, err := file.Stat()
			require.NoError(t, err)
			require.Equal(t, kind == "directory", file.Type().IsDir())
			require.Equal(t, info.IsDir(), info.Mode().IsDir())
			require.Equal(t, info.Mode().Type(), file.Type())
			require.Zero(t, file.Type()&goFs.ModePerm, "DirEntry.Type contains type bits only")
		})
	}
}

func TestFS_MkdirTempReportsUnsupportedWithoutRequests(t *testing.T) {
	for _, test := range []struct {
		name string
		dir  string
		path string
	}{
		{name: "default directory", path: "job-*"},
		{name: "explicit directory", dir: "parent", path: "parent/job-*"},
	} {
		t.Run(test.name, func(t *testing.T) {
			var requests atomic.Int32
			config := files_sdk.Config{}.Init().SetCustomClient(&http.Client{Transport: ownedHTTPResponseTransport(func(*http.Request) (*http.Response, error) {
				requests.Add(1)
				return nil, errors.New("unexpected API request")
			})})
			remote := (&FS{Root: "root"}).Init(config, false)

			path, err := remote.MkdirTemp(test.dir, "job-*")

			require.Empty(t, path)
			require.ErrorIs(t, err, errors.ErrUnsupported)
			var pathErr *goFs.PathError
			require.ErrorAs(t, err, &pathErr)
			require.Equal(t, "mkdirtemp", pathErr.Op)
			require.Equal(t, test.path, pathErr.Path)
			require.ErrorContains(t, err, "atomic exclusive remote directory allocation is unavailable")
			require.Zero(t, requests.Load())
		})
	}
}

func TestReadDirFileValuesShareTheirOpenDirectoryCursor(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `[{"path":"a","display_name":"a","type":"file"},{"path":"b","display_name":"b","type":"file"}]`)
	}))
	t.Cleanup(server.Close)
	fs := (&FS{}).Init(files_sdk.Config{EndpointOverride: server.URL}.Init(), false)
	file := (&File{File: &files_sdk.File{Type: "directory"}, FS: fs}).Init()
	value := ReadDirFile{File: file}
	copyOfValue := value
	var directory goFs.ReadDirFile = value
	var copiedDirectory goFs.ReadDirFile = copyOfValue

	first, err := directory.ReadDir(1)
	require.NoError(t, err)
	require.Len(t, first, 1)
	require.Equal(t, "a", first[0].Name())
	second, err := copiedDirectory.ReadDir(1)
	require.NoError(t, err)
	require.Len(t, second, 1)
	require.Equal(t, "b", second[0].Name())
	for _, handle := range []goFs.ReadDirFile{directory, copiedDirectory} {
		entries, err := handle.ReadDir(1)
		require.ErrorIs(t, err, io.EOF)
		require.Empty(t, entries)
	}
	require.EqualValues(t, 1, requests.Load(), "copies continue one directory traversal")
}
