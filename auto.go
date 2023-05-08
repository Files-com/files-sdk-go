package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Auto struct {
	Dynamic json.RawMessage `json:"dynamic,omitempty" path:"dynamic"`
}

// Identifier no path or id

type AutoCollection []Auto

func (a *Auto) UnmarshalJSON(data []byte) error {
	type auto Auto
	var v auto
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = Auto(v)
	return nil
}

func (a *AutoCollection) UnmarshalJSON(data []byte) error {
	type autos AutoCollection
	var v autos
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = AutoCollection(v)
	return nil
}

func (a *AutoCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
