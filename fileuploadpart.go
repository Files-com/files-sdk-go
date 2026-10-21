package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// FileUploadPart is a Files.com API resource.
type FileUploadPart struct {
	// Content-Type and File to send.
	Send interface{} `json:"send,omitempty" path:"send,omitempty" url:"send,omitempty"`
	// Type of upload.
	Action string `json:"action,omitempty" path:"action,omitempty" url:"action,omitempty"`
	// If `true`, this file exists and you may wish to ask the user for overwrite confirmation.
	AskAboutOverwrites *bool `json:"ask_about_overwrites,omitempty" path:"ask_about_overwrites,omitempty" url:"ask_about_overwrites,omitempty"`
	// Number of parts in the upload.
	AvailableParts int64 `json:"available_parts,omitempty" path:"available_parts,omitempty" url:"available_parts,omitempty"`
	// Date/time of when this Upload part expires and the URL cannot be used any more.
	Expires string `json:"expires,omitempty" path:"expires,omitempty" url:"expires,omitempty"`
	// Additional upload headers to provide as part of the upload.
	Headers interface{} `json:"headers,omitempty" path:"headers,omitempty" url:"headers,omitempty"`
	// HTTP Method to use for uploading the part, usually `PUT`.
	HttpMethod string `json:"http_method,omitempty" path:"http_method,omitempty" url:"http_method,omitempty"`
	// Size in bytes for this part.
	NextPartsize int64 `json:"next_partsize,omitempty" path:"next_partsize,omitempty" url:"next_partsize,omitempty"`
	// If `true`, multiple parts may be uploaded in parallel. If `false`, be sure to only upload one part at a time, in order.
	ParallelParts *bool `json:"parallel_parts,omitempty" path:"parallel_parts,omitempty" url:"parallel_parts,omitempty"`
	// If `true`, parts may be retried. If `false`, a part cannot be retried and the upload should be restarted.
	RetryParts *bool `json:"retry_parts,omitempty" path:"retry_parts,omitempty" url:"retry_parts,omitempty"`
	// Additional HTTP parameters to send with the upload.
	Parameters interface{} `json:"parameters,omitempty" path:"parameters,omitempty" url:"parameters,omitempty"`
	// Number of this upload part.
	PartNumber int64 `json:"part_number,omitempty" path:"part_number,omitempty" url:"part_number,omitempty"`
	// Size in bytes for the next upload part.
	Partsize int64 `json:"partsize,omitempty" path:"partsize,omitempty" url:"partsize,omitempty"`
	// New file path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.
	Path string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	// Reference name for this upload part.
	Ref string `json:"ref,omitempty" path:"ref,omitempty" url:"ref,omitempty"`
	// URI to upload this part to.
	UploadUri string `json:"upload_uri,omitempty" path:"upload_uri,omitempty" url:"upload_uri,omitempty"`
	// Optional direct connection information for direct Agent transfer attempts.
	DirectConnectionInfo DirectConnectionInfo `json:"direct_connection_info,omitempty" path:"direct_connection_info,omitempty" url:"direct_connection_info,omitempty"`
}

// Identifier returns the resource path.
func (f FileUploadPart) Identifier() interface{} {
	return f.Path
}

// FileUploadPartCollection is a list of FileUploadPart resources.
type FileUploadPartCollection []FileUploadPart

const UploadPartExpires = time.Minute * 15
const UploadObjectExpires = time.Minute * (24 * 3)

func (f FileUploadPart) ExpiresTime() time.Time {
	partExpires, _ := time.Parse(time.RFC3339, f.Expires)
	return partExpires
}

// UploadExpires only valid on first part request
func (f FileUploadPart) UploadExpires() time.Time {
	if f.ExpiresTime().IsZero() {
		return time.Time{}
	}
	return f.ExpiresTime().Add(-(UploadPartExpires)).Add(UploadObjectExpires)
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (f *FileUploadPart) UnmarshalJSON(data []byte) error {
	type fileUploadPart FileUploadPart
	var v fileUploadPart
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*f = FileUploadPart(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (f *FileUploadPartCollection) UnmarshalJSON(data []byte) error {
	type fileUploadParts FileUploadPartCollection
	var v fileUploadParts
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*f = FileUploadPartCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (f *FileUploadPartCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*f))
	for i, v := range *f {
		ret[i] = v
	}

	return &ret
}
