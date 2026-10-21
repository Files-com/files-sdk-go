# DirectConnectionInfo

`files_sdk.DirectConnectionInfo` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/directconnectioninfo`, named `direct_connection_info`.

This package has no API operations of its own.

## Model

`files_sdk.DirectConnectionInfo` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Version` | `int64` | `version` | Direct connection information schema version. |
| `ServerName` | `string` | `server_name` | TLS server name (SNI and Host header) for the Agent's direct transfer listener. |
| `Addresses` | `[]string` | `addresses` | Validated ip:port candidates that may be dialed over TCP+TLS for this transfer. |
| `DirectUri` | `string` | `direct_uri` | Signed HTTPS URI for direct Agent transfer traffic. |
| `CaPem` | `string` | `ca_pem` | PEM-encoded CA certificate used to verify the Agent's direct transfer TLS certificate. |
