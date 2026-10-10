package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// KeyLifecycleRule is a Files.com API resource.
type KeyLifecycleRule struct {
	Id                   int64  `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	KeyType              string `json:"key_type,omitempty" path:"key_type,omitempty" url:"key_type,omitempty"`
	InactivityDays       int64  `json:"inactivity_days,omitempty" path:"inactivity_days,omitempty" url:"inactivity_days,omitempty"`
	ExpirationDays       int64  `json:"expiration_days,omitempty" path:"expiration_days,omitempty" url:"expiration_days,omitempty"`
	ApplyToAllWorkspaces *bool  `json:"apply_to_all_workspaces,omitempty" path:"apply_to_all_workspaces,omitempty" url:"apply_to_all_workspaces,omitempty"`
	Name                 string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	WorkspaceId          int64  `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
}

// Identifier returns the resource ID.
func (k KeyLifecycleRule) Identifier() interface{} {
	return k.Id
}

// KeyLifecycleRuleCollection is a list of KeyLifecycleRule resources.
type KeyLifecycleRuleCollection []KeyLifecycleRule

// KeyLifecycleRuleKeyTypeEnum is a string value for key_type.
// Enum lists the values documented by the API.
type KeyLifecycleRuleKeyTypeEnum string

// String returns the API parameter value.
func (u KeyLifecycleRuleKeyTypeEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u KeyLifecycleRuleKeyTypeEnum) Enum() map[string]KeyLifecycleRuleKeyTypeEnum {
	return map[string]KeyLifecycleRuleKeyTypeEnum{
		"gpg": KeyLifecycleRuleKeyTypeEnum("gpg"),
		"ssh": KeyLifecycleRuleKeyTypeEnum("ssh"),
		"api": KeyLifecycleRuleKeyTypeEnum("api"),
	}
}

// KeyLifecycleRuleListParams contains the request parameters for GET /key_lifecycle_rules.
type KeyLifecycleRuleListParams struct {
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	ListParams
}

// KeyLifecycleRuleFindParams contains the request parameters for GET /key_lifecycle_rules/{id}.
type KeyLifecycleRuleFindParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// KeyLifecycleRuleCreateParams contains the request parameters for POST /key_lifecycle_rules.
type KeyLifecycleRuleCreateParams struct {
	ApplyToAllWorkspaces *bool                       `url:"apply_to_all_workspaces,omitempty" json:"apply_to_all_workspaces,omitempty" path:"apply_to_all_workspaces"`
	ExpirationDays       int64                       `url:"expiration_days,omitempty" json:"expiration_days,omitempty" path:"expiration_days"`
	KeyType              KeyLifecycleRuleKeyTypeEnum `url:"key_type,omitempty" json:"key_type,omitempty" path:"key_type"`
	InactivityDays       int64                       `url:"inactivity_days,omitempty" json:"inactivity_days,omitempty" path:"inactivity_days"`
	Name                 string                      `url:"name,omitempty" json:"name,omitempty" path:"name"`
	WorkspaceId          int64                       `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// KeyLifecycleRuleUpdateParams contains the request parameters for PATCH /key_lifecycle_rules/{id}.
type KeyLifecycleRuleUpdateParams struct {
	Id                   int64                       `url:"-,omitempty" json:"-,omitempty" path:"id"`
	ApplyToAllWorkspaces *bool                       `url:"apply_to_all_workspaces,omitempty" json:"apply_to_all_workspaces,omitempty" path:"apply_to_all_workspaces"`
	ExpirationDays       int64                       `url:"expiration_days,omitempty" json:"expiration_days,omitempty" path:"expiration_days"`
	KeyType              KeyLifecycleRuleKeyTypeEnum `url:"key_type,omitempty" json:"key_type,omitempty" path:"key_type"`
	InactivityDays       int64                       `url:"inactivity_days,omitempty" json:"inactivity_days,omitempty" path:"inactivity_days"`
	Name                 string                      `url:"name,omitempty" json:"name,omitempty" path:"name"`
	WorkspaceId          int64                       `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// KeyLifecycleRuleDeleteParams contains the request parameters for DELETE /key_lifecycle_rules/{id}.
type KeyLifecycleRuleDeleteParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (k *KeyLifecycleRule) UnmarshalJSON(data []byte) error {
	type keyLifecycleRule KeyLifecycleRule
	var v keyLifecycleRule
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*k = KeyLifecycleRule(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (k *KeyLifecycleRuleCollection) UnmarshalJSON(data []byte) error {
	type keyLifecycleRules KeyLifecycleRuleCollection
	var v keyLifecycleRules
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*k = KeyLifecycleRuleCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (k *KeyLifecycleRuleCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*k))
	for i, v := range *k {
		ret[i] = v
	}

	return &ret
}
