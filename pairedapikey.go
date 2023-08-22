package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type PairedApiKey struct {
	Server       string `json:"server,omitempty" path:"server,omitempty" url:"server,omitempty"`
	Username     string `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	Password     string `json:"password,omitempty" path:"password,omitempty" url:"password,omitempty"`
	UserUsername string `json:"user_username,omitempty" path:"user_username,omitempty" url:"user_username,omitempty"`
	Nickname     string `json:"nickname,omitempty" path:"nickname,omitempty" url:"nickname,omitempty"`
}

// Identifier no path or id

type PairedApiKeyCollection []PairedApiKey

func (p *PairedApiKey) UnmarshalJSON(data []byte) error {
	type pairedApiKey PairedApiKey
	var v pairedApiKey
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*p = PairedApiKey(v)
	return nil
}

func (p *PairedApiKeyCollection) UnmarshalJSON(data []byte) error {
	type pairedApiKeys PairedApiKeyCollection
	var v pairedApiKeys
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*p = PairedApiKeyCollection(v)
	return nil
}

func (p *PairedApiKeyCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
