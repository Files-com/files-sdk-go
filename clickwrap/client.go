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

func (c *Client) List(params files_sdk.ClickwrapListParams) *Iter {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/clickwraps"

	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, i.ExportParams())
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
	i.ListParams = &params
	return i
}

func List(params files_sdk.ClickwrapListParams) *Iter {
	return (&Client{}).List(params)
}

func (c *Client) Find(params files_sdk.ClickwrapFindParams) (files_sdk.Clickwrap, error) {
	clickwrap := files_sdk.Clickwrap{}
	path := "/clickwraps/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	data, res, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
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
	data, res, err := files_sdk.Call("POST", c.Config, path, lib.ExportParams(params))
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
	path := "/clickwraps/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	data, res, err := files_sdk.Call("PATCH", c.Config, path, lib.ExportParams(params))
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
	path := "/clickwraps/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	data, res, err := files_sdk.Call("DELETE", c.Config, path, lib.ExportParams(params))
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
