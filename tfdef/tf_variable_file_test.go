package tfdef

import (
	"encoding/json"
	"reflect"
	"testing"
)

type test struct {
	name  string
	input Property
	want  TFVariable
}

func TestNewTFVariable(t *testing.T) {
	tests := []test{
		{
			name:  "scalars",
			input: Property{Type: "number"},
			want:  TFVariable{Type: "number"},
		},
		{
			name:  "list of scalars",
			input: Property{Type: "array", Items: PropertyItemsType{Type: "string"}},
			want:  TFVariable{Type: "list(string)"},
		},
		{
			name:  "list of any",
			input: Property{Type: "array"},
			want:  TFVariable{Type: "list(any)"},
		},
		{
			name:  "object w/ scalars",
			input: Property{Type: "object", Properties: PropertiesMap{"street_number": Property{Type: "number"}, "street_name": Property{Type: "string"}}},
			want:  TFVariable{Type: "object({street_name = string, street_number = number})"},
		},
		{
			name: "complex objects",
			input: Property{
				Type: "object",
				Properties: PropertiesMap{
					"name": Property{Type: "string"},
					"children": Property{
						Type: "array",
						Items: PropertyItemsType{
							Type: "object",
							Properties: PropertiesMap{
								"name": Property{Type: "string"},
							},
						},
					},
				},
			},
			want: TFVariable{Type: "object({children = list(object({name = string})), name = string})"},
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
			input: TFVariableFile{Variable: map[string]TFVariable{"name": TFVariable{Type: "string"}}},
			want:  `{"variable":{"name":{"type":"string"}}}`,
		},
		{
			name:  "Multiple variables",
			input: TFVariableFile{Variable: map[string]TFVariable{"name": TFVariable{Type: "string"}, "age": TFVariable{Type: "number"}}},
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
