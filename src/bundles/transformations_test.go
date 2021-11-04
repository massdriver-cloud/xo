package bundles_test

import (
	"reflect"
	"testing"
	"xo/src/bundles"
	"xo/src/jsonschema"
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
		Artifacts:   jsonschema.OrderedJSON{},
		Connections: jsonschema.OrderedJSON{},
		Params: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
			{Key: "properties", Value: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
				{Key: "set_id_test", Value: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
					{Key: "type", Value: "array"},
					{Key: "items", Value: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
						{Key: "type", Value: "object"},
						{Key: "properties", Value: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
							{Key: "foo", Value: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
								{Key: "type", Value: "string"},
							})},
							{Key: "md_set_id", Value: jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
								{Key: "type", Value: "string"},
							})},
						})},
						{Key: "additionalProperties", Value: false},
					})},
				})},
			})},
		}),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
