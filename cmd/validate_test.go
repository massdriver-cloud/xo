package cmd

import (
	"os"
	"path"
	"testing"
)

func TestValidateJSONDocument(t *testing.T) {
	cwd, _ := os.Getwd()
	schema := path.Join(cwd, "testdata/valid-schema.json")
	got, _ := Validate(schema, "testdata/valid-document.json")
	want := true

	if got != want {
		t.Errorf("got %t want %t", got, want)
	}
}

func TestInvalidateJSONDocument(t *testing.T) {
	got, _ := Validate("testdata/valid-schema.json", "testdata/invalid-document.json")
	want := false

	if got != want {
		t.Errorf("got %t want %t", got, want)
	}
}
