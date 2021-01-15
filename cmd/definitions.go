package cmd

import (
	"fmt"
	"io/ioutil"

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
	definitionsCmd.Flags().StringP("output", "o", "variables.tf.json", "Output path. Use - for STDOUT")
}

func runDefinitions(cmd *cobra.Command, args []string) error {
	provisioner := args[0]
	outputPath, _ := cmd.Flags().GetString("output")
	schema, _ := cmd.Flags().GetString("schema")
	var compiled string
	var err error

	switch provisioner {
	case "terraform":
		compiled, err = tfdef.Compile(schema)
		if err != nil {
			return err
		}
	default:
		err := fmt.Errorf("Unsupported argument %s the single argument 'terraform' is supported", provisioner)
		return err
	}

	return writeVariableFile(compiled, outputPath)
}

func writeVariableFile(compiled string, outPath string) error {
	switch outPath {
	case "-":
		fmt.Println(compiled)
		return nil
	default:
		data := []byte(compiled)
		err := ioutil.WriteFile(outPath, data, 0644)
		return err
	}
}
