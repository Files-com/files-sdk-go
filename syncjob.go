package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type SyncJob struct {
	QueuedAt             time.Time `json:"queued_at,omitempty"`
	UpdatedAt            time.Time `json:"updated_at,omitempty"`
	Status               string    `json:"status,omitempty"`
	RegionalWorkerStatus string    `json:"regional_worker_status,omitempty"`
	Uuid                 string    `json:"uuid,omitempty"`
	FolderBehaviorId     int64     `json:"folder_behavior_id,omitempty"`
}

type SyncJobCollection []SyncJob

type SyncJobListParams struct {
	Cursor  string `url:"cursor,omitempty" required:"false"`
	PerPage int64  `url:"per_page,omitempty" required:"false"`
	lib.ListParams
}

func (s *SyncJob) UnmarshalJSON(data []byte) error {
	type syncJob SyncJob
	var v syncJob
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*s = SyncJob(v)
	return nil
}

func (s *SyncJobCollection) UnmarshalJSON(data []byte) error {
	type syncJobs []SyncJob
	var v syncJobs
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*s = SyncJobCollection(v)
	return nil
}

func (s *SyncJobCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
