package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// WebhookTest is a Files.com API resource.
type WebhookTest struct {
	// Status HTTP code.
	Code int64 `json:"code,omitempty" path:"code,omitempty" url:"code,omitempty"`
	// Error message.
	Message string `json:"message,omitempty" path:"message,omitempty" url:"message,omitempty"`
	// Status message.
	Status string `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	// Additional data.
	Data Auto `json:"data,omitempty" path:"data,omitempty" url:"data,omitempty"`
	// The success status of the webhook test.
	Success *bool `json:"success,omitempty" path:"success,omitempty" url:"success,omitempty"`
	// URL for testing the webhook.
	Url string `json:"url,omitempty" path:"url,omitempty" url:"url,omitempty"`
	// HTTP method(GET or POST).
	Method string `json:"method,omitempty" path:"method,omitempty" url:"method,omitempty"`
	// HTTP encoding method. Can be JSON, XML, or RAW (form data).
	Encoding string `json:"encoding,omitempty" path:"encoding,omitempty" url:"encoding,omitempty"`
	// Additional request headers.
	Headers interface{} `json:"headers,omitempty" path:"headers,omitempty" url:"headers,omitempty"`
	// Additional body parameters.
	Body interface{} `json:"body,omitempty" path:"body,omitempty" url:"body,omitempty"`
	// raw body text.
	RawBody string `json:"raw_body,omitempty" path:"raw_body,omitempty" url:"raw_body,omitempty"`
	// Send the file data as the request body?
	FileAsBody *bool `json:"file_as_body,omitempty" path:"file_as_body,omitempty" url:"file_as_body,omitempty"`
	// Send the file data as a named parameter in the request POST body.
	FileFormField string `json:"file_form_field,omitempty" path:"file_form_field,omitempty" url:"file_form_field,omitempty"`
	// action for test body.
	Action string `json:"action,omitempty" path:"action,omitempty" url:"action,omitempty"`
	// Use dedicated IPs for sending the webhook?
	UseDedicatedIps *bool `json:"use_dedicated_ips,omitempty" path:"use_dedicated_ips,omitempty" url:"use_dedicated_ips,omitempty"`
}

// Identifier no path or id

// WebhookTestCollection is a list of WebhookTest resources.
type WebhookTestCollection []WebhookTest

// WebhookTestCreateParams contains the request parameters for POST /webhook_tests.
type WebhookTestCreateParams struct {
	// URL for testing the webhook. Required.
	Url string `url:"url" json:"url" path:"url"`
	// HTTP method(GET or POST).
	Method string `url:"method,omitempty" json:"method,omitempty" path:"method"`
	// HTTP encoding method. Can be JSON, XML, or RAW (form data).
	Encoding string `url:"encoding,omitempty" json:"encoding,omitempty" path:"encoding"`
	// Additional request headers.
	Headers interface{} `url:"headers,omitempty" json:"headers,omitempty" path:"headers"`
	// Additional body parameters.
	Body interface{} `url:"body,omitempty" json:"body,omitempty" path:"body"`
	// raw body text.
	RawBody string `url:"raw_body,omitempty" json:"raw_body,omitempty" path:"raw_body"`
	// Send the file data as the request body?
	FileAsBody *bool `url:"file_as_body,omitempty" json:"file_as_body,omitempty" path:"file_as_body"`
	// Send the file data as a named parameter in the request POST body.
	FileFormField string `url:"file_form_field,omitempty" json:"file_form_field,omitempty" path:"file_form_field"`
	// action for test body.
	Action string `url:"action,omitempty" json:"action,omitempty" path:"action"`
	// Use dedicated IPs for sending the webhook?
	UseDedicatedIps *bool `url:"use_dedicated_ips,omitempty" json:"use_dedicated_ips,omitempty" path:"use_dedicated_ips"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (w *WebhookTest) UnmarshalJSON(data []byte) error {
	type webhookTest WebhookTest
	var v webhookTest
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*w = WebhookTest(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (w *WebhookTestCollection) UnmarshalJSON(data []byte) error {
	type webhookTests WebhookTestCollection
	var v webhookTests
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*w = WebhookTestCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (w *WebhookTestCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*w))
	for i, v := range *w {
		ret[i] = v
	}

	return &ret
}
