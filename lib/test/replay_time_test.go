package test_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/downloadurl"
	"github.com/Files-com/files-sdk-go/v3/file"
	fixturetest "github.com/Files-com/files-sdk-go/v3/lib/test"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/stretchr/testify/require"
)

func TestReplayRelocatesTransferTimesWithoutChangingSignedQueryOrFileBytes(t *testing.T) {
	t.Setenv("GITLAB", "1")
	t.Chdir(t.TempDir())
	const api = "https://api.example.test/api/rest/v1/files/"
	const captured = "Tue, 25 Oct 2022 23:42:33 GMT"
	const originalDate = "20221025T234233Z"
	const query = "z=%2B&note=%26X-Amz-Date=19900101T000000Z&X-Amz-Date=" + originalDate + "&X-Amz-Expires=60&a=keep%2fraw&X-Amz-Signature=fixture"
	const active = "https://storage.example.test/api/rest/v1/archive.json?" + query
	const expired = "https://storage.example.test/expired?X-Amz-Date=20221025T234203Z&X-Amz-Expires=10&X-Amz-Signature=expired-fixture"
	const payload = `{"expires":"2000-01-01T00:00:00Z","example":"X-Amz-Date=19900101T000000Z"}`
	c := cassette.New(filepath.Join("fixtures", "timeline"))
	add := func(uri, body string) {
		c.AddInteraction(&cassette.Interaction{Request: cassette.Request{Method: "GET", URL: uri}, Response: cassette.Response{Code: 200, Status: "200 OK", Headers: http.Header{"Date": []string{captured}, "Content-Type": []string{"application/json"}}, Body: body}})
	}
	metadata, _ := json.Marshal(map[string]string{"download_uri": active, "expires": "2022-10-25T23:43:18Z"})
	add(api+"file", string(metadata))
	add(active, payload)
	metadata, _ = json.Marshal(map[string]string{"download_uri": expired, "expires": "2022-10-25T23:42:13Z"})
	add(api+"expired", string(metadata))
	require.NoError(t, c.Save())
	before, err := os.ReadFile(c.File)
	require.NoError(t, err)
	config, recorder, err := fixturetest.CreateConfigWithConfig("timeline", files_sdk.Config{EndpointOverride: "https://api.example.test"})
	require.NoError(t, err)
	defer recorder.Stop()
	config.DisableDirectTransfers = true
	load := func(name string) (string, string, time.Time) {
		request, err := http.NewRequest("GET", api+name, nil)
		require.NoError(t, err)
		response, err := files_sdk.WrapRequestOptions(config, request)
		require.NoError(t, err)
		var data struct {
			URI     string `json:"download_uri"`
			Expires string `json:"expires"`
		}
		require.NoError(t, json.NewDecoder(response.Body).Decode(&data))
		require.NoError(t, response.Body.Close())
		date, err := http.ParseTime(response.Header.Get("Date"))
		require.NoError(t, err)
		return data.URI, data.Expires, date
	}
	uri, expires, date := load("file")
	partExpiry, err := time.Parse(time.RFC3339, expires)
	require.NoError(t, err)
	require.Equal(t, 45*time.Second, partExpiry.Sub(date))
	parsed, err := downloadurl.New(uri)
	require.NoError(t, err)
	_, valid := parsed.Valid(0)
	require.True(t, valid)
	require.Equal(t, time.Minute, parsed.Time.Sub(date))
	var sentQuery string
	var output []byte
	_, err = (&file.Client{Config: config}).Download(files_sdk.FileDownloadParams{File: files_sdk.File{Path: "file", DownloadUri: uri}},
		files_sdk.RequestOption(func(request *http.Request) error { sentQuery = request.URL.RawQuery; return nil }),
		files_sdk.ResponseBodyOption(func(body io.ReadCloser) error {
			defer body.Close()
			var err error
			output, err = io.ReadAll(body)
			return err
		}))
	require.NoError(t, err)
	require.Equal(t, payload, string(output))
	requestURL, err := url.Parse(uri)
	require.NoError(t, err)
	newDate := requestURL.Query().Get("X-Amz-Date")
	require.NotEqual(t, originalDate, newDate)
	require.Equal(t, strings.Replace(query, originalDate, "DATE", 1), strings.Replace(sentQuery, newDate, "DATE", 1))
	uri, expires, date = load("expired")
	parsed, err = downloadurl.New(uri)
	require.NoError(t, err)
	_, valid = parsed.Valid(0)
	require.False(t, valid)
	partExpiry, err = time.Parse(time.RFC3339, expires)
	require.NoError(t, err)
	require.Equal(t, -20*time.Second, partExpiry.Sub(date))
	require.NoError(t, recorder.Stop())
	after, err := os.ReadFile(c.File)
	require.NoError(t, err)
	require.Equal(t, before, after)
}
