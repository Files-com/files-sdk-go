package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type SsoStrategy struct {
	Protocol                       string `json:"protocol,omitempty" path:"protocol"`
	Provider                       string `json:"provider,omitempty" path:"provider"`
	Label                          string `json:"label,omitempty" path:"label"`
	LogoUrl                        string `json:"logo_url,omitempty" path:"logo_url"`
	Id                             int64  `json:"id,omitempty" path:"id"`
	SamlProviderCertFingerprint    string `json:"saml_provider_cert_fingerprint,omitempty" path:"saml_provider_cert_fingerprint"`
	SamlProviderIssuerUrl          string `json:"saml_provider_issuer_url,omitempty" path:"saml_provider_issuer_url"`
	SamlProviderMetadataContent    string `json:"saml_provider_metadata_content,omitempty" path:"saml_provider_metadata_content"`
	SamlProviderMetadataUrl        string `json:"saml_provider_metadata_url,omitempty" path:"saml_provider_metadata_url"`
	SamlProviderSloTargetUrl       string `json:"saml_provider_slo_target_url,omitempty" path:"saml_provider_slo_target_url"`
	SamlProviderSsoTargetUrl       string `json:"saml_provider_sso_target_url,omitempty" path:"saml_provider_sso_target_url"`
	ScimAuthenticationMethod       string `json:"scim_authentication_method,omitempty" path:"scim_authentication_method"`
	ScimUsername                   string `json:"scim_username,omitempty" path:"scim_username"`
	ScimOauthAccessToken           string `json:"scim_oauth_access_token,omitempty" path:"scim_oauth_access_token"`
	ScimOauthAccessTokenExpiresAt  string `json:"scim_oauth_access_token_expires_at,omitempty" path:"scim_oauth_access_token_expires_at"`
	Subdomain                      string `json:"subdomain,omitempty" path:"subdomain"`
	ProvisionUsers                 *bool  `json:"provision_users,omitempty" path:"provision_users"`
	ProvisionGroups                *bool  `json:"provision_groups,omitempty" path:"provision_groups"`
	DeprovisionUsers               *bool  `json:"deprovision_users,omitempty" path:"deprovision_users"`
	DeprovisionGroups              *bool  `json:"deprovision_groups,omitempty" path:"deprovision_groups"`
	DeprovisionBehavior            string `json:"deprovision_behavior,omitempty" path:"deprovision_behavior"`
	ProvisionGroupDefault          string `json:"provision_group_default,omitempty" path:"provision_group_default"`
	ProvisionGroupExclusion        string `json:"provision_group_exclusion,omitempty" path:"provision_group_exclusion"`
	ProvisionGroupInclusion        string `json:"provision_group_inclusion,omitempty" path:"provision_group_inclusion"`
	ProvisionGroupRequired         string `json:"provision_group_required,omitempty" path:"provision_group_required"`
	ProvisionEmailSignupGroups     string `json:"provision_email_signup_groups,omitempty" path:"provision_email_signup_groups"`
	ProvisionSiteAdminGroups       string `json:"provision_site_admin_groups,omitempty" path:"provision_site_admin_groups"`
	ProvisionAttachmentsPermission *bool  `json:"provision_attachments_permission,omitempty" path:"provision_attachments_permission"`
	ProvisionDavPermission         *bool  `json:"provision_dav_permission,omitempty" path:"provision_dav_permission"`
	ProvisionFtpPermission         *bool  `json:"provision_ftp_permission,omitempty" path:"provision_ftp_permission"`
	ProvisionSftpPermission        *bool  `json:"provision_sftp_permission,omitempty" path:"provision_sftp_permission"`
	ProvisionTimeZone              string `json:"provision_time_zone,omitempty" path:"provision_time_zone"`
	ProvisionCompany               string `json:"provision_company,omitempty" path:"provision_company"`
	LdapBaseDn                     string `json:"ldap_base_dn,omitempty" path:"ldap_base_dn"`
	LdapDomain                     string `json:"ldap_domain,omitempty" path:"ldap_domain"`
	Enabled                        *bool  `json:"enabled,omitempty" path:"enabled"`
	LdapHost                       string `json:"ldap_host,omitempty" path:"ldap_host"`
	LdapHost2                      string `json:"ldap_host_2,omitempty" path:"ldap_host_2"`
	LdapHost3                      string `json:"ldap_host_3,omitempty" path:"ldap_host_3"`
	LdapPort                       int64  `json:"ldap_port,omitempty" path:"ldap_port"`
	LdapSecure                     *bool  `json:"ldap_secure,omitempty" path:"ldap_secure"`
	LdapUsername                   string `json:"ldap_username,omitempty" path:"ldap_username"`
	LdapUsernameField              string `json:"ldap_username_field,omitempty" path:"ldap_username_field"`
}

func (s SsoStrategy) Identifier() interface{} {
	return s.Id
}

type SsoStrategyCollection []SsoStrategy

type SsoStrategyListParams struct {
	ListParams
}

type SsoStrategyFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

// Synchronize provisioning data with the SSO remote server
type SsoStrategySyncParams struct {
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
