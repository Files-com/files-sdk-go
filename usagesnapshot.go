package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// UsageSnapshot is a Files.com API resource.
type UsageSnapshot struct {
	// Usage snapshot ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Usage snapshot start date/time.
	StartAt *time.Time `json:"start_at,omitempty" path:"start_at,omitempty" url:"start_at,omitempty"`
	// Usage snapshot end date/time.
	EndAt *time.Time `json:"end_at,omitempty" path:"end_at,omitempty" url:"end_at,omitempty"`
	// Highest user count number in time period.
	HighWaterUserCount int64 `json:"high_water_user_count,omitempty" path:"high_water_user_count,omitempty" url:"high_water_user_count,omitempty"`
	// Total Transformation and AI credits used in this usage period. Updated daily. Does not include credits consumed today.
	TransformationCreditsUsage int64 `json:"transformation_credits_usage,omitempty" path:"transformation_credits_usage,omitempty" url:"transformation_credits_usage,omitempty"`
	// Current total Storage Usage GB as of end date (not necessarily high water mark, which is used for billing).
	CurrentStorage int64 `json:"current_storage,omitempty" path:"current_storage,omitempty" url:"current_storage,omitempty"`
	// Highest Storage Usage GB recorded in time period (used for billing).
	HighWaterStorage int64 `json:"high_water_storage,omitempty" path:"high_water_storage,omitempty" url:"high_water_storage,omitempty"`
	// Storage Usage for root folder as of end date (not necessarily high water mark, which is used for billing).
	RootStorage int64 `json:"root_storage,omitempty" path:"root_storage,omitempty" url:"root_storage,omitempty"`
	// Storage Usage for files that are deleted but uploaded within last 30 days as of end date (not necessarily high water mark, which is used for billing).
	DeletedFilesCountedInMinimum int64 `json:"deleted_files_counted_in_minimum,omitempty" path:"deleted_files_counted_in_minimum,omitempty" url:"deleted_files_counted_in_minimum,omitempty"`
	// Storage Usage for files that are deleted but retained as backups as of end date (not necessarily high water mark, which is used for billing).
	DeletedFilesStorage int64 `json:"deleted_files_storage,omitempty" path:"deleted_files_storage,omitempty" url:"deleted_files_storage,omitempty"`
	// Storage + Transfer Usage - Total Billable amount.
	TotalBillableUsage int64 `json:"total_billable_usage,omitempty" path:"total_billable_usage,omitempty" url:"total_billable_usage,omitempty"`
	// Transfer usage for period - Total Billable amount.
	TotalBillableTransferUsage int64 `json:"total_billable_transfer_usage,omitempty" path:"total_billable_transfer_usage,omitempty" url:"total_billable_transfer_usage,omitempty"`
	// Transfer Usage for period - Outbound GB from Files Native Storage.
	BytesSent int64 `json:"bytes_sent,omitempty" path:"bytes_sent,omitempty" url:"bytes_sent,omitempty"`
	// Transfer Usage for period - Inbound GB to Remote Servers (Sync/Mount).
	SyncBytesReceived int64 `json:"sync_bytes_received,omitempty" path:"sync_bytes_received,omitempty" url:"sync_bytes_received,omitempty"`
	// Transfer Usage for period - Outbound GB from Remote Servers (Sync/Mount).
	SyncBytesSent int64 `json:"sync_bytes_sent,omitempty" path:"sync_bytes_sent,omitempty" url:"sync_bytes_sent,omitempty"`
	// Storage Usage - map of root folders to their usage as of end date (not necessarily high water mark, which is used for billing).
	UsageByTopLevelDir []UsageByTopLevelDir `json:"usage_by_top_level_dir,omitempty" path:"usage_by_top_level_dir,omitempty" url:"usage_by_top_level_dir,omitempty"`
}

// Identifier returns the resource ID.
func (u UsageSnapshot) Identifier() interface{} {
	return u.Id
}

// UsageSnapshotCollection is a list of UsageSnapshot resources.
type UsageSnapshotCollection []UsageSnapshot

// UsageSnapshotListParams contains the request parameters for GET /usage_snapshots.
type UsageSnapshotListParams struct {
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (u *UsageSnapshot) UnmarshalJSON(data []byte) error {
	type usageSnapshot UsageSnapshot
	var v usageSnapshot
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*u = UsageSnapshot(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (u *UsageSnapshotCollection) UnmarshalJSON(data []byte) error {
	type usageSnapshots UsageSnapshotCollection
	var v usageSnapshots
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*u = UsageSnapshotCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (u *UsageSnapshotCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*u))
	for i, v := range *u {
		ret[i] = v
	}

	return &ret
}
