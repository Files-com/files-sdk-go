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
	*lib.Iter
}

func (i *Iter) UserRequest() files_sdk.UserRequest {
	return i.Current().(files_sdk.UserRequest)
}

func (c *Client) List(ctx context.Context, params files_sdk.UserRequestListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/user_requests", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.UserRequestCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.UserRequestListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.UserRequestFindParams) (userRequest files_sdk.UserRequest, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/user_requests/{id}", Params: params, Entity: &userRequest})
	return
}

func Find(ctx context.Context, params files_sdk.UserRequestFindParams) (userRequest files_sdk.UserRequest, err error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.UserRequestCreateParams) (userRequest files_sdk.UserRequest, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/user_requests", Params: params, Entity: &userRequest})
	return
}

func Create(ctx context.Context, params files_sdk.UserRequestCreateParams) (userRequest files_sdk.UserRequest, err error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.UserRequestDeleteParams) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/user_requests/{id}", Params: params, Entity: nil})
	return
}

func Delete(ctx context.Context, params files_sdk.UserRequestDeleteParams) (err error) {
	return (&Client{}).Delete(ctx, params)
}
