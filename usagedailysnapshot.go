package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/lpar/date"
)

type UsageDailySnapshot struct {
	Id                           int64     `json:"id,omitempty"`
	Date                         date.Date `json:"date,omitempty"`
	ApiUsageAvailable            *bool     `json:"api_usage_available,omitempty"`
	ReadApiUsage                 int64     `json:"read_api_usage,omitempty"`
	WriteApiUsage                int64     `json:"write_api_usage,omitempty"`
	UserCount                    int64     `json:"user_count,omitempty"`
	CurrentStorage               int64     `json:"current_storage,omitempty"`
	DeletedFilesStorage          int64     `json:"deleted_files_storage,omitempty"`
	DeletedFilesCountedInMinimum int64     `json:"deleted_files_counted_in_minimum,omitempty"`
	RootStorage                  int64     `json:"root_storage,omitempty"`
	UsageByTopLevelDir           []string  `json:"usage_by_top_level_dir,omitempty"`
}

type UsageDailySnapshotCollection []UsageDailySnapshot

type UsageDailySnapshotListParams struct {
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

func (u *UsageDailySnapshot) UnmarshalJSON(data []byte) error {
	type usageDailySnapshot UsageDailySnapshot
	var v usageDailySnapshot
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*u = UsageDailySnapshot(v)
	return nil
}

func (u *UsageDailySnapshotCollection) UnmarshalJSON(data []byte) error {
	type usageDailySnapshots []UsageDailySnapshot
	var v usageDailySnapshots
	if err := json.Unmarshal(data, &v); err != nil {
		return err
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
