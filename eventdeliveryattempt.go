package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// EventDeliveryAttempt is a Files.com API resource.
type EventDeliveryAttempt struct {
	Id                  int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	EventRecordId       int64      `json:"event_record_id,omitempty" path:"event_record_id,omitempty" url:"event_record_id,omitempty"`
	EventSubscriptionId int64      `json:"event_subscription_id,omitempty" path:"event_subscription_id,omitempty" url:"event_subscription_id,omitempty"`
	EventTargetId       int64      `json:"event_target_id,omitempty" path:"event_target_id,omitempty" url:"event_target_id,omitempty"`
	WorkspaceId         int64      `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	Status              string     `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	AttemptNumber       int64      `json:"attempt_number,omitempty" path:"attempt_number,omitempty" url:"attempt_number,omitempty"`
	ResponseCode        int64      `json:"response_code,omitempty" path:"response_code,omitempty" url:"response_code,omitempty"`
	ErrorMessage        string     `json:"error_message,omitempty" path:"error_message,omitempty" url:"error_message,omitempty"`
	ResponseBody        string     `json:"response_body,omitempty" path:"response_body,omitempty" url:"response_body,omitempty"`
	LatencyMs           int64      `json:"latency_ms,omitempty" path:"latency_ms,omitempty" url:"latency_ms,omitempty"`
	DeliveredAt         *time.Time `json:"delivered_at,omitempty" path:"delivered_at,omitempty" url:"delivered_at,omitempty"`
	LastAttemptedAt     *time.Time `json:"last_attempted_at,omitempty" path:"last_attempted_at,omitempty" url:"last_attempted_at,omitempty"`
	NextAttemptAt       *time.Time `json:"next_attempt_at,omitempty" path:"next_attempt_at,omitempty" url:"next_attempt_at,omitempty"`
	CreatedAt           *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
}

// Identifier returns the resource ID.
func (e EventDeliveryAttempt) Identifier() interface{} {
	return e.Id
}

// EventDeliveryAttemptCollection is a list of EventDeliveryAttempt resources.
type EventDeliveryAttemptCollection []EventDeliveryAttempt

// EventDeliveryAttemptListParams contains the request parameters for GET /event_delivery_attempts.
type EventDeliveryAttemptListParams struct {
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	ListParams
}

// EventDeliveryAttemptFindParams contains the request parameters for GET /event_delivery_attempts/{id}.
type EventDeliveryAttemptFindParams struct {
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
