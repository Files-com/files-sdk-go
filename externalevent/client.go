package external_event

import (
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

func (i *Iter) ExternalEvent() files_sdk.ExternalEvent {
	return i.Current().(files_sdk.ExternalEvent)
}

func (c *Client) List(params files_sdk.ExternalEventListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/external_events"
	i.ListParams = &params
	list := files_sdk.ExternalEventCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func List(params files_sdk.ExternalEventListParams) (*Iter, error) {
	return (&Client{}).List(params)
}
