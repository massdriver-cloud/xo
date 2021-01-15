package cmd

import (
	"fmt"
	"io/ioutil"

	"xo/tfdef"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// compileCmd represents the compile command
var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compile JSON Schema to provisioner variable definition file.",
	Long:  ``,
	RunE:  runCompile,
}

func init() {
	provisionerCmd.AddCommand(compileCmd)
	compileCmd.Flags().StringP("schema", "s", "./schema.json", "Path to JSON Schema")
	compileCmd.Flags().StringP("output", "o", "./variables.tf.json", "Output path. Use - for STDOUT")
}

func runCompile(cmd *cobra.Command, args []string) error {
	var compiled string
	var err error

	provisioner := args[0]
	schema, _ := cmd.Flags().GetString("schema")
	outputPath, _ := cmd.Flags().GetString("output")

	log.Debug().
		Str("provisioner", provisioner).
		Str("schemaPath", schema).Msg("Compiling schema.")

	switch provisioner {
	case "terraform":
		compiled, err = tfdef.Compile(schema)
		if err != nil {
			return err
		}
	default:
		err := fmt.Errorf("Unsupported argument %s the single argument 'terraform' is supported", provisioner)
		log.Error().Err(err).Msg("Compilation failed.")
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
