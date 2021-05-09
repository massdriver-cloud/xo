package cmd

import (
	"os"
	"xo/src/massdriver"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var artifactCmd = &cobra.Command{
	Use:   "artifact",
	Short: "Manage Massdriver artifacts",
	Long:  ``,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("artifact called")
	//},
}

var artifactUploadCmd = &cobra.Command{
	Use:                   "upload -i [id] -t [token]",
	Short:                 "Upload artifact to Massdriver",
	Long:                  ``,
	RunE:                  RunArtifactUpload,
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.AddCommand(artifactCmd)
	artifactCmd.AddCommand(artifactUploadCmd)
	artifactUploadCmd.Flags().StringP("deployment-id", "i", os.Getenv("MASSDRIVER_DEPLOYMENT_ID"), "Deployment ID")
	artifactUploadCmd.Flags().StringP("token", "t", os.Getenv("MASSDRIVER_TOKEN"), "Secure token to authenticate with Massdriver")
	artifactUploadCmd.Flags().StringP("file", "f", "./artifact.json", "JSON formatted artifact file to upload")
	artifactUploadCmd.MarkFlagRequired("id")
	artifactUploadCmd.MarkFlagRequired("token")
}

func RunArtifactUpload(cmd *cobra.Command, args []string) error {
	id, _ := cmd.Flags().GetString("deployment-id")
	token, _ := cmd.Flags().GetString("token")
	file, _ := cmd.Flags().GetString("file")

	logger.Info("Uploading artifact file",
		zap.String("deployment", id),
	)
	err := massdriver.UploadArtifactFile(file, id, token)
	return err
}
