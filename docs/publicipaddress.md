# PublicIpAddress

`files_sdk.PublicIpAddress` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/publicipaddress`, named `public_ip_address`.

This package has no API operations of its own.

## Model

`files_sdk.PublicIpAddress` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `IpAddress` | `string` | `ip_address` | The public IP address. |
| `ServerName` | `string` | `server_name` | The name of the frontend server. |
| `FtpEnabled` | `*bool` | `ftp_enabled` |  |
| `SftpEnabled` | `*bool` | `sftp_enabled` |  |
