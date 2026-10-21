package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// ExpectationIncident is a Files.com API resource.
type ExpectationIncident struct {
	// Expectation Incident ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Workspace ID. `0` means the default workspace.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Expectation ID.
	ExpectationId int64 `json:"expectation_id,omitempty" path:"expectation_id,omitempty" url:"expectation_id,omitempty"`
	// Incident status.
	Status string `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	// When the incident was opened.
	OpenedAt *time.Time `json:"opened_at,omitempty" path:"opened_at,omitempty" url:"opened_at,omitempty"`
	// When the most recent failing evaluation contributing to the incident occurred.
	LastFailedAt *time.Time `json:"last_failed_at,omitempty" path:"last_failed_at,omitempty" url:"last_failed_at,omitempty"`
	// When the incident was acknowledged.
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty" path:"acknowledged_at,omitempty" url:"acknowledged_at,omitempty"`
	// When the current snooze expires.
	SnoozedUntil *time.Time `json:"snoozed_until,omitempty" path:"snoozed_until,omitempty" url:"snoozed_until,omitempty"`
	// When the incident was resolved.
	ResolvedAt *time.Time `json:"resolved_at,omitempty" path:"resolved_at,omitempty" url:"resolved_at,omitempty"`
	// Evaluation that first opened the incident.
	OpenedByEvaluationId int64 `json:"opened_by_evaluation_id,omitempty" path:"opened_by_evaluation_id,omitempty" url:"opened_by_evaluation_id,omitempty"`
	// Most recent evaluation linked to the incident.
	LastEvaluationId int64 `json:"last_evaluation_id,omitempty" path:"last_evaluation_id,omitempty" url:"last_evaluation_id,omitempty"`
	// Evaluation that resolved the incident.
	ResolvedByEvaluationId int64 `json:"resolved_by_evaluation_id,omitempty" path:"resolved_by_evaluation_id,omitempty" url:"resolved_by_evaluation_id,omitempty"`
	// Compact incident summary payload.
	Summary interface{} `json:"summary,omitempty" path:"summary,omitempty" url:"summary,omitempty"`
	// Creation time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Last update time.
	UpdatedAt *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
}

// Identifier returns the resource ID.
func (e ExpectationIncident) Identifier() interface{} {
	return e.Id
}

// ExpectationIncidentCollection is a list of ExpectationIncident resources.
type ExpectationIncidentCollection []ExpectationIncident

// ExpectationIncidentListParams contains the request parameters for GET /expectation_incidents.
type ExpectationIncidentListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `created_at` or `expectation_id`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `expectation_id` and `workspace_id`. Valid field combinations are `[ workspace_id, expectation_id ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// ExpectationIncidentFindParams contains the request parameters for GET /expectation_incidents/{id}.
type ExpectationIncidentFindParams struct {
	// Expectation Incident ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// ExpectationIncidentResolveParams contains the request parameters for POST /expectation_incidents/{id}/resolve.
//
// Resolve an expectation incident
type ExpectationIncidentResolveParams struct {
	// Expectation Incident ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// ExpectationIncidentSnoozeParams contains the request parameters for POST /expectation_incidents/{id}/snooze.
//
// Snooze an expectation incident until a specified time
type ExpectationIncidentSnoozeParams struct {
	// Expectation Incident ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Time until which the incident should remain snoozed. Required.
	SnoozedUntil *time.Time `url:"snoozed_until" json:"snoozed_until" path:"snoozed_until"`
}

// ExpectationIncidentAcknowledgeParams contains the request parameters for POST /expectation_incidents/{id}/acknowledge.
//
// Acknowledge an expectation incident
type ExpectationIncidentAcknowledgeParams struct {
	// Expectation Incident ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (e *ExpectationIncident) UnmarshalJSON(data []byte) error {
	type expectationIncident ExpectationIncident
	var v expectationIncident
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*e = ExpectationIncident(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (e *ExpectationIncidentCollection) UnmarshalJSON(data []byte) error {
	type expectationIncidents ExpectationIncidentCollection
	var v expectationIncidents
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*e = ExpectationIncidentCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (e *ExpectationIncidentCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*e))
	for i, v := range *e {
		ret[i] = v
	}

	return &ret
}
