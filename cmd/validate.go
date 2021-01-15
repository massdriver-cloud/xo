package cmd

import (
	"errors"
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
	RunE:  runValidate,
}

func runValidate(cmd *cobra.Command, args []string) error {
	schema, _ := cmd.Flags().GetString("schema")
	document, _ := cmd.Flags().GetString("document")

	result, err := Validate(schema, document)
	fmt.Printf("Document valid? %t", result)
	return err
}

func init() {
	schemaCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringP("document", "d", "document.json", "Path to JSON document")
}

// Validate the input object against the schema
func Validate(schemaPath string, documentPath string) (bool, error) {
	sl := schemaloader.Load(schemaPath)
	dl := schemaloader.Load(documentPath)

	result, err := gojsonschema.Validate(sl, dl)
	if err != nil {
		return false, err
	}

	if result.Valid() {
		return true, nil
	} else {
		msg := "The document is not valid. see errors :\n"
		for _, desc := range result.Errors() {
			msg = msg + fmt.Sprintf("- %s\n", desc)
		}

		err = errors.New(msg)
		return false, err
	}
}
