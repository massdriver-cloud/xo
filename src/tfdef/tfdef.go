package tfdef

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/xeipuuv/gojsonschema"
)

// Compile a JSON Schema to Terraform Variable Definition JSON
func Compile(path string) string {
	properties := getProperties(path)

	variableDef := map[string]interface{}{
		"variable": properties,
	}

	// TODO handle this :D
	result, _ := json.Marshal(variableDef)

	return string(result)
}

func getProperties(path string) interface{} {
	schemaLoader := gojsonschema.NewReferenceLoader(path)

	// TODO handle this :D
	schemaSrc, schemaErr := schemaLoader.LoadJSON()
	src, _ := schemaLoader.JsonReference()
	fmt.Printf("Src: %+v\n", src.String())

	// fmt.Printf("Schema: %+v\n", schemaSrc)
	// fmt.Printf("schemaErr: %+v\n", schemaErr)
	_ = schemaErr

	// TODO handle this :D
	schema, _ := schemaSrc.(map[string]interface{})
	propertiesRaw := schema["properties"].(map[string]interface{})

	properties := map[string]Property{}

	for k, v := range propertiesRaw {
		var property Property
		err := mapstructure.Decode(v, &property)
		if err != nil {
			// TODO handle this :D
			panic(err)
		}
		properties[k] = property
	}

	return properties
}
