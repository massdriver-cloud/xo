package tfdef

// PropertyItemType is the type for JSON Schema arrays
type PropertyItemType struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties,omitempty"`
}

// Property is a single JSON Schema property field
type Property struct {
	Type       string              `json:"type"`
	Items      PropertyItemType    `json:"items,omitempty"`
	Properties map[string]Property `json:"properties,omitempty"`
}

// Schema is a flimsy representation of a JSON Schema.
// It provides just enough structure to get type information.
type Schema struct {
	Properties map[string]Property `json:"properties"`
}
