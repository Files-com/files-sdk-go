package inbox_registration

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

func (i *Iter) InboxRegistration() files_sdk.InboxRegistration {
	return i.Current().(files_sdk.InboxRegistration)
}

func (c *Client) List(params files_sdk.InboxRegistrationListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/inbox_registrations"
	i.ListParams = &params
	list := files_sdk.InboxRegistrationCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func List(params files_sdk.InboxRegistrationListParams) (*Iter, error) {
	return (&Client{}).List(params)
}
