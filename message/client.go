package message

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

func (i *Iter) Message() files_sdk.Message {
	return i.Current().(files_sdk.Message)
}

func (c *Client) List(ctx context.Context, params files_sdk.MessageListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/messages", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.MessageCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list, opts...)
	return i, nil
}

func List(ctx context.Context, params files_sdk.MessageListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(ctx, params, opts...)
}

func (c *Client) Find(ctx context.Context, params files_sdk.MessageFindParams, opts ...files_sdk.RequestResponseOption) (message files_sdk.Message, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/messages/{id}", Params: params, Entity: &message}, opts...)
	return
}

func Find(ctx context.Context, params files_sdk.MessageFindParams, opts ...files_sdk.RequestResponseOption) (message files_sdk.Message, err error) {
	return (&Client{}).Find(ctx, params, opts...)
}

func (c *Client) Create(ctx context.Context, params files_sdk.MessageCreateParams, opts ...files_sdk.RequestResponseOption) (message files_sdk.Message, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/messages", Params: params, Entity: &message}, opts...)
	return
}

func Create(ctx context.Context, params files_sdk.MessageCreateParams, opts ...files_sdk.RequestResponseOption) (message files_sdk.Message, err error) {
	return (&Client{}).Create(ctx, params, opts...)
}

func (c *Client) Update(ctx context.Context, params files_sdk.MessageUpdateParams, opts ...files_sdk.RequestResponseOption) (message files_sdk.Message, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/messages/{id}", Params: params, Entity: &message}, opts...)
	return
}

func Update(ctx context.Context, params files_sdk.MessageUpdateParams, opts ...files_sdk.RequestResponseOption) (message files_sdk.Message, err error) {
	return (&Client{}).Update(ctx, params, opts...)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.MessageDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/messages/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(ctx context.Context, params files_sdk.MessageDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(ctx, params, opts...)
}
