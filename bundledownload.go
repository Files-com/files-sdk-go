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
	CreatedAt          time.Time          `json:"created_at,omitempty"`
}

type BundleDownloadCollection []BundleDownload

type BundleDownloadListParams struct {
	Cursor               string          `url:"cursor,omitempty" required:"false"`
	PerPage              int64           `url:"per_page,omitempty" required:"false"`
	SortBy               json.RawMessage `url:"sort_by,omitempty" required:"false"`
	Filter               json.RawMessage `url:"filter,omitempty" required:"false"`
	FilterGt             json.RawMessage `url:"filter_gt,omitempty" required:"false"`
	FilterGteq           json.RawMessage `url:"filter_gteq,omitempty" required:"false"`
	FilterLike           json.RawMessage `url:"filter_like,omitempty" required:"false"`
	FilterLt             json.RawMessage `url:"filter_lt,omitempty" required:"false"`
	FilterLteq           json.RawMessage `url:"filter_lteq,omitempty" required:"false"`
	BundleId             int64           `url:"bundle_id,omitempty" required:"false"`
	BundleRegistrationId int64           `url:"bundle_registration_id,omitempty" required:"false"`
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
