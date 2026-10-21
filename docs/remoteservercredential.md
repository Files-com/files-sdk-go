# RemoteServerCredential

`files_sdk.RemoteServerCredential` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/remoteservercredential`, named `remote_server_credential`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	remote_server_credential "github.com/Files-com/files-sdk-go/v3/remoteservercredential"
)
```

## Operations

### List

`GET /remote_server_credentials`

List Remote Server Credentials

```go
func (c *Client) List(params files_sdk.RemoteServerCredentialListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.RemoteServerCredentialListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.RemoteServerCredential` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.RemoteServerCredentialListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `name`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id` and `name`. Valid field combinations are `[ workspace_id, name ]`. |
| `FilterPrefix` | `interface{}` | `filter_prefix` | no | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `name`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := remote_server_credential.List(files_sdk.RemoteServerCredentialListParams{})
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

`GET /remote_server_credentials/{id}`

Show Remote Server Credential

```go
func (c *Client) Find(params files_sdk.RemoteServerCredentialFindParams, opts ...files_sdk.RequestResponseOption) (remoteServerCredential files_sdk.RemoteServerCredential, err error)
func Find(params files_sdk.RemoteServerCredentialFindParams, opts ...files_sdk.RequestResponseOption) (remoteServerCredential files_sdk.RemoteServerCredential, err error)
```

Parameters, `files_sdk.RemoteServerCredentialFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Remote Server Credential ID. |

Example:

```go
result, err := remote_server_credential.Find(files_sdk.RemoteServerCredentialFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /remote_server_credentials`

Create Remote Server Credential

```go
func (c *Client) Create(params files_sdk.RemoteServerCredentialCreateParams, opts ...files_sdk.RequestResponseOption) (remoteServerCredential files_sdk.RemoteServerCredential, err error)
func Create(params files_sdk.RemoteServerCredentialCreateParams, opts ...files_sdk.RequestResponseOption) (remoteServerCredential files_sdk.RemoteServerCredential, err error)
```

Parameters, `files_sdk.RemoteServerCredentialCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Name` | `string` | `name` | no | Internal name for your reference |
| `Description` | `string` | `description` | no | Internal description for your reference |
| `ServerType` | `files_sdk.RemoteServerCredentialServerTypeEnum` | `server_type` | no | Remote server type. Remote Server Credentials are only valid for a single type of Remote Server. Values: `ftp`, `sftp`, `s3`, `google_cloud_storage`, `webdav`, `wasabi`, `backblaze_b2`, `one_drive`, `box`, `dropbox`, `google_drive`, `azure`, `sharepoint`, `s3_compatible`, `azure_files`, `files_agent`, `filebase`, `cloudflare`, `linode`, `files_com`. |
| `AwsAccessKey` | `string` | `aws_access_key` | no | AWS Access Key. |
| `S3AssumeRoleArn` | `string` | `s3_assume_role_arn` | no | AWS IAM Role ARN for AssumeRole authentication. |
| `S3AssumeRoleDurationSeconds` | `int64` | `s3_assume_role_duration_seconds` | no | Session duration in seconds for AssumeRole authentication (900-43200). |
| `CloudflareAccessKey` | `string` | `cloudflare_access_key` | no | Cloudflare: Access Key. |
| `FilebaseAccessKey` | `string` | `filebase_access_key` | no | Filebase: Access Key. |
| `GoogleCloudStorageS3CompatibleAccessKey` | `string` | `google_cloud_storage_s3_compatible_access_key` | no | Google Cloud Storage: S3-compatible Access Key. |
| `LinodeAccessKey` | `string` | `linode_access_key` | no | Linode: Access Key |
| `S3CompatibleAccessKey` | `string` | `s3_compatible_access_key` | no | S3-compatible: Access Key |
| `SharepointClientId` | `string` | `sharepoint_client_id` | no | SharePoint: Microsoft Entra application client ID for app-only authentication. |
| `SharepointTenantId` | `string` | `sharepoint_tenant_id` | no | SharePoint: Microsoft Entra tenant ID for app-only authentication. |
| `Username` | `string` | `username` | no | Remote server username. |
| `WasabiAccessKey` | `string` | `wasabi_access_key` | no | Wasabi: Access Key. |
| `Password` | `string` | `password` | no | Password, if needed. |
| `PrivateKey` | `string` | `private_key` | no | Private key, if needed. |
| `PrivateKeyPassphrase` | `string` | `private_key_passphrase` | no | Passphrase for private key if needed. |
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
| `SharepointClientCertificate` | `string` | `sharepoint_client_certificate` | no | SharePoint: PEM-encoded certificate and unencrypted private key for app-only authentication. |
| `SharepointClientSecret` | `string` | `sharepoint_client_secret` | no | SharePoint: Microsoft Entra application client secret for app-only authentication. |
| `WasabiSecretKey` | `string` | `wasabi_secret_key` | no | Wasabi: Secret Key |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID (0 for default workspace) |
| `CopyValuesFromCredentialId` | `int64` | `copy_values_from_credential_id` | no | ID of Remote Server Credential to copy omitted values from. |

Example:

```go
result, err := remote_server_credential.Create(files_sdk.RemoteServerCredentialCreateParams{})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /remote_server_credentials/{id}`

Update Remote Server Credential

```go
func (c *Client) Update(params files_sdk.RemoteServerCredentialUpdateParams, opts ...files_sdk.RequestResponseOption) (remoteServerCredential files_sdk.RemoteServerCredential, err error)
func Update(params files_sdk.RemoteServerCredentialUpdateParams, opts ...files_sdk.RequestResponseOption) (remoteServerCredential files_sdk.RemoteServerCredential, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (remoteServerCredential files_sdk.RemoteServerCredential, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (remoteServerCredential files_sdk.RemoteServerCredential, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.RemoteServerCredentialUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Remote Server Credential ID. |
| `Name` | `string` | `name` | no | Internal name for your reference |
| `Description` | `string` | `description` | no | Internal description for your reference |
| `ServerType` | `files_sdk.RemoteServerCredentialServerTypeEnum` | `server_type` | no | Remote server type. Remote Server Credentials are only valid for a single type of Remote Server. Values: `ftp`, `sftp`, `s3`, `google_cloud_storage`, `webdav`, `wasabi`, `backblaze_b2`, `one_drive`, `box`, `dropbox`, `google_drive`, `azure`, `sharepoint`, `s3_compatible`, `azure_files`, `files_agent`, `filebase`, `cloudflare`, `linode`, `files_com`. |
| `AwsAccessKey` | `string` | `aws_access_key` | no | AWS Access Key. |
| `S3AssumeRoleArn` | `string` | `s3_assume_role_arn` | no | AWS IAM Role ARN for AssumeRole authentication. |
| `S3AssumeRoleDurationSeconds` | `int64` | `s3_assume_role_duration_seconds` | no | Session duration in seconds for AssumeRole authentication (900-43200). |
| `CloudflareAccessKey` | `string` | `cloudflare_access_key` | no | Cloudflare: Access Key. |
| `FilebaseAccessKey` | `string` | `filebase_access_key` | no | Filebase: Access Key. |
| `GoogleCloudStorageS3CompatibleAccessKey` | `string` | `google_cloud_storage_s3_compatible_access_key` | no | Google Cloud Storage: S3-compatible Access Key. |
| `LinodeAccessKey` | `string` | `linode_access_key` | no | Linode: Access Key |
| `S3CompatibleAccessKey` | `string` | `s3_compatible_access_key` | no | S3-compatible: Access Key |
| `SharepointClientId` | `string` | `sharepoint_client_id` | no | SharePoint: Microsoft Entra application client ID for app-only authentication. |
| `SharepointTenantId` | `string` | `sharepoint_tenant_id` | no | SharePoint: Microsoft Entra tenant ID for app-only authentication. |
| `Username` | `string` | `username` | no | Remote server username. |
| `WasabiAccessKey` | `string` | `wasabi_access_key` | no | Wasabi: Access Key. |
| `Password` | `string` | `password` | no | Password, if needed. |
| `PrivateKey` | `string` | `private_key` | no | Private key, if needed. |
| `PrivateKeyPassphrase` | `string` | `private_key_passphrase` | no | Passphrase for private key if needed. |
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
| `SharepointClientCertificate` | `string` | `sharepoint_client_certificate` | no | SharePoint: PEM-encoded certificate and unencrypted private key for app-only authentication. |
| `SharepointClientSecret` | `string` | `sharepoint_client_secret` | no | SharePoint: Microsoft Entra application client secret for app-only authentication. |
| `WasabiSecretKey` | `string` | `wasabi_secret_key` | no | Wasabi: Secret Key |

Example:

```go
result, err := remote_server_credential.Update(files_sdk.RemoteServerCredentialUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := remote_server_credential.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /remote_server_credentials/{id}`

Delete Remote Server Credential

```go
func (c *Client) Delete(params files_sdk.RemoteServerCredentialDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.RemoteServerCredentialDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.RemoteServerCredentialDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Remote Server Credential ID. |

Example:

```go
err := remote_server_credential.Delete(files_sdk.RemoteServerCredentialDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.RemoteServerCredential` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Remote Server Credential ID |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID (0 for default workspace) |
| `Name` | `string` | `name` | Internal name for your reference |
| `Description` | `string` | `description` | Internal description for your reference |
| `ServerType` | `string` | `server_type` | Remote server type. Remote Server Credentials are only valid for a single type of Remote Server. |
| `AwsAccessKey` | `string` | `aws_access_key` | AWS Access Key. |
| `S3AssumeRoleArn` | `string` | `s3_assume_role_arn` | AWS IAM Role ARN for AssumeRole authentication. |
| `S3AssumeRoleDurationSeconds` | `int64` | `s3_assume_role_duration_seconds` | Session duration in seconds for AssumeRole authentication (900-43200). |
| `S3AssumeRoleExternalId` | `string` | `s3_assume_role_external_id` | External ID for AssumeRole authentication. |
| `GoogleCloudStorageS3CompatibleAccessKey` | `string` | `google_cloud_storage_s3_compatible_access_key` | Google Cloud Storage: S3-compatible Access Key. |
| `WasabiAccessKey` | `string` | `wasabi_access_key` | Wasabi: Access Key. |
| `S3CompatibleAccessKey` | `string` | `s3_compatible_access_key` | S3-compatible: Access Key |
| `FilebaseAccessKey` | `string` | `filebase_access_key` | Filebase: Access Key. |
| `CloudflareAccessKey` | `string` | `cloudflare_access_key` | Cloudflare: Access Key. |
| `LinodeAccessKey` | `string` | `linode_access_key` | Linode: Access Key |
| `SharepointTenantId` | `string` | `sharepoint_tenant_id` | SharePoint: Microsoft Entra tenant ID for app-only authentication. |
| `SharepointClientId` | `string` | `sharepoint_client_id` | SharePoint: Microsoft Entra application client ID for app-only authentication. |
| `SharepointAppCredentialType` | `string` | `sharepoint_app_credential_type` | SharePoint: App-only credential type. Either secret or certificate. |
| `Username` | `string` | `username` | Remote server username. |
| `Password` | `string` | `password` | Password, if needed. |
| `PrivateKey` | `string` | `private_key` | Private key, if needed. |
| `PrivateKeyPassphrase` | `string` | `private_key_passphrase` | Passphrase for private key if needed. |
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
| `SharepointClientCertificate` | `string` | `sharepoint_client_certificate` | SharePoint: PEM-encoded certificate and unencrypted private key for app-only authentication. |
| `SharepointClientSecret` | `string` | `sharepoint_client_secret` | SharePoint: Microsoft Entra application client secret for app-only authentication. |
| `WasabiSecretKey` | `string` | `wasabi_secret_key` | Wasabi: Secret Key |
| `CopyValuesFromCredentialId` | `int64` | `copy_values_from_credential_id` | ID of Remote Server Credential to copy omitted values from. |
