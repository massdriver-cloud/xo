package bundles_test

import (
	"fmt"
	"testing"
	"xo/src/bundles"
)

func TestTransformations(t *testing.T) {
	var got, _ = bundles.ParseBundle("./testdata/transformations.yaml")
	var want = bundles.Bundle{
		Uuid:        "FC2C7101-86A6-437B-B8C2-A2391FE8C847",
		Schema:      "draft-07",
		Type:        "aws-vpc",
		Title:       "AWS VPC",
		Description: "Something",
		Provisioner: "terraform",
		Artifacts:   map[string]interface{}{},
		Connections: map[string]interface{}{},
		Params: map[string]interface{}{
			"properties": map[string]interface{}{
				"set_id_test": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type":                 "object",
						"additionalProperties": false,
						"required":             []string{"md_set_id"},
						"properties": map[string]interface{}{
							"foo": map[string]interface{}{
								"type": "string",
							},
							"md_set_id": map[string]interface{}{
								"type": "string",
							},
						},
					},
				},
			},
		},
	}

	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
