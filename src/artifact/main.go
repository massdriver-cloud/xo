package artifact

import (
	"errors"
	"fmt"
	"xo/src/bundle"
)

type artifactMetadata struct {
	Field              string `json:"field"`
	ProviderResourceID string `json:"provider_resource_id"`
	Type               string `json:"type"`
	Name               string `json:"name"`
}

func getArtifactTypeFromBundle(bun *bundle.Bundle, field string) (string, error) {
	properties, exists := bun.Artifacts["properties"].(map[string]interface{})
	if !exists {
		return "", errors.New("malformed artifacts specification: no properties")
	}

	artifactSpec, exists := properties[field].(map[string]interface{})
	if !exists {
		return "", fmt.Errorf("artifact %s does not exist in specification", field)
	}

	artifactType, exists := artifactSpec["$ref"].(string)
	if !exists {
		return "", fmt.Errorf("artifact %s does not exist in specification", field)
	}

	return artifactType, nil
}
