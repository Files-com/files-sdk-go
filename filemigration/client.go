package file_migration

import (
	"context"
	"strconv"

	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Find(ctx context.Context, params files_sdk.FileMigrationFindParams) (files_sdk.FileMigration, error) {
	fileMigration := files_sdk.FileMigration{}
	if params.Id == 0 {
		return fileMigration, lib.CreateError(params, "Id")
	}
	path := "/file_migrations/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return fileMigration, err
	}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return fileMigration, err
	}
	if res.StatusCode == 204 {
		return fileMigration, nil
	}
	if err := fileMigration.UnmarshalJSON(*data); err != nil {
		return fileMigration, err
	}

	return fileMigration, nil
}

func Find(ctx context.Context, params files_sdk.FileMigrationFindParams) (files_sdk.FileMigration, error) {
	return (&Client{}).Find(ctx, params)
}
