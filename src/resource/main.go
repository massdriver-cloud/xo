package resource

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"xo/src/bundle"

	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/provisioning/resources"
)

type ResourceService interface {
	CreateResource(ctx context.Context, input *resources.ResourceInput) (*resources.Resource, error)
	DeleteResource(ctx context.Context, id string) error
	// ListResources lists the resources currently produced by the deployment's instance.
	ListResources(ctx context.Context) ([]resources.Resource, error)
}

type Result struct {
	Published []string
	Deleted   []string
}

func findStep(b bundle.Bundle, path string) (*bundle.BundleStep, error) {
	step, ok := b.Step(path)
	if !ok {
		paths := make([]string, 0, len(b.Steps))
		for _, s := range b.Steps {
			paths = append(paths, s.Path)
		}
		return nil, fmt.Errorf("step %q not found in massdriver.yaml (steps: %s)", path, strings.Join(paths, ", "))
	}
	return step, nil
}

func resourceID(instanceID string, res resources.Resource) string {
	if res.ID != "" {
		return res.ID
	}
	return instanceID + "." + res.Field
}

func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
