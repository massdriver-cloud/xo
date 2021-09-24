package terraform

import (
	"encoding/json"
	"io"
	"os"
	"path"
	"xo/src/bundles"
	"xo/src/jsonschema"
)

func GenerateFiles(baseDir string) error {
	massdriverVariables := map[string]interface{}{
		"variable": map[string]interface{}{
			"md_name_prefix": map[string]string{
				"type": "string",
			},
			"md_default_tags": map[string]string{
				"type": "map",
			},
		},
	}

	err := os.MkdirAll(path.Join(baseDir, "src"), 0755)
	if err != nil {
		return err
	}

	paramsVariablesFile, err := os.Create(path.Join(baseDir, "src", "_params_variables.tf.json"))
	if err != nil {
		return err
	}
	err = Compile(path.Join(baseDir, bundles.ParamsSchemaFilename), paramsVariablesFile)
	if err != nil {
		return err
	}

	connectionsVariablesFile, err := os.Create(path.Join(baseDir, "src", "_connections_variables.tf.json"))
	if err != nil {
		return err
	}
	err = Compile(path.Join(baseDir, bundles.ConnectionsSchemaFilename), connectionsVariablesFile)
	if err != nil {
		return err
	}

	massdriverVariablesFile, err := os.Create(path.Join(baseDir, "src", "_md_variables.tf.json"))
	if err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(massdriverVariables, "", "  ")
	if err != nil {
		return err
	}
	_, err = massdriverVariablesFile.Write(bytes)

	return err
}

// Compile a JSON Schema to Terraform Variable Definition JSON
func Compile(path string, out io.Writer) error {
	vars, err := getVars(path)
	if err != nil {
		return err
	}

	// You can't have an empty variable block, so if there are no vars return an empty json block
	if len(vars) == 0 {
		out.Write([]byte("{}"))
		return nil
	}

	variableFile := TFVariableFile{Variable: vars}

	bytes, err := json.MarshalIndent(variableFile, "", "  ")
	if err != nil {
		return err
	}

	_, err = out.Write(bytes)

	return err
}

func getVars(path string) (map[string]TFVariable, error) {
	variables := map[string]TFVariable{}
	schema, err := jsonschema.GetJSONSchema(path)
	if err != nil {
		return variables, err
	}

	for name, prop := range schema.Properties {
		variables[name] = NewTFVariable(prop)
	}
	return variables, nil
}
