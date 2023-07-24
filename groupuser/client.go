package group_user

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
	listquery "github.com/Files-com/files-sdk-go/v2/listquery"
)

type Client struct {
	files_sdk.Config
}

type Iter struct {
	*files_sdk.Iter
	*Client
}

func (i *Iter) Reload(opts ...files_sdk.RequestResponseOption) files_sdk.IterI {
	return &Iter{Iter: i.Iter.Reload(opts...).(*files_sdk.Iter), Client: i.Client}
}

func (i *Iter) GroupUser() files_sdk.GroupUser {
	return i.Current().(files_sdk.GroupUser)
}

func (c *Client) List(params files_sdk.GroupUserListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/group_users", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.GroupUserCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.GroupUserListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) Create(params files_sdk.GroupUserCreateParams, opts ...files_sdk.RequestResponseOption) (groupUser files_sdk.GroupUser, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/group_users", Params: params, Entity: &groupUser}, opts...)
	return
}

func Create(params files_sdk.GroupUserCreateParams, opts ...files_sdk.RequestResponseOption) (groupUser files_sdk.GroupUser, err error) {
	return (&Client{}).Create(params, opts...)
}

func (c *Client) Update(params files_sdk.GroupUserUpdateParams, opts ...files_sdk.RequestResponseOption) (groupUser files_sdk.GroupUser, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/group_users/{id}", Params: params, Entity: &groupUser}, opts...)
	return
}

func Update(params files_sdk.GroupUserUpdateParams, opts ...files_sdk.RequestResponseOption) (groupUser files_sdk.GroupUser, err error) {
	return (&Client{}).Update(params, opts...)
}

func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (groupUser files_sdk.GroupUser, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/group_users/{id}", Params: params, Entity: &groupUser}, opts...)
	return
}

func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (groupUser files_sdk.GroupUser, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}

func (c *Client) Delete(params files_sdk.GroupUserDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/group_users/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(params files_sdk.GroupUserDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
