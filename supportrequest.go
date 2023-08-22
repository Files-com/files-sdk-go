package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/lpar/date"
)

type SupportRequest struct {
	Id                    int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Subject               string     `json:"subject,omitempty" path:"subject,omitempty" url:"subject,omitempty"`
	Comment               string     `json:"comment,omitempty" path:"comment,omitempty" url:"comment,omitempty"`
	CreatedAt             *date.Date `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	AccessUntil           *date.Date `json:"access_until,omitempty" path:"access_until,omitempty" url:"access_until,omitempty"`
	CustomerSuccessAccess string     `json:"customer_success_access,omitempty" path:"customer_success_access,omitempty" url:"customer_success_access,omitempty"`
	Priority              string     `json:"priority,omitempty" path:"priority,omitempty" url:"priority,omitempty"`
	Name                  string     `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	PhoneNumber           string     `json:"phone_number,omitempty" path:"phone_number,omitempty" url:"phone_number,omitempty"`
	AccessReset           *bool      `json:"access_reset,omitempty" path:"access_reset,omitempty" url:"access_reset,omitempty"`
	Email                 string     `json:"email,omitempty" path:"email,omitempty" url:"email,omitempty"`
	AttachmentsFiles      []string   `json:"attachments_files,omitempty" path:"attachments_files,omitempty" url:"attachments_files,omitempty"`
}

func (s SupportRequest) Identifier() interface{} {
	return s.Id
}

type SupportRequestCollection []SupportRequest

type SupportRequestCustomerSuccessAccessEnum string

func (u SupportRequestCustomerSuccessAccessEnum) String() string {
	return string(u)
}

func (u SupportRequestCustomerSuccessAccessEnum) Enum() map[string]SupportRequestCustomerSuccessAccessEnum {
	return map[string]SupportRequestCustomerSuccessAccessEnum{
		"no":                           SupportRequestCustomerSuccessAccessEnum("no"),
		"readonly":                     SupportRequestCustomerSuccessAccessEnum("readonly"),
		"full":                         SupportRequestCustomerSuccessAccessEnum("full"),
		"full_plus_remote_credentials": SupportRequestCustomerSuccessAccessEnum("full_plus_remote_credentials"),
	}
}

type SupportRequestPriorityEnum string

func (u SupportRequestPriorityEnum) String() string {
	return string(u)
}

func (u SupportRequestPriorityEnum) Enum() map[string]SupportRequestPriorityEnum {
	return map[string]SupportRequestPriorityEnum{
		"low":      SupportRequestPriorityEnum("low"),
		"normal":   SupportRequestPriorityEnum("normal"),
		"high":     SupportRequestPriorityEnum("high"),
		"urgent":   SupportRequestPriorityEnum("urgent"),
		"critical": SupportRequestPriorityEnum("critical"),
	}
}

type SupportRequestListParams struct {
	Action string                 `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	SortBy map[string]interface{} `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	ListParams
}

type SupportRequestCreateParams struct {
	CustomerSuccessAccess SupportRequestCustomerSuccessAccessEnum `url:"customer_success_access,omitempty" required:"false" json:"customer_success_access,omitempty" path:"customer_success_access"`
	AccessReset           *bool                                   `url:"access_reset,omitempty" required:"false" json:"access_reset,omitempty" path:"access_reset"`
	Email                 string                                  `url:"email,omitempty" required:"true" json:"email,omitempty" path:"email"`
	Subject               string                                  `url:"subject,omitempty" required:"true" json:"subject,omitempty" path:"subject"`
	Comment               string                                  `url:"comment,omitempty" required:"true" json:"comment,omitempty" path:"comment"`
	Priority              SupportRequestPriorityEnum              `url:"priority,omitempty" required:"false" json:"priority,omitempty" path:"priority"`
	PhoneNumber           string                                  `url:"phone_number,omitempty" required:"false" json:"phone_number,omitempty" path:"phone_number"`
	Name                  string                                  `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	AttachmentsFiles      []string                                `url:"attachments_files,omitempty" required:"false" json:"attachments_files,omitempty" path:"attachments_files"`
}

type SupportRequestUpdateParams struct {
	Id                    int64                                   `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
	CustomerSuccessAccess SupportRequestCustomerSuccessAccessEnum `url:"customer_success_access,omitempty" required:"false" json:"customer_success_access,omitempty" path:"customer_success_access"`
	AccessReset           *bool                                   `url:"access_reset,omitempty" required:"false" json:"access_reset,omitempty" path:"access_reset"`
}

func (s *SupportRequest) UnmarshalJSON(data []byte) error {
	type supportRequest SupportRequest
	var v supportRequest
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = SupportRequest(v)
	return nil
}

func (s *SupportRequestCollection) UnmarshalJSON(data []byte) error {
	type supportRequests SupportRequestCollection
	var v supportRequests
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SupportRequestCollection(v)
	return nil
}

func (s *SupportRequestCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
