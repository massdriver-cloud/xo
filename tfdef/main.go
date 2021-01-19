package tfdef

import (
	"encoding/json"
	"xo/schemaloader"
)

// Compile a JSON Schema to Terraform Variable Definition JSON
func Compile(path string) (string, error) {
	vars, err := getVars(path)
	if err != nil {
		return "", err
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

func getJSONSchema(path string) (Schema, error) {
	schema := Schema{}
	sl := schemaloader.Load(path)

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
