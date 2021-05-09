package bundles

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
)

var ArtifactPath = "./definitions/artifacts"
var SpecPath = "./definitions/specs"

// relativeFilePathPattern only accepts relative file path prefixes "./" and "../"
var relativeFilePathPattern = regexp.MustCompile(`^(\.\/|\.\.\/)`)
var artifactPattern = regexp.MustCompile("^artifact://([a-z0-9-]+)")
var specPattern = regexp.MustCompile("^spec://([a-z0-9-]+)")

func Hydrate(any interface{}) interface{} {
	val := getValue(any)

	switch val.Kind() {
	case reflect.String:
		if artifactPattern.MatchString(val.String()) {
			artifact, err := readArtifactRef(val.String())
			maybePanic(err)

			// TODO: Do we want to recursively hydrate specs. We could replace $ref's w/ specs/artifacts
			// and fully hydrate a snapshot of the entire JSON Schema into one file for the bundle... which
			// would stop any weirdness in file changes between deploys/caching
			return Hydrate(artifact)
		} else if specPattern.MatchString(val.String()) {
			spec, err := readSpecRef(val.String())
			maybePanic(err)

			// TODO: Do we want to recursively hydrate specs. We could replace $ref's w/ specs/artifacts
			// and fully hydrate a snapshot of the entire JSON Schema into one file for the bundle... which
			// would stop any weirdness in file changes between deploys/caching
			return Hydrate(spec)
		} else {
			return val.String()
		}
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

					// TODO: this isn't deterministic...
					// if key a exists in both 'new' and ref'd, then depending on the order
					// of the key being processed it may or may not be overwritten by the ref (which we _dont_ want)

					// TODO: When recursively resolving, we dont know _where_ we started
					// for the relative path...
					// if we ref ./path/a.json it refs "./b.json"
					// this will look in b.json ont path/b.json
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

func readSpecRef(ref string) (map[string]interface{}, error) {
	refBytes := ([]byte(ref))
	m := specPattern.FindSubmatch(refBytes)

	filename := string(m[1])
	filepath := fmt.Sprintf("%s/%s.json", SpecPath, filename)
	return readJsonFile((filepath))
}

func readArtifactRef(ref string) (map[string]interface{}, error) {
	refBytes := ([]byte(ref))
	m := artifactPattern.FindSubmatch(refBytes)

	filename := string(m[1])
	filepath := fmt.Sprintf("%s/%s.json", ArtifactPath, filename)
	return readJsonFile((filepath))
}
