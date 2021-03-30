package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type Lock struct {
	Path     string `json:"path,omitempty"`
	Timeout  int    `json:"timeout,omitempty"`
	Depth    string `json:"depth,omitempty"`
	Owner    string `json:"owner,omitempty"`
	Scope    string `json:"scope,omitempty"`
	Token    string `json:"token,omitempty"`
	Type     string `json:"type,omitempty"`
	UserId   int64  `json:"user_id,omitempty"`
	Username string `json:"username,omitempty"`
}

type LockCollection []Lock

type LockListForParams struct {
	Cursor          string `url:"cursor,omitempty" required:"false"`
	PerPage         int    `url:"per_page,omitempty" required:"false"`
	Path            string `url:"-,omitempty" required:"true"`
	IncludeChildren *bool  `url:"include_children,omitempty" required:"false"`
	lib.ListParams
}

type LockCreateParams struct {
	Path    string `url:"-,omitempty" required:"true"`
	Timeout int    `url:"timeout,omitempty" required:"false"`
}

type LockDeleteParams struct {
	Path  string `url:"-,omitempty" required:"true"`
	Token string `url:"token,omitempty" required:"true"`
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
