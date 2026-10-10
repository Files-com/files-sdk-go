package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Image is a Files.com API resource.
type Image struct {
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Uri  string `json:"uri,omitempty" path:"uri,omitempty" url:"uri,omitempty"`
}

// Identifier no path or id

// ImageCollection is a list of Image resources.
type ImageCollection []Image

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (i *Image) UnmarshalJSON(data []byte) error {
	type image Image
	var v image
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*i = Image(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (i *ImageCollection) UnmarshalJSON(data []byte) error {
	type images ImageCollection
	var v images
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*i = ImageCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (i *ImageCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*i))
	for i, v := range *i {
		ret[i] = v
	}

	return &ret
}
