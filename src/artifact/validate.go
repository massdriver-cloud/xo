package artifact

import (
	"encoding/json"
	"errors"
	"io"
	"xo/src/jsonschema"

	"github.com/xeipuuv/gojsonschema"
)

type artifactSchema struct {
	Properties map[string]interface{} `json:"properties"`
}

func Validate(field string, artifactIn, schemasIn io.Reader) (bool, error) {

	artifactBytes, err := io.ReadAll(artifactIn)
	if err != nil {
		return false, err
	}
	schemaBytes, err := io.ReadAll(schemasIn)
	if err != nil {
		return false, err
	}

	var schemaObj artifactSchema
	err = json.Unmarshal(schemaBytes, &schemaObj)
	if err != nil {
		return false, err
	}
	specificSchema, exists := schemaObj.Properties[field]
	if !exists {
		return false, errors.New(`artifact validation failed: field "` + field + `" does not exist in schema`)
	}

	sl := gojsonschema.NewGoLoader(specificSchema.(map[string]interface{}))
	dl := gojsonschema.NewBytesLoader(artifactBytes)

	return jsonschema.Validate(sl, dl)
}
