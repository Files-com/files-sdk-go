package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type RemoteServer struct {
	Id                                int64  `json:"id,omitempty"`
	AuthenticationMethod              string `json:"authentication_method,omitempty"`
	Hostname                          string `json:"hostname,omitempty"`
	RemoteHomePath                    string `json:"remote_home_path,omitempty"`
	Name                              string `json:"name,omitempty"`
	Port                              int64  `json:"port,omitempty"`
	MaxConnections                    int64  `json:"max_connections,omitempty"`
	S3Bucket                          string `json:"s3_bucket,omitempty"`
	S3Region                          string `json:"s3_region,omitempty"`
	ServerCertificate                 string `json:"server_certificate,omitempty"`
	ServerHostKey                     string `json:"server_host_key,omitempty"`
	ServerType                        string `json:"server_type,omitempty"`
	Ssl                               string `json:"ssl,omitempty"`
	Username                          string `json:"username,omitempty"`
	GoogleCloudStorageBucket          string `json:"google_cloud_storage_bucket,omitempty"`
	GoogleCloudStorageProjectId       string `json:"google_cloud_storage_project_id,omitempty"`
	BackblazeB2S3Endpoint             string `json:"backblaze_b2_s3_endpoint,omitempty"`
	BackblazeB2Bucket                 string `json:"backblaze_b2_bucket,omitempty"`
	WasabiBucket                      string `json:"wasabi_bucket,omitempty"`
	WasabiRegion                      string `json:"wasabi_region,omitempty"`
	RackspaceUsername                 string `json:"rackspace_username,omitempty"`
	RackspaceRegion                   string `json:"rackspace_region,omitempty"`
	RackspaceContainer                string `json:"rackspace_container,omitempty"`
	AuthSetupLink                     string `json:"auth_setup_link,omitempty"`
	AuthStatus                        string `json:"auth_status,omitempty"`
	AuthAccountName                   string `json:"auth_account_name,omitempty"`
	OneDriveAccountType               string `json:"one_drive_account_type,omitempty"`
	AzureBlobStorageAccount           string `json:"azure_blob_storage_account,omitempty"`
	AzureBlobStorageContainer         string `json:"azure_blob_storage_container,omitempty"`
	S3CompatibleBucket                string `json:"s3_compatible_bucket,omitempty"`
	S3CompatibleRegion                string `json:"s3_compatible_region,omitempty"`
	S3CompatibleEndpoint              string `json:"s3_compatible_endpoint,omitempty"`
	AwsAccessKey                      string `json:"aws_access_key,omitempty"`
	AwsSecretKey                      string `json:"aws_secret_key,omitempty"`
	Password                          string `json:"password,omitempty"`
	PrivateKey                        string `json:"private_key,omitempty"`
	SslCertificate                    string `json:"ssl_certificate,omitempty"`
	GoogleCloudStorageCredentialsJson string `json:"google_cloud_storage_credentials_json,omitempty"`
	WasabiAccessKey                   string `json:"wasabi_access_key,omitempty"`
	WasabiSecretKey                   string `json:"wasabi_secret_key,omitempty"`
	BackblazeB2KeyId                  string `json:"backblaze_b2_key_id,omitempty"`
	BackblazeB2ApplicationKey         string `json:"backblaze_b2_application_key,omitempty"`
	RackspaceApiKey                   string `json:"rackspace_api_key,omitempty"`
	ResetAuthentication               *bool  `json:"reset_authentication,omitempty"`
	AzureBlobStorageAccessKey         string `json:"azure_blob_storage_access_key,omitempty"`
	S3CompatibleAccessKey             string `json:"s3_compatible_access_key,omitempty"`
	S3CompatibleSecretKey             string `json:"s3_compatible_secret_key,omitempty"`
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

type RemoteServerListParams struct {
	Cursor  string `url:"cursor,omitempty" required:"false"`
	PerPage int64  `url:"per_page,omitempty" required:"false"`
	lib.ListParams
}

type RemoteServerFindParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

type RemoteServerCreateParams struct {
	AwsAccessKey                      string                              `url:"aws_access_key,omitempty" required:"false"`
	AwsSecretKey                      string                              `url:"aws_secret_key,omitempty" required:"false"`
	Password                          string                              `url:"password,omitempty" required:"false"`
	PrivateKey                        string                              `url:"private_key,omitempty" required:"false"`
	SslCertificate                    string                              `url:"ssl_certificate,omitempty" required:"false"`
	GoogleCloudStorageCredentialsJson string                              `url:"google_cloud_storage_credentials_json,omitempty" required:"false"`
	WasabiAccessKey                   string                              `url:"wasabi_access_key,omitempty" required:"false"`
	WasabiSecretKey                   string                              `url:"wasabi_secret_key,omitempty" required:"false"`
	BackblazeB2KeyId                  string                              `url:"backblaze_b2_key_id,omitempty" required:"false"`
	BackblazeB2ApplicationKey         string                              `url:"backblaze_b2_application_key,omitempty" required:"false"`
	RackspaceApiKey                   string                              `url:"rackspace_api_key,omitempty" required:"false"`
	ResetAuthentication               *bool                               `url:"reset_authentication,omitempty" required:"false"`
	AzureBlobStorageAccessKey         string                              `url:"azure_blob_storage_access_key,omitempty" required:"false"`
	Hostname                          string                              `url:"hostname,omitempty" required:"false"`
	Name                              string                              `url:"name,omitempty" required:"false"`
	MaxConnections                    int64                               `url:"max_connections,omitempty" required:"false"`
	Port                              int64                               `url:"port,omitempty" required:"false"`
	S3Bucket                          string                              `url:"s3_bucket,omitempty" required:"false"`
	S3Region                          string                              `url:"s3_region,omitempty" required:"false"`
	ServerCertificate                 RemoteServerServerCertificateEnum   `url:"server_certificate,omitempty" required:"false"`
	ServerHostKey                     string                              `url:"server_host_key,omitempty" required:"false"`
	ServerType                        RemoteServerServerTypeEnum          `url:"server_type,omitempty" required:"false"`
	Ssl                               RemoteServerSslEnum                 `url:"ssl,omitempty" required:"false"`
	Username                          string                              `url:"username,omitempty" required:"false"`
	GoogleCloudStorageBucket          string                              `url:"google_cloud_storage_bucket,omitempty" required:"false"`
	GoogleCloudStorageProjectId       string                              `url:"google_cloud_storage_project_id,omitempty" required:"false"`
	BackblazeB2Bucket                 string                              `url:"backblaze_b2_bucket,omitempty" required:"false"`
	BackblazeB2S3Endpoint             string                              `url:"backblaze_b2_s3_endpoint,omitempty" required:"false"`
	WasabiBucket                      string                              `url:"wasabi_bucket,omitempty" required:"false"`
	WasabiRegion                      string                              `url:"wasabi_region,omitempty" required:"false"`
	RackspaceUsername                 string                              `url:"rackspace_username,omitempty" required:"false"`
	RackspaceRegion                   string                              `url:"rackspace_region,omitempty" required:"false"`
	RackspaceContainer                string                              `url:"rackspace_container,omitempty" required:"false"`
	OneDriveAccountType               RemoteServerOneDriveAccountTypeEnum `url:"one_drive_account_type,omitempty" required:"false"`
	AzureBlobStorageAccount           string                              `url:"azure_blob_storage_account,omitempty" required:"false"`
	AzureBlobStorageContainer         string                              `url:"azure_blob_storage_container,omitempty" required:"false"`
	S3CompatibleBucket                string                              `url:"s3_compatible_bucket,omitempty" required:"false"`
	S3CompatibleRegion                string                              `url:"s3_compatible_region,omitempty" required:"false"`
	S3CompatibleEndpoint              string                              `url:"s3_compatible_endpoint,omitempty" required:"false"`
	S3CompatibleAccessKey             string                              `url:"s3_compatible_access_key,omitempty" required:"false"`
	S3CompatibleSecretKey             string                              `url:"s3_compatible_secret_key,omitempty" required:"false"`
}

type RemoteServerUpdateParams struct {
	Id                                int64                               `url:"-,omitempty" required:"true"`
	AwsAccessKey                      string                              `url:"aws_access_key,omitempty" required:"false"`
	AwsSecretKey                      string                              `url:"aws_secret_key,omitempty" required:"false"`
	Password                          string                              `url:"password,omitempty" required:"false"`
	PrivateKey                        string                              `url:"private_key,omitempty" required:"false"`
	SslCertificate                    string                              `url:"ssl_certificate,omitempty" required:"false"`
	GoogleCloudStorageCredentialsJson string                              `url:"google_cloud_storage_credentials_json,omitempty" required:"false"`
	WasabiAccessKey                   string                              `url:"wasabi_access_key,omitempty" required:"false"`
	WasabiSecretKey                   string                              `url:"wasabi_secret_key,omitempty" required:"false"`
	BackblazeB2KeyId                  string                              `url:"backblaze_b2_key_id,omitempty" required:"false"`
	BackblazeB2ApplicationKey         string                              `url:"backblaze_b2_application_key,omitempty" required:"false"`
	RackspaceApiKey                   string                              `url:"rackspace_api_key,omitempty" required:"false"`
	ResetAuthentication               *bool                               `url:"reset_authentication,omitempty" required:"false"`
	AzureBlobStorageAccessKey         string                              `url:"azure_blob_storage_access_key,omitempty" required:"false"`
	Hostname                          string                              `url:"hostname,omitempty" required:"false"`
	Name                              string                              `url:"name,omitempty" required:"false"`
	MaxConnections                    int64                               `url:"max_connections,omitempty" required:"false"`
	Port                              int64                               `url:"port,omitempty" required:"false"`
	S3Bucket                          string                              `url:"s3_bucket,omitempty" required:"false"`
	S3Region                          string                              `url:"s3_region,omitempty" required:"false"`
	ServerCertificate                 RemoteServerServerCertificateEnum   `url:"server_certificate,omitempty" required:"false"`
	ServerHostKey                     string                              `url:"server_host_key,omitempty" required:"false"`
	ServerType                        RemoteServerServerTypeEnum          `url:"server_type,omitempty" required:"false"`
	Ssl                               RemoteServerSslEnum                 `url:"ssl,omitempty" required:"false"`
	Username                          string                              `url:"username,omitempty" required:"false"`
	GoogleCloudStorageBucket          string                              `url:"google_cloud_storage_bucket,omitempty" required:"false"`
	GoogleCloudStorageProjectId       string                              `url:"google_cloud_storage_project_id,omitempty" required:"false"`
	BackblazeB2Bucket                 string                              `url:"backblaze_b2_bucket,omitempty" required:"false"`
	BackblazeB2S3Endpoint             string                              `url:"backblaze_b2_s3_endpoint,omitempty" required:"false"`
	WasabiBucket                      string                              `url:"wasabi_bucket,omitempty" required:"false"`
	WasabiRegion                      string                              `url:"wasabi_region,omitempty" required:"false"`
	RackspaceUsername                 string                              `url:"rackspace_username,omitempty" required:"false"`
	RackspaceRegion                   string                              `url:"rackspace_region,omitempty" required:"false"`
	RackspaceContainer                string                              `url:"rackspace_container,omitempty" required:"false"`
	OneDriveAccountType               RemoteServerOneDriveAccountTypeEnum `url:"one_drive_account_type,omitempty" required:"false"`
	AzureBlobStorageAccount           string                              `url:"azure_blob_storage_account,omitempty" required:"false"`
	AzureBlobStorageContainer         string                              `url:"azure_blob_storage_container,omitempty" required:"false"`
	S3CompatibleBucket                string                              `url:"s3_compatible_bucket,omitempty" required:"false"`
	S3CompatibleRegion                string                              `url:"s3_compatible_region,omitempty" required:"false"`
	S3CompatibleEndpoint              string                              `url:"s3_compatible_endpoint,omitempty" required:"false"`
	S3CompatibleAccessKey             string                              `url:"s3_compatible_access_key,omitempty" required:"false"`
	S3CompatibleSecretKey             string                              `url:"s3_compatible_secret_key,omitempty" required:"false"`
}

type RemoteServerDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

func (r *RemoteServer) UnmarshalJSON(data []byte) error {
	type remoteServer RemoteServer
	var v remoteServer
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*r = RemoteServer(v)
	return nil
}

func (r *RemoteServerCollection) UnmarshalJSON(data []byte) error {
	type remoteServers []RemoteServer
	var v remoteServers
	if err := json.Unmarshal(data, &v); err != nil {
		return err
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
