package tfdef

import (
	"encoding/json"
	"reflect"
	"testing"
)

type test struct {
	input Property
	want  TFVariable
}

func TestNewTFVariable(t *testing.T) {
	tests := []test{
		// scalar
		{
			input: Property{Type: "number"},
			want:  TFVariable{Type: "number"},
		},
		// list of scalars
		{
			input: Property{Type: "array", Items: PropertyItemType{Type: "string"}},
			want:  TFVariable{Type: "list(string)"},
		},
		// list of any
		{
			input: Property{Type: "array"},
			want:  TFVariable{Type: "list(any)"},
		},
		// scalar objects
		{
			input: Property{Type: "object", Properties: map[string]Property{"street_number": Property{Type: "number"}, "street_name": Property{Type: "string"}}},
			want:  TFVariable{Type: "object(street_name = string, street_number = number)"},
		},
		// complex objects
		{
			input: Property{
				Type: "object",
				Properties: map[string]Property{
					"name": Property{Type: "string"},
					"children": Property{
						Type: "array",
						Items: PropertyItemType{
							Type: "object",
							Properties: map[string]Property{
								"name": Property{Type: "string"},
							},
						},
					},
				},
			},
			want: TFVariable{Type: "object(children = list(object(name = string), name = string)"},
		},
	}

	for _, tc := range tests {
		got := NewTFVariable(tc.input)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}

func TestTFVariableFileJSONEncoding(t *testing.T) {
	type test struct {
		input TFVariableFile
		want  string
	}

	tests := []test{
		{
			// A single variable
			input: TFVariableFile{Variable: map[string]TFVariable{"name": TFVariable{Type: "string"}}},
			want:  `{"variable":{"name":{"type":"string"}}}`,
		},
		{
			// Multiple variables
			input: TFVariableFile{Variable: map[string]TFVariable{"name": TFVariable{Type: "string"}, "age": TFVariable{Type: "number"}}},
			want:  `{"variable":{"age":{"type":"number"},"name":{"type":"string"}}}`,
		},
	}

	for _, tc := range tests {
		bytes, _ := json.Marshal(tc.input)
		got := string(bytes)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}
