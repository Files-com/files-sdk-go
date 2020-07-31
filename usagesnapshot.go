package files_sdk

import (
  lib "github.com/Files-com/files-sdk-go/lib"
  "encoding/json"
  "time"
)

type UsageSnapshot struct {
  Id int `json:"id,omitempty"`
  StartAt time.Time `json:"start_at,omitempty"`
  EndAt time.Time `json:"end_at,omitempty"`
  CreatedAt time.Time `json:"created_at,omitempty"`
  CurrentStorage float32 `json:"current_storage,omitempty"`
  HighWaterStorage float32 `json:"high_water_storage,omitempty"`
  TotalDownloads int `json:"total_downloads,omitempty"`
  TotalUploads int `json:"total_uploads,omitempty"`
  UpdatedAt time.Time `json:"updated_at,omitempty"`
  UsageByTopLevelDir json.RawMessage `json:"usage_by_top_level_dir,omitempty"`
  RootStorage float32 `json:"root_storage,omitempty"`
  DeletedFilesCountedInMinimum float32 `json:"deleted_files_counted_in_minimum,omitempty"`
  DeletedFilesStorage float32 `json:"deleted_files_storage,omitempty"`
}

type UsageSnapshotCollection []UsageSnapshot

type UsageSnapshotListParams struct {
  Page int `url:"page,omitempty"`
  PerPage int `url:"per_page,omitempty"`
  Action string `url:"action,omitempty"`
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

