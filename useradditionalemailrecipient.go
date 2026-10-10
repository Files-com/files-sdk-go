package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// UserAdditionalEmailRecipient is a Files.com API resource.
type UserAdditionalEmailRecipient struct {
	Id          int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	UserId      int64      `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	WorkspaceId int64      `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	Email       string     `json:"email,omitempty" path:"email,omitempty" url:"email,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
}

// Identifier returns the resource ID.
func (u UserAdditionalEmailRecipient) Identifier() interface{} {
	return u.Id
}

// UserAdditionalEmailRecipientCollection is a list of UserAdditionalEmailRecipient resources.
type UserAdditionalEmailRecipientCollection []UserAdditionalEmailRecipient

// UserAdditionalEmailRecipientListParams contains the request parameters for GET /user_additional_email_recipients.
type UserAdditionalEmailRecipientListParams struct {
	UserId       int64       `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	SortBy       interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	Filter       interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	FilterPrefix interface{} `url:"filter_prefix,omitempty" json:"filter_prefix,omitempty" path:"filter_prefix"`
	ListParams
}

// UserAdditionalEmailRecipientFindParams contains the request parameters for GET /user_additional_email_recipients/{id}.
type UserAdditionalEmailRecipientFindParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UserAdditionalEmailRecipientCreateParams contains the request parameters for POST /user_additional_email_recipients.
type UserAdditionalEmailRecipientCreateParams struct {
	UserId int64  `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	Email  string `url:"email" json:"email" path:"email"`
}

// UserAdditionalEmailRecipientUpdateParams contains the request parameters for PATCH /user_additional_email_recipients/{id}.
type UserAdditionalEmailRecipientUpdateParams struct {
	Id    int64  `url:"-,omitempty" json:"-,omitempty" path:"id"`
	Email string `url:"email,omitempty" json:"email,omitempty" path:"email"`
}

// UserAdditionalEmailRecipientDeleteParams contains the request parameters for DELETE /user_additional_email_recipients/{id}.
type UserAdditionalEmailRecipientDeleteParams struct {
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
