package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type Message struct {
	Id        int64    `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Subject   string   `json:"subject,omitempty" path:"subject,omitempty" url:"subject,omitempty"`
	Body      string   `json:"body,omitempty" path:"body,omitempty" url:"body,omitempty"`
	Comments  []string `json:"comments,omitempty" path:"comments,omitempty" url:"comments,omitempty"`
	UserId    int64    `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	ProjectId int64    `json:"project_id,omitempty" path:"project_id,omitempty" url:"project_id,omitempty"`
}

func (m Message) Identifier() interface{} {
	return m.Id
}

type MessageCollection []Message

type MessageListParams struct {
	UserId    int64 `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	ProjectId int64 `url:"project_id,omitempty" required:"true" json:"project_id,omitempty" path:"project_id"`
	ListParams
}

type MessageFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type MessageCreateParams struct {
	UserId    int64  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	ProjectId int64  `url:"project_id,omitempty" required:"true" json:"project_id,omitempty" path:"project_id"`
	Subject   string `url:"subject,omitempty" required:"true" json:"subject,omitempty" path:"subject"`
	Body      string `url:"body,omitempty" required:"true" json:"body,omitempty" path:"body"`
}

type MessageUpdateParams struct {
	Id        int64  `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
	ProjectId int64  `url:"project_id,omitempty" required:"true" json:"project_id,omitempty" path:"project_id"`
	Subject   string `url:"subject,omitempty" required:"true" json:"subject,omitempty" path:"subject"`
	Body      string `url:"body,omitempty" required:"true" json:"body,omitempty" path:"body"`
}

type MessageDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

func (m *Message) UnmarshalJSON(data []byte) error {
	type message Message
	var v message
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*m = Message(v)
	return nil
}

func (m *MessageCollection) UnmarshalJSON(data []byte) error {
	type messages MessageCollection
	var v messages
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
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
