package public_key

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

func (i *Iter) PublicKey() files_sdk.PublicKey {
	return i.Current().(files_sdk.PublicKey)
}

func (c *Client) List(ctx context.Context, params files_sdk.PublicKeyListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/public_keys", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.PublicKeyCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.PublicKeyListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.PublicKeyFindParams) (publicKey files_sdk.PublicKey, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/public_keys/{id}", Params: params, Entity: &publicKey})
	return
}

func Find(ctx context.Context, params files_sdk.PublicKeyFindParams) (publicKey files_sdk.PublicKey, err error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.PublicKeyCreateParams) (publicKey files_sdk.PublicKey, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/public_keys", Params: params, Entity: &publicKey})
	return
}

func Create(ctx context.Context, params files_sdk.PublicKeyCreateParams) (publicKey files_sdk.PublicKey, err error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.PublicKeyUpdateParams) (publicKey files_sdk.PublicKey, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/public_keys/{id}", Params: params, Entity: &publicKey})
	return
}

func Update(ctx context.Context, params files_sdk.PublicKeyUpdateParams) (publicKey files_sdk.PublicKey, err error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.PublicKeyDeleteParams) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/public_keys/{id}", Params: params, Entity: nil})
	return
}

func Delete(ctx context.Context, params files_sdk.PublicKeyDeleteParams) (err error) {
	return (&Client{}).Delete(ctx, params)
}
