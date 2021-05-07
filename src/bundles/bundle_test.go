package bundles_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"
	"xo/src/bundles"
)

func init() {
	bundles.ArtifactPath = "./testdata/artifacts"
	bundles.SpecPath = "./testdata/specs"
}

func TestBuild(t *testing.T) {
	var bundle = bundles.ParseBundle("./testdata/bundle.Build/bundle.yaml")
	bundle.Build("/tmp")
}

func TestBuildSchema(t *testing.T) {
	var bundle = bundles.ParseBundle("./testdata/bundle.Build/bundle.yaml")
	var inputIo bytes.Buffer

	bundles.BuildSchema(bundle.Inputs, bundle.Metadata(), &inputIo)
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
	var got = bundles.ParseBundle("./testdata/bundle.yaml")
	var want = bundles.Bundle{
		Uuid:        "FC2C7101-86A6-437B-B8C2-A2391FE8C847",
		Schema:      "draft-07",
		Slug:        "aws-vpc",
		Title:       "AWS VPC",
		Description: "Something",
		Artifacts: map[string]interface{}{
			"items": map[string]interface{}{
				"anyOf": []interface{}{
					map[string]interface{}{
						"id": "fake-schema-id",
					},
					map[string]interface{}{
						"id": "fake-schema-id",
					},
				},
			},
		},
		Inputs: map[string]interface{}{
			"allOf": []interface{}{
				map[string]interface{}{
					"id": "fake-schema-id",
				},
				map[string]interface{}{
					"properties": map[string]interface{}{
						"specs/kubernetes": map[string]interface{}{
							"version": "1.15",
						},
					},
				},
			},
			"properties": map[string]interface{}{
				"specs": map[string]interface{}{
					"properties": map[string]interface{}{
						"platform_version": map[string]interface{}{
							"enum": []interface{}{
								"eks1",
								"eks2",
							},
							"type": "string",
						},
					},
					"type": "object",
				},
			},
			"required": []interface{}{
				"specs",
			},
		},
		Connections: map[string]interface{}{
			"properties": map[string]interface{}{
				"default": map[string]interface{}{
					"id": "fake-schema-id",
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
