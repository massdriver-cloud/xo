package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

type TestCase struct {
	Name     string
	Input    interface{}
	Expected interface{}
}

func TestHydrate(t *testing.T) {
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
			got := Hydrate(test.Input)

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
