package cmd

import (
	"fmt"
	"xo/schemaloader"

	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validates an input object against a schema",
	Long:  `Validates a JSON object against a JSON Schema.`,
	Run: func(cmd *cobra.Command, args []string) {
		schema, _ := cmd.Flags().GetString("schema")
		document, _ := cmd.Flags().GetString("document")

		// TODO: Handle error
		Validate(schema, document)
	},
}

func init() {
	schemaCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringP("document", "d", "document.json", "Path to JSON document")
}

// Validate the input object against the schema
func Validate(schemaPath string, documentPath string) bool {
	sl := schemaloader.Load(schemaPath)
	dl := schemaloader.Load(documentPath)

	result, err := gojsonschema.Validate(sl, dl)
	maybeHandleError(err)

	if result.Valid() {
		fmt.Printf("The document is valid\n")
		return true
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
		return false
	}
}

func maybeHandleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
