package bundles_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"xo/src/bundles"
)

type TestCase struct {
	Name     string
	Input    interface{}
	Expected interface{}
}

func TestHydrate(t *testing.T) {
	bundles.ArtifactPath = "./testdata/artifacts"
	bundles.SpecPath = "./testdata/specs"

	cases := []TestCase{
		{
			Name:  "Hydrates a shallow map with an artifact ref",
			Input: jsonDecode(`{"key": "artifact://aws-example"}`),
			Expected: map[string]map[string]string{
				"key": {
					"id": "fake-schema-id",
				},
			},
		},
		{
			Name:  "Hydrates a $ref",
			Input: jsonDecode(`{"$ref": "./testdata/artifacts/aws-example.json"}`),
			Expected: map[string]string{
				"id": "fake-schema-id",
			},
		},
		{
			Name:  "Hydrates a nested $ref",
			Input: jsonDecode(`{"key": {"$ref": "./testdata/artifacts/aws-example.json"}}`),
			Expected: map[string]map[string]string{
				"key": {
					"id": "fake-schema-id",
				},
			},
		},
		{
			Name:  "Hydrates a $ref recursively",
			Input: jsonDecode(`{"$ref": "./testdata/artifacts/ref-aws-example.json"}`),
			Expected: map[string]map[string]string{
				"properties": {
					"id": "fake-schema-id",
				},
			},
		},
		{
			Name:  "Does not hydrate HTTPS refs",
			Input: jsonDecode(`{"$ref": "https://elsewhere.com/schema.json"}`),
			Expected: map[string]string{
				"$ref": "https://elsewhere.com/schema.json",
			},
		},
		{
			Name:  "Does not hydrate fragment (#) refs",
			Input: jsonDecode(`{"$ref": "#/its-in-this-file"}`),
			Expected: map[string]string{
				"$ref": "#/its-in-this-file",
			},
		},
		{
			Name:  "Hydrates a shallow map with an spec ref",
			Input: jsonDecode(`{"key": "spec://kubernetes"}`),
			Expected: map[string]map[string]string{
				"key": {
					"version": "1.15",
				},
			},
		},
		{
			Name:  "Map with arbiratry values",
			Input: jsonDecode(`{"s": "just-a-string", "i": 3, "key": "artifact://aws-example"}`),
			Expected: map[string]interface{}{
				"s": "just-a-string",
				"i": 3,
				"key": map[string]interface{}{
					"id": "fake-schema-id",
				},
			},
		},
		{
			Name:  "Nested map",
			Input: jsonDecode(`{"parent": {"key": "artifact://aws-example"}}`),
			Expected: map[string]interface{}{
				"parent": map[string]interface{}{
					"key": map[string]interface{}{
						"id": "fake-schema-id",
					},
				},
			},
		},
		{
			Name:  "Lists",
			Input: jsonDecode(`{"list": ["string", {"key": "artifact://aws-example"}]}`),
			Expected: map[string]interface{}{
				"list": []interface{}{
					"string",
					map[string]interface{}{
						"key": map[string]interface{}{
							"id": "fake-schema-id",
						},
					},
				},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			got := bundles.Hydrate(test.Input)

			if fmt.Sprint(got) != fmt.Sprint(test.Expected) {
				t.Errorf("got %v, want %v", got, test.Expected)
			}

		})
	}
}

func jsonDecode(data string) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal([]byte(data), &result)
	return result
}
