package terraform

import (
	"fmt"
	"sort"
	"strings"

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
	t := convertPropertyToType(p.Type, p.Properties, p.Items)
	return TFVariable{Type: t}
}

func convertPropertyToType(pType string, pProperties jsonschema.PropertiesMap, pItems jsonschema.PropertyItemsType) string {
	switch pType {
	case "array":
		var t string
		switch pItems.Type {
		case "":
			t = "any"
		case "array":
			t = convertArray()
		case "object":
			t = convertObject(pItems.Properties)
		default:
			t = pItems.Type
		}

		return fmt.Sprintf("list(%s)", t)
	case "object":
		return convertObject(pProperties)
	default:
		return convertScalar(pType)
	}
}

func convertArray() string {
	// t = convertPropertyToType(pItems.Type, pItems.Properties, PropertyItemsType{})
	//TODO
	return "convertArray - not implemented"
}

func convertObject(pProperties jsonschema.PropertiesMap) string {
	var types []string
	for name, prop := range pProperties {
		subType := convertPropertyToType(prop.Type, prop.Properties, prop.Items)
		subTypeDeclaration := fmt.Sprintf("%s = %s", name, subType)
		types = append(types, subTypeDeclaration)
	}
	sort.Strings(types)
	strTypes := strings.Join(types, ", ")
	return fmt.Sprintf("object({%s})", strTypes)
}

func convertScalar(pType string) string {
	switch pType {
	// json-schema calls it boolean, terraform calls it bool
	case "boolean":
		return "bool"
	default:
		return pType
	}
}
