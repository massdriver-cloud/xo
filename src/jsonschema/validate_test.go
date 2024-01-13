package jsonschema_test

import (
	"testing"
	"xo/src/jsonschema"
)

func TestValidateJSONDocument(t *testing.T) {
	schema := jsonschema.Load("testdata/schema.json")
	document := jsonschema.Load("testdata/valid-document.json")

	got, _ := jsonschema.Validate(schema, document)
	want := true

	if got != want {
		t.Errorf("got %t want %t", got, want)
	}
}

func TestInvalidateJSONDocument(t *testing.T) {
	schema := jsonschema.Load("testdata/valid-schema.json")
	document := jsonschema.Load("testdata/invalid-document.json")

	got, _ := jsonschema.Validate(schema, document)
	want := false

	if got != want {
		t.Errorf("got %t want %t", got, want)
	}
}
