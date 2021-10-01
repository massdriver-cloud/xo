package bundles

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"xo/src/jsonschema"

	"gopkg.in/yaml.v2"
)

const ArtifactsSchemaFilename = "schema-artifacts.json"
const ConnectionsSchemaFilename = "schema-connections.json"
const ParamsSchemaFilename = "schema-params.json"

const idUrlPattern = "https://massdriver.sh/schemas/bundles/%s/schema-%s.json"
const jsonSchemaUrlPattern = "http://json-schema.org/%s/schema"

type Bundle struct {
	Uuid        string      `json:"uuid"`
	Schema      string      `json:"schema"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Provisioner string      `json:"provisioner"`
	Type        string      `json:"type"`
	Artifacts   jsonschema.OrderedJSON `json:"artifacts"`
	Params      jsonschema.OrderedJSON `json:"params"`
	Connections jsonschema.OrderedJSON `json:"connections"`
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

	hydratedArtifacts, err := jsonschema.Hydrate(bundle.Artifacts, cwd)
	if err != nil {
		return bundle, err
	}

	bundle.Artifacts = hydratedArtifacts.(jsonschema.OrderedJSON)

	hydratedParams, err := jsonschema.Hydrate(bundle.Params, cwd)
	if err != nil {
		return bundle, err
	}
	bundle.Params = hydratedParams.(jsonschema.OrderedJSON)

	hydratedConnections, err := jsonschema.Hydrate(bundle.Connections, cwd)
	if err != nil {
		return bundle, err
	}
	bundle.Connections = hydratedConnections.(jsonschema.OrderedJSON)

	return bundle, nil
}

// Metadata returns common metadata fields for each JSON Schema
func (b *Bundle) Metadata(schemaType string) jsonschema.OrderedJSON {
	return jsonschema.OrderedJSON([]jsonschema.OrderedJSONElement{
		{Key: "$schema", Value: generateSchemaUrl(b.Schema)},
		{Key: "$id", Value: generateIdUrl(b.Type, schemaType)},
		{Key: "title", Value: b.Title},
		{Key: "description", Value: b.Description},
	})
}

func createFile(dir string, fileName string) (*os.File, error) {
	return os.Create(path.Join(dir, fileName))
}

// Build generates all bundle files in the given bundle
func (b *Bundle) GenerateSchemas(dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	paramsSchemaFile, err := createFile(dir, ParamsSchemaFilename)
	if err != nil {
		return err
	}

	connectionsSchemaFile, err := createFile(dir, ConnectionsSchemaFilename)
	if err != nil {
		return err
	}

	artifactsSchemaFile, err := createFile(dir, ArtifactsSchemaFilename)
	if err != nil {
		return err
	}

	err = GenerateSchema(b.Params, b.Metadata("params"), paramsSchemaFile)
	if err != nil {
		return err
	}
	err = GenerateSchema(b.Connections, b.Metadata("connections"), connectionsSchemaFile)
	if err != nil {
		return err
	}
	err = GenerateSchema(b.Artifacts, b.Metadata("artifacts"), artifactsSchemaFile)
	if err != nil {
		return err
	}

	err = paramsSchemaFile.Close()
	if err != nil {
		return err
	}
	err = connectionsSchemaFile.Close()
	if err != nil {
		return err
	}
	err = artifactsSchemaFile.Close()
	if err != nil {
		return err
	}

	return nil
}

// generateSchema generates a specific schema-*.json file
func GenerateSchema(schema jsonschema.OrderedJSON, metadata jsonschema.OrderedJSON, buffer io.Writer) error {
	var err error
	var mergedSchema = jsonschema.OrderedJSON(append(metadata, schema...))

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
