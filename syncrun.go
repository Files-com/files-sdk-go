package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type SyncRun struct {
	Id                 int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	SyncId             int64      `json:"sync_id,omitempty" path:"sync_id,omitempty" url:"sync_id,omitempty"`
	SiteId             int64      `json:"site_id,omitempty" path:"site_id,omitempty" url:"site_id,omitempty"`
	Status             string     `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	RemoteServerType   string     `json:"remote_server_type,omitempty" path:"remote_server_type,omitempty" url:"remote_server_type,omitempty"`
	Body               string     `json:"body,omitempty" path:"body,omitempty" url:"body,omitempty"`
	EventErrors        []string   `json:"event_errors,omitempty" path:"event_errors,omitempty" url:"event_errors,omitempty"`
	BytesSynced        int64      `json:"bytes_synced,omitempty" path:"bytes_synced,omitempty" url:"bytes_synced,omitempty"`
	ComparedFiles      int64      `json:"compared_files,omitempty" path:"compared_files,omitempty" url:"compared_files,omitempty"`
	ComparedFolders    int64      `json:"compared_folders,omitempty" path:"compared_folders,omitempty" url:"compared_folders,omitempty"`
	ErroredFiles       int64      `json:"errored_files,omitempty" path:"errored_files,omitempty" url:"errored_files,omitempty"`
	SuccessfulFiles    int64      `json:"successful_files,omitempty" path:"successful_files,omitempty" url:"successful_files,omitempty"`
	Runtime            string     `json:"runtime,omitempty" path:"runtime,omitempty" url:"runtime,omitempty"`
	S3BodyPath         string     `json:"s3_body_path,omitempty" path:"s3_body_path,omitempty" url:"s3_body_path,omitempty"`
	S3InternalBodyPath string     `json:"s3_internal_body_path,omitempty" path:"s3_internal_body_path,omitempty" url:"s3_internal_body_path,omitempty"`
	CompletedAt        *time.Time `json:"completed_at,omitempty" path:"completed_at,omitempty" url:"completed_at,omitempty"`
	Notified           *bool      `json:"notified,omitempty" path:"notified,omitempty" url:"notified,omitempty"`
	CreatedAt          *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
}

func (s SyncRun) Identifier() interface{} {
	return s.Id
}

type SyncRunCollection []SyncRun

type SyncRunListParams struct {
	UserId int64                  `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	SortBy map[string]interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	Filter SyncRun                `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	SyncId int64                  `url:"sync_id" json:"sync_id" path:"sync_id"`
	ListParams
}

type SyncRunFindParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

func (s *SyncRun) UnmarshalJSON(data []byte) error {
	type syncRun SyncRun
	var v syncRun
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = SyncRun(v)
	return nil
}

func (s *SyncRunCollection) UnmarshalJSON(data []byte) error {
	type syncRuns SyncRunCollection
	var v syncRuns
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SyncRunCollection(v)
	return nil
}

func (s *SyncRunCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
