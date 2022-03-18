package files_sdk

import (
	"encoding/json"
)

type WebhookTest struct {
	Code          int64           `json:"code,omitempty"`
	Message       string          `json:"message,omitempty"`
	Status        string          `json:"status,omitempty"`
	Data          Auto            `json:"data,omitempty"`
	Success       *bool           `json:"success,omitempty"`
	Url           string          `json:"url,omitempty"`
	Method        string          `json:"method,omitempty"`
	Encoding      string          `json:"encoding,omitempty"`
	Headers       json.RawMessage `json:"headers,omitempty"`
	Body          json.RawMessage `json:"body,omitempty"`
	RawBody       string          `json:"raw_body,omitempty"`
	FileAsBody    *bool           `json:"file_as_body,omitempty"`
	FileFormField string          `json:"file_form_field,omitempty"`
	Action        string          `json:"action,omitempty"`
}

type WebhookTestCollection []WebhookTest

type WebhookTestCreateParams struct {
	Url           string          `url:"url,omitempty" required:"true" json:"url,omitempty"`
	Method        string          `url:"method,omitempty" required:"false" json:"method,omitempty"`
	Encoding      string          `url:"encoding,omitempty" required:"false" json:"encoding,omitempty"`
	Headers       json.RawMessage `url:"headers,omitempty" required:"false" json:"headers,omitempty"`
	Body          json.RawMessage `url:"body,omitempty" required:"false" json:"body,omitempty"`
	RawBody       string          `url:"raw_body,omitempty" required:"false" json:"raw_body,omitempty"`
	FileAsBody    *bool           `url:"file_as_body,omitempty" required:"false" json:"file_as_body,omitempty"`
	FileFormField string          `url:"file_form_field,omitempty" required:"false" json:"file_form_field,omitempty"`
	Action        string          `url:"action,omitempty" required:"false" json:"action,omitempty"`
}

func (w *WebhookTest) UnmarshalJSON(data []byte) error {
	type webhookTest WebhookTest
	var v webhookTest
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*w = WebhookTest(v)
	return nil
}

func (w *WebhookTestCollection) UnmarshalJSON(data []byte) error {
	type webhookTests []WebhookTest
	var v webhookTests
	if err := json.Unmarshal(data, &v); err != nil {
		return err
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
