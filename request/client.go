package request

import (
	"context"
	"strconv"

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

func (i *Iter) Request() files_sdk.Request {
	return i.Current().(files_sdk.Request)
}

func (c *Client) List(ctx context.Context, params files_sdk.RequestListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/requests"
	i.ListParams = &params
	list := files_sdk.RequestCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.RequestListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) GetFolder(ctx context.Context, params files_sdk.RequestGetFolderParams) (files_sdk.RequestCollection, error) {
	requestCollection := files_sdk.RequestCollection{}
	path := lib.BuildPath("/requests/folders/", params.Path)
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return requestCollection, err
	}
	if res.StatusCode == 204 {
		return requestCollection, nil
	}
	if err := requestCollection.UnmarshalJSON(*data); err != nil {
		return requestCollection, err
	}

	return requestCollection, nil
}

func GetFolder(ctx context.Context, params files_sdk.RequestGetFolderParams) (files_sdk.RequestCollection, error) {
	return (&Client{}).GetFolder(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.RequestCreateParams) (files_sdk.Request, error) {
	request := files_sdk.Request{}
	path := "/requests"
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return request, err
	}
	if res.StatusCode == 204 {
		return request, nil
	}
	if err := request.UnmarshalJSON(*data); err != nil {
		return request, err
	}

	return request, nil
}

func Create(ctx context.Context, params files_sdk.RequestCreateParams) (files_sdk.Request, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.RequestDeleteParams) (files_sdk.Request, error) {
	request := files_sdk.Request{}
	if params.Id == 0 {
		return request, lib.CreateError(params, "Id")
	}
	path := "/requests/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return request, err
	}
	if res.StatusCode == 204 {
		return request, nil
	}
	if err := request.UnmarshalJSON(*data); err != nil {
		return request, err
	}

	return request, nil
}

func Delete(ctx context.Context, params files_sdk.RequestDeleteParams) (files_sdk.Request, error) {
	return (&Client{}).Delete(ctx, params)
}
