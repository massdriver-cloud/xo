package cmd

import (
	"github.com/spf13/cobra"
)

var deploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "Manage Massdriver deployments",
	Long:  ``,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("schema called")
	// },
}

func init() {
	rootCmd.AddCommand(deploymentCmd)
}
