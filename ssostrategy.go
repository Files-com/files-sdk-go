package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// SsoStrategy is a Files.com API resource.
type SsoStrategy struct {
	// SSO Protocol.
	Protocol string `json:"protocol,omitempty" path:"protocol,omitempty" url:"protocol,omitempty"`
	// Provider name.
	Provider string `json:"provider,omitempty" path:"provider,omitempty" url:"provider,omitempty"`
	// Custom label for the SSO provider on the login page.
	Label string `json:"label,omitempty" path:"label,omitempty" url:"label,omitempty"`
	// URL holding a custom logo for the SSO provider on the login page.
	LogoUrl string `json:"logo_url,omitempty" path:"logo_url,omitempty" url:"logo_url,omitempty"`
	// ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Is strategy enabled? This may become automatically set to `false` after a high number and duration of failures.
	Enabled *bool `json:"enabled,omitempty" path:"enabled,omitempty" url:"enabled,omitempty"`
	// Count of users with this SSO Strategy.
	UserCount int64 `json:"user_count,omitempty" path:"user_count,omitempty" url:"user_count,omitempty"`
	// Identity provider sha256 cert fingerprint if saml_provider_metadata_url is not available.
	SamlProviderCertFingerprint string `json:"saml_provider_cert_fingerprint,omitempty" path:"saml_provider_cert_fingerprint,omitempty" url:"saml_provider_cert_fingerprint,omitempty"`
	// Identity provider issuer url.
	SamlProviderIssuerUrl string `json:"saml_provider_issuer_url,omitempty" path:"saml_provider_issuer_url,omitempty" url:"saml_provider_issuer_url,omitempty"`
	// Custom identity provider metadata.
	SamlProviderMetadataContent string `json:"saml_provider_metadata_content,omitempty" path:"saml_provider_metadata_content,omitempty" url:"saml_provider_metadata_content,omitempty"`
	// Metadata URL for the SAML identity provider.
	SamlProviderMetadataUrl string `json:"saml_provider_metadata_url,omitempty" path:"saml_provider_metadata_url,omitempty" url:"saml_provider_metadata_url,omitempty"`
	// Identity provider SLO endpoint.
	SamlProviderSloTargetUrl string `json:"saml_provider_slo_target_url,omitempty" path:"saml_provider_slo_target_url,omitempty" url:"saml_provider_slo_target_url,omitempty"`
	// Identity provider SSO endpoint if saml_provider_metadata_url is not available.
	SamlProviderSsoTargetUrl string `json:"saml_provider_sso_target_url,omitempty" path:"saml_provider_sso_target_url,omitempty" url:"saml_provider_sso_target_url,omitempty"`
	// SCIM authentication type.
	ScimAuthenticationMethod string `json:"scim_authentication_method,omitempty" path:"scim_authentication_method,omitempty" url:"scim_authentication_method,omitempty"`
	// SCIM username.
	ScimUsername string `json:"scim_username,omitempty" path:"scim_username,omitempty" url:"scim_username,omitempty"`
	// SCIM OAuth Access Token. Only returned when the token is first issued or reset.
	ScimOauthAccessToken string `json:"scim_oauth_access_token,omitempty" path:"scim_oauth_access_token,omitempty" url:"scim_oauth_access_token,omitempty"`
	// SCIM OAuth Access Token Expiration Time.
	ScimOauthAccessTokenExpiresAt string `json:"scim_oauth_access_token_expires_at,omitempty" path:"scim_oauth_access_token_expires_at,omitempty" url:"scim_oauth_access_token_expires_at,omitempty"`
	// Subdomain or domain name for your auth provider. Example: `https://[subdomain].okta.com/`.
	Subdomain string `json:"subdomain,omitempty" path:"subdomain,omitempty" url:"subdomain,omitempty"`
	// Auto-provision users?
	ProvisionUsers *bool `json:"provision_users,omitempty" path:"provision_users,omitempty" url:"provision_users,omitempty"`
	// Auto-provision group membership based on group memberships on the SSO side?
	ProvisionGroups *bool `json:"provision_groups,omitempty" path:"provision_groups,omitempty" url:"provision_groups,omitempty"`
	// Auto-deprovision users?
	DeprovisionUsers *bool `json:"deprovision_users,omitempty" path:"deprovision_users,omitempty" url:"deprovision_users,omitempty"`
	// Auto-deprovision group membership based on group memberships on the SSO side?
	DeprovisionGroups *bool `json:"deprovision_groups,omitempty" path:"deprovision_groups,omitempty" url:"deprovision_groups,omitempty"`
	// Method used for deprovisioning users.
	DeprovisionBehavior string `json:"deprovision_behavior,omitempty" path:"deprovision_behavior,omitempty" url:"deprovision_behavior,omitempty"`
	// Comma-separated list of group names for groups to automatically add all auto-provisioned users to.
	ProvisionGroupDefault string `json:"provision_group_default,omitempty" path:"provision_group_default,omitempty" url:"provision_group_default,omitempty"`
	// Comma-separated list of group names for groups (with optional wildcards) that will be excluded from auto-provisioning.
	ProvisionGroupExclusion string `json:"provision_group_exclusion,omitempty" path:"provision_group_exclusion,omitempty" url:"provision_group_exclusion,omitempty"`
	// Comma-separated list of group names for groups (with optional wildcards) that will be auto-provisioned.
	ProvisionGroupInclusion string `json:"provision_group_inclusion,omitempty" path:"provision_group_inclusion,omitempty" url:"provision_group_inclusion,omitempty"`
	// Comma or newline separated list of group names (with optional wildcards) to require membership for user provisioning.
	ProvisionGroupRequired string `json:"provision_group_required,omitempty" path:"provision_group_required,omitempty" url:"provision_group_required,omitempty"`
	// Comma-separated list of group names whose members will be created with email_signup authentication.
	ProvisionEmailSignupGroups string `json:"provision_email_signup_groups,omitempty" path:"provision_email_signup_groups,omitempty" url:"provision_email_signup_groups,omitempty"`
	// Comma-separated list of group names used to grant or remove Read-Only Site Admin access for users on this SSO provider. When this list is nonempty and a directory sync changes a user's group membership, membership in at least one listed group grants this access; membership in none removes it, including access granted manually. Directory syncs include SCIM requests from any provider on the site, LDAP syncs, and sign-ins that set the user's groups from the identity provider. Membership changes made any other way, including through the web app, API, SDKs, or CLI, do not affect administrator access.
	ProvisionReadonlySiteAdminGroups string `json:"provision_readonly_site_admin_groups,omitempty" path:"provision_readonly_site_admin_groups,omitempty" url:"provision_readonly_site_admin_groups,omitempty"`
	// Comma-separated list of group names used to grant or remove Site Admin access for users on this SSO provider. When this list is nonempty and a directory sync changes a user's group membership, membership in at least one listed group grants this access; membership in none removes it, including access granted manually. Directory syncs include SCIM requests from any provider on the site, LDAP syncs, and sign-ins that set the user's groups from the identity provider. Membership changes made any other way, including through the web app, API, SDKs, or CLI, do not affect administrator access.
	ProvisionSiteAdminGroups string `json:"provision_site_admin_groups,omitempty" path:"provision_site_admin_groups,omitempty" url:"provision_site_admin_groups,omitempty"`
	// Comma-separated list of group names whose members will be provisioned as Group Admins.
	ProvisionGroupAdminGroups      string `json:"provision_group_admin_groups,omitempty" path:"provision_group_admin_groups,omitempty" url:"provision_group_admin_groups,omitempty"`
	ProvisionAttachmentsPermission *bool  `json:"provision_attachments_permission,omitempty" path:"provision_attachments_permission,omitempty" url:"provision_attachments_permission,omitempty"`
	// Auto-provisioned users get WebDAV permission?
	ProvisionDavPermission *bool `json:"provision_dav_permission,omitempty" path:"provision_dav_permission,omitempty" url:"provision_dav_permission,omitempty"`
	// Auto-provisioned users get FTP permission?
	ProvisionFtpPermission *bool `json:"provision_ftp_permission,omitempty" path:"provision_ftp_permission,omitempty" url:"provision_ftp_permission,omitempty"`
	// Auto-provisioned users get SFTP permission?
	ProvisionSftpPermission *bool `json:"provision_sftp_permission,omitempty" path:"provision_sftp_permission,omitempty" url:"provision_sftp_permission,omitempty"`
	// Default time zone for auto provisioned users.
	ProvisionTimeZone string `json:"provision_time_zone,omitempty" path:"provision_time_zone,omitempty" url:"provision_time_zone,omitempty"`
	// Default company for auto provisioned users.
	ProvisionCompany string `json:"provision_company,omitempty" path:"provision_company,omitempty" url:"provision_company,omitempty"`
	// 2FA required setting for auto provisioned users. `use_system_setting` uses the site-wide setting, including SSO exemptions. `always_require` and `never_require` override the site-wide setting when user-level overrides are allowed.
	ProvisionRequire2fa string `json:"provision_require_2fa,omitempty" path:"provision_require_2fa,omitempty" url:"provision_require_2fa,omitempty"`
	// File System layout to use for auto provisioned users.
	ProvisionFilesystemLayout string `json:"provision_filesystem_layout,omitempty" path:"provision_filesystem_layout,omitempty" url:"provision_filesystem_layout,omitempty"`
	// URL-friendly, unique identifier for Azure SAML configuration.
	ProviderIdentifier string `json:"provider_identifier,omitempty" path:"provider_identifier,omitempty" url:"provider_identifier,omitempty"`
	// Base DN for looking up users in LDAP server.
	LdapBaseDn string `json:"ldap_base_dn,omitempty" path:"ldap_base_dn,omitempty" url:"ldap_base_dn,omitempty"`
	// Domain name that will be appended to LDAP usernames.
	LdapDomain string `json:"ldap_domain,omitempty" path:"ldap_domain,omitempty" url:"ldap_domain,omitempty"`
	// Should this strategy be displayed on the login page?
	DisplayOnLoginPage *bool `json:"display_on_login_page,omitempty" path:"display_on_login_page,omitempty" url:"display_on_login_page,omitempty"`
	// LDAP host.
	LdapHost string `json:"ldap_host,omitempty" path:"ldap_host,omitempty" url:"ldap_host,omitempty"`
	// LDAP backup host.
	LdapHost2 string `json:"ldap_host_2,omitempty" path:"ldap_host_2,omitempty" url:"ldap_host_2,omitempty"`
	// LDAP backup host.
	LdapHost3 string `json:"ldap_host_3,omitempty" path:"ldap_host_3,omitempty" url:"ldap_host_3,omitempty"`
	// LDAP port.
	LdapPort int64 `json:"ldap_port,omitempty" path:"ldap_port,omitempty" url:"ldap_port,omitempty"`
	// Use LDAP server settings for scheduled provisioning while using this SSO provider for authentication?
	LdapProvisioningEnabled *bool `json:"ldap_provisioning_enabled,omitempty" path:"ldap_provisioning_enabled,omitempty" url:"ldap_provisioning_enabled,omitempty"`
	// Use secure LDAP?
	LdapSecure *bool `json:"ldap_secure,omitempty" path:"ldap_secure,omitempty" url:"ldap_secure,omitempty"`
	// How to validate the LDAP server certificate. `require_match` validates the certificate chain and hostname; `allow_any` disables certificate validation.
	LdapServerCertificate string `json:"ldap_server_certificate,omitempty" path:"ldap_server_certificate,omitempty" url:"ldap_server_certificate,omitempty"`
	// LDAP server type.
	LdapType string `json:"ldap_type,omitempty" path:"ldap_type,omitempty" url:"ldap_type,omitempty"`
	// Username for signing in to LDAP server.
	LdapUsername string `json:"ldap_username,omitempty" path:"ldap_username,omitempty" url:"ldap_username,omitempty"`
	// LDAP username field.
	LdapUsernameField string `json:"ldap_username_field,omitempty" path:"ldap_username_field,omitempty" url:"ldap_username_field,omitempty"`
}

// Identifier returns the resource ID.
func (s SsoStrategy) Identifier() interface{} {
	return s.Id
}

// SsoStrategyCollection is a list of SsoStrategy resources.
type SsoStrategyCollection []SsoStrategy

// SsoStrategyListParams contains the request parameters for GET /sso_strategies.
type SsoStrategyListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are .
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// SsoStrategyFindParams contains the request parameters for GET /sso_strategies/{id}.
type SsoStrategyFindParams struct {
	// Sso Strategy ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// SsoStrategySyncParams contains the request parameters for POST /sso_strategies/{id}/sync.
//
// Synchronize provisioning data with the SSO remote server
type SsoStrategySyncParams struct {
	// Sso Strategy ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (s *SsoStrategy) UnmarshalJSON(data []byte) error {
	type ssoStrategy SsoStrategy
	var v ssoStrategy
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = SsoStrategy(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (s *SsoStrategyCollection) UnmarshalJSON(data []byte) error {
	type ssoStrategys SsoStrategyCollection
	var v ssoStrategys
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SsoStrategyCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (s *SsoStrategyCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
