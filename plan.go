package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Plan struct {
	Id                        int64  `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	ActivationCost            string `json:"activation_cost,omitempty" path:"activation_cost,omitempty" url:"activation_cost,omitempty"`
	AddonDescription          string `json:"addon_description,omitempty" path:"addon_description,omitempty" url:"addon_description,omitempty"`
	Annually                  string `json:"annually,omitempty" path:"annually,omitempty" url:"annually,omitempty"`
	AnnuallyAddon             string `json:"annually_addon,omitempty" path:"annually_addon,omitempty" url:"annually_addon,omitempty"`
	ChildSites                int64  `json:"child_sites,omitempty" path:"child_sites,omitempty" url:"child_sites,omitempty"`
	Currency                  string `json:"currency,omitempty" path:"currency,omitempty" url:"currency,omitempty"`
	DedicatedIp               *bool  `json:"dedicated_ip,omitempty" path:"dedicated_ip,omitempty" url:"dedicated_ip,omitempty"`
	DedicatedIps              int64  `json:"dedicated_ips,omitempty" path:"dedicated_ips,omitempty" url:"dedicated_ips,omitempty"`
	DomainCount               int64  `json:"domain_count,omitempty" path:"domain_count,omitempty" url:"domain_count,omitempty"`
	FeatureBundleEca          *bool  `json:"feature_bundle_eca,omitempty" path:"feature_bundle_eca,omitempty" url:"feature_bundle_eca,omitempty"`
	FeatureBundlePower        *bool  `json:"feature_bundle_power,omitempty" path:"feature_bundle_power,omitempty" url:"feature_bundle_power,omitempty"`
	FeatureBundlePremier      *bool  `json:"feature_bundle_premier,omitempty" path:"feature_bundle_premier,omitempty" url:"feature_bundle_premier,omitempty"`
	FeatureBundleStarter      *bool  `json:"feature_bundle_starter,omitempty" path:"feature_bundle_starter,omitempty" url:"feature_bundle_starter,omitempty"`
	Monthly                   string `json:"monthly,omitempty" path:"monthly,omitempty" url:"monthly,omitempty"`
	MonthlyAddon              string `json:"monthly_addon,omitempty" path:"monthly_addon,omitempty" url:"monthly_addon,omitempty"`
	Name                      string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	PreviewPageLimit          int64  `json:"preview_page_limit,omitempty" path:"preview_page_limit,omitempty" url:"preview_page_limit,omitempty"`
	RegionsIncluded           int64  `json:"regions_included,omitempty" path:"regions_included,omitempty" url:"regions_included,omitempty"`
	RemoteSyncInterval        int64  `json:"remote_sync_interval,omitempty" path:"remote_sync_interval,omitempty" url:"remote_sync_interval,omitempty"`
	StagingSites              int64  `json:"staging_sites,omitempty" path:"staging_sites,omitempty" url:"staging_sites,omitempty"`
	UserCost                  string `json:"user_cost,omitempty" path:"user_cost,omitempty" url:"user_cost,omitempty"`
	UsageCost                 string `json:"usage_cost,omitempty" path:"usage_cost,omitempty" url:"usage_cost,omitempty"`
	UsageIncluded             string `json:"usage_included,omitempty" path:"usage_included,omitempty" url:"usage_included,omitempty"`
	Users                     int64  `json:"users,omitempty" path:"users,omitempty" url:"users,omitempty"`
	AdvancedBehaviors         *bool  `json:"advanced_behaviors,omitempty" path:"advanced_behaviors,omitempty" url:"advanced_behaviors,omitempty"`
	AuthGoogle                *bool  `json:"auth_google,omitempty" path:"auth_google,omitempty" url:"auth_google,omitempty"`
	AuthOauth                 *bool  `json:"auth_oauth,omitempty" path:"auth_oauth,omitempty" url:"auth_oauth,omitempty"`
	AuthOauthCustom           *bool  `json:"auth_oauth_custom,omitempty" path:"auth_oauth_custom,omitempty" url:"auth_oauth_custom,omitempty"`
	AuthUserCount             int64  `json:"auth_user_count,omitempty" path:"auth_user_count,omitempty" url:"auth_user_count,omitempty"`
	Automations               *bool  `json:"automations,omitempty" path:"automations,omitempty" url:"automations,omitempty"`
	CustomNamespace           *bool  `json:"custom_namespace,omitempty" path:"custom_namespace,omitempty" url:"custom_namespace,omitempty"`
	CustomSmtp                *bool  `json:"custom_smtp,omitempty" path:"custom_smtp,omitempty" url:"custom_smtp,omitempty"`
	Domain                    *bool  `json:"domain,omitempty" path:"domain,omitempty" url:"domain,omitempty"`
	ExtendedFolderPermissions *bool  `json:"extended_folder_permissions,omitempty" path:"extended_folder_permissions,omitempty" url:"extended_folder_permissions,omitempty"`
	FtpSftpWebdav             *bool  `json:"ftp_sftp_webdav,omitempty" path:"ftp_sftp_webdav,omitempty" url:"ftp_sftp_webdav,omitempty"`
	Gpg                       *bool  `json:"gpg,omitempty" path:"gpg,omitempty" url:"gpg,omitempty"`
	GroupAdminsEnabled        *bool  `json:"group_admins_enabled,omitempty" path:"group_admins_enabled,omitempty" url:"group_admins_enabled,omitempty"`
	GroupNotifications        *bool  `json:"group_notifications,omitempty" path:"group_notifications,omitempty" url:"group_notifications,omitempty"`
	Hipaa                     *bool  `json:"hipaa,omitempty" path:"hipaa,omitempty" url:"hipaa,omitempty"`
	Ldap                      *bool  `json:"ldap,omitempty" path:"ldap,omitempty" url:"ldap,omitempty"`
	LegalFlexibility          *bool  `json:"legal_flexibility,omitempty" path:"legal_flexibility,omitempty" url:"legal_flexibility,omitempty"`
	RemoteSyncFtp             *bool  `json:"remote_sync_ftp,omitempty" path:"remote_sync_ftp,omitempty" url:"remote_sync_ftp,omitempty"`
	Require2fa                *bool  `json:"require_2fa,omitempty" path:"require_2fa,omitempty" url:"require_2fa,omitempty"`
	SecurityOptOut            *bool  `json:"security_opt_out,omitempty" path:"security_opt_out,omitempty" url:"security_opt_out,omitempty"`
	WatermarkImages           *bool  `json:"watermark_images,omitempty" path:"watermark_images,omitempty" url:"watermark_images,omitempty"`
	Webhooks                  *bool  `json:"webhooks,omitempty" path:"webhooks,omitempty" url:"webhooks,omitempty"`
	WebhooksSns               *bool  `json:"webhooks_sns,omitempty" path:"webhooks_sns,omitempty" url:"webhooks_sns,omitempty"`
}

func (p Plan) Identifier() interface{} {
	return p.Id
}

type PlanCollection []Plan

type PlanListParams struct {
	Action   string `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	Currency string `url:"currency,omitempty" required:"false" json:"currency,omitempty" path:"currency"`
	ListParams
}

func (p *Plan) UnmarshalJSON(data []byte) error {
	type plan Plan
	var v plan
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*p = Plan(v)
	return nil
}

func (p *PlanCollection) UnmarshalJSON(data []byte) error {
	type plans PlanCollection
	var v plans
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*p = PlanCollection(v)
	return nil
}

func (p *PlanCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
