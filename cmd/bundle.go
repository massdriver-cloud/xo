package cmd

import (
	"xo/src/bundle"
	"xo/src/massdriver"
	"xo/src/telemetry"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

var bundleCmd = &cobra.Command{
	Use:   "bundle",
	Short: "Bundle development tools",
	Long:  ``,
}

var bundlePullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls a bundle from S3",
	RunE:  runBundlePull,
}

func init() {
	rootCmd.AddCommand(bundleCmd)

	bundleCmd.AddCommand(bundlePullCmd)
}

func runBundlePull(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runBundlePull")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	client, initErr := massdriver.InitializeMassdriverClient()
	if initErr != nil {
		log.Error().Err(initErr).Msg("an error occurred while initializing massdriver client")
		span.RecordError(initErr)
		span.SetStatus(codes.Error, initErr.Error())
		return initErr
	}

	log.Info().Msg("pulling bundle")
	pullErr := bundle.Pull(ctx, client)
	if pullErr != nil {
		log.Error().Err(pullErr).Msg("an error occurred while pulling bundle")
		span.RecordError(pullErr)
		span.SetStatus(codes.Error, pullErr.Error())
		return pullErr
	}
	log.Info().Msg("bundle pulled")

	return nil
}
