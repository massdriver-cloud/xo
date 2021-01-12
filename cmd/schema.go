package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// schemaCmd represents the schema command
var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("schema called")
	},
}

func init() {
	rootCmd.AddCommand(schemaCmd)
	validateCmd.PersistentFlags().StringP("schema", "s", "schema.json", "Path to JSON Schema")
}
