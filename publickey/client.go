package public_key

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

func (i *Iter) PublicKey() files_sdk.PublicKey {
	return i.Current().(files_sdk.PublicKey)
}

func (c *Client) List(params files_sdk.PublicKeyListParams) (*Iter, error) {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/public_keys"
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
		list := files_sdk.PublicKeyCollection{}
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

func List(params files_sdk.PublicKeyListParams) (*Iter, error) {
	return (&Client{}).List(params)
}

func (c *Client) Find(params files_sdk.PublicKeyFindParams) (files_sdk.PublicKey, error) {
	publicKey := files_sdk.PublicKey{}
	if params.Id == 0 {
		return publicKey, lib.CreateError(params, "Id")
	}
	path := "/public_keys/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return publicKey, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParms)
	if err != nil {
		return publicKey, err
	}
	if res.StatusCode == 204 {
		return publicKey, nil
	}
	if err := publicKey.UnmarshalJSON(*data); err != nil {
		return publicKey, err
	}

	return publicKey, nil
}

func Find(params files_sdk.PublicKeyFindParams) (files_sdk.PublicKey, error) {
	return (&Client{}).Find(params)
}

func (c *Client) Create(params files_sdk.PublicKeyCreateParams) (files_sdk.PublicKey, error) {
	publicKey := files_sdk.PublicKey{}
	path := "/public_keys"
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return publicKey, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParms)
	if err != nil {
		return publicKey, err
	}
	if res.StatusCode == 204 {
		return publicKey, nil
	}
	if err := publicKey.UnmarshalJSON(*data); err != nil {
		return publicKey, err
	}

	return publicKey, nil
}

func Create(params files_sdk.PublicKeyCreateParams) (files_sdk.PublicKey, error) {
	return (&Client{}).Create(params)
}

func (c *Client) Update(params files_sdk.PublicKeyUpdateParams) (files_sdk.PublicKey, error) {
	publicKey := files_sdk.PublicKey{}
	if params.Id == 0 {
		return publicKey, lib.CreateError(params, "Id")
	}
	path := "/public_keys/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return publicKey, err
	}
	data, res, err := files_sdk.Call("PATCH", c.Config, path, exportedParms)
	if err != nil {
		return publicKey, err
	}
	if res.StatusCode == 204 {
		return publicKey, nil
	}
	if err := publicKey.UnmarshalJSON(*data); err != nil {
		return publicKey, err
	}

	return publicKey, nil
}

func Update(params files_sdk.PublicKeyUpdateParams) (files_sdk.PublicKey, error) {
	return (&Client{}).Update(params)
}

func (c *Client) Delete(params files_sdk.PublicKeyDeleteParams) (files_sdk.PublicKey, error) {
	publicKey := files_sdk.PublicKey{}
	if params.Id == 0 {
		return publicKey, lib.CreateError(params, "Id")
	}
	path := "/public_keys/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return publicKey, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParms)
	if err != nil {
		return publicKey, err
	}
	if res.StatusCode == 204 {
		return publicKey, nil
	}
	if err := publicKey.UnmarshalJSON(*data); err != nil {
		return publicKey, err
	}

	return publicKey, nil
}

func Delete(params files_sdk.PublicKeyDeleteParams) (files_sdk.PublicKey, error) {
	return (&Client{}).Delete(params)
}
