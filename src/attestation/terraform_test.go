package attestation

import (
	"testing"
)

func TestInventoryFromTerraformState(t *testing.T) {
	stateJSON := []byte(`{
  "version": 4,
  "resources": [
    {
      "type": "aws_db_instance",
      "name": "main",
      "provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
      "mode": "managed",
      "instances": [
        {
          "attributes": {
            "id": "my-db-instance",
            "name": "production-db",
            "tags": {
              "Environment": "production",
              "ManagedBy": "massdriver"
            }
          }
        }
      ]
    },
    {
      "type": "aws_security_group",
      "name": "db_sg",
      "provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
      "mode": "managed",
      "instances": [
        {
          "attributes": {
            "id": "sg-123456",
            "name": "db-security-group"
          }
        }
      ]
    },
    {
      "type": "aws_availability_zones",
      "name": "available",
      "provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
      "mode": "data",
      "instances": [
        {
          "attributes": {
            "id": "us-east-1",
            "names": ["us-east-1a", "us-east-1b"]
          }
        }
      ]
    }
  ]
}`)

	account := AccountInfo{
		Cloud:     "aws",
		AccountID: "123456789012",
		Region:    "us-east-1",
	}

	pred, err := InventoryFromTerraformState(
		stateJSON,
		"deploy-123",
		"my-bundle",
		"v1.0.0",
		"my-project",
		"production",
		account,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if pred.DeploymentID != "deploy-123" {
		t.Errorf("expected deploymentID 'deploy-123', got %s", pred.DeploymentID)
	}

	if pred.BundleName != "my-bundle" {
		t.Errorf("expected bundleName 'my-bundle', got %s", pred.BundleName)
	}

	// Should have 2 managed resources (data sources excluded)
	if len(pred.Resources) != 2 {
		t.Errorf("expected 2 resources, got %d", len(pred.Resources))
	}

	// Verify first resource
	if len(pred.Resources) > 0 {
		res := pred.Resources[0]
		if res.Type != "aws_db_instance" {
			t.Errorf("expected type 'aws_db_instance', got %s", res.Type)
		}
		if res.ID != "my-db-instance" {
			t.Errorf("expected ID 'my-db-instance', got %s", res.ID)
		}
		if res.Name != "production-db" {
			t.Errorf("expected name 'production-db', got %s", res.Name)
		}
		if res.Tags["Environment"] != "production" {
			t.Errorf("expected Environment tag 'production', got %s", res.Tags["Environment"])
		}
	}

	// Verify second resource (without tags)
	if len(pred.Resources) > 1 {
		res := pred.Resources[1]
		if res.Type != "aws_security_group" {
			t.Errorf("expected type 'aws_security_group', got %s", res.Type)
		}
		if res.ID != "sg-123456" {
			t.Errorf("expected ID 'sg-123456', got %s", res.ID)
		}
	}

	if pred.Producer.Tool != "terraform" {
		t.Errorf("expected producer tool 'terraform', got %s", pred.Producer.Tool)
	}
}

func TestInventoryFromTerraformState_InvalidJSON(t *testing.T) {
	stateJSON := []byte(`{invalid json`)

	account := AccountInfo{
		Cloud:     "aws",
		AccountID: "123456789012",
		Region:    "us-east-1",
	}

	_, err := InventoryFromTerraformState(
		stateJSON,
		"deploy-123",
		"my-bundle",
		"v1.0.0",
		"my-project",
		"production",
		account,
	)

	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestInventoryFromTerraformState_EmptyState(t *testing.T) {
	stateJSON := []byte(`{
  "version": 4,
  "resources": []
}`)

	account := AccountInfo{
		Cloud:     "aws",
		AccountID: "123456789012",
		Region:    "us-east-1",
	}

	pred, err := InventoryFromTerraformState(
		stateJSON,
		"deploy-123",
		"my-bundle",
		"v1.0.0",
		"my-project",
		"production",
		account,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(pred.Resources) != 0 {
		t.Errorf("expected 0 resources, got %d", len(pred.Resources))
	}
}
