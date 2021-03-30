package public_key

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

func (i *Iter) PublicKey() files_sdk.PublicKey {
	return i.Current().(files_sdk.PublicKey)
}

func (c *Client) List(params files_sdk.PublicKeyListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/public_keys"
	i.ListParams = &params
	list := files_sdk.PublicKeyCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
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
	path := "/public_keys/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return publicKey, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParams)
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
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return publicKey, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
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
	path := "/public_keys/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return publicKey, err
	}
	data, res, err := files_sdk.Call("PATCH", c.Config, path, exportedParams)
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
	path := "/public_keys/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return publicKey, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParams)
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
