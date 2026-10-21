package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// EventSubscription is a Files.com API resource.
type EventSubscription struct {
	// Event Subscription ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Event Channel ID.
	EventChannelId int64 `json:"event_channel_id,omitempty" path:"event_channel_id,omitempty" url:"event_channel_id,omitempty"`
	// Workspace ID. 0 means the default workspace or site-wide.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// If true, this default-workspace subscription applies to events from all workspaces.
	ApplyToAllWorkspaces *bool `json:"apply_to_all_workspaces,omitempty" path:"apply_to_all_workspaces,omitempty" url:"apply_to_all_workspaces,omitempty"`
	// Event Subscription name.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Custom subject line to use for notification emails.
	Subject string `json:"subject,omitempty" path:"subject,omitempty" url:"subject,omitempty"`
	// Custom message to include in notification emails.
	Message string `json:"message,omitempty" path:"message,omitempty" url:"message,omitempty"`
	// If true, notification email bodies contain only the custom message, omitting event details and the review button. Requires a custom message, defaults to false, and does not affect non-email targets.
	MessageOnly *bool `json:"message_only,omitempty" path:"message_only,omitempty" url:"message_only,omitempty"`
	// Whether this Event Subscription can dispatch events.
	Enabled *bool `json:"enabled,omitempty" path:"enabled,omitempty" url:"enabled,omitempty"`
	// Event type strings matched by this subscription. Blank means all event types. Valid values: `automation_run.canceled.v1`, `automation_run.failure.v1`, `automation_run.failure_will_retry.v1`, `automation_run.partial_failure.v1`, `automation_run.partial_failure_will_retry.v1`, `automation_run.skipped.v1`, `automation_run.success.v1`, `expectation_evaluation.invalid.v1`, `expectation_evaluation.late.v1`, `expectation_evaluation.missing.v1`, `expectation_evaluation.success.v1`, `expectation_incident.acknowledged.v1`, `expectation_incident.open.v1`, `expectation_incident.resolved.v1`, `expectation_incident.snoozed.v1`, `external_event.client_log.failure.v1`, `external_event.client_log.partial_failure.v1`, `external_event.client_log.skipped.v1`, `external_event.client_log.success.v1`, `pending_work_event.failure.v1`, `pending_work_event.partial_failure.v1`, `pending_work_event.skipped.v1`, `pending_work_event.success.v1`, `siem_http_destination_event.failure.v1`, `siem_http_destination_event.partial_failure.v1`, `siem_http_destination_event.skipped.v1`, `siem_http_destination_event.success.v1`, `sso_event.ldap_login.failure.v1`, `sso_event.ldap_login.partial_failure.v1`, `sso_event.ldap_login.skipped.v1`, `sso_event.ldap_login.success.v1`, `sso_event.ldap_sync.failure.v1`, `sso_event.ldap_sync.partial_failure.v1`, `sso_event.ldap_sync.skipped.v1`, `sso_event.ldap_sync.success.v1`, `sso_event.saml_login.failure.v1`, `sso_event.saml_login.partial_failure.v1`, `sso_event.saml_login.skipped.v1`, `sso_event.saml_login.success.v1`, `sync_run.failure.v1`, `sync_run.partial_failure.v1`, `sync_run.skipped.v1`, `sync_run.success.v1`, `user_security_event.lockout.v1`.
	EventTypes []string `json:"event_types,omitempty" path:"event_types,omitempty" url:"event_types,omitempty"`
	// Conditions on event fields that an event must meet to be delivered, keyed by dot-notation field path. All conditions must match.
	Filter interface{} `json:"filter,omitempty" path:"filter,omitempty" url:"filter,omitempty"`
	// Event Subscription delivery policy.
	DeliveryPolicy interface{} `json:"delivery_policy,omitempty" path:"delivery_policy,omitempty" url:"delivery_policy,omitempty"`
	// Event Target IDs this subscription sends to.
	EventTargetIds []int64 `json:"event_target_ids,omitempty" path:"event_target_ids,omitempty" url:"event_target_ids,omitempty"`
	// Event Subscription create date/time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Event Subscription update date/time.
	UpdatedAt *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
}

// Identifier returns the resource ID.
func (e EventSubscription) Identifier() interface{} {
	return e.Id
}

// EventSubscriptionCollection is a list of EventSubscription resources.
type EventSubscriptionCollection []EventSubscription

// EventSubscriptionListParams contains the request parameters for GET /event_subscriptions.
type EventSubscriptionListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `name`, `enabled`, `event_channel_id` or `workspace_id`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `enabled`, `event_channel_id` or `workspace_id`. Valid field combinations are `[ enabled, event_channel_id ]`, `[ workspace_id, enabled ]`, `[ workspace_id, event_channel_id ]` or `[ workspace_id, enabled, event_channel_id ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// EventSubscriptionFindParams contains the request parameters for GET /event_subscriptions/{id}.
type EventSubscriptionFindParams struct {
	// Event Subscription ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// EventSubscriptionCreateParams contains the request parameters for POST /event_subscriptions.
type EventSubscriptionCreateParams struct {
	// Event Channel ID.
	EventChannelId int64 `url:"event_channel_id,omitempty" json:"event_channel_id,omitempty" path:"event_channel_id"`
	// Workspace ID. 0 means the default workspace or site-wide.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
	// If true, this default-workspace subscription applies to events from all workspaces.
	ApplyToAllWorkspaces *bool `url:"apply_to_all_workspaces,omitempty" json:"apply_to_all_workspaces,omitempty" path:"apply_to_all_workspaces"`
	// Event Subscription name. Required.
	Name string `url:"name" json:"name" path:"name"`
	// Custom subject line to use for notification emails.
	Subject string `url:"subject,omitempty" json:"subject,omitempty" path:"subject"`
	// Custom message to include in notification emails.
	Message string `url:"message,omitempty" json:"message,omitempty" path:"message"`
	// If true, notification email bodies contain only the custom message, omitting event details and the review button. Requires a custom message, defaults to false, and does not affect non-email targets.
	MessageOnly *bool `url:"message_only,omitempty" json:"message_only,omitempty" path:"message_only"`
	// Whether this Event Subscription can dispatch events.
	Enabled *bool `url:"enabled,omitempty" json:"enabled,omitempty" path:"enabled"`
	// Event type strings matched by this subscription. Blank means all event types. Valid values: `automation_run.canceled.v1`, `automation_run.failure.v1`, `automation_run.failure_will_retry.v1`, `automation_run.partial_failure.v1`, `automation_run.partial_failure_will_retry.v1`, `automation_run.skipped.v1`, `automation_run.success.v1`, `expectation_evaluation.invalid.v1`, `expectation_evaluation.late.v1`, `expectation_evaluation.missing.v1`, `expectation_evaluation.success.v1`, `expectation_incident.acknowledged.v1`, `expectation_incident.open.v1`, `expectation_incident.resolved.v1`, `expectation_incident.snoozed.v1`, `external_event.client_log.failure.v1`, `external_event.client_log.partial_failure.v1`, `external_event.client_log.skipped.v1`, `external_event.client_log.success.v1`, `pending_work_event.failure.v1`, `pending_work_event.partial_failure.v1`, `pending_work_event.skipped.v1`, `pending_work_event.success.v1`, `siem_http_destination_event.failure.v1`, `siem_http_destination_event.partial_failure.v1`, `siem_http_destination_event.skipped.v1`, `siem_http_destination_event.success.v1`, `sso_event.ldap_login.failure.v1`, `sso_event.ldap_login.partial_failure.v1`, `sso_event.ldap_login.skipped.v1`, `sso_event.ldap_login.success.v1`, `sso_event.ldap_sync.failure.v1`, `sso_event.ldap_sync.partial_failure.v1`, `sso_event.ldap_sync.skipped.v1`, `sso_event.ldap_sync.success.v1`, `sso_event.saml_login.failure.v1`, `sso_event.saml_login.partial_failure.v1`, `sso_event.saml_login.skipped.v1`, `sso_event.saml_login.success.v1`, `sync_run.failure.v1`, `sync_run.partial_failure.v1`, `sync_run.skipped.v1`, `sync_run.success.v1`, `user_security_event.lockout.v1`.
	EventTypes []string `url:"event_types,omitempty" json:"event_types,omitempty" path:"event_types"`
	// Conditions on event fields that an event must meet to be delivered, keyed by dot-notation field path. All conditions must match.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// Event Subscription delivery policy.
	DeliveryPolicy interface{} `url:"delivery_policy,omitempty" json:"delivery_policy,omitempty" path:"delivery_policy"`
	// Event Target IDs this subscription sends to.
	EventTargetIds []int64 `url:"event_target_ids,omitempty" json:"event_target_ids,omitempty" path:"event_target_ids"`
}

// EventSubscriptionUpdateParams contains the request parameters for PATCH /event_subscriptions/{id}.
type EventSubscriptionUpdateParams struct {
	// Event Subscription ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Event Channel ID.
	EventChannelId int64 `url:"event_channel_id,omitempty" json:"event_channel_id,omitempty" path:"event_channel_id"`
	// Workspace ID. 0 means the default workspace or site-wide.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
	// If true, this default-workspace subscription applies to events from all workspaces.
	ApplyToAllWorkspaces *bool `url:"apply_to_all_workspaces,omitempty" json:"apply_to_all_workspaces,omitempty" path:"apply_to_all_workspaces"`
	// Event Subscription name.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Custom subject line to use for notification emails.
	Subject string `url:"subject,omitempty" json:"subject,omitempty" path:"subject"`
	// Custom message to include in notification emails.
	Message string `url:"message,omitempty" json:"message,omitempty" path:"message"`
	// If true, notification email bodies contain only the custom message, omitting event details and the review button. Requires a custom message, defaults to false, and does not affect non-email targets.
	MessageOnly *bool `url:"message_only,omitempty" json:"message_only,omitempty" path:"message_only"`
	// Whether this Event Subscription can dispatch events.
	Enabled *bool `url:"enabled,omitempty" json:"enabled,omitempty" path:"enabled"`
	// Event type strings matched by this subscription. Blank means all event types. Valid values: `automation_run.canceled.v1`, `automation_run.failure.v1`, `automation_run.failure_will_retry.v1`, `automation_run.partial_failure.v1`, `automation_run.partial_failure_will_retry.v1`, `automation_run.skipped.v1`, `automation_run.success.v1`, `expectation_evaluation.invalid.v1`, `expectation_evaluation.late.v1`, `expectation_evaluation.missing.v1`, `expectation_evaluation.success.v1`, `expectation_incident.acknowledged.v1`, `expectation_incident.open.v1`, `expectation_incident.resolved.v1`, `expectation_incident.snoozed.v1`, `external_event.client_log.failure.v1`, `external_event.client_log.partial_failure.v1`, `external_event.client_log.skipped.v1`, `external_event.client_log.success.v1`, `pending_work_event.failure.v1`, `pending_work_event.partial_failure.v1`, `pending_work_event.skipped.v1`, `pending_work_event.success.v1`, `siem_http_destination_event.failure.v1`, `siem_http_destination_event.partial_failure.v1`, `siem_http_destination_event.skipped.v1`, `siem_http_destination_event.success.v1`, `sso_event.ldap_login.failure.v1`, `sso_event.ldap_login.partial_failure.v1`, `sso_event.ldap_login.skipped.v1`, `sso_event.ldap_login.success.v1`, `sso_event.ldap_sync.failure.v1`, `sso_event.ldap_sync.partial_failure.v1`, `sso_event.ldap_sync.skipped.v1`, `sso_event.ldap_sync.success.v1`, `sso_event.saml_login.failure.v1`, `sso_event.saml_login.partial_failure.v1`, `sso_event.saml_login.skipped.v1`, `sso_event.saml_login.success.v1`, `sync_run.failure.v1`, `sync_run.partial_failure.v1`, `sync_run.skipped.v1`, `sync_run.success.v1`, `user_security_event.lockout.v1`.
	EventTypes []string `url:"event_types,omitempty" json:"event_types,omitempty" path:"event_types"`
	// Conditions on event fields that an event must meet to be delivered, keyed by dot-notation field path. All conditions must match.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// Event Subscription delivery policy.
	DeliveryPolicy interface{} `url:"delivery_policy,omitempty" json:"delivery_policy,omitempty" path:"delivery_policy"`
	// Event Target IDs this subscription sends to.
	EventTargetIds []int64 `url:"event_target_ids,omitempty" json:"event_target_ids,omitempty" path:"event_target_ids"`
}

// EventSubscriptionDeleteParams contains the request parameters for DELETE /event_subscriptions/{id}.
type EventSubscriptionDeleteParams struct {
	// Event Subscription ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (e *EventSubscription) UnmarshalJSON(data []byte) error {
	type eventSubscription EventSubscription
	var v eventSubscription
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*e = EventSubscription(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (e *EventSubscriptionCollection) UnmarshalJSON(data []byte) error {
	type eventSubscriptions EventSubscriptionCollection
	var v eventSubscriptions
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*e = EventSubscriptionCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (e *EventSubscriptionCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*e))
	for i, v := range *e {
		ret[i] = v
	}

	return &ret
}
