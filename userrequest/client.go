package user_request

import (
	"context"

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

func (i *Iter) UserRequest() files_sdk.UserRequest {
	return i.Current().(files_sdk.UserRequest)
}

func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.UserRequestFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(context.Background(), params, opts...)
}

func (c *Client) List(ctx context.Context, params files_sdk.UserRequestListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/user_requests", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.UserRequestCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list, opts...)
	return i, nil
}

func List(ctx context.Context, params files_sdk.UserRequestListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(ctx, params, opts...)
}

func (c *Client) Find(ctx context.Context, params files_sdk.UserRequestFindParams, opts ...files_sdk.RequestResponseOption) (userRequest files_sdk.UserRequest, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/user_requests/{id}", Params: params, Entity: &userRequest}, opts...)
	return
}

func Find(ctx context.Context, params files_sdk.UserRequestFindParams, opts ...files_sdk.RequestResponseOption) (userRequest files_sdk.UserRequest, err error) {
	return (&Client{}).Find(ctx, params, opts...)
}

func (c *Client) Create(ctx context.Context, params files_sdk.UserRequestCreateParams, opts ...files_sdk.RequestResponseOption) (userRequest files_sdk.UserRequest, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/user_requests", Params: params, Entity: &userRequest}, opts...)
	return
}

func Create(ctx context.Context, params files_sdk.UserRequestCreateParams, opts ...files_sdk.RequestResponseOption) (userRequest files_sdk.UserRequest, err error) {
	return (&Client{}).Create(ctx, params, opts...)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.UserRequestDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/user_requests/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(ctx context.Context, params files_sdk.UserRequestDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(ctx, params, opts...)
}
