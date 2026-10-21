package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// SyncRun is a Files.com API resource.
type SyncRun struct {
	// SyncRun ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Log or summary body for this run.
	Body string `json:"body,omitempty" path:"body,omitempty" url:"body,omitempty"`
	// Total bytes synced in this run.
	BytesSynced int64 `json:"bytes_synced,omitempty" path:"bytes_synced,omitempty" url:"bytes_synced,omitempty"`
	// Number of files compared.
	ComparedFiles int64 `json:"compared_files,omitempty" path:"compared_files,omitempty" url:"compared_files,omitempty"`
	// Number of folders compared.
	ComparedFolders int64 `json:"compared_folders,omitempty" path:"compared_folders,omitempty" url:"compared_folders,omitempty"`
	// When this run was completed.
	CompletedAt *time.Time `json:"completed_at,omitempty" path:"completed_at,omitempty" url:"completed_at,omitempty"`
	// When this run was created.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Destination remote server type, if any.
	DestRemoteServerType string `json:"dest_remote_server_type,omitempty" path:"dest_remote_server_type,omitempty" url:"dest_remote_server_type,omitempty"`
	// Whether this run was a dry run (no actual changes made).
	DryRun *bool `json:"dry_run,omitempty" path:"dry_run,omitempty" url:"dry_run,omitempty"`
	// Number of files that errored.
	ErroredFiles int64 `json:"errored_files,omitempty" path:"errored_files,omitempty" url:"errored_files,omitempty"`
	// Estimated bytes count for this run.
	EstimatedBytesCount int64 `json:"estimated_bytes_count,omitempty" path:"estimated_bytes_count,omitempty" url:"estimated_bytes_count,omitempty"`
	// Array of errors encountered during the run.
	EventErrors []string `json:"event_errors,omitempty" path:"event_errors,omitempty" url:"event_errors,omitempty"`
	// Link to external log file.
	LogUrl string `json:"log_url,omitempty" path:"log_url,omitempty" url:"log_url,omitempty"`
	// Total runtime in seconds.
	Runtime float64 `json:"runtime,omitempty" path:"runtime,omitempty" url:"runtime,omitempty"`
	// Site ID.
	SiteId int64 `json:"site_id,omitempty" path:"site_id,omitempty" url:"site_id,omitempty"`
	// Workspace ID.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Source remote server type, if any.
	SrcRemoteServerType string `json:"src_remote_server_type,omitempty" path:"src_remote_server_type,omitempty" url:"src_remote_server_type,omitempty"`
	// Status of the sync run (success, failure, partial_failure, in_progress, skipped).
	Status string `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	// Number of files successfully synced.
	SuccessfulFiles int64 `json:"successful_files,omitempty" path:"successful_files,omitempty" url:"successful_files,omitempty"`
	// ID of the Sync this run belongs to.
	SyncId int64 `json:"sync_id,omitempty" path:"sync_id,omitempty" url:"sync_id,omitempty"`
	// Name of the Sync this run belongs to.
	SyncName string `json:"sync_name,omitempty" path:"sync_name,omitempty" url:"sync_name,omitempty"`
	// When this run was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
	// Array of in-progress file transfers with progress data. Only present when the sync run status is in_progress.
	LiveTransfers []SyncRunLiveTransfer `json:"live_transfers,omitempty" path:"live_transfers,omitempty" url:"live_transfers,omitempty"`
}

// Identifier returns the resource ID.
func (s SyncRun) Identifier() interface{} {
	return s.Id
}

// SyncRunCollection is a list of SyncRun resources.
type SyncRunCollection []SyncRun

// SyncRunListParams contains the request parameters for GET /sync_runs.
type SyncRunListParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `site_id`, `workspace_id`, `sync_id`, `created_at` or `status`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at`, `status`, `dry_run`, `workspace_id`, `src_remote_server_type`, `dest_remote_server_type` or `sync_id`. Valid field combinations are `[ status, created_at ]`, `[ workspace_id, created_at ]`, `[ src_remote_server_type, created_at ]`, `[ dest_remote_server_type, created_at ]`, `[ sync_id, created_at ]`, `[ workspace_id, status ]`, `[ src_remote_server_type, status ]`, `[ dest_remote_server_type, status ]`, `[ sync_id, status ]`, `[ workspace_id, src_remote_server_type ]`, `[ workspace_id, dest_remote_server_type ]`, `[ workspace_id, sync_id ]`, `[ workspace_id, status, created_at ]`, `[ src_remote_server_type, status, created_at ]`, `[ dest_remote_server_type, status, created_at ]`, `[ sync_id, status, created_at ]`, `[ workspace_id, src_remote_server_type, created_at ]`, `[ workspace_id, dest_remote_server_type, created_at ]`, `[ workspace_id, sync_id, created_at ]`, `[ workspace_id, src_remote_server_type, status ]`, `[ workspace_id, dest_remote_server_type, status ]`, `[ workspace_id, sync_id, status ]`, `[ workspace_id, src_remote_server_type, status, created_at ]`, `[ workspace_id, dest_remote_server_type, status, created_at ]` or `[ workspace_id, sync_id, status, created_at ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`.
	FilterGt interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	// If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`.
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	// If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`.
	FilterLt interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	// If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`.
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// SyncRunFindParams contains the request parameters for GET /sync_runs/{id}.
type SyncRunFindParams struct {
	// Sync Run ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (s *SyncRun) UnmarshalJSON(data []byte) error {
	type syncRun SyncRun
	var v syncRun
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = SyncRun(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (s *SyncRunCollection) UnmarshalJSON(data []byte) error {
	type syncRuns SyncRunCollection
	var v syncRuns
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SyncRunCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (s *SyncRunCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
