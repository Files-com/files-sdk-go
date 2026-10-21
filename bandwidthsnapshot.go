package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// BandwidthSnapshot is a Files.com API resource.
type BandwidthSnapshot struct {
	// Site bandwidth ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Site bandwidth report bytes received.
	BytesReceived int64 `json:"bytes_received,omitempty" path:"bytes_received,omitempty" url:"bytes_received,omitempty"`
	// Site bandwidth report bytes sent.
	BytesSent int64 `json:"bytes_sent,omitempty" path:"bytes_sent,omitempty" url:"bytes_sent,omitempty"`
	// Site sync bandwidth report bytes received.
	SyncBytesReceived int64 `json:"sync_bytes_received,omitempty" path:"sync_bytes_received,omitempty" url:"sync_bytes_received,omitempty"`
	// Site sync bandwidth report bytes sent.
	SyncBytesSent int64 `json:"sync_bytes_sent,omitempty" path:"sync_bytes_sent,omitempty" url:"sync_bytes_sent,omitempty"`
	// Site bandwidth report get requests.
	RequestsGet int64 `json:"requests_get,omitempty" path:"requests_get,omitempty" url:"requests_get,omitempty"`
	// Site bandwidth report put requests.
	RequestsPut int64 `json:"requests_put,omitempty" path:"requests_put,omitempty" url:"requests_put,omitempty"`
	// Site bandwidth report other requests.
	RequestsOther int64 `json:"requests_other,omitempty" path:"requests_other,omitempty" url:"requests_other,omitempty"`
	// Time the site bandwidth report was logged.
	LoggedAt *time.Time `json:"logged_at,omitempty" path:"logged_at,omitempty" url:"logged_at,omitempty"`
}

// Identifier returns the resource ID.
func (b BandwidthSnapshot) Identifier() interface{} {
	return b.Id
}

// BandwidthSnapshotCollection is a list of BandwidthSnapshot resources.
type BandwidthSnapshotCollection []BandwidthSnapshot

// BandwidthSnapshotListParams contains the request parameters for GET /bandwidth_snapshots.
type BandwidthSnapshotListParams struct {
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
func (b *BandwidthSnapshot) UnmarshalJSON(data []byte) error {
	type bandwidthSnapshot BandwidthSnapshot
	var v bandwidthSnapshot
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*b = BandwidthSnapshot(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (b *BandwidthSnapshotCollection) UnmarshalJSON(data []byte) error {
	type bandwidthSnapshots BandwidthSnapshotCollection
	var v bandwidthSnapshots
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*b = BandwidthSnapshotCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (b *BandwidthSnapshotCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*b))
	for i, v := range *b {
		ret[i] = v
	}

	return &ret
}
