package automation

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

func (i *Iter) Automation() files_sdk.Automation {
	return i.Current().(files_sdk.Automation)
}

func (c *Client) List(ctx context.Context, params files_sdk.AutomationListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/automations"
	i.ListParams = &params
	list := files_sdk.AutomationCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.AutomationListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.AutomationFindParams) (files_sdk.Automation, error) {
	automation := files_sdk.Automation{}
	if params.Id == 0 {
		return automation, lib.CreateError(params, "Id")
	}
	path := "/automations/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return automation, err
	}
	if res.StatusCode == 204 {
		return automation, nil
	}

	return automation, automation.UnmarshalJSON(*data)
}

func Find(ctx context.Context, params files_sdk.AutomationFindParams) (files_sdk.Automation, error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.AutomationCreateParams) (files_sdk.Automation, error) {
	automation := files_sdk.Automation{}
	path := "/automations"
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return automation, err
	}
	if res.StatusCode == 204 {
		return automation, nil
	}

	return automation, automation.UnmarshalJSON(*data)
}

func Create(ctx context.Context, params files_sdk.AutomationCreateParams) (files_sdk.Automation, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.AutomationUpdateParams) (files_sdk.Automation, error) {
	automation := files_sdk.Automation{}
	if params.Id == 0 {
		return automation, lib.CreateError(params, "Id")
	}
	path := "/automations/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return automation, err
	}
	if res.StatusCode == 204 {
		return automation, nil
	}

	return automation, automation.UnmarshalJSON(*data)
}

func Update(ctx context.Context, params files_sdk.AutomationUpdateParams) (files_sdk.Automation, error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.AutomationDeleteParams) error {
	automation := files_sdk.Automation{}
	if params.Id == 0 {
		return lib.CreateError(params, "Id")
	}
	path := "/automations/" + strconv.FormatInt(params.Id, 10) + ""
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

	return automation.UnmarshalJSON(*data)
}

func Delete(ctx context.Context, params files_sdk.AutomationDeleteParams) error {
	return (&Client{}).Delete(ctx, params)
}
