package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/lpar/date"
)

type UsageDailySnapshot struct {
	Id                           int64           `json:"id,omitempty" path:"id"`
	Date                         *date.Date      `json:"date,omitempty" path:"date"`
	ApiUsageAvailable            *bool           `json:"api_usage_available,omitempty" path:"api_usage_available"`
	ReadApiUsage                 int64           `json:"read_api_usage,omitempty" path:"read_api_usage"`
	WriteApiUsage                int64           `json:"write_api_usage,omitempty" path:"write_api_usage"`
	UserCount                    int64           `json:"user_count,omitempty" path:"user_count"`
	CurrentStorage               int64           `json:"current_storage,omitempty" path:"current_storage"`
	DeletedFilesStorage          int64           `json:"deleted_files_storage,omitempty" path:"deleted_files_storage"`
	DeletedFilesCountedInMinimum int64           `json:"deleted_files_counted_in_minimum,omitempty" path:"deleted_files_counted_in_minimum"`
	RootStorage                  int64           `json:"root_storage,omitempty" path:"root_storage"`
	UsageByTopLevelDir           json.RawMessage `json:"usage_by_top_level_dir,omitempty" path:"usage_by_top_level_dir"`
}

type UsageDailySnapshotCollection []UsageDailySnapshot

type UsageDailySnapshotListParams struct {
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty" path:"filter_lteq"`
	lib.ListParams
}

func (u *UsageDailySnapshot) UnmarshalJSON(data []byte) error {
	type usageDailySnapshot UsageDailySnapshot
	var v usageDailySnapshot
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*u = UsageDailySnapshot(v)
	return nil
}

func (u *UsageDailySnapshotCollection) UnmarshalJSON(data []byte) error {
	type usageDailySnapshots UsageDailySnapshotCollection
	var v usageDailySnapshots
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*u = UsageDailySnapshotCollection(v)
	return nil
}

func (u *UsageDailySnapshotCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*u))
	for i, v := range *u {
		ret[i] = v
	}

	return &ret
}
