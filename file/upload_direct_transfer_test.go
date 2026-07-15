package file

import (
	"context"
	"errors"
	"io"
	"log"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestUploadIORequestsDirectConnectionInfoByDefaultAndOmitsItWhenDisabled(t *testing.T) {
	for _, test := range []struct {
		name     string
		disabled bool
		expected *bool
	}{
		{name: "default", expected: lib.Bool(true)},
		{name: "disabled", disabled: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			server := (&MockAPIServer{T: t}).Do()
			defer server.Shutdown()
			client := server.Client()
			client.Config.DisableDirectTransfers = test.disabled
			var directConnectionInfo *bool
			server.MockRoute("/api/rest/v1/file_actions/begin_upload/upload.bin", func(_ *gin.Context, model interface{}) bool {
				directConnectionInfo = model.(files_sdk.FileBeginUploadParams).WithDirectConnectionInfo
				return false
			})

			_, err := (&uploadIO{Client: client, uploadV2: true}).startUpload(context.Background(), files_sdk.FileBeginUploadParams{
				Path:         "upload.bin",
				MkdirParents: lib.Bool(true),
			})

			require.NoError(t, err)
			require.Equal(t, test.expected, directConnectionInfo)
		})
	}
}

func TestUploadIODirectTransferDisableSuppressesLaterParts(t *testing.T) {
	connectionInfo := testDirectConnectionInfo("/uploads?jwt=upload-token")
	u := &uploadIO{
		Client: &Client{
			Config: files_sdk.Config{
				Logger: log.New(io.Discard, "", 0),
			}.Init(),
		},
		Path: "/upload.bin",
	}

	require.True(t, u.directTransferUploadAvailable(connectionInfo))

	firstDisable := u.disableDirectTransfersForUpload("direct_request_failed", errors.New("dial failed"), 1)

	require.True(t, firstDisable)
	require.False(t, u.directTransferUploadAvailable(connectionInfo))

	secondDisable := u.disableDirectTransfersForUpload("direct_request_failed", errors.New("still failed"), 2)

	require.False(t, secondDisable)
}

func TestUploadIODirectTransfersDisabled(t *testing.T) {
	u := &uploadIO{Client: &Client{Config: files_sdk.Config{DisableDirectTransfers: true}.Init()}}

	require.False(t, u.directTransferUploadAvailable(testDirectConnectionInfo("/uploads?jwt=upload-token")))
}

func testDirectConnectionInfo(path string) files_sdk.DirectConnectionInfo {
	return files_sdk.DirectConnectionInfo{
		Version:    1,
		ServerName: "agent-123.agents.files.internal",
		Addresses:  []string{"tcp://127.0.0.1:4001"},
		DirectUri:  "https://agent-123.agents.files.internal:4001" + path,
		CaPem:      "test-ca",
	}
}
