package cmd

import (
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

func init() {
	rootCmd.AddCommand(bundleCmd)
	bundleCmd.AddCommand(bundleBuildCmd)
	bundleBuildCmd.Flags().StringP("output", "o", ".", "Path to output directory.")
}

func runBundleBuild(cmd *cobra.Command, args []string) error {
	var err error
	path := args[0]
	output, err := cmd.Flags().GetString("output")

	if err != nil {
		return err
	}

	logger.Info("building bundle",
		zap.String("source", path),
		zap.String("output", output),
	)

	bundle := bundles.ParseBundle(path)
	bundle.Build(output)

	logger.Info("bundle built",
		zap.String("source", path),
		zap.String("output", output),
	)

	return err
}
