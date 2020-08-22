package files_sdk

import (
	"encoding/json"
)

type FilePartUpload struct {
	Send               json.RawMessage `json:"send,omitempty"`
	Action             string          `json:"action,omitempty"`
	AskAboutOverwrites *bool           `json:"ask_about_overwrites,omitempty"`
	AvailableParts     int             `json:"available_parts,omitempty"`
	Expires            string          `json:"expires,omitempty"`
	Headers            json.RawMessage `json:"headers,omitempty"`
	HttpMethod         string          `json:"http_method,omitempty"`
	NextPartsize       int             `json:"next_partsize,omitempty"`
	ParallelParts      *bool           `json:"parallel_parts,omitempty"`
	Parameters         json.RawMessage `json:"parameters,omitempty"`
	PartNumber         int             `json:"part_number,omitempty"`
	Partsize           int             `json:"partsize,omitempty"`
	Path               string          `json:"path,omitempty"`
	Ref                string          `json:"ref,omitempty"`
	UploadUri          string          `json:"upload_uri,omitempty"`
}

type FilePartUploadCollection []FilePartUpload

func (f *FilePartUpload) UnmarshalJSON(data []byte) error {
	type filePartUpload FilePartUpload
	var v filePartUpload
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*f = FilePartUpload(v)
	return nil
}

func (f *FilePartUploadCollection) UnmarshalJSON(data []byte) error {
	type filePartUploads []FilePartUpload
	var v filePartUploads
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*f = FilePartUploadCollection(v)
	return nil
}
