package terraform

import (
	"testing"
	"xo/src/massdriver"
)

func TestGenerateJSONBackendS3Config(t *testing.T) {
	spec := massdriver.Specification{
		S3StateBucket:             "bucket",
		OrganizationID:            "org",
		PackageID:                 "pkg",
		S3StateRegion:             "region",
		DynamoDBStateLockTableArn: "arn:aws:dynamodb:us-west-2:111111111111:table/dynamoDbTable",
	}
	got, _ := GenerateJSONBackendS3Config(&spec, "step")
	want := doc(`
	{
		"terraform": {
			"backend": {
				"s3": {
					"bucket": "bucket",
					"key": "org/pkg/step.tfstate",
					"region": "region",
					"dynamodb_table": "dynamoDbTable"
				}
			}
		}
	}
`)

	if string(got) != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestGetS3StateKey(t *testing.T) {
	got := GetS3StateKey("org", "pkg", "step")
	want := "org/pkg/step.tfstate"

	if string(got) != want {
		t.Errorf("got %s want %s", got, want)
	}
}
