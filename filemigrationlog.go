package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type FileMigrationLog struct {
	Timestamp       *time.Time `json:"timestamp,omitempty" path:"timestamp,omitempty" url:"timestamp,omitempty"`
	FileMigrationId int64      `json:"file_migration_id,omitempty" path:"file_migration_id,omitempty" url:"file_migration_id,omitempty"`
	DestPath        string     `json:"dest_path,omitempty" path:"dest_path,omitempty" url:"dest_path,omitempty"`
	ErrorType       string     `json:"error_type,omitempty" path:"error_type,omitempty" url:"error_type,omitempty"`
	Message         string     `json:"message,omitempty" path:"message,omitempty" url:"message,omitempty"`
	Operation       string     `json:"operation,omitempty" path:"operation,omitempty" url:"operation,omitempty"`
	Path            string     `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	Status          string     `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
}

func (f FileMigrationLog) Identifier() interface{} {
	return f.Path
}

type FileMigrationLogCollection []FileMigrationLog

type FileMigrationLogListParams struct {
	Action       string                 `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	Filter       FileMigrationLog       `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterPrefix map[string]interface{} `url:"filter_prefix,omitempty" required:"false" json:"filter_prefix,omitempty" path:"filter_prefix"`
	ListParams
}

func (f *FileMigrationLog) UnmarshalJSON(data []byte) error {
	type fileMigrationLog FileMigrationLog
	var v fileMigrationLog
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*f = FileMigrationLog(v)
	return nil
}

func (f *FileMigrationLogCollection) UnmarshalJSON(data []byte) error {
	type fileMigrationLogs FileMigrationLogCollection
	var v fileMigrationLogs
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*f = FileMigrationLogCollection(v)
	return nil
}

func (f *FileMigrationLogCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*f))
	for i, v := range *f {
		ret[i] = v
	}

	return &ret
}
