package automation

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

func (i *Iter) Automation() files_sdk.Automation {
	return i.Current().(files_sdk.Automation)
}

func (c *Client) List(ctx context.Context, params files_sdk.AutomationListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/automations"
	i.ListParams = &params
	list := files_sdk.AutomationCollection{}
	i.Query = listquery.Build(ctx, i, c.Config, path, &list)
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
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return automation, err
	}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return automation, err
	}
	if res.StatusCode == 204 {
		return automation, nil
	}
	if err := automation.UnmarshalJSON(*data); err != nil {
		return automation, err
	}

	return automation, nil
}

func Find(ctx context.Context, params files_sdk.AutomationFindParams) (files_sdk.Automation, error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.AutomationCreateParams) (files_sdk.Automation, error) {
	automation := files_sdk.Automation{}
	path := "/automations"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return automation, err
	}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return automation, err
	}
	if res.StatusCode == 204 {
		return automation, nil
	}
	if err := automation.UnmarshalJSON(*data); err != nil {
		return automation, err
	}

	return automation, nil
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
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return automation, err
	}
	data, res, err := files_sdk.Call(ctx, "PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return automation, err
	}
	if res.StatusCode == 204 {
		return automation, nil
	}
	if err := automation.UnmarshalJSON(*data); err != nil {
		return automation, err
	}

	return automation, nil
}

func Update(ctx context.Context, params files_sdk.AutomationUpdateParams) (files_sdk.Automation, error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.AutomationDeleteParams) (files_sdk.Automation, error) {
	automation := files_sdk.Automation{}
	if params.Id == 0 {
		return automation, lib.CreateError(params, "Id")
	}
	path := "/automations/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return automation, err
	}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return automation, err
	}
	if res.StatusCode == 204 {
		return automation, nil
	}
	if err := automation.UnmarshalJSON(*data); err != nil {
		return automation, err
	}

	return automation, nil
}

func Delete(ctx context.Context, params files_sdk.AutomationDeleteParams) (files_sdk.Automation, error) {
	return (&Client{}).Delete(ctx, params)
}
