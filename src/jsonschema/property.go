package jsonschema

// PropertiesMap is a named map of Property
type PropertiesMap map[string]Property

// PropertyItemsType is the type for JSON Schema arrays
type PropertyItemsType struct {
	Type       string        `json:"type"`
	Properties PropertiesMap `json:"properties,omitempty"`
}

type GenerateAuthFile struct {
	Format   string  `json:"format"`
	Template *string `json:"template,omitempty"`
}

// Property is a single JSON Schema property field
type Property struct {
	Type             string            `json:"type"`
	Items            PropertyItemsType `json:"items,omitempty"`
	Properties       PropertiesMap     `json:"properties,omitempty"`
	GenerateAuthFile *GenerateAuthFile `json:"md.generateAuthFile,omitempty"`
}

// Schema is a flimsy representation of a JSON Schema.
// It provides just enough structure to get type information.
type Schema struct {
	Properties PropertiesMap `json:"properties"`
}
