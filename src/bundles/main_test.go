package bundles_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"xo/src/bundles"
)

func TestGenerateSchemas(t *testing.T) {
	var bundle, _ = bundles.ParseBundle("./testdata/bundle.Build/bundle.yaml")
	_ = bundle.GenerateSchemas("./tmp/")

	gotDir, _ := os.ReadDir("./tmp")
	got := []string{}

	for _, dirEntry := range gotDir {
		got = append(got, dirEntry.Name())
	}

	want := []string{"schema-artifacts.json", "schema-connections.json", "schema-params.json", "schema-ui.json"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	defer os.RemoveAll("./tmp")
}

func TestGenerateSchema(t *testing.T) {
	var bundle, _ = bundles.ParseBundle("./testdata/bundle.Build/bundle.yaml")
	var inputIo bytes.Buffer

	bundles.GenerateSchema(bundle.Params, bundle.Metadata("params"), &inputIo)
	var gotJson = &map[string]interface{}{}
	_ = json.Unmarshal(inputIo.Bytes(), gotJson)

	wantBytes, _ := ioutil.ReadFile("./testdata/bundle.Build/schema-params.json")
	var wantJson = &map[string]interface{}{}
	_ = json.Unmarshal(wantBytes, wantJson)

	if !reflect.DeepEqual(gotJson, wantJson) {
		t.Errorf("got %v, want %v", gotJson, wantJson)
	}
}

func TestParseBundle(t *testing.T) {
	var got, _ = bundles.ParseBundle("./testdata/bundle.yaml")
	var want = bundles.Bundle{
		Schema:      "draft-07",
		Name:        "aws-vpc",
		Description: "Something",
		Type:        "bundle",
		Artifacts:   map[string]interface{}{},
		Params: map[string]interface{}{
			"properties": map[string]interface{}{
				"name": map[string]interface{}{
					"type":  "string",
					"title": "Name",
				},
				"age": map[string]interface{}{
					"type":  "integer",
					"title": "Age",
				},
			},
			"required": []interface{}{
				"name",
			},
		},
		Connections: map[string]interface{}{
			"required": []interface{}{
				"default",
			},
			"properties": map[string]interface{}{
				"default": map[string]interface{}{
					"type":  "string",
					"title": "Default credential",
				},
			},
		},
		Ui: map[string]interface{}{},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
