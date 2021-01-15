package schemaloader

import (
	"os"
	"path"
	"regexp"

	"github.com/xeipuuv/gojsonschema"
)

// Load a JSON Schema with our without a path prefix
func Load(path string) gojsonschema.JSONLoader {
	loaderPrefixPattern := regexp.MustCompile(`^(file|http|https)://`)
	var ref string
	if loaderPrefixPattern.MatchString(path) {
		ref = expandPath(path)
	} else {
		ref = "file://" + expandPath(path)
	}

	return gojsonschema.NewReferenceLoader(ref)
}

func expandPath(p string) string {
	if path.IsAbs(p) {
		return p
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}

	return path.Join(cwd, p)
}
