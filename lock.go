package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Lock struct {
	Path                 string `json:"path,omitempty"`
	Timeout              int64  `json:"timeout,omitempty"`
	Depth                string `json:"depth,omitempty"`
	Recursive            *bool  `json:"recursive,omitempty"`
	Owner                string `json:"owner,omitempty"`
	Scope                string `json:"scope,omitempty"`
	Exclusive            *bool  `json:"exclusive,omitempty"`
	Token                string `json:"token,omitempty"`
	Type                 string `json:"type,omitempty"`
	AllowAccessByAnyUser *bool  `json:"allow_access_by_any_user,omitempty"`
	UserId               int64  `json:"user_id,omitempty"`
	Username             string `json:"username,omitempty"`
}

type LockCollection []Lock

type LockListForParams struct {
	Cursor          string `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage         int64  `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	Path            string `url:"-,omitempty" required:"true" json:"-,omitempty"`
	IncludeChildren *bool  `url:"include_children,omitempty" required:"false" json:"include_children,omitempty"`
	lib.ListParams
}

type LockCreateParams struct {
	Path                 string `url:"-,omitempty" required:"true" json:"-,omitempty"`
	AllowAccessByAnyUser *bool  `url:"allow_access_by_any_user,omitempty" required:"false" json:"allow_access_by_any_user,omitempty"`
	Exclusive            *bool  `url:"exclusive,omitempty" required:"false" json:"exclusive,omitempty"`
	Recursive            string `url:"recursive,omitempty" required:"false" json:"recursive,omitempty"`
	Timeout              int64  `url:"timeout,omitempty" required:"false" json:"timeout,omitempty"`
}

type LockDeleteParams struct {
	Path  string `url:"-,omitempty" required:"true" json:"-,omitempty"`
	Token string `url:"token,omitempty" required:"true" json:"token,omitempty"`
}

func (l *Lock) UnmarshalJSON(data []byte) error {
	type lock Lock
	var v lock
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*l = Lock(v)
	return nil
}

func (l *LockCollection) UnmarshalJSON(data []byte) error {
	type locks []Lock
	var v locks
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*l = LockCollection(v)
	return nil
}

func (l *LockCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*l))
	for i, v := range *l {
		ret[i] = v
	}

	return &ret
}
