package bundle

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

func (i *Iter) Bundle() files_sdk.Bundle {
	return i.Current().(files_sdk.Bundle)
}

func (c *Client) List(params files_sdk.BundleListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/bundles"
	i.ListParams = &params
	list := files_sdk.BundleCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
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
	path := "/bundles/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return bundle, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return bundle, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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
	path := "/bundles/" + strconv.FormatInt(params.Id, 10) + "/share"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return bundle, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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
	path := "/bundles/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return bundle, err
	}
	data, res, err := files_sdk.Call("PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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
	path := "/bundles/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return bundle, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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
