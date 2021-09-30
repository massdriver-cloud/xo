package terraform

import (
	"xo/src/jsonschema"
)

// TFVariableFile is a representation of a variables.tf file in JSON format
type TFVariableFile struct {
	Variable map[string]TFVariable `json:"variable"`
}

// TFVariable is a representation of a terraform HCL "variable"
type TFVariable struct {
	Type string `json:"type"`
}

// NewTFVariable creates a new TFVariable from a JSON Schema Property
func NewTFVariable(p jsonschema.Property) TFVariable {
	t := convertPropertyToType(p)
	return TFVariable{Type: t}
}

func convertPropertyToType(p jsonschema.Property) string {
	switch p.Type {
	case "array":
		return convertArray(p)
	case "object":
		return convertObject(p)
	default:
		return convertScalar(p.Type)
	}
}

func convertObject(prop jsonschema.Property) string {
	_ = prop
	// See: https://github.com/massdriver-cloud/xo/issues/44
	return "any"
}

func convertArray(prop jsonschema.Property) string {
	_ = prop
	// See: https://github.com/massdriver-cloud/xo/issues/44
	return "any"
}

func convertScalar(pType string) string {
	switch pType {
	case "boolean":
		return "bool"
	case "integer":
		return "number"
	default:
		return pType
	}
}
