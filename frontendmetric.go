package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type FrontendMetric struct {
	MetricType string `json:"metric_type,omitempty" path:"metric_type,omitempty" url:"metric_type,omitempty"`
	Subkey     string `json:"subkey,omitempty" path:"subkey,omitempty" url:"subkey,omitempty"`
	Ms         int64  `json:"ms,omitempty" path:"ms,omitempty" url:"ms,omitempty"`
}

// Identifier no path or id

type FrontendMetricCollection []FrontendMetric

type FrontendMetricMetricTypeEnum string

func (u FrontendMetricMetricTypeEnum) String() string {
	return string(u)
}

func (u FrontendMetricMetricTypeEnum) Enum() map[string]FrontendMetricMetricTypeEnum {
	return map[string]FrontendMetricMetricTypeEnum{
		"increment": FrontendMetricMetricTypeEnum("increment"),
		"timing":    FrontendMetricMetricTypeEnum("timing"),
	}
}

type FrontendMetricCreateParams struct {
	MetricType FrontendMetricMetricTypeEnum `url:"metric_type,omitempty" required:"true" json:"metric_type,omitempty" path:"metric_type"`
	Subkey     string                       `url:"subkey,omitempty" required:"true" json:"subkey,omitempty" path:"subkey"`
	Ms         int64                        `url:"ms,omitempty" required:"false" json:"ms,omitempty" path:"ms"`
}

func (f *FrontendMetric) UnmarshalJSON(data []byte) error {
	type frontendMetric FrontendMetric
	var v frontendMetric
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*f = FrontendMetric(v)
	return nil
}

func (f *FrontendMetricCollection) UnmarshalJSON(data []byte) error {
	type frontendMetrics FrontendMetricCollection
	var v frontendMetrics
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*f = FrontendMetricCollection(v)
	return nil
}

func (f *FrontendMetricCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*f))
	for i, v := range *f {
		ret[i] = v
	}

	return &ret
}
