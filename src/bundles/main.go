package bundles

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const idUrlPattern = "https://massdriver.sh/schemas/bundles/%s/schema-%s.json"
const jsonSchemaUrlPattern = "http://json-schema.org/%s/schema"

type Bundle struct {
	Uuid        string      `json:"uuid"`
	Schema      string      `json:"schema"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Type        string      `json:"type"`
	Artifacts   OrderedJSON `json:"artifacts"`
	Params      OrderedJSON `json:"params"`
	Connections OrderedJSON `json:"connections"`
}

// ParseBundle parses a bundle from a YAML file
// bundle, err := ParseBundle("./bundle.yaml")
// Generate the files in this directory
// bundle.Build(".")
func ParseBundle(path string) (Bundle, error) {
	bundle := Bundle{}
	cwd := filepath.Dir(path)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return bundle, err
	}

	err = yaml.Unmarshal([]byte(data), &bundle)
	if err != nil {
		return bundle, err
	}

	hydratedArtifacts, err := Hydrate(bundle.Artifacts, cwd)
	if err != nil {
		return bundle, err
	}

	bundle.Artifacts = hydratedArtifacts.(OrderedJSON)

	hydratedParams, err := Hydrate(bundle.Params, cwd)
	if err != nil {
		return bundle, err
	}
	bundle.Params = hydratedParams.(OrderedJSON)

	hydratedConnections, err := Hydrate(bundle.Connections, cwd)
	if err != nil {
		return bundle, err
	}
	bundle.Connections = hydratedConnections.(OrderedJSON)

	return bundle, nil
}

// Metadata returns common metadata fields for each JSON Schema
func (b *Bundle) Metadata(schemaType string) OrderedJSON {
	return OrderedJSON([]OrderedJSONElement{
		{Key: "$schema", Value: generateSchemaUrl(b.Schema)},
		{Key: "$id", Value: generateIdUrl(b.Type, schemaType)},
		{Key: "title", Value: b.Title},
		{Key: "description", Value: b.Description},
	})
}

func createFile(dir string, fileName string) (*os.File, error) {
	filePath := fmt.Sprintf("%s/schema-%s.json", dir, fileName)
	return os.Create(filePath)
}

// Build generates all bundle files in the given directory
func (b *Bundle) Build(dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	paramsSchemaFile, err := createFile(dir, "params")
	if err != nil {
		return err
	}

	connectionsSchemaFile, err := createFile(dir, "connections")
	if err != nil {
		return err
	}

	artifactsSchemaFile, err := createFile(dir, "artifacts")
	if err != nil {
		return err
	}

	err = BuildSchema(b.Params, b.Metadata("params"), paramsSchemaFile)
	if err != nil {
		return err
	}
	err = BuildSchema(b.Connections, b.Metadata("connections"), connectionsSchemaFile)
	if err != nil {
		return err
	}
	err = BuildSchema(b.Artifacts, b.Metadata("artifacts"), artifactsSchemaFile)
	if err != nil {
		return err
	}

	defer paramsSchemaFile.Close()
	defer connectionsSchemaFile.Close()
	defer artifactsSchemaFile.Close()

	return nil
}

// BuildSchema generates schema-*.json files
func BuildSchema(schema OrderedJSON, metadata OrderedJSON, buffer io.Writer) error {
	var err error
	var mergedSchema = OrderedJSON(append(metadata, schema...))

	json, err := json.Marshal(mergedSchema)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(buffer, string(json))
	if err != nil {
		return err
	}

	return nil
}

func generateIdUrl(mdType string, schemaType string) string {
	return fmt.Sprintf(idUrlPattern, mdType, schemaType)
}

func generateSchemaUrl(schema string) string {
	return fmt.Sprintf(jsonSchemaUrlPattern, schema)
}
