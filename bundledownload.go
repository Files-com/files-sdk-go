package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type BundleDownload struct {
	BundleRegistration BundleRegistration `json:"bundle_registration,omitempty"`
	DownloadMethod     string             `json:"download_method,omitempty"`
	Path               string             `json:"path,omitempty"`
	CreatedAt          *time.Time         `json:"created_at,omitempty"`
}

type BundleDownloadCollection []BundleDownload

type BundleDownloadListParams struct {
	Cursor               string          `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage              int64           `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	SortBy               json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty"`
	Filter               json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty"`
	FilterGt             json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty"`
	FilterGteq           json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty"`
	FilterLike           json.RawMessage `url:"filter_like,omitempty" required:"false" json:"filter_like,omitempty"`
	FilterLt             json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty"`
	FilterLteq           json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty"`
	BundleId             int64           `url:"bundle_id,omitempty" required:"false" json:"bundle_id,omitempty"`
	BundleRegistrationId int64           `url:"bundle_registration_id,omitempty" required:"false" json:"bundle_registration_id,omitempty"`
	lib.ListParams
}

func (b *BundleDownload) UnmarshalJSON(data []byte) error {
	type bundleDownload BundleDownload
	var v bundleDownload
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*b = BundleDownload(v)
	return nil
}

func (b *BundleDownloadCollection) UnmarshalJSON(data []byte) error {
	type bundleDownloads []BundleDownload
	var v bundleDownloads
	if err := json.Unmarshal(data, &v); err != nil {
		return err
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
