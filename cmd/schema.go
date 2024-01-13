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

var schemaDereferenceCmd = &cobra.Command{
	Use:   "dereference",
	Short: "Dereferences a schema, resolving all local $ref's",
	Long:  ``,
	RunE:  runSchemaDereference,
}

func init() {
	rootCmd.AddCommand(schemaCmd)
	schemaCmd.AddCommand(schemaValidateCmd)
	schemaCmd.AddCommand(schemaDereferenceCmd)

	schemaValidateCmd.Flags().StringP("document", "d", "document.json", "Path to JSON document")
	schemaValidateCmd.Flags().StringP("schema", "s", "./schema.json", "Path to JSON Schema")

	schemaDereferenceCmd.Flags().StringP("schema", "s", "./schema.json", "Path to JSON Schema")
	schemaDereferenceCmd.Flags().StringP("dir", "d", ".", "Path to output directory")
}

func runSchemaValidate(cmd *cobra.Command, args []string) error {
	schemaPath, _ := cmd.Flags().GetString("schema")
	documentPath, _ := cmd.Flags().GetString("document")

	schema := jsonschema.Load(schemaPath)
	document := jsonschema.Load(documentPath)

	_, err := jsonschema.Validate(schema, document)
	return err
}

func runSchemaDereference(cmd *cobra.Command, args []string) error {
	schema, _ := cmd.Flags().GetString("schema")
	dir, _ := cmd.Flags().GetString("dir")
	err := jsonschema.WriteDereferencedSchema(schema, dir)
	return err
}
