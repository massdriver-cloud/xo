package attestation

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const PredicateTypeInventory = "https://massdriver.cloud/attestations/inventory/v1"
const ArtifactTypeInventory = "application/vnd.massdriver.attestation.inventory.v1+json"

type InventoryPredicate struct {
	DeploymentID  string              `json:"deploymentId"`
	BundleName    string              `json:"bundleName"`
	BundleVersion string              `json:"bundleVersion"`
	Project       string              `json:"project"`
	Environment   string              `json:"environment"`
	Account       AccountInfo         `json:"account"`
	Resources     []InventoryResource `json:"resources"`
	GeneratedAt   time.Time           `json:"generatedAt"`
	Producer      ProducerInfo        `json:"producer"`
}

type AccountInfo struct {
	Cloud     string `json:"cloud"` // aws / azure / gcp
	AccountID string `json:"accountId"`
	Region    string `json:"region"`
}

type InventoryResource struct {
	Type string            `json:"type"` // e.g. aws:rds:db-instance
	ID   string            `json:"id"`
	Name string            `json:"name"`
	Tags map[string]string `json:"tags,omitempty"`
}

type ProducerInfo struct {
	Tool    string `json:"tool"`
	Version string `json:"version"`
}

func NewStatementFromInventory(bundleDigest string, pred InventoryPredicate) (*Statement, error) {
	// Validate required fields
	if bundleDigest == "" {
		return nil, fmt.Errorf("bundleDigest is required")
	}
	if pred.DeploymentID == "" {
		return nil, fmt.Errorf("deploymentID is required")
	}
	if pred.BundleName == "" {
		return nil, fmt.Errorf("bundleName is required")
	}

	predBytes, err := json.Marshal(pred)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal predicate: %w", err)
	}

	return &Statement{
		Type:          StatementTypeV1,
		PredicateType: PredicateTypeInventory,
		Predicate:     predBytes,
		Subject: []Subject{
			{
				Name: "bundle",
				Digest: map[string]string{
					"sha256": strings.TrimPrefix(bundleDigest, "sha256:"),
				},
			},
		},
	}, nil
}
