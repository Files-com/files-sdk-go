package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// SiemHttpDestination is a Files.com API resource.
type SiemHttpDestination struct {
	// SIEM HTTP Destination ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Name for this Destination.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Destination Type.
	DestinationType string `json:"destination_type,omitempty" path:"destination_type,omitempty" url:"destination_type,omitempty"`
	// Destination Url.
	DestinationUrl string `json:"destination_url,omitempty" path:"destination_url,omitempty" url:"destination_url,omitempty"`
	// Applicable only for destination type: file. Destination folder path on Files.com.
	FileDestinationPath string `json:"file_destination_path,omitempty" path:"file_destination_path,omitempty" url:"file_destination_path,omitempty"`
	// Applicable only for destination type: file. Generated file format.
	FileFormat string `json:"file_format,omitempty" path:"file_format,omitempty" url:"file_format,omitempty"`
	// Applicable only for destination type: file. Interval, in minutes, between file deliveries.
	FileIntervalMinutes int64 `json:"file_interval_minutes,omitempty" path:"file_interval_minutes,omitempty" url:"file_interval_minutes,omitempty"`
	// Additional HTTP Headers included in calls to the destination URL.
	AdditionalHeaders interface{} `json:"additional_headers,omitempty" path:"additional_headers,omitempty" url:"additional_headers,omitempty"`
	// Whether this SIEM HTTP Destination is currently being sent to or not.
	SendingActive *bool `json:"sending_active,omitempty" path:"sending_active,omitempty" url:"sending_active,omitempty"`
	// Applicable only for destination type: generic. Indicates the type of HTTP body. Can be json_newline or json_array. json_newline is multiple log entries as JSON separated by newlines. json_array is a single JSON array containing multiple log entries as JSON.
	GenericPayloadType string `json:"generic_payload_type,omitempty" path:"generic_payload_type,omitempty" url:"generic_payload_type,omitempty"`
	// Applicable only for destination types: splunk, splunk_compatible. Authentication token for the destination.
	SplunkTokenMasked string `json:"splunk_token_masked,omitempty" path:"splunk_token_masked,omitempty" url:"splunk_token_masked,omitempty"`
	// Applicable only for destination type: crowdstrike. Authentication token provided by Crowdstrike.
	CrowdstrikeTokenMasked string `json:"crowdstrike_token_masked,omitempty" path:"crowdstrike_token_masked,omitempty" url:"crowdstrike_token_masked,omitempty"`
	// Applicable only for destination types: azure, azure_legacy. Immutable ID of the Data Collection Rule.
	AzureDcrImmutableId string `json:"azure_dcr_immutable_id,omitempty" path:"azure_dcr_immutable_id,omitempty" url:"azure_dcr_immutable_id,omitempty"`
	// Applicable only for destination type: azure. Name of the stream in the DCR that represents the destination table.
	AzureStreamName string `json:"azure_stream_name,omitempty" path:"azure_stream_name,omitempty" url:"azure_stream_name,omitempty"`
	// Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Tenant ID.
	AzureOauthClientCredentialsTenantId string `json:"azure_oauth_client_credentials_tenant_id,omitempty" path:"azure_oauth_client_credentials_tenant_id,omitempty" url:"azure_oauth_client_credentials_tenant_id,omitempty"`
	// Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Client ID.
	AzureOauthClientCredentialsClientId string `json:"azure_oauth_client_credentials_client_id,omitempty" path:"azure_oauth_client_credentials_client_id,omitempty" url:"azure_oauth_client_credentials_client_id,omitempty"`
	// Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Client Secret.
	AzureOauthClientCredentialsClientSecretMasked string `json:"azure_oauth_client_credentials_client_secret_masked,omitempty" path:"azure_oauth_client_credentials_client_secret_masked,omitempty" url:"azure_oauth_client_credentials_client_secret_masked,omitempty"`
	// Applicable only for destination type: qradar. Basic auth username provided by QRadar.
	QradarUsername string `json:"qradar_username,omitempty" path:"qradar_username,omitempty" url:"qradar_username,omitempty"`
	// Applicable only for destination type: qradar. Basic auth password provided by QRadar.
	QradarPasswordMasked string `json:"qradar_password_masked,omitempty" path:"qradar_password_masked,omitempty" url:"qradar_password_masked,omitempty"`
	// Applicable only for destination type: solar_winds. Authentication token provided by Solar Winds.
	SolarWindsTokenMasked string `json:"solar_winds_token_masked,omitempty" path:"solar_winds_token_masked,omitempty" url:"solar_winds_token_masked,omitempty"`
	// Applicable only for destination type: new_relic. API key provided by New Relic.
	NewRelicApiKeyMasked string `json:"new_relic_api_key_masked,omitempty" path:"new_relic_api_key_masked,omitempty" url:"new_relic_api_key_masked,omitempty"`
	// Applicable only for destination type: datadog. API key provided by Datadog.
	DatadogApiKeyMasked string `json:"datadog_api_key_masked,omitempty" path:"datadog_api_key_masked,omitempty" url:"datadog_api_key_masked,omitempty"`
	// Whether or not sending is enabled for action logs.
	ActionSendEnabled *bool `json:"action_send_enabled,omitempty" path:"action_send_enabled,omitempty" url:"action_send_enabled,omitempty"`
	// Number of log entries sent for the lifetime of this destination.
	ActionEntriesSent int64 `json:"action_entries_sent,omitempty" path:"action_entries_sent,omitempty" url:"action_entries_sent,omitempty"`
	// Whether or not sending is enabled for sftp_action logs.
	SftpActionSendEnabled *bool `json:"sftp_action_send_enabled,omitempty" path:"sftp_action_send_enabled,omitempty" url:"sftp_action_send_enabled,omitempty"`
	// Number of log entries sent for the lifetime of this destination.
	SftpActionEntriesSent int64 `json:"sftp_action_entries_sent,omitempty" path:"sftp_action_entries_sent,omitempty" url:"sftp_action_entries_sent,omitempty"`
	// Whether or not sending is enabled for ftp_action logs.
	FtpActionSendEnabled *bool `json:"ftp_action_send_enabled,omitempty" path:"ftp_action_send_enabled,omitempty" url:"ftp_action_send_enabled,omitempty"`
	// Number of log entries sent for the lifetime of this destination.
	FtpActionEntriesSent int64 `json:"ftp_action_entries_sent,omitempty" path:"ftp_action_entries_sent,omitempty" url:"ftp_action_entries_sent,omitempty"`
	// Whether or not sending is enabled for web_dav_action logs.
	WebDavActionSendEnabled *bool `json:"web_dav_action_send_enabled,omitempty" path:"web_dav_action_send_enabled,omitempty" url:"web_dav_action_send_enabled,omitempty"`
	// Number of log entries sent for the lifetime of this destination.
	WebDavActionEntriesSent int64 `json:"web_dav_action_entries_sent,omitempty" path:"web_dav_action_entries_sent,omitempty" url:"web_dav_action_entries_sent,omitempty"`
	// Whether or not sending is enabled for sync logs.
	SyncSendEnabled *bool `json:"sync_send_enabled,omitempty" path:"sync_send_enabled,omitempty" url:"sync_send_enabled,omitempty"`
	// Number of log entries sent for the lifetime of this destination.
	SyncEntriesSent int64 `json:"sync_entries_sent,omitempty" path:"sync_entries_sent,omitempty" url:"sync_entries_sent,omitempty"`
	// Whether or not sending is enabled for outbound_connection logs.
	OutboundConnectionSendEnabled *bool `json:"outbound_connection_send_enabled,omitempty" path:"outbound_connection_send_enabled,omitempty" url:"outbound_connection_send_enabled,omitempty"`
	// Number of log entries sent for the lifetime of this destination.
	OutboundConnectionEntriesSent int64 `json:"outbound_connection_entries_sent,omitempty" path:"outbound_connection_entries_sent,omitempty" url:"outbound_connection_entries_sent,omitempty"`
	// Whether or not sending is enabled for automation logs.
	AutomationSendEnabled *bool `json:"automation_send_enabled,omitempty" path:"automation_send_enabled,omitempty" url:"automation_send_enabled,omitempty"`
	// Number of log entries sent for the lifetime of this destination.
	AutomationEntriesSent int64 `json:"automation_entries_sent,omitempty" path:"automation_entries_sent,omitempty" url:"automation_entries_sent,omitempty"`
	// Whether or not sending is enabled for api_request logs.
	ApiRequestSendEnabled *bool `json:"api_request_send_enabled,omitempty" path:"api_request_send_enabled,omitempty" url:"api_request_send_enabled,omitempty"`
	// Number of log entries sent for the lifetime of this destination.
	ApiRequestEntriesSent int64 `json:"api_request_entries_sent,omitempty" path:"api_request_entries_sent,omitempty" url:"api_request_entries_sent,omitempty"`
	// Whether or not sending is enabled for public_hosting_request logs.
	PublicHostingRequestSendEnabled *bool `json:"public_hosting_request_send_enabled,omitempty" path:"public_hosting_request_send_enabled,omitempty" url:"public_hosting_request_send_enabled,omitempty"`
	// Number of log entries sent for the lifetime of this destination.
	PublicHostingRequestEntriesSent int64 `json:"public_hosting_request_entries_sent,omitempty" path:"public_hosting_request_entries_sent,omitempty" url:"public_hosting_request_entries_sent,omitempty"`
	// Whether or not sending is enabled for email logs.
	EmailSendEnabled *bool `json:"email_send_enabled,omitempty" path:"email_send_enabled,omitempty" url:"email_send_enabled,omitempty"`
	// Number of log entries sent for the lifetime of this destination.
	EmailEntriesSent int64 `json:"email_entries_sent,omitempty" path:"email_entries_sent,omitempty" url:"email_entries_sent,omitempty"`
	// Whether or not sending is enabled for exavault_api_request logs.
	ExavaultApiRequestSendEnabled *bool `json:"exavault_api_request_send_enabled,omitempty" path:"exavault_api_request_send_enabled,omitempty" url:"exavault_api_request_send_enabled,omitempty"`
	// Number of log entries sent for the lifetime of this destination.
	ExavaultApiRequestEntriesSent int64 `json:"exavault_api_request_entries_sent,omitempty" path:"exavault_api_request_entries_sent,omitempty" url:"exavault_api_request_entries_sent,omitempty"`
	// Whether or not sending is enabled for settings_change logs.
	SettingsChangeSendEnabled *bool `json:"settings_change_send_enabled,omitempty" path:"settings_change_send_enabled,omitempty" url:"settings_change_send_enabled,omitempty"`
	// Number of log entries sent for the lifetime of this destination.
	SettingsChangeEntriesSent int64 `json:"settings_change_entries_sent,omitempty" path:"settings_change_entries_sent,omitempty" url:"settings_change_entries_sent,omitempty"`
	// Type of URL that was last called. Can be `destination_url` or `azure_oauth_client_credentials_url`.
	LastHttpCallTargetType string `json:"last_http_call_target_type,omitempty" path:"last_http_call_target_type,omitempty" url:"last_http_call_target_type,omitempty"`
	// Was the last HTTP call made successful?
	LastHttpCallSuccess *bool `json:"last_http_call_success,omitempty" path:"last_http_call_success,omitempty" url:"last_http_call_success,omitempty"`
	// Last HTTP Call Response Code.
	LastHttpCallResponseCode int64 `json:"last_http_call_response_code,omitempty" path:"last_http_call_response_code,omitempty" url:"last_http_call_response_code,omitempty"`
	// Last HTTP Call Response Body. Large responses are truncated.
	LastHttpCallResponseBody string `json:"last_http_call_response_body,omitempty" path:"last_http_call_response_body,omitempty" url:"last_http_call_response_body,omitempty"`
	// Last HTTP Call Error Message if applicable.
	LastHttpCallErrorMessage string `json:"last_http_call_error_message,omitempty" path:"last_http_call_error_message,omitempty" url:"last_http_call_error_message,omitempty"`
	// Time of Last HTTP Call.
	LastHttpCallTime string `json:"last_http_call_time,omitempty" path:"last_http_call_time,omitempty" url:"last_http_call_time,omitempty"`
	// Duration of the last HTTP Call in milliseconds.
	LastHttpCallDurationMs int64 `json:"last_http_call_duration_ms,omitempty" path:"last_http_call_duration_ms,omitempty" url:"last_http_call_duration_ms,omitempty"`
	// Time of Most Recent Successful HTTP Call.
	MostRecentHttpCallSuccessTime string `json:"most_recent_http_call_success_time,omitempty" path:"most_recent_http_call_success_time,omitempty" url:"most_recent_http_call_success_time,omitempty"`
	// Connection Test Entry.
	ConnectionTestEntry string `json:"connection_test_entry,omitempty" path:"connection_test_entry,omitempty" url:"connection_test_entry,omitempty"`
	// Applicable only for destination types: splunk, splunk_compatible. Authentication token for the destination.
	SplunkToken string `json:"splunk_token,omitempty" path:"splunk_token,omitempty" url:"splunk_token,omitempty"`
	// Applicable only for destination type: crowdstrike. Authentication token provided by Crowdstrike.
	CrowdstrikeToken string `json:"crowdstrike_token,omitempty" path:"crowdstrike_token,omitempty" url:"crowdstrike_token,omitempty"`
	// Applicable only for destination type: azure. Client Credentials OAuth Client Secret.
	AzureOauthClientCredentialsClientSecret string `json:"azure_oauth_client_credentials_client_secret,omitempty" path:"azure_oauth_client_credentials_client_secret,omitempty" url:"azure_oauth_client_credentials_client_secret,omitempty"`
	// Applicable only for destination type: qradar. Basic auth password provided by QRadar.
	QradarPassword string `json:"qradar_password,omitempty" path:"qradar_password,omitempty" url:"qradar_password,omitempty"`
	// Applicable only for destination type: solar_winds. Authentication token provided by Solar Winds.
	SolarWindsToken string `json:"solar_winds_token,omitempty" path:"solar_winds_token,omitempty" url:"solar_winds_token,omitempty"`
	// Applicable only for destination type: new_relic. API key provided by New Relic.
	NewRelicApiKey string `json:"new_relic_api_key,omitempty" path:"new_relic_api_key,omitempty" url:"new_relic_api_key,omitempty"`
	// Applicable only for destination type: datadog. API key provided by Datadog.
	DatadogApiKey string `json:"datadog_api_key,omitempty" path:"datadog_api_key,omitempty" url:"datadog_api_key,omitempty"`
}

// Identifier returns the resource ID.
func (s SiemHttpDestination) Identifier() interface{} {
	return s.Id
}

// SiemHttpDestinationCollection is a list of SiemHttpDestination resources.
type SiemHttpDestinationCollection []SiemHttpDestination

// SiemHttpDestinationGenericPayloadTypeEnum is a string value for generic_payload_type.
// Enum lists the values documented by the API.
type SiemHttpDestinationGenericPayloadTypeEnum string

// String returns the API parameter value.
func (u SiemHttpDestinationGenericPayloadTypeEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u SiemHttpDestinationGenericPayloadTypeEnum) Enum() map[string]SiemHttpDestinationGenericPayloadTypeEnum {
	return map[string]SiemHttpDestinationGenericPayloadTypeEnum{
		"json_newline": SiemHttpDestinationGenericPayloadTypeEnum("json_newline"),
		"json_array":   SiemHttpDestinationGenericPayloadTypeEnum("json_array"),
	}
}

// SiemHttpDestinationFileFormatEnum is a string value for file_format.
// Enum lists the values documented by the API.
type SiemHttpDestinationFileFormatEnum string

// String returns the API parameter value.
func (u SiemHttpDestinationFileFormatEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u SiemHttpDestinationFileFormatEnum) Enum() map[string]SiemHttpDestinationFileFormatEnum {
	return map[string]SiemHttpDestinationFileFormatEnum{
		"json": SiemHttpDestinationFileFormatEnum("json"),
		"csv":  SiemHttpDestinationFileFormatEnum("csv"),
	}
}

// SiemHttpDestinationDestinationTypeEnum is a string value for destination_type.
// Enum lists the values documented by the API.
type SiemHttpDestinationDestinationTypeEnum string

// String returns the API parameter value.
func (u SiemHttpDestinationDestinationTypeEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u SiemHttpDestinationDestinationTypeEnum) Enum() map[string]SiemHttpDestinationDestinationTypeEnum {
	return map[string]SiemHttpDestinationDestinationTypeEnum{
		"generic":           SiemHttpDestinationDestinationTypeEnum("generic"),
		"splunk":            SiemHttpDestinationDestinationTypeEnum("splunk"),
		"azure_legacy":      SiemHttpDestinationDestinationTypeEnum("azure_legacy"),
		"qradar":            SiemHttpDestinationDestinationTypeEnum("qradar"),
		"sumo":              SiemHttpDestinationDestinationTypeEnum("sumo"),
		"rapid7":            SiemHttpDestinationDestinationTypeEnum("rapid7"),
		"solar_winds":       SiemHttpDestinationDestinationTypeEnum("solar_winds"),
		"new_relic":         SiemHttpDestinationDestinationTypeEnum("new_relic"),
		"datadog":           SiemHttpDestinationDestinationTypeEnum("datadog"),
		"azure":             SiemHttpDestinationDestinationTypeEnum("azure"),
		"file":              SiemHttpDestinationDestinationTypeEnum("file"),
		"crowdstrike":       SiemHttpDestinationDestinationTypeEnum("crowdstrike"),
		"splunk_compatible": SiemHttpDestinationDestinationTypeEnum("splunk_compatible"),
	}
}

// SiemHttpDestinationListParams contains the request parameters for GET /siem_http_destinations.
type SiemHttpDestinationListParams struct {
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// SiemHttpDestinationFindParams contains the request parameters for GET /siem_http_destinations/{id}.
type SiemHttpDestinationFindParams struct {
	// Siem Http Destination ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// SiemHttpDestinationCreateParams contains the request parameters for POST /siem_http_destinations.
type SiemHttpDestinationCreateParams struct {
	// Name for this Destination.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Additional HTTP Headers included in calls to the destination URL.
	AdditionalHeaders interface{} `url:"additional_headers,omitempty" json:"additional_headers,omitempty" path:"additional_headers"`
	// Whether this SIEM HTTP Destination is currently being sent to or not.
	SendingActive *bool `url:"sending_active,omitempty" json:"sending_active,omitempty" path:"sending_active"`
	// Applicable only for destination type: generic. Indicates the type of HTTP body. Can be json_newline or json_array. json_newline is multiple log entries as JSON separated by newlines. json_array is a single JSON array containing multiple log entries as JSON. Values: json_newline, json_array.
	GenericPayloadType SiemHttpDestinationGenericPayloadTypeEnum `url:"generic_payload_type,omitempty" json:"generic_payload_type,omitempty" path:"generic_payload_type"`
	// Applicable only for destination type: file. Destination folder path on Files.com.
	FileDestinationPath string `url:"file_destination_path,omitempty" json:"file_destination_path,omitempty" path:"file_destination_path"`
	// Applicable only for destination type: file. Generated file format. Values: json, csv.
	FileFormat SiemHttpDestinationFileFormatEnum `url:"file_format,omitempty" json:"file_format,omitempty" path:"file_format"`
	// Applicable only for destination type: file. Interval, in minutes, between file deliveries. Valid values are 5, 10, 15, 20, 30, 60, 90, 180, 240, 360.
	FileIntervalMinutes int64 `url:"file_interval_minutes,omitempty" json:"file_interval_minutes,omitempty" path:"file_interval_minutes"`
	// Applicable only for destination types: splunk, splunk_compatible. Authentication token for the destination.
	SplunkToken string `url:"splunk_token,omitempty" json:"splunk_token,omitempty" path:"splunk_token"`
	// Applicable only for destination type: crowdstrike. Authentication token provided by Crowdstrike.
	CrowdstrikeToken string `url:"crowdstrike_token,omitempty" json:"crowdstrike_token,omitempty" path:"crowdstrike_token"`
	// Applicable only for destination types: azure, azure_legacy. Immutable ID of the Data Collection Rule.
	AzureDcrImmutableId string `url:"azure_dcr_immutable_id,omitempty" json:"azure_dcr_immutable_id,omitempty" path:"azure_dcr_immutable_id"`
	// Applicable only for destination type: azure. Name of the stream in the DCR that represents the destination table.
	AzureStreamName string `url:"azure_stream_name,omitempty" json:"azure_stream_name,omitempty" path:"azure_stream_name"`
	// Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Tenant ID.
	AzureOauthClientCredentialsTenantId string `url:"azure_oauth_client_credentials_tenant_id,omitempty" json:"azure_oauth_client_credentials_tenant_id,omitempty" path:"azure_oauth_client_credentials_tenant_id"`
	// Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Client ID.
	AzureOauthClientCredentialsClientId string `url:"azure_oauth_client_credentials_client_id,omitempty" json:"azure_oauth_client_credentials_client_id,omitempty" path:"azure_oauth_client_credentials_client_id"`
	// Applicable only for destination type: azure. Client Credentials OAuth Client Secret.
	AzureOauthClientCredentialsClientSecret string `url:"azure_oauth_client_credentials_client_secret,omitempty" json:"azure_oauth_client_credentials_client_secret,omitempty" path:"azure_oauth_client_credentials_client_secret"`
	// Applicable only for destination type: qradar. Basic auth username provided by QRadar.
	QradarUsername string `url:"qradar_username,omitempty" json:"qradar_username,omitempty" path:"qradar_username"`
	// Applicable only for destination type: qradar. Basic auth password provided by QRadar.
	QradarPassword string `url:"qradar_password,omitempty" json:"qradar_password,omitempty" path:"qradar_password"`
	// Applicable only for destination type: solar_winds. Authentication token provided by Solar Winds.
	SolarWindsToken string `url:"solar_winds_token,omitempty" json:"solar_winds_token,omitempty" path:"solar_winds_token"`
	// Applicable only for destination type: new_relic. API key provided by New Relic.
	NewRelicApiKey string `url:"new_relic_api_key,omitempty" json:"new_relic_api_key,omitempty" path:"new_relic_api_key"`
	// Applicable only for destination type: datadog. API key provided by Datadog.
	DatadogApiKey string `url:"datadog_api_key,omitempty" json:"datadog_api_key,omitempty" path:"datadog_api_key"`
	// Whether or not sending is enabled for action logs.
	ActionSendEnabled *bool `url:"action_send_enabled,omitempty" json:"action_send_enabled,omitempty" path:"action_send_enabled"`
	// Whether or not sending is enabled for sftp_action logs.
	SftpActionSendEnabled *bool `url:"sftp_action_send_enabled,omitempty" json:"sftp_action_send_enabled,omitempty" path:"sftp_action_send_enabled"`
	// Whether or not sending is enabled for ftp_action logs.
	FtpActionSendEnabled *bool `url:"ftp_action_send_enabled,omitempty" json:"ftp_action_send_enabled,omitempty" path:"ftp_action_send_enabled"`
	// Whether or not sending is enabled for web_dav_action logs.
	WebDavActionSendEnabled *bool `url:"web_dav_action_send_enabled,omitempty" json:"web_dav_action_send_enabled,omitempty" path:"web_dav_action_send_enabled"`
	// Whether or not sending is enabled for sync logs.
	SyncSendEnabled *bool `url:"sync_send_enabled,omitempty" json:"sync_send_enabled,omitempty" path:"sync_send_enabled"`
	// Whether or not sending is enabled for outbound_connection logs.
	OutboundConnectionSendEnabled *bool `url:"outbound_connection_send_enabled,omitempty" json:"outbound_connection_send_enabled,omitempty" path:"outbound_connection_send_enabled"`
	// Whether or not sending is enabled for automation logs.
	AutomationSendEnabled *bool `url:"automation_send_enabled,omitempty" json:"automation_send_enabled,omitempty" path:"automation_send_enabled"`
	// Whether or not sending is enabled for api_request logs.
	ApiRequestSendEnabled *bool `url:"api_request_send_enabled,omitempty" json:"api_request_send_enabled,omitempty" path:"api_request_send_enabled"`
	// Whether or not sending is enabled for public_hosting_request logs.
	PublicHostingRequestSendEnabled *bool `url:"public_hosting_request_send_enabled,omitempty" json:"public_hosting_request_send_enabled,omitempty" path:"public_hosting_request_send_enabled"`
	// Whether or not sending is enabled for email logs.
	EmailSendEnabled *bool `url:"email_send_enabled,omitempty" json:"email_send_enabled,omitempty" path:"email_send_enabled"`
	// Whether or not sending is enabled for exavault_api_request logs.
	ExavaultApiRequestSendEnabled *bool `url:"exavault_api_request_send_enabled,omitempty" json:"exavault_api_request_send_enabled,omitempty" path:"exavault_api_request_send_enabled"`
	// Whether or not sending is enabled for settings_change logs.
	SettingsChangeSendEnabled *bool `url:"settings_change_send_enabled,omitempty" json:"settings_change_send_enabled,omitempty" path:"settings_change_send_enabled"`
	// Destination Type. Values: generic, splunk, azure_legacy, qradar, sumo, rapid7, solar_winds, new_relic, datadog, azure, file, crowdstrike, splunk_compatible. Required.
	DestinationType SiemHttpDestinationDestinationTypeEnum `url:"destination_type" json:"destination_type" path:"destination_type"`
	// Destination Url.
	DestinationUrl string `url:"destination_url,omitempty" json:"destination_url,omitempty" path:"destination_url"`
}

// SiemHttpDestinationSendTestEntryParams contains the request parameters for POST /siem_http_destinations/send_test_entry.
type SiemHttpDestinationSendTestEntryParams struct {
	// SIEM HTTP Destination ID.
	SiemHttpDestinationId int64 `url:"siem_http_destination_id,omitempty" json:"siem_http_destination_id,omitempty" path:"siem_http_destination_id"`
	// Destination Type. Values: generic, splunk, azure_legacy, qradar, sumo, rapid7, solar_winds, new_relic, datadog, azure, file, crowdstrike, splunk_compatible.
	DestinationType SiemHttpDestinationDestinationTypeEnum `url:"destination_type,omitempty" json:"destination_type,omitempty" path:"destination_type"`
	// Destination Url.
	DestinationUrl string `url:"destination_url,omitempty" json:"destination_url,omitempty" path:"destination_url"`
	// Name for this Destination.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Additional HTTP Headers included in calls to the destination URL.
	AdditionalHeaders interface{} `url:"additional_headers,omitempty" json:"additional_headers,omitempty" path:"additional_headers"`
	// Whether this SIEM HTTP Destination is currently being sent to or not.
	SendingActive *bool `url:"sending_active,omitempty" json:"sending_active,omitempty" path:"sending_active"`
	// Applicable only for destination type: generic. Indicates the type of HTTP body. Can be json_newline or json_array. json_newline is multiple log entries as JSON separated by newlines. json_array is a single JSON array containing multiple log entries as JSON. Values: json_newline, json_array.
	GenericPayloadType SiemHttpDestinationGenericPayloadTypeEnum `url:"generic_payload_type,omitempty" json:"generic_payload_type,omitempty" path:"generic_payload_type"`
	// Applicable only for destination type: file. Destination folder path on Files.com.
	FileDestinationPath string `url:"file_destination_path,omitempty" json:"file_destination_path,omitempty" path:"file_destination_path"`
	// Applicable only for destination type: file. Generated file format. Values: json, csv.
	FileFormat SiemHttpDestinationFileFormatEnum `url:"file_format,omitempty" json:"file_format,omitempty" path:"file_format"`
	// Applicable only for destination type: file. Interval, in minutes, between file deliveries. Valid values are 5, 10, 15, 20, 30, 60, 90, 180, 240, 360.
	FileIntervalMinutes int64 `url:"file_interval_minutes,omitempty" json:"file_interval_minutes,omitempty" path:"file_interval_minutes"`
	// Applicable only for destination types: splunk, splunk_compatible. Authentication token for the destination.
	SplunkToken string `url:"splunk_token,omitempty" json:"splunk_token,omitempty" path:"splunk_token"`
	// Applicable only for destination type: crowdstrike. Authentication token provided by Crowdstrike.
	CrowdstrikeToken string `url:"crowdstrike_token,omitempty" json:"crowdstrike_token,omitempty" path:"crowdstrike_token"`
	// Applicable only for destination types: azure, azure_legacy. Immutable ID of the Data Collection Rule.
	AzureDcrImmutableId string `url:"azure_dcr_immutable_id,omitempty" json:"azure_dcr_immutable_id,omitempty" path:"azure_dcr_immutable_id"`
	// Applicable only for destination type: azure. Name of the stream in the DCR that represents the destination table.
	AzureStreamName string `url:"azure_stream_name,omitempty" json:"azure_stream_name,omitempty" path:"azure_stream_name"`
	// Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Tenant ID.
	AzureOauthClientCredentialsTenantId string `url:"azure_oauth_client_credentials_tenant_id,omitempty" json:"azure_oauth_client_credentials_tenant_id,omitempty" path:"azure_oauth_client_credentials_tenant_id"`
	// Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Client ID.
	AzureOauthClientCredentialsClientId string `url:"azure_oauth_client_credentials_client_id,omitempty" json:"azure_oauth_client_credentials_client_id,omitempty" path:"azure_oauth_client_credentials_client_id"`
	// Applicable only for destination type: azure. Client Credentials OAuth Client Secret.
	AzureOauthClientCredentialsClientSecret string `url:"azure_oauth_client_credentials_client_secret,omitempty" json:"azure_oauth_client_credentials_client_secret,omitempty" path:"azure_oauth_client_credentials_client_secret"`
	// Applicable only for destination type: qradar. Basic auth username provided by QRadar.
	QradarUsername string `url:"qradar_username,omitempty" json:"qradar_username,omitempty" path:"qradar_username"`
	// Applicable only for destination type: qradar. Basic auth password provided by QRadar.
	QradarPassword string `url:"qradar_password,omitempty" json:"qradar_password,omitempty" path:"qradar_password"`
	// Applicable only for destination type: solar_winds. Authentication token provided by Solar Winds.
	SolarWindsToken string `url:"solar_winds_token,omitempty" json:"solar_winds_token,omitempty" path:"solar_winds_token"`
	// Applicable only for destination type: new_relic. API key provided by New Relic.
	NewRelicApiKey string `url:"new_relic_api_key,omitempty" json:"new_relic_api_key,omitempty" path:"new_relic_api_key"`
	// Applicable only for destination type: datadog. API key provided by Datadog.
	DatadogApiKey string `url:"datadog_api_key,omitempty" json:"datadog_api_key,omitempty" path:"datadog_api_key"`
	// Whether or not sending is enabled for action logs.
	ActionSendEnabled *bool `url:"action_send_enabled,omitempty" json:"action_send_enabled,omitempty" path:"action_send_enabled"`
	// Whether or not sending is enabled for sftp_action logs.
	SftpActionSendEnabled *bool `url:"sftp_action_send_enabled,omitempty" json:"sftp_action_send_enabled,omitempty" path:"sftp_action_send_enabled"`
	// Whether or not sending is enabled for ftp_action logs.
	FtpActionSendEnabled *bool `url:"ftp_action_send_enabled,omitempty" json:"ftp_action_send_enabled,omitempty" path:"ftp_action_send_enabled"`
	// Whether or not sending is enabled for web_dav_action logs.
	WebDavActionSendEnabled *bool `url:"web_dav_action_send_enabled,omitempty" json:"web_dav_action_send_enabled,omitempty" path:"web_dav_action_send_enabled"`
	// Whether or not sending is enabled for sync logs.
	SyncSendEnabled *bool `url:"sync_send_enabled,omitempty" json:"sync_send_enabled,omitempty" path:"sync_send_enabled"`
	// Whether or not sending is enabled for outbound_connection logs.
	OutboundConnectionSendEnabled *bool `url:"outbound_connection_send_enabled,omitempty" json:"outbound_connection_send_enabled,omitempty" path:"outbound_connection_send_enabled"`
	// Whether or not sending is enabled for automation logs.
	AutomationSendEnabled *bool `url:"automation_send_enabled,omitempty" json:"automation_send_enabled,omitempty" path:"automation_send_enabled"`
	// Whether or not sending is enabled for api_request logs.
	ApiRequestSendEnabled *bool `url:"api_request_send_enabled,omitempty" json:"api_request_send_enabled,omitempty" path:"api_request_send_enabled"`
	// Whether or not sending is enabled for public_hosting_request logs.
	PublicHostingRequestSendEnabled *bool `url:"public_hosting_request_send_enabled,omitempty" json:"public_hosting_request_send_enabled,omitempty" path:"public_hosting_request_send_enabled"`
	// Whether or not sending is enabled for email logs.
	EmailSendEnabled *bool `url:"email_send_enabled,omitempty" json:"email_send_enabled,omitempty" path:"email_send_enabled"`
	// Whether or not sending is enabled for exavault_api_request logs.
	ExavaultApiRequestSendEnabled *bool `url:"exavault_api_request_send_enabled,omitempty" json:"exavault_api_request_send_enabled,omitempty" path:"exavault_api_request_send_enabled"`
	// Whether or not sending is enabled for settings_change logs.
	SettingsChangeSendEnabled *bool `url:"settings_change_send_enabled,omitempty" json:"settings_change_send_enabled,omitempty" path:"settings_change_send_enabled"`
}

// SiemHttpDestinationUpdateParams contains the request parameters for PATCH /siem_http_destinations/{id}.
type SiemHttpDestinationUpdateParams struct {
	// Siem Http Destination ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Name for this Destination.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Additional HTTP Headers included in calls to the destination URL.
	AdditionalHeaders interface{} `url:"additional_headers,omitempty" json:"additional_headers,omitempty" path:"additional_headers"`
	// Whether this SIEM HTTP Destination is currently being sent to or not.
	SendingActive *bool `url:"sending_active,omitempty" json:"sending_active,omitempty" path:"sending_active"`
	// Applicable only for destination type: generic. Indicates the type of HTTP body. Can be json_newline or json_array. json_newline is multiple log entries as JSON separated by newlines. json_array is a single JSON array containing multiple log entries as JSON. Values: json_newline, json_array.
	GenericPayloadType SiemHttpDestinationGenericPayloadTypeEnum `url:"generic_payload_type,omitempty" json:"generic_payload_type,omitempty" path:"generic_payload_type"`
	// Applicable only for destination type: file. Destination folder path on Files.com.
	FileDestinationPath string `url:"file_destination_path,omitempty" json:"file_destination_path,omitempty" path:"file_destination_path"`
	// Applicable only for destination type: file. Generated file format. Values: json, csv.
	FileFormat SiemHttpDestinationFileFormatEnum `url:"file_format,omitempty" json:"file_format,omitempty" path:"file_format"`
	// Applicable only for destination type: file. Interval, in minutes, between file deliveries. Valid values are 5, 10, 15, 20, 30, 60, 90, 180, 240, 360.
	FileIntervalMinutes int64 `url:"file_interval_minutes,omitempty" json:"file_interval_minutes,omitempty" path:"file_interval_minutes"`
	// Applicable only for destination types: splunk, splunk_compatible. Authentication token for the destination.
	SplunkToken string `url:"splunk_token,omitempty" json:"splunk_token,omitempty" path:"splunk_token"`
	// Applicable only for destination type: crowdstrike. Authentication token provided by Crowdstrike.
	CrowdstrikeToken string `url:"crowdstrike_token,omitempty" json:"crowdstrike_token,omitempty" path:"crowdstrike_token"`
	// Applicable only for destination types: azure, azure_legacy. Immutable ID of the Data Collection Rule.
	AzureDcrImmutableId string `url:"azure_dcr_immutable_id,omitempty" json:"azure_dcr_immutable_id,omitempty" path:"azure_dcr_immutable_id"`
	// Applicable only for destination type: azure. Name of the stream in the DCR that represents the destination table.
	AzureStreamName string `url:"azure_stream_name,omitempty" json:"azure_stream_name,omitempty" path:"azure_stream_name"`
	// Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Tenant ID.
	AzureOauthClientCredentialsTenantId string `url:"azure_oauth_client_credentials_tenant_id,omitempty" json:"azure_oauth_client_credentials_tenant_id,omitempty" path:"azure_oauth_client_credentials_tenant_id"`
	// Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Client ID.
	AzureOauthClientCredentialsClientId string `url:"azure_oauth_client_credentials_client_id,omitempty" json:"azure_oauth_client_credentials_client_id,omitempty" path:"azure_oauth_client_credentials_client_id"`
	// Applicable only for destination type: azure. Client Credentials OAuth Client Secret.
	AzureOauthClientCredentialsClientSecret string `url:"azure_oauth_client_credentials_client_secret,omitempty" json:"azure_oauth_client_credentials_client_secret,omitempty" path:"azure_oauth_client_credentials_client_secret"`
	// Applicable only for destination type: qradar. Basic auth username provided by QRadar.
	QradarUsername string `url:"qradar_username,omitempty" json:"qradar_username,omitempty" path:"qradar_username"`
	// Applicable only for destination type: qradar. Basic auth password provided by QRadar.
	QradarPassword string `url:"qradar_password,omitempty" json:"qradar_password,omitempty" path:"qradar_password"`
	// Applicable only for destination type: solar_winds. Authentication token provided by Solar Winds.
	SolarWindsToken string `url:"solar_winds_token,omitempty" json:"solar_winds_token,omitempty" path:"solar_winds_token"`
	// Applicable only for destination type: new_relic. API key provided by New Relic.
	NewRelicApiKey string `url:"new_relic_api_key,omitempty" json:"new_relic_api_key,omitempty" path:"new_relic_api_key"`
	// Applicable only for destination type: datadog. API key provided by Datadog.
	DatadogApiKey string `url:"datadog_api_key,omitempty" json:"datadog_api_key,omitempty" path:"datadog_api_key"`
	// Whether or not sending is enabled for action logs.
	ActionSendEnabled *bool `url:"action_send_enabled,omitempty" json:"action_send_enabled,omitempty" path:"action_send_enabled"`
	// Whether or not sending is enabled for sftp_action logs.
	SftpActionSendEnabled *bool `url:"sftp_action_send_enabled,omitempty" json:"sftp_action_send_enabled,omitempty" path:"sftp_action_send_enabled"`
	// Whether or not sending is enabled for ftp_action logs.
	FtpActionSendEnabled *bool `url:"ftp_action_send_enabled,omitempty" json:"ftp_action_send_enabled,omitempty" path:"ftp_action_send_enabled"`
	// Whether or not sending is enabled for web_dav_action logs.
	WebDavActionSendEnabled *bool `url:"web_dav_action_send_enabled,omitempty" json:"web_dav_action_send_enabled,omitempty" path:"web_dav_action_send_enabled"`
	// Whether or not sending is enabled for sync logs.
	SyncSendEnabled *bool `url:"sync_send_enabled,omitempty" json:"sync_send_enabled,omitempty" path:"sync_send_enabled"`
	// Whether or not sending is enabled for outbound_connection logs.
	OutboundConnectionSendEnabled *bool `url:"outbound_connection_send_enabled,omitempty" json:"outbound_connection_send_enabled,omitempty" path:"outbound_connection_send_enabled"`
	// Whether or not sending is enabled for automation logs.
	AutomationSendEnabled *bool `url:"automation_send_enabled,omitempty" json:"automation_send_enabled,omitempty" path:"automation_send_enabled"`
	// Whether or not sending is enabled for api_request logs.
	ApiRequestSendEnabled *bool `url:"api_request_send_enabled,omitempty" json:"api_request_send_enabled,omitempty" path:"api_request_send_enabled"`
	// Whether or not sending is enabled for public_hosting_request logs.
	PublicHostingRequestSendEnabled *bool `url:"public_hosting_request_send_enabled,omitempty" json:"public_hosting_request_send_enabled,omitempty" path:"public_hosting_request_send_enabled"`
	// Whether or not sending is enabled for email logs.
	EmailSendEnabled *bool `url:"email_send_enabled,omitempty" json:"email_send_enabled,omitempty" path:"email_send_enabled"`
	// Whether or not sending is enabled for exavault_api_request logs.
	ExavaultApiRequestSendEnabled *bool `url:"exavault_api_request_send_enabled,omitempty" json:"exavault_api_request_send_enabled,omitempty" path:"exavault_api_request_send_enabled"`
	// Whether or not sending is enabled for settings_change logs.
	SettingsChangeSendEnabled *bool `url:"settings_change_send_enabled,omitempty" json:"settings_change_send_enabled,omitempty" path:"settings_change_send_enabled"`
	// Destination Type. Values: generic, splunk, azure_legacy, qradar, sumo, rapid7, solar_winds, new_relic, datadog, azure, file, crowdstrike, splunk_compatible.
	DestinationType SiemHttpDestinationDestinationTypeEnum `url:"destination_type,omitempty" json:"destination_type,omitempty" path:"destination_type"`
	// Destination Url.
	DestinationUrl string `url:"destination_url,omitempty" json:"destination_url,omitempty" path:"destination_url"`
}

// SiemHttpDestinationDeleteParams contains the request parameters for DELETE /siem_http_destinations/{id}.
type SiemHttpDestinationDeleteParams struct {
	// Siem Http Destination ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (s *SiemHttpDestination) UnmarshalJSON(data []byte) error {
	type siemHttpDestination SiemHttpDestination
	var v siemHttpDestination
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = SiemHttpDestination(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (s *SiemHttpDestinationCollection) UnmarshalJSON(data []byte) error {
	type siemHttpDestinations SiemHttpDestinationCollection
	var v siemHttpDestinations
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SiemHttpDestinationCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (s *SiemHttpDestinationCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
