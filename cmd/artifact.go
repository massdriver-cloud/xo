package cmd

import (
	"errors"
	"fmt"
	"os"
	"xo/src/massdriver"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var artifactUploadLong = `
	Uploads artifact data about a deployment to Massdriver.

	This is intended to be used as a step in workflow execution to update 
	metadata after provisioning.
	`
var artifactUploadExamples = `
	# Upload artifact (deployment-id and token in environment)
	xo artifact upload

	# Upload artifact manually specifying deployment-id and token
	xo artifact upload -i <deployment-id> -t <token>

	# Upload artifacts in custom file
	xo artifact upload -f /tmp/custom-artifacts.json
	`

var artifactCmd = &cobra.Command{
	Use:   "artifact",
	Short: "Manage Massdriver artifacts",
	Long:  ``,
}

var artifactUploadCmd = &cobra.Command{
	Use:                   "upload",
	Short:                 "Upload artifact to Massdriver",
	Long:                  artifactUploadLong,
	Example:               artifactUploadExamples,
	RunE:                  RunArtifactUpload,
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.AddCommand(artifactCmd)
	artifactCmd.AddCommand(artifactUploadCmd)
	artifactUploadCmd.Flags().StringP("deployment-id", "i", os.Getenv("MASSDRIVER_DEPLOYMENT_ID"), "Massdriver Deployment ID. Defaults to value in MASSDRIVER_DEPLOYMENT_ID environment variable.")
	artifactUploadCmd.Flags().StringP("token", "t", os.Getenv("MASSDRIVER_TOKEN"), "Secure token to authenticate with Massdriver. Defaults to value in MASSDRIVER_TOKEN environment variable.")
	artifactUploadCmd.Flags().StringP("file", "f", "./artifact.json", "JSON formatted artifact file to upload")
}

func RunArtifactUpload(cmd *cobra.Command, args []string) error {
	id, _ := cmd.Flags().GetString("deployment-id")
	token, _ := cmd.Flags().GetString("token")
	file, _ := cmd.Flags().GetString("file")

	if id == "" || token == "" {
		cmd.Help()
		fmt.Println("\nERROR: Both deployment-id and token must be set (by flags or environment variable)")
		return errors.New("both deployment-id and token must be set (by flags or environment variable)")
	}

	logger.Info("uploading artifact file", zap.String("deployment", id))
	err := massdriver.UploadArtifactFile(file, id, token)
	if err != nil {
		logger.Error("an error occurred while uploading artifact files", zap.String("deployment", id), zap.Error(err))
		return err
	}
	logger.Info("artifact uploaded", zap.String("deployment", id))

	return nil
}
