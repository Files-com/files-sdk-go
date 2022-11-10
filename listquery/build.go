package listquery

import (
	"context"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/lib"
)

type List interface {
	UnmarshalJSON(data []byte) error
	ToSlice() *[]interface{}
}

func Build(ctx context.Context, config files_sdk.Config, path string, list List) func(params lib.Values) (*[]interface{}, string, error) {
	return func(params lib.Values) (*[]interface{}, string, error) {
		defaultValue := make([]interface{}, 0)
		data, res, err := files_sdk.Call(ctx, "GET", config, path, params)
		defer func() {
			if res != nil && res.Body != nil {
				res.Body.Close()
			}
		}()
		if err != nil {
			return &defaultValue, "", err
		}

		if err := lib.ResponseErrors(res, lib.NonOkError, lib.NonJSONError); err != nil {
			return &defaultValue, "", err
		}

		if err := list.UnmarshalJSON(*data); err != nil {
			return &defaultValue, res.Header.Get("X-Files-Cursor"), err
		}
		return list.ToSlice(), res.Header.Get("X-Files-Cursor"), nil
	}
}
