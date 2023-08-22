package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type MessageComment struct {
	Id        int64    `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Body      string   `json:"body,omitempty" path:"body,omitempty" url:"body,omitempty"`
	Reactions []string `json:"reactions,omitempty" path:"reactions,omitempty" url:"reactions,omitempty"`
	UserId    int64    `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
}

func (m MessageComment) Identifier() interface{} {
	return m.Id
}

type MessageCommentCollection []MessageComment

type MessageCommentListParams struct {
	UserId    int64  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	Action    string `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	MessageId int64  `url:"message_id,omitempty" required:"true" json:"message_id,omitempty" path:"message_id"`
	ListParams
}

type MessageCommentFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type MessageCommentCreateParams struct {
	UserId int64  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	Body   string `url:"body,omitempty" required:"true" json:"body,omitempty" path:"body"`
}

type MessageCommentUpdateParams struct {
	Id   int64  `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
	Body string `url:"body,omitempty" required:"true" json:"body,omitempty" path:"body"`
}

type MessageCommentDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

func (m *MessageComment) UnmarshalJSON(data []byte) error {
	type messageComment MessageComment
	var v messageComment
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*m = MessageComment(v)
	return nil
}

func (m *MessageCommentCollection) UnmarshalJSON(data []byte) error {
	type messageComments MessageCommentCollection
	var v messageComments
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
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
