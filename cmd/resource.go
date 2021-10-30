package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var resourceUpdateLong = `
	Sends massdriver and update about the status of a resource during a deployment.

	In order to update users about resources that are provisioned as part of a bundle
	we want to send status updates as they occur.
	`
var resourceUpdateExamples = `
	# Update resource (variabes in environment)
	xo resource update

	# Update resource manually specifying command-line variables
	xo resource update -i <deployment-id> -t <token> -r <resource-id> -y <resource-type> -s <resource-status>
	`

var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "Manage Massdriver bundle resources",
	Long:  ``,
}

var resourceUpdateCmd = &cobra.Command{
	Use:                   "update",
	Short:                 "Send Massdriver and update regarding resource status",
	Long:                  resourceUpdateLong,
	Example:               resourceUpdateExamples,
	RunE:                  RunResourceUpdate,
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.AddCommand(resourceCmd)
	resourceCmd.AddCommand(resourceUpdateCmd)
	resourceUpdateCmd.Flags().StringP("deployment-id", "i", os.Getenv("MASSDRIVER_DEPLOYMENT_ID"), "Massdriver Deployment ID. Defaults to value in MASSDRIVER_DEPLOYMENT_ID environment variable.")
	resourceUpdateCmd.Flags().StringP("token", "t", os.Getenv("MASSDRIVER_TOKEN"), "Secure token to authenticate with Massdriver. Defaults to value in MASSDRIVER_TOKEN environment variable.")
	resourceUpdateCmd.Flags().StringP("resource-id", "r", os.Getenv("RESOURCE_ID"), "ID of resource. Defaults to value in RESOURCE_ID environment variable.")
	resourceUpdateCmd.Flags().StringP("resource-type", "y", os.Getenv("RESOURCE_TYPE"), "Type name of resource. Defaults to value in RESOURCE_TYPE environment variable.")
	resourceUpdateCmd.Flags().StringP("resource-status", "s", os.Getenv("RESOURCE_STATUS"), "Status code of resource. Defaults to value in RESOURCE_STATUS environment variable.")
}

func RunResourceUpdate(cmd *cobra.Command, args []string) error {
	deploymentId, _ := cmd.Flags().GetString("deployment-id")
	token, _ := cmd.Flags().GetString("token")
	resourceId, _ := cmd.Flags().GetString("resource-id")
	resourceType, _ := cmd.Flags().GetString("resource-type")
	resourceStatus, _ := cmd.Flags().GetString("resource-status")

	if deploymentId == "" || token == "" {
		cmd.Help()
		fmt.Println("both deployment-id and token must be set (by flags or environment variable)")
		return errors.New("ERROR: Both deployment-id and token must be set (by flags or environment variable)")
	}

	log.Info().
		Str("deployment", deploymentId).
		Str("resource-id", resourceId).
		Str("resource-type", resourceType).
		Str("resource-status", resourceStatus).
		Msg("sending resource update to Massdriver")
	// err := massdriver.UpdateResource(deploymentId, token, resourceId, resourceType, resourceStatus)
	// if err != nil {
	// 	log.Error().Err(err).
	// 		Str("deployment", deploymentId).
	// 		Str("resource-id", resourceId).
	// 		Str("resource-type", resourceType).
	// 		Str("resource-status", resourceStatus).
	// 		Msg("an error occurred while sending resource update to Massdriver")
	// 	return err
	// }
	log.Info().
		Str("deployment", deploymentId).
		Str("resource-id", resourceId).
		Str("resource-type", resourceType).
		Str("resource-status", resourceStatus).
		Msg("resource update complete")

	return nil
}
