package artifact

import (
	"context"
	"errors"
	"fmt"
	"xo/src/bundle"

	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/services/artifacts"
)

type ArtifactService interface {
	CreateArtifact(ctx context.Context, a *artifacts.Artifact) (*artifacts.Artifact, error)
	DeleteArtifact(ctx context.Context, id, field string) error
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
