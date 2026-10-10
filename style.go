package files_sdk

import (
	"encoding/json"
	"io"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Style is a Files.com API resource.
type Style struct {
	Id            int64     `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Path          string    `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	Logo          Image     `json:"logo,omitempty" path:"logo,omitempty" url:"logo,omitempty"`
	LogoClickHref string    `json:"logo_click_href,omitempty" path:"logo_click_href,omitempty" url:"logo_click_href,omitempty"`
	Thumbnail     Image     `json:"thumbnail,omitempty" path:"thumbnail,omitempty" url:"thumbnail,omitempty"`
	File          io.Reader `json:"file,omitempty" path:"file,omitempty" url:"file,omitempty"`
}

// Identifier returns the resource ID.
func (s Style) Identifier() interface{} {
	return s.Id
}

// StyleCollection is a list of Style resources.
type StyleCollection []Style

// StyleFindParams contains the request parameters for GET /styles/{path}.
type StyleFindParams struct {
	Path string `url:"-,omitempty" json:"-,omitempty" path:"path"`
}

// StyleUpdateParams contains the request parameters for PATCH /styles/{path}.
type StyleUpdateParams struct {
	Path          string    `url:"-,omitempty" json:"-,omitempty" path:"path"`
	File          io.Writer `url:"file,omitempty" json:"file,omitempty" path:"file"`
	LogoClickHref string    `url:"logo_click_href,omitempty" json:"logo_click_href,omitempty" path:"logo_click_href"`
}

// StyleDeleteParams contains the request parameters for DELETE /styles/{path}.
type StyleDeleteParams struct {
	Path string `url:"-,omitempty" json:"-,omitempty" path:"path"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (s *Style) UnmarshalJSON(data []byte) error {
	type style Style
	var v style
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = Style(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (s *StyleCollection) UnmarshalJSON(data []byte) error {
	type styles StyleCollection
	var v styles
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = StyleCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (s *StyleCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
