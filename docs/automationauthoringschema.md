# AutomationAuthoringSchema

`files_sdk.AutomationAuthoringSchema` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/automationauthoringschema`, named `automation_authoring_schema`.

This package has no API operations of its own.

## Model

`files_sdk.AutomationAuthoringSchema` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `DefinitionSchema` | `interface{}` | `definition_schema` | JSON Schema for active Automation v2 graph definitions. |
| `ErrorFamilies` | `[]map[string]interface{}` | `error_families` | Typed error families accepted by Automation v2 on_error rules. |
| `Nodes` | `[]map[string]interface{}` | `nodes` | Active Automation v2 node authoring metadata. |
| `SchemaUrl` | `string` | `schema_url` | Stable public URL for the Automation v2 graph definition JSON Schema. |
