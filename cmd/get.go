package cmd

import (
	"xo/massdriver"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:                   "get -i [id] -t [token]",
	Short:                 "Fetch object from Massdriver",
	Long:                  ``,
	RunE:                  RunGet,
	DisableFlagsInUseLine: true,
}

func init() {
	deploymentCmd.AddCommand(getCmd)
	getCmd.Flags().StringP("id", "i", "", "ID of resource to fetch")
	getCmd.Flags().StringP("token", "t", "", "Secure token to authenticate with Massdriver")
	getCmd.MarkFlagRequired("id")
	getCmd.MarkFlagRequired("token")
}

func RunGet(cmd *cobra.Command, args []string) error {
	var err error
	if cmd.Parent().Name() == "deployment" {
		id, _ := cmd.Flags().GetString("id")
		token, _ := cmd.Flags().GetString("token")
		dep, depErr := massdriver.GetDeployment(id, token)
		_ = dep
		return depErr
	}
	return err
}
