package file_comment

import (
	"context"
	"strconv"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
	listquery "github.com/Files-com/files-sdk-go/v2/listquery"
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

func (c *Client) ListFor(ctx context.Context, params files_sdk.FileCommentListForParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := lib.BuildPath("/file_comments/files/", params.Path)
	i.ListParams = &params
	list := files_sdk.FileCommentCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func ListFor(ctx context.Context, params files_sdk.FileCommentListForParams) (*Iter, error) {
	return (&Client{}).ListFor(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.FileCommentCreateParams) (files_sdk.FileComment, error) {
	fileComment := files_sdk.FileComment{}
	path := "/file_comments"
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return fileComment, err
	}
	if res.StatusCode == 204 {
		return fileComment, nil
	}

	return fileComment, fileComment.UnmarshalJSON(*data)
}

func Create(ctx context.Context, params files_sdk.FileCommentCreateParams) (files_sdk.FileComment, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.FileCommentUpdateParams) (files_sdk.FileComment, error) {
	fileComment := files_sdk.FileComment{}
	if params.Id == 0 {
		return fileComment, lib.CreateError(params, "Id")
	}
	path := "/file_comments/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return fileComment, err
	}
	if res.StatusCode == 204 {
		return fileComment, nil
	}

	return fileComment, fileComment.UnmarshalJSON(*data)
}

func Update(ctx context.Context, params files_sdk.FileCommentUpdateParams) (files_sdk.FileComment, error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.FileCommentDeleteParams) error {
	fileComment := files_sdk.FileComment{}
	if params.Id == 0 {
		return lib.CreateError(params, "Id")
	}
	path := "/file_comments/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return err
	}
	if res.StatusCode == 204 {
		return nil
	}

	return fileComment.UnmarshalJSON(*data)
}

func Delete(ctx context.Context, params files_sdk.FileCommentDeleteParams) error {
	return (&Client{}).Delete(ctx, params)
}
