package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type UsageSnapshot struct {
	Id                           int64           `json:"id,omitempty" path:"id"`
	StartAt                      *time.Time      `json:"start_at,omitempty" path:"start_at"`
	EndAt                        *time.Time      `json:"end_at,omitempty" path:"end_at"`
	HighWaterUserCount           string          `json:"high_water_user_count,omitempty" path:"high_water_user_count"`
	CurrentStorage               string          `json:"current_storage,omitempty" path:"current_storage"`
	HighWaterStorage             string          `json:"high_water_storage,omitempty" path:"high_water_storage"`
	UsageByTopLevelDir           json.RawMessage `json:"usage_by_top_level_dir,omitempty" path:"usage_by_top_level_dir"`
	RootStorage                  string          `json:"root_storage,omitempty" path:"root_storage"`
	DeletedFilesCountedInMinimum string          `json:"deleted_files_counted_in_minimum,omitempty" path:"deleted_files_counted_in_minimum"`
	DeletedFilesStorage          string          `json:"deleted_files_storage,omitempty" path:"deleted_files_storage"`
	TotalBillableUsage           string          `json:"total_billable_usage,omitempty" path:"total_billable_usage"`
	TotalBillableTransferUsage   string          `json:"total_billable_transfer_usage,omitempty" path:"total_billable_transfer_usage"`
	BytesSent                    string          `json:"bytes_sent,omitempty" path:"bytes_sent"`
	SyncBytesReceived            string          `json:"sync_bytes_received,omitempty" path:"sync_bytes_received"`
	SyncBytesSent                string          `json:"sync_bytes_sent,omitempty" path:"sync_bytes_sent"`
}

func (u UsageSnapshot) Identifier() interface{} {
	return u.Id
}

type UsageSnapshotCollection []UsageSnapshot

type UsageSnapshotListParams struct {
	ListParams
}

func (u *UsageSnapshot) UnmarshalJSON(data []byte) error {
	type usageSnapshot UsageSnapshot
	var v usageSnapshot
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*u = UsageSnapshot(v)
	return nil
}

func (u *UsageSnapshotCollection) UnmarshalJSON(data []byte) error {
	type usageSnapshots UsageSnapshotCollection
	var v usageSnapshots
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*u = UsageSnapshotCollection(v)
	return nil
}

func (u *UsageSnapshotCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*u))
	for i, v := range *u {
		ret[i] = v
	}

	return &ret
}
