package sync_job

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

func (i *Iter) SyncJob() files_sdk.SyncJob {
	return i.Current().(files_sdk.SyncJob)
}

func (c *Client) List(params files_sdk.SyncJobListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/sync_jobs"
	i.ListParams = &params
	list := files_sdk.SyncJobCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func List(params files_sdk.SyncJobListParams) (*Iter, error) {
	return (&Client{}).List(params)
}
