package jsonschema

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/xeipuuv/gojsonschema"
)

// Validate the input object against the schema
func Validate(schema, document gojsonschema.JSONLoader) (bool, error) {
	log.Debug().Msg("Validating schema.")

	result, err := gojsonschema.Validate(schema, document)
	if err != nil {
		log.Error().Err(err).Msg("Validator failed.")
		return false, err
	}

	if !result.Valid() {
		msg := "The document failed validation:\nErrors:\n"
		for _, desc := range result.Errors() {
			msg = msg + fmt.Sprintf("\t- %s\n", desc)
		}

		err = errors.New(msg)
		log.Error().Err(err).Msg("Validation failed.")
		return false, err
	}

	return true, nil
}
