package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Lock is a Files.com API resource.
type Lock struct {
	// Path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.
	Path string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	// Lock timeout in seconds.
	Timeout int64  `json:"timeout,omitempty" path:"timeout,omitempty" url:"timeout,omitempty"`
	Depth   string `json:"depth,omitempty" path:"depth,omitempty" url:"depth,omitempty"`
	// Does lock apply to subfolders?
	Recursive *bool `json:"recursive,omitempty" path:"recursive,omitempty" url:"recursive,omitempty"`
	// Arbitrary descriptive label for the lock. Does not change the lock creator or permissions.
	Owner string `json:"owner,omitempty" path:"owner,omitempty" url:"owner,omitempty"`
	Scope string `json:"scope,omitempty" path:"scope,omitempty" url:"scope,omitempty"`
	// Is lock exclusive?
	Exclusive *bool `json:"exclusive,omitempty" path:"exclusive,omitempty" url:"exclusive,omitempty"`
	// Lock token. Use to release lock.
	Token string `json:"token,omitempty" path:"token,omitempty" url:"token,omitempty"`
	Type  string `json:"type,omitempty" path:"type,omitempty" url:"type,omitempty"`
	// Can lock be modified by users other than its creator?
	AllowAccessByAnyUser *bool `json:"allow_access_by_any_user,omitempty" path:"allow_access_by_any_user,omitempty" url:"allow_access_by_any_user,omitempty"`
	// Lock creator user ID.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	// Lock creator username.
	Username string `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	// Require this existing, unexpired token before refreshing or replacing a lock. Set token to the same value to refresh, or a different value to replace.
	ExpectedToken string `json:"expected_token,omitempty" path:"expected_token,omitempty" url:"expected_token,omitempty"`
}

// Identifier returns the resource path.
func (l Lock) Identifier() interface{} {
	return l.Path
}

// LockCollection is a list of Lock resources.
type LockCollection []Lock

// LockListForParams contains the request parameters for GET /locks/{path}.
type LockListForParams struct {
	// Path to operate on. Required; sent in the request path.
	Path string `url:"-,omitempty" json:"-,omitempty" path:"path"`
	// Include locks from children objects?
	IncludeChildren *bool `url:"include_children,omitempty" json:"include_children,omitempty" path:"include_children"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// LockCreateParams contains the request parameters for POST /locks/{path}.
type LockCreateParams struct {
	// Path. Required; sent in the request path.
	Path string `url:"-,omitempty" json:"-,omitempty" path:"path"`
	// Lock token. With expected_token, use the same value to refresh or a different value to replace the existing token.
	Token string `url:"token,omitempty" json:"token,omitempty" path:"token"`
	// Require this existing, unexpired token before refreshing or replacing a lock. Set token to the same value to refresh, or a different value to replace.
	ExpectedToken string `url:"expected_token,omitempty" json:"expected_token,omitempty" path:"expected_token"`
	// Can lock be modified by users other than its creator?
	AllowAccessByAnyUser *bool `url:"allow_access_by_any_user,omitempty" json:"allow_access_by_any_user,omitempty" path:"allow_access_by_any_user"`
	// Is lock exclusive?
	Exclusive *bool `url:"exclusive,omitempty" json:"exclusive,omitempty" path:"exclusive"`
	// Does lock apply to subfolders?
	Recursive *bool `url:"recursive,omitempty" json:"recursive,omitempty" path:"recursive"`
	// Arbitrary descriptive label for the lock. Does not change the lock creator or permissions.
	Owner string `url:"owner,omitempty" json:"owner,omitempty" path:"owner"`
	// Lock timeout in seconds.
	Timeout int64 `url:"timeout,omitempty" json:"timeout,omitempty" path:"timeout"`
}

// LockDeleteParams contains the request parameters for DELETE /locks/{path}.
type LockDeleteParams struct {
	// Path. Required; sent in the request path.
	Path string `url:"-,omitempty" json:"-,omitempty" path:"path"`
	// Lock token. Required.
	Token string `url:"token" json:"token" path:"token"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (l *Lock) UnmarshalJSON(data []byte) error {
	type lock Lock
	var v lock
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*l = Lock(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (l *LockCollection) UnmarshalJSON(data []byte) error {
	type locks LockCollection
	var v locks
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*l = LockCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (l *LockCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*l))
	for i, v := range *l {
		ret[i] = v
	}

	return &ret
}
