package tfdef

import (
	"encoding/json"
	"testing"
)

// Helper function for asserting json serde matches
func doc(str string) string {
	b := []byte(str)

	jsonMap := make(map[string](interface{}))
	json.Unmarshal([]byte(b), &jsonMap)

	outBytes, _ := json.Marshal(jsonMap)
	return string(outBytes)
}

// https://github.com/xeipuuv/gojsonschema#loading-local-schemas
// This test is failing because the library doesnt automatically
// resolve $refs until a document is validated. You can trick it into
// doing it w/ the last example mentioned in the above link, but
// we will need to have an idea of how we are doing that in massdriver-bundles
// first. I assume we'll end up treating the bundle's JSON Schema as the main
// and ref loading a single 'definitions' JSON Schema that has all of our
// secrets and connections
// func TestCompileRemoteRefSchemas(t *testing.T) {
// 	got := Compile("file://./testdata/remote-ref-schema.json")
// 	want := doc(`
// 	{
// 		"variable": {
// 			"local": {
// 				"type": "string"
// 			},
// 			"remote": {
// 				"type": "string"
// 			}
// 		}
// 	}
// `)

// 	if got != want {
// 		t.Errorf("got %s want %s", got, want)
// 	}
// }

func TestCompileArrayTypes(t *testing.T) {
	got := Compile("file://./testdata/array-types-schema.json")
	want := doc(`
	{
		"variable": {
			"favNumbers": {
				"type": "list(number)"
			},
			"favStrings": {
				"type": "list(string)"
			},
			"favThings": {
				"type": "list(any)"
			}
		}
	}	
`)

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestCompileObjectTypes(t *testing.T) {
	got := Compile("file://./testdata/object-types-schema.json")
	want := doc(`
	{
		"variable": {
			"address": {
				"type": "object(city = string, state = string, street_address = string)"
			}
		}
	}	
`)

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestNestedObjectTypes(t *testing.T) {
	got := Compile("file://./testdata/nested-object-schema.json")
	want := doc(`
	{
		"variable": {
			"person": {
				"type": "object(children = list(object(name = string), name = string)"
			}
		}
	}
`)

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestCompileScalarTypes(t *testing.T) {
	got := Compile("file://./testdata/scalar-types-schema.json")
	want := doc(`
	{
		"variable": {
			"name": {
				"type": "string"
			},
			"age": {
				"type": "integer"
			},
			"active": {
				"type": "boolean"
			},
			"height": {
				"type": "number"
			}
		}
	}	
`)

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
