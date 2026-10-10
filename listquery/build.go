package listquery

import (
	"net/http"
	"slices"
	"sync"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/lib"
)

type List interface {
	UnmarshalJSON(data []byte) error
	ToSlice() *[]interface{}
}

func Build(config files_sdk.Config, path string, list List, opts ...files_sdk.RequestResponseOption) func(params lib.Values, laterOpts ...files_sdk.RequestResponseOption) (*[]interface{}, string, error) {
	// Reloaded iterators share this query, so calls can overlap. Each call
	// gets its own options slice, and responses decode into list one at a time.
	var decoding sync.Mutex
	return func(params lib.Values, laterOpts ...files_sdk.RequestResponseOption) (*[]interface{}, string, error) {
		defaultValue := make([]interface{}, 0)
		data, res, err := files_sdk.Call("GET", config, path, params, slices.Concat(opts, laterOpts)...)
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

		decoding.Lock()
		defer decoding.Unlock()
		if err := list.UnmarshalJSON(*data); err != nil {
			return &defaultValue, nextCursor(res.Header), err
		}
		return list.ToSlice(), nextCursor(res.Header), nil
	}
}

// nextCursor returns the cursor for the next page: the documented X-Files-Cursor-Next header, or the legacy
// X-Files-Cursor header when Next is absent or empty.
func nextCursor(header http.Header) string {
	if cursor := header.Get("X-Files-Cursor-Next"); cursor != "" {
		return cursor
	}
	return header.Get("X-Files-Cursor")
}
