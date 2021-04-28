package behavior

import (
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

func (i *Iter) Behavior() files_sdk.Behavior {
	return i.Current().(files_sdk.Behavior)
}

func (c *Client) List(params files_sdk.BehaviorListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/behaviors"
	i.ListParams = &params
	list := files_sdk.BehaviorCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func List(params files_sdk.BehaviorListParams) (*Iter, error) {
	return (&Client{}).List(params)
}

func (c *Client) Find(params files_sdk.BehaviorFindParams) (files_sdk.Behavior, error) {
	behavior := files_sdk.Behavior{}
	if params.Id == 0 {
		return behavior, lib.CreateError(params, "Id")
	}
	path := "/behaviors/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return behavior, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
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

func Find(params files_sdk.BehaviorFindParams) (files_sdk.Behavior, error) {
	return (&Client{}).Find(params)
}

func (c *Client) ListFor(params files_sdk.BehaviorListForParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := lib.BuildPath("/behaviors/folders/", params.Path)
	i.ListParams = &params
	list := files_sdk.BehaviorCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func ListFor(params files_sdk.BehaviorListForParams) (*Iter, error) {
	return (&Client{}).ListFor(params)
}

func (c *Client) Create(params files_sdk.BehaviorCreateParams) (files_sdk.Behavior, error) {
	behavior := files_sdk.Behavior{}
	path := "/behaviors"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return behavior, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
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

func Create(params files_sdk.BehaviorCreateParams) (files_sdk.Behavior, error) {
	return (&Client{}).Create(params)
}

func (c *Client) WebhookTest(params files_sdk.BehaviorWebhookTestParams) (files_sdk.Behavior, error) {
	behavior := files_sdk.Behavior{}
	path := "/behaviors/webhook/test"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return behavior, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
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

func WebhookTest(params files_sdk.BehaviorWebhookTestParams) (files_sdk.Behavior, error) {
	return (&Client{}).WebhookTest(params)
}

func (c *Client) Update(params files_sdk.BehaviorUpdateParams) (files_sdk.Behavior, error) {
	behavior := files_sdk.Behavior{}
	if params.Id == 0 {
		return behavior, lib.CreateError(params, "Id")
	}
	path := "/behaviors/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return behavior, err
	}
	data, res, err := files_sdk.Call("PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
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

func Update(params files_sdk.BehaviorUpdateParams) (files_sdk.Behavior, error) {
	return (&Client{}).Update(params)
}

func (c *Client) Delete(params files_sdk.BehaviorDeleteParams) (files_sdk.Behavior, error) {
	behavior := files_sdk.Behavior{}
	if params.Id == 0 {
		return behavior, lib.CreateError(params, "Id")
	}
	path := "/behaviors/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return behavior, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
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

func Delete(params files_sdk.BehaviorDeleteParams) (files_sdk.Behavior, error) {
	return (&Client{}).Delete(params)
}
