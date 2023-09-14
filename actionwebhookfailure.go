package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type ActionWebhookFailure struct {
}

// Identifier no path or id

type ActionWebhookFailureCollection []ActionWebhookFailure

// retry Action Webhook Failure
type ActionWebhookFailureRetryParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

func (a *ActionWebhookFailure) UnmarshalJSON(data []byte) error {
	type actionWebhookFailure ActionWebhookFailure
	var v actionWebhookFailure
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = ActionWebhookFailure(v)
	return nil
}

func (a *ActionWebhookFailureCollection) UnmarshalJSON(data []byte) error {
	type actionWebhookFailures ActionWebhookFailureCollection
	var v actionWebhookFailures
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = ActionWebhookFailureCollection(v)
	return nil
}

func (a *ActionWebhookFailureCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
