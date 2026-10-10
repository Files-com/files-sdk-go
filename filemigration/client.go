// Package file_migration provides the Files.com FileMigration API client.
package file_migration

import (
	"context"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Client calls the FileMigration API using its embedded Config.
type Client struct {
	files_sdk.Config
}

// Find calls GET /file_migrations/{id}.
func (c *Client) Find(params files_sdk.FileMigrationFindParams, opts ...files_sdk.RequestResponseOption) (fileMigration files_sdk.FileMigration, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/file_migrations/{id}", Params: params, Entity: &fileMigration}, opts...)
	return
}

// Find calls Client.Find using the default configuration.
func Find(params files_sdk.FileMigrationFindParams, opts ...files_sdk.RequestResponseOption) (fileMigration files_sdk.FileMigration, err error) {
	return (&Client{}).Find(params, opts...)
}

// Wait polls the migration once per second until it completes, fails, or a request
// returns an error. status is called after each successful poll and must be non-nil.
// WithContext controls cancellation. A failed migration is returned with nil error
// when the status request itself succeeds; inspect the returned Status.
func (c *Client) Wait(fileAction files_sdk.FileAction, status func(files_sdk.FileMigration), opts ...files_sdk.RequestResponseOption) (files_sdk.FileMigration, error) {
	var err error
	var migration files_sdk.FileMigration
	migration.Status = fileAction.Status
	migration.Id = fileAction.FileMigrationId
	ctx := files_sdk.ContextOption(opts)

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
		select {
		case <-ctx.Done():
			return migration, ctx.Err()
		case <-time.After(time.Second * 1):
		}
	}
}

// LogIterator returns an iterator over the migration log using ctx for cancellation.
func (c *Client) LogIterator(ctx context.Context, f files_sdk.FileMigration) files_sdk.IterI {
	return files_sdk.FilesMigrationLogIter{FileMigration: f, Context: ctx, Config: c.Config}.Init()
}
