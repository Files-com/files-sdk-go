package message_comment

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

func (i *Iter) MessageComment() files_sdk.MessageComment {
	return i.Current().(files_sdk.MessageComment)
}

func (c *Client) List(ctx context.Context, params files_sdk.MessageCommentListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/message_comments"
	i.ListParams = &params
	list := files_sdk.MessageCommentCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.MessageCommentListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.MessageCommentFindParams) (files_sdk.MessageComment, error) {
	messageComment := files_sdk.MessageComment{}
	if params.Id == 0 {
		return messageComment, lib.CreateError(params, "Id")
	}
	path := "/message_comments/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return messageComment, err
	}
	if res.StatusCode == 204 {
		return messageComment, nil
	}
	if err := messageComment.UnmarshalJSON(*data); err != nil {
		return messageComment, err
	}

	return messageComment, nil
}

func Find(ctx context.Context, params files_sdk.MessageCommentFindParams) (files_sdk.MessageComment, error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.MessageCommentCreateParams) (files_sdk.MessageComment, error) {
	messageComment := files_sdk.MessageComment{}
	path := "/message_comments"
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return messageComment, err
	}
	if res.StatusCode == 204 {
		return messageComment, nil
	}
	if err := messageComment.UnmarshalJSON(*data); err != nil {
		return messageComment, err
	}

	return messageComment, nil
}

func Create(ctx context.Context, params files_sdk.MessageCommentCreateParams) (files_sdk.MessageComment, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.MessageCommentUpdateParams) (files_sdk.MessageComment, error) {
	messageComment := files_sdk.MessageComment{}
	if params.Id == 0 {
		return messageComment, lib.CreateError(params, "Id")
	}
	path := "/message_comments/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return messageComment, err
	}
	if res.StatusCode == 204 {
		return messageComment, nil
	}
	if err := messageComment.UnmarshalJSON(*data); err != nil {
		return messageComment, err
	}

	return messageComment, nil
}

func Update(ctx context.Context, params files_sdk.MessageCommentUpdateParams) (files_sdk.MessageComment, error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.MessageCommentDeleteParams) (files_sdk.MessageComment, error) {
	messageComment := files_sdk.MessageComment{}
	if params.Id == 0 {
		return messageComment, lib.CreateError(params, "Id")
	}
	path := "/message_comments/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return messageComment, err
	}
	if res.StatusCode == 204 {
		return messageComment, nil
	}
	if err := messageComment.UnmarshalJSON(*data); err != nil {
		return messageComment, err
	}

	return messageComment, nil
}

func Delete(ctx context.Context, params files_sdk.MessageCommentDeleteParams) (files_sdk.MessageComment, error) {
	return (&Client{}).Delete(ctx, params)
}
