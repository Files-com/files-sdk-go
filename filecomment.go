package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type FileComment struct {
	Id        int64               `json:"id,omitempty"`
	Body      string              `json:"body,omitempty"`
	Reactions FileCommentReaction `json:"reactions,omitempty"`
	Path      string              `json:"path,omitempty"`
}

type FileCommentCollection []FileComment

type FileCommentListForParams struct {
	Cursor  string `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage int64  `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	Path    string `url:"-,omitempty" required:"true" json:"-,omitempty"`
	lib.ListParams
}

type FileCommentCreateParams struct {
	Body string `url:"body,omitempty" required:"true" json:"body,omitempty"`
	Path string `url:"path,omitempty" required:"true" json:"path,omitempty"`
}

type FileCommentUpdateParams struct {
	Id   int64  `url:"-,omitempty" required:"true" json:"-,omitempty"`
	Body string `url:"body,omitempty" required:"true" json:"body,omitempty"`
}

type FileCommentDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
}

func (f *FileComment) UnmarshalJSON(data []byte) error {
	type fileComment FileComment
	var v fileComment
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*f = FileComment(v)
	return nil
}

func (f *FileCommentCollection) UnmarshalJSON(data []byte) error {
	type fileComments []FileComment
	var v fileComments
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*f = FileCommentCollection(v)
	return nil
}

func (f *FileCommentCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*f))
	for i, v := range *f {
		ret[i] = v
	}

	return &ret
}
