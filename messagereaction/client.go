package message_reaction

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

func (i *Iter) MessageReaction() files_sdk.MessageReaction {
	return i.Current().(files_sdk.MessageReaction)
}

func (c *Client) List(params files_sdk.MessageReactionListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/message_reactions"
	i.ListParams = &params
	list := files_sdk.MessageReactionCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func List(params files_sdk.MessageReactionListParams) (*Iter, error) {
	return (&Client{}).List(params)
}

func (c *Client) Find(params files_sdk.MessageReactionFindParams) (files_sdk.MessageReaction, error) {
	messageReaction := files_sdk.MessageReaction{}
	if params.Id == 0 {
		return messageReaction, lib.CreateError(params, "Id")
	}
	path := "/message_reactions/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return messageReaction, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return messageReaction, err
	}
	if res.StatusCode == 204 {
		return messageReaction, nil
	}
	if err := messageReaction.UnmarshalJSON(*data); err != nil {
		return messageReaction, err
	}

	return messageReaction, nil
}

func Find(params files_sdk.MessageReactionFindParams) (files_sdk.MessageReaction, error) {
	return (&Client{}).Find(params)
}

func (c *Client) Create(params files_sdk.MessageReactionCreateParams) (files_sdk.MessageReaction, error) {
	messageReaction := files_sdk.MessageReaction{}
	path := "/message_reactions"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return messageReaction, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return messageReaction, err
	}
	if res.StatusCode == 204 {
		return messageReaction, nil
	}
	if err := messageReaction.UnmarshalJSON(*data); err != nil {
		return messageReaction, err
	}

	return messageReaction, nil
}

func Create(params files_sdk.MessageReactionCreateParams) (files_sdk.MessageReaction, error) {
	return (&Client{}).Create(params)
}

func (c *Client) Delete(params files_sdk.MessageReactionDeleteParams) (files_sdk.MessageReaction, error) {
	messageReaction := files_sdk.MessageReaction{}
	if params.Id == 0 {
		return messageReaction, lib.CreateError(params, "Id")
	}
	path := "/message_reactions/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return messageReaction, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return messageReaction, err
	}
	if res.StatusCode == 204 {
		return messageReaction, nil
	}
	if err := messageReaction.UnmarshalJSON(*data); err != nil {
		return messageReaction, err
	}

	return messageReaction, nil
}

func Delete(params files_sdk.MessageReactionDeleteParams) (files_sdk.MessageReaction, error) {
	return (&Client{}).Delete(params)
}
