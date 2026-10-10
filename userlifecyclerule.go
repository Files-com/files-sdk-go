package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// UserLifecycleRule is a Files.com API resource.
type UserLifecycleRule struct {
	Id                   int64   `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	AuthenticationMethod string  `json:"authentication_method,omitempty" path:"authentication_method,omitempty" url:"authentication_method,omitempty"`
	GroupIds             []int64 `json:"group_ids,omitempty" path:"group_ids,omitempty" url:"group_ids,omitempty"`
	Action               string  `json:"action,omitempty" path:"action,omitempty" url:"action,omitempty"`
	InactivityDays       int64   `json:"inactivity_days,omitempty" path:"inactivity_days,omitempty" url:"inactivity_days,omitempty"`
	IncludeFolderAdmins  *bool   `json:"include_folder_admins,omitempty" path:"include_folder_admins,omitempty" url:"include_folder_admins,omitempty"`
	IncludeSiteAdmins    *bool   `json:"include_site_admins,omitempty" path:"include_site_admins,omitempty" url:"include_site_admins,omitempty"`
	ApplyToAllWorkspaces *bool   `json:"apply_to_all_workspaces,omitempty" path:"apply_to_all_workspaces,omitempty" url:"apply_to_all_workspaces,omitempty"`
	Name                 string  `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	NotifyUsers          *bool   `json:"notify_users,omitempty" path:"notify_users,omitempty" url:"notify_users,omitempty"`
	PartnerTag           string  `json:"partner_tag,omitempty" path:"partner_tag,omitempty" url:"partner_tag,omitempty"`
	SiteId               int64   `json:"site_id,omitempty" path:"site_id,omitempty" url:"site_id,omitempty"`
	WorkspaceId          int64   `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	UserState            string  `json:"user_state,omitempty" path:"user_state,omitempty" url:"user_state,omitempty"`
	UserTag              string  `json:"user_tag,omitempty" path:"user_tag,omitempty" url:"user_tag,omitempty"`
}

// Identifier returns the resource ID.
func (u UserLifecycleRule) Identifier() interface{} {
	return u.Id
}

// UserLifecycleRuleCollection is a list of UserLifecycleRule resources.
type UserLifecycleRuleCollection []UserLifecycleRule

// UserLifecycleRuleActionEnum is a string value for action.
// Enum lists the values documented by the API.
type UserLifecycleRuleActionEnum string

// String returns the API parameter value.
func (u UserLifecycleRuleActionEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u UserLifecycleRuleActionEnum) Enum() map[string]UserLifecycleRuleActionEnum {
	return map[string]UserLifecycleRuleActionEnum{
		"disable": UserLifecycleRuleActionEnum("disable"),
		"delete":  UserLifecycleRuleActionEnum("delete"),
	}
}

// UserLifecycleRuleAuthenticationMethodEnum is a string value for authentication_method.
// Enum lists the values documented by the API.
type UserLifecycleRuleAuthenticationMethodEnum string

// String returns the API parameter value.
func (u UserLifecycleRuleAuthenticationMethodEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u UserLifecycleRuleAuthenticationMethodEnum) Enum() map[string]UserLifecycleRuleAuthenticationMethodEnum {
	return map[string]UserLifecycleRuleAuthenticationMethodEnum{
		"all":                         UserLifecycleRuleAuthenticationMethodEnum("all"),
		"password":                    UserLifecycleRuleAuthenticationMethodEnum("password"),
		"sso":                         UserLifecycleRuleAuthenticationMethodEnum("sso"),
		"none":                        UserLifecycleRuleAuthenticationMethodEnum("none"),
		"email_signup":                UserLifecycleRuleAuthenticationMethodEnum("email_signup"),
		"password_with_imported_hash": UserLifecycleRuleAuthenticationMethodEnum("password_with_imported_hash"),
		"password_and_ssh_key":        UserLifecycleRuleAuthenticationMethodEnum("password_and_ssh_key"),
		"all_non_sso":                 UserLifecycleRuleAuthenticationMethodEnum("all_non_sso"),
	}
}

// UserLifecycleRuleUserStateEnum is a string value for user_state.
// Enum lists the values documented by the API.
type UserLifecycleRuleUserStateEnum string

// String returns the API parameter value.
func (u UserLifecycleRuleUserStateEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u UserLifecycleRuleUserStateEnum) Enum() map[string]UserLifecycleRuleUserStateEnum {
	return map[string]UserLifecycleRuleUserStateEnum{
		"inactive": UserLifecycleRuleUserStateEnum("inactive"),
		"disabled": UserLifecycleRuleUserStateEnum("disabled"),
	}
}

// UserLifecycleRuleListParams contains the request parameters for GET /user_lifecycle_rules.
type UserLifecycleRuleListParams struct {
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	ListParams
}

// UserLifecycleRuleFindParams contains the request parameters for GET /user_lifecycle_rules/{id}.
type UserLifecycleRuleFindParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UserLifecycleRuleCreateParams contains the request parameters for POST /user_lifecycle_rules.
type UserLifecycleRuleCreateParams struct {
	Action               UserLifecycleRuleActionEnum               `url:"action,omitempty" json:"action,omitempty" path:"action"`
	ApplyToAllWorkspaces *bool                                     `url:"apply_to_all_workspaces,omitempty" json:"apply_to_all_workspaces,omitempty" path:"apply_to_all_workspaces"`
	AuthenticationMethod UserLifecycleRuleAuthenticationMethodEnum `url:"authentication_method,omitempty" json:"authentication_method,omitempty" path:"authentication_method"`
	GroupIds             []int64                                   `url:"group_ids,omitempty" json:"group_ids,omitempty" path:"group_ids"`
	InactivityDays       int64                                     `url:"inactivity_days,omitempty" json:"inactivity_days,omitempty" path:"inactivity_days"`
	IncludeSiteAdmins    *bool                                     `url:"include_site_admins,omitempty" json:"include_site_admins,omitempty" path:"include_site_admins"`
	IncludeFolderAdmins  *bool                                     `url:"include_folder_admins,omitempty" json:"include_folder_admins,omitempty" path:"include_folder_admins"`
	Name                 string                                    `url:"name,omitempty" json:"name,omitempty" path:"name"`
	NotifyUsers          *bool                                     `url:"notify_users,omitempty" json:"notify_users,omitempty" path:"notify_users"`
	PartnerTag           string                                    `url:"partner_tag,omitempty" json:"partner_tag,omitempty" path:"partner_tag"`
	UserState            UserLifecycleRuleUserStateEnum            `url:"user_state,omitempty" json:"user_state,omitempty" path:"user_state"`
	UserTag              string                                    `url:"user_tag,omitempty" json:"user_tag,omitempty" path:"user_tag"`
	WorkspaceId          int64                                     `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// UserLifecycleRuleUpdateParams contains the request parameters for PATCH /user_lifecycle_rules/{id}.
type UserLifecycleRuleUpdateParams struct {
	Id                   int64                                     `url:"-,omitempty" json:"-,omitempty" path:"id"`
	Action               UserLifecycleRuleActionEnum               `url:"action,omitempty" json:"action,omitempty" path:"action"`
	ApplyToAllWorkspaces *bool                                     `url:"apply_to_all_workspaces,omitempty" json:"apply_to_all_workspaces,omitempty" path:"apply_to_all_workspaces"`
	AuthenticationMethod UserLifecycleRuleAuthenticationMethodEnum `url:"authentication_method,omitempty" json:"authentication_method,omitempty" path:"authentication_method"`
	GroupIds             []int64                                   `url:"group_ids,omitempty" json:"group_ids,omitempty" path:"group_ids"`
	InactivityDays       int64                                     `url:"inactivity_days,omitempty" json:"inactivity_days,omitempty" path:"inactivity_days"`
	IncludeSiteAdmins    *bool                                     `url:"include_site_admins,omitempty" json:"include_site_admins,omitempty" path:"include_site_admins"`
	IncludeFolderAdmins  *bool                                     `url:"include_folder_admins,omitempty" json:"include_folder_admins,omitempty" path:"include_folder_admins"`
	Name                 string                                    `url:"name,omitempty" json:"name,omitempty" path:"name"`
	NotifyUsers          *bool                                     `url:"notify_users,omitempty" json:"notify_users,omitempty" path:"notify_users"`
	PartnerTag           string                                    `url:"partner_tag,omitempty" json:"partner_tag,omitempty" path:"partner_tag"`
	UserState            UserLifecycleRuleUserStateEnum            `url:"user_state,omitempty" json:"user_state,omitempty" path:"user_state"`
	UserTag              string                                    `url:"user_tag,omitempty" json:"user_tag,omitempty" path:"user_tag"`
	WorkspaceId          int64                                     `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// UserLifecycleRuleDeleteParams contains the request parameters for DELETE /user_lifecycle_rules/{id}.
type UserLifecycleRuleDeleteParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (u *UserLifecycleRule) UnmarshalJSON(data []byte) error {
	type userLifecycleRule UserLifecycleRule
	var v userLifecycleRule
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*u = UserLifecycleRule(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (u *UserLifecycleRuleCollection) UnmarshalJSON(data []byte) error {
	type userLifecycleRules UserLifecycleRuleCollection
	var v userLifecycleRules
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*u = UserLifecycleRuleCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (u *UserLifecycleRuleCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*u))
	for i, v := range *u {
		ret[i] = v
	}

	return &ret
}
