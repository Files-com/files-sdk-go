package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type ApiKey struct {
	Id               int64     `json:"id,omitempty"`
	DescriptiveLabel string    `json:"descriptive_label,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	ExpiresAt        time.Time `json:"expires_at,omitempty"`
	Key              string    `json:"key,omitempty"`
	LastUseAt        time.Time `json:"last_use_at,omitempty"`
	Name             string    `json:"name,omitempty"`
	Path             string    `json:"path,omitempty"`
	PermissionSet    string    `json:"permission_set,omitempty"`
	Platform         string    `json:"platform,omitempty"`
	UserId           int64     `json:"user_id,omitempty"`
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
	UserId     int64           `url:"user_id,omitempty" required:"false" json:"user_id,omitempty"`
	Cursor     string          `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage    int64           `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false" json:"filter_like,omitempty"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty"`
	lib.ListParams
}

type ApiKeyFindParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
}

type ApiKeyCreateParams struct {
	UserId        int64                   `url:"user_id,omitempty" required:"false" json:"user_id,omitempty"`
	Name          string                  `url:"name,omitempty" required:"false" json:"name,omitempty"`
	ExpiresAt     time.Time               `url:"expires_at,omitempty" required:"false" json:"expires_at,omitempty"`
	PermissionSet ApiKeyPermissionSetEnum `url:"permission_set,omitempty" required:"false" json:"permission_set,omitempty"`
	Path          string                  `url:"path,omitempty" required:"false" json:"path,omitempty"`
}

type ApiKeyUpdateCurrentParams struct {
	ExpiresAt     time.Time               `url:"expires_at,omitempty" required:"false" json:"expires_at,omitempty"`
	Name          string                  `url:"name,omitempty" required:"false" json:"name,omitempty"`
	PermissionSet ApiKeyPermissionSetEnum `url:"permission_set,omitempty" required:"false" json:"permission_set,omitempty"`
}

type ApiKeyUpdateParams struct {
	Id            int64                   `url:"-,omitempty" required:"true" json:"-,omitempty"`
	Name          string                  `url:"name,omitempty" required:"false" json:"name,omitempty"`
	ExpiresAt     time.Time               `url:"expires_at,omitempty" required:"false" json:"expires_at,omitempty"`
	PermissionSet ApiKeyPermissionSetEnum `url:"permission_set,omitempty" required:"false" json:"permission_set,omitempty"`
}

type ApiKeyDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
}

func (a *ApiKey) UnmarshalJSON(data []byte) error {
	type apiKey ApiKey
	var v apiKey
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = ApiKey(v)
	return nil
}

func (a *ApiKeyCollection) UnmarshalJSON(data []byte) error {
	type apiKeys []ApiKey
	var v apiKeys
	if err := json.Unmarshal(data, &v); err != nil {
		return err
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
