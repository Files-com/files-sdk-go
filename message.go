package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Message struct {
	Id        int64          `json:"id,omitempty"`
	Subject   string         `json:"subject,omitempty"`
	Body      string         `json:"body,omitempty"`
	Comments  MessageComment `json:"comments,omitempty"`
	UserId    int64          `json:"user_id,omitempty"`
	ProjectId int64          `json:"project_id,omitempty"`
}

type MessageCollection []Message

type MessageListParams struct {
	UserId    int64  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty"`
	Cursor    string `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage   int64  `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	ProjectId int64  `url:"project_id,omitempty" required:"true" json:"project_id,omitempty"`
	lib.ListParams
}

type MessageFindParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
}

type MessageCreateParams struct {
	UserId    int64  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty"`
	ProjectId int64  `url:"project_id,omitempty" required:"true" json:"project_id,omitempty"`
	Subject   string `url:"subject,omitempty" required:"true" json:"subject,omitempty"`
	Body      string `url:"body,omitempty" required:"true" json:"body,omitempty"`
}

type MessageUpdateParams struct {
	Id        int64  `url:"-,omitempty" required:"true" json:"-,omitempty"`
	ProjectId int64  `url:"project_id,omitempty" required:"true" json:"project_id,omitempty"`
	Subject   string `url:"subject,omitempty" required:"true" json:"subject,omitempty"`
	Body      string `url:"body,omitempty" required:"true" json:"body,omitempty"`
}

type MessageDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
}

func (m *Message) UnmarshalJSON(data []byte) error {
	type message Message
	var v message
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*m = Message(v)
	return nil
}

func (m *MessageCollection) UnmarshalJSON(data []byte) error {
	type messages []Message
	var v messages
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*m = MessageCollection(v)
	return nil
}

func (m *MessageCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*m))
	for i, v := range *m {
		ret[i] = v
	}

	return &ret
}
