package files_sdk

import (
	"encoding/json"
)

type Preview struct {
	Id          int64  `json:"id,omitempty"`
	Status      string `json:"status,omitempty"`
	DownloadUri string `json:"download_uri,omitempty"`
	Type        string `json:"type,omitempty"`
	Size        int    `json:"size,omitempty"`
}

type PreviewCollection []Preview

func (p *Preview) UnmarshalJSON(data []byte) error {
	type preview Preview
	var v preview
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*p = Preview(v)
	return nil
}

func (p *PreviewCollection) UnmarshalJSON(data []byte) error {
	type previews []Preview
	var v previews
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*p = PreviewCollection(v)
	return nil
}

func (p *PreviewCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
