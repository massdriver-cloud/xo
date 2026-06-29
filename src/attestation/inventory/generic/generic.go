// Package generic builds inventory resources from a caller-supplied list, for
// provisioners that don't have a built-in parser. With no input, it yields no
// resources (an empty inventory).
package generic

import (
	"encoding/json"
	"fmt"

	"xo/src/attestation/inventory"

	v1 "github.com/in-toto/attestation/go/v1"
)

// resource is the Massdriver resource schema a custom provisioner emits.
type resource struct {
	URI        string            `json:"uri"`
	Name       string            `json:"name,omitempty"`
	Type       string            `json:"type,omitempty"`
	Digest     map[string]string `json:"digest,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// Extractor reads a JSON array of resources (the generic schema).
type Extractor struct{}

func (Extractor) Resources(input []byte, attributes map[string]string) ([]*v1.ResourceDescriptor, error) {
	if len(input) == 0 {
		return nil, nil // no resources supplied; the inventory is empty
	}

	var items []resource
	if err := json.Unmarshal(input, &items); err != nil {
		return nil, fmt.Errorf("failed to parse resources file: %w", err)
	}

	var resources []*v1.ResourceDescriptor
	for _, item := range items {
		if item.URI == "" {
			return nil, fmt.Errorf("each resource requires a uri")
		}
		digest := item.Digest
		if len(digest) == 0 {
			digest = inventory.IdentityDigest(item.URI)
		}

		annotations := map[string]string{}
		if item.Type != "" {
			annotations["type"] = item.Type
		}
		for k, v := range item.Attributes {
			annotations[k] = v
		}
		for k, v := range attributes {
			annotations[k] = v
		}

		r, err := inventory.NewResource(item.URI, item.Name, digest, annotations)
		if err != nil {
			return nil, err
		}
		resources = append(resources, r)
	}

	return resources, nil
}
