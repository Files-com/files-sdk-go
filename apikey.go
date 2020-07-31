package files_sdk

import (
  lib "github.com/Files-com/files-sdk-go/lib"
  "encoding/json"
  "time"
)

type ApiKey struct {
  Id int `json:"id,omitempty"`
  DescriptiveLabel string `json:"descriptive_label,omitempty"`
  CreatedAt time.Time `json:"created_at,omitempty"`
  ExpiresAt time.Time `json:"expires_at,omitempty"`
  Key string `json:"key,omitempty"`
  LastUseAt time.Time `json:"last_use_at,omitempty"`
  Name string `json:"name,omitempty"`
  Path string `json:"path,omitempty"`
  PermissionSet string `json:"permission_set,omitempty"`
  Platform string `json:"platform,omitempty"`
  UserId int `json:"user_id,omitempty"`
}

type ApiKeyCollection []ApiKey

type ApiKeyListParams struct {
  UserId int `url:"user_id,omitempty"`
  Page int `url:"page,omitempty"`
  PerPage int `url:"per_page,omitempty"`
  Action string `url:"action,omitempty"`
  Cursor string `url:"cursor,omitempty"`
  SortBy json.RawMessage `url:"sort_by,omitempty"`
  Filter json.RawMessage `url:"filter,omitempty"`
  FilterGt json.RawMessage `url:"filter_gt,omitempty"`
  FilterGteq json.RawMessage `url:"filter_gteq,omitempty"`
  FilterLike json.RawMessage `url:"filter_like,omitempty"`
  FilterLt json.RawMessage `url:"filter_lt,omitempty"`
  FilterLteq json.RawMessage `url:"filter_lteq,omitempty"`
  lib.ListParams
}

type ApiKeyFindParams struct {
  Id int `url:"-,omitempty"`
}

type ApiKeyCreateParams struct {
  UserId int `url:"user_id,omitempty"`
  Name string `url:"name,omitempty"`
  ExpiresAt string `url:"expires_at,omitempty"`
  PermissionSet string `url:"permission_set,omitempty"`
  Path string `url:"path,omitempty"`
}

type ApiKeyUpdateParams struct {
  Id int `url:"-,omitempty"`
  Name string `url:"name,omitempty"`
  ExpiresAt string `url:"expires_at,omitempty"`
  PermissionSet string `url:"permission_set,omitempty"`
}

type ApiKeyUpdateCurrentParams struct {
  ExpiresAt string `url:"expires_at,omitempty"`
  Name string `url:"name,omitempty"`
  PermissionSet string `url:"permission_set,omitempty"`
}

type ApiKeyDeleteParams struct {
  Id int `url:"-,omitempty"`
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

