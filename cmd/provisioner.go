package cmd

import (
	"xo/src/massdriver"
	tf "xo/src/provisioners/terraform"
	"xo/src/telemetry"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
)

var provisionerCmd = &cobra.Command{
	Use:   "provisioner",
	Short: "Manage provisioners",
	Long:  ``,
}

var provisionerTerraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Commands specific to terraform provisioner",
	Long:  ``,
}

var provisionerTerraformBackendCmd = &cobra.Command{
	Use:   "backend",
	Short: "Generate a terraform backend config",
	Long:  ``,
}

var provisionerTerraformBackendHTTPCmd = &cobra.Command{
	Use:   "http",
	Short: "Generate a terraform HTTP backend config",
	Long:  ``,
	RunE:  runProvisionerTerraformBackendHTTP,
}

func init() {
	rootCmd.AddCommand(provisionerCmd)

	provisionerCmd.AddCommand(provisionerTerraformCmd)

	provisionerTerraformCmd.AddCommand(provisionerTerraformBackendCmd)
	provisionerTerraformBackendCmd.PersistentFlags().StringP("output", "o", "./backend.tf.json", "Output file path")

	provisionerTerraformBackendCmd.AddCommand(provisionerTerraformBackendHTTPCmd)
	provisionerTerraformBackendHTTPCmd.Flags().StringP("step", "s", "", "Bundle Step")
	provisionerTerraformBackendHTTPCmd.MarkFlagRequired("step")
}

func runProvisionerTerraformBackendHTTP(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runProvisionerTerraformBackendHTTP")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	output, _ := cmd.Flags().GetString("output")
	step, _ := cmd.Flags().GetString("step")

	spec, specErr := massdriver.GetSpecification()
	if specErr != nil {
		return telemetry.LogError(span, specErr, "an error occurred while extracting Massdriver specification")
	}

	log.Info().Msg("Generating state file...")

	generateErr := tf.GenerateBackendHTTPFile(ctx, output, spec, step)
	if generateErr != nil {
		return telemetry.LogError(span, generateErr, "an error occurred while generating backend file")
	}
	log.Info().Msg("State file generated")
	return nil
}
