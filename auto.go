package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Auto is a Files.com API resource.
type Auto struct {
	Dynamic interface{} `json:"dynamic,omitempty" path:"dynamic,omitempty" url:"dynamic,omitempty"`
}

// Identifier no path or id

// AutoCollection is a list of Auto resources.
type AutoCollection []Auto

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (a *Auto) UnmarshalJSON(data []byte) error {
	type auto Auto
	var v auto
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = Auto(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (a *AutoCollection) UnmarshalJSON(data []byte) error {
	type autos AutoCollection
	var v autos
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = AutoCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (a *AutoCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
