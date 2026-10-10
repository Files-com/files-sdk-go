package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// ZipListEntry is a Files.com API resource.
type ZipListEntry struct {
	Path string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	Size int64  `json:"size,omitempty" path:"size,omitempty" url:"size,omitempty"`
}

// Identifier returns the resource path.
func (z ZipListEntry) Identifier() interface{} {
	return z.Path
}

// ZipListEntryCollection is a list of ZipListEntry resources.
type ZipListEntryCollection []ZipListEntry

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (z *ZipListEntry) UnmarshalJSON(data []byte) error {
	type zipListEntry ZipListEntry
	var v zipListEntry
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*z = ZipListEntry(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (z *ZipListEntryCollection) UnmarshalJSON(data []byte) error {
	type zipListEntrys ZipListEntryCollection
	var v zipListEntrys
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*z = ZipListEntryCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (z *ZipListEntryCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*z))
	for i, v := range *z {
		ret[i] = v
	}

	return &ret
}
