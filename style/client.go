package style

import (
	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Find(params files_sdk.StyleFindParams) (files_sdk.Style, error) {
	style := files_sdk.Style{}
	path := lib.BuildPath("/styles/", params.Path)
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return style, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParams)
	if err != nil {
		return style, err
	}
	if res.StatusCode == 204 {
		return style, nil
	}
	if err := style.UnmarshalJSON(*data); err != nil {
		return style, err
	}

	return style, nil
}

func Find(params files_sdk.StyleFindParams) (files_sdk.Style, error) {
	return (&Client{}).Find(params)
}

func (c *Client) Update(params files_sdk.StyleUpdateParams) (files_sdk.Style, error) {
	style := files_sdk.Style{}
	path := lib.BuildPath("/styles/", params.Path)
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return style, err
	}
	data, res, err := files_sdk.Call("PATCH", c.Config, path, exportedParams)
	if err != nil {
		return style, err
	}
	if res.StatusCode == 204 {
		return style, nil
	}
	if err := style.UnmarshalJSON(*data); err != nil {
		return style, err
	}

	return style, nil
}

func Update(params files_sdk.StyleUpdateParams) (files_sdk.Style, error) {
	return (&Client{}).Update(params)
}

func (c *Client) Delete(params files_sdk.StyleDeleteParams) (files_sdk.Style, error) {
	style := files_sdk.Style{}
	path := lib.BuildPath("/styles/", params.Path)
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return style, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParams)
	if err != nil {
		return style, err
	}
	if res.StatusCode == 204 {
		return style, nil
	}
	if err := style.UnmarshalJSON(*data); err != nil {
		return style, err
	}

	return style, nil
}

func Delete(params files_sdk.StyleDeleteParams) (files_sdk.Style, error) {
	return (&Client{}).Delete(params)
}
