package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// InboxUpload is a Files.com API resource.
type InboxUpload struct {
	InboxRegistration InboxRegistration `json:"inbox_registration,omitempty" path:"inbox_registration,omitempty" url:"inbox_registration,omitempty"`
	Path              string            `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	CreatedAt         *time.Time        `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
}

// Identifier returns the resource path.
func (i InboxUpload) Identifier() interface{} {
	return i.Path
}

// InboxUploadCollection is a list of InboxUpload resources.
type InboxUploadCollection []InboxUpload

// InboxUploadListParams contains the request parameters for GET /inbox_uploads.
type InboxUploadListParams struct {
	SortBy     interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	Filter     interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	FilterGt   interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterLt   interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (i *InboxUpload) UnmarshalJSON(data []byte) error {
	type inboxUpload InboxUpload
	var v inboxUpload
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*i = InboxUpload(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (i *InboxUploadCollection) UnmarshalJSON(data []byte) error {
	type inboxUploads InboxUploadCollection
	var v inboxUploads
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*i = InboxUploadCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (i *InboxUploadCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*i))
	for i, v := range *i {
		ret[i] = v
	}

	return &ret
}
