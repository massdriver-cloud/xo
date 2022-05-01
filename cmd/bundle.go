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
}

func runBundleBuild(cmd *cobra.Command, args []string) error {
	var err error
	bundlePath := args[0]

	// default the output to the path of the bundle.yaml file
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		log.Error().Err(err).Str("bundle", bundlePath).Msg("an error occurred while building bundle")
		return err
	}
	if output == "" {
		output = filepath.Dir(bundlePath)
	}

	log.Info().Str("bundle", bundlePath).Msg("building bundle")

	bundle, err := bundles.ParseBundle(bundlePath)
	if err != nil {
		log.Error().Err(err).Str("bundle", bundlePath).Msg("an error occurred while building bundle")
		return err
	}

	err = bundle.GenerateSchemas(output)
	if err != nil {
		log.Error().Err(err).Str("bundle", bundlePath).Msg("an error occurred while generating bundle schema files")
		return err
	}

	for _, step := range bundle.Steps {
		switch step.Provisioner {
		case "terraform":
			err = terraform.GenerateFiles(output, step.Path)
			if err != nil {
				log.Error().Err(err).Str("bundle", bundlePath).Str("provisioner", step.Provisioner).Msg("an error occurred while generating provisioner files")
				return err
			}
		case "exec":
			// No-op (Golang doesn't not fallthrough unless explicitly stated)
		default:
			log.Error().Str("bundle", bundlePath).Msg("unknown provisioner: " + step.Provisioner)
			return fmt.Errorf("unknown provisioner: %v", step.Provisioner)
		}
	}

	log.Info().Str("bundle", bundlePath).Str("output", output).Msg("bundle built")

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

	templateData := generator.TemplateData{
		BundleDir:   bundleDir,
		TemplateDir: templateDir,
		Type:        "bundle",
	}

	err = generator.RunPrompt(&templateData)
	if err != nil {
		return err
	}

	err = generator.Generate(&templateData)
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

	bundleId := os.Getenv("MASSDRIVER_BUNDLE_ID")
	if bundleId == "" {
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

	log.Info().Msg("pulling bundle")
	err := bundles.Pull(ctx, bundleBucket, organizationId, bundleId)
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while pulling bundle")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	log.Info().Msg("bundle pulled")

	return nil
}
