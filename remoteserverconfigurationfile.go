package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// RemoteServerConfigurationFile is a Files.com API resource.
type RemoteServerConfigurationFile struct {
	// The remote server ID of the agent.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// The permission set for the agent ['read_write', 'read_only', 'write_only'].
	PermissionSet string `json:"permission_set,omitempty" path:"permission_set,omitempty" url:"permission_set,omitempty"`
	// The private key for the agent.
	PrivateKey string `json:"private_key,omitempty" path:"private_key,omitempty" url:"private_key,omitempty"`
	// Files.com subdomain site name.
	Subdomain string `json:"subdomain,omitempty" path:"subdomain,omitempty" url:"subdomain,omitempty"`
	// The root directory for the agent.
	Root string `json:"root,omitempty" path:"root,omitempty" url:"root,omitempty"`
	// Optional customer-facing name for this Agent installation.
	NodeName string `json:"node_name,omitempty" path:"node_name,omitempty" url:"node_name,omitempty"`
	// Follow symlinks when traversing directories.
	FollowLinks *bool `json:"follow_links,omitempty" path:"follow_links,omitempty" url:"follow_links,omitempty"`
	// Preferred network protocol ['udp', 'tcp'] (default udp).
	PreferProtocol string `json:"prefer_protocol,omitempty" path:"prefer_protocol,omitempty" url:"prefer_protocol,omitempty"`
	// DNS lookup method ['auto','doh','system'] (default auto).
	Dns string `json:"dns,omitempty" path:"dns,omitempty" url:"dns,omitempty"`
	// Proxy all outbound traffic through files.com proxy server.
	ProxyAllOutbound *bool `json:"proxy_all_outbound,omitempty" path:"proxy_all_outbound,omitempty" url:"proxy_all_outbound,omitempty"`
	// Custom site endpoint URL.
	EndpointOverride string `json:"endpoint_override,omitempty" path:"endpoint_override,omitempty" url:"endpoint_override,omitempty"`
	// Log file name and location.
	LogFile string `json:"log_file,omitempty" path:"log_file,omitempty" url:"log_file,omitempty"`
	// Log level for the agent logs ['debug', 'info', 'warn', 'error', 'fatal'] (default info).
	LogLevel string `json:"log_level,omitempty" path:"log_level,omitempty" url:"log_level,omitempty"`
	// Log route for agent logs. (default 5).
	LogRotateNum int64 `json:"log_rotate_num,omitempty" path:"log_rotate_num,omitempty" url:"log_rotate_num,omitempty"`
	// Log route size in MB for agent logs. (default 20).
	LogRotateSize int64 `json:"log_rotate_size,omitempty" path:"log_rotate_size,omitempty" url:"log_rotate_size,omitempty"`
	// Maximum number of concurrent jobs (default 500).
	OverrideMaxConcurrentJobs int64 `json:"override_max_concurrent_jobs,omitempty" path:"override_max_concurrent_jobs,omitempty" url:"override_max_concurrent_jobs,omitempty"`
	// Graceful shutdown timeout in seconds (default 15).
	GracefulShutdownTimeout int64 `json:"graceful_shutdown_timeout,omitempty" path:"graceful_shutdown_timeout,omitempty" url:"graceful_shutdown_timeout,omitempty"`
	// File transfer (upload/download) rate limit `<limit>-<period>`, with the given periods: * 'S': second * 'M': minute * 'H': hour * 'D': day Examples: * 5 requests/second: '5-S' * 10 requests/minute: '10-M' * 1000 requests/hour: '1000-H' * 2000 requests/day: '2000-D'.
	TransferRateLimit string `json:"transfer_rate_limit,omitempty" path:"transfer_rate_limit,omitempty" url:"transfer_rate_limit,omitempty"`
	// Auto update policy ['manual_trigger', 'critical_only', 'always', 'never'] (default always).
	AutoUpdatePolicy string `json:"auto_update_policy,omitempty" path:"auto_update_policy,omitempty" url:"auto_update_policy,omitempty"`
}

// Identifier returns the resource ID.
func (r RemoteServerConfigurationFile) Identifier() interface{} {
	return r.Id
}

// RemoteServerConfigurationFileCollection is a list of RemoteServerConfigurationFile resources.
type RemoteServerConfigurationFileCollection []RemoteServerConfigurationFile

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (r *RemoteServerConfigurationFile) UnmarshalJSON(data []byte) error {
	type remoteServerConfigurationFile RemoteServerConfigurationFile
	var v remoteServerConfigurationFile
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*r = RemoteServerConfigurationFile(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (r *RemoteServerConfigurationFileCollection) UnmarshalJSON(data []byte) error {
	type remoteServerConfigurationFiles RemoteServerConfigurationFileCollection
	var v remoteServerConfigurationFiles
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*r = RemoteServerConfigurationFileCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (r *RemoteServerConfigurationFileCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*r))
	for i, v := range *r {
		ret[i] = v
	}

	return &ret
}
