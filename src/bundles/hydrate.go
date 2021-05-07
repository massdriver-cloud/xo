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
var artifactPattern = regexp.MustCompile("^artifact://([a-z0-9-]+)")
var specPattern = regexp.MustCompile("^spec://([a-z0-9-]+)")

func Hydrate(any interface{}) interface{} {
	val := getValue(any)

	switch val.Kind() {
	case reflect.String:
		if artifactPattern.MatchString(val.String()) {
			artifact, err := readArtifactRef(val.String())
			if err != nil {
				panic(err)
			}
			// TODO: Do we want to recursively hydrate specs. We could replace $ref's w/ specs/artifacts
			// and fully hydrate a snapshot of the entire JSON Schema into one file for the bundle... which
			// would stop any weirdness in file changes between deploys/caching
			return artifact
		} else if specPattern.MatchString(val.String()) {
			spec, err := readSpecRef(val.String())
			if err != nil {
				panic(err)
			}
			// TODO: Do we want to recursively hydrate specs. We could replace $ref's w/ specs/artifacts
			// and fully hydrate a snapshot of the entire JSON Schema into one file for the bundle... which
			// would stop any weirdness in file changes between deploys/caching
			return spec
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
		for _, key := range val.MapKeys() {
			hydratedVal := Hydrate(val.MapIndex(key).Interface())
			newMap[key.String()] = hydratedVal
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

func readSpecRef(ref string) (map[string]interface{}, error) {
	refBytes := ([]byte(ref))
	m := specPattern.FindSubmatch(refBytes)

	filename := string(m[1])
	filepath := fmt.Sprintf("%s/%s.json", SpecPath, filename)
	data, err := ioutil.ReadFile(filepath)

	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, err
	}

	return result, err
}

func readArtifactRef(ref string) (map[string]interface{}, error) {
	refBytes := ([]byte(ref))
	m := artifactPattern.FindSubmatch(refBytes)

	filename := string(m[1])
	filepath := fmt.Sprintf("%s/%s.json", ArtifactPath, filename)
	data, err := ioutil.ReadFile(filepath)

	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, err
	}

	return result, err
}
