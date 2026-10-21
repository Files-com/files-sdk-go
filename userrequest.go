package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// UserRequest is a Files.com API resource.
type UserRequest struct {
	// ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// User's full name.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// User email address.
	Email string `json:"email,omitempty" path:"email,omitempty" url:"email,omitempty"`
	// Details of the user's request.
	Details string `json:"details,omitempty" path:"details,omitempty" url:"details,omitempty"`
	// User's company name.
	Company string `json:"company,omitempty" path:"company,omitempty" url:"company,omitempty"`
}

// Identifier returns the resource ID.
func (u UserRequest) Identifier() interface{} {
	return u.Id
}

// UserRequestCollection is a list of UserRequest resources.
type UserRequestCollection []UserRequest

// UserRequestListParams contains the request parameters for GET /user_requests.
type UserRequestListParams struct {
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// UserRequestFindParams contains the request parameters for GET /user_requests/{id}.
type UserRequestFindParams struct {
	// User Request ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UserRequestCreateParams contains the request parameters for POST /user_requests.
type UserRequestCreateParams struct {
	// Name of user requested. Required.
	Name string `url:"name" json:"name" path:"name"`
	// Email of user requested. Required.
	Email string `url:"email" json:"email" path:"email"`
	// Details of the user request. Required.
	Details string `url:"details" json:"details" path:"details"`
	// Company of the user requested.
	Company string `url:"company,omitempty" json:"company,omitempty" path:"company"`
}

// UserRequestDeleteParams contains the request parameters for DELETE /user_requests/{id}.
type UserRequestDeleteParams struct {
	// User Request ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (u *UserRequest) UnmarshalJSON(data []byte) error {
	type userRequest UserRequest
	var v userRequest
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*u = UserRequest(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (u *UserRequestCollection) UnmarshalJSON(data []byte) error {
	type userRequests UserRequestCollection
	var v userRequests
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*u = UserRequestCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (u *UserRequestCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*u))
	for i, v := range *u {
		ret[i] = v
	}

	return &ret
}
