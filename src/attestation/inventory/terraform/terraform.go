// Package terraform extracts inventory assets from `terraform show -json`
// output (also covers OpenTofu's `tofu show -json`, same format).
package terraform

import (
	"encoding/json"
	"fmt"
	"strings"

	"xo/src/attestation/inventory"

	tfjson "github.com/hashicorp/terraform-json"
	v1 "github.com/in-toto/attestation/go/v1"
)

// Extractor reads `terraform show -json` output (parsed with
// hashicorp/terraform-json — the supported, versioned format).
type Extractor struct{}

func (Extractor) Assets(showJSON []byte, attributes map[string]string) ([]*v1.ResourceDescriptor, error) {
	var state tfjson.State
	if err := json.Unmarshal(showJSON, &state); err != nil {
		return nil, fmt.Errorf("failed to parse terraform show output: %w", err)
	}
	if state.Values == nil || state.Values.RootModule == nil {
		return nil, nil
	}

	var assets []*v1.ResourceDescriptor
	if err := collectModule(state.Values.RootModule, attributes, &assets); err != nil {
		return nil, err
	}
	return assets, nil
}

// collectModule appends a module's managed resources (and its children's) as
// inventory assets.
func collectModule(module *tfjson.StateModule, attributes map[string]string, assets *[]*v1.ResourceDescriptor) error {
	for _, resource := range module.Resources {
		// Only managed resources are products of the apply; skip data sources.
		if resource.Mode != tfjson.ManagedResourceMode {
			continue
		}
		id, _ := resource.AttributeValues["id"].(string)
		// Only include resources with a stable identifier.
		if id == "" {
			continue
		}
		name, _ := resource.AttributeValues["name"].(string)

		digest, err := inventory.ConfigDigest(resource.AttributeValues)
		if err != nil {
			return fmt.Errorf("failed to digest resource config: %w", err)
		}

		asset, err := inventory.NewAsset(id, name, digest, withType(normalizeType(resource.Type), attributes))
		if err != nil {
			return err
		}
		*assets = append(*assets, asset)
	}

	for _, child := range module.ChildModules {
		if err := collectModule(child, attributes, assets); err != nil {
			return err
		}
	}
	return nil
}

// withType merges the normalized resource type with the Massdriver-assigned
// attributes into one annotations map.
func withType(resourceType string, attributes map[string]string) map[string]string {
	annotations := map[string]string{"type": resourceType}
	for k, v := range attributes {
		annotations[k] = v
	}
	return annotations
}

// normalizeType maps a raw provisioner type to a cloud-agnostic scheme:
// "aws_db_instance" -> "aws:db-instance". Best-effort.
func normalizeType(raw string) string {
	provider, rest, found := strings.Cut(raw, "_")
	if !found {
		return raw
	}
	return provider + ":" + strings.ReplaceAll(rest, "_", "-")
}
