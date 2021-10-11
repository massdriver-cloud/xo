package terraform

import (
	"encoding/json"
	"reflect"
	"testing"
	"xo/src/jsonschema"
)

type test struct {
	name  string
	input jsonschema.Property
	want  TFVariable
}

func TestNewTFVariable(t *testing.T) {
	tests := []test{
		{
			name:  "scalars",
			input: jsonschema.Property{Type: "number"},
			want:  TFVariable{Type: "number"},
		},
		{
			name:  "list of scalars",
			input: jsonschema.Property{Type: "array", Items: &jsonschema.Property{Type: "string"}},
			want:  TFVariable{Type: "any"},
		},
		{
			name:  "list of any",
			input: jsonschema.Property{Type: "array", Items: &jsonschema.Property{}},
			want:  TFVariable{Type: "any"},
		},
		{
			name:  "maps",
			input: jsonschema.Property{Type: "object", AdditionalProperties: true},
			want:  TFVariable{Type: "any"},
		},
		{
			name:  "object w/ scalars",
			input: jsonschema.Property{Type: "object", Properties: jsonschema.PropertiesMap{"street_number": jsonschema.Property{Type: "number"}, "street_name": jsonschema.Property{Type: "string"}}},
			want:  TFVariable{Type: "any"},
		},
		{
			name: "complex objects",
			input: jsonschema.Property{
				Type: "object",
				Properties: jsonschema.PropertiesMap{
					"name": jsonschema.Property{Type: "string"},
					"children": jsonschema.Property{
						Type: "array",
						Items: &jsonschema.Property{
							Type: "object",
							Properties: jsonschema.PropertiesMap{
								"name": jsonschema.Property{Type: "string"},
							},
						},
					},
				},
			},
			want: TFVariable{Type: "any"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := NewTFVariable(tc.input)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestTFVariableFileJSONEncoding(t *testing.T) {
	type test struct {
		name  string
		input TFVariableFile
		want  string
	}

	tests := []test{
		{
			name:  "A single variable",
			input: TFVariableFile{Variable: map[string]TFVariable{"name": {Type: "string"}}},
			want:  `{"variable":{"name":{"type":"string"}}}`,
		},
		{
			name:  "Multiple variables",
			input: TFVariableFile{Variable: map[string]TFVariable{"name": {Type: "string"}, "age": {Type: "number"}}},
			want:  `{"variable":{"age":{"type":"number"},"name":{"type":"string"}}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bytes, _ := json.Marshal(tc.input)
			got := string(bytes)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
