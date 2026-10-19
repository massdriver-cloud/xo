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
	// Resources holds per-resource configuration for the resources this step
	// produces, keyed by resource field.
	Resources map[string]StepResource `json:"resources,omitempty" yaml:"resources,omitempty"`
}

type StepResource struct {
	// Template is an inline jq template used to build the resource's payload
	// from the step's inputs and outputs.
	Template string `json:"template,omitempty" yaml:"template,omitempty"`
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

type Resource struct {
	ResourceType string `json:"resource_type" yaml:"resource_type"`
	Required     bool   `json:"required" yaml:"required"`
}
type Dependency struct {
	ResourceType string `json:"resource_type" yaml:"resource_type"`
	Required     bool   `json:"required" yaml:"required"`
}

type Bundle struct {
	Schema       string                 `json:"schema" yaml:"schema"`
	Name         string                 `json:"name" yaml:"name"`
	Description  string                 `json:"description" yaml:"description"`
	Type         string                 `json:"type" yaml:"type"`
	App          *AppBlock              `json:"app" yaml:"app"`
	Steps        []BundleStep           `json:"steps" yaml:"steps"`
	Resources    map[string]Resource    `json:"resources" yaml:"resources"`
	Dependencies map[string]Dependency  `json:"dependencies" yaml:"dependencies"`
	Artifacts    map[string]interface{} `json:"artifacts" yaml:"artifacts"`
	Params       map[string]interface{} `json:"params" yaml:"params"`
	Connections  map[string]interface{} `json:"connections" yaml:"connections"`
	Ui           map[string]interface{} `json:"ui" yaml:"ui"`
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

// Step returns the step with the given path.
func (b *Bundle) Step(path string) (*BundleStep, bool) {
	for i := range b.Steps {
		if b.Steps[i].Path == path {
			return &b.Steps[i], true
		}
	}
	return nil, false
}

// ResourceFields returns the set of resource fields the bundle declares, from
// both the `resources` block and the legacy `artifacts` schema properties.
func (b *Bundle) ResourceFields() map[string]bool {
	fields := map[string]bool{}
	for field := range b.Resources {
		fields[field] = true
	}
	if props, ok := b.Artifacts["properties"].(map[string]interface{}); ok {
		for field := range props {
			fields[field] = true
		}
	}
	return fields
}
