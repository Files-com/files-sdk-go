package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// AiAssistantPersonality is a Files.com API resource.
type AiAssistantPersonality struct {
	// AI Assistant Personality ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Workspace ID. `0` means the default workspace.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// AI Assistant Personality name.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// System prompt injected into the in-app AI Assistant.
	SystemPrompt string `json:"system_prompt,omitempty" path:"system_prompt,omitempty" url:"system_prompt,omitempty"`
	// Whether this personality is the default personality for the Workspace.
	UseByDefault *bool `json:"use_by_default,omitempty" path:"use_by_default,omitempty" url:"use_by_default,omitempty"`
	// If true, this default-workspace personality can apply to users in all workspaces.
	ApplyToAllWorkspaces *bool `json:"apply_to_all_workspaces,omitempty" path:"apply_to_all_workspaces,omitempty" url:"apply_to_all_workspaces,omitempty"`
	// Creation time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Last update time.
	UpdatedAt *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
}

// Identifier returns the resource ID.
func (a AiAssistantPersonality) Identifier() interface{} {
	return a.Id
}

// AiAssistantPersonalityCollection is a list of AiAssistantPersonality resources.
type AiAssistantPersonalityCollection []AiAssistantPersonality

// AiAssistantPersonalityListParams contains the request parameters for GET /ai_assistant_personalities.
type AiAssistantPersonalityListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `id`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// AiAssistantPersonalityFindParams contains the request parameters for GET /ai_assistant_personalities/{id}.
type AiAssistantPersonalityFindParams struct {
	// Ai Assistant Personality ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// AiAssistantPersonalityCreateParams contains the request parameters for POST /ai_assistant_personalities.
type AiAssistantPersonalityCreateParams struct {
	// If true, this default-workspace personality can apply to users in all workspaces.
	ApplyToAllWorkspaces *bool `url:"apply_to_all_workspaces,omitempty" json:"apply_to_all_workspaces,omitempty" path:"apply_to_all_workspaces"`
	// AI Assistant Personality name. Required.
	Name string `url:"name" json:"name" path:"name"`
	// System prompt injected into the in-app AI Assistant. Required.
	SystemPrompt string `url:"system_prompt" json:"system_prompt" path:"system_prompt"`
	// Whether this personality is the default personality for the Workspace.
	UseByDefault *bool `url:"use_by_default,omitempty" json:"use_by_default,omitempty" path:"use_by_default"`
	// Workspace ID. `0` means the default workspace.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// AiAssistantPersonalityUpdateParams contains the request parameters for PATCH /ai_assistant_personalities/{id}.
type AiAssistantPersonalityUpdateParams struct {
	// Ai Assistant Personality ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// If true, this default-workspace personality can apply to users in all workspaces.
	ApplyToAllWorkspaces *bool `url:"apply_to_all_workspaces,omitempty" json:"apply_to_all_workspaces,omitempty" path:"apply_to_all_workspaces"`
	// AI Assistant Personality name.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// System prompt injected into the in-app AI Assistant.
	SystemPrompt string `url:"system_prompt,omitempty" json:"system_prompt,omitempty" path:"system_prompt"`
	// Whether this personality is the default personality for the Workspace.
	UseByDefault *bool `url:"use_by_default,omitempty" json:"use_by_default,omitempty" path:"use_by_default"`
	// Workspace ID. `0` means the default workspace.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// AiAssistantPersonalityDeleteParams contains the request parameters for DELETE /ai_assistant_personalities/{id}.
type AiAssistantPersonalityDeleteParams struct {
	// Ai Assistant Personality ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (a *AiAssistantPersonality) UnmarshalJSON(data []byte) error {
	type aiAssistantPersonality AiAssistantPersonality
	var v aiAssistantPersonality
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = AiAssistantPersonality(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (a *AiAssistantPersonalityCollection) UnmarshalJSON(data []byte) error {
	type aiAssistantPersonalitys AiAssistantPersonalityCollection
	var v aiAssistantPersonalitys
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = AiAssistantPersonalityCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (a *AiAssistantPersonalityCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
