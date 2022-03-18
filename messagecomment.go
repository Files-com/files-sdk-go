package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type MessageComment struct {
	Id        int64                  `json:"id,omitempty"`
	Body      string                 `json:"body,omitempty"`
	Reactions MessageCommentReaction `json:"reactions,omitempty"`
	UserId    int64                  `json:"user_id,omitempty"`
}

type MessageCommentCollection []MessageComment

type MessageCommentListParams struct {
	UserId    int64  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty"`
	Cursor    string `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage   int64  `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	MessageId int64  `url:"message_id,omitempty" required:"true" json:"message_id,omitempty"`
	lib.ListParams
}

type MessageCommentFindParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
}

type MessageCommentCreateParams struct {
	UserId int64  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty"`
	Body   string `url:"body,omitempty" required:"true" json:"body,omitempty"`
}

type MessageCommentUpdateParams struct {
	Id   int64  `url:"-,omitempty" required:"true" json:"-,omitempty"`
	Body string `url:"body,omitempty" required:"true" json:"body,omitempty"`
}

type MessageCommentDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
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

func (m *MessageCommentCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*m))
	for i, v := range *m {
		ret[i] = v
	}

	return &ret
}
