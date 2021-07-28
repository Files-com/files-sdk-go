package group

import (
	"context"
	"strconv"

	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
	listquery "github.com/Files-com/files-sdk-go/listquery"
)

type Client struct {
	files_sdk.Config
}

type Iter struct {
	*lib.Iter
}

func (i *Iter) Group() files_sdk.Group {
	return i.Current().(files_sdk.Group)
}

func (c *Client) List(ctx context.Context, params files_sdk.GroupListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/groups"
	i.ListParams = &params
	list := files_sdk.GroupCollection{}
	i.Query = listquery.Build(ctx, i, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.GroupListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.GroupFindParams) (files_sdk.Group, error) {
	group := files_sdk.Group{}
	if params.Id == 0 {
		return group, lib.CreateError(params, "Id")
	}
	path := "/groups/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return group, err
	}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return group, err
	}
	if res.StatusCode == 204 {
		return group, nil
	}
	if err := group.UnmarshalJSON(*data); err != nil {
		return group, err
	}

	return group, nil
}

func Find(ctx context.Context, params files_sdk.GroupFindParams) (files_sdk.Group, error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.GroupCreateParams) (files_sdk.Group, error) {
	group := files_sdk.Group{}
	path := "/groups"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return group, err
	}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return group, err
	}
	if res.StatusCode == 204 {
		return group, nil
	}
	if err := group.UnmarshalJSON(*data); err != nil {
		return group, err
	}

	return group, nil
}

func Create(ctx context.Context, params files_sdk.GroupCreateParams) (files_sdk.Group, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.GroupUpdateParams) (files_sdk.Group, error) {
	group := files_sdk.Group{}
	if params.Id == 0 {
		return group, lib.CreateError(params, "Id")
	}
	path := "/groups/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return group, err
	}
	data, res, err := files_sdk.Call(ctx, "PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return group, err
	}
	if res.StatusCode == 204 {
		return group, nil
	}
	if err := group.UnmarshalJSON(*data); err != nil {
		return group, err
	}

	return group, nil
}

func Update(ctx context.Context, params files_sdk.GroupUpdateParams) (files_sdk.Group, error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.GroupDeleteParams) (files_sdk.Group, error) {
	group := files_sdk.Group{}
	if params.Id == 0 {
		return group, lib.CreateError(params, "Id")
	}
	path := "/groups/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return group, err
	}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return group, err
	}
	if res.StatusCode == 204 {
		return group, nil
	}
	if err := group.UnmarshalJSON(*data); err != nil {
		return group, err
	}

	return group, nil
}

func Delete(ctx context.Context, params files_sdk.GroupDeleteParams) (files_sdk.Group, error) {
	return (&Client{}).Delete(ctx, params)
}
