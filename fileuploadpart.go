package files_sdk

import (
	"encoding/json"
)

type FileUploadPart struct {
	Send               json.RawMessage `json:"send,omitempty"`
	Action             string          `json:"action,omitempty"`
	AskAboutOverwrites *bool           `json:"ask_about_overwrites,omitempty"`
	AvailableParts     int64           `json:"available_parts,omitempty"`
	Expires            string          `json:"expires,omitempty"`
	Headers            json.RawMessage `json:"headers,omitempty"`
	HttpMethod         string          `json:"http_method,omitempty"`
	NextPartsize       int64           `json:"next_partsize,omitempty"`
	ParallelParts      *bool           `json:"parallel_parts,omitempty"`
	Parameters         json.RawMessage `json:"parameters,omitempty"`
	PartNumber         int64           `json:"part_number,omitempty"`
	Partsize           int64           `json:"partsize,omitempty"`
	Path               string          `json:"path,omitempty"`
	Ref                string          `json:"ref,omitempty"`
	UploadUri          string          `json:"upload_uri,omitempty"`
}

type FileUploadPartCollection []FileUploadPart

func (f *FileUploadPart) UnmarshalJSON(data []byte) error {
	type fileUploadPart FileUploadPart
	var v fileUploadPart
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*f = FileUploadPart(v)
	return nil
}

func (f *FileUploadPartCollection) UnmarshalJSON(data []byte) error {
	type fileUploadParts []FileUploadPart
	var v fileUploadParts
	if err := json.Unmarshal(data, &v); err != nil {
		return err
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
