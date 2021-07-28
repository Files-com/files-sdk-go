package message

import (
	"context"
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

func (i *Iter) Message() files_sdk.Message {
	return i.Current().(files_sdk.Message)
}

func (c *Client) List(ctx context.Context, params files_sdk.MessageListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/messages"
	i.ListParams = &params
	list := files_sdk.MessageCollection{}
	i.Query = listquery.Build(ctx, i, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.MessageListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.MessageFindParams) (files_sdk.Message, error) {
	message := files_sdk.Message{}
	if params.Id == 0 {
		return message, lib.CreateError(params, "Id")
	}
	path := "/messages/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return message, err
	}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return message, err
	}
	if res.StatusCode == 204 {
		return message, nil
	}
	if err := message.UnmarshalJSON(*data); err != nil {
		return message, err
	}

	return message, nil
}

func Find(ctx context.Context, params files_sdk.MessageFindParams) (files_sdk.Message, error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.MessageCreateParams) (files_sdk.Message, error) {
	message := files_sdk.Message{}
	path := "/messages"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return message, err
	}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return message, err
	}
	if res.StatusCode == 204 {
		return message, nil
	}
	if err := message.UnmarshalJSON(*data); err != nil {
		return message, err
	}

	return message, nil
}

func Create(ctx context.Context, params files_sdk.MessageCreateParams) (files_sdk.Message, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.MessageUpdateParams) (files_sdk.Message, error) {
	message := files_sdk.Message{}
	if params.Id == 0 {
		return message, lib.CreateError(params, "Id")
	}
	path := "/messages/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return message, err
	}
	data, res, err := files_sdk.Call(ctx, "PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return message, err
	}
	if res.StatusCode == 204 {
		return message, nil
	}
	if err := message.UnmarshalJSON(*data); err != nil {
		return message, err
	}

	return message, nil
}

func Update(ctx context.Context, params files_sdk.MessageUpdateParams) (files_sdk.Message, error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.MessageDeleteParams) (files_sdk.Message, error) {
	message := files_sdk.Message{}
	if params.Id == 0 {
		return message, lib.CreateError(params, "Id")
	}
	path := "/messages/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return message, err
	}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return message, err
	}
	if res.StatusCode == 204 {
		return message, nil
	}
	if err := message.UnmarshalJSON(*data); err != nil {
		return message, err
	}

	return message, nil
}

func Delete(ctx context.Context, params files_sdk.MessageDeleteParams) (files_sdk.Message, error) {
	return (&Client{}).Delete(ctx, params)
}
