package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type FileComment struct {
	Id        int64               `json:"id,omitempty" path:"id"`
	Body      string              `json:"body,omitempty" path:"body"`
	Reactions FileCommentReaction `json:"reactions,omitempty" path:"reactions"`
	Path      string              `json:"path,omitempty" path:"path"`
}

type FileCommentCollection []FileComment

type FileCommentListForParams struct {
	Path string `url:"-,omitempty" required:"true" json:"-,omitempty" path:"path"`
	lib.ListParams
}

type FileCommentCreateParams struct {
	Body string `url:"body,omitempty" required:"true" json:"body,omitempty" path:"body"`
	Path string `url:"path,omitempty" required:"true" json:"path,omitempty" path:"path"`
}

type FileCommentUpdateParams struct {
	Id   int64  `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
	Body string `url:"body,omitempty" required:"true" json:"body,omitempty" path:"body"`
}

type FileCommentDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
}

func (f *FileComment) UnmarshalJSON(data []byte) error {
	type fileComment FileComment
	var v fileComment
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*f = FileComment(v)
	return nil
}

func (f *FileCommentCollection) UnmarshalJSON(data []byte) error {
	type fileComments FileCommentCollection
	var v fileComments
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
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
