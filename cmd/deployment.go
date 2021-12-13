package cmd

import (
	"errors"
	"fmt"
	"os"
	"xo/src/massdriver"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var deploymentStartLong = `
	Fetches metadata about a Massdriver deployment and writes the data to files.

	Specifically this is fetching the params and connections and writing them to
	params.tfvars.json and connections.tfvars.json, respectively. This is intended
	to be used as a step in workflow execution to gather resources for the provisioner.
	`
var deploymentStartExamples = `
	# Get deployment (deployment-id and token in environment)
	xo deployment get

	# Get deployment manually specifying deployment-id and token
	xo deployment get -i <deployment-id> -t <token>

	# Get deployment and write files to /tmp
	xo deployment get -d /tmp
	`

var deploymentCompleteLong = `
	Uploads artifact data about a deployment to Massdriver.

	This is intended to be used as a step in workflow execution to update 
	metadata after provisioning.
	`
var deploymentCompleteExamples = `
	# Upload artifact (deployment-id and token in environment)
	xo deployment complete

	# Upload artifact manually specifying deployment-id and token
	xo deployment complete -i <deployment-id> -t <token>

	# Upload artifacts in custom file
	xo deployment complete -f /tmp/custom-artifacts.json
	`

var deploymentFailLong = `
	Reports the deployment has failed to Massdriver
	`

var deploymentFailExamples = `
	# Fail deployment (deployment-id and token in environment)
	xo deployment fail

	# Fail deployment specifying deployment-id and token
	xo deployment fail -i <deployment-id> -t <token>
	`

var deploymentDestroyedLong = `
	Reports the deployment has been destroyed to Massdriver
	`

var deploymentDestroyedExamples = `
	# Destroyed deployment (deployment-id and token in environment)
	xo deployment destroyed

	# Destroyed deployment specifying deployment-id and token
	xo deployment destroyed -i <deployment-id> -t <token>
	`

var deploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "Manage Massdriver deployments",
	Long:  ``,
}

var deploymentStartCmd = &cobra.Command{
	Use:                   "start",
	Short:                 "Start Massdriver deployment",
	Long:                  deploymentStartLong,
	Example:               deploymentStartExamples,
	RunE:                  RunDeploymentStart,
	DisableFlagsInUseLine: true,
}

var deploymentCompleteCmd = &cobra.Command{
	Use:                   "complete",
	Short:                 "Complete Massdriver deployment",
	Long:                  deploymentCompleteLong,
	Example:               deploymentCompleteExamples,
	RunE:                  RunDeploymentComplete,
	DisableFlagsInUseLine: true,
}

var deploymentFailCmd = &cobra.Command{
	Use:                   "fail",
	Short:                 "Report Massdriver deployment has failed",
	Long:                  deploymentFailLong,
	Example:               deploymentFailExamples,
	RunE:                  RunDeploymentFail,
	DisableFlagsInUseLine: true,
}

var deploymentDestroyedCmd = &cobra.Command{
	Use:                   "destroyed",
	Short:                 "Report Massdriver deployment has been destroyed",
	Long:                  deploymentDestroyedLong,
	Example:               deploymentDestroyedExamples,
	RunE:                  RunDeploymentDestroyed,
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.AddCommand(deploymentCmd)

	deploymentCmd.AddCommand(deploymentStartCmd)
	deploymentStartCmd.Flags().StringP("deployment-id", "i", os.Getenv("MASSDRIVER_DEPLOYMENT_ID"), "Massdriver Deployment ID. Defaults to value in MASSDRIVER_DEPLOYMENT_ID environment variable.")
	deploymentStartCmd.Flags().StringP("token", "t", os.Getenv("MASSDRIVER_TOKEN"), "Secure token to authenticate with Massdriver. Defaults to value in MASSDRIVER_TOKEN environment variable.")
	deploymentStartCmd.Flags().StringP("out", "o", ".", "Destination path to write deployment json files. Defaults to current directory")

	deploymentCmd.AddCommand(deploymentCompleteCmd)
	deploymentCompleteCmd.Flags().StringP("deployment-id", "i", os.Getenv("MASSDRIVER_DEPLOYMENT_ID"), "Massdriver Deployment ID. Defaults to value in MASSDRIVER_DEPLOYMENT_ID environment variable.")
	deploymentCompleteCmd.Flags().StringP("token", "t", os.Getenv("MASSDRIVER_TOKEN"), "Secure token to authenticate with Massdriver. Defaults to value in MASSDRIVER_TOKEN environment variable.")
	deploymentCompleteCmd.Flags().StringP("artifacts", "f", "./artifact.json", "Path to JSON formatted artifact file to upload. Defaults to ./artifact.json")

	deploymentCmd.AddCommand(deploymentFailCmd)
	deploymentFailCmd.Flags().StringP("deployment-id", "i", os.Getenv("MASSDRIVER_DEPLOYMENT_ID"), "Massdriver Deployment ID. Defaults to value in MASSDRIVER_DEPLOYMENT_ID environment variable.")
	deploymentFailCmd.Flags().StringP("token", "t", os.Getenv("MASSDRIVER_TOKEN"), "Secure token to authenticate with Massdriver. Defaults to value in MASSDRIVER_TOKEN environment variable.")

	deploymentCmd.AddCommand(deploymentDestroyedCmd)
	deploymentDestroyedCmd.Flags().StringP("deployment-id", "i", os.Getenv("MASSDRIVER_DEPLOYMENT_ID"), "Massdriver Deployment ID. Defaults to value in MASSDRIVER_DEPLOYMENT_ID environment variable.")
	deploymentDestroyedCmd.Flags().StringP("token", "t", os.Getenv("MASSDRIVER_TOKEN"), "Secure token to authenticate with Massdriver. Defaults to value in MASSDRIVER_TOKEN environment variable.")
}

func RunDeploymentStart(cmd *cobra.Command, args []string) error {
	id, _ := cmd.Flags().GetString("deployment-id")
	token, _ := cmd.Flags().GetString("token")
	out, _ := cmd.Flags().GetString("out")

	if id == "" || token == "" {
		cmd.Help()
		fmt.Println("both deployment-id and token must be set (by flags or environment variable)")
		return errors.New("ERROR: Both deployment-id and token must be set (by flags or environment variable)")
	}

	log.Info().Str("deployment", id).Msg("getting deployment from Massdriver")
	err := massdriver.StartDeployment(id, token, out)
	if err != nil {
		log.Error().Err(err).Str("deployment", id).Msg("an error occurred while getting deployment from Massdriver")
		return err
	}

	log.Info().Str("deployment", id).Msg("deployment get complete")

	return nil
}

func RunDeploymentComplete(cmd *cobra.Command, args []string) error {
	id, _ := cmd.Flags().GetString("deployment-id")
	token, _ := cmd.Flags().GetString("token")
	artifacts, _ := cmd.Flags().GetString("artifacts")

	if id == "" || token == "" {
		cmd.Help()
		fmt.Println("\nERROR: Both deployment-id and token must be set (by flags or environment variable)")
		return errors.New("both deployment-id and token must be set (by flags or environment variable)")
	}

	log.Info().Str("deployment", id).Msg("uploading artifact file to Massdriver")
	err := massdriver.UploadArtifactFile(artifacts, id, token)
	if err != nil {
		log.Error().Err(err).Str("deployment", id).Msg("an error occurred while uploading artifact files")
		return err
	}
	log.Info().Str("deployment", id).Msg("artifact uploaded")

	return nil
}

func RunDeploymentFail(cmd *cobra.Command, args []string) error {
	id, _ := cmd.Flags().GetString("deployment-id")
	token, _ := cmd.Flags().GetString("token")

	if id == "" || token == "" {
		cmd.Help()
		fmt.Println("\nERROR: Both deployment-id and token must be set (by flags or environment variable)")
		return errors.New("both deployment-id and token must be set (by flags or environment variable)")
	}

	log.Info().Str("deployment", id).Msg("reporting deployed has failed to Massdriver")
	err := massdriver.FailDeployment(id, token)
	if err != nil {
		log.Error().Err(err).Str("deployment", id).Msg("an error occurred while reporting deployed has failed")
		return err
	}
	log.Info().Str("deployment", id).Msg("failed deployment reported")

	return nil
}

func RunDeploymentDestroyed(cmd *cobra.Command, args []string) error {
	id, _ := cmd.Flags().GetString("deployment-id")
	token, _ := cmd.Flags().GetString("token")

	if id == "" || token == "" {
		cmd.Help()
		fmt.Println("\nERROR: Both deployment-id and token must be set (by flags or environment variable)")
		return errors.New("both deployment-id and token must be set (by flags or environment variable)")
	}

	log.Info().Str("deployment", id).Msg("reporting deployment has been destroyed to Massdriver")
	err := massdriver.DestroyedDeployment(id, token)
	if err != nil {
		log.Error().Err(err).Str("deployment", id).Msg("an error occurred while reporting deployed has been destroyed")
		return err
	}
	log.Info().Str("deployment", id).Msg("destroyed deployment reported")

	return nil
}
