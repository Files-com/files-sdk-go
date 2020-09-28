package api_key

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

func (i *Iter) ApiKey() files_sdk.ApiKey {
	return i.Current().(files_sdk.ApiKey)
}

func (c *Client) List(params files_sdk.ApiKeyListParams) (*Iter, error) {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/api_keys"
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
		list := files_sdk.ApiKeyCollection{}
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

func List(params files_sdk.ApiKeyListParams) (*Iter, error) {
	return (&Client{}).List(params)
}

func (c *Client) FindCurrent(params files_sdk.ApiKeyFindCurrentParams) (files_sdk.ApiKey, error) {
	apiKey := files_sdk.ApiKey{}
	path := "/api_key"
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return apiKey, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParms)
	if err != nil {
		return apiKey, err
	}
	if res.StatusCode == 204 {
		return apiKey, nil
	}
	if err := apiKey.UnmarshalJSON(*data); err != nil {
		return apiKey, err
	}

	return apiKey, nil
}

func FindCurrent(params files_sdk.ApiKeyFindCurrentParams) (files_sdk.ApiKey, error) {
	return (&Client{}).FindCurrent(params)
}

func (c *Client) Find(params files_sdk.ApiKeyFindParams) (files_sdk.ApiKey, error) {
	apiKey := files_sdk.ApiKey{}
	if params.Id == 0 {
		return apiKey, lib.CreateError(params, "Id")
	}
	path := "/api_keys/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return apiKey, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParms)
	if err != nil {
		return apiKey, err
	}
	if res.StatusCode == 204 {
		return apiKey, nil
	}
	if err := apiKey.UnmarshalJSON(*data); err != nil {
		return apiKey, err
	}

	return apiKey, nil
}

func Find(params files_sdk.ApiKeyFindParams) (files_sdk.ApiKey, error) {
	return (&Client{}).Find(params)
}

func (c *Client) Create(params files_sdk.ApiKeyCreateParams) (files_sdk.ApiKey, error) {
	apiKey := files_sdk.ApiKey{}
	path := "/api_keys"
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return apiKey, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParms)
	if err != nil {
		return apiKey, err
	}
	if res.StatusCode == 204 {
		return apiKey, nil
	}
	if err := apiKey.UnmarshalJSON(*data); err != nil {
		return apiKey, err
	}

	return apiKey, nil
}

func Create(params files_sdk.ApiKeyCreateParams) (files_sdk.ApiKey, error) {
	return (&Client{}).Create(params)
}

func (c *Client) UpdateCurrent(params files_sdk.ApiKeyUpdateCurrentParams) (files_sdk.ApiKey, error) {
	apiKey := files_sdk.ApiKey{}
	path := "/api_key"
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return apiKey, err
	}
	data, res, err := files_sdk.Call("PATCH", c.Config, path, exportedParms)
	if err != nil {
		return apiKey, err
	}
	if res.StatusCode == 204 {
		return apiKey, nil
	}
	if err := apiKey.UnmarshalJSON(*data); err != nil {
		return apiKey, err
	}

	return apiKey, nil
}

func UpdateCurrent(params files_sdk.ApiKeyUpdateCurrentParams) (files_sdk.ApiKey, error) {
	return (&Client{}).UpdateCurrent(params)
}

func (c *Client) Update(params files_sdk.ApiKeyUpdateParams) (files_sdk.ApiKey, error) {
	apiKey := files_sdk.ApiKey{}
	if params.Id == 0 {
		return apiKey, lib.CreateError(params, "Id")
	}
	path := "/api_keys/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return apiKey, err
	}
	data, res, err := files_sdk.Call("PATCH", c.Config, path, exportedParms)
	if err != nil {
		return apiKey, err
	}
	if res.StatusCode == 204 {
		return apiKey, nil
	}
	if err := apiKey.UnmarshalJSON(*data); err != nil {
		return apiKey, err
	}

	return apiKey, nil
}

func Update(params files_sdk.ApiKeyUpdateParams) (files_sdk.ApiKey, error) {
	return (&Client{}).Update(params)
}

func (c *Client) DeleteCurrent(params files_sdk.ApiKeyDeleteCurrentParams) (files_sdk.ApiKey, error) {
	apiKey := files_sdk.ApiKey{}
	path := "/api_key"
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return apiKey, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParms)
	if err != nil {
		return apiKey, err
	}
	if res.StatusCode == 204 {
		return apiKey, nil
	}
	if err := apiKey.UnmarshalJSON(*data); err != nil {
		return apiKey, err
	}

	return apiKey, nil
}

func DeleteCurrent(params files_sdk.ApiKeyDeleteCurrentParams) (files_sdk.ApiKey, error) {
	return (&Client{}).DeleteCurrent(params)
}

func (c *Client) Delete(params files_sdk.ApiKeyDeleteParams) (files_sdk.ApiKey, error) {
	apiKey := files_sdk.ApiKey{}
	if params.Id == 0 {
		return apiKey, lib.CreateError(params, "Id")
	}
	path := "/api_keys/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return apiKey, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParms)
	if err != nil {
		return apiKey, err
	}
	if res.StatusCode == 204 {
		return apiKey, nil
	}
	if err := apiKey.UnmarshalJSON(*data); err != nil {
		return apiKey, err
	}

	return apiKey, nil
}

func Delete(params files_sdk.ApiKeyDeleteParams) (files_sdk.ApiKey, error) {
	return (&Client{}).Delete(params)
}
