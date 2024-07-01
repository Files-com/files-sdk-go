package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type FtpActionLog struct {
	Timestamp       *time.Time `json:"timestamp,omitempty" path:"timestamp,omitempty" url:"timestamp,omitempty"`
	RemoteIp        string     `json:"remote_ip,omitempty" path:"remote_ip,omitempty" url:"remote_ip,omitempty"`
	ServerIp        string     `json:"server_ip,omitempty" path:"server_ip,omitempty" url:"server_ip,omitempty"`
	Username        string     `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	SessionUuid     string     `json:"session_uuid,omitempty" path:"session_uuid,omitempty" url:"session_uuid,omitempty"`
	SeqId           int64      `json:"seq_id,omitempty" path:"seq_id,omitempty" url:"seq_id,omitempty"`
	AuthCiphers     string     `json:"auth_ciphers,omitempty" path:"auth_ciphers,omitempty" url:"auth_ciphers,omitempty"`
	ActionType      string     `json:"action_type,omitempty" path:"action_type,omitempty" url:"action_type,omitempty"`
	Path            string     `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	TruePath        string     `json:"true_path,omitempty" path:"true_path,omitempty" url:"true_path,omitempty"`
	Name            string     `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Cmd             string     `json:"cmd,omitempty" path:"cmd,omitempty" url:"cmd,omitempty"`
	Param           string     `json:"param,omitempty" path:"param,omitempty" url:"param,omitempty"`
	ResponseCode    int64      `json:"responseCode,omitempty" path:"responseCode,omitempty" url:"responseCode,omitempty"`
	ResponseMessage string     `json:"responseMessage,omitempty" path:"responseMessage,omitempty" url:"responseMessage,omitempty"`
	EntriesReturned int64      `json:"entries_returned,omitempty" path:"entries_returned,omitempty" url:"entries_returned,omitempty"`
	Success         *bool      `json:"success,omitempty" path:"success,omitempty" url:"success,omitempty"`
	DurationMs      int64      `json:"duration_ms,omitempty" path:"duration_ms,omitempty" url:"duration_ms,omitempty"`
}

func (f FtpActionLog) Identifier() interface{} {
	return f.Path
}

type FtpActionLogCollection []FtpActionLog

type FtpActionLogListParams struct {
	Action       string                 `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	Filter       FtpActionLog           `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterPrefix map[string]interface{} `url:"filter_prefix,omitempty" required:"false" json:"filter_prefix,omitempty" path:"filter_prefix"`
	ListParams
}

func (f *FtpActionLog) UnmarshalJSON(data []byte) error {
	type ftpActionLog FtpActionLog
	var v ftpActionLog
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*f = FtpActionLog(v)
	return nil
}

func (f *FtpActionLogCollection) UnmarshalJSON(data []byte) error {
	type ftpActionLogs FtpActionLogCollection
	var v ftpActionLogs
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*f = FtpActionLogCollection(v)
	return nil
}

func (f *FtpActionLogCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*f))
	for i, v := range *f {
		ret[i] = v
	}

	return &ret
}
