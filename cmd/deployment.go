package cmd

import (
	"fmt"
	"os"
	"xo/src/massdriver"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var deploymentGetLong = `
	Fetches metadata about a Massdriver deployment and writes the data to files.

	Specifically this is fetching the params and connections and writing them to
	params.tfvars.json and connections.tfvars.json, respectively. This is intended
	to be used as a step in workflow execution to gather resources for the provisioner.
	`
var deploymentGetExamples = `
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

var deploymentGetCmd = &cobra.Command{
	Use:                   "get",
	Short:                 "Fetch deployment from Massdriver",
	Long:                  deploymentGetLong,
	Example:               deploymentGetExamples,
	Run:                   RunDeploymentGet,
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.AddCommand(deploymentCmd)
	deploymentCmd.AddCommand(deploymentGetCmd)
	deploymentGetCmd.Flags().StringP("deployment-id", "i", os.Getenv("MASSDRIVER_DEPLOYMENT_ID"), "Massdriver Deployment ID. Defaults to value in MASSDRIVER_DEPLOYMENT_ID environment variable.")
	deploymentGetCmd.Flags().StringP("token", "t", os.Getenv("MASSDRIVER_TOKEN"), "Secure token to authenticate with Massdriver. Defaults to value in MASSDRIVER_TOKEN environment variable.")
	deploymentGetCmd.Flags().StringP("dest", "d", ".", "Destination path to write deployment json files. Defaults to current directory")
}

func RunDeploymentGet(cmd *cobra.Command, args []string) {
	id, _ := cmd.Flags().GetString("deployment-id")
	token, _ := cmd.Flags().GetString("token")
	dest, _ := cmd.Flags().GetString("dest")

	if id == "" || token == "" {
		cmd.Help()
		fmt.Println("\nERROR: Both deployment-id and token must be set (by flags or environment variable)")
		os.Exit(1)
	}

	logger.Info("getting deployment from Massdriver", zap.String("deployment", id))
	dep, err := massdriver.GetDeployment(id, token)
	if err != nil {
		logger.Error("an error occurred while getting deployment from Massdriver", zap.String("deployment", id), zap.Error(err))
		os.Exit(1)
	}

	logger.Info("writing deployment to file", zap.String("deployment", id))
	err = massdriver.WriteDeploymentToFile(dep, dest)
	if err != nil {
		logger.Error("an error occurred while writing deployment files", zap.String("deployment", id), zap.Error(err))
		os.Exit(1)
	}
	logger.Info("deployment get complete")
}
