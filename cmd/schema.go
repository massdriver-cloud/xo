package cmd

import (
	"xo/src/jsonschema"

	"github.com/spf13/cobra"
)

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Manage JSON Schemas",
	Long:  ``,
}

var schemaValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validates an input JSON object against a JSON schema",
	Long:  ``,
	RunE:  runSchemaValidate,
}

func init() {
	rootCmd.AddCommand(schemaCmd)
	schemaCmd.AddCommand(schemaValidateCmd)

	schemaValidateCmd.Flags().StringP("document", "d", "document.json", "Path to JSON document")
	schemaValidateCmd.Flags().StringP("schema", "s", "./schema.json", "Path to JSON Schema")
}

func runSchemaValidate(cmd *cobra.Command, args []string) error {
	schemaPath, _ := cmd.Flags().GetString("schema")
	documentPath, _ := cmd.Flags().GetString("document")

	schema := jsonschema.Load(schemaPath)
	document := jsonschema.Load(documentPath)

	_, err := jsonschema.Validate(schema, document)
	return err
}
