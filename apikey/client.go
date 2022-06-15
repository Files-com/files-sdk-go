package api_key

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

func (i *Iter) ApiKey() files_sdk.ApiKey {
	return i.Current().(files_sdk.ApiKey)
}

func (c *Client) List(ctx context.Context, params files_sdk.ApiKeyListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/api_keys"
	i.ListParams = &params
	list := files_sdk.ApiKeyCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.ApiKeyListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) FindCurrent(ctx context.Context) (files_sdk.ApiKey, error) {
	apiKey := files_sdk.ApiKey{}
	path := "/api_key"
	exportedParams := lib.Params{Params: lib.Interface()}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return apiKey, err
	}
	if res.StatusCode == 204 {
		return apiKey, nil
	}

	return apiKey, apiKey.UnmarshalJSON(*data)
}

func FindCurrent(ctx context.Context) (files_sdk.ApiKey, error) {
	return (&Client{}).FindCurrent(ctx)
}

func (c *Client) Find(ctx context.Context, params files_sdk.ApiKeyFindParams) (files_sdk.ApiKey, error) {
	apiKey := files_sdk.ApiKey{}
	if params.Id == 0 {
		return apiKey, lib.CreateError(params, "Id")
	}
	path := "/api_keys/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return apiKey, err
	}
	if res.StatusCode == 204 {
		return apiKey, nil
	}

	return apiKey, apiKey.UnmarshalJSON(*data)
}

func Find(ctx context.Context, params files_sdk.ApiKeyFindParams) (files_sdk.ApiKey, error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.ApiKeyCreateParams) (files_sdk.ApiKey, error) {
	apiKey := files_sdk.ApiKey{}
	path := "/api_keys"
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return apiKey, err
	}
	if res.StatusCode == 204 {
		return apiKey, nil
	}

	return apiKey, apiKey.UnmarshalJSON(*data)
}

func Create(ctx context.Context, params files_sdk.ApiKeyCreateParams) (files_sdk.ApiKey, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) UpdateCurrent(ctx context.Context, params files_sdk.ApiKeyUpdateCurrentParams) (files_sdk.ApiKey, error) {
	apiKey := files_sdk.ApiKey{}
	path := "/api_key"
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return apiKey, err
	}
	if res.StatusCode == 204 {
		return apiKey, nil
	}

	return apiKey, apiKey.UnmarshalJSON(*data)
}

func UpdateCurrent(ctx context.Context, params files_sdk.ApiKeyUpdateCurrentParams) (files_sdk.ApiKey, error) {
	return (&Client{}).UpdateCurrent(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.ApiKeyUpdateParams) (files_sdk.ApiKey, error) {
	apiKey := files_sdk.ApiKey{}
	if params.Id == 0 {
		return apiKey, lib.CreateError(params, "Id")
	}
	path := "/api_keys/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return apiKey, err
	}
	if res.StatusCode == 204 {
		return apiKey, nil
	}

	return apiKey, apiKey.UnmarshalJSON(*data)
}

func Update(ctx context.Context, params files_sdk.ApiKeyUpdateParams) (files_sdk.ApiKey, error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) DeleteCurrent(ctx context.Context) error {
	apiKey := files_sdk.ApiKey{}
	path := "/api_key"
	exportedParams := lib.Params{Params: lib.Interface()}
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

	return apiKey.UnmarshalJSON(*data)
}

func DeleteCurrent(ctx context.Context) error {
	return (&Client{}).DeleteCurrent(ctx)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.ApiKeyDeleteParams) error {
	apiKey := files_sdk.ApiKey{}
	if params.Id == 0 {
		return lib.CreateError(params, "Id")
	}
	path := "/api_keys/" + strconv.FormatInt(params.Id, 10) + ""
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

	return apiKey.UnmarshalJSON(*data)
}

func Delete(ctx context.Context, params files_sdk.ApiKeyDeleteParams) error {
	return (&Client{}).Delete(ctx, params)
}
