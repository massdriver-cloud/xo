package cmd

import (
	"fmt"
	"xo/src/bundles"

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
	Use:   "generate [Bundle name]",
	Short: "Generates a new bundle",
	RunE:  runBundleGenerate,
}

func init() {
	rootCmd.AddCommand(bundleCmd)
	bundleCmd.AddCommand(bundleBuildCmd)
	bundleBuildCmd.Flags().StringP("output", "o", ".", "Path to output directory.")
	bundleCmd.AddCommand(bundleGenerateCmd)
	bundleGenerateCmd.Flags().StringP("template-dir", "t", "./xo-bundle-template", "Path to template directory")
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
	bundle.Build(output)

	logger.Info("bundle built",
		zap.String("bundle", path),
		zap.String("output", output),
	)

	return err
}

func runBundleGenerate(cmd *cobra.Command, args []string) error {
	fmt.Println(args)
	return nil
}
