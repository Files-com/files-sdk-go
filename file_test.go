package files_sdk

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// FileDownloadParams.File is metadata the SDK already has about the file,
// such as a listing entry it is about to download. It is never part of the
// request.
func TestFileDownloadParams_EncodeTheRequestButNeverTheKnownFile(t *testing.T) {
	params := FileDownloadParams{
		Path:         "reports/q1.txt",
		PreviewSize:  "large",
		WithPreviews: lib.Bool(true),
		File: File{
			Path:        "reports/q1.txt",
			DisplayName: "known-display-name.txt",
			Md5:         "known-md5",
			DownloadUri: "https://storage.example.invalid/q1?signature=known-signature",
		},
	}

	// A JSON round trip of FileDownloadParams keeps its own request fields,
	// its path included.
	encoded, err := json.Marshal(params)
	require.NoError(t, err)
	for _, known := range []string{"known-display-name.txt", "known-md5", "known-signature"} {
		assert.NotContains(t, string(encoded), known, "the known File is not encoded")
	}
	var decoded FileDownloadParams
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	want := params
	want.File = File{}
	assert.Equal(t, want, decoded)

	// The request puts the path in the URL and sends only the ordinary fields
	// as query values.
	values, err := lib.Params{Params: params}.ToValues()
	require.NoError(t, err)
	assert.Equal(t, url.Values{"preview_size": {"large"}, "with_previews": {"true"}}, values)
	path, err := lib.BuildPath("/files/{path}", params)
	require.NoError(t, err)
	assert.Equal(t, "/files/reports/q1.txt", path)
}
