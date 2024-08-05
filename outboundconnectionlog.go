package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type OutboundConnectionLog struct {
	Timestamp          *time.Time `json:"timestamp,omitempty" path:"timestamp,omitempty" url:"timestamp,omitempty"`
	Path               string     `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	ClientIp           string     `json:"client_ip,omitempty" path:"client_ip,omitempty" url:"client_ip,omitempty"`
	SrcRemoteServerId  string     `json:"src_remote_server_id,omitempty" path:"src_remote_server_id,omitempty" url:"src_remote_server_id,omitempty"`
	DestRemoteServerId string     `json:"dest_remote_server_id,omitempty" path:"dest_remote_server_id,omitempty" url:"dest_remote_server_id,omitempty"`
	Operation          string     `json:"operation,omitempty" path:"operation,omitempty" url:"operation,omitempty"`
	ErrorMessage       string     `json:"error_message,omitempty" path:"error_message,omitempty" url:"error_message,omitempty"`
	ErrorOperation     string     `json:"error_operation,omitempty" path:"error_operation,omitempty" url:"error_operation,omitempty"`
	ErrorType          string     `json:"error_type,omitempty" path:"error_type,omitempty" url:"error_type,omitempty"`
	Status             string     `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	DurationMs         int64      `json:"duration_ms,omitempty" path:"duration_ms,omitempty" url:"duration_ms,omitempty"`
	BytesUploaded      int64      `json:"bytes_uploaded,omitempty" path:"bytes_uploaded,omitempty" url:"bytes_uploaded,omitempty"`
	BytesDownloaded    int64      `json:"bytes_downloaded,omitempty" path:"bytes_downloaded,omitempty" url:"bytes_downloaded,omitempty"`
	ListCount          int64      `json:"list_count,omitempty" path:"list_count,omitempty" url:"list_count,omitempty"`
}

func (o OutboundConnectionLog) Identifier() interface{} {
	return o.Path
}

type OutboundConnectionLogCollection []OutboundConnectionLog

type OutboundConnectionLogListParams struct {
	Filter       OutboundConnectionLog  `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterPrefix map[string]interface{} `url:"filter_prefix,omitempty" required:"false" json:"filter_prefix,omitempty" path:"filter_prefix"`
	ListParams
}

func (o *OutboundConnectionLog) UnmarshalJSON(data []byte) error {
	type outboundConnectionLog OutboundConnectionLog
	var v outboundConnectionLog
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*o = OutboundConnectionLog(v)
	return nil
}

func (o *OutboundConnectionLogCollection) UnmarshalJSON(data []byte) error {
	type outboundConnectionLogs OutboundConnectionLogCollection
	var v outboundConnectionLogs
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*o = OutboundConnectionLogCollection(v)
	return nil
}

func (o *OutboundConnectionLogCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*o))
	for i, v := range *o {
		ret[i] = v
	}

	return &ret
}
