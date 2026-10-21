package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// EmailIncomingMessage is a Files.com API resource.
type EmailIncomingMessage struct {
	// Id of the Email Incoming Message.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Id of the Inbox associated with this message.
	InboxId int64 `json:"inbox_id,omitempty" path:"inbox_id,omitempty" url:"inbox_id,omitempty"`
	// Sender of the email.
	Sender string `json:"sender,omitempty" path:"sender,omitempty" url:"sender,omitempty"`
	// Sender name.
	SenderName string `json:"sender_name,omitempty" path:"sender_name,omitempty" url:"sender_name,omitempty"`
	// Status of the message.
	Status string `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	// Body of the email.
	Body string `json:"body,omitempty" path:"body,omitempty" url:"body,omitempty"`
	// Message describing the failure.
	Message string `json:"message,omitempty" path:"message,omitempty" url:"message,omitempty"`
	// Message creation date/time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Title of the Inbox associated with this message.
	InboxTitle string `json:"inbox_title,omitempty" path:"inbox_title,omitempty" url:"inbox_title,omitempty"`
}

// Identifier returns the resource ID.
func (e EmailIncomingMessage) Identifier() interface{} {
	return e.Id
}

// EmailIncomingMessageCollection is a list of EmailIncomingMessage resources.
type EmailIncomingMessageCollection []EmailIncomingMessage

// EmailIncomingMessageListParams contains the request parameters for GET /email_incoming_messages.
type EmailIncomingMessageListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`, `sender`, `status` or `inbox_id`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at`, `inbox_id`, `sender` or `status`. Valid field combinations are `[ inbox_id, created_at ]`, `[ sender, created_at ]`, `[ status, created_at ]`, `[ inbox_id, status ]`, `[ status, sender ]`, `[ inbox_id, status, created_at ]`, `[ status, sender, created_at ]`, `[ inbox_id, status, sender ]` or `[ inbox_id, status, sender, created_at ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`.
	FilterGt interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	// If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`.
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	// If set, return records where the specified field is prefixed by the supplied value. Valid fields are `sender`.
	FilterPrefix interface{} `url:"filter_prefix,omitempty" json:"filter_prefix,omitempty" path:"filter_prefix"`
	// If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`.
	FilterLt interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	// If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`.
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
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
