package files_sdk

import (
	"encoding/json"
	lib "github.com/Files-com/files-sdk-go/lib"
)

type MessageComment struct {
	Id        int      `json:"id,omitempty"`
	Body      string   `json:"body,omitempty"`
	Reactions []string `json:"reactions,omitempty"`
	UserId    int      `json:"user_id,omitempty"`
}

type MessageCommentCollection []MessageComment

type MessageCommentListParams struct {
	UserId    int    `url:"user_id,omitempty"`
	Page      int    `url:"page,omitempty"`
	PerPage   int    `url:"per_page,omitempty"`
	Action    string `url:"action,omitempty"`
	MessageId int    `url:"message_id,omitempty"`
	lib.ListParams
}

type MessageCommentFindParams struct {
	Id int `url:"-,omitempty"`
}

type MessageCommentCreateParams struct {
	UserId int    `url:"user_id,omitempty"`
	Body   string `url:"body,omitempty"`
}

type MessageCommentUpdateParams struct {
	Id   int    `url:"-,omitempty"`
	Body string `url:"body,omitempty"`
}

type MessageCommentDeleteParams struct {
	Id int `url:"-,omitempty"`
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
