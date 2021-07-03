package terraform

import (
	"testing"
)

func TestGenerateJSONBackendS3Config(t *testing.T) {
	got, _ := GenerateJSONBackendS3Config("bucket", "mrn", "region", "dynamoDbTable", "sharedCredFile", "profile")
	want := doc(`
	{
		"terraform": {
			"backend": {
				"s3": {
					"bucket": "bucket",
					"key": "mrn",
					"region": "region",
					"dynamodb_table": "dynamoDbTable",
					"shared_credentials_file": "sharedCredFile",
					"profile": "profile"
				}
			}
		}
	}
`)

	if string(got) != want {
		t.Errorf("got %s want %s", got, want)
	}
}
