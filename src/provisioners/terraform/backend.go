package terraform

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"xo/src/massdriver"
	"xo/src/telemetry"

	"go.opentelemetry.io/otel"
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
	Bucket        string `json:"bucket"`
	DynamoDBTable string `json:"dynamodb_table,omitempty"`
	Key           string `json:"key"`
	Region        string `json:"region"`
}

func GenerateBackendS3File(ctx context.Context, output string, spec *massdriver.Specification, bundleStep string) error {
	_, span := otel.Tracer("xo").Start(ctx, "GenerateBackendS3File")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	outputHandle, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outputHandle.Close()

	config, err := GenerateJSONBackendS3Config(spec, bundleStep)
	if err != nil {
		return err
	}

	return writeBackend(config, outputHandle)
}

func GenerateJSONBackendS3Config(spec *massdriver.Specification, bundleStep string) ([]byte, error) {
	s3bb := new(S3BackendBlock)
	s3bb.Bucket = spec.S3StateBucket
	s3bb.Key = GetS3StateKey(spec.OrganizationID, spec.PackageID, bundleStep)
	s3bb.Region = spec.S3StateRegion
	s3bb.DynamoDBTable = spec.DynamoDBStateLockTable

	topBlock := &TopLevelBlock{Terraform: &TerraformBlock{BackendBlock: &BackendBlock{S3BackendBlock: s3bb}}}

	return json.MarshalIndent(topBlock, "", "  ")
}

func GetS3StateKey(organizationID, packageID, bundleStep string) string {
	return path.Join(organizationID, packageID, fmt.Sprintf("%s.tfstate", bundleStep))
}

func writeBackend(config []byte, out io.Writer) error {
	_, err := out.Write(config)
	return err
}
