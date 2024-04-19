package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type ExavaultApiRequestLog struct {
	Timestamp     *time.Time `json:"timestamp,omitempty" path:"timestamp,omitempty" url:"timestamp,omitempty"`
	Endpoint      string     `json:"endpoint,omitempty" path:"endpoint,omitempty" url:"endpoint,omitempty"`
	Version       string     `json:"version,omitempty" path:"version,omitempty" url:"version,omitempty"`
	RequestIp     string     `json:"request_ip,omitempty" path:"request_ip,omitempty" url:"request_ip,omitempty"`
	RequestMethod string     `json:"request_method,omitempty" path:"request_method,omitempty" url:"request_method,omitempty"`
	ErrorType     string     `json:"error_type,omitempty" path:"error_type,omitempty" url:"error_type,omitempty"`
	ErrorMessage  string     `json:"error_message,omitempty" path:"error_message,omitempty" url:"error_message,omitempty"`
	UserAgent     string     `json:"user_agent,omitempty" path:"user_agent,omitempty" url:"user_agent,omitempty"`
	ResponseCode  int64      `json:"response_code,omitempty" path:"response_code,omitempty" url:"response_code,omitempty"`
	DurationMs    int64      `json:"duration_ms,omitempty" path:"duration_ms,omitempty" url:"duration_ms,omitempty"`
}

// Identifier no path or id

type ExavaultApiRequestLogCollection []ExavaultApiRequestLog

type ExavaultApiRequestLogListParams struct {
	Filter       ExavaultApiRequestLog  `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterPrefix map[string]interface{} `url:"filter_prefix,omitempty" required:"false" json:"filter_prefix,omitempty" path:"filter_prefix"`
	ListParams
}

func (e *ExavaultApiRequestLog) UnmarshalJSON(data []byte) error {
	type exavaultApiRequestLog ExavaultApiRequestLog
	var v exavaultApiRequestLog
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*e = ExavaultApiRequestLog(v)
	return nil
}

func (e *ExavaultApiRequestLogCollection) UnmarshalJSON(data []byte) error {
	type exavaultApiRequestLogs ExavaultApiRequestLogCollection
	var v exavaultApiRequestLogs
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*e = ExavaultApiRequestLogCollection(v)
	return nil
}

func (e *ExavaultApiRequestLogCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*e))
	for i, v := range *e {
		ret[i] = v
	}

	return &ret
}
