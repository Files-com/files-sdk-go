package automation

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

func (i *Iter) Automation() files_sdk.Automation {
	return i.Current().(files_sdk.Automation)
}

func (c *Client) List(params files_sdk.AutomationListParams) (*Iter, error) {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/automations"
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
		list := files_sdk.AutomationCollection{}
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

func List(params files_sdk.AutomationListParams) (*Iter, error) {
	return (&Client{}).List(params)
}

func (c *Client) Find(params files_sdk.AutomationFindParams) (files_sdk.Automation, error) {
	automation := files_sdk.Automation{}
	if params.Id == 0 {
		return automation, lib.CreateError(params, "Id")
	}
	path := "/automations/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return automation, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParms)
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

func Find(params files_sdk.AutomationFindParams) (files_sdk.Automation, error) {
	return (&Client{}).Find(params)
}

func (c *Client) Create(params files_sdk.AutomationCreateParams) (files_sdk.Automation, error) {
	automation := files_sdk.Automation{}
	path := "/automations"
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return automation, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParms)
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

func Create(params files_sdk.AutomationCreateParams) (files_sdk.Automation, error) {
	return (&Client{}).Create(params)
}

func (c *Client) Update(params files_sdk.AutomationUpdateParams) (files_sdk.Automation, error) {
	automation := files_sdk.Automation{}
	if params.Id == 0 {
		return automation, lib.CreateError(params, "Id")
	}
	path := "/automations/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return automation, err
	}
	data, res, err := files_sdk.Call("PATCH", c.Config, path, exportedParms)
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

func Update(params files_sdk.AutomationUpdateParams) (files_sdk.Automation, error) {
	return (&Client{}).Update(params)
}

func (c *Client) Delete(params files_sdk.AutomationDeleteParams) (files_sdk.Automation, error) {
	automation := files_sdk.Automation{}
	if params.Id == 0 {
		return automation, lib.CreateError(params, "Id")
	}
	path := "/automations/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return automation, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParms)
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

func Delete(params files_sdk.AutomationDeleteParams) (files_sdk.Automation, error) {
	return (&Client{}).Delete(params)
}
