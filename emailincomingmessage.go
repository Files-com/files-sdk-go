package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// EmailIncomingMessage is a Files.com API resource.
type EmailIncomingMessage struct {
	Id         int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	InboxId    int64      `json:"inbox_id,omitempty" path:"inbox_id,omitempty" url:"inbox_id,omitempty"`
	Sender     string     `json:"sender,omitempty" path:"sender,omitempty" url:"sender,omitempty"`
	SenderName string     `json:"sender_name,omitempty" path:"sender_name,omitempty" url:"sender_name,omitempty"`
	Status     string     `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	Body       string     `json:"body,omitempty" path:"body,omitempty" url:"body,omitempty"`
	Message    string     `json:"message,omitempty" path:"message,omitempty" url:"message,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	InboxTitle string     `json:"inbox_title,omitempty" path:"inbox_title,omitempty" url:"inbox_title,omitempty"`
}

// Identifier returns the resource ID.
func (e EmailIncomingMessage) Identifier() interface{} {
	return e.Id
}

// EmailIncomingMessageCollection is a list of EmailIncomingMessage resources.
type EmailIncomingMessageCollection []EmailIncomingMessage

// EmailIncomingMessageListParams contains the request parameters for GET /email_incoming_messages.
type EmailIncomingMessageListParams struct {
	SortBy       interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	Filter       interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	FilterGt     interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq   interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterPrefix interface{} `url:"filter_prefix,omitempty" json:"filter_prefix,omitempty" path:"filter_prefix"`
	FilterLt     interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq   interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (e *EmailIncomingMessage) UnmarshalJSON(data []byte) error {
	type emailIncomingMessage EmailIncomingMessage
	var v emailIncomingMessage
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*e = EmailIncomingMessage(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (e *EmailIncomingMessageCollection) UnmarshalJSON(data []byte) error {
	type emailIncomingMessages EmailIncomingMessageCollection
	var v emailIncomingMessages
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*e = EmailIncomingMessageCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (e *EmailIncomingMessageCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*e))
	for i, v := range *e {
		ret[i] = v
	}

	return &ret
}
