# RemoteServer

`files_sdk.RemoteServer` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/remoteserver`, named `remote_server`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	remote_server "github.com/Files-com/files-sdk-go/v3/remoteserver"
)
```

## Operations

### List

`GET /remote_servers`

List Remote Servers

```go
func (c *Client) List(params files_sdk.RemoteServerListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.RemoteServerListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.RemoteServer` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.RemoteServerListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `name`, `server_type`, `backblaze_b2_bucket`, `google_cloud_storage_bucket`, `wasabi_bucket`, `s3_bucket`, `azure_blob_storage_container`, `azure_files_storage_share_name`, `s3_compatible_bucket`, `filebase_bucket`, `cloudflare_bucket` or `linode_bucket`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `name`, `server_type`, `workspace_id`, `backblaze_b2_bucket`, `google_cloud_storage_bucket`, `wasabi_bucket`, `s3_bucket`, `azure_blob_storage_container`, `azure_files_storage_share_name`, `s3_compatible_bucket`, `filebase_bucket`, `cloudflare_bucket` or `linode_bucket`. Valid field combinations are `[ server_type, name ]`, `[ workspace_id, name ]`, `[ backblaze_b2_bucket, name ]`, `[ google_cloud_storage_bucket, name ]`, `[ wasabi_bucket, name ]`, `[ s3_bucket, name ]`, `[ azure_blob_storage_container, name ]`, `[ azure_files_storage_share_name, name ]`, `[ s3_compatible_bucket, name ]`, `[ filebase_bucket, name ]`, `[ cloudflare_bucket, name ]`, `[ linode_bucket, name ]`, `[ workspace_id, server_type ]` or `[ workspace_id, server_type, name ]`. |
| `FilterPrefix` | `interface{}` | `filter_prefix` | no | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `name`, `backblaze_b2_bucket`, `google_cloud_storage_bucket`, `wasabi_bucket`, `s3_bucket`, `azure_blob_storage_container`, `azure_files_storage_share_name`, `s3_compatible_bucket`, `filebase_bucket`, `cloudflare_bucket` or `linode_bucket`. Valid field combinations are `[ backblaze_b2_bucket, name ]`, `[ google_cloud_storage_bucket, name ]`, `[ wasabi_bucket, name ]`, `[ s3_bucket, name ]`, `[ azure_blob_storage_container, name ]`, `[ azure_files_storage_share_name, name ]`, `[ s3_compatible_bucket, name ]`, `[ filebase_bucket, name ]`, `[ cloudflare_bucket, name ]` or `[ linode_bucket, name ]`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := remote_server.List(files_sdk.RemoteServerListParams{})
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

`GET /remote_servers/{id}`

Show Remote Server

```go
func (c *Client) Find(params files_sdk.RemoteServerFindParams, opts ...files_sdk.RequestResponseOption) (remoteServer files_sdk.RemoteServer, err error)
func Find(params files_sdk.RemoteServerFindParams, opts ...files_sdk.RequestResponseOption) (remoteServer files_sdk.RemoteServer, err error)
```

Parameters, `files_sdk.RemoteServerFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Remote Server ID. |

Example:

```go
result, err := remote_server.Find(files_sdk.RemoteServerFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### AgentNodes

`GET /remote_servers/{id}/agent_nodes`

List Files.com Agent nodes

```go
func (c *Client) AgentNodes(params files_sdk.RemoteServerAgentNodesParams, opts ...files_sdk.RequestResponseOption) (agentNode files_sdk.AgentNode, err error)
func AgentNodes(params files_sdk.RemoteServerAgentNodesParams, opts ...files_sdk.RequestResponseOption) (agentNode files_sdk.AgentNode, err error)
```

Parameters, `files_sdk.RemoteServerAgentNodesParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Remote Server ID. |

Example:

```go
result, err := remote_server.AgentNodes(files_sdk.RemoteServerAgentNodesParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### FindConfigurationFile

`GET /remote_servers/{id}/configuration_file`

Download configuration file (required for some Remote Server integrations, such as the Files.com Agent)

```go
func (c *Client) FindConfigurationFile(params files_sdk.RemoteServerFindConfigurationFileParams, opts ...files_sdk.RequestResponseOption) (remoteServerConfigurationFile files_sdk.RemoteServerConfigurationFile, err error)
func FindConfigurationFile(params files_sdk.RemoteServerFindConfigurationFileParams, opts ...files_sdk.RequestResponseOption) (remoteServerConfigurationFile files_sdk.RemoteServerConfigurationFile, err error)
```

Parameters, `files_sdk.RemoteServerFindConfigurationFileParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Remote Server ID. |

Example:

```go
result, err := remote_server.FindConfigurationFile(files_sdk.RemoteServerFindConfigurationFileParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /remote_servers`

Create Remote Server

```go
func (c *Client) Create(params files_sdk.RemoteServerCreateParams, opts ...files_sdk.RequestResponseOption) (remoteServer files_sdk.RemoteServer, err error)
func Create(params files_sdk.RemoteServerCreateParams, opts ...files_sdk.RequestResponseOption) (remoteServer files_sdk.RemoteServer, err error)
```

Parameters, `files_sdk.RemoteServerCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `Password` | `string` | `password` | no | Password, if needed. |
| `PrivateKey` | `string` | `private_key` | no | Private key, if needed. |
| `PrivateKeyPassphrase` | `string` | `private_key_passphrase` | no | Passphrase for private key if needed. |
| `ResetAuthentication` | `*bool` | `reset_authentication` | no | Reset authenticated account? |
| `SharepointClientCertificate` | `string` | `sharepoint_client_certificate` | no | SharePoint: PEM-encoded certificate and unencrypted private key for app-only authentication. |
| `SharepointClientSecret` | `string` | `sharepoint_client_secret` | no | SharePoint: Microsoft Entra application client secret for app-only authentication. |
| `SslCertificate` | `string` | `ssl_certificate` | no | SSL client certificate. |
| `AwsSecretKey` | `string` | `aws_secret_key` | no | AWS: secret key. |
| `AzureBlobStorageAccessKey` | `string` | `azure_blob_storage_access_key` | no | Azure Blob Storage: Access Key |
| `AzureBlobStorageSasToken` | `string` | `azure_blob_storage_sas_token` | no | Azure Blob Storage: Shared Access Signature (SAS) token |
| `AzureFilesStorageAccessKey` | `string` | `azure_files_storage_access_key` | no | Azure File Storage: Access Key |
| `AzureFilesStorageSasToken` | `string` | `azure_files_storage_sas_token` | no | Azure File Storage: Shared Access Signature (SAS) token |
| `BackblazeB2ApplicationKey` | `string` | `backblaze_b2_application_key` | no | Backblaze B2 Cloud Storage: applicationKey |
| `BackblazeB2KeyId` | `string` | `backblaze_b2_key_id` | no | Backblaze B2 Cloud Storage: keyID |
| `CloudflareSecretKey` | `string` | `cloudflare_secret_key` | no | Cloudflare: Secret Key |
| `FilebaseSecretKey` | `string` | `filebase_secret_key` | no | Filebase: Secret Key |
| `GoogleCloudStorageCredentialsJson` | `string` | `google_cloud_storage_credentials_json` | no | Google Cloud Storage: JSON file that contains the private key. To generate see https://cloud.google.com/storage/docs/json_api/v1/how-tos/authorizing#APIKey |
| `GoogleCloudStorageS3CompatibleSecretKey` | `string` | `google_cloud_storage_s3_compatible_secret_key` | no | Google Cloud Storage: S3-compatible secret key |
| `LinodeSecretKey` | `string` | `linode_secret_key` | no | Linode: Secret Key |
| `S3CompatibleSecretKey` | `string` | `s3_compatible_secret_key` | no | S3-compatible: Secret Key |
| `WasabiSecretKey` | `string` | `wasabi_secret_key` | no | Wasabi: Secret Key |
| `AllowRelativePaths` | `*bool` | `allow_relative_paths` | no | Allow relative paths in SFTP. If true, paths will not be forced to be absolute, allowing operations relative to the user's home directory. |
| `AwsAccessKey` | `string` | `aws_access_key` | no | AWS Access Key. |
| `AzureBlobStorageAccount` | `string` | `azure_blob_storage_account` | no | Azure Blob Storage: Account name |
| `AzureBlobStorageContainer` | `string` | `azure_blob_storage_container` | no | Azure Blob Storage: Container name |
| `AzureBlobStorageDnsSuffix` | `string` | `azure_blob_storage_dns_suffix` | no | Azure Blob Storage: Custom DNS suffix |
| `AzureBlobStorageHierarchicalNamespace` | `*bool` | `azure_blob_storage_hierarchical_namespace` | no | Azure Blob Storage: Does the storage account has hierarchical namespace feature enabled? |
| `AzureFilesStorageAccount` | `string` | `azure_files_storage_account` | no | Azure Files: Storage Account name |
| `AzureFilesStorageDnsSuffix` | `string` | `azure_files_storage_dns_suffix` | no | Azure Files: Custom DNS suffix |
| `AzureFilesStorageShareName` | `string` | `azure_files_storage_share_name` | no | Azure Files: Storage Share name |
| `BackblazeB2Bucket` | `string` | `backblaze_b2_bucket` | no | Backblaze B2 Cloud Storage: Bucket name |
| `BackblazeB2S3Endpoint` | `string` | `backblaze_b2_s3_endpoint` | no | Backblaze B2 Cloud Storage: S3 Endpoint |
| `BufferUploads` | `files_sdk.RemoteServerBufferUploadsEnum` | `buffer_uploads` | no | If set to always, uploads to this server will be uploaded first to Files.com before being sent to the remote server. This can improve performance in certain access patterns, such as high-latency connections. It will cause data to be temporarily stored in Files.com. If set to auto, we will perform this optimization if we believe it to be a benefit in a given situation. Values: `auto`, `always`, `never`. |
| `CloudflareAccessKey` | `string` | `cloudflare_access_key` | no | Cloudflare: Access Key. |
| `CloudflareBucket` | `string` | `cloudflare_bucket` | no | Cloudflare: Bucket name |
| `CloudflareEndpoint` | `string` | `cloudflare_endpoint` | no | Cloudflare: endpoint |
| `Description` | `string` | `description` | no | Internal description for your reference |
| `DropboxTeams` | `*bool` | `dropbox_teams` | no | Dropbox: If true, list Team folders in root? |
| `EnableDedicatedIps` | `*bool` | `enable_dedicated_ips` | no | `true` if remote server only accepts connections from dedicated IPs |
| `FilebaseAccessKey` | `string` | `filebase_access_key` | no | Filebase: Access Key. |
| `FilebaseBucket` | `string` | `filebase_bucket` | no | Filebase: Bucket name |
| `FilesApiKey` | `string` | `files_api_key` | no | Files.com direct link: API key used once to pair the remote server. |
| `FilesAgentPermissionSet` | `files_sdk.RemoteServerFilesAgentPermissionSetEnum` | `files_agent_permission_set` | no | Agent file permissions: `read_only`, `write_only`, or `read_write`. Files.com writes this value as `permission_set` in the configuration file generated for this Remote Server. Setting it on Create or Update changes the generated file, not a running Agent. The Agent enforces the permissions in its local configuration file, controlled by the host's system administrator, and reports them at registration, replacing the stored value. Values: `read_write`, `read_only`, `write_only`. |
| `FilesAgentRoot` | `string` | `files_agent_root` | no | Agent local root path. Files.com writes this value as `root` in the configuration file generated for this Remote Server. Setting it on Create or Update changes the generated file, not a running Agent. The Agent enforces the root in its local configuration file, controlled by the host's system administrator, and reports it at registration, replacing the stored value. |
| `FilesAgentVersion` | `string` | `files_agent_version` | no | Files.com Agent version reported by the Agent at each registration, replacing the stored value. This field is not included in the generated Agent configuration file. Setting it on Create or Update does not update the running Agent. |
| `OutboundAgentId` | `int64` | `outbound_agent_id` | no | Route traffic to outbound on a files-agent |
| `CustomDomainId` | `int64` | `custom_domain_id` | no | Custom Domain ID whose dedicated IP addresses are selected when this Remote Server uses dedicated IPs. Must be available to this Remote Server's workspace. Requires enable_dedicated_ips and cannot be combined with an outbound Agent. Set to null to use the site's default dedicated IPs. |
| `GoogleCloudStorageAuthenticationMethod` | `files_sdk.RemoteServerGoogleCloudStorageAuthenticationMethodEnum` | `google_cloud_storage_authentication_method` | no | Google Cloud Storage: Authentication method. Can be json, hmac, or oauth. Values: `json`, `hmac`, `oauth`. |
| `GoogleCloudStorageBucket` | `string` | `google_cloud_storage_bucket` | no | Google Cloud Storage: Bucket Name |
| `GoogleCloudStorageOauthScope` | `string` | `google_cloud_storage_oauth_scope` | no | Google Cloud Storage: OAuth scope. Can be https://www.googleapis.com/auth/devstorage.read_only or https://www.googleapis.com/auth/devstorage.read_write. |
| `GoogleCloudStorageProjectId` | `string` | `google_cloud_storage_project_id` | no | Google Cloud Storage: Project ID |
| `GoogleCloudStorageS3CompatibleAccessKey` | `string` | `google_cloud_storage_s3_compatible_access_key` | no | Google Cloud Storage: S3-compatible Access Key. |
| `Hostname` | `string` | `hostname` | no | Hostname or IP address. For Agent Remote Servers, the Agent reports the hostname it connects from at each registration, replacing the stored value. This field is not included in the generated Agent configuration file. Setting it on Create or Update does not change the running Agent. |
| `LinodeAccessKey` | `string` | `linode_access_key` | no | Linode: Access Key |
| `LinodeBucket` | `string` | `linode_bucket` | no | Linode: Bucket name |
| `LinodeRegion` | `string` | `linode_region` | no | Linode: region |
| `MaxConnections` | `int64` | `max_connections` | no | Max number of parallel connections. Ignored for S3 connections (we will parallelize these as much as possible). |
| `Name` | `string` | `name` | no | Internal name for your reference |
| `OneDriveAccountType` | `files_sdk.RemoteServerOneDriveAccountTypeEnum` | `one_drive_account_type` | no | OneDrive: Either personal or business_other account types Values: `personal`, `business_other`. |
| `PinToSiteRegion` | `*bool` | `pin_to_site_region` | no | If true, we will ensure that all communications with this remote server are made through the primary region of the site. This setting can also be overridden by a site-wide setting which will force it to true. |
| `Port` | `int64` | `port` | no | Port for remote server. For Agent Remote Servers, the Agent reports its downstream proxy port at each registration, replacing the stored value. This field is not included in the generated Agent configuration file. Setting it on Create or Update does not change the running Agent. |
| `UploadStagingPath` | `string` | `upload_staging_path` | no | Upload staging path. Applies to SFTP only. If a path is provided here, files will first be uploaded to this path on the remote folder and the moved into the final correct path via an SFTP move command. This is required by some remote MFT systems to emulate atomic uploads, which are otherwise not supoprted by SFTP. |
| `RemoteServerCredentialId` | `int64` | `remote_server_credential_id` | no | ID of Remote Server Credential, if applicable. |
| `S3AssumeRoleArn` | `string` | `s3_assume_role_arn` | no | AWS IAM Role ARN for AssumeRole authentication. |
| `S3AssumeRoleDurationSeconds` | `int64` | `s3_assume_role_duration_seconds` | no | Session duration in seconds for AssumeRole authentication (900-43200). |
| `S3Bucket` | `string` | `s3_bucket` | no | S3 bucket name |
| `S3CompatibleAccessKey` | `string` | `s3_compatible_access_key` | no | S3-compatible: Access Key |
| `S3CompatibleBucket` | `string` | `s3_compatible_bucket` | no | S3-compatible: Bucket name |
| `S3CompatibleEndpoint` | `string` | `s3_compatible_endpoint` | no | S3-compatible: endpoint |
| `S3CompatibleRegion` | `string` | `s3_compatible_region` | no | S3-compatible: region |
| `S3CompatibleVirtualHostedStyle` | `*bool` | `s3_compatible_virtual_hosted_style` | no | S3-compatible: If true, use virtual-hosted-style URLs instead of path-style URLs |
| `S3KmsKeyId` | `string` | `s3_kms_key_id` | no | ARN of the AWS KMS key that encrypts files written to the bucket (SSE-KMS). Leave blank unless your bucket requires a specific key. |
| `S3Region` | `string` | `s3_region` | no | S3 region |
| `ServerCertificate` | `files_sdk.RemoteServerServerCertificateEnum` | `server_certificate` | no | Remote server certificate Values: `require_match`, `allow_any`. |
| `ServerHostKey` | `string` | `server_host_key` | no | Pinned SSH host key or OpenSSH host certificate for SFTP. If omitted, Files.com detects and stores a host key, preferring plain keys over certificates. With `server_certificate=require_match` (the default), the server must present the exact pinned key or certificate and prove it holds the matching private key. A pinned certificate is compared in full, so renewal can require updating `server_host_key` even when its underlying key is unchanged. Files.com does not check certificate CA signatures, principals, or validity periods. Certificate expiration alone does not end the pin. Update `server_host_key` to replace the pin. |
| `ServerType` | `files_sdk.RemoteServerServerTypeEnum` | `server_type` | no | Remote server type. Values: `ftp`, `sftp`, `s3`, `google_cloud_storage`, `webdav`, `wasabi`, `backblaze_b2`, `one_drive`, `box`, `dropbox`, `google_drive`, `azure`, `sharepoint`, `s3_compatible`, `azure_files`, `files_agent`, `filebase`, `cloudflare`, `linode`, `files_com`. |
| `SharepointClientId` | `string` | `sharepoint_client_id` | no | SharePoint: Microsoft Entra application client ID for app-only authentication. |
| `SharepointSiteUrl` | `string` | `sharepoint_site_url` | no | SharePoint: Site URL to scope app-only authentication to a single site. Leave blank to browse all sites. |
| `SharepointTenantId` | `string` | `sharepoint_tenant_id` | no | SharePoint: Microsoft Entra tenant ID for app-only authentication. |
| `Ssl` | `files_sdk.RemoteServerSslEnum` | `ssl` | no | Should we require SSL? Values: `if_available`, `require`, `require_implicit`, `never`. |
| `Username` | `string` | `username` | no | Remote server username. |
| `WasabiAccessKey` | `string` | `wasabi_access_key` | no | Wasabi: Access Key. |
| `WasabiBucket` | `string` | `wasabi_bucket` | no | Wasabi: Bucket name |
| `WasabiRegion` | `string` | `wasabi_region` | no | Wasabi: Region |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID (0 for default workspace) |

Example:

```go
result, err := remote_server.Create(files_sdk.RemoteServerCreateParams{})
if err != nil {
	return err
}
fmt.Println(result)
```

### AgentPushUpdate

`POST /remote_servers/{id}/agent_push_update`

Push update to Files Agent

```go
func (c *Client) AgentPushUpdate(params files_sdk.RemoteServerAgentPushUpdateParams, opts ...files_sdk.RequestResponseOption) (agentPushUpdate files_sdk.AgentPushUpdate, err error)
func AgentPushUpdate(params files_sdk.RemoteServerAgentPushUpdateParams, opts ...files_sdk.RequestResponseOption) (agentPushUpdate files_sdk.AgentPushUpdate, err error)
func (c *Client) AgentPushUpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (agentPushUpdate files_sdk.AgentPushUpdate, err error)
func AgentPushUpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (agentPushUpdate files_sdk.AgentPushUpdate, err error)
```

`AgentPushUpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.RemoteServerAgentPushUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Remote Server ID. |

Example:

```go
result, err := remote_server.AgentPushUpdate(files_sdk.RemoteServerAgentPushUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := remote_server.AgentPushUpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /remote_servers/{id}`

Update Remote Server

```go
func (c *Client) Update(params files_sdk.RemoteServerUpdateParams, opts ...files_sdk.RequestResponseOption) (remoteServer files_sdk.RemoteServer, err error)
func Update(params files_sdk.RemoteServerUpdateParams, opts ...files_sdk.RequestResponseOption) (remoteServer files_sdk.RemoteServer, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (remoteServer files_sdk.RemoteServer, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (remoteServer files_sdk.RemoteServer, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.RemoteServerUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Remote Server ID. |
| `Password` | `string` | `password` | no | Password, if needed. |
| `PrivateKey` | `string` | `private_key` | no | Private key, if needed. |
| `PrivateKeyPassphrase` | `string` | `private_key_passphrase` | no | Passphrase for private key if needed. |
| `ResetAuthentication` | `*bool` | `reset_authentication` | no | Reset authenticated account? |
| `SharepointClientCertificate` | `string` | `sharepoint_client_certificate` | no | SharePoint: PEM-encoded certificate and unencrypted private key for app-only authentication. |
| `SharepointClientSecret` | `string` | `sharepoint_client_secret` | no | SharePoint: Microsoft Entra application client secret for app-only authentication. |
| `SslCertificate` | `string` | `ssl_certificate` | no | SSL client certificate. |
| `AwsSecretKey` | `string` | `aws_secret_key` | no | AWS: secret key. |
| `AzureBlobStorageAccessKey` | `string` | `azure_blob_storage_access_key` | no | Azure Blob Storage: Access Key |
| `AzureBlobStorageSasToken` | `string` | `azure_blob_storage_sas_token` | no | Azure Blob Storage: Shared Access Signature (SAS) token |
| `AzureFilesStorageAccessKey` | `string` | `azure_files_storage_access_key` | no | Azure File Storage: Access Key |
| `AzureFilesStorageSasToken` | `string` | `azure_files_storage_sas_token` | no | Azure File Storage: Shared Access Signature (SAS) token |
| `BackblazeB2ApplicationKey` | `string` | `backblaze_b2_application_key` | no | Backblaze B2 Cloud Storage: applicationKey |
| `BackblazeB2KeyId` | `string` | `backblaze_b2_key_id` | no | Backblaze B2 Cloud Storage: keyID |
| `CloudflareSecretKey` | `string` | `cloudflare_secret_key` | no | Cloudflare: Secret Key |
| `FilebaseSecretKey` | `string` | `filebase_secret_key` | no | Filebase: Secret Key |
| `GoogleCloudStorageCredentialsJson` | `string` | `google_cloud_storage_credentials_json` | no | Google Cloud Storage: JSON file that contains the private key. To generate see https://cloud.google.com/storage/docs/json_api/v1/how-tos/authorizing#APIKey |
| `GoogleCloudStorageS3CompatibleSecretKey` | `string` | `google_cloud_storage_s3_compatible_secret_key` | no | Google Cloud Storage: S3-compatible secret key |
| `LinodeSecretKey` | `string` | `linode_secret_key` | no | Linode: Secret Key |
| `S3CompatibleSecretKey` | `string` | `s3_compatible_secret_key` | no | S3-compatible: Secret Key |
| `WasabiSecretKey` | `string` | `wasabi_secret_key` | no | Wasabi: Secret Key |
| `AllowRelativePaths` | `*bool` | `allow_relative_paths` | no | Allow relative paths in SFTP. If true, paths will not be forced to be absolute, allowing operations relative to the user's home directory. |
| `AwsAccessKey` | `string` | `aws_access_key` | no | AWS Access Key. |
| `AzureBlobStorageAccount` | `string` | `azure_blob_storage_account` | no | Azure Blob Storage: Account name |
| `AzureBlobStorageContainer` | `string` | `azure_blob_storage_container` | no | Azure Blob Storage: Container name |
| `AzureBlobStorageDnsSuffix` | `string` | `azure_blob_storage_dns_suffix` | no | Azure Blob Storage: Custom DNS suffix |
| `AzureBlobStorageHierarchicalNamespace` | `*bool` | `azure_blob_storage_hierarchical_namespace` | no | Azure Blob Storage: Does the storage account has hierarchical namespace feature enabled? |
| `AzureFilesStorageAccount` | `string` | `azure_files_storage_account` | no | Azure Files: Storage Account name |
| `AzureFilesStorageDnsSuffix` | `string` | `azure_files_storage_dns_suffix` | no | Azure Files: Custom DNS suffix |
| `AzureFilesStorageShareName` | `string` | `azure_files_storage_share_name` | no | Azure Files: Storage Share name |
| `BackblazeB2Bucket` | `string` | `backblaze_b2_bucket` | no | Backblaze B2 Cloud Storage: Bucket name |
| `BackblazeB2S3Endpoint` | `string` | `backblaze_b2_s3_endpoint` | no | Backblaze B2 Cloud Storage: S3 Endpoint |
| `BufferUploads` | `files_sdk.RemoteServerBufferUploadsEnum` | `buffer_uploads` | no | If set to always, uploads to this server will be uploaded first to Files.com before being sent to the remote server. This can improve performance in certain access patterns, such as high-latency connections. It will cause data to be temporarily stored in Files.com. If set to auto, we will perform this optimization if we believe it to be a benefit in a given situation. Values: `auto`, `always`, `never`. |
| `CloudflareAccessKey` | `string` | `cloudflare_access_key` | no | Cloudflare: Access Key. |
| `CloudflareBucket` | `string` | `cloudflare_bucket` | no | Cloudflare: Bucket name |
| `CloudflareEndpoint` | `string` | `cloudflare_endpoint` | no | Cloudflare: endpoint |
| `Description` | `string` | `description` | no | Internal description for your reference |
| `DropboxTeams` | `*bool` | `dropbox_teams` | no | Dropbox: If true, list Team folders in root? |
| `EnableDedicatedIps` | `*bool` | `enable_dedicated_ips` | no | `true` if remote server only accepts connections from dedicated IPs |
| `FilebaseAccessKey` | `string` | `filebase_access_key` | no | Filebase: Access Key. |
| `FilebaseBucket` | `string` | `filebase_bucket` | no | Filebase: Bucket name |
| `FilesApiKey` | `string` | `files_api_key` | no | Files.com direct link: API key used once to pair the remote server. |
| `FilesAgentPermissionSet` | `files_sdk.RemoteServerFilesAgentPermissionSetEnum` | `files_agent_permission_set` | no | Agent file permissions: `read_only`, `write_only`, or `read_write`. Files.com writes this value as `permission_set` in the configuration file generated for this Remote Server. Setting it on Create or Update changes the generated file, not a running Agent. The Agent enforces the permissions in its local configuration file, controlled by the host's system administrator, and reports them at registration, replacing the stored value. Values: `read_write`, `read_only`, `write_only`. |
| `FilesAgentRoot` | `string` | `files_agent_root` | no | Agent local root path. Files.com writes this value as `root` in the configuration file generated for this Remote Server. Setting it on Create or Update changes the generated file, not a running Agent. The Agent enforces the root in its local configuration file, controlled by the host's system administrator, and reports it at registration, replacing the stored value. |
| `FilesAgentVersion` | `string` | `files_agent_version` | no | Files.com Agent version reported by the Agent at each registration, replacing the stored value. This field is not included in the generated Agent configuration file. Setting it on Create or Update does not update the running Agent. |
| `OutboundAgentId` | `int64` | `outbound_agent_id` | no | Route traffic to outbound on a files-agent |
| `CustomDomainId` | `int64` | `custom_domain_id` | no | Custom Domain ID whose dedicated IP addresses are selected when this Remote Server uses dedicated IPs. Must be available to this Remote Server's workspace. Requires enable_dedicated_ips and cannot be combined with an outbound Agent. Set to null to use the site's default dedicated IPs. |
| `GoogleCloudStorageAuthenticationMethod` | `files_sdk.RemoteServerGoogleCloudStorageAuthenticationMethodEnum` | `google_cloud_storage_authentication_method` | no | Google Cloud Storage: Authentication method. Can be json, hmac, or oauth. Values: `json`, `hmac`, `oauth`. |
| `GoogleCloudStorageBucket` | `string` | `google_cloud_storage_bucket` | no | Google Cloud Storage: Bucket Name |
| `GoogleCloudStorageOauthScope` | `string` | `google_cloud_storage_oauth_scope` | no | Google Cloud Storage: OAuth scope. Can be https://www.googleapis.com/auth/devstorage.read_only or https://www.googleapis.com/auth/devstorage.read_write. |
| `GoogleCloudStorageProjectId` | `string` | `google_cloud_storage_project_id` | no | Google Cloud Storage: Project ID |
| `GoogleCloudStorageS3CompatibleAccessKey` | `string` | `google_cloud_storage_s3_compatible_access_key` | no | Google Cloud Storage: S3-compatible Access Key. |
| `Hostname` | `string` | `hostname` | no | Hostname or IP address. For Agent Remote Servers, the Agent reports the hostname it connects from at each registration, replacing the stored value. This field is not included in the generated Agent configuration file. Setting it on Create or Update does not change the running Agent. |
| `LinodeAccessKey` | `string` | `linode_access_key` | no | Linode: Access Key |
| `LinodeBucket` | `string` | `linode_bucket` | no | Linode: Bucket name |
| `LinodeRegion` | `string` | `linode_region` | no | Linode: region |
| `MaxConnections` | `int64` | `max_connections` | no | Max number of parallel connections. Ignored for S3 connections (we will parallelize these as much as possible). |
| `Name` | `string` | `name` | no | Internal name for your reference |
| `OneDriveAccountType` | `files_sdk.RemoteServerOneDriveAccountTypeEnum` | `one_drive_account_type` | no | OneDrive: Either personal or business_other account types Values: `personal`, `business_other`. |
| `PinToSiteRegion` | `*bool` | `pin_to_site_region` | no | If true, we will ensure that all communications with this remote server are made through the primary region of the site. This setting can also be overridden by a site-wide setting which will force it to true. |
| `Port` | `int64` | `port` | no | Port for remote server. For Agent Remote Servers, the Agent reports its downstream proxy port at each registration, replacing the stored value. This field is not included in the generated Agent configuration file. Setting it on Create or Update does not change the running Agent. |
| `UploadStagingPath` | `string` | `upload_staging_path` | no | Upload staging path. Applies to SFTP only. If a path is provided here, files will first be uploaded to this path on the remote folder and the moved into the final correct path via an SFTP move command. This is required by some remote MFT systems to emulate atomic uploads, which are otherwise not supoprted by SFTP. |
| `RemoteServerCredentialId` | `int64` | `remote_server_credential_id` | no | ID of Remote Server Credential, if applicable. |
| `S3AssumeRoleArn` | `string` | `s3_assume_role_arn` | no | AWS IAM Role ARN for AssumeRole authentication. |
| `S3AssumeRoleDurationSeconds` | `int64` | `s3_assume_role_duration_seconds` | no | Session duration in seconds for AssumeRole authentication (900-43200). |
| `S3Bucket` | `string` | `s3_bucket` | no | S3 bucket name |
| `S3CompatibleAccessKey` | `string` | `s3_compatible_access_key` | no | S3-compatible: Access Key |
| `S3CompatibleBucket` | `string` | `s3_compatible_bucket` | no | S3-compatible: Bucket name |
| `S3CompatibleEndpoint` | `string` | `s3_compatible_endpoint` | no | S3-compatible: endpoint |
| `S3CompatibleRegion` | `string` | `s3_compatible_region` | no | S3-compatible: region |
| `S3CompatibleVirtualHostedStyle` | `*bool` | `s3_compatible_virtual_hosted_style` | no | S3-compatible: If true, use virtual-hosted-style URLs instead of path-style URLs |
| `S3KmsKeyId` | `string` | `s3_kms_key_id` | no | ARN of the AWS KMS key that encrypts files written to the bucket (SSE-KMS). Leave blank unless your bucket requires a specific key. |
| `S3Region` | `string` | `s3_region` | no | S3 region |
| `ServerCertificate` | `files_sdk.RemoteServerServerCertificateEnum` | `server_certificate` | no | Remote server certificate Values: `require_match`, `allow_any`. |
| `ServerHostKey` | `string` | `server_host_key` | no | Pinned SSH host key or OpenSSH host certificate for SFTP. If omitted, Files.com detects and stores a host key, preferring plain keys over certificates. With `server_certificate=require_match` (the default), the server must present the exact pinned key or certificate and prove it holds the matching private key. A pinned certificate is compared in full, so renewal can require updating `server_host_key` even when its underlying key is unchanged. Files.com does not check certificate CA signatures, principals, or validity periods. Certificate expiration alone does not end the pin. Update `server_host_key` to replace the pin. |
| `ServerType` | `files_sdk.RemoteServerServerTypeEnum` | `server_type` | no | Remote server type. Values: `ftp`, `sftp`, `s3`, `google_cloud_storage`, `webdav`, `wasabi`, `backblaze_b2`, `one_drive`, `box`, `dropbox`, `google_drive`, `azure`, `sharepoint`, `s3_compatible`, `azure_files`, `files_agent`, `filebase`, `cloudflare`, `linode`, `files_com`. |
| `SharepointClientId` | `string` | `sharepoint_client_id` | no | SharePoint: Microsoft Entra application client ID for app-only authentication. |
| `SharepointSiteUrl` | `string` | `sharepoint_site_url` | no | SharePoint: Site URL to scope app-only authentication to a single site. Leave blank to browse all sites. |
| `SharepointTenantId` | `string` | `sharepoint_tenant_id` | no | SharePoint: Microsoft Entra tenant ID for app-only authentication. |
| `Ssl` | `files_sdk.RemoteServerSslEnum` | `ssl` | no | Should we require SSL? Values: `if_available`, `require`, `require_implicit`, `never`. |
| `Username` | `string` | `username` | no | Remote server username. |
| `WasabiAccessKey` | `string` | `wasabi_access_key` | no | Wasabi: Access Key. |
| `WasabiBucket` | `string` | `wasabi_bucket` | no | Wasabi: Bucket name |
| `WasabiRegion` | `string` | `wasabi_region` | no | Wasabi: Region |

Example:

```go
result, err := remote_server.Update(files_sdk.RemoteServerUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := remote_server.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /remote_servers/{id}`

Delete Remote Server

```go
func (c *Client) Delete(params files_sdk.RemoteServerDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.RemoteServerDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.RemoteServerDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Remote Server ID. |

Example:

```go
err := remote_server.Delete(files_sdk.RemoteServerDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.RemoteServer` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Remote Server ID |
| `Disabled` | `*bool` | `disabled` | If true, this Remote Server is disabled. Updating it clears this flag, except for retired Agent v1 records, which remain disabled. |
| `AuthenticationMethod` | `string` | `authentication_method` | Type of authentication method to use |
| `Hostname` | `string` | `hostname` | Hostname or IP address. For Agent Remote Servers, the Agent reports the hostname it connects from at each registration, replacing the stored value. This field is not included in the generated Agent configuration file. Setting it on Create or Update does not change the running Agent. |
| `RemoteHomePath` | `string` | `remote_home_path` | Initial home folder on remote server |
| `UploadStagingPath` | `string` | `upload_staging_path` | Upload staging path. Applies to SFTP only. If a path is provided here, files will first be uploaded to this path on the remote folder and the moved into the final correct path via an SFTP move command. This is required by some remote MFT systems to emulate atomic uploads, which are otherwise not supoprted by SFTP. |
| `AllowRelativePaths` | `*bool` | `allow_relative_paths` | Allow relative paths in SFTP. If true, paths will not be forced to be absolute, allowing operations relative to the user's home directory. |
| `Name` | `string` | `name` | Internal name for your reference |
| `Description` | `string` | `description` | Internal description for your reference |
| `Port` | `int64` | `port` | Port for remote server. For Agent Remote Servers, the Agent reports its downstream proxy port at each registration, replacing the stored value. This field is not included in the generated Agent configuration file. Setting it on Create or Update does not change the running Agent. |
| `BufferUploads` | `string` | `buffer_uploads` | If set to always, uploads to this server will be uploaded first to Files.com before being sent to the remote server. This can improve performance in certain access patterns, such as high-latency connections. It will cause data to be temporarily stored in Files.com. If set to auto, we will perform this optimization if we believe it to be a benefit in a given situation. |
| `MaxConnections` | `int64` | `max_connections` | Max number of parallel connections. Ignored for S3 connections (we will parallelize these as much as possible). |
| `PinToSiteRegion` | `*bool` | `pin_to_site_region` | If true, we will ensure that all communications with this remote server are made through the primary region of the site. This setting can also be overridden by a site-wide setting which will force it to true. |
| `PinnedRegion` | `string` | `pinned_region` | If set, all communications with this remote server are made through the provided region. |
| `RemoteServerCredentialId` | `int64` | `remote_server_credential_id` | ID of Remote Server Credential, if applicable. |
| `S3Bucket` | `string` | `s3_bucket` | S3 bucket name |
| `S3Region` | `string` | `s3_region` | S3 region |
| `S3KmsKeyId` | `string` | `s3_kms_key_id` | ARN of the AWS KMS key that encrypts files written to the bucket (SSE-KMS). Leave blank unless your bucket requires a specific key. |
| `AwsAccessKey` | `string` | `aws_access_key` | AWS Access Key. |
| `S3AssumeRoleArn` | `string` | `s3_assume_role_arn` | AWS IAM Role ARN for AssumeRole authentication. |
| `S3AssumeRoleDurationSeconds` | `int64` | `s3_assume_role_duration_seconds` | Session duration in seconds for AssumeRole authentication (900-43200). |
| `S3AssumeRoleExternalId` | `string` | `s3_assume_role_external_id` | External ID for AssumeRole authentication. |
| `ServerCertificate` | `string` | `server_certificate` | Remote server certificate |
| `ServerHostKey` | `string` | `server_host_key` | Pinned SSH host key or OpenSSH host certificate for SFTP. If omitted, Files.com detects and stores a host key, preferring plain keys over certificates. With `server_certificate=require_match` (the default), the server must present the exact pinned key or certificate and prove it holds the matching private key. A pinned certificate is compared in full, so renewal can require updating `server_host_key` even when its underlying key is unchanged. Files.com does not check certificate CA signatures, principals, or validity periods. Certificate expiration alone does not end the pin. Update `server_host_key` to replace the pin. |
| `ServerType` | `string` | `server_type` | Remote server type. |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID (0 for default workspace) |
| `Ssl` | `string` | `ssl` | Should we require SSL? |
| `Username` | `string` | `username` | Remote server username. |
| `GoogleCloudStorageBucket` | `string` | `google_cloud_storage_bucket` | Google Cloud Storage: Bucket Name |
| `GoogleCloudStorageAuthenticationMethod` | `string` | `google_cloud_storage_authentication_method` | Google Cloud Storage: Authentication method. Can be json, hmac, or oauth. |
| `GoogleCloudStorageOauthScope` | `string` | `google_cloud_storage_oauth_scope` | Google Cloud Storage: OAuth scope. Can be https://www.googleapis.com/auth/devstorage.read_only or https://www.googleapis.com/auth/devstorage.read_write. |
| `GoogleCloudStorageProjectId` | `string` | `google_cloud_storage_project_id` | Google Cloud Storage: Project ID |
| `GoogleCloudStorageS3CompatibleAccessKey` | `string` | `google_cloud_storage_s3_compatible_access_key` | Google Cloud Storage: S3-compatible Access Key. |
| `BackblazeB2S3Endpoint` | `string` | `backblaze_b2_s3_endpoint` | Backblaze B2 Cloud Storage: S3 Endpoint |
| `BackblazeB2Bucket` | `string` | `backblaze_b2_bucket` | Backblaze B2 Cloud Storage: Bucket name |
| `WasabiBucket` | `string` | `wasabi_bucket` | Wasabi: Bucket name |
| `WasabiRegion` | `string` | `wasabi_region` | Wasabi: Region |
| `WasabiAccessKey` | `string` | `wasabi_access_key` | Wasabi: Access Key. |
| `AuthStatus` | `string` | `auth_status` | Either `in_setup` or `complete` |
| `AuthAccountName` | `string` | `auth_account_name` | Describes the authorized account |
| `OneDriveAccountType` | `string` | `one_drive_account_type` | OneDrive: Either personal or business_other account types |
| `SharepointTenantId` | `string` | `sharepoint_tenant_id` | SharePoint: Microsoft Entra tenant ID for app-only authentication. |
| `SharepointClientId` | `string` | `sharepoint_client_id` | SharePoint: Microsoft Entra application client ID for app-only authentication. |
| `SharepointAppAuthentication` | `*bool` | `sharepoint_app_authentication` | SharePoint: If true, this remote server uses Microsoft Entra app-only authentication. |
| `SharepointAppCredentialType` | `string` | `sharepoint_app_credential_type` | SharePoint: App-only credential type. Either secret or certificate. |
| `SharepointSiteUrl` | `string` | `sharepoint_site_url` | SharePoint: Site URL to scope app-only authentication to a single site. Leave blank to browse all sites. |
| `AzureBlobStorageAccount` | `string` | `azure_blob_storage_account` | Azure Blob Storage: Account name |
| `AzureBlobStorageContainer` | `string` | `azure_blob_storage_container` | Azure Blob Storage: Container name |
| `AzureBlobStorageHierarchicalNamespace` | `*bool` | `azure_blob_storage_hierarchical_namespace` | Azure Blob Storage: Does the storage account has hierarchical namespace feature enabled? |
| `AzureBlobStorageDnsSuffix` | `string` | `azure_blob_storage_dns_suffix` | Azure Blob Storage: Custom DNS suffix |
| `AzureFilesStorageAccount` | `string` | `azure_files_storage_account` | Azure Files: Storage Account name |
| `AzureFilesStorageShareName` | `string` | `azure_files_storage_share_name` | Azure Files: Storage Share name |
| `AzureFilesStorageDnsSuffix` | `string` | `azure_files_storage_dns_suffix` | Azure Files: Custom DNS suffix |
| `S3CompatibleBucket` | `string` | `s3_compatible_bucket` | S3-compatible: Bucket name |
| `S3CompatibleEndpoint` | `string` | `s3_compatible_endpoint` | S3-compatible: endpoint |
| `S3CompatibleRegion` | `string` | `s3_compatible_region` | S3-compatible: region |
| `S3CompatibleVirtualHostedStyle` | `*bool` | `s3_compatible_virtual_hosted_style` | S3-compatible: If true, use virtual-hosted-style URLs instead of path-style URLs |
| `S3CompatibleAccessKey` | `string` | `s3_compatible_access_key` | S3-compatible: Access Key |
| `EnableDedicatedIps` | `*bool` | `enable_dedicated_ips` | `true` if remote server only accepts connections from dedicated IPs |
| `CustomDomainId` | `int64` | `custom_domain_id` | Custom Domain ID whose dedicated IP addresses are selected when this Remote Server uses dedicated IPs. Must be available to this Remote Server's workspace. Requires enable_dedicated_ips and cannot be combined with an outbound Agent. Set to null to use the site's default dedicated IPs. |
| `OutboundIpAddresses` | `[]string` | `outbound_ip_addresses` | Current eligible public IP addresses for the selected Custom Domain. Any address in this list may originate a connection. Empty when no domain is selected or its configuration is unavailable. Only included in responses for a single Remote Server. |
| `FilesAgentPermissionSet` | `string` | `files_agent_permission_set` | Agent file permissions: `read_only`, `write_only`, or `read_write`. Files.com writes this value as `permission_set` in the configuration file generated for this Remote Server. Setting it on Create or Update changes the generated file, not a running Agent. The Agent enforces the permissions in its local configuration file, controlled by the host's system administrator, and reports them at registration, replacing the stored value. |
| `FilesAgentRoot` | `string` | `files_agent_root` | Agent local root path. Files.com writes this value as `root` in the configuration file generated for this Remote Server. Setting it on Create or Update changes the generated file, not a running Agent. The Agent enforces the root in its local configuration file, controlled by the host's system administrator, and reports it at registration, replacing the stored value. |
| `FilesAgentVersion` | `string` | `files_agent_version` | Files.com Agent version reported by the Agent at each registration, replacing the stored value. This field is not included in the generated Agent configuration file. Setting it on Create or Update does not update the running Agent. |
| `FilesAgentUpToDate` | `*bool` | `files_agent_up_to_date` | If true, the Files Agent is up to date. |
| `FilesAgentLatestVersion` | `string` | `files_agent_latest_version` | Latest available Files Agent version |
| `FilesAgentSupportsPushUpdates` | `*bool` | `files_agent_supports_push_updates` | Files Agent supports receiving push updates |
| `DirectTransferAvailable` | `*bool` | `direct_transfer_available` | Whether the Files Agent Proxy recently validated a direct transfer connection. `true` means a direct connection was recently validated (actual availability can vary by client network), `false` means direct transfers are enabled but not currently available, and `null` means direct transfers are disabled or unsupported. Only provided for a connected Files Agent when showing a single Remote Server. |
| `OutboundAgentId` | `int64` | `outbound_agent_id` | Route traffic to outbound on a files-agent |
| `FilebaseBucket` | `string` | `filebase_bucket` | Filebase: Bucket name |
| `FilebaseAccessKey` | `string` | `filebase_access_key` | Filebase: Access Key. |
| `FilesApiKeyPrefix` | `string` | `files_api_key_prefix` | Files.com direct link: paired API key prefix. |
| `CloudflareBucket` | `string` | `cloudflare_bucket` | Cloudflare: Bucket name |
| `CloudflareAccessKey` | `string` | `cloudflare_access_key` | Cloudflare: Access Key. |
| `CloudflareEndpoint` | `string` | `cloudflare_endpoint` | Cloudflare: endpoint |
| `DropboxTeams` | `*bool` | `dropbox_teams` | Dropbox: If true, list Team folders in root? |
| `LinodeBucket` | `string` | `linode_bucket` | Linode: Bucket name |
| `LinodeAccessKey` | `string` | `linode_access_key` | Linode: Access Key |
| `LinodeRegion` | `string` | `linode_region` | Linode: region |
| `SupportsVersioning` | `*bool` | `supports_versioning` | If true, this remote server supports file versioning. This value is determined automatically by Files.com. |
| `UserId` | `int64` | `user_id` | User ID. Provide a value of `0` to operate the current session's user. |
| `Password` | `string` | `password` | Password, if needed. |
| `PrivateKey` | `string` | `private_key` | Private key, if needed. |
| `PrivateKeyPassphrase` | `string` | `private_key_passphrase` | Passphrase for private key if needed. |
| `ResetAuthentication` | `*bool` | `reset_authentication` | Reset authenticated account? |
| `SharepointClientCertificate` | `string` | `sharepoint_client_certificate` | SharePoint: PEM-encoded certificate and unencrypted private key for app-only authentication. |
| `SharepointClientSecret` | `string` | `sharepoint_client_secret` | SharePoint: Microsoft Entra application client secret for app-only authentication. |
| `SslCertificate` | `string` | `ssl_certificate` | SSL client certificate. |
| `AwsSecretKey` | `string` | `aws_secret_key` | AWS: secret key. |
| `AzureBlobStorageAccessKey` | `string` | `azure_blob_storage_access_key` | Azure Blob Storage: Access Key |
| `AzureBlobStorageSasToken` | `string` | `azure_blob_storage_sas_token` | Azure Blob Storage: Shared Access Signature (SAS) token |
| `AzureFilesStorageAccessKey` | `string` | `azure_files_storage_access_key` | Azure File Storage: Access Key |
| `AzureFilesStorageSasToken` | `string` | `azure_files_storage_sas_token` | Azure File Storage: Shared Access Signature (SAS) token |
| `BackblazeB2ApplicationKey` | `string` | `backblaze_b2_application_key` | Backblaze B2 Cloud Storage: applicationKey |
| `BackblazeB2KeyId` | `string` | `backblaze_b2_key_id` | Backblaze B2 Cloud Storage: keyID |
| `CloudflareSecretKey` | `string` | `cloudflare_secret_key` | Cloudflare: Secret Key |
| `FilebaseSecretKey` | `string` | `filebase_secret_key` | Filebase: Secret Key |
| `GoogleCloudStorageCredentialsJson` | `string` | `google_cloud_storage_credentials_json` | Google Cloud Storage: JSON file that contains the private key. To generate see https://cloud.google.com/storage/docs/json_api/v1/how-tos/authorizing#APIKey |
| `GoogleCloudStorageS3CompatibleSecretKey` | `string` | `google_cloud_storage_s3_compatible_secret_key` | Google Cloud Storage: S3-compatible secret key |
| `LinodeSecretKey` | `string` | `linode_secret_key` | Linode: Secret Key |
| `S3CompatibleSecretKey` | `string` | `s3_compatible_secret_key` | S3-compatible: Secret Key |
| `WasabiSecretKey` | `string` | `wasabi_secret_key` | Wasabi: Secret Key |
| `FilesApiKey` | `string` | `files_api_key` | Files.com direct link: API key used once to pair the remote server. |
