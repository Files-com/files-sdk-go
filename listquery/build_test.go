package listquery

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildFollowsNextCursorToLaterPageError(t *testing.T) {
	const cursor = "2:a+b&=opaque"
	firstPageHeaders := map[string]http.Header{
		"documented next cursor":     {"X-Files-Cursor-Next": {cursor}},
		"legacy cursor":              {"X-Files-Cursor": {cursor}},
		"next preferred over legacy": {"X-Files-Cursor-Next": {cursor}, "X-Files-Cursor": {"stale-legacy-cursor"}},
	}
	for name, firstPageHeader := range firstPageHeaders {
		t.Run(name, func(t *testing.T) {
			var requests []*http.Request
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests = append(requests, r.Clone(r.Context()))
				w.Header().Set("Content-Type", "application/json")
				if len(requests) == 1 {
					for key, values := range firstPageHeader {
						w.Header()[key] = values
					}
					w.Write([]byte(`[{"id":73},{"id":74}]`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"error":"Later page was not found","http-code":404,"title":"Not Found","type":"not-found"}`))
			}))
			defer server.Close()

			config := files_sdk.Config{APIKey: "paging-fixture-key", WorkspaceId: 7, EndpointOverride: server.URL}.Init()
			params := struct {
				UserId int64 `url:"user_id,omitempty"`
				files_sdk.ListParams
			}{UserId: 17, ListParams: files_sdk.ListParams{PerPage: 2}}
			it := &files_sdk.Iter{Query: Build(config, "/bundles", &idList{}), ListParams: &params}

			var ids []int64
			for it.Next() {
				ids = append(ids, it.Current().(idRecord).Id)
			}

			assert.Equal(t, []int64{73, 74}, ids, "records from the first page are kept")
			assert.Error(t, it.Err(), "the later page's failure is reported")
			require.Len(t, requests, 2)
			assert.False(t, requests[0].URL.Query().Has("cursor"))
			assert.Equal(t, cursor, requests[1].URL.Query().Get("cursor"))
			for _, request := range requests {
				assert.Equal(t, "2", request.URL.Query().Get("per_page"))
				assert.Equal(t, "17", request.URL.Query().Get("user_id"))
				assert.Equal(t, "paging-fixture-key", request.Header.Get("X-FilesAPI-Key"))
				assert.Equal(t, "7", request.Header.Get("X-Files-Workspace-Id"))
			}
		})
	}
}

type idRecord struct {
	Id int64 `json:"id"`
}

type idList []idRecord

func (l *idList) UnmarshalJSON(data []byte) error {
	var records []idRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return err
	}
	*l = records
	return nil
}

func (l *idList) ToSlice() *[]interface{} {
	values := make([]interface{}, len(*l))
	for i, record := range *l {
		values[i] = record
	}
	return &values
}
