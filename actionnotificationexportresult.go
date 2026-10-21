package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// ActionNotificationExportResult is a Files.com API resource.
type ActionNotificationExportResult struct {
	// Notification ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// When the notification was sent.
	CreatedAt int64 `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// HTTP status code returned in the webhook response.
	Status int64 `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	// A message indicating the overall status of the webhook notification.
	Message string `json:"message,omitempty" path:"message,omitempty" url:"message,omitempty"`
	// `true` if the webhook succeeded by receiving a 200 or 204 response.
	Success *bool `json:"success,omitempty" path:"success,omitempty" url:"success,omitempty"`
	// A JSON-encoded string with headers that were sent with the webhook.
	RequestHeaders string `json:"request_headers,omitempty" path:"request_headers,omitempty" url:"request_headers,omitempty"`
	// The HTTP verb used to perform the webhook.
	RequestMethod string `json:"request_method,omitempty" path:"request_method,omitempty" url:"request_method,omitempty"`
	// The webhook request URL.
	RequestUrl string `json:"request_url,omitempty" path:"request_url,omitempty" url:"request_url,omitempty"`
	// The path to the actual file that triggered this notification. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.
	Path string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	// The folder associated with the triggering action for this notification.
	Folder string `json:"folder,omitempty" path:"folder,omitempty" url:"folder,omitempty"`
}

// Identifier returns the resource ID.
func (a ActionNotificationExportResult) Identifier() interface{} {
	return a.Id
}

// ActionNotificationExportResultCollection is a list of ActionNotificationExportResult resources.
type ActionNotificationExportResultCollection []ActionNotificationExportResult

// ActionNotificationExportResultListParams contains the request parameters for GET /action_notification_export_results.
type ActionNotificationExportResultListParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// ID of the associated action notification export. Required.
	ActionNotificationExportId int64 `url:"action_notification_export_id" json:"action_notification_export_id" path:"action_notification_export_id"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (a *ActionNotificationExportResult) UnmarshalJSON(data []byte) error {
	type actionNotificationExportResult ActionNotificationExportResult
	var v actionNotificationExportResult
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = ActionNotificationExportResult(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (a *ActionNotificationExportResultCollection) UnmarshalJSON(data []byte) error {
	type actionNotificationExportResults ActionNotificationExportResultCollection
	var v actionNotificationExportResults
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = ActionNotificationExportResultCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (a *ActionNotificationExportResultCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
