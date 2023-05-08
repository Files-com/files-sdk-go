package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type WebhookTest struct {
	Code          int64           `json:"code,omitempty" path:"code"`
	Message       string          `json:"message,omitempty" path:"message"`
	Status        string          `json:"status,omitempty" path:"status"`
	Data          Auto            `json:"data,omitempty" path:"data"`
	Success       *bool           `json:"success,omitempty" path:"success"`
	Url           string          `json:"url,omitempty" path:"url"`
	Method        string          `json:"method,omitempty" path:"method"`
	Encoding      string          `json:"encoding,omitempty" path:"encoding"`
	Headers       json.RawMessage `json:"headers,omitempty" path:"headers"`
	Body          json.RawMessage `json:"body,omitempty" path:"body"`
	RawBody       string          `json:"raw_body,omitempty" path:"raw_body"`
	FileAsBody    *bool           `json:"file_as_body,omitempty" path:"file_as_body"`
	FileFormField string          `json:"file_form_field,omitempty" path:"file_form_field"`
	Action        string          `json:"action,omitempty" path:"action"`
}

// Identifier no path or id

type WebhookTestCollection []WebhookTest

type WebhookTestCreateParams struct {
	Url           string          `url:"url,omitempty" required:"true" json:"url,omitempty" path:"url"`
	Method        string          `url:"method,omitempty" required:"false" json:"method,omitempty" path:"method"`
	Encoding      string          `url:"encoding,omitempty" required:"false" json:"encoding,omitempty" path:"encoding"`
	Headers       json.RawMessage `url:"headers,omitempty" required:"false" json:"headers,omitempty" path:"headers"`
	Body          json.RawMessage `url:"body,omitempty" required:"false" json:"body,omitempty" path:"body"`
	RawBody       string          `url:"raw_body,omitempty" required:"false" json:"raw_body,omitempty" path:"raw_body"`
	FileAsBody    *bool           `url:"file_as_body,omitempty" required:"false" json:"file_as_body,omitempty" path:"file_as_body"`
	FileFormField string          `url:"file_form_field,omitempty" required:"false" json:"file_form_field,omitempty" path:"file_form_field"`
	Action        string          `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
}

func (w *WebhookTest) UnmarshalJSON(data []byte) error {
	type webhookTest WebhookTest
	var v webhookTest
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*w = WebhookTest(v)
	return nil
}

func (w *WebhookTestCollection) UnmarshalJSON(data []byte) error {
	type webhookTests WebhookTestCollection
	var v webhookTests
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*w = WebhookTestCollection(v)
	return nil
}

func (w *WebhookTestCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*w))
	for i, v := range *w {
		ret[i] = v
	}

	return &ret
}
