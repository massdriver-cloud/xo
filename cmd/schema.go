package cmd

import (
	"errors"
	"fmt"
	"xo/src/schemaloader"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
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
	schema, _ := cmd.Flags().GetString("schema")
	document, _ := cmd.Flags().GetString("document")
	_, err := Validate(schema, document)
	return err
}

// Validate the input object against the schema
func Validate(schemaPath string, documentPath string) (bool, error) {
	log.Debug().
		Str("schemaPath", schemaPath).
		Str("documentPath", documentPath).Msg("Validating schema.")

	sl := schemaloader.Load(schemaPath)
	dl := schemaloader.Load(documentPath)

	result, err := gojsonschema.Validate(sl, dl)
	if err != nil {
		log.Error().Err(err).Msg("Validator failed.")
		return false, err
	}

	if !result.Valid() {
		msg := "The document is not valid. see errors :\n"
		for _, desc := range result.Errors() {
			msg = msg + fmt.Sprintf("- %s\n", desc)
		}

		err = errors.New(msg)
		log.Error().Err(err).Msg("Validation failed.")
		return false, err
	}

	return true, nil
}
