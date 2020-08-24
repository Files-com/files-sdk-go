package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type RemoteServer struct {
	Id                                int64  `json:"id,omitempty"`
	AuthenticationMethod              string `json:"authentication_method,omitempty"`
	Hostname                          string `json:"hostname,omitempty"`
	Name                              string `json:"name,omitempty"`
	Port                              int    `json:"port,omitempty"`
	MaxConnections                    int    `json:"max_connections,omitempty"`
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
	AwsAccessKey                      string `json:"aws_access_key,omitempty"`
	AwsSecretKey                      string `json:"aws_secret_key,omitempty"`
	Password                          string `json:"password,omitempty"`
	PrivateKey                        string `json:"private_key,omitempty"`
	GoogleCloudStorageCredentialsJson string `json:"google_cloud_storage_credentials_json,omitempty"`
	WasabiAccessKey                   string `json:"wasabi_access_key,omitempty"`
	WasabiSecretKey                   string `json:"wasabi_secret_key,omitempty"`
	BackblazeB2KeyId                  string `json:"backblaze_b2_key_id,omitempty"`
	BackblazeB2ApplicationKey         string `json:"backblaze_b2_application_key,omitempty"`
	RackspaceApiKey                   string `json:"rackspace_api_key,omitempty"`
	ResetAuthentication               *bool  `json:"reset_authentication,omitempty"`
	AzureBlobStorageAccessKey         string `json:"azure_blob_storage_access_key,omitempty"`
}

type RemoteServerCollection []RemoteServer

type RemoteServerListParams struct {
	Page    int    `url:"page,omitempty"`
	PerPage int    `url:"per_page,omitempty"`
	Action  string `url:"action,omitempty"`
	Cursor  string `url:"cursor,omitempty"`
	lib.ListParams
}

type RemoteServerFindParams struct {
	Id int64 `url:"-,omitempty"`
}

type RemoteServerCreateParams struct {
	AwsAccessKey                      string `url:"aws_access_key,omitempty"`
	AwsSecretKey                      string `url:"aws_secret_key,omitempty"`
	Password                          string `url:"password,omitempty"`
	PrivateKey                        string `url:"private_key,omitempty"`
	GoogleCloudStorageCredentialsJson string `url:"google_cloud_storage_credentials_json,omitempty"`
	WasabiAccessKey                   string `url:"wasabi_access_key,omitempty"`
	WasabiSecretKey                   string `url:"wasabi_secret_key,omitempty"`
	BackblazeB2KeyId                  string `url:"backblaze_b2_key_id,omitempty"`
	BackblazeB2ApplicationKey         string `url:"backblaze_b2_application_key,omitempty"`
	RackspaceApiKey                   string `url:"rackspace_api_key,omitempty"`
	ResetAuthentication               *bool  `url:"reset_authentication,omitempty"`
	AzureBlobStorageAccessKey         string `url:"azure_blob_storage_access_key,omitempty"`
	Hostname                          string `url:"hostname,omitempty"`
	Name                              string `url:"name,omitempty"`
	MaxConnections                    int    `url:"max_connections,omitempty"`
	Port                              int    `url:"port,omitempty"`
	S3Bucket                          string `url:"s3_bucket,omitempty"`
	S3Region                          string `url:"s3_region,omitempty"`
	ServerCertificate                 string `url:"server_certificate,omitempty"`
	ServerHostKey                     string `url:"server_host_key,omitempty"`
	ServerType                        string `url:"server_type,omitempty"`
	Ssl                               string `url:"ssl,omitempty"`
	Username                          string `url:"username,omitempty"`
	GoogleCloudStorageBucket          string `url:"google_cloud_storage_bucket,omitempty"`
	GoogleCloudStorageProjectId       string `url:"google_cloud_storage_project_id,omitempty"`
	BackblazeB2Bucket                 string `url:"backblaze_b2_bucket,omitempty"`
	BackblazeB2S3Endpoint             string `url:"backblaze_b2_s3_endpoint,omitempty"`
	WasabiBucket                      string `url:"wasabi_bucket,omitempty"`
	WasabiRegion                      string `url:"wasabi_region,omitempty"`
	RackspaceUsername                 string `url:"rackspace_username,omitempty"`
	RackspaceRegion                   string `url:"rackspace_region,omitempty"`
	RackspaceContainer                string `url:"rackspace_container,omitempty"`
	OneDriveAccountType               string `url:"one_drive_account_type,omitempty"`
	AzureBlobStorageAccount           string `url:"azure_blob_storage_account,omitempty"`
	AzureBlobStorageContainer         string `url:"azure_blob_storage_container,omitempty"`
}

type RemoteServerUpdateParams struct {
	Id                                int64  `url:"-,omitempty"`
	AwsAccessKey                      string `url:"aws_access_key,omitempty"`
	AwsSecretKey                      string `url:"aws_secret_key,omitempty"`
	Password                          string `url:"password,omitempty"`
	PrivateKey                        string `url:"private_key,omitempty"`
	GoogleCloudStorageCredentialsJson string `url:"google_cloud_storage_credentials_json,omitempty"`
	WasabiAccessKey                   string `url:"wasabi_access_key,omitempty"`
	WasabiSecretKey                   string `url:"wasabi_secret_key,omitempty"`
	BackblazeB2KeyId                  string `url:"backblaze_b2_key_id,omitempty"`
	BackblazeB2ApplicationKey         string `url:"backblaze_b2_application_key,omitempty"`
	RackspaceApiKey                   string `url:"rackspace_api_key,omitempty"`
	ResetAuthentication               *bool  `url:"reset_authentication,omitempty"`
	AzureBlobStorageAccessKey         string `url:"azure_blob_storage_access_key,omitempty"`
	Hostname                          string `url:"hostname,omitempty"`
	Name                              string `url:"name,omitempty"`
	MaxConnections                    int    `url:"max_connections,omitempty"`
	Port                              int    `url:"port,omitempty"`
	S3Bucket                          string `url:"s3_bucket,omitempty"`
	S3Region                          string `url:"s3_region,omitempty"`
	ServerCertificate                 string `url:"server_certificate,omitempty"`
	ServerHostKey                     string `url:"server_host_key,omitempty"`
	ServerType                        string `url:"server_type,omitempty"`
	Ssl                               string `url:"ssl,omitempty"`
	Username                          string `url:"username,omitempty"`
	GoogleCloudStorageBucket          string `url:"google_cloud_storage_bucket,omitempty"`
	GoogleCloudStorageProjectId       string `url:"google_cloud_storage_project_id,omitempty"`
	BackblazeB2Bucket                 string `url:"backblaze_b2_bucket,omitempty"`
	BackblazeB2S3Endpoint             string `url:"backblaze_b2_s3_endpoint,omitempty"`
	WasabiBucket                      string `url:"wasabi_bucket,omitempty"`
	WasabiRegion                      string `url:"wasabi_region,omitempty"`
	RackspaceUsername                 string `url:"rackspace_username,omitempty"`
	RackspaceRegion                   string `url:"rackspace_region,omitempty"`
	RackspaceContainer                string `url:"rackspace_container,omitempty"`
	OneDriveAccountType               string `url:"one_drive_account_type,omitempty"`
	AzureBlobStorageAccount           string `url:"azure_blob_storage_account,omitempty"`
	AzureBlobStorageContainer         string `url:"azure_blob_storage_container,omitempty"`
}

type RemoteServerDeleteParams struct {
	Id int64 `url:"-,omitempty"`
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
