package folder

import (
	"bytes"
	"log"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/lib"
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
			files_sdk.Config{},
			args{params: files_sdk.FolderListForParams{WithPreviews: lib.Bool(true)}, opts: []files_sdk.RequestResponseOption{}},
			"with_preview",
		},
		{
			"with path it send fields",
			files_sdk.Config{},
			args{params: files_sdk.FolderListForParams{Path: "anything", WithPreviews: lib.Bool(true)}, opts: []files_sdk.RequestResponseOption{}},
			"with_preview",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.Config.Debug = true
			var buf bytes.Buffer
			logger := log.New(&buf, "InMemoryLogger: ", log.LstdFlags)

			tt.Config.SetLogger(logger)
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
