package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// BundleDownload is a Files.com API resource.
type BundleDownload struct {
	BundleRegistration BundleRegistration `json:"bundle_registration,omitempty" path:"bundle_registration,omitempty" url:"bundle_registration,omitempty"`
	DownloadMethod     string             `json:"download_method,omitempty" path:"download_method,omitempty" url:"download_method,omitempty"`
	Path               string             `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	WorkspaceId        int64              `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	CreatedAt          *time.Time         `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
}

// Identifier returns the resource path.
func (b BundleDownload) Identifier() interface{} {
	return b.Path
}

// BundleDownloadCollection is a list of BundleDownload resources.
type BundleDownloadCollection []BundleDownload

// BundleDownloadListParams contains the request parameters for GET /bundle_downloads.
type BundleDownloadListParams struct {
	SortBy               interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	Filter               interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	FilterGt             interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq           interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterLt             interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq           interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	BundleId             int64       `url:"bundle_id,omitempty" json:"bundle_id,omitempty" path:"bundle_id"`
	BundleRegistrationId int64       `url:"bundle_registration_id,omitempty" json:"bundle_registration_id,omitempty" path:"bundle_registration_id"`
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (b *BundleDownload) UnmarshalJSON(data []byte) error {
	type bundleDownload BundleDownload
	var v bundleDownload
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*b = BundleDownload(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (b *BundleDownloadCollection) UnmarshalJSON(data []byte) error {
	type bundleDownloads BundleDownloadCollection
	var v bundleDownloads
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*b = BundleDownloadCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (b *BundleDownloadCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*b))
	for i, v := range *b {
		ret[i] = v
	}

	return &ret
}
