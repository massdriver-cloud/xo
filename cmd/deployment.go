package cmd

import (
	"os"
	"xo/src/massdriver"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var deploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "Manage Massdriver deployments",
	Long:  ``,
}

var deploymentGetCmd = &cobra.Command{
	Use:                   "get -i [id]",
	Short:                 "Fetch deployment from Massdriver",
	Long:                  ``,
	RunE:                  RunDeploymentGet,
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.AddCommand(deploymentCmd)
	deploymentCmd.AddCommand(deploymentGetCmd)
	deploymentGetCmd.Flags().StringP("deployment-id", "i", os.Getenv("MASSDRIVER_DEPLOYMENT_ID"), "Massdriver Deployment ID")
	deploymentGetCmd.Flags().StringP("token", "t", os.Getenv("MASSDRIVER_TOKEN"), "Secure token to authenticate with Massdriver")
	deploymentGetCmd.Flags().StringP("dest", "d", ".", "Destination path to write deployment json files")
}

func RunDeploymentGet(cmd *cobra.Command, args []string) error {
	id, _ := cmd.Flags().GetString("deployment-id")
	token, _ := cmd.Flags().GetString("token")
	dest, _ := cmd.Flags().GetString("dest")

	logger.Info("getting deployment from massdriver", zap.String("deployment", id))
	dep, err := massdriver.GetDeployment(id, token)
	if err != nil {
		return err
	}

	logger.Info("writing deployment to file", zap.String("deployment", id))
	err = massdriver.WriteDeploymentToFile(dep, dest)
	return err
}
