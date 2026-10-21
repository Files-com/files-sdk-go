# RemoteServerConfigurationFile

`files_sdk.RemoteServerConfigurationFile` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/remoteserverconfigurationfile`, named `remote_server_configuration_file`.

This package has no API operations of its own.

## Model

`files_sdk.RemoteServerConfigurationFile` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | The remote server ID of the agent |
| `PermissionSet` | `string` | `permission_set` | The permission set for the agent ['read_write', 'read_only', 'write_only'] |
| `PrivateKey` | `string` | `private_key` | The private key for the agent |
| `Subdomain` | `string` | `subdomain` | Files.com subdomain site name |
| `Root` | `string` | `root` | The root directory for the agent |
| `NodeName` | `string` | `node_name` | Optional customer-facing name for this Agent installation |
| `FollowLinks` | `*bool` | `follow_links` | Follow symlinks when traversing directories |
| `PreferProtocol` | `string` | `prefer_protocol` | Preferred network protocol ['udp', 'tcp'] (default udp) |
| `Dns` | `string` | `dns` | DNS lookup method ['auto','doh','system'] (default auto) |
| `ProxyAllOutbound` | `*bool` | `proxy_all_outbound` | Proxy all outbound traffic through files.com proxy server |
| `EndpointOverride` | `string` | `endpoint_override` | Custom site endpoint URL |
| `LogFile` | `string` | `log_file` | Log file name and location |
| `LogLevel` | `string` | `log_level` | Log level for the agent logs ['debug', 'info', 'warn', 'error', 'fatal'] (default info) |
| `LogRotateNum` | `int64` | `log_rotate_num` | Log route for agent logs. (default 5) |
| `LogRotateSize` | `int64` | `log_rotate_size` | Log route size in MB for agent logs. (default 20) |
| `OverrideMaxConcurrentJobs` | `int64` | `override_max_concurrent_jobs` | Maximum number of concurrent jobs (default 500) |
| `GracefulShutdownTimeout` | `int64` | `graceful_shutdown_timeout` | Graceful shutdown timeout in seconds (default 15) |
| `TransferRateLimit` | `string` | `transfer_rate_limit` | File transfer (upload/download) rate limit `<limit>-<period>`, with the given periods: * 'S': second * 'M': minute * 'H': hour * 'D': day Examples: * 5 requests/second: '5-S' * 10 requests/minute: '10-M' * 1000 requests/hour: '1000-H' * 2000 requests/day: '2000-D' |
| `AutoUpdatePolicy` | `string` | `auto_update_policy` | Auto update policy ['manual_trigger', 'critical_only', 'always', 'never'] (default always) |
