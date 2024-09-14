package bundle

import (
	"io/ioutil"
	"path/filepath"
	"xo/src/jsonschema"

	"gopkg.in/yaml.v3"
)

const ArtifactsSchemaFilename = "schema-artifacts.json"
const ConnectionsSchemaFilename = "schema-connections.json"
const ParamsSchemaFilename = "schema-params.json"
const UiSchemaFilename = "schema-ui.json"

type BundleStep struct {
	Path        string `json:"path" yaml:"path"`
	Provisioner string `json:"provisioner" yaml:"provisioner"`
}
type SecretsBlock struct {
	Required    bool   `json:"required" yaml:"required"`
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
}
type AppBlock struct {
	Envs     map[string]string       `json:"envs" yaml:"envs"`
	Policies []string                `json:"policies" yaml:"policies"`
	Secrets  map[string]SecretsBlock `json:"secrets" yaml:"secrets"`
}

type Bundle struct {
	Schema      string                 `json:"schema" yaml:"schema"`
	Name        string                 `json:"name" yaml:"name"`
	Description string                 `json:"description" yaml:"description"`
	Type        string                 `json:"type" yaml:"type"`
	App         *AppBlock              `json:"app" yaml:"app"`
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
