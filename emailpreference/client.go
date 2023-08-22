package email_preference

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Get(params files_sdk.EmailPreferenceGetParams, opts ...files_sdk.RequestResponseOption) (emailPreference files_sdk.EmailPreference, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/email_preferences/{token}", Params: params, Entity: &emailPreference}, opts...)
	return
}

func Get(params files_sdk.EmailPreferenceGetParams, opts ...files_sdk.RequestResponseOption) (emailPreference files_sdk.EmailPreference, err error) {
	return (&Client{}).Get(params, opts...)
}

func (c *Client) Update(params files_sdk.EmailPreferenceUpdateParams, opts ...files_sdk.RequestResponseOption) (emailPreference files_sdk.EmailPreference, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/email_preferences/{token}", Params: params, Entity: &emailPreference}, opts...)
	return
}

func Update(params files_sdk.EmailPreferenceUpdateParams, opts ...files_sdk.RequestResponseOption) (emailPreference files_sdk.EmailPreference, err error) {
	return (&Client{}).Update(params, opts...)
}

func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (emailPreference files_sdk.EmailPreference, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/email_preferences/{token}", Params: params, Entity: &emailPreference}, opts...)
	return
}

func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (emailPreference files_sdk.EmailPreference, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}
