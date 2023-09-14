package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type RemoteServerConfigurationFile struct {
	Id            int64  `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	PermissionSet string `json:"permission_set,omitempty" path:"permission_set,omitempty" url:"permission_set,omitempty"`
	ApiToken      string `json:"api_token,omitempty" path:"api_token,omitempty" url:"api_token,omitempty"`
	Root          string `json:"root,omitempty" path:"root,omitempty" url:"root,omitempty"`
	Port          int64  `json:"port,omitempty" path:"port,omitempty" url:"port,omitempty"`
	Hostname      string `json:"hostname,omitempty" path:"hostname,omitempty" url:"hostname,omitempty"`
	PublicKey     string `json:"public_key,omitempty" path:"public_key,omitempty" url:"public_key,omitempty"`
	PrivateKey    string `json:"private_key,omitempty" path:"private_key,omitempty" url:"private_key,omitempty"`
	Status        string `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	ConfigVersion string `json:"config_version,omitempty" path:"config_version,omitempty" url:"config_version,omitempty"`
	ServerHostKey string `json:"server_host_key,omitempty" path:"server_host_key,omitempty" url:"server_host_key,omitempty"`
	Subdomain     string `json:"subdomain,omitempty" path:"subdomain,omitempty" url:"subdomain,omitempty"`
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
