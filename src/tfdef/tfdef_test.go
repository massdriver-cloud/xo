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
func TestConvertRemoteRefSchemas(t *testing.T) {
	got := Convert("file://./testdata/remote-ref-schema.json")
	want := doc(`
	{
		"variable": {
			"local": {
				"type": "string"
			},
			"remote": {
				"type": "string"
			}
		}
	}
`)

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestConvertScalarTypes(t *testing.T) {
	got := Convert("file://./testdata/scalar-types-schema.json")
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
