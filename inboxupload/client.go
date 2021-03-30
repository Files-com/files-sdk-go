package inbox_upload

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

func (i *Iter) InboxUpload() files_sdk.InboxUpload {
	return i.Current().(files_sdk.InboxUpload)
}

func (c *Client) List(params files_sdk.InboxUploadListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/inbox_uploads"
	i.ListParams = &params
	list := files_sdk.InboxUploadCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func List(params files_sdk.InboxUploadListParams) (*Iter, error) {
	return (&Client{}).List(params)
}
