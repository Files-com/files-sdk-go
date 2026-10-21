package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// AutomationLog is a Files.com API resource.
type AutomationLog struct {
	// Start Time of Action. Deprecrated: Use created_at.
	Timestamp *time.Time `json:"timestamp,omitempty" path:"timestamp,omitempty" url:"timestamp,omitempty"`
	// Automation ID.
	AutomationId int64 `json:"automation_id,omitempty" path:"automation_id,omitempty" url:"automation_id,omitempty"`
	// Automation Run ID.
	AutomationRunId int64 `json:"automation_run_id,omitempty" path:"automation_run_id,omitempty" url:"automation_run_id,omitempty"`
	// Logical item correlation ID.
	CorrelationId string `json:"correlation_id,omitempty" path:"correlation_id,omitempty" url:"correlation_id,omitempty"`
	// Destination path, for moves and copies.
	DestPath string `json:"dest_path,omitempty" path:"dest_path,omitempty" url:"dest_path,omitempty"`
	// Typed error, if applicable.
	ErrorType string `json:"error_type,omitempty" path:"error_type,omitempty" url:"error_type,omitempty"`
	// ID of the Automation whose definition owns this node.
	SourceAutomationId int64 `json:"source_automation_id,omitempty" path:"source_automation_id,omitempty" url:"source_automation_id,omitempty"`
	// Pinned version of the Automation whose definition owns this node.
	SourceAutomationVersion int64 `json:"source_automation_version,omitempty" path:"source_automation_version,omitempty" url:"source_automation_version,omitempty"`
	// Message.
	Message string `json:"message,omitempty" path:"message,omitempty" url:"message,omitempty"`
	// Node ID.
	NodeId string `json:"node_id,omitempty" path:"node_id,omitempty" url:"node_id,omitempty"`
	// Node type.
	NodeType string `json:"node_type,omitempty" path:"node_type,omitempty" url:"node_type,omitempty"`
	// Operation type.
	Operation string `json:"operation,omitempty" path:"operation,omitempty" url:"operation,omitempty"`
	// Node outlet used by the item.
	Output string `json:"output,omitempty" path:"output,omitempty" url:"output,omitempty"`
	// File path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.
	Path string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	// Status.
	Status string `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	// Start Time of Action.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
}

// Identifier returns the resource path.
func (a AutomationLog) Identifier() interface{} {
	return a.Path
}

// AutomationLogCollection is a list of AutomationLog resources.
type AutomationLogCollection []AutomationLog

// AutomationLogListParams contains the request parameters for GET /automation_logs.
type AutomationLogListParams struct {
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `automation_id`, `automation_run_id`, `correlation_id`, `error_type`, `node_id`, `node_type`, `operation`, `output`, `path`, `source_automation_id`, `source_automation_version`, `status` or `created_at`. Valid field combinations are every combination of one or more of the 13 fields `automation_id`, `automation_run_id`, `correlation_id`, `error_type`, `node_id`, `node_type`, `operation`, `output`, `path`, `source_automation_id`, `source_automation_version`, `status` and `created_at`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is greater than the supplied value. Valid fields are `source_automation_version` and `created_at`. Valid field combinations are every combination of one or more of the 13 fields `automation_id`, `automation_run_id`, `correlation_id`, `error_type`, `node_id`, `node_type`, `operation`, `output`, `path`, `source_automation_id`, `source_automation_version`, `status` and `created_at`.
	FilterGt interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	// If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `source_automation_version` and `created_at`. Valid field combinations are every combination of one or more of the 13 fields `automation_id`, `automation_run_id`, `correlation_id`, `error_type`, `node_id`, `node_type`, `operation`, `output`, `path`, `source_automation_id`, `source_automation_version`, `status` and `created_at`.
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	// If set, return records where the specified field is prefixed by the supplied value. Valid fields are `correlation_id`, `error_type`, `node_id`, `node_type`, `output` or `path`. Valid field combinations are every combination of one or more of the 13 fields `automation_id`, `automation_run_id`, `correlation_id`, `error_type`, `node_id`, `node_type`, `operation`, `output`, `path`, `source_automation_id`, `source_automation_version`, `status` and `created_at`.
	FilterPrefix interface{} `url:"filter_prefix,omitempty" json:"filter_prefix,omitempty" path:"filter_prefix"`
	// If set, return records where the specified field is less than the supplied value. Valid fields are `source_automation_version` and `created_at`. Valid field combinations are every combination of one or more of the 13 fields `automation_id`, `automation_run_id`, `correlation_id`, `error_type`, `node_id`, `node_type`, `operation`, `output`, `path`, `source_automation_id`, `source_automation_version`, `status` and `created_at`.
	FilterLt interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	// If set, return records where the specified field is less than or equal the supplied value. Valid fields are `source_automation_version` and `created_at`. Valid field combinations are every combination of one or more of the 13 fields `automation_id`, `automation_run_id`, `correlation_id`, `error_type`, `node_id`, `node_type`, `operation`, `output`, `path`, `source_automation_id`, `source_automation_version`, `status` and `created_at`.
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (a *AutomationLog) UnmarshalJSON(data []byte) error {
	type automationLog AutomationLog
	var v automationLog
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = AutomationLog(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (a *AutomationLogCollection) UnmarshalJSON(data []byte) error {
	type automationLogs AutomationLogCollection
	var v automationLogs
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = AutomationLogCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (a *AutomationLogCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
