package clickwrap

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

func (i *Iter) Clickwrap() files_sdk.Clickwrap {
	return i.Current().(files_sdk.Clickwrap)
}

func (c *Client) List(params files_sdk.ClickwrapListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/clickwraps"
	i.ListParams = &params
	list := files_sdk.ClickwrapCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
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
	path := "/clickwraps/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return clickwrap, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return clickwrap, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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
	path := "/clickwraps/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return clickwrap, err
	}
	data, res, err := files_sdk.Call("PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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
	path := "/clickwraps/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return clickwrap, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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
