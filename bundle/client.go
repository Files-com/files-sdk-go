package bundle

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

func (i *Iter) Bundle() files_sdk.Bundle {
	return i.Current().(files_sdk.Bundle)
}

func (c *Client) List(params files_sdk.BundleListParams) (*Iter, error) {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/bundles"
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
		list := files_sdk.BundleCollection{}
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

func List(params files_sdk.BundleListParams) (*Iter, error) {
	return (&Client{}).List(params)
}

func (c *Client) Find(params files_sdk.BundleFindParams) (files_sdk.Bundle, error) {
	bundle := files_sdk.Bundle{}
	if params.Id == 0 {
		return bundle, lib.CreateError(params, "Id")
	}
	path := "/bundles/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return bundle, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParms)
	if err != nil {
		return bundle, err
	}
	if res.StatusCode == 204 {
		return bundle, nil
	}
	if err := bundle.UnmarshalJSON(*data); err != nil {
		return bundle, err
	}

	return bundle, nil
}

func Find(params files_sdk.BundleFindParams) (files_sdk.Bundle, error) {
	return (&Client{}).Find(params)
}

func (c *Client) Create(params files_sdk.BundleCreateParams) (files_sdk.Bundle, error) {
	bundle := files_sdk.Bundle{}
	path := "/bundles"
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return bundle, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParms)
	if err != nil {
		return bundle, err
	}
	if res.StatusCode == 204 {
		return bundle, nil
	}
	if err := bundle.UnmarshalJSON(*data); err != nil {
		return bundle, err
	}

	return bundle, nil
}

func Create(params files_sdk.BundleCreateParams) (files_sdk.Bundle, error) {
	return (&Client{}).Create(params)
}

func (c *Client) Share(params files_sdk.BundleShareParams) (files_sdk.Bundle, error) {
	bundle := files_sdk.Bundle{}
	if params.Id == 0 {
		return bundle, lib.CreateError(params, "Id")
	}
	path := "/bundles/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + "/share"
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return bundle, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParms)
	if err != nil {
		return bundle, err
	}
	if res.StatusCode == 204 {
		return bundle, nil
	}
	if err := bundle.UnmarshalJSON(*data); err != nil {
		return bundle, err
	}

	return bundle, nil
}

func Share(params files_sdk.BundleShareParams) (files_sdk.Bundle, error) {
	return (&Client{}).Share(params)
}

func (c *Client) Update(params files_sdk.BundleUpdateParams) (files_sdk.Bundle, error) {
	bundle := files_sdk.Bundle{}
	if params.Id == 0 {
		return bundle, lib.CreateError(params, "Id")
	}
	path := "/bundles/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return bundle, err
	}
	data, res, err := files_sdk.Call("PATCH", c.Config, path, exportedParms)
	if err != nil {
		return bundle, err
	}
	if res.StatusCode == 204 {
		return bundle, nil
	}
	if err := bundle.UnmarshalJSON(*data); err != nil {
		return bundle, err
	}

	return bundle, nil
}

func Update(params files_sdk.BundleUpdateParams) (files_sdk.Bundle, error) {
	return (&Client{}).Update(params)
}

func (c *Client) Delete(params files_sdk.BundleDeleteParams) (files_sdk.Bundle, error) {
	bundle := files_sdk.Bundle{}
	if params.Id == 0 {
		return bundle, lib.CreateError(params, "Id")
	}
	path := "/bundles/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return bundle, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParms)
	if err != nil {
		return bundle, err
	}
	if res.StatusCode == 204 {
		return bundle, nil
	}
	if err := bundle.UnmarshalJSON(*data); err != nil {
		return bundle, err
	}

	return bundle, nil
}

func Delete(params files_sdk.BundleDeleteParams) (files_sdk.Bundle, error) {
	return (&Client{}).Delete(params)
}
