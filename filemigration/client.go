package file_migration

import (
	"context"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Find(ctx context.Context, params files_sdk.FileMigrationFindParams) (fileMigration files_sdk.FileMigration, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/file_migrations/{id}", Params: params, Entity: &fileMigration})
	return
}

func Find(ctx context.Context, params files_sdk.FileMigrationFindParams) (fileMigration files_sdk.FileMigration, err error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Wait(ctx context.Context, fileAction files_sdk.FileAction, status func(files_sdk.FileMigration)) (files_sdk.FileMigration, error) {
	var err error
	var migration files_sdk.FileMigration
	migration.Status = fileAction.Status
	migration.Id = fileAction.FileMigrationId
	if migration.Status == "completed" || migration.Status == "failed" || err != nil {
		return migration, nil
	}
	for {
		migration, err = c.Find(ctx, files_sdk.FileMigrationFindParams{Id: fileAction.FileMigrationId})
		if err == nil {
			status(migration)
		}
		if migration.Status == "completed" || migration.Status == "failed" || err != nil {
			return migration, err
		}
		time.Sleep(time.Second * 1)
	}
}

func (c *Client) LogIterator(ctx context.Context, f files_sdk.FileMigration) lib.IterI {
	return files_sdk.FilesMigrationLogIter{FileMigration: f, Context: ctx, Config: c.Config}.Init()
}
