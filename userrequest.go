package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type UserRequest struct {
	Id      int64  `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
	Details string `json:"details,omitempty"`
}

type UserRequestCollection []UserRequest

type UserRequestListParams struct {
	Cursor  string `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage int64  `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	lib.ListParams
}

type UserRequestFindParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
}

type UserRequestCreateParams struct {
	Name    string `url:"name,omitempty" required:"true" json:"name,omitempty"`
	Email   string `url:"email,omitempty" required:"true" json:"email,omitempty"`
	Details string `url:"details,omitempty" required:"true" json:"details,omitempty"`
}

type UserRequestDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
}

func (u *UserRequest) UnmarshalJSON(data []byte) error {
	type userRequest UserRequest
	var v userRequest
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*u = UserRequest(v)
	return nil
}

func (u *UserRequestCollection) UnmarshalJSON(data []byte) error {
	type userRequests []UserRequest
	var v userRequests
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*u = UserRequestCollection(v)
	return nil
}

func (u *UserRequestCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*u))
	for i, v := range *u {
		ret[i] = v
	}

	return &ret
}
