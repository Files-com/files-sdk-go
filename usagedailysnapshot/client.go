package usage_daily_snapshot

import (
  lib "github.com/Files-com/files-sdk-go/lib"
  files_sdk "github.com/Files-com/files-sdk-go"
)

type Client struct {
	files_sdk.Config
}

type Iter struct {
	*lib.Iter
}

func (i *Iter) UsageDailySnapshot() files_sdk.UsageDailySnapshot {
	return i.Current().(files_sdk.UsageDailySnapshot)
}

func (c *Client) List(params files_sdk.UsageDailySnapshotListParams) *Iter {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/usage_daily_snapshots"

	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, i.ExportParams())
		defaultValue := make([]interface{}, 0)
        if err != nil {
          return &defaultValue, "", err
        }
		list := files_sdk.UsageDailySnapshotCollection{}
		if err := list.UnmarshalJSON(*data); err != nil {
          return &defaultValue, "", err
        }

		ret := make([]interface{}, len(list))
		for i, v := range list {
			ret[i] = v
		}
		cursor := res.Header.Get("X-Files-Cursor")
		return &ret, cursor, nil
	}
	i.ListParams = &params
	return i
}

func List(params files_sdk.UsageDailySnapshotListParams) *Iter {
  client := Client{}
  return client.List (params)
}
