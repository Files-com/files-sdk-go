package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// RemoteServer is a Files.com API resource.
type RemoteServer struct {
	// Remote Server ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// If true, this Remote Server is disabled. Updating it clears this flag, except for retired Agent v1 records, which remain disabled.
	Disabled *bool `json:"disabled,omitempty" path:"disabled,omitempty" url:"disabled,omitempty"`
	// Type of authentication method to use.
	AuthenticationMethod string `json:"authentication_method,omitempty" path:"authentication_method,omitempty" url:"authentication_method,omitempty"`
	// Hostname or IP address. For Agent Remote Servers, the Agent reports the hostname it connects from at each registration, replacing the stored value. This field is not included in the generated Agent configuration file. Setting it on Create or Update does not change the running Agent.
	Hostname string `json:"hostname,omitempty" path:"hostname,omitempty" url:"hostname,omitempty"`
	// Initial home folder on remote server.
	RemoteHomePath string `json:"remote_home_path,omitempty" path:"remote_home_path,omitempty" url:"remote_home_path,omitempty"`
	// Upload staging path. Applies to SFTP only. If a path is provided here, files will first be uploaded to this path on the remote folder and the moved into the final correct path via an SFTP move command. This is required by some remote MFT systems to emulate atomic uploads, which are otherwise not supoprted by SFTP.
	UploadStagingPath string `json:"upload_staging_path,omitempty" path:"upload_staging_path,omitempty" url:"upload_staging_path,omitempty"`
	// Allow relative paths in SFTP. If true, paths will not be forced to be absolute, allowing operations relative to the user's home directory.
	AllowRelativePaths *bool `json:"allow_relative_paths,omitempty" path:"allow_relative_paths,omitempty" url:"allow_relative_paths,omitempty"`
	// Internal name for your reference.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Internal description for your reference.
	Description string `json:"description,omitempty" path:"description,omitempty" url:"description,omitempty"`
	// Port for remote server. For Agent Remote Servers, the Agent reports its downstream proxy port at each registration, replacing the stored value. This field is not included in the generated Agent configuration file. Setting it on Create or Update does not change the running Agent.
	Port int64 `json:"port,omitempty" path:"port,omitempty" url:"port,omitempty"`
	// If set to always, uploads to this server will be uploaded first to Files.com before being sent to the remote server. This can improve performance in certain access patterns, such as high-latency connections. It will cause data to be temporarily stored in Files.com. If set to auto, we will perform this optimization if we believe it to be a benefit in a given situation.
	BufferUploads string `json:"buffer_uploads,omitempty" path:"buffer_uploads,omitempty" url:"buffer_uploads,omitempty"`
	// Max number of parallel connections. Ignored for S3 connections (we will parallelize these as much as possible).
	MaxConnections int64 `json:"max_connections,omitempty" path:"max_connections,omitempty" url:"max_connections,omitempty"`
	// If true, we will ensure that all communications with this remote server are made through the primary region of the site. This setting can also be overridden by a site-wide setting which will force it to true.
	PinToSiteRegion *bool `json:"pin_to_site_region,omitempty" path:"pin_to_site_region,omitempty" url:"pin_to_site_region,omitempty"`
	// If set, all communications with this remote server are made through the provided region.
	PinnedRegion string `json:"pinned_region,omitempty" path:"pinned_region,omitempty" url:"pinned_region,omitempty"`
	// ID of Remote Server Credential, if applicable.
	RemoteServerCredentialId int64 `json:"remote_server_credential_id,omitempty" path:"remote_server_credential_id,omitempty" url:"remote_server_credential_id,omitempty"`
	// S3 bucket name.
	S3Bucket string `json:"s3_bucket,omitempty" path:"s3_bucket,omitempty" url:"s3_bucket,omitempty"`
	// S3 region.
	S3Region string `json:"s3_region,omitempty" path:"s3_region,omitempty" url:"s3_region,omitempty"`
	// ARN of the AWS KMS key that encrypts files written to the bucket (SSE-KMS). Leave blank unless your bucket requires a specific key.
	S3KmsKeyId string `json:"s3_kms_key_id,omitempty" path:"s3_kms_key_id,omitempty" url:"s3_kms_key_id,omitempty"`
	// AWS Access Key.
	AwsAccessKey string `json:"aws_access_key,omitempty" path:"aws_access_key,omitempty" url:"aws_access_key,omitempty"`
	// AWS IAM Role ARN for AssumeRole authentication.
	S3AssumeRoleArn string `json:"s3_assume_role_arn,omitempty" path:"s3_assume_role_arn,omitempty" url:"s3_assume_role_arn,omitempty"`
	// Session duration in seconds for AssumeRole authentication (900-43200).
	S3AssumeRoleDurationSeconds int64 `json:"s3_assume_role_duration_seconds,omitempty" path:"s3_assume_role_duration_seconds,omitempty" url:"s3_assume_role_duration_seconds,omitempty"`
	// External ID for AssumeRole authentication.
	S3AssumeRoleExternalId string `json:"s3_assume_role_external_id,omitempty" path:"s3_assume_role_external_id,omitempty" url:"s3_assume_role_external_id,omitempty"`
	// Remote server certificate.
	ServerCertificate string `json:"server_certificate,omitempty" path:"server_certificate,omitempty" url:"server_certificate,omitempty"`
	// Pinned SSH host key or OpenSSH host certificate for SFTP. If omitted, Files.com detects and stores a host key, preferring plain keys over certificates. With `server_certificate=require_match` (the default), the server must present the exact pinned key or certificate and prove it holds the matching private key. A pinned certificate is compared in full, so renewal can require updating `server_host_key` even when its underlying key is unchanged. Files.com does not check certificate CA signatures, principals, or validity periods. Certificate expiration alone does not end the pin. Update `server_host_key` to replace the pin.
	ServerHostKey string `json:"server_host_key,omitempty" path:"server_host_key,omitempty" url:"server_host_key,omitempty"`
	// Remote server type.
	ServerType string `json:"server_type,omitempty" path:"server_type,omitempty" url:"server_type,omitempty"`
	// Workspace ID (0 for default workspace).
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Should we require SSL?
	Ssl string `json:"ssl,omitempty" path:"ssl,omitempty" url:"ssl,omitempty"`
	// Remote server username.
	Username string `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	// Google Cloud Storage: Bucket Name.
	GoogleCloudStorageBucket string `json:"google_cloud_storage_bucket,omitempty" path:"google_cloud_storage_bucket,omitempty" url:"google_cloud_storage_bucket,omitempty"`
	// Google Cloud Storage: Authentication method. Can be json, hmac, or oauth.
	GoogleCloudStorageAuthenticationMethod string `json:"google_cloud_storage_authentication_method,omitempty" path:"google_cloud_storage_authentication_method,omitempty" url:"google_cloud_storage_authentication_method,omitempty"`
	// Google Cloud Storage: OAuth scope. Can be https://www.googleapis.com/auth/devstorage.read_only or https://www.googleapis.com/auth/devstorage.read_write.
	GoogleCloudStorageOauthScope string `json:"google_cloud_storage_oauth_scope,omitempty" path:"google_cloud_storage_oauth_scope,omitempty" url:"google_cloud_storage_oauth_scope,omitempty"`
	// Google Cloud Storage: Project ID.
	GoogleCloudStorageProjectId string `json:"google_cloud_storage_project_id,omitempty" path:"google_cloud_storage_project_id,omitempty" url:"google_cloud_storage_project_id,omitempty"`
	// Google Cloud Storage: S3-compatible Access Key.
	GoogleCloudStorageS3CompatibleAccessKey string `json:"google_cloud_storage_s3_compatible_access_key,omitempty" path:"google_cloud_storage_s3_compatible_access_key,omitempty" url:"google_cloud_storage_s3_compatible_access_key,omitempty"`
	// Backblaze B2 Cloud Storage: S3 Endpoint.
	BackblazeB2S3Endpoint string `json:"backblaze_b2_s3_endpoint,omitempty" path:"backblaze_b2_s3_endpoint,omitempty" url:"backblaze_b2_s3_endpoint,omitempty"`
	// Backblaze B2 Cloud Storage: Bucket name.
	BackblazeB2Bucket string `json:"backblaze_b2_bucket,omitempty" path:"backblaze_b2_bucket,omitempty" url:"backblaze_b2_bucket,omitempty"`
	// Wasabi: Bucket name.
	WasabiBucket string `json:"wasabi_bucket,omitempty" path:"wasabi_bucket,omitempty" url:"wasabi_bucket,omitempty"`
	// Wasabi: Region.
	WasabiRegion string `json:"wasabi_region,omitempty" path:"wasabi_region,omitempty" url:"wasabi_region,omitempty"`
	// Wasabi: Access Key.
	WasabiAccessKey string `json:"wasabi_access_key,omitempty" path:"wasabi_access_key,omitempty" url:"wasabi_access_key,omitempty"`
	// Either `in_setup` or `complete`.
	AuthStatus string `json:"auth_status,omitempty" path:"auth_status,omitempty" url:"auth_status,omitempty"`
	// Describes the authorized account.
	AuthAccountName string `json:"auth_account_name,omitempty" path:"auth_account_name,omitempty" url:"auth_account_name,omitempty"`
	// OneDrive: Either personal or business_other account types.
	OneDriveAccountType string `json:"one_drive_account_type,omitempty" path:"one_drive_account_type,omitempty" url:"one_drive_account_type,omitempty"`
	// SharePoint: Microsoft Entra tenant ID for app-only authentication.
	SharepointTenantId string `json:"sharepoint_tenant_id,omitempty" path:"sharepoint_tenant_id,omitempty" url:"sharepoint_tenant_id,omitempty"`
	// SharePoint: Microsoft Entra application client ID for app-only authentication.
	SharepointClientId string `json:"sharepoint_client_id,omitempty" path:"sharepoint_client_id,omitempty" url:"sharepoint_client_id,omitempty"`
	// SharePoint: If true, this remote server uses Microsoft Entra app-only authentication.
	SharepointAppAuthentication *bool `json:"sharepoint_app_authentication,omitempty" path:"sharepoint_app_authentication,omitempty" url:"sharepoint_app_authentication,omitempty"`
	// SharePoint: App-only credential type. Either secret or certificate.
	SharepointAppCredentialType string `json:"sharepoint_app_credential_type,omitempty" path:"sharepoint_app_credential_type,omitempty" url:"sharepoint_app_credential_type,omitempty"`
	// SharePoint: Site URL to scope app-only authentication to a single site. Leave blank to browse all sites.
	SharepointSiteUrl string `json:"sharepoint_site_url,omitempty" path:"sharepoint_site_url,omitempty" url:"sharepoint_site_url,omitempty"`
	// Azure Blob Storage: Account name.
	AzureBlobStorageAccount string `json:"azure_blob_storage_account,omitempty" path:"azure_blob_storage_account,omitempty" url:"azure_blob_storage_account,omitempty"`
	// Azure Blob Storage: Container name.
	AzureBlobStorageContainer string `json:"azure_blob_storage_container,omitempty" path:"azure_blob_storage_container,omitempty" url:"azure_blob_storage_container,omitempty"`
	// Azure Blob Storage: Does the storage account has hierarchical namespace feature enabled?
	AzureBlobStorageHierarchicalNamespace *bool `json:"azure_blob_storage_hierarchical_namespace,omitempty" path:"azure_blob_storage_hierarchical_namespace,omitempty" url:"azure_blob_storage_hierarchical_namespace,omitempty"`
	// Azure Blob Storage: Custom DNS suffix.
	AzureBlobStorageDnsSuffix string `json:"azure_blob_storage_dns_suffix,omitempty" path:"azure_blob_storage_dns_suffix,omitempty" url:"azure_blob_storage_dns_suffix,omitempty"`
	// Azure Files: Storage Account name.
	AzureFilesStorageAccount string `json:"azure_files_storage_account,omitempty" path:"azure_files_storage_account,omitempty" url:"azure_files_storage_account,omitempty"`
	// Azure Files: Storage Share name.
	AzureFilesStorageShareName string `json:"azure_files_storage_share_name,omitempty" path:"azure_files_storage_share_name,omitempty" url:"azure_files_storage_share_name,omitempty"`
	// Azure Files: Custom DNS suffix.
	AzureFilesStorageDnsSuffix string `json:"azure_files_storage_dns_suffix,omitempty" path:"azure_files_storage_dns_suffix,omitempty" url:"azure_files_storage_dns_suffix,omitempty"`
	// S3-compatible: Bucket name.
	S3CompatibleBucket string `json:"s3_compatible_bucket,omitempty" path:"s3_compatible_bucket,omitempty" url:"s3_compatible_bucket,omitempty"`
	// S3-compatible: endpoint.
	S3CompatibleEndpoint string `json:"s3_compatible_endpoint,omitempty" path:"s3_compatible_endpoint,omitempty" url:"s3_compatible_endpoint,omitempty"`
	// S3-compatible: region.
	S3CompatibleRegion string `json:"s3_compatible_region,omitempty" path:"s3_compatible_region,omitempty" url:"s3_compatible_region,omitempty"`
	// S3-compatible: If true, use virtual-hosted-style URLs instead of path-style URLs.
	S3CompatibleVirtualHostedStyle *bool `json:"s3_compatible_virtual_hosted_style,omitempty" path:"s3_compatible_virtual_hosted_style,omitempty" url:"s3_compatible_virtual_hosted_style,omitempty"`
	// S3-compatible: Access Key.
	S3CompatibleAccessKey string `json:"s3_compatible_access_key,omitempty" path:"s3_compatible_access_key,omitempty" url:"s3_compatible_access_key,omitempty"`
	// `true` if remote server only accepts connections from dedicated IPs.
	EnableDedicatedIps *bool `json:"enable_dedicated_ips,omitempty" path:"enable_dedicated_ips,omitempty" url:"enable_dedicated_ips,omitempty"`
	// Custom Domain ID whose dedicated IP addresses are selected when this Remote Server uses dedicated IPs. Must be available to this Remote Server's workspace. Requires enable_dedicated_ips and cannot be combined with an outbound Agent. Set to null to use the site's default dedicated IPs.
	CustomDomainId int64 `json:"custom_domain_id,omitempty" path:"custom_domain_id,omitempty" url:"custom_domain_id,omitempty"`
	// Current eligible public IP addresses for the selected Custom Domain. Any address in this list may originate a connection. Empty when no domain is selected or its configuration is unavailable. Only included in responses for a single Remote Server.
	OutboundIpAddresses []string `json:"outbound_ip_addresses,omitempty" path:"outbound_ip_addresses,omitempty" url:"outbound_ip_addresses,omitempty"`
	// Agent file permissions: `read_only`, `write_only`, or `read_write`. Files.com writes this value as `permission_set` in the configuration file generated for this Remote Server. Setting it on Create or Update changes the generated file, not a running Agent. The Agent enforces the permissions in its local configuration file, controlled by the host's system administrator, and reports them at registration, replacing the stored value.
	FilesAgentPermissionSet string `json:"files_agent_permission_set,omitempty" path:"files_agent_permission_set,omitempty" url:"files_agent_permission_set,omitempty"`
	// Agent local root path. Files.com writes this value as `root` in the configuration file generated for this Remote Server. Setting it on Create or Update changes the generated file, not a running Agent. The Agent enforces the root in its local configuration file, controlled by the host's system administrator, and reports it at registration, replacing the stored value.
	FilesAgentRoot string `json:"files_agent_root,omitempty" path:"files_agent_root,omitempty" url:"files_agent_root,omitempty"`
	// Files.com Agent version reported by the Agent at each registration, replacing the stored value. This field is not included in the generated Agent configuration file. Setting it on Create or Update does not update the running Agent.
	FilesAgentVersion string `json:"files_agent_version,omitempty" path:"files_agent_version,omitempty" url:"files_agent_version,omitempty"`
	// If true, the Files Agent is up to date.
	FilesAgentUpToDate *bool `json:"files_agent_up_to_date,omitempty" path:"files_agent_up_to_date,omitempty" url:"files_agent_up_to_date,omitempty"`
	// Latest available Files Agent version.
	FilesAgentLatestVersion string `json:"files_agent_latest_version,omitempty" path:"files_agent_latest_version,omitempty" url:"files_agent_latest_version,omitempty"`
	// Files Agent supports receiving push updates.
	FilesAgentSupportsPushUpdates *bool `json:"files_agent_supports_push_updates,omitempty" path:"files_agent_supports_push_updates,omitempty" url:"files_agent_supports_push_updates,omitempty"`
	// Whether the Files Agent Proxy recently validated a direct transfer connection. `true` means a direct connection was recently validated (actual availability can vary by client network), `false` means direct transfers are enabled but not currently available, and `null` means direct transfers are disabled or unsupported. Only provided for a connected Files Agent when showing a single Remote Server.
	DirectTransferAvailable *bool `json:"direct_transfer_available,omitempty" path:"direct_transfer_available,omitempty" url:"direct_transfer_available,omitempty"`
	// Route traffic to outbound on a files-agent.
	OutboundAgentId int64 `json:"outbound_agent_id,omitempty" path:"outbound_agent_id,omitempty" url:"outbound_agent_id,omitempty"`
	// Filebase: Bucket name.
	FilebaseBucket string `json:"filebase_bucket,omitempty" path:"filebase_bucket,omitempty" url:"filebase_bucket,omitempty"`
	// Filebase: Access Key.
	FilebaseAccessKey string `json:"filebase_access_key,omitempty" path:"filebase_access_key,omitempty" url:"filebase_access_key,omitempty"`
	// Files.com direct link: paired API key prefix.
	FilesApiKeyPrefix string `json:"files_api_key_prefix,omitempty" path:"files_api_key_prefix,omitempty" url:"files_api_key_prefix,omitempty"`
	// Cloudflare: Bucket name.
	CloudflareBucket string `json:"cloudflare_bucket,omitempty" path:"cloudflare_bucket,omitempty" url:"cloudflare_bucket,omitempty"`
	// Cloudflare: Access Key.
	CloudflareAccessKey string `json:"cloudflare_access_key,omitempty" path:"cloudflare_access_key,omitempty" url:"cloudflare_access_key,omitempty"`
	// Cloudflare: endpoint.
	CloudflareEndpoint string `json:"cloudflare_endpoint,omitempty" path:"cloudflare_endpoint,omitempty" url:"cloudflare_endpoint,omitempty"`
	// Dropbox: If true, list Team folders in root?
	DropboxTeams *bool `json:"dropbox_teams,omitempty" path:"dropbox_teams,omitempty" url:"dropbox_teams,omitempty"`
	// Linode: Bucket name.
	LinodeBucket string `json:"linode_bucket,omitempty" path:"linode_bucket,omitempty" url:"linode_bucket,omitempty"`
	// Linode: Access Key.
	LinodeAccessKey string `json:"linode_access_key,omitempty" path:"linode_access_key,omitempty" url:"linode_access_key,omitempty"`
	// Linode: region.
	LinodeRegion string `json:"linode_region,omitempty" path:"linode_region,omitempty" url:"linode_region,omitempty"`
	// If true, this remote server supports file versioning. This value is determined automatically by Files.com.
	SupportsVersioning *bool `json:"supports_versioning,omitempty" path:"supports_versioning,omitempty" url:"supports_versioning,omitempty"`
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	// Password, if needed.
	Password string `json:"password,omitempty" path:"password,omitempty" url:"password,omitempty"`
	// Private key, if needed.
	PrivateKey string `json:"private_key,omitempty" path:"private_key,omitempty" url:"private_key,omitempty"`
	// Passphrase for private key if needed.
	PrivateKeyPassphrase string `json:"private_key_passphrase,omitempty" path:"private_key_passphrase,omitempty" url:"private_key_passphrase,omitempty"`
	// Reset authenticated account?
	ResetAuthentication *bool `json:"reset_authentication,omitempty" path:"reset_authentication,omitempty" url:"reset_authentication,omitempty"`
	// SharePoint: PEM-encoded certificate and unencrypted private key for app-only authentication.
	SharepointClientCertificate string `json:"sharepoint_client_certificate,omitempty" path:"sharepoint_client_certificate,omitempty" url:"sharepoint_client_certificate,omitempty"`
	// SharePoint: Microsoft Entra application client secret for app-only authentication.
	SharepointClientSecret string `json:"sharepoint_client_secret,omitempty" path:"sharepoint_client_secret,omitempty" url:"sharepoint_client_secret,omitempty"`
	// SSL client certificate.
	SslCertificate string `json:"ssl_certificate,omitempty" path:"ssl_certificate,omitempty" url:"ssl_certificate,omitempty"`
	// AWS: secret key.
	AwsSecretKey string `json:"aws_secret_key,omitempty" path:"aws_secret_key,omitempty" url:"aws_secret_key,omitempty"`
	// Azure Blob Storage: Access Key.
	AzureBlobStorageAccessKey string `json:"azure_blob_storage_access_key,omitempty" path:"azure_blob_storage_access_key,omitempty" url:"azure_blob_storage_access_key,omitempty"`
	// Azure Blob Storage: Shared Access Signature (SAS) token.
	AzureBlobStorageSasToken string `json:"azure_blob_storage_sas_token,omitempty" path:"azure_blob_storage_sas_token,omitempty" url:"azure_blob_storage_sas_token,omitempty"`
	// Azure File Storage: Access Key.
	AzureFilesStorageAccessKey string `json:"azure_files_storage_access_key,omitempty" path:"azure_files_storage_access_key,omitempty" url:"azure_files_storage_access_key,omitempty"`
	// Azure File Storage: Shared Access Signature (SAS) token.
	AzureFilesStorageSasToken string `json:"azure_files_storage_sas_token,omitempty" path:"azure_files_storage_sas_token,omitempty" url:"azure_files_storage_sas_token,omitempty"`
	// Backblaze B2 Cloud Storage: applicationKey.
	BackblazeB2ApplicationKey string `json:"backblaze_b2_application_key,omitempty" path:"backblaze_b2_application_key,omitempty" url:"backblaze_b2_application_key,omitempty"`
	// Backblaze B2 Cloud Storage: keyID.
	BackblazeB2KeyId string `json:"backblaze_b2_key_id,omitempty" path:"backblaze_b2_key_id,omitempty" url:"backblaze_b2_key_id,omitempty"`
	// Cloudflare: Secret Key.
	CloudflareSecretKey string `json:"cloudflare_secret_key,omitempty" path:"cloudflare_secret_key,omitempty" url:"cloudflare_secret_key,omitempty"`
	// Filebase: Secret Key.
	FilebaseSecretKey string `json:"filebase_secret_key,omitempty" path:"filebase_secret_key,omitempty" url:"filebase_secret_key,omitempty"`
	// Google Cloud Storage: JSON file that contains the private key. To generate see https://cloud.google.com/storage/docs/json_api/v1/how-tos/authorizing#APIKey.
	GoogleCloudStorageCredentialsJson string `json:"google_cloud_storage_credentials_json,omitempty" path:"google_cloud_storage_credentials_json,omitempty" url:"google_cloud_storage_credentials_json,omitempty"`
	// Google Cloud Storage: S3-compatible secret key.
	GoogleCloudStorageS3CompatibleSecretKey string `json:"google_cloud_storage_s3_compatible_secret_key,omitempty" path:"google_cloud_storage_s3_compatible_secret_key,omitempty" url:"google_cloud_storage_s3_compatible_secret_key,omitempty"`
	// Linode: Secret Key.
	LinodeSecretKey string `json:"linode_secret_key,omitempty" path:"linode_secret_key,omitempty" url:"linode_secret_key,omitempty"`
	// S3-compatible: Secret Key.
	S3CompatibleSecretKey string `json:"s3_compatible_secret_key,omitempty" path:"s3_compatible_secret_key,omitempty" url:"s3_compatible_secret_key,omitempty"`
	// Wasabi: Secret Key.
	WasabiSecretKey string `json:"wasabi_secret_key,omitempty" path:"wasabi_secret_key,omitempty" url:"wasabi_secret_key,omitempty"`
	// Files.com direct link: API key used once to pair the remote server.
	FilesApiKey string `json:"files_api_key,omitempty" path:"files_api_key,omitempty" url:"files_api_key,omitempty"`
}

// Identifier returns the resource ID.
func (r RemoteServer) Identifier() interface{} {
	return r.Id
}

// RemoteServerCollection is a list of RemoteServer resources.
type RemoteServerCollection []RemoteServer

// RemoteServerBufferUploadsEnum is a string value for buffer_uploads.
// Enum lists the values documented by the API.
type RemoteServerBufferUploadsEnum string

// String returns the API parameter value.
func (u RemoteServerBufferUploadsEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u RemoteServerBufferUploadsEnum) Enum() map[string]RemoteServerBufferUploadsEnum {
	return map[string]RemoteServerBufferUploadsEnum{
		"auto":   RemoteServerBufferUploadsEnum("auto"),
		"always": RemoteServerBufferUploadsEnum("always"),
		"never":  RemoteServerBufferUploadsEnum("never"),
	}
}

// RemoteServerFilesAgentPermissionSetEnum is a string value for files_agent_permission_set.
// Enum lists the values documented by the API.
type RemoteServerFilesAgentPermissionSetEnum string

// String returns the API parameter value.
func (u RemoteServerFilesAgentPermissionSetEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u RemoteServerFilesAgentPermissionSetEnum) Enum() map[string]RemoteServerFilesAgentPermissionSetEnum {
	return map[string]RemoteServerFilesAgentPermissionSetEnum{
		"read_write": RemoteServerFilesAgentPermissionSetEnum("read_write"),
		"read_only":  RemoteServerFilesAgentPermissionSetEnum("read_only"),
		"write_only": RemoteServerFilesAgentPermissionSetEnum("write_only"),
	}
}

// RemoteServerGoogleCloudStorageAuthenticationMethodEnum is a string value for google_cloud_storage_authentication_method.
// Enum lists the values documented by the API.
type RemoteServerGoogleCloudStorageAuthenticationMethodEnum string

// String returns the API parameter value.
func (u RemoteServerGoogleCloudStorageAuthenticationMethodEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u RemoteServerGoogleCloudStorageAuthenticationMethodEnum) Enum() map[string]RemoteServerGoogleCloudStorageAuthenticationMethodEnum {
	return map[string]RemoteServerGoogleCloudStorageAuthenticationMethodEnum{
		"json":  RemoteServerGoogleCloudStorageAuthenticationMethodEnum("json"),
		"hmac":  RemoteServerGoogleCloudStorageAuthenticationMethodEnum("hmac"),
		"oauth": RemoteServerGoogleCloudStorageAuthenticationMethodEnum("oauth"),
	}
}

// RemoteServerOneDriveAccountTypeEnum is a string value for one_drive_account_type.
// Enum lists the values documented by the API.
type RemoteServerOneDriveAccountTypeEnum string

// String returns the API parameter value.
func (u RemoteServerOneDriveAccountTypeEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u RemoteServerOneDriveAccountTypeEnum) Enum() map[string]RemoteServerOneDriveAccountTypeEnum {
	return map[string]RemoteServerOneDriveAccountTypeEnum{
		"personal":       RemoteServerOneDriveAccountTypeEnum("personal"),
		"business_other": RemoteServerOneDriveAccountTypeEnum("business_other"),
	}
}

// RemoteServerServerCertificateEnum is a string value for server_certificate.
// Enum lists the values documented by the API.
type RemoteServerServerCertificateEnum string

// String returns the API parameter value.
func (u RemoteServerServerCertificateEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u RemoteServerServerCertificateEnum) Enum() map[string]RemoteServerServerCertificateEnum {
	return map[string]RemoteServerServerCertificateEnum{
		"require_match": RemoteServerServerCertificateEnum("require_match"),
		"allow_any":     RemoteServerServerCertificateEnum("allow_any"),
	}
}

// RemoteServerServerTypeEnum is a string value for server_type.
// Enum lists the values documented by the API.
type RemoteServerServerTypeEnum string

// String returns the API parameter value.
func (u RemoteServerServerTypeEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u RemoteServerServerTypeEnum) Enum() map[string]RemoteServerServerTypeEnum {
	return map[string]RemoteServerServerTypeEnum{
		"ftp":                  RemoteServerServerTypeEnum("ftp"),
		"sftp":                 RemoteServerServerTypeEnum("sftp"),
		"s3":                   RemoteServerServerTypeEnum("s3"),
		"google_cloud_storage": RemoteServerServerTypeEnum("google_cloud_storage"),
		"webdav":               RemoteServerServerTypeEnum("webdav"),
		"wasabi":               RemoteServerServerTypeEnum("wasabi"),
		"backblaze_b2":         RemoteServerServerTypeEnum("backblaze_b2"),
		"one_drive":            RemoteServerServerTypeEnum("one_drive"),
		"box":                  RemoteServerServerTypeEnum("box"),
		"dropbox":              RemoteServerServerTypeEnum("dropbox"),
		"google_drive":         RemoteServerServerTypeEnum("google_drive"),
		"azure":                RemoteServerServerTypeEnum("azure"),
		"sharepoint":           RemoteServerServerTypeEnum("sharepoint"),
		"s3_compatible":        RemoteServerServerTypeEnum("s3_compatible"),
		"azure_files":          RemoteServerServerTypeEnum("azure_files"),
		"files_agent":          RemoteServerServerTypeEnum("files_agent"),
		"filebase":             RemoteServerServerTypeEnum("filebase"),
		"cloudflare":           RemoteServerServerTypeEnum("cloudflare"),
		"linode":               RemoteServerServerTypeEnum("linode"),
		"files_com":            RemoteServerServerTypeEnum("files_com"),
	}
}

// RemoteServerSslEnum is a string value for ssl.
// Enum lists the values documented by the API.
type RemoteServerSslEnum string

// String returns the API parameter value.
func (u RemoteServerSslEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u RemoteServerSslEnum) Enum() map[string]RemoteServerSslEnum {
	return map[string]RemoteServerSslEnum{
		"if_available":     RemoteServerSslEnum("if_available"),
		"require":          RemoteServerSslEnum("require"),
		"require_implicit": RemoteServerSslEnum("require_implicit"),
		"never":            RemoteServerSslEnum("never"),
	}
}

// RemoteServerListParams contains the request parameters for GET /remote_servers.
type RemoteServerListParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `name`, `server_type`, `backblaze_b2_bucket`, `google_cloud_storage_bucket`, `wasabi_bucket`, `s3_bucket`, `azure_blob_storage_container`, `azure_files_storage_share_name`, `s3_compatible_bucket`, `filebase_bucket`, `cloudflare_bucket` or `linode_bucket`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `name`, `server_type`, `workspace_id`, `backblaze_b2_bucket`, `google_cloud_storage_bucket`, `wasabi_bucket`, `s3_bucket`, `azure_blob_storage_container`, `azure_files_storage_share_name`, `s3_compatible_bucket`, `filebase_bucket`, `cloudflare_bucket` or `linode_bucket`. Valid field combinations are `[ server_type, name ]`, `[ workspace_id, name ]`, `[ backblaze_b2_bucket, name ]`, `[ google_cloud_storage_bucket, name ]`, `[ wasabi_bucket, name ]`, `[ s3_bucket, name ]`, `[ azure_blob_storage_container, name ]`, `[ azure_files_storage_share_name, name ]`, `[ s3_compatible_bucket, name ]`, `[ filebase_bucket, name ]`, `[ cloudflare_bucket, name ]`, `[ linode_bucket, name ]`, `[ workspace_id, server_type ]` or `[ workspace_id, server_type, name ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is prefixed by the supplied value. Valid fields are `name`, `backblaze_b2_bucket`, `google_cloud_storage_bucket`, `wasabi_bucket`, `s3_bucket`, `azure_blob_storage_container`, `azure_files_storage_share_name`, `s3_compatible_bucket`, `filebase_bucket`, `cloudflare_bucket` or `linode_bucket`. Valid field combinations are `[ backblaze_b2_bucket, name ]`, `[ google_cloud_storage_bucket, name ]`, `[ wasabi_bucket, name ]`, `[ s3_bucket, name ]`, `[ azure_blob_storage_container, name ]`, `[ azure_files_storage_share_name, name ]`, `[ s3_compatible_bucket, name ]`, `[ filebase_bucket, name ]`, `[ cloudflare_bucket, name ]` or `[ linode_bucket, name ]`.
	FilterPrefix interface{} `url:"filter_prefix,omitempty" json:"filter_prefix,omitempty" path:"filter_prefix"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// RemoteServerFindParams contains the request parameters for GET /remote_servers/{id}.
type RemoteServerFindParams struct {
	// Remote Server ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// RemoteServerAgentNodesParams contains the request parameters for GET /remote_servers/{id}/agent_nodes.
//
// List Files.com Agent nodes
type RemoteServerAgentNodesParams struct {
	// Remote Server ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// RemoteServerFindConfigurationFileParams contains the request parameters for GET /remote_servers/{id}/configuration_file.
type RemoteServerFindConfigurationFileParams struct {
	// Remote Server ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// RemoteServerCreateParams contains the request parameters for POST /remote_servers.
type RemoteServerCreateParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// Password, if needed.
	Password string `url:"password,omitempty" json:"password,omitempty" path:"password"`
	// Private key, if needed.
	PrivateKey string `url:"private_key,omitempty" json:"private_key,omitempty" path:"private_key"`
	// Passphrase for private key if needed.
	PrivateKeyPassphrase string `url:"private_key_passphrase,omitempty" json:"private_key_passphrase,omitempty" path:"private_key_passphrase"`
	// Reset authenticated account?
	ResetAuthentication *bool `url:"reset_authentication,omitempty" json:"reset_authentication,omitempty" path:"reset_authentication"`
	// SharePoint: PEM-encoded certificate and unencrypted private key for app-only authentication.
	SharepointClientCertificate string `url:"sharepoint_client_certificate,omitempty" json:"sharepoint_client_certificate,omitempty" path:"sharepoint_client_certificate"`
	// SharePoint: Microsoft Entra application client secret for app-only authentication.
	SharepointClientSecret string `url:"sharepoint_client_secret,omitempty" json:"sharepoint_client_secret,omitempty" path:"sharepoint_client_secret"`
	// SSL client certificate.
	SslCertificate string `url:"ssl_certificate,omitempty" json:"ssl_certificate,omitempty" path:"ssl_certificate"`
	// AWS: secret key.
	AwsSecretKey string `url:"aws_secret_key,omitempty" json:"aws_secret_key,omitempty" path:"aws_secret_key"`
	// Azure Blob Storage: Access Key.
	AzureBlobStorageAccessKey string `url:"azure_blob_storage_access_key,omitempty" json:"azure_blob_storage_access_key,omitempty" path:"azure_blob_storage_access_key"`
	// Azure Blob Storage: Shared Access Signature (SAS) token.
	AzureBlobStorageSasToken string `url:"azure_blob_storage_sas_token,omitempty" json:"azure_blob_storage_sas_token,omitempty" path:"azure_blob_storage_sas_token"`
	// Azure File Storage: Access Key.
	AzureFilesStorageAccessKey string `url:"azure_files_storage_access_key,omitempty" json:"azure_files_storage_access_key,omitempty" path:"azure_files_storage_access_key"`
	// Azure File Storage: Shared Access Signature (SAS) token.
	AzureFilesStorageSasToken string `url:"azure_files_storage_sas_token,omitempty" json:"azure_files_storage_sas_token,omitempty" path:"azure_files_storage_sas_token"`
	// Backblaze B2 Cloud Storage: applicationKey.
	BackblazeB2ApplicationKey string `url:"backblaze_b2_application_key,omitempty" json:"backblaze_b2_application_key,omitempty" path:"backblaze_b2_application_key"`
	// Backblaze B2 Cloud Storage: keyID.
	BackblazeB2KeyId string `url:"backblaze_b2_key_id,omitempty" json:"backblaze_b2_key_id,omitempty" path:"backblaze_b2_key_id"`
	// Cloudflare: Secret Key.
	CloudflareSecretKey string `url:"cloudflare_secret_key,omitempty" json:"cloudflare_secret_key,omitempty" path:"cloudflare_secret_key"`
	// Filebase: Secret Key.
	FilebaseSecretKey string `url:"filebase_secret_key,omitempty" json:"filebase_secret_key,omitempty" path:"filebase_secret_key"`
	// Google Cloud Storage: JSON file that contains the private key. To generate see https://cloud.google.com/storage/docs/json_api/v1/how-tos/authorizing#APIKey.
	GoogleCloudStorageCredentialsJson string `url:"google_cloud_storage_credentials_json,omitempty" json:"google_cloud_storage_credentials_json,omitempty" path:"google_cloud_storage_credentials_json"`
	// Google Cloud Storage: S3-compatible secret key.
	GoogleCloudStorageS3CompatibleSecretKey string `url:"google_cloud_storage_s3_compatible_secret_key,omitempty" json:"google_cloud_storage_s3_compatible_secret_key,omitempty" path:"google_cloud_storage_s3_compatible_secret_key"`
	// Linode: Secret Key.
	LinodeSecretKey string `url:"linode_secret_key,omitempty" json:"linode_secret_key,omitempty" path:"linode_secret_key"`
	// S3-compatible: Secret Key.
	S3CompatibleSecretKey string `url:"s3_compatible_secret_key,omitempty" json:"s3_compatible_secret_key,omitempty" path:"s3_compatible_secret_key"`
	// Wasabi: Secret Key.
	WasabiSecretKey string `url:"wasabi_secret_key,omitempty" json:"wasabi_secret_key,omitempty" path:"wasabi_secret_key"`
	// Allow relative paths in SFTP. If true, paths will not be forced to be absolute, allowing operations relative to the user's home directory.
	AllowRelativePaths *bool `url:"allow_relative_paths,omitempty" json:"allow_relative_paths,omitempty" path:"allow_relative_paths"`
	// AWS Access Key.
	AwsAccessKey string `url:"aws_access_key,omitempty" json:"aws_access_key,omitempty" path:"aws_access_key"`
	// Azure Blob Storage: Account name.
	AzureBlobStorageAccount string `url:"azure_blob_storage_account,omitempty" json:"azure_blob_storage_account,omitempty" path:"azure_blob_storage_account"`
	// Azure Blob Storage: Container name.
	AzureBlobStorageContainer string `url:"azure_blob_storage_container,omitempty" json:"azure_blob_storage_container,omitempty" path:"azure_blob_storage_container"`
	// Azure Blob Storage: Custom DNS suffix.
	AzureBlobStorageDnsSuffix string `url:"azure_blob_storage_dns_suffix,omitempty" json:"azure_blob_storage_dns_suffix,omitempty" path:"azure_blob_storage_dns_suffix"`
	// Azure Blob Storage: Does the storage account has hierarchical namespace feature enabled?
	AzureBlobStorageHierarchicalNamespace *bool `url:"azure_blob_storage_hierarchical_namespace,omitempty" json:"azure_blob_storage_hierarchical_namespace,omitempty" path:"azure_blob_storage_hierarchical_namespace"`
	// Azure Files: Storage Account name.
	AzureFilesStorageAccount string `url:"azure_files_storage_account,omitempty" json:"azure_files_storage_account,omitempty" path:"azure_files_storage_account"`
	// Azure Files: Custom DNS suffix.
	AzureFilesStorageDnsSuffix string `url:"azure_files_storage_dns_suffix,omitempty" json:"azure_files_storage_dns_suffix,omitempty" path:"azure_files_storage_dns_suffix"`
	// Azure Files: Storage Share name.
	AzureFilesStorageShareName string `url:"azure_files_storage_share_name,omitempty" json:"azure_files_storage_share_name,omitempty" path:"azure_files_storage_share_name"`
	// Backblaze B2 Cloud Storage: Bucket name.
	BackblazeB2Bucket string `url:"backblaze_b2_bucket,omitempty" json:"backblaze_b2_bucket,omitempty" path:"backblaze_b2_bucket"`
	// Backblaze B2 Cloud Storage: S3 Endpoint.
	BackblazeB2S3Endpoint string `url:"backblaze_b2_s3_endpoint,omitempty" json:"backblaze_b2_s3_endpoint,omitempty" path:"backblaze_b2_s3_endpoint"`
	// If set to always, uploads to this server will be uploaded first to Files.com before being sent to the remote server. This can improve performance in certain access patterns, such as high-latency connections. It will cause data to be temporarily stored in Files.com. If set to auto, we will perform this optimization if we believe it to be a benefit in a given situation. Values: auto, always, never.
	BufferUploads RemoteServerBufferUploadsEnum `url:"buffer_uploads,omitempty" json:"buffer_uploads,omitempty" path:"buffer_uploads"`
	// Cloudflare: Access Key.
	CloudflareAccessKey string `url:"cloudflare_access_key,omitempty" json:"cloudflare_access_key,omitempty" path:"cloudflare_access_key"`
	// Cloudflare: Bucket name.
	CloudflareBucket string `url:"cloudflare_bucket,omitempty" json:"cloudflare_bucket,omitempty" path:"cloudflare_bucket"`
	// Cloudflare: endpoint.
	CloudflareEndpoint string `url:"cloudflare_endpoint,omitempty" json:"cloudflare_endpoint,omitempty" path:"cloudflare_endpoint"`
	// Internal description for your reference.
	Description string `url:"description,omitempty" json:"description,omitempty" path:"description"`
	// Dropbox: If true, list Team folders in root?
	DropboxTeams *bool `url:"dropbox_teams,omitempty" json:"dropbox_teams,omitempty" path:"dropbox_teams"`
	// `true` if remote server only accepts connections from dedicated IPs.
	EnableDedicatedIps *bool `url:"enable_dedicated_ips,omitempty" json:"enable_dedicated_ips,omitempty" path:"enable_dedicated_ips"`
	// Filebase: Access Key.
	FilebaseAccessKey string `url:"filebase_access_key,omitempty" json:"filebase_access_key,omitempty" path:"filebase_access_key"`
	// Filebase: Bucket name.
	FilebaseBucket string `url:"filebase_bucket,omitempty" json:"filebase_bucket,omitempty" path:"filebase_bucket"`
	// Files.com direct link: API key used once to pair the remote server.
	FilesApiKey string `url:"files_api_key,omitempty" json:"files_api_key,omitempty" path:"files_api_key"`
	// Agent file permissions: `read_only`, `write_only`, or `read_write`. Files.com writes this value as `permission_set` in the configuration file generated for this Remote Server. Setting it on Create or Update changes the generated file, not a running Agent. The Agent enforces the permissions in its local configuration file, controlled by the host's system administrator, and reports them at registration, replacing the stored value. Values: read_write, read_only, write_only.
	FilesAgentPermissionSet RemoteServerFilesAgentPermissionSetEnum `url:"files_agent_permission_set,omitempty" json:"files_agent_permission_set,omitempty" path:"files_agent_permission_set"`
	// Agent local root path. Files.com writes this value as `root` in the configuration file generated for this Remote Server. Setting it on Create or Update changes the generated file, not a running Agent. The Agent enforces the root in its local configuration file, controlled by the host's system administrator, and reports it at registration, replacing the stored value.
	FilesAgentRoot string `url:"files_agent_root,omitempty" json:"files_agent_root,omitempty" path:"files_agent_root"`
	// Files.com Agent version reported by the Agent at each registration, replacing the stored value. This field is not included in the generated Agent configuration file. Setting it on Create or Update does not update the running Agent.
	FilesAgentVersion string `url:"files_agent_version,omitempty" json:"files_agent_version,omitempty" path:"files_agent_version"`
	// Route traffic to outbound on a files-agent.
	OutboundAgentId int64 `url:"outbound_agent_id,omitempty" json:"outbound_agent_id,omitempty" path:"outbound_agent_id"`
	// Custom Domain ID whose dedicated IP addresses are selected when this Remote Server uses dedicated IPs. Must be available to this Remote Server's workspace. Requires enable_dedicated_ips and cannot be combined with an outbound Agent. Set to null to use the site's default dedicated IPs.
	CustomDomainId int64 `url:"custom_domain_id,omitempty" json:"custom_domain_id,omitempty" path:"custom_domain_id"`
	// Google Cloud Storage: Authentication method. Can be json, hmac, or oauth. Values: json, hmac, oauth.
	GoogleCloudStorageAuthenticationMethod RemoteServerGoogleCloudStorageAuthenticationMethodEnum `url:"google_cloud_storage_authentication_method,omitempty" json:"google_cloud_storage_authentication_method,omitempty" path:"google_cloud_storage_authentication_method"`
	// Google Cloud Storage: Bucket Name.
	GoogleCloudStorageBucket string `url:"google_cloud_storage_bucket,omitempty" json:"google_cloud_storage_bucket,omitempty" path:"google_cloud_storage_bucket"`
	// Google Cloud Storage: OAuth scope. Can be https://www.googleapis.com/auth/devstorage.read_only or https://www.googleapis.com/auth/devstorage.read_write.
	GoogleCloudStorageOauthScope string `url:"google_cloud_storage_oauth_scope,omitempty" json:"google_cloud_storage_oauth_scope,omitempty" path:"google_cloud_storage_oauth_scope"`
	// Google Cloud Storage: Project ID.
	GoogleCloudStorageProjectId string `url:"google_cloud_storage_project_id,omitempty" json:"google_cloud_storage_project_id,omitempty" path:"google_cloud_storage_project_id"`
	// Google Cloud Storage: S3-compatible Access Key.
	GoogleCloudStorageS3CompatibleAccessKey string `url:"google_cloud_storage_s3_compatible_access_key,omitempty" json:"google_cloud_storage_s3_compatible_access_key,omitempty" path:"google_cloud_storage_s3_compatible_access_key"`
	// Hostname or IP address. For Agent Remote Servers, the Agent reports the hostname it connects from at each registration, replacing the stored value. This field is not included in the generated Agent configuration file. Setting it on Create or Update does not change the running Agent.
	Hostname string `url:"hostname,omitempty" json:"hostname,omitempty" path:"hostname"`
	// Linode: Access Key.
	LinodeAccessKey string `url:"linode_access_key,omitempty" json:"linode_access_key,omitempty" path:"linode_access_key"`
	// Linode: Bucket name.
	LinodeBucket string `url:"linode_bucket,omitempty" json:"linode_bucket,omitempty" path:"linode_bucket"`
	// Linode: region.
	LinodeRegion string `url:"linode_region,omitempty" json:"linode_region,omitempty" path:"linode_region"`
	// Max number of parallel connections. Ignored for S3 connections (we will parallelize these as much as possible).
	MaxConnections int64 `url:"max_connections,omitempty" json:"max_connections,omitempty" path:"max_connections"`
	// Internal name for your reference.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// OneDrive: Either personal or business_other account types. Values: personal, business_other.
	OneDriveAccountType RemoteServerOneDriveAccountTypeEnum `url:"one_drive_account_type,omitempty" json:"one_drive_account_type,omitempty" path:"one_drive_account_type"`
	// If true, we will ensure that all communications with this remote server are made through the primary region of the site. This setting can also be overridden by a site-wide setting which will force it to true.
	PinToSiteRegion *bool `url:"pin_to_site_region,omitempty" json:"pin_to_site_region,omitempty" path:"pin_to_site_region"`
	// Port for remote server. For Agent Remote Servers, the Agent reports its downstream proxy port at each registration, replacing the stored value. This field is not included in the generated Agent configuration file. Setting it on Create or Update does not change the running Agent.
	Port int64 `url:"port,omitempty" json:"port,omitempty" path:"port"`
	// Upload staging path. Applies to SFTP only. If a path is provided here, files will first be uploaded to this path on the remote folder and the moved into the final correct path via an SFTP move command. This is required by some remote MFT systems to emulate atomic uploads, which are otherwise not supoprted by SFTP.
	UploadStagingPath string `url:"upload_staging_path,omitempty" json:"upload_staging_path,omitempty" path:"upload_staging_path"`
	// ID of Remote Server Credential, if applicable.
	RemoteServerCredentialId int64 `url:"remote_server_credential_id,omitempty" json:"remote_server_credential_id,omitempty" path:"remote_server_credential_id"`
	// AWS IAM Role ARN for AssumeRole authentication.
	S3AssumeRoleArn string `url:"s3_assume_role_arn,omitempty" json:"s3_assume_role_arn,omitempty" path:"s3_assume_role_arn"`
	// Session duration in seconds for AssumeRole authentication (900-43200).
	S3AssumeRoleDurationSeconds int64 `url:"s3_assume_role_duration_seconds,omitempty" json:"s3_assume_role_duration_seconds,omitempty" path:"s3_assume_role_duration_seconds"`
	// S3 bucket name.
	S3Bucket string `url:"s3_bucket,omitempty" json:"s3_bucket,omitempty" path:"s3_bucket"`
	// S3-compatible: Access Key.
	S3CompatibleAccessKey string `url:"s3_compatible_access_key,omitempty" json:"s3_compatible_access_key,omitempty" path:"s3_compatible_access_key"`
	// S3-compatible: Bucket name.
	S3CompatibleBucket string `url:"s3_compatible_bucket,omitempty" json:"s3_compatible_bucket,omitempty" path:"s3_compatible_bucket"`
	// S3-compatible: endpoint.
	S3CompatibleEndpoint string `url:"s3_compatible_endpoint,omitempty" json:"s3_compatible_endpoint,omitempty" path:"s3_compatible_endpoint"`
	// S3-compatible: region.
	S3CompatibleRegion string `url:"s3_compatible_region,omitempty" json:"s3_compatible_region,omitempty" path:"s3_compatible_region"`
	// S3-compatible: If true, use virtual-hosted-style URLs instead of path-style URLs.
	S3CompatibleVirtualHostedStyle *bool `url:"s3_compatible_virtual_hosted_style,omitempty" json:"s3_compatible_virtual_hosted_style,omitempty" path:"s3_compatible_virtual_hosted_style"`
	// ARN of the AWS KMS key that encrypts files written to the bucket (SSE-KMS). Leave blank unless your bucket requires a specific key.
	S3KmsKeyId string `url:"s3_kms_key_id,omitempty" json:"s3_kms_key_id,omitempty" path:"s3_kms_key_id"`
	// S3 region.
	S3Region string `url:"s3_region,omitempty" json:"s3_region,omitempty" path:"s3_region"`
	// Remote server certificate. Values: require_match, allow_any.
	ServerCertificate RemoteServerServerCertificateEnum `url:"server_certificate,omitempty" json:"server_certificate,omitempty" path:"server_certificate"`
	// Pinned SSH host key or OpenSSH host certificate for SFTP. If omitted, Files.com detects and stores a host key, preferring plain keys over certificates. With `server_certificate=require_match` (the default), the server must present the exact pinned key or certificate and prove it holds the matching private key. A pinned certificate is compared in full, so renewal can require updating `server_host_key` even when its underlying key is unchanged. Files.com does not check certificate CA signatures, principals, or validity periods. Certificate expiration alone does not end the pin. Update `server_host_key` to replace the pin.
	ServerHostKey string `url:"server_host_key,omitempty" json:"server_host_key,omitempty" path:"server_host_key"`
	// Remote server type. Values: ftp, sftp, s3, google_cloud_storage, webdav, wasabi, backblaze_b2, one_drive, box, dropbox, google_drive, azure, sharepoint, s3_compatible, azure_files, files_agent, filebase, cloudflare, linode, files_com.
	ServerType RemoteServerServerTypeEnum `url:"server_type,omitempty" json:"server_type,omitempty" path:"server_type"`
	// SharePoint: Microsoft Entra application client ID for app-only authentication.
	SharepointClientId string `url:"sharepoint_client_id,omitempty" json:"sharepoint_client_id,omitempty" path:"sharepoint_client_id"`
	// SharePoint: Site URL to scope app-only authentication to a single site. Leave blank to browse all sites.
	SharepointSiteUrl string `url:"sharepoint_site_url,omitempty" json:"sharepoint_site_url,omitempty" path:"sharepoint_site_url"`
	// SharePoint: Microsoft Entra tenant ID for app-only authentication.
	SharepointTenantId string `url:"sharepoint_tenant_id,omitempty" json:"sharepoint_tenant_id,omitempty" path:"sharepoint_tenant_id"`
	// Should we require SSL? Values: if_available, require, require_implicit, never.
	Ssl RemoteServerSslEnum `url:"ssl,omitempty" json:"ssl,omitempty" path:"ssl"`
	// Remote server username.
	Username string `url:"username,omitempty" json:"username,omitempty" path:"username"`
	// Wasabi: Access Key.
	WasabiAccessKey string `url:"wasabi_access_key,omitempty" json:"wasabi_access_key,omitempty" path:"wasabi_access_key"`
	// Wasabi: Bucket name.
	WasabiBucket string `url:"wasabi_bucket,omitempty" json:"wasabi_bucket,omitempty" path:"wasabi_bucket"`
	// Wasabi: Region.
	WasabiRegion string `url:"wasabi_region,omitempty" json:"wasabi_region,omitempty" path:"wasabi_region"`
	// Workspace ID (0 for default workspace).
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// RemoteServerAgentPushUpdateParams contains the request parameters for POST /remote_servers/{id}/agent_push_update.
//
// Push update to Files Agent
type RemoteServerAgentPushUpdateParams struct {
	// Remote Server ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// RemoteServerUpdateParams contains the request parameters for PATCH /remote_servers/{id}.
type RemoteServerUpdateParams struct {
	// Remote Server ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Password, if needed.
	Password string `url:"password,omitempty" json:"password,omitempty" path:"password"`
	// Private key, if needed.
	PrivateKey string `url:"private_key,omitempty" json:"private_key,omitempty" path:"private_key"`
	// Passphrase for private key if needed.
	PrivateKeyPassphrase string `url:"private_key_passphrase,omitempty" json:"private_key_passphrase,omitempty" path:"private_key_passphrase"`
	// Reset authenticated account?
	ResetAuthentication *bool `url:"reset_authentication,omitempty" json:"reset_authentication,omitempty" path:"reset_authentication"`
	// SharePoint: PEM-encoded certificate and unencrypted private key for app-only authentication.
	SharepointClientCertificate string `url:"sharepoint_client_certificate,omitempty" json:"sharepoint_client_certificate,omitempty" path:"sharepoint_client_certificate"`
	// SharePoint: Microsoft Entra application client secret for app-only authentication.
	SharepointClientSecret string `url:"sharepoint_client_secret,omitempty" json:"sharepoint_client_secret,omitempty" path:"sharepoint_client_secret"`
	// SSL client certificate.
	SslCertificate string `url:"ssl_certificate,omitempty" json:"ssl_certificate,omitempty" path:"ssl_certificate"`
	// AWS: secret key.
	AwsSecretKey string `url:"aws_secret_key,omitempty" json:"aws_secret_key,omitempty" path:"aws_secret_key"`
	// Azure Blob Storage: Access Key.
	AzureBlobStorageAccessKey string `url:"azure_blob_storage_access_key,omitempty" json:"azure_blob_storage_access_key,omitempty" path:"azure_blob_storage_access_key"`
	// Azure Blob Storage: Shared Access Signature (SAS) token.
	AzureBlobStorageSasToken string `url:"azure_blob_storage_sas_token,omitempty" json:"azure_blob_storage_sas_token,omitempty" path:"azure_blob_storage_sas_token"`
	// Azure File Storage: Access Key.
	AzureFilesStorageAccessKey string `url:"azure_files_storage_access_key,omitempty" json:"azure_files_storage_access_key,omitempty" path:"azure_files_storage_access_key"`
	// Azure File Storage: Shared Access Signature (SAS) token.
	AzureFilesStorageSasToken string `url:"azure_files_storage_sas_token,omitempty" json:"azure_files_storage_sas_token,omitempty" path:"azure_files_storage_sas_token"`
	// Backblaze B2 Cloud Storage: applicationKey.
	BackblazeB2ApplicationKey string `url:"backblaze_b2_application_key,omitempty" json:"backblaze_b2_application_key,omitempty" path:"backblaze_b2_application_key"`
	// Backblaze B2 Cloud Storage: keyID.
	BackblazeB2KeyId string `url:"backblaze_b2_key_id,omitempty" json:"backblaze_b2_key_id,omitempty" path:"backblaze_b2_key_id"`
	// Cloudflare: Secret Key.
	CloudflareSecretKey string `url:"cloudflare_secret_key,omitempty" json:"cloudflare_secret_key,omitempty" path:"cloudflare_secret_key"`
	// Filebase: Secret Key.
	FilebaseSecretKey string `url:"filebase_secret_key,omitempty" json:"filebase_secret_key,omitempty" path:"filebase_secret_key"`
	// Google Cloud Storage: JSON file that contains the private key. To generate see https://cloud.google.com/storage/docs/json_api/v1/how-tos/authorizing#APIKey.
	GoogleCloudStorageCredentialsJson string `url:"google_cloud_storage_credentials_json,omitempty" json:"google_cloud_storage_credentials_json,omitempty" path:"google_cloud_storage_credentials_json"`
	// Google Cloud Storage: S3-compatible secret key.
	GoogleCloudStorageS3CompatibleSecretKey string `url:"google_cloud_storage_s3_compatible_secret_key,omitempty" json:"google_cloud_storage_s3_compatible_secret_key,omitempty" path:"google_cloud_storage_s3_compatible_secret_key"`
	// Linode: Secret Key.
	LinodeSecretKey string `url:"linode_secret_key,omitempty" json:"linode_secret_key,omitempty" path:"linode_secret_key"`
	// S3-compatible: Secret Key.
	S3CompatibleSecretKey string `url:"s3_compatible_secret_key,omitempty" json:"s3_compatible_secret_key,omitempty" path:"s3_compatible_secret_key"`
	// Wasabi: Secret Key.
	WasabiSecretKey string `url:"wasabi_secret_key,omitempty" json:"wasabi_secret_key,omitempty" path:"wasabi_secret_key"`
	// Allow relative paths in SFTP. If true, paths will not be forced to be absolute, allowing operations relative to the user's home directory.
	AllowRelativePaths *bool `url:"allow_relative_paths,omitempty" json:"allow_relative_paths,omitempty" path:"allow_relative_paths"`
	// AWS Access Key.
	AwsAccessKey string `url:"aws_access_key,omitempty" json:"aws_access_key,omitempty" path:"aws_access_key"`
	// Azure Blob Storage: Account name.
	AzureBlobStorageAccount string `url:"azure_blob_storage_account,omitempty" json:"azure_blob_storage_account,omitempty" path:"azure_blob_storage_account"`
	// Azure Blob Storage: Container name.
	AzureBlobStorageContainer string `url:"azure_blob_storage_container,omitempty" json:"azure_blob_storage_container,omitempty" path:"azure_blob_storage_container"`
	// Azure Blob Storage: Custom DNS suffix.
	AzureBlobStorageDnsSuffix string `url:"azure_blob_storage_dns_suffix,omitempty" json:"azure_blob_storage_dns_suffix,omitempty" path:"azure_blob_storage_dns_suffix"`
	// Azure Blob Storage: Does the storage account has hierarchical namespace feature enabled?
	AzureBlobStorageHierarchicalNamespace *bool `url:"azure_blob_storage_hierarchical_namespace,omitempty" json:"azure_blob_storage_hierarchical_namespace,omitempty" path:"azure_blob_storage_hierarchical_namespace"`
	// Azure Files: Storage Account name.
	AzureFilesStorageAccount string `url:"azure_files_storage_account,omitempty" json:"azure_files_storage_account,omitempty" path:"azure_files_storage_account"`
	// Azure Files: Custom DNS suffix.
	AzureFilesStorageDnsSuffix string `url:"azure_files_storage_dns_suffix,omitempty" json:"azure_files_storage_dns_suffix,omitempty" path:"azure_files_storage_dns_suffix"`
	// Azure Files: Storage Share name.
	AzureFilesStorageShareName string `url:"azure_files_storage_share_name,omitempty" json:"azure_files_storage_share_name,omitempty" path:"azure_files_storage_share_name"`
	// Backblaze B2 Cloud Storage: Bucket name.
	BackblazeB2Bucket string `url:"backblaze_b2_bucket,omitempty" json:"backblaze_b2_bucket,omitempty" path:"backblaze_b2_bucket"`
	// Backblaze B2 Cloud Storage: S3 Endpoint.
	BackblazeB2S3Endpoint string `url:"backblaze_b2_s3_endpoint,omitempty" json:"backblaze_b2_s3_endpoint,omitempty" path:"backblaze_b2_s3_endpoint"`
	// If set to always, uploads to this server will be uploaded first to Files.com before being sent to the remote server. This can improve performance in certain access patterns, such as high-latency connections. It will cause data to be temporarily stored in Files.com. If set to auto, we will perform this optimization if we believe it to be a benefit in a given situation. Values: auto, always, never.
	BufferUploads RemoteServerBufferUploadsEnum `url:"buffer_uploads,omitempty" json:"buffer_uploads,omitempty" path:"buffer_uploads"`
	// Cloudflare: Access Key.
	CloudflareAccessKey string `url:"cloudflare_access_key,omitempty" json:"cloudflare_access_key,omitempty" path:"cloudflare_access_key"`
	// Cloudflare: Bucket name.
	CloudflareBucket string `url:"cloudflare_bucket,omitempty" json:"cloudflare_bucket,omitempty" path:"cloudflare_bucket"`
	// Cloudflare: endpoint.
	CloudflareEndpoint string `url:"cloudflare_endpoint,omitempty" json:"cloudflare_endpoint,omitempty" path:"cloudflare_endpoint"`
	// Internal description for your reference.
	Description string `url:"description,omitempty" json:"description,omitempty" path:"description"`
	// Dropbox: If true, list Team folders in root?
	DropboxTeams *bool `url:"dropbox_teams,omitempty" json:"dropbox_teams,omitempty" path:"dropbox_teams"`
	// `true` if remote server only accepts connections from dedicated IPs.
	EnableDedicatedIps *bool `url:"enable_dedicated_ips,omitempty" json:"enable_dedicated_ips,omitempty" path:"enable_dedicated_ips"`
	// Filebase: Access Key.
	FilebaseAccessKey string `url:"filebase_access_key,omitempty" json:"filebase_access_key,omitempty" path:"filebase_access_key"`
	// Filebase: Bucket name.
	FilebaseBucket string `url:"filebase_bucket,omitempty" json:"filebase_bucket,omitempty" path:"filebase_bucket"`
	// Files.com direct link: API key used once to pair the remote server.
	FilesApiKey string `url:"files_api_key,omitempty" json:"files_api_key,omitempty" path:"files_api_key"`
	// Agent file permissions: `read_only`, `write_only`, or `read_write`. Files.com writes this value as `permission_set` in the configuration file generated for this Remote Server. Setting it on Create or Update changes the generated file, not a running Agent. The Agent enforces the permissions in its local configuration file, controlled by the host's system administrator, and reports them at registration, replacing the stored value. Values: read_write, read_only, write_only.
	FilesAgentPermissionSet RemoteServerFilesAgentPermissionSetEnum `url:"files_agent_permission_set,omitempty" json:"files_agent_permission_set,omitempty" path:"files_agent_permission_set"`
	// Agent local root path. Files.com writes this value as `root` in the configuration file generated for this Remote Server. Setting it on Create or Update changes the generated file, not a running Agent. The Agent enforces the root in its local configuration file, controlled by the host's system administrator, and reports it at registration, replacing the stored value.
	FilesAgentRoot string `url:"files_agent_root,omitempty" json:"files_agent_root,omitempty" path:"files_agent_root"`
	// Files.com Agent version reported by the Agent at each registration, replacing the stored value. This field is not included in the generated Agent configuration file. Setting it on Create or Update does not update the running Agent.
	FilesAgentVersion string `url:"files_agent_version,omitempty" json:"files_agent_version,omitempty" path:"files_agent_version"`
	// Route traffic to outbound on a files-agent.
	OutboundAgentId int64 `url:"outbound_agent_id,omitempty" json:"outbound_agent_id,omitempty" path:"outbound_agent_id"`
	// Custom Domain ID whose dedicated IP addresses are selected when this Remote Server uses dedicated IPs. Must be available to this Remote Server's workspace. Requires enable_dedicated_ips and cannot be combined with an outbound Agent. Set to null to use the site's default dedicated IPs.
	CustomDomainId int64 `url:"custom_domain_id,omitempty" json:"custom_domain_id,omitempty" path:"custom_domain_id"`
	// Google Cloud Storage: Authentication method. Can be json, hmac, or oauth. Values: json, hmac, oauth.
	GoogleCloudStorageAuthenticationMethod RemoteServerGoogleCloudStorageAuthenticationMethodEnum `url:"google_cloud_storage_authentication_method,omitempty" json:"google_cloud_storage_authentication_method,omitempty" path:"google_cloud_storage_authentication_method"`
	// Google Cloud Storage: Bucket Name.
	GoogleCloudStorageBucket string `url:"google_cloud_storage_bucket,omitempty" json:"google_cloud_storage_bucket,omitempty" path:"google_cloud_storage_bucket"`
	// Google Cloud Storage: OAuth scope. Can be https://www.googleapis.com/auth/devstorage.read_only or https://www.googleapis.com/auth/devstorage.read_write.
	GoogleCloudStorageOauthScope string `url:"google_cloud_storage_oauth_scope,omitempty" json:"google_cloud_storage_oauth_scope,omitempty" path:"google_cloud_storage_oauth_scope"`
	// Google Cloud Storage: Project ID.
	GoogleCloudStorageProjectId string `url:"google_cloud_storage_project_id,omitempty" json:"google_cloud_storage_project_id,omitempty" path:"google_cloud_storage_project_id"`
	// Google Cloud Storage: S3-compatible Access Key.
	GoogleCloudStorageS3CompatibleAccessKey string `url:"google_cloud_storage_s3_compatible_access_key,omitempty" json:"google_cloud_storage_s3_compatible_access_key,omitempty" path:"google_cloud_storage_s3_compatible_access_key"`
	// Hostname or IP address. For Agent Remote Servers, the Agent reports the hostname it connects from at each registration, replacing the stored value. This field is not included in the generated Agent configuration file. Setting it on Create or Update does not change the running Agent.
	Hostname string `url:"hostname,omitempty" json:"hostname,omitempty" path:"hostname"`
	// Linode: Access Key.
	LinodeAccessKey string `url:"linode_access_key,omitempty" json:"linode_access_key,omitempty" path:"linode_access_key"`
	// Linode: Bucket name.
	LinodeBucket string `url:"linode_bucket,omitempty" json:"linode_bucket,omitempty" path:"linode_bucket"`
	// Linode: region.
	LinodeRegion string `url:"linode_region,omitempty" json:"linode_region,omitempty" path:"linode_region"`
	// Max number of parallel connections. Ignored for S3 connections (we will parallelize these as much as possible).
	MaxConnections int64 `url:"max_connections,omitempty" json:"max_connections,omitempty" path:"max_connections"`
	// Internal name for your reference.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// OneDrive: Either personal or business_other account types. Values: personal, business_other.
	OneDriveAccountType RemoteServerOneDriveAccountTypeEnum `url:"one_drive_account_type,omitempty" json:"one_drive_account_type,omitempty" path:"one_drive_account_type"`
	// If true, we will ensure that all communications with this remote server are made through the primary region of the site. This setting can also be overridden by a site-wide setting which will force it to true.
	PinToSiteRegion *bool `url:"pin_to_site_region,omitempty" json:"pin_to_site_region,omitempty" path:"pin_to_site_region"`
	// Port for remote server. For Agent Remote Servers, the Agent reports its downstream proxy port at each registration, replacing the stored value. This field is not included in the generated Agent configuration file. Setting it on Create or Update does not change the running Agent.
	Port int64 `url:"port,omitempty" json:"port,omitempty" path:"port"`
	// Upload staging path. Applies to SFTP only. If a path is provided here, files will first be uploaded to this path on the remote folder and the moved into the final correct path via an SFTP move command. This is required by some remote MFT systems to emulate atomic uploads, which are otherwise not supoprted by SFTP.
	UploadStagingPath string `url:"upload_staging_path,omitempty" json:"upload_staging_path,omitempty" path:"upload_staging_path"`
	// ID of Remote Server Credential, if applicable.
	RemoteServerCredentialId int64 `url:"remote_server_credential_id,omitempty" json:"remote_server_credential_id,omitempty" path:"remote_server_credential_id"`
	// AWS IAM Role ARN for AssumeRole authentication.
	S3AssumeRoleArn string `url:"s3_assume_role_arn,omitempty" json:"s3_assume_role_arn,omitempty" path:"s3_assume_role_arn"`
	// Session duration in seconds for AssumeRole authentication (900-43200).
	S3AssumeRoleDurationSeconds int64 `url:"s3_assume_role_duration_seconds,omitempty" json:"s3_assume_role_duration_seconds,omitempty" path:"s3_assume_role_duration_seconds"`
	// S3 bucket name.
	S3Bucket string `url:"s3_bucket,omitempty" json:"s3_bucket,omitempty" path:"s3_bucket"`
	// S3-compatible: Access Key.
	S3CompatibleAccessKey string `url:"s3_compatible_access_key,omitempty" json:"s3_compatible_access_key,omitempty" path:"s3_compatible_access_key"`
	// S3-compatible: Bucket name.
	S3CompatibleBucket string `url:"s3_compatible_bucket,omitempty" json:"s3_compatible_bucket,omitempty" path:"s3_compatible_bucket"`
	// S3-compatible: endpoint.
	S3CompatibleEndpoint string `url:"s3_compatible_endpoint,omitempty" json:"s3_compatible_endpoint,omitempty" path:"s3_compatible_endpoint"`
	// S3-compatible: region.
	S3CompatibleRegion string `url:"s3_compatible_region,omitempty" json:"s3_compatible_region,omitempty" path:"s3_compatible_region"`
	// S3-compatible: If true, use virtual-hosted-style URLs instead of path-style URLs.
	S3CompatibleVirtualHostedStyle *bool `url:"s3_compatible_virtual_hosted_style,omitempty" json:"s3_compatible_virtual_hosted_style,omitempty" path:"s3_compatible_virtual_hosted_style"`
	// ARN of the AWS KMS key that encrypts files written to the bucket (SSE-KMS). Leave blank unless your bucket requires a specific key.
	S3KmsKeyId string `url:"s3_kms_key_id,omitempty" json:"s3_kms_key_id,omitempty" path:"s3_kms_key_id"`
	// S3 region.
	S3Region string `url:"s3_region,omitempty" json:"s3_region,omitempty" path:"s3_region"`
	// Remote server certificate. Values: require_match, allow_any.
	ServerCertificate RemoteServerServerCertificateEnum `url:"server_certificate,omitempty" json:"server_certificate,omitempty" path:"server_certificate"`
	// Pinned SSH host key or OpenSSH host certificate for SFTP. If omitted, Files.com detects and stores a host key, preferring plain keys over certificates. With `server_certificate=require_match` (the default), the server must present the exact pinned key or certificate and prove it holds the matching private key. A pinned certificate is compared in full, so renewal can require updating `server_host_key` even when its underlying key is unchanged. Files.com does not check certificate CA signatures, principals, or validity periods. Certificate expiration alone does not end the pin. Update `server_host_key` to replace the pin.
	ServerHostKey string `url:"server_host_key,omitempty" json:"server_host_key,omitempty" path:"server_host_key"`
	// Remote server type. Values: ftp, sftp, s3, google_cloud_storage, webdav, wasabi, backblaze_b2, one_drive, box, dropbox, google_drive, azure, sharepoint, s3_compatible, azure_files, files_agent, filebase, cloudflare, linode, files_com.
	ServerType RemoteServerServerTypeEnum `url:"server_type,omitempty" json:"server_type,omitempty" path:"server_type"`
	// SharePoint: Microsoft Entra application client ID for app-only authentication.
	SharepointClientId string `url:"sharepoint_client_id,omitempty" json:"sharepoint_client_id,omitempty" path:"sharepoint_client_id"`
	// SharePoint: Site URL to scope app-only authentication to a single site. Leave blank to browse all sites.
	SharepointSiteUrl string `url:"sharepoint_site_url,omitempty" json:"sharepoint_site_url,omitempty" path:"sharepoint_site_url"`
	// SharePoint: Microsoft Entra tenant ID for app-only authentication.
	SharepointTenantId string `url:"sharepoint_tenant_id,omitempty" json:"sharepoint_tenant_id,omitempty" path:"sharepoint_tenant_id"`
	// Should we require SSL? Values: if_available, require, require_implicit, never.
	Ssl RemoteServerSslEnum `url:"ssl,omitempty" json:"ssl,omitempty" path:"ssl"`
	// Remote server username.
	Username string `url:"username,omitempty" json:"username,omitempty" path:"username"`
	// Wasabi: Access Key.
	WasabiAccessKey string `url:"wasabi_access_key,omitempty" json:"wasabi_access_key,omitempty" path:"wasabi_access_key"`
	// Wasabi: Bucket name.
	WasabiBucket string `url:"wasabi_bucket,omitempty" json:"wasabi_bucket,omitempty" path:"wasabi_bucket"`
	// Wasabi: Region.
	WasabiRegion string `url:"wasabi_region,omitempty" json:"wasabi_region,omitempty" path:"wasabi_region"`
}

// RemoteServerDeleteParams contains the request parameters for DELETE /remote_servers/{id}.
type RemoteServerDeleteParams struct {
	// Remote Server ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (r *RemoteServer) UnmarshalJSON(data []byte) error {
	type remoteServer RemoteServer
	var v remoteServer
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*r = RemoteServer(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (r *RemoteServerCollection) UnmarshalJSON(data []byte) error {
	type remoteServers RemoteServerCollection
	var v remoteServers
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*r = RemoteServerCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (r *RemoteServerCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*r))
	for i, v := range *r {
		ret[i] = v
	}

	return &ret
}
