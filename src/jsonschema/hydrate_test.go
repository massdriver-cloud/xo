package jsonschema_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"xo/src/jsonschema"
)

type TestCase struct {
	Name     string
	Input    jsonschema.OrderedJSON
	Expected jsonschema.OrderedJSON
}

func TestHydrate(t *testing.T) {
	cases := []TestCase{
		{
			Name:     "Hydrates a $ref",
			Input:    jsonDecode(`{"$ref": "./testdata/artifacts/aws-example.json"}`),
			Expected: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{{Key: "id", Value: "fake-schema-id"}}),
		},
		{
			Name:  "Hydrates a $ref alongside arbitrary values",
			Input: jsonDecode(`{"foo": true, "bar": {}, "$ref": "./testdata/artifacts/aws-example.json"}`),
			Expected: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
				{Key: "foo", Value: true},
				{Key: "bar", Value: jsonschema.OrderedJSON{}},
				{Key: "id", Value: "fake-schema-id"},
			}),
		},
		{
			Name:  "Hydrates a nested $ref",
			Input: jsonDecode(`{"key": {"$ref": "./testdata/artifacts/aws-example.json"}}`),
			Expected: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
				{Key: "key", Value: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
					{Key: "id", Value: "fake-schema-id"},
				})},
			}),
		},
		{
			Name:     "Does not hydrate HTTPS refs",
			Input:    jsonDecode(`{"$ref": "https://elsewhere.com/schema.json"}`),
			Expected: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{{Key: "$ref", Value: "https://elsewhere.com/schema.json"}}),
		},
		{
			Name:     "Does not hydrate fragment (#) refs",
			Input:    jsonDecode(`{"$ref": "#/its-in-this-file"}`),
			Expected: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{{Key: "$ref", Value: "#/its-in-this-file"}}),
		},
		{
			Name:  "Hydrates $refs in a list",
			Input: jsonDecode(`{"list": ["string", {"$ref": "./testdata/artifacts/aws-example.json"}]}`),
			Expected: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
				{Key: "list", Value: []interface{}{
					"string",
					jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{{Key: "id", Value: "fake-schema-id"}}),
				}},
			}),
		},
		{
			Name:  "Hydrates a $ref deterministically (keys outside of ref always win)",
			Input: jsonDecode(`{"conflictingKey": "not-from-ref", "$ref": "./testdata/artifacts/conflicting-keys.json"}`),
			Expected: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
				{Key: "conflictingKey", Value: "not-from-ref"},
				{Key: "nonConflictKey", Value: "from-ref"},
			}),
		},
		{
			Name:  "Hydrates a $ref recursively",
			Input: jsonDecode(`{"$ref": "./testdata/artifacts/ref-aws-example.json"}`),
			Expected: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
				{Key: "properties", Value: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
					{Key: "id", Value: "fake-schema-id"},
				})},
			}),
		},
		{
			Name:  "Hydrates a $ref recursively",
			Input: jsonDecode(`{"$ref": "./testdata/artifacts/ref-lower-dir-aws-example.json"}`),
			Expected: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
				{Key: "properties", Value: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
					{Key: "id", Value: "fake-schema-id"},
				})},
			}),
		},
		{
			Name:  `Adds "additionalProperties":false to object types`,
			Input: jsonDecode(`{"properties": {"a": "b"}, "type": "object"}`),
			Expected: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
				{Key: "properties", Value: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
					{Key: "a", Value: "b"},
				})},
				{Key: "type", Value: "object"},
				{Key: "additionalProperties", Value: false},
			}),
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			hydrated, _ := jsonschema.Hydrate(test.Input, ".")
			got := hydrated.(jsonschema.OrderedJSON)

			if fmt.Sprint(got) != fmt.Sprint(test.Expected) {
				t.Errorf("got %v, want %v", got, test.Expected)
			}
		})
	}
}

func jsonDecode(data string) jsonschema.OrderedJSON {
	var result jsonschema.OrderedJSON
	json.Unmarshal([]byte(data), &result)
	return result
}
