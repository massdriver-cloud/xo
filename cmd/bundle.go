package cmd

import (
	"xo/src/bundles"
	"xo/src/generator"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
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

func init() {
	rootCmd.AddCommand(bundleCmd)
	bundleCmd.AddCommand(bundleBuildCmd)
	bundleBuildCmd.Flags().StringP("output", "o", ".", "Path to output directory.")
	bundleCmd.AddCommand(bundleGenerateCmd)
	bundleGenerateCmd.Flags().StringP("template-dir", "t", "./generators/xo-bundle-template", "Path to template directory")
	bundleGenerateCmd.Flags().StringP("bundle-dir", "b", "./bundles", "Path to bundle directory")
}

func runBundleBuild(cmd *cobra.Command, args []string) error {
	var err error
	path := args[0]
	output, err := cmd.Flags().GetString("output")

	if err != nil {
		logger.Error("an error occurred while building bundle", zap.String("bundle", path), zap.Error(err))
		return err
	}

	logger.Info("building bundle",
		zap.String("bundle", path),
	)

	bundle, err := bundles.ParseBundle(path)
	if err != nil {
		logger.Error("an error occurred while building bundle", zap.String("bundle", path), zap.Error(err))
		return err
	}

	err = bundle.Build(output)
	if err != nil {
		logger.Error("an error occurred while building bundle", zap.String("bundle", path), zap.Error(err))
		return err
	}

	logger.Info("bundle built",
		zap.String("bundle", path),
		zap.String("output", output),
	)

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
