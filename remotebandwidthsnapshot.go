package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// RemoteBandwidthSnapshot is a Files.com API resource.
type RemoteBandwidthSnapshot struct {
	// Site bandwidth ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Site sync bandwidth report bytes received.
	SyncBytesReceived int64 `json:"sync_bytes_received,omitempty" path:"sync_bytes_received,omitempty" url:"sync_bytes_received,omitempty"`
	// Site sync bandwidth report bytes sent.
	SyncBytesSent int64 `json:"sync_bytes_sent,omitempty" path:"sync_bytes_sent,omitempty" url:"sync_bytes_sent,omitempty"`
	// Time the site bandwidth report was logged.
	LoggedAt *time.Time `json:"logged_at,omitempty" path:"logged_at,omitempty" url:"logged_at,omitempty"`
	// ID of related Remote Server.
	RemoteServerId int64 `json:"remote_server_id,omitempty" path:"remote_server_id,omitempty" url:"remote_server_id,omitempty"`
}

// Identifier returns the resource ID.
func (r RemoteBandwidthSnapshot) Identifier() interface{} {
	return r.Id
}

// RemoteBandwidthSnapshotCollection is a list of RemoteBandwidthSnapshot resources.
type RemoteBandwidthSnapshotCollection []RemoteBandwidthSnapshot

// RemoteBandwidthSnapshotListParams contains the request parameters for GET /remote_bandwidth_snapshots.
type RemoteBandwidthSnapshotListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `logged_at`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `logged_at`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is greater than the supplied value. Valid fields are `logged_at`.
	FilterGt interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	// If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `logged_at`.
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	// If set, return records where the specified field is less than the supplied value. Valid fields are `logged_at`.
	FilterLt interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	// If set, return records where the specified field is less than or equal the supplied value. Valid fields are `logged_at`.
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
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
