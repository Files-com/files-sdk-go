package file_migration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/lib/test"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func CreateClient(fixture string) (client *Client, r *recorder.Recorder, err error) {
	client = &Client{}
	client.Config, r, err = test.CreateConfig(fixture)

	return client, r, err
}

func TestClient_Wait(t *testing.T) {
	client, r, err := CreateClient("TestClient_Wait")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	assert := assert.New(t)
	var migrationPassedToFunc files_sdk.FileMigration
	statusFunc := func(migration files_sdk.FileMigration) {
		migrationPassedToFunc = migration
	}
	//
	fileAction := files_sdk.FileAction{Status: "in_progress", FileMigrationId: 11}

	migration, err := client.Wait(fileAction, statusFunc)

	assert.Error(err, "Not Found - `Not Found`")
	assert.Equal("", migration.Status)
	assert.Equal(int64(0), migration.Id)
	//
	fileAction = files_sdk.FileAction{Status: "in_progress", FileMigrationId: 12}

	migration, err = client.Wait(fileAction, statusFunc)

	assert.NoError(err)
	assert.Equal("completed", migrationPassedToFunc.Status)
	assert.Equal("completed", migration.Status)
	assert.Equal(int64(12), migration.Id)
	assert.Equal(int64(12), migrationPassedToFunc.Id)
	assert.Equal("video.mp4", migrationPassedToFunc.Path)

	//
	fileAction = files_sdk.FileAction{Status: "completed", FileMigrationId: 13}

	migration, err = client.Wait(fileAction, statusFunc)

	assert.NoError(err)
	assert.Equal("completed", migration.Status)
	assert.Equal(int64(13), migration.Id)

	//
	fileAction = files_sdk.FileAction{FileMigrationId: 14}

	migration, err = client.Wait(fileAction, statusFunc)

	assert.NoError(err)
	assert.Equal("failed", migration.Status)
	assert.Equal(int64(14), migration.Id)

	//
	fileAction = files_sdk.FileAction{FileMigrationId: 15}

	migration, err = client.Wait(fileAction, statusFunc)

	assert.NoError(err)
	assert.Equal("completed", migration.Status)
	assert.Equal(int64(15), migration.Id)
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
