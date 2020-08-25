package clickwrap

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

func (i *Iter) Clickwrap() files_sdk.Clickwrap {
	return i.Current().(files_sdk.Clickwrap)
}

func (c *Client) List(params files_sdk.ClickwrapListParams) (*Iter, error) {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/clickwraps"
	i.ListParams = &params
	exportParams, err := i.ExportParams()
	if err != nil {
		return i, err
	}
	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, exportParams)
		defaultValue := make([]interface{}, 0)
		if err != nil {
			return &defaultValue, "", err
		}
		list := files_sdk.ClickwrapCollection{}
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
	return i, nil
}

func List(params files_sdk.ClickwrapListParams) (*Iter, error) {
	return (&Client{}).List(params)
}

func (c *Client) Find(params files_sdk.ClickwrapFindParams) (files_sdk.Clickwrap, error) {
	clickwrap := files_sdk.Clickwrap{}
	if params.Id == 0 {
		return clickwrap, lib.CreateError(params, "Id")
	}
	path := "/clickwraps/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return clickwrap, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParms)
	if err != nil {
		return clickwrap, err
	}
	if res.StatusCode == 204 {
		return clickwrap, nil
	}
	if err := clickwrap.UnmarshalJSON(*data); err != nil {
		return clickwrap, err
	}

	return clickwrap, nil
}

func Find(params files_sdk.ClickwrapFindParams) (files_sdk.Clickwrap, error) {
	return (&Client{}).Find(params)
}

func (c *Client) Create(params files_sdk.ClickwrapCreateParams) (files_sdk.Clickwrap, error) {
	clickwrap := files_sdk.Clickwrap{}
	path := "/clickwraps"
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return clickwrap, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParms)
	if err != nil {
		return clickwrap, err
	}
	if res.StatusCode == 204 {
		return clickwrap, nil
	}
	if err := clickwrap.UnmarshalJSON(*data); err != nil {
		return clickwrap, err
	}

	return clickwrap, nil
}

func Create(params files_sdk.ClickwrapCreateParams) (files_sdk.Clickwrap, error) {
	return (&Client{}).Create(params)
}

func (c *Client) Update(params files_sdk.ClickwrapUpdateParams) (files_sdk.Clickwrap, error) {
	clickwrap := files_sdk.Clickwrap{}
	if params.Id == 0 {
		return clickwrap, lib.CreateError(params, "Id")
	}
	path := "/clickwraps/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return clickwrap, err
	}
	data, res, err := files_sdk.Call("PATCH", c.Config, path, exportedParms)
	if err != nil {
		return clickwrap, err
	}
	if res.StatusCode == 204 {
		return clickwrap, nil
	}
	if err := clickwrap.UnmarshalJSON(*data); err != nil {
		return clickwrap, err
	}

	return clickwrap, nil
}

func Update(params files_sdk.ClickwrapUpdateParams) (files_sdk.Clickwrap, error) {
	return (&Client{}).Update(params)
}

func (c *Client) Delete(params files_sdk.ClickwrapDeleteParams) (files_sdk.Clickwrap, error) {
	clickwrap := files_sdk.Clickwrap{}
	if params.Id == 0 {
		return clickwrap, lib.CreateError(params, "Id")
	}
	path := "/clickwraps/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return clickwrap, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParms)
	if err != nil {
		return clickwrap, err
	}
	if res.StatusCode == 204 {
		return clickwrap, nil
	}
	if err := clickwrap.UnmarshalJSON(*data); err != nil {
		return clickwrap, err
	}

	return clickwrap, nil
}

func Delete(params files_sdk.ClickwrapDeleteParams) (files_sdk.Clickwrap, error) {
	return (&Client{}).Delete(params)
}
