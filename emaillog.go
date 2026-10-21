package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// EmailLog is a Files.com API resource.
type EmailLog struct {
	// Start Time of Action. Deprecrated: Use created_at.
	Timestamp *time.Time `json:"timestamp,omitempty" path:"timestamp,omitempty" url:"timestamp,omitempty"`
	// Log Message.
	Message string `json:"message,omitempty" path:"message,omitempty" url:"message,omitempty"`
	// Status of E-Mail delivery.
	Status string `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	// Subject line of E-Mail.
	Subject string `json:"subject,omitempty" path:"subject,omitempty" url:"subject,omitempty"`
	// To field of E-Mail.
	To string `json:"to,omitempty" path:"to,omitempty" url:"to,omitempty"`
	// CC field of E-Mail.
	Cc string `json:"cc,omitempty" path:"cc,omitempty" url:"cc,omitempty"`
	// How was the email delivered? `customer_smtp` or `files.com`.
	DeliveryMethod string `json:"delivery_method,omitempty" path:"delivery_method,omitempty" url:"delivery_method,omitempty"`
	// Customer SMTP Hostname used.
	SmtpHostname string `json:"smtp_hostname,omitempty" path:"smtp_hostname,omitempty" url:"smtp_hostname,omitempty"`
	// Customer SMTP IP address as resolved for use (useful for troubleshooting DNS issues with customer SMTP).
	SmtpIp string `json:"smtp_ip,omitempty" path:"smtp_ip,omitempty" url:"smtp_ip,omitempty"`
	// Start Time of Action.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
}

// Identifier no path or id

// EmailLogCollection is a list of EmailLog resources.
type EmailLogCollection []EmailLog

// EmailLogListParams contains the request parameters for GET /email_logs.
type EmailLogListParams struct {
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `status` and `created_at`. Valid field combinations are `[ status ]`, `[ created_at ]` or `[ status, created_at ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. Valid field combinations are `[ status ]`, `[ created_at ]` or `[ status, created_at ]`.
	FilterGt interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	// If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. Valid field combinations are `[ status ]`, `[ created_at ]` or `[ status, created_at ]`.
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	// If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. Valid field combinations are `[ status ]`, `[ created_at ]` or `[ status, created_at ]`.
	FilterLt interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	// If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. Valid field combinations are `[ status ]`, `[ created_at ]` or `[ status, created_at ]`.
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (e *EmailLog) UnmarshalJSON(data []byte) error {
	type emailLog EmailLog
	var v emailLog
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*e = EmailLog(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (e *EmailLogCollection) UnmarshalJSON(data []byte) error {
	type emailLogs EmailLogCollection
	var v emailLogs
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*e = EmailLogCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (e *EmailLogCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*e))
	for i, v := range *e {
		ret[i] = v
	}

	return &ret
}
