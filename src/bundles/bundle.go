package bundles

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

const idUrlPattern = "https://massdriver.sh/schemas/bundles/%s/schema-inputs.json"
const jsonSchemaUrlPattern = "http://json-schema.org/%s/schema"

type WeakSchema map[string]interface{}
type Bundle struct {
	Uuid        string     `json:"uuid"`
	Schema      string     `json:"schema"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Slug        string     `json:"slug"`
	Artifacts   WeakSchema `json:"artifacts"`
	Inputs      WeakSchema `json:"inputs"`
	Connections WeakSchema `json:"connections"`
}

// ParseBundle parses a bundle from a YAML file
// bundle := ParseBundle("./bundle.yaml")
// Generate the files in this directory
// bundle.Build(".")
func ParseBundle(path string) Bundle {
	bundle := Bundle{}

	data, err := ioutil.ReadFile(path)
	checkErr(err)

	err = yaml.Unmarshal([]byte(data), &bundle)
	checkErr(err)

	hydratedArtifacts := Hydrate(bundle.Artifacts)
	bundle.Artifacts = hydratedArtifacts.(map[string]interface{})

	hydratedInputs := Hydrate(bundle.Inputs)
	bundle.Inputs = hydratedInputs.(map[string]interface{})

	hydratedConnections := Hydrate(bundle.Connections)
	bundle.Connections = hydratedConnections.(map[string]interface{})

	return bundle
}

// Metadata returns common metadata fields for each JSON Schema
func (b *Bundle) Metadata() map[string]string {
	return map[string]string{
		"$schema":     generateSchemaUrl(b.Schema),
		"$id":         generateIdUrl(b.Slug),
		"title":       b.Title,
		"description": b.Description,
	}
}

func createFile(dir string, fileName string) *os.File {
	filePath := fmt.Sprintf("%s/schema-%s.json", dir, fileName)
	f, err := os.Create(filePath)
	checkErr(err)

	return f
}

// Build generates all bundle files in the given directory
func (b *Bundle) Build(dir string) {
	err := os.MkdirAll(dir, 0755)
	checkErr(err)

	inputsSchemaFile := createFile(dir, "inputs")
	connectionsSchemaFile := createFile(dir, "connections")
	artifactsSchemaFile := createFile(dir, "artifacts")

	// TODO: connect to build cmd and run in aws-vpc!

	BuildSchema(b.Inputs, b.Metadata(), inputsSchemaFile)
	BuildSchema(b.Connections, b.Metadata(), connectionsSchemaFile)
	BuildSchema(b.Artifacts, b.Metadata(), artifactsSchemaFile)

	defer inputsSchemaFile.Close()
	defer connectionsSchemaFile.Close()
	defer artifactsSchemaFile.Close()
}

// BuildSchema generates schema-*.json files
func BuildSchema(schema WeakSchema, metadata map[string]string, buffer io.Writer) {
	var err error
	var mergedSchema = mergeMaps(schema, metadata)

	json, err := json.Marshal(mergedSchema)
	checkErr(err)

	_, err = fmt.Fprint(buffer, string(json))
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func mergeMaps(a map[string]interface{}, b map[string]string) map[string]interface{} {
	for k, v := range b {
		a[k] = v
	}

	return a
}

func generateIdUrl(slug string) string {
	return fmt.Sprintf(idUrlPattern, slug)
}

func generateSchemaUrl(schema string) string {
	return fmt.Sprintf(jsonSchemaUrlPattern, schema)
}
