package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type BandwidthSnapshot struct {
	Id                int64      `json:"id,omitempty" path:"id"`
	BytesReceived     string     `json:"bytes_received,omitempty" path:"bytes_received"`
	BytesSent         string     `json:"bytes_sent,omitempty" path:"bytes_sent"`
	SyncBytesReceived string     `json:"sync_bytes_received,omitempty" path:"sync_bytes_received"`
	SyncBytesSent     string     `json:"sync_bytes_sent,omitempty" path:"sync_bytes_sent"`
	RequestsGet       string     `json:"requests_get,omitempty" path:"requests_get"`
	RequestsPut       string     `json:"requests_put,omitempty" path:"requests_put"`
	RequestsOther     string     `json:"requests_other,omitempty" path:"requests_other"`
	LoggedAt          *time.Time `json:"logged_at,omitempty" path:"logged_at"`
}

type BandwidthSnapshotCollection []BandwidthSnapshot

type BandwidthSnapshotListParams struct {
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty" path:"filter_lteq"`
	lib.ListParams
}

func (b *BandwidthSnapshot) UnmarshalJSON(data []byte) error {
	type bandwidthSnapshot BandwidthSnapshot
	var v bandwidthSnapshot
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*b = BandwidthSnapshot(v)
	return nil
}

func (b *BandwidthSnapshotCollection) UnmarshalJSON(data []byte) error {
	type bandwidthSnapshots BandwidthSnapshotCollection
	var v bandwidthSnapshots
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*b = BandwidthSnapshotCollection(v)
	return nil
}

func (b *BandwidthSnapshotCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*b))
	for i, v := range *b {
		ret[i] = v
	}

	return &ret
}
