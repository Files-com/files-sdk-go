package files_sdk

import (
	"encoding/json"
)

type ActionWebhookFailure struct {
}

type ActionWebhookFailureCollection []ActionWebhookFailure

// retry Action Webhook Failure
type ActionWebhookFailureRetryParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

func (a *ActionWebhookFailure) UnmarshalJSON(data []byte) error {
	type actionWebhookFailure ActionWebhookFailure
	var v actionWebhookFailure
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = ActionWebhookFailure(v)
	return nil
}

func (a *ActionWebhookFailureCollection) UnmarshalJSON(data []byte) error {
	type actionWebhookFailures []ActionWebhookFailure
	var v actionWebhookFailures
	if err := json.Unmarshal(data, &v); err != nil {
		return err
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
