package files_sdk

import (
	"encoding/json"
	"io"
)

type Style struct {
	Id        int64     `json:"id,omitempty"`
	Path      string    `json:"path,omitempty"`
	Logo      string    `json:"logo,omitempty"`
	Thumbnail string    `json:"thumbnail,omitempty"`
	File      io.Reader `json:"file,omitempty"`
}

type StyleCollection []Style

type StyleFindParams struct {
	Path string `url:"-,omitempty" required:"true"`
}

type StyleUpdateParams struct {
	Path string    `url:"-,omitempty" required:"true"`
	File io.Writer `url:"file,omitempty" required:"true"`
}

type StyleDeleteParams struct {
	Path string `url:"-,omitempty" required:"true"`
}

func (s *Style) UnmarshalJSON(data []byte) error {
	type style Style
	var v style
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*s = Style(v)
	return nil
}

func (s *StyleCollection) UnmarshalJSON(data []byte) error {
	type styles []Style
	var v styles
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*s = StyleCollection(v)
	return nil
}

func (s *StyleCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
