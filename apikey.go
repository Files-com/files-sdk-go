package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type ApiKey struct {
	Id               int64      `json:"id,omitempty" path:"id"`
	DescriptiveLabel string     `json:"descriptive_label,omitempty" path:"descriptive_label"`
	Description      string     `json:"description,omitempty" path:"description"`
	CreatedAt        *time.Time `json:"created_at,omitempty" path:"created_at"`
	ExpiresAt        *time.Time `json:"expires_at,omitempty" path:"expires_at"`
	Key              string     `json:"key,omitempty" path:"key"`
	LastUseAt        *time.Time `json:"last_use_at,omitempty" path:"last_use_at"`
	Name             string     `json:"name,omitempty" path:"name"`
	Path             string     `json:"path,omitempty" path:"path"`
	PermissionSet    string     `json:"permission_set,omitempty" path:"permission_set"`
	Platform         string     `json:"platform,omitempty" path:"platform"`
	Url              string     `json:"url,omitempty" path:"url"`
	UserId           int64      `json:"user_id,omitempty" path:"user_id"`
}

type ApiKeyCollection []ApiKey

type ApiKeyPermissionSetEnum string

func (u ApiKeyPermissionSetEnum) String() string {
	return string(u)
}

func (u ApiKeyPermissionSetEnum) Enum() map[string]ApiKeyPermissionSetEnum {
	return map[string]ApiKeyPermissionSetEnum{
		"none":               ApiKeyPermissionSetEnum("none"),
		"full":               ApiKeyPermissionSetEnum("full"),
		"desktop_app":        ApiKeyPermissionSetEnum("desktop_app"),
		"sync_app":           ApiKeyPermissionSetEnum("sync_app"),
		"office_integration": ApiKeyPermissionSetEnum("office_integration"),
		"mobile_app":         ApiKeyPermissionSetEnum("mobile_app"),
	}
}

type ApiKeyListParams struct {
	UserId     int64           `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty" path:"filter_lteq"`
	lib.ListParams
}

type ApiKeyFindParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
}

type ApiKeyCreateParams struct {
	UserId        int64                   `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	Name          string                  `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	Description   string                  `url:"description,omitempty" required:"false" json:"description,omitempty" path:"description"`
	ExpiresAt     *time.Time              `url:"expires_at,omitempty" required:"false" json:"expires_at,omitempty" path:"expires_at"`
	PermissionSet ApiKeyPermissionSetEnum `url:"permission_set,omitempty" required:"false" json:"permission_set,omitempty" path:"permission_set"`
	Path          string                  `url:"path,omitempty" required:"false" json:"path,omitempty" path:"path"`
}

type ApiKeyUpdateCurrentParams struct {
	ExpiresAt     *time.Time              `url:"expires_at,omitempty" required:"false" json:"expires_at,omitempty" path:"expires_at"`
	Name          string                  `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	PermissionSet ApiKeyPermissionSetEnum `url:"permission_set,omitempty" required:"false" json:"permission_set,omitempty" path:"permission_set"`
}

type ApiKeyUpdateParams struct {
	Id            int64                   `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
	Name          string                  `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	Description   string                  `url:"description,omitempty" required:"false" json:"description,omitempty" path:"description"`
	ExpiresAt     *time.Time              `url:"expires_at,omitempty" required:"false" json:"expires_at,omitempty" path:"expires_at"`
	PermissionSet ApiKeyPermissionSetEnum `url:"permission_set,omitempty" required:"false" json:"permission_set,omitempty" path:"permission_set"`
}

type ApiKeyDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
}

func (a *ApiKey) UnmarshalJSON(data []byte) error {
	type apiKey ApiKey
	var v apiKey
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = ApiKey(v)
	return nil
}

func (a *ApiKeyCollection) UnmarshalJSON(data []byte) error {
	type apiKeys ApiKeyCollection
	var v apiKeys
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = ApiKeyCollection(v)
	return nil
}

func (a *ApiKeyCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
