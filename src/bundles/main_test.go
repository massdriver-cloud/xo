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

	want := []string{"schema-artifacts.json", "schema-connections.json", "schema-params.json"}

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
		Uuid:        "FC2C7101-86A6-437B-B8C2-A2391FE8C847",
		Schema:      "draft-07",
		Type:        "aws-vpc",
		Title:       "AWS VPC",
		Description: "Something",
		Provisioner: "terraform",
		Artifacts:   bundles.OrderedJSON{},
		Params: bundles.OrderedJSON([]bundles.OrderedJSONElement{
			{Key: "required", Value: []interface{}{"name"}},
			{Key: "properties", Value: bundles.OrderedJSON([]bundles.OrderedJSONElement{
				{Key: "name", Value: bundles.OrderedJSON([]bundles.OrderedJSONElement{
					{Key: "type", Value: "string"},
					{Key: "title", Value: "Name"},
				})},
				{Key: "age", Value: bundles.OrderedJSON([]bundles.OrderedJSONElement{
					{Key: "type", Value: "integer"},
					{Key: "title", Value: "Age"},
				})},
			})},
		}),
		Connections: bundles.OrderedJSON([]bundles.OrderedJSONElement{
			{Key: "required", Value: []interface{}{"default"}},
			{Key: "properties", Value: bundles.OrderedJSON([]bundles.OrderedJSONElement{
				{Key: "default", Value: bundles.OrderedJSON([]bundles.OrderedJSONElement{
					{Key: "type", Value: "string"},
					{Key: "title", Value: "Default credential"},
				})},
			})},
		}),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
