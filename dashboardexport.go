package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type DashboardExport struct {
	Id           int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	StartAt      *time.Time `json:"start_at,omitempty" path:"start_at,omitempty" url:"start_at,omitempty"`
	EndAt        *time.Time `json:"end_at,omitempty" path:"end_at,omitempty" url:"end_at,omitempty"`
	DeliveredAt  *time.Time `json:"delivered_at,omitempty" path:"delivered_at,omitempty" url:"delivered_at,omitempty"`
	ExportStatus string     `json:"export_status,omitempty" path:"export_status,omitempty" url:"export_status,omitempty"`
	Resolution   int64      `json:"resolution,omitempty" path:"resolution,omitempty" url:"resolution,omitempty"`
	Series       []string   `json:"series,omitempty" path:"series,omitempty" url:"series,omitempty"`
	UserId       int64      `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
}

func (d DashboardExport) Identifier() interface{} {
	return d.Id
}

type DashboardExportCollection []DashboardExport

type DashboardExportFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type DashboardExportCreateParams struct {
	UserId  int64                    `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	StartAt *time.Time               `url:"start_at,omitempty" required:"true" json:"start_at,omitempty" path:"start_at"`
	EndAt   *time.Time               `url:"end_at,omitempty" required:"true" json:"end_at,omitempty" path:"end_at"`
	Series  []map[string]interface{} `url:"series,omitempty" required:"true" json:"series,omitempty" path:"series"`
}

func (d *DashboardExport) UnmarshalJSON(data []byte) error {
	type dashboardExport DashboardExport
	var v dashboardExport
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*d = DashboardExport(v)
	return nil
}

func (d *DashboardExportCollection) UnmarshalJSON(data []byte) error {
	type dashboardExports DashboardExportCollection
	var v dashboardExports
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*d = DashboardExportCollection(v)
	return nil
}

func (d *DashboardExportCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*d))
	for i, v := range *d {
		ret[i] = v
	}

	return &ret
}
