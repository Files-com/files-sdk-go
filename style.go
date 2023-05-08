package files_sdk

import (
	"encoding/json"
	"io"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Style struct {
	Id        int64     `json:"id,omitempty" path:"id"`
	Path      string    `json:"path,omitempty" path:"path"`
	Logo      Image     `json:"logo,omitempty" path:"logo"`
	Thumbnail Image     `json:"thumbnail,omitempty" path:"thumbnail"`
	File      io.Reader `json:"file,omitempty" path:"file"`
}

func (s Style) Identifier() interface{} {
	return s.Id
}

type StyleCollection []Style

type StyleFindParams struct {
	Path string `url:"-,omitempty" required:"false" json:"-,omitempty" path:"path"`
}

type StyleUpdateParams struct {
	Path string    `url:"-,omitempty" required:"false" json:"-,omitempty" path:"path"`
	File io.Writer `url:"file,omitempty" required:"true" json:"file,omitempty" path:"file"`
}

type StyleDeleteParams struct {
	Path string `url:"-,omitempty" required:"false" json:"-,omitempty" path:"path"`
}

func (s *Style) UnmarshalJSON(data []byte) error {
	type style Style
	var v style
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = Style(v)
	return nil
}

func (s *StyleCollection) UnmarshalJSON(data []byte) error {
	type styles StyleCollection
	var v styles
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
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
