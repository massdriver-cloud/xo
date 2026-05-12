package attestation

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewStatementFromInventory_Success(t *testing.T) {
	pred := InventoryPredicate{
		DeploymentID:  "deploy-123",
		BundleName:    "test-bundle",
		BundleVersion: "v1.0.0",
		Project:       "test-project",
		Environment:   "production",
		Account: AccountInfo{
			Cloud:     "aws",
			AccountID: "123456789012",
			Region:    "us-east-1",
		},
		Resources: []InventoryResource{
			{
				Type: "aws:rds:db-instance",
				ID:   "my-db-instance",
				Name: "production-db",
				Tags: map[string]string{
					"Environment": "production",
				},
			},
		},
		GeneratedAt: time.Now(),
		Producer: ProducerInfo{
			Tool:    "xo",
			Version: "1.0.0",
		},
	}

	stmt, err := NewStatementFromInventory("sha256:abc123", pred)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stmt.Type != StatementTypeV1 {
		t.Errorf("expected Type %s, got %s", StatementTypeV1, stmt.Type)
	}

	if stmt.PredicateType != PredicateTypeInventory {
		t.Errorf("expected PredicateType %s, got %s", PredicateTypeInventory, stmt.PredicateType)
	}

	if len(stmt.Subject) != 1 {
		t.Fatalf("expected 1 subject, got %d", len(stmt.Subject))
	}

	if stmt.Subject[0].Name != "bundle" {
		t.Errorf("expected subject name 'bundle', got %s", stmt.Subject[0].Name)
	}

	if stmt.Subject[0].Digest["sha256"] != "abc123" {
		t.Errorf("expected digest 'abc123', got %s", stmt.Subject[0].Digest["sha256"])
	}

	// Verify predicate can be unmarshaled
	var unmarshaledPred InventoryPredicate
	if err := json.Unmarshal(stmt.Predicate, &unmarshaledPred); err != nil {
		t.Fatalf("failed to unmarshal predicate: %v", err)
	}

	if unmarshaledPred.DeploymentID != pred.DeploymentID {
		t.Errorf("expected DeploymentID %s, got %s", pred.DeploymentID, unmarshaledPred.DeploymentID)
	}
}

func TestNewStatementFromInventory_ValidationErrors(t *testing.T) {
	tests := []struct {
		name         string
		bundleDigest string
		pred         InventoryPredicate
		wantErr      string
	}{
		{
			name:         "empty bundle digest",
			bundleDigest: "",
			pred: InventoryPredicate{
				DeploymentID: "deploy-123",
				BundleName:   "test-bundle",
			},
			wantErr: "bundleDigest is required",
		},
		{
			name:         "empty deployment ID",
			bundleDigest: "sha256:abc123",
			pred: InventoryPredicate{
				BundleName: "test-bundle",
			},
			wantErr: "deploymentID is required",
		},
		{
			name:         "empty bundle name",
			bundleDigest: "sha256:abc123",
			pred: InventoryPredicate{
				DeploymentID: "deploy-123",
			},
			wantErr: "bundleName is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewStatementFromInventory(tt.bundleDigest, tt.pred)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("expected error %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestNewStatementFromInventory_DigestPrefixStripping(t *testing.T) {
	pred := InventoryPredicate{
		DeploymentID: "deploy-123",
		BundleName:   "test-bundle",
	}

	tests := []struct {
		name           string
		input          string
		expectedDigest string
	}{
		{
			name:           "with sha256 prefix",
			input:          "sha256:abc123",
			expectedDigest: "abc123",
		},
		{
			name:           "without prefix",
			input:          "abc123",
			expectedDigest: "abc123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := NewStatementFromInventory(tt.input, pred)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if stmt.Subject[0].Digest["sha256"] != tt.expectedDigest {
				t.Errorf("expected digest %s, got %s", tt.expectedDigest, stmt.Subject[0].Digest["sha256"])
			}
		})
	}
}

func TestInventoryPredicate_JSONSerialization(t *testing.T) {
	pred := InventoryPredicate{
		DeploymentID:  "deploy-123",
		BundleName:    "test-bundle",
		BundleVersion: "v1.0.0",
		Project:       "test-project",
		Environment:   "production",
		Account: AccountInfo{
			Cloud:     "aws",
			AccountID: "123456789012",
			Region:    "us-east-1",
		},
		Resources: []InventoryResource{
			{
				Type: "aws:rds:db-instance",
				ID:   "db-123",
				Name: "my-db",
			},
		},
		GeneratedAt: time.Date(2025, 12, 10, 0, 0, 0, 0, time.UTC),
		Producer: ProducerInfo{
			Tool:    "xo",
			Version: "1.0.0",
		},
	}

	data, err := json.Marshal(pred)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled InventoryPredicate
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if unmarshaled.DeploymentID != pred.DeploymentID {
		t.Errorf("expected DeploymentID %s, got %s", pred.DeploymentID, unmarshaled.DeploymentID)
	}

	if len(unmarshaled.Resources) != len(pred.Resources) {
		t.Errorf("expected %d resources, got %d", len(pred.Resources), len(unmarshaled.Resources))
	}
}
