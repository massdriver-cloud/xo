package bundles_test

import (
	"fmt"
	"io/ioutil"
	"testing"
	"xo/src/bundles"

	"gopkg.in/yaml.v3"
)

func TestTransformations(t *testing.T) {
	type testData struct {
		name           string
		schemaPath     string
		transformation func(map[string]interface{}) error
		expected       map[string]interface{}
	}
	tests := []testData{
		{
			name:           "md_set_id",
			schemaPath:     "./testdata/transformation-md_set_id.yaml",
			transformation: bundles.AddSetIdToObjectArrays,
			expected: map[string]interface{}{
				"params": map[string]interface{}{
					"properties": map[string]interface{}{
						"set_id": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type":     "object",
								"required": []string{"md_set_id"},
								"properties": map[string]interface{}{
									"foo": map[string]interface{}{
										"type": "string",
									},
									"md_set_id": map[string]interface{}{
										"type": "string",
									},
								},
							},
						},
						"no_set_id": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"bar": map[string]interface{}{
									"type": "string",
								},
							},
						},
					},
				},
			},
		},
		{
			name:           "additionalProperties",
			schemaPath:     "./testdata/transformation-additional_properties.yaml",
			transformation: bundles.DisableAdditionalPropertiesInObjects,
			expected: map[string]interface{}{
				"params": map[string]interface{}{
					"properties": map[string]interface{}{
						"addPropFalse": map[string]interface{}{
							"type":                 "object",
							"additionalProperties": false,
						},
						"addPropTrue": map[string]interface{}{
							"type": "object",
							"anyOf": []interface{}{
								"lol",
								"rofl",
							},
							"additionalProperties": true,
						},
						"addPropExists": map[string]interface{}{
							"type":                 "object",
							"additionalProperties": true,
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := ioutil.ReadFile(tc.schemaPath)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			got := map[string]interface{}{}

			err = yaml.Unmarshal([]byte(data), &got)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			err = bundles.ApplyTransformations(got, []func(map[string]interface{}) error{tc.transformation})
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if fmt.Sprint(got) != fmt.Sprint(tc.expected) {
				t.Errorf("got %v, want %v", got, tc.expected)
			}
		})
	}
}
