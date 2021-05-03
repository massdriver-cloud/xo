package cmd

import (
	"github.com/spf13/cobra"
)

var artifactCmd = &cobra.Command{
	Use:   "artifact",
	Short: "Manage Massdriver deployments",
	Long:  ``,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("artifact called")
	//},
}

func init() {
	rootCmd.AddCommand(artifactCmd)
}
