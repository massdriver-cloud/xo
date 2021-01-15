package cmd

import (
	"github.com/spf13/cobra"
)

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Manage JSON Schemas",
	Long:  ``,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("schema called")
	// },
}

func init() {
	rootCmd.AddCommand(schemaCmd)
	validateCmd.PersistentFlags().StringP("schema", "s", "./schema.json", "Path to JSON Schema")
}
