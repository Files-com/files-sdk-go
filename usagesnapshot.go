package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type UsageSnapshot struct {
	Id                           int64           `json:"id,omitempty"`
	StartAt                      time.Time       `json:"start_at,omitempty"`
	EndAt                        time.Time       `json:"end_at,omitempty"`
	CreatedAt                    time.Time       `json:"created_at,omitempty"`
	HighWaterUserCount           float32         `json:"high_water_user_count,omitempty"`
	CurrentStorage               float32         `json:"current_storage,omitempty"`
	HighWaterStorage             float32         `json:"high_water_storage,omitempty"`
	TotalDownloads               int64           `json:"total_downloads,omitempty"`
	TotalUploads                 int64           `json:"total_uploads,omitempty"`
	UpdatedAt                    time.Time       `json:"updated_at,omitempty"`
	UsageByTopLevelDir           json.RawMessage `json:"usage_by_top_level_dir,omitempty"`
	RootStorage                  float32         `json:"root_storage,omitempty"`
	DeletedFilesCountedInMinimum float32         `json:"deleted_files_counted_in_minimum,omitempty"`
	DeletedFilesStorage          float32         `json:"deleted_files_storage,omitempty"`
	TotalBillableUsage           float32         `json:"total_billable_usage,omitempty"`
	TotalBillableTransferUsage   float32         `json:"total_billable_transfer_usage,omitempty"`
	BytesSent                    float32         `json:"bytes_sent,omitempty"`
	SyncBytesReceived            float32         `json:"sync_bytes_received,omitempty"`
	SyncBytesSent                float32         `json:"sync_bytes_sent,omitempty"`
}

type UsageSnapshotCollection []UsageSnapshot

type UsageSnapshotListParams struct {
	Cursor  string `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage int64  `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	lib.ListParams
}

func (u *UsageSnapshot) UnmarshalJSON(data []byte) error {
	type usageSnapshot UsageSnapshot
	var v usageSnapshot
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*u = UsageSnapshot(v)
	return nil
}

func (u *UsageSnapshotCollection) UnmarshalJSON(data []byte) error {
	type usageSnapshots []UsageSnapshot
	var v usageSnapshots
	if err := json.Unmarshal(data, &v); err != nil {
		return err
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
