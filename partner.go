package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Partner is a Files.com API resource.
type Partner struct {
	// Allow Partner Admins to change Two-Factor Authentication requirements for Partner Users.
	AllowBypassing2faPolicies *bool `json:"allow_bypassing_2fa_policies,omitempty" path:"allow_bypassing_2fa_policies,omitempty" url:"allow_bypassing_2fa_policies,omitempty"`
	// A list of allowed IPs for this Partner. Newline delimited. Partner User IP access is allowed when the IP matches the Partner, User, or Site allowed IP lists.
	AllowedIps string `json:"allowed_ips,omitempty" path:"allowed_ips,omitempty" url:"allowed_ips,omitempty"`
	// Allow Partner Admins to change or reset credentials for users belonging to this Partner.
	AllowCredentialChanges *bool `json:"allow_credential_changes,omitempty" path:"allow_credential_changes,omitempty" url:"allow_credential_changes,omitempty"`
	// Allow Partner Admins to provide GPG keys.
	AllowProvidingGpgKeys *bool `json:"allow_providing_gpg_keys,omitempty" path:"allow_providing_gpg_keys,omitempty" url:"allow_providing_gpg_keys,omitempty"`
	// Allow Partner Admins to create users.
	AllowUserCreation *bool `json:"allow_user_creation,omitempty" path:"allow_user_creation,omitempty" url:"allow_user_creation,omitempty"`
	// When `true`, emails sent to Partner users are copied to the responsible User or Group.
	CcEmailsToResponsibleParty *bool `json:"cc_emails_to_responsible_party,omitempty" path:"cc_emails_to_responsible_party,omitempty" url:"cc_emails_to_responsible_party,omitempty"`
	// Approved Connected Sites relationships for this Partner, in both directions. Empty when this Partner has no connections. Read-only.
	Connections []PartnerConnection `json:"connections,omitempty" path:"connections,omitempty" url:"connections,omitempty"`
	// The unique ID of the Partner.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// AI Assistant Personality ID assigned to this Partner, if any. Users in the Partner inherit it unless a direct per-user assignment overrides it.
	AiAssistantPersonalityId int64 `json:"ai_assistant_personality_id,omitempty" path:"ai_assistant_personality_id,omitempty" url:"ai_assistant_personality_id,omitempty"`
	// ID of the Workspace associated with this Partner.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// The name of the Partner.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Notes about this Partner.
	Notes string `json:"notes,omitempty" path:"notes,omitempty" url:"notes,omitempty"`
	// Array of User IDs that are Partner Admins for this Partner.
	PartnerAdminIds []int64 `json:"partner_admin_ids,omitempty" path:"partner_admin_ids,omitempty" url:"partner_admin_ids,omitempty"`
	// ID of the Partner Channel Template assigned to this Partner.
	PartnerChannelTemplateId int64 `json:"partner_channel_template_id,omitempty" path:"partner_channel_template_id,omitempty" url:"partner_channel_template_id,omitempty"`
	// This site's role for this Partner in Connected Sites relationships. `host` is a Partner this site configured. `guest` is a Partner created by approving another site's connection request; it has no root folder and cannot hold users, permissions, or Partner Channels, or host a connection. `host_and_guest` is a configured Partner that is also the guest side of a connection. Promote a `guest` Partner by setting this to `host_and_guest` together with a `root_folder`.
	PartnershipRole string `json:"partnership_role,omitempty" path:"partnership_role,omitempty" url:"partnership_role,omitempty"`
	// ID of the Group responsible for this Partner.
	ResponsibleGroupId int64 `json:"responsible_group_id,omitempty" path:"responsible_group_id,omitempty" url:"responsible_group_id,omitempty"`
	// ID of the User responsible for this Partner.
	ResponsibleUserId int64 `json:"responsible_user_id,omitempty" path:"responsible_user_id,omitempty" url:"responsible_user_id,omitempty"`
	// The root folder path for this Partner.
	RootFolder string `json:"root_folder,omitempty" path:"root_folder,omitempty" url:"root_folder,omitempty"`
	// Show Partner users a simplified home page built from this Partner's Channels.
	ShowPartnerChannelHomePage *bool `json:"show_partner_channel_home_page,omitempty" path:"show_partner_channel_home_page,omitempty" url:"show_partner_channel_home_page,omitempty"`
	// Comma-separated list of Tags for this Partner. Tags are used for other features, such as UserLifecycleRules, which can target specific tags. Tags must only contain lowercase letters, numbers, and hyphens.
	Tags string `json:"tags,omitempty" path:"tags,omitempty" url:"tags,omitempty"`
	// Array of User IDs that belong to this Partner.
	UserIds []int64 `json:"user_ids,omitempty" path:"user_ids,omitempty" url:"user_ids,omitempty"`
}

// Identifier returns the resource ID.
func (p Partner) Identifier() interface{} {
	return p.Id
}

// PartnerCollection is a list of Partner resources.
type PartnerCollection []Partner

// PartnerPartnershipRoleEnum is a string value for partnership_role.
// Enum lists the values documented by the API.
type PartnerPartnershipRoleEnum string

// String returns the API parameter value.
func (u PartnerPartnershipRoleEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u PartnerPartnershipRoleEnum) Enum() map[string]PartnerPartnershipRoleEnum {
	return map[string]PartnerPartnershipRoleEnum{
		"host":           PartnerPartnershipRoleEnum("host"),
		"guest":          PartnerPartnershipRoleEnum("guest"),
		"host_and_guest": PartnerPartnershipRoleEnum("host_and_guest"),
	}
}

// PartnerListParams contains the request parameters for GET /partners.
type PartnerListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `name`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// PartnerFindParams contains the request parameters for GET /partners/{id}.
type PartnerFindParams struct {
	// Partner ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// PartnerCreateParams contains the request parameters for POST /partners.
type PartnerCreateParams struct {
	// AI Assistant Personality ID assigned to this Partner, if any. Users in the Partner inherit it unless a direct per-user assignment overrides it.
	AiAssistantPersonalityId int64 `url:"ai_assistant_personality_id,omitempty" json:"ai_assistant_personality_id,omitempty" path:"ai_assistant_personality_id"`
	// A list of allowed IPs for this Partner. Newline delimited. Partner User IP access is allowed when the IP matches the Partner, User, or Site allowed IP lists.
	AllowedIps string `url:"allowed_ips,omitempty" json:"allowed_ips,omitempty" path:"allowed_ips"`
	// Allow Partner Admins to change Two-Factor Authentication requirements for Partner Users.
	AllowBypassing2faPolicies *bool `url:"allow_bypassing_2fa_policies,omitempty" json:"allow_bypassing_2fa_policies,omitempty" path:"allow_bypassing_2fa_policies"`
	// Allow Partner Admins to change or reset credentials for users belonging to this Partner.
	AllowCredentialChanges *bool `url:"allow_credential_changes,omitempty" json:"allow_credential_changes,omitempty" path:"allow_credential_changes"`
	// Allow Partner Admins to provide GPG keys.
	AllowProvidingGpgKeys *bool `url:"allow_providing_gpg_keys,omitempty" json:"allow_providing_gpg_keys,omitempty" path:"allow_providing_gpg_keys"`
	// Allow Partner Admins to create users.
	AllowUserCreation *bool `url:"allow_user_creation,omitempty" json:"allow_user_creation,omitempty" path:"allow_user_creation"`
	// When `true`, emails sent to Partner users are copied to the responsible User or Group.
	CcEmailsToResponsibleParty *bool `url:"cc_emails_to_responsible_party,omitempty" json:"cc_emails_to_responsible_party,omitempty" path:"cc_emails_to_responsible_party"`
	// Notes about this Partner.
	Notes string `url:"notes,omitempty" json:"notes,omitempty" path:"notes"`
	// ID of the Partner Channel Template assigned to this Partner.
	PartnerChannelTemplateId int64 `url:"partner_channel_template_id,omitempty" json:"partner_channel_template_id,omitempty" path:"partner_channel_template_id"`
	// ID of the Group responsible for this Partner.
	ResponsibleGroupId int64 `url:"responsible_group_id,omitempty" json:"responsible_group_id,omitempty" path:"responsible_group_id"`
	// ID of the User responsible for this Partner.
	ResponsibleUserId int64 `url:"responsible_user_id,omitempty" json:"responsible_user_id,omitempty" path:"responsible_user_id"`
	// Show Partner users a simplified home page built from this Partner's Channels.
	ShowPartnerChannelHomePage *bool `url:"show_partner_channel_home_page,omitempty" json:"show_partner_channel_home_page,omitempty" path:"show_partner_channel_home_page"`
	// Comma-separated list of Tags for this Partner. Tags are used for other features, such as UserLifecycleRules, which can target specific tags. Tags must only contain lowercase letters, numbers, and hyphens.
	Tags string `url:"tags,omitempty" json:"tags,omitempty" path:"tags"`
	// The name of the Partner. Required.
	Name string `url:"name" json:"name" path:"name"`
	// The root folder path for this Partner. Required.
	RootFolder string `url:"root_folder" json:"root_folder" path:"root_folder"`
	// ID of the Workspace associated with this Partner.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// PartnerUpdateParams contains the request parameters for PATCH /partners/{id}.
type PartnerUpdateParams struct {
	// Partner ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// AI Assistant Personality ID assigned to this Partner, if any. Users in the Partner inherit it unless a direct per-user assignment overrides it.
	AiAssistantPersonalityId int64 `url:"ai_assistant_personality_id,omitempty" json:"ai_assistant_personality_id,omitempty" path:"ai_assistant_personality_id"`
	// A list of allowed IPs for this Partner. Newline delimited. Partner User IP access is allowed when the IP matches the Partner, User, or Site allowed IP lists.
	AllowedIps string `url:"allowed_ips,omitempty" json:"allowed_ips,omitempty" path:"allowed_ips"`
	// Allow Partner Admins to change Two-Factor Authentication requirements for Partner Users.
	AllowBypassing2faPolicies *bool `url:"allow_bypassing_2fa_policies,omitempty" json:"allow_bypassing_2fa_policies,omitempty" path:"allow_bypassing_2fa_policies"`
	// Allow Partner Admins to change or reset credentials for users belonging to this Partner.
	AllowCredentialChanges *bool `url:"allow_credential_changes,omitempty" json:"allow_credential_changes,omitempty" path:"allow_credential_changes"`
	// Allow Partner Admins to provide GPG keys.
	AllowProvidingGpgKeys *bool `url:"allow_providing_gpg_keys,omitempty" json:"allow_providing_gpg_keys,omitempty" path:"allow_providing_gpg_keys"`
	// Allow Partner Admins to create users.
	AllowUserCreation *bool `url:"allow_user_creation,omitempty" json:"allow_user_creation,omitempty" path:"allow_user_creation"`
	// When `true`, emails sent to Partner users are copied to the responsible User or Group.
	CcEmailsToResponsibleParty *bool `url:"cc_emails_to_responsible_party,omitempty" json:"cc_emails_to_responsible_party,omitempty" path:"cc_emails_to_responsible_party"`
	// Notes about this Partner.
	Notes string `url:"notes,omitempty" json:"notes,omitempty" path:"notes"`
	// ID of the Partner Channel Template assigned to this Partner.
	PartnerChannelTemplateId int64 `url:"partner_channel_template_id,omitempty" json:"partner_channel_template_id,omitempty" path:"partner_channel_template_id"`
	// ID of the Group responsible for this Partner.
	ResponsibleGroupId int64 `url:"responsible_group_id,omitempty" json:"responsible_group_id,omitempty" path:"responsible_group_id"`
	// ID of the User responsible for this Partner.
	ResponsibleUserId int64 `url:"responsible_user_id,omitempty" json:"responsible_user_id,omitempty" path:"responsible_user_id"`
	// Show Partner users a simplified home page built from this Partner's Channels.
	ShowPartnerChannelHomePage *bool `url:"show_partner_channel_home_page,omitempty" json:"show_partner_channel_home_page,omitempty" path:"show_partner_channel_home_page"`
	// Comma-separated list of Tags for this Partner. Tags are used for other features, such as UserLifecycleRules, which can target specific tags. Tags must only contain lowercase letters, numbers, and hyphens.
	Tags string `url:"tags,omitempty" json:"tags,omitempty" path:"tags"`
	// The name of the Partner.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// This site's role for this Partner in Connected Sites relationships. `host` is a Partner this site configured. `guest` is a Partner created by approving another site's connection request; it has no root folder and cannot hold users, permissions, or Partner Channels, or host a connection. `host_and_guest` is a configured Partner that is also the guest side of a connection. Promote a `guest` Partner by setting this to `host_and_guest` together with a `root_folder`. Values: host, guest, host_and_guest.
	PartnershipRole PartnerPartnershipRoleEnum `url:"partnership_role,omitempty" json:"partnership_role,omitempty" path:"partnership_role"`
	// The root folder path for this Partner.
	RootFolder string `url:"root_folder,omitempty" json:"root_folder,omitempty" path:"root_folder"`
}

// PartnerDeleteParams contains the request parameters for DELETE /partners/{id}.
type PartnerDeleteParams struct {
	// Partner ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (p *Partner) UnmarshalJSON(data []byte) error {
	type partner Partner
	var v partner
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*p = Partner(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (p *PartnerCollection) UnmarshalJSON(data []byte) error {
	type partners PartnerCollection
	var v partners
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*p = PartnerCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (p *PartnerCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
