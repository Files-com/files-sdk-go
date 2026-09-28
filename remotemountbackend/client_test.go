package remote_mount_backend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type recordedRequest struct {
	method, path, apiKey string
	body                 map[string]interface{}
}

// newDecimalServer records each request and answers with a backend whose
// decimal fields are strings, as the API returns them.
func newDecimalServer(t *testing.T) (*Client, func() []recordedRequest) {
	var mu sync.Mutex
	var requests []recordedRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.UseNumber() // compare number tokens, not float64 conversions
		var body map[string]interface{}
		assert.NoError(t, decoder.Decode(&body))
		mu.Lock()
		requests = append(requests, recordedRequest{r.Method, r.URL.Path, r.Header.Get("X-FilesAPI-Key"), body})
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":7,"min_free_cpu":"1.0049999999999999999999999999","min_free_mem":"0"}`)
	}))
	t.Cleanup(server.Close)
	client := &Client{Config: files_sdk.Config{APIKey: "decimal-key", EndpointOverride: server.URL}.Init()}
	return client, func() []recordedRequest {
		mu.Lock()
		defer mu.Unlock()
		return append([]recordedRequest(nil), requests...)
	}
}

func TestClient_CreateAndUpdateSendExactDecimalText(t *testing.T) {
	client, requests := newDecimalServer(t)
	exact := "1.0049999999999999999999999999"

	backend, err := client.Create(files_sdk.RemoteMountBackendCreateParams{
		MinFreeCpuDecimal: lib.String(exact),
		MinFreeMem:        12.5,
		CanaryFilePath:    "canary.txt", RemoteServerMountId: 2, RemoteServerId: 3,
	}, files_sdk.WithContext(context.Background()))
	require.NoError(t, err)
	assert.Equal(t, exact, backend.MinFreeCpu, "responses keep the API's decimal string")

	_, err = client.Update(files_sdk.RemoteMountBackendUpdateParams{Id: 7, MinFreeCpu: 1.005, MinFreeMemDecimal: lib.String("0")})
	require.NoError(t, err)

	_, err = client.UpdateWithMap(map[string]interface{}{"id": 7, "min_free_cpu": exact})
	require.NoError(t, err)

	assert.Equal(t, []recordedRequest{
		{"POST", "/api/rest/v1/remote_mount_backends", "decimal-key", map[string]interface{}{
			"min_free_cpu": exact, "min_free_mem": json.Number("12.5"),
			"canary_file_path": "canary.txt", "remote_server_mount_id": json.Number("2"), "remote_server_id": json.Number("3"),
		}},
		{"PATCH", "/api/rest/v1/remote_mount_backends/7", "decimal-key", map[string]interface{}{
			"min_free_cpu": json.Number("1.005"), "min_free_mem": "0",
		}},
		{"PATCH", "/api/rest/v1/remote_mount_backends/7", "decimal-key", map[string]interface{}{
			"id": json.Number("7"), "min_free_cpu": exact,
		}},
	}, requests())
}

func TestClient_CreateAndUpdateRejectInvalidDecimalsWithoutSending(t *testing.T) {
	client, requests := newDecimalServer(t)
	required := files_sdk.RemoteMountBackendCreateParams{CanaryFilePath: "canary.txt", RemoteServerMountId: 2, RemoteServerId: 3}

	for name, call := range map[string]func() error{
		"create with float and decimal for one field": func() error {
			params := required
			params.MinFreeCpu, params.MinFreeCpuDecimal = 1.5, lib.String("1.5")
			_, err := client.Create(params)
			return err
		},
		"create with a non-finite float": func() error {
			params := required
			params.MinFreeMem = math.Inf(1)
			_, err := client.Create(params)
			return err
		},
		"update with hexadecimal text": func() error {
			_, err := client.Update(files_sdk.RemoteMountBackendUpdateParams{Id: 7, MinFreeMemDecimal: lib.String("0x1p0")})
			return err
		},
		"update with empty text": func() error {
			_, err := client.Update(files_sdk.RemoteMountBackendUpdateParams{Id: 7, MinFreeCpuDecimal: lib.String("")})
			return err
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Error(t, call())
		})
	}
	assert.Empty(t, requests())
}
