package files_sdk

import (
	"encoding/json"
	"io"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type CrashReport struct {
	Id              int64     `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Build           string    `json:"build,omitempty" path:"build,omitempty" url:"build,omitempty"`
	Platform        string    `json:"platform,omitempty" path:"platform,omitempty" url:"platform,omitempty"`
	ProductName     string    `json:"product_name,omitempty" path:"product_name,omitempty" url:"product_name,omitempty"`
	Version         string    `json:"version,omitempty" path:"version,omitempty" url:"version,omitempty"`
	Comment         string    `json:"comment,omitempty" path:"comment,omitempty" url:"comment,omitempty"`
	Email           string    `json:"email,omitempty" path:"email,omitempty" url:"email,omitempty"`
	PlatformVersion string    `json:"platform_version,omitempty" path:"platform_version,omitempty" url:"platform_version,omitempty"`
	ReleaseChannel  string    `json:"release_channel,omitempty" path:"release_channel,omitempty" url:"release_channel,omitempty"`
	DumpFile        io.Reader `json:"dump_file,omitempty" path:"dump_file,omitempty" url:"dump_file,omitempty"`
	LogFile         io.Reader `json:"log_file,omitempty" path:"log_file,omitempty" url:"log_file,omitempty"`
}

func (c CrashReport) Identifier() interface{} {
	return c.Id
}

type CrashReportCollection []CrashReport

type CrashReportCreateParams struct {
	Build           string    `url:"build,omitempty" required:"true" json:"build,omitempty" path:"build"`
	Platform        string    `url:"platform,omitempty" required:"true" json:"platform,omitempty" path:"platform"`
	ProductName     string    `url:"product_name,omitempty" required:"true" json:"product_name,omitempty" path:"product_name"`
	Version         string    `url:"version,omitempty" required:"true" json:"version,omitempty" path:"version"`
	Comment         string    `url:"comment,omitempty" required:"false" json:"comment,omitempty" path:"comment"`
	Email           string    `url:"email,omitempty" required:"false" json:"email,omitempty" path:"email"`
	PlatformVersion string    `url:"platform_version,omitempty" required:"false" json:"platform_version,omitempty" path:"platform_version"`
	ReleaseChannel  string    `url:"release_channel,omitempty" required:"false" json:"release_channel,omitempty" path:"release_channel"`
	DumpFile        io.Writer `url:"dump_file,omitempty" required:"false" json:"dump_file,omitempty" path:"dump_file"`
	LogFile         io.Writer `url:"log_file,omitempty" required:"false" json:"log_file,omitempty" path:"log_file"`
}

func (c *CrashReport) UnmarshalJSON(data []byte) error {
	type crashReport CrashReport
	var v crashReport
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*c = CrashReport(v)
	return nil
}

func (c *CrashReportCollection) UnmarshalJSON(data []byte) error {
	type crashReports CrashReportCollection
	var v crashReports
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*c = CrashReportCollection(v)
	return nil
}

func (c *CrashReportCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*c))
	for i, v := range *c {
		ret[i] = v
	}

	return &ret
}
