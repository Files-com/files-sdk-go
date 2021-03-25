package files_sdk

import (
	"encoding/json"
)

type WebhookTest struct {
	Code     int             `json:"code,omitempty"`
	Message  string          `json:"message,omitempty"`
	Status   string          `json:"status,omitempty"`
	Data     string          `json:"data,omitempty"`
	Success  *bool           `json:"success,omitempty"`
	Url      string          `json:"url,omitempty"`
	Method   string          `json:"method,omitempty"`
	Encoding string          `json:"encoding,omitempty"`
	Headers  json.RawMessage `json:"headers,omitempty"`
	Body     json.RawMessage `json:"body,omitempty"`
	Action   string          `json:"action,omitempty"`
}

type WebhookTestCollection []WebhookTest

type WebhookTestCreateParams struct {
	Url      string          `url:"url,omitempty" required:"true"`
	Method   string          `url:"method,omitempty" required:"false"`
	Encoding string          `url:"encoding,omitempty" required:"false"`
	Headers  json.RawMessage `url:"headers,omitempty" required:"false"`
	Body     json.RawMessage `url:"body,omitempty" required:"false"`
	Action   string          `url:"action,omitempty" required:"false"`
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
