package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type EmailLog struct {
	Timestamp      *time.Time `json:"timestamp,omitempty" path:"timestamp,omitempty" url:"timestamp,omitempty"`
	Message        string     `json:"message,omitempty" path:"message,omitempty" url:"message,omitempty"`
	Status         string     `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	Subject        string     `json:"subject,omitempty" path:"subject,omitempty" url:"subject,omitempty"`
	To             string     `json:"to,omitempty" path:"to,omitempty" url:"to,omitempty"`
	Cc             string     `json:"cc,omitempty" path:"cc,omitempty" url:"cc,omitempty"`
	DeliveryMethod string     `json:"delivery_method,omitempty" path:"delivery_method,omitempty" url:"delivery_method,omitempty"`
	SmtpHostname   string     `json:"smtp_hostname,omitempty" path:"smtp_hostname,omitempty" url:"smtp_hostname,omitempty"`
	SmtpIp         string     `json:"smtp_ip,omitempty" path:"smtp_ip,omitempty" url:"smtp_ip,omitempty"`
}

// Identifier no path or id

type EmailLogCollection []EmailLog

type EmailLogListParams struct {
	Action       string                 `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	Filter       EmailLog               `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterPrefix map[string]interface{} `url:"filter_prefix,omitempty" required:"false" json:"filter_prefix,omitempty" path:"filter_prefix"`
	ListParams
}

func (e *EmailLog) UnmarshalJSON(data []byte) error {
	type emailLog EmailLog
	var v emailLog
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*e = EmailLog(v)
	return nil
}

func (e *EmailLogCollection) UnmarshalJSON(data []byte) error {
	type emailLogs EmailLogCollection
	var v emailLogs
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*e = EmailLogCollection(v)
	return nil
}

func (e *EmailLogCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*e))
	for i, v := range *e {
		ret[i] = v
	}

	return &ret
}
