package cmd

import (
	"os"
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
	Run:   runBundleBuild,
}

func init() {
	rootCmd.AddCommand(bundleCmd)
	bundleCmd.AddCommand(bundleBuildCmd)
	bundleBuildCmd.Flags().StringP("output", "o", ".", "Path to output directory.")
}

func runBundleBuild(cmd *cobra.Command, args []string) {
	var err error
	path := args[0]
	output, err := cmd.Flags().GetString("output")

	if err != nil {
		logger.Error("an error occurred while building bundle", zap.String("bundle", path), zap.Error(err))
		os.Exit(1)
	}

	logger.Info("building bundle",
		zap.String("bundle", path)
	)

	bundle, err := bundles.ParseBundle(path)
	if err != nil {
		logger.Error("an error occurred while building bundle", zap.String("bundle", path), zap.Error(err))
		os.Exit(1)
	}
	bundle.Build(output)

	logger.Info("bundle built",
		zap.String("bundle", path),
		zap.String("output", output),
	)

	return err
}
