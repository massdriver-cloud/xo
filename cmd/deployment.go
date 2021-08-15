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

var deploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "Manage Massdriver deployments",
	Long:  ``,
}

var deploymentStartCmd = &cobra.Command{
	Use:                   "start",
	Short:                 "Start deployment from Massdriver",
	Long:                  deploymentStartLong,
	Example:               deploymentStartExamples,
	RunE:                  RunDeploymentStart,
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.AddCommand(deploymentCmd)
	deploymentCmd.AddCommand(deploymentStartCmd)
	deploymentStartCmd.Flags().StringP("deployment-id", "i", os.Getenv("MASSDRIVER_DEPLOYMENT_ID"), "Massdriver Deployment ID. Defaults to value in MASSDRIVER_DEPLOYMENT_ID environment variable.")
	deploymentStartCmd.Flags().StringP("token", "t", os.Getenv("MASSDRIVER_TOKEN"), "Secure token to authenticate with Massdriver. Defaults to value in MASSDRIVER_TOKEN environment variable.")
	deploymentStartCmd.Flags().StringP("out", "o", ".", "Destination path to write deployment json files. Defaults to current directory")
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
	dep, err := massdriver.StartDeployment(id, token)
	if err != nil {
		log.Error().Err(err).Str("deployment", id).Msg("an error occurred while getting deployment from Massdriver")
		return err
	}

	log.Info().Str("deployment", id).Msg("writing deployment to file")
	err = massdriver.WriteDeploymentToFile(dep, out)
	if err != nil {
		log.Error().Err(err).Str("deployment", id).Msg("an error occurred while writing deployment files")
		return err
	}
	log.Info().Str("deployment", id).Msg("deployment get complete")

	return nil
}
