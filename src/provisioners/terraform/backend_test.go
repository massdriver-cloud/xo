package terraform

import (
	"testing"
	"xo/src/massdriver"

	"github.com/stretchr/testify/require"
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

	require.JSONEq(t, string(got), want)
}

func TestGenerateJSONBackendHTTPConfig(t *testing.T) {
	spec := massdriver.Specification{
		DeploymentID: "depId",
		Token:        "token",
		PackageID:    "pkgId",
	}
	got, _ := GenerateJSONBackendHTTPConfig(&spec, "step")
	want := doc(`
	{
		"terraform": {
			"backend": {
				"http": {
					"username": "depId",
					"password": "token",
					"address": "https://api.massdriver.cloud/state/pkgId/step",
					"lock_address": "https://api.massdriver.cloud/state/pkgId/step",
					"unlock_address": "https://api.massdriver.cloud/state/pkgId/step"
				}
			}
		}
	}
`)

	require.JSONEq(t, string(got), want)
}

func TestGetS3StateKey(t *testing.T) {
	got := GetS3StateKey("org", "pkg", "step")
	want := "org/pkg/step.tfstate"

	if string(got) != want {
		t.Errorf("got %s want %s", got, want)
	}
}
