package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/lib"
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
	}
}

type ApiKeyListParams struct {
	UserId     int64           `url:"user_id,omitempty" required:"false"`
	Cursor     string          `url:"cursor,omitempty" required:"false"`
	PerPage    int64           `url:"per_page,omitempty" required:"false"`
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false"`
	lib.ListParams
}

type ApiKeyFindParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

type ApiKeyCreateParams struct {
	UserId        int64                   `url:"user_id,omitempty" required:"false"`
	Name          string                  `url:"name,omitempty" required:"false"`
	ExpiresAt     time.Time               `url:"expires_at,omitempty" required:"false"`
	PermissionSet ApiKeyPermissionSetEnum `url:"permission_set,omitempty" required:"false"`
	Path          string                  `url:"path,omitempty" required:"false"`
}

type ApiKeyUpdateCurrentParams struct {
	ExpiresAt     time.Time               `url:"expires_at,omitempty" required:"false"`
	Name          string                  `url:"name,omitempty" required:"false"`
	PermissionSet ApiKeyPermissionSetEnum `url:"permission_set,omitempty" required:"false"`
}

type ApiKeyUpdateParams struct {
	Id            int64                   `url:"-,omitempty" required:"true"`
	Name          string                  `url:"name,omitempty" required:"false"`
	ExpiresAt     time.Time               `url:"expires_at,omitempty" required:"false"`
	PermissionSet ApiKeyPermissionSetEnum `url:"permission_set,omitempty" required:"false"`
}

type ApiKeyDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
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
