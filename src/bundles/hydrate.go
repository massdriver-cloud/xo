package bundles

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"regexp"
)

// relativeFilePathPattern only accepts relative file path prefixes "./" and "../"
var relativeFilePathPattern = regexp.MustCompile(`^(\.\/|\.\.\/)`)

func Hydrate(any interface{}, cwd string) interface{} {
	_ = cwd
	val := getValue(any)

	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		hydratedList := make([]interface{}, 0)
		for i := 0; i < val.Len(); i++ {
			hydratedVal := Hydrate(val.Index(i).Interface(), cwd)
			hydratedList = append(hydratedList, hydratedVal)
		}
		return hydratedList
	case reflect.Map:
		schemaInterface := val.Interface()
		schema := schemaInterface.(map[string]interface{})
		hydratedSchema := map[string]interface{}{}

		// if this part of the schema has a $ref that is a local file, read it and make it
		// the map that we hydrate into. This causes any keys in the ref'ing object to override anything in the ref'd object
		// which adheres to the JSON Schema spec.
		if schemaRefInterface, ok := schema["$ref"]; ok {
			schemaRefPath := schemaRefInterface.(string)
			if relativeFilePathPattern.MatchString(schemaRefPath) {
				// Build up the path from where the dir current schema was read
				schemaRefAbsPath, err := filepath.Abs(filepath.Join(cwd, schemaRefPath))
				maybePanic(err)

				schemaRefDir := filepath.Dir(schemaRefAbsPath)
				referencedSchema, err := readJsonFile(schemaRefAbsPath)
				maybePanic(err)
				// Remove it if, so it doesn't get rehydrated below
				delete(schema, "$ref")

				// TODO: Why won't Hydrate accept referencedSchema.(interface{})
				for k, v := range referencedSchema {
					hydratedSchema[k] = Hydrate(v.(interface{}), schemaRefDir)
				}
			}
		}

		for key, value := range schema {
			var valueInterface = value.(interface{})
			hydratedSchema[key] = Hydrate(valueInterface, cwd)
		}
		return hydratedSchema
	default:
		return any
	}
}

func getValue(any interface{}) reflect.Value {
	val := reflect.ValueOf(any)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return val
}

func readJsonFile(filepath string) (map[string]interface{}, error) {
	var result map[string]interface{}
	data, err := ioutil.ReadFile(filepath)
	maybePanic(err)
	err = json.Unmarshal([]byte(data), &result)
	maybePanic(err)

	return result, err
}
