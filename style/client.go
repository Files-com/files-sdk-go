package style

import (
	"context"

	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Find(ctx context.Context, params files_sdk.StyleFindParams) (files_sdk.Style, error) {
	style := files_sdk.Style{}
	path := lib.BuildPath("/styles/", params.Path)
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return style, err
	}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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

func Find(ctx context.Context, params files_sdk.StyleFindParams) (files_sdk.Style, error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.StyleUpdateParams) (files_sdk.Style, error) {
	style := files_sdk.Style{}
	path := lib.BuildPath("/styles/", params.Path)
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return style, err
	}
	data, res, err := files_sdk.Call(ctx, "PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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

func Update(ctx context.Context, params files_sdk.StyleUpdateParams) (files_sdk.Style, error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.StyleDeleteParams) (files_sdk.Style, error) {
	style := files_sdk.Style{}
	path := lib.BuildPath("/styles/", params.Path)
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return style, err
	}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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

func Delete(ctx context.Context, params files_sdk.StyleDeleteParams) (files_sdk.Style, error) {
	return (&Client{}).Delete(ctx, params)
}
