package files_sdk

import (
	"encoding/json"
)

type Image struct {
	Name string `json:"name,omitempty"`
	Uri  string `json:"uri,omitempty"`
}

type ImageCollection []Image

func (i *Image) UnmarshalJSON(data []byte) error {
	type image Image
	var v image
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = Image(v)
	return nil
}

func (i *ImageCollection) UnmarshalJSON(data []byte) error {
	type images []Image
	var v images
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = ImageCollection(v)
	return nil
}

func (i *ImageCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*i))
	for i, v := range *i {
		ret[i] = v
	}

	return &ret
}
