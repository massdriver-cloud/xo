package tfdef

import (
	"fmt"
	"sort"
	"strings"
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
func NewTFVariable(p Property) TFVariable {
	t := convertPropertyToType(p)
	return TFVariable{Type: t}
}

func convertPropertyToType(p Property) string {
	switch p.Type {
	case "array":
		t := p.Items.Type
		if len(t) == 0 {
			t = "any"
		}
		return fmt.Sprintf("list(%s)", t)
	case "object":
		var types []string
		for name, prop := range p.Properties {
			subType := fmt.Sprintf("%s = %s", name, convertPropertyToType(prop))
			types = append(types, subType)
		}
		sort.Strings(types)
		strTypes := strings.Join(types, ", ")
		return fmt.Sprintf("object(%s)", strTypes)
	default:
		return p.Type
	}
}
