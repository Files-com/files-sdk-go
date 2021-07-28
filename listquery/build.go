package listquery

import (
	"context"
	"net/url"

	files_sdk "github.com/Files-com/files-sdk-go"
)

type List interface {
	UnmarshalJSON(data []byte) error
	ToSlice() *[]interface{}
}

type ExportParams interface {
	ExportParams() (url.Values, error)
}

func Build(ctx context.Context, i ExportParams, config files_sdk.Config, path string, list List) func() (*[]interface{}, string, error) {
	return func() (*[]interface{}, string, error) {
		defaultValue := make([]interface{}, 0)
		exportParams, err := i.ExportParams()
		if err != nil {
			return &defaultValue, "", err
		}
		data, res, err := files_sdk.Call(ctx, "GET", config, path, exportParams)
		defer func() {
			if res != nil && res.Body != nil {
				res.Body.Close()
			}
		}()
		if err != nil {
			return &defaultValue, "", err
		}
		if err := list.UnmarshalJSON(*data); err != nil {
			return &defaultValue, "", err
		}
		return list.ToSlice(), res.Header.Get("X-Files-Cursor"), nil
	}
}
