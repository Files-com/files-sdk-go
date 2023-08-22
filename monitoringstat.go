package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type MonitoringStat struct {
	Alerts string `json:"alerts,omitempty" path:"alerts,omitempty" url:"alerts,omitempty"`
	Info   string `json:"info,omitempty" path:"info,omitempty" url:"info,omitempty"`
}

// Identifier no path or id

type MonitoringStatCollection []MonitoringStat

type MonitoringStatListParams struct {
	Action string `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	ListParams
}

func (m *MonitoringStat) UnmarshalJSON(data []byte) error {
	type monitoringStat MonitoringStat
	var v monitoringStat
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*m = MonitoringStat(v)
	return nil
}

func (m *MonitoringStatCollection) UnmarshalJSON(data []byte) error {
	type monitoringStats MonitoringStatCollection
	var v monitoringStats
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*m = MonitoringStatCollection(v)
	return nil
}

func (m *MonitoringStatCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*m))
	for i, v := range *m {
		ret[i] = v
	}

	return &ret
}
