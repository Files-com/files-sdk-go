package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type RemoteServerConfigurationFile struct {
	Id                      int64  `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	PermissionSet           string `json:"permission_set,omitempty" path:"permission_set,omitempty" url:"permission_set,omitempty"`
	PrivateKey              string `json:"private_key,omitempty" path:"private_key,omitempty" url:"private_key,omitempty"`
	Subdomain               string `json:"subdomain,omitempty" path:"subdomain,omitempty" url:"subdomain,omitempty"`
	Root                    string `json:"root,omitempty" path:"root,omitempty" url:"root,omitempty"`
	FollowLinks             *bool  `json:"follow_links,omitempty" path:"follow_links,omitempty" url:"follow_links,omitempty"`
	PreferProtocol          string `json:"prefer_protocol,omitempty" path:"prefer_protocol,omitempty" url:"prefer_protocol,omitempty"`
	Dns                     string `json:"dns,omitempty" path:"dns,omitempty" url:"dns,omitempty"`
	ProxyAllOutbound        *bool  `json:"proxy_all_outbound,omitempty" path:"proxy_all_outbound,omitempty" url:"proxy_all_outbound,omitempty"`
	EndpointOverride        string `json:"endpoint_override,omitempty" path:"endpoint_override,omitempty" url:"endpoint_override,omitempty"`
	LogFile                 string `json:"log_file,omitempty" path:"log_file,omitempty" url:"log_file,omitempty"`
	LogLevel                string `json:"log_level,omitempty" path:"log_level,omitempty" url:"log_level,omitempty"`
	LogRotateNum            int64  `json:"log_rotate_num,omitempty" path:"log_rotate_num,omitempty" url:"log_rotate_num,omitempty"`
	LogRotateSize           int64  `json:"log_rotate_size,omitempty" path:"log_rotate_size,omitempty" url:"log_rotate_size,omitempty"`
	MaxConcurrentJobs       int64  `json:"max_concurrent_jobs,omitempty" path:"max_concurrent_jobs,omitempty" url:"max_concurrent_jobs,omitempty"`
	GracefulShutdownTimeout int64  `json:"graceful_shutdown_timeout,omitempty" path:"graceful_shutdown_timeout,omitempty" url:"graceful_shutdown_timeout,omitempty"`
	TransferRateLimit       string `json:"transfer_rate_limit,omitempty" path:"transfer_rate_limit,omitempty" url:"transfer_rate_limit,omitempty"`
}

func (r RemoteServerConfigurationFile) Identifier() interface{} {
	return r.Id
}

type RemoteServerConfigurationFileCollection []RemoteServerConfigurationFile

func (r *RemoteServerConfigurationFile) UnmarshalJSON(data []byte) error {
	type remoteServerConfigurationFile RemoteServerConfigurationFile
	var v remoteServerConfigurationFile
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*r = RemoteServerConfigurationFile(v)
	return nil
}

func (r *RemoteServerConfigurationFileCollection) UnmarshalJSON(data []byte) error {
	type remoteServerConfigurationFiles RemoteServerConfigurationFileCollection
	var v remoteServerConfigurationFiles
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*r = RemoteServerConfigurationFileCollection(v)
	return nil
}

func (r *RemoteServerConfigurationFileCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*r))
	for i, v := range *r {
		ret[i] = v
	}

	return &ret
}
