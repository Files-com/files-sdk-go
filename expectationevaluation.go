package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// ExpectationEvaluation is a Files.com API resource.
type ExpectationEvaluation struct {
	// ExpectationEvaluation ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Workspace ID. `0` means the default workspace.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Expectation ID.
	ExpectationId int64 `json:"expectation_id,omitempty" path:"expectation_id,omitempty" url:"expectation_id,omitempty"`
	// Evaluation status.
	Status string `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	// How the evaluation window was opened.
	OpenedVia string `json:"opened_via,omitempty" path:"opened_via,omitempty" url:"opened_via,omitempty"`
	// When the evaluation row was opened.
	OpenedAt *time.Time `json:"opened_at,omitempty" path:"opened_at,omitempty" url:"opened_at,omitempty"`
	// Logical start of the candidate window.
	WindowStartAt *time.Time `json:"window_start_at,omitempty" path:"window_start_at,omitempty" url:"window_start_at,omitempty"`
	// Actual candidate cutoff boundary for the window.
	WindowEndAt *time.Time `json:"window_end_at,omitempty" path:"window_end_at,omitempty" url:"window_end_at,omitempty"`
	// Logical due boundary for schedule-driven windows.
	DeadlineAt *time.Time `json:"deadline_at,omitempty" path:"deadline_at,omitempty" url:"deadline_at,omitempty"`
	// Logical cutoff for late acceptance, when present.
	LateAcceptanceCutoffAt *time.Time `json:"late_acceptance_cutoff_at,omitempty" path:"late_acceptance_cutoff_at,omitempty" url:"late_acceptance_cutoff_at,omitempty"`
	// Hard stop after which the window may not remain open.
	HardCloseAt *time.Time `json:"hard_close_at,omitempty" path:"hard_close_at,omitempty" url:"hard_close_at,omitempty"`
	// When the evaluation row was finalized.
	ClosedAt *time.Time `json:"closed_at,omitempty" path:"closed_at,omitempty" url:"closed_at,omitempty"`
	// Captured evidence for files that matched the window.
	MatchedFiles []map[string]interface{} `json:"matched_files,omitempty" path:"matched_files,omitempty" url:"matched_files,omitempty"`
	// Captured evidence for required files that were missing.
	MissingFiles []map[string]interface{} `json:"missing_files,omitempty" path:"missing_files,omitempty" url:"missing_files,omitempty"`
	// Captured criteria failures for the window, including file-scoped FTS errors and returned validation details.
	CriteriaErrors []string `json:"criteria_errors,omitempty" path:"criteria_errors,omitempty" url:"criteria_errors,omitempty"`
	// Compact evaluator summary payload.
	Summary interface{} `json:"summary,omitempty" path:"summary,omitempty" url:"summary,omitempty"`
	// Creation time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Last update time.
	UpdatedAt *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
}

// Identifier returns the resource ID.
func (e ExpectationEvaluation) Identifier() interface{} {
	return e.Id
}

// ExpectationEvaluationCollection is a list of ExpectationEvaluation resources.
type ExpectationEvaluationCollection []ExpectationEvaluation

// ExpectationEvaluationListParams contains the request parameters for GET /expectation_evaluations.
type ExpectationEvaluationListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `created_at` or `expectation_id`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `expectation_id` and `workspace_id`. Valid field combinations are `[ workspace_id, expectation_id ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// ExpectationEvaluationFindParams contains the request parameters for GET /expectation_evaluations/{id}.
type ExpectationEvaluationFindParams struct {
	// Expectation Evaluation ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (e *ExpectationEvaluation) UnmarshalJSON(data []byte) error {
	type expectationEvaluation ExpectationEvaluation
	var v expectationEvaluation
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*e = ExpectationEvaluation(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (e *ExpectationEvaluationCollection) UnmarshalJSON(data []byte) error {
	type expectationEvaluations ExpectationEvaluationCollection
	var v expectationEvaluations
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*e = ExpectationEvaluationCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (e *ExpectationEvaluationCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*e))
	for i, v := range *e {
		ret[i] = v
	}

	return &ret
}
