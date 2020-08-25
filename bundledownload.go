package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type BundleDownload struct {
	DownloadMethod string    `json:"download_method,omitempty"`
	Path           string    `json:"path,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
}

type BundleDownloadCollection []BundleDownload

type BundleDownloadListParams struct {
	Page                 int    `url:"page,omitempty" required:"false"`
	PerPage              int    `url:"per_page,omitempty" required:"false"`
	Action               string `url:"action,omitempty" required:"false"`
	Cursor               string `url:"cursor,omitempty" required:"false"`
	BundleRegistrationId int64  `url:"bundle_registration_id,omitempty" required:"true"`
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
