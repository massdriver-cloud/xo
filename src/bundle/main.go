package bundle

import (
	"os"

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

func ParseBundle(path string) (Bundle, error) {
	bundle := Bundle{}

	data, err := os.ReadFile(path)
	if err != nil {
		return bundle, err
	}

	err = yaml.Unmarshal([]byte(data), &bundle)
	if err != nil {
		return bundle, err
	}

	// Check and initialize any nil maps in the bundle
	if bundle.Artifacts == nil {
		bundle.Artifacts = map[string]interface{}{}
	}
	if bundle.Params == nil {
		bundle.Params = map[string]interface{}{}
	}
	if bundle.Connections == nil {
		bundle.Connections = map[string]interface{}{}
	}
	if bundle.Ui == nil {
		bundle.Ui = map[string]interface{}{}
	}

	return bundle, nil
}
