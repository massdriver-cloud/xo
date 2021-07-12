package bundles

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"regexp"

	"github.com/rs/zerolog/log"
)

// relativeFilePathPattern only accepts relative file path prefixes "./" and "../"
var relativeFilePathPattern = regexp.MustCompile(`^(\.\/|\.\.\/)`)

type hydratedOrderedJSON OrderedJSON

func Hydrate(any interface{}, cwd string) (interface{}, error) {
	val := getValue(any)

	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		elem := reflect.TypeOf(any).Elem()
		hydratedList := reflect.Zero(reflect.SliceOf(elem))
		for i := 0; i < val.Len(); i++ {
			hydratedVal, err := Hydrate(val.Index(i).Interface(), cwd)
			if err != nil {
				return hydratedList, err
			}
			// if we hydrated a the, we need some special logic. What comes back is an array,
			// which we need to extract the fields and promote them to the current level. We
			// also need to check for keys
			if reflect.TypeOf(hydratedVal) == reflect.TypeOf(hydratedOrderedJSON{}) {
				keys := getKeys(val)
				hoj := hydratedVal.(hydratedOrderedJSON)
				for _, v := range hoj {
					key := v.Key
					collision := false
					for _, k := range keys {
						if key == k {
							collision = true
						}
					}
					if !collision {
						hydratedList = reflect.Append(hydratedList, getValue(v))
					}
				}
			} else {
				hydratedList = reflect.Append(hydratedList, getValue(hydratedVal))
			}
		}
		if elem == reflect.TypeOf(OrderedJSONElement{}) {
			out := OrderedJSON{}
			for i := 0; i < hydratedList.Len(); i++ {
				out = append(out, hydratedList.Index(i).Interface().(OrderedJSONElement))
			}
			return out, nil
		}
		return hydratedList.Interface(), nil
	// As of right now, the only structs we should receive are the OrderedJSONElement structs
	// so we can make some assumptions about how to cast the object and extract data
	case reflect.Struct:
		schemaInterface := val.Interface()
		oje := schemaInterface.(OrderedJSONElement)
		hydratedSchema := hydratedOrderedJSON{}

		// we know that all the Keys should be strings, cuz JSON...
		key := oje.Key.(string)
		if key == "$ref" {
			schemaRefPath := oje.Value.(string)
			if relativeFilePathPattern.MatchString(schemaRefPath) {
				// Build up the path from where the dir current schema was read
				schemaRefAbsPath, err := filepath.Abs(filepath.Join(cwd, schemaRefPath))
				if err != nil {
					return hydratedSchema, err
				}

				schemaRefDir := filepath.Dir(schemaRefAbsPath)
				referencedSchema, err := readJsonFile(schemaRefAbsPath)
				if err != nil {
					return hydratedSchema, err
				}

				for _, v := range referencedSchema {
					hydratedValue, err := Hydrate(v, schemaRefDir)
					if err != nil {
						return hydratedSchema, err
					}
					hydratedSchema = append(hydratedSchema, hydratedValue.(OrderedJSONElement))
				}
				return hydratedSchema, nil
			} else {
				// We can remove this log statement, but I think its useful to alert the user instead of skipping
				log.Warn().Msg("Unable to dereference schema $ref path: " + schemaRefPath)
			}
		}

		hydratedValue, err := Hydrate(oje.Value, cwd)
		if err != nil {
			return hydratedSchema, err
		}
		oje.Value = hydratedValue

		return oje, nil

	// case reflect.Map:
	// 	schemaInterface := val.Interface()
	// 	schema := schemaInterface.(map[string]interface{})
	// 	hydratedSchema := map[string]interface{}{}

	// 	// if this part of the schema has a $ref that is a local file, read it and make it
	// 	// the map that we hydrate into. This causes any keys in the ref'ing object to override anything in the ref'd object
	// 	// which adheres to the JSON Schema spec.
	// 	if schemaRefInterface, ok := schema["$ref"]; ok {
	// 		schemaRefPath := schemaRefInterface.(string)
	// 		if relativeFilePathPattern.MatchString(schemaRefPath) {
	// 			// Build up the path from where the dir current schema was read
	// 			schemaRefAbsPath, err := filepath.Abs(filepath.Join(cwd, schemaRefPath))
	// 			if err != nil {
	// 				return hydratedSchema, err
	// 			}

	// 			schemaRefDir := filepath.Dir(schemaRefAbsPath)
	// 			referencedSchema, err := readJsonFile(schemaRefAbsPath)
	// 			if err != nil {
	// 				return hydratedSchema, err
	// 			}

	// 			// Remove it if, so it doesn't get rehydrated below
	// 			delete(schema, "$ref")

	// 			for k, v := range referencedSchema {
	// 				hydratedValue, err := Hydrate(v.(interface{}), schemaRefDir)
	// 				if err != nil {
	// 					return hydratedSchema, err
	// 				}
	// 				hydratedSchema[k] = hydratedValue
	// 			}
	// 		}
	// 	}

	// 	for key, value := range schema {
	// 		var valueInterface = value.(interface{})
	// 		hydratedValue, err := Hydrate(valueInterface, cwd)
	// 		if err != nil {
	// 			return hydratedSchema, err
	// 		}
	// 		hydratedSchema[key] = hydratedValue
	// 	}

	// 	return hydratedSchema, nil
	default:
		return any, nil
	}
}

func getValue(any interface{}) reflect.Value {
	val := reflect.ValueOf(any)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return val
}

// func readJsonFile(filepath string) (map[string]interface{}, error) {
// 	var result map[string]interface{}
// 	data, err := ioutil.ReadFile(filepath)
// 	if err != nil {
// 		return result, err
// 	}
// 	err = json.Unmarshal([]byte(data), &result)

// 	return result, err
// }

func readJsonFile(filepath string) (OrderedJSON, error) {
	var result OrderedJSON
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal([]byte(data), &result)

	return result, err
}

func maybePanic(err error) {
	if err != nil {
		panic(err)
	}
}

func getKeys(val reflect.Value) []string {
	var keys []string
	for i := 0; i < val.Len(); i++ {
		oje := val.Index(i).Interface().(OrderedJSONElement)
		keys = append(keys, oje.Key.(string))
	}
	return keys
}
