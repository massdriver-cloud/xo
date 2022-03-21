package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"xo/src/bundles"
	"xo/src/generator"
	"xo/src/provisioners/terraform"

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

var bundleBuildCmd = &cobra.Command{
	Use:   "build [Path to bundle.yaml]",
	Short: "Builds bundle JSON Schemas",
	Args:  cobra.ExactArgs(1),
	RunE:  runBundleBuild,
}

var bundleGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates a new bundle",
	RunE:  runBundleGenerate,
}

var bundlePullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls a bundle from S3",
	RunE:  runBundlePull,
}

func init() {
	rootCmd.AddCommand(bundleCmd)

	bundleCmd.AddCommand(bundleBuildCmd)
	bundleBuildCmd.Flags().StringP("output", "o", "", "Path to output directory (default is bundle.yaml directory)")

	bundleCmd.AddCommand(bundleGenerateCmd)
	bundleGenerateCmd.Flags().StringP("template-dir", "t", "./generators/xo-bundle-template", "Path to template directory")
	bundleGenerateCmd.Flags().StringP("bundle-dir", "b", "./bundles", "Path to bundle directory")

	bundleCmd.AddCommand(bundlePullCmd)
	// bundlePullCmd.Flags().StringP("bucket", "b", "", "Bundle bucket")
	// bundlePullCmd.Flags().StringP("type", "t", "public", "Bundle type (public or private)")
	// bundlePullCmd.Flags().StringP("organization", "o", "", "Organization ID (required if private)")
	// bundlePullCmd.Flags().StringP("name", "n", "", "Bundle name")
	// bundlePullCmd.MarkFlagRequired("name")
}

func runBundleBuild(cmd *cobra.Command, args []string) error {
	var err error
	path := args[0]

	// default the output to the path of the bundle.yaml file
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		log.Error().Err(err).Str("bundle", path).Msg("an error occurred while building bundle")
		return err
	}
	if output == "" {
		output = filepath.Dir(path)
	}

	log.Info().Str("bundle", path).Msg("building bundle")

	bundle, err := bundles.ParseBundle(path)
	if err != nil {
		log.Error().Err(err).Str("bundle", path).Msg("an error occurred while building bundle")
		return err
	}

	err = bundle.GenerateSchemas(output)
	if err != nil {
		log.Error().Err(err).Str("bundle", path).Msg("an error occurred while generating bundle schema files")
		return err
	}

	switch bundle.Provisioner {
	case "terraform":
		err = terraform.GenerateFiles(output)
		if err != nil {
			log.Error().Err(err).Str("bundle", path).Str("provisioner", bundle.Provisioner).Msg("an error occurred while generating provisioner files")
			return err
		}
	default:
		log.Error().Str("bundle", path).Msg("unknown provisioner: " + bundle.Provisioner)
		return fmt.Errorf("unknown provisioner: %v", bundle.Provisioner)
	}

	log.Info().Str("bundle", path).Str("output", output).Msg("bundle built")

	return err
}

func runBundleGenerate(cmd *cobra.Command, args []string) error {
	var err error

	bundleDir, err := cmd.Flags().GetString("bundle-dir")
	if err != nil {
		return err
	}

	templateDir, err := cmd.Flags().GetString("template-dir")
	if err != nil {
		return err
	}

	templateData := &generator.TemplateData{
		BundleDir:   bundleDir,
		TemplateDir: templateDir,
	}

	err = generator.RunPrompt(templateData)
	if err != nil {
		return err
	}

	err = generator.Generate(*templateData)
	if err != nil {
		return err
	}

	return nil
}

func runBundlePull(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(context.Background(), "RunDeploymentStatus")
	defer span.End()

	bundleBucket := os.Getenv("MASSDRIVER_BUNDLE_BUCKET")
	if bundleBucket == "" {
		err := errors.New("MASSDRIVER_BUNDLE_BUCKET environment variable must be set")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	bundleAccess := os.Getenv("MASSDRIVER_BUNDLE_ACCESS")
	if bundleAccess == "" {
		err := errors.New("MASSDRIVER_BUNDLE_ACCESS environment variable must be set")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	bundleName := os.Getenv("MASSDRIVER_BUNDLE_NAME")
	if bundleName == "" {
		err := errors.New("MASSDRIVER_BUNDLE_NAME environment variable must be set")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	organizationId := os.Getenv("MASSDRIVER_ORGANIZATION_ID")
	if organizationId == "" {
		err := errors.New("MASSDRIVER_ORGANIZATION_ID environment variable must be set")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	err := bundles.Pull(ctx, bundleBucket, bundleAccess, bundleName, organizationId)
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while pulling bundle: " + bundleName)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
