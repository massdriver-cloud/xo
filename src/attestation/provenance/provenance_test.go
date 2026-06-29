package provenance

import (
	"encoding/json"
	"testing"

	v1 "github.com/in-toto/attestation/go/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

func samplePredicate() Predicate {
	return Predicate{
		BuildType:            BuildType,
		ExternalParameters:   map[string]any{"instance": "inst-1"},
		ResolvedDependencies: []*v1.ResourceDescriptor{{Uri: "pkg:bundle/my-rds-bundle@v1.2.3"}},
		Builder:              BuilderID,
		InvocationID:         "deploy-123",
	}
}

func TestNewStatement_Success(t *testing.T) {
	stmt, err := NewStatement("massdriver://deploy-123", samplePredicate())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stmt.PredicateType != PredicateType {
		t.Errorf("expected PredicateType %s, got %s", PredicateType, stmt.PredicateType)
	}
	if stmt.PredicateType != "https://slsa.dev/provenance/v1" {
		t.Errorf("expected SLSA provenance predicate type, got %s", stmt.PredicateType)
	}
	if len(stmt.Subject) != 1 || stmt.Subject[0].Uri != "massdriver://deploy-123" {
		t.Fatalf("expected the deployment as subject, got %+v", stmt.Subject)
	}
	if len(stmt.Subject[0].Digest["sha256"]) != 64 {
		t.Errorf("expected a sha256 digest of the deployment uri, got %q", stmt.Subject[0].Digest["sha256"])
	}
	if err := stmt.Validate(); err != nil {
		t.Fatalf("statement failed in-toto validation: %v", err)
	}

	b, err := protojson.Marshal(stmt)
	if err != nil {
		t.Fatalf("failed to marshal statement: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	pred := m["predicate"].(map[string]any)
	if _, ok := pred["buildDefinition"]; !ok {
		t.Error("predicate should contain buildDefinition")
	}
	buildDef := pred["buildDefinition"].(map[string]any)
	if _, ok := buildDef["internalParameters"]; ok {
		t.Error("provenance should not carry internalParameters; it is provisioner-free")
	}
}

func TestNewStatement_ValidationErrors(t *testing.T) {
	tests := []struct {
		name       string
		subjectURI string
		pred       Predicate
		wantErr    string
	}{
		{"empty subject uri", "", samplePredicate(), "subject URI is required"},
		{"missing deployment id", "massdriver://x", Predicate{}, "deployment ID is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewStatement(tt.subjectURI, tt.pred)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("expected error %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}
