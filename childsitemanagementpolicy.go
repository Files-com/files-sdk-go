package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// ChildSiteManagementPolicy is a Files.com API resource.
type ChildSiteManagementPolicy struct {
	// Policy ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Type of policy. Valid values: `settings`.
	PolicyType string `json:"policy_type,omitempty" path:"policy_type,omitempty" url:"policy_type,omitempty"`
	// Name for this policy.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Description for this policy.
	Description string `json:"description,omitempty" path:"description,omitempty" url:"description,omitempty"`
	// Policy configuration data. On a Child Site, contains only managed setting names with null values. Settings policies accept site settings plus an optional `folder_behaviors` array for parent-managed root behaviors on child sites. For more information, refer to the Value Hash section of the developer documentation.
	Value interface{} `json:"value,omitempty" path:"value,omitempty" url:"value,omitempty"`
	// IDs of child sites that this policy has been applied to. This field is read-only.
	AppliedChildSiteIds []int64 `json:"applied_child_site_ids,omitempty" path:"applied_child_site_ids,omitempty" url:"applied_child_site_ids,omitempty"`
	// IDs of child sites excluded from this default policy.
	SkipChildSiteIds []int64 `json:"skip_child_site_ids,omitempty" path:"skip_child_site_ids,omitempty" url:"skip_child_site_ids,omitempty"`
	// IDs of child sites explicitly assigned to this non-default policy.
	ChildSiteIds []int64 `json:"child_site_ids,omitempty" path:"child_site_ids,omitempty" url:"child_site_ids,omitempty"`
	// Whether this policy applies to child sites not explicitly assigned to another policy.
	DefaultPolicy *bool `json:"default_policy,omitempty" path:"default_policy,omitempty" url:"default_policy,omitempty"`
	// When this policy was created.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// When this policy was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
}

// Identifier returns the resource ID.
func (c ChildSiteManagementPolicy) Identifier() interface{} {
	return c.Id
}

// ChildSiteManagementPolicyCollection is a list of ChildSiteManagementPolicy resources.
type ChildSiteManagementPolicyCollection []ChildSiteManagementPolicy

// ChildSiteManagementPolicyPolicyTypeEnum is a string value for policy_type.
// Enum lists the values documented by the API.
type ChildSiteManagementPolicyPolicyTypeEnum string

// String returns the API parameter value.
func (u ChildSiteManagementPolicyPolicyTypeEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u ChildSiteManagementPolicyPolicyTypeEnum) Enum() map[string]ChildSiteManagementPolicyPolicyTypeEnum {
	return map[string]ChildSiteManagementPolicyPolicyTypeEnum{
		"settings": ChildSiteManagementPolicyPolicyTypeEnum("settings"),
	}
}

// ChildSiteManagementPolicyListParams contains the request parameters for GET /child_site_management_policies.
type ChildSiteManagementPolicyListParams struct {
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// ChildSiteManagementPolicyFindParams contains the request parameters for GET /child_site_management_policies/{id}.
type ChildSiteManagementPolicyFindParams struct {
	// Child Site Management Policy ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// ChildSiteManagementPolicyCreateParams contains the request parameters for POST /child_site_management_policies.
type ChildSiteManagementPolicyCreateParams struct {
	// Policy configuration data. Attributes differ by policy type. For more information, refer to the Value Hash section of the developer documentation.
	Value interface{} `url:"value,omitempty" json:"value,omitempty" path:"value"`
	// IDs of child sites excluded from this default policy.
	SkipChildSiteIds []int64 `url:"skip_child_site_ids,omitempty" json:"skip_child_site_ids,omitempty" path:"skip_child_site_ids"`
	// IDs of child sites explicitly assigned to this non-default policy.
	ChildSiteIds []int64 `url:"child_site_ids,omitempty" json:"child_site_ids,omitempty" path:"child_site_ids"`
	// Whether this policy applies to child sites not explicitly assigned to another policy.
	DefaultPolicy *bool `url:"default_policy,omitempty" json:"default_policy,omitempty" path:"default_policy"`
	// Type of policy. Valid values: `settings`. Values: settings. Required.
	PolicyType ChildSiteManagementPolicyPolicyTypeEnum `url:"policy_type" json:"policy_type" path:"policy_type"`
	// Name for this policy.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Description for this policy.
	Description string `url:"description,omitempty" json:"description,omitempty" path:"description"`
}

// ChildSiteManagementPolicyUpdateParams contains the request parameters for PATCH /child_site_management_policies/{id}.
type ChildSiteManagementPolicyUpdateParams struct {
	// Child Site Management Policy ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Policy configuration data. Attributes differ by policy type. For more information, refer to the Value Hash section of the developer documentation.
	Value interface{} `url:"value,omitempty" json:"value,omitempty" path:"value"`
	// IDs of child sites excluded from this default policy.
	SkipChildSiteIds []int64 `url:"skip_child_site_ids,omitempty" json:"skip_child_site_ids,omitempty" path:"skip_child_site_ids"`
	// IDs of child sites explicitly assigned to this non-default policy.
	ChildSiteIds []int64 `url:"child_site_ids,omitempty" json:"child_site_ids,omitempty" path:"child_site_ids"`
	// Whether this policy applies to child sites not explicitly assigned to another policy.
	DefaultPolicy *bool `url:"default_policy,omitempty" json:"default_policy,omitempty" path:"default_policy"`
	// Type of policy. Valid values: `settings`. Values: settings.
	PolicyType ChildSiteManagementPolicyPolicyTypeEnum `url:"policy_type,omitempty" json:"policy_type,omitempty" path:"policy_type"`
	// Name for this policy.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Description for this policy.
	Description string `url:"description,omitempty" json:"description,omitempty" path:"description"`
}

// ChildSiteManagementPolicyDeleteParams contains the request parameters for DELETE /child_site_management_policies/{id}.
type ChildSiteManagementPolicyDeleteParams struct {
	// Child Site Management Policy ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (c *ChildSiteManagementPolicy) UnmarshalJSON(data []byte) error {
	type childSiteManagementPolicy ChildSiteManagementPolicy
	var v childSiteManagementPolicy
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*c = ChildSiteManagementPolicy(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (c *ChildSiteManagementPolicyCollection) UnmarshalJSON(data []byte) error {
	type childSiteManagementPolicys ChildSiteManagementPolicyCollection
	var v childSiteManagementPolicys
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*c = ChildSiteManagementPolicyCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (c *ChildSiteManagementPolicyCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*c))
	for i, v := range *c {
		ret[i] = v
	}

	return &ret
}
