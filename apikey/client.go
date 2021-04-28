package api_key

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

func (i *Iter) ApiKey() files_sdk.ApiKey {
	return i.Current().(files_sdk.ApiKey)
}

func (c *Client) List(params files_sdk.ApiKeyListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/api_keys"
	i.ListParams = &params
	list := files_sdk.ApiKeyCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func List(params files_sdk.ApiKeyListParams) (*Iter, error) {
	return (&Client{}).List(params)
}

func (c *Client) FindCurrent() (files_sdk.ApiKey, error) {
	apiKey := files_sdk.ApiKey{}
	path := "/api_key"
	exportedParams, err := lib.ExportParams(lib.Interface())
	if err != nil {
		return apiKey, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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

func FindCurrent() (files_sdk.ApiKey, error) {
	return (&Client{}).FindCurrent()
}

func (c *Client) Find(params files_sdk.ApiKeyFindParams) (files_sdk.ApiKey, error) {
	apiKey := files_sdk.ApiKey{}
	if params.Id == 0 {
		return apiKey, lib.CreateError(params, "Id")
	}
	path := "/api_keys/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return apiKey, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return apiKey, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return apiKey, err
	}
	data, res, err := files_sdk.Call("PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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
	path := "/api_keys/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return apiKey, err
	}
	data, res, err := files_sdk.Call("PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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

func (c *Client) DeleteCurrent() (files_sdk.ApiKey, error) {
	apiKey := files_sdk.ApiKey{}
	path := "/api_key"
	exportedParams, err := lib.ExportParams(lib.Interface())
	if err != nil {
		return apiKey, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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

func DeleteCurrent() (files_sdk.ApiKey, error) {
	return (&Client{}).DeleteCurrent()
}

func (c *Client) Delete(params files_sdk.ApiKeyDeleteParams) (files_sdk.ApiKey, error) {
	apiKey := files_sdk.ApiKey{}
	if params.Id == 0 {
		return apiKey, lib.CreateError(params, "Id")
	}
	path := "/api_keys/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return apiKey, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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
