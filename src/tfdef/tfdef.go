package tfdef

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/xeipuuv/gojsonschema"
)

// Property is a single JSON Schema property field
type Property struct {
	Type string `json:"type"`
}

// Convert a JSON Schema to Terraform Variable Definition JSON
func Convert(mainPath string) string {
	properties := getProperties(mainPath)

	fmt.Printf("Properties: %+v\n", properties)

	variableDef := map[string]interface{}{
		"variable": properties,
	}

	// TODO handle this :D
	result, _ := json.Marshal(variableDef)
	return string(result)
}

func getProperties(path string) interface{} {
	schemaLoader := gojsonschema.NewReferenceLoader(path)

	sl := gojsonschema.NewSchemaLoader()
	what, serr := sl.Compile(schemaLoader)
	fmt.Printf("what: %+v\n", what)
	fmt.Printf("serr: %+v\n", serr)

	// TODO handle this :D
	schemaSrc, schemaErr := schemaLoader.LoadJSON()

	fmt.Printf("Schema: %+v\n", schemaSrc)
	fmt.Printf("schemaErr: %+v\n", schemaErr)

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
