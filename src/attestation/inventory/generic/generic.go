// Package generic builds inventory assets from a caller-supplied list, for
// provisioners that don't have a built-in parser. With no input, it yields no
// assets (an empty inventory).
package generic

import (
	"encoding/json"
	"fmt"

	"xo/src/attestation/inventory"

	v1 "github.com/in-toto/attestation/go/v1"
)

// asset is the Massdriver asset schema a custom provisioner emits.
type asset struct {
	URI        string            `json:"uri"`
	Name       string            `json:"name,omitempty"`
	Type       string            `json:"type,omitempty"`
	Digest     map[string]string `json:"digest,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// Extractor reads a JSON array of assets (the generic schema).
type Extractor struct{}

func (Extractor) Assets(input []byte, attributes map[string]string) ([]*v1.ResourceDescriptor, error) {
	if len(input) == 0 {
		return nil, nil // no assets supplied; the inventory is empty
	}

	var items []asset
	if err := json.Unmarshal(input, &items); err != nil {
		return nil, fmt.Errorf("failed to parse assets file: %w", err)
	}

	var assets []*v1.ResourceDescriptor
	for _, item := range items {
		if item.URI == "" {
			return nil, fmt.Errorf("each asset requires a uri")
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

		a, err := inventory.NewAsset(item.URI, item.Name, digest, annotations)
		if err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}

	return assets, nil
}
