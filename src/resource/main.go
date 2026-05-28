package resource

import (
	"context"
	"errors"
	"fmt"
	"xo/src/bundle"

	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/provisioning/resources"
)

type ResourceService interface {
	CreateResource(ctx context.Context, r *resources.Resource) (*resources.Resource, error)
	DeleteResource(ctx context.Context, id string) error
}

func getResourceTypeFromBundle(bun *bundle.Bundle, field string) (string, error) {
	properties, exists := bun.Artifacts["properties"].(map[string]interface{})
	if !exists {
		return "", errors.New("malformed resources specification: no properties")
	}

	resourceSpec, exists := properties[field].(map[string]interface{})
	if !exists {
		return "", fmt.Errorf("resource %s does not exist in specification", field)
	}

	resourceType, exists := resourceSpec["$ref"].(string)
	if !exists {
		return "", fmt.Errorf("resource %s does not exist in specification", field)
	}

	return resourceType, nil
}
