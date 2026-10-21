package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// UserAdditionalEmailRecipient is a Files.com API resource.
type UserAdditionalEmailRecipient struct {
	// User additional email recipient ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// User ID.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	// Workspace ID (0 for default workspace).
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Additional email recipient address.
	Email string `json:"email,omitempty" path:"email,omitempty" url:"email,omitempty"`
	// Created at date/time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
}

// Identifier returns the resource ID.
func (u UserAdditionalEmailRecipient) Identifier() interface{} {
	return u.Id
}

// UserAdditionalEmailRecipientCollection is a list of UserAdditionalEmailRecipient resources.
type UserAdditionalEmailRecipientCollection []UserAdditionalEmailRecipient

// UserAdditionalEmailRecipientListParams contains the request parameters for GET /user_additional_email_recipients.
type UserAdditionalEmailRecipientListParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `email`, `user_id` or `workspace_id`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `email`, `workspace_id` or `user_id`. Valid field combinations are `[ workspace_id, email ]`, `[ user_id, email ]`, `[ workspace_id, user_id ]` or `[ workspace_id, user_id, email ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is prefixed by the supplied value. Valid fields are `email`.
	FilterPrefix interface{} `url:"filter_prefix,omitempty" json:"filter_prefix,omitempty" path:"filter_prefix"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// UserAdditionalEmailRecipientFindParams contains the request parameters for GET /user_additional_email_recipients/{id}.
type UserAdditionalEmailRecipientFindParams struct {
	// User Additional Email Recipient ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UserAdditionalEmailRecipientCreateParams contains the request parameters for POST /user_additional_email_recipients.
type UserAdditionalEmailRecipientCreateParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// Additional email recipient address. Required.
	Email string `url:"email" json:"email" path:"email"`
}

// UserAdditionalEmailRecipientUpdateParams contains the request parameters for PATCH /user_additional_email_recipients/{id}.
type UserAdditionalEmailRecipientUpdateParams struct {
	// User Additional Email Recipient ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Additional email recipient address.
	Email string `url:"email,omitempty" json:"email,omitempty" path:"email"`
}

// UserAdditionalEmailRecipientDeleteParams contains the request parameters for DELETE /user_additional_email_recipients/{id}.
type UserAdditionalEmailRecipientDeleteParams struct {
	// User Additional Email Recipient ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (u *UserAdditionalEmailRecipient) UnmarshalJSON(data []byte) error {
	type userAdditionalEmailRecipient UserAdditionalEmailRecipient
	var v userAdditionalEmailRecipient
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*u = UserAdditionalEmailRecipient(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (u *UserAdditionalEmailRecipientCollection) UnmarshalJSON(data []byte) error {
	type userAdditionalEmailRecipients UserAdditionalEmailRecipientCollection
	var v userAdditionalEmailRecipients
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*u = UserAdditionalEmailRecipientCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (u *UserAdditionalEmailRecipientCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*u))
	for i, v := range *u {
		ret[i] = v
	}

	return &ret
}
