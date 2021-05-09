package bundles

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"regexp"
)

// relativeFilePathPattern only accepts relative file path prefixes "./" and "../"
var relativeFilePathPattern = regexp.MustCompile(`^(\.\/|\.\.\/)`)

func Hydrate(any interface{}) interface{} {
	val := getValue(any)

	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		newList := make([]interface{}, 0)
		for i := 0; i < val.Len(); i++ {
			hydratedVal := Hydrate(val.Index(i).Interface())
			newList = append(newList, hydratedVal)
		}
		return newList
	case reflect.Map:
		newMap := map[string]interface{}{}
		for _, keyInterface := range val.MapKeys() {
			var key = keyInterface.String()
			var valueInterface = val.MapIndex(keyInterface).Interface()

			if key == "$ref" {
				var refPath = getValue(valueInterface).String()
				if relativeFilePathPattern.MatchString(refPath) {
					jsonObject, err := readJsonFile(refPath)
					maybePanic(err)

					// non-deterministic...
					// i think we can fix this by merging new into ref'd and setting result as new
					for k, v := range jsonObject {
						newMap[k] = Hydrate(v.(interface{}))
					}
				} else {
					newMap[key] = Hydrate(valueInterface)
				}
			} else {
				newMap[key] = Hydrate(valueInterface)
			}
		}
		return newMap
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
