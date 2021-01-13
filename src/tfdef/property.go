package tfdef

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Property is a single JSON Schema property field
type Property struct {
	Type  string `json:"type"`
	Items struct {
		Type string `json:"type"`
	} `json:"items,omitempty"`
	Properties map[string]Property `json:"properties,omitempty"`
}

func (p Property) MarshalJSON() ([]byte, error) {
	return json.Marshal(TFVariable{
		Type: p.CastType(),
	})
}

func (p Property) CastType() string {
	switch p.Type {
	case "array":
		return fmt.Sprintf("list(%s)", p.listType())
	case "object":
		var types []string
		for name, prop := range p.Properties {
			subType := fmt.Sprintf("%s = %s", name, prop.CastType())
			types = append(types, subType)
		}
		sort.Strings(types)
		strTypes := strings.Join(types, ", ")
		return fmt.Sprintf("object(%s)", strTypes)
	default:
		return p.Type
	}
}

func (p Property) listType() string {
	t := p.Items.Type
	if len(t) == 0 {
		t = "any"
	}

	return t
}
