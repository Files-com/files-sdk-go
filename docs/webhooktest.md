# WebhookTest

`files_sdk.WebhookTest` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/webhooktest`, named `webhooktest`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/webhooktest"
)
```

## Operations

### Create

`POST /webhook_tests`

Create Webhook Test

```go
func (c *Client) Create(params files_sdk.WebhookTestCreateParams, opts ...files_sdk.RequestResponseOption) (webhookTest files_sdk.WebhookTest, err error)
func Create(params files_sdk.WebhookTestCreateParams, opts ...files_sdk.RequestResponseOption) (webhookTest files_sdk.WebhookTest, err error)
```

Parameters, `files_sdk.WebhookTestCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Url` | `string` | `url` | yes | URL for testing the webhook. |
| `Method` | `string` | `method` | no | HTTP method(GET or POST). |
| `Encoding` | `string` | `encoding` | no | HTTP encoding method. Can be JSON, XML, or RAW (form data). |
| `Headers` | `interface{}` | `headers` | no | Additional request headers. |
| `Body` | `interface{}` | `body` | no | Additional body parameters. |
| `RawBody` | `string` | `raw_body` | no | raw body text |
| `FileAsBody` | `*bool` | `file_as_body` | no | Send the file data as the request body? |
| `FileFormField` | `string` | `file_form_field` | no | Send the file data as a named parameter in the request POST body |
| `Action` | `string` | `action` | no | action for test body |
| `UseDedicatedIps` | `*bool` | `use_dedicated_ips` | no | Use dedicated IPs for sending the webhook? |

Example:

```go
result, err := webhooktest.Create(files_sdk.WebhookTestCreateParams{
	Url: "https://www.site.com/...",
})
if err != nil {
	return err
}
fmt.Println(result)
```

## Model

`files_sdk.WebhookTest` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Code` | `int64` | `code` | Status HTTP code |
| `Message` | `string` | `message` | Error message |
| `Status` | `string` | `status` | Status message |
| `Data` | `files_sdk.Auto` | `data` | Additional data |
| `Success` | `*bool` | `success` | The success status of the webhook test |
| `Url` | `string` | `url` | URL for testing the webhook. |
| `Method` | `string` | `method` | HTTP method(GET or POST). |
| `Encoding` | `string` | `encoding` | HTTP encoding method. Can be JSON, XML, or RAW (form data). |
| `Headers` | `interface{}` | `headers` | Additional request headers. |
| `Body` | `interface{}` | `body` | Additional body parameters. |
| `RawBody` | `string` | `raw_body` | raw body text |
| `FileAsBody` | `*bool` | `file_as_body` | Send the file data as the request body? |
| `FileFormField` | `string` | `file_form_field` | Send the file data as a named parameter in the request POST body |
| `Action` | `string` | `action` | action for test body |
| `UseDedicatedIps` | `*bool` | `use_dedicated_ips` | Use dedicated IPs for sending the webhook? |
