package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Setting struct {
}

// Identifier no path or id

type SettingCollection []Setting

type SettingListParams struct {
	ListParams
}

type SettingGetDomainParams struct {
	Domain string `url:"domain,omitempty" required:"true" json:"domain,omitempty" path:"domain"`
}

func (s *Setting) UnmarshalJSON(data []byte) error {
	type setting Setting
	var v setting
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = Setting(v)
	return nil
}

func (s *SettingCollection) UnmarshalJSON(data []byte) error {
	type settings SettingCollection
	var v settings
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SettingCollection(v)
	return nil
}

func (s *SettingCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
