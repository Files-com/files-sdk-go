package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type RemoteServer struct {
	Id                                int64  `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Disabled                          *bool  `json:"disabled,omitempty" path:"disabled,omitempty" url:"disabled,omitempty"`
	AuthenticationMethod              string `json:"authentication_method,omitempty" path:"authentication_method,omitempty" url:"authentication_method,omitempty"`
	Hostname                          string `json:"hostname,omitempty" path:"hostname,omitempty" url:"hostname,omitempty"`
	RemoteHomePath                    string `json:"remote_home_path,omitempty" path:"remote_home_path,omitempty" url:"remote_home_path,omitempty"`
	Name                              string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Port                              int64  `json:"port,omitempty" path:"port,omitempty" url:"port,omitempty"`
	MaxConnections                    int64  `json:"max_connections,omitempty" path:"max_connections,omitempty" url:"max_connections,omitempty"`
	PinToSiteRegion                   *bool  `json:"pin_to_site_region,omitempty" path:"pin_to_site_region,omitempty" url:"pin_to_site_region,omitempty"`
	PinnedRegion                      string `json:"pinned_region,omitempty" path:"pinned_region,omitempty" url:"pinned_region,omitempty"`
	S3Bucket                          string `json:"s3_bucket,omitempty" path:"s3_bucket,omitempty" url:"s3_bucket,omitempty"`
	S3Region                          string `json:"s3_region,omitempty" path:"s3_region,omitempty" url:"s3_region,omitempty"`
	AwsAccessKey                      string `json:"aws_access_key,omitempty" path:"aws_access_key,omitempty" url:"aws_access_key,omitempty"`
	ServerCertificate                 string `json:"server_certificate,omitempty" path:"server_certificate,omitempty" url:"server_certificate,omitempty"`
	ServerHostKey                     string `json:"server_host_key,omitempty" path:"server_host_key,omitempty" url:"server_host_key,omitempty"`
	ServerType                        string `json:"server_type,omitempty" path:"server_type,omitempty" url:"server_type,omitempty"`
	Ssl                               string `json:"ssl,omitempty" path:"ssl,omitempty" url:"ssl,omitempty"`
	Username                          string `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	GoogleCloudStorageBucket          string `json:"google_cloud_storage_bucket,omitempty" path:"google_cloud_storage_bucket,omitempty" url:"google_cloud_storage_bucket,omitempty"`
	GoogleCloudStorageProjectId       string `json:"google_cloud_storage_project_id,omitempty" path:"google_cloud_storage_project_id,omitempty" url:"google_cloud_storage_project_id,omitempty"`
	BackblazeB2S3Endpoint             string `json:"backblaze_b2_s3_endpoint,omitempty" path:"backblaze_b2_s3_endpoint,omitempty" url:"backblaze_b2_s3_endpoint,omitempty"`
	BackblazeB2Bucket                 string `json:"backblaze_b2_bucket,omitempty" path:"backblaze_b2_bucket,omitempty" url:"backblaze_b2_bucket,omitempty"`
	WasabiBucket                      string `json:"wasabi_bucket,omitempty" path:"wasabi_bucket,omitempty" url:"wasabi_bucket,omitempty"`
	WasabiRegion                      string `json:"wasabi_region,omitempty" path:"wasabi_region,omitempty" url:"wasabi_region,omitempty"`
	WasabiAccessKey                   string `json:"wasabi_access_key,omitempty" path:"wasabi_access_key,omitempty" url:"wasabi_access_key,omitempty"`
	RackspaceUsername                 string `json:"rackspace_username,omitempty" path:"rackspace_username,omitempty" url:"rackspace_username,omitempty"`
	RackspaceRegion                   string `json:"rackspace_region,omitempty" path:"rackspace_region,omitempty" url:"rackspace_region,omitempty"`
	RackspaceContainer                string `json:"rackspace_container,omitempty" path:"rackspace_container,omitempty" url:"rackspace_container,omitempty"`
	AuthSetupLink                     string `json:"auth_setup_link,omitempty" path:"auth_setup_link,omitempty" url:"auth_setup_link,omitempty"`
	AuthStatus                        string `json:"auth_status,omitempty" path:"auth_status,omitempty" url:"auth_status,omitempty"`
	AuthAccountName                   string `json:"auth_account_name,omitempty" path:"auth_account_name,omitempty" url:"auth_account_name,omitempty"`
	OneDriveAccountType               string `json:"one_drive_account_type,omitempty" path:"one_drive_account_type,omitempty" url:"one_drive_account_type,omitempty"`
	AzureBlobStorageAccount           string `json:"azure_blob_storage_account,omitempty" path:"azure_blob_storage_account,omitempty" url:"azure_blob_storage_account,omitempty"`
	AzureBlobStorageSasToken          string `json:"azure_blob_storage_sas_token,omitempty" path:"azure_blob_storage_sas_token,omitempty" url:"azure_blob_storage_sas_token,omitempty"`
	AzureBlobStorageContainer         string `json:"azure_blob_storage_container,omitempty" path:"azure_blob_storage_container,omitempty" url:"azure_blob_storage_container,omitempty"`
	AzureFilesStorageAccount          string `json:"azure_files_storage_account,omitempty" path:"azure_files_storage_account,omitempty" url:"azure_files_storage_account,omitempty"`
	AzureFilesStorageSasToken         string `json:"azure_files_storage_sas_token,omitempty" path:"azure_files_storage_sas_token,omitempty" url:"azure_files_storage_sas_token,omitempty"`
	AzureFilesStorageShareName        string `json:"azure_files_storage_share_name,omitempty" path:"azure_files_storage_share_name,omitempty" url:"azure_files_storage_share_name,omitempty"`
	S3CompatibleBucket                string `json:"s3_compatible_bucket,omitempty" path:"s3_compatible_bucket,omitempty" url:"s3_compatible_bucket,omitempty"`
	S3CompatibleEndpoint              string `json:"s3_compatible_endpoint,omitempty" path:"s3_compatible_endpoint,omitempty" url:"s3_compatible_endpoint,omitempty"`
	S3CompatibleRegion                string `json:"s3_compatible_region,omitempty" path:"s3_compatible_region,omitempty" url:"s3_compatible_region,omitempty"`
	S3CompatibleAccessKey             string `json:"s3_compatible_access_key,omitempty" path:"s3_compatible_access_key,omitempty" url:"s3_compatible_access_key,omitempty"`
	EnableDedicatedIps                *bool  `json:"enable_dedicated_ips,omitempty" path:"enable_dedicated_ips,omitempty" url:"enable_dedicated_ips,omitempty"`
	FilesAgentPermissionSet           string `json:"files_agent_permission_set,omitempty" path:"files_agent_permission_set,omitempty" url:"files_agent_permission_set,omitempty"`
	FilesAgentRoot                    string `json:"files_agent_root,omitempty" path:"files_agent_root,omitempty" url:"files_agent_root,omitempty"`
	FilesAgentApiToken                string `json:"files_agent_api_token,omitempty" path:"files_agent_api_token,omitempty" url:"files_agent_api_token,omitempty"`
	FilebaseBucket                    string `json:"filebase_bucket,omitempty" path:"filebase_bucket,omitempty" url:"filebase_bucket,omitempty"`
	FilebaseAccessKey                 string `json:"filebase_access_key,omitempty" path:"filebase_access_key,omitempty" url:"filebase_access_key,omitempty"`
	CloudflareBucket                  string `json:"cloudflare_bucket,omitempty" path:"cloudflare_bucket,omitempty" url:"cloudflare_bucket,omitempty"`
	CloudflareAccessKey               string `json:"cloudflare_access_key,omitempty" path:"cloudflare_access_key,omitempty" url:"cloudflare_access_key,omitempty"`
	CloudflareEndpoint                string `json:"cloudflare_endpoint,omitempty" path:"cloudflare_endpoint,omitempty" url:"cloudflare_endpoint,omitempty"`
	DropboxTeams                      *bool  `json:"dropbox_teams,omitempty" path:"dropbox_teams,omitempty" url:"dropbox_teams,omitempty"`
	LinodeBucket                      string `json:"linode_bucket,omitempty" path:"linode_bucket,omitempty" url:"linode_bucket,omitempty"`
	LinodeAccessKey                   string `json:"linode_access_key,omitempty" path:"linode_access_key,omitempty" url:"linode_access_key,omitempty"`
	LinodeRegion                      string `json:"linode_region,omitempty" path:"linode_region,omitempty" url:"linode_region,omitempty"`
	AwsSecretKey                      string `json:"aws_secret_key,omitempty" path:"aws_secret_key,omitempty" url:"aws_secret_key,omitempty"`
	Password                          string `json:"password,omitempty" path:"password,omitempty" url:"password,omitempty"`
	PrivateKey                        string `json:"private_key,omitempty" path:"private_key,omitempty" url:"private_key,omitempty"`
	PrivateKeyPassphrase              string `json:"private_key_passphrase,omitempty" path:"private_key_passphrase,omitempty" url:"private_key_passphrase,omitempty"`
	SslCertificate                    string `json:"ssl_certificate,omitempty" path:"ssl_certificate,omitempty" url:"ssl_certificate,omitempty"`
	GoogleCloudStorageCredentialsJson string `json:"google_cloud_storage_credentials_json,omitempty" path:"google_cloud_storage_credentials_json,omitempty" url:"google_cloud_storage_credentials_json,omitempty"`
	WasabiSecretKey                   string `json:"wasabi_secret_key,omitempty" path:"wasabi_secret_key,omitempty" url:"wasabi_secret_key,omitempty"`
	BackblazeB2KeyId                  string `json:"backblaze_b2_key_id,omitempty" path:"backblaze_b2_key_id,omitempty" url:"backblaze_b2_key_id,omitempty"`
	BackblazeB2ApplicationKey         string `json:"backblaze_b2_application_key,omitempty" path:"backblaze_b2_application_key,omitempty" url:"backblaze_b2_application_key,omitempty"`
	RackspaceApiKey                   string `json:"rackspace_api_key,omitempty" path:"rackspace_api_key,omitempty" url:"rackspace_api_key,omitempty"`
	ResetAuthentication               *bool  `json:"reset_authentication,omitempty" path:"reset_authentication,omitempty" url:"reset_authentication,omitempty"`
	AzureBlobStorageAccessKey         string `json:"azure_blob_storage_access_key,omitempty" path:"azure_blob_storage_access_key,omitempty" url:"azure_blob_storage_access_key,omitempty"`
	AzureFilesStorageAccessKey        string `json:"azure_files_storage_access_key,omitempty" path:"azure_files_storage_access_key,omitempty" url:"azure_files_storage_access_key,omitempty"`
	S3CompatibleSecretKey             string `json:"s3_compatible_secret_key,omitempty" path:"s3_compatible_secret_key,omitempty" url:"s3_compatible_secret_key,omitempty"`
	FilebaseSecretKey                 string `json:"filebase_secret_key,omitempty" path:"filebase_secret_key,omitempty" url:"filebase_secret_key,omitempty"`
	CloudflareSecretKey               string `json:"cloudflare_secret_key,omitempty" path:"cloudflare_secret_key,omitempty" url:"cloudflare_secret_key,omitempty"`
	LinodeSecretKey                   string `json:"linode_secret_key,omitempty" path:"linode_secret_key,omitempty" url:"linode_secret_key,omitempty"`
}

func (r RemoteServer) Identifier() interface{} {
	return r.Id
}

type RemoteServerCollection []RemoteServer

type RemoteServerServerCertificateEnum string

func (u RemoteServerServerCertificateEnum) String() string {
	return string(u)
}

func (u RemoteServerServerCertificateEnum) Enum() map[string]RemoteServerServerCertificateEnum {
	return map[string]RemoteServerServerCertificateEnum{
		"require_match": RemoteServerServerCertificateEnum("require_match"),
		"allow_any":     RemoteServerServerCertificateEnum("allow_any"),
	}
}

type RemoteServerServerTypeEnum string

func (u RemoteServerServerTypeEnum) String() string {
	return string(u)
}

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
		"rackspace":            RemoteServerServerTypeEnum("rackspace"),
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
	}
}

type RemoteServerSslEnum string

func (u RemoteServerSslEnum) String() string {
	return string(u)
}

func (u RemoteServerSslEnum) Enum() map[string]RemoteServerSslEnum {
	return map[string]RemoteServerSslEnum{
		"if_available":     RemoteServerSslEnum("if_available"),
		"require":          RemoteServerSslEnum("require"),
		"require_implicit": RemoteServerSslEnum("require_implicit"),
		"never":            RemoteServerSslEnum("never"),
	}
}

type RemoteServerOneDriveAccountTypeEnum string

func (u RemoteServerOneDriveAccountTypeEnum) String() string {
	return string(u)
}

func (u RemoteServerOneDriveAccountTypeEnum) Enum() map[string]RemoteServerOneDriveAccountTypeEnum {
	return map[string]RemoteServerOneDriveAccountTypeEnum{
		"personal":       RemoteServerOneDriveAccountTypeEnum("personal"),
		"business_other": RemoteServerOneDriveAccountTypeEnum("business_other"),
	}
}

type RemoteServerFilesAgentPermissionSetEnum string

func (u RemoteServerFilesAgentPermissionSetEnum) String() string {
	return string(u)
}

func (u RemoteServerFilesAgentPermissionSetEnum) Enum() map[string]RemoteServerFilesAgentPermissionSetEnum {
	return map[string]RemoteServerFilesAgentPermissionSetEnum{
		"read_write": RemoteServerFilesAgentPermissionSetEnum("read_write"),
		"read_only":  RemoteServerFilesAgentPermissionSetEnum("read_only"),
		"write_only": RemoteServerFilesAgentPermissionSetEnum("write_only"),
	}
}

type RemoteServerListParams struct {
	ListParams
}

type RemoteServerFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type RemoteServerFindConfigurationFileParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type RemoteServerCreateParams struct {
	AwsAccessKey                      string                                  `url:"aws_access_key,omitempty" required:"false" json:"aws_access_key,omitempty" path:"aws_access_key"`
	AwsSecretKey                      string                                  `url:"aws_secret_key,omitempty" required:"false" json:"aws_secret_key,omitempty" path:"aws_secret_key"`
	Password                          string                                  `url:"password,omitempty" required:"false" json:"password,omitempty" path:"password"`
	PrivateKey                        string                                  `url:"private_key,omitempty" required:"false" json:"private_key,omitempty" path:"private_key"`
	PrivateKeyPassphrase              string                                  `url:"private_key_passphrase,omitempty" required:"false" json:"private_key_passphrase,omitempty" path:"private_key_passphrase"`
	SslCertificate                    string                                  `url:"ssl_certificate,omitempty" required:"false" json:"ssl_certificate,omitempty" path:"ssl_certificate"`
	GoogleCloudStorageCredentialsJson string                                  `url:"google_cloud_storage_credentials_json,omitempty" required:"false" json:"google_cloud_storage_credentials_json,omitempty" path:"google_cloud_storage_credentials_json"`
	WasabiAccessKey                   string                                  `url:"wasabi_access_key,omitempty" required:"false" json:"wasabi_access_key,omitempty" path:"wasabi_access_key"`
	WasabiSecretKey                   string                                  `url:"wasabi_secret_key,omitempty" required:"false" json:"wasabi_secret_key,omitempty" path:"wasabi_secret_key"`
	BackblazeB2KeyId                  string                                  `url:"backblaze_b2_key_id,omitempty" required:"false" json:"backblaze_b2_key_id,omitempty" path:"backblaze_b2_key_id"`
	BackblazeB2ApplicationKey         string                                  `url:"backblaze_b2_application_key,omitempty" required:"false" json:"backblaze_b2_application_key,omitempty" path:"backblaze_b2_application_key"`
	RackspaceApiKey                   string                                  `url:"rackspace_api_key,omitempty" required:"false" json:"rackspace_api_key,omitempty" path:"rackspace_api_key"`
	ResetAuthentication               *bool                                   `url:"reset_authentication,omitempty" required:"false" json:"reset_authentication,omitempty" path:"reset_authentication"`
	AzureBlobStorageAccessKey         string                                  `url:"azure_blob_storage_access_key,omitempty" required:"false" json:"azure_blob_storage_access_key,omitempty" path:"azure_blob_storage_access_key"`
	AzureFilesStorageAccessKey        string                                  `url:"azure_files_storage_access_key,omitempty" required:"false" json:"azure_files_storage_access_key,omitempty" path:"azure_files_storage_access_key"`
	Hostname                          string                                  `url:"hostname,omitempty" required:"false" json:"hostname,omitempty" path:"hostname"`
	Name                              string                                  `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	MaxConnections                    int64                                   `url:"max_connections,omitempty" required:"false" json:"max_connections,omitempty" path:"max_connections"`
	PinToSiteRegion                   *bool                                   `url:"pin_to_site_region,omitempty" required:"false" json:"pin_to_site_region,omitempty" path:"pin_to_site_region"`
	Port                              int64                                   `url:"port,omitempty" required:"false" json:"port,omitempty" path:"port"`
	S3Bucket                          string                                  `url:"s3_bucket,omitempty" required:"false" json:"s3_bucket,omitempty" path:"s3_bucket"`
	S3Region                          string                                  `url:"s3_region,omitempty" required:"false" json:"s3_region,omitempty" path:"s3_region"`
	ServerCertificate                 RemoteServerServerCertificateEnum       `url:"server_certificate,omitempty" required:"false" json:"server_certificate,omitempty" path:"server_certificate"`
	ServerHostKey                     string                                  `url:"server_host_key,omitempty" required:"false" json:"server_host_key,omitempty" path:"server_host_key"`
	ServerType                        RemoteServerServerTypeEnum              `url:"server_type,omitempty" required:"false" json:"server_type,omitempty" path:"server_type"`
	Ssl                               RemoteServerSslEnum                     `url:"ssl,omitempty" required:"false" json:"ssl,omitempty" path:"ssl"`
	Username                          string                                  `url:"username,omitempty" required:"false" json:"username,omitempty" path:"username"`
	GoogleCloudStorageBucket          string                                  `url:"google_cloud_storage_bucket,omitempty" required:"false" json:"google_cloud_storage_bucket,omitempty" path:"google_cloud_storage_bucket"`
	GoogleCloudStorageProjectId       string                                  `url:"google_cloud_storage_project_id,omitempty" required:"false" json:"google_cloud_storage_project_id,omitempty" path:"google_cloud_storage_project_id"`
	BackblazeB2Bucket                 string                                  `url:"backblaze_b2_bucket,omitempty" required:"false" json:"backblaze_b2_bucket,omitempty" path:"backblaze_b2_bucket"`
	BackblazeB2S3Endpoint             string                                  `url:"backblaze_b2_s3_endpoint,omitempty" required:"false" json:"backblaze_b2_s3_endpoint,omitempty" path:"backblaze_b2_s3_endpoint"`
	WasabiBucket                      string                                  `url:"wasabi_bucket,omitempty" required:"false" json:"wasabi_bucket,omitempty" path:"wasabi_bucket"`
	WasabiRegion                      string                                  `url:"wasabi_region,omitempty" required:"false" json:"wasabi_region,omitempty" path:"wasabi_region"`
	RackspaceUsername                 string                                  `url:"rackspace_username,omitempty" required:"false" json:"rackspace_username,omitempty" path:"rackspace_username"`
	RackspaceRegion                   string                                  `url:"rackspace_region,omitempty" required:"false" json:"rackspace_region,omitempty" path:"rackspace_region"`
	RackspaceContainer                string                                  `url:"rackspace_container,omitempty" required:"false" json:"rackspace_container,omitempty" path:"rackspace_container"`
	OneDriveAccountType               RemoteServerOneDriveAccountTypeEnum     `url:"one_drive_account_type,omitempty" required:"false" json:"one_drive_account_type,omitempty" path:"one_drive_account_type"`
	AzureBlobStorageAccount           string                                  `url:"azure_blob_storage_account,omitempty" required:"false" json:"azure_blob_storage_account,omitempty" path:"azure_blob_storage_account"`
	AzureBlobStorageContainer         string                                  `url:"azure_blob_storage_container,omitempty" required:"false" json:"azure_blob_storage_container,omitempty" path:"azure_blob_storage_container"`
	AzureBlobStorageSasToken          string                                  `url:"azure_blob_storage_sas_token,omitempty" required:"false" json:"azure_blob_storage_sas_token,omitempty" path:"azure_blob_storage_sas_token"`
	AzureFilesStorageAccount          string                                  `url:"azure_files_storage_account,omitempty" required:"false" json:"azure_files_storage_account,omitempty" path:"azure_files_storage_account"`
	AzureFilesStorageShareName        string                                  `url:"azure_files_storage_share_name,omitempty" required:"false" json:"azure_files_storage_share_name,omitempty" path:"azure_files_storage_share_name"`
	AzureFilesStorageSasToken         string                                  `url:"azure_files_storage_sas_token,omitempty" required:"false" json:"azure_files_storage_sas_token,omitempty" path:"azure_files_storage_sas_token"`
	S3CompatibleBucket                string                                  `url:"s3_compatible_bucket,omitempty" required:"false" json:"s3_compatible_bucket,omitempty" path:"s3_compatible_bucket"`
	S3CompatibleEndpoint              string                                  `url:"s3_compatible_endpoint,omitempty" required:"false" json:"s3_compatible_endpoint,omitempty" path:"s3_compatible_endpoint"`
	S3CompatibleRegion                string                                  `url:"s3_compatible_region,omitempty" required:"false" json:"s3_compatible_region,omitempty" path:"s3_compatible_region"`
	EnableDedicatedIps                *bool                                   `url:"enable_dedicated_ips,omitempty" required:"false" json:"enable_dedicated_ips,omitempty" path:"enable_dedicated_ips"`
	S3CompatibleAccessKey             string                                  `url:"s3_compatible_access_key,omitempty" required:"false" json:"s3_compatible_access_key,omitempty" path:"s3_compatible_access_key"`
	S3CompatibleSecretKey             string                                  `url:"s3_compatible_secret_key,omitempty" required:"false" json:"s3_compatible_secret_key,omitempty" path:"s3_compatible_secret_key"`
	FilesAgentRoot                    string                                  `url:"files_agent_root,omitempty" required:"false" json:"files_agent_root,omitempty" path:"files_agent_root"`
	FilesAgentPermissionSet           RemoteServerFilesAgentPermissionSetEnum `url:"files_agent_permission_set,omitempty" required:"false" json:"files_agent_permission_set,omitempty" path:"files_agent_permission_set"`
	FilebaseAccessKey                 string                                  `url:"filebase_access_key,omitempty" required:"false" json:"filebase_access_key,omitempty" path:"filebase_access_key"`
	FilebaseSecretKey                 string                                  `url:"filebase_secret_key,omitempty" required:"false" json:"filebase_secret_key,omitempty" path:"filebase_secret_key"`
	FilebaseBucket                    string                                  `url:"filebase_bucket,omitempty" required:"false" json:"filebase_bucket,omitempty" path:"filebase_bucket"`
	CloudflareAccessKey               string                                  `url:"cloudflare_access_key,omitempty" required:"false" json:"cloudflare_access_key,omitempty" path:"cloudflare_access_key"`
	CloudflareSecretKey               string                                  `url:"cloudflare_secret_key,omitempty" required:"false" json:"cloudflare_secret_key,omitempty" path:"cloudflare_secret_key"`
	CloudflareBucket                  string                                  `url:"cloudflare_bucket,omitempty" required:"false" json:"cloudflare_bucket,omitempty" path:"cloudflare_bucket"`
	CloudflareEndpoint                string                                  `url:"cloudflare_endpoint,omitempty" required:"false" json:"cloudflare_endpoint,omitempty" path:"cloudflare_endpoint"`
	DropboxTeams                      *bool                                   `url:"dropbox_teams,omitempty" required:"false" json:"dropbox_teams,omitempty" path:"dropbox_teams"`
	LinodeAccessKey                   string                                  `url:"linode_access_key,omitempty" required:"false" json:"linode_access_key,omitempty" path:"linode_access_key"`
	LinodeSecretKey                   string                                  `url:"linode_secret_key,omitempty" required:"false" json:"linode_secret_key,omitempty" path:"linode_secret_key"`
	LinodeBucket                      string                                  `url:"linode_bucket,omitempty" required:"false" json:"linode_bucket,omitempty" path:"linode_bucket"`
	LinodeRegion                      string                                  `url:"linode_region,omitempty" required:"false" json:"linode_region,omitempty" path:"linode_region"`
}

// Post local changes, check in, and download configuration file (used by some Remote Server integrations, such as the Files.com Agent)
type RemoteServerConfigurationFileParams struct {
	Id            int64  `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
	ApiToken      string `url:"api_token,omitempty" required:"false" json:"api_token,omitempty" path:"api_token"`
	PermissionSet string `url:"permission_set,omitempty" required:"false" json:"permission_set,omitempty" path:"permission_set"`
	Root          string `url:"root,omitempty" required:"false" json:"root,omitempty" path:"root"`
	Hostname      string `url:"hostname,omitempty" required:"false" json:"hostname,omitempty" path:"hostname"`
	Port          int64  `url:"port,omitempty" required:"false" json:"port,omitempty" path:"port"`
	Status        string `url:"status,omitempty" required:"false" json:"status,omitempty" path:"status"`
	ConfigVersion string `url:"config_version,omitempty" required:"false" json:"config_version,omitempty" path:"config_version"`
	PrivateKey    string `url:"private_key,omitempty" required:"false" json:"private_key,omitempty" path:"private_key"`
	PublicKey     string `url:"public_key,omitempty" required:"false" json:"public_key,omitempty" path:"public_key"`
	ServerHostKey string `url:"server_host_key,omitempty" required:"false" json:"server_host_key,omitempty" path:"server_host_key"`
	Subdomain     string `url:"subdomain,omitempty" required:"false" json:"subdomain,omitempty" path:"subdomain"`
}

type RemoteServerUpdateParams struct {
	Id                                int64                                   `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
	AwsAccessKey                      string                                  `url:"aws_access_key,omitempty" required:"false" json:"aws_access_key,omitempty" path:"aws_access_key"`
	AwsSecretKey                      string                                  `url:"aws_secret_key,omitempty" required:"false" json:"aws_secret_key,omitempty" path:"aws_secret_key"`
	Password                          string                                  `url:"password,omitempty" required:"false" json:"password,omitempty" path:"password"`
	PrivateKey                        string                                  `url:"private_key,omitempty" required:"false" json:"private_key,omitempty" path:"private_key"`
	PrivateKeyPassphrase              string                                  `url:"private_key_passphrase,omitempty" required:"false" json:"private_key_passphrase,omitempty" path:"private_key_passphrase"`
	SslCertificate                    string                                  `url:"ssl_certificate,omitempty" required:"false" json:"ssl_certificate,omitempty" path:"ssl_certificate"`
	GoogleCloudStorageCredentialsJson string                                  `url:"google_cloud_storage_credentials_json,omitempty" required:"false" json:"google_cloud_storage_credentials_json,omitempty" path:"google_cloud_storage_credentials_json"`
	WasabiAccessKey                   string                                  `url:"wasabi_access_key,omitempty" required:"false" json:"wasabi_access_key,omitempty" path:"wasabi_access_key"`
	WasabiSecretKey                   string                                  `url:"wasabi_secret_key,omitempty" required:"false" json:"wasabi_secret_key,omitempty" path:"wasabi_secret_key"`
	BackblazeB2KeyId                  string                                  `url:"backblaze_b2_key_id,omitempty" required:"false" json:"backblaze_b2_key_id,omitempty" path:"backblaze_b2_key_id"`
	BackblazeB2ApplicationKey         string                                  `url:"backblaze_b2_application_key,omitempty" required:"false" json:"backblaze_b2_application_key,omitempty" path:"backblaze_b2_application_key"`
	RackspaceApiKey                   string                                  `url:"rackspace_api_key,omitempty" required:"false" json:"rackspace_api_key,omitempty" path:"rackspace_api_key"`
	ResetAuthentication               *bool                                   `url:"reset_authentication,omitempty" required:"false" json:"reset_authentication,omitempty" path:"reset_authentication"`
	AzureBlobStorageAccessKey         string                                  `url:"azure_blob_storage_access_key,omitempty" required:"false" json:"azure_blob_storage_access_key,omitempty" path:"azure_blob_storage_access_key"`
	AzureFilesStorageAccessKey        string                                  `url:"azure_files_storage_access_key,omitempty" required:"false" json:"azure_files_storage_access_key,omitempty" path:"azure_files_storage_access_key"`
	Hostname                          string                                  `url:"hostname,omitempty" required:"false" json:"hostname,omitempty" path:"hostname"`
	Name                              string                                  `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	MaxConnections                    int64                                   `url:"max_connections,omitempty" required:"false" json:"max_connections,omitempty" path:"max_connections"`
	PinToSiteRegion                   *bool                                   `url:"pin_to_site_region,omitempty" required:"false" json:"pin_to_site_region,omitempty" path:"pin_to_site_region"`
	Port                              int64                                   `url:"port,omitempty" required:"false" json:"port,omitempty" path:"port"`
	S3Bucket                          string                                  `url:"s3_bucket,omitempty" required:"false" json:"s3_bucket,omitempty" path:"s3_bucket"`
	S3Region                          string                                  `url:"s3_region,omitempty" required:"false" json:"s3_region,omitempty" path:"s3_region"`
	ServerCertificate                 RemoteServerServerCertificateEnum       `url:"server_certificate,omitempty" required:"false" json:"server_certificate,omitempty" path:"server_certificate"`
	ServerHostKey                     string                                  `url:"server_host_key,omitempty" required:"false" json:"server_host_key,omitempty" path:"server_host_key"`
	ServerType                        RemoteServerServerTypeEnum              `url:"server_type,omitempty" required:"false" json:"server_type,omitempty" path:"server_type"`
	Ssl                               RemoteServerSslEnum                     `url:"ssl,omitempty" required:"false" json:"ssl,omitempty" path:"ssl"`
	Username                          string                                  `url:"username,omitempty" required:"false" json:"username,omitempty" path:"username"`
	GoogleCloudStorageBucket          string                                  `url:"google_cloud_storage_bucket,omitempty" required:"false" json:"google_cloud_storage_bucket,omitempty" path:"google_cloud_storage_bucket"`
	GoogleCloudStorageProjectId       string                                  `url:"google_cloud_storage_project_id,omitempty" required:"false" json:"google_cloud_storage_project_id,omitempty" path:"google_cloud_storage_project_id"`
	BackblazeB2Bucket                 string                                  `url:"backblaze_b2_bucket,omitempty" required:"false" json:"backblaze_b2_bucket,omitempty" path:"backblaze_b2_bucket"`
	BackblazeB2S3Endpoint             string                                  `url:"backblaze_b2_s3_endpoint,omitempty" required:"false" json:"backblaze_b2_s3_endpoint,omitempty" path:"backblaze_b2_s3_endpoint"`
	WasabiBucket                      string                                  `url:"wasabi_bucket,omitempty" required:"false" json:"wasabi_bucket,omitempty" path:"wasabi_bucket"`
	WasabiRegion                      string                                  `url:"wasabi_region,omitempty" required:"false" json:"wasabi_region,omitempty" path:"wasabi_region"`
	RackspaceUsername                 string                                  `url:"rackspace_username,omitempty" required:"false" json:"rackspace_username,omitempty" path:"rackspace_username"`
	RackspaceRegion                   string                                  `url:"rackspace_region,omitempty" required:"false" json:"rackspace_region,omitempty" path:"rackspace_region"`
	RackspaceContainer                string                                  `url:"rackspace_container,omitempty" required:"false" json:"rackspace_container,omitempty" path:"rackspace_container"`
	OneDriveAccountType               RemoteServerOneDriveAccountTypeEnum     `url:"one_drive_account_type,omitempty" required:"false" json:"one_drive_account_type,omitempty" path:"one_drive_account_type"`
	AzureBlobStorageAccount           string                                  `url:"azure_blob_storage_account,omitempty" required:"false" json:"azure_blob_storage_account,omitempty" path:"azure_blob_storage_account"`
	AzureBlobStorageContainer         string                                  `url:"azure_blob_storage_container,omitempty" required:"false" json:"azure_blob_storage_container,omitempty" path:"azure_blob_storage_container"`
	AzureBlobStorageSasToken          string                                  `url:"azure_blob_storage_sas_token,omitempty" required:"false" json:"azure_blob_storage_sas_token,omitempty" path:"azure_blob_storage_sas_token"`
	AzureFilesStorageAccount          string                                  `url:"azure_files_storage_account,omitempty" required:"false" json:"azure_files_storage_account,omitempty" path:"azure_files_storage_account"`
	AzureFilesStorageShareName        string                                  `url:"azure_files_storage_share_name,omitempty" required:"false" json:"azure_files_storage_share_name,omitempty" path:"azure_files_storage_share_name"`
	AzureFilesStorageSasToken         string                                  `url:"azure_files_storage_sas_token,omitempty" required:"false" json:"azure_files_storage_sas_token,omitempty" path:"azure_files_storage_sas_token"`
	S3CompatibleBucket                string                                  `url:"s3_compatible_bucket,omitempty" required:"false" json:"s3_compatible_bucket,omitempty" path:"s3_compatible_bucket"`
	S3CompatibleEndpoint              string                                  `url:"s3_compatible_endpoint,omitempty" required:"false" json:"s3_compatible_endpoint,omitempty" path:"s3_compatible_endpoint"`
	S3CompatibleRegion                string                                  `url:"s3_compatible_region,omitempty" required:"false" json:"s3_compatible_region,omitempty" path:"s3_compatible_region"`
	EnableDedicatedIps                *bool                                   `url:"enable_dedicated_ips,omitempty" required:"false" json:"enable_dedicated_ips,omitempty" path:"enable_dedicated_ips"`
	S3CompatibleAccessKey             string                                  `url:"s3_compatible_access_key,omitempty" required:"false" json:"s3_compatible_access_key,omitempty" path:"s3_compatible_access_key"`
	S3CompatibleSecretKey             string                                  `url:"s3_compatible_secret_key,omitempty" required:"false" json:"s3_compatible_secret_key,omitempty" path:"s3_compatible_secret_key"`
	FilesAgentRoot                    string                                  `url:"files_agent_root,omitempty" required:"false" json:"files_agent_root,omitempty" path:"files_agent_root"`
	FilesAgentPermissionSet           RemoteServerFilesAgentPermissionSetEnum `url:"files_agent_permission_set,omitempty" required:"false" json:"files_agent_permission_set,omitempty" path:"files_agent_permission_set"`
	FilebaseAccessKey                 string                                  `url:"filebase_access_key,omitempty" required:"false" json:"filebase_access_key,omitempty" path:"filebase_access_key"`
	FilebaseSecretKey                 string                                  `url:"filebase_secret_key,omitempty" required:"false" json:"filebase_secret_key,omitempty" path:"filebase_secret_key"`
	FilebaseBucket                    string                                  `url:"filebase_bucket,omitempty" required:"false" json:"filebase_bucket,omitempty" path:"filebase_bucket"`
	CloudflareAccessKey               string                                  `url:"cloudflare_access_key,omitempty" required:"false" json:"cloudflare_access_key,omitempty" path:"cloudflare_access_key"`
	CloudflareSecretKey               string                                  `url:"cloudflare_secret_key,omitempty" required:"false" json:"cloudflare_secret_key,omitempty" path:"cloudflare_secret_key"`
	CloudflareBucket                  string                                  `url:"cloudflare_bucket,omitempty" required:"false" json:"cloudflare_bucket,omitempty" path:"cloudflare_bucket"`
	CloudflareEndpoint                string                                  `url:"cloudflare_endpoint,omitempty" required:"false" json:"cloudflare_endpoint,omitempty" path:"cloudflare_endpoint"`
	DropboxTeams                      *bool                                   `url:"dropbox_teams,omitempty" required:"false" json:"dropbox_teams,omitempty" path:"dropbox_teams"`
	LinodeAccessKey                   string                                  `url:"linode_access_key,omitempty" required:"false" json:"linode_access_key,omitempty" path:"linode_access_key"`
	LinodeSecretKey                   string                                  `url:"linode_secret_key,omitempty" required:"false" json:"linode_secret_key,omitempty" path:"linode_secret_key"`
	LinodeBucket                      string                                  `url:"linode_bucket,omitempty" required:"false" json:"linode_bucket,omitempty" path:"linode_bucket"`
	LinodeRegion                      string                                  `url:"linode_region,omitempty" required:"false" json:"linode_region,omitempty" path:"linode_region"`
}

type RemoteServerDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

func (r *RemoteServer) UnmarshalJSON(data []byte) error {
	type remoteServer RemoteServer
	var v remoteServer
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*r = RemoteServer(v)
	return nil
}

func (r *RemoteServerCollection) UnmarshalJSON(data []byte) error {
	type remoteServers RemoteServerCollection
	var v remoteServers
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*r = RemoteServerCollection(v)
	return nil
}

func (r *RemoteServerCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*r))
	for i, v := range *r {
		ret[i] = v
	}

	return &ret
}
