package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// definitionsCmd represents the definitions command
var definitionsCmd = &cobra.Command{
	Use:   "definitions",
	Short: "Generate provisioner variable definition files from a schema",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: run,
}

func init() {
	provisionerCmd.AddCommand(definitionsCmd)
	validateCmd.Flags().StringP("schema", "s", "schema.json", "Path to JSON Schema")
}

func run(cmd *cobra.Command, args []string) {
	// TODO: only supports 1 arg currently, "terraform"
	fmt.Printf("%v", args)
	fmt.Println("definitions called")
}
