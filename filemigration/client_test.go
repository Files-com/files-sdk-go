package file_migration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The variables targets/mock-server/with-scenario.rb gives a test it runs against a simulator it started: the
// simulator's origin, schema and instance, and the synthetic API key to send it.
var simulatorVariables = []string{
	"FILES_MOCK_SERVER_URL",
	"FILES_MOCK_SERVER_HOST",
	"FILES_MOCK_SERVER_PORT",
	"FILES_MOCK_SCHEMA_SHA256",
	"FILES_MOCK_INSTANCE",
	"FILES_MOCK_API_KEY",
}

// newSimulatorClient returns a client whose requests go only to the Files.com mock server's simulator that
// targets/mock-server/with-scenario.rb started for this test, with the synthetic API key it gave. First it checks that
// the launcher's variables name that simulator on a loopback address and that its readiness report shows it ready, in
// simulation, at contract 3, with the launcher's schema and instance; only then does it reset the simulator with
// reset, the body of POST /__files_mock/v1/reset. Without the launcher the test is skipped, unless
// FILES_API_SCENARIOS is "required"; a missing or invalid variable fails it. It never falls back to recorded
// responses, Files.com or another API key.
func newSimulatorClient(t *testing.T, reset map[string]any) *Client {
	t.Helper()
	launcher := map[string]string{}
	var missing []string
	for _, name := range simulatorVariables {
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
		why := "Run this test through targets/mock-server/with-scenario.rb, which starts the Files.com mock server's simulator for it; it never uses Files.com."
		if os.Getenv("FILES_API_SCENARIOS") == "required" {
			t.Fatalf("FILES_API_SCENARIOS is required, so this test may not be skipped: %s", why)
		}
		t.Skip(why)
	}
	if len(missing) > 0 {
		t.Fatalf("with-scenario.rb sets %s, but %s is missing.", strings.Join(simulatorVariables, ", "), strings.Join(missing, ", "))
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
	}
	if err := json.Unmarshal(body, &reported); err != nil {
		t.Fatalf("%s answered its readiness report with something other than the simulator's JSON: %v", origin, err)
	}
	if reported.Status != "ready" || reported.Mode != "simulation" || reported.ContractVersion != 3 ||
		reported.SchemaSHA256 != launcher["FILES_MOCK_SCHEMA_SHA256"] || reported.Instance != launcher["FILES_MOCK_INSTANCE"] {
		t.Fatalf("%s is not the simulator with-scenario.rb started: it reports %+v, not ready, simulation, contract 3, schema %s and instance %s.", origin, reported, launcher["FILES_MOCK_SCHEMA_SHA256"], launcher["FILES_MOCK_INSTANCE"])
	}

	encoded, err := json.Marshal(reset)
	require.NoError(t, err)
	response, err = control.Post(origin+"/__files_mock/v1/reset", "application/json", bytes.NewReader(encoded))
	if err != nil {
		t.Fatalf("Could not reset the simulator at %s: %v", origin, err)
	}
	answer, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	response.Body.Close()
	if err != nil || response.StatusCode != http.StatusOK {
		t.Fatalf("The simulator at %s answered its reset with %d: %s (%v)", origin, response.StatusCode, answer, err)
	}
	return &Client{Config: files_sdk.Config{APIKey: launcher["FILES_MOCK_API_KEY"], EndpointOverride: origin}.Init()}
}

func TestClient_Wait(t *testing.T) {
	// Under the pending file_actions profile a copy answers with the FileMigration that will make it, as Files.com
	// answers one it enqueues. The simulator reports each migration processing when it is first polled and applies its
	// copy the next time, against the files there are then.
	client := newSimulatorClient(t, map[string]any{
		"profile":  map[string]any{"file_actions": "pending"},
		"fixtures": map[string]any{"files": []map[string]any{{"path": "video.mp4", "text": "video"}}},
	})
	fileClient := file.Client{Config: client.Config}
	pendingCopy := func(destination string) files_sdk.FileAction {
		t.Helper()
		fileAction, err := fileClient.Copy(files_sdk.FileCopyParams{Path: "video.mp4", Destination: destination})
		require.NoError(t, err)
		require.Equal(t, "pending", fileAction.Status)
		require.NotZero(t, fileAction.FileMigrationId)
		return fileAction
	}

	assert := assert.New(t)
	var migrationPassedToFunc files_sdk.FileMigration
	statusFunc := func(migration files_sdk.FileMigration) {
		migrationPassedToFunc = migration
	}
	// No migration has been made yet, so there is none with this ID.
	fileAction := files_sdk.FileAction{Status: "in_progress", FileMigrationId: 11}

	migration, err := client.Wait(fileAction, statusFunc)

	assert.Error(err, "Not Found - `Not Found`")
	assert.Equal("", migration.Status)
	assert.Equal(int64(0), migration.Id)
	//
	fileAction = pendingCopy("video-2.mp4")
	// A second copy to video-2.mp4, made before the first has been, so its migration fails once the first has made it.
	conflictingCopy := pendingCopy("video-2.mp4")

	migration, err = client.Wait(fileAction, statusFunc)

	assert.NoError(err)
	assert.Equal("completed", migrationPassedToFunc.Status)
	assert.Equal("completed", migration.Status)
	assert.Equal(fileAction.FileMigrationId, migration.Id)
	assert.Equal(fileAction.FileMigrationId, migrationPassedToFunc.Id)
	assert.Equal("video.mp4", migrationPassedToFunc.Path)

	// An action that has completed is not polled, so a migration ID the simulator never made is not looked up.
	fileAction = files_sdk.FileAction{Status: "completed", FileMigrationId: 13}

	migration, err = client.Wait(fileAction, statusFunc)

	assert.NoError(err)
	assert.Equal("completed", migration.Status)
	assert.Equal(int64(13), migration.Id)

	//
	fileAction = files_sdk.FileAction{FileMigrationId: conflictingCopy.FileMigrationId}

	migration, err = client.Wait(fileAction, statusFunc)

	assert.NoError(err)
	assert.Equal("failed", migration.Status)
	assert.Equal(fileAction.FileMigrationId, migration.Id)

	//
	fileAction = files_sdk.FileAction{FileMigrationId: pendingCopy("video-3.mp4").FileMigrationId}

	migration, err = client.Wait(fileAction, statusFunc)

	assert.NoError(err)
	assert.Equal("completed", migration.Status)
	assert.Equal(fileAction.FileMigrationId, migration.Id)
}

func TestClient_WaitContextCancellationInterruptsPollDelay(t *testing.T) {
	var requestCount atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requestCount.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":16,"path":"video.mp4","operation":"copy","status":"in_progress"}`))
	}))
	defer server.Close()

	client := &Client{Config: files_sdk.Config{EndpointOverride: server.URL}.Init()}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	waitError := make(chan error, 1)
	firstStatusReceived := make(chan struct{})
	var notifyFirstStatus sync.Once
	go func() {
		_, err := client.Wait(files_sdk.FileAction{Status: "in_progress", FileMigrationId: 16}, func(files_sdk.FileMigration) {
			notifyFirstStatus.Do(func() {
				close(firstStatusReceived)
			})
		}, files_sdk.WithContext(ctx))
		waitError <- err
	}()

	select {
	case <-firstStatusReceived:
	case <-time.After(time.Second):
		t.Fatal("Wait did not receive its initial migration status")
	}

	time.Sleep(50 * time.Millisecond)
	requestsBeforeCancel := requestCount.Load()
	cancel()

	select {
	case err := <-waitError:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("Wait did not return after its context was canceled")
	}
	assert.Equal(t, int64(1), requestsBeforeCancel, "Wait polled again before the one-second interval elapsed")
}
