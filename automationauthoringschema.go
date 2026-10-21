package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// AutomationAuthoringSchema is a Files.com API resource.
type AutomationAuthoringSchema struct {
	// JSON Schema for active Automation v2 graph definitions.
	DefinitionSchema interface{} `json:"definition_schema,omitempty" path:"definition_schema,omitempty" url:"definition_schema,omitempty"`
	// Typed error families accepted by Automation v2 on_error rules.
	ErrorFamilies []map[string]interface{} `json:"error_families,omitempty" path:"error_families,omitempty" url:"error_families,omitempty"`
	// Active Automation v2 node authoring metadata.
	Nodes []map[string]interface{} `json:"nodes,omitempty" path:"nodes,omitempty" url:"nodes,omitempty"`
	// Stable public URL for the Automation v2 graph definition JSON Schema.
	SchemaUrl string `json:"schema_url,omitempty" path:"schema_url,omitempty" url:"schema_url,omitempty"`
}

// Identifier no path or id

// AutomationAuthoringSchemaCollection is a list of AutomationAuthoringSchema resources.
type AutomationAuthoringSchemaCollection []AutomationAuthoringSchema

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (a *AutomationAuthoringSchema) UnmarshalJSON(data []byte) error {
	type automationAuthoringSchema AutomationAuthoringSchema
	var v automationAuthoringSchema
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = AutomationAuthoringSchema(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (a *AutomationAuthoringSchemaCollection) UnmarshalJSON(data []byte) error {
	type automationAuthoringSchemas AutomationAuthoringSchemaCollection
	var v automationAuthoringSchemas
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = AutomationAuthoringSchemaCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (a *AutomationAuthoringSchemaCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
