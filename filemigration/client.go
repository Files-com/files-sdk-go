package file_migration

import (
	"context"
	"time"

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

func (i *Iter) FileMigration() files_sdk.FileMigration {
	return i.Current().(files_sdk.FileMigration)
}

func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.FileMigrationFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

func (c *Client) List(params files_sdk.FileMigrationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/file_migrations", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.FileMigrationCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.FileMigrationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) Find(params files_sdk.FileMigrationFindParams, opts ...files_sdk.RequestResponseOption) (fileMigration files_sdk.FileMigration, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/file_migrations/{id}", Params: params, Entity: &fileMigration}, opts...)
	return
}

func Find(params files_sdk.FileMigrationFindParams, opts ...files_sdk.RequestResponseOption) (fileMigration files_sdk.FileMigration, err error) {
	return (&Client{}).Find(params, opts...)
}

func (c *Client) Wait(fileAction files_sdk.FileAction, status func(files_sdk.FileMigration), opts ...files_sdk.RequestResponseOption) (files_sdk.FileMigration, error) {
	var err error
	var migration files_sdk.FileMigration
	migration.Status = fileAction.Status
	migration.Id = fileAction.FileMigrationId
	if migration.Status == "completed" || migration.Status == "failed" || err != nil {
		return migration, nil
	}
	for {
		migration, err = c.Find(files_sdk.FileMigrationFindParams{Id: fileAction.FileMigrationId}, opts...)
		if err == nil {
			status(migration)
		}
		if migration.Status == "completed" || migration.Status == "failed" || err != nil {
			return migration, err
		}
		time.Sleep(time.Second * 1)
	}
}

func (c *Client) LogIterator(ctx context.Context, f files_sdk.FileMigration) files_sdk.IterI {
	return files_sdk.FilesMigrationLogIter{FileMigration: f, Context: ctx, Config: c.Config}.Init()
}
