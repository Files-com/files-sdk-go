package lock

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

func (i *Iter) Lock() files_sdk.Lock {
	return i.Current().(files_sdk.Lock)
}

func (c *Client) ListFor(ctx context.Context, params files_sdk.LockListForParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := lib.BuildPath("/locks/", params.Path)
	i.ListParams = &params
	list := files_sdk.LockCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func ListFor(ctx context.Context, params files_sdk.LockListForParams) (*Iter, error) {
	return (&Client{}).ListFor(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.LockCreateParams) (files_sdk.Lock, error) {
	lock := files_sdk.Lock{}
	path := lib.BuildPath("/locks/", params.Path)
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return lock, err
	}
	if res.StatusCode == 204 {
		return lock, nil
	}

	return lock, lock.UnmarshalJSON(*data)
}

func Create(ctx context.Context, params files_sdk.LockCreateParams) (files_sdk.Lock, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.LockDeleteParams) error {
	lock := files_sdk.Lock{}
	path := lib.BuildPath("/locks/", params.Path)
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return err
	}
	if res.StatusCode == 204 {
		return nil
	}

	return lock.UnmarshalJSON(*data)
}

func Delete(ctx context.Context, params files_sdk.LockDeleteParams) error {
	return (&Client{}).Delete(ctx, params)
}
