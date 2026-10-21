package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/lpar/calendar"
)

// UsageDailySnapshot is a Files.com API resource.
type UsageDailySnapshot struct {
	// ID of the usage record.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// The date of this usage record.
	Date *calendar.Date `json:"date,omitempty" path:"date,omitempty" url:"date,omitempty"`
	// True if the API usage fields `read_api_usage` and `write_api_usage` can be relied upon. If this is false, we suggest hiding that value from any UI.
	ApiUsageAvailable *bool `json:"api_usage_available,omitempty" path:"api_usage_available,omitempty" url:"api_usage_available,omitempty"`
	// Read API Calls used on this day. Note: only updated for days before the current day.
	ReadApiUsage int64 `json:"read_api_usage,omitempty" path:"read_api_usage,omitempty" url:"read_api_usage,omitempty"`
	// Write API Calls used on this day. Note: only updated for days before the current day.
	WriteApiUsage int64 `json:"write_api_usage,omitempty" path:"write_api_usage,omitempty" url:"write_api_usage,omitempty"`
	// Transformation and AI credits used on this day. Note: only updated for days before the current day.
	TransformationCreditsUsage int64 `json:"transformation_credits_usage,omitempty" path:"transformation_credits_usage,omitempty" url:"transformation_credits_usage,omitempty"`
	// Number of billable users as of this day.
	UserCount int64 `json:"user_count,omitempty" path:"user_count,omitempty" url:"user_count,omitempty"`
	// GB of Files Native Storage used on this day.
	CurrentStorage int64 `json:"current_storage,omitempty" path:"current_storage,omitempty" url:"current_storage,omitempty"`
	// GB of Files Native Storage used on this day for files that have been deleted and are stored as backups.
	DeletedFilesStorage int64 `json:"deleted_files_storage,omitempty" path:"deleted_files_storage,omitempty" url:"deleted_files_storage,omitempty"`
	// GB of Files Native Storage used on this day for files that have been permanently deleted but were uploaded less than 30 days ago, and are still billable.
	DeletedFilesCountedInMinimum int64 `json:"deleted_files_counted_in_minimum,omitempty" path:"deleted_files_counted_in_minimum,omitempty" url:"deleted_files_counted_in_minimum,omitempty"`
	// GB of Files Native Storage used for the root folder. Included here because this value will not be part of `usage_by_top_level_dir`.
	RootStorage int64 `json:"root_storage,omitempty" path:"root_storage,omitempty" url:"root_storage,omitempty"`
	// Usage broken down by each top-level folder.
	UsageByTopLevelDir []UsageByTopLevelDir `json:"usage_by_top_level_dir,omitempty" path:"usage_by_top_level_dir,omitempty" url:"usage_by_top_level_dir,omitempty"`
}

// Identifier returns the resource ID.
func (u UsageDailySnapshot) Identifier() interface{} {
	return u.Id
}

// UsageDailySnapshotCollection is a list of UsageDailySnapshot resources.
type UsageDailySnapshotCollection []UsageDailySnapshot

// UsageDailySnapshotListParams contains the request parameters for GET /usage_daily_snapshots.
type UsageDailySnapshotListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `date`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `date` and `usage_snapshot_id`. Valid field combinations are `[ usage_snapshot_id, date ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is greater than the supplied value. Valid fields are `date`.
	FilterGt interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	// If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `date`.
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	// If set, return records where the specified field is less than the supplied value. Valid fields are `date`.
	FilterLt interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	// If set, return records where the specified field is less than or equal the supplied value. Valid fields are `date`.
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (u *UsageDailySnapshot) UnmarshalJSON(data []byte) error {
	type usageDailySnapshot UsageDailySnapshot
	var v usageDailySnapshot
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*u = UsageDailySnapshot(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (u *UsageDailySnapshotCollection) UnmarshalJSON(data []byte) error {
	type usageDailySnapshots UsageDailySnapshotCollection
	var v usageDailySnapshots
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*u = UsageDailySnapshotCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (u *UsageDailySnapshotCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*u))
	for i, v := range *u {
		ret[i] = v
	}

	return &ret
}
