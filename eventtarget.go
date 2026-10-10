package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// EventTarget is a Files.com API resource.
type EventTarget struct {
	Id                   int64       `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Name                 string      `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	TargetType           string      `json:"target_type,omitempty" path:"target_type,omitempty" url:"target_type,omitempty"`
	WorkspaceId          int64       `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	ApplyToAllWorkspaces *bool       `json:"apply_to_all_workspaces,omitempty" path:"apply_to_all_workspaces,omitempty" url:"apply_to_all_workspaces,omitempty"`
	Enabled              *bool       `json:"enabled,omitempty" path:"enabled,omitempty" url:"enabled,omitempty"`
	Config               interface{} `json:"config,omitempty" path:"config,omitempty" url:"config,omitempty"`
	DeliveryPolicy       interface{} `json:"delivery_policy,omitempty" path:"delivery_policy,omitempty" url:"delivery_policy,omitempty"`
	CreatedAt            *time.Time  `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	UpdatedAt            *time.Time  `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
}

// Identifier returns the resource ID.
func (e EventTarget) Identifier() interface{} {
	return e.Id
}

// EventTargetCollection is a list of EventTarget resources.
type EventTargetCollection []EventTarget

// EventTargetTargetTypeEnum is a string value for target_type.
// Enum lists the values documented by the API.
type EventTargetTargetTypeEnum string

// String returns the API parameter value.
func (u EventTargetTargetTypeEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u EventTargetTargetTypeEnum) Enum() map[string]EventTargetTargetTypeEnum {
	return map[string]EventTargetTargetTypeEnum{
		"email":         EventTargetTargetTypeEnum("email"),
		"webhook":       EventTargetTargetTypeEnum("webhook"),
		"slack_webhook": EventTargetTargetTypeEnum("slack_webhook"),
		"teams_webhook": EventTargetTargetTypeEnum("teams_webhook"),
		"amazon_sns":    EventTargetTargetTypeEnum("amazon_sns"),
		"google_pubsub": EventTargetTargetTypeEnum("google_pubsub"),
		"folder":        EventTargetTargetTypeEnum("folder"),
	}
}

// EventTargetListParams contains the request parameters for GET /event_targets.
type EventTargetListParams struct {
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	ListParams
}

// EventTargetFindParams contains the request parameters for GET /event_targets/{id}.
type EventTargetFindParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// EventTargetCreateParams contains the request parameters for POST /event_targets.
type EventTargetCreateParams struct {
	Name                 string                    `url:"name" json:"name" path:"name"`
	WorkspaceId          int64                     `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
	ApplyToAllWorkspaces *bool                     `url:"apply_to_all_workspaces,omitempty" json:"apply_to_all_workspaces,omitempty" path:"apply_to_all_workspaces"`
	Enabled              *bool                     `url:"enabled,omitempty" json:"enabled,omitempty" path:"enabled"`
	Config               interface{}               `url:"config" json:"config" path:"config"`
	DeliveryPolicy       interface{}               `url:"delivery_policy,omitempty" json:"delivery_policy,omitempty" path:"delivery_policy"`
	TargetType           EventTargetTargetTypeEnum `url:"target_type" json:"target_type" path:"target_type"`
}

// EventTargetUpdateParams contains the request parameters for PATCH /event_targets/{id}.
type EventTargetUpdateParams struct {
	Id                   int64       `url:"-,omitempty" json:"-,omitempty" path:"id"`
	Name                 string      `url:"name,omitempty" json:"name,omitempty" path:"name"`
	WorkspaceId          int64       `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
	ApplyToAllWorkspaces *bool       `url:"apply_to_all_workspaces,omitempty" json:"apply_to_all_workspaces,omitempty" path:"apply_to_all_workspaces"`
	Enabled              *bool       `url:"enabled,omitempty" json:"enabled,omitempty" path:"enabled"`
	Config               interface{} `url:"config,omitempty" json:"config,omitempty" path:"config"`
	DeliveryPolicy       interface{} `url:"delivery_policy,omitempty" json:"delivery_policy,omitempty" path:"delivery_policy"`
}

// EventTargetDeleteParams contains the request parameters for DELETE /event_targets/{id}.
type EventTargetDeleteParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (e *EventTarget) UnmarshalJSON(data []byte) error {
	type eventTarget EventTarget
	var v eventTarget
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*e = EventTarget(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (e *EventTargetCollection) UnmarshalJSON(data []byte) error {
	type eventTargets EventTargetCollection
	var v eventTargets
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*e = EventTargetCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (e *EventTargetCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*e))
	for i, v := range *e {
		ret[i] = v
	}

	return &ret
}
