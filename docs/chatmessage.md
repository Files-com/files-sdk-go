# ChatMessage

`files_sdk.ChatMessage` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/chatmessage`, named `chat_message`.

This package has no API operations of its own.

## Model

`files_sdk.ChatMessage` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Chat Message ID. |
| `Role` | `string` | `role` | Message role. |
| `Content` | `string` | `content` | Message content. |
| `CreatedAt` | `*time.Time` | `created_at` | Message creation date/time. |
