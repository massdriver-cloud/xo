package cmd

import (
	"fmt"
	"os"
	"path"

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
	schemaLoader := getLoader(schemaPath)
	documentLoader := getLoader(documentPath)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
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

func getLoader(path string) gojsonschema.JSONLoader {
	prefix := "file://"
	ref := prefix + expandPath(path)
	return gojsonschema.NewReferenceLoader(ref)
}

func expandPath(p string) string {
	if path.IsAbs(p) {
		return p
	}

	cwd, err := os.Getwd()
	maybeHandleError(err)
	return path.Join(cwd, p)
}

func maybeHandleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
