package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// AutomationRun is a Files.com API resource.
type AutomationRun struct {
	// ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// ID of the associated Automation.
	AutomationId int64 `json:"automation_id,omitempty" path:"automation_id,omitempty" url:"automation_id,omitempty"`
	// ID of the immutable Automation version pinned by this run.
	AutomationVersionId int64 `json:"automation_version_id,omitempty" path:"automation_version_id,omitempty" url:"automation_version_id,omitempty"`
	// Pinned Automation v2 definition version.
	Version int64 `json:"version,omitempty" path:"version,omitempty" url:"version,omitempty"`
	// Workspace ID.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Date/time at which cancellation was requested.
	CancelRequestedAt *time.Time `json:"cancel_requested_at,omitempty" path:"cancel_requested_at,omitempty" url:"cancel_requested_at,omitempty"`
	// Automation run completion/failure date/time.
	CompletedAt *time.Time `json:"completed_at,omitempty" path:"completed_at,omitempty" url:"completed_at,omitempty"`
	// Automation run start date/time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// If set, this automation will be retried at this date/time due to `failure` or `partial_failure`.
	RetryAt *time.Time `json:"retry_at,omitempty" path:"retry_at,omitempty" url:"retry_at,omitempty"`
	// If set, this Automation run was retried due to `failure` or `partial_failure`.
	RetriedAt *time.Time `json:"retried_at,omitempty" path:"retried_at,omitempty" url:"retried_at,omitempty"`
	// ID of the run that is or will be retrying this run.
	RetriedInRunId int64 `json:"retried_in_run_id,omitempty" path:"retried_in_run_id,omitempty" url:"retried_in_run_id,omitempty"`
	// ID of the original run that this run is retrying.
	RetryOfRunId int64 `json:"retry_of_run_id,omitempty" path:"retry_of_run_id,omitempty" url:"retry_of_run_id,omitempty"`
	// ID of the run whose persisted node outputs this run reused.
	RerunOfRunId int64 `json:"rerun_of_run_id,omitempty" path:"rerun_of_run_id,omitempty" url:"rerun_of_run_id,omitempty"`
	// Node at which this run resumed execution.
	RerunFromNodeId string `json:"rerun_from_node_id,omitempty" path:"rerun_from_node_id,omitempty" url:"rerun_from_node_id,omitempty"`
	// Automation run runtime.
	Runtime float64 `json:"runtime,omitempty" path:"runtime,omitempty" url:"runtime,omitempty"`
	// The status of the AutomationRun. One of `queued`, `running`, `success`, `partial_failure`, `failure`, `skipped`, or `canceled`.
	Status string `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	// Count of successful operations.
	SuccessfulOperations int64 `json:"successful_operations,omitempty" path:"successful_operations,omitempty" url:"successful_operations,omitempty"`
	// Count of failed operations.
	FailedOperations int64 `json:"failed_operations,omitempty" path:"failed_operations,omitempty" url:"failed_operations,omitempty"`
	// Automation definition snapshot pinned by this run. For performance reasons, this is not provided when listing Automation runs.
	Definition interface{} `json:"definition,omitempty" path:"definition,omitempty" url:"definition,omitempty"`
	// Status and execution stage for each node in this run. For performance reasons, this is not provided when listing Automation runs.
	NodeStates interface{} `json:"node_states,omitempty" path:"node_states,omitempty" url:"node_states,omitempty"`
	// Execution status, timing, and bounded output summaries for each node. For performance reasons, this is not provided when listing Automation runs.
	ExecutionNodes []AutomationExecutionNode `json:"execution_nodes,omitempty" path:"execution_nodes,omitempty" url:"execution_nodes,omitempty"`
	// Link to the run journal artifact.
	JournalUrl string `json:"journal_url,omitempty" path:"journal_url,omitempty" url:"journal_url,omitempty"`
	// Link to status messages log file.
	StatusMessagesUrl string `json:"status_messages_url,omitempty" path:"status_messages_url,omitempty" url:"status_messages_url,omitempty"`
}

// Identifier returns the resource ID.
func (a AutomationRun) Identifier() interface{} {
	return a.Id
}

// AutomationRunCollection is a list of AutomationRun resources.
type AutomationRunCollection []AutomationRun

// AutomationRunListParams contains the request parameters for GET /automation_runs.
type AutomationRunListParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `automation_id`, `created_at` or `status`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `status`, `workspace_id` or `automation_id`. Valid field combinations are `[ workspace_id, status ]`, `[ automation_id, status ]`, `[ workspace_id, automation_id ]` or `[ workspace_id, automation_id, status ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ID of the associated Automation. Required.
	AutomationId int64 `url:"automation_id" json:"automation_id" path:"automation_id"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// AutomationRunFindParams contains the request parameters for GET /automation_runs/{id}.
type AutomationRunFindParams struct {
	// Automation Run ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// AutomationRunFindNodeParams contains the request parameters for GET /automation_runs/{id}/node.
type AutomationRunFindNodeParams struct {
	// Automation Run ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Node ID from the pinned Automation definition. Required.
	NodeId string `url:"node_id" json:"node_id" path:"node_id"`
}

// AutomationRunCancelParams contains the request parameters for POST /automation_runs/{id}/cancel.
//
// Cancel Automation Run
type AutomationRunCancelParams struct {
	// Automation Run ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// AutomationRunRerunParams contains the request parameters for POST /automation_runs/{id}/rerun.
//
// Re-run Automation from Node
type AutomationRunRerunParams struct {
	// Automation Run ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Node ID at which execution should resume. Required.
	NodeId string `url:"node_id" json:"node_id" path:"node_id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (a *AutomationRun) UnmarshalJSON(data []byte) error {
	type automationRun AutomationRun
	var v automationRun
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = AutomationRun(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (a *AutomationRunCollection) UnmarshalJSON(data []byte) error {
	type automationRuns AutomationRunCollection
	var v automationRuns
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = AutomationRunCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (a *AutomationRunCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
