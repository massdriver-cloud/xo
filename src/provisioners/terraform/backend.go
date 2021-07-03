package terraform

import (
	"encoding/json"
	"io"
	"os"
)

type TopLevelBlock struct {
	Terraform *TerraformBlock `json:"terraform,omitempty"`
}

type TerraformBlock struct {
	RequiredVersion string        `json:"required_version,omitempty"`
	BackendBlock    *BackendBlock `json:"backend,omitempty"`
}

type BackendBlock struct {
	S3BackendBlock *S3BackendBlock `json:"s3,omitempty"`
}

type S3BackendBlock struct {
	Bucket                string `json:"bucket"`
	DynamoDBTable         string `json:"dynamodb_table,omitempty"`
	Key                   string `json:"key"`
	Profile               string `json:"profile,omitempty"`
	Region                string `json:"region"`
	SharedCredentialsFile string `json:"shared_credentials_file,omitempty"`
}

func GenerateBackendS3File(output string, bucket string, mrn string, region string, dynamoDbTable string, sharedCredFile string, profile string) error {
	outputHandle, err := os.OpenFile(output, os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	config, err := GenerateJSONBackendS3Config(bucket, mrn, region, dynamoDbTable, sharedCredFile, profile)
	if err != nil {
		return err
	}

	return writeBackend(config, outputHandle)
}

func GenerateJSONBackendS3Config(bucket string, mrn string, region string, dynamoDbTable string, sharedCredFile string, profile string) ([]byte, error) {
	s3bb := new(S3BackendBlock)
	s3bb.Bucket = bucket
	s3bb.Key = mrn
	s3bb.Region = region
	s3bb.DynamoDBTable = dynamoDbTable
	s3bb.SharedCredentialsFile = sharedCredFile
	s3bb.Profile = profile

	topBlock := &TopLevelBlock{Terraform: &TerraformBlock{BackendBlock: &BackendBlock{S3BackendBlock: s3bb}}}

	return json.MarshalIndent(topBlock, "", "  ")
}

func writeBackend(config []byte, file io.Writer) error {
	_, err := file.Write(config)
	return err
}
