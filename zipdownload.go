package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type ZipDownload struct {
	DownloadUri            string   `json:"download_uri,omitempty" path:"download_uri,omitempty" url:"download_uri,omitempty"`
	Paths                  []string `json:"paths,omitempty" path:"paths,omitempty" url:"paths,omitempty"`
	BundleRegistrationCode string   `json:"bundle_registration_code,omitempty" path:"bundle_registration_code,omitempty" url:"bundle_registration_code,omitempty"`
	EncodedPaths           []string `json:"encoded_paths,omitempty" path:"encoded_paths,omitempty" url:"encoded_paths,omitempty"`
}

// Identifier no path or id

type ZipDownloadCollection []ZipDownload

type ZipDownloadCreateParams struct {
	Paths                  []string `url:"paths,omitempty" required:"true" json:"paths,omitempty" path:"paths"`
	BundleRegistrationCode string   `url:"bundle_registration_code,omitempty" required:"false" json:"bundle_registration_code,omitempty" path:"bundle_registration_code"`
	EncodedPaths           []string `url:"encoded_paths,omitempty" required:"false" json:"encoded_paths,omitempty" path:"encoded_paths"`
}

func (z *ZipDownload) UnmarshalJSON(data []byte) error {
	type zipDownload ZipDownload
	var v zipDownload
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*z = ZipDownload(v)
	return nil
}

func (z *ZipDownloadCollection) UnmarshalJSON(data []byte) error {
	type zipDownloads ZipDownloadCollection
	var v zipDownloads
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*z = ZipDownloadCollection(v)
	return nil
}

func (z *ZipDownloadCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*z))
	for i, v := range *z {
		ret[i] = v
	}

	return &ret
}
