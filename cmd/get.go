package cmd

import (
	"xo/massdriver"

	"github.com/spf13/cobra"
)

var deploymentGetCmd = &cobra.Command{
	Use:                   "get -i [id] -t [token]",
	Short:                 "Fetch object from Massdriver",
	Long:                  ``,
	RunE:                  RunDeploymentGet,
	DisableFlagsInUseLine: true,
}

func init() {
	deploymentCmd.AddCommand(deploymentGetCmd)
	deploymentGetCmd.Flags().StringP("id", "i", "", "ID of resource to fetch")
	deploymentGetCmd.Flags().StringP("token", "t", "", "Secure token to authenticate with Massdriver")
	deploymentGetCmd.Flags().StringP("dest", "d", ".", "Destination path to write deployment json files")
	deploymentGetCmd.MarkFlagRequired("id")
	deploymentGetCmd.MarkFlagRequired("token")
}

func RunDeploymentGet(cmd *cobra.Command, args []string) error {
	id, _ := cmd.Flags().GetString("id")
	token, _ := cmd.Flags().GetString("token")
	dest, _ := cmd.Flags().GetString("dest")

	dep, err := massdriver.GetDeployment(id, token)
	if err != nil {
		return err
	}

	err = massdriver.WriteDeploymentToFile(dep, dest)
	return err
}
