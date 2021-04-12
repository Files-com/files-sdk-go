package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type Message struct {
	Id        int64  `json:"id,omitempty"`
	Subject   string `json:"subject,omitempty"`
	Body      string `json:"body,omitempty"`
	Comments  string `json:"comments,omitempty"`
	UserId    int64  `json:"user_id,omitempty"`
	ProjectId int64  `json:"project_id,omitempty"`
}

type MessageCollection []Message

type MessageListParams struct {
	UserId    int64  `url:"user_id,omitempty" required:"false"`
	Cursor    string `url:"cursor,omitempty" required:"false"`
	PerPage   int64  `url:"per_page,omitempty" required:"false"`
	ProjectId int64  `url:"project_id,omitempty" required:"true"`
	lib.ListParams
}

type MessageFindParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

type MessageCreateParams struct {
	UserId    int64  `url:"user_id,omitempty" required:"false"`
	ProjectId int64  `url:"project_id,omitempty" required:"true"`
	Subject   string `url:"subject,omitempty" required:"true"`
	Body      string `url:"body,omitempty" required:"true"`
}

type MessageUpdateParams struct {
	Id        int64  `url:"-,omitempty" required:"true"`
	ProjectId int64  `url:"project_id,omitempty" required:"true"`
	Subject   string `url:"subject,omitempty" required:"true"`
	Body      string `url:"body,omitempty" required:"true"`
}

type MessageDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
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
