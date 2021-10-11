package jsonschema_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"xo/src/jsonschema"

	"gopkg.in/yaml.v2"
)

func TestMarshalJSON(t *testing.T) {
	type test struct {
		name  string
		input jsonschema.OrderedJSON
		want  string
	}

	tests := []test{
		{
			name: "normal test",
			input: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
				{Key: "foo", Value: "bar"},
				{Key: "name", Value: "John Doe"},
				{Key: "address", Value: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
					{Key: "street", Value: "123 E 3rd St"},
					{Key: "city", Value: "Denver"},
					{Key: "test", Value: []interface{}{"a", []interface{}{1, 2}, "b", jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
						{Key: "county", Value: "Jefferson"},
						{Key: "district", Value: 20},
					})}},
					{Key: "state", Value: "CO"},
					{Key: "zip", Value: 81526},
				})},
				{Key: "nestedArrayMap", Value: []interface{}{[]interface{}{jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
					{Key: "key", Value: "value"},
				})}}},
				{Key: "nestedArrayScalar", Value: []interface{}{[]interface{}{3}}},
				{Key: "anotherTest", Value: []interface{}{1, 2, 3, 4}},
				{Key: "emptyArray", Value: []interface{}{}},
			}),
			want: `{"foo":"bar","name":"John Doe","address":{"street":"123 E 3rd St","city":"Denver","test":["a",[1,2],"b",{"county":"Jefferson","district":20}],"state":"CO","zip":81526},"nestedArrayMap":[[{"key":"value"}]],"nestedArrayScalar":[[3]],"anotherTest":[1,2,3,4],"emptyArray":[]}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := json.Marshal(tc.input)

			if fmt.Sprint(string(got)) != fmt.Sprint(tc.want) {
				t.Errorf("got %v, want %v", string(got), tc.want)
			}
		})
	}
}

func TestUnMarshalJSON(t *testing.T) {
	type test struct {
		name  string
		input string
		want  jsonschema.OrderedJSON
	}

	tests := []test{
		{
			name: "normal test",
			input: `{
    "foo": "bar",
    "name": "John Doe",
    "address": {
        "street": "123 E 3rd St",
        "city": "Denver",
        "test": ["a", [1, 2], "b", {
            "county": "Jefferson",
            "district": 20
        }],
        "state": "CO",
        "zip": 81526
	},
	"nestedArrayMap": [[{"key": "value"}]],
	"nestedArrayScalar": [[3]],
    "anotherTest": [1, 2, 3, 4],
	"emptyArray": []
}`,
			want: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
				{Key: "foo", Value: "bar"},
				{Key: "name", Value: "John Doe"},
				{Key: "address", Value: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
					{Key: "street", Value: "123 E 3rd St"},
					{Key: "city", Value: "Denver"},
					{Key: "test", Value: []interface{}{"a", []interface{}{1, 2}, "b", jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
						{Key: "county", Value: "Jefferson"},
						{Key: "district", Value: 20},
					})}},
					{Key: "state", Value: "CO"},
					{Key: "zip", Value: 81526},
				})},
				{Key: "nestedArrayMap", Value: []interface{}{[]interface{}{jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
					{Key: "key", Value: "value"},
				})}}},
				{Key: "nestedArrayScalar", Value: []interface{}{[]interface{}{3}}},
				{Key: "anotherTest", Value: []interface{}{1, 2, 3, 4}},
				{Key: "emptyArray", Value: []interface{}{}},
			}),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := jsonschema.OrderedJSON{}
			json.Unmarshal([]byte(tc.input), &got)

			if fmt.Sprint(got) != fmt.Sprint(tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestUnMarshalYAML(t *testing.T) {
	type test struct {
		name  string
		input string
		want  jsonschema.OrderedJSON
	}

	tests := []test{
		{
			name: "normal test",
			input: `foo: bar
name: John Doe
address:
    street: 123 E 3rd St
    city: Denver
    test:
      - a
      - [1, 2]
      - b
      - county: Jefferson
        district: 20
    state: CO
    zip: 81526
nestedArrayMap:
  - - key: value
nestedArrayScalar:
  - - 3
anotherTest: [1, 2, 3, 4]
emptyArray: []`,
			want: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
				{Key: "foo", Value: "bar"},
				{Key: "name", Value: "John Doe"},
				{Key: "address", Value: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
					{Key: "street", Value: "123 E 3rd St"},
					{Key: "city", Value: "Denver"},
					{Key: "test", Value: []interface{}{"a", []interface{}{1, 2}, "b", jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
						{Key: "county", Value: "Jefferson"},
						{Key: "district", Value: 20},
					})}},
					{Key: "state", Value: "CO"},
					{Key: "zip", Value: 81526},
				})},
				{Key: "nestedArrayMap", Value: []interface{}{[]interface{}{jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
					{Key: "key", Value: "value"},
				})}}},
				{Key: "nestedArrayScalar", Value: []interface{}{[]interface{}{3}}},
				{Key: "anotherTest", Value: []interface{}{1, 2, 3, 4}},
				{Key: "emptyArray", Value: []interface{}{}},
			}),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := jsonschema.OrderedJSON{}
			yaml.Unmarshal([]byte(tc.input), &got)

			if fmt.Sprint(got) != fmt.Sprint(tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
