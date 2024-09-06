package terraform

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"xo/src/massdriver"
	"xo/src/telemetry"

	"github.com/massdriver-cloud/terraform-config-inspect/tfconfig"
	"github.com/rs/zerolog/log"
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
	S3BackendBlock   *S3BackendBlock   `json:"s3,omitempty"`
	HTTPBackendBlock *HTTPBackendBlock `json:"http,omitempty"`
}

type S3BackendBlock struct {
	Bucket        string `json:"bucket"`
	DynamoDBTable string `json:"dynamodb_table,omitempty"`
	Key           string `json:"key"`
	Region        string `json:"region"`
}

type HTTPBackendBlock struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Address       string `json:"address"`
	LockAddress   string `json:"lock_address"`
	UnlockAddress string `json:"unlock_address"`
}

func GenerateBackendS3File(ctx context.Context, output string, spec *massdriver.Specification, bundleStep string) error {
	_, span := otel.Tracer("xo").Start(ctx, "GenerateBackendS3File")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	skip, err := hasExistingBackendConfig(".")
	if err != nil {
		return err
	}
	if skip {
		log.Info().Msg("Existing backend configuration detected. Skipping generation.")
		return nil
	}

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
	s3bb.DynamoDBTable = getDynamoDBTableNameFromARN(spec.DynamoDBStateLockTableArn)

	topBlock := &TopLevelBlock{Terraform: &TerraformBlock{BackendBlock: &BackendBlock{S3BackendBlock: s3bb}}}

	return json.MarshalIndent(topBlock, "", "  ")
}

func GetS3StateKey(organizationID, packageID, bundleStep string) string {
	return path.Join(organizationID, packageID, fmt.Sprintf("%s.tfstate", bundleStep))
}

func getDynamoDBTableNameFromARN(dynamoDBTableARN string) string {
	return strings.Split(dynamoDBTableARN, ":table/")[1]
}

func writeBackend(config []byte, out io.Writer) error {
	_, err := out.Write(config)
	return err
}

func hasExistingBackendConfig(path string) (bool, error) {
	module, diag := tfconfig.LoadModule(path)
	if diag != nil && diag.HasErrors() {
		return false, diag
	}

	hasBackendConfig := module.Backend != nil

	return hasBackendConfig, nil
}

func GenerateBackendHTTPFile(ctx context.Context, output string, spec *massdriver.Specification, bundleStep string) error {
	_, span := otel.Tracer("xo").Start(ctx, "GenerateBackendHTTPFile")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	skip, err := hasExistingBackendConfig(".")
	if err != nil {
		return err
	}
	if skip {
		log.Info().Msg("Existing backend configuration detected. Skipping generation.")
		return nil
	}

	outputHandle, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outputHandle.Close()

	config, err := GenerateJSONBackendHTTPConfig(spec, bundleStep)
	if err != nil {
		return err
	}

	return writeBackend(config, outputHandle)
}

func GenerateJSONBackendHTTPConfig(spec *massdriver.Specification, bundleStep string) ([]byte, error) {
	httpbb := new(HTTPBackendBlock)

	httpbb.Username = spec.OrganizationID
	httpbb.Password = spec.Token
	httpbb.Address = fmt.Sprintf("https://api.massdriver.cloud/state/%s/%s", spec.PackageID, bundleStep)
	httpbb.LockAddress = httpbb.Address
	httpbb.UnlockAddress = httpbb.Address

	topBlock := &TopLevelBlock{Terraform: &TerraformBlock{BackendBlock: &BackendBlock{HTTPBackendBlock: httpbb}}}

	return json.MarshalIndent(topBlock, "", "  ")
}
