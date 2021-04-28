package automation

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

func (i *Iter) Automation() files_sdk.Automation {
	return i.Current().(files_sdk.Automation)
}

func (c *Client) List(params files_sdk.AutomationListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/automations"
	i.ListParams = &params
	list := files_sdk.AutomationCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
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
	path := "/automations/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return automation, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParams)
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

func Find(params files_sdk.AutomationFindParams) (files_sdk.Automation, error) {
	return (&Client{}).Find(params)
}

func (c *Client) Create(params files_sdk.AutomationCreateParams) (files_sdk.Automation, error) {
	automation := files_sdk.Automation{}
	path := "/automations"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return automation, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
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

func Create(params files_sdk.AutomationCreateParams) (files_sdk.Automation, error) {
	return (&Client{}).Create(params)
}

func (c *Client) Update(params files_sdk.AutomationUpdateParams) (files_sdk.Automation, error) {
	automation := files_sdk.Automation{}
	if params.Id == 0 {
		return automation, lib.CreateError(params, "Id")
	}
	path := "/automations/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return automation, err
	}
	data, res, err := files_sdk.Call("PATCH", c.Config, path, exportedParams)
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

func Update(params files_sdk.AutomationUpdateParams) (files_sdk.Automation, error) {
	return (&Client{}).Update(params)
}

func (c *Client) Delete(params files_sdk.AutomationDeleteParams) (files_sdk.Automation, error) {
	automation := files_sdk.Automation{}
	if params.Id == 0 {
		return automation, lib.CreateError(params, "Id")
	}
	path := "/automations/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return automation, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParams)
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

func Delete(params files_sdk.AutomationDeleteParams) (files_sdk.Automation, error) {
	return (&Client{}).Delete(params)
}
