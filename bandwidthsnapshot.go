package files_sdk

import (
	"encoding/json"
	lib "github.com/Files-com/files-sdk-go/lib"
	"time"
)

type BandwidthSnapshot struct {
	Id            int64     `json:"id,omitempty"`
	BytesReceived float32   `json:"bytes_received,omitempty"`
	BytesSent     float32   `json:"bytes_sent,omitempty"`
	RequestsGet   float32   `json:"requests_get,omitempty"`
	RequestsPut   float32   `json:"requests_put,omitempty"`
	RequestsOther float32   `json:"requests_other,omitempty"`
	LoggedAt      time.Time `json:"logged_at,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}

type BandwidthSnapshotCollection []BandwidthSnapshot

type BandwidthSnapshotListParams struct {
	Page       int             `url:"page,omitempty"`
	PerPage    int             `url:"per_page,omitempty"`
	Action     string          `url:"action,omitempty"`
	Cursor     string          `url:"cursor,omitempty"`
	SortBy     json.RawMessage `url:"sort_by,omitempty"`
	Filter     json.RawMessage `url:"filter,omitempty"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty"`
	FilterLike json.RawMessage `url:"filter_like,omitempty"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty"`
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
