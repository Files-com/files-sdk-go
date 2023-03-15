package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type RemoteBandwidthSnapshot struct {
	Id                int64      `json:"id,omitempty" path:"id"`
	SyncBytesReceived string     `json:"sync_bytes_received,omitempty" path:"sync_bytes_received"`
	SyncBytesSent     string     `json:"sync_bytes_sent,omitempty" path:"sync_bytes_sent"`
	LoggedAt          *time.Time `json:"logged_at,omitempty" path:"logged_at"`
	RemoteServerId    int64      `json:"remote_server_id,omitempty" path:"remote_server_id"`
}

type RemoteBandwidthSnapshotCollection []RemoteBandwidthSnapshot

type RemoteBandwidthSnapshotListParams struct {
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty" path:"filter_lteq"`
	lib.ListParams
}

func (r *RemoteBandwidthSnapshot) UnmarshalJSON(data []byte) error {
	type remoteBandwidthSnapshot RemoteBandwidthSnapshot
	var v remoteBandwidthSnapshot
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*r = RemoteBandwidthSnapshot(v)
	return nil
}

func (r *RemoteBandwidthSnapshotCollection) UnmarshalJSON(data []byte) error {
	type remoteBandwidthSnapshots RemoteBandwidthSnapshotCollection
	var v remoteBandwidthSnapshots
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
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
