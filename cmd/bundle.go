package cmd

import (
	"fmt"
	"path/filepath"
	"xo/src/bundles"
	"xo/src/generator"
	"xo/src/provisioners/terraform"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
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
	bundlePullCmd.Flags().StringP("bucket", "b", "xo-prod-bundlebucket-0000", "Bundle bucket")
	bundlePullCmd.Flags().StringP("key", "k", "k8s-application-aws.zip", "Path to bundle directory")
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
	bucket, err := cmd.Flags().GetString("bucket")
	if err != nil {
		return err
	}
	key, err := cmd.Flags().GetString("key")
	if err != nil {
		return err
	}

	err = bundles.Pull(bucket, key)
	if err != nil {
		return err
	}

	return nil
}
