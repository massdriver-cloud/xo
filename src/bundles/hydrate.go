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
		hydratedList := []interface{}{}
		for i := 0; i < val.Len(); i++ {
			hydratedVal, err := Hydrate(val.Index(i).Interface(), cwd)
			if err != nil {
				return hydratedList, err
			}
			// if we hydrated something then we need some special logic. What comes back is a custom array type
			// (hydratedOrderedJSON), but we don't want to embed the array. We want to extract the elements and
			// place them alongside the current level. We also need to check for duplicate keys and ignore them
			if reflect.TypeOf(hydratedVal) == reflect.TypeOf(hydratedOrderedJSON{}) {
				keys := getKeys(&val)
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
						hydratedList = append(hydratedList, getValue(v).Interface())
					}
				}
			} else {
				hydratedList = append(hydratedList, getValue(hydratedVal).Interface())
			}
		}
		if elem == reflect.TypeOf(OrderedJSONElement{}) {
			out := OrderedJSON{}
			for i := 0; i < len(hydratedList); i++ {
				out = append(out, hydratedList[i].(OrderedJSONElement))
			}
			return out, nil
		}
		return hydratedList, nil
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

func readJsonFile(filepath string) (OrderedJSON, error) {
	var result OrderedJSON
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal([]byte(data), &result)

	return result, err
}

// utility function to extract the "keys" from an OrderedJSONElement Array
func getKeys(val *reflect.Value) []string {
	var keys []string
	for i := 0; i < (*val).Len(); i++ {
		oje := (*val).Index(i).Interface().(OrderedJSONElement)
		keys = append(keys, oje.Key.(string))
	}
	return keys
}
