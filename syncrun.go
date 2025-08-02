package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type SyncRun struct {
	Id                   int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	SyncId               int64      `json:"sync_id,omitempty" path:"sync_id,omitempty" url:"sync_id,omitempty"`
	SiteId               int64      `json:"site_id,omitempty" path:"site_id,omitempty" url:"site_id,omitempty"`
	Status               string     `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	SrcRemoteServerType  string     `json:"src_remote_server_type,omitempty" path:"src_remote_server_type,omitempty" url:"src_remote_server_type,omitempty"`
	DestRemoteServerType string     `json:"dest_remote_server_type,omitempty" path:"dest_remote_server_type,omitempty" url:"dest_remote_server_type,omitempty"`
	Body                 string     `json:"body,omitempty" path:"body,omitempty" url:"body,omitempty"`
	EventErrors          []string   `json:"event_errors,omitempty" path:"event_errors,omitempty" url:"event_errors,omitempty"`
	ComparedFiles        int64      `json:"compared_files,omitempty" path:"compared_files,omitempty" url:"compared_files,omitempty"`
	ComparedFolders      int64      `json:"compared_folders,omitempty" path:"compared_folders,omitempty" url:"compared_folders,omitempty"`
	ErroredFiles         int64      `json:"errored_files,omitempty" path:"errored_files,omitempty" url:"errored_files,omitempty"`
	SuccessfulFiles      int64      `json:"successful_files,omitempty" path:"successful_files,omitempty" url:"successful_files,omitempty"`
	Runtime              string     `json:"runtime,omitempty" path:"runtime,omitempty" url:"runtime,omitempty"`
	LogUrl               string     `json:"log_url,omitempty" path:"log_url,omitempty" url:"log_url,omitempty"`
	CompletedAt          *time.Time `json:"completed_at,omitempty" path:"completed_at,omitempty" url:"completed_at,omitempty"`
	Notified             *bool      `json:"notified,omitempty" path:"notified,omitempty" url:"notified,omitempty"`
	DryRun               *bool      `json:"dry_run,omitempty" path:"dry_run,omitempty" url:"dry_run,omitempty"`
	BytesSynced          int64      `json:"bytes_synced,omitempty" path:"bytes_synced,omitempty" url:"bytes_synced,omitempty"`
	EstimatedBytesCount  int64      `json:"estimated_bytes_count,omitempty" path:"estimated_bytes_count,omitempty" url:"estimated_bytes_count,omitempty"`
	CreatedAt            *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	UpdatedAt            *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
}

func (s SyncRun) Identifier() interface{} {
	return s.Id
}

type SyncRunCollection []SyncRun

type SyncRunListParams struct {
	UserId int64                  `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	SortBy map[string]interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	Filter SyncRun                `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
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
