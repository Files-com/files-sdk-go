package behavior

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

func (i *Iter) Behavior() files_sdk.Behavior {
	return i.Current().(files_sdk.Behavior)
}

func (c *Client) List(ctx context.Context, params files_sdk.BehaviorListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/behaviors"
	i.ListParams = &params
	list := files_sdk.BehaviorCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.BehaviorListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.BehaviorFindParams) (files_sdk.Behavior, error) {
	behavior := files_sdk.Behavior{}
	if params.Id == 0 {
		return behavior, lib.CreateError(params, "Id")
	}
	path := "/behaviors/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return behavior, err
	}
	if res.StatusCode == 204 {
		return behavior, nil
	}
	if err := behavior.UnmarshalJSON(*data); err != nil {
		return behavior, err
	}

	return behavior, nil
}

func Find(ctx context.Context, params files_sdk.BehaviorFindParams) (files_sdk.Behavior, error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) ListFor(ctx context.Context, params files_sdk.BehaviorListForParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := lib.BuildPath("/behaviors/folders/", params.Path)
	i.ListParams = &params
	list := files_sdk.BehaviorCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func ListFor(ctx context.Context, params files_sdk.BehaviorListForParams) (*Iter, error) {
	return (&Client{}).ListFor(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.BehaviorCreateParams) (files_sdk.Behavior, error) {
	behavior := files_sdk.Behavior{}
	path := "/behaviors"
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return behavior, err
	}
	if res.StatusCode == 204 {
		return behavior, nil
	}
	if err := behavior.UnmarshalJSON(*data); err != nil {
		return behavior, err
	}

	return behavior, nil
}

func Create(ctx context.Context, params files_sdk.BehaviorCreateParams) (files_sdk.Behavior, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) WebhookTest(ctx context.Context, params files_sdk.BehaviorWebhookTestParams) (files_sdk.Behavior, error) {
	behavior := files_sdk.Behavior{}
	path := "/behaviors/webhook/test"
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return behavior, err
	}
	if res.StatusCode == 204 {
		return behavior, nil
	}
	if err := behavior.UnmarshalJSON(*data); err != nil {
		return behavior, err
	}

	return behavior, nil
}

func WebhookTest(ctx context.Context, params files_sdk.BehaviorWebhookTestParams) (files_sdk.Behavior, error) {
	return (&Client{}).WebhookTest(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.BehaviorUpdateParams) (files_sdk.Behavior, error) {
	behavior := files_sdk.Behavior{}
	if params.Id == 0 {
		return behavior, lib.CreateError(params, "Id")
	}
	path := "/behaviors/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return behavior, err
	}
	if res.StatusCode == 204 {
		return behavior, nil
	}
	if err := behavior.UnmarshalJSON(*data); err != nil {
		return behavior, err
	}

	return behavior, nil
}

func Update(ctx context.Context, params files_sdk.BehaviorUpdateParams) (files_sdk.Behavior, error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.BehaviorDeleteParams) (files_sdk.Behavior, error) {
	behavior := files_sdk.Behavior{}
	if params.Id == 0 {
		return behavior, lib.CreateError(params, "Id")
	}
	path := "/behaviors/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return behavior, err
	}
	if res.StatusCode == 204 {
		return behavior, nil
	}
	if err := behavior.UnmarshalJSON(*data); err != nil {
		return behavior, err
	}

	return behavior, nil
}

func Delete(ctx context.Context, params files_sdk.BehaviorDeleteParams) (files_sdk.Behavior, error) {
	return (&Client{}).Delete(ctx, params)
}
