package jsonschema

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

const filePrefix = "file://"

var loaderPrefixPattern = regexp.MustCompile(`^(file|http|https)://`)

// Load a JSON file with or without a path prefix
func Load(path string) (gojsonschema.JSONLoader, error) {
	switch ext := filepath.Ext(path); ext {
	case ".json":
		var ref string
		if loaderPrefixPattern.MatchString(path) {
			ref = path
		} else {
			ref = filePrefix + path
		}
		return gojsonschema.NewReferenceLoader(ref), nil
	case ".yaml":
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		b := map[string]interface{}{}
		err = yaml.Unmarshal(bytes, b)
		if err != nil {
			return nil, err
		}
		return gojsonschema.NewGoLoader(b), nil
	default:
		return nil, errors.New("Unsupported file type (only yaml and json supported): " + ext)
	}
}
