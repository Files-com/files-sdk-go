package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// AiAssistantPersonality is a Files.com API resource.
type AiAssistantPersonality struct {
	Id                   int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	WorkspaceId          int64      `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	Name                 string     `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	SystemPrompt         string     `json:"system_prompt,omitempty" path:"system_prompt,omitempty" url:"system_prompt,omitempty"`
	UseByDefault         *bool      `json:"use_by_default,omitempty" path:"use_by_default,omitempty" url:"use_by_default,omitempty"`
	ApplyToAllWorkspaces *bool      `json:"apply_to_all_workspaces,omitempty" path:"apply_to_all_workspaces,omitempty" url:"apply_to_all_workspaces,omitempty"`
	CreatedAt            *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	UpdatedAt            *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
}

// Identifier returns the resource ID.
func (a AiAssistantPersonality) Identifier() interface{} {
	return a.Id
}

// AiAssistantPersonalityCollection is a list of AiAssistantPersonality resources.
type AiAssistantPersonalityCollection []AiAssistantPersonality

// AiAssistantPersonalityListParams contains the request parameters for GET /ai_assistant_personalities.
type AiAssistantPersonalityListParams struct {
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	ListParams
}

// AiAssistantPersonalityFindParams contains the request parameters for GET /ai_assistant_personalities/{id}.
type AiAssistantPersonalityFindParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// AiAssistantPersonalityCreateParams contains the request parameters for POST /ai_assistant_personalities.
type AiAssistantPersonalityCreateParams struct {
	ApplyToAllWorkspaces *bool  `url:"apply_to_all_workspaces,omitempty" json:"apply_to_all_workspaces,omitempty" path:"apply_to_all_workspaces"`
	Name                 string `url:"name" json:"name" path:"name"`
	SystemPrompt         string `url:"system_prompt" json:"system_prompt" path:"system_prompt"`
	UseByDefault         *bool  `url:"use_by_default,omitempty" json:"use_by_default,omitempty" path:"use_by_default"`
	WorkspaceId          int64  `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// AiAssistantPersonalityUpdateParams contains the request parameters for PATCH /ai_assistant_personalities/{id}.
type AiAssistantPersonalityUpdateParams struct {
	Id                   int64  `url:"-,omitempty" json:"-,omitempty" path:"id"`
	ApplyToAllWorkspaces *bool  `url:"apply_to_all_workspaces,omitempty" json:"apply_to_all_workspaces,omitempty" path:"apply_to_all_workspaces"`
	Name                 string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	SystemPrompt         string `url:"system_prompt,omitempty" json:"system_prompt,omitempty" path:"system_prompt"`
	UseByDefault         *bool  `url:"use_by_default,omitempty" json:"use_by_default,omitempty" path:"use_by_default"`
	WorkspaceId          int64  `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// AiAssistantPersonalityDeleteParams contains the request parameters for DELETE /ai_assistant_personalities/{id}.
type AiAssistantPersonalityDeleteParams struct {
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
