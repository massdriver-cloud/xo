package cmd

import (
	"testing"
)

func TestValidateJSONDocument(t *testing.T) {
	got := Validate("testdata/valid-schema.json", "testdata/valid-document.json")
	want := true

	if got != want {
		t.Errorf("got %t want %t", got, want)
	}
}

func TestInvalidateJSONDocument(t *testing.T) {
	got := Validate("testdata/valid-schema.json", "testdata/invalid-document.json")
	want := true

	if got != want {
		t.Errorf("got %t want %t", got, want)
	}
}

// func TestValidateYAMLDocument() {}
