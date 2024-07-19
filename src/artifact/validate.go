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

func Validate(field string, artifact []byte, schemasIn io.Reader) (bool, error) {

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
	dl := gojsonschema.NewBytesLoader(artifact)

	return jsonschema.Validate(sl, dl)
}
