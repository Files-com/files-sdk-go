package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// UserCipherUse is a Files.com API resource.
type UserCipherUse struct {
	// UserCipherUse ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// ID of the user who performed this access.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	// Username of the user who performed this access.
	Username string `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	// The protocol and cipher employed.
	ProtocolCipher string `json:"protocol_cipher,omitempty" path:"protocol_cipher,omitempty" url:"protocol_cipher,omitempty"`
	// The earliest recorded use of this combination of interface and protocol and cipher (for this user).
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Is this cipher considered insecure?
	Insecure *bool `json:"insecure,omitempty" path:"insecure,omitempty" url:"insecure,omitempty"`
	// The interface accessed.
	Interface string `json:"interface,omitempty" path:"interface,omitempty" url:"interface,omitempty"`
	// The most recent use of this combination of interface and protocol and cipher (for this user).
	UpdatedAt *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
}

// Identifier returns the resource ID.
func (u UserCipherUse) Identifier() interface{} {
	return u.Id
}

// UserCipherUseCollection is a list of UserCipherUse resources.
type UserCipherUseCollection []UserCipherUse

// UserCipherUseListParams contains the request parameters for GET /user_cipher_uses.
type UserCipherUseListParams struct {
	// User ID. If provided, will return uses for this user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `updated_at`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `insecure` and `updated_at`. Valid field combinations are `[ insecure, updated_at ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is greater than the supplied value. Valid fields are `updated_at`.
	FilterGt interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	// If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `updated_at`.
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	// If set, return records where the specified field is less than the supplied value. Valid fields are `updated_at`.
	FilterLt interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	// If set, return records where the specified field is less than or equal the supplied value. Valid fields are `updated_at`.
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (u *UserCipherUse) UnmarshalJSON(data []byte) error {
	type userCipherUse UserCipherUse
	var v userCipherUse
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*u = UserCipherUse(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (u *UserCipherUseCollection) UnmarshalJSON(data []byte) error {
	type userCipherUses UserCipherUseCollection
	var v userCipherUses
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*u = UserCipherUseCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (u *UserCipherUseCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*u))
	for i, v := range *u {
		ret[i] = v
	}

	return &ret
}
