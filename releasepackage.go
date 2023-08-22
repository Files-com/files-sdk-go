package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type ReleasePackage struct {
	PackageLink string `json:"package_link,omitempty" path:"package_link,omitempty" url:"package_link,omitempty"`
	Platform    string `json:"platform,omitempty" path:"platform,omitempty" url:"platform,omitempty"`
}

// Identifier no path or id

type ReleasePackageCollection []ReleasePackage

func (r *ReleasePackage) UnmarshalJSON(data []byte) error {
	type releasePackage ReleasePackage
	var v releasePackage
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*r = ReleasePackage(v)
	return nil
}

func (r *ReleasePackageCollection) UnmarshalJSON(data []byte) error {
	type releasePackages ReleasePackageCollection
	var v releasePackages
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*r = ReleasePackageCollection(v)
	return nil
}

func (r *ReleasePackageCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*r))
	for i, v := range *r {
		ret[i] = v
	}

	return &ret
}
