# FileAction

`files_sdk.FileAction` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/fileaction`, named `file_action`.

This package has no API operations of its own.

## Model

`files_sdk.FileAction` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Status` | `string` | `status` | Status of file operation. |
| `FileMigrationId` | `int64` | `file_migration_id` | If status is pending, this is the id of the File Migration to check for status updates. Delete actions include it whenever the deletion runs as a File Migration, including while it is processing. |
