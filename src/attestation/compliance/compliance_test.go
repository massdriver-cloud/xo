package compliance

import (
	"encoding/json"
	"testing"

	"xo/src/attestation"

	"google.golang.org/protobuf/encoding/protojson"
)

func TestNewStatement_Success(t *testing.T) {
	pred := Predicate{
		DeploymentContext: attestation.DeploymentContext{
			DeploymentID: "deploy-123",
			InstanceID:   "inst-1",
			Bundle:       attestation.BundleRef{Name: "my-rds-bundle", Version: "v1.2.3"},
			Producer:     attestation.ProducerInfo{Tool: "xo"},
		},
		Scanners: []Scanner{{Name: "checkov", Version: "3.x"}},
		Summary:  Summary{Passed: 10, Failed: 1},
		Results:  json.RawMessage(`{"runs":[]}`),
	}

	stmt, err := NewStatement("massdriver://deploy-123", pred)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stmt.PredicateType != PredicateType {
		t.Errorf("expected %s, got %s", PredicateType, stmt.PredicateType)
	}
	if err := stmt.Validate(); err != nil {
		t.Fatalf("statement failed in-toto validation: %v", err)
	}

	b, err := protojson.Marshal(stmt)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	predicate := m["predicate"].(map[string]any)
	if predicate["deploymentId"] != "deploy-123" {
		t.Errorf("expected deploymentId 'deploy-123', got %v", predicate["deploymentId"])
	}
	summary := predicate["summary"].(map[string]any)
	if summary["failed"].(float64) != 1 {
		t.Errorf("expected 1 failed, got %v", summary["failed"])
	}
}

func TestNewStatement_ValidationErrors(t *testing.T) {
	if _, err := NewStatement("", Predicate{}); err == nil {
		t.Error("expected error for empty subject URI")
	}
	if _, err := NewStatement("massdriver://x", Predicate{}); err == nil {
		t.Error("expected error for missing deployment ID")
	}
}
