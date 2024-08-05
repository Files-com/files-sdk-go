package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type PublicHostingRequestLog struct {
	Timestamp    *time.Time `json:"timestamp,omitempty" path:"timestamp,omitempty" url:"timestamp,omitempty"`
	RemoteIp     string     `json:"remote_ip,omitempty" path:"remote_ip,omitempty" url:"remote_ip,omitempty"`
	ServerIp     string     `json:"server_ip,omitempty" path:"server_ip,omitempty" url:"server_ip,omitempty"`
	Hostname     string     `json:"hostname,omitempty" path:"hostname,omitempty" url:"hostname,omitempty"`
	Path         string     `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	ResponseCode int64      `json:"responseCode,omitempty" path:"responseCode,omitempty" url:"responseCode,omitempty"`
	Success      *bool      `json:"success,omitempty" path:"success,omitempty" url:"success,omitempty"`
	DurationMs   int64      `json:"duration_ms,omitempty" path:"duration_ms,omitempty" url:"duration_ms,omitempty"`
}

func (p PublicHostingRequestLog) Identifier() interface{} {
	return p.Path
}

type PublicHostingRequestLogCollection []PublicHostingRequestLog

type PublicHostingRequestLogListParams struct {
	Filter       PublicHostingRequestLog `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterPrefix map[string]interface{}  `url:"filter_prefix,omitempty" required:"false" json:"filter_prefix,omitempty" path:"filter_prefix"`
	ListParams
}

func (p *PublicHostingRequestLog) UnmarshalJSON(data []byte) error {
	type publicHostingRequestLog PublicHostingRequestLog
	var v publicHostingRequestLog
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*p = PublicHostingRequestLog(v)
	return nil
}

func (p *PublicHostingRequestLogCollection) UnmarshalJSON(data []byte) error {
	type publicHostingRequestLogs PublicHostingRequestLogCollection
	var v publicHostingRequestLogs
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*p = PublicHostingRequestLogCollection(v)
	return nil
}

func (p *PublicHostingRequestLogCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
