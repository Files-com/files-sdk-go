package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// EventDeliveryAttempt is a Files.com API resource.
type EventDeliveryAttempt struct {
	// Event Delivery Attempt ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Event Record ID.
	EventRecordId int64 `json:"event_record_id,omitempty" path:"event_record_id,omitempty" url:"event_record_id,omitempty"`
	// Event Subscription ID.
	EventSubscriptionId int64 `json:"event_subscription_id,omitempty" path:"event_subscription_id,omitempty" url:"event_subscription_id,omitempty"`
	// Event Target ID.
	EventTargetId int64 `json:"event_target_id,omitempty" path:"event_target_id,omitempty" url:"event_target_id,omitempty"`
	// Shared summary ID when this event is delivered in a Slack or Teams subfolder summary. All events in a summary share its delivery outcome.
	ActivityNotificationBatchId int64 `json:"activity_notification_batch_id,omitempty" path:"activity_notification_batch_id,omitempty" url:"activity_notification_batch_id,omitempty"`
	// Workspace ID. 0 means the default workspace or site-wide.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Delivery status.
	Status string `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	// Number of delivery attempts made.
	AttemptNumber int64 `json:"attempt_number,omitempty" path:"attempt_number,omitempty" url:"attempt_number,omitempty"`
	// HTTP response code, if applicable.
	ResponseCode int64 `json:"response_code,omitempty" path:"response_code,omitempty" url:"response_code,omitempty"`
	// Delivery error message, if applicable.
	ErrorMessage string `json:"error_message,omitempty" path:"error_message,omitempty" url:"error_message,omitempty"`
	// Delivery response body, if applicable.
	ResponseBody string `json:"response_body,omitempty" path:"response_body,omitempty" url:"response_body,omitempty"`
	// Delivery latency in milliseconds.
	LatencyMs int64 `json:"latency_ms,omitempty" path:"latency_ms,omitempty" url:"latency_ms,omitempty"`
	// Successful delivery date/time.
	DeliveredAt *time.Time `json:"delivered_at,omitempty" path:"delivered_at,omitempty" url:"delivered_at,omitempty"`
	// Most recent attempt date/time.
	LastAttemptedAt *time.Time `json:"last_attempted_at,omitempty" path:"last_attempted_at,omitempty" url:"last_attempted_at,omitempty"`
	// Next scheduled attempt date/time.
	NextAttemptAt *time.Time `json:"next_attempt_at,omitempty" path:"next_attempt_at,omitempty" url:"next_attempt_at,omitempty"`
	// Delivery Attempt create date/time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
}

// Identifier returns the resource ID.
func (e EventDeliveryAttempt) Identifier() interface{} {
	return e.Id
}

// EventDeliveryAttemptCollection is a list of EventDeliveryAttempt resources.
type EventDeliveryAttemptCollection []EventDeliveryAttempt

// EventDeliveryAttemptListParams contains the request parameters for GET /event_delivery_attempts.
type EventDeliveryAttemptListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`, `status`, `event_record_id`, `event_target_id` or `workspace_id`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `status`, `workspace_id`, `event_record_id` or `event_target_id`. Valid field combinations are `[ workspace_id, status ]`, `[ workspace_id, event_record_id ]` or `[ workspace_id, event_target_id ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// EventDeliveryAttemptFindParams contains the request parameters for GET /event_delivery_attempts/{id}.
type EventDeliveryAttemptFindParams struct {
	// Event Delivery Attempt ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (e *EventDeliveryAttempt) UnmarshalJSON(data []byte) error {
	type eventDeliveryAttempt EventDeliveryAttempt
	var v eventDeliveryAttempt
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*e = EventDeliveryAttempt(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (e *EventDeliveryAttemptCollection) UnmarshalJSON(data []byte) error {
	type eventDeliveryAttempts EventDeliveryAttemptCollection
	var v eventDeliveryAttempts
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*e = EventDeliveryAttemptCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (e *EventDeliveryAttemptCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*e))
	for i, v := range *e {
		ret[i] = v
	}

	return &ret
}
