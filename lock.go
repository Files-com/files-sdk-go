package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Lock struct {
	Path                 string `json:"path,omitempty" path:"path"`
	Timeout              int64  `json:"timeout,omitempty" path:"timeout"`
	Depth                string `json:"depth,omitempty" path:"depth"`
	Recursive            *bool  `json:"recursive,omitempty" path:"recursive"`
	Owner                string `json:"owner,omitempty" path:"owner"`
	Scope                string `json:"scope,omitempty" path:"scope"`
	Exclusive            *bool  `json:"exclusive,omitempty" path:"exclusive"`
	Token                string `json:"token,omitempty" path:"token"`
	Type                 string `json:"type,omitempty" path:"type"`
	AllowAccessByAnyUser *bool  `json:"allow_access_by_any_user,omitempty" path:"allow_access_by_any_user"`
	UserId               int64  `json:"user_id,omitempty" path:"user_id"`
	Username             string `json:"username,omitempty" path:"username"`
}

func (l Lock) Identifier() interface{} {
	return l.Path
}

type LockCollection []Lock

type LockListForParams struct {
	Path            string `url:"-,omitempty" required:"false" json:"-,omitempty" path:"path"`
	IncludeChildren *bool  `url:"include_children,omitempty" required:"false" json:"include_children,omitempty" path:"include_children"`
	ListParams
}

type LockCreateParams struct {
	Path                 string `url:"-,omitempty" required:"false" json:"-,omitempty" path:"path"`
	AllowAccessByAnyUser *bool  `url:"allow_access_by_any_user,omitempty" required:"false" json:"allow_access_by_any_user,omitempty" path:"allow_access_by_any_user"`
	Exclusive            *bool  `url:"exclusive,omitempty" required:"false" json:"exclusive,omitempty" path:"exclusive"`
	Recursive            string `url:"recursive,omitempty" required:"false" json:"recursive,omitempty" path:"recursive"`
	Timeout              int64  `url:"timeout,omitempty" required:"false" json:"timeout,omitempty" path:"timeout"`
}

type LockDeleteParams struct {
	Path  string `url:"-,omitempty" required:"false" json:"-,omitempty" path:"path"`
	Token string `url:"token,omitempty" required:"true" json:"token,omitempty" path:"token"`
}

func (l *Lock) UnmarshalJSON(data []byte) error {
	type lock Lock
	var v lock
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*l = Lock(v)
	return nil
}

func (l *LockCollection) UnmarshalJSON(data []byte) error {
	type locks LockCollection
	var v locks
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
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
