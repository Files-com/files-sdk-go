package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// EventRecord is a Files.com API resource.
type EventRecord struct {
	// Event Record ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Workspace ID. 0 means the default workspace or site-wide.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Stable event UUID.
	EventUuid string `json:"event_uuid,omitempty" path:"event_uuid,omitempty" url:"event_uuid,omitempty"`
	// Versioned event type string. Valid values: `automation_run.canceled.v1`, `automation_run.failure.v1`, `automation_run.failure_will_retry.v1`, `automation_run.partial_failure.v1`, `automation_run.partial_failure_will_retry.v1`, `automation_run.skipped.v1`, `automation_run.success.v1`, `expectation_evaluation.invalid.v1`, `expectation_evaluation.late.v1`, `expectation_evaluation.missing.v1`, `expectation_evaluation.success.v1`, `expectation_incident.acknowledged.v1`, `expectation_incident.open.v1`, `expectation_incident.resolved.v1`, `expectation_incident.snoozed.v1`, `external_event.client_log.failure.v1`, `external_event.client_log.partial_failure.v1`, `external_event.client_log.skipped.v1`, `external_event.client_log.success.v1`, `pending_work_event.failure.v1`, `pending_work_event.partial_failure.v1`, `pending_work_event.skipped.v1`, `pending_work_event.success.v1`, `siem_http_destination_event.failure.v1`, `siem_http_destination_event.partial_failure.v1`, `siem_http_destination_event.skipped.v1`, `siem_http_destination_event.success.v1`, `sso_event.ldap_login.failure.v1`, `sso_event.ldap_login.partial_failure.v1`, `sso_event.ldap_login.skipped.v1`, `sso_event.ldap_login.success.v1`, `sso_event.ldap_sync.failure.v1`, `sso_event.ldap_sync.partial_failure.v1`, `sso_event.ldap_sync.skipped.v1`, `sso_event.ldap_sync.success.v1`, `sso_event.saml_login.failure.v1`, `sso_event.saml_login.partial_failure.v1`, `sso_event.saml_login.skipped.v1`, `sso_event.saml_login.success.v1`, `sync_run.failure.v1`, `sync_run.partial_failure.v1`, `sync_run.skipped.v1`, `sync_run.success.v1`, `user_security_event.lockout.v1`.
	EventType string `json:"event_type,omitempty" path:"event_type,omitempty" url:"event_type,omitempty"`
	// Event severity.
	Severity string `json:"severity,omitempty" path:"severity,omitempty" url:"severity,omitempty"`
	// Source record type.
	SourceType string `json:"source_type,omitempty" path:"source_type,omitempty" url:"source_type,omitempty"`
	// Source record ID.
	SourceId int64 `json:"source_id,omitempty" path:"source_id,omitempty" url:"source_id,omitempty"`
	// Event occurrence date/time.
	OccurredAt *time.Time `json:"occurred_at,omitempty" path:"occurred_at,omitempty" url:"occurred_at,omitempty"`
	// Human-readable event title.
	HumanTitle string `json:"human_title,omitempty" path:"human_title,omitempty" url:"human_title,omitempty"`
	// Human-readable event summary.
	HumanSummary string `json:"human_summary,omitempty" path:"human_summary,omitempty" url:"human_summary,omitempty"`
	// Human-readable event detail fields.
	HumanFields []map[string]interface{} `json:"human_fields,omitempty" path:"human_fields,omitempty" url:"human_fields,omitempty"`
	// Actor associated with the event.
	Actor interface{} `json:"actor,omitempty" path:"actor,omitempty" url:"actor,omitempty"`
	// Resources associated with the event.
	Resources []map[string]interface{} `json:"resources,omitempty" path:"resources,omitempty" url:"resources,omitempty"`
	// Event payload.
	Payload interface{} `json:"payload,omitempty" path:"payload,omitempty" url:"payload,omitempty"`
	// Event Record create date/time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
}

// Identifier returns the resource ID.
func (e EventRecord) Identifier() interface{} {
	return e.Id
}

// EventRecordCollection is a list of EventRecord resources.
type EventRecordCollection []EventRecord

// EventRecordListParams contains the request parameters for GET /event_records.
type EventRecordListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `event_type`, `created_at` or `workspace_id`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at`, `event_type` or `workspace_id`. Valid field combinations are `[ event_type, created_at ]`, `[ workspace_id, created_at ]`, `[ workspace_id, event_type ]` or `[ workspace_id, event_type, created_at ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`.
	FilterGt interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	// If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`.
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	// If set, return records where the specified field is prefixed by the supplied value. Valid fields are `event_type`.
	FilterPrefix interface{} `url:"filter_prefix,omitempty" json:"filter_prefix,omitempty" path:"filter_prefix"`
	// If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`.
	FilterLt interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	// If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`.
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// EventRecordFindParams contains the request parameters for GET /event_records/{id}.
type EventRecordFindParams struct {
	// Event Record ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (e *EventRecord) UnmarshalJSON(data []byte) error {
	type eventRecord EventRecord
	var v eventRecord
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*e = EventRecord(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (e *EventRecordCollection) UnmarshalJSON(data []byte) error {
	type eventRecords EventRecordCollection
	var v eventRecords
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*e = EventRecordCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (e *EventRecordCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*e))
	for i, v := range *e {
		ret[i] = v
	}

	return &ret
}
