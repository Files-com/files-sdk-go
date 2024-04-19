package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type SftpActionLog struct {
	Timestamp            *time.Time `json:"timestamp,omitempty" path:"timestamp,omitempty" url:"timestamp,omitempty"`
	RemoteIp             string     `json:"remote_ip,omitempty" path:"remote_ip,omitempty" url:"remote_ip,omitempty"`
	User                 string     `json:"user,omitempty" path:"user,omitempty" url:"user,omitempty"`
	SessionUid           string     `json:"session_uid,omitempty" path:"session_uid,omitempty" url:"session_uid,omitempty"`
	SeqId                int64      `json:"seq_id,omitempty" path:"seq_id,omitempty" url:"seq_id,omitempty"`
	AuthMethod           string     `json:"auth_method,omitempty" path:"auth_method,omitempty" url:"auth_method,omitempty"`
	AuthCiphers          string     `json:"auth_ciphers,omitempty" path:"auth_ciphers,omitempty" url:"auth_ciphers,omitempty"`
	ActionType           string     `json:"action_type,omitempty" path:"action_type,omitempty" url:"action_type,omitempty"`
	Path                 string     `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	TruePath             string     `json:"true_path,omitempty" path:"true_path,omitempty" url:"true_path,omitempty"`
	Name                 string     `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Message              string     `json:"message,omitempty" path:"message,omitempty" url:"message,omitempty"`
	FailureReasonType    string     `json:"failure_reason_type,omitempty" path:"failure_reason_type,omitempty" url:"failure_reason_type,omitempty"`
	FailureReasonMessage string     `json:"failure_reason_message,omitempty" path:"failure_reason_message,omitempty" url:"failure_reason_message,omitempty"`
	Md5                  string     `json:"md5,omitempty" path:"md5,omitempty" url:"md5,omitempty"`
	Flags                string     `json:"flags,omitempty" path:"flags,omitempty" url:"flags,omitempty"`
	Handle               string     `json:"handle,omitempty" path:"handle,omitempty" url:"handle,omitempty"`
	Attrs                string     `json:"attrs,omitempty" path:"attrs,omitempty" url:"attrs,omitempty"`
	Size                 string     `json:"size,omitempty" path:"size,omitempty" url:"size,omitempty"`
	Offset               string     `json:"offset,omitempty" path:"offset,omitempty" url:"offset,omitempty"`
	Length               string     `json:"length,omitempty" path:"length,omitempty" url:"length,omitempty"`
	DataLength           string     `json:"data_length,omitempty" path:"data_length,omitempty" url:"data_length,omitempty"`
	Success              string     `json:"success,omitempty" path:"success,omitempty" url:"success,omitempty"`
	DurationMs           int64      `json:"duration_ms,omitempty" path:"duration_ms,omitempty" url:"duration_ms,omitempty"`
}

func (s SftpActionLog) Identifier() interface{} {
	return s.Path
}

type SftpActionLogCollection []SftpActionLog

type SftpActionLogListParams struct {
	Filter       SftpActionLog          `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterPrefix map[string]interface{} `url:"filter_prefix,omitempty" required:"false" json:"filter_prefix,omitempty" path:"filter_prefix"`
	ListParams
}

func (s *SftpActionLog) UnmarshalJSON(data []byte) error {
	type sftpActionLog SftpActionLog
	var v sftpActionLog
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = SftpActionLog(v)
	return nil
}

func (s *SftpActionLogCollection) UnmarshalJSON(data []byte) error {
	type sftpActionLogs SftpActionLogCollection
	var v sftpActionLogs
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SftpActionLogCollection(v)
	return nil
}

func (s *SftpActionLogCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
