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

func TestBuild(t *testing.T) {
	var bundle = bundles.ParseBundle("./testdata/bundle.Build/bundle.yaml")
	bundle.Build("./tmp/")

	// TODO: assert the files are there :D

	defer os.RemoveAll("./tmp")
}

func TestBuildSchema(t *testing.T) {
	var bundle = bundles.ParseBundle("./testdata/bundle.Build/bundle.yaml")
	var inputIo bytes.Buffer

	bundles.BuildSchema(bundle.Inputs, bundle.Metadata("inputs"), &inputIo)
	var gotJson = &map[string]interface{}{}
	_ = json.Unmarshal(inputIo.Bytes(), gotJson)

	wantBytes, _ := ioutil.ReadFile("./testdata/bundle.Build/schema-inputs.json")
	var wantJson = &map[string]interface{}{}
	_ = json.Unmarshal(wantBytes, wantJson)

	if !reflect.DeepEqual(gotJson, wantJson) {
		t.Errorf("got %v, want %v", gotJson, wantJson)
	}
}

func TestParseBundle(t *testing.T) {
	// TODO: add some $refs to the bundle.yaml for a better test
	var got = bundles.ParseBundle("./testdata/bundle.yaml")
	var want = bundles.Bundle{
		Uuid:        "FC2C7101-86A6-437B-B8C2-A2391FE8C847",
		Schema:      "draft-07",
		Slug:        "aws-vpc",
		Title:       "AWS VPC",
		Description: "Something",
		Artifacts:   map[string]interface{}{},
		Inputs: map[string]interface{}{
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
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
