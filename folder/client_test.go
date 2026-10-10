package folder

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_ListFor(t *testing.T) {
	type args struct {
		params files_sdk.FolderListForParams
		opts   []files_sdk.RequestResponseOption
	}
	tests := []struct {
		name string
		files_sdk.Config
		args        args
		debugOutput string
	}{
		{
			"without path it send fields",
			files_sdk.Config{}.Init(),
			args{params: files_sdk.FolderListForParams{WithPreviews: lib.Bool(true)}, opts: []files_sdk.RequestResponseOption{}},
			"with_preview",
		},
		{
			"with path it send fields",
			files_sdk.Config{}.Init(),
			args{params: files_sdk.FolderListForParams{Path: "anything", WithPreviews: lib.Bool(true)}, opts: []files_sdk.RequestResponseOption{}},
			"with_preview",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.Config.Debug = true
			var buf bytes.Buffer
			logger := log.New(&buf, "InMemoryLogger: ", log.LstdFlags)

			tt.Config.Logger = logger
			c := &Client{
				Config: tt.Config,
			}

			it, err := c.ListFor(tt.args.params, tt.args.opts...)
			require.NoError(t, err)
			it.GetPage()
			assert.Contains(t, buf.String(), tt.debugOutput)
		})
	}
}

func TestClient_ListFor_All(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Get("cursor") == "" {
			w.Header().Set("X-Files-Cursor-Next", "page-2")
			w.Write([]byte(`[{"path":"a"},{"path":"b"}]`))
			return
		}
		w.Write([]byte(`[{"path":"c"}]`))
	}))
	defer server.Close()
	c := &Client{Config: files_sdk.Config{EndpointOverride: server.URL}.Init()}

	it, err := c.ListFor(files_sdk.FolderListForParams{})
	require.NoError(t, err)
	files := it.All()
	assert.Equal(t, int32(0), requests.Load(), "creating the iterator requests nothing")

	// A folder listing yields files.
	var paths []string
	for file, err := range files {
		require.NoError(t, err)
		paths = append(paths, file.Path)
		if len(paths) == 1 {
			break
		}
	}
	assert.Equal(t, int32(1), requests.Load(), "a break requests no further pages")
	for file, err := range it.All() {
		require.NoError(t, err)
		paths = append(paths, file.Path)
	}
	assert.Equal(t, []string{"a", "b", "c"}, paths, "the next loop continues after the break")
	assert.Equal(t, int32(2), requests.Load())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	canceled, err := c.ListFor(files_sdk.FolderListForParams{}, files_sdk.WithContext(ctx))
	require.NoError(t, err)
	var errs []error
	for file, err := range canceled.All() {
		assert.Empty(t, file.Path)
		errs = append(errs, err)
	}
	require.Len(t, errs, 1)
	assert.ErrorIs(t, errs[0], context.Canceled)
}
