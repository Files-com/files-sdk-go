package as2_key

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

func (i *Iter) As2Key() files_sdk.As2Key {
	return i.Current().(files_sdk.As2Key)
}

func (c *Client) List(params files_sdk.As2KeyListParams) *Iter {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/as2_keys"

	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, i.ExportParams())
		defaultValue := make([]interface{}, 0)
		if err != nil {
			return &defaultValue, "", err
		}
		list := files_sdk.As2KeyCollection{}
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

func List(params files_sdk.As2KeyListParams) *Iter {
	return (&Client{}).List(params)
}

func (c *Client) Find(params files_sdk.As2KeyFindParams) (files_sdk.As2Key, error) {
	as2Key := files_sdk.As2Key{}
	path := "/as2_keys/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	data, res, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
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
	data, res, err := files_sdk.Call("POST", c.Config, path, lib.ExportParams(params))
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
	path := "/as2_keys/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	data, res, err := files_sdk.Call("PATCH", c.Config, path, lib.ExportParams(params))
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
	path := "/as2_keys/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	data, res, err := files_sdk.Call("DELETE", c.Config, path, lib.ExportParams(params))
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
