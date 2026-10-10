package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Clickwrap is a Files.com API resource.
type Clickwrap struct {
	Id             int64  `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Name           string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Body           string `json:"body,omitempty" path:"body,omitempty" url:"body,omitempty"`
	UseWithUsers   string `json:"use_with_users,omitempty" path:"use_with_users,omitempty" url:"use_with_users,omitempty"`
	UseWithBundles string `json:"use_with_bundles,omitempty" path:"use_with_bundles,omitempty" url:"use_with_bundles,omitempty"`
	UseWithInboxes string `json:"use_with_inboxes,omitempty" path:"use_with_inboxes,omitempty" url:"use_with_inboxes,omitempty"`
}

// Identifier returns the resource ID.
func (c Clickwrap) Identifier() interface{} {
	return c.Id
}

// ClickwrapCollection is a list of Clickwrap resources.
type ClickwrapCollection []Clickwrap

// ClickwrapUseWithBundlesEnum is a string value for use_with_bundles.
// Enum lists the values documented by the API.
type ClickwrapUseWithBundlesEnum string

// String returns the API parameter value.
func (u ClickwrapUseWithBundlesEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u ClickwrapUseWithBundlesEnum) Enum() map[string]ClickwrapUseWithBundlesEnum {
	return map[string]ClickwrapUseWithBundlesEnum{
		"none":                   ClickwrapUseWithBundlesEnum("none"),
		"available":              ClickwrapUseWithBundlesEnum("available"),
		"require":                ClickwrapUseWithBundlesEnum("require"),
		"available_to_all_users": ClickwrapUseWithBundlesEnum("available_to_all_users"),
	}
}

// ClickwrapUseWithInboxesEnum is a string value for use_with_inboxes.
// Enum lists the values documented by the API.
type ClickwrapUseWithInboxesEnum string

// String returns the API parameter value.
func (u ClickwrapUseWithInboxesEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u ClickwrapUseWithInboxesEnum) Enum() map[string]ClickwrapUseWithInboxesEnum {
	return map[string]ClickwrapUseWithInboxesEnum{
		"none":                   ClickwrapUseWithInboxesEnum("none"),
		"available":              ClickwrapUseWithInboxesEnum("available"),
		"require":                ClickwrapUseWithInboxesEnum("require"),
		"available_to_all_users": ClickwrapUseWithInboxesEnum("available_to_all_users"),
	}
}

// ClickwrapUseWithUsersEnum is a string value for use_with_users.
// Enum lists the values documented by the API.
type ClickwrapUseWithUsersEnum string

// String returns the API parameter value.
func (u ClickwrapUseWithUsersEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u ClickwrapUseWithUsersEnum) Enum() map[string]ClickwrapUseWithUsersEnum {
	return map[string]ClickwrapUseWithUsersEnum{
		"none":                     ClickwrapUseWithUsersEnum("none"),
		"require":                  ClickwrapUseWithUsersEnum("require"),
		"require_all_users_once":   ClickwrapUseWithUsersEnum("require_all_users_once"),
		"require_all_users_always": ClickwrapUseWithUsersEnum("require_all_users_always"),
	}
}

// ClickwrapListParams contains the request parameters for GET /clickwraps.
type ClickwrapListParams struct {
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	ListParams
}

// ClickwrapFindParams contains the request parameters for GET /clickwraps/{id}.
type ClickwrapFindParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// ClickwrapCreateParams contains the request parameters for POST /clickwraps.
type ClickwrapCreateParams struct {
	Name           string                      `url:"name,omitempty" json:"name,omitempty" path:"name"`
	Body           string                      `url:"body,omitempty" json:"body,omitempty" path:"body"`
	UseWithBundles ClickwrapUseWithBundlesEnum `url:"use_with_bundles,omitempty" json:"use_with_bundles,omitempty" path:"use_with_bundles"`
	UseWithInboxes ClickwrapUseWithInboxesEnum `url:"use_with_inboxes,omitempty" json:"use_with_inboxes,omitempty" path:"use_with_inboxes"`
	UseWithUsers   ClickwrapUseWithUsersEnum   `url:"use_with_users,omitempty" json:"use_with_users,omitempty" path:"use_with_users"`
}

// ClickwrapUpdateParams contains the request parameters for PATCH /clickwraps/{id}.
type ClickwrapUpdateParams struct {
	Id             int64                       `url:"-,omitempty" json:"-,omitempty" path:"id"`
	Name           string                      `url:"name,omitempty" json:"name,omitempty" path:"name"`
	Body           string                      `url:"body,omitempty" json:"body,omitempty" path:"body"`
	UseWithBundles ClickwrapUseWithBundlesEnum `url:"use_with_bundles,omitempty" json:"use_with_bundles,omitempty" path:"use_with_bundles"`
	UseWithInboxes ClickwrapUseWithInboxesEnum `url:"use_with_inboxes,omitempty" json:"use_with_inboxes,omitempty" path:"use_with_inboxes"`
	UseWithUsers   ClickwrapUseWithUsersEnum   `url:"use_with_users,omitempty" json:"use_with_users,omitempty" path:"use_with_users"`
}

// ClickwrapDeleteParams contains the request parameters for DELETE /clickwraps/{id}.
type ClickwrapDeleteParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (c *Clickwrap) UnmarshalJSON(data []byte) error {
	type clickwrap Clickwrap
	var v clickwrap
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*c = Clickwrap(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (c *ClickwrapCollection) UnmarshalJSON(data []byte) error {
	type clickwraps ClickwrapCollection
	var v clickwraps
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*c = ClickwrapCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (c *ClickwrapCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*c))
	for i, v := range *c {
		ret[i] = v
	}

	return &ret
}
