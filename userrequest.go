package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type UserRequest struct {
	Id      int64  `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
	Details string `json:"details,omitempty"`
}

type UserRequestCollection []UserRequest

type UserRequestListParams struct {
	Page    int    `url:"page,omitempty"`
	PerPage int    `url:"per_page,omitempty"`
	Action  string `url:"action,omitempty"`
	Cursor  string `url:"cursor,omitempty"`
	lib.ListParams
}

type UserRequestFindParams struct {
	Id int64 `url:"-,omitempty"`
}

type UserRequestCreateParams struct {
	Name    string `url:"name,omitempty"`
	Email   string `url:"email,omitempty"`
	Details string `url:"details,omitempty"`
}

type UserRequestDeleteParams struct {
	Id int64 `url:"-,omitempty"`
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
