package tfdef

import (
	"encoding/json"

	"github.com/xeipuuv/gojsonschema"
)

// Compile a JSON Schema to Terraform Variable Definition JSON
func Compile(path string) string {
	variableFile := TFVariableFile{Variable: getVars(path)}

	// TODO handle this :D
	result, _ := json.Marshal(variableFile)
	return string(result)
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
	schemaLoader := gojsonschema.NewReferenceLoader(path)

	// TODO handle this :D
	schemaSrc, schemaErr := schemaLoader.LoadJSON()

	// fmt.Printf("Schema: %+v\n", schemaSrc)
	// fmt.Printf("schemaErr: %+v\n", schemaErr)
	_ = schemaErr
	byteData, _ := json.Marshal(schemaSrc)

	schema := Schema{}
	json.Unmarshal(byteData, &schema)
	return schema
}
