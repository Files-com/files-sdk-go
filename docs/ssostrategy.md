# SsoStrategy

`files_sdk.SsoStrategy` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/ssostrategy`, named `sso_strategy`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	sso_strategy "github.com/Files-com/files-sdk-go/v3/ssostrategy"
)
```

## Operations

### List

`GET /sso_strategies`

List SSO Strategies

```go
func (c *Client) List(params files_sdk.SsoStrategyListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.SsoStrategyListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.SsoStrategy` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.SsoStrategyListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are . |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := sso_strategy.List(files_sdk.SsoStrategyListParams{})
if err != nil {
	return err
}
for item, err := range it.All() {
	if err != nil {
		return err
	}
	fmt.Println(item)
}
```

### Find

`GET /sso_strategies/{id}`

Show SSO Strategy

```go
func (c *Client) Find(params files_sdk.SsoStrategyFindParams, opts ...files_sdk.RequestResponseOption) (ssoStrategy files_sdk.SsoStrategy, err error)
func Find(params files_sdk.SsoStrategyFindParams, opts ...files_sdk.RequestResponseOption) (ssoStrategy files_sdk.SsoStrategy, err error)
```

Parameters, `files_sdk.SsoStrategyFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Sso Strategy ID. |

Example:

```go
result, err := sso_strategy.Find(files_sdk.SsoStrategyFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Sync

`POST /sso_strategies/{id}/sync`

Synchronize provisioning data with the SSO remote server

```go
func (c *Client) Sync(params files_sdk.SsoStrategySyncParams, opts ...files_sdk.RequestResponseOption) (err error)
func Sync(params files_sdk.SsoStrategySyncParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.SsoStrategySyncParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Sso Strategy ID. |

Example:

```go
err := sso_strategy.Sync(files_sdk.SsoStrategySyncParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.SsoStrategy` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Protocol` | `string` | `protocol` | SSO Protocol |
| `Provider` | `string` | `provider` | Provider name |
| `Label` | `string` | `label` | Custom label for the SSO provider on the login page. |
| `LogoUrl` | `string` | `logo_url` | URL holding a custom logo for the SSO provider on the login page. |
| `Id` | `int64` | `id` | ID |
| `Enabled` | `*bool` | `enabled` | Is strategy enabled? This may become automatically set to `false` after a high number and duration of failures. |
| `UserCount` | `int64` | `user_count` | Count of users with this SSO Strategy |
| `SamlProviderCertFingerprint` | `string` | `saml_provider_cert_fingerprint` | Identity provider sha256 cert fingerprint if saml_provider_metadata_url is not available. |
| `SamlProviderIssuerUrl` | `string` | `saml_provider_issuer_url` | Identity provider issuer url |
| `SamlProviderMetadataContent` | `string` | `saml_provider_metadata_content` | Custom identity provider metadata |
| `SamlProviderMetadataUrl` | `string` | `saml_provider_metadata_url` | Metadata URL for the SAML identity provider |
| `SamlProviderSloTargetUrl` | `string` | `saml_provider_slo_target_url` | Identity provider SLO endpoint |
| `SamlProviderSsoTargetUrl` | `string` | `saml_provider_sso_target_url` | Identity provider SSO endpoint if saml_provider_metadata_url is not available. |
| `ScimAuthenticationMethod` | `string` | `scim_authentication_method` | SCIM authentication type. |
| `ScimUsername` | `string` | `scim_username` | SCIM username. |
| `ScimOauthAccessToken` | `string` | `scim_oauth_access_token` | SCIM OAuth Access Token. Only returned when the token is first issued or reset. |
| `ScimOauthAccessTokenExpiresAt` | `string` | `scim_oauth_access_token_expires_at` | SCIM OAuth Access Token Expiration Time. |
| `Subdomain` | `string` | `subdomain` | Subdomain or domain name for your auth provider. Example: `https://[subdomain].okta.com/` |
| `ProvisionUsers` | `*bool` | `provision_users` | Auto-provision users? |
| `ProvisionGroups` | `*bool` | `provision_groups` | Auto-provision group membership based on group memberships on the SSO side? |
| `DeprovisionUsers` | `*bool` | `deprovision_users` | Auto-deprovision users? |
| `DeprovisionGroups` | `*bool` | `deprovision_groups` | Auto-deprovision group membership based on group memberships on the SSO side? |
| `DeprovisionBehavior` | `string` | `deprovision_behavior` | Method used for deprovisioning users. |
| `ProvisionGroupDefault` | `string` | `provision_group_default` | Comma-separated list of group names for groups to automatically add all auto-provisioned users to. |
| `ProvisionGroupExclusion` | `string` | `provision_group_exclusion` | Comma-separated list of group names for groups (with optional wildcards) that will be excluded from auto-provisioning. |
| `ProvisionGroupInclusion` | `string` | `provision_group_inclusion` | Comma-separated list of group names for groups (with optional wildcards) that will be auto-provisioned. |
| `ProvisionGroupRequired` | `string` | `provision_group_required` | Comma or newline separated list of group names (with optional wildcards) to require membership for user provisioning. |
| `ProvisionEmailSignupGroups` | `string` | `provision_email_signup_groups` | Comma-separated list of group names whose members will be created with email_signup authentication. |
| `ProvisionReadonlySiteAdminGroups` | `string` | `provision_readonly_site_admin_groups` | Comma-separated list of group names used to grant or remove Read-Only Site Admin access for users on this SSO provider. When this list is nonempty and a directory sync changes a user's group membership, membership in at least one listed group grants this access; membership in none removes it, including access granted manually. Directory syncs include SCIM requests from any provider on the site, LDAP syncs, and sign-ins that set the user's groups from the identity provider. Membership changes made any other way, including through the web app, API, SDKs, or CLI, do not affect administrator access. |
| `ProvisionSiteAdminGroups` | `string` | `provision_site_admin_groups` | Comma-separated list of group names used to grant or remove Site Admin access for users on this SSO provider. When this list is nonempty and a directory sync changes a user's group membership, membership in at least one listed group grants this access; membership in none removes it, including access granted manually. Directory syncs include SCIM requests from any provider on the site, LDAP syncs, and sign-ins that set the user's groups from the identity provider. Membership changes made any other way, including through the web app, API, SDKs, or CLI, do not affect administrator access. |
| `ProvisionGroupAdminGroups` | `string` | `provision_group_admin_groups` | Comma-separated list of group names whose members will be provisioned as Group Admins. |
| `ProvisionAttachmentsPermission` | `*bool` | `provision_attachments_permission` |  |
| `ProvisionDavPermission` | `*bool` | `provision_dav_permission` | Auto-provisioned users get WebDAV permission? |
| `ProvisionFtpPermission` | `*bool` | `provision_ftp_permission` | Auto-provisioned users get FTP permission? |
| `ProvisionSftpPermission` | `*bool` | `provision_sftp_permission` | Auto-provisioned users get SFTP permission? |
| `ProvisionTimeZone` | `string` | `provision_time_zone` | Default time zone for auto provisioned users. |
| `ProvisionCompany` | `string` | `provision_company` | Default company for auto provisioned users. |
| `ProvisionRequire2fa` | `string` | `provision_require_2fa` | 2FA required setting for auto provisioned users. `use_system_setting` uses the site-wide setting, including SSO exemptions. `always_require` and `never_require` override the site-wide setting when user-level overrides are allowed. |
| `ProvisionFilesystemLayout` | `string` | `provision_filesystem_layout` | File System layout to use for auto provisioned users. |
| `ProviderIdentifier` | `string` | `provider_identifier` | URL-friendly, unique identifier for Azure SAML configuration |
| `LdapBaseDn` | `string` | `ldap_base_dn` | Base DN for looking up users in LDAP server |
| `LdapDomain` | `string` | `ldap_domain` | Domain name that will be appended to LDAP usernames |
| `DisplayOnLoginPage` | `*bool` | `display_on_login_page` | Should this strategy be displayed on the login page? |
| `LdapHost` | `string` | `ldap_host` | LDAP host |
| `LdapHost2` | `string` | `ldap_host_2` | LDAP backup host |
| `LdapHost3` | `string` | `ldap_host_3` | LDAP backup host |
| `LdapPort` | `int64` | `ldap_port` | LDAP port |
| `LdapProvisioningEnabled` | `*bool` | `ldap_provisioning_enabled` | Use LDAP server settings for scheduled provisioning while using this SSO provider for authentication? |
| `LdapSecure` | `*bool` | `ldap_secure` | Use secure LDAP? |
| `LdapServerCertificate` | `string` | `ldap_server_certificate` | How to validate the LDAP server certificate. `require_match` validates the certificate chain and hostname; `allow_any` disables certificate validation. |
| `LdapType` | `string` | `ldap_type` | LDAP server type |
| `LdapUsername` | `string` | `ldap_username` | Username for signing in to LDAP server. |
| `LdapUsernameField` | `string` | `ldap_username_field` | LDAP username field |
