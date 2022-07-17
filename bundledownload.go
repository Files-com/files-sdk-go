package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type BundleDownload struct {
	BundleRegistration BundleRegistration `json:"bundle_registration,omitempty" path:"bundle_registration"`
	DownloadMethod     string             `json:"download_method,omitempty" path:"download_method"`
	Path               string             `json:"path,omitempty" path:"path"`
	CreatedAt          *time.Time         `json:"created_at,omitempty" path:"created_at"`
}

type BundleDownloadCollection []BundleDownload

type BundleDownloadListParams struct {
	SortBy               json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter               json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterGt             json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq           json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterLike           json.RawMessage `url:"filter_like,omitempty" required:"false" json:"filter_like,omitempty" path:"filter_like"`
	FilterLt             json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq           json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty" path:"filter_lteq"`
	BundleId             int64           `url:"bundle_id,omitempty" required:"false" json:"bundle_id,omitempty" path:"bundle_id"`
	BundleRegistrationId int64           `url:"bundle_registration_id,omitempty" required:"false" json:"bundle_registration_id,omitempty" path:"bundle_registration_id"`
	lib.ListParams
}

func (b *BundleDownload) UnmarshalJSON(data []byte) error {
	type bundleDownload BundleDownload
	var v bundleDownload
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*b = BundleDownload(v)
	return nil
}

func (b *BundleDownloadCollection) UnmarshalJSON(data []byte) error {
	type bundleDownloads BundleDownloadCollection
	var v bundleDownloads
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*b = BundleDownloadCollection(v)
	return nil
}

func (b *BundleDownloadCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*b))
	for i, v := range *b {
		ret[i] = v
	}

	return &ret
}
