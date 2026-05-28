package resource_test

import (
	"bytes"
	"os"
	"testing"
	"xo/src/resource"
)

func TestValidate(t *testing.T) {
	type testData struct {
		name       string
		field      string
		schemaPath string
		resource   []byte
		want       bool
	}
	tests := []testData{
		{
			name:       "pass",
			field:      "one",
			schemaPath: "testdata/schema-artifacts.json",
			resource:   []byte(`{"data":{"foo":{"bar":"baz"}},"specs":{"hello":"world"}}`),
			want:       true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			schemasBytes, err := os.ReadFile(tc.schemaPath)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}
			schemasBuffer := bytes.NewBuffer(schemasBytes)

			got, err := resource.Validate(tc.field, tc.resource, schemasBuffer)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}
			if got != tc.want {
				t.Fatalf("want: %v, got: %v", tc.want, got)
			}
		})
	}
}
