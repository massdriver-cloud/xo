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

	// TODO handle this :D
	result, err := json.MarshalIndent(variableFile, "", "  ")
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func getVars(path string) (map[string]TFVariable, error) {
	variables := map[string]TFVariable{}
	schema, err := getJSONSchema(path)
	if err != nil {
		return variables, err
	}

	for name, prop := range schema.Properties {
		variables[name] = NewTFVariable(prop)
	}
	return variables, nil
}

func getJSONSchema(path string) (jsonschema.Schema, error) {
	schema := jsonschema.Schema{}
	sl := jsonschema.Load(path)

	schemaSrc, err := sl.LoadJSON()
	if err != nil {
		return schema, err
	}

	byteData, err := json.Marshal(schemaSrc)
	if err != nil {
		return schema, err
	}

	json.Unmarshal(byteData, &schema)
	return schema, nil
}
