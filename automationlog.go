package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type AutomationLog struct {
	Timestamp       *time.Time `json:"timestamp,omitempty" path:"timestamp,omitempty" url:"timestamp,omitempty"`
	AutomationId    int64      `json:"automation_id,omitempty" path:"automation_id,omitempty" url:"automation_id,omitempty"`
	AutomationRunId int64      `json:"automation_run_id,omitempty" path:"automation_run_id,omitempty" url:"automation_run_id,omitempty"`
	DestPath        string     `json:"dest_path,omitempty" path:"dest_path,omitempty" url:"dest_path,omitempty"`
	ErrorType       string     `json:"error_type,omitempty" path:"error_type,omitempty" url:"error_type,omitempty"`
	Message         string     `json:"message,omitempty" path:"message,omitempty" url:"message,omitempty"`
	Operation       string     `json:"operation,omitempty" path:"operation,omitempty" url:"operation,omitempty"`
	Path            string     `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	Status          string     `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
}

func (a AutomationLog) Identifier() interface{} {
	return a.Path
}

type AutomationLogCollection []AutomationLog

type AutomationLogListParams struct {
	Action       string                 `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	Filter       AutomationLog          `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterPrefix map[string]interface{} `url:"filter_prefix,omitempty" required:"false" json:"filter_prefix,omitempty" path:"filter_prefix"`
	ListParams
}

func (a *AutomationLog) UnmarshalJSON(data []byte) error {
	type automationLog AutomationLog
	var v automationLog
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = AutomationLog(v)
	return nil
}

func (a *AutomationLogCollection) UnmarshalJSON(data []byte) error {
	type automationLogs AutomationLogCollection
	var v automationLogs
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = AutomationLogCollection(v)
	return nil
}

func (a *AutomationLogCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
