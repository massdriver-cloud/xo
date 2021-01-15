package cmd

import (
	"fmt"

	"xo/tfdef"

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
	RunE: runDefinitions,
}

func init() {
	provisionerCmd.AddCommand(definitionsCmd)
	definitionsCmd.Flags().StringP("schema", "s", "schema.json", "Path to JSON Schema")
}

func runDefinitions(cmd *cobra.Command, args []string) error {
	provisioner := args[0]
	switch provisioner {
	case "terraform":
		schema, _ := cmd.Flags().GetString("schema")
		compiled, _ := tfdef.Compile(schema)
		fmt.Println(compiled)
		return nil
	default:
		err := fmt.Errorf("Unsupported argument %s the single argument 'terraform' is supported", provisioner)
		return err
	}
}
