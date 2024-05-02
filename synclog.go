package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type SyncLog struct {
	Timestamp       *time.Time `json:"timestamp,omitempty" path:"timestamp,omitempty" url:"timestamp,omitempty"`
	SyncId          int64      `json:"sync_id,omitempty" path:"sync_id,omitempty" url:"sync_id,omitempty"`
	ExternalEventId int64      `json:"external_event_id,omitempty" path:"external_event_id,omitempty" url:"external_event_id,omitempty"`
	ErrorType       string     `json:"error_type,omitempty" path:"error_type,omitempty" url:"error_type,omitempty"`
	Message         string     `json:"message,omitempty" path:"message,omitempty" url:"message,omitempty"`
	Operation       string     `json:"operation,omitempty" path:"operation,omitempty" url:"operation,omitempty"`
	Path            string     `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	Size            string     `json:"size,omitempty" path:"size,omitempty" url:"size,omitempty"`
	FileType        string     `json:"file_type,omitempty" path:"file_type,omitempty" url:"file_type,omitempty"`
	Status          string     `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
}

func (s SyncLog) Identifier() interface{} {
	return s.Path
}

type SyncLogCollection []SyncLog

type SyncLogListParams struct {
	Filter       SyncLog                `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterPrefix map[string]interface{} `url:"filter_prefix,omitempty" required:"false" json:"filter_prefix,omitempty" path:"filter_prefix"`
	ListParams
}

func (s *SyncLog) UnmarshalJSON(data []byte) error {
	type syncLog SyncLog
	var v syncLog
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = SyncLog(v)
	return nil
}

func (s *SyncLogCollection) UnmarshalJSON(data []byte) error {
	type syncLogs SyncLogCollection
	var v syncLogs
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SyncLogCollection(v)
	return nil
}

func (s *SyncLogCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
