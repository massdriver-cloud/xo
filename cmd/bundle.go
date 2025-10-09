package cmd

import (
	"fmt"
	"xo/src/bundle"
	"xo/src/telemetry"

	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/client"
	sdkbundle "github.com/massdriver-cloud/massdriver-sdk-go/massdriver/platform/bundle"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"oras.land/oras-go/v2/content/file"
)

var bundleCmd = &cobra.Command{
	Use:   "bundle",
	Short: "Bundle development tools",
	Long:  ``,
}

var bundlePullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls a bundle from Massdriver",
	RunE:  runBundlePull,
}

func init() {
	rootCmd.AddCommand(bundleCmd)

	bundleCmd.AddCommand(bundlePullCmd)
	bundlePullCmd.Flags().StringP("version", "v", "latest", "Bundle version (defaults to 'latest')")
	bundlePullCmd.Flags().StringP("name", "n", "", "Bundle name")
	viper.BindPFlag("bundle.version", bundlePullCmd.Flags().Lookup("version"))
	viper.BindPFlag("bundle.name", bundlePullCmd.Flags().Lookup("name"))
}

func runBundlePull(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runBundlePull")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	bundleName := viper.GetString("bundle.name")
	if bundleName == "" {
		return fmt.Errorf("required flag bundleName must be set via flag or environment variable")
	}
	bundleVersion := viper.GetString("bundle.version")

	mdClient, clientErr := client.New()
	if clientErr != nil {
		return clientErr
	}

	repo, repoErr := sdkbundle.GetBundleRepository(mdClient, bundleName)
	if repoErr != nil {
		return repoErr
	}

	fileStore, fileErr := file.New("bundle")
	if fileErr != nil {
		return fileErr
	}
	defer fileStore.Close()

	log.Info().Msgf("pulling bundle %s:%s...", bundleName, bundleVersion)
	desc, pullErr := bundle.Pull(ctx, repo, fileStore, bundleVersion)
	if pullErr != nil {
		return telemetry.LogError(span, pullErr, "an error occurred while pulling bundle")
	}
	log.Info().Msg("bundle pulled")
	log.Debug().Msg("bundle digest: " + desc.Digest.String())

	return nil
}
