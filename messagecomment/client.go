package message_comment

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

func (i *Iter) MessageComment() files_sdk.MessageComment {
	return i.Current().(files_sdk.MessageComment)
}

func (c *Client) List(params files_sdk.MessageCommentListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/message_comments"
	i.ListParams = &params
	list := files_sdk.MessageCommentCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func List(params files_sdk.MessageCommentListParams) (*Iter, error) {
	return (&Client{}).List(params)
}

func (c *Client) Find(params files_sdk.MessageCommentFindParams) (files_sdk.MessageComment, error) {
	messageComment := files_sdk.MessageComment{}
	if params.Id == 0 {
		return messageComment, lib.CreateError(params, "Id")
	}
	path := "/message_comments/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return messageComment, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParams)
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

func Find(params files_sdk.MessageCommentFindParams) (files_sdk.MessageComment, error) {
	return (&Client{}).Find(params)
}

func (c *Client) Create(params files_sdk.MessageCommentCreateParams) (files_sdk.MessageComment, error) {
	messageComment := files_sdk.MessageComment{}
	path := "/message_comments"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return messageComment, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
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

func Create(params files_sdk.MessageCommentCreateParams) (files_sdk.MessageComment, error) {
	return (&Client{}).Create(params)
}

func (c *Client) Update(params files_sdk.MessageCommentUpdateParams) (files_sdk.MessageComment, error) {
	messageComment := files_sdk.MessageComment{}
	if params.Id == 0 {
		return messageComment, lib.CreateError(params, "Id")
	}
	path := "/message_comments/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return messageComment, err
	}
	data, res, err := files_sdk.Call("PATCH", c.Config, path, exportedParams)
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

func Update(params files_sdk.MessageCommentUpdateParams) (files_sdk.MessageComment, error) {
	return (&Client{}).Update(params)
}

func (c *Client) Delete(params files_sdk.MessageCommentDeleteParams) (files_sdk.MessageComment, error) {
	messageComment := files_sdk.MessageComment{}
	if params.Id == 0 {
		return messageComment, lib.CreateError(params, "Id")
	}
	path := "/message_comments/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return messageComment, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParams)
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

func Delete(params files_sdk.MessageCommentDeleteParams) (files_sdk.MessageComment, error) {
	return (&Client{}).Delete(params)
}
