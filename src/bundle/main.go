package bundle

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"xo/src/jsonschema"

	"gopkg.in/yaml.v3"
)

const ArtifactsSchemaFilename = "schema-artifacts.json"
const ConnectionsSchemaFilename = "schema-connections.json"
const ParamsSchemaFilename = "schema-params.json"
const UiSchemaFilename = "schema-ui.json"

const idUrlPattern = "https://schemas.massdriver.cloud/schemas/bundles/%s/schema-%s.json"
const jsonSchemaUrlPattern = "http://json-schema.org/%s/schema"

type BundleStep struct {
	Path        string `json:"path" yaml:"path"`
	Provisioner string `json:"provisioner" yaml:"provisioner"`
}

type Bundle struct {
	Schema      string                 `json:"schema" yaml:"schema"`
	Name        string                 `json:"name" yaml:"name"`
	Description string                 `json:"description" yaml:"description"`
	Type        string                 `json:"type" yaml:"type"`
	Steps       []BundleStep           `json:"steps" yaml:"steps"`
	Artifacts   map[string]interface{} `json:"artifacts" yaml:"artifacts"`
	Params      map[string]interface{} `json:"params" yaml:"params"`
	Connections map[string]interface{} `json:"connections" yaml:"connections"`
	Ui          map[string]interface{} `json:"ui" yaml:"ui"`
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
	bundle.Artifacts = hydratedArtifacts.(map[string]interface{})
	err = ApplyTransformations(bundle.Artifacts, artifactsTransformations)
	if err != nil {
		return bundle, err
	}

	hydratedParams, err := jsonschema.Hydrate(bundle.Params, cwd)
	if err != nil {
		return bundle, err
	}
	bundle.Params = hydratedParams.(map[string]interface{})
	err = ApplyTransformations(bundle.Params, paramsTransformations)
	if err != nil {
		return bundle, err
	}

	hydratedConnections, err := jsonschema.Hydrate(bundle.Connections, cwd)
	if err != nil {
		return bundle, err
	}
	bundle.Connections = hydratedConnections.(map[string]interface{})
	err = ApplyTransformations(bundle.Connections, connectionsTransformations)
	if err != nil {
		return bundle, err
	}

	hydratedUi, err := jsonschema.Hydrate(bundle.Ui, cwd)
	if err != nil {
		return bundle, err
	}
	bundle.Ui = hydratedUi.(map[string]interface{})
	err = ApplyTransformations(bundle.Ui, uiTransformations)
	if err != nil {
		return bundle, err
	}

	return bundle, nil
}

// Metadata returns common metadata fields for each JSON Schema
func (b *Bundle) Metadata(schemaType string) map[string]string {
	return map[string]string{
		"$schema":     generateSchemaUrl(b.Schema),
		"$id":         generateIdUrl(b.Name, schemaType),
		"name":        b.Name,
		"description": b.Description,
	}
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

	uiSchemaFile, err := createFile(dir, UiSchemaFilename)
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

	emptyMetadata := make(map[string]string)
	err = GenerateSchema(b.Ui, emptyMetadata, uiSchemaFile)
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
func GenerateSchema(schema map[string]interface{}, metadata map[string]string, buffer io.Writer) error {
	var err error
	var mergedSchema = mergeMaps(schema, metadata)

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

func mergeMaps(a map[string]interface{}, b map[string]string) map[string]interface{} {
	for k, v := range b {
		a[k] = v
	}

	return a
}

func generateIdUrl(mdName string, schemaType string) string {
	return fmt.Sprintf(idUrlPattern, mdName, schemaType)
}

func generateSchemaUrl(schema string) string {
	return fmt.Sprintf(jsonSchemaUrlPattern, schema)
}
