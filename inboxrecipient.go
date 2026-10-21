package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// InboxRecipient is a Files.com API resource.
type InboxRecipient struct {
	// The recipient's company.
	Company string `json:"company,omitempty" path:"company,omitempty" url:"company,omitempty"`
	// The recipient's name.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// A note sent to the recipient with the inbox.
	Note string `json:"note,omitempty" path:"note,omitempty" url:"note,omitempty"`
	// The recipient's email address.
	Recipient string `json:"recipient,omitempty" path:"recipient,omitempty" url:"recipient,omitempty"`
	// When the Inbox was shared with this recipient.
	SentAt *time.Time `json:"sent_at,omitempty" path:"sent_at,omitempty" url:"sent_at,omitempty"`
	// Inbox to share.
	InboxId int64 `json:"inbox_id,omitempty" path:"inbox_id,omitempty" url:"inbox_id,omitempty"`
	// Set to true to share the link with the recipient upon creation.
	ShareAfterCreate *bool `json:"share_after_create,omitempty" path:"share_after_create,omitempty" url:"share_after_create,omitempty"`
}

// Identifier no path or id

// InboxRecipientCollection is a list of InboxRecipient resources.
type InboxRecipientCollection []InboxRecipient

// InboxRecipientListParams contains the request parameters for GET /inbox_recipients.
type InboxRecipientListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are .
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `has_registrations`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// List recipients for the inbox with this ID. Required.
	InboxId int64 `url:"inbox_id" json:"inbox_id" path:"inbox_id"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// InboxRecipientCreateParams contains the request parameters for POST /inbox_recipients.
type InboxRecipientCreateParams struct {
	// Inbox to share. Required.
	InboxId int64 `url:"inbox_id" json:"inbox_id" path:"inbox_id"`
	// Email address to share this inbox with. Required.
	Recipient string `url:"recipient" json:"recipient" path:"recipient"`
	// Name of recipient.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Company of recipient.
	Company string `url:"company,omitempty" json:"company,omitempty" path:"company"`
	// Note to include in email.
	Note string `url:"note,omitempty" json:"note,omitempty" path:"note"`
	// Set to true to share the link with the recipient upon creation.
	ShareAfterCreate *bool `url:"share_after_create,omitempty" json:"share_after_create,omitempty" path:"share_after_create"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (i *InboxRecipient) UnmarshalJSON(data []byte) error {
	type inboxRecipient InboxRecipient
	var v inboxRecipient
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*i = InboxRecipient(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (i *InboxRecipientCollection) UnmarshalJSON(data []byte) error {
	type inboxRecipients InboxRecipientCollection
	var v inboxRecipients
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*i = InboxRecipientCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (i *InboxRecipientCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*i))
	for i, v := range *i {
		ret[i] = v
	}

	return &ret
}
