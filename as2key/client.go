package as2_key

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

func (i *Iter) As2Key() files_sdk.As2Key {
	return i.Current().(files_sdk.As2Key)
}

func (c *Client) List(params files_sdk.As2KeyListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/as2_keys"
	i.ListParams = &params
	list := files_sdk.As2KeyCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func List(params files_sdk.As2KeyListParams) (*Iter, error) {
	return (&Client{}).List(params)
}

func (c *Client) Find(params files_sdk.As2KeyFindParams) (files_sdk.As2Key, error) {
	as2Key := files_sdk.As2Key{}
	if params.Id == 0 {
		return as2Key, lib.CreateError(params, "Id")
	}
	path := "/as2_keys/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return as2Key, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return as2Key, err
	}
	if res.StatusCode == 204 {
		return as2Key, nil
	}
	if err := as2Key.UnmarshalJSON(*data); err != nil {
		return as2Key, err
	}

	return as2Key, nil
}

func Find(params files_sdk.As2KeyFindParams) (files_sdk.As2Key, error) {
	return (&Client{}).Find(params)
}

func (c *Client) Create(params files_sdk.As2KeyCreateParams) (files_sdk.As2Key, error) {
	as2Key := files_sdk.As2Key{}
	path := "/as2_keys"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return as2Key, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return as2Key, err
	}
	if res.StatusCode == 204 {
		return as2Key, nil
	}
	if err := as2Key.UnmarshalJSON(*data); err != nil {
		return as2Key, err
	}

	return as2Key, nil
}

func Create(params files_sdk.As2KeyCreateParams) (files_sdk.As2Key, error) {
	return (&Client{}).Create(params)
}

func (c *Client) Update(params files_sdk.As2KeyUpdateParams) (files_sdk.As2Key, error) {
	as2Key := files_sdk.As2Key{}
	if params.Id == 0 {
		return as2Key, lib.CreateError(params, "Id")
	}
	path := "/as2_keys/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return as2Key, err
	}
	data, res, err := files_sdk.Call("PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return as2Key, err
	}
	if res.StatusCode == 204 {
		return as2Key, nil
	}
	if err := as2Key.UnmarshalJSON(*data); err != nil {
		return as2Key, err
	}

	return as2Key, nil
}

func Update(params files_sdk.As2KeyUpdateParams) (files_sdk.As2Key, error) {
	return (&Client{}).Update(params)
}

func (c *Client) Delete(params files_sdk.As2KeyDeleteParams) (files_sdk.As2Key, error) {
	as2Key := files_sdk.As2Key{}
	if params.Id == 0 {
		return as2Key, lib.CreateError(params, "Id")
	}
	path := "/as2_keys/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return as2Key, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return as2Key, err
	}
	if res.StatusCode == 204 {
		return as2Key, nil
	}
	if err := as2Key.UnmarshalJSON(*data); err != nil {
		return as2Key, err
	}

	return as2Key, nil
}

func Delete(params files_sdk.As2KeyDeleteParams) (files_sdk.As2Key, error) {
	return (&Client{}).Delete(params)
}
