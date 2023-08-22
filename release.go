package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Release struct {
	Version               string   `json:"version,omitempty" path:"version,omitempty" url:"version,omitempty"`
	Description           string   `json:"description,omitempty" path:"description,omitempty" url:"description,omitempty"`
	NativeReleasePackages []string `json:"native_release_packages,omitempty" path:"native_release_packages,omitempty" url:"native_release_packages,omitempty"`
	Title                 string   `json:"title,omitempty" path:"title,omitempty" url:"title,omitempty"`
}

// Identifier no path or id

type ReleaseCollection []Release

type ReleaseGetLatestParams struct {
	Platform string `url:"platform,omitempty" required:"false" json:"platform,omitempty" path:"platform"`
}

func (r *Release) UnmarshalJSON(data []byte) error {
	type release Release
	var v release
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*r = Release(v)
	return nil
}

func (r *ReleaseCollection) UnmarshalJSON(data []byte) error {
	type releases ReleaseCollection
	var v releases
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*r = ReleaseCollection(v)
	return nil
}

func (r *ReleaseCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*r))
	for i, v := range *r {
		ret[i] = v
	}

	return &ret
}
