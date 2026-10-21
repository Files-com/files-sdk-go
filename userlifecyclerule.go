package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// UserLifecycleRule is a Files.com API resource.
type UserLifecycleRule struct {
	// User Lifecycle Rule ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// User authentication method for which the rule will apply. Use `all_non_sso` to target every non-SSO authentication method with one rule.
	AuthenticationMethod string `json:"authentication_method,omitempty" path:"authentication_method,omitempty" url:"authentication_method,omitempty"`
	// Array of Group IDs to which the rule applies. If empty or not set, the rule applies to all users.
	GroupIds []int64 `json:"group_ids,omitempty" path:"group_ids,omitempty" url:"group_ids,omitempty"`
	// Action to take on inactive users (disable or delete).
	Action string `json:"action,omitempty" path:"action,omitempty" url:"action,omitempty"`
	// Number of days of inactivity before the rule applies.
	InactivityDays int64 `json:"inactivity_days,omitempty" path:"inactivity_days,omitempty" url:"inactivity_days,omitempty"`
	// If true, the rule will apply to folder admins.
	IncludeFolderAdmins *bool `json:"include_folder_admins,omitempty" path:"include_folder_admins,omitempty" url:"include_folder_admins,omitempty"`
	// If true, the rule includes Site Administrators, who always belong to the Default Workspace. Can only be enabled when `workspace_id` is `0`.
	IncludeSiteAdmins *bool `json:"include_site_admins,omitempty" path:"include_site_admins,omitempty" url:"include_site_admins,omitempty"`
	// If true, a Default Workspace rule also applies to users in all Custom Workspaces. Can only be enabled when `workspace_id` is `0`.
	ApplyToAllWorkspaces *bool `json:"apply_to_all_workspaces,omitempty" path:"apply_to_all_workspaces,omitempty" url:"apply_to_all_workspaces,omitempty"`
	// User Lifecycle Rule name.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// If true, users will be emailed before the rule disables or deletes them.
	NotifyUsers *bool `json:"notify_users,omitempty" path:"notify_users,omitempty" url:"notify_users,omitempty"`
	// If provided, only users belonging to Partners with this tag at the Partner level will be affected by the rule. Tags must only contain lowercase letters, numbers, and hyphens.
	PartnerTag string `json:"partner_tag,omitempty" path:"partner_tag,omitempty" url:"partner_tag,omitempty"`
	// Site ID.
	SiteId int64 `json:"site_id,omitempty" path:"site_id,omitempty" url:"site_id,omitempty"`
	// Workspace whose users the rule applies to. `0` means the Default Workspace. A Custom Workspace rule applies only to users who belong to that Workspace, regardless of access granted to other users.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// State of the users to apply the rule to (inactive or disabled).
	UserState string `json:"user_state,omitempty" path:"user_state,omitempty" url:"user_state,omitempty"`
	// If provided, only users with this tag will be affected by the rule. Tags must only contain lowercase letters, numbers, and hyphens.
	UserTag string `json:"user_tag,omitempty" path:"user_tag,omitempty" url:"user_tag,omitempty"`
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
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `site_id` and `workspace_id`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// UserLifecycleRuleFindParams contains the request parameters for GET /user_lifecycle_rules/{id}.
type UserLifecycleRuleFindParams struct {
	// User Lifecycle Rule ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UserLifecycleRuleCreateParams contains the request parameters for POST /user_lifecycle_rules.
type UserLifecycleRuleCreateParams struct {
	// Action to take on inactive users (disable or delete). Values: disable, delete.
	Action UserLifecycleRuleActionEnum `url:"action,omitempty" json:"action,omitempty" path:"action"`
	// If true, a Default Workspace rule also applies to users in all Custom Workspaces. Can only be enabled when `workspace_id` is `0`.
	ApplyToAllWorkspaces *bool `url:"apply_to_all_workspaces,omitempty" json:"apply_to_all_workspaces,omitempty" path:"apply_to_all_workspaces"`
	// User authentication method for which the rule will apply. Use `all_non_sso` to target every non-SSO authentication method with one rule. Values: all, password, sso, none, email_signup, password_with_imported_hash, password_and_ssh_key, all_non_sso.
	AuthenticationMethod UserLifecycleRuleAuthenticationMethodEnum `url:"authentication_method,omitempty" json:"authentication_method,omitempty" path:"authentication_method"`
	// Array of Group IDs to which the rule applies. If empty or not set, the rule applies to all users.
	GroupIds []int64 `url:"group_ids,omitempty" json:"group_ids,omitempty" path:"group_ids"`
	// Number of days of inactivity before the rule applies.
	InactivityDays int64 `url:"inactivity_days,omitempty" json:"inactivity_days,omitempty" path:"inactivity_days"`
	// If true, the rule includes Site Administrators, who always belong to the Default Workspace. Can only be enabled when `workspace_id` is `0`.
	IncludeSiteAdmins *bool `url:"include_site_admins,omitempty" json:"include_site_admins,omitempty" path:"include_site_admins"`
	// If true, the rule will apply to folder admins.
	IncludeFolderAdmins *bool `url:"include_folder_admins,omitempty" json:"include_folder_admins,omitempty" path:"include_folder_admins"`
	// User Lifecycle Rule name.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// If true, users will be emailed before the rule disables or deletes them.
	NotifyUsers *bool `url:"notify_users,omitempty" json:"notify_users,omitempty" path:"notify_users"`
	// If provided, only users belonging to Partners with this tag at the Partner level will be affected by the rule. Tags must only contain lowercase letters, numbers, and hyphens.
	PartnerTag string `url:"partner_tag,omitempty" json:"partner_tag,omitempty" path:"partner_tag"`
	// State of the users to apply the rule to (inactive or disabled). Values: inactive, disabled.
	UserState UserLifecycleRuleUserStateEnum `url:"user_state,omitempty" json:"user_state,omitempty" path:"user_state"`
	// If provided, only users with this tag will be affected by the rule. Tags must only contain lowercase letters, numbers, and hyphens.
	UserTag string `url:"user_tag,omitempty" json:"user_tag,omitempty" path:"user_tag"`
	// Workspace whose users the rule applies to. `0` means the Default Workspace. A Custom Workspace rule applies only to users who belong to that Workspace, regardless of access granted to other users.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// UserLifecycleRuleUpdateParams contains the request parameters for PATCH /user_lifecycle_rules/{id}.
type UserLifecycleRuleUpdateParams struct {
	// User Lifecycle Rule ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Action to take on inactive users (disable or delete). Values: disable, delete.
	Action UserLifecycleRuleActionEnum `url:"action,omitempty" json:"action,omitempty" path:"action"`
	// If true, a Default Workspace rule also applies to users in all Custom Workspaces. Can only be enabled when `workspace_id` is `0`.
	ApplyToAllWorkspaces *bool `url:"apply_to_all_workspaces,omitempty" json:"apply_to_all_workspaces,omitempty" path:"apply_to_all_workspaces"`
	// User authentication method for which the rule will apply. Use `all_non_sso` to target every non-SSO authentication method with one rule. Values: all, password, sso, none, email_signup, password_with_imported_hash, password_and_ssh_key, all_non_sso.
	AuthenticationMethod UserLifecycleRuleAuthenticationMethodEnum `url:"authentication_method,omitempty" json:"authentication_method,omitempty" path:"authentication_method"`
	// Array of Group IDs to which the rule applies. If empty or not set, the rule applies to all users.
	GroupIds []int64 `url:"group_ids,omitempty" json:"group_ids,omitempty" path:"group_ids"`
	// Number of days of inactivity before the rule applies.
	InactivityDays int64 `url:"inactivity_days,omitempty" json:"inactivity_days,omitempty" path:"inactivity_days"`
	// If true, the rule includes Site Administrators, who always belong to the Default Workspace. Can only be enabled when `workspace_id` is `0`.
	IncludeSiteAdmins *bool `url:"include_site_admins,omitempty" json:"include_site_admins,omitempty" path:"include_site_admins"`
	// If true, the rule will apply to folder admins.
	IncludeFolderAdmins *bool `url:"include_folder_admins,omitempty" json:"include_folder_admins,omitempty" path:"include_folder_admins"`
	// User Lifecycle Rule name.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// If true, users will be emailed before the rule disables or deletes them.
	NotifyUsers *bool `url:"notify_users,omitempty" json:"notify_users,omitempty" path:"notify_users"`
	// If provided, only users belonging to Partners with this tag at the Partner level will be affected by the rule. Tags must only contain lowercase letters, numbers, and hyphens.
	PartnerTag string `url:"partner_tag,omitempty" json:"partner_tag,omitempty" path:"partner_tag"`
	// State of the users to apply the rule to (inactive or disabled). Values: inactive, disabled.
	UserState UserLifecycleRuleUserStateEnum `url:"user_state,omitempty" json:"user_state,omitempty" path:"user_state"`
	// If provided, only users with this tag will be affected by the rule. Tags must only contain lowercase letters, numbers, and hyphens.
	UserTag string `url:"user_tag,omitempty" json:"user_tag,omitempty" path:"user_tag"`
	// Workspace whose users the rule applies to. `0` means the Default Workspace. A Custom Workspace rule applies only to users who belong to that Workspace, regardless of access granted to other users.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// UserLifecycleRuleDeleteParams contains the request parameters for DELETE /user_lifecycle_rules/{id}.
type UserLifecycleRuleDeleteParams struct {
	// User Lifecycle Rule ID. Required; sent in the request path.
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
