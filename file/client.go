package file

import (
	"io"

	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Download(params files_sdk.FileDownloadParams) (files_sdk.File, error) {
	file := files_sdk.File{}
	path := lib.BuildPath("/files/", params.Path)
	exportParams, err := lib.ExportParams(params)
	if err != nil {
		return file, err
	}
	data, _, err := files_sdk.Call("GET", c.Config, path, exportParams)
	if err != nil {
		return file, err
	}
	if err := file.UnmarshalJSON(*data); err != nil {
		return file, err
	}

	resp, err := c.Config.GetHttpClient().Get(file.DownloadUri)
	if err != nil {
		return file, err
	}
	_, err = io.Copy(params.Writer, resp.Body)
	if err != nil {
		return file, err
	}

	return file, nil
}

func Download(params files_sdk.FileDownloadParams) (files_sdk.File, error) {
	client := Client{}
	return client.Download(params)
}

func (c *Client) Create(params files_sdk.FileCreateParams) (files_sdk.File, error) {
	file := files_sdk.File{}
	path := lib.BuildPath("/files/", params.Path)
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return file, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
	if err != nil {
		return file, err
	}
	if res.StatusCode == 204 {
		return file, nil
	}
	if err := file.UnmarshalJSON(*data); err != nil {
		return file, err
	}

	return file, nil
}

func Create(params files_sdk.FileCreateParams) (files_sdk.File, error) {
	return (&Client{}).Create(params)
}

func (c *Client) Update(params files_sdk.FileUpdateParams) (files_sdk.File, error) {
	file := files_sdk.File{}
	path := lib.BuildPath("/files/", params.Path)
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return file, err
	}
	data, res, err := files_sdk.Call("PATCH", c.Config, path, exportedParams)
	if err != nil {
		return file, err
	}
	if res.StatusCode == 204 {
		return file, nil
	}
	if err := file.UnmarshalJSON(*data); err != nil {
		return file, err
	}

	return file, nil
}

func Update(params files_sdk.FileUpdateParams) (files_sdk.File, error) {
	return (&Client{}).Update(params)
}

func (c *Client) Delete(params files_sdk.FileDeleteParams) (files_sdk.File, error) {
	file := files_sdk.File{}
	path := lib.BuildPath("/files/", params.Path)
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return file, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParams)
	if err != nil {
		return file, err
	}
	if res.StatusCode == 204 {
		return file, nil
	}
	if err := file.UnmarshalJSON(*data); err != nil {
		return file, err
	}

	return file, nil
}

func Delete(params files_sdk.FileDeleteParams) (files_sdk.File, error) {
	return (&Client{}).Delete(params)
}
