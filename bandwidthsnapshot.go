package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type BandwidthSnapshot struct {
	Id                int64     `json:"id,omitempty"`
	BytesReceived     float32   `json:"bytes_received,omitempty"`
	BytesSent         float32   `json:"bytes_sent,omitempty"`
	SyncBytesReceived float32   `json:"sync_bytes_received,omitempty"`
	SyncBytesSent     float32   `json:"sync_bytes_sent,omitempty"`
	RequestsGet       float32   `json:"requests_get,omitempty"`
	RequestsPut       float32   `json:"requests_put,omitempty"`
	RequestsOther     float32   `json:"requests_other,omitempty"`
	LoggedAt          time.Time `json:"logged_at,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitempty"`
}

type BandwidthSnapshotCollection []BandwidthSnapshot

type BandwidthSnapshotListParams struct {
	Cursor     string          `url:"cursor,omitempty" required:"false"`
	PerPage    int             `url:"per_page,omitempty" required:"false"`
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false"`
	lib.ListParams
}

func (b *BandwidthSnapshot) UnmarshalJSON(data []byte) error {
	type bandwidthSnapshot BandwidthSnapshot
	var v bandwidthSnapshot
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*b = BandwidthSnapshot(v)
	return nil
}

func (b *BandwidthSnapshotCollection) UnmarshalJSON(data []byte) error {
	type bandwidthSnapshots []BandwidthSnapshot
	var v bandwidthSnapshots
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*b = BandwidthSnapshotCollection(v)
	return nil
}
