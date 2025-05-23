package terraform

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
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
	HTTPBackendBlock *HTTPBackendBlock `json:"http,omitempty"`
}

type HTTPBackendBlock struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Address  string `json:"address"`
	// LockAddress   string `json:"lock_address"`
	// UnlockAddress string `json:"unlock_address"`
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

	httpbb.Username = spec.DeploymentID
	httpbb.Password = spec.Token
	httpbb.Address = fmt.Sprintf("https://api.massdriver.cloud/state/%s/%s", getPackageNameShort(spec.PackageName), bundleStep)
	// httpbb.LockAddress = httpbb.Address
	// httpbb.UnlockAddress = httpbb.Address

	topBlock := &TopLevelBlock{Terraform: &TerraformBlock{BackendBlock: &BackendBlock{HTTPBackendBlock: httpbb}}}

	return json.MarshalIndent(topBlock, "", "  ")
}

func getPackageNameShort(packageId string) string {
	parts := strings.Split(packageId, "-")
	return strings.Join(parts[:len(parts)-1], "-")
}
