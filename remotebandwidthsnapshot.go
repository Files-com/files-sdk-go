package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type RemoteBandwidthSnapshot struct {
	Id                int64      `json:"id,omitempty"`
	SyncBytesReceived float32    `json:"sync_bytes_received,omitempty"`
	SyncBytesSent     float32    `json:"sync_bytes_sent,omitempty"`
	LoggedAt          *time.Time `json:"logged_at,omitempty"`
	RemoteServerId    int64      `json:"remote_server_id,omitempty"`
}

type RemoteBandwidthSnapshotCollection []RemoteBandwidthSnapshot

type RemoteBandwidthSnapshotListParams struct {
	Cursor     string          `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage    int64           `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false" json:"filter_like,omitempty"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty"`
	lib.ListParams
}

func (r *RemoteBandwidthSnapshot) UnmarshalJSON(data []byte) error {
	type remoteBandwidthSnapshot RemoteBandwidthSnapshot
	var v remoteBandwidthSnapshot
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*r = RemoteBandwidthSnapshot(v)
	return nil
}

func (r *RemoteBandwidthSnapshotCollection) UnmarshalJSON(data []byte) error {
	type remoteBandwidthSnapshots []RemoteBandwidthSnapshot
	var v remoteBandwidthSnapshots
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*r = RemoteBandwidthSnapshotCollection(v)
	return nil
}

func (r *RemoteBandwidthSnapshotCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*r))
	for i, v := range *r {
		ret[i] = v
	}

	return &ret
}
