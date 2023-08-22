package files_sdk

import (
	"encoding/json"
	"io"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type SsoStrategy struct {
	Protocol                       string    `json:"protocol,omitempty" path:"protocol,omitempty" url:"protocol,omitempty"`
	Provider                       string    `json:"provider,omitempty" path:"provider,omitempty" url:"provider,omitempty"`
	Label                          string    `json:"label,omitempty" path:"label,omitempty" url:"label,omitempty"`
	LogoUrl                        string    `json:"logo_url,omitempty" path:"logo_url,omitempty" url:"logo_url,omitempty"`
	Id                             int64     `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	SamlProviderCertFingerprint    string    `json:"saml_provider_cert_fingerprint,omitempty" path:"saml_provider_cert_fingerprint,omitempty" url:"saml_provider_cert_fingerprint,omitempty"`
	SamlProviderIssuerUrl          string    `json:"saml_provider_issuer_url,omitempty" path:"saml_provider_issuer_url,omitempty" url:"saml_provider_issuer_url,omitempty"`
	SamlProviderMetadataContent    string    `json:"saml_provider_metadata_content,omitempty" path:"saml_provider_metadata_content,omitempty" url:"saml_provider_metadata_content,omitempty"`
	SamlProviderMetadataUrl        string    `json:"saml_provider_metadata_url,omitempty" path:"saml_provider_metadata_url,omitempty" url:"saml_provider_metadata_url,omitempty"`
	SamlProviderSloTargetUrl       string    `json:"saml_provider_slo_target_url,omitempty" path:"saml_provider_slo_target_url,omitempty" url:"saml_provider_slo_target_url,omitempty"`
	SamlProviderSsoTargetUrl       string    `json:"saml_provider_sso_target_url,omitempty" path:"saml_provider_sso_target_url,omitempty" url:"saml_provider_sso_target_url,omitempty"`
	ScimAuthenticationMethod       string    `json:"scim_authentication_method,omitempty" path:"scim_authentication_method,omitempty" url:"scim_authentication_method,omitempty"`
	ScimUsername                   string    `json:"scim_username,omitempty" path:"scim_username,omitempty" url:"scim_username,omitempty"`
	ScimOauthAccessToken           string    `json:"scim_oauth_access_token,omitempty" path:"scim_oauth_access_token,omitempty" url:"scim_oauth_access_token,omitempty"`
	ScimOauthAccessTokenExpiresAt  string    `json:"scim_oauth_access_token_expires_at,omitempty" path:"scim_oauth_access_token_expires_at,omitempty" url:"scim_oauth_access_token_expires_at,omitempty"`
	Subdomain                      string    `json:"subdomain,omitempty" path:"subdomain,omitempty" url:"subdomain,omitempty"`
	ProvisionUsers                 *bool     `json:"provision_users,omitempty" path:"provision_users,omitempty" url:"provision_users,omitempty"`
	ProvisionGroups                *bool     `json:"provision_groups,omitempty" path:"provision_groups,omitempty" url:"provision_groups,omitempty"`
	DeprovisionUsers               *bool     `json:"deprovision_users,omitempty" path:"deprovision_users,omitempty" url:"deprovision_users,omitempty"`
	DeprovisionGroups              *bool     `json:"deprovision_groups,omitempty" path:"deprovision_groups,omitempty" url:"deprovision_groups,omitempty"`
	DeprovisionBehavior            string    `json:"deprovision_behavior,omitempty" path:"deprovision_behavior,omitempty" url:"deprovision_behavior,omitempty"`
	ProvisionGroupDefault          string    `json:"provision_group_default,omitempty" path:"provision_group_default,omitempty" url:"provision_group_default,omitempty"`
	ProvisionGroupExclusion        string    `json:"provision_group_exclusion,omitempty" path:"provision_group_exclusion,omitempty" url:"provision_group_exclusion,omitempty"`
	ProvisionGroupInclusion        string    `json:"provision_group_inclusion,omitempty" path:"provision_group_inclusion,omitempty" url:"provision_group_inclusion,omitempty"`
	ProvisionGroupRequired         string    `json:"provision_group_required,omitempty" path:"provision_group_required,omitempty" url:"provision_group_required,omitempty"`
	ProvisionEmailSignupGroups     string    `json:"provision_email_signup_groups,omitempty" path:"provision_email_signup_groups,omitempty" url:"provision_email_signup_groups,omitempty"`
	ProvisionSiteAdminGroups       string    `json:"provision_site_admin_groups,omitempty" path:"provision_site_admin_groups,omitempty" url:"provision_site_admin_groups,omitempty"`
	ProvisionGroupAdminGroups      string    `json:"provision_group_admin_groups,omitempty" path:"provision_group_admin_groups,omitempty" url:"provision_group_admin_groups,omitempty"`
	ProvisionAttachmentsPermission *bool     `json:"provision_attachments_permission,omitempty" path:"provision_attachments_permission,omitempty" url:"provision_attachments_permission,omitempty"`
	ProvisionDavPermission         *bool     `json:"provision_dav_permission,omitempty" path:"provision_dav_permission,omitempty" url:"provision_dav_permission,omitempty"`
	ProvisionFtpPermission         *bool     `json:"provision_ftp_permission,omitempty" path:"provision_ftp_permission,omitempty" url:"provision_ftp_permission,omitempty"`
	ProvisionSftpPermission        *bool     `json:"provision_sftp_permission,omitempty" path:"provision_sftp_permission,omitempty" url:"provision_sftp_permission,omitempty"`
	ProvisionTimeZone              string    `json:"provision_time_zone,omitempty" path:"provision_time_zone,omitempty" url:"provision_time_zone,omitempty"`
	ProvisionCompany               string    `json:"provision_company,omitempty" path:"provision_company,omitempty" url:"provision_company,omitempty"`
	LdapBaseDn                     string    `json:"ldap_base_dn,omitempty" path:"ldap_base_dn,omitempty" url:"ldap_base_dn,omitempty"`
	LdapDomain                     string    `json:"ldap_domain,omitempty" path:"ldap_domain,omitempty" url:"ldap_domain,omitempty"`
	Enabled                        *bool     `json:"enabled,omitempty" path:"enabled,omitempty" url:"enabled,omitempty"`
	LdapHost                       string    `json:"ldap_host,omitempty" path:"ldap_host,omitempty" url:"ldap_host,omitempty"`
	LdapHost2                      string    `json:"ldap_host_2,omitempty" path:"ldap_host_2,omitempty" url:"ldap_host_2,omitempty"`
	LdapHost3                      string    `json:"ldap_host_3,omitempty" path:"ldap_host_3,omitempty" url:"ldap_host_3,omitempty"`
	LdapPort                       int64     `json:"ldap_port,omitempty" path:"ldap_port,omitempty" url:"ldap_port,omitempty"`
	LdapSecure                     *bool     `json:"ldap_secure,omitempty" path:"ldap_secure,omitempty" url:"ldap_secure,omitempty"`
	LdapUsername                   string    `json:"ldap_username,omitempty" path:"ldap_username,omitempty" url:"ldap_username,omitempty"`
	LdapUsernameField              string    `json:"ldap_username_field,omitempty" path:"ldap_username_field,omitempty" url:"ldap_username_field,omitempty"`
	ClientId                       string    `json:"client_id,omitempty" path:"client_id,omitempty" url:"client_id,omitempty"`
	ClientSecret                   string    `json:"client_secret,omitempty" path:"client_secret,omitempty" url:"client_secret,omitempty"`
	ScimPassword                   string    `json:"scim_password,omitempty" path:"scim_password,omitempty" url:"scim_password,omitempty"`
	ResetScimOauthAccessToken      *bool     `json:"reset_scim_oauth_access_token,omitempty" path:"reset_scim_oauth_access_token,omitempty" url:"reset_scim_oauth_access_token,omitempty"`
	LogoFile                       io.Reader `json:"logo_file,omitempty" path:"logo_file,omitempty" url:"logo_file,omitempty"`
	LogoDelete                     *bool     `json:"logo_delete,omitempty" path:"logo_delete,omitempty" url:"logo_delete,omitempty"`
	LdapPassword                   string    `json:"ldap_password,omitempty" path:"ldap_password,omitempty" url:"ldap_password,omitempty"`
}

func (s SsoStrategy) Identifier() interface{} {
	return s.Id
}

type SsoStrategyCollection []SsoStrategy

type SsoStrategyScimAuthenticationMethodEnum string

func (u SsoStrategyScimAuthenticationMethodEnum) String() string {
	return string(u)
}

func (u SsoStrategyScimAuthenticationMethodEnum) Enum() map[string]SsoStrategyScimAuthenticationMethodEnum {
	return map[string]SsoStrategyScimAuthenticationMethodEnum{
		"none":  SsoStrategyScimAuthenticationMethodEnum("none"),
		"basic": SsoStrategyScimAuthenticationMethodEnum("basic"),
		"token": SsoStrategyScimAuthenticationMethodEnum("token"),
	}
}

type SsoStrategyDeprovisionBehaviorEnum string

func (u SsoStrategyDeprovisionBehaviorEnum) String() string {
	return string(u)
}

func (u SsoStrategyDeprovisionBehaviorEnum) Enum() map[string]SsoStrategyDeprovisionBehaviorEnum {
	return map[string]SsoStrategyDeprovisionBehaviorEnum{
		"disable": SsoStrategyDeprovisionBehaviorEnum("disable"),
		"delete":  SsoStrategyDeprovisionBehaviorEnum("delete"),
	}
}

type SsoStrategyLdapUsernameFieldEnum string

func (u SsoStrategyLdapUsernameFieldEnum) String() string {
	return string(u)
}

func (u SsoStrategyLdapUsernameFieldEnum) Enum() map[string]SsoStrategyLdapUsernameFieldEnum {
	return map[string]SsoStrategyLdapUsernameFieldEnum{
		"sAMAccountName":    SsoStrategyLdapUsernameFieldEnum("sAMAccountName"),
		"userPrincipalName": SsoStrategyLdapUsernameFieldEnum("userPrincipalName"),
	}
}

type SsoStrategyListParams struct {
	Action string `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	ListParams
}

type SsoStrategyFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type SsoStrategyCreateParams struct {
	Provider                       string                                  `url:"provider,omitempty" required:"false" json:"provider,omitempty" path:"provider"`
	Subdomain                      string                                  `url:"subdomain,omitempty" required:"false" json:"subdomain,omitempty" path:"subdomain"`
	ClientId                       string                                  `url:"client_id,omitempty" required:"false" json:"client_id,omitempty" path:"client_id"`
	ClientSecret                   string                                  `url:"client_secret,omitempty" required:"false" json:"client_secret,omitempty" path:"client_secret"`
	SamlProviderMetadataContent    string                                  `url:"saml_provider_metadata_content,omitempty" required:"false" json:"saml_provider_metadata_content,omitempty" path:"saml_provider_metadata_content"`
	SamlProviderMetadataUrl        string                                  `url:"saml_provider_metadata_url,omitempty" required:"false" json:"saml_provider_metadata_url,omitempty" path:"saml_provider_metadata_url"`
	SamlProviderCertFingerprint    string                                  `url:"saml_provider_cert_fingerprint,omitempty" required:"false" json:"saml_provider_cert_fingerprint,omitempty" path:"saml_provider_cert_fingerprint"`
	SamlProviderIssuerUrl          string                                  `url:"saml_provider_issuer_url,omitempty" required:"false" json:"saml_provider_issuer_url,omitempty" path:"saml_provider_issuer_url"`
	SamlProviderSloTargetUrl       string                                  `url:"saml_provider_slo_target_url,omitempty" required:"false" json:"saml_provider_slo_target_url,omitempty" path:"saml_provider_slo_target_url"`
	SamlProviderSsoTargetUrl       string                                  `url:"saml_provider_sso_target_url,omitempty" required:"false" json:"saml_provider_sso_target_url,omitempty" path:"saml_provider_sso_target_url"`
	ScimAuthenticationMethod       SsoStrategyScimAuthenticationMethodEnum `url:"scim_authentication_method,omitempty" required:"false" json:"scim_authentication_method,omitempty" path:"scim_authentication_method"`
	ScimPassword                   string                                  `url:"scim_password,omitempty" required:"false" json:"scim_password,omitempty" path:"scim_password"`
	ScimUsername                   string                                  `url:"scim_username,omitempty" required:"false" json:"scim_username,omitempty" path:"scim_username"`
	ScimOauthAccessTokenExpiresAt  *time.Time                              `url:"scim_oauth_access_token_expires_at,omitempty" required:"false" json:"scim_oauth_access_token_expires_at,omitempty" path:"scim_oauth_access_token_expires_at"`
	ResetScimOauthAccessToken      *bool                                   `url:"reset_scim_oauth_access_token,omitempty" required:"false" json:"reset_scim_oauth_access_token,omitempty" path:"reset_scim_oauth_access_token"`
	Protocol                       string                                  `url:"protocol,omitempty" required:"false" json:"protocol,omitempty" path:"protocol"`
	ProvisionUsers                 *bool                                   `url:"provision_users,omitempty" required:"false" json:"provision_users,omitempty" path:"provision_users"`
	ProvisionGroups                *bool                                   `url:"provision_groups,omitempty" required:"false" json:"provision_groups,omitempty" path:"provision_groups"`
	DeprovisionUsers               *bool                                   `url:"deprovision_users,omitempty" required:"false" json:"deprovision_users,omitempty" path:"deprovision_users"`
	DeprovisionGroups              *bool                                   `url:"deprovision_groups,omitempty" required:"false" json:"deprovision_groups,omitempty" path:"deprovision_groups"`
	DeprovisionBehavior            SsoStrategyDeprovisionBehaviorEnum      `url:"deprovision_behavior,omitempty" required:"false" json:"deprovision_behavior,omitempty" path:"deprovision_behavior"`
	ProvisionGroupDefault          string                                  `url:"provision_group_default,omitempty" required:"false" json:"provision_group_default,omitempty" path:"provision_group_default"`
	ProvisionGroupExclusion        string                                  `url:"provision_group_exclusion,omitempty" required:"false" json:"provision_group_exclusion,omitempty" path:"provision_group_exclusion"`
	ProvisionGroupInclusion        string                                  `url:"provision_group_inclusion,omitempty" required:"false" json:"provision_group_inclusion,omitempty" path:"provision_group_inclusion"`
	ProvisionGroupRequired         string                                  `url:"provision_group_required,omitempty" required:"false" json:"provision_group_required,omitempty" path:"provision_group_required"`
	ProvisionAttachmentsPermission *bool                                   `url:"provision_attachments_permission,omitempty" required:"false" json:"provision_attachments_permission,omitempty" path:"provision_attachments_permission"`
	ProvisionDavPermission         *bool                                   `url:"provision_dav_permission,omitempty" required:"false" json:"provision_dav_permission,omitempty" path:"provision_dav_permission"`
	ProvisionFtpPermission         *bool                                   `url:"provision_ftp_permission,omitempty" required:"false" json:"provision_ftp_permission,omitempty" path:"provision_ftp_permission"`
	ProvisionSftpPermission        *bool                                   `url:"provision_sftp_permission,omitempty" required:"false" json:"provision_sftp_permission,omitempty" path:"provision_sftp_permission"`
	ProvisionEmailSignupGroups     string                                  `url:"provision_email_signup_groups,omitempty" required:"false" json:"provision_email_signup_groups,omitempty" path:"provision_email_signup_groups"`
	ProvisionSiteAdminGroups       string                                  `url:"provision_site_admin_groups,omitempty" required:"false" json:"provision_site_admin_groups,omitempty" path:"provision_site_admin_groups"`
	ProvisionGroupAdminGroups      string                                  `url:"provision_group_admin_groups,omitempty" required:"false" json:"provision_group_admin_groups,omitempty" path:"provision_group_admin_groups"`
	ProvisionTimeZone              string                                  `url:"provision_time_zone,omitempty" required:"false" json:"provision_time_zone,omitempty" path:"provision_time_zone"`
	ProvisionCompany               string                                  `url:"provision_company,omitempty" required:"false" json:"provision_company,omitempty" path:"provision_company"`
	Label                          string                                  `url:"label,omitempty" required:"false" json:"label,omitempty" path:"label"`
	LogoFile                       io.Writer                               `url:"logo_file,omitempty" required:"false" json:"logo_file,omitempty" path:"logo_file"`
	LogoDelete                     *bool                                   `url:"logo_delete,omitempty" required:"false" json:"logo_delete,omitempty" path:"logo_delete"`
	LdapBaseDn                     string                                  `url:"ldap_base_dn,omitempty" required:"false" json:"ldap_base_dn,omitempty" path:"ldap_base_dn"`
	LdapDomain                     string                                  `url:"ldap_domain,omitempty" required:"false" json:"ldap_domain,omitempty" path:"ldap_domain"`
	LdapHost                       string                                  `url:"ldap_host,omitempty" required:"false" json:"ldap_host,omitempty" path:"ldap_host"`
	LdapHost2                      string                                  `url:"ldap_host_2,omitempty" required:"false" json:"ldap_host_2,omitempty" path:"ldap_host_2"`
	LdapHost3                      string                                  `url:"ldap_host_3,omitempty" required:"false" json:"ldap_host_3,omitempty" path:"ldap_host_3"`
	LdapPort                       int64                                   `url:"ldap_port,omitempty" required:"false" json:"ldap_port,omitempty" path:"ldap_port"`
	LdapSecure                     *bool                                   `url:"ldap_secure,omitempty" required:"false" json:"ldap_secure,omitempty" path:"ldap_secure"`
	LdapUsername                   string                                  `url:"ldap_username,omitempty" required:"false" json:"ldap_username,omitempty" path:"ldap_username"`
	LdapPassword                   string                                  `url:"ldap_password,omitempty" required:"false" json:"ldap_password,omitempty" path:"ldap_password"`
	LdapUsernameField              SsoStrategyLdapUsernameFieldEnum        `url:"ldap_username_field,omitempty" required:"false" json:"ldap_username_field,omitempty" path:"ldap_username_field"`
	Enabled                        *bool                                   `url:"enabled,omitempty" required:"false" json:"enabled,omitempty" path:"enabled"`
}

// Synchronize provisioning data with the SSO remote server
type SsoStrategySyncParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type SsoStrategyUpdateParams struct {
	Id                             int64                                   `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
	Provider                       string                                  `url:"provider,omitempty" required:"false" json:"provider,omitempty" path:"provider"`
	Subdomain                      string                                  `url:"subdomain,omitempty" required:"false" json:"subdomain,omitempty" path:"subdomain"`
	ClientId                       string                                  `url:"client_id,omitempty" required:"false" json:"client_id,omitempty" path:"client_id"`
	ClientSecret                   string                                  `url:"client_secret,omitempty" required:"false" json:"client_secret,omitempty" path:"client_secret"`
	SamlProviderMetadataContent    string                                  `url:"saml_provider_metadata_content,omitempty" required:"false" json:"saml_provider_metadata_content,omitempty" path:"saml_provider_metadata_content"`
	SamlProviderMetadataUrl        string                                  `url:"saml_provider_metadata_url,omitempty" required:"false" json:"saml_provider_metadata_url,omitempty" path:"saml_provider_metadata_url"`
	SamlProviderCertFingerprint    string                                  `url:"saml_provider_cert_fingerprint,omitempty" required:"false" json:"saml_provider_cert_fingerprint,omitempty" path:"saml_provider_cert_fingerprint"`
	SamlProviderIssuerUrl          string                                  `url:"saml_provider_issuer_url,omitempty" required:"false" json:"saml_provider_issuer_url,omitempty" path:"saml_provider_issuer_url"`
	SamlProviderSloTargetUrl       string                                  `url:"saml_provider_slo_target_url,omitempty" required:"false" json:"saml_provider_slo_target_url,omitempty" path:"saml_provider_slo_target_url"`
	SamlProviderSsoTargetUrl       string                                  `url:"saml_provider_sso_target_url,omitempty" required:"false" json:"saml_provider_sso_target_url,omitempty" path:"saml_provider_sso_target_url"`
	ScimAuthenticationMethod       SsoStrategyScimAuthenticationMethodEnum `url:"scim_authentication_method,omitempty" required:"false" json:"scim_authentication_method,omitempty" path:"scim_authentication_method"`
	ScimPassword                   string                                  `url:"scim_password,omitempty" required:"false" json:"scim_password,omitempty" path:"scim_password"`
	ScimUsername                   string                                  `url:"scim_username,omitempty" required:"false" json:"scim_username,omitempty" path:"scim_username"`
	ScimOauthAccessTokenExpiresAt  *time.Time                              `url:"scim_oauth_access_token_expires_at,omitempty" required:"false" json:"scim_oauth_access_token_expires_at,omitempty" path:"scim_oauth_access_token_expires_at"`
	ResetScimOauthAccessToken      *bool                                   `url:"reset_scim_oauth_access_token,omitempty" required:"false" json:"reset_scim_oauth_access_token,omitempty" path:"reset_scim_oauth_access_token"`
	Protocol                       string                                  `url:"protocol,omitempty" required:"false" json:"protocol,omitempty" path:"protocol"`
	ProvisionUsers                 *bool                                   `url:"provision_users,omitempty" required:"false" json:"provision_users,omitempty" path:"provision_users"`
	ProvisionGroups                *bool                                   `url:"provision_groups,omitempty" required:"false" json:"provision_groups,omitempty" path:"provision_groups"`
	DeprovisionUsers               *bool                                   `url:"deprovision_users,omitempty" required:"false" json:"deprovision_users,omitempty" path:"deprovision_users"`
	DeprovisionGroups              *bool                                   `url:"deprovision_groups,omitempty" required:"false" json:"deprovision_groups,omitempty" path:"deprovision_groups"`
	DeprovisionBehavior            SsoStrategyDeprovisionBehaviorEnum      `url:"deprovision_behavior,omitempty" required:"false" json:"deprovision_behavior,omitempty" path:"deprovision_behavior"`
	ProvisionGroupDefault          string                                  `url:"provision_group_default,omitempty" required:"false" json:"provision_group_default,omitempty" path:"provision_group_default"`
	ProvisionGroupExclusion        string                                  `url:"provision_group_exclusion,omitempty" required:"false" json:"provision_group_exclusion,omitempty" path:"provision_group_exclusion"`
	ProvisionGroupInclusion        string                                  `url:"provision_group_inclusion,omitempty" required:"false" json:"provision_group_inclusion,omitempty" path:"provision_group_inclusion"`
	ProvisionGroupRequired         string                                  `url:"provision_group_required,omitempty" required:"false" json:"provision_group_required,omitempty" path:"provision_group_required"`
	ProvisionAttachmentsPermission *bool                                   `url:"provision_attachments_permission,omitempty" required:"false" json:"provision_attachments_permission,omitempty" path:"provision_attachments_permission"`
	ProvisionDavPermission         *bool                                   `url:"provision_dav_permission,omitempty" required:"false" json:"provision_dav_permission,omitempty" path:"provision_dav_permission"`
	ProvisionFtpPermission         *bool                                   `url:"provision_ftp_permission,omitempty" required:"false" json:"provision_ftp_permission,omitempty" path:"provision_ftp_permission"`
	ProvisionSftpPermission        *bool                                   `url:"provision_sftp_permission,omitempty" required:"false" json:"provision_sftp_permission,omitempty" path:"provision_sftp_permission"`
	ProvisionEmailSignupGroups     string                                  `url:"provision_email_signup_groups,omitempty" required:"false" json:"provision_email_signup_groups,omitempty" path:"provision_email_signup_groups"`
	ProvisionSiteAdminGroups       string                                  `url:"provision_site_admin_groups,omitempty" required:"false" json:"provision_site_admin_groups,omitempty" path:"provision_site_admin_groups"`
	ProvisionGroupAdminGroups      string                                  `url:"provision_group_admin_groups,omitempty" required:"false" json:"provision_group_admin_groups,omitempty" path:"provision_group_admin_groups"`
	ProvisionTimeZone              string                                  `url:"provision_time_zone,omitempty" required:"false" json:"provision_time_zone,omitempty" path:"provision_time_zone"`
	ProvisionCompany               string                                  `url:"provision_company,omitempty" required:"false" json:"provision_company,omitempty" path:"provision_company"`
	Label                          string                                  `url:"label,omitempty" required:"false" json:"label,omitempty" path:"label"`
	LogoFile                       io.Writer                               `url:"logo_file,omitempty" required:"false" json:"logo_file,omitempty" path:"logo_file"`
	LogoDelete                     *bool                                   `url:"logo_delete,omitempty" required:"false" json:"logo_delete,omitempty" path:"logo_delete"`
	LdapBaseDn                     string                                  `url:"ldap_base_dn,omitempty" required:"false" json:"ldap_base_dn,omitempty" path:"ldap_base_dn"`
	LdapDomain                     string                                  `url:"ldap_domain,omitempty" required:"false" json:"ldap_domain,omitempty" path:"ldap_domain"`
	LdapHost                       string                                  `url:"ldap_host,omitempty" required:"false" json:"ldap_host,omitempty" path:"ldap_host"`
	LdapHost2                      string                                  `url:"ldap_host_2,omitempty" required:"false" json:"ldap_host_2,omitempty" path:"ldap_host_2"`
	LdapHost3                      string                                  `url:"ldap_host_3,omitempty" required:"false" json:"ldap_host_3,omitempty" path:"ldap_host_3"`
	LdapPort                       int64                                   `url:"ldap_port,omitempty" required:"false" json:"ldap_port,omitempty" path:"ldap_port"`
	LdapSecure                     *bool                                   `url:"ldap_secure,omitempty" required:"false" json:"ldap_secure,omitempty" path:"ldap_secure"`
	LdapUsername                   string                                  `url:"ldap_username,omitempty" required:"false" json:"ldap_username,omitempty" path:"ldap_username"`
	LdapPassword                   string                                  `url:"ldap_password,omitempty" required:"false" json:"ldap_password,omitempty" path:"ldap_password"`
	LdapUsernameField              SsoStrategyLdapUsernameFieldEnum        `url:"ldap_username_field,omitempty" required:"false" json:"ldap_username_field,omitempty" path:"ldap_username_field"`
	Enabled                        *bool                                   `url:"enabled,omitempty" required:"false" json:"enabled,omitempty" path:"enabled"`
}

type SsoStrategyDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

func (s *SsoStrategy) UnmarshalJSON(data []byte) error {
	type ssoStrategy SsoStrategy
	var v ssoStrategy
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = SsoStrategy(v)
	return nil
}

func (s *SsoStrategyCollection) UnmarshalJSON(data []byte) error {
	type ssoStrategys SsoStrategyCollection
	var v ssoStrategys
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SsoStrategyCollection(v)
	return nil
}

func (s *SsoStrategyCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
