package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// RemoteBandwidthSnapshot is a Files.com API resource.
type RemoteBandwidthSnapshot struct {
	Id                int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	SyncBytesReceived int64      `json:"sync_bytes_received,omitempty" path:"sync_bytes_received,omitempty" url:"sync_bytes_received,omitempty"`
	SyncBytesSent     int64      `json:"sync_bytes_sent,omitempty" path:"sync_bytes_sent,omitempty" url:"sync_bytes_sent,omitempty"`
	LoggedAt          *time.Time `json:"logged_at,omitempty" path:"logged_at,omitempty" url:"logged_at,omitempty"`
	RemoteServerId    int64      `json:"remote_server_id,omitempty" path:"remote_server_id,omitempty" url:"remote_server_id,omitempty"`
}

// Identifier returns the resource ID.
func (r RemoteBandwidthSnapshot) Identifier() interface{} {
	return r.Id
}

// RemoteBandwidthSnapshotCollection is a list of RemoteBandwidthSnapshot resources.
type RemoteBandwidthSnapshotCollection []RemoteBandwidthSnapshot

// RemoteBandwidthSnapshotListParams contains the request parameters for GET /remote_bandwidth_snapshots.
type RemoteBandwidthSnapshotListParams struct {
	SortBy     interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	Filter     interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	FilterGt   interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterLt   interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (r *RemoteBandwidthSnapshot) UnmarshalJSON(data []byte) error {
	type remoteBandwidthSnapshot RemoteBandwidthSnapshot
	var v remoteBandwidthSnapshot
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*r = RemoteBandwidthSnapshot(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (r *RemoteBandwidthSnapshotCollection) UnmarshalJSON(data []byte) error {
	type remoteBandwidthSnapshots RemoteBandwidthSnapshotCollection
	var v remoteBandwidthSnapshots
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*r = RemoteBandwidthSnapshotCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (r *RemoteBandwidthSnapshotCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*r))
	for i, v := range *r {
		ret[i] = v
	}

	return &ret
}
