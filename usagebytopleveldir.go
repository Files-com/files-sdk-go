package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// UsageByTopLevelDir is a Files.com API resource.
type UsageByTopLevelDir struct {
	Dir   string `json:"dir,omitempty" path:"dir,omitempty" url:"dir,omitempty"`
	Size  int64  `json:"size,omitempty" path:"size,omitempty" url:"size,omitempty"`
	Count int64  `json:"count,omitempty" path:"count,omitempty" url:"count,omitempty"`
}

// Identifier no path or id

// UsageByTopLevelDirCollection is a list of UsageByTopLevelDir resources.
type UsageByTopLevelDirCollection []UsageByTopLevelDir

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (u *UsageByTopLevelDir) UnmarshalJSON(data []byte) error {
	type usageByTopLevelDir UsageByTopLevelDir
	var v usageByTopLevelDir
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*u = UsageByTopLevelDir(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (u *UsageByTopLevelDirCollection) UnmarshalJSON(data []byte) error {
	type usageByTopLevelDirs UsageByTopLevelDirCollection
	var v usageByTopLevelDirs
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*u = UsageByTopLevelDirCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (u *UsageByTopLevelDirCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*u))
	for i, v := range *u {
		ret[i] = v
	}

	return &ret
}
