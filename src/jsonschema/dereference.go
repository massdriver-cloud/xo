package jsonschema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

var schemaTypePattern = regexp.MustCompile(`^.*\/(.*).json$`)

// A RefdSchema is a JSON Schema that may contain $ref
type RefdSchema struct {
	SchemaId   string
	Definition interface{}
}

func WriteDereferencedSchema(schemaFilePath string, outDir string) error {
	dereferencedSchema := RefdSchema{}
	orderedJson := OrderedJSON{}
	cwd := filepath.Dir(schemaFilePath)
	data, err := ioutil.ReadFile(schemaFilePath)
	if err != nil {
		return err
	}

	json.Unmarshal(data, &orderedJson)
	definition, err := Hydrate(orderedJson, cwd)
	if err != nil {
		return err
	}
	dereferencedSchema.Definition = definition

	for _, ele := range orderedJson {
		if ele.Key == "$id" {
			dereferencedSchema.SchemaId = ele.Value.(string)
			break
		}
	}

	schemaFileName := fmt.Sprintf("%s.compiled.json", dereferencedSchema.Type())
	path := path.Join(outDir, schemaFileName)

	file, err := os.Create(path)
	defer file.Close()

	if err != nil {
		return err
	}

	json, err := json.Marshal(dereferencedSchema.Definition)
	if err != nil {
		return err
	}

	_, err = file.Write(json)

	return err
}

func (r *RefdSchema) Type() string {
	bytes := []byte(r.SchemaId)
	match := schemaTypePattern.FindSubmatch(bytes)[1]
	artifactType := string(match)

	return artifactType
}
