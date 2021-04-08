package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type SsoStrategy struct {
	Protocol                       string `json:"protocol,omitempty"`
	Provider                       string `json:"provider,omitempty"`
	Label                          string `json:"label,omitempty"`
	LogoUrl                        string `json:"logo_url,omitempty"`
	Id                             int64  `json:"id,omitempty"`
	SamlProviderCertFingerprint    string `json:"saml_provider_cert_fingerprint,omitempty"`
	SamlProviderIssuerUrl          string `json:"saml_provider_issuer_url,omitempty"`
	SamlProviderMetadataContent    string `json:"saml_provider_metadata_content,omitempty"`
	SamlProviderMetadataUrl        string `json:"saml_provider_metadata_url,omitempty"`
	SamlProviderSloTargetUrl       string `json:"saml_provider_slo_target_url,omitempty"`
	SamlProviderSsoTargetUrl       string `json:"saml_provider_sso_target_url,omitempty"`
	ScimAuthenticationMethod       string `json:"scim_authentication_method,omitempty"`
	ScimUsername                   string `json:"scim_username,omitempty"`
	ScimOauthAccessToken           string `json:"scim_oauth_access_token,omitempty"`
	ScimOauthAccessTokenExpiresAt  string `json:"scim_oauth_access_token_expires_at,omitempty"`
	Subdomain                      string `json:"subdomain,omitempty"`
	ProvisionUsers                 *bool  `json:"provision_users,omitempty"`
	ProvisionGroups                *bool  `json:"provision_groups,omitempty"`
	DeprovisionUsers               *bool  `json:"deprovision_users,omitempty"`
	DeprovisionGroups              *bool  `json:"deprovision_groups,omitempty"`
	DeprovisionBehavior            string `json:"deprovision_behavior,omitempty"`
	ProvisionGroupDefault          string `json:"provision_group_default,omitempty"`
	ProvisionGroupExclusion        string `json:"provision_group_exclusion,omitempty"`
	ProvisionGroupInclusion        string `json:"provision_group_inclusion,omitempty"`
	ProvisionGroupRequired         string `json:"provision_group_required,omitempty"`
	ProvisionSiteAdminGroups       string `json:"provision_site_admin_groups,omitempty"`
	ProvisionAttachmentsPermission *bool  `json:"provision_attachments_permission,omitempty"`
	ProvisionDavPermission         *bool  `json:"provision_dav_permission,omitempty"`
	ProvisionFtpPermission         *bool  `json:"provision_ftp_permission,omitempty"`
	ProvisionSftpPermission        *bool  `json:"provision_sftp_permission,omitempty"`
	ProvisionTimeZone              string `json:"provision_time_zone,omitempty"`
	ProvisionCompany               string `json:"provision_company,omitempty"`
	LdapBaseDn                     string `json:"ldap_base_dn,omitempty"`
	LdapDomain                     string `json:"ldap_domain,omitempty"`
	Enabled                        *bool  `json:"enabled,omitempty"`
	LdapHost                       string `json:"ldap_host,omitempty"`
	LdapHost2                      string `json:"ldap_host_2,omitempty"`
	LdapHost3                      string `json:"ldap_host_3,omitempty"`
	LdapPort                       int    `json:"ldap_port,omitempty"`
	LdapSecure                     *bool  `json:"ldap_secure,omitempty"`
	LdapUsername                   string `json:"ldap_username,omitempty"`
	LdapUsernameField              string `json:"ldap_username_field,omitempty"`
}

type SsoStrategyCollection []SsoStrategy

type SsoStrategyListParams struct {
	Cursor  string `url:"cursor,omitempty" required:"false"`
	PerPage int    `url:"per_page,omitempty" required:"false"`
	lib.ListParams
}

type SsoStrategyFindParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

func (s *SsoStrategy) UnmarshalJSON(data []byte) error {
	type ssoStrategy SsoStrategy
	var v ssoStrategy
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*s = SsoStrategy(v)
	return nil
}

func (s *SsoStrategyCollection) UnmarshalJSON(data []byte) error {
	type ssoStrategys []SsoStrategy
	var v ssoStrategys
	if err := json.Unmarshal(data, &v); err != nil {
		return err
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
