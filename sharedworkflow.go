package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// SharedWorkflow is a Files.com API resource.
type SharedWorkflow struct {
	// Shared Workflow ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Name displayed to recipients.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Description displayed to recipients.
	Description string `json:"description,omitempty" path:"description,omitempty" url:"description,omitempty"`
	// Revision required when submitting a Shared Workflow run.
	SubmissionRevision string `json:"submission_revision,omitempty" path:"submission_revision,omitempty" url:"submission_revision,omitempty"`
	// Current custom questions in display order.
	FormFields []FormField `json:"form_fields,omitempty" path:"form_fields,omitempty" url:"form_fields,omitempty"`
}

// Identifier returns the resource ID.
func (s SharedWorkflow) Identifier() interface{} {
	return s.Id
}

// SharedWorkflowCollection is a list of SharedWorkflow resources.
type SharedWorkflowCollection []SharedWorkflow

// SharedWorkflowListParams contains the request parameters for GET /shared_workflows.
type SharedWorkflowListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are .
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// SharedWorkflowFindParams contains the request parameters for GET /shared_workflows/{id}.
type SharedWorkflowFindParams struct {
	// Shared Workflow ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// SharedWorkflowFindRunParams contains the request parameters for GET /shared_workflows/{id}/runs/{run_id}.
type SharedWorkflowFindRunParams struct {
	// Shared Workflow ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// ID of the run returned by the Shared Workflow submission. Required; sent in the request path.
	RunId int64 `url:"-,omitempty" json:"-,omitempty" path:"run_id"`
}

// SharedWorkflowRunWorkflowParams contains the request parameters for POST /shared_workflows/{id}/run_workflow.
//
// Run a Shared Workflow
type SharedWorkflowRunWorkflowParams struct {
	// Shared Workflow ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Revision returned when opening the Shared Workflow. Required.
	SubmissionRevision string `url:"submission_revision" json:"submission_revision" path:"submission_revision"`
	// Client-generated UUID used to retry the same submission without starting another run. Required.
	SubmissionId string `url:"submission_id" json:"submission_id" path:"submission_id"`
	// Path of a file the recipient can read, including a file uploaded to a Snapshot created for this workflow. Required.
	FilePath string `url:"file_path" json:"file_path" path:"file_path"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (s *SharedWorkflow) UnmarshalJSON(data []byte) error {
	type sharedWorkflow SharedWorkflow
	var v sharedWorkflow
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = SharedWorkflow(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (s *SharedWorkflowCollection) UnmarshalJSON(data []byte) error {
	type sharedWorkflows SharedWorkflowCollection
	var v sharedWorkflows
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SharedWorkflowCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (s *SharedWorkflowCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
