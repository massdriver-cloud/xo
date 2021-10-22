package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"xo/src/massdriver"
	"xo/src/provisioners"

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
	xo artifact upload

	# Upload artifact manually specifying deployment-id and token
	xo artifact upload -i <deployment-id> -t <token>

	# Upload artifacts in custom file
	xo artifact upload -f /tmp/custom-artifacts.json
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

func init() {
	rootCmd.AddCommand(deploymentCmd)

	deploymentCmd.AddCommand(deploymentStartCmd)
	deploymentStartCmd.Flags().StringP("bundle", "b", ".", "Path to the bundle to execute. Defaults to current directory")
	deploymentStartCmd.Flags().StringP("deployment-id", "i", os.Getenv("MASSDRIVER_DEPLOYMENT_ID"), "Massdriver Deployment ID. Defaults to value in MASSDRIVER_DEPLOYMENT_ID environment variable.")
	deploymentStartCmd.Flags().StringP("token", "t", os.Getenv("MASSDRIVER_TOKEN"), "Secure token to authenticate with Massdriver. Defaults to value in MASSDRIVER_TOKEN environment variable.")
	deploymentStartCmd.Flags().StringP("out", "o", ".", "Destination path to write deployment json files. Defaults to current directory")

	deploymentCmd.AddCommand(deploymentCompleteCmd)
	deploymentCompleteCmd.Flags().StringP("deployment-id", "i", os.Getenv("MASSDRIVER_DEPLOYMENT_ID"), "Massdriver Deployment ID. Defaults to value in MASSDRIVER_DEPLOYMENT_ID environment variable.")
	deploymentCompleteCmd.Flags().StringP("token", "t", os.Getenv("MASSDRIVER_TOKEN"), "Secure token to authenticate with Massdriver. Defaults to value in MASSDRIVER_TOKEN environment variable.")
	deploymentCompleteCmd.Flags().StringP("artifacts", "f", "./artifact.json", "Path to JSON formatted artifact file to upload. Defaults to ./artifact.json")
}

func RunDeploymentStart(cmd *cobra.Command, args []string) error {
	bundle, _ := cmd.Flags().GetString("bundle")
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

	log.Info().Str("deployment", id).Msg("generating auth files")
	schemaPath := path.Join(bundle, "/schema-connections.json")
	connectionsPath := path.Join(out, massdriver.ConnectionsFileName)
	authPath := path.Join(out, "auth")
	if _, err := os.Stat(authPath); os.IsNotExist(err) {
		err := os.Mkdir(authPath, 0777)
		if err != nil {
			log.Error().Err(err).Str("deployment", id).Msg("an error occurred while creating auth directory")
			return err
		}
	}
	err = provisioners.GenerateAuthFiles(schemaPath, connectionsPath, authPath)
	if err != nil {
		log.Error().Err(err).Str("deployment", id).Msg("an error occurred while generating auth files")
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
