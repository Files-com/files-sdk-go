package external_event

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

func (i *Iter) ExternalEvent() files_sdk.ExternalEvent {
	return i.Current().(files_sdk.ExternalEvent)
}

func (c *Client) List(ctx context.Context, params files_sdk.ExternalEventListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/external_events"
	i.ListParams = &params
	list := files_sdk.ExternalEventCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.ExternalEventListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.ExternalEventFindParams) (files_sdk.ExternalEvent, error) {
	externalEvent := files_sdk.ExternalEvent{}
	if params.Id == 0 {
		return externalEvent, lib.CreateError(params, "Id")
	}
	path := "/external_events/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return externalEvent, err
	}
	if res.StatusCode == 204 {
		return externalEvent, nil
	}

	return externalEvent, externalEvent.UnmarshalJSON(*data)
}

func Find(ctx context.Context, params files_sdk.ExternalEventFindParams) (files_sdk.ExternalEvent, error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.ExternalEventCreateParams) (files_sdk.ExternalEvent, error) {
	externalEvent := files_sdk.ExternalEvent{}
	path := "/external_events"
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return externalEvent, err
	}
	if res.StatusCode == 204 {
		return externalEvent, nil
	}

	return externalEvent, externalEvent.UnmarshalJSON(*data)
}

func Create(ctx context.Context, params files_sdk.ExternalEventCreateParams) (files_sdk.ExternalEvent, error) {
	return (&Client{}).Create(ctx, params)
}
