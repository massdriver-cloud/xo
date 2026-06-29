package inventory

import (
	"encoding/json"
	"testing"

	"xo/src/attestation"

	v1 "github.com/in-toto/attestation/go/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

func samplePredicate() Predicate {
	res, _ := NewAsset(
		"arn:aws:rds:::db:prod", "production-db",
		map[string]string{"sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		map[string]string{"type": "aws:db-instance", "md:instance": "inst-1"},
	)
	return Predicate{
		DeploymentContext: attestation.DeploymentContext{
			DeploymentID: "deploy-123",
			InstanceID:   "inst-1",
			Bundle:       attestation.BundleRef{Name: "my-rds-bundle", Version: "v1.2.3"},
			Producer:     attestation.ProducerInfo{Tool: "xo"},
		},
		Provisioner: "terraform",
		Assets:      []*v1.ResourceDescriptor{res},
	}
}

func TestNewStatement_Success(t *testing.T) {
	stmt, err := NewStatement("massdriver://deploy-123", samplePredicate())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stmt.PredicateType != PredicateType {
		t.Errorf("expected %s, got %s", PredicateType, stmt.PredicateType)
	}
	if len(stmt.Subject) != 1 || stmt.Subject[0].Uri != "massdriver://deploy-123" {
		t.Fatalf("expected the deployment as subject, got %+v", stmt.Subject)
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
	pred := m["predicate"].(map[string]any)
	if pred["deploymentId"] != "deploy-123" {
		t.Errorf("expected deploymentId 'deploy-123', got %v", pred["deploymentId"])
	}
	if pred["provisioner"] != "terraform" {
		t.Errorf("expected provisioner 'terraform', got %v", pred["provisioner"])
	}
	assets, ok := pred["assets"].([]any)
	if !ok || len(assets) != 1 {
		t.Fatalf("expected 1 asset in predicate, got %v", pred["assets"])
	}
	first := assets[0].(map[string]any)
	if first["uri"] != "arn:aws:rds:::db:prod" {
		t.Errorf("expected asset uri preserved, got %v", first["uri"])
	}
	ann := first["annotations"].(map[string]any)
	if ann["type"] != "aws:db-instance" {
		t.Errorf("expected asset type annotation preserved, got %v", ann["type"])
	}
}

func TestNewStatement_EmptyAssets(t *testing.T) {
	pred := Predicate{DeploymentContext: attestation.DeploymentContext{DeploymentID: "deploy-123"}}
	stmt, err := NewStatement("massdriver://deploy-123", pred)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, _ := protojson.Marshal(stmt)
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	assets := m["predicate"].(map[string]any)["assets"]
	if got, ok := assets.([]any); !ok || len(got) != 0 {
		t.Errorf("expected an empty assets array, got %v", assets)
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
