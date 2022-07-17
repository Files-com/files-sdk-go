package folder

import (
	"context"

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

func (i *Iter) Folder() files_sdk.Folder {
	return i.Current().(files_sdk.Folder)
}

func (c *Client) ListFor(ctx context.Context, params files_sdk.FolderListForParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/folders/{path}", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.FolderCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func ListFor(ctx context.Context, params files_sdk.FolderListForParams) (*Iter, error) {
	return (&Client{}).ListFor(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.FolderCreateParams) (file files_sdk.File, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/folders/{path}", Params: params, Entity: &file})
	return
}

func Create(ctx context.Context, params files_sdk.FolderCreateParams) (file files_sdk.File, err error) {
	return (&Client{}).Create(ctx, params)
}
