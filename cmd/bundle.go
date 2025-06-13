package cmd

import (
	"fmt"
	"os"
	"xo/src/bundle"
	"xo/src/massdriver"
	"xo/src/telemetry"

	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/client"
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

var bundlePullv0Cmd = &cobra.Command{
	Use:   "pullv0",
	Short: "Pulls a bundle from Massdriver",
	RunE:  runBundlePullv0,
}

var bundlePullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls a bundle from Massdriver",
	RunE:  runBundlePull,
}

func init() {
	rootCmd.AddCommand(bundleCmd)

	bundleCmd.AddCommand(bundlePullv0Cmd)

	bundleCmd.AddCommand(bundlePullCmd)
	bundlePullCmd.Flags().StringP("organization", "o", "", "Organization slug")
	bundlePullCmd.Flags().StringP("tag", "t", "latest", "Bundle tag (defaults to 'latest')")
	bundlePullCmd.Flags().StringP("name", "n", "", "Bundle name")
	viper.BindPFlag("bundle.tag", bundlePullCmd.Flags().Lookup("tag"))
	viper.BindPFlag("bundle.name", bundlePullCmd.Flags().Lookup("name"))
	viper.BindPFlag("organization.slug", bundlePullCmd.Flags().Lookup("organization"))
}

func runBundlePullv0(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runBundlePull")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	client, initErr := massdriver.InitializeMassdriverClient()
	if initErr != nil {
		return telemetry.LogError(span, initErr, "an error occurred while initializing massdriver client")
	}

	outFile, openErr := os.OpenFile("bundle.tar.gz", os.O_CREATE|os.O_WRONLY, 0644)
	if openErr != nil {
		return telemetry.LogError(span, openErr, "unable to open bundle.tar.gz file")
	}
	defer outFile.Close()

	log.Info().Msg("pulling bundle...")
	pullErr := bundle.PullV0(ctx, client, outFile)
	if pullErr != nil {
		return telemetry.LogError(span, pullErr, "an error occurred while pulling bundle")
	}
	log.Info().Msg("bundle pulled")

	return nil
}

func runBundlePull(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runBundlePull")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	organizationSlug := viper.GetString("organization.slug")
	if organizationSlug == "" {
		return fmt.Errorf("required flag organization must be set via flag or environment variable")
	}
	bundleName := viper.GetString("bundle.name")
	if bundleName == "" {
		return fmt.Errorf("required flag bundleName must be set via flag or environment variable")
	}
	tag := viper.GetString("bundle.tag")

	mdClient, clientErr := client.New()
	if clientErr != nil {
		return clientErr
	}

	repo, repoErr := bundle.GetRepo(mdClient, organizationSlug, bundleName)
	if repoErr != nil {
		return repoErr
	}

	fileStore, fileErr := file.New("bundle")
	if fileErr != nil {
		return fileErr
	}
	defer fileStore.Close()

	log.Info().Msg("pulling bundle...")
	desc, pullErr := bundle.PullV1(ctx, repo, fileStore, tag)
	if pullErr != nil {
		return telemetry.LogError(span, pullErr, "an error occurred while pulling bundle")
	}
	log.Info().Msg("bundle pulled")
	log.Debug().Msg("bundle digest: " + desc.Digest.String())

	return nil
}
