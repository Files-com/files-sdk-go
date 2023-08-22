package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type MonitoringStats struct {
	Alerts string `json:"alerts,omitempty" path:"alerts,omitempty" url:"alerts,omitempty"`
	Info   string `json:"info,omitempty" path:"info,omitempty" url:"info,omitempty"`
}

// Identifier no path or id

type MonitoringStatsCollection []MonitoringStats

func (m *MonitoringStats) UnmarshalJSON(data []byte) error {
	type monitoringStats MonitoringStats
	var v monitoringStats
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*m = MonitoringStats(v)
	return nil
}

func (m *MonitoringStatsCollection) UnmarshalJSON(data []byte) error {
	type monitoringStatss MonitoringStatsCollection
	var v monitoringStatss
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*m = MonitoringStatsCollection(v)
	return nil
}

func (m *MonitoringStatsCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*m))
	for i, v := range *m {
		ret[i] = v
	}

	return &ret
}
