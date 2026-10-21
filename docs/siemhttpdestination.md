# SiemHttpDestination

`files_sdk.SiemHttpDestination` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/siemhttpdestination`, named `siem_http_destination`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	siem_http_destination "github.com/Files-com/files-sdk-go/v3/siemhttpdestination"
)
```

## Operations

### List

`GET /siem_http_destinations`

List SIEM HTTP Destinations

```go
func (c *Client) List(params files_sdk.SiemHttpDestinationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.SiemHttpDestinationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.SiemHttpDestination` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.SiemHttpDestinationListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := siem_http_destination.List(files_sdk.SiemHttpDestinationListParams{})
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

`GET /siem_http_destinations/{id}`

Show SIEM HTTP Destination

```go
func (c *Client) Find(params files_sdk.SiemHttpDestinationFindParams, opts ...files_sdk.RequestResponseOption) (siemHttpDestination files_sdk.SiemHttpDestination, err error)
func Find(params files_sdk.SiemHttpDestinationFindParams, opts ...files_sdk.RequestResponseOption) (siemHttpDestination files_sdk.SiemHttpDestination, err error)
```

Parameters, `files_sdk.SiemHttpDestinationFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Siem Http Destination ID. |

Example:

```go
result, err := siem_http_destination.Find(files_sdk.SiemHttpDestinationFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /siem_http_destinations`

Create SIEM HTTP Destination

```go
func (c *Client) Create(params files_sdk.SiemHttpDestinationCreateParams, opts ...files_sdk.RequestResponseOption) (siemHttpDestination files_sdk.SiemHttpDestination, err error)
func Create(params files_sdk.SiemHttpDestinationCreateParams, opts ...files_sdk.RequestResponseOption) (siemHttpDestination files_sdk.SiemHttpDestination, err error)
```

Parameters, `files_sdk.SiemHttpDestinationCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Name` | `string` | `name` | no | Name for this Destination |
| `AdditionalHeaders` | `interface{}` | `additional_headers` | no | Additional HTTP Headers included in calls to the destination URL |
| `SendingActive` | `*bool` | `sending_active` | no | Whether this SIEM HTTP Destination is currently being sent to or not |
| `GenericPayloadType` | `files_sdk.SiemHttpDestinationGenericPayloadTypeEnum` | `generic_payload_type` | no | Applicable only for destination type: generic. Indicates the type of HTTP body. Can be json_newline or json_array. json_newline is multiple log entries as JSON separated by newlines. json_array is a single JSON array containing multiple log entries as JSON. Values: `json_newline`, `json_array`. |
| `FileDestinationPath` | `string` | `file_destination_path` | no | Applicable only for destination type: file. Destination folder path on Files.com. |
| `FileFormat` | `files_sdk.SiemHttpDestinationFileFormatEnum` | `file_format` | no | Applicable only for destination type: file. Generated file format. Values: `json`, `csv`. |
| `FileIntervalMinutes` | `int64` | `file_interval_minutes` | no | Applicable only for destination type: file. Interval, in minutes, between file deliveries. Valid values are 5, 10, 15, 20, 30, 60, 90, 180, 240, 360. |
| `SplunkToken` | `string` | `splunk_token` | no | Applicable only for destination types: splunk, splunk_compatible. Authentication token for the destination. |
| `CrowdstrikeToken` | `string` | `crowdstrike_token` | no | Applicable only for destination type: crowdstrike. Authentication token provided by Crowdstrike. |
| `AzureDcrImmutableId` | `string` | `azure_dcr_immutable_id` | no | Applicable only for destination types: azure, azure_legacy. Immutable ID of the Data Collection Rule. |
| `AzureStreamName` | `string` | `azure_stream_name` | no | Applicable only for destination type: azure. Name of the stream in the DCR that represents the destination table. |
| `AzureOauthClientCredentialsTenantId` | `string` | `azure_oauth_client_credentials_tenant_id` | no | Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Tenant ID. |
| `AzureOauthClientCredentialsClientId` | `string` | `azure_oauth_client_credentials_client_id` | no | Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Client ID. |
| `AzureOauthClientCredentialsClientSecret` | `string` | `azure_oauth_client_credentials_client_secret` | no | Applicable only for destination type: azure. Client Credentials OAuth Client Secret. |
| `QradarUsername` | `string` | `qradar_username` | no | Applicable only for destination type: qradar. Basic auth username provided by QRadar. |
| `QradarPassword` | `string` | `qradar_password` | no | Applicable only for destination type: qradar. Basic auth password provided by QRadar. |
| `SolarWindsToken` | `string` | `solar_winds_token` | no | Applicable only for destination type: solar_winds. Authentication token provided by Solar Winds. |
| `NewRelicApiKey` | `string` | `new_relic_api_key` | no | Applicable only for destination type: new_relic. API key provided by New Relic. |
| `DatadogApiKey` | `string` | `datadog_api_key` | no | Applicable only for destination type: datadog. API key provided by Datadog. |
| `ActionSendEnabled` | `*bool` | `action_send_enabled` | no | Whether or not sending is enabled for action logs. |
| `SftpActionSendEnabled` | `*bool` | `sftp_action_send_enabled` | no | Whether or not sending is enabled for sftp_action logs. |
| `FtpActionSendEnabled` | `*bool` | `ftp_action_send_enabled` | no | Whether or not sending is enabled for ftp_action logs. |
| `WebDavActionSendEnabled` | `*bool` | `web_dav_action_send_enabled` | no | Whether or not sending is enabled for web_dav_action logs. |
| `SyncSendEnabled` | `*bool` | `sync_send_enabled` | no | Whether or not sending is enabled for sync logs. |
| `OutboundConnectionSendEnabled` | `*bool` | `outbound_connection_send_enabled` | no | Whether or not sending is enabled for outbound_connection logs. |
| `AutomationSendEnabled` | `*bool` | `automation_send_enabled` | no | Whether or not sending is enabled for automation logs. |
| `ApiRequestSendEnabled` | `*bool` | `api_request_send_enabled` | no | Whether or not sending is enabled for api_request logs. |
| `PublicHostingRequestSendEnabled` | `*bool` | `public_hosting_request_send_enabled` | no | Whether or not sending is enabled for public_hosting_request logs. |
| `EmailSendEnabled` | `*bool` | `email_send_enabled` | no | Whether or not sending is enabled for email logs. |
| `ExavaultApiRequestSendEnabled` | `*bool` | `exavault_api_request_send_enabled` | no | Whether or not sending is enabled for exavault_api_request logs. |
| `SettingsChangeSendEnabled` | `*bool` | `settings_change_send_enabled` | no | Whether or not sending is enabled for settings_change logs. |
| `DestinationType` | `files_sdk.SiemHttpDestinationDestinationTypeEnum` | `destination_type` | yes | Destination Type Values: `generic`, `splunk`, `azure_legacy`, `qradar`, `sumo`, `rapid7`, `solar_winds`, `new_relic`, `datadog`, `azure`, `file`, `crowdstrike`, `splunk_compatible`. |
| `DestinationUrl` | `string` | `destination_url` | no | Destination Url |

Example:

```go
result, err := siem_http_destination.Create(files_sdk.SiemHttpDestinationCreateParams{
	DestinationType: "generic",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### SendTestEntry

`POST /siem_http_destinations/send_test_entry`

send_test_entry SIEM HTTP Destination

```go
func (c *Client) SendTestEntry(params files_sdk.SiemHttpDestinationSendTestEntryParams, opts ...files_sdk.RequestResponseOption) (err error)
func SendTestEntry(params files_sdk.SiemHttpDestinationSendTestEntryParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.SiemHttpDestinationSendTestEntryParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SiemHttpDestinationId` | `int64` | `siem_http_destination_id` | no | SIEM HTTP Destination ID |
| `DestinationType` | `files_sdk.SiemHttpDestinationDestinationTypeEnum` | `destination_type` | no | Destination Type Values: `generic`, `splunk`, `azure_legacy`, `qradar`, `sumo`, `rapid7`, `solar_winds`, `new_relic`, `datadog`, `azure`, `file`, `crowdstrike`, `splunk_compatible`. |
| `DestinationUrl` | `string` | `destination_url` | no | Destination Url |
| `Name` | `string` | `name` | no | Name for this Destination |
| `AdditionalHeaders` | `interface{}` | `additional_headers` | no | Additional HTTP Headers included in calls to the destination URL |
| `SendingActive` | `*bool` | `sending_active` | no | Whether this SIEM HTTP Destination is currently being sent to or not |
| `GenericPayloadType` | `files_sdk.SiemHttpDestinationGenericPayloadTypeEnum` | `generic_payload_type` | no | Applicable only for destination type: generic. Indicates the type of HTTP body. Can be json_newline or json_array. json_newline is multiple log entries as JSON separated by newlines. json_array is a single JSON array containing multiple log entries as JSON. Values: `json_newline`, `json_array`. |
| `FileDestinationPath` | `string` | `file_destination_path` | no | Applicable only for destination type: file. Destination folder path on Files.com. |
| `FileFormat` | `files_sdk.SiemHttpDestinationFileFormatEnum` | `file_format` | no | Applicable only for destination type: file. Generated file format. Values: `json`, `csv`. |
| `FileIntervalMinutes` | `int64` | `file_interval_minutes` | no | Applicable only for destination type: file. Interval, in minutes, between file deliveries. Valid values are 5, 10, 15, 20, 30, 60, 90, 180, 240, 360. |
| `SplunkToken` | `string` | `splunk_token` | no | Applicable only for destination types: splunk, splunk_compatible. Authentication token for the destination. |
| `CrowdstrikeToken` | `string` | `crowdstrike_token` | no | Applicable only for destination type: crowdstrike. Authentication token provided by Crowdstrike. |
| `AzureDcrImmutableId` | `string` | `azure_dcr_immutable_id` | no | Applicable only for destination types: azure, azure_legacy. Immutable ID of the Data Collection Rule. |
| `AzureStreamName` | `string` | `azure_stream_name` | no | Applicable only for destination type: azure. Name of the stream in the DCR that represents the destination table. |
| `AzureOauthClientCredentialsTenantId` | `string` | `azure_oauth_client_credentials_tenant_id` | no | Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Tenant ID. |
| `AzureOauthClientCredentialsClientId` | `string` | `azure_oauth_client_credentials_client_id` | no | Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Client ID. |
| `AzureOauthClientCredentialsClientSecret` | `string` | `azure_oauth_client_credentials_client_secret` | no | Applicable only for destination type: azure. Client Credentials OAuth Client Secret. |
| `QradarUsername` | `string` | `qradar_username` | no | Applicable only for destination type: qradar. Basic auth username provided by QRadar. |
| `QradarPassword` | `string` | `qradar_password` | no | Applicable only for destination type: qradar. Basic auth password provided by QRadar. |
| `SolarWindsToken` | `string` | `solar_winds_token` | no | Applicable only for destination type: solar_winds. Authentication token provided by Solar Winds. |
| `NewRelicApiKey` | `string` | `new_relic_api_key` | no | Applicable only for destination type: new_relic. API key provided by New Relic. |
| `DatadogApiKey` | `string` | `datadog_api_key` | no | Applicable only for destination type: datadog. API key provided by Datadog. |
| `ActionSendEnabled` | `*bool` | `action_send_enabled` | no | Whether or not sending is enabled for action logs. |
| `SftpActionSendEnabled` | `*bool` | `sftp_action_send_enabled` | no | Whether or not sending is enabled for sftp_action logs. |
| `FtpActionSendEnabled` | `*bool` | `ftp_action_send_enabled` | no | Whether or not sending is enabled for ftp_action logs. |
| `WebDavActionSendEnabled` | `*bool` | `web_dav_action_send_enabled` | no | Whether or not sending is enabled for web_dav_action logs. |
| `SyncSendEnabled` | `*bool` | `sync_send_enabled` | no | Whether or not sending is enabled for sync logs. |
| `OutboundConnectionSendEnabled` | `*bool` | `outbound_connection_send_enabled` | no | Whether or not sending is enabled for outbound_connection logs. |
| `AutomationSendEnabled` | `*bool` | `automation_send_enabled` | no | Whether or not sending is enabled for automation logs. |
| `ApiRequestSendEnabled` | `*bool` | `api_request_send_enabled` | no | Whether or not sending is enabled for api_request logs. |
| `PublicHostingRequestSendEnabled` | `*bool` | `public_hosting_request_send_enabled` | no | Whether or not sending is enabled for public_hosting_request logs. |
| `EmailSendEnabled` | `*bool` | `email_send_enabled` | no | Whether or not sending is enabled for email logs. |
| `ExavaultApiRequestSendEnabled` | `*bool` | `exavault_api_request_send_enabled` | no | Whether or not sending is enabled for exavault_api_request logs. |
| `SettingsChangeSendEnabled` | `*bool` | `settings_change_send_enabled` | no | Whether or not sending is enabled for settings_change logs. |

Example:

```go
err := siem_http_destination.SendTestEntry(files_sdk.SiemHttpDestinationSendTestEntryParams{})
if err != nil {
	return err
}
```

### Update

`PATCH /siem_http_destinations/{id}`

Update SIEM HTTP Destination

```go
func (c *Client) Update(params files_sdk.SiemHttpDestinationUpdateParams, opts ...files_sdk.RequestResponseOption) (siemHttpDestination files_sdk.SiemHttpDestination, err error)
func Update(params files_sdk.SiemHttpDestinationUpdateParams, opts ...files_sdk.RequestResponseOption) (siemHttpDestination files_sdk.SiemHttpDestination, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (siemHttpDestination files_sdk.SiemHttpDestination, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (siemHttpDestination files_sdk.SiemHttpDestination, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.SiemHttpDestinationUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Siem Http Destination ID. |
| `Name` | `string` | `name` | no | Name for this Destination |
| `AdditionalHeaders` | `interface{}` | `additional_headers` | no | Additional HTTP Headers included in calls to the destination URL |
| `SendingActive` | `*bool` | `sending_active` | no | Whether this SIEM HTTP Destination is currently being sent to or not |
| `GenericPayloadType` | `files_sdk.SiemHttpDestinationGenericPayloadTypeEnum` | `generic_payload_type` | no | Applicable only for destination type: generic. Indicates the type of HTTP body. Can be json_newline or json_array. json_newline is multiple log entries as JSON separated by newlines. json_array is a single JSON array containing multiple log entries as JSON. Values: `json_newline`, `json_array`. |
| `FileDestinationPath` | `string` | `file_destination_path` | no | Applicable only for destination type: file. Destination folder path on Files.com. |
| `FileFormat` | `files_sdk.SiemHttpDestinationFileFormatEnum` | `file_format` | no | Applicable only for destination type: file. Generated file format. Values: `json`, `csv`. |
| `FileIntervalMinutes` | `int64` | `file_interval_minutes` | no | Applicable only for destination type: file. Interval, in minutes, between file deliveries. Valid values are 5, 10, 15, 20, 30, 60, 90, 180, 240, 360. |
| `SplunkToken` | `string` | `splunk_token` | no | Applicable only for destination types: splunk, splunk_compatible. Authentication token for the destination. |
| `CrowdstrikeToken` | `string` | `crowdstrike_token` | no | Applicable only for destination type: crowdstrike. Authentication token provided by Crowdstrike. |
| `AzureDcrImmutableId` | `string` | `azure_dcr_immutable_id` | no | Applicable only for destination types: azure, azure_legacy. Immutable ID of the Data Collection Rule. |
| `AzureStreamName` | `string` | `azure_stream_name` | no | Applicable only for destination type: azure. Name of the stream in the DCR that represents the destination table. |
| `AzureOauthClientCredentialsTenantId` | `string` | `azure_oauth_client_credentials_tenant_id` | no | Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Tenant ID. |
| `AzureOauthClientCredentialsClientId` | `string` | `azure_oauth_client_credentials_client_id` | no | Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Client ID. |
| `AzureOauthClientCredentialsClientSecret` | `string` | `azure_oauth_client_credentials_client_secret` | no | Applicable only for destination type: azure. Client Credentials OAuth Client Secret. |
| `QradarUsername` | `string` | `qradar_username` | no | Applicable only for destination type: qradar. Basic auth username provided by QRadar. |
| `QradarPassword` | `string` | `qradar_password` | no | Applicable only for destination type: qradar. Basic auth password provided by QRadar. |
| `SolarWindsToken` | `string` | `solar_winds_token` | no | Applicable only for destination type: solar_winds. Authentication token provided by Solar Winds. |
| `NewRelicApiKey` | `string` | `new_relic_api_key` | no | Applicable only for destination type: new_relic. API key provided by New Relic. |
| `DatadogApiKey` | `string` | `datadog_api_key` | no | Applicable only for destination type: datadog. API key provided by Datadog. |
| `ActionSendEnabled` | `*bool` | `action_send_enabled` | no | Whether or not sending is enabled for action logs. |
| `SftpActionSendEnabled` | `*bool` | `sftp_action_send_enabled` | no | Whether or not sending is enabled for sftp_action logs. |
| `FtpActionSendEnabled` | `*bool` | `ftp_action_send_enabled` | no | Whether or not sending is enabled for ftp_action logs. |
| `WebDavActionSendEnabled` | `*bool` | `web_dav_action_send_enabled` | no | Whether or not sending is enabled for web_dav_action logs. |
| `SyncSendEnabled` | `*bool` | `sync_send_enabled` | no | Whether or not sending is enabled for sync logs. |
| `OutboundConnectionSendEnabled` | `*bool` | `outbound_connection_send_enabled` | no | Whether or not sending is enabled for outbound_connection logs. |
| `AutomationSendEnabled` | `*bool` | `automation_send_enabled` | no | Whether or not sending is enabled for automation logs. |
| `ApiRequestSendEnabled` | `*bool` | `api_request_send_enabled` | no | Whether or not sending is enabled for api_request logs. |
| `PublicHostingRequestSendEnabled` | `*bool` | `public_hosting_request_send_enabled` | no | Whether or not sending is enabled for public_hosting_request logs. |
| `EmailSendEnabled` | `*bool` | `email_send_enabled` | no | Whether or not sending is enabled for email logs. |
| `ExavaultApiRequestSendEnabled` | `*bool` | `exavault_api_request_send_enabled` | no | Whether or not sending is enabled for exavault_api_request logs. |
| `SettingsChangeSendEnabled` | `*bool` | `settings_change_send_enabled` | no | Whether or not sending is enabled for settings_change logs. |
| `DestinationType` | `files_sdk.SiemHttpDestinationDestinationTypeEnum` | `destination_type` | no | Destination Type Values: `generic`, `splunk`, `azure_legacy`, `qradar`, `sumo`, `rapid7`, `solar_winds`, `new_relic`, `datadog`, `azure`, `file`, `crowdstrike`, `splunk_compatible`. |
| `DestinationUrl` | `string` | `destination_url` | no | Destination Url |

Example:

```go
result, err := siem_http_destination.Update(files_sdk.SiemHttpDestinationUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := siem_http_destination.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /siem_http_destinations/{id}`

Delete SIEM HTTP Destination

```go
func (c *Client) Delete(params files_sdk.SiemHttpDestinationDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.SiemHttpDestinationDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.SiemHttpDestinationDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Siem Http Destination ID. |

Example:

```go
err := siem_http_destination.Delete(files_sdk.SiemHttpDestinationDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.SiemHttpDestination` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | SIEM HTTP Destination ID |
| `Name` | `string` | `name` | Name for this Destination |
| `DestinationType` | `string` | `destination_type` | Destination Type |
| `DestinationUrl` | `string` | `destination_url` | Destination Url |
| `FileDestinationPath` | `string` | `file_destination_path` | Applicable only for destination type: file. Destination folder path on Files.com. |
| `FileFormat` | `string` | `file_format` | Applicable only for destination type: file. Generated file format. |
| `FileIntervalMinutes` | `int64` | `file_interval_minutes` | Applicable only for destination type: file. Interval, in minutes, between file deliveries. |
| `AdditionalHeaders` | `interface{}` | `additional_headers` | Additional HTTP Headers included in calls to the destination URL |
| `SendingActive` | `*bool` | `sending_active` | Whether this SIEM HTTP Destination is currently being sent to or not |
| `GenericPayloadType` | `string` | `generic_payload_type` | Applicable only for destination type: generic. Indicates the type of HTTP body. Can be json_newline or json_array. json_newline is multiple log entries as JSON separated by newlines. json_array is a single JSON array containing multiple log entries as JSON. |
| `SplunkTokenMasked` | `string` | `splunk_token_masked` | Applicable only for destination types: splunk, splunk_compatible. Authentication token for the destination. |
| `CrowdstrikeTokenMasked` | `string` | `crowdstrike_token_masked` | Applicable only for destination type: crowdstrike. Authentication token provided by Crowdstrike. |
| `AzureDcrImmutableId` | `string` | `azure_dcr_immutable_id` | Applicable only for destination types: azure, azure_legacy. Immutable ID of the Data Collection Rule. |
| `AzureStreamName` | `string` | `azure_stream_name` | Applicable only for destination type: azure. Name of the stream in the DCR that represents the destination table. |
| `AzureOauthClientCredentialsTenantId` | `string` | `azure_oauth_client_credentials_tenant_id` | Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Tenant ID. |
| `AzureOauthClientCredentialsClientId` | `string` | `azure_oauth_client_credentials_client_id` | Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Client ID. |
| `AzureOauthClientCredentialsClientSecretMasked` | `string` | `azure_oauth_client_credentials_client_secret_masked` | Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Client Secret. |
| `QradarUsername` | `string` | `qradar_username` | Applicable only for destination type: qradar. Basic auth username provided by QRadar. |
| `QradarPasswordMasked` | `string` | `qradar_password_masked` | Applicable only for destination type: qradar. Basic auth password provided by QRadar. |
| `SolarWindsTokenMasked` | `string` | `solar_winds_token_masked` | Applicable only for destination type: solar_winds. Authentication token provided by Solar Winds. |
| `NewRelicApiKeyMasked` | `string` | `new_relic_api_key_masked` | Applicable only for destination type: new_relic. API key provided by New Relic. |
| `DatadogApiKeyMasked` | `string` | `datadog_api_key_masked` | Applicable only for destination type: datadog. API key provided by Datadog. |
| `ActionSendEnabled` | `*bool` | `action_send_enabled` | Whether or not sending is enabled for action logs. |
| `ActionEntriesSent` | `int64` | `action_entries_sent` | Number of log entries sent for the lifetime of this destination. |
| `SftpActionSendEnabled` | `*bool` | `sftp_action_send_enabled` | Whether or not sending is enabled for sftp_action logs. |
| `SftpActionEntriesSent` | `int64` | `sftp_action_entries_sent` | Number of log entries sent for the lifetime of this destination. |
| `FtpActionSendEnabled` | `*bool` | `ftp_action_send_enabled` | Whether or not sending is enabled for ftp_action logs. |
| `FtpActionEntriesSent` | `int64` | `ftp_action_entries_sent` | Number of log entries sent for the lifetime of this destination. |
| `WebDavActionSendEnabled` | `*bool` | `web_dav_action_send_enabled` | Whether or not sending is enabled for web_dav_action logs. |
| `WebDavActionEntriesSent` | `int64` | `web_dav_action_entries_sent` | Number of log entries sent for the lifetime of this destination. |
| `SyncSendEnabled` | `*bool` | `sync_send_enabled` | Whether or not sending is enabled for sync logs. |
| `SyncEntriesSent` | `int64` | `sync_entries_sent` | Number of log entries sent for the lifetime of this destination. |
| `OutboundConnectionSendEnabled` | `*bool` | `outbound_connection_send_enabled` | Whether or not sending is enabled for outbound_connection logs. |
| `OutboundConnectionEntriesSent` | `int64` | `outbound_connection_entries_sent` | Number of log entries sent for the lifetime of this destination. |
| `AutomationSendEnabled` | `*bool` | `automation_send_enabled` | Whether or not sending is enabled for automation logs. |
| `AutomationEntriesSent` | `int64` | `automation_entries_sent` | Number of log entries sent for the lifetime of this destination. |
| `ApiRequestSendEnabled` | `*bool` | `api_request_send_enabled` | Whether or not sending is enabled for api_request logs. |
| `ApiRequestEntriesSent` | `int64` | `api_request_entries_sent` | Number of log entries sent for the lifetime of this destination. |
| `PublicHostingRequestSendEnabled` | `*bool` | `public_hosting_request_send_enabled` | Whether or not sending is enabled for public_hosting_request logs. |
| `PublicHostingRequestEntriesSent` | `int64` | `public_hosting_request_entries_sent` | Number of log entries sent for the lifetime of this destination. |
| `EmailSendEnabled` | `*bool` | `email_send_enabled` | Whether or not sending is enabled for email logs. |
| `EmailEntriesSent` | `int64` | `email_entries_sent` | Number of log entries sent for the lifetime of this destination. |
| `ExavaultApiRequestSendEnabled` | `*bool` | `exavault_api_request_send_enabled` | Whether or not sending is enabled for exavault_api_request logs. |
| `ExavaultApiRequestEntriesSent` | `int64` | `exavault_api_request_entries_sent` | Number of log entries sent for the lifetime of this destination. |
| `SettingsChangeSendEnabled` | `*bool` | `settings_change_send_enabled` | Whether or not sending is enabled for settings_change logs. |
| `SettingsChangeEntriesSent` | `int64` | `settings_change_entries_sent` | Number of log entries sent for the lifetime of this destination. |
| `LastHttpCallTargetType` | `string` | `last_http_call_target_type` | Type of URL that was last called. Can be `destination_url` or `azure_oauth_client_credentials_url` |
| `LastHttpCallSuccess` | `*bool` | `last_http_call_success` | Was the last HTTP call made successful? |
| `LastHttpCallResponseCode` | `int64` | `last_http_call_response_code` | Last HTTP Call Response Code |
| `LastHttpCallResponseBody` | `string` | `last_http_call_response_body` | Last HTTP Call Response Body. Large responses are truncated. |
| `LastHttpCallErrorMessage` | `string` | `last_http_call_error_message` | Last HTTP Call Error Message if applicable |
| `LastHttpCallTime` | `string` | `last_http_call_time` | Time of Last HTTP Call |
| `LastHttpCallDurationMs` | `int64` | `last_http_call_duration_ms` | Duration of the last HTTP Call in milliseconds |
| `MostRecentHttpCallSuccessTime` | `string` | `most_recent_http_call_success_time` | Time of Most Recent Successful HTTP Call |
| `ConnectionTestEntry` | `string` | `connection_test_entry` | Connection Test Entry |
| `SplunkToken` | `string` | `splunk_token` | Applicable only for destination types: splunk, splunk_compatible. Authentication token for the destination. |
| `CrowdstrikeToken` | `string` | `crowdstrike_token` | Applicable only for destination type: crowdstrike. Authentication token provided by Crowdstrike. |
| `AzureOauthClientCredentialsClientSecret` | `string` | `azure_oauth_client_credentials_client_secret` | Applicable only for destination type: azure. Client Credentials OAuth Client Secret. |
| `QradarPassword` | `string` | `qradar_password` | Applicable only for destination type: qradar. Basic auth password provided by QRadar. |
| `SolarWindsToken` | `string` | `solar_winds_token` | Applicable only for destination type: solar_winds. Authentication token provided by Solar Winds. |
| `NewRelicApiKey` | `string` | `new_relic_api_key` | Applicable only for destination type: new_relic. API key provided by New Relic. |
| `DatadogApiKey` | `string` | `datadog_api_key` | Applicable only for destination type: datadog. API key provided by Datadog. |
