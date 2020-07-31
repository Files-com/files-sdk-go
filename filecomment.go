package files_sdk

import (
  lib "github.com/Files-com/files-sdk-go/lib"
  "encoding/json"
)

type FileComment struct {
  Id int `json:"id,omitempty"`
  Body string `json:"body,omitempty"`
  Reactions []string `json:"reactions,omitempty"`
  Path string `json:"path,omitempty"`
}

type FileCommentCollection []FileComment

type FileCommentListForParams struct {
  Page int `url:"page,omitempty"`
  PerPage int `url:"per_page,omitempty"`
  Action string `url:"action,omitempty"`
  Path string `url:"-,omitempty"`
  lib.ListParams
}

type FileCommentCreateParams struct {
  Body string `url:"body,omitempty"`
  Path string `url:"path,omitempty"`
}

type FileCommentUpdateParams struct {
  Id int `url:"-,omitempty"`
  Body string `url:"body,omitempty"`
}

type FileCommentDeleteParams struct {
  Id int `url:"-,omitempty"`
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

