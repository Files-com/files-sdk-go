package file_comment

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
	*files_sdk.Iter
	*Client
}

func (i *Iter) Reload(opts ...files_sdk.RequestResponseOption) files_sdk.IterI {
	return &Iter{Iter: i.Iter.Reload(opts...).(*files_sdk.Iter), Client: i.Client}
}

func (i *Iter) FileComment() files_sdk.FileComment {
	return i.Current().(files_sdk.FileComment)
}

func (i *Iter) Iterate(identifier interface{}, opts ...files_sdk.RequestResponseOption) (files_sdk.IterI, error) {
	params := files_sdk.FileCommentListForParams{}
	if path, ok := identifier.(string); ok {
		params.Path = path
	}
	return i.Client.ListFor(context.Background(), params, opts...)
}

func (c *Client) ListFor(ctx context.Context, params files_sdk.FileCommentListForParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/file_comments/files/{path}", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.FileCommentCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list, opts...)
	return i, nil
}

func ListFor(ctx context.Context, params files_sdk.FileCommentListForParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).ListFor(ctx, params, opts...)
}

func (c *Client) Create(ctx context.Context, params files_sdk.FileCommentCreateParams, opts ...files_sdk.RequestResponseOption) (fileComment files_sdk.FileComment, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/file_comments", Params: params, Entity: &fileComment}, opts...)
	return
}

func Create(ctx context.Context, params files_sdk.FileCommentCreateParams, opts ...files_sdk.RequestResponseOption) (fileComment files_sdk.FileComment, err error) {
	return (&Client{}).Create(ctx, params, opts...)
}

func (c *Client) Update(ctx context.Context, params files_sdk.FileCommentUpdateParams, opts ...files_sdk.RequestResponseOption) (fileComment files_sdk.FileComment, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/file_comments/{id}", Params: params, Entity: &fileComment}, opts...)
	return
}

func Update(ctx context.Context, params files_sdk.FileCommentUpdateParams, opts ...files_sdk.RequestResponseOption) (fileComment files_sdk.FileComment, err error) {
	return (&Client{}).Update(ctx, params, opts...)
}

func (c *Client) UpdateWithMap(ctx context.Context, params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (fileComment files_sdk.FileComment, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/file_comments/{id}", Params: params, Entity: &fileComment}, opts...)
	return
}

func UpdateWithMap(ctx context.Context, params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (fileComment files_sdk.FileComment, err error) {
	return (&Client{}).UpdateWithMap(ctx, params, opts...)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.FileCommentDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/file_comments/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(ctx context.Context, params files_sdk.FileCommentDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(ctx, params, opts...)
}
