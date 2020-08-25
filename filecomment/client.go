package file_comment

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

func (i *Iter) FileComment() files_sdk.FileComment {
	return i.Current().(files_sdk.FileComment)
}

func (c *Client) ListFor(params files_sdk.FileCommentListForParams) (*Iter, error) {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/file_comments/files/" + lib.QueryEscape(params.Path) + ""
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
		list := files_sdk.FileCommentCollection{}
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

func ListFor(params files_sdk.FileCommentListForParams) (*Iter, error) {
	return (&Client{}).ListFor(params)
}

func (c *Client) Create(params files_sdk.FileCommentCreateParams) (files_sdk.FileComment, error) {
	fileComment := files_sdk.FileComment{}
	path := "/file_comments"
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return fileComment, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParms)
	if err != nil {
		return fileComment, err
	}
	if res.StatusCode == 204 {
		return fileComment, nil
	}
	if err := fileComment.UnmarshalJSON(*data); err != nil {
		return fileComment, err
	}

	return fileComment, nil
}

func Create(params files_sdk.FileCommentCreateParams) (files_sdk.FileComment, error) {
	return (&Client{}).Create(params)
}

func (c *Client) Update(params files_sdk.FileCommentUpdateParams) (files_sdk.FileComment, error) {
	fileComment := files_sdk.FileComment{}
	if params.Id == 0 {
		return fileComment, lib.CreateError(params, "Id")
	}
	path := "/file_comments/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return fileComment, err
	}
	data, res, err := files_sdk.Call("PATCH", c.Config, path, exportedParms)
	if err != nil {
		return fileComment, err
	}
	if res.StatusCode == 204 {
		return fileComment, nil
	}
	if err := fileComment.UnmarshalJSON(*data); err != nil {
		return fileComment, err
	}

	return fileComment, nil
}

func Update(params files_sdk.FileCommentUpdateParams) (files_sdk.FileComment, error) {
	return (&Client{}).Update(params)
}

func (c *Client) Delete(params files_sdk.FileCommentDeleteParams) (files_sdk.FileComment, error) {
	fileComment := files_sdk.FileComment{}
	if params.Id == 0 {
		return fileComment, lib.CreateError(params, "Id")
	}
	path := "/file_comments/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return fileComment, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParms)
	if err != nil {
		return fileComment, err
	}
	if res.StatusCode == 204 {
		return fileComment, nil
	}
	if err := fileComment.UnmarshalJSON(*data); err != nil {
		return fileComment, err
	}

	return fileComment, nil
}

func Delete(params files_sdk.FileCommentDeleteParams) (files_sdk.FileComment, error) {
	return (&Client{}).Delete(params)
}
