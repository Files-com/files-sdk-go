package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type EmailFeedback struct {
}

// Identifier no path or id

type EmailFeedbackCollection []EmailFeedback

type FeedbackParam struct {
	Email  string `url:"email,omitempty" json:"email,omitempty" path:"email"`
	Reason string `url:"reason,omitempty" json:"reason,omitempty" path:"reason"`
}

type EmailFeedbackCreateParams struct {
	FeedbackParam FeedbackParam `url:"feedback,omitempty" required:"false" json:"feedback,omitempty" path:"feedback"`
}

func (e *EmailFeedback) UnmarshalJSON(data []byte) error {
	type emailFeedback EmailFeedback
	var v emailFeedback
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*e = EmailFeedback(v)
	return nil
}

func (e *EmailFeedbackCollection) UnmarshalJSON(data []byte) error {
	type emailFeedbacks EmailFeedbackCollection
	var v emailFeedbacks
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*e = EmailFeedbackCollection(v)
	return nil
}

func (e *EmailFeedbackCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*e))
	for i, v := range *e {
		ret[i] = v
	}

	return &ret
}
