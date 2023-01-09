package cmd

import (
	"xo/src/massdriver"
	"xo/src/telemetry"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Manage secrets",
	Long:  ``,
}

var secretsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get secrets",
	RunE:  runSecretsGet,
}

func init() {
	rootCmd.AddCommand(secretsCmd)

	secretsCmd.AddCommand(secretsGetCmd)
}

func runSecretsGet(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runSecretsGet")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	mdClient, err := massdriver.InitializeMassdriverClient()
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while sending deployment status event: deploymentStatus")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	log.Info().Msg("Fetching secrets.")

	mdClient.GetSecrets(ctx)

	log.Info().Msg("Secrets retreived, writing to file.")

	return nil
}
