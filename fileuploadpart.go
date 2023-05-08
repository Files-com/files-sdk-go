package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type FileUploadPart struct {
	Send               json.RawMessage `json:"send,omitempty" path:"send"`
	Action             string          `json:"action,omitempty" path:"action"`
	AskAboutOverwrites *bool           `json:"ask_about_overwrites,omitempty" path:"ask_about_overwrites"`
	AvailableParts     int64           `json:"available_parts,omitempty" path:"available_parts"`
	Expires            string          `json:"expires,omitempty" path:"expires"`
	Headers            json.RawMessage `json:"headers,omitempty" path:"headers"`
	HttpMethod         string          `json:"http_method,omitempty" path:"http_method"`
	NextPartsize       int64           `json:"next_partsize,omitempty" path:"next_partsize"`
	ParallelParts      *bool           `json:"parallel_parts,omitempty" path:"parallel_parts"`
	RetryParts         *bool           `json:"retry_parts,omitempty" path:"retry_parts"`
	Parameters         json.RawMessage `json:"parameters,omitempty" path:"parameters"`
	PartNumber         int64           `json:"part_number,omitempty" path:"part_number"`
	Partsize           int64           `json:"partsize,omitempty" path:"partsize"`
	Path               string          `json:"path,omitempty" path:"path"`
	Ref                string          `json:"ref,omitempty" path:"ref"`
	UploadUri          string          `json:"upload_uri,omitempty" path:"upload_uri"`
}

func (f FileUploadPart) Identifier() interface{} {
	return f.Path
}

type FileUploadPartCollection []FileUploadPart

func (f *FileUploadPart) UnmarshalJSON(data []byte) error {
	type fileUploadPart FileUploadPart
	var v fileUploadPart
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*f = FileUploadPart(v)
	return nil
}

func (f *FileUploadPartCollection) UnmarshalJSON(data []byte) error {
	type fileUploadParts FileUploadPartCollection
	var v fileUploadParts
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*f = FileUploadPartCollection(v)
	return nil
}

func (f *FileUploadPartCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*f))
	for i, v := range *f {
		ret[i] = v
	}

	return &ret
}
