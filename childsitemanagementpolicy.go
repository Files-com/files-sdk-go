package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// ChildSiteManagementPolicy is a Files.com API resource.
type ChildSiteManagementPolicy struct {
	Id                  int64       `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	PolicyType          string      `json:"policy_type,omitempty" path:"policy_type,omitempty" url:"policy_type,omitempty"`
	Name                string      `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Description         string      `json:"description,omitempty" path:"description,omitempty" url:"description,omitempty"`
	Value               interface{} `json:"value,omitempty" path:"value,omitempty" url:"value,omitempty"`
	AppliedChildSiteIds []int64     `json:"applied_child_site_ids,omitempty" path:"applied_child_site_ids,omitempty" url:"applied_child_site_ids,omitempty"`
	SkipChildSiteIds    []int64     `json:"skip_child_site_ids,omitempty" path:"skip_child_site_ids,omitempty" url:"skip_child_site_ids,omitempty"`
	ChildSiteIds        []int64     `json:"child_site_ids,omitempty" path:"child_site_ids,omitempty" url:"child_site_ids,omitempty"`
	DefaultPolicy       *bool       `json:"default_policy,omitempty" path:"default_policy,omitempty" url:"default_policy,omitempty"`
	CreatedAt           *time.Time  `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	UpdatedAt           *time.Time  `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
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
	ListParams
}

// ChildSiteManagementPolicyFindParams contains the request parameters for GET /child_site_management_policies/{id}.
type ChildSiteManagementPolicyFindParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// ChildSiteManagementPolicyCreateParams contains the request parameters for POST /child_site_management_policies.
type ChildSiteManagementPolicyCreateParams struct {
	Value            interface{}                             `url:"value,omitempty" json:"value,omitempty" path:"value"`
	SkipChildSiteIds []int64                                 `url:"skip_child_site_ids,omitempty" json:"skip_child_site_ids,omitempty" path:"skip_child_site_ids"`
	ChildSiteIds     []int64                                 `url:"child_site_ids,omitempty" json:"child_site_ids,omitempty" path:"child_site_ids"`
	DefaultPolicy    *bool                                   `url:"default_policy,omitempty" json:"default_policy,omitempty" path:"default_policy"`
	PolicyType       ChildSiteManagementPolicyPolicyTypeEnum `url:"policy_type" json:"policy_type" path:"policy_type"`
	Name             string                                  `url:"name,omitempty" json:"name,omitempty" path:"name"`
	Description      string                                  `url:"description,omitempty" json:"description,omitempty" path:"description"`
}

// ChildSiteManagementPolicyUpdateParams contains the request parameters for PATCH /child_site_management_policies/{id}.
type ChildSiteManagementPolicyUpdateParams struct {
	Id               int64                                   `url:"-,omitempty" json:"-,omitempty" path:"id"`
	Value            interface{}                             `url:"value,omitempty" json:"value,omitempty" path:"value"`
	SkipChildSiteIds []int64                                 `url:"skip_child_site_ids,omitempty" json:"skip_child_site_ids,omitempty" path:"skip_child_site_ids"`
	ChildSiteIds     []int64                                 `url:"child_site_ids,omitempty" json:"child_site_ids,omitempty" path:"child_site_ids"`
	DefaultPolicy    *bool                                   `url:"default_policy,omitempty" json:"default_policy,omitempty" path:"default_policy"`
	PolicyType       ChildSiteManagementPolicyPolicyTypeEnum `url:"policy_type,omitempty" json:"policy_type,omitempty" path:"policy_type"`
	Name             string                                  `url:"name,omitempty" json:"name,omitempty" path:"name"`
	Description      string                                  `url:"description,omitempty" json:"description,omitempty" path:"description"`
}

// ChildSiteManagementPolicyDeleteParams contains the request parameters for DELETE /child_site_management_policies/{id}.
type ChildSiteManagementPolicyDeleteParams struct {
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
