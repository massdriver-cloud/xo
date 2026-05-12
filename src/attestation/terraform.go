package attestation

import (
	"encoding/json"
	"fmt"
	"time"
)

// TerraformResource represents a resource in Terraform state
type TerraformResource struct {
	Type      string              `json:"type"`
	Name      string              `json:"name"`
	Provider  string              `json:"provider"`
	Instances []TerraformInstance `json:"instances"`
	Mode      string              `json:"mode"` // "managed" or "data"
}

// TerraformInstance represents a resource instance
type TerraformInstance struct {
	Attributes map[string]interface{} `json:"attributes"`
}

// TerraformState represents simplified Terraform state structure
type TerraformState struct {
	Version   int                 `json:"version"`
	Resources []TerraformResource `json:"resources"`
}

// InventoryFromTerraformState converts Terraform state JSON to an InventoryPredicate
// This is a helper function to generate inventory attestations from Terraform deployments
func InventoryFromTerraformState(
	stateJSON []byte,
	deploymentID string,
	bundleName string,
	bundleVersion string,
	project string,
	environment string,
	account AccountInfo,
) (*InventoryPredicate, error) {
	var state TerraformState
	if err := json.Unmarshal(stateJSON, &state); err != nil {
		return nil, fmt.Errorf("failed to parse terraform state: %w", err)
	}

	var resources []InventoryResource
	for _, tfResource := range state.Resources {
		// Skip data sources, only include managed resources
		if tfResource.Mode != "managed" {
			continue
		}

		for _, instance := range tfResource.Instances {
			resource := InventoryResource{
				Type: tfResource.Type,
				Tags: make(map[string]string),
			}

			// Extract common attributes
			if id, ok := instance.Attributes["id"].(string); ok {
				resource.ID = id
			}
			if name, ok := instance.Attributes["name"].(string); ok {
				resource.Name = name
			}

			// Extract tags if present
			if tags, ok := instance.Attributes["tags"].(map[string]interface{}); ok {
				for k, v := range tags {
					if strVal, ok := v.(string); ok {
						resource.Tags[k] = strVal
					}
				}
			}

			// Only add resources with valid IDs
			if resource.ID != "" {
				resources = append(resources, resource)
			}
		}
	}

	return &InventoryPredicate{
		DeploymentID:  deploymentID,
		BundleName:    bundleName,
		BundleVersion: bundleVersion,
		Project:       project,
		Environment:   environment,
		Account:       account,
		Resources:     resources,
		GeneratedAt:   time.Now(),
		Producer: ProducerInfo{
			Tool:    "terraform",
			Version: "unknown", // TODO: extract from state if available
		},
	}, nil
}
