package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// BundleRecipient is a Files.com API resource.
type BundleRecipient struct {
	// The recipient's company.
	Company string `json:"company,omitempty" path:"company,omitempty" url:"company,omitempty"`
	// The recipient's name.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// A note sent to the recipient with the bundle.
	Note string `json:"note,omitempty" path:"note,omitempty" url:"note,omitempty"`
	// The recipient's email address.
	Recipient string `json:"recipient,omitempty" path:"recipient,omitempty" url:"recipient,omitempty"`
	// When the Bundle was shared with this recipient.
	SentAt *time.Time `json:"sent_at,omitempty" path:"sent_at,omitempty" url:"sent_at,omitempty"`
	// Workspace ID. `0` means the default workspace.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	// Bundle to share.
	BundleId int64 `json:"bundle_id,omitempty" path:"bundle_id,omitempty" url:"bundle_id,omitempty"`
	// Set to true to share the link with the recipient upon creation.
	ShareAfterCreate *bool `json:"share_after_create,omitempty" path:"share_after_create,omitempty" url:"share_after_create,omitempty"`
}

// Identifier no path or id

// BundleRecipientCollection is a list of BundleRecipient resources.
type BundleRecipientCollection []BundleRecipient

// BundleRecipientListParams contains the request parameters for GET /bundle_recipients.
type BundleRecipientListParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `has_registrations`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// List recipients for the bundle with this ID. Required.
	BundleId int64 `url:"bundle_id" json:"bundle_id" path:"bundle_id"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// BundleRecipientCreateParams contains the request parameters for POST /bundle_recipients.
type BundleRecipientCreateParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// Bundle to share. Required.
	BundleId int64 `url:"bundle_id" json:"bundle_id" path:"bundle_id"`
	// Email addresses to share this bundle with. Required.
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
func (b *BundleRecipient) UnmarshalJSON(data []byte) error {
	type bundleRecipient BundleRecipient
	var v bundleRecipient
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*b = BundleRecipient(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (b *BundleRecipientCollection) UnmarshalJSON(data []byte) error {
	type bundleRecipients BundleRecipientCollection
	var v bundleRecipients
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*b = BundleRecipientCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (b *BundleRecipientCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*b))
	for i, v := range *b {
		ret[i] = v
	}

	return &ret
}
