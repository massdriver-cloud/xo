package bundles

import (
	"fmt"
	"os"
	"path"
	"xo/src/jsonschema"
	// have to use v3 here so it decodes to map[string]interface{} instead of map[interface{}]interface{} https://github.com/go-yaml/yaml/issues/139
)

func LintBundle(bundlePath string) (bool, error) {
	// Make sure bundle directory exists
	bundleFolderExists, err := lintDirExists(bundlePath)
	if err != nil {
		return false, err
	} else if !bundleFolderExists {
		fmt.Println("Specified bundle directory does not exist")
		return false, nil
	}

	// Make sure bundle.yaml exists and is valid
	bundleYamlPath := path.Join(bundlePath, "bundle.yaml")
	bundleSchemaPath := "./bundle-schema.json"
	fileValid, err := lintValidateJson(bundleSchemaPath, bundleYamlPath)
	if err != nil {
		return false, err
	} else if !fileValid {
		return false, nil
	}

	// Make sure params, connections and artifacts schemas exist and are valid
	for _, v := range []string{"schema-artifacts.json", "schema-connections.json", "schema-params.json"} {
		filePath := path.Join(bundlePath, v)
		schemaPath := "./meta-schema-v7.json"
		fileValid, err := lintValidateJson(schemaPath, filePath)
		if err != nil {
			return false, err
		} else if !fileValid {
			return false, nil
		}
	}

	// Make sure schema.stories.js exists and is valid
	storiesPath := path.Join(bundlePath, "schema.stories.js")
	storiesExists, err := lintFileExists(storiesPath)
	if err != nil {
		return false, err
	} else if !storiesExists {
		return false, nil
	}

	return true, nil
}

func lintDirExists(dirPath string) (bool, error) {
	info, err := os.Stat(dirPath)
	if err == nil && info.Mode().IsDir() {
		return true, nil
	}
	if os.IsNotExist(err) || !info.IsDir() {
		fmt.Println("Bundle invalid: directory " + dirPath + " doesn't exist")
		return false, nil
	}
	return false, err
}

func lintFileExists(filePath string) (bool, error) {
	info, err := os.Stat(filePath)
	if err == nil && !info.Mode().IsDir() {
		return true, nil
	}
	if os.IsNotExist(err) || info.IsDir() {
		fmt.Println("Bundle invalid: File " + filePath + " doesn't exist")
		return false, nil
	}
	return false, err
}

func lintValidateJson(schemaPath string, jsonPath string) (bool, error) {
	fileExists, err := lintFileExists(jsonPath)
	if err != nil {
		return false, err
	} else if !fileExists {
		return false, nil
	}
	fileValid, err := jsonschema.Validate(schemaPath, jsonPath)
	if err != nil {
		return false, err
	} else if !fileValid {
		return false, nil
	}
	return true, nil
}
