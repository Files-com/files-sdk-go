package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type FileComment struct {
	Id        int64    `json:"id,omitempty"`
	Body      string   `json:"body,omitempty"`
	Reactions []string `json:"reactions,omitempty"`
	Path      string   `json:"path,omitempty"`
}

type FileCommentCollection []FileComment

type FileCommentListForParams struct {
	Cursor  string `url:"cursor,omitempty" required:"false"`
	PerPage int    `url:"per_page,omitempty" required:"false"`
	Path    string `url:"-,omitempty" required:"true"`
	lib.ListParams
}

type FileCommentCreateParams struct {
	Body string `url:"body,omitempty" required:"true"`
	Path string `url:"path,omitempty" required:"true"`
}

type FileCommentUpdateParams struct {
	Id   int64  `url:"-,omitempty" required:"true"`
	Body string `url:"body,omitempty" required:"true"`
}

type FileCommentDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
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
