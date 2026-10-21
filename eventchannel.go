package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// EventChannel is a Files.com API resource.
type EventChannel struct {
	// Event Channel ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Event Channel name.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Workspace ID. 0 means the default workspace.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Event Channel description.
	Description string `json:"description,omitempty" path:"description,omitempty" url:"description,omitempty"`
	// Whether this Event Channel can dispatch events.
	Enabled *bool `json:"enabled,omitempty" path:"enabled,omitempty" url:"enabled,omitempty"`
	// Whether this Event Channel is the default destination for newly published events.
	DefaultChannel *bool `json:"default_channel,omitempty" path:"default_channel,omitempty" url:"default_channel,omitempty"`
	// Event Channel create date/time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Event Channel update date/time.
	UpdatedAt *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
}

// Identifier returns the resource ID.
func (e EventChannel) Identifier() interface{} {
	return e.Id
}

// EventChannelCollection is a list of EventChannel resources.
type EventChannelCollection []EventChannel

// EventChannelListParams contains the request parameters for GET /event_channels.
type EventChannelListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `name`, `enabled`, `default_channel` or `workspace_id`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `enabled`, `default_channel` or `workspace_id`. Valid field combinations are `[ workspace_id, enabled ]` and `[ workspace_id, default_channel ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// EventChannelFindParams contains the request parameters for GET /event_channels/{id}.
type EventChannelFindParams struct {
	// Event Channel ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// EventChannelCreateParams contains the request parameters for POST /event_channels.
type EventChannelCreateParams struct {
	// Event Channel name. Required.
	Name string `url:"name" json:"name" path:"name"`
	// Workspace ID. 0 means the default workspace.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
	// Event Channel description.
	Description string `url:"description,omitempty" json:"description,omitempty" path:"description"`
	// Whether this Event Channel can dispatch events.
	Enabled *bool `url:"enabled,omitempty" json:"enabled,omitempty" path:"enabled"`
	// Whether this Event Channel is the default destination for newly published events.
	DefaultChannel *bool `url:"default_channel,omitempty" json:"default_channel,omitempty" path:"default_channel"`
}

// EventChannelUpdateParams contains the request parameters for PATCH /event_channels/{id}.
type EventChannelUpdateParams struct {
	// Event Channel ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Event Channel name.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Workspace ID. 0 means the default workspace.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
	// Event Channel description.
	Description string `url:"description,omitempty" json:"description,omitempty" path:"description"`
	// Whether this Event Channel can dispatch events.
	Enabled *bool `url:"enabled,omitempty" json:"enabled,omitempty" path:"enabled"`
	// Whether this Event Channel is the default destination for newly published events.
	DefaultChannel *bool `url:"default_channel,omitempty" json:"default_channel,omitempty" path:"default_channel"`
}

// EventChannelDeleteParams contains the request parameters for DELETE /event_channels/{id}.
type EventChannelDeleteParams struct {
	// Event Channel ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (e *EventChannel) UnmarshalJSON(data []byte) error {
	type eventChannel EventChannel
	var v eventChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*e = EventChannel(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (e *EventChannelCollection) UnmarshalJSON(data []byte) error {
	type eventChannels EventChannelCollection
	var v eventChannels
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*e = EventChannelCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (e *EventChannelCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*e))
	for i, v := range *e {
		ret[i] = v
	}

	return &ret
}
