package bundles_test

import (
	"reflect"
	"testing"
	"xo/src/bundles"
)

// TODO func TestBuildBundle(bundle,)
// 	// var inputsIo bytes.Buffer
// var connectionsIo bytes.Buffer
// var artifactsIo bytes.Buffer

func TestParseBundle(t *testing.T) {
	bundles.ArtifactPath = "./testdata/artifacts"
	bundles.SpecPath = "./testdata/specs"

	var got = bundles.ParseBundle("./testdata/bundle.yaml")
	var want = bundles.Bundle{
		Schema:      "draft-07",
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
