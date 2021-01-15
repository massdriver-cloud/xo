package schemaloader

import (
	"testing"
)

type test struct {
	name  string
	input string
	want  string
}

func TestLoad(t *testing.T) {
	tests := []test{
		{
			name:  "Without prefix",
			input: "./testdata/schema.json",
			want:  "https://example.com/person.schema.json",
		},
		{
			name:  "With prefix",
			input: "file://./testdata/schema.json",
			want:  "https://example.com/person.schema.json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sl := Load(tc.input)
			schema, _ := sl.LoadJSON()
			got := schema.(map[string]interface{})["$id"]

			if got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
