package cmd

import (
	"xo/massdriver"

	"github.com/spf13/cobra"
)

var artifactUploadCmd = &cobra.Command{
	Use:                   "upload -i [id] -t [token]",
	Short:                 "Upload object to Massdriver",
	Long:                  ``,
	RunE:                  RunArtifactUpload,
	DisableFlagsInUseLine: true,
}

func init() {
	artifactCmd.AddCommand(artifactUploadCmd)
	artifactUploadCmd.Flags().StringP("id", "i", "", "Deployment ID")
	artifactUploadCmd.Flags().StringP("token", "t", "", "Secure token to authenticate with Massdriver")
	artifactUploadCmd.Flags().StringP("file", "f", "./artifact.json", "JSON formatted artifact file to upload")
	artifactUploadCmd.MarkFlagRequired("id")
	artifactUploadCmd.MarkFlagRequired("token")
}

func RunArtifactUpload(cmd *cobra.Command, args []string) error {
	id, _ := cmd.Flags().GetString("id")
	token, _ := cmd.Flags().GetString("token")
	file, _ := cmd.Flags().GetString("file")

	err := massdriver.UploadArtifactFile(file, id, token)
	return err
}
