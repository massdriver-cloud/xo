package cmd

import (
	"errors"
	"fmt"
	"xo/schemaloader"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validates an input JSON object against a JSON schema",
	Long:  ``,
	RunE:  runValidate,
}

func runValidate(cmd *cobra.Command, args []string) error {
	schema, _ := cmd.Flags().GetString("schema")
	document, _ := cmd.Flags().GetString("document")
	_, err := Validate(schema, document)
	return err
}

func init() {
	schemaCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringP("document", "d", "document.json", "Path to JSON document")
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
