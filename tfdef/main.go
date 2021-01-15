package tfdef

import (
	"encoding/json"
	"xo/schemaloader"
)

// Compile a JSON Schema to Terraform Variable Definition JSON
func Compile(path string) (string, error) {

	variableFile := TFVariableFile{Variable: getVars(path)}

	// TODO handle this :D
	result, _ := json.Marshal(variableFile)
	return string(result), nil
}

func getVars(path string) map[string]TFVariable {
	schema := getJSONSchema(path)
	variables := map[string]TFVariable{}
	for name, prop := range schema.Properties {
		variables[name] = NewTFVariable(prop)
	}
	return variables
}

func getJSONSchema(path string) Schema {
	sl := schemaloader.Load(path)

	// TODO handle this :D
	schemaSrc, _ := sl.LoadJSON()
	byteData, _ := json.Marshal(schemaSrc)

	schema := Schema{}
	json.Unmarshal(byteData, &schema)
	return schema
}
