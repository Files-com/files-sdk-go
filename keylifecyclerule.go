package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// KeyLifecycleRule is a Files.com API resource.
type KeyLifecycleRule struct {
	// Key Lifecycle Rule ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Key type for which the rule will apply (gpg, ssh, or api).
	KeyType string `json:"key_type,omitempty" path:"key_type,omitempty" url:"key_type,omitempty"`
	// Number of days of inactivity before the rule applies.
	InactivityDays int64 `json:"inactivity_days,omitempty" path:"inactivity_days,omitempty" url:"inactivity_days,omitempty"`
	// Number of days after creation before an SSH key expires. Applies only to SSH keys.
	ExpirationDays int64 `json:"expiration_days,omitempty" path:"expiration_days,omitempty" url:"expiration_days,omitempty"`
	// If true, a default-workspace rule also applies to keys in all workspaces.
	ApplyToAllWorkspaces *bool `json:"apply_to_all_workspaces,omitempty" path:"apply_to_all_workspaces,omitempty" url:"apply_to_all_workspaces,omitempty"`
	// Key Lifecycle Rule name.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Workspace ID. `0` means the default workspace.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
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
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `key_type`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `key_type` and `workspace_id`. Valid field combinations are `[ workspace_id, key_type ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// KeyLifecycleRuleFindParams contains the request parameters for GET /key_lifecycle_rules/{id}.
type KeyLifecycleRuleFindParams struct {
	// Key Lifecycle Rule ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// KeyLifecycleRuleCreateParams contains the request parameters for POST /key_lifecycle_rules.
type KeyLifecycleRuleCreateParams struct {
	// If true, a default-workspace rule also applies to keys in all workspaces.
	ApplyToAllWorkspaces *bool `url:"apply_to_all_workspaces,omitempty" json:"apply_to_all_workspaces,omitempty" path:"apply_to_all_workspaces"`
	// Number of days after creation before an SSH key expires. Applies only to SSH keys.
	ExpirationDays int64 `url:"expiration_days,omitempty" json:"expiration_days,omitempty" path:"expiration_days"`
	// Key type for which the rule will apply (gpg, ssh, or api). Values: gpg, ssh, api.
	KeyType KeyLifecycleRuleKeyTypeEnum `url:"key_type,omitempty" json:"key_type,omitempty" path:"key_type"`
	// Number of days of inactivity before the rule applies.
	InactivityDays int64 `url:"inactivity_days,omitempty" json:"inactivity_days,omitempty" path:"inactivity_days"`
	// Key Lifecycle Rule name.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Workspace ID. `0` means the default workspace.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// KeyLifecycleRuleUpdateParams contains the request parameters for PATCH /key_lifecycle_rules/{id}.
type KeyLifecycleRuleUpdateParams struct {
	// Key Lifecycle Rule ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// If true, a default-workspace rule also applies to keys in all workspaces.
	ApplyToAllWorkspaces *bool `url:"apply_to_all_workspaces,omitempty" json:"apply_to_all_workspaces,omitempty" path:"apply_to_all_workspaces"`
	// Number of days after creation before an SSH key expires. Applies only to SSH keys.
	ExpirationDays int64 `url:"expiration_days,omitempty" json:"expiration_days,omitempty" path:"expiration_days"`
	// Key type for which the rule will apply (gpg, ssh, or api). Values: gpg, ssh, api.
	KeyType KeyLifecycleRuleKeyTypeEnum `url:"key_type,omitempty" json:"key_type,omitempty" path:"key_type"`
	// Number of days of inactivity before the rule applies.
	InactivityDays int64 `url:"inactivity_days,omitempty" json:"inactivity_days,omitempty" path:"inactivity_days"`
	// Key Lifecycle Rule name.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Workspace ID. `0` means the default workspace.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// KeyLifecycleRuleDeleteParams contains the request parameters for DELETE /key_lifecycle_rules/{id}.
type KeyLifecycleRuleDeleteParams struct {
	// Key Lifecycle Rule ID. Required; sent in the request path.
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
