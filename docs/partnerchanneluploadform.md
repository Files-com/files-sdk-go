# PartnerChannelUploadForm

`files_sdk.PartnerChannelUploadForm` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/partnerchanneluploadform`, named `partner_channel_upload_form`.

This package has no API operations of its own.

## Model

`files_sdk.PartnerChannelUploadForm` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Path` | `string` | `path` | Upload destination path relative to the From Partner folder. Must be a configured leaf managed folder path, or an empty string for the From Partner folder when no managed folders are configured. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `FormFieldSetId` | `int64` | `form_field_set_id` | ID of the Form Field Set answered for uploads to this destination. Must be in the same Workspace and contain at least one custom field. |
