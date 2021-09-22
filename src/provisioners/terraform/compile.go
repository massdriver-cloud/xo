package terraform

import (
	"encoding/json"
	"xo/src/jsonschema"
)

// Compile a JSON Schema to Terraform Variable Definition JSON
func Compile(path string) (string, error) {
	vars, err := getVars(path)
	if err != nil {
		return "", err
	}

	// You can't have an empty variable block, so if there are no vars return an empty json block
	if len(vars) == 0 {
		return "{}", nil
	}

	variableFile := TFVariableFile{Variable: vars}

	result, err := json.MarshalIndent(variableFile, "", "  ")
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func getVars(path string) (map[string]TFVariable, error) {
	variables := map[string]TFVariable{}
	schema, err := jsonschema.GetJSONSchema(path)
	if err != nil {
		return variables, err
	}

	for name, prop := range schema.Properties {
		variables[name] = NewTFVariable(prop)
	}
	return variables, nil
}
