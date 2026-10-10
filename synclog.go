package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// SyncLog is a Files.com API resource.
type SyncLog struct {
	Timestamp       *time.Time `json:"timestamp,omitempty" path:"timestamp,omitempty" url:"timestamp,omitempty"`
	SyncId          int64      `json:"sync_id,omitempty" path:"sync_id,omitempty" url:"sync_id,omitempty"`
	ExternalEventId int64      `json:"external_event_id,omitempty" path:"external_event_id,omitempty" url:"external_event_id,omitempty"`
	SyncRunId       int64      `json:"sync_run_id,omitempty" path:"sync_run_id,omitempty" url:"sync_run_id,omitempty"`
	ErrorType       string     `json:"error_type,omitempty" path:"error_type,omitempty" url:"error_type,omitempty"`
	Message         string     `json:"message,omitempty" path:"message,omitempty" url:"message,omitempty"`
	Operation       string     `json:"operation,omitempty" path:"operation,omitempty" url:"operation,omitempty"`
	Path            string     `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	Size            int64      `json:"size,omitempty" path:"size,omitempty" url:"size,omitempty"`
	FileType        string     `json:"file_type,omitempty" path:"file_type,omitempty" url:"file_type,omitempty"`
	Status          string     `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
}

// Identifier returns the resource path.
func (s SyncLog) Identifier() interface{} {
	return s.Path
}

// SyncLogCollection is a list of SyncLog resources.
type SyncLogCollection []SyncLog

// SyncLogListParams contains the request parameters for GET /sync_logs.
type SyncLogListParams struct {
	Filter     interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	FilterGt   interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterLt   interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (s *SyncLog) UnmarshalJSON(data []byte) error {
	type syncLog SyncLog
	var v syncLog
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = SyncLog(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (s *SyncLogCollection) UnmarshalJSON(data []byte) error {
	type syncLogs SyncLogCollection
	var v syncLogs
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SyncLogCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (s *SyncLogCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
