package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type RemoteServerConfigurationFile struct {
	Id            int64  `json:"id,omitempty" path:"id"`
	PermissionSet string `json:"permission_set,omitempty" path:"permission_set"`
	ApiToken      string `json:"api_token,omitempty" path:"api_token"`
	Root          string `json:"root,omitempty" path:"root"`
	Port          int64  `json:"port,omitempty" path:"port"`
	Hostname      string `json:"hostname,omitempty" path:"hostname"`
	PublicKey     string `json:"public_key,omitempty" path:"public_key"`
	PrivateKey    string `json:"private_key,omitempty" path:"private_key"`
	Status        string `json:"status,omitempty" path:"status"`
	ConfigVersion string `json:"config_version,omitempty" path:"config_version"`
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
