package automation_run

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

func (i *Iter) AutomationRun() files_sdk.AutomationRun {
	return i.Current().(files_sdk.AutomationRun)
}

func (c *Client) List(ctx context.Context, params files_sdk.AutomationRunListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/automation_runs"
	i.ListParams = &params
	list := files_sdk.AutomationRunCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.AutomationRunListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.AutomationRunFindParams) (files_sdk.AutomationRun, error) {
	automationRun := files_sdk.AutomationRun{}
	if params.Id == 0 {
		return automationRun, lib.CreateError(params, "Id")
	}
	path := "/automation_runs/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return automationRun, err
	}
	if res.StatusCode == 204 {
		return automationRun, nil
	}
	if err := automationRun.UnmarshalJSON(*data); err != nil {
		return automationRun, err
	}

	return automationRun, nil
}

func Find(ctx context.Context, params files_sdk.AutomationRunFindParams) (files_sdk.AutomationRun, error) {
	return (&Client{}).Find(ctx, params)
}
