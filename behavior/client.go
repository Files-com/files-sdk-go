package behavior

import (
	"strconv"

	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
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

func (c *Client) List(params files_sdk.BehaviorListParams) *Iter {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/behaviors"

	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, i.ExportParams())
		defaultValue := make([]interface{}, 0)
		if err != nil {
			return &defaultValue, "", err
		}
		list := files_sdk.BehaviorCollection{}
		if err := list.UnmarshalJSON(*data); err != nil {
			return &defaultValue, "", err
		}

		ret := make([]interface{}, len(list))
		for i, v := range list {
			ret[i] = v
		}
		cursor := res.Header.Get("X-Files-Cursor")
		return &ret, cursor, nil
	}
	i.ListParams = &params
	return i
}

func List(params files_sdk.BehaviorListParams) *Iter {
	return (&Client{}).List(params)
}

func (c *Client) Find(params files_sdk.BehaviorFindParams) (files_sdk.Behavior, error) {
	behavior := files_sdk.Behavior{}
	path := "/behaviors/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	data, res, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
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

func (c *Client) ListFor(params files_sdk.BehaviorListForParams) *Iter {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/behaviors/folders/" + lib.QueryEscape(params.Path) + ""

	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, i.ExportParams())
		defaultValue := make([]interface{}, 0)
		if err != nil {
			return &defaultValue, "", err
		}
		list := files_sdk.BehaviorCollection{}
		if err := list.UnmarshalJSON(*data); err != nil {
			return &defaultValue, "", err
		}

		ret := make([]interface{}, len(list))
		for i, v := range list {
			ret[i] = v
		}
		cursor := res.Header.Get("X-Files-Cursor")
		return &ret, cursor, nil
	}
	i.ListParams = &params
	return i
}

func ListFor(params files_sdk.BehaviorListForParams) *Iter {
	return (&Client{}).ListFor(params)
}

func (c *Client) Create(params files_sdk.BehaviorCreateParams) (files_sdk.Behavior, error) {
	behavior := files_sdk.Behavior{}
	path := "/behaviors"
	data, res, err := files_sdk.Call("POST", c.Config, path, lib.ExportParams(params))
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
	data, res, err := files_sdk.Call("POST", c.Config, path, lib.ExportParams(params))
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
	path := "/behaviors/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	data, res, err := files_sdk.Call("PATCH", c.Config, path, lib.ExportParams(params))
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
	path := "/behaviors/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	data, res, err := files_sdk.Call("DELETE", c.Config, path, lib.ExportParams(params))
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
