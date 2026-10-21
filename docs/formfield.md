# FormField

`files_sdk.FormField` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/formfield`, named `form_field`.

This package has no API operations of its own.

## Model

`files_sdk.FormField` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Form field id |
| `Label` | `string` | `label` | Label to be displayed |
| `Required` | `*bool` | `required` | Is this a required field? |
| `HelpText` | `string` | `help_text` | Help text to be displayed |
| `FieldType` | `string` | `field_type` | Type of Field |
| `OptionsForSelect` | `[]string` | `options_for_select` | Options to display for radio and dropdown |
| `DefaultOption` | `string` | `default_option` | Default option for radio and dropdown |
| `FormFieldSetId` | `int64` | `form_field_set_id` | Form field set id |
