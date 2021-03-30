package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type MessageComment struct {
	Id        int64    `json:"id,omitempty"`
	Body      string   `json:"body,omitempty"`
	Reactions []string `json:"reactions,omitempty"`
	UserId    int64    `json:"user_id,omitempty"`
}

type MessageCommentCollection []MessageComment

type MessageCommentListParams struct {
	UserId    int64  `url:"user_id,omitempty" required:"false"`
	Cursor    string `url:"cursor,omitempty" required:"false"`
	PerPage   int    `url:"per_page,omitempty" required:"false"`
	MessageId int64  `url:"message_id,omitempty" required:"true"`
	lib.ListParams
}

type MessageCommentFindParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

type MessageCommentCreateParams struct {
	UserId int64  `url:"user_id,omitempty" required:"false"`
	Body   string `url:"body,omitempty" required:"true"`
}

type MessageCommentUpdateParams struct {
	Id   int64  `url:"-,omitempty" required:"true"`
	Body string `url:"body,omitempty" required:"true"`
}

type MessageCommentDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

func (m *MessageComment) UnmarshalJSON(data []byte) error {
	type messageComment MessageComment
	var v messageComment
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*m = MessageComment(v)
	return nil
}

func (m *MessageCommentCollection) UnmarshalJSON(data []byte) error {
	type messageComments []MessageComment
	var v messageComments
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*m = MessageCommentCollection(v)
	return nil
}
