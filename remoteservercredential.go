package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// RemoteServerCredential is a Files.com API resource.
type RemoteServerCredential struct {
	// Remote Server Credential ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Workspace ID (0 for default workspace).
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Internal name for your reference.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Internal description for your reference.
	Description string `json:"description,omitempty" path:"description,omitempty" url:"description,omitempty"`
	// Remote server type. Remote Server Credentials are only valid for a single type of Remote Server.
	ServerType string `json:"server_type,omitempty" path:"server_type,omitempty" url:"server_type,omitempty"`
	// AWS Access Key.
	AwsAccessKey string `json:"aws_access_key,omitempty" path:"aws_access_key,omitempty" url:"aws_access_key,omitempty"`
	// AWS IAM Role ARN for AssumeRole authentication.
	S3AssumeRoleArn string `json:"s3_assume_role_arn,omitempty" path:"s3_assume_role_arn,omitempty" url:"s3_assume_role_arn,omitempty"`
	// Session duration in seconds for AssumeRole authentication (900-43200).
	S3AssumeRoleDurationSeconds int64 `json:"s3_assume_role_duration_seconds,omitempty" path:"s3_assume_role_duration_seconds,omitempty" url:"s3_assume_role_duration_seconds,omitempty"`
	// External ID for AssumeRole authentication.
	S3AssumeRoleExternalId string `json:"s3_assume_role_external_id,omitempty" path:"s3_assume_role_external_id,omitempty" url:"s3_assume_role_external_id,omitempty"`
	// Google Cloud Storage: S3-compatible Access Key.
	GoogleCloudStorageS3CompatibleAccessKey string `json:"google_cloud_storage_s3_compatible_access_key,omitempty" path:"google_cloud_storage_s3_compatible_access_key,omitempty" url:"google_cloud_storage_s3_compatible_access_key,omitempty"`
	// Wasabi: Access Key.
	WasabiAccessKey string `json:"wasabi_access_key,omitempty" path:"wasabi_access_key,omitempty" url:"wasabi_access_key,omitempty"`
	// S3-compatible: Access Key.
	S3CompatibleAccessKey string `json:"s3_compatible_access_key,omitempty" path:"s3_compatible_access_key,omitempty" url:"s3_compatible_access_key,omitempty"`
	// Filebase: Access Key.
	FilebaseAccessKey string `json:"filebase_access_key,omitempty" path:"filebase_access_key,omitempty" url:"filebase_access_key,omitempty"`
	// Cloudflare: Access Key.
	CloudflareAccessKey string `json:"cloudflare_access_key,omitempty" path:"cloudflare_access_key,omitempty" url:"cloudflare_access_key,omitempty"`
	// Linode: Access Key.
	LinodeAccessKey string `json:"linode_access_key,omitempty" path:"linode_access_key,omitempty" url:"linode_access_key,omitempty"`
	// SharePoint: Microsoft Entra tenant ID for app-only authentication.
	SharepointTenantId string `json:"sharepoint_tenant_id,omitempty" path:"sharepoint_tenant_id,omitempty" url:"sharepoint_tenant_id,omitempty"`
	// SharePoint: Microsoft Entra application client ID for app-only authentication.
	SharepointClientId string `json:"sharepoint_client_id,omitempty" path:"sharepoint_client_id,omitempty" url:"sharepoint_client_id,omitempty"`
	// SharePoint: App-only credential type. Either secret or certificate.
	SharepointAppCredentialType string `json:"sharepoint_app_credential_type,omitempty" path:"sharepoint_app_credential_type,omitempty" url:"sharepoint_app_credential_type,omitempty"`
	// Remote server username.
	Username string `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	// Password, if needed.
	Password string `json:"password,omitempty" path:"password,omitempty" url:"password,omitempty"`
	// Private key, if needed.
	PrivateKey string `json:"private_key,omitempty" path:"private_key,omitempty" url:"private_key,omitempty"`
	// Passphrase for private key if needed.
	PrivateKeyPassphrase string `json:"private_key_passphrase,omitempty" path:"private_key_passphrase,omitempty" url:"private_key_passphrase,omitempty"`
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
	// SharePoint: PEM-encoded certificate and unencrypted private key for app-only authentication.
	SharepointClientCertificate string `json:"sharepoint_client_certificate,omitempty" path:"sharepoint_client_certificate,omitempty" url:"sharepoint_client_certificate,omitempty"`
	// SharePoint: Microsoft Entra application client secret for app-only authentication.
	SharepointClientSecret string `json:"sharepoint_client_secret,omitempty" path:"sharepoint_client_secret,omitempty" url:"sharepoint_client_secret,omitempty"`
	// Wasabi: Secret Key.
	WasabiSecretKey string `json:"wasabi_secret_key,omitempty" path:"wasabi_secret_key,omitempty" url:"wasabi_secret_key,omitempty"`
	// ID of Remote Server Credential to copy omitted values from.
	CopyValuesFromCredentialId int64 `json:"copy_values_from_credential_id,omitempty" path:"copy_values_from_credential_id,omitempty" url:"copy_values_from_credential_id,omitempty"`
}

// Identifier returns the resource ID.
func (r RemoteServerCredential) Identifier() interface{} {
	return r.Id
}

// RemoteServerCredentialCollection is a list of RemoteServerCredential resources.
type RemoteServerCredentialCollection []RemoteServerCredential

// RemoteServerCredentialServerTypeEnum is a string value for server_type.
// Enum lists the values documented by the API.
type RemoteServerCredentialServerTypeEnum string

// String returns the API parameter value.
func (u RemoteServerCredentialServerTypeEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u RemoteServerCredentialServerTypeEnum) Enum() map[string]RemoteServerCredentialServerTypeEnum {
	return map[string]RemoteServerCredentialServerTypeEnum{
		"ftp":                  RemoteServerCredentialServerTypeEnum("ftp"),
		"sftp":                 RemoteServerCredentialServerTypeEnum("sftp"),
		"s3":                   RemoteServerCredentialServerTypeEnum("s3"),
		"google_cloud_storage": RemoteServerCredentialServerTypeEnum("google_cloud_storage"),
		"webdav":               RemoteServerCredentialServerTypeEnum("webdav"),
		"wasabi":               RemoteServerCredentialServerTypeEnum("wasabi"),
		"backblaze_b2":         RemoteServerCredentialServerTypeEnum("backblaze_b2"),
		"one_drive":            RemoteServerCredentialServerTypeEnum("one_drive"),
		"box":                  RemoteServerCredentialServerTypeEnum("box"),
		"dropbox":              RemoteServerCredentialServerTypeEnum("dropbox"),
		"google_drive":         RemoteServerCredentialServerTypeEnum("google_drive"),
		"azure":                RemoteServerCredentialServerTypeEnum("azure"),
		"sharepoint":           RemoteServerCredentialServerTypeEnum("sharepoint"),
		"s3_compatible":        RemoteServerCredentialServerTypeEnum("s3_compatible"),
		"azure_files":          RemoteServerCredentialServerTypeEnum("azure_files"),
		"files_agent":          RemoteServerCredentialServerTypeEnum("files_agent"),
		"filebase":             RemoteServerCredentialServerTypeEnum("filebase"),
		"cloudflare":           RemoteServerCredentialServerTypeEnum("cloudflare"),
		"linode":               RemoteServerCredentialServerTypeEnum("linode"),
		"files_com":            RemoteServerCredentialServerTypeEnum("files_com"),
	}
}

// RemoteServerCredentialListParams contains the request parameters for GET /remote_server_credentials.
type RemoteServerCredentialListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `name`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id` and `name`. Valid field combinations are `[ workspace_id, name ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is prefixed by the supplied value. Valid fields are `name`.
	FilterPrefix interface{} `url:"filter_prefix,omitempty" json:"filter_prefix,omitempty" path:"filter_prefix"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// RemoteServerCredentialFindParams contains the request parameters for GET /remote_server_credentials/{id}.
type RemoteServerCredentialFindParams struct {
	// Remote Server Credential ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// RemoteServerCredentialCreateParams contains the request parameters for POST /remote_server_credentials.
type RemoteServerCredentialCreateParams struct {
	// Internal name for your reference.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Internal description for your reference.
	Description string `url:"description,omitempty" json:"description,omitempty" path:"description"`
	// Remote server type. Remote Server Credentials are only valid for a single type of Remote Server. Values: ftp, sftp, s3, google_cloud_storage, webdav, wasabi, backblaze_b2, one_drive, box, dropbox, google_drive, azure, sharepoint, s3_compatible, azure_files, files_agent, filebase, cloudflare, linode, files_com.
	ServerType RemoteServerCredentialServerTypeEnum `url:"server_type,omitempty" json:"server_type,omitempty" path:"server_type"`
	// AWS Access Key.
	AwsAccessKey string `url:"aws_access_key,omitempty" json:"aws_access_key,omitempty" path:"aws_access_key"`
	// AWS IAM Role ARN for AssumeRole authentication.
	S3AssumeRoleArn string `url:"s3_assume_role_arn,omitempty" json:"s3_assume_role_arn,omitempty" path:"s3_assume_role_arn"`
	// Session duration in seconds for AssumeRole authentication (900-43200).
	S3AssumeRoleDurationSeconds int64 `url:"s3_assume_role_duration_seconds,omitempty" json:"s3_assume_role_duration_seconds,omitempty" path:"s3_assume_role_duration_seconds"`
	// Cloudflare: Access Key.
	CloudflareAccessKey string `url:"cloudflare_access_key,omitempty" json:"cloudflare_access_key,omitempty" path:"cloudflare_access_key"`
	// Filebase: Access Key.
	FilebaseAccessKey string `url:"filebase_access_key,omitempty" json:"filebase_access_key,omitempty" path:"filebase_access_key"`
	// Google Cloud Storage: S3-compatible Access Key.
	GoogleCloudStorageS3CompatibleAccessKey string `url:"google_cloud_storage_s3_compatible_access_key,omitempty" json:"google_cloud_storage_s3_compatible_access_key,omitempty" path:"google_cloud_storage_s3_compatible_access_key"`
	// Linode: Access Key.
	LinodeAccessKey string `url:"linode_access_key,omitempty" json:"linode_access_key,omitempty" path:"linode_access_key"`
	// S3-compatible: Access Key.
	S3CompatibleAccessKey string `url:"s3_compatible_access_key,omitempty" json:"s3_compatible_access_key,omitempty" path:"s3_compatible_access_key"`
	// SharePoint: Microsoft Entra application client ID for app-only authentication.
	SharepointClientId string `url:"sharepoint_client_id,omitempty" json:"sharepoint_client_id,omitempty" path:"sharepoint_client_id"`
	// SharePoint: Microsoft Entra tenant ID for app-only authentication.
	SharepointTenantId string `url:"sharepoint_tenant_id,omitempty" json:"sharepoint_tenant_id,omitempty" path:"sharepoint_tenant_id"`
	// Remote server username.
	Username string `url:"username,omitempty" json:"username,omitempty" path:"username"`
	// Wasabi: Access Key.
	WasabiAccessKey string `url:"wasabi_access_key,omitempty" json:"wasabi_access_key,omitempty" path:"wasabi_access_key"`
	// Password, if needed.
	Password string `url:"password,omitempty" json:"password,omitempty" path:"password"`
	// Private key, if needed.
	PrivateKey string `url:"private_key,omitempty" json:"private_key,omitempty" path:"private_key"`
	// Passphrase for private key if needed.
	PrivateKeyPassphrase string `url:"private_key_passphrase,omitempty" json:"private_key_passphrase,omitempty" path:"private_key_passphrase"`
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
	// SharePoint: PEM-encoded certificate and unencrypted private key for app-only authentication.
	SharepointClientCertificate string `url:"sharepoint_client_certificate,omitempty" json:"sharepoint_client_certificate,omitempty" path:"sharepoint_client_certificate"`
	// SharePoint: Microsoft Entra application client secret for app-only authentication.
	SharepointClientSecret string `url:"sharepoint_client_secret,omitempty" json:"sharepoint_client_secret,omitempty" path:"sharepoint_client_secret"`
	// Wasabi: Secret Key.
	WasabiSecretKey string `url:"wasabi_secret_key,omitempty" json:"wasabi_secret_key,omitempty" path:"wasabi_secret_key"`
	// Workspace ID (0 for default workspace).
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
	// ID of Remote Server Credential to copy omitted values from.
	CopyValuesFromCredentialId int64 `url:"copy_values_from_credential_id,omitempty" json:"copy_values_from_credential_id,omitempty" path:"copy_values_from_credential_id"`
}

// RemoteServerCredentialUpdateParams contains the request parameters for PATCH /remote_server_credentials/{id}.
type RemoteServerCredentialUpdateParams struct {
	// Remote Server Credential ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Internal name for your reference.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Internal description for your reference.
	Description string `url:"description,omitempty" json:"description,omitempty" path:"description"`
	// Remote server type. Remote Server Credentials are only valid for a single type of Remote Server. Values: ftp, sftp, s3, google_cloud_storage, webdav, wasabi, backblaze_b2, one_drive, box, dropbox, google_drive, azure, sharepoint, s3_compatible, azure_files, files_agent, filebase, cloudflare, linode, files_com.
	ServerType RemoteServerCredentialServerTypeEnum `url:"server_type,omitempty" json:"server_type,omitempty" path:"server_type"`
	// AWS Access Key.
	AwsAccessKey string `url:"aws_access_key,omitempty" json:"aws_access_key,omitempty" path:"aws_access_key"`
	// AWS IAM Role ARN for AssumeRole authentication.
	S3AssumeRoleArn string `url:"s3_assume_role_arn,omitempty" json:"s3_assume_role_arn,omitempty" path:"s3_assume_role_arn"`
	// Session duration in seconds for AssumeRole authentication (900-43200).
	S3AssumeRoleDurationSeconds int64 `url:"s3_assume_role_duration_seconds,omitempty" json:"s3_assume_role_duration_seconds,omitempty" path:"s3_assume_role_duration_seconds"`
	// Cloudflare: Access Key.
	CloudflareAccessKey string `url:"cloudflare_access_key,omitempty" json:"cloudflare_access_key,omitempty" path:"cloudflare_access_key"`
	// Filebase: Access Key.
	FilebaseAccessKey string `url:"filebase_access_key,omitempty" json:"filebase_access_key,omitempty" path:"filebase_access_key"`
	// Google Cloud Storage: S3-compatible Access Key.
	GoogleCloudStorageS3CompatibleAccessKey string `url:"google_cloud_storage_s3_compatible_access_key,omitempty" json:"google_cloud_storage_s3_compatible_access_key,omitempty" path:"google_cloud_storage_s3_compatible_access_key"`
	// Linode: Access Key.
	LinodeAccessKey string `url:"linode_access_key,omitempty" json:"linode_access_key,omitempty" path:"linode_access_key"`
	// S3-compatible: Access Key.
	S3CompatibleAccessKey string `url:"s3_compatible_access_key,omitempty" json:"s3_compatible_access_key,omitempty" path:"s3_compatible_access_key"`
	// SharePoint: Microsoft Entra application client ID for app-only authentication.
	SharepointClientId string `url:"sharepoint_client_id,omitempty" json:"sharepoint_client_id,omitempty" path:"sharepoint_client_id"`
	// SharePoint: Microsoft Entra tenant ID for app-only authentication.
	SharepointTenantId string `url:"sharepoint_tenant_id,omitempty" json:"sharepoint_tenant_id,omitempty" path:"sharepoint_tenant_id"`
	// Remote server username.
	Username string `url:"username,omitempty" json:"username,omitempty" path:"username"`
	// Wasabi: Access Key.
	WasabiAccessKey string `url:"wasabi_access_key,omitempty" json:"wasabi_access_key,omitempty" path:"wasabi_access_key"`
	// Password, if needed.
	Password string `url:"password,omitempty" json:"password,omitempty" path:"password"`
	// Private key, if needed.
	PrivateKey string `url:"private_key,omitempty" json:"private_key,omitempty" path:"private_key"`
	// Passphrase for private key if needed.
	PrivateKeyPassphrase string `url:"private_key_passphrase,omitempty" json:"private_key_passphrase,omitempty" path:"private_key_passphrase"`
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
	// SharePoint: PEM-encoded certificate and unencrypted private key for app-only authentication.
	SharepointClientCertificate string `url:"sharepoint_client_certificate,omitempty" json:"sharepoint_client_certificate,omitempty" path:"sharepoint_client_certificate"`
	// SharePoint: Microsoft Entra application client secret for app-only authentication.
	SharepointClientSecret string `url:"sharepoint_client_secret,omitempty" json:"sharepoint_client_secret,omitempty" path:"sharepoint_client_secret"`
	// Wasabi: Secret Key.
	WasabiSecretKey string `url:"wasabi_secret_key,omitempty" json:"wasabi_secret_key,omitempty" path:"wasabi_secret_key"`
}

// RemoteServerCredentialDeleteParams contains the request parameters for DELETE /remote_server_credentials/{id}.
type RemoteServerCredentialDeleteParams struct {
	// Remote Server Credential ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (r *RemoteServerCredential) UnmarshalJSON(data []byte) error {
	type remoteServerCredential RemoteServerCredential
	var v remoteServerCredential
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*r = RemoteServerCredential(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (r *RemoteServerCredentialCollection) UnmarshalJSON(data []byte) error {
	type remoteServerCredentials RemoteServerCredentialCollection
	var v remoteServerCredentials
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*r = RemoteServerCredentialCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (r *RemoteServerCredentialCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*r))
	for i, v := range *r {
		ret[i] = v
	}

	return &ret
}
